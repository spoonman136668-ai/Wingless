package unitary

import (
	"fmt"
	"math"
)

const RealOrthogonalEquivalenceSchema = "wingless.real-orthogonal-equivalence.v1"
const realLatentDimension = fullLatentDimension * 2

type realLatentState []float64
type realLatentMatrix [][]float64

type realObservable struct {
	re [][]float64
	im [][]float64
}

type RealOrthogonalStaticResult struct {
	ComplexHeldOutAccuracy      float64 `json:"complex_heldout_accuracy"`
	RealHeldOutAccuracy         float64 `json:"real_heldout_accuracy"`
	DecisionDisagreements       int     `json:"decision_disagreements"`
	MaxComplexNormDrift         float64 `json:"max_complex_norm_drift"`
	MaxRealNormDrift            float64 `json:"max_real_norm_drift"`
	MaxStateEquivalenceError    float64 `json:"max_state_equivalence_error"`
	MaxFeatureEquivalenceError  float64 `json:"max_feature_equivalence_error"`
}

type RealOrthogonalIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxRealNormDrift        float64 `json:"max_real_norm_drift"`
}

type RealOrthogonalDiagnosis struct {
	BaselineReproductionPass  bool    `json:"baseline_reproduction_pass"`
	OrthogonalityPass          bool    `json:"orthogonality_pass"`
	StateEquivalencePass       bool    `json:"state_equivalence_pass"`
	FeatureEquivalencePass     bool    `json:"feature_equivalence_pass"`
	DecisionEquivalencePass    bool    `json:"decision_equivalence_pass"`
	MutableIntegrationPass     bool    `json:"mutable_integration_pass"`
	MaxOrthogonalityError      float64 `json:"max_orthogonality_error"`
	MaxStateEquivalenceError   float64 `json:"max_state_equivalence_error"`
	MaxFeatureEquivalenceError float64 `json:"max_feature_equivalence_error"`
	ComplexHeldOutAccuracy     float64 `json:"complex_heldout_accuracy"`
	RealHeldOutAccuracy        float64 `json:"real_heldout_accuracy"`
	DecisionDisagreements      int     `json:"decision_disagreements"`
}

type RealOrthogonalEquivalenceProbeResult struct {
	Schema                              string                    `json:"schema"`
	Experiment                          string                    `json:"experiment"`
	ComplexLatentDimension              int                       `json:"complex_latent_dimension"`
	RealLatentDimension                 int                       `json:"real_latent_dimension"`
	EqualRealScalarDegreesOfFreedom     bool                      `json:"equal_real_scalar_degrees_of_freedom"`
	ComplexTransportUnitary             bool                      `json:"complex_transport_unitary"`
	RealTransportOrthogonal             bool                      `json:"real_transport_orthogonal"`
	ExactRealificationComparator        bool                      `json:"exact_realification_comparator"`
	RuntimeRealArithmeticOnly           bool                      `json:"runtime_real_arithmetic_only"`
	IndependentRealEncoderLearned       bool                      `json:"independent_real_encoder_learned"`
	EncoderRealifiedFromComplexConstruction bool                  `json:"encoder_realified_from_complex_construction"`
	TransportRealifiedFromComplexConstruction bool                `json:"transport_realified_from_complex_construction"`
	ObservableBankFrozenFromUP28        bool                      `json:"observable_bank_frozen_from_up28"`
	DecoderFrozenAcrossRepresentations  bool                      `json:"decoder_frozen_across_representations"`
	RuntimeObservableCount              int                       `json:"runtime_observable_count"`
	RuntimePrototypeLookup              bool                      `json:"runtime_prototype_lookup"`
	ExplicitDepthProvided               bool                      `json:"explicit_depth_provided"`
	Static                              RealOrthogonalStaticResult `json:"static_equivalence"`
	Integration                         RealOrthogonalIntegration `json:"real_mutable_integration"`
	Diagnosis                           RealOrthogonalDiagnosis   `json:"diagnosis"`
}

func realifyState(state State) (realLatentState, error) {
	if len(state) != fullLatentDimension {
		return nil, fmt.Errorf("realify state dimension=%d want=%d", len(state), fullLatentDimension)
	}
	out := make(realLatentState, realLatentDimension)
	for i, value := range state {
		out[2*i] = real(value)
		out[2*i+1] = imag(value)
	}
	return out, nil
}

func complexifyRealState(state realLatentState) (State, error) {
	if len(state) != realLatentDimension {
		return nil, fmt.Errorf("complexify real state dimension mismatch")
	}
	out := make(State, fullLatentDimension)
	for i := 0; i < fullLatentDimension; i++ {
		out[i] = complex(state[2*i], state[2*i+1])
	}
	return out, nil
}

