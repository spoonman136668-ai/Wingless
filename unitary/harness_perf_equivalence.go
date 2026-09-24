package unitary

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"
)

const HarnessPerfEquivalenceSchema = "wingless.harness-perf-equivalence.v1"

const harnessPerfMinimumSpeedup = 1.05

type up44PerfContext struct {
	selectionStates       []State
	selectionTables       []memoryTable
	relationProbabilities func([]float64) ([]float64, error)
}

type HarnessPerfEquivalenceResult struct {
	Schema                   string  `json:"schema"`
	Experiment               string  `json:"experiment"`
	ReferenceUntouched       bool    `json:"reference_untouched"`
	OnlyInvariantReuse       bool    `json:"only_invariant_reuse"`
	ReferenceSHA256          string  `json:"reference_sha256"`
	OptimizedSHA256          string  `json:"optimized_sha256"`
	SecondOptimizedSHA256    string  `json:"second_optimized_sha256"`
	ExactEquivalence         bool    `json:"exact_equivalence"`
	OptimizedDeterministic   bool    `json:"optimized_deterministic"`
	ReferenceDurationMillis  int64   `json:"reference_duration_millis"`
	OptimizedDurationMillis  int64   `json:"optimized_duration_millis"`
	SecondOptimizedMillis    int64   `json:"second_optimized_duration_millis"`
	Speedup                  float64 `json:"speedup"`
	MinimumSpeedup           float64 `json:"minimum_speedup"`
	SpeedupGatePass          bool    `json:"speedup_gate_pass"`
	Accepted                 bool    `json:"accepted"`
	SelectedCheckpointStep   int     `json:"selected_checkpoint_step"`
	SelectedObjective        float64 `json:"selected_objective"`
	SelectedObjectiveGain    float64 `json:"selected_objective_gain"`
	SelectedCapacity         int     `json:"selected_capacity"`
	SelectedNearestGap       float64 `json:"selected_nearest_gap"`
	SelectedHeldOutAccuracy  float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy   float64 `json:"selected_commit_accuracy"`
	SelectedFinalAccuracy    float64 `json:"selected_final_accuracy"`
	SelectedRelationAccuracy float64 `json:"selected_relation_accuracy"`
}

func buildUP44PerfContext(
	fitTables []memoryTable,
	mixer latentMatrix,
	memoryNoise float64,
) (up44PerfContext, error) {
	selectionStates, selectionTables, err :=
		taskSelectedTrainingStates(
			fitTables, mixer, memoryNoise, 2,
		)
	if err != nil {
		return up44PerfContext{}, err
	}

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return up44PerfContext{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return up44PerfContext{}, err
	}

	return up44PerfContext{
		selectionStates: append([]State(nil), selectionStates...),
		selectionTables: append([]memoryTable(nil), selectionTables...),
		relationProbabilities: relationHead.probabilities,
	}, nil
}

