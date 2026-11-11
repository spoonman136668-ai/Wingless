package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const CommutingWeylSchema = "wingless.unitary-commuting-weyl.v1"

const (
	weylMultiplicityDimension = compositeChannels
	weylObservableCount       = weylMultiplicityDimension * weylMultiplicityDimension
	weylRawFeatureDim         = weylObservableCount * 2
	weylQuadraticFeatureDim   = weylRawFeatureDim + (weylRawFeatureDim*(weylRawFeatureDim+1))/2
	weylRidgeLambda           = 1e-6
)

type weylRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type WeylStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MeanPhaseCosine  float64   `json:"mean_phase_cosine"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type WeylIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type WeylDiagnosis struct {
	CommutatorPass          bool    `json:"commutator_pass"`
	FeatureInvariancePass   bool    `json:"feature_invariance_pass"`
	PhaseCodeLearningPass   bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass         bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass  bool    `json:"mutable_integration_pass"`
	MaxCommutatorEntryError float64 `json:"max_commutator_entry_error"`
	MaxFeatureDrift         float64 `json:"max_feature_drift"`
	TrainAccuracy           float64 `json:"train_accuracy"`
	HeldOutAccuracy         float64 `json:"held_out_accuracy"`
	MatchedControlAccuracy  float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine     float64 `json:"mean_held_phase_cosine"`
}

type CommutingWeylProbeResult struct {
	Schema                          string            `json:"schema"`
	Experiment                      string            `json:"experiment"`
	LatentDimension                 int               `json:"latent_dimension"`
	RuntimeStateObjects             int               `json:"runtime_state_objects"`
	VisibleChannelBlocks            bool              `json:"visible_channel_blocks"`
	FullCoordinateMixing            bool              `json:"full_coordinate_mixing"`
	TransportMultiplicityDimension  int               `json:"transport_multiplicity_dimension"`
	CommutingObservableAlgebra      bool              `json:"commuting_observable_algebra"`
	ObservableAlgebraHandConstructed bool             `json:"observable_algebra_hand_constructed"`
	WeylObservableCount             int               `json:"weyl_observable_count"`
	KnownFullMixerExposedToLearner  bool              `json:"known_full_mixer_exposed_to_learner"`
	RuntimeUnmixApplied             bool              `json:"runtime_unmix_applied"`
	LearnedDemixerUsed              bool              `json:"learned_demixer_used"`
	PhaseAlphabetSupervision        bool              `json:"phase_alphabet_supervision"`
	PhaseCodeDimension              int               `json:"phase_code_dimension"`
	RawFeatureDimension             int               `json:"raw_feature_dimension"`
	QuadraticFeatureDimension       int               `json:"quadratic_feature_dimension"`
	RuntimePrototypeLookup          bool              `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool              `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided           bool              `json:"explicit_depth_provided"`
	GlobalPhaseNuisance             bool              `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude            float64           `json:"memory_noise_amplitude"`
	TrainingDepths                  []int             `json:"training_depths"`
	HeldOutDepths                   []int             `json:"held_out_depths"`
	TrainTables                     int               `json:"train_tables"`
	HeldOutTables                   int               `json:"held_out_tables"`
	Unitary                         WeylStaticResult   `json:"unitary"`
	MatchedControl                  WeylStaticResult   `json:"matched_control"`
	Integration                     WeylIntegration    `json:"unitary_mutable_integration"`
	Diagnosis                       WeylDiagnosis      `json:"diagnosis"`
}

type weylRow struct {
	features []float64
	table    memoryTable
}

func compositeWeylOperator(shift, phaseIndex int) (latentMatrix, error) {
	if shift < 0 || shift >= weylMultiplicityDimension ||
		phaseIndex < 0 || phaseIndex >= weylMultiplicityDimension {
		return nil, fmt.Errorf("weyl index out of range")
	}
	out := make(latentMatrix, fullLatentDimension)
	for row := range out {
		out[row] = make([]complex128, fullLatentDimension)
	}
	for channel := 0; channel < weylMultiplicityDimension; channel++ {
		target := (channel + shift) % weylMultiplicityDimension
		phase := cmplx.Rect(
			1,
			2*math.Pi*float64(phaseIndex*channel)/
				float64(weylMultiplicityDimension),
		)
		for coordinate := 0; coordinate < compositeChannelDim; coordinate++ {
			row := target*compositeChannelDim + coordinate
			column := channel*compositeChannelDim + coordinate
			out[row][column] = phase
		}
	}
	return out, nil
}

