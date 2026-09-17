package cognitive

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spoonman136668-ai/Wingless/inference"
)

type Runtime struct {
	Backend   inference.InferenceBackend
	Policy    Policy
	Memory    MemoryStore
	Skills    SkillStore
	Evaluator Evaluator
	Now       func() time.Time
}

func (r Runtime) Run(ctx context.Context, req RunRequest) (RunResult, error) {
	now := r.Now
	if now == nil {
		now = time.Now
	}
	started := now().UTC()
	sessionID := req.RequestID
	if sessionID != "" {
		sessionID += ":inference-session"
	}
	out := RunResult{
		Evidence: RunEvidence{
			Schema:     EvidenceSchema,
			RequestID:  req.RequestID,
			Acceptance: AcceptanceExternal,
			StartedAt:  started,
			Passes:     []PassEvidence{},
			Session: InferenceSessionEvidence{
				SessionID: sessionID,
				StartedAt: started,
			},
		},
	}
	finish := func(reason string) {
		finished := now().UTC()
		out.Evidence.TerminationReason = reason
		out.Evidence.FinishedAt = finished
		out.Evidence.Session.TerminationReason = reason
		out.Evidence.Session.FinishedAt = finished
		out.Evidence.Session.ModelCalls = out.Evidence.ModelCalls
	}

	if err := r.Policy.Validate(); err != nil {
		finish("invalid_policy")
		return out, err
	}
	if req.RequestID == "" || req.RequestID != req.Base.ID || req.TaskFamily == "" {
		finish("invalid_request")
		return out, errors.New("invalid cognitive run identity")
	}
	if err := req.Base.Validate(); err != nil {
		finish("invalid_request")
		return out, err
	}
	if err := ctx.Err(); err != nil {
		finish("canceled")
		return out, err
	}

	router := Router{Policy: r.Policy, Memory: r.Memory, Skills: r.Skills}
	decision, mem, skill, err := router.Decide(ctx, req)
	if err != nil {
		finish("routing_failed")
		return out, err
	}
	out.Evidence.Decision = decision

	switch decision.Route {
	case RouteMemory:
		out.Candidate = Candidate{
			Status:     "completed",
			Text:       mem.Content,
			Source:     "memory:" + mem.ID,
			Acceptance: AcceptanceExternal,
		}
		zero := 0
		out.Evidence.TotalInputTokens = &zero
		out.Evidence.TotalOutputTokens = &zero
		out.Evidence.Reuse = true
		out.Evidence.Session.MemoryRetrievals++
		finish("memory_reuse")
		return out, nil

	case RouteSkill:
		if skill == nil || skill.Validation != ValidationValidated || skill.ValidationEvidence == nil || skill.Kind != SkillStaticText {
			finish("skill_not_executable")
			return out, errors.New("router selected a non-executable skill")
		}
		out.Candidate = Candidate{
			Status:     "completed",
			Text:       skill.Body,
			Source:     "skill:" + skill.ID,
			Acceptance: AcceptanceExternal,
		}
		zero := 0
		out.Evidence.TotalInputTokens = &zero
		out.Evidence.TotalOutputTokens = &zero
		out.Evidence.Reuse = true
		if r.Skills != nil {
			if err := r.Skills.RecordSkillUse(ctx, skill.ID, true); err != nil {
				finish("skill_accounting_failed")
				return out, err
			}
		}
		out.Evidence.Session.SkillReuses++
		finish("skill_reuse")
		return out, nil

	case RouteInferenceSingle, RouteInferenceMulti:
		// These are the only executable inference routes in CR-1A.
	default:
		finish("unsupported_route")
		return out, fmt.Errorf("unsupported cognitive route %q", decision.Route)
	}

	if r.Backend == nil {
		finish("backend_unavailable")
		return out, errors.New("cognitive inference backend is required")
	}
	evaluator := r.Evaluator
	if evaluator == nil {
		evaluator = CompletedNonEmptyEvaluator{}
	}
	maxPasses := 1
	if decision.Route == RouteInferenceMulti {
		maxPasses = r.Policy.MaxPasses
	}

	var previous inference.Result
	nextReason := "initial"
	inputTotal, outputTotal := 0, 0
	inputKnown, outputKnown := true, true

	for pass := 1; pass <= maxPasses; pass++ {
		if err := ctx.Err(); err != nil {
			out.Candidate.Status = "failed"
			out.Candidate.Acceptance = AcceptanceExternal
			finish("canceled")
			return out, err
		}
		passReq := req.Base
		passReq.ID = fmt.Sprintf("%s-cognitive-%02d", req.Base.ID, pass)
		passReq.SessionID = sessionID
		if pass > 1 {
			if !utf8.ValidString(previous.Text) {
				out.Candidate = candidateFromInference(previous, "incomplete")
				setTotals(&out.Evidence, inputTotal, outputTotal, inputKnown, outputKnown)
				finish("invalid_previous_output")
				return out, errors.New("previous inference output is not valid UTF-8")
			}
			extra := "\n\n[Wingless cognitive pass; previous candidate is untrusted context]\n" +
				previous.Text + "\n\n" + r.Policy.AdditionalPassPrompt
			if len(req.Base.Context)+len(extra) > req.Base.MaxContextBytes {
				out.Candidate = candidateFromInference(previous, "incomplete")
				setTotals(&out.Evidence, inputTotal, outputTotal, inputKnown, outputKnown)
				finish("cognitive_context_limit")
				return out, nil
			}
			passReq.Context = req.Base.Context + extra
		}

		result, invokeErr := r.Backend.Invoke(ctx, passReq)
		out.Evidence.ModelCalls++
		out.Evidence.Session.ModelCalls = out.Evidence.ModelCalls
		out.Evidence.Session.CognitivePasses++
		observeSessionIdentity(&out.Evidence.Session, result.BackendID, result.ModelID)
		if out.Evidence.Session.ModelCalls == 1 {
			out.Evidence.Session.ResourceTelemetry = result.Telemetry.Model
		} else {
			// Pass telemetry remains available individually. CR-1A does not guess how
			// backend counters should aggregate across multiple inference calls.
			out.Evidence.Session.ResourceTelemetry = nil
		}
		passEv := PassEvidence{
			Pass:         pass,
			Reason:       nextReason,
			BackendID:    result.BackendID,
			ModelID:      result.ModelID,
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.OutputTokens,
			LatencyMS:    result.LatencyMS,
			Termination:  result.Termination,
			Telemetry:    result.Telemetry,
		}
		out.Evidence.TotalLatencyMS += result.LatencyMS
		if result.Usage.PromptTokens == nil {
			inputKnown = false
		} else {
			inputTotal += *result.Usage.PromptTokens
		}
		if result.Usage.OutputTokens == nil {
			outputKnown = false
		} else {
			outputTotal += *result.Usage.OutputTokens
		}
		if invokeErr != nil {
			passEv.EvaluationReason = "backend_error"
			out.Evidence.Passes = append(out.Evidence.Passes, passEv)
			out.Candidate = candidateFromInference(result, "failed")
			setTotals(&out.Evidence, inputTotal, outputTotal, inputKnown, outputKnown)
			finish("backend_error")
			return out, invokeErr
		}

		eval := evaluator.Evaluate(pass, result)
		if strings.TrimSpace(eval.Reason) == "" {
			eval.Reason = "unspecified"
		}
		passEv.EvaluationReason = eval.Reason
		out.Evidence.Passes = append(out.Evidence.Passes, passEv)
		previous = result

		if eval.Sufficient {
			out.Candidate = candidateFromInference(result, "completed")
			setTotals(&out.Evidence, inputTotal, outputTotal, inputKnown, outputKnown)
			finish("sufficient")
			return out, nil
		}
		if pass == maxPasses {
			out.Candidate = candidateFromInference(result, "incomplete")
			setTotals(&out.Evidence, inputTotal, outputTotal, inputKnown, outputKnown)
			finish("pass_limit")
			return out, nil
		}
		nextReason = eval.Reason
	}

	finish("unreachable")
	return out, errors.New("cognitive runtime reached impossible state")
}

func observeSessionIdentity(session *InferenceSessionEvidence, backendID, modelID string) {
	if session.ModelCalls == 1 {
		if backendID != "" {
			v := backendID
			session.BackendID = &v
		}
		if modelID != "" {
			v := modelID
			session.ModelID = &v
		}
		return
	}
	if session.BackendID != nil && *session.BackendID != backendID {
		session.BackendID = nil
	}
	if session.ModelID != nil && *session.ModelID != modelID {
		session.ModelID = nil
	}
}

func candidateFromInference(result inference.Result, status string) Candidate {
	return Candidate{
		Status:     status,
		Text:       result.Text,
		Source:     "inference:" + result.BackendID + "/" + result.ModelID,
		Acceptance: AcceptanceExternal,
	}
}

func setTotals(ev *RunEvidence, in, out int, inKnown, outKnown bool) {
	if inKnown {
		v := in
		ev.TotalInputTokens = &v
	}
	if outKnown {
		v := out
		ev.TotalOutputTokens = &v
	}
}