func realifyMatrix(matrix latentMatrix) (realLatentMatrix, error) {
	if len(matrix) != fullLatentDimension {
		return nil, fmt.Errorf("realify matrix dimension mismatch")
	}
	out := make(realLatentMatrix, realLatentDimension)
	for row := range out {
		out[row] = make([]float64, realLatentDimension)
	}
	for row := 0; row < fullLatentDimension; row++ {
		if len(matrix[row]) != fullLatentDimension {
			return nil, fmt.Errorf("realify matrix row dimension mismatch")
		}
		for column := 0; column < fullLatentDimension; column++ {
			value := matrix[row][column]
			a, b := real(value), imag(value)
			out[2*row][2*column] = a
			out[2*row][2*column+1] = -b
			out[2*row+1][2*column] = b
			out[2*row+1][2*column+1] = a
		}
	}
	return out, nil
}

func realMatrixVector(
	matrix realLatentMatrix,
	state realLatentState,
) (realLatentState, error) {
	if len(matrix) != realLatentDimension || len(state) != realLatentDimension {
		return nil, fmt.Errorf("real matrix-vector dimension mismatch")
	}
	out := make(realLatentState, realLatentDimension)
	for row := 0; row < realLatentDimension; row++ {
		if len(matrix[row]) != realLatentDimension {
			return nil, fmt.Errorf("real matrix row dimension mismatch")
		}
		var value float64
		for column := 0; column < realLatentDimension; column++ {
			value += matrix[row][column] * state[column]
		}
		if !finite(value) {
			return nil, fmt.Errorf("real matrix-vector result non-finite")
		}
		out[row] = value
	}
	return out, nil
}

func realNormSquared(state realLatentState) (float64, error) {
	if len(state) != realLatentDimension {
		return 0, fmt.Errorf("real norm dimension mismatch")
	}
	var sum float64
	for _, value := range state {
		sum += value * value
	}
	if !finite(sum) {
		return 0, fmt.Errorf("real norm non-finite")
	}
	return sum, nil
}

