package unitary

import (
	"fmt"
	"math"
)

const CompositeStateSchema = "wingless.unitary-composite-state.v1"

const (
	compositeChannels   = 6
	compositeChannelDim = 16
	compositeDimension  = compositeChannels * compositeChannelDim
)

type CompositeStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type CompositeIntegrationResult struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxCompositeNormDrift   float64 `json:"max_composite_norm_drift"`
}

type CompositeDiagnosis struct {
	ObservableEquivalencePass bool    `json:"observable_equivalence_pass"`
	StaticFoldPass            bool    `json:"static_fold_pass"`
	MutableIntegrationPass    bool    `json:"mutable_integration_pass"`
	MaxObservableError        float64 `json:"max_observable_error"`
	UnitaryHeldOutAccuracy    float64 `json:"unitary_heldout_accuracy"`
	MatchedControlAccuracy    float64 `json:"matched_control_accuracy"`
	UnitaryMaxNormDrift       float64 `json:"unitary_max_norm_drift"`
}

type CompositeStateProbeResult struct {
	Schema                       string                     `json:"schema"`
	Experiment                   string                     `json:"experiment"`
	BaseDimension                int                        `json:"base_dimension"`
	CompositeDimension           int                        `json:"composite_dimension"`
	LogicalChannelsBeforeFold    int                        `json:"logical_channels_before_fold"`
	RuntimeStateObjects          int                        `json:"runtime_state_objects"`
	MemoryChannelIncluded        bool                       `json:"memory_channel_included"`
	FrameSideVectorsMaterialized bool                       `json:"frame_side_vectors_materialized_after_pack"`
	ChannelSubspacesExplicit     bool                       `json:"channel_subspaces_explicit"`
	RuntimePrototypeLookup       bool                       `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout       bool                       `json:"explicit_inverse_readout"`
	ExplicitDepthProvided        bool                       `json:"explicit_depth_provided"`
	GlobalPhaseNuisance          bool                       `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude         float64                    `json:"memory_noise_amplitude"`
	TrainTables                  int                        `json:"train_tables"`
	HeldOutTables                int                        `json:"held_out_tables"`
	TrainDepths                  []int                      `json:"train_depths"`
	HeldOutDepths                []int                      `json:"held_out_depths"`
	Unitary                      CompositeStaticResult      `json:"unitary"`
	MatchedControl               CompositeStaticResult      `json:"matched_control"`
	Integration                  CompositeIntegrationResult `json:"unitary_mutable_integration"`
	Diagnosis                    CompositeDiagnosis         `json:"diagnosis"`
}

func up15LearnedAnchorParams() AnchorParams {
	return AnchorParams{
		Real: [16]float64{
			0.48458911391697557,
			0.18515591060129735,
			0.2212503413102914,
			0.18384846600205368,
			-0.004967519674986197,
			0.021123557568011823,
			0.06426801319857534,
			0.05499528700223274,
			0.07957246207769357,
			0.09487651115948935,
			0.12384475097548428,
			0.13598986980354394,
			0.02441203804962919,
			0.03660745991244771,
			0.06935852132972692,
			0.03381937769823711,
		},
		Imag: [16]float64{
			2.168404344971009e-19,
			0.09270321850258899,
			0.06564044160408754,
			0.07342942636536494,
			0.20380974300750404,
			0.17712022583374157,
			0.182726193077273,
			0.17872695708093903,
			0.21426772370877345,
			0.19806920429532657,
			0.2123323901119704,
			0.2160470445062389,
			0.2728607185614191,
			0.23673448246080692,
			0.2598822188682088,
			0.22799882501716598,
		},
	}
}

func qualifiedLearnedFrameBank() (frameBank, error) {
	return makeTaskLearnedAnchorBank(up15LearnedAnchorParams())
}

func packCompositeState(memory State, bank frameBank) (State, error) {
	if err := validateState(memory); err != nil {
		return nil, err
	}
	if len(memory) != compositeChannelDim {
		return nil, fmt.Errorf("composite memory dimension=%d want=%d", len(memory), compositeChannelDim)
	}
	states := []State{
		memory,
		bank.anchor,
		bank.phase[0],
		bank.phase[1],
		bank.phase[2],
		bank.phase[3],
	}
	scale := complex(1/math.Sqrt(float64(compositeChannels)), 0)
	out := make(State, 0, compositeDimension)
	for channel, state := range states {
		if err := validateState(state); err != nil {
			return nil, fmt.Errorf("composite channel %d: %w", channel, err)
		}
		if len(state) != compositeChannelDim {
			return nil, fmt.Errorf(
				"composite channel %d dimension=%d want=%d",
				channel, len(state), compositeChannelDim,
			)
		}
		normed, err := Normalize(state)
		if err != nil {
			return nil, err
		}
		for _, value := range normed {
			out = append(out, value*scale)
		}
	}
	return out, nil
}

