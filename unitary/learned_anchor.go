package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const LearnedAnchorSchema = "wingless.unitary-learned-anchor.v1"

type AnchorParams struct {
	Real [16]float64 `json:"real"`
	Imag [16]float64 `json:"imag"`
}

type LearnedAnchorTrace struct {
	Outer              int          `json:"outer"`
	Loss               float64      `json:"loss"`
	PhaseConcentration float64      `json:"table_phase_concentration"`
	Params             AnchorParams `json:"params"`
}

type LearnedAnchorStaticResult struct {
	Name             string    `json:"name"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
}

type LearnedAnchorIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type LearnedAnchorDiagnosis struct {
	AnchorLearningPass          bool    `json:"anchor_learning_pass"`
	UnseenDepthPass             bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass      bool    `json:"mutable_integration_pass"`
	InitialCapacityAccuracy     float64 `json:"initial_capacity_accuracy"`
	FinalTrainAccuracy          float64 `json:"final_train_accuracy"`
	FinalHeldOutAccuracy        float64 `json:"final_heldout_accuracy"`
	MatchedControlAccuracy      float64 `json:"matched_control_accuracy"`
	InitialPhaseConcentration   float64 `json:"initial_phase_concentration"`
	LearnedPhaseConcentration   float64 `json:"learned_phase_concentration"`
	AnchorStateL2Shift          float64 `json:"anchor_state_l2_shift"`
}

type LearnedAnchorProbeResult struct {
	Schema                  string                   `json:"schema"`
	Experiment              string                   `json:"experiment"`
	Dimension               int                      `json:"dimension"`
	Entities                int                      `json:"entities"`
	ValuesPerEntity         int                      `json:"values_per_entity"`
	CodeRank                int                      `json:"code_rank"`
	PilotStates             int                      `json:"pilot_states"`
	PilotComplexScalars     int                      `json:"pilot_complex_scalars"`
	FeaturesPerEntity       int                      `json:"features_per_entity"`
	RuntimePrototypeLookup  bool                     `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout  bool                     `json:"explicit_inverse_readout"`
	ExplicitDepthProvided   bool                     `json:"explicit_depth_provided"`
	GlobalPhaseNuisance     bool                     `json:"global_phase_nuisance"`
	PhaseAlphabetLearned    bool                     `json:"phase_alphabet_learned"`
	EntitySupportLearned    bool                     `json:"entity_support_learned"`
	AnchorLearned           bool                     `json:"anchor_learned"`
	MemoryNoiseAmplitude    float64                  `json:"memory_noise_amplitude"`
	TrainTables             int                      `json:"train_tables"`
	HeldOutTables           int                      `json:"held_out_tables"`
	TrainDepths             []int                    `json:"train_depths"`
	HeldOutDepths           []int                    `json:"held_out_depths"`
	FixedPhases             []float64                `json:"fixed_phases_from_up13"`
	FixedSupport            SupportMatrix            `json:"fixed_support_from_up14"`
	InitialAnchor           AnchorParams             `json:"initial_anchor"`
	LearnedAnchor           AnchorParams             `json:"learned_anchor"`
	InitialLoss             float64                  `json:"initial_loss"`
	FinalLoss               float64                  `json:"final_loss"`
	Trace                   []LearnedAnchorTrace     `json:"trace"`
	Unitary                 LearnedAnchorStaticResult `json:"unitary"`
	MatchedControl          LearnedAnchorStaticResult `json:"matched_control"`
	Integration             LearnedAnchorIntegration `json:"unitary_mutable_integration"`
	Diagnosis               LearnedAnchorDiagnosis   `json:"diagnosis"`
}

func up14LearnedSupportMatrix() SupportMatrix {
	return SupportMatrix{
		{0.9998840473574464, -0.013665962603931539, 0.0045093021552159, -0.004979909665126646},
		{0.006964963205448895, 0.9997658581171776, -0.009805998087682916, 0.017987791200561012},
		{0.002464228556381514, -0.002966576558468774, 0.9999072087395744, -0.013065255901730488},
		{-0.012433390948591277, 0.014598737076627758, -0.009138487445338501, 0.999774362400086},
	}
}

func initialLearnedAnchorParams() AnchorParams {
	var out AnchorParams
	for i := 0; i < 16; i++ {
		amp := 0.65 + 0.05*float64((i*7+3)%7)
		phase := -2.15 + 0.43*float64((i*5+2)%13)
		out.Real[i] = amp * math.Cos(phase)
		out.Imag[i] = amp * math.Sin(phase)
	}
	normalized, _ := canonicalizeAnchorParams(out)
	return normalized
}

