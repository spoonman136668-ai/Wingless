package unitary

import (
	"fmt"
	"math"
)

const FullRankIntegrationSchema = "wingless.unitary-fullrank-integration.v1"

type FullRankStaticResult struct {
	Name                string    `json:"name"`
	Path                string    `json:"path"`
	TrainTables         int       `json:"train_tables"`
	HeldOutTables       int       `json:"held_out_tables"`
	TrainTrials         int       `json:"train_trials_per_table_depth"`
	HeldOutTrials       int       `json:"held_out_trials_per_table_depth"`
	OptimizerSteps      int       `json:"optimizer_steps"`
	TrainAccuracy       float64   `json:"train_accuracy"`
	HeldOutAccuracy     float64   `json:"held_out_accuracy"`
	PerEntityTrain      []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut    []float64 `json:"per_entity_held_out_accuracy"`
}

type FullRankIntegrationResult struct {
	Name                    string  `json:"name"`
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type FullRankDiagnosis struct {
	StaticSaturationPass       bool    `json:"static_saturation_pass"`
	MutableIntegrationPass     bool    `json:"mutable_integration_pass"`
	UP11Rank4NoisyAccuracy     float64 `json:"up11_rank4_noisy_accuracy"`
	SaturatedStaticAccuracy    float64 `json:"saturated_static_accuracy"`
	StaticAccuracyGain         float64 `json:"static_accuracy_gain"`
	MatchedControlAccuracy     float64 `json:"matched_control_accuracy"`
	DecoderDataLimitSupported  bool    `json:"decoder_data_limit_supported"`
	NonlinearReadoutIndicated  bool    `json:"nonlinear_readout_indicated"`
}

type FullRankProbeResult struct {
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
	MemoryNoiseAmplitude    float64                   `json:"memory_noise_amplitude"`
	FullTrainPoolTables     int                       `json:"full_train_pool_tables"`
	FullHeldOutPoolTables   int                       `json:"full_held_out_pool_tables"`
	MinTrainMarginalCount   int                       `json:"min_train_marginal_count"`
	MinHeldOutMarginalCount int                       `json:"min_held_out_marginal_count"`
	TrainDepths             []int                     `json:"train_depths"`
	HeldOutDepths           []int                     `json:"held_out_depths"`
	UnitaryStatic           FullRankStaticResult      `json:"unitary_static"`
	ControlStatic           FullRankStaticResult      `json:"control_static"`
	Integration             FullRankIntegrationResult `json:"unitary_mutable_integration"`
	Diagnosis               FullRankDiagnosis         `json:"diagnosis"`
}

func fullObserverTablePool(train bool) []memoryTable {
	out := make([]memoryTable, 0, 128)
	for _, table := range allMemoryTables() {
		if balancedObserverTrainTable(table) == train {
			out = append(out, table)
		}
	}
	return out
}

func buildFullRankSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	if trials < 1 {
		return nil, fmt.Errorf("full-rank trials must be positive")
	}
	const rank = 4
	baseBank, err := makePilotRankBank(rank)
	if err != nil {
		return nil, err
	}
	transformed := make(map[int]pilotRankBank, len(depths))
	for _, depth := range depths {
		bank, err := transportPilotRankBank(baseBank, block, depth, apply)
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
				state, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, err
				}
				state = rotateGlobalPhase(
					state,
					math.Mod(0.223*float64(seed+1), 2*math.Pi),
				)
				forward, err := apply(state, block, depth)
				if err != nil {
					return nil, err
				}
				features, err := pilotRankFeatures(forward, transformed[depth])
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

func trainFullRankHeads(
	name, path string,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trainTrials, heldTrials, steps int,
) ([4]linearSoftmaxHead, FullRankStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildFullRankSamples(
			trainTables,
			trainDepths,
			entity,
			block,
			apply,
			memoryNoise,
			trainTrials,
			0,
		)
		if err != nil {
			return heads, FullRankStaticResult{}, err
		}
		heldSamples, err := buildFullRankSamples(
			heldTables,
			heldDepths,
			entity,
			block,
			apply,
			memoryNoise,
			heldTrials,
			7000000,
		)
		if err != nil {
			return heads, FullRankStaticResult{}, err
		}

		head, _, err := trainLinearSoftmax(
			trainSamples,
			4,
			8,
			steps,
			1.0,
		)
		if err != nil {
			return heads, FullRankStaticResult{}, err
		}
		heads[entity] = head

		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return heads, FullRankStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return heads, FullRankStaticResult{}, err
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

	return heads, FullRankStaticResult{
		Name:             name,
		Path:             path,
		TrainTables:      len(trainTables),
		HeldOutTables:    len(heldTables),
		TrainTrials:      trainTrials,
		HeldOutTrials:    heldTrials,
		OptimizerSteps:   steps,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func decodeFullRankTable(
	memory State,
	bank pilotRankBank,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	features, err := pilotRankFeatures(memory, bank)
	if err != nil {
		return memoryTable{}, [4][]float64{}, 0, err
	}
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
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

func makeFullRankScenario(seed int, initial memoryTable, writes int, depths []int) memoryScenario {
	ops := make([]memoryWrite, 0, writes)
	for j := 0; j < writes; j++ {
		entity := (seed + 3*j + j/2) % 4
		value := (seed*2 + j*j + entity + 3) % 4
		gap := depths[(seed*3+j)%len(depths)]
		ops = append(ops, memoryWrite{
			entity: entity,
			value:  value,
			gap:    gap,
		})
	}
	return memoryScenario{
		initial:  initial,
		writes:   ops,
		queryA:   (seed + 2) % 4,
		queryB:   (seed + 1) % 4,
		finalGap: depths[(seed+writes*2)%len(depths)],
	}
}

func runFullRankIntegration(
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (FullRankIntegrationResult, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return FullRankIntegrationResult{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples,
		4,
		16,
		600,
		1.0,
	)
	if err != nil {
		return FullRankIntegrationResult{}, err
	}

	baseBank, err := makePilotRankBank(4)
	if err != nil {
		return FullRankIntegrationResult{}, err
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
				return FullRankIntegrationResult{}, err
			}
			seed := 9000000 + scenarioIndex*10000 + writeIndex*31
			state, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			state = rotateGlobalPhase(
				state,
				math.Mod(0.223*float64(seed+1), 2*math.Pi),
			)
			forward, err := applyStressUnitary(state, block, write.gap)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			frame, err := transportPilotRankBank(
				baseBank,
				block,
				write.gap,
				applyStressUnitary,
			)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			decoded, _, margin, err := decodeFullRankTable(forward, frame, heads)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			if margin < minValueMargin {
				minValueMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}

			pathTable, err = applyMemoryWrite(
				decoded,
				write.entity,
				write.value,
			)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable,
				write.entity,
				write.value,
			)
			if err != nil {
				return FullRankIntegrationResult{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		seed := 9000000 + scenarioIndex*10000 + 9999
		state, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		state = rotateGlobalPhase(
			state,
			math.Mod(0.223*float64(seed+1), 2*math.Pi),
		)
		forward, err := applyStressUnitary(state, block, scenario.finalGap)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		frame, err := transportPilotRankBank(
			baseBank,
			block,
			scenario.finalGap,
			applyStressUnitary,
		)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		decoded, distributions, margin, err := decodeFullRankTable(
			forward,
			frame,
			heads,
		)
		if err != nil {
			return FullRankIntegrationResult{}, err
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
			return FullRankIntegrationResult{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable,
			scenario.queryA,
			scenario.queryB,
		)
		if err != nil {
			return FullRankIntegrationResult{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return FullRankIntegrationResult{
		Name:                    "unitary_rank4_mutable_integration",
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

// RunUP12 keeps the full rank-4 internal frame required by UP-11 and asks
// whether the remaining rank-4 errors are due to sparse decoder training.
//
// Unlike UP-11's 32/32 table selection with one noisy sample per table/depth,
// UP-12 uses the complete balanced 128-table training pool, the complete
// disjoint 128-table held-out pool, four independent noisy trials for every
// training table/depth, two held-out trials, and a longer deterministic
// optimizer budget.
//
// If the saturated linear decoder reaches >=0.99 unseen-depth accuracy, the
// rank-4 representation is treated as sufficient under the current noise
// budget and is placed into a 48-scenario mutable read/write integration.
func RunUP12() (FullRankProbeResult, error) {
	const (
		memoryNoise  = 0.05
		trainTrials  = 4
		heldTrials   = 2
		steps        = 1200
		up11Accuracy = 0.95703125
	)
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	if len(trainTables) != 128 || len(heldTables) != 128 {
		return FullRankProbeResult{}, fmt.Errorf(
			"full observer pool sizes train=%d held=%d want=128/128",
			len(trainTables),
			len(heldTables),
		)
	}
	minTrain := observerMarginalMinimum(trainTables)
	minHeld := observerMarginalMinimum(heldTables)
	if minTrain <= 0 || minHeld <= 0 {
		return FullRankProbeResult{}, fmt.Errorf("full observer pool lost a marginal")
	}

	block := stressProgram()
	unitaryHeads, unitaryStatic, err := trainFullRankHeads(
		"unitary_rank4_saturated_linear",
		"unitary",
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressUnitary,
		memoryNoise,
		trainTrials,
		heldTrials,
		steps,
	)
	if err != nil {
		return FullRankProbeResult{}, err
	}

	_, controlStatic, err := trainFullRankHeads(
		"non_unitary_rank4_saturated_linear",
		"non_unitary_matched",
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressNonUnitary,
		memoryNoise,
		trainTrials,
		heldTrials,
		steps,
	)
	if err != nil {
		return FullRankProbeResult{}, err
	}

	integration, err := runFullRankIntegration(
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return FullRankProbeResult{}, err
	}

	staticPass := unitaryStatic.HeldOutAccuracy >= 0.99
	integrationPass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95
	gain := unitaryStatic.HeldOutAccuracy - up11Accuracy

	return FullRankProbeResult{
		Schema:                  FullRankIntegrationSchema,
		Experiment:              "UP-12-fullrank-decoder-saturation-integration",
		Dimension:               16,
		Entities:                4,
		ValuesPerEntity:         4,
		CodeRank:                4,
		PilotStates:             5,
		PilotComplexScalars:     80,
		FeaturesPerEntity:       8,
		RuntimePrototypeLookup:  false,
		ExplicitInverseReadout:  false,
		ExplicitDepthProvided:   false,
		GlobalPhaseNuisance:     true,
		MemoryNoiseAmplitude:    memoryNoise,
		FullTrainPoolTables:     len(trainTables),
		FullHeldOutPoolTables:   len(heldTables),
		MinTrainMarginalCount:   minTrain,
		MinHeldOutMarginalCount: minHeld,
		TrainDepths:             append([]int(nil), trainDepths...),
		HeldOutDepths:           append([]int(nil), heldDepths...),
		UnitaryStatic:           unitaryStatic,
		ControlStatic:           controlStatic,
		Integration:             integration,
		Diagnosis: FullRankDiagnosis{
			StaticSaturationPass:      staticPass,
			MutableIntegrationPass:    integrationPass,
			UP11Rank4NoisyAccuracy:    up11Accuracy,
			SaturatedStaticAccuracy:   unitaryStatic.HeldOutAccuracy,
			StaticAccuracyGain:        gain,
			MatchedControlAccuracy:    controlStatic.HeldOutAccuracy,
			DecoderDataLimitSupported: staticPass && gain >= 0.02,
			NonlinearReadoutIndicated: !staticPass,
		},
	}, nil
}
