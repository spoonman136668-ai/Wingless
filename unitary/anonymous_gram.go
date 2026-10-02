package unitary

import (
	"fmt"
	"math"
)

const AnonymousGramSchema = "wingless.unitary-anonymous-gram.v1"

const (
	anonymousGramLinearDim    = 36
	anonymousGramQuadraticDim = 702
)

type AnonymousGramStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type AnonymousGramIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type AnonymousGramDiagnosis struct {
	FeatureInvariancePass       bool    `json:"feature_invariance_pass"`
	DirectQuadraticReadoutPass  bool    `json:"direct_quadratic_readout_pass"`
	UnseenDepthPass             bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass      bool    `json:"mutable_integration_pass"`
	LinearHeldOutAccuracy       float64 `json:"linear_heldout_accuracy"`
	QuadraticTrainAccuracy      float64 `json:"quadratic_train_accuracy"`
	QuadraticHeldOutAccuracy    float64 `json:"quadratic_heldout_accuracy"`
	MatchedControlAccuracy      float64 `json:"matched_control_accuracy"`
	MaxUnitaryFeatureDrift      float64 `json:"max_unitary_feature_drift"`
	MaxNonUnitaryFeatureDrift   float64 `json:"max_nonunitary_feature_drift"`
}

type AnonymousGramProbeResult struct {
	Schema                         string                    `json:"schema"`
	Experiment                     string                    `json:"experiment"`
	CompositeDimension             int                       `json:"composite_dimension"`
	RuntimeStateObjects            int                       `json:"runtime_state_objects"`
	AnonymousChannels              int                       `json:"anonymous_channels"`
	SemanticChannelLocationsKnown  bool                      `json:"semantic_channel_locations_known"`
	KnownMixerExposedToLearner     bool                      `json:"known_mixer_exposed_to_learner"`
	RuntimeMixerInverseApplied     bool                      `json:"runtime_mixer_inverse_applied"`
	LearnedDemixerUsed             bool                      `json:"learned_demixer_used"`
	DirectAnonymousReadout         bool                      `json:"direct_anonymous_readout"`
	RuntimePrototypeLookup         bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                     `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided          bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance            bool                      `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude           float64                   `json:"memory_noise_amplitude"`
	CanonicalTrainingOnly          bool                      `json:"canonical_training_only"`
	CanonicalTrainTrials           int                       `json:"canonical_train_trials_per_table"`
	LinearFeatureDimension         int                       `json:"linear_feature_dimension"`
	QuadraticFeatureDimension      int                       `json:"quadratic_feature_dimension"`
	TrainTables                    int                       `json:"train_tables"`
	HeldOutTables                  int                       `json:"held_out_tables"`
	HeldOutDepths                  []int                     `json:"held_out_depths"`
	LinearUnitary                  AnonymousGramStaticResult `json:"linear_unitary"`
	QuadraticUnitary               AnonymousGramStaticResult `json:"quadratic_unitary"`
	QuadraticMatchedControl        AnonymousGramStaticResult `json:"quadratic_matched_control"`
	Integration                    AnonymousGramIntegration  `json:"unitary_mutable_integration"`
	Diagnosis                      AnonymousGramDiagnosis    `json:"diagnosis"`
}

func anonymousGramFeatures(state State) ([]float64, error) {
	if len(state) != compositeDimension {
		return nil, fmt.Errorf(
			"anonymous gram state dimension=%d want=%d",
			len(state), compositeDimension,
		)
	}

	var channels [compositeChannels]State
	for channel := 0; channel < compositeChannels; channel++ {
		slice, err := compositeChannel(state, channel)
		if err != nil {
			return nil, err
		}
		channels[channel] = slice
	}

	// Each packed channel was scaled by 1/sqrt(6). Multiplying pairwise
	// inner products by 6 restores O(1) source-scale correlations without
	// using the hidden channel mixer.
	scale := float64(compositeChannels)
	features := make([]float64, 0, anonymousGramLinearDim)

	for channel := 0; channel < compositeChannels; channel++ {
		value, err := stateInner(channels[channel], channels[channel])
		if err != nil {
			return nil, err
		}
		features = append(features, scale*real(value))
	}

	for first := 0; first < compositeChannels; first++ {
		for second := first + 1; second < compositeChannels; second++ {
			value, err := stateInner(channels[first], channels[second])
			if err != nil {
				return nil, err
			}
			features = append(
				features,
				scale*real(value),
				scale*imag(value),
			)
		}
	}

	if len(features) != anonymousGramLinearDim {
		return nil, fmt.Errorf(
			"anonymous gram feature dimension=%d want=%d",
			len(features), anonymousGramLinearDim,
		)
	}
	for _, value := range features {
		if !finite(value) {
			return nil, fmt.Errorf("anonymous gram feature is non-finite")
		}
	}
	return features, nil
}

