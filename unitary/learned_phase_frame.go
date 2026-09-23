package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const LearnedPhaseFrameSchema = "wingless.unitary-learned-phase-frame.v1"

type LearnedPhaseTrace struct {
	Outer  int       `json:"outer"`
	Phases []float64 `json:"phases"`
	Loss   float64   `json:"loss"`
}

type LearnedPhaseStaticResult struct {
	Name                string    `json:"name"`
	TrainAccuracy       float64   `json:"train_accuracy"`
	HeldOutAccuracy     float64   `json:"held_out_accuracy"`
	PerEntityTrain      []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut    []float64 `json:"per_entity_held_out_accuracy"`
}

type LearnedPhaseIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type LearnedPhaseDiagnosis struct {
	PhaseLearningPass       bool    `json:"phase_learning_pass"`
	UnseenDepthPass         bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass  bool    `json:"mutable_integration_pass"`
	InitialCapacityAccuracy float64 `json:"initial_capacity_accuracy"`
	FinalTrainAccuracy      float64 `json:"final_train_accuracy"`
	FinalHeldOutAccuracy    float64 `json:"final_heldout_accuracy"`
	MatchedControlAccuracy  float64 `json:"matched_control_accuracy"`
	InitialMinSeparation    float64 `json:"initial_min_phase_separation"`
	LearnedMinSeparation    float64 `json:"learned_min_phase_separation"`
}

