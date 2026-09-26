package unitary

import (
	"fmt"
	"math"
)

const LearnedSupportFrameSchema = "wingless.unitary-learned-support-frame.v1"

type SupportMatrix [4][4]float64

type LearnedSupportTrace struct {
	Outer   int           `json:"outer"`
	Matrix  SupportMatrix `json:"matrix"`
	Loss    float64       `json:"loss"`
	Purities []float64    `json:"target_support_purities"`
}

type LearnedSupportStaticResult struct {
	Name             string    `json:"name"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
}

type LearnedSupportIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type LearnedSupportDiagnosis struct {
	SupportLearningPass       bool    `json:"support_learning_pass"`
	UnseenDepthPass           bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass    bool    `json:"mutable_integration_pass"`
	InitialCapacityAccuracy   float64 `json:"initial_capacity_accuracy"`
	FinalTrainAccuracy        float64 `json:"final_train_accuracy"`
	FinalHeldOutAccuracy      float64 `json:"final_heldout_accuracy"`
	MatchedControlAccuracy    float64 `json:"matched_control_accuracy"`
	InitialMeanTargetPurity   float64 `json:"initial_mean_target_purity"`
	LearnedMeanTargetPurity   float64 `json:"learned_mean_target_purity"`
	MinimumLearnedTargetPurity float64 `json:"minimum_learned_target_purity"`
}

type LearnedSupportFrameProbeResult struct {
	Schema                  string                    `json:"schema"`
	Experiment              string                    `json:"experiment"`
	Dimension               int                       `json:"dimension"`
	Entities                int                       `json:"entities"`
	ValuesPerEntity         int                       `json:"values_per_entity"`
	CodeRank                int                       `json:"code_rank"`
	PilotStates             int                       `json:"pilot_states"`
	PilotComplexScalars     int                       `json:"pilot_complex_scalars"`
	FeaturesPerEntity       int                       `json:"features_per_entity"`
	RuntimePrototypeLookup  bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout  bool                      `json:"explicit_inverse_readout"`
	ExplicitDepthProvided   bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance     bool                      `json:"global_phase_nuisance"`
	PhaseAlphabetLearned    bool                      `json:"phase_alphabet_learned"`
	EntitySupportLearned    bool                      `json:"entity_support_learned"`
	AnchorLearned           bool                      `json:"anchor_learned"`
	MemoryNoiseAmplitude    float64                   `json:"memory_noise_amplitude"`
	TrainTables             int                       `json:"train_tables"`
	HeldOutTables           int                       `json:"held_out_tables"`
	TrainDepths             []int                     `json:"train_depths"`
	HeldOutDepths           []int                     `json:"held_out_depths"`
	FixedPhases             []float64                 `json:"fixed_phases_from_up13"`
	InitialSupport          SupportMatrix             `json:"initial_support"`
	LearnedSupport          SupportMatrix             `json:"learned_support"`
	InitialLoss             float64                   `json:"initial_loss"`
	FinalLoss               float64                   `json:"final_loss"`
	Trace                   []LearnedSupportTrace     `json:"trace"`
	Unitary                 LearnedSupportStaticResult `json:"unitary"`
	MatchedControl          LearnedSupportStaticResult `json:"matched_control"`
	Integration             LearnedSupportIntegration `json:"unitary_mutable_integration"`
	Diagnosis               LearnedSupportDiagnosis   `json:"diagnosis"`
}

func up13LearnedPhases() []float64 {
	return []float64{
		0,
		-1.7700278912430298,
		1.1656136621055615,
		2.5917033093902164,
	}
}

func initialSupportMatrix() SupportMatrix {
	var matrix SupportMatrix
	for row := 0; row < 4; row++ {
		for column := 0; column < 4; column++ {
			matrix[row][column] = 1
		}
		matrix[row][row] += 0.04
	}
	return matrix
}

func normalizeSupportRow(row *[4]float64) error {
	if row == nil {
		return fmt.Errorf("support row is nil")
	}
	var norm2 float64
	for _, value := range row {
		if !finite(value) {
			return fmt.Errorf("support row contains non-finite value")
		}
		norm2 += value * value
	}
	if norm2 <= 1e-18 {
		return fmt.Errorf("support row norm too small")
	}
	scale := 1 / math.Sqrt(norm2)
	for i := range row {
		row[i] *= scale
	}
	return nil
}

func normalizedSupportMatrix(matrix SupportMatrix) (SupportMatrix, error) {
	out := matrix
	for row := 0; row < 4; row++ {
		if err := normalizeSupportRow(&out[row]); err != nil {
			return SupportMatrix{}, err
		}
	}
	return out, nil
}

func supportPurities(matrix SupportMatrix) []float64 {
	out := make([]float64, 4)
	for row := 0; row < 4; row++ {
		var total float64
		for column := 0; column < 4; column++ {
			total += matrix[row][column] * matrix[row][column]
		}
		if total > 0 {
			out[row] =
				matrix[row][row] * matrix[row][row] / total
		}
	}
	return out
}

func meanAndMin(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	minimum := values[0]
	var total float64
	for _, value := range values {
		total += value
		if value < minimum {
			minimum = value
		}
	}
	return total / float64(len(values)), minimum
}

func makeLearnedSupportBank(
	phases []float64,
	matrix SupportMatrix,
) (frameBank, error) {
	matrix, err := normalizedSupportMatrix(matrix)
	if err != nil {
		return frameBank{}, err
	}
	base, err := makeLearnedPhaseBank(phases)
	if err != nil {
		return frameBank{}, err
	}

	var pilots [4]State
	for row := 0; row < 4; row++ {
		pilot := make(State, 16)
		for entity := 0; entity < 4; entity++ {
			coefficient := complex(matrix[row][entity], 0)
			for i := range pilot {
				pilot[i] += coefficient * base.phase[entity][i]
			}
		}
		pilot, err = Normalize(pilot)
		if err != nil {
			return frameBank{}, err
		}
		pilots[row] = pilot
	}
	return frameBank{
		anchor: base.anchor,
		phase:  pilots,
	}, nil
}

func buildSupportCanonicalSamples(
	tables []memoryTable,
	entity int,
	phases []float64,
	matrix SupportMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	bank, err := makeLearnedSupportBank(phases, matrix)
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
			seed := seedOffset +
				memoryTableIndex(table)*100 +
				entity*17 +
				trial
			state, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return nil, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.181*float64(seed+1), 2*math.Pi),
			)
			features, err := frameFeature(
				state, bank.anchor, bank.phase[entity],
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
	return samples, nil
}

func supportHeadLoss(
	head linearSoftmaxHead,
	tables []memoryTable,
	entity int,
	phases []float64,
	matrix SupportMatrix,
	memoryNoise float64,
	trials int,
) (float64, error) {
	samples, err := buildSupportCanonicalSamples(
		tables, entity, phases, matrix,
		memoryNoise, trials, 0,
	)
	if err != nil {
		return 0, err
	}
	loss, _, err := evaluateHead(head, samples)
	return loss, err
}

func supportHeadsLoss(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	phases []float64,
	matrix SupportMatrix,
	memoryNoise float64,
	trials int,
) (float64, error) {
	var total float64
	for entity := 0; entity < 4; entity++ {
		loss, err := supportHeadLoss(
			heads[entity], tables, entity,
			phases, matrix, memoryNoise, trials,
		)
		if err != nil {
			return 0, err
		}
		total += loss
	}
	return total / 4, nil
}

func fitSupportHeads(
	tables []memoryTable,
	phases []float64,
	matrix SupportMatrix,
	memoryNoise float64,
	trials int,
	steps int,
) ([4]linearSoftmaxHead, float64, error) {
	var heads [4]linearSoftmaxHead
	var average float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildSupportCanonicalSamples(
			tables, entity, phases, matrix,
			memoryNoise, trials, 0,
		)
		if err != nil {
			return heads, 0, err
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, 0, err
		}
		heads[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return heads, 0, err
		}
		average += accuracy
	}
	return heads, average / 4, nil
}

func learnSupportMatrix(
	tables []memoryTable,
	phases []float64,
	memoryNoise float64,
) (
	initial SupportMatrix,
	learned SupportMatrix,
	initialLoss float64,
	finalLoss float64,
	trace []LearnedSupportTrace,
	err error,
) {
	const (
		trials            = 2
		outerSteps        = 100
		headStepsPerOuter = 10
		headLearningRate  = 1.0
		supportLearningRate = 0.35
		supportEpsilon      = 1e-4
	)
	matrix, err := normalizedSupportMatrix(initialSupportMatrix())
	if err != nil {
		return SupportMatrix{}, SupportMatrix{}, 0, 0, nil, err
	}
	initial = matrix

	var heads [4]linearSoftmaxHead
	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(4, 2)
	}

	for entity := 0; entity < 4; entity++ {
		samples, e := buildSupportCanonicalSamples(
			tables, entity, phases, matrix,
			memoryNoise, trials, 0,
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
	initialLoss, err = supportHeadsLoss(
		heads, tables, phases, matrix,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}

	for outer := 1; outer <= outerSteps; outer++ {
		for entity := 0; entity < 4; entity++ {
			samples, e := buildSupportCanonicalSamples(
				tables, entity, phases, matrix,
				memoryNoise, trials, 0,
			)
			if e != nil {
				err = e
				return
			}
			for step := 0; step < headStepsPerOuter; step++ {
				if e := updateSoftmaxHead(
					&heads[entity], samples,
					headLearningRate,
				); e != nil {
					err = e
					return
				}
			}
		}

		for row := 0; row < 4; row++ {
			gradient := [4]float64{}
			for column := 0; column < 4; column++ {
				plus := matrix
				minus := matrix
				plus[row][column] += supportEpsilon
				minus[row][column] -= supportEpsilon

				plusLoss, e := supportHeadLoss(
					heads[row], tables, row, phases,
					plus, memoryNoise, trials,
				)
				if e != nil {
					err = e
					return
				}
				minusLoss, e := supportHeadLoss(
					heads[row], tables, row, phases,
					minus, memoryNoise, trials,
				)
				if e != nil {
					err = e
					return
				}
				gradient[column] =
					(plusLoss - minusLoss) /
						(2 * supportEpsilon)
				if !finite(gradient[column]) {
					err = fmt.Errorf(
						"support gradient row=%d col=%d non-finite",
						row, column,
					)
					return
				}
			}
			for column := 0; column < 4; column++ {
				matrix[row][column] -=
					supportLearningRate * gradient[column]
			}
			if e := normalizeSupportRow(&matrix[row]); e != nil {
				err = e
				return
			}
		}

		if outer == 1 || outer == 10 || outer == 20 ||
			outer == 40 || outer == 60 ||
			outer == 80 || outer == 100 {
			loss, e := supportHeadsLoss(
				heads, tables, phases, matrix,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			trace = append(trace, LearnedSupportTrace{
				Outer:    outer,
				Matrix:   matrix,
				Loss:     loss,
				Purities: supportPurities(matrix),
			})
		}
	}

	finalLoss, err = supportHeadsLoss(
		heads, tables, phases, matrix,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}
	learned = matrix
	return
}

func buildSupportTransportSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	phases []float64,
	matrix SupportMatrix,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	baseBank, err := makeLearnedSupportBank(phases, matrix)
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
					math.Mod(0.181*float64(seed+1), 2*math.Pi),
				)
				forward, err := apply(state, block, depth)
				if err != nil {
					return nil, err
				}
				bank := transformed[depth]
				features, err := frameFeature(
					forward,
					bank.anchor,
					bank.phase[entity],
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

func trainSupportTransportHeads(
	name string,
	phases []float64,
	matrix SupportMatrix,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	steps int,
) ([4]linearSoftmaxHead, LearnedSupportStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildSupportTransportSamples(
			trainTables, trainDepths, entity,
			phases, matrix, block, apply,
			memoryNoise, 4, 0,
		)
		if err != nil {
			return heads, LearnedSupportStaticResult{}, err
		}
		heldSamples, err := buildSupportTransportSamples(
			heldTables, heldDepths, entity,
			phases, matrix, block, apply,
			memoryNoise, 2, 7000000,
		)
		if err != nil {
			return heads, LearnedSupportStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, LearnedSupportStaticResult{}, err
		}
		heads[entity] = head
		_, trainAccuracy, err := evaluateHead(
			head, trainSamples,
		)
		if err != nil {
			return heads, LearnedSupportStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(
			head, heldSamples,
		)
		if err != nil {
			return heads, LearnedSupportStaticResult{}, err
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

	return heads, LearnedSupportStaticResult{
		Name:             name,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func decodeSupportTable(
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
		probabilities, err :=
			heads[entity].probabilities(features)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		value, margin, err :=
			classAndMargin(probabilities)
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

func runLearnedSupportIntegration(
	phases []float64,
	matrix SupportMatrix,
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (LearnedSupportIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return LearnedSupportIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return LearnedSupportIntegration{}, err
	}
	baseBank, err := makeLearnedSupportBank(phases, matrix)
	if err != nil {
		return LearnedSupportIntegration{}, err
	}

	var commitCorrect, commitTotal int
	var finalCorrect, relationCorrect int
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
				return LearnedSupportIntegration{}, err
			}
			seed := 13000000 +
				scenarioIndex*10000 +
				writeIndex*31
			state, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return LearnedSupportIntegration{}, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.181*float64(seed+1), 2*math.Pi),
			)
			forward, err := applyStressUnitary(
				state, block, write.gap,
			)
			if err != nil {
				return LearnedSupportIntegration{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return LearnedSupportIntegration{}, err
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
				return LearnedSupportIntegration{}, err
			}
			decoded, _, margin, err := decodeSupportTable(
				forward, bank, heads,
			)
			if err != nil {
				return LearnedSupportIntegration{}, err
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
				return LearnedSupportIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return LearnedSupportIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		seed := 13000000 + scenarioIndex*10000 + 9999
		state, err := perturbMemory(
			canonical, seed, memoryNoise,
		)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		state = rotateGlobalPhase(
			state,
			math.Mod(0.181*float64(seed+1), 2*math.Pi),
		)
		forward, err := applyStressUnitary(
			state, block, scenario.finalGap,
		)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return LearnedSupportIntegration{}, err
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
			return LearnedSupportIntegration{}, err
		}
		decoded, distributions, margin, err := decodeSupportTable(
			forward, bank, heads,
		)
		if err != nil {
			return LearnedSupportIntegration{}, err
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
			return LearnedSupportIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return LearnedSupportIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return LearnedSupportIntegration{
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

// RunUP14 learns which entity block each pilot should reference.
//
// The rank-4 structure, UP-13 learned phase alphabet, and global anchor remain
// fixed. Each pilot begins as an almost-uniform dense mixture of all four
// entity-local phase templates. Task loss then updates a real 4x4 support
// matrix. No target support matrix is supplied.
//
// Because each entity head sees only its corresponding pilot/anchor
// correlation, a mixed pilot is contaminated by three independently varying
// entities. Successful optimization must make each row selective enough to
// recover its queried entity.
func RunUP14() (LearnedSupportFrameProbeResult, error) {
	const memoryNoise = 0.05
	phases := up13LearnedPhases()
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	initialMatrix, err := normalizedSupportMatrix(
		initialSupportMatrix(),
	)
	if err != nil {
		return LearnedSupportFrameProbeResult{}, err
	}
	_, initialCapacity, err := fitSupportHeads(
		trainTables, phases, initialMatrix,
		memoryNoise, 2, 1000,
	)
	if err != nil {
		return LearnedSupportFrameProbeResult{}, err
	}

	initial, learned, initialLoss, finalLoss, trace, err :=
		learnSupportMatrix(
			trainTables, phases, memoryNoise,
		)
	if err != nil {
		return LearnedSupportFrameProbeResult{}, err
	}

	block := stressProgram()
	unitaryHeads, unitaryResult, err :=
		trainSupportTransportHeads(
			"unitary_learned_support_frame",
			phases,
			learned,
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
		return LearnedSupportFrameProbeResult{}, err
	}

	_, controlResult, err :=
		trainSupportTransportHeads(
			"non_unitary_learned_support_frame",
			phases,
			learned,
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
		return LearnedSupportFrameProbeResult{}, err
	}

	integration, err := runLearnedSupportIntegration(
		phases,
		learned,
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return LearnedSupportFrameProbeResult{}, err
	}

	initialPurities := supportPurities(initial)
	learnedPurities := supportPurities(learned)
	initialMean, _ := meanAndMin(initialPurities)
	learnedMean, learnedMin := meanAndMin(learnedPurities)

	supportPass :=
		initialCapacity < 0.70 &&
			unitaryResult.TrainAccuracy >= 0.99 &&
			learnedMean >= 0.85 &&
			learnedMin >= 0.75 &&
			finalLoss < initialLoss
	unseenPass := unitaryResult.HeldOutAccuracy >= 0.99
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95

	return LearnedSupportFrameProbeResult{
		Schema:                 LearnedSupportFrameSchema,
		Experiment:             "UP-14-task-learned-pilot-support",
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
		PhaseAlphabetLearned:   false,
		EntitySupportLearned:   true,
		AnchorLearned:          false,
		MemoryNoiseAmplitude:   memoryNoise,
		TrainTables:            len(trainTables),
		HeldOutTables:          len(heldTables),
		TrainDepths:            append([]int(nil), trainDepths...),
		HeldOutDepths:          append([]int(nil), heldDepths...),
		FixedPhases:            append([]float64(nil), phases...),
		InitialSupport:         initial,
		LearnedSupport:         learned,
		InitialLoss:            initialLoss,
		FinalLoss:              finalLoss,
		Trace:                  trace,
		Unitary:                unitaryResult,
		MatchedControl:         controlResult,
		Integration:            integration,
		Diagnosis: LearnedSupportDiagnosis{
			SupportLearningPass:        supportPass,
			UnseenDepthPass:            unseenPass,
			MutableIntegrationPass:     mutablePass,
			InitialCapacityAccuracy:    initialCapacity,
			FinalTrainAccuracy:         unitaryResult.TrainAccuracy,
			FinalHeldOutAccuracy:       unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:     controlResult.HeldOutAccuracy,
			InitialMeanTargetPurity:    initialMean,
			LearnedMeanTargetPurity:    learnedMean,
			MinimumLearnedTargetPurity: learnedMin,
		},
	}, nil
}
