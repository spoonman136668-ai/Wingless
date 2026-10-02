package unitary

import (
	"fmt"
	"math/cmplx"
	"sort"
)

const ObserverAblationSchema = "wingless.unitary-observer-ablation.v1"

type ObserverAblationCell struct {
	Name               string    `json:"name"`
	Path               string    `json:"path"`
	Feature             string    `json:"feature"`
	TrainDepths         []int     `json:"train_depths"`
	HeldOutDepths       []int     `json:"held_out_depths"`
	TrainAccuracy       float64   `json:"train_accuracy"`
	HeldOutAccuracy     float64   `json:"held_out_accuracy"`
	PerEntityTrain      []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut    []float64 `json:"per_entity_held_out_accuracy"`
}

type ObserverAblationDiagnosis struct {
	BalancedSplitValid          bool    `json:"balanced_split_valid"`
	FixedCoherencePass          bool    `json:"fixed_coherence_pass"`
	FixedMagnitudePass          bool    `json:"fixed_magnitude_pass"`
	CrossDepthCoherencePass     bool    `json:"cross_depth_coherence_pass"`
	ControlFixedCoherencePass   bool    `json:"control_fixed_coherence_pass"`
	ControlCrossCoherencePass   bool    `json:"control_cross_coherence_pass"`
	MeasurementLossSupported    bool    `json:"measurement_loss_supported"`
	FrameAmbiguitySupported     bool    `json:"frame_ambiguity_supported"`
	FixedCoherenceAdvantage     float64 `json:"fixed_coherence_advantage"`
}

type ObserverAblationProbeResult struct {
	Schema                    string                    `json:"schema"`
	Experiment                string                    `json:"experiment"`
	Dimension                 int                       `json:"dimension"`
	Entities                  int                       `json:"entities"`
	ValuesPerEntity           int                       `json:"values_per_entity"`
	FullTrainPoolTables       int                       `json:"full_train_pool_tables"`
	FullHeldOutPoolTables     int                       `json:"full_held_out_pool_tables"`
	SelectedTrainTables       int                       `json:"selected_train_tables"`
	SelectedHeldOutTables     int                       `json:"selected_held_out_tables"`
	MinTrainMarginalCount     int                       `json:"min_train_marginal_count"`
	MinHeldOutMarginalCount   int                       `json:"min_held_out_marginal_count"`
	NoiseAmplitude            float64                   `json:"noise_amplitude"`
	OptimizerSteps            int                       `json:"optimizer_steps"`
	LearningRate              float64                   `json:"learning_rate"`
	EntityConditioning        string                    `json:"entity_conditioning"`
	CoherenceRepresentation   string                    `json:"coherence_representation"`
	Cells                     []ObserverAblationCell    `json:"cells"`
	Diagnosis                 ObserverAblationDiagnosis `json:"diagnosis"`
}

func balancedObserverTrainTable(table memoryTable) bool {
	sum := 0
	for _, value := range table {
		sum += value
	}
	return sum%2 == 0
}

func observerTableHash(table memoryTable) int {
	return (memoryTableIndex(table)*73 + 19) % 256
}

func selectObserverTables(train bool, count int) ([]memoryTable, error) {
	if count <= 0 {
		return nil, fmt.Errorf("observer table count must be positive")
	}
	var pool []memoryTable
	for _, table := range allMemoryTables() {
		if balancedObserverTrainTable(table) == train {
			pool = append(pool, table)
		}
	}
	sort.Slice(pool, func(i, j int) bool {
		hi := observerTableHash(pool[i])
		hj := observerTableHash(pool[j])
		if hi == hj {
			return memoryTableIndex(pool[i]) < memoryTableIndex(pool[j])
		}
		return hi < hj
	})
	if len(pool) < count {
		return nil, fmt.Errorf("observer table pool too small")
	}
	return append([]memoryTable(nil), pool[:count]...), nil
}

func observerMarginalMinimum(tables []memoryTable) int {
	minimum := len(tables) + 1
	for entity := 0; entity < 4; entity++ {
		for value := 0; value < 4; value++ {
			count := 0
			for _, table := range tables {
				if table[entity] == value {
					count++
				}
			}
			if count < minimum {
				minimum = count
			}
		}
	}
	if minimum == len(tables)+1 {
		return 0
	}
	return minimum
}

type observerFeature func(State) ([]float64, error)

func observerMagnitudeFeatures(state State) ([]float64, error) {
	return Probabilities(state)
}

