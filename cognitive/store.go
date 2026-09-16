package cognitive

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	maxMemoryRecords      = 4096
	maxSkillRecords       = 2048
	maxMemoryKeyBytes     = 512
	maxMemoryContentBytes = 64 << 10
	maxProvenanceBytes    = 1024
	maxSkillSourceBytes   = 1024
	maxTaskFamilyBytes    = 256
	maxSkillBodyBytes     = 64 << 10
	maxSkillListItems     = 32
	maxSkillListItemBytes = 2048
	maxValidatorIDBytes   = 256
)

type FileStore struct {
	root string
	now  func() time.Time
	mu   sync.Mutex
}

func NewFileStore(root string) (*FileStore, error) {
	if strings.TrimSpace(root) == "" {
		return nil, errors.New("cognitive store root is required")
	}
	return &FileStore{root: root, now: time.Now}, nil
}

func (s *FileStore) PutMemory(ctx context.Context, rec MemoryRecord) (MemoryRecord, error) {
	if err := ctx.Err(); err != nil {
		return MemoryRecord{}, err
	}
	if !rec.Class.Valid() || strings.TrimSpace(rec.Key) == "" || strings.TrimSpace(rec.Provenance) == "" || rec.Version < 1 || !utf8.ValidString(rec.Key) || !utf8.ValidString(rec.Content) || !utf8.ValidString(rec.Provenance) {
		return MemoryRecord{}, errors.New("invalid memory record")
	}
	if len(rec.Key) > maxMemoryKeyBytes || len(rec.Content) > maxMemoryContentBytes || len(rec.Provenance) > maxProvenanceBytes {
		return MemoryRecord{}, errors.New("memory record exceeds hard storage bound")
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = s.now().UTC()
	}
	want, err := memoryID(rec)
	if err != nil {
		return MemoryRecord{}, err
	}
	if rec.ID != "" && rec.ID != want {
		return MemoryRecord{}, errors.New("memory id does not match content")
	}
	rec.ID = want
	s.mu.Lock()
	defer s.mu.Unlock()
	dir := filepath.Join(s.root, "memory")
	target := filepath.Join(dir, rec.ID+".json")
	if existing, ok, err := existingMemory(target, rec.ID); err != nil {
		return MemoryRecord{}, err
	} else if ok {
		return existing, nil
	}
	if err := enforceRecordCapacity(dir, maxMemoryRecords); err != nil {
		return MemoryRecord{}, err
	}
	if err := s.writeJSON(dir, rec.ID+".json", rec, true); err != nil {
		return MemoryRecord{}, err
	}
	return rec, nil
}

func (s *FileStore) LookupMemory(ctx context.Context, class MemoryClass, key string) (MemoryRecord, bool, error) {
	if err := ctx.Err(); err != nil {
		return MemoryRecord{}, false, err
	}
	if !class.Valid() || strings.TrimSpace(key) == "" {
		return MemoryRecord{}, false, nil
	}
	files, err := filepath.Glob(filepath.Join(s.root, "memory", "*.json"))
	if err != nil {
		return MemoryRecord{}, false, err
	}
	var matches []MemoryRecord
	for _, path := range files {
		var rec MemoryRecord
		if err := readJSON(path, &rec); err != nil {
			return MemoryRecord{}, false, fmt.Errorf("read memory record: %w", err)
		}
		if rec.Class == class && rec.Key == key {
			matches = append(matches, rec)
		}
	}
	if len(matches) == 0 {
		return MemoryRecord{}, false, nil
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Version == matches[j].Version {
			return matches[i].ID < matches[j].ID
		}
		return matches[i].Version > matches[j].Version
	})
	return matches[0], true, nil
}