func sqrtFreeRunningProbabilisticSignalsCached(
	ctx up44PerfContext,
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	tables []memoryTable,
	depths []int,
	memoryNoise float64,
) (float64, float64, error) {
	if len(tables) == 0 || len(depths) == 0 {
		return 0, 0, fmt.Errorf(
			"sqrt free-running validation set empty",
		)
	}

	scenarios := richDirectScenarioCount
	if len(tables) < scenarios {
		scenarios = len(tables)
	}

	var valueProbabilityTotal float64
	var valueProbabilityCount int
	var relationProbabilityTotal float64
	var relationProbabilityCount int

	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeFullRankScenario(
			scenarioIndex,
			tables[(scenarioIndex*7)%len(tables)],
			richDirectWrites,
			depths,
		)
		trueTable := scenario.initial
		softMemory, err := encodeMemory(scenario.initial)
		if err != nil {
			return 0, 0, err
		}

		for writeIndex, write := range scenario.writes {
			seed := 44000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(
				softMemory, seed, memoryNoise,
			)
			if err != nil {
				return 0, 0, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return 0, 0, err
			}
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return 0, 0, fmt.Errorf(
					"missing sqrt free-running mutable depth=%d",
					write.gap,
				)
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil {
				return 0, 0, err
			}
			_, distributions, _, err := decodeBreadthTable(
				state, observables, regressors, classifiers,
			)
			if err != nil {
				return 0, 0, err
			}
			for currentEntity := 0; currentEntity < 4; currentEntity++ {
				valueProbabilityTotal +=
					distributions[currentEntity][trueTable[currentEntity]]
				valueProbabilityCount++
			}

			softMemory, err = sqrtSoftMemoryAfterCommandedWrite(
				distributions, write.entity, write.value,
			)
			if err != nil {
				return 0, 0, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return 0, 0, err
			}
		}

		seed := 44000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(
			softMemory, seed, memoryNoise,
		)
		if err != nil {
			return 0, 0, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1), 2*math.Pi),
		)
		state, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return 0, 0, err
		}
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return 0, 0, fmt.Errorf(
				"missing sqrt free-running final depth=%d",
				scenario.finalGap,
			)
		}
		state, err = latentMatrixVector(operator, state)
		if err != nil {
			return 0, 0, err
		}
		_, distributions, _, err := decodeBreadthTable(
			state, observables, regressors, classifiers,
		)
		if err != nil {
			return 0, 0, err
		}
		for currentEntity := 0; currentEntity < 4; currentEntity++ {
			valueProbabilityTotal +=
				distributions[currentEntity][trueTable[currentEntity]]
			valueProbabilityCount++
		}

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil {
			return 0, 0, err
		}
		relationProbabilities, err :=
			ctx.relationProbabilities(relationInput)
		if err != nil {
			return 0, 0, err
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return 0, 0, err
		}
		relationProbabilityTotal +=
			relationProbabilities[wantRelation]
		relationProbabilityCount++
	}

	if valueProbabilityCount == 0 || relationProbabilityCount == 0 {
		return 0, 0, fmt.Errorf(
			"sqrt free-running objective produced no samples",
		)
	}
	return valueProbabilityTotal / float64(valueProbabilityCount),
		relationProbabilityTotal / float64(relationProbabilityCount),
		nil
}