func fullLatentWeylObservables(mixer latentMatrix) ([]latentMatrix, error) {
	adjoint, err := latentAdjoint(mixer)
	if err != nil {
		return nil, err
	}
	out := make([]latentMatrix, 0, weylObservableCount)
	for shift := 0; shift < weylMultiplicityDimension; shift++ {
		for phaseIndex := 0; phaseIndex < weylMultiplicityDimension; phaseIndex++ {
			raw, err := compositeWeylOperator(shift, phaseIndex)
			if err != nil {
				return nil, err
			}
			left, err := latentMatrixMultiply(mixer, raw)
			if err != nil {
				return nil, err
			}
			full, err := latentMatrixMultiply(left, adjoint)
			if err != nil {
				return nil, err
			}
			out = append(out, full)
		}
	}
	return out, nil
}

func weylRawFeatures(state State, observables []latentMatrix) ([]float64, error) {
	if len(state) != fullLatentDimension {
		return nil, fmt.Errorf("weyl state dimension mismatch")
	}
	if len(observables) != weylObservableCount {
		return nil, fmt.Errorf("weyl observable count=%d want=%d", len(observables), weylObservableCount)
	}
	out := make([]float64, 0, weylRawFeatureDim)
	for _, observable := range observables {
		applied, err := latentMatrixVector(observable, state)
		if err != nil {
			return nil, err
		}
		value, err := stateInner(state, applied)
		if err != nil {
			return nil, err
		}
		out = append(out, real(value), imag(value))
	}
	for _, value := range out {
		if !finite(value) {
			return nil, fmt.Errorf("weyl raw feature is non-finite")
		}
	}
	return out, nil
}

func weylQuadraticFeatures(state State, observables []latentMatrix) ([]float64, error) {
	raw, err := weylRawFeatures(state, observables)
	if err != nil {
		return nil, err
	}
	out := make([]float64, 0, weylQuadraticFeatureDim)
	out = append(out, raw...)
	for first := 0; first < len(raw); first++ {
		for second := first; second < len(raw); second++ {
			value := raw[first] * raw[second]
			if !finite(value) {
				return nil, fmt.Errorf("weyl quadratic feature is non-finite")
			}
			out = append(out, value)
		}
	}
	if len(out) != weylQuadraticFeatureDim {
		return nil, fmt.Errorf("weyl quadratic feature dimension=%d want=%d", len(out), weylQuadraticFeatureDim)
	}
	return out, nil
}

func maxWeylCommutatorEntry(observables []latentMatrix, step latentMatrix) (float64, error) {
	maximum := 0.0
	for _, observable := range observables {
		left, err := latentMatrixMultiply(observable, step)
		if err != nil {
			return 0, err
		}
		right, err := latentMatrixMultiply(step, observable)
		if err != nil {
			return 0, err
		}
		for row := 0; row < fullLatentDimension; row++ {
			for column := 0; column < fullLatentDimension; column++ {
				delta := cmplx.Abs(left[row][column] - right[row][column])
				if delta > maximum {
					maximum = delta
				}
			}
		}
	}
	return maximum, nil
}

func weylFeatureDrift(
	mixer latentMatrix,
	observables []latentMatrix,
	depthOperator latentMatrix,
) (float64, error) {
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 242424, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.613)
	state, err := fullLatentEncode(memory, mixer)
	if err != nil {
		return 0, err
	}
	before, err := weylQuadraticFeatures(state, observables)
	if err != nil {
		return 0, err
	}
	evolved, err := latentMatrixVector(depthOperator, state)
	if err != nil {
		return 0, err
	}
	after, err := weylQuadraticFeatures(evolved, observables)
	if err != nil {
		return 0, err
	}
	maximum := 0.0
	for i := range before {
		delta := math.Abs(before[i] - after[i])
		if delta > maximum {
			maximum = delta
		}
	}
	return maximum, nil
}

func buildWeylRows(
	tables []memoryTable,
	depths []int,
	operators map[int]latentMatrix,
	mixer latentMatrix,
	observables []latentMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]weylRow, float64, error) {
	rows := make([]weylRow, 0, len(tables)*len(depths)*trials)
	var maxNormDrift float64
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			operator, ok := operators[depth]
			if !ok {
				return nil, 0, fmt.Errorf("missing weyl operator depth=%d", depth)
			}
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(memory, math.Mod(0.263*float64(seed+1), 2*math.Pi))
				state, err := fullLatentEncode(memory, mixer)
				if err != nil {
					return nil, 0, err
				}
				state, err = latentMatrixVector(operator, state)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(state)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}
				features, err := weylQuadraticFeatures(state, observables)
				if err != nil {
					return nil, 0, err
				}
				rows = append(rows, weylRow{features: features, table: table})
			}
		}
	}
	return rows, maxNormDrift, nil
}

