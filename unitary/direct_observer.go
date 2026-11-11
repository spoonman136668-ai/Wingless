package unitary

import (
	"fmt"
	"math"
)

const DirectObserverProbeSchema = "wingless.unitary-direct-observer-probe.v1"

type DirectObserverHeadMetric struct {
	Name             string  `json:"name"`
	TrainAccuracy    float64 `json:"train_accuracy"`
	TrainLoss        float64 `json:"train_loss"`
	StaticHeldOutAcc float64 `json:"static_held_out_accuracy"`
}

type DirectObserverPathResult struct {
	Name                    string  `json:"name"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxForwardNormDrift     float64 `json:"max_forward_norm_drift"`
	HypothesisPass          bool    `json:"hypothesis_pass"`
}

type DirectObserverProbeResult struct {
	Schema                 string                   `json:"schema"`
	Experiment             string                   `json:"experiment"`
	Dimension              int                      `json:"dimension"`
	Entities               int                      `json:"entities"`
	ValuesPerEntity        int                      `json:"values_per_entity"`
	RuntimePrototypeLookup bool                     `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout bool                     `json:"explicit_inverse_readout"`
	DepthFeatureProvided   bool                     `json:"depth_feature_provided"`
	TrainTables            int                      `json:"train_tables"`
	HeldOutTables          int                      `json:"held_out_tables"`
	TrainDepths            []int                    `json:"train_depths"`
	HeldOutDepths          []int                    `json:"held_out_depths"`
	Scenarios              int                      `json:"scenarios"`
	WritesPerScenario      int                      `json:"writes_per_scenario"`
	NoiseAmplitude         float64                  `json:"noise_amplitude"`
	UnitaryObserver        DirectObserverHeadMetric `json:"unitary_observer"`
	ControlObserver        DirectObserverHeadMetric `json:"control_observer"`
	RelationHead           LearnedHeadMetric        `json:"relation_head"`
	Unitary                DirectObserverPathResult `json:"unitary"`
	NonUnitary             DirectObserverPathResult `json:"non_unitary_matched"`
}

func memoryTableIndex(table memoryTable) int {
	return (((table[0]*4)+table[1])*4+table[2])*4 + table[3]
}

func directTrainTable(table memoryTable) bool {
	return memoryTableIndex(table)%2 == 0
}

func directFeatures(state State, entity int) ([]float64, error) {
	if entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("direct observer entity out of range")
	}
	probabilities, err := Probabilities(state)
	if err != nil {
		return nil, err
	}
	if len(probabilities) != 16 {
		return nil, fmt.Errorf("direct observer state dimension=%d want=16", len(probabilities))
	}
	features := make([]float64, 20)
	copy(features, probabilities)
	features[16+entity] = 1
	return features, nil
}

func buildDirectObserverSamples(
	block []Coupling,
	depths []int,
	apply stressApply,
	wantTrainTables bool,
	noiseAmplitude float64,
	trials int,
) ([]headSample, int, error) {
	if trials < 1 {
		return nil, 0, fmt.Errorf("direct observer trials must be positive")
	}
	var samples []headSample
	tableCount := 0
	for _, table := range allMemoryTables() {
		if directTrainTable(table) != wantTrainTables {
			continue
		}
		tableCount++
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			for trial := 0; trial < trials; trial++ {
				seed := memoryTableIndex(table)*10000 + depthIndex*100 + trial
				noisy, err := perturbMemory(canonical, seed, noiseAmplitude)
				if err != nil {
					return nil, 0, err
				}
				forward, err := apply(noisy, block, depth)
				if err != nil {
					return nil, 0, err
				}
				for entity := 0; entity < 4; entity++ {
					features, err := directFeatures(forward, entity)
					if err != nil {
						return nil, 0, err
					}
					samples = append(samples, headSample{
						features: features,
						target:   table[entity],
					})
				}
			}
		}
	}
	return samples, tableCount, nil
}