func anchorStateFromParams(params AnchorParams) (State, error) {
	state := make(State, 16)
	for i := range state {
		state[i] = complex(params.Real[i], params.Imag[i])
	}
	state, err := Normalize(state)
	if err != nil {
		return nil, err
	}
	if cmplx.Abs(state[0]) <= 1e-15 {
		return nil, fmt.Errorf("anchor gauge coordinate is too small")
	}
	gauge := cmplx.Rect(1, -cmplx.Phase(state[0]))
	for i := range state {
		state[i] *= gauge
	}
	if real(state[0]) < 0 {
		for i := range state {
			state[i] = -state[i]
		}
	}
	return state, nil
}

func paramsFromAnchorState(state State) (AnchorParams, error) {
	state, err := Normalize(state)
	if err != nil {
		return AnchorParams{}, err
	}
	var out AnchorParams
	for i := range state {
		out.Real[i] = real(state[i])
		out.Imag[i] = imag(state[i])
	}
	return out, nil
}

func canonicalizeAnchorParams(params AnchorParams) (AnchorParams, error) {
	state, err := anchorStateFromParams(params)
	if err != nil {
		return AnchorParams{}, err
	}
	return paramsFromAnchorState(state)
}

func makeTaskLearnedAnchorBank(params AnchorParams) (frameBank, error) {
	bank, err := makeLearnedSupportBank(
		up13LearnedPhases(),
		up14LearnedSupportMatrix(),
	)
	if err != nil {
		return frameBank{}, err
	}
	anchor, err := anchorStateFromParams(params)
	if err != nil {
		return frameBank{}, err
	}
	bank.anchor = anchor
	return bank, nil
}

func anchorPhaseConcentration(params AnchorParams) (float64, error) {
	anchor, err := anchorStateFromParams(params)
	if err != nil {
		return 0, err
	}
	var sum complex128
	count := 0
	for _, table := range allMemoryTables() {
		memory, err := encodeMemory(table)
		if err != nil {
			return 0, err
		}
		correlation, err := stateInner(anchor, memory)
		if err != nil {
			return 0, err
		}
		magnitude := cmplx.Abs(correlation)
		if magnitude <= 1e-15 {
			continue
		}
		sum += correlation / complex(magnitude, 0)
		count++
	}
	if count == 0 {
		return 0, fmt.Errorf("anchor has no measurable table correlations")
	}
	return cmplx.Abs(sum / complex(float64(count), 0)), nil
}

func anchorStateShift(a, b AnchorParams) (float64, error) {
	sa, err := anchorStateFromParams(a)
	if err != nil {
		return 0, err
	}
	sb, err := anchorStateFromParams(b)
	if err != nil {
		return 0, err
	}
	return L2Distance(sa, sb)
}

func buildAnchorCanonicalSamples(
	tables []memoryTable,
	entity int,
	params AnchorParams,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	if entity < 0 || entity >= 4 || trials < 1 {
		return nil, fmt.Errorf("invalid learned-anchor sample request")
	}
	bank, err := makeTaskLearnedAnchorBank(params)
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
				math.Mod(0.197*float64(seed+1), 2*math.Pi),
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

func anchorHeadsLoss(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	params AnchorParams,
	memoryNoise float64,
	trials int,
) (float64, error) {
	var total float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildAnchorCanonicalSamples(
			tables, entity, params,
			memoryNoise, trials, 0,
		)
		if err != nil {
			return 0, err
		}
		loss, _, err := evaluateHead(
			heads[entity], samples,
		)
		if err != nil {
			return 0, err
		}
		total += loss
	}
	return total / 4, nil
}

func fitAnchorHeads(
	tables []memoryTable,
	params AnchorParams,
	memoryNoise float64,
	trials int,
	steps int,
) ([4]linearSoftmaxHead, float64, error) {
	var heads [4]linearSoftmaxHead
	var average float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildAnchorCanonicalSamples(
			tables, entity, params,
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
		_, accuracy, err := evaluateHead(
			head, samples,
		)
		if err != nil {
			return heads, 0, err
		}
		average += accuracy
	}
	return heads, average / 4, nil
}