func trainWeylRegressors(rows []weylRow, lambda float64) ([4]weylRegressor, error) {
	var regressors [4]weylRegressor
	if len(rows) == 0 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid weyl ridge configuration")
	}
	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)
	for row := 0; row < n; row++ {
		if len(rows[row].features) != weylQuadraticFeatureDim {
			return regressors, fmt.Errorf("weyl feature dimension mismatch")
		}
		gram[row] = make([]float64, n)
		right[row] = make([]float64, outputs)
		for entity := 0; entity < 4; entity++ {
			target, err := learnedPhaseTarget(rows[row].table[entity])
			if err != nil {
				return regressors, err
			}
			right[row][entity*2] = target[0]
			right[row][entity*2+1] = target[1]
		}
	}
	for first := 0; first < n; first++ {
		value := 1.0 + lambda
		for feature := 0; feature < weylQuadraticFeatureDim; feature++ {
			value += rows[first].features[feature] * rows[first].features[feature]
		}
		gram[first][first] = value
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < weylQuadraticFeatureDim; feature++ {
				value += rows[first].features[feature] * rows[second].features[feature]
			}
			if !finite(value) {
				return regressors, fmt.Errorf("weyl gram value is non-finite")
			}
			gram[first][second] = value
			gram[second][first] = value
		}
	}
	alpha, err := solveDenseMultiple(gram, right)
	if err != nil {
		return regressors, err
	}
	for entity := 0; entity < 4; entity++ {
		for output := 0; output < 2; output++ {
			regressors[entity].weights[output] = make([]float64, weylQuadraticFeatureDim)
		}
	}
	for row := 0; row < n; row++ {
		for entity := 0; entity < 4; entity++ {
			for output := 0; output < 2; output++ {
				coefficient := alpha[row][entity*2+output]
				regressors[entity].bias[output] += coefficient
				for feature, value := range rows[row].features {
					regressors[entity].weights[output][feature] += coefficient * value
				}
			}
		}
	}
	return regressors, nil
}

func (reg weylRegressor) predict(features []float64) ([]float64, error) {
	if len(features) != weylQuadraticFeatureDim {
		return nil, fmt.Errorf("weyl prediction feature mismatch")
	}
	out := []float64{reg.bias[0], reg.bias[1]}
	for output := 0; output < 2; output++ {
		if len(reg.weights[output]) != weylQuadraticFeatureDim {
			return nil, fmt.Errorf("weyl regressor weight mismatch")
		}
		for i, value := range features {
			out[output] += reg.weights[output][i] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("weyl prediction is non-finite")
		}
	}
	return out, nil
}

func trainWeylModel(rows []weylRow) (
	[4]weylRegressor,
	[4]linearSoftmaxHead,
	WeylStaticResult,
	error,
) {
	var classifiers [4]linearSoftmaxHead
	regressors, err := trainWeylRegressors(rows, weylRidgeLambda)
	if err != nil {
		return regressors, classifiers, WeylStaticResult{}, err
	}
	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return regressors, classifiers, WeylStaticResult{}, err
			}
			samples = append(samples, headSample{features: phase, target: row.table[entity]})
		}
		head, _, err := trainLinearSoftmax(samples, 4, 2, 1200, 1.0)
		if err != nil {
			return regressors, classifiers, WeylStaticResult{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, WeylStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}
	average := 0.0
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4
	return regressors, classifiers, WeylStaticResult{
		Name: "unitary_full_latent_commuting_weyl_phase_code",
		Path: "unitary_commuting_weyl_algebra",
		TrainAccuracy: average,
		PerEntityTrain: perTrain,
	}, nil
}