func trainDirectObserver(
	name string,
	block []Coupling,
	trainDepths, heldDepths []int,
	apply stressApply,
	noiseAmplitude float64,
) (linearSoftmaxHead, DirectObserverHeadMetric, int, int, error) {
	trainSamples, trainTables, err := buildDirectObserverSamples(
		block, trainDepths, apply, true, noiseAmplitude, 2,
	)
	if err != nil {
		return linearSoftmaxHead{}, DirectObserverHeadMetric{}, 0, 0, err
	}
	heldSamples, heldTables, err := buildDirectObserverSamples(
		block, heldDepths, apply, false, noiseAmplitude, 2,
	)
	if err != nil {
		return linearSoftmaxHead{}, DirectObserverHeadMetric{}, 0, 0, err
	}

	head, trainMetric, err := trainLinearSoftmax(trainSamples, 4, 20, 500, 1.0)
	if err != nil {
		return linearSoftmaxHead{}, DirectObserverHeadMetric{}, 0, 0, err
	}
	_, heldAccuracy, err := evaluateHead(head, heldSamples)
	if err != nil {
		return linearSoftmaxHead{}, DirectObserverHeadMetric{}, 0, 0, err
	}

	return head, DirectObserverHeadMetric{
		Name:             name,
		TrainAccuracy:    trainMetric.TrainAccuracy,
		TrainLoss:        trainMetric.TrainLoss,
		StaticHeldOutAcc: heldAccuracy,
	}, trainTables, heldTables, nil
}

func decodeDirectTable(state State, head linearSoftmaxHead) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := directFeatures(state, entity)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err := head.probabilities(features)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		value, margin, err := classAndMargin(probabilities)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		table[entity] = value
		distributions[entity] = probabilities
		if margin < minMargin {
			minMargin = margin
		}
	}
	return table, distributions, minMargin, nil
}

func directHeldOutInitial(seed int) memoryTable {
	// Force an odd table index so the initial table is outside observer training.
	return memoryTable{
		(seed + 1) % 4,
		(seed*2 + 3) % 4,
		(seed*3 + 2) % 4,
		1,
	}
}

func makeDirectObserverScenario(seed, writes int, heldDepths []int) (memoryScenario, error) {
	initial := directHeldOutInitial(seed)
	if directTrainTable(initial) {
		return memoryScenario{}, fmt.Errorf("direct observer scenario initial table leaked into training split")
	}
	ops := make([]memoryWrite, 0, writes)
	for j := 0; j < writes; j++ {
		entity := (seed*3 + j + j/2) % 4
		value := (seed + j*j + 3*j + entity) % 4
		gap := heldDepths[(seed+j*2)%len(heldDepths)]
		ops = append(ops, memoryWrite{entity: entity, value: value, gap: gap})
	}
	return memoryScenario{
		initial:  initial,
		writes:   ops,
		queryA:   (seed + 3) % 4,
		queryB:   (seed + 1) % 4,
		finalGap: heldDepths[(seed+writes)%len(heldDepths)],
	}, nil
}

