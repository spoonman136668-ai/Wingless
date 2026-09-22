package unitary

import "math"

const ProbeSchema = "wingless.unitary-probe.v1"

type ProbeCase struct {
	ID            string `json:"id"`
	ExpectedRoute int    `json:"expected_route"`
	BaselineRoute int    `json:"baseline_route"`
	UnitaryRoute  int    `json:"unitary_route"`
}

type ProbeResult struct {
	Schema            string      `json:"schema"`
	Experiment        string      `json:"experiment"`
	Cases             []ProbeCase `json:"cases"`
	BaselineAccuracy  float64     `json:"baseline_accuracy"`
	UnitaryAccuracy   float64     `json:"unitary_accuracy"`
	MaxNormDrift      float64     `json:"max_norm_drift"`
	MaxRoundTripError float64     `json:"max_round_trip_error"`
}

// RunUP0 tests one narrow claim: a norm-preserving propagator can turn relative
// phase/sign information that is invisible to direct magnitude measurement into
// a deterministic routing decision through interference.
//
// This is a primitive proof only. It does not claim language-model parity,
// quantum advantage, or frontier reasoning.
func RunUP0() (ProbeResult, error) {
	r := 1 / math.Sqrt2
	cases := []struct {
		id       string
		state    State
		expected int
	}{
		{"in_phase", State{complex(r, 0), complex(r, 0)}, 1},
		{"opposed_phase", State{complex(r, 0), complex(-r, 0)}, 0},
		{"global_i_in_phase", State{complex(0, r), complex(0, r)}, 1},
		{"global_i_opposed_phase", State{complex(0, r), complex(0, -r)}, 0},
	}

	// exp(-i theta sigma_y), theta=pi/4. For the two input coordinates this
	// maps in-phase and opposed-phase states onto different basis routes.
	program := []Coupling{{A: 0, B: 1, Theta: math.Pi / 4}}

	result := ProbeResult{
		Schema:     ProbeSchema,
		Experiment: "UP-0-interference-routing",
		Cases:      make([]ProbeCase, 0, len(cases)),
	}

	var baselineHits, unitaryHits int
	for _, tc := range cases {
		baselineRoute, err := ArgMax(tc.state)
		if err != nil {
			return ProbeResult{}, err
		}

		beforeNorm, err := NormSquared(tc.state)
		if err != nil {
			return ProbeResult{}, err
		}
		evolved, err := Propagate(tc.state, program)
		if err != nil {
			return ProbeResult{}, err
		}
		afterNorm, err := NormSquared(evolved)
		if err != nil {
			return ProbeResult{}, err
		}
		unitaryRoute, err := ArgMax(evolved)
		if err != nil {
			return ProbeResult{}, err
		}
		recovered, err := Propagate(evolved, Invert(program))
		if err != nil {
			return ProbeResult{}, err
		}
		roundTrip, err := L2Distance(tc.state, recovered)
		if err != nil {
			return ProbeResult{}, err
		}

		if baselineRoute == tc.expected {
			baselineHits++
		}
		if unitaryRoute == tc.expected {
			unitaryHits++
		}
		drift := math.Abs(afterNorm - beforeNorm)
		if drift > result.MaxNormDrift {
			result.MaxNormDrift = drift
		}
		if roundTrip > result.MaxRoundTripError {
			result.MaxRoundTripError = roundTrip
		}
		result.Cases = append(result.Cases, ProbeCase{
			ID:            tc.id,
			ExpectedRoute: tc.expected,
			BaselineRoute: baselineRoute,
			UnitaryRoute:  unitaryRoute,
		})
	}

	result.BaselineAccuracy = float64(baselineHits) / float64(len(cases))
	result.UnitaryAccuracy = float64(unitaryHits) / float64(len(cases))
	return result, nil
}
