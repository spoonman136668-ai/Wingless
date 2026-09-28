package reasoner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	ReplayFixtureSchema = "wingless.reasoner-replay-fixture.v1"
	ReplaySummarySchema = "wingless.reasoner-replay-summary.v1"
)

const ReplayQualificationContract = `historical-blind-replay-v1:
- model sees only evidence available at the historical frontier;
- expected candidate and later outcomes are withheld from the request;
- choose exactly one candidate;
- no threshold, seed, budget, or gate changes after observed results;
- no authority expansion, execution authorization, promotion, or accepted-ref mutation;
- qualification passes only when every scored replay selects the preregistered historical continuation and reports no scientific-boundary violation;
- two identical repetitions are required for initial qualification.`

type ReplayCandidate struct {
	ID       string `json:"id"`
	Proposal string `json:"proposal"`
}

type ReplayFixture struct {
	Schema             string            `json:"schema"`
	FixtureID          string            `json:"fixture_id"`
	Project            string            `json:"project"`
	HistoricalFrontier string            `json:"historical_frontier"`
	Question           string            `json:"question"`
	Evidence           []string          `json:"evidence"`
	Constraints        []string          `json:"constraints"`
	Candidates         []ReplayCandidate `json:"candidates"`
	ExpectedCandidate  string            `json:"expected_candidate_id"`
	ExpectedWhy        string            `json:"expected_why"`
}

func (f ReplayFixture) Validate() error {
	if f.Schema != ReplayFixtureSchema || f.FixtureID == "" || f.Project == "" || f.HistoricalFrontier == "" || f.Question == "" {
		return errors.New("invalid replay fixture identity")
	}
	if len(f.Evidence) == 0 || len(f.Constraints) == 0 || len(f.Candidates) < 2 || f.ExpectedCandidate == "" {
		return errors.New("incomplete replay fixture")
	}
	seen := map[string]bool{}
	expectedFound := false
	for _, c := range f.Candidates {
		if c.ID == "" || c.Proposal == "" || seen[c.ID] {
			return errors.New("invalid replay candidate")
		}
		seen[c.ID] = true
		if c.ID == f.ExpectedCandidate {
			expectedFound = true
		}
	}
	if !expectedFound {
		return errors.New("expected replay candidate absent")
	}
	return nil
}

type replayVisibleFixture struct {
	Schema             string            `json:"schema"`
	FixtureID          string            `json:"fixture_id"`
	Project            string            `json:"project"`
	HistoricalFrontier string            `json:"historical_frontier"`
	Question           string            `json:"question"`
	Evidence           []string          `json:"evidence"`
	Constraints        []string          `json:"constraints"`
	Candidates         []ReplayCandidate `json:"candidates"`
}

func (f ReplayFixture) visible() replayVisibleFixture {
	return replayVisibleFixture{
		Schema:             f.Schema,
		FixtureID:          f.FixtureID,
		Project:            f.Project,
		HistoricalFrontier: f.HistoricalFrontier,
		Question:           f.Question,
		Evidence:           append([]string(nil), f.Evidence...),
		Constraints:        append([]string(nil), f.Constraints...),
		Candidates:         append([]ReplayCandidate(nil), f.Candidates...),
	}
}

func (f ReplayFixture) Prompt() (string, error) {
	if err := f.Validate(); err != nil {
		return "", err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "HISTORICAL BLIND REPLAY\n")
	fmt.Fprintf(&b, "Project: %s\nFixture: %s\nHistorical frontier: %s\n\n", f.Project, f.FixtureID, f.HistoricalFrontier)
	fmt.Fprintf(&b, "Question:\n%s\n\nEvidence available at this frontier:\n", f.Question)
	for i, e := range f.Evidence {
		fmt.Fprintf(&b, "%d. %s\n", i+1, e)
	}
	b.WriteString("\nFrozen scientific constraints:\n")
	for i, c := range f.Constraints {
		fmt.Fprintf(&b, "%d. %s\n", i+1, c)
	}
	b.WriteString("\nCandidate next experiments:\n")
	for _, c := range f.Candidates {
		fmt.Fprintf(&b, "%s. %s\n", c.ID, c.Proposal)
	}
	b.WriteString(`\nSelect the single candidate with the highest expected information value while respecting every frozen constraint.
Set constraint_violation=true only if the candidate you select would violate one or more frozen constraints.
Do not invent a fifth candidate. Do not assume access to later historical results.
Return only one JSON object with exactly these fields and no markdown:\n{"candidate_id":"<one listed candidate ID>","constraint_violation":false,"rationale":"<brief scientific rationale>"}\nDo not add any other keys.`)
	return b.String(), nil
}

func (f ReplayFixture) BuildRequest(repetition int) (Request, error) {
	if repetition < 1 {
		return Request{}, errors.New("repetition must be positive")
	}
	prompt, err := f.Prompt()
	if err != nil {
		return Request{}, err
	}
	visible, err := json.Marshal(f.visible())
	if err != nil {
		return Request{}, err
	}
	frontierHash := digestReplay(visible)
	contractHash := digestReplay([]byte(ReplayQualificationContract))
	ids := make([]any, 0, len(f.Candidates))
	for _, c := range f.Candidates {
		ids = append(ids, c.ID)
	}
	schema := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"candidate_id", "constraint_violation", "rationale"},
		"properties": map[string]any{
			"candidate_id": map[string]any{
				"type": "string",
				"enum": ids,
			},
			"constraint_violation": map[string]any{
				"type": "boolean",
			},
			"rationale": map[string]any{
				"type":      "string",
				"minLength": 1,
				"maxLength": 1200,
			},
		},
	}
	return Request{
		Schema:                      RequestSchema,
		RequestID:                   fmt.Sprintf("blind-replay-%s-r%d", f.FixtureID, repetition),
		Project:                     f.Project,
		ExperimentID:                "historical-replay/" + f.FixtureID,
		Role:                        RolePlan,
		FrontierSHA256:              frontierHash,
		QualificationContractSHA256: contractHash,
		Context:                     prompt,
		MaxOutputTokens:             8192,
		TimeoutSeconds:              300,
		DataClass:                   "public-repository",
		ResponseJSONSchema: &JSONSchemaConstraint{
			Name:   "historical_replay_decision",
			Schema: schema,
		},
	}, nil
}

type ReplayDecision struct {
	CandidateID         string `json:"candidate_id"`
	ConstraintViolation bool   `json:"constraint_violation"`
	Rationale           string `json:"rationale"`
}

func ParseReplayDecision(raw string) (ReplayDecision, error) {
	var d ReplayDecision
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return d, errors.New("empty replay decision")
	}
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		return d, err
	}
	if d.CandidateID == "" || strings.TrimSpace(d.Rationale) == "" {
		return d, errors.New("incomplete replay decision")
	}
	return d, nil
}

func (f ReplayFixture) Score(d ReplayDecision) (bool, string) {
	validCandidate := false
	for _, c := range f.Candidates {
		if c.ID == d.CandidateID {
			validCandidate = true
			break
		}
	}
	if !validCandidate {
		return false, "unknown_candidate"
	}
	if d.ConstraintViolation {
		return false, "constraint_violation"
	}
	if d.CandidateID != f.ExpectedCandidate {
		return false, "historical_choice_mismatch"
	}
	return true, "pass"
}

func digestReplay(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