func maxRealOrthogonalityError(matrix realLatentMatrix) (float64, error) {
	if len(matrix) != realLatentDimension {
		return 0, fmt.Errorf("orthogonality matrix dimension mismatch")
	}
	maximum := 0.0
	for first := 0; first < realLatentDimension; first++ {
		for second := first; second < realLatentDimension; second++ {
			var dot float64
			for row := 0; row < realLatentDimension; row++ {
				if len(matrix[row]) != realLatentDimension {
					return 0, fmt.Errorf("orthogonality row dimension mismatch")
				}
				dot += matrix[row][first] * matrix[row][second]
			}
			want := 0.0
			if first == second {
				want = 1
			}
			delta := math.Abs(dot - want)
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum, nil
}

func realifyObservable(observable latentMatrix) (realObservable, error) {
	if len(observable) != fullLatentDimension {
		return realObservable{}, fmt.Errorf("real observable dimension mismatch")
	}
	out := realObservable{
		re: make([][]float64, fullLatentDimension),
		im: make([][]float64, fullLatentDimension),
	}
	for row := 0; row < fullLatentDimension; row++ {
		if len(observable[row]) != fullLatentDimension {
			return realObservable{}, fmt.Errorf("real observable row dimension mismatch")
		}
		out.re[row] = make([]float64, fullLatentDimension)
		out.im[row] = make([]float64, fullLatentDimension)
		for column := 0; column < fullLatentDimension; column++ {
			out.re[row][column] = real(observable[row][column])
			out.im[row][column] = imag(observable[row][column])
		}
	}
	return out, nil
}

func realifyObservableBank(
	observables []latentMatrix,
) ([]realObservable, error) {
	out := make([]realObservable, len(observables))
	for i, observable := range observables {
		converted, err := realifyObservable(observable)
		if err != nil {
			return nil, err
		}
		out[i] = converted
	}
	return out, nil
}

func realObservableExpectation(
	state realLatentState,
	observable realObservable,
) (float64, float64, error) {
	if len(state) != realLatentDimension ||
		len(observable.re) != fullLatentDimension ||
		len(observable.im) != fullLatentDimension {
		return 0, 0, fmt.Errorf("real observable expectation dimension mismatch")
	}
	var expectationReal, expectationImag float64
	for row := 0; row < fullLatentDimension; row++ {
		xr := state[2*row]
		xi := state[2*row+1]
		var yr, yi float64
		for column := 0; column < fullLatentDimension; column++ {
			if len(observable.re[row]) != fullLatentDimension ||
				len(observable.im[row]) != fullLatentDimension {
				return 0, 0, fmt.Errorf("real observable coefficient row mismatch")
			}
			a := observable.re[row][column]
			b := observable.im[row][column]
			cr := state[2*column]
			ci := state[2*column+1]
			yr += a*cr - b*ci
			yi += b*cr + a*ci
		}
		expectationReal += xr*yr + xi*yi
		expectationImag += xr*yi - xi*yr
	}
	if !finite(expectationReal) || !finite(expectationImag) {
		return 0, 0, fmt.Errorf("real observable expectation non-finite")
	}
	return expectationReal, expectationImag, nil
}

func realComparatorQuadraticFeatures(
	state realLatentState,
	observables []realObservable,
) ([]float64, error) {
	raw := make([]float64, 0, len(observables)*2)
	for _, observable := range observables {
		re, im, err := realObservableExpectation(state, observable)
		if err != nil {
			return nil, err
		}
		raw = append(raw, re, im)
	}
	expected := len(raw) + (len(raw)*(len(raw)+1))/2
	out := make([]float64, 0, expected)
	out = append(out, raw...)
	for first := 0; first < len(raw); first++ {
		for second := first; second < len(raw); second++ {
			out = append(out, raw[first]*raw[second])
		}
	}
	return out, nil
}

func up28FrozenSelectedIndices() []int {
	return []int{
		117,17,31,79,3,121,116,48,36,76,86,25,92,7,112,93,
		58,103,18,108,32,91,90,105,71,119,80,26,120,104,78,27,
		38,0,34,99,87,40,20,57,68,98,11,23,85,24,52,8,
		41,118,69,100,33,4,77,124,64,84,74,49,96,65,13,62,
	}
}

func frozenUP28ObservableBank(
	unitaryStep latentMatrix,
) ([]latentMatrix, error) {
	candidates, _, err := discoverCommutingObservables(
		unitaryStep, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return nil, err
	}
	indices := up28FrozenSelectedIndices()
	out := make([]latentMatrix, len(indices))
	for i, index := range indices {
		if index < 0 || index >= len(candidates) {
			return nil, fmt.Errorf("frozen UP28 observable index out of range=%d", index)
		}
		out[i] = candidates[index]
	}
	return out, nil
}

func realifyDepthOperators(
	complexOps map[int]latentMatrix,
) (map[int]realLatentMatrix, error) {
	out := make(map[int]realLatentMatrix, len(complexOps))
	for depth, op := range complexOps {
		realOp, err := realifyMatrix(op)
		if err != nil {
			return nil, err
		}
		out[depth] = realOp
	}
	return out, nil
}

func buildRealComparatorRows(
	tables []memoryTable,
	depths []int,
	complexOps map[int]latentMatrix,
	realOps map[int]realLatentMatrix,
	mixer latentMatrix,
	complexObservables []latentMatrix,
	realObservables []realObservable,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]breadthRow, float64, float64, float64, error) {
	rows := make([]breadthRow, 0, len(tables)*len(depths)*trials)
	var maxRealNormDrift, maxStateError, maxFeatureError float64

	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, 0, 0, err
		}
		for depthIndex, depth := range depths {
			complexOp, ok := complexOps[depth]
			if !ok {
				return nil, 0, 0, 0, fmt.Errorf("missing complex comparator depth=%d", depth)
			}
			realOp, ok := realOps[depth]
			if !ok {
				return nil, 0, 0, 0, fmt.Errorf("missing real comparator depth=%d", depth)
			}
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.271*float64(seed+1), 2*math.Pi),
				)
				initialComplex, err := fullLatentEncode(memory, mixer)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				evolvedComplex, err := latentMatrixVector(complexOp, initialComplex)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				initialReal, err := realifyState(initialComplex)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				evolvedReal, err := realMatrixVector(realOp, initialReal)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				realNorm, err := realNormSquared(evolvedReal)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				if drift := math.Abs(realNorm - 1); drift > maxRealNormDrift {
					maxRealNormDrift = drift
				}
				converted, err := complexifyRealState(evolvedReal)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				stateError, err := L2Distance(converted, evolvedComplex)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				if stateError > maxStateError {
					maxStateError = stateError
				}

				complexFeatures, err := breadthQuadraticFeatures(
					evolvedComplex, complexObservables,
				)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				realFeatures, err := realComparatorQuadraticFeatures(
					evolvedReal, realObservables,
				)
				if err != nil {
					return nil, 0, 0, 0, err
				}
				if len(realFeatures) != len(complexFeatures) {
					return nil, 0, 0, 0, fmt.Errorf("real/complex feature dimension mismatch")
				}
				for i := range complexFeatures {
					if delta := math.Abs(realFeatures[i] - complexFeatures[i]); delta > maxFeatureError {
						maxFeatureError = delta
					}
				}
				rows = append(rows, breadthRow{
					features: realFeatures,
					table: table,
				})
			}
		}
	}
	return rows, maxRealNormDrift, maxStateError, maxFeatureError, nil
}