func anonymousQuadraticFeatures(state State) ([]float64, error) {
	base, err := anonymousGramFeatures(state)
	if err != nil {
		return nil, err
	}
	out := make([]float64, 0, anonymousGramQuadraticDim)
	out = append(out, base...)
	for first := 0; first < len(base); first++ {
		for second := first; second < len(base); second++ {
			value := base[first] * base[second]
			if !finite(value) {
				return nil, fmt.Errorf("anonymous quadratic feature is non-finite")
			}
			out = append(out, value)
		}
	}
	if len(out) != anonymousGramQuadraticDim {
		return nil, fmt.Errorf(
			"anonymous quadratic feature dimension=%d want=%d",
			len(out), anonymousGramQuadraticDim,
		)
	}
	return out, nil
}

func anonymousDirectFeatures(state State, quadratic bool) ([]float64, error) {
	if quadratic {
		return anonymousQuadraticFeatures(state)
	}
	return anonymousGramFeatures(state)
}

func buildAnonymousCanonicalSamples(
	tables []memoryTable,
	entity int,
	quadratic bool,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for trial := 0; trial < trials; trial++ {
			seed := seedOffset +
				memoryTableIndex(table)*1000 +
				entity*101 +
				trial*17
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.211*float64(seed+1), 2*math.Pi),
			)
			mixed, err := blindMixedComposite(memory)
			if err != nil {
				return nil, err
			}
			features, err := anonymousDirectFeatures(mixed, quadratic)
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

func buildAnonymousTransportSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	quadratic bool,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, float64, error) {
	var samples []headSample
	var maxNormDrift float64

	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					entity*101 +
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.211*float64(seed+1), 2*math.Pi),
				)
				mixed, err := blindMixedComposite(memory)
				if err != nil {
					return nil, 0, err
				}
				mixed, err = evolveCompositeState(mixed, block, depth, apply)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(mixed)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}

				features, err := anonymousDirectFeatures(mixed, quadratic)
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
	return samples, maxNormDrift, nil
}