func compositeChannel(state State, channel int) (State, error) {
	if len(state) != compositeDimension {
		return nil, fmt.Errorf(
			"composite dimension=%d want=%d",
			len(state), compositeDimension,
		)
	}
	if channel < 0 || channel >= compositeChannels {
		return nil, fmt.Errorf("composite channel=%d out of range", channel)
	}
	start := channel * compositeChannelDim
	end := start + compositeChannelDim
	return append(State(nil), state[start:end]...), nil
}

func evolveCompositeState(
	state State,
	block []Coupling,
	depth int,
	apply stressApply,
) (State, error) {
	if len(state) != compositeDimension {
		return nil, fmt.Errorf(
			"composite dimension=%d want=%d",
			len(state), compositeDimension,
		)
	}
	out := make(State, compositeDimension)
	for channel := 0; channel < compositeChannels; channel++ {
		slice, err := compositeChannel(state, channel)
		if err != nil {
			return nil, err
		}
		evolved, err := apply(slice, block, depth)
		if err != nil {
			return nil, err
		}
		copy(
			out[channel*compositeChannelDim:(channel+1)*compositeChannelDim],
			evolved,
		)
	}
	return out, nil
}

func compositeEntityFeature(state State, entity int) ([]float64, error) {
	if entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("composite entity=%d out of range", entity)
	}
	memory, err := compositeChannel(state, 0)
	if err != nil {
		return nil, err
	}
	anchor, err := compositeChannel(state, 1)
	if err != nil {
		return nil, err
	}
	pilot, err := compositeChannel(state, 2+entity)
	if err != nil {
		return nil, err
	}
	return frameFeature(memory, anchor, pilot)
}

func compositeObservableEquivalence(depth int) (float64, error) {
	table := memoryTable{3, 1, 2, 0}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 161616, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.173)

	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return 0, err
	}
	block := stressProgram()

	separateMemory, err := applyStressUnitary(memory, block, depth)
	if err != nil {
		return 0, err
	}
	separateBank, err := transportFrameBank(
		bank, block, depth, applyStressUnitary, true,
	)
	if err != nil {
		return 0, err
	}

	composite, err := packCompositeState(memory, bank)
	if err != nil {
		return 0, err
	}
	composite, err = evolveCompositeState(
		composite, block, depth, applyStressUnitary,
	)
	if err != nil {
		return 0, err
	}

	var maxError float64
	for entity := 0; entity < 4; entity++ {
		separate, err := frameFeature(
			separateMemory,
			separateBank.anchor,
			separateBank.phase[entity],
		)
		if err != nil {
			return 0, err
		}
		folded, err := compositeEntityFeature(composite, entity)
		if err != nil {
			return 0, err
		}
		for i := range separate {
			delta := math.Abs(separate[i] - folded[i])
			if delta > maxError {
				maxError = delta
			}
		}
	}
	return maxError, nil
}

func buildCompositeSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, float64, error) {
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return nil, 0, err
	}
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
				memory, err := perturbMemory(
					canonical, seed, memoryNoise,
				)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.197*float64(seed+1), 2*math.Pi),
				)
				composite, err := packCompositeState(memory, bank)
				if err != nil {
					return nil, 0, err
				}
				composite, err = evolveCompositeState(
					composite, block, depth, apply,
				)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(composite)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}
				features, err := compositeEntityFeature(
					composite, entity,
				)
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

func trainCompositeHeads(
	name, path string,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	steps int,
) ([4]linearSoftmaxHead, CompositeStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)
	var maxNormDrift float64

	for entity := 0; entity < 4; entity++ {
		trainSamples, trainDrift, err := buildCompositeSamples(
			trainTables, trainDepths, entity,
			block, apply, memoryNoise, 4, 0,
		)
		if err != nil {
			return heads, CompositeStaticResult{}, err
		}
		heldSamples, heldDrift, err := buildCompositeSamples(
			heldTables, heldDepths, entity,
			block, apply, memoryNoise, 2, 7000000,
		)
		if err != nil {
			return heads, CompositeStaticResult{}, err
		}
		if trainDrift > maxNormDrift {
			maxNormDrift = trainDrift
		}
		if heldDrift > maxNormDrift {
			maxNormDrift = heldDrift
		}

		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, CompositeStaticResult{}, err
		}
		heads[entity] = head
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return heads, CompositeStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return heads, CompositeStaticResult{}, err
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

	return heads, CompositeStaticResult{
		Name:             name,
		Path:             path,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
		MaxNormDrift:     maxNormDrift,
	}, nil
}

