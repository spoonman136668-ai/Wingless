package unitary

import (
	"fmt"
	"math"
)

const UP52CMemoryLoadSchema = "wingless.up52c-memory-load-interference.v1"

type UP52CLoadPoint struct {
	WritesPerScenario int                   `json:"writes_per_scenario"`
	Integration       DiscoveredIntegration `json:"integration"`
	HardGate          bool                  `json:"hard_gate"`
}

type UP52CMemoryLoadResult struct {
	Schema             string                `json:"schema"`
	Experiment         string                `json:"experiment"`
	SourceUP51CSeal    string                `json:"source_up51c_seal"`
	MemoryNoise        float64               `json:"memory_noise"`
	HeldDepths         []int                 `json:"held_depths"`
	WriteLevels        []int                 `json:"write_levels"`
	StaticHeldAccuracy float64               `json:"static_held_accuracy"`
	SelectedIndices    []int                 `json:"selected_indices"`
	Points             []UP52CLoadPoint      `json:"points"`
	LastPassingWrites  int                   `json:"last_passing_writes"`
	FirstFailingWrites int                   `json:"first_failing_writes"`
}

func runUP52CIntegration(
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
) (DiscoveredIntegration, error) {
	const scenarios = 48

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return DiscoveredIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return DiscoveredIntegration{}, err
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
				return DiscoveredIntegration{}, err
			}
			seed := 52000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return DiscoveredIntegration{}, fmt.Errorf(
					"missing UP52C mutable depth=%d", write.gap,
				)
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			norm2, err := NormSquared(state)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			if drift := math.Abs(norm2 - 1); drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeBreadthTable(
				state, observables, regressors, classifiers,
			)
			if err != nil {
				return DiscoveredIntegration{}, err
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
				return DiscoveredIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		seed := 52000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1), 2*math.Pi),
		)
		state, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return DiscoveredIntegration{}, fmt.Errorf(
				"missing UP52C final depth=%d", scenario.finalGap,
			)
		}
		state, err = latentMatrixVector(operator, state)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		norm2, err := NormSquared(state)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		if drift := math.Abs(norm2 - 1); drift > maxNormDrift {
			maxNormDrift = drift
		}

		decoded, distributions, margin, err := decodeBreadthTable(
			state, observables, regressors, classifiers,
		)
		if err != nil {
			return DiscoveredIntegration{}, err
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
			return DiscoveredIntegration{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return DiscoveredIntegration{
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

func RunUP52C() (UP52CMemoryLoadResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32,128,512,1024}
	allDepths := []int{0,32,128,512,1024}
	writeLevels := []int{16,32,64,128,256}

	mixer := fullLatentMixer()
	offsets := []float64{0,0,0,0,0,0}
	step, err := fullLatentDoseStep(mixer, offsets)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}
	ops, err := latentDepthOperators(step, allDepths)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	candidates, _, err := discoverCommutingObservables(
		step, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}
	selectedIndices, selectedObservables, err := selectInteractionRelevantObservables(
		selectionStates, selectionTables, candidates, interactionRuntimeCount,
	)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}

	regressors, classifiers, static, err := trainAndEvaluateMultiplicityArm(
		"up52c_full_capacity",
		"up52c_memory_load",
		trainTables, heldTables,
		trainDepths, heldDepths,
		ops, mixer, selectedObservables, memoryNoise,
	)
	if err != nil {
		return UP52CMemoryLoadResult{}, err
	}

	result := UP52CMemoryLoadResult{
		Schema: UP52CMemoryLoadSchema,
		Experiment: "UP-52C-memory-load-interference",
		SourceUP51CSeal: "c4fb28b3852e916d419d77cc85a014b9f95946b9",
		MemoryNoise: memoryNoise,
		HeldDepths: append([]int(nil),heldDepths...),
		WriteLevels: append([]int(nil),writeLevels...),
		StaticHeldAccuracy: static.HeldOutAccuracy,
		SelectedIndices: append([]int(nil),selectedIndices...),
	}
	for _, writes := range writeLevels {
		integration, err := runUP52CIntegration(
			mixer, selectedObservables, ops,
			regressors, classifiers,
			heldTables, heldDepths, memoryNoise, writes,
		)
		if err != nil {
			return UP52CMemoryLoadResult{}, err
		}
		gate := static.HeldOutAccuracy >= 0.99 &&
			integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-10

		result.Points = append(result.Points, UP52CLoadPoint{
			WritesPerScenario:writes,
			Integration:integration,
			HardGate:gate,
		})
		if gate {
			result.LastPassingWrites = writes
		} else if result.FirstFailingWrites == 0 {
			result.FirstFailingWrites = writes
		}
	}
	return result,nil
}