func countRealDecisionDisagreements(
	complexRows, realRows []breadthRow,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
) (int, error) {
	if len(complexRows) != len(realRows) {
		return 0, fmt.Errorf("decision comparison row mismatch")
	}
	var disagreements int
	for row := range complexRows {
		for entity := 0; entity < 4; entity++ {
			complexPhase, err := regressors[entity].predict(complexRows[row].features)
			if err != nil {
				return 0, err
			}
			realPhase, err := regressors[entity].predict(realRows[row].features)
			if err != nil {
				return 0, err
			}
			complexProb, err := classifiers[entity].probabilities(complexPhase)
			if err != nil {
				return 0, err
			}
			realProb, err := classifiers[entity].probabilities(realPhase)
			if err != nil {
				return 0, err
			}
			complexClass, _, err := classAndMargin(complexProb)
			if err != nil {
				return 0, err
			}
			realClass, _, err := classAndMargin(realProb)
			if err != nil {
				return 0, err
			}
			if complexClass != realClass {
				disagreements++
			}
		}
	}
	return disagreements, nil
}

func decodeRealComparatorTable(
	state realLatentState,
	observables []realObservable,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	features, err := realComparatorQuadraticFeatures(state, observables)
	if err != nil {
		return table, distributions, 0, err
	}
	for entity := 0; entity < 4; entity++ {
		phase, err := regressors[entity].predict(features)
		if err != nil {
			return table, distributions, 0, err
		}
		probabilities, err := classifiers[entity].probabilities(phase)
		if err != nil {
			return table, distributions, 0, err
		}
		value, margin, err := classAndMargin(probabilities)
		if err != nil {
			return table, distributions, 0, err
		}
		table[entity] = value
		distributions[entity] = probabilities
		if margin < minMargin {
			minMargin = margin
		}
	}
	return table, distributions, minMargin, nil
}

func runRealComparatorIntegration(
	mixer latentMatrix,
	realObservables []realObservable,
	realOps map[int]realLatentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (RealOrthogonalIntegration, error) {
	const scenarios, writes = 48, 16
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return RealOrthogonalIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return RealOrthogonalIntegration{}, err
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
				return RealOrthogonalIntegration{}, err
			}
			seed := 34000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			complexState, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
			state, err := realifyState(complexState)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
			op, ok := realOps[write.gap]
			if !ok {
				return RealOrthogonalIntegration{}, fmt.Errorf("missing real mutable depth=%d", write.gap)
			}
			state, err = realMatrixVector(op, state)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
			norm2, err := realNormSquared(state)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
			if drift := math.Abs(norm2 - 1); drift > maxNormDrift {
				maxNormDrift = drift
			}
			decoded, _, margin, err := decodeRealComparatorTable(
				state, realObservables, regressors, classifiers,
			)
			if err != nil {
				return RealOrthogonalIntegration{}, err
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
				return RealOrthogonalIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return RealOrthogonalIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		seed := 34000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1), 2*math.Pi),
		)
		complexState, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		state, err := realifyState(complexState)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		op, ok := realOps[scenario.finalGap]
		if !ok {
			return RealOrthogonalIntegration{}, fmt.Errorf("missing real final depth=%d", scenario.finalGap)
		}
		state, err = realMatrixVector(op, state)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		norm2, err := realNormSquared(state)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		if drift := math.Abs(norm2 - 1); drift > maxNormDrift {
			maxNormDrift = drift
		}
		decoded, distributions, margin, err := decodeRealComparatorTable(
			state, realObservables, regressors, classifiers,
		)
		if err != nil {
			return RealOrthogonalIntegration{}, err
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
			return RealOrthogonalIntegration{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return RealOrthogonalIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return RealOrthogonalIntegration{
		Scenarios: scenarios,
		WritesPerScenario: writes,
		CommitDecodeAccuracy: float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect) / float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(scenarios),
		MinValueMargin: minValueMargin,
		MinRelationMargin: minRelationMargin,
		MaxRealNormDrift: maxNormDrift,
	}, nil
}