func trainAnonymousCanonicalHeads(
	name, path string,
	tables []memoryTable,
	quadratic bool,
	memoryNoise float64,
	trials int,
	steps int,
	learningRate float64,
) ([4]linearSoftmaxHead, AnonymousGramStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	featureDim := anonymousGramLinearDim
	if quadratic {
		featureDim = anonymousGramQuadraticDim
	}

	for entity := 0; entity < 4; entity++ {
		samples, err := buildAnonymousCanonicalSamples(
			tables, entity, quadratic, memoryNoise, trials, 0,
		)
		if err != nil {
			return heads, AnonymousGramStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, featureDim, steps, learningRate,
		)
		if err != nil {
			return heads, AnonymousGramStaticResult{}, err
		}
		heads[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return heads, AnonymousGramStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}

	var average float64
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return heads, AnonymousGramStaticResult{
		Name:           name,
		Path:           path,
		TrainAccuracy:  average,
		PerEntityTrain: perTrain,
	}, nil
}

func evaluateAnonymousHeads(
	base AnonymousGramStaticResult,
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	depths []int,
	quadratic bool,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) (AnonymousGramStaticResult, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var maxNormDrift float64

	for entity := 0; entity < 4; entity++ {
		samples, drift, err := buildAnonymousTransportSamples(
			tables, depths, entity, quadratic,
			block, apply, memoryNoise, trials, seedOffset,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		_, accuracy, err := evaluateHead(heads[entity], samples)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		result.PerEntityHeldOut[entity] = accuracy
	}

	for _, accuracy := range result.PerEntityHeldOut {
		result.HeldOutAccuracy += accuracy
	}
	result.HeldOutAccuracy /= 4
	result.MaxNormDrift = maxNormDrift
	return result, nil
}

func anonymousFeatureDrift(apply stressApply, depth int) (float64, error) {
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 191919, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.417)

	mixed, err := blindMixedComposite(memory)
	if err != nil {
		return 0, err
	}
	before, err := anonymousQuadraticFeatures(mixed)
	if err != nil {
		return 0, err
	}
	afterState, err := evolveCompositeState(
		mixed, stressProgram(), depth, apply,
	)
	if err != nil {
		return 0, err
	}
	after, err := anonymousQuadraticFeatures(afterState)
	if err != nil {
		return 0, err
	}

	var maximum float64
	for i := range before {
		delta := math.Abs(before[i] - after[i])
		if delta > maximum {
			maximum = delta
		}
	}
	return maximum, nil
}

func decodeAnonymousDirectTable(
	state State,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)

	features, err := anonymousQuadraticFeatures(state)
	if err != nil {
		return memoryTable{}, distributions, 0, err
	}

	for entity := 0; entity < 4; entity++ {
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

func runAnonymousGramIntegration(
	heads [4]linearSoftmaxHead,
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
			seed := 21000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return AnonymousGramIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.211*float64(seed+1), 2*math.Pi),
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

			decoded, _, margin, err := decodeAnonymousDirectTable(
				mixed, heads,
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
		seed := 21000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return AnonymousGramIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.211*float64(seed+1), 2*math.Pi),
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

		decoded, distributions, margin, err := decodeAnonymousDirectTable(
			mixed, heads,
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
		relationProbabilities, err := relationHead.probabilities(relationInput)
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

// RunUP19 tests whether reconstructing the historical semantic channel basis is
// unnecessary. It exposes only the anonymous mixed state's complete pairwise
// Gram geometry to learned heads.
//
// The quadratic feature lift is motivated by the existing frame observable:
// after any fixed unknown linear channel mixing, each historical channel
// correlation is linear in the anonymous Gram entries, while the established
// phase-comparison numerator is quadratic in those correlations. A complete
// degree-2 lift therefore tests direct representability without a demixer.
func RunUP19() (AnonymousGramProbeResult, error) {
	const (
		memoryNoise = 0.05
		trainTrials = 2
		heldTrials  = 2
	)
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	heldDepths := []int{32, 128, 512, 1024}
	block := stressProgram()

	linearHeads, linearBase, err := trainAnonymousCanonicalHeads(
		"unitary_anonymous_gram_linear",
		"unitary_linear_baseline",
		trainTables,
		false,
		memoryNoise,
		trainTrials,
		600,
		1.0,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}
	linearUnitary, err := evaluateAnonymousHeads(
		linearBase,
		linearHeads,
		heldTables,
		heldDepths,
		false,
		block,
		applyStressUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}

	quadraticHeads, quadraticBase, err := trainAnonymousCanonicalHeads(
		"unitary_anonymous_gram_quadratic",
		"unitary_direct_quadratic",
		trainTables,
		true,
		memoryNoise,
		trainTrials,
		400,
		0.8,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}
	quadraticUnitary, err := evaluateAnonymousHeads(
		quadraticBase,
		quadraticHeads,
		heldTables,
		heldDepths,
		true,
		block,
		applyStressUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}

	quadraticControl, err := evaluateAnonymousHeads(
		AnonymousGramStaticResult{
			Name:            "non_unitary_anonymous_gram_quadratic",
			Path:            "non_unitary_matched_same_heads",
			TrainAccuracy:   quadraticBase.TrainAccuracy,
			PerEntityTrain:  append([]float64(nil), quadraticBase.PerEntityTrain...),
		},
		quadraticHeads,
		heldTables,
		heldDepths,
		true,
		block,
		applyStressNonUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}

	unitaryFeatureDrift, err := anonymousFeatureDrift(
		applyStressUnitary, 1024,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}
	nonUnitaryFeatureDrift, err := anonymousFeatureDrift(
		applyStressNonUnitary, 1024,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}

	integration, err := runAnonymousGramIntegration(
		quadraticHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return AnonymousGramProbeResult{}, err
	}

	invariancePass := unitaryFeatureDrift <= 1e-12
	directPass := quadraticUnitary.TrainAccuracy >= 0.99
	unseenPass :=
		quadraticUnitary.HeldOutAccuracy >= 0.99 &&
			quadraticUnitary.MaxNormDrift <= 1e-12
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-12

	return AnonymousGramProbeResult{
		Schema:                          AnonymousGramSchema,
		Experiment:                      "UP-19-direct-anonymous-gram-readout",
		CompositeDimension:              compositeDimension,
		RuntimeStateObjects:             1,
		AnonymousChannels:               compositeChannels,
		SemanticChannelLocationsKnown:   false,
		KnownMixerExposedToLearner:      false,
		RuntimeMixerInverseApplied:      false,
		LearnedDemixerUsed:              false,
		DirectAnonymousReadout:          true,
		RuntimePrototypeLookup:          false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided:           false,
		GlobalPhaseNuisance:             true,
		MemoryNoiseAmplitude:            memoryNoise,
		CanonicalTrainingOnly:           true,
		CanonicalTrainTrials:            trainTrials,
		LinearFeatureDimension:          anonymousGramLinearDim,
		QuadraticFeatureDimension:       anonymousGramQuadraticDim,
		TrainTables:                     len(trainTables),
		HeldOutTables:                   len(heldTables),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		LinearUnitary:                    linearUnitary,
		QuadraticUnitary:                 quadraticUnitary,
		QuadraticMatchedControl:          quadraticControl,
		Integration:                      integration,
		Diagnosis: AnonymousGramDiagnosis{
			FeatureInvariancePass:      invariancePass,
			DirectQuadraticReadoutPass: directPass,
			UnseenDepthPass:            unseenPass,
			MutableIntegrationPass:     mutablePass,
			LinearHeldOutAccuracy:      linearUnitary.HeldOutAccuracy,
			QuadraticTrainAccuracy:     quadraticUnitary.TrainAccuracy,
			QuadraticHeldOutAccuracy:   quadraticUnitary.HeldOutAccuracy,
			MatchedControlAccuracy:     quadraticControl.HeldOutAccuracy,
			MaxUnitaryFeatureDrift:     unitaryFeatureDrift,
			MaxNonUnitaryFeatureDrift:  nonUnitaryFeatureDrift,
		},
	}, nil
}
