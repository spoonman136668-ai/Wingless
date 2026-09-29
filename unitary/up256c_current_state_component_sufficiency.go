package unitary

import "fmt"

const UP256CCurrentComponentSchema = "wingless.up256c-current-state-component-sufficiency.v1"

type UP256CResult struct {
	Schema                       string            `json:"schema"`
	Experiment                   string            `json:"experiment"`
	SourceUP255CSeal             string            `json:"source_up255c_seal"`
	PairedInitialStates          int               `json:"paired_initial_states"`
	ParentClassification         string            `json:"parent_classification"`
	ParentCurrentMixedGroups     int               `json:"parent_current_mixed_groups"`
	DiagnosticOnly               bool              `json:"diagnostic_only"`
	InterventionChanged          bool              `json:"intervention_changed"`
	ClassifierTrainingUsed       bool              `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool              `json:"adaptive_feature_selection_used"`
	LiveActivation               bool              `json:"live_activation"`
	Partitions                   []UP255CPartition `json:"partitions"`
	Classification               string            `json:"classification"`
}

func RunUP256C() (UP256CResult, error) {
	p, err := RunUP255C()
	if err != nil {
		return UP256CResult{}, err
	}
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	static, events, signatures, joint, classes := []string{}, []string{}, []string{}, []string{}, []string{}
	r := UP256CResult{
		Schema:               UP256CCurrentComponentSchema,
		Experiment:           "UP-256C-current-state-component-sufficiency",
		SourceUP255CSeal:     "789b8c380043343b7e751b4ddd5f9ca21c04c223",
		ParentClassification: p.Classification,
		DiagnosticOnly:       true,
	}
	if len(p.Partitions) > 1 {
		r.ParentCurrentMixedGroups = p.Partitions[1].MixedGroups
	}
	for _, co := range cohorts {
		for _, hand := range hands {
			x, e, ok := up161cTriggerState(hand, co)
			if !ok {
				continue
			}
			s := up223cSnapshot(x, e, 0)
			a, er := up237cRun(x, e, 22, "alternating_shield", "fixed_offset_refresh")
			if er != nil {
				return UP256CResult{}, er
			}
			b, er := up237cRun(x, e, 0, "fixed_offset_refresh", "alternating_shield")
			if er != nil {
				return UP256CResult{}, er
			}
			cl := up246cClass(a, b)
			sk := up255cStatic(s)
			ek := fmt.Sprintf("%s|ae=%t|be=%t", sk, a.EventPresent, b.EventPresent)
			ck := fmt.Sprintf("%s|ac=%s|bc=%s", sk, a.CurrentSignature, b.CurrentSignature)
			jk := fmt.Sprintf("%s|ae=%t|ac=%s|be=%t|bc=%s", sk, a.EventPresent, a.CurrentSignature, b.EventPresent, b.CurrentSignature)
			static = append(static, sk)
			events = append(events, ek)
			signatures = append(signatures, ck)
			joint = append(joint, jk)
			classes = append(classes, cl)
			r.PairedInitialStates++
		}
	}
	ps := up255cPart("STATIC9", static, classes)
	pe := up255cPart("STATIC9_EVENTS", events, classes)
	pc := up255cPart("STATIC9_SIGNATURES", signatures, classes)
	pj := up255cPart("STATIC9_EVENTS_SIGNATURES", joint, classes)
	r.Partitions = []UP255CPartition{ps, pe, pc, pj}

	anchor := p.Classification == "CURRENT_STATE_RESOLVES_ALIASING" &&
		len(p.Partitions) > 1 &&
		ps.MixedGroups == p.Partitions[0].MixedGroups &&
		pj.MixedGroups == p.Partitions[1].MixedGroups &&
		p.Partitions[1].MixedGroups == 0
	if !anchor {
		r.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if pe.MixedGroups == 0 && pc.MixedGroups == 0 {
		r.Classification = "BOTH_INDIVIDUALLY_SUFFICIENT"
	} else if pe.MixedGroups == 0 {
		r.Classification = "EVENTS_SUFFICIENT"
	} else if pc.MixedGroups == 0 {
		r.Classification = "CURRENT_SIGNATURE_SUFFICIENT"
	} else if pj.MixedGroups == 0 {
		r.Classification = "JOINT_ONLY"
	} else if pj.MixedGroups < ps.MixedGroups {
		r.Classification = "CURRENT_COMPONENTS_REDUCE_NOT_RESOLVE"
	} else {
		r.Classification = "CURRENT_COMPONENTS_DO_NOT_REDUCE"
	}
	return r, nil
}