func evaluateWeylModel(
	base WeylStaticResult,
	regressors [4]weylRegressor,
	classifiers [4]linearSoftmaxHead,
	rows []weylRow,
	maxNormDrift float64,
) (WeylStaticResult, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var cosineTotal float64
	var cosineCount int
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return WeylStaticResult{}, err
			}
			target, err := learnedPhaseTarget(row.table[entity])
			if err != nil {
				return WeylStaticResult{}, err
			}
			norm := math.Hypot(phase[0], phase[1])
			if norm > 1e-15 {
				cosine := (phase[0]*target[0] + phase[1]*target[1]) / norm
				if cosine > 1 { cosine = 1 }
				if cosine < -1 { cosine = -1 }
				cosineTotal += cosine
				cosineCount++
			}
			samples = append(samples, headSample{features: phase, target: row.table[entity]})
		}
		_, accuracy, err := evaluateHead(classifiers[entity], samples)
		if err != nil {
			return WeylStaticResult{}, err
		}
		result.PerEntityHeldOut[entity] = accuracy
	}
	for _, accuracy := range result.PerEntityHeldOut {
		result.HeldOutAccuracy += accuracy
	}
	result.HeldOutAccuracy /= 4
	result.MaxNormDrift = maxNormDrift
	if cosineCount > 0 {
		result.MeanPhaseCosine = cosineTotal / float64(cosineCount)
	}
	return result, nil
}

func decodeWeylTable(
	state State,
	observables []latentMatrix,
	regressors [4]weylRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	features, err := weylQuadraticFeatures(state, observables)
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

func runWeylIntegration(
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]weylRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (WeylIntegration, error) {
	const scenarios, writes = 48, 16
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil { return WeylIntegration{}, err }
	relationHead, _, err := trainLinearSoftmax(relationSamples, 4, 16, 600, 1.0)
	if err != nil { return WeylIntegration{}, err }

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
			if err != nil { return WeylIntegration{}, err }
			seed := 31000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil { return WeylIntegration{}, err }
			memory = rotateGlobalPhase(memory, math.Mod(0.263*float64(seed+1), 2*math.Pi))
			state, err := fullLatentEncode(memory, mixer)
			if err != nil { return WeylIntegration{}, err }
			operator, ok := unitaryOps[write.gap]
			if !ok { return WeylIntegration{}, fmt.Errorf("missing mutable weyl depth=%d", write.gap) }
			state, err = latentMatrixVector(operator, state)
			if err != nil { return WeylIntegration{}, err }
			norm2, err := NormSquared(state)
			if err != nil { return WeylIntegration{}, err }
			drift := math.Abs(norm2-1)
			if drift > maxNormDrift { maxNormDrift = drift }

			decoded, _, margin, err := decodeWeylTable(state, observables, regressors, classifiers)
			if err != nil { return WeylIntegration{}, err }
			if margin < minValueMargin { minValueMargin = margin }
			commitTotal++
			if decoded == trueTable { commitCorrect++ }
			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil { return WeylIntegration{}, err }
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil { return WeylIntegration{}, err }
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil { return WeylIntegration{}, err }
		seed := 31000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil { return WeylIntegration{}, err }
		memory = rotateGlobalPhase(memory, math.Mod(0.263*float64(seed+1), 2*math.Pi))
		state, err := fullLatentEncode(memory, mixer)
		if err != nil { return WeylIntegration{}, err }
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok { return WeylIntegration{}, fmt.Errorf("missing final weyl depth=%d", scenario.finalGap) }
		state, err = latentMatrixVector(operator, state)
		if err != nil { return WeylIntegration{}, err }
		norm2, err := NormSquared(state)
		if err != nil { return WeylIntegration{}, err }
		drift := math.Abs(norm2-1)
		if drift > maxNormDrift { maxNormDrift = drift }

		decoded, distributions, margin, err := decodeWeylTable(state, observables, regressors, classifiers)
		if err != nil { return WeylIntegration{}, err }
		if margin < minValueMargin { minValueMargin = margin }
		if decoded == trueTable { finalCorrect++ }

		relationInput, err := relationFeatures(distributions[scenario.queryA], distributions[scenario.queryB])
		if err != nil { return WeylIntegration{}, err }
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil { return WeylIntegration{}, err }
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil { return WeylIntegration{}, err }
		if relationMargin < minRelationMargin { minRelationMargin = relationMargin }
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil { return WeylIntegration{}, err }
		if gotRelation == wantRelation { relationCorrect++ }
	}

	return WeylIntegration{
		Scenarios: scenarios,
		WritesPerScenario: writes,
		CommitDecodeAccuracy: float64(commitCorrect)/float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect)/float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect)/float64(scenarios),
		MinValueMargin: minValueMargin,
		MinRelationMargin: minRelationMargin,
		MaxNormDrift: maxNormDrift,
	}, nil
}

