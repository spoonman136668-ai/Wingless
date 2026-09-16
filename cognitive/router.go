package cognitive

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type Router struct {
	Policy Policy
	Memory MemoryStore
	Skills SkillStore
}

func (r Router) Decide(ctx context.Context, req RunRequest) (RoutingDecision, *MemoryRecord, *Skill, error) {
	if err := r.Policy.Validate(); err != nil {
		return RoutingDecision{}, nil, nil, err
	}
	if err := ctx.Err(); err != nil {
		return RoutingDecision{}, nil, nil, err
	}
	if req.RequestID == "" || req.RequestID != req.Base.ID || req.TaskFamily == "" {
		return RoutingDecision{}, nil, nil, errors.New("invalid cognitive run identity")
	}

	decision := RoutingDecision{
		Policy:     r.Policy.Version,
		RequestID:  req.RequestID,
		TaskFamily: req.TaskFamily,
	}

	if r.Policy.EnableMemory && r.Memory != nil && req.MemoryKey != "" && req.MemoryClass.Valid() {
		mem, ok, err := r.Memory.LookupMemory(ctx, req.MemoryClass, req.MemoryKey)
		if err != nil {
			return RoutingDecision{}, nil, nil, err
		}
		if ok {
			decision.Route = RouteMemory
			decision.Reason = "exact_memory_match"
			decision.MemoryIDs = []string{mem.ID}
			finalizeDecision(&decision, req)
			return decision, &mem, nil, nil
		}
	}

	if r.Policy.EnableSkills && r.Skills != nil {
		skills, err := r.Skills.ListValidatedSkills(ctx, req.TaskFamily)
		if err != nil {
			return RoutingDecision{}, nil, nil, err
		}
		for i := range skills {
			if skills[i].Kind != SkillStaticText {
				continue
			}
			decision.Route = RouteSkill
			decision.Reason = "validated_static_skill"
			decision.SkillIDs = []string{skills[i].ID}
			finalizeDecision(&decision, req)
			return decision, nil, &skills[i], nil
		}
	}

	if req.AllowMultiPass && r.Policy.MaxPasses > 1 {
		decision.Route = RouteInferenceMulti
		decision.Reason = "novel_or_unresolved_task"
	} else {
		decision.Route = RouteInferenceSingle
		decision.Reason = "single_pass_policy"
	}
	finalizeDecision(&decision, req)
	return decision, nil, nil, nil
}

func finalizeDecision(d *RoutingDecision, req RunRequest) {
	basis := struct {
		Policy         string      `json:"policy"`
		RequestID      string      `json:"request_id"`
		TaskFamily     string      `json:"task_family"`
		MemoryClass    MemoryClass `json:"memory_class"`
		MemoryKey      string      `json:"memory_key"`
		AllowMultiPass bool        `json:"allow_multi_pass"`
		BaseContext    string      `json:"base_context"`
		BaseMaxOutput  int         `json:"base_max_output"`
		Route          Route       `json:"route"`
		Reason         string      `json:"reason"`
		MemoryIDs      []string    `json:"memory_ids"`
		SkillIDs       []string    `json:"skill_ids"`
	}{
		d.Policy, req.RequestID, req.TaskFamily, req.MemoryClass, req.MemoryKey,
		req.AllowMultiPass, req.Base.Context, req.Base.MaxOutputTokens,
		d.Route, d.Reason, d.MemoryIDs, d.SkillIDs,
	}
	b, _ := json.Marshal(basis)
	sum := sha256.Sum256(b)
	d.InputDigest = hex.EncodeToString(sum[:])
	idSum := sha256.Sum256(append([]byte("decision:"), b...))
	d.ID = hex.EncodeToString(idSum[:])
}