func evaluateSqrtFreeRunningDirectOffsetsCached(
	ctx up44PerfContext,
	name string,
	offsets []float64,
	mixer latentMatrix,
	fitTables, validationTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (RichDirectEvaluation, error) {
	step, err := fullLatentDoseStep(mixer, offsets)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	realStep, err := realifyMatrix(step)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	orthError, err := maxRealOrthogonalityError(realStep)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	ops, err := latentDepthOperators(step, allDepths)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	candidates, _, err := discoverCommutingObservables(
		step, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	selectionStates := ctx.selectionStates
	selectionTables := ctx.selectionTables
	selectedIndices, selectedObservables, err :=
		selectInteractionRelevantObservables(
			selectionStates, selectionTables,
			candidates, interactionRuntimeCount,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(
		selectedObservables, step,
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, selectedObservables, ops[1024],
	)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	regressors, classifiers, static, err :=
		trainAndEvaluateMultiplicityArm(
			name,
			"sqrt_free_running_"+name,
			fitTables, validationTables,
			trainDepths, heldDepths,
			ops, mixer, selectedObservables, memoryNoise,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	valueProb, relationProb, err :=
		sqrtFreeRunningProbabilisticSignalsCached(
			ctx,
			mixer, selectedObservables, ops,
			regressors, classifiers,
			validationTables, heldDepths, memoryNoise,
		)
	if err != nil {
		return RichDirectEvaluation{}, err
	}

	phaseScore := 0.5 * (static.MeanPhaseCosine + 1.0)
	if phaseScore < 0 {
		phaseScore = 0
	}
	if phaseScore > 1 {
		phaseScore = 1
	}
	taskScore := harmonicMean3(
		phaseScore, valueProb, relationProb,
	)
	softCapacity, err := directSoftCapacity(offsets)
	if err != nil {
		return RichDirectEvaluation{}, err
	}
	resourcePenalty :=
		directOffsetResourcePrice * softCapacity / 36.0

	return RichDirectEvaluation{
		Name:            name,
		Offsets:         append([]float64(nil), offsets...),
		SoftCapacity:    softCapacity,
		ResourcePenalty: resourcePenalty,
		Signals: RichSmoothSignals{
			NormalizedPhaseScore:           phaseScore,
			MeanCorrectValueProbability:    valueProb,
			MeanCorrectRelationProbability: relationProb,
			HarmonicTaskScore:              taskScore,
		},
		Objective:                 taskScore - resourcePenalty,
		OrthogonalityError:       orthError,
		DiscoveryCommutatorError: commutatorError,
		FeatureDrift:             featureDrift,
		Static:                   static,
		SelectedIndices:          append([]int(nil), selectedIndices...),
	}, nil
}


func runUP44Cached() (SqrtProbabilityTaskBridgeProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables :=
		splitTaskAllocationTrainingPool(allTrain)

	ctx, err := buildUP44PerfContext(
		fitTables, mixer, memoryNoise,
	)
	if err != nil {
		return SqrtProbabilityTaskBridgeProbeResult{}, err
	}

	minFit := minimumTableMarginalCount(fitTables)
	minValidation := minimumTableMarginalCount(validationTables)
	innerBalanced := minFit >= 12 && minValidation >= 4

	currentOffsets := continuousFusionInitialOffsets()
	initial, err := evaluateSqrtFreeRunningDirectOffsetsCached(
		ctx,
		"sqrt_task_initial",
		currentOffsets,
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return SqrtProbabilityTaskBridgeProbeResult{}, err
	}
	if !richEvaluationStructurallyValid(initial) {
		return SqrtProbabilityTaskBridgeProbeResult{}, fmt.Errorf(
			"sqrt task initial evaluation structurally invalid",
		)
	}

	initialGroups, initialCapacity, initialGap, err :=
		reversibleProvisionalGeometry(currentOffsets)
	if err != nil {
		return SqrtProbabilityTaskBridgeProbeResult{}, err
	}
	selected := RollingTaskCheckpoint{
		Step:                0,
		Offsets:             append([]float64(nil), currentOffsets...),
		ProvisionalGroups:   initialGroups,
		ProvisionalCapacity: initialCapacity,
		NearestOffsetGap:    initialGap,
		Evaluation:          initial,
	}

	observations := make(
		[]rollingDirectionalObservation, 0,
		rollingGradientBufferSize,
	)
	steps := make([]RollingTaskStep, 0, rollingTaskMaximumSteps)
	checkpoints := make(
		[]RollingTaskCheckpoint, 0,
		len(rollingTaskCheckpointSteps),
	)

	fullRankSteps := 0
	firstProvisionalFusionStep := -1
	maximumProvisionalCapacity := initialCapacity

	for step := 1; step <= rollingTaskMaximumSteps; step++ {
		union := newFusionUnion()
		direction := deterministicSPSADirection(
			step, union.groups(),
		)
		plusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			directOffsetPerturbation, &union,
		)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		minusOffsets, err := perturbDirectOffsets(
			currentOffsets, direction,
			-directOffsetPerturbation, &union,
		)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}

		plus, err := evaluateSqrtFreeRunningDirectOffsetsCached(
		ctx,
			fmt.Sprintf("sqrt_task_plus_%02d", step),
			plusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		minus, err := evaluateSqrtFreeRunningDirectOffsetsCached(
		ctx,
			fmt.Sprintf("sqrt_task_minus_%02d", step),
			minusOffsets,
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		if !richEvaluationStructurallyValid(plus) ||
			!richEvaluationStructurallyValid(minus) {
			return SqrtProbabilityTaskBridgeProbeResult{}, fmt.Errorf(
				"sqrt task perturbation structurally invalid step=%d",
				step,
			)
		}

		derivative :=
			(plus.Objective - minus.Objective) /
				(2 * directOffsetPerturbation)
		rankOne := directGradientEstimate(
			plus.Objective, minus.Objective, direction,
		)
		observations = append(
			observations,
			rollingDirectionalObservation{
				direction: append([]float64(nil), direction...),
				derivative: derivative,
			},
		)
		if len(observations) > rollingGradientBufferSize {
			observations = observations[
				len(observations)-rollingGradientBufferSize:
			]
		}

		gradient := rankOne
		reconstructed, fullRank, err :=
			rollingFullRankGradient(observations)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		if fullRank {
			gradient = reconstructed
			fullRankSteps++
		}

		nextOffsets, update, err :=
			applyCalibrationGradientNoSticky(
				currentOffsets, gradient,
			)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		groups, capacity, nearest, err :=
			reversibleProvisionalGeometry(nextOffsets)
		if err != nil {
			return SqrtProbabilityTaskBridgeProbeResult{}, err
		}
		if capacity > 6 && firstProvisionalFusionStep < 0 {
			firstProvisionalFusionStep = step
		}
		if capacity > maximumProvisionalCapacity {
			maximumProvisionalCapacity = capacity
		}

		steps = append(steps, RollingTaskStep{
			Step:                   step,
			Direction:              append([]float64(nil), direction...),
			DirectionalDerivative:  derivative,
			FullRankReconstruction: fullRank,
			UpdateMaxAbs:           maxAbsFloat64(update),
			OutputOffsets:          append([]float64(nil), nextOffsets...),
			ProvisionalGroups:      groups,
			ProvisionalCapacity:    capacity,
			NearestOffsetGap:       nearest,
		})
		currentOffsets = nextOffsets

		if rollingTaskCheckpointStep(step) {
			evaluation, err :=
				evaluateSqrtFreeRunningDirectOffsetsCached(
		ctx,
					fmt.Sprintf(
						"sqrt_task_checkpoint_%02d",
						step,
					),
					currentOffsets,
					mixer,
					fitTables, validationTables,
					trainDepths, heldDepths, allDepths,
					memoryNoise,
				)
			if err != nil {
				return SqrtProbabilityTaskBridgeProbeResult{}, err
			}
			if !richEvaluationStructurallyValid(evaluation) {
				return SqrtProbabilityTaskBridgeProbeResult{}, fmt.Errorf(
					"sqrt task checkpoint structurally invalid step=%d",
					step,
				)
			}
			checkpoint := RollingTaskCheckpoint{
				Step:                step,
				Offsets:             append([]float64(nil), currentOffsets...),
				ProvisionalGroups:   groups,
				ProvisionalCapacity: capacity,
				NearestOffsetGap:    nearest,
				Evaluation:          evaluation,
			}
			checkpoints = append(checkpoints, checkpoint)
			if evaluation.Objective > selected.Evaluation.Objective ||
				(evaluation.Objective == selected.Evaluation.Objective &&
					evaluation.SoftCapacity <
						selected.Evaluation.SoftCapacity) {
				selected = checkpoint
			}
		}
	}

	selectedFinal, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf(
			"selected_sqrt_task_checkpoint_%02d",
			selected.Step,
		),
		selected.Offsets,
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return SqrtProbabilityTaskBridgeProbeResult{}, err
	}
	fullControl, err := evaluateMultiplicityDoseArm(
		"full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return SqrtProbabilityTaskBridgeProbeResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedFinal.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedFinal.Integration.CommitDecodeAccuracy
	capabilityWithoutExactFusion :=
		selectedFinal.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedFinal.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedFinal.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedFinal.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	objectiveGain :=
		selected.Evaluation.Objective - initial.Objective

	return SqrtProbabilityTaskBridgeProbeResult{
		Schema:                           SqrtProbabilityTaskBridgeSchema,
		Experiment:                       "UP-44-sqrt-probability-soft-memory-task-bridge",
		KnownTargetUsed:                  false,
		RollingFullRankDirectionalMemory: true,
		RollingBufferSize:                rollingGradientBufferSize,
		StickyProjectionUsed:             false,
		TeacherForcedSmoothMutable:       false,
		FreeRunningProbabilisticMutable:  true,
		SoftMemoryAmplitudeEncoding:      "0.5*sqrt(probability)",
		GlobalSoftMemoryRenormalization:  false,
		SoftWriteBoundary:                sqrtProbabilitySoftWriteBoundary,
		SelectionUsesHeldOutData:         false,
		MaximumSteps:                     rollingTaskMaximumSteps,
		CheckpointSteps:                  append([]int(nil), rollingTaskCheckpointSteps...),
		Perturbation:                     directOffsetPerturbation,
		LearningRate:                     directOffsetLearningRate,
		MaximumCoordinateUpdate:          directOffsetMaxUpdate,
		ExactFusionTolerance:             continuousFusionTolerance,
		ResourcePrice:                    directOffsetResourcePrice,
		SoftCapacityTau:                  directOffsetSoftCapacityTau,
		FitTables:                        len(fitTables),
		ValidationTables:                 len(validationTables),
		TrueHeldOutTables:                len(trueHeld),
		Initial:                          initial,
		Steps:                            steps,
		Checkpoints:                      checkpoints,
		SelectedCheckpoint:               selected,
		SelectedFinalArm:                 selectedFinal,
		FullCapacityControl:              fullControl,
		Diagnosis: RollingTaskBridgeDiagnosis{
			InnerSplitBalancedPass:            innerBalanced,
			SelectionUsesTrueHeldOut:          false,
			FullRankReconstructionSteps:       fullRankSteps,
			SelectedCheckpointStep:            selected.Step,
			InitialObjective:                  initial.Objective,
			SelectedObjective:                 selected.Evaluation.Objective,
			SelectedObjectiveGain:             objectiveGain,
			ObjectiveImproved:                 objectiveGain > 0,
			FirstProvisionalFusionStep:        firstProvisionalFusionStep,
			MaximumProvisionalCapacity:        maximumProvisionalCapacity,
			SelectedProvisionalCapacity:       selected.ProvisionalCapacity,
			SelectedNearestOffsetGap:          selected.NearestOffsetGap,
			TaskGradientEnteredFusionBasin:    maximumProvisionalCapacity > 6,
			SelectedHeldOutAccuracy:           selectedFinal.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:            selectedFinal.Integration.CommitDecodeAccuracy,
			SelectedFinalAccuracy:             selectedFinal.Integration.ExactFinalTableAccuracy,
			SelectedRelationAccuracy:          selectedFinal.Integration.RelationalQueryAccuracy,
			HeldOutRetentionDelta:             heldDelta,
			CommitRetentionDelta:              commitDelta,
			CapabilityGatesWithoutExactFusion: capabilityWithoutExactFusion,
		},
	}, nil
}


func hashUP44Result(
	result SqrtProbabilityTaskBridgeProbeResult,
) (string, error) {
	payload, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func RunHarnessPerfEquivalence() (HarnessPerfEquivalenceResult, error) {
	referenceStart := time.Now()
	reference, err := RunUP44()
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}
	referenceDuration := time.Since(referenceStart)

	optimizedStart := time.Now()
	optimized, err := runUP44Cached()
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}
	optimizedDuration := time.Since(optimizedStart)

	secondStart := time.Now()
	secondOptimized, err := runUP44Cached()
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}
	secondDuration := time.Since(secondStart)

	exact := reflect.DeepEqual(reference, optimized)
	deterministic := reflect.DeepEqual(optimized, secondOptimized)

	referenceHash, err := hashUP44Result(reference)
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}
	optimizedHash, err := hashUP44Result(optimized)
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}
	secondHash, err := hashUP44Result(secondOptimized)
	if err != nil {
		return HarnessPerfEquivalenceResult{}, err
	}

	speedup := 0.0
	if optimizedDuration > 0 {
		speedup =
			float64(referenceDuration) /
				float64(optimizedDuration)
	}
	speedPass := speedup >= harnessPerfMinimumSpeedup
	accepted := exact && deterministic && speedPass

	return HarnessPerfEquivalenceResult{
		Schema:                  HarnessPerfEquivalenceSchema,
		Experiment:              "HARNESS-PERF-R1-UP44-INVARIANT-REUSE",
		ReferenceUntouched:      true,
		OnlyInvariantReuse:      true,
		ReferenceSHA256:         referenceHash,
		OptimizedSHA256:         optimizedHash,
		SecondOptimizedSHA256:   secondHash,
		ExactEquivalence:        exact,
		OptimizedDeterministic:  deterministic,
		ReferenceDurationMillis: referenceDuration.Milliseconds(),
		OptimizedDurationMillis: optimizedDuration.Milliseconds(),
		SecondOptimizedMillis:   secondDuration.Milliseconds(),
		Speedup:                 speedup,
		MinimumSpeedup:          harnessPerfMinimumSpeedup,
		SpeedupGatePass:         speedPass,
		Accepted:                accepted,
		SelectedCheckpointStep:
			optimized.Diagnosis.SelectedCheckpointStep,
		SelectedObjective:
			optimized.Diagnosis.SelectedObjective,
		SelectedObjectiveGain:
			optimized.Diagnosis.SelectedObjectiveGain,
		SelectedCapacity:
			optimized.Diagnosis.SelectedProvisionalCapacity,
		SelectedNearestGap:
			optimized.Diagnosis.SelectedNearestOffsetGap,
		SelectedHeldOutAccuracy:
			optimized.Diagnosis.SelectedHeldOutAccuracy,
		SelectedCommitAccuracy:
			optimized.Diagnosis.SelectedCommitAccuracy,
		SelectedFinalAccuracy:
			optimized.Diagnosis.SelectedFinalAccuracy,
		SelectedRelationAccuracy:
			optimized.Diagnosis.SelectedRelationAccuracy,
	}, nil
}
