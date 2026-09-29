package unitary

const UPLM6CSchema = "wingless.up-lm6c-prepressure-deficit-law.v1"

type UPLM6CRow struct {
	Profile         string `json:"profile"`
	Rotation        int    `json:"rotation"`
	Permutation     string `json:"permutation"`
	Arm             int    `json:"arm"`
	Target          int    `json:"target"`
	PreCountdown    int    `json:"pre_countdown"`
	Writes          int    `json:"writes"`
	Deficit         int    `json:"deficit"`
	WriteError      int    `json:"write_error"`
	PostCountdown   int    `json:"post_countdown"`
	PostTargetError int    `json:"post_target_error"`
}

type UPLM6CResult struct {
	Schema                 string         `json:"schema"`
	Experiment             string         `json:"experiment"`
	SourceUPLM6BSeal       string         `json:"source_up_lm6b_seal"`
	ParentClassification   string         `json:"parent_classification"`
	Rows                   []UPLM6CRow    `json:"rows"`
	ExactLawRows           int            `json:"exact_law_rows"`
	MismatchRows           int            `json:"mismatch_rows"`
	MinWriteError          int            `json:"min_write_error"`
	MaxWriteError          int            `json:"max_write_error"`
	MinPostTargetError     int            `json:"min_post_target_error"`
	MaxPostTargetError     int            `json:"max_post_target_error"`
	MismatchByProfile      map[string]int `json:"mismatch_by_profile"`
	LiveActivation         bool           `json:"live_activation"`
	PolicyChanged          bool           `json:"policy_changed"`
	Classification         string         `json:"classification"`
}

func RunUPLM6C() (UPLM6CResult, error) {
	p, err := RunUPLM6B()
	if err != nil {
		return UPLM6CResult{}, err
	}
	r := UPLM6CResult{
		Schema:               UPLM6CSchema,
		Experiment:           "UP-LM6C-prepressure-deficit-law",
		SourceUPLM6BSeal:     "44c2c687b090b7ef3a4063227fdbd546f604d1ef",
		ParentClassification: p.Classification,
		MismatchByProfile:    map[string]int{},
	}
	first := true
	allWritesExact := true
	allPostExact := true
	for _, x := range p.Rows {
		deficit := x.PreCountdown - x.Target
		if deficit < 0 {
			deficit = 0
		}
		we := x.Writes - deficit
		pe := x.PostCountdown - x.Target
		row := UPLM6CRow{
			Profile: x.Profile, Rotation: x.Rotation, Permutation: x.Permutation,
			Arm: x.Arm, Target: x.Target, PreCountdown: x.PreCountdown,
			Writes: x.Writes, Deficit: deficit, WriteError: we,
			PostCountdown: x.PostCountdown, PostTargetError: pe,
		}
		r.Rows = append(r.Rows, row)
		if first {
			r.MinWriteError, r.MaxWriteError = we, we
			r.MinPostTargetError, r.MaxPostTargetError = pe, pe
			first = false
		} else {
			if we < r.MinWriteError { r.MinWriteError = we }
			if we > r.MaxWriteError { r.MaxWriteError = we }
			if pe < r.MinPostTargetError { r.MinPostTargetError = pe }
			if pe > r.MaxPostTargetError { r.MaxPostTargetError = pe }
		}
		if we == 0 && pe == 0 {
			r.ExactLawRows++
		} else {
			r.MismatchRows++
			r.MismatchByProfile[x.Profile]++
		}
		if we != 0 { allWritesExact = false }
		if pe != 0 { allPostExact = false }
	}
	if p.Classification != "TARGET_DEPENDENT_EXTENSION" {
		r.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if allWritesExact && allPostExact {
		r.Classification = "EXACT_DEFICIT_LAW"
	} else if allWritesExact {
		r.Classification = "WRITE_DEFICIT_EXACT_POST_RESIDUAL"
	} else if !allWritesExact {
		r.Classification = "TARGET_DEPENDENT_RESIDUAL"
	} else {
		r.Classification = "OTHER_VALID_PATTERN"
	}
	return r, nil
}