func RunUP24() (CommutingWeylProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	observables, err := fullLatentWeylObservables(mixer)
	if err != nil { return CommutingWeylProbeResult{}, err }

	unitaryStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { return CommutingWeylProbeResult{}, err }
	controlStep, err := conjugatedLatentStep(mixer, applyStressNonUnitary)
	if err != nil { return CommutingWeylProbeResult{}, err }

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil { return CommutingWeylProbeResult{}, err }
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil { return CommutingWeylProbeResult{}, err }

	commutatorError, err := maxWeylCommutatorEntry(observables, unitaryStep)
	if err != nil { return CommutingWeylProbeResult{}, err }
	featureDrift, err := weylFeatureDrift(mixer, observables, unitaryOps[1024])
	if err != nil { return CommutingWeylProbeResult{}, err }

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	trainRows, trainNormDrift, err := buildWeylRows(
		trainTables, trainDepths, unitaryOps,
		mixer, observables, memoryNoise, 2, 0,
	)
	if err != nil { return CommutingWeylProbeResult{}, err }

	regressors, classifiers, base, err := trainWeylModel(trainRows)
	if err != nil { return CommutingWeylProbeResult{}, err }
	base.MaxNormDrift = trainNormDrift

	heldRows, heldDrift, err := buildWeylRows(
		heldTables, heldDepths, unitaryOps,
		mixer, observables, memoryNoise, 2, 7000000,
	)
	if err != nil { return CommutingWeylProbeResult{}, err }
	unitaryResult, err := evaluateWeylModel(base, regressors, classifiers, heldRows, heldDrift)
	if err != nil { return CommutingWeylProbeResult{}, err }

	controlRows, controlDrift, err := buildWeylRows(
		heldTables, heldDepths, controlOps,
		mixer, observables, memoryNoise, 2, 7000000,
	)
	if err != nil { return CommutingWeylProbeResult{}, err }
	controlResult, err := evaluateWeylModel(
		WeylStaticResult{
			Name: "non_unitary_full_latent_commuting_weyl_phase_code",
			Path: "non_unitary_same_model",
			TrainAccuracy: base.TrainAccuracy,
			PerEntityTrain: append([]float64(nil), base.PerEntityTrain...),
		},
		regressors, classifiers, controlRows, controlDrift,
	)
	if err != nil { return CommutingWeylProbeResult{}, err }

	integration, err := runWeylIntegration(
		mixer, observables, unitaryOps,
		regressors, classifiers,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil { return CommutingWeylProbeResult{}, err }

	commutatorPass := commutatorError <= 1e-10
	invariancePass := featureDrift <= 1e-9
	learningPass := unitaryResult.TrainAccuracy >= 0.99
	unseenPass := unitaryResult.HeldOutAccuracy >= 0.99 && unitaryResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
		integration.ExactFinalTableAccuracy >= 0.95 &&
		integration.RelationalQueryAccuracy >= 0.95 &&
		integration.MaxNormDrift <= 1e-10

	return CommutingWeylProbeResult{
		Schema: CommutingWeylSchema,
		Experiment: "UP-24-full-latent-commuting-weyl-algebra",
		LatentDimension: fullLatentDimension,
		RuntimeStateObjects: 1,
		VisibleChannelBlocks: false,
		FullCoordinateMixing: true,
		TransportMultiplicityDimension: weylMultiplicityDimension,
		CommutingObservableAlgebra: true,
		ObservableAlgebraHandConstructed: true,
		WeylObservableCount: weylObservableCount,
		KnownFullMixerExposedToLearner: false,
		RuntimeUnmixApplied: false,
		LearnedDemixerUsed: false,
		PhaseAlphabetSupervision: true,
		PhaseCodeDimension: 2,
		RawFeatureDimension: weylRawFeatureDim,
		QuadraticFeatureDimension: weylQuadraticFeatureDim,
		RuntimePrototypeLookup: false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided: false,
		GlobalPhaseNuisance: true,
		MemoryNoiseAmplitude: memoryNoise,
		TrainingDepths: trainDepths,
		HeldOutDepths: heldDepths,
		TrainTables: len(trainTables),
		HeldOutTables: len(heldTables),
		Unitary: unitaryResult,
		MatchedControl: controlResult,
		Integration: integration,
		Diagnosis: WeylDiagnosis{
			CommutatorPass: commutatorPass,
			FeatureInvariancePass: invariancePass,
			PhaseCodeLearningPass: learningPass,
			UnseenDepthPass: unseenPass,
			MutableIntegrationPass: mutablePass,
			MaxCommutatorEntryError: commutatorError,
			MaxFeatureDrift: featureDrift,
			TrainAccuracy: unitaryResult.TrainAccuracy,
			HeldOutAccuracy: unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy: controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine: unitaryResult.MeanPhaseCosine,
		},
	}, nil
}
