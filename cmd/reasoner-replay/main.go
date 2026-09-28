package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spoonman136668-ai/Wingless/reasoner"
)

type observation struct {
	FixtureID                   string `json:"fixture_id"`
	Repetition                  int    `json:"repetition"`
	ExpectedCandidate           string `json:"expected_candidate"`
	SelectedCandidate           string `json:"selected_candidate"`
	ConstraintViolation         bool   `json:"constraint_violation"`
	Score                       string `json:"score"`
	Passed                      bool   `json:"passed"`
	Provider                    string `json:"provider,omitempty"`
	ReturnedModel               string `json:"returned_model"`
	QualificationIdentitySHA256 string `json:"qualification_identity_sha256"`
	ResponseSHA256              string `json:"response_sha256"`
	LatencyMS                   int64  `json:"latency_ms"`
	Rationale                   string `json:"rationale"`
}

type summary struct {
	Schema         string        `json:"schema"`
	Model          string        `json:"model"`
	FixtureCount   int           `json:"fixture_count"`
	Repetitions    int           `json:"repetitions"`
	TotalDecisions int           `json:"total_decisions"`
	Correct        int           `json:"correct"`
	Violations     int           `json:"constraint_violations"`
	Passed         bool          `json:"passed"`
	StartedAtUTC   string        `json:"started_at_utc"`
	CompletedAtUTC string        `json:"completed_at_utc"`
	Observations   []observation `json:"observations"`
}

func main() {
	fixtureDir := flag.String("fixtures", "reasoner/fixtures/history", "historical replay fixture directory")
	outPath := flag.String("out", "evidence/nemotron-reasoner-replay.json", "qualification summary path")
	repeats := flag.Int("repeats", 2, "identical blind replay repetitions")
	flag.Parse()

	if os.Getenv("WINGLESS_REMOTE_REASONER_ENABLE") != "1" {
		die(errors.New("remote reasoner disabled"))
	}
	if *repeats != 2 {
		die(errors.New("initial qualification requires exactly two repetitions"))
	}
	key := os.Getenv("WINGLESS_REASONER_OPENROUTER_API_KEY")
	if strings.TrimSpace(key) == "" {
		die(errors.New("dedicated research OpenRouter key missing"))
	}

	fixtures, err := loadFixtures(*fixtureDir)
	if err != nil {
		die(err)
	}
	if len(fixtures) < 5 {
		die(fmt.Errorf("at least five blind replay fixtures required, got %d", len(fixtures)))
	}
	projects := map[string]bool{"Wingless": false, "Yggdrasil": false}
	for _, fixture := range fixtures {
		if _, required := projects[fixture.Project]; required {
			projects[fixture.Project] = true
		}
	}
	for project, present := range projects {
		if !present {
			die(fmt.Errorf("required replay project missing: %s", project))
		}
	}

	cfg := reasoner.DefaultConfig()
	client, err := reasoner.NewClient(key, cfg)
	if err != nil {
		die(err)
	}

	started := time.Now().UTC()
	s := summary{
		Schema:       reasoner.ReplaySummarySchema,
		Model:        cfg.Model,
		FixtureCount: len(fixtures),
		Repetitions:  *repeats,
		StartedAtUTC: started.Format(time.RFC3339Nano),
	}

	for rep := 1; rep <= *repeats; rep++ {
		for _, fixture := range fixtures {
			req, err := fixture.BuildRequest(rep)
			if err != nil {
				die(fmt.Errorf("%s request: %w", fixture.FixtureID, err))
			}
			result, err := client.Invoke(context.Background(), req)
			if err != nil {
				die(fmt.Errorf("%s repetition %d inference: %w", fixture.FixtureID, rep, err))
			}
			decision, err := reasoner.ParseReplayDecision(result.Text)
			if err != nil {
				die(fmt.Errorf("%s repetition %d parse: %w", fixture.FixtureID, rep, err))
			}
			passed, score := fixture.Score(decision)
			if passed {
				s.Correct++
			}
			if decision.ConstraintViolation {
				s.Violations++
			}
			s.TotalDecisions++
			s.Observations = append(s.Observations, observation{
				FixtureID:                   fixture.FixtureID,
				Repetition:                  rep,
				ExpectedCandidate:           fixture.ExpectedCandidate,
				SelectedCandidate:           decision.CandidateID,
				ConstraintViolation:         decision.ConstraintViolation,
				Score:                       score,
				Passed:                      passed,
				Provider:                    result.Provider,
				ReturnedModel:               result.ReturnedModel,
				QualificationIdentitySHA256: result.QualificationIdentitySHA256,
				ResponseSHA256:              result.ResponseSHA256,
				LatencyMS:                   result.LatencyMS,
				Rationale:                   decision.Rationale,
			})
			fmt.Printf("REPLAY fixture=%s repetition=%d selected=%s expected=%s score=%s provider=%s\n",
				fixture.FixtureID, rep, decision.CandidateID, fixture.ExpectedCandidate, score, result.Provider)
		}
	}

	s.CompletedAtUTC = time.Now().UTC().Format(time.RFC3339Nano)
	s.Passed = s.TotalDecisions == len(fixtures)*(*repeats) &&
		s.Correct == s.TotalDecisions &&
		s.Violations == 0

	if err := writeSummary(*outPath, s); err != nil {
		die(err)
	}
	if !s.Passed {
		fmt.Printf("NEMOTRON_REASONER_REPLAY_FAIL correct=%d total=%d violations=%d\n", s.Correct, s.TotalDecisions, s.Violations)
		os.Exit(1)
	}
	fmt.Printf("NEMOTRON_REASONER_REPLAY_PASS correct=%d total=%d violations=%d\n", s.Correct, s.TotalDecisions, s.Violations)
}

func loadFixtures(dir string) ([]reasoner.ReplayFixture, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(paths)
	fixtures := make([]reasoner.ReplayFixture, 0, len(paths))
	ids := map[string]bool{}
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var fixture reasoner.ReplayFixture
		if err := json.Unmarshal(raw, &fixture); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if err := fixture.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if ids[fixture.FixtureID] {
			return nil, fmt.Errorf("duplicate fixture id %s", fixture.FixtureID)
		}
		ids[fixture.FixtureID] = true
		fixtures = append(fixtures, fixture)
	}
	return fixtures, nil
}

func writeSummary(path string, value summary) error {
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".nemotron-replay-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(raw); err != nil {
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
	return os.Rename(name, path)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "reasoner-replay:", err)
	os.Exit(1)
}