func RunUP29() (RealOrthogonalEquivalenceProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	unitaryStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	complexObservables, err := frozenUP28ObservableBank(unitaryStep)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	realObservables, err := realifyObservableBank(complexObservables)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	realStep, err := realifyMatrix(unitaryStep)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	orthogonalityError, err := maxRealOrthogonalityError(realStep)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	complexOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	realOps, err := realifyDepthOperators(complexOps)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	trainRows, trainComplexDrift, err := buildBreadthRows(
		trainTables, trainDepths, complexOps,
		mixer, complexObservables, memoryNoise, 2, 0,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	regressors, classifiers, base, err := trainBreadthModel(
		trainRows, len(complexObservables),
		"up28_frozen_complex_baseline", "complex_baseline",
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	base.MaxNormDrift = trainComplexDrift

	complexHeldRows, complexHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, complexOps,
		mixer, complexObservables, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	complexResult, err := evaluateBreadthModel(
		base, regressors, classifiers,
		complexHeldRows, complexHeldDrift,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	realHeldRows, realHeldDrift, stateError, featureError, err :=
		buildRealComparatorRows(
			heldTables, heldDepths,
			complexOps, realOps,
			mixer, complexObservables, realObservables,
			memoryNoise, 2, 7000000,
		)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}
	realBase := base
	realBase.Name = "real_orthogonal_equivalent"
	realBase.Path = "real_192d_exact_realification"
	realResult, err := evaluateBreadthModel(
		realBase, regressors, classifiers,
		realHeldRows, realHeldDrift,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	disagreements, err := countRealDecisionDisagreements(
		complexHeldRows, realHeldRows,
		regressors, classifiers,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	integration, err := runRealComparatorIntegration(
		mixer, realObservables, realOps,
		regressors, classifiers,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil {
		return RealOrthogonalEquivalenceProbeResult{}, err
	}

	baselinePass := complexResult.HeldOutAccuracy == 1
	orthogonalityPass := orthogonalityError <= 1e-10
	statePass := stateError <= 1e-9
	featurePass := featureError <= 1e-8
	decisionPass :=
		realResult.HeldOutAccuracy == complexResult.HeldOutAccuracy &&
			disagreements == 0
	mutablePass :=
		integration.CommitDecodeAccuracy == 1 &&
			integration.ExactFinalTableAccuracy == 1 &&
			integration.RelationalQueryAccuracy == 1 &&
			integration.MaxRealNormDrift <= 1e-10

	return RealOrthogonalEquivalenceProbeResult{
		Schema: RealOrthogonalEquivalenceSchema,
		Experiment: "UP-29-exact-real-orthogonal-equivalence",
		ComplexLatentDimension: fullLatentDimension,
		RealLatentDimension: realLatentDimension,
		EqualRealScalarDegreesOfFreedom: true,
		ComplexTransportUnitary: true,
		RealTransportOrthogonal: true,
		ExactRealificationComparator: true,
		RuntimeRealArithmeticOnly: true,
		IndependentRealEncoderLearned: false,
		EncoderRealifiedFromComplexConstruction: true,
		TransportRealifiedFromComplexConstruction: true,
		ObservableBankFrozenFromUP28: true,
		DecoderFrozenAcrossRepresentations: true,
		RuntimeObservableCount: len(complexObservables),
		RuntimePrototypeLookup: false,
		ExplicitDepthProvided: false,
		Static: RealOrthogonalStaticResult{
			ComplexHeldOutAccuracy: complexResult.HeldOutAccuracy,
			RealHeldOutAccuracy: realResult.HeldOutAccuracy,
			DecisionDisagreements: disagreements,
			MaxComplexNormDrift: complexHeldDrift,
			MaxRealNormDrift: realHeldDrift,
			MaxStateEquivalenceError: stateError,
			MaxFeatureEquivalenceError: featureError,
		},
		Integration: integration,
		Diagnosis: RealOrthogonalDiagnosis{
			BaselineReproductionPass: baselinePass,
			OrthogonalityPass: orthogonalityPass,
			StateEquivalencePass: statePass,
			FeatureEquivalencePass: featurePass,
			DecisionEquivalencePass: decisionPass,
			MutableIntegrationPass: mutablePass,
			MaxOrthogonalityError: orthogonalityError,
			MaxStateEquivalenceError: stateError,
			MaxFeatureEquivalenceError: featureError,
			ComplexHeldOutAccuracy: complexResult.HeldOutAccuracy,
			RealHeldOutAccuracy: realResult.HeldOutAccuracy,
			DecisionDisagreements: disagreements,
		},
	}, nil
}
