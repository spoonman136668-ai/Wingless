package unitary

import (
	"fmt"
	"math"
)

const PhaseCodeSchema = "wingless.unitary-phase-code.v1"

type phaseCodeRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type PhaseCodeDiagnosis struct {
	PhaseCodeLearningPass    bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	TrainAccuracy            float64 `json:"train_accuracy"`
	HeldOutAccuracy          float64 `json:"held_out_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine      float64 `json:"mean_held_phase_cosine"`
	MaxUnitaryFeatureDrift   float64 `json:"max_unitary_feature_drift"`
}

type PhaseCodeProbeResult struct {
	Schema                          string                    `json:"schema"`
	Experiment                      string                    `json:"experiment"`
	CompositeDimension              int                       `json:"composite_dimension"`
	RuntimeStateObjects             int                       `json:"runtime_state_objects"`
	AnonymousChannels               int                       `json:"anonymous_channels"`
	SemanticChannelLocationsKnown   bool                      `json:"semantic_channel_locations_known"`
	KnownMixerExposedToLearner      bool                      `json:"known_mixer_exposed_to_learner"`
	OracleUsedForLearning           bool                      `json:"oracle_used_for_learning"`
	RuntimeMixerInverseApplied      bool                      `json:"runtime_mixer_inverse_applied"`
	LearnedDemixerUsed              bool                      `json:"learned_demixer_used"`
	DirectAnonymousReadout          bool                      `json:"direct_anonymous_readout"`
	PhaseAlphabetSupervision        bool                      `json:"phase_alphabet_supervision"`
	PhaseCodeDimension              int                       `json:"phase_code_dimension"`
	RuntimePrototypeLookup          bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                      `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided           bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance             bool                      `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude            float64                   `json:"memory_noise_amplitude"`
	QuadraticFeatureDimension       int                       `json:"quadratic_feature_dimension"`
	RidgeLambda                     float64                   `json:"ridge_lambda"`
	TrainingTrials                  int                       `json:"training_trials_per_table"`
	TrainTables                     int                       `json:"train_tables"`
	HeldOutTables                   int                       `json:"held_out_tables"`
	HeldOutDepths                   []int                     `json:"held_out_depths"`
	Unitary                         AnonymousGramStaticResult `json:"unitary"`
	MatchedControl                  AnonymousGramStaticResult `json:"matched_control"`
	Integration                     AnonymousGramIntegration  `json:"unitary_mutable_integration"`
	Diagnosis                       PhaseCodeDiagnosis        `json:"diagnosis"`
}

func learnedPhaseTarget(value int) ([]float64, error) {
	phases := up13LearnedPhases()
	if value < 0 || value >= len(phases) {
		return nil, fmt.Errorf("phase target value=%d out of range", value)
	}
	return []float64{
		math.Cos(phases[value]),
		math.Sin(phases[value]),
	}, nil
}

func (regressor phaseCodeRegressor) predict(features []float64) ([]float64, error) {
	if len(features) != anonymousGramQuadraticDim {
		return nil, fmt.Errorf("phase regressor feature dimension mismatch")
	}
	out := []float64{regressor.bias[0], regressor.bias[1]}
	for output := 0; output < 2; output++ {
		if len(regressor.weights[output]) != anonymousGramQuadraticDim {
			return nil, fmt.Errorf("phase regressor weight dimension mismatch")
		}
		for feature, value := range features {
			out[output] += regressor.weights[output][feature] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("phase regressor output is non-finite")
		}
	}
	return out, nil
}

func trainPhaseCodeRegressors(
	rows []ridgeAnonymousRow,
	lambda float64,
) ([4]phaseCodeRegressor, error) {
	var regressors [4]phaseCodeRegressor
	if len(rows) == 0 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid phase-code regression configuration")
	}

	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)

	for row := 0; row < n; row++ {
		if len(rows[row].features) != anonymousGramQuadraticDim {
			return regressors, fmt.Errorf("phase-code feature dimension mismatch")
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
		for feature := 0; feature < anonymousGramQuadraticDim; feature++ {
			value += rows[first].features[feature] *
				rows[first].features[feature]
		}
		gram[first][first] = value

		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < anonymousGramQuadraticDim; feature++ {
				value += rows[first].features[feature] *
					rows[second].features[feature]
			}
			if !finite(value) {
				return regressors, fmt.Errorf("phase-code gram value is non-finite")
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
			regressors[entity].weights[output] =
				make([]float64, anonymousGramQuadraticDim)
		}
	}

	for row := 0; row < n; row++ {
		for entity := 0; entity < 4; entity++ {
			for output := 0; output < 2; output++ {
				coefficient := alpha[row][entity*2+output]
				regressors[entity].bias[output] += coefficient
				for feature, value := range rows[row].features {
					regressors[entity].weights[output][feature] +=
						coefficient * value
				}
			}
		}
	}

	return regressors, nil
}

func phaseCodeSamplesFromRows(
	rows []ridgeAnonymousRow,
	entity int,
	regressor phaseCodeRegressor,
) ([]headSample, error) {
	samples := make([]headSample, 0, len(rows))
	for _, row := range rows {
		features, err := regressor.predict(row.features)
		if err != nil {
			return nil, err
		}
		samples = append(samples, headSample{
			features: features,
			target:   row.table[entity],
		})
	}
	return samples, nil
}

func trainPhaseCodeModel(
	rows []ridgeAnonymousRow,
	lambda float64,
) (
	[4]phaseCodeRegressor,
	[4]linearSoftmaxHead,
	AnonymousGramStaticResult,
	error,
) {
	var classifiers [4]linearSoftmaxHead

	regressors, err := trainPhaseCodeRegressors(rows, lambda)
	if err != nil {
		return regressors, classifiers, AnonymousGramStaticResult{}, err
	}

	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples, err := phaseCodeSamplesFromRows(
			rows, entity, regressors[entity],
		)
		if err != nil {
			return regressors, classifiers, AnonymousGramStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, 2, 1200, 1.0,
		)
		if err != nil {
			return regressors, classifiers, AnonymousGramStaticResult{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, AnonymousGramStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}

	var average float64
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return regressors, classifiers, AnonymousGramStaticResult{
		Name:           "unitary_anonymous_phase_code",
		Path:           "unitary_phase_code_bottleneck",
		TrainAccuracy:  average,
		PerEntityTrain: perTrain,
	}, nil
}

func transformAnonymousSamplesToPhaseCode(
	samples []headSample,
	regressor phaseCodeRegressor,
) ([]headSample, error) {
	out := make([]headSample, 0, len(samples))
	for _, sample := range samples {
		features, err := regressor.predict(sample.features)
		if err != nil {
			return nil, err
		}
		out = append(out, headSample{
			features: features,
			target:   sample.target,
		})
	}
	return out, nil
}

func evaluatePhaseCodeModel(
	base AnonymousGramStaticResult,
	regressors [4]phaseCodeRegressor,
	classifiers [4]linearSoftmaxHead,
	tables []memoryTable,
	depths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) (AnonymousGramStaticResult, float64, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var maxNormDrift float64
	var cosineTotal float64
	var cosineCount int

	for entity := 0; entity < 4; entity++ {
		rawSamples, drift, err := buildAnonymousTransportSamples(
			tables,
			depths,
			entity,
			true,
			block,
			apply,
			memoryNoise,
			trials,
			seedOffset,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, 0, err
		}
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		phaseSamples, err := transformAnonymousSamplesToPhaseCode(
			rawSamples, regressors[entity],
		)
		if err != nil {
			return AnonymousGramStaticResult{}, 0, err
		}
		_, accuracy, err := evaluateHead(
			classifiers[entity], phaseSamples,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, 0, err
		}
		result.PerEntityHeldOut[entity] = accuracy

		for _, sample := range phaseSamples {
			target, err := learnedPhaseTarget(sample.target)
			if err != nil {
				return AnonymousGramStaticResult{}, 0, err
			}
			predNorm := math.Hypot(sample.features[0], sample.features[1])
			if predNorm <= 1e-15 {
				continue
			}
			cosine := (sample.features[0]*target[0] +
				sample.features[1]*target[1]) / predNorm
			if cosine > 1 {
				cosine = 1
			}
			if cosine < -1 {
				cosine = -1
			}
			cosineTotal += cosine
			cosineCount++
		}
	}

	for _, accuracy := range result.PerEntityHeldOut {
		result.HeldOutAccuracy += accuracy
	}
	result.HeldOutAccuracy /= 4
	result.MaxNormDrift = maxNormDrift

	meanCosine := 0.0
	if cosineCount > 0 {
		meanCosine = cosineTotal / float64(cosineCount)
	}
	return result, meanCosine, nil
}

func decodePhaseCodeTable(
	state State,
	regressors [4]phaseCodeRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)

	features, err := anonymousQuadraticFeatures(state)
	if err != nil {
		return memoryTable{}, distributions, 0, err
	}

	for entity := 0; entity < 4; entity++ {
		phaseFeatures, err := regressors[entity].predict(features)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err :=
			classifiers[entity].probabilities(phaseFeatures)
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

func runPhaseCodeIntegration(
	regressors [4]phaseCodeRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (AnonymousGramIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return AnonymousGramIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return AnonymousGramIntegration{}, err
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
				return AnonymousGramIntegration{}, err
			}
			seed := 25000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.229*float64(seed+1), 2*math.Pi),
			)
			mixed, err := blindMixedComposite(memory)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
			mixed, err = evolveCompositeState(
				mixed, block, write.gap, applyStressUnitary,
			)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
			norm2, err := NormSquared(mixed)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodePhaseCodeTable(
				mixed, regressors, classifiers,
			)
			if err != nil {
				return AnonymousGramIntegration{}, err
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
				return AnonymousGramIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		seed := 25000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.229*float64(seed+1), 2*math.Pi),
		)
		mixed, err := blindMixedComposite(memory)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		mixed, err = evolveCompositeState(
			mixed, block, scenario.finalGap, applyStressUnitary,
		)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		norm2, err := NormSquared(mixed)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		decoded, distributions, margin, err := decodePhaseCodeTable(
			mixed, regressors, classifiers,
		)
		if err != nil {
			return AnonymousGramIntegration{}, err
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
			return AnonymousGramIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return AnonymousGramIntegration{
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

func RunUP21() (PhaseCodeProbeResult, error) {
	const (
		memoryNoise = 0.05
		trainTrials = 4
		heldTrials  = 2
	)
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	heldDepths := []int{32, 128, 512, 1024}
	block := stressProgram()

	rows, err := buildRidgeAnonymousRows(
		trainTables, memoryNoise, trainTrials, 0,
	)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	regressors, classifiers, base, err :=
		trainPhaseCodeModel(rows, ridgeAnonymousLambda)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	unitaryResult, meanCosine, err := evaluatePhaseCodeModel(
		base,
		regressors,
		classifiers,
		heldTables,
		heldDepths,
		block,
		applyStressUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	controlResult, _, err := evaluatePhaseCodeModel(
		AnonymousGramStaticResult{
			Name:           "non_unitary_anonymous_phase_code",
			Path:           "non_unitary_matched_same_model",
			TrainAccuracy:  base.TrainAccuracy,
			PerEntityTrain: append([]float64(nil), base.PerEntityTrain...),
		},
		regressors,
		classifiers,
		heldTables,
		heldDepths,
		block,
		applyStressNonUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	featureDrift, err := anonymousFeatureDrift(
		applyStressUnitary, 1024,
	)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	integration, err := runPhaseCodeIntegration(
		regressors,
		classifiers,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return PhaseCodeProbeResult{}, err
	}

	learningPass := unitaryResult.TrainAccuracy >= 0.99
	unseenPass :=
		unitaryResult.HeldOutAccuracy >= 0.99 &&
			unitaryResult.MaxNormDrift <= 1e-12
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-12

	return PhaseCodeProbeResult{
		Schema:                          PhaseCodeSchema,
		Experiment:                      "UP-21-anonymous-phase-code-bottleneck",
		CompositeDimension:              compositeDimension,
		RuntimeStateObjects:             1,
		AnonymousChannels:               compositeChannels,
		SemanticChannelLocationsKnown:   false,
		KnownMixerExposedToLearner:      false,
		OracleUsedForLearning:           false,
		RuntimeMixerInverseApplied:      false,
		LearnedDemixerUsed:              false,
		DirectAnonymousReadout:          true,
		PhaseAlphabetSupervision:        true,
		PhaseCodeDimension:              2,
		RuntimePrototypeLookup:          false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided:           false,
		GlobalPhaseNuisance:             true,
		MemoryNoiseAmplitude:            memoryNoise,
		QuadraticFeatureDimension:       anonymousGramQuadraticDim,
		RidgeLambda:                     ridgeAnonymousLambda,
		TrainingTrials:                  trainTrials,
		TrainTables:                     len(trainTables),
		HeldOutTables:                   len(heldTables),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		Unitary:                         unitaryResult,
		MatchedControl:                  controlResult,
		Integration:                     integration,
		Diagnosis: PhaseCodeDiagnosis{
			PhaseCodeLearningPass:  learningPass,
			UnseenDepthPass:        unseenPass,
			MutableIntegrationPass: mutablePass,
			TrainAccuracy:          unitaryResult.TrainAccuracy,
			HeldOutAccuracy:        unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy: controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine:    meanCosine,
			MaxUnitaryFeatureDrift: featureDrift,
		},
	}, nil
}