func (s *FileStore) DeleteMemory(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !validDigest(id) {
		return errors.New("invalid memory id")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(filepath.Join(s.root, "memory", id+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *FileStore) PutSkillCandidate(ctx context.Context, skill Skill) (Skill, error) {
	if err := ctx.Err(); err != nil {
		return Skill{}, err
	}
	if strings.TrimSpace(skill.Source) == "" || strings.TrimSpace(skill.TaskFamily) == "" || skill.Version < 1 || strings.TrimSpace(skill.Body) == "" || !utf8.ValidString(skill.Source) || !utf8.ValidString(skill.TaskFamily) || !utf8.ValidString(skill.Body) {
		return Skill{}, errors.New("invalid skill candidate")
	}
	if len(skill.Source) > maxSkillSourceBytes || len(skill.TaskFamily) > maxTaskFamilyBytes || len(skill.Body) > maxSkillBodyBytes || !boundedStrings(skill.Inputs) || !boundedStrings(skill.ExpectedOutputs) {
		return Skill{}, errors.New("skill candidate exceeds hard storage bound")
	}
	switch skill.Kind {
	case SkillStaticText, SkillProcedure, SkillTemplate:
	default:
		return Skill{}, errors.New("unsupported skill kind")
	}
	if skill.Validation == ValidationValidated || skill.ValidationEvidence != nil {
		return Skill{}, errors.New("candidate cannot self-declare validation")
	}
	skill.Validation = ValidationCandidate
	skill.UsageCount = 0
	skill.SuccessCount = 0
	skill.FailureCount = 0
	skill.LastUseAt = nil
	if skill.CreatedAt.IsZero() {
		skill.CreatedAt = s.now().UTC()
	}
	want, err := skillID(skill)
	if err != nil {
		return Skill{}, err
	}
	if skill.ID != "" && skill.ID != want {
		return Skill{}, errors.New("skill id does not match content")
	}
	skill.ID = want
	s.mu.Lock()
	defer s.mu.Unlock()
	dir := filepath.Join(s.root, "skills")
	target := filepath.Join(dir, skill.ID+".json")
	if existing, ok, err := existingSkill(target, skill.ID); err != nil {
		return Skill{}, err
	} else if ok {
		return existing, nil
	}
	if err := enforceRecordCapacity(dir, maxSkillRecords); err != nil {
		return Skill{}, err
	}
	if err := s.writeJSON(dir, skill.ID+".json", skill, true); err != nil {
		return Skill{}, err
	}
	return skill, nil
}

func (s *FileStore) PromoteSkill(ctx context.Context, id string, ev ValidationEvidence) (Skill, error) {
	if err := ctx.Err(); err != nil {
		return Skill{}, err
	}
	if !validDigest(id) || strings.TrimSpace(ev.ValidatorID) == "" || len(ev.ValidatorID) > maxValidatorIDBytes || !utf8.ValidString(ev.ValidatorID) || !validDigest(ev.EvidenceSHA256) {
		return Skill{}, errors.New("invalid deterministic validation evidence")
	}
	if ev.ValidatedAt.IsZero() {
		ev.ValidatedAt = s.now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	path := filepath.Join(s.root, "skills", id+".json")
	var skill Skill
	if err := readJSON(path, &skill); err != nil {
		return Skill{}, err
	}
	if skill.Validation != ValidationCandidate {
		return Skill{}, errors.New("only candidate skills may be promoted")
	}
	skill.Validation = ValidationValidated
	skill.ValidationEvidence = &ev
	if err := s.writeJSON(filepath.Join(s.root, "skills"), id+".json", skill, false); err != nil {
		return Skill{}, err
	}
	return skill, nil
}

func (s *FileStore) ListValidatedSkills(ctx context.Context, taskFamily string) ([]Skill, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(taskFamily) == "" {
		return nil, nil
	}
	files, err := filepath.Glob(filepath.Join(s.root, "skills", "*.json"))
	if err != nil {
		return nil, err
	}
	var out []Skill
	for _, path := range files {
		var skill Skill
		if err := readJSON(path, &skill); err != nil {
			return nil, fmt.Errorf("read skill record: %w", err)
		}
		if skill.TaskFamily == taskFamily && skill.Validation == ValidationValidated && skill.ValidationEvidence != nil {
			out = append(out, skill)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Version == out[j].Version {
			return out[i].ID < out[j].ID
		}
		return out[i].Version > out[j].Version
	})
	return out, nil
}

func (s *FileStore) RecordSkillUse(ctx context.Context, id string, success bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !validDigest(id) {
		return errors.New("invalid skill id")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	path := filepath.Join(s.root, "skills", id+".json")
	var skill Skill
	if err := readJSON(path, &skill); err != nil {
		return err
	}
	if skill.Validation != ValidationValidated || skill.ValidationEvidence == nil {
		return errors.New("unvalidated skill cannot be used")
	}
	skill.UsageCount++
	if success {
		skill.SuccessCount++
	} else {
		skill.FailureCount++
	}
	now := s.now().UTC()
	skill.LastUseAt = &now
	return s.writeJSON(filepath.Join(s.root, "skills"), id+".json", skill, false)
}

func boundedStrings(values []string) bool {
	if len(values) > maxSkillListItems {
		return false
	}
	for _, value := range values {
		if !utf8.ValidString(value) || len(value) > maxSkillListItemBytes {
			return false
		}
	}
	return true
}

func enforceRecordCapacity(dir string, max int) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}
	if len(files) >= max {
		return errors.New("cognitive store record capacity reached")
	}
	return nil
}

func existingMemory(path, id string) (MemoryRecord, bool, error) {
	var rec MemoryRecord
	if err := readJSON(path, &rec); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return MemoryRecord{}, false, nil
		}
		return MemoryRecord{}, false, err
	}
	want, err := memoryID(rec)
	if err != nil {
		return MemoryRecord{}, false, err
	}
	if rec.ID != id || want != id {
		return MemoryRecord{}, false, errors.New("content-addressed memory record tamper detected")
	}
	return rec, true, nil
}

func existingSkill(path, id string) (Skill, bool, error) {
	var skill Skill
	if err := readJSON(path, &skill); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Skill{}, false, nil
		}
		return Skill{}, false, err
	}
	want, err := skillID(skill)
	if err != nil {
		return Skill{}, false, err
	}
	if skill.ID != id || want != id {
		return Skill{}, false, errors.New("content-addressed skill record tamper detected")
	}
	return skill, true, nil
}

