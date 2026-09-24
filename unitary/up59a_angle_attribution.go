package unitary

import "math"

const UP59AAngleAttributionSchema = "wingless.up59a-angle-component-attribution.v1"

type UP59AAngleArm struct {
	Name       string    `json:"name"`
	Parameters []float64 `json:"parameters"`
}
type UP59AAngleMetric struct {
	Arm              string  `json:"arm"`
	Length           int     `json:"length"`
	Accuracy         float64 `json:"accuracy"`
	MaxNormDrift     float64 `json:"max_norm_drift"`
	MaxRoundTripError float64 `json:"max_round_trip_error"`
	Gate             bool    `json:"gate"`
}
type UP59AAngleAttributionResult struct {
	Schema            string             `json:"schema"`
	Experiment        string             `json:"experiment"`
	SourceUP58ASeal   string             `json:"source_up58a_seal"`
	TrainingChanged   bool               `json:"training_changed"`
	LearnedParameters []float64          `json:"learned_parameters"`
	ExactParameters   []float64          `json:"exact_parameters"`
	Lengths           []int              `json:"lengths"`
	Arms              []UP59AAngleArm    `json:"arms"`
	Metrics           []UP59AAngleMetric `json:"metrics"`
}

func RunUP59A() (UP59AAngleAttributionResult, error) {
	lengths := []int{16, 24, 48, 128}
	train, err := up56aTrainSamples()
	if err != nil { return UP59AAngleAttributionResult{}, err }
	referenceHeld, err := up56aHeldSamples([]int{16})
	if err != nil { return UP59AAngleAttributionResult{}, err }
	learned, err := up56aTrain("unitary_angle_attribution_source", train, referenceHeld, up56aApplyUnitary, true)
	if err != nil { return UP59AAngleAttributionResult{}, err }

	exact := []float64{math.Pi/2, math.Pi/2}
	arms := []UP59AAngleArm{
		{Name:"learned_both", Parameters:append([]float64(nil), learned.Parameters...)},
		{Name:"exact_swap_learned_double", Parameters:[]float64{math.Pi/2, learned.Parameters[1]}},
		{Name:"learned_swap_exact_double", Parameters:[]float64{learned.Parameters[0], math.Pi/2}},
		{Name:"exact_both", Parameters:append([]float64(nil), exact...)},
	}
	result := UP59AAngleAttributionResult{
		Schema:UP59AAngleAttributionSchema,
		Experiment:"UP-59A-angle-component-attribution",
		SourceUP58ASeal:"0bf110e1938fdde619364e26fb73f98e63e4a6e0",
		TrainingChanged:false,
		LearnedParameters:append([]float64(nil), learned.Parameters...),
		ExactParameters:exact,
		Lengths:append([]int(nil), lengths...),
		Arms:arms,
	}
	for _, arm := range arms {
		for _, length := range lengths {
			held, err := up56aHeldSamples([]int{length})
			if err != nil { return UP59AAngleAttributionResult{}, err }
			long := up57aLongOnly(held)
			_, accuracy, norm, roundTrip, _, err := up56aEvaluate(arm.Parameters, long, up56aApplyUnitary, true)
			if err != nil { return UP59AAngleAttributionResult{}, err }
			result.Metrics = append(result.Metrics, UP59AAngleMetric{
				Arm:arm.Name, Length:length, Accuracy:accuracy,
				MaxNormDrift:norm, MaxRoundTripError:roundTrip,
				Gate:accuracy >= 0.98 && norm <= 1e-9 && roundTrip <= 1e-8,
			})
		}
	}
	return result, nil
}
