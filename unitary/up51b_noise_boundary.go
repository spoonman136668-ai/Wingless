package unitary

const UP51BNoiseBoundarySchema = "wingless.up51b-noise-boundary.v1"

type UP51BNoisePoint struct {
	MemoryNoise float64             `json:"memory_noise"`
	Selected    MultiplicityDoseArm `json:"selected"`
	FullControl MultiplicityDoseArm `json:"full_control"`
	HeldDelta   float64             `json:"held_delta"`
	CommitDelta float64             `json:"commit_delta"`
	HardGate    bool                `json:"hard_gate"`
}

type UP51BNoiseBoundaryResult struct {
	Schema          string              `json:"schema"`
	Experiment      string              `json:"experiment"`
	SourceUP50BSeal string              `json:"source_up50b_seal"`
	FrozenPair      []int               `json:"frozen_pair"`
	HeldDepths      []int               `json:"held_depths"`
	NoiseLevels     []float64           `json:"noise_levels"`
	SelectionRun    bool                `json:"selection_run"`
	Points          []UP51BNoisePoint   `json:"points"`
	LastPassingNoise float64            `json:"last_passing_noise"`
	FirstFailingNoise float64           `json:"first_failing_noise"`
}

func RunUP51B() (UP51BNoiseBoundaryResult, error) {
	offsets, err := up50bFusedOffsets()
	if err != nil {
		return UP51BNoiseBoundaryResult{}, err
	}
	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}
	noiseLevels := []float64{0.05, 0.055, 0.06, 0.065, 0.07, 0.075, 0.08}

	result := UP51BNoiseBoundaryResult{
		Schema: UP51BNoiseBoundarySchema,
		Experiment: "UP-51B-pair05-noise-boundary",
		SourceUP50BSeal: "d74bd3a457b59bbd21532dd72c40f46e36786023",
		FrozenPair: []int{0, 5},
		HeldDepths: append([]int(nil), heldDepths...),
		NoiseLevels: append([]float64(nil), noiseLevels...),
		SelectionRun: false,
	}

	for _, noise := range noiseLevels {
		selected, err := evaluateMultiplicityDoseArm(
			"up51b_selected",
			append([]float64(nil), offsets...),
			mixer, allTrain, trueHeld,
			trainDepths, heldDepths, allDepths, noise,
		)
		if err != nil {
			return UP51BNoiseBoundaryResult{}, err
		}
		full, err := evaluateMultiplicityDoseArm(
			"up51b_full",
			[]float64{0, 0, 0, 0, 0, 0},
			mixer, allTrain, trueHeld,
			trainDepths, heldDepths, allDepths, noise,
		)
		if err != nil {
			return UP51BNoiseBoundaryResult{}, err
		}
		gate := up49bHardGate(selected, full)
		point := UP51BNoisePoint{
			MemoryNoise: noise,
			Selected: selected,
			FullControl: full,
			HeldDelta: full.Static.HeldOutAccuracy - selected.Static.HeldOutAccuracy,
			CommitDelta: full.Integration.CommitDecodeAccuracy - selected.Integration.CommitDecodeAccuracy,
			HardGate: gate,
		}
		result.Points = append(result.Points, point)
		if gate {
			result.LastPassingNoise = noise
		} else if result.FirstFailingNoise == 0 {
			result.FirstFailingNoise = noise
		}
	}

	return result, nil
}