// observerCoherenceFeatures returns the complete normalized pure-state density
// representation: 16 diagonal probabilities plus real/imaginary parts of all
// 120 unique off-diagonal coherences. It contains 256 real features and is
// invariant to multiplication of the complete state by a global phase.
func observerCoherenceFeatures(state State) ([]float64, error) {
	normalized, err := Normalize(state)
	if err != nil {
		return nil, err
	}
	features := make([]float64, 0, 256)
	for i := range normalized {
		features = append(features, cmplx.Abs(normalized[i])*cmplx.Abs(normalized[i]))
	}
	for i := 0; i < len(normalized); i++ {
		for j := i + 1; j < len(normalized); j++ {
			value := cmplx.Conj(normalized[i]) * normalized[j]
			features = append(features, real(value), imag(value))
		}
	}
	if len(features) != 256 {
		return nil, fmt.Errorf("coherence feature length=%d want=256", len(features))
	}
	return features, nil
}

func buildObserverEntitySamples(
	tables []memoryTable,
	depths []int,
	entity int,
	block []Coupling,
	apply stressApply,
	feature observerFeature,
	noiseAmplitude float64,
	seedOffset int,
) ([]headSample, error) {
	if entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("observer entity out of range")
	}
	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for depthIndex, depth := range depths {
			seed := seedOffset + memoryTableIndex(table)*10000 + depthIndex*101
			noisy, err := perturbMemory(canonical, seed, noiseAmplitude)
			if err != nil {
				return nil, err
			}
			forward, err := apply(noisy, block, depth)
			if err != nil {
				return nil, err
			}
			features, err := feature(forward)
			if err != nil {
				return nil, err
			}
			samples = append(samples, headSample{
				features: features,
				target:   table[entity],
			})
		}
	}
	return samples, nil
}

