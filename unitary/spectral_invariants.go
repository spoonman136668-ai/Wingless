package unitary

import (
	"fmt"
	"math"
)

const SpectralInvariantSchema = "wingless.unitary-spectral-invariants.v1"

const (
	spectralMomentCount = 16
	spectralFeatureDim   = spectralMomentCount * 2
	spectralRidgeLambda  = 1e-6
)

type spectralRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type SpectralStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MeanPhaseCosine  float64   `json:"mean_phase_cosine"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type SpectralIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type SpectralDiagnosis struct {
	MomentInvariancePass   bool    `json:"moment_invariance_pass"`
	SpectralLearningPass   bool    `json:"spectral_learning_pass"`
	UnseenDepthPass        bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass bool    `json:"mutable_integration_pass"`
	MaxMomentDrift         float64 `json:"max_moment_drift"`
	TrainAccuracy          float64 `json:"train_accuracy"`
	HeldOutAccuracy        float64 `json:"held_out_accuracy"`
	MatchedControlAccuracy float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine    float64 `json:"mean_held_phase_cosine"`
}

type SpectralInvariantProbeResult struct {
	Schema                          string              `json:"schema"`
	Experiment                      string              `json:"experiment"`
	LatentDimension                 int                 `json:"latent_dimension"`
	RuntimeStateObjects             int                 `json:"runtime_state_objects"`
	VisibleChannelBlocks            bool                `json:"visible_channel_blocks"`
	FullCoordinateMixing            bool                `json:"full_coordinate_mixing"`
	TransportOperatorUsedAsReference bool               `json:"transport_operator_used_as_reference"`
	KnownFullMixerExposedToLearner  bool                `json:"known_full_mixer_exposed_to_learner"`
	RuntimeUnmixApplied             bool                `json:"runtime_unmix_applied"`
	LearnedDemixerUsed              bool                `json:"learned_demixer_used"`
	PhaseAlphabetSupervision        bool                `json:"phase_alphabet_supervision"`
	PhaseCodeDimension              int                 `json:"phase_code_dimension"`
	SpectralMomentCount             int                 `json:"spectral_moment_count"`
	SpectralFeatureDimension        int                 `json:"spectral_feature_dimension"`
	RuntimePrototypeLookup          bool                `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided           bool                `json:"explicit_depth_provided"`
	GlobalPhaseNuisance             bool                `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude            float64             `json:"memory_noise_amplitude"`
	TrainingDepths                  []int               `json:"training_depths"`
	HeldOutDepths                   []int               `json:"held_out_depths"`
	TrainTables                     int                 `json:"train_tables"`
	HeldOutTables                   int                 `json:"held_out_tables"`
	Unitary                         SpectralStaticResult `json:"unitary"`
	MatchedControl                  SpectralStaticResult `json:"matched_control"`
	Integration                     SpectralIntegration  `json:"unitary_mutable_integration"`
	Diagnosis                       SpectralDiagnosis    `json:"diagnosis"`
}

type spectralRow struct {
	features []float64
	table    memoryTable
}

func spectralMomentFeatures(
	state State,
	step latentMatrix,
	count int,
) ([]float64, error) {
	if len(state) != fullLatentDimension || len(step) != fullLatentDimension {
		return nil, fmt.Errorf("spectral moment dimension mismatch")
	}
	if count < 1 {
		return nil, fmt.Errorf("spectral moment count must be positive")
	}

	current := append(State(nil), state...)
	out := make([]float64, 0, count*2)
	for k := 0; k < count; k++ {
		next, err := latentMatrixVector(step, current)
		if err != nil {
			return nil, err
		}
		current = next
		value, err := stateInner(state, current)
		if err != nil {
			return nil, err
		}
		out = append(out, real(value), imag(value))
	}
	for _, value := range out {
		if !finite(value) {
			return nil, fmt.Errorf("spectral moment is non-finite")
		}
	}
	return out, nil
}

func spectralMomentDrift(
	mixer latentMatrix,
	step latentMatrix,
	depthOp latentMatrix,
) (float64, error) {
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 232323, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.471)

	state, err := fullLatentEncode(memory, mixer)
	if err != nil {
		return 0, err
	}
	before, err := spectralMomentFeatures(state, step, spectralMomentCount)
	if err != nil {
		return 0, err
	}
	evolved, err := latentMatrixVector(depthOp, state)
	if err != nil {
		return 0, err
	}
	after, err := spectralMomentFeatures(evolved, step, spectralMomentCount)
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

func buildSpectralRows(
	tables []memoryTable,
	depths []int,
	operators map[int]latentMatrix,
	step latentMatrix,
	mixer latentMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]spectralRow, float64, error) {
	rows := make([]spectralRow, 0, len(tables)*len(depths)*trials)
	var maxNormDrift float64

	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			operator, ok := operators[depth]
			if !ok {
				return nil, 0, fmt.Errorf("missing spectral operator depth=%d", depth)
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
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.251*float64(seed+1), 2*math.Pi),
				)
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
				features, err := spectralMomentFeatures(
					state, step, spectralMomentCount,
				)
				if err != nil {
					return nil, 0, err
				}
				rows = append(rows, spectralRow{
					features: features,
					table:    table,
				})
			}
		}
	}
	return rows, maxNormDrift, nil
}

func trainSpectralRegressors(
	rows []spectralRow,
	lambda float64,
) ([4]spectralRegressor, error) {
	var regressors [4]spectralRegressor
	if len(rows) == 0 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid spectral ridge configuration")
	}

	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)
	for row := 0; row < n; row++ {
		if len(rows[row].features) != spectralFeatureDim {
			return regressors, fmt.Errorf("spectral feature dimension mismatch")
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
		for feature := 0; feature < spectralFeatureDim; feature++ {
			value += rows[first].features[feature] * rows[first].features[feature]
		}
		gram[first][first] = value
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < spectralFeatureDim; feature++ {
				value += rows[first].features[feature] * rows[second].features[feature]
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
			regressors[entity].weights[output] = make([]float64, spectralFeatureDim)
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

func (reg spectralRegressor) predict(features []float64) ([]float64, error) {
	if len(features) != spectralFeatureDim {
		return nil, fmt.Errorf("spectral prediction feature mismatch")
	}
	out := []float64{reg.bias[0], reg.bias[1]}
	for output := 0; output < 2; output++ {
		if len(reg.weights[output]) != spectralFeatureDim {
			return nil, fmt.Errorf("spectral regressor weight mismatch")
		}
		for i, value := range features {
			out[output] += reg.weights[output][i] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("spectral prediction non-finite")
		}
	}
	return out, nil
}

func trainSpectralModel(
	rows []spectralRow,
	lambda float64,
	name, path string,
) (
	[4]spectralRegressor,
	[4]linearSoftmaxHead,
	SpectralStaticResult,
	error,
) {
	var classifiers [4]linearSoftmaxHead
	regressors, err := trainSpectralRegressors(rows, lambda)
	if err != nil {
		return regressors, classifiers, SpectralStaticResult{}, err
	}

	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return regressors, classifiers, SpectralStaticResult{}, err
			}
			samples = append(samples, headSample{
				features: phase,
				target:   row.table[entity],
			})
		}
		head, _, err := trainLinearSoftmax(samples, 4, 2, 1200, 1.0)
		if err != nil {
			return regressors, classifiers, SpectralStaticResult{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, SpectralStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}

	average := 0.0
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return regressors, classifiers, SpectralStaticResult{
		Name:           name,
		Path:           path,
		TrainAccuracy:  average,
		PerEntityTrain: perTrain,
	}, nil
}

func evaluateSpectralModel(
	base SpectralStaticResult,
	regressors [4]spectralRegressor,
	classifiers [4]linearSoftmaxHead,
	rows []spectralRow,
	maxNormDrift float64,
) (SpectralStaticResult, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var cosineTotal float64
	var cosineCount int

	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return SpectralStaticResult{}, err
			}
			target, err := learnedPhaseTarget(row.table[entity])
			if err != nil {
				return SpectralStaticResult{}, err
			}
			norm := math.Hypot(phase[0], phase[1])
			if norm > 1e-15 {
				cosine := (phase[0]*target[0]+phase[1]*target[1])/norm
				if cosine > 1 { cosine = 1 }
				if cosine < -1 { cosine = -1 }
				cosineTotal += cosine
				cosineCount++
			}
			samples = append(samples, headSample{
				features: phase,
				target: row.table[entity],
			})
		}
		_, accuracy, err := evaluateHead(classifiers[entity], samples)
		if err != nil {
			return SpectralStaticResult{}, err
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

func decodeSpectralTable(
	state State,
	step latentMatrix,
	regressors [4]spectralRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	features, err := spectralMomentFeatures(state, step, spectralMomentCount)
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

func runSpectralIntegration(
	mixer latentMatrix,
	step latentMatrix,
	ops map[int]latentMatrix,
	regressors [4]spectralRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (SpectralIntegration, error) {
	const scenarios, writes = 48, 16
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil { return SpectralIntegration{}, err }
	relationHead, _, err := trainLinearSoftmax(relationSamples, 4, 16, 600, 1.0)
	if err != nil { return SpectralIntegration{}, err }

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
			if err != nil { return SpectralIntegration{}, err }
			seed := 29000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil { return SpectralIntegration{}, err }
			memory = rotateGlobalPhase(memory, math.Mod(0.251*float64(seed+1), 2*math.Pi))
			state, err := fullLatentEncode(memory, mixer)
			if err != nil { return SpectralIntegration{}, err }
			state, err = latentMatrixVector(ops[write.gap], state)
			if err != nil { return SpectralIntegration{}, err }
			norm2, err := NormSquared(state)
			if err != nil { return SpectralIntegration{}, err }
			drift := math.Abs(norm2-1)
			if drift > maxNormDrift { maxNormDrift = drift }

			decoded, _, margin, err := decodeSpectralTable(state, step, regressors, classifiers)
			if err != nil { return SpectralIntegration{}, err }
			if margin < minValueMargin { minValueMargin = margin }
			commitTotal++
			if decoded == trueTable { commitCorrect++ }
			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil { return SpectralIntegration{}, err }
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil { return SpectralIntegration{}, err }
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil { return SpectralIntegration{}, err }
		seed := 29000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil { return SpectralIntegration{}, err }
		memory = rotateGlobalPhase(memory, math.Mod(0.251*float64(seed+1), 2*math.Pi))
		state, err := fullLatentEncode(memory, mixer)
		if err != nil { return SpectralIntegration{}, err }
		state, err = latentMatrixVector(ops[scenario.finalGap], state)
		if err != nil { return SpectralIntegration{}, err }
		norm2, err := NormSquared(state)
		if err != nil { return SpectralIntegration{}, err }
		drift := math.Abs(norm2-1)
		if drift > maxNormDrift { maxNormDrift = drift }

		decoded, distributions, margin, err := decodeSpectralTable(state, step, regressors, classifiers)
		if err != nil { return SpectralIntegration{}, err }
		if margin < minValueMargin { minValueMargin = margin }
		if decoded == trueTable { finalCorrect++ }

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil { return SpectralIntegration{}, err }
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil { return SpectralIntegration{}, err }
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil { return SpectralIntegration{}, err }
		if relationMargin < minRelationMargin { minRelationMargin = relationMargin }
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil { return SpectralIntegration{}, err }
		if gotRelation == wantRelation { relationCorrect++ }
	}

	return SpectralIntegration{
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

func RunUP23() (SpectralInvariantProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	mixer := fullLatentMixer()

	unitaryStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	controlStep, err := conjugatedLatentStep(mixer, applyStressNonUnitary)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	controlOps, err := latentDepthOperators(controlStep, allDepths)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	trainRows, _, err := buildSpectralRows(
		trainTables, trainDepths, unitaryOps, unitaryStep,
		mixer, memoryNoise, 4, 0,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	unitRegs, unitHeads, unitBase, err := trainSpectralModel(
		trainRows, spectralRidgeLambda,
		"unitary_transport_spectral_phase_code",
		"unitary_spectral_invariants",
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	heldRows, unitDrift, err := buildSpectralRows(
		heldTables, heldDepths, unitaryOps, unitaryStep,
		mixer, memoryNoise, 2, 7000000,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	unitResult, err := evaluateSpectralModel(
		unitBase, unitRegs, unitHeads, heldRows, unitDrift,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	controlTrainRows, _, err := buildSpectralRows(
		trainTables, trainDepths, controlOps, controlStep,
		mixer, memoryNoise, 4, 0,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	controlRegs, controlHeads, controlBase, err := trainSpectralModel(
		controlTrainRows, spectralRidgeLambda,
		"non_unitary_transport_spectral_phase_code",
		"non_unitary_spectral_invariants",
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	controlHeldRows, controlDrift, err := buildSpectralRows(
		heldTables, heldDepths, controlOps, controlStep,
		mixer, memoryNoise, 2, 7000000,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }
	controlResult, err := evaluateSpectralModel(
		controlBase, controlRegs, controlHeads, controlHeldRows, controlDrift,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	momentDrift, err := spectralMomentDrift(
		mixer, unitaryStep, unitaryOps[1024],
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	integration, err := runSpectralIntegration(
		mixer, unitaryStep, unitaryOps, unitRegs, unitHeads,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil { return SpectralInvariantProbeResult{}, err }

	invariancePass := momentDrift <= 1e-10
	learningPass := unitResult.TrainAccuracy >= 0.99
	unseenPass :=
		unitResult.HeldOutAccuracy >= 0.99 &&
		unitResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
		integration.ExactFinalTableAccuracy >= 0.95 &&
		integration.RelationalQueryAccuracy >= 0.95 &&
		integration.MaxNormDrift <= 1e-10

	return SpectralInvariantProbeResult{
		Schema: SpectralInvariantSchema,
		Experiment: "UP-23-full-latent-transport-spectral-invariants",
		LatentDimension: fullLatentDimension,
		RuntimeStateObjects: 1,
		VisibleChannelBlocks: false,
		FullCoordinateMixing: true,
		TransportOperatorUsedAsReference: true,
		KnownFullMixerExposedToLearner: false,
		RuntimeUnmixApplied: false,
		LearnedDemixerUsed: false,
		PhaseAlphabetSupervision: true,
		PhaseCodeDimension: 2,
		SpectralMomentCount: spectralMomentCount,
		SpectralFeatureDimension: spectralFeatureDim,
		RuntimePrototypeLookup: false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided: false,
		GlobalPhaseNuisance: true,
		MemoryNoiseAmplitude: memoryNoise,
		TrainingDepths: trainDepths,
		HeldOutDepths: heldDepths,
		TrainTables: len(trainTables),
		HeldOutTables: len(heldTables),
		Unitary: unitResult,
		MatchedControl: controlResult,
		Integration: integration,
		Diagnosis: SpectralDiagnosis{
			MomentInvariancePass: invariancePass,
			SpectralLearningPass: learningPass,
			UnseenDepthPass: unseenPass,
			MutableIntegrationPass: mutablePass,
			MaxMomentDrift: momentDrift,
			TrainAccuracy: unitResult.TrainAccuracy,
			HeldOutAccuracy: unitResult.HeldOutAccuracy,
			MatchedControlAccuracy: controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine: unitResult.MeanPhaseCosine,
		},
	}, nil
}