func memoryID(rec MemoryRecord) (string, error) {
	b, err := json.Marshal(struct {
		Version    int         `json:"version"`
		Class      MemoryClass `json:"class"`
		Key        string      `json:"key"`
		Content    string      `json:"content"`
		Provenance string      `json:"provenance"`
	}{rec.Version, rec.Class, rec.Key, rec.Content, rec.Provenance})
	if err != nil {
		return "", err
	}
	return digest(b), nil
}

func skillID(skill Skill) (string, error) {
	b, err := json.Marshal(struct {
		Version         int       `json:"version"`
		Source          string    `json:"source"`
		TaskFamily      string    `json:"task_family"`
		Inputs          []string  `json:"inputs"`
		ExpectedOutputs []string  `json:"expected_outputs"`
		Kind            SkillKind `json:"kind"`
		Body            string    `json:"body"`
	}{skill.Version, skill.Source, skill.TaskFamily, skill.Inputs, skill.ExpectedOutputs, skill.Kind, skill.Body})
	if err != nil {
		return "", err
	}
	return digest(b), nil
}

func digest(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func validDigest(s string) bool {
	if len(s) != sha256.Size*2 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil && s == strings.ToLower(s)
}

func readJSON(path string, dst any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		if err == nil {
			return errors.New("trailing JSON value")
		}
		return err
	}
	return nil
}

func (s *FileStore) writeJSON(dir, name string, value any, createOnly bool) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	target := filepath.Join(dir, name)
	if createOnly {
		if existing, err := os.ReadFile(target); err == nil {
			if bytes.Equal(bytes.TrimSpace(existing), b) {
				return nil
			}
			return errors.New("content-addressed record collision or tamper detected")
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	tmp, err := os.CreateTemp(dir, ".cognitive-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(append(b, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if !createOnly {
		_ = os.Remove(target)
	}
	return os.Rename(tmpName, target)
}