type LearnedPhaseFrameProbeResult struct {
	Schema                  string                  `json:"schema"`
	Experiment              string                  `json:"experiment"`
	Dimension               int                     `json:"dimension"`
	Entities                int                     `json:"entities"`
	ValuesPerEntity         int                     `json:"values_per_entity"`
	CodeRank                int                     `json:"code_rank"`
	PilotStates             int                     `json:"pilot_states"`
	PilotComplexScalars     int                     `json:"pilot_complex_scalars"`
	FeaturesPerEntity       int                     `json:"features_per_entity"`
	RuntimePrototypeLookup  bool                    `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout  bool                    `json:"explicit_inverse_readout"`
	ExplicitDepthProvided   bool                    `json:"explicit_depth_provided"`
	GlobalPhaseNuisance     bool                    `json:"global_phase_nuisance"`
	PhaseAlphabetLearned    bool                    `json:"phase_alphabet_learned"`
	EntitySupportLearned    bool                    `json:"entity_support_learned"`
	AnchorLearned           bool                    `json:"anchor_learned"`
	MemoryNoiseAmplitude    float64                 `json:"memory_noise_amplitude"`
	TrainTables             int                     `json:"train_tables"`
	HeldOutTables           int                     `json:"held_out_tables"`
	TrainDepths             []int                   `json:"train_depths"`
	HeldOutDepths           []int                   `json:"held_out_depths"`
	InitialPhases           []float64               `json:"initial_phases"`
	LearnedPhases           []float64               `json:"learned_phases"`
	InitialLoss             float64                 `json:"initial_loss"`
	FinalLoss               float64                 `json:"final_loss"`
	Trace                   []LearnedPhaseTrace     `json:"trace"`
	Unitary                 LearnedPhaseStaticResult `json:"unitary"`
	MatchedControl          LearnedPhaseStaticResult `json:"matched_control"`
	Integration             LearnedPhaseIntegration `json:"unitary_mutable_integration"`
	Diagnosis               LearnedPhaseDiagnosis   `json:"diagnosis"`
}

func makeLearnedPhaseBank(phases []float64) (frameBank, error) {
	if len(phases) != 4 {
		return frameBank{}, fmt.Errorf("phase alphabet length=%d want=4", len(phases))
	}
	base, err := makeFrameBank()
	if err != nil {
		return frameBank{}, err
	}
	var pilots [4]State
	for entity := 0; entity < 4; entity++ {
		pilot := make(State, 16)
		for value := 0; value < 4; value++ {
			pilot[entity*4+value] = cmplx.Rect(0.5, phases[value])
		}
		pilot, err = Normalize(pilot)
		if err != nil {
			return frameBank{}, err
		}
		pilots[entity] = pilot
	}
	return frameBank{
		anchor: base.anchor,
		phase:  pilots,
	}, nil
}

func wrapPhase(value float64) float64 {
	for value <= -math.Pi {
		value += 2 * math.Pi
	}
	for value > math.Pi {
		value -= 2 * math.Pi
	}
	return value
}

func minPhaseSeparation(phases []float64) float64 {
	if len(phases) < 2 {
		return 0
	}
	minimum := math.Inf(1)
	for i := 0; i < len(phases); i++ {
		for j := i + 1; j < len(phases); j++ {
			delta := math.Abs(wrapPhase(phases[i] - phases[j]))
			if delta < minimum {
				minimum = delta
			}
		}
	}
	return minimum
}

func buildPhaseCanonicalSamples(
	tables []memoryTable,
	entity int,
	phases []float64,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	if entity < 0 || entity >= 4 || trials < 1 {
		return nil, fmt.Errorf("invalid phase canonical sample request")
	}
	bank, err := makeLearnedPhaseBank(phases)
	if err != nil {
		return nil, err
	}
	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for trial := 0; trial < trials; trial++ {
			seed := seedOffset + memoryTableIndex(table)*100 + entity*17 + trial
			state, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.173*float64(seed+1), 2*math.Pi),
			)
			features, err := frameFeature(state, bank.anchor, bank.phase[entity])
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

func updateSoftmaxHead(
	head *linearSoftmaxHead,
	samples []headSample,
	learningRate float64,
) error {
	if head == nil || len(samples) == 0 || learningRate <= 0 {
		return fmt.Errorf("invalid softmax-head update")
	}
	classes := len(head.weights)
	features := len(head.weights[0])
	gradWeights := make([][]float64, classes)
	for i := range gradWeights {
		gradWeights[i] = make([]float64, features)
	}
	gradBias := make([]float64, classes)

	for _, sample := range samples {
		if len(sample.features) != features {
			return fmt.Errorf("softmax update feature mismatch")
		}
		probabilities, err := head.probabilities(sample.features)
		if err != nil {
			return err
		}
		for class := 0; class < classes; class++ {
			delta := probabilities[class]
			if class == sample.target {
				delta -= 1
			}
			gradBias[class] += delta
			for j, feature := range sample.features {
				gradWeights[class][j] += delta * feature
			}
		}
	}

	scale := learningRate / float64(len(samples))
	for class := 0; class < classes; class++ {
		head.bias[class] -= scale * gradBias[class]
		for j := 0; j < features; j++ {
			head.weights[class][j] -= scale * gradWeights[class][j]
		}
	}
	return nil
}

func phaseHeadsLoss(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	phases []float64,
	memoryNoise float64,
	trials int,
) (float64, error) {
	var total float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildPhaseCanonicalSamples(
			tables,
			entity,
			phases,
			memoryNoise,
			trials,
			0,
		)
		if err != nil {
			return 0, err
		}
		loss, _, err := evaluateHead(heads[entity], samples)
		if err != nil {
			return 0, err
		}
		total += loss
	}
	return total / 4, nil
}

func fitPhaseHeads(
	tables []memoryTable,
	phases []float64,
	memoryNoise float64,
	trials int,
	steps int,
) ([4]linearSoftmaxHead, []float64, float64, error) {
	var heads [4]linearSoftmaxHead
	perEntity := make([]float64, 4)
	var average float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildPhaseCanonicalSamples(
			tables,
			entity,
			phases,
			memoryNoise,
			trials,
			0,
		)
		if err != nil {
			return heads, nil, 0, err
		}
		head, _, err := trainLinearSoftmax(samples, 4, 2, steps, 1.0)
		if err != nil {
			return heads, nil, 0, err
		}
		heads[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return heads, nil, 0, err
		}
		perEntity[entity] = accuracy
		average += accuracy
	}
	return heads, perEntity, average / 4, nil
}

func learnPhaseAlphabet(
	tables []memoryTable,
	memoryNoise float64,
) (
	initial []float64,
	learned []float64,
	initialLoss float64,
	finalLoss float64,
	trace []LearnedPhaseTrace,
	err error,
) {
	const (
		trials            = 2
		outerSteps        = 80
		headStepsPerOuter = 10
		headLearningRate  = 1.0
		phaseLearningRate = 0.25
		phaseEpsilon      = 1e-4
	)
	phases := []float64{0, 0.12, 0.24, 0.36}
	initial = append([]float64(nil), phases...)

	var heads [4]linearSoftmaxHead
	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(4, 2)
	}

	// Warm the heads once so the phase gradient is nonzero.
	for entity := 0; entity < 4; entity++ {
		samples, e := buildPhaseCanonicalSamples(
			tables, entity, phases, memoryNoise, trials, 0,
		)
		if e != nil {
			err = e
			return
		}
		for step := 0; step < headStepsPerOuter; step++ {
			if e := updateSoftmaxHead(
				&heads[entity], samples, headLearningRate,
			); e != nil {
				err = e
				return
			}
		}
	}
	initialLoss, err = phaseHeadsLoss(
		heads, tables, phases, memoryNoise, trials,
	)
	if err != nil {
		return
	}

	for outer := 1; outer <= outerSteps; outer++ {
		for entity := 0; entity < 4; entity++ {
			samples, e := buildPhaseCanonicalSamples(
				tables, entity, phases, memoryNoise, trials, 0,
			)
			if e != nil {
				err = e
				return
			}
			for step := 0; step < headStepsPerOuter; step++ {
				if e := updateSoftmaxHead(
					&heads[entity], samples, headLearningRate,
				); e != nil {
					err = e
					return
				}
			}
		}

		gradient := make([]float64, 4)
		for index := 1; index < 4; index++ {
			plus := append([]float64(nil), phases...)
			minus := append([]float64(nil), phases...)
			plus[index] += phaseEpsilon
			minus[index] -= phaseEpsilon

			plusLoss, e := phaseHeadsLoss(
				heads, tables, plus, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			minusLoss, e := phaseHeadsLoss(
				heads, tables, minus, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			gradient[index] =
				(plusLoss - minusLoss) / (2 * phaseEpsilon)
			if !finite(gradient[index]) {
				err = fmt.Errorf("phase gradient %d is not finite", index)
				return
			}
		}

		for index := 1; index < 4; index++ {
			phases[index] -= phaseLearningRate * gradient[index]
			phases[index] = wrapPhase(phases[index])
		}

		if outer == 1 || outer == 10 || outer == 20 ||
			outer == 40 || outer == 60 || outer == 80 {
			loss, e := phaseHeadsLoss(
				heads, tables, phases, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			trace = append(trace, LearnedPhaseTrace{
				Outer:  outer,
				Phases: append([]float64(nil), phases...),
				Loss:   loss,
			})
		}
	}

	finalLoss, err = phaseHeadsLoss(
		heads, tables, phases, memoryNoise, trials,
	)
	if err != nil {
		return
	}
	learned = append([]float64(nil), phases...)
	return
}

func buildPhaseTransportSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	phases []float64,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	baseBank, err := makeLearnedPhaseBank(phases)
	if err != nil {
		return nil, err
	}
	transformed := make(map[int]frameBank, len(depths))
	for _, depth := range depths {
		bank, err := transportFrameBank(
			baseBank, block, depth, apply, true,
		)
		if err != nil {
			return nil, err
		}
		transformed[depth] = bank
	}

	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for depthIndex, depth := range depths {
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					entity*101 +
					trial*17
				state, err := perturbMemory(
					canonical, seed, memoryNoise,
				)
				if err != nil {
					return nil, err
				}
				state = rotateGlobalPhase(
					state,
					math.Mod(0.173*float64(seed+1), 2*math.Pi),
				)
				forward, err := apply(state, block, depth)
				if err != nil {
					return nil, err
				}
				bank := transformed[depth]
				features, err := frameFeature(
					forward, bank.anchor, bank.phase[entity],
				)
				if err != nil {
					return nil, err
				}
				samples = append(samples, headSample{
					features: features,
					target:   table[entity],
				})
			}
		}
	}
	return samples, nil
}

func trainPhaseTransportHeads(
	name, path string,
	phases []float64,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	steps int,
) ([4]linearSoftmaxHead, LearnedPhaseStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildPhaseTransportSamples(
			trainTables, trainDepths, entity, phases,
			block, apply, memoryNoise, 4, 0,
		)
		if err != nil {
			return heads, LearnedPhaseStaticResult{}, err
		}
		heldSamples, err := buildPhaseTransportSamples(
			heldTables, heldDepths, entity, phases,
			block, apply, memoryNoise, 2, 7000000,
		)
		if err != nil {
			return heads, LearnedPhaseStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, LearnedPhaseStaticResult{}, err
		}
		heads[entity] = head
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return heads, LearnedPhaseStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return heads, LearnedPhaseStaticResult{}, err
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

	return heads, LearnedPhaseStaticResult{
		Name:             name,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func decodeLearnedPhaseTable(
	state State,
	bank frameBank,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := frameFeature(
			state, bank.anchor, bank.phase[entity],
		)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err := heads[entity].probabilities(features)
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

func runLearnedPhaseIntegration(
	phases []float64,
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (LearnedPhaseIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return LearnedPhaseIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return LearnedPhaseIntegration{}, err
	}

	baseBank, err := makeLearnedPhaseBank(phases)
	if err != nil {
		return LearnedPhaseIntegration{}, err
	}

	var commitCorrect, commitTotal, finalCorrect, relationCorrect int
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)
	var maxNormDrift float64

	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeFullRankScenario(
			scenarioIndex,
			heldTables[(scenarioIndex*7)%len(heldTables)],
			writes,
			depths,
		)
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			seed := 11000000 + scenarioIndex*10000 + writeIndex*31
			state, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.173*float64(seed+1), 2*math.Pi),
			)
			forward, err := applyStressUnitary(
				state, block, write.gap,
			)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			bank, err := transportFrameBank(
				baseBank, block, write.gap,
				applyStressUnitary, true,
			)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			decoded, _, margin, err := decodeLearnedPhaseTable(
				forward, bank, heads,
			)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			if margin < minValueMargin {
				minValueMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}

			pathTable, err = applyMemoryWrite(
				decoded, write.entity, write.value,
			)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return LearnedPhaseIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		seed := 11000000 + scenarioIndex*10000 + 9999
		state, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		state = rotateGlobalPhase(
			state,
			math.Mod(0.173*float64(seed+1), 2*math.Pi),
		)
		forward, err := applyStressUnitary(
			state, block, scenario.finalGap,
		)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		bank, err := transportFrameBank(
			baseBank, block, scenario.finalGap,
			applyStressUnitary, true,
		)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		decoded, distributions, margin, err := decodeLearnedPhaseTable(
			forward, bank, heads,
		)
		if err != nil {
			return LearnedPhaseIntegration{}, err
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
			return LearnedPhaseIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return LearnedPhaseIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return LearnedPhaseIntegration{
		Scenarios:               scenarios,
		WritesPerScenario:       writes,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect) / float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(scenarios),
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
		MaxNormDrift:            maxNormDrift,
	}, nil
}

// RunUP13 learns the four-value phase alphabet from task loss.
//
// The rank, anchor, and entity-local support pattern remain fixed. The four
// value phases begin deliberately clustered at 0, 0.12, 0.24, 0.36 radians.
// Three relative phase parameters are learned by alternating:
//   - deterministic softmax-head gradient steps;
//   - central-difference phase gradients through the actual classification loss.
//
// After phase learning, fresh saturated decoders are trained and evaluated at
// unseen transport depths, then used inside the mutable-memory workload.
func RunUP13() (LearnedPhaseFrameProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	initial := []float64{0, 0.12, 0.24, 0.36}
	_, _, initialCapacity, err := fitPhaseHeads(
		trainTables, initial, memoryNoise, 2, 1000,
	)
	if err != nil {
		return LearnedPhaseFrameProbeResult{}, err
	}

	initialPhases, learnedPhases, initialLoss, finalLoss, trace, err :=
		learnPhaseAlphabet(trainTables, memoryNoise)
	if err != nil {
		return LearnedPhaseFrameProbeResult{}, err
	}

	block := stressProgram()
	unitaryHeads, unitaryResult, err := trainPhaseTransportHeads(
		"unitary_learned_phase_frame",
		"unitary",
		learnedPhases,
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressUnitary,
		memoryNoise,
		1200,
	)
	if err != nil {
		return LearnedPhaseFrameProbeResult{}, err
	}

	_, controlResult, err := trainPhaseTransportHeads(
		"non_unitary_learned_phase_frame",
		"non_unitary_matched",
		learnedPhases,
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressNonUnitary,
		memoryNoise,
		1200,
	)
	if err != nil {
		return LearnedPhaseFrameProbeResult{}, err
	}

	integration, err := runLearnedPhaseIntegration(
		learnedPhases,
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return LearnedPhaseFrameProbeResult{}, err
	}

	initialSep := minPhaseSeparation(initialPhases)
	learnedSep := minPhaseSeparation(learnedPhases)
	phasePass :=
		initialCapacity < 0.70 &&
			unitaryResult.TrainAccuracy >= 0.99 &&
			learnedSep >= 0.50 &&
			finalLoss < initialLoss
	unseenPass := unitaryResult.HeldOutAccuracy >= 0.99
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95

	return LearnedPhaseFrameProbeResult{
		Schema:                 LearnedPhaseFrameSchema,
		Experiment:             "UP-13-task-learned-phase-frame",
		Dimension:              16,
		Entities:               4,
		ValuesPerEntity:        4,
		CodeRank:               4,
		PilotStates:            5,
		PilotComplexScalars:    80,
		FeaturesPerEntity:      2,
		RuntimePrototypeLookup: false,
		ExplicitInverseReadout: false,
		ExplicitDepthProvided:  false,
		GlobalPhaseNuisance:    true,
		PhaseAlphabetLearned:   true,
		EntitySupportLearned:   false,
		AnchorLearned:          false,
		MemoryNoiseAmplitude:   memoryNoise,
		TrainTables:            len(trainTables),
		HeldOutTables:          len(heldTables),
		TrainDepths:            append([]int(nil), trainDepths...),
		HeldOutDepths:          append([]int(nil), heldDepths...),
		InitialPhases:          initialPhases,
		LearnedPhases:          learnedPhases,
		InitialLoss:            initialLoss,
		FinalLoss:              finalLoss,
		Trace:                  trace,
		Unitary:                unitaryResult,
		MatchedControl:         controlResult,
		Integration:            integration,
		Diagnosis: LearnedPhaseDiagnosis{
			PhaseLearningPass:       phasePass,
			UnseenDepthPass:         unseenPass,
			MutableIntegrationPass:  mutablePass,
			InitialCapacityAccuracy: initialCapacity,
			FinalTrainAccuracy:      unitaryResult.TrainAccuracy,
			FinalHeldOutAccuracy:    unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:  controlResult.HeldOutAccuracy,
			InitialMinSeparation:    initialSep,
			LearnedMinSeparation:    learnedSep,
		},
	}, nil
}
