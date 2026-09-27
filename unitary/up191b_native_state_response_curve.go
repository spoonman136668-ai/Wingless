package unitary

import (
	"math"
	"sort"
)

const UP191BResponseCurveSchema = "wingless.up191b-native-state-response-curve.v1"

type UP191BProfile struct {
	Name string `json:"name"`
	SurfaceIndices []int `json:"surface_indices"`
}

type UP191BMetric struct {
	Profile string `json:"profile"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	PredictedActualRatio float64 `json:"predicted_actual_ratio"`
	AbsoluteMassRatioError float64 `json:"absolute_mass_ratio_error"`
	RequiredCalibrationFactor float64 `json:"required_calibration_factor"`
	TotalSlots int `json:"total_slots"`
}

type UP191BCountSummary struct {
	NativeCorrectCount int `json:"native_correct_count"`
	Observations int `json:"observations"`
	DistinctProfiles int `json:"distinct_profiles"`
	MinRequiredFactor float64 `json:"min_required_factor"`
	MaxRequiredFactor float64 `json:"max_required_factor"`
	RequiredFactorSpread float64 `json:"required_factor_spread"`
}

type UP191BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP190BSeal string `json:"source_up190b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Profiles []UP191BProfile `json:"profiles"`
	NewCorrectionFit bool `json:"new_correction_fit"`
	CorrectionApplied bool `json:"correction_applied"`
	AdaptiveProfileSearchUsed bool `json:"adaptive_profile_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	PearsonNativeCountRequiredFactor float64 `json:"pearson_native_count_required_factor"`
	Metrics []UP191BMetric `json:"metrics"`
	CountSummaries []UP191BCountSummary `json:"count_summaries"`
}

func up191bMetric(model up184bModel, profile string, phase int, indices []int) UP191BMetric {
	native := up190bNativeCorrectCount(phase, indices)
	rows := up190bRows(phase, indices)
	predicted := 0.0
	actual := 0
	for _, row := range rows {
		predicted += up188bP(model, row)
		if row.positive {
			actual++
		}
	}
	ratio := 0.0
	required := 0.0
	if actual > 0 {
		ratio = predicted / float64(actual)
	}
	if predicted > 0 {
		required = float64(actual) / predicted
	}
	return UP191BMetric{
		Profile: profile,
		Phase: phase,
		NativeCorrectCount: native,
		ActualCrossings: actual,
		PredictedCrossingsSum: predicted,
		PredictedActualRatio: ratio,
		AbsoluteMassRatioError: math.Abs(ratio - 1),
		RequiredCalibrationFactor: required,
		TotalSlots: len(rows),
	}
}

func up191bPearson(metrics []UP191BMetric) float64 {
	if len(metrics) < 2 {
		return 0
	}
	mx, my := 0.0, 0.0
	for _, m := range metrics {
		mx += float64(m.NativeCorrectCount)
		my += m.RequiredCalibrationFactor
	}
	n := float64(len(metrics))
	mx /= n
	my /= n
	num, dx, dy := 0.0, 0.0, 0.0
	for _, m := range metrics {
		x := float64(m.NativeCorrectCount) - mx
		y := m.RequiredCalibrationFactor - my
		num += x * y
		dx += x * x
		dy += y * y
	}
	if dx == 0 || dy == 0 {
		return 0
	}
	return num / math.Sqrt(dx*dy)
}

func RunUP191B() (UP191BResult, error) {
	cal := []int{26, 27, 28, 29, 30}
	eval := []int{43, 44, 45}
	profiles := []UP191BProfile{
		{Name: "store_forward", SurfaceIndices: []int{0, 1, 2, 3}},
		{Name: "store_reverse", SurfaceIndices: []int{3, 2, 1, 0}},
		{Name: "observe_forward", SurfaceIndices: []int{5, 6, 7, 8}},
		{Name: "observe_reverse", SurfaceIndices: []int{8, 7, 6, 5}},
		{Name: "mixed_forward", SurfaceIndices: []int{0, 5, 1, 6}},
		{Name: "mixed_reverse", SurfaceIndices: []int{6, 1, 5, 0}},
	}
	model := up184bBuild("pooled_26_30", cal)
	res := UP191BResult{
		Schema: UP191BResponseCurveSchema,
		Experiment: "UP-191B-native-state-response-curve",
		SourceUP190BSeal: "bae5008f361cca459694cffc82fd2302162920fb",
		CalibrationPhases: cal,
		EvaluationPhases: eval,
		Profiles: profiles,
		NewCorrectionFit: false,
		CorrectionApplied: false,
		AdaptiveProfileSearchUsed: false,
		PhaseInputUsed: false,
		ParityInputUsed: false,
		MaintenanceTriggered: false,
	}
	type acc struct {
		n int
		min float64
		max float64
		profiles map[string]bool
	}
	groups := map[int]*acc{}
	for _, p := range profiles {
		for _, phase := range eval {
			m := up191bMetric(model, p.Name, phase, p.SurfaceIndices)
			res.Metrics = append(res.Metrics, m)
			a := groups[m.NativeCorrectCount]
			if a == nil {
				a = &acc{min: m.RequiredCalibrationFactor, max: m.RequiredCalibrationFactor, profiles: map[string]bool{}}
				groups[m.NativeCorrectCount] = a
			}
			a.n++
			a.profiles[m.Profile] = true
			if m.RequiredCalibrationFactor < a.min {
				a.min = m.RequiredCalibrationFactor
			}
			if m.RequiredCalibrationFactor > a.max {
				a.max = m.RequiredCalibrationFactor
			}
		}
	}
	counts := make([]int, 0, len(groups))
	for c := range groups {
		counts = append(counts, c)
	}
	sort.Ints(counts)
	for _, c := range counts {
		a := groups[c]
		res.CountSummaries = append(res.CountSummaries, UP191BCountSummary{
			NativeCorrectCount: c,
			Observations: a.n,
			DistinctProfiles: len(a.profiles),
			MinRequiredFactor: a.min,
			MaxRequiredFactor: a.max,
			RequiredFactorSpread: a.max - a.min,
		})
	}
	res.PearsonNativeCountRequiredFactor = up191bPearson(res.Metrics)
	return res, nil
}