func runObserverAblationCell(
	name, path, featureName string,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	feature observerFeature,
	noiseAmplitude float64,
	steps int,
	learningRate float64,
) (ObserverAblationCell, error) {
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildObserverEntitySamples(
			trainTables, trainDepths, entity, block, apply, feature,
			noiseAmplitude, 0,
		)
		if err != nil {
			return ObserverAblationCell{}, err
		}
		heldSamples, err := buildObserverEntitySamples(
			heldTables, heldDepths, entity, block, apply, feature,
			noiseAmplitude, 700000,
		)
		if err != nil {
			return ObserverAblationCell{}, err
		}
		if len(trainSamples) == 0 || len(heldSamples) == 0 {
			return ObserverAblationCell{}, fmt.Errorf("observer ablation cell has empty samples")
		}
		featureCount := len(trainSamples[0].features)
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, featureCount, steps, learningRate,
		)
		if err != nil {
			return ObserverAblationCell{}, err
		}
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return ObserverAblationCell{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return ObserverAblationCell{}, err
		}
		perTrain[entity] = trainAccuracy
		perHeld[entity] = heldAccuracy
	}

	var trainAverage, heldAverage float64
	for entity := 0; entity < 4; entity++ {
		trainAverage += perTrain[entity]
		heldAverage += perHeld[entity]
	}
	trainAverage /= 4
	heldAverage /= 4

	return ObserverAblationCell{
		Name:            name,
		Path:            path,
		Feature:         featureName,
		TrainDepths:     append([]int(nil), trainDepths...),
		HeldOutDepths:   append([]int(nil), heldDepths...),
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy: heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func findObserverCell(cells []ObserverAblationCell, name string) (ObserverAblationCell, error) {
	for _, cell := range cells {
		if cell.Name == name {
			return cell, nil
		}
	}
	return ObserverAblationCell{}, fmt.Errorf("observer ablation cell %q missing", name)
}

// RunUP8 corrects the two UP-7 confounds before drawing conclusions.
//
//  1. The 128/128 split pool is based on sum(values) parity, so every entity
//     and every value appears in both partitions.
//  2. Four separately trained entity-conditioned heads are used. The observer
//     no longer relies on a linear entity one-hot that can only alter bias.
//
// It then compares magnitude-only measurement with a full global-phase-
// invariant coherence representation at a fixed frame and across unseen
// propagation depths.
func RunUP8() (ObserverAblationProbeResult, error) {
	const (
		selectedTables = 64
		noiseAmplitude = 0.05
		steps          = 100
		learningRate   = 1.0
	)
	fixedDepth := []int{128}
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}

	trainTables, err := selectObserverTables(true, selectedTables)
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	heldTables, err := selectObserverTables(false, selectedTables)
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	minTrain := observerMarginalMinimum(trainTables)
	minHeld := observerMarginalMinimum(heldTables)
	if minTrain <= 0 || minHeld <= 0 {
		return ObserverAblationProbeResult{}, fmt.Errorf("balanced observer split lost a marginal label")
	}

	block := stressProgram()
	specs := []struct {
		name        string
		path        string
		featureName string
		trainDepths []int
		heldDepths  []int
		apply       stressApply
		feature     observerFeature
	}{
		{"unitary_magnitude_fixed", "unitary", "magnitude_only", fixedDepth, fixedDepth, applyStressUnitary, observerMagnitudeFeatures},
		{"unitary_coherence_fixed", "unitary", "phase_coherence", fixedDepth, fixedDepth, applyStressUnitary, observerCoherenceFeatures},
		{"unitary_magnitude_unseen_depth", "unitary", "magnitude_only", trainDepths, heldDepths, applyStressUnitary, observerMagnitudeFeatures},
		{"unitary_coherence_unseen_depth", "unitary", "phase_coherence", trainDepths, heldDepths, applyStressUnitary, observerCoherenceFeatures},
		{"control_coherence_fixed", "non_unitary_matched", "phase_coherence", fixedDepth, fixedDepth, applyStressNonUnitary, observerCoherenceFeatures},
		{"control_coherence_unseen_depth", "non_unitary_matched", "phase_coherence", trainDepths, heldDepths, applyStressNonUnitary, observerCoherenceFeatures},
	}

	cells := make([]ObserverAblationCell, 0, len(specs))
	for _, spec := range specs {
		cell, err := runObserverAblationCell(
			spec.name,
			spec.path,
			spec.featureName,
			trainTables,
			heldTables,
			spec.trainDepths,
			spec.heldDepths,
			block,
			spec.apply,
			spec.feature,
			noiseAmplitude,
			steps,
			learningRate,
		)
		if err != nil {
			return ObserverAblationProbeResult{}, err
		}
		cells = append(cells, cell)
	}

	magFixed, err := findObserverCell(cells, "unitary_magnitude_fixed")
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	cohFixed, err := findObserverCell(cells, "unitary_coherence_fixed")
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	cohCross, err := findObserverCell(cells, "unitary_coherence_unseen_depth")
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	controlFixed, err := findObserverCell(cells, "control_coherence_fixed")
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}
	controlCross, err := findObserverCell(cells, "control_coherence_unseen_depth")
	if err != nil {
		return ObserverAblationProbeResult{}, err
	}

	advantage := cohFixed.HeldOutAccuracy - magFixed.HeldOutAccuracy
	diagnosis := ObserverAblationDiagnosis{
		BalancedSplitValid:        minTrain > 0 && minHeld > 0,
		FixedCoherencePass:        cohFixed.HeldOutAccuracy >= 0.95,
		FixedMagnitudePass:        magFixed.HeldOutAccuracy >= 0.95,
		CrossDepthCoherencePass:   cohCross.HeldOutAccuracy >= 0.95,
		ControlFixedCoherencePass: controlFixed.HeldOutAccuracy >= 0.95,
		ControlCrossCoherencePass: controlCross.HeldOutAccuracy >= 0.95,
		MeasurementLossSupported:  advantage >= 0.15,
		FrameAmbiguitySupported:   cohFixed.HeldOutAccuracy >= 0.95 && cohCross.HeldOutAccuracy < 0.95,
		FixedCoherenceAdvantage:   advantage,
	}

	return ObserverAblationProbeResult{
		Schema:                  ObserverAblationSchema,
		Experiment:              "UP-8-confound-controlled-observer-ablation",
		Dimension:               16,
		Entities:                4,
		ValuesPerEntity:         4,
		FullTrainPoolTables:     128,
		FullHeldOutPoolTables:   128,
		SelectedTrainTables:     len(trainTables),
		SelectedHeldOutTables:   len(heldTables),
		MinTrainMarginalCount:   minTrain,
		MinHeldOutMarginalCount: minHeld,
		NoiseAmplitude:          noiseAmplitude,
		OptimizerSteps:          steps,
		LearningRate:            learningRate,
		EntityConditioning:      "four_separate_entity_heads",
		CoherenceRepresentation: "normalized_density_matrix_global_phase_invariant",
		Cells:                   cells,
		Diagnosis:               diagnosis,
	}, nil
}
