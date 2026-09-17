package nr1

import (
	"errors"
	"fmt"
)

const (
	QualifiedQwen3CoderLayers          = 48
	QualifiedQwen3CoderExpertsPerLayer = 128
	QualifiedQwen3CoderExpertsPerToken = 8
)

type traceTokenKey struct {
	session  string
	workload string
	task     string
	token    int64
}

type traceWorkloadKey struct {
	session  string
	workload string
	task     string
}

type traceTokenState struct {
	layerMask uint64
}

type traceWorkloadRange struct {
	minToken int64
	maxToken int64
	count    int64
	set      bool
}

// ValidateQualifiedQwen3CoderTrace rejects telemetry that cannot be a complete
// passive routing observation of the hash-pinned Qwen3-Coder-30B-A3B target.
// This is intentionally stricter than TraceEvent.Validate, which remains a
// reusable wire-format validator for future NR experiments.
func ValidateQualifiedQwen3CoderTrace(events []TraceEvent) error {
	if len(events) == 0 {
		return errors.New("NR1_QWEN_TRACE_EMPTY")
	}

	const expectedMask uint64 = (uint64(1) << QualifiedQwen3CoderLayers) - 1
	tokens := make(map[traceTokenKey]traceTokenState, len(events)/QualifiedQwen3CoderLayers+1)

	for i, event := range events {
		if err := event.Validate(); err != nil {
			return fmt.Errorf("NR1_QWEN_TRACE_EVENT_INVALID index=%d: %w", i, err)
		}
		if event.Layer >= QualifiedQwen3CoderLayers {
			return fmt.Errorf("NR1_QWEN_TRACE_LAYER_INVALID index=%d layer=%d", i, event.Layer)
		}
		if len(event.Experts) != QualifiedQwen3CoderExpertsPerToken {
			return fmt.Errorf("NR1_QWEN_TRACE_TOPK_INVALID index=%d got=%d want=%d",
				i, len(event.Experts), QualifiedQwen3CoderExpertsPerToken)
		}
		for _, expert := range event.Experts {
			if expert < 0 || expert >= QualifiedQwen3CoderExpertsPerLayer {
				return fmt.Errorf("NR1_QWEN_TRACE_EXPERT_ID_INVALID index=%d expert=%d", i, expert)
			}
		}

		key := traceTokenKey{
			session: event.SessionID,
			workload: event.WorkloadID,
			task: event.TaskFamily,
			token: event.TokenIndex,
		}
		state := tokens[key]
		bit := uint64(1) << uint(event.Layer)
		if state.layerMask&bit != 0 {
			return fmt.Errorf("NR1_QWEN_TRACE_DUPLICATE_LAYER token=%d layer=%d workload=%s",
				event.TokenIndex, event.Layer, event.WorkloadID)
		}
		state.layerMask |= bit
		tokens[key] = state
	}

	ranges := make(map[traceWorkloadKey]traceWorkloadRange)
	for key, state := range tokens {
		if state.layerMask != expectedMask {
			return fmt.Errorf("NR1_QWEN_TRACE_INCOMPLETE_LAYERS token=%d workload=%s mask=%012x want=%012x",
				key.token, key.workload, state.layerMask, expectedMask)
		}
		workloadKey := traceWorkloadKey{session: key.session, workload: key.workload, task: key.task}
		r := ranges[workloadKey]
		if !r.set {
			r.minToken = key.token
			r.maxToken = key.token
			r.set = true
		} else {
			if key.token < r.minToken {
				r.minToken = key.token
			}
			if key.token > r.maxToken {
				r.maxToken = key.token
			}
		}
		r.count++
		ranges[workloadKey] = r
	}

	for key, r := range ranges {
		if !r.set || r.minToken != 0 || r.maxToken+1 != r.count {
			return fmt.Errorf("NR1_QWEN_TRACE_TOKEN_SEQUENCE_INVALID workload=%s min=%d max=%d count=%d",
				key.workload, r.minToken, r.maxToken, r.count)
		}
	}
	return nil
}