func decodeCompositeTable(
	state State,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := compositeEntityFeature(state, entity)
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

func runCompositeIntegration(
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (CompositeIntegrationResult, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return CompositeIntegrationResult{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return CompositeIntegrationResult{}, err
	}
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return CompositeIntegrationResult{}, err
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
				return CompositeIntegrationResult{}, err
			}
			seed := 17000000 +
				scenarioIndex*10000 +
				writeIndex*31
			memory, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return CompositeIntegrationResult{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.197*float64(seed+1), 2*math.Pi),
			)
			composite, err := packCompositeState(memory, bank)
			if err != nil {
				return CompositeIntegrationResult{}, err
			}
			composite, err = evolveCompositeState(
				composite, block, write.gap,
				applyStressUnitary,
			)
			if err != nil {
				return CompositeIntegrationResult{}, err
			}
			norm2, err := NormSquared(composite)
			if err != nil {
				return CompositeIntegrationResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeCompositeTable(
				composite, heads,
			)
			if err != nil {
				return CompositeIntegrationResult{}, err
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
				return CompositeIntegrationResult{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return CompositeIntegrationResult{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		seed := 17000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(
			canonical, seed, memoryNoise,
		)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.197*float64(seed+1), 2*math.Pi),
		)
		composite, err := packCompositeState(memory, bank)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		composite, err = evolveCompositeState(
			composite, block, scenario.finalGap,
			applyStressUnitary,
		)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		norm2, err := NormSquared(composite)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		decoded, distributions, margin, err := decodeCompositeTable(
			composite, heads,
		)
		if err != nil {
			return CompositeIntegrationResult{}, err
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
			return CompositeIntegrationResult{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return CompositeIntegrationResult{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return CompositeIntegrationResult{
		Scenarios:               scenarios,
		WritesPerScenario:       writes,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect) / float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(scenarios),
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
		MaxCompositeNormDrift:   maxNormDrift,
	}, nil
}

// RunUP16 folds memory plus the complete learned five-state frame into one
// 96-dimensional recurrent complex state.
//
// Logical subspaces remain explicit in this experiment. The scientific question
// is narrower: can one state object carry and co-evolve the memory/reference
// system without any separate runtime frame vectors while preserving the exact
// gauge-invariant observables, unseen-depth decoding, and mutable behavior?
func RunUP16() (CompositeStateProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	block := stressProgram()

	maxObservableError := 0.0
	for _, depth := range []int{1, 8, 128, 1024} {
		errValue, err := compositeObservableEquivalence(depth)
		if err != nil {
			return CompositeStateProbeResult{}, err
		}
		if errValue > maxObservableError {
			maxObservableError = errValue
		}
	}

	unitaryHeads, unitaryResult, err := trainCompositeHeads(
		"unitary_single_composite_state",
		"unitary",
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
		return CompositeStateProbeResult{}, err
	}

	_, controlResult, err := trainCompositeHeads(
		"non_unitary_single_composite_state",
		"non_unitary_matched",
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
		return CompositeStateProbeResult{}, err
	}

	integration, err := runCompositeIntegration(
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return CompositeStateProbeResult{}, err
	}

	equivalencePass := maxObservableError <= 1e-12
	staticPass :=
		equivalencePass &&
			unitaryResult.HeldOutAccuracy >= 0.99 &&
			unitaryResult.MaxNormDrift <= 1e-12
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxCompositeNormDrift <= 1e-12

	return CompositeStateProbeResult{
		Schema:                       CompositeStateSchema,
		Experiment:                   "UP-16-single-composite-recurrent-state",
		BaseDimension:                compositeChannelDim,
		CompositeDimension:           compositeDimension,
		LogicalChannelsBeforeFold:    compositeChannels,
		RuntimeStateObjects:          1,
		MemoryChannelIncluded:        true,
		FrameSideVectorsMaterialized: false,
		ChannelSubspacesExplicit:     true,
		RuntimePrototypeLookup:       false,
		ExplicitInverseReadout:       false,
		ExplicitDepthProvided:        false,
		GlobalPhaseNuisance:          true,
		MemoryNoiseAmplitude:         memoryNoise,
		TrainTables:                  len(trainTables),
		HeldOutTables:                len(heldTables),
		TrainDepths:                  append([]int(nil), trainDepths...),
		HeldOutDepths:                append([]int(nil), heldDepths...),
		Unitary:                      unitaryResult,
		MatchedControl:               controlResult,
		Integration:                  integration,
		Diagnosis: CompositeDiagnosis{
			ObservableEquivalencePass: equivalencePass,
			StaticFoldPass:            staticPass,
			MutableIntegrationPass:    mutablePass,
			MaxObservableError:        maxObservableError,
			UnitaryHeldOutAccuracy:    unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:    controlResult.HeldOutAccuracy,
			UnitaryMaxNormDrift:       unitaryResult.MaxNormDrift,
		},
	}, nil
}