func learnAnchor(
	tables []memoryTable,
	memoryNoise float64,
) (
	initial AnchorParams,
	learned AnchorParams,
	initialLoss float64,
	finalLoss float64,
	trace []LearnedAnchorTrace,
	err error,
) {
	const (
		trials              = 2
		outerSteps          = 80
		headStepsPerOuter   = 10
		headLearningRate    = 1.0
		anchorLearningRate  = 0.18
		anchorEpsilon       = 1e-4
	)
	params := initialLearnedAnchorParams()
	initial = params

	var heads [4]linearSoftmaxHead
	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(4, 2)
	}

	for entity := 0; entity < 4; entity++ {
		samples, e := buildAnchorCanonicalSamples(
			tables, entity, params,
			memoryNoise, trials, 0,
		)
		if e != nil {
			err = e
			return
		}
		for step := 0; step < headStepsPerOuter; step++ {
			if e := updateSoftmaxHead(
				&heads[entity],
				samples,
				headLearningRate,
			); e != nil {
				err = e
				return
			}
		}
	}
	initialLoss, err = anchorHeadsLoss(
		heads, tables, params,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}

	for outer := 1; outer <= outerSteps; outer++ {
		for entity := 0; entity < 4; entity++ {
			samples, e := buildAnchorCanonicalSamples(
				tables, entity, params,
				memoryNoise, trials, 0,
			)
			if e != nil {
				err = e
				return
			}
			for step := 0; step < headStepsPerOuter; step++ {
				if e := updateSoftmaxHead(
					&heads[entity],
					samples,
					headLearningRate,
				); e != nil {
					err = e
					return
				}
			}
		}

		var gradReal, gradImag [16]float64
		for index := 0; index < 16; index++ {
			plus := params
			minus := params
			plus.Real[index] += anchorEpsilon
			minus.Real[index] -= anchorEpsilon
			plusLoss, e := anchorHeadsLoss(
				heads, tables, plus,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			minusLoss, e := anchorHeadsLoss(
				heads, tables, minus,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			gradReal[index] =
				(plusLoss - minusLoss) /
					(2 * anchorEpsilon)

			plus = params
			minus = params
			plus.Imag[index] += anchorEpsilon
			minus.Imag[index] -= anchorEpsilon
			plusLoss, e = anchorHeadsLoss(
				heads, tables, plus,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			minusLoss, e = anchorHeadsLoss(
				heads, tables, minus,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			gradImag[index] =
				(plusLoss - minusLoss) /
					(2 * anchorEpsilon)

			if !finite(gradReal[index]) ||
				!finite(gradImag[index]) {
				err = fmt.Errorf(
					"anchor gradient index=%d non-finite",
					index,
				)
				return
			}
		}

		for index := 0; index < 16; index++ {
			params.Real[index] -=
				anchorLearningRate * gradReal[index]
			params.Imag[index] -=
				anchorLearningRate * gradImag[index]
		}
		params, err = canonicalizeAnchorParams(params)
		if err != nil {
			return
		}

		if outer == 1 || outer == 10 ||
			outer == 20 || outer == 40 ||
			outer == 60 || outer == 80 {
			loss, e := anchorHeadsLoss(
				heads, tables, params,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			concentration, e :=
				anchorPhaseConcentration(params)
			if e != nil {
				err = e
				return
			}
			trace = append(trace, LearnedAnchorTrace{
				Outer:              outer,
				Loss:               loss,
				PhaseConcentration: concentration,
				Params:             params,
			})
		}
	}

	finalLoss, err = anchorHeadsLoss(
		heads, tables, params,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}
	learned = params
	return
}

func buildAnchorTransportSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	params AnchorParams,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	baseBank, err := makeTaskLearnedAnchorBank(params)
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
					math.Mod(0.197*float64(seed+1), 2*math.Pi),
				)
				forward, err := apply(
					state, block, depth,
				)
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

func trainAnchorTransportHeads(
	name string,
	params AnchorParams,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	steps int,
) ([4]linearSoftmaxHead, LearnedAnchorStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildAnchorTransportSamples(
			trainTables, trainDepths, entity,
			params, block, apply,
			memoryNoise, 4, 0,
		)
		if err != nil {
			return heads, LearnedAnchorStaticResult{}, err
		}
		heldSamples, err := buildAnchorTransportSamples(
			heldTables, heldDepths, entity,
			params, block, apply,
			memoryNoise, 2, 7000000,
		)
		if err != nil {
			return heads, LearnedAnchorStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, LearnedAnchorStaticResult{}, err
		}
		heads[entity] = head
		_, trainAccuracy, err := evaluateHead(
			head, trainSamples,
		)
		if err != nil {
			return heads, LearnedAnchorStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(
			head, heldSamples,
		)
		if err != nil {
			return heads, LearnedAnchorStaticResult{}, err
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

	return heads, LearnedAnchorStaticResult{
		Name:             name,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func decodeAnchorTable(
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

func runLearnedAnchorIntegration(
	params AnchorParams,
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (LearnedAnchorIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return LearnedAnchorIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return LearnedAnchorIntegration{}, err
	}
	baseBank, err := makeTaskLearnedAnchorBank(params)
	if err != nil {
		return LearnedAnchorIntegration{}, err
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
				return LearnedAnchorIntegration{}, err
			}
			seed := 15000000 +
				scenarioIndex*10000 +
				writeIndex*31
			state, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return LearnedAnchorIntegration{}, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.197*float64(seed+1), 2*math.Pi),
			)
			forward, err := applyStressUnitary(
				state, block, write.gap,
			)
			if err != nil {
				return LearnedAnchorIntegration{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return LearnedAnchorIntegration{}, err
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
				return LearnedAnchorIntegration{}, err
			}
			decoded, _, margin, err := decodeAnchorTable(
				forward, bank, heads,
			)
			if err != nil {
				return LearnedAnchorIntegration{}, err
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
				return LearnedAnchorIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return LearnedAnchorIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		seed := 15000000 + scenarioIndex*10000 + 9999
		state, err := perturbMemory(
			canonical, seed, memoryNoise,
		)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		state = rotateGlobalPhase(
			state,
			math.Mod(0.197*float64(seed+1), 2*math.Pi),
		)
		forward, err := applyStressUnitary(
			state, block, scenario.finalGap,
		)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return LearnedAnchorIntegration{}, err
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
			return LearnedAnchorIntegration{}, err
		}
		decoded, distributions, margin, err := decodeAnchorTable(
			forward, bank, heads,
		)
		if err != nil {
			return LearnedAnchorIntegration{}, err
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
			return LearnedAnchorIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return LearnedAnchorIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return LearnedAnchorIntegration{
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

// RunUP15 learns the global anchor from task loss.
//
// The rank-4 structure, UP-13 learned value phases, and UP-14 learned pilot
// support matrix remain fixed. The anchor is a full dense 16-D complex state.
// It begins from a deterministic irregular complex vector, not the uniform
// anchor used by UP-9 through UP-14.
//
// All 32 real degrees of freedom are updated by central-difference gradients
// of average classification cross-entropy. State normalization removes scale;
// after every update a deterministic gauge makes coordinate zero real-positive.
// No target anchor is supplied.
func RunUP15() (LearnedAnchorProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	initial := initialLearnedAnchorParams()
	_, initialCapacity, err := fitAnchorHeads(
		trainTables, initial,
		memoryNoise, 2, 1000,
	)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}
	initialConcentration, err :=
		anchorPhaseConcentration(initial)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}

	initialAnchor, learnedAnchor, initialLoss, finalLoss, trace, err :=
		learnAnchor(trainTables, memoryNoise)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}
	learnedConcentration, err :=
		anchorPhaseConcentration(learnedAnchor)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}
	shift, err := anchorStateShift(
		initialAnchor, learnedAnchor,
	)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}

	block := stressProgram()
	unitaryHeads, unitaryResult, err :=
		trainAnchorTransportHeads(
			"unitary_task_learned_anchor",
			learnedAnchor,
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
		return LearnedAnchorProbeResult{}, err
	}

	_, controlResult, err :=
		trainAnchorTransportHeads(
			"non_unitary_task_learned_anchor",
			learnedAnchor,
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
		return LearnedAnchorProbeResult{}, err
	}

	integration, err := runLearnedAnchorIntegration(
		learnedAnchor,
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return LearnedAnchorProbeResult{}, err
	}

	anchorPass :=
		initialCapacity < 0.70 &&
			unitaryResult.TrainAccuracy >= 0.99 &&
			finalLoss < initialLoss &&
			shift >= 0.20
	unseenPass := unitaryResult.HeldOutAccuracy >= 0.99
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95

	return LearnedAnchorProbeResult{
		Schema:                 LearnedAnchorSchema,
		Experiment:             "UP-15-task-learned-global-anchor",
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
		EntitySupportLearned:   false,
		AnchorLearned:          true,
		MemoryNoiseAmplitude:   memoryNoise,
		TrainTables:            len(trainTables),
		HeldOutTables:          len(heldTables),
		TrainDepths:            append([]int(nil), trainDepths...),
		HeldOutDepths:          append([]int(nil), heldDepths...),
		FixedPhases:            append([]float64(nil), up13LearnedPhases()...),
		FixedSupport:           up14LearnedSupportMatrix(),
		InitialAnchor:          initialAnchor,
		LearnedAnchor:          learnedAnchor,
		InitialLoss:            initialLoss,
		FinalLoss:              finalLoss,
		Trace:                  trace,
		Unitary:                unitaryResult,
		MatchedControl:         controlResult,
		Integration:            integration,
		Diagnosis: LearnedAnchorDiagnosis{
			AnchorLearningPass:        anchorPass,
			UnseenDepthPass:           unseenPass,
			MutableIntegrationPass:    mutablePass,
			InitialCapacityAccuracy:   initialCapacity,
			FinalTrainAccuracy:        unitaryResult.TrainAccuracy,
			FinalHeldOutAccuracy:      unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:    controlResult.HeldOutAccuracy,
			InitialPhaseConcentration: initialConcentration,
			LearnedPhaseConcentration: learnedConcentration,
			AnchorStateL2Shift:        shift,
		},
	}, nil
}