func runDirectObserverPath(
	name string,
	scenarios []memoryScenario,
	block []Coupling,
	apply stressApply,
	observer linearSoftmaxHead,
	relationHead linearSoftmaxHead,
	noiseAmplitude float64,
	staticHeldOutAccuracy float64,
) (DirectObserverPathResult, error) {
	var commitCorrect, commitTotal int
	var finalCorrect, relationCorrect int
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)
	var maxNormDrift float64

	for scenarioIndex, scenario := range scenarios {
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			noisy, err := perturbMemory(canonical, 900000+scenarioIndex*1000+writeIndex, noiseAmplitude)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			forward, err := apply(noisy, block, write.gap)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeDirectTable(forward, observer)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			if margin < minValueMargin {
				minValueMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}

			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return DirectObserverPathResult{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		noisy, err := perturbMemory(canonical, 900000+scenarioIndex*1000+999, noiseAmplitude)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		forward, err := apply(noisy, block, scenario.finalGap)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		decoded, distributions, margin, err := decodeDirectTable(forward, observer)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		if margin < minValueMargin {
			minValueMargin = margin
		}
		if decoded == trueTable {
			finalCorrect++
		}

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil {
			return DirectObserverPathResult{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	if commitTotal == 0 || len(scenarios) == 0 || !finite(minValueMargin) || !finite(minRelationMargin) {
		return DirectObserverPathResult{}, fmt.Errorf("direct observer produced no evaluable results")
	}

	commitAccuracy := float64(commitCorrect) / float64(commitTotal)
	finalAccuracy := float64(finalCorrect) / float64(len(scenarios))
	relationAccuracy := float64(relationCorrect) / float64(len(scenarios))
	hypothesisPass := staticHeldOutAccuracy >= 0.95 &&
		commitAccuracy >= 0.95 &&
		finalAccuracy >= 0.90 &&
		relationAccuracy >= 0.90

	return DirectObserverPathResult{
		Name:                    name,
		CommitDecodeAccuracy:    commitAccuracy,
		ExactFinalTableAccuracy: finalAccuracy,
		RelationalQueryAccuracy: relationAccuracy,
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
		MaxForwardNormDrift:     maxNormDrift,
		HypothesisPass:          hypothesisPass,
	}, nil
}

// RunUP7 removes explicit inverse transport from observation.
//
// Each path gets an independently trained but architecture- and budget-matched
// linear-softmax observer. The observer receives only normalized transported
// coordinate probabilities plus an entity selector. It receives no transport
// depth, no inverse state, and no prototype table.
//
// Training uses one half of the 256 memory tables and six transport depths.
// Static evaluation uses the disjoint table half and four unseen depths. The
// learned observer is then used inside unseen mutable write programs.
func RunUP7() (DirectObserverProbeResult, error) {
	const (
		scenarioCount       = 32
		writesPerScenario   = 12
		noiseAmplitude      = 0.05
	)
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	block := stressProgram()

	unitaryObserver, unitaryMetric, trainTables, heldTables, err := trainDirectObserver(
		"unitary_direct_observer",
		block,
		trainDepths,
		heldDepths,
		applyStressUnitary,
		noiseAmplitude,
	)
	if err != nil {
		return DirectObserverProbeResult{}, err
	}
	controlObserver, controlMetric, controlTrainTables, controlHeldTables, err := trainDirectObserver(
		"non_unitary_direct_observer",
		block,
		trainDepths,
		heldDepths,
		applyStressNonUnitary,
		noiseAmplitude,
	)
	if err != nil {
		return DirectObserverProbeResult{}, err
	}
	if trainTables != controlTrainTables || heldTables != controlHeldTables {
		return DirectObserverProbeResult{}, fmt.Errorf("direct observer split mismatch")
	}

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return DirectObserverProbeResult{}, err
	}
	relationHead, relationMetric, err := trainLinearSoftmax(relationSamples, 4, 16, 400, 1.0)
	if err != nil {
		return DirectObserverProbeResult{}, err
	}
	relationMetric.Name = "shared_relation_head"

	scenarios := make([]memoryScenario, 0, scenarioCount)
	for seed := 0; seed < scenarioCount; seed++ {
		scenario, err := makeDirectObserverScenario(seed, writesPerScenario, heldDepths)
		if err != nil {
			return DirectObserverProbeResult{}, err
		}
		scenarios = append(scenarios, scenario)
	}

	unitaryResult, err := runDirectObserverPath(
		"unitary_transport",
		scenarios,
		block,
		applyStressUnitary,
		unitaryObserver,
		relationHead,
		noiseAmplitude,
		unitaryMetric.StaticHeldOutAcc,
	)
	if err != nil {
		return DirectObserverProbeResult{}, err
	}
	controlResult, err := runDirectObserverPath(
		"non_unitary_matched_transport",
		scenarios,
		block,
		applyStressNonUnitary,
		controlObserver,
		relationHead,
		noiseAmplitude,
		controlMetric.StaticHeldOutAcc,
	)
	if err != nil {
		return DirectObserverProbeResult{}, err
	}

	return DirectObserverProbeResult{
		Schema:                 DirectObserverProbeSchema,
		Experiment:             "UP-7-direct-transported-state-observer",
		Dimension:              16,
		Entities:               4,
		ValuesPerEntity:        4,
		RuntimePrototypeLookup: false,
		ExplicitInverseReadout: false,
		DepthFeatureProvided:   false,
		TrainTables:            trainTables,
		HeldOutTables:          heldTables,
		TrainDepths:            append([]int(nil), trainDepths...),
		HeldOutDepths:          append([]int(nil), heldDepths...),
		Scenarios:              scenarioCount,
		WritesPerScenario:      writesPerScenario,
		NoiseAmplitude:         noiseAmplitude,
		UnitaryObserver:        unitaryMetric,
		ControlObserver:        controlMetric,
		RelationHead:           relationMetric,
		Unitary:                unitaryResult,
		NonUnitary:             controlResult,
	}, nil
}
