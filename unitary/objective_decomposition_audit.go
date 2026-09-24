package unitary

const ObjectiveDecompositionAuditSchema =
	"wingless.objective-decomposition-audit.v1"

var up45SelectedStep24Offsets = []float64{
	-0.05841298929708432,
	-0.03887820934987345,
	0.0051658188506906524,
	-0.03702321187042985,
	0.07759973733795059,
	0.05154885432874637,
}

var up45ExactFusionStep29Offsets = []float64{
	-0.045664698831175514,
	-0.04501676812740249,
	0.007218659931181362,
	-0.045664698831175514,
	0.07894943527059196,
	0.05017807058798021,
}

type ObjectiveComponentDelta struct {
	PhaseScore        float64 `json:"phase_score"`
	ValueProbability  float64 `json:"value_probability"`
	RelationProbability float64 `json:"relation_probability"`
	HarmonicTaskScore float64 `json:"harmonic_task_score"`
	SoftCapacity      float64 `json:"soft_capacity"`
	ResourcePenalty  float64 `json:"resource_penalty"`
	Objective         float64 `json:"objective"`
}

type HardCapabilityDelta struct {
	HeldOutAccuracy      float64 `json:"heldout_accuracy"`
	CommitAccuracy       float64 `json:"commit_accuracy"`
	FinalTableAccuracy   float64 `json:"final_table_accuracy"`
	RelationAccuracy     float64 `json:"relation_accuracy"`
}

type ObjectiveDecompositionAuditResult struct {
	Schema                    string               `json:"schema"`
	Experiment                string               `json:"experiment"`
	SourceUP45Seal            string               `json:"source_up45_seal"`
	OptimizerRun              bool                 `json:"optimizer_run"`
	LambdaChanged             bool                 `json:"lambda_changed"`
	SelectionUsesHeldOutData  bool                 `json:"selection_uses_heldout_data"`
	StatePairFrozenFromUP45   bool                 `json:"state_pair_frozen_from_up45"`
	SelectedStep              int                  `json:"selected_step"`
	ExactFusionStep           int                  `json:"exact_fusion_step"`
	SelectedExactCapacity     int                  `json:"selected_exact_capacity"`
	FusedExactCapacity        int                  `json:"fused_exact_capacity"`
	SelectedSmooth            RichDirectEvaluation `json:"selected_smooth"`
	FusedSmooth               RichDirectEvaluation `json:"fused_smooth"`
	SmoothDeltaFusedMinusSelected ObjectiveComponentDelta `json:"smooth_delta_fused_minus_selected"`
	SelectedHard              MultiplicityDoseArm  `json:"selected_hard"`
	FusedHard                 MultiplicityDoseArm  `json:"fused_hard"`
	HardDeltaFusedMinusSelected HardCapabilityDelta `json:"hard_delta_fused_minus_selected"`
}

func RunUP46ObjectiveDecompositionAudit() (
	ObjectiveDecompositionAuditResult,
	error,
) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables :=
		splitTaskAllocationTrainingPool(allTrain)

	selectedSmooth, err := evaluateSqrtFreeRunningDirectOffsets(
		"up46_selected_step24",
		append([]float64(nil), up45SelectedStep24Offsets...),
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}

	fusedSmooth, err := evaluateSqrtFreeRunningDirectOffsets(
		"up46_exact_fusion_step29",
		append([]float64(nil), up45ExactFusionStep29Offsets...),
		mixer,
		fitTables, validationTables,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}

	selectedHard, err := evaluateMultiplicityDoseArm(
		"up46_selected_step24_hard",
		append([]float64(nil), up45SelectedStep24Offsets...),
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}

	fusedHard, err := evaluateMultiplicityDoseArm(
		"up46_exact_fusion_step29_hard",
		append([]float64(nil), up45ExactFusionStep29Offsets...),
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}

	selectedGroups, selectedCapacity, err :=
		exactOffsetGroups(up45SelectedStep24Offsets)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}
	_ = selectedGroups
	fusedGroups, fusedCapacity, err :=
		exactOffsetGroups(up45ExactFusionStep29Offsets)
	if err != nil {
		return ObjectiveDecompositionAuditResult{}, err
	}
	_ = fusedGroups

	return ObjectiveDecompositionAuditResult{
		Schema:                   ObjectiveDecompositionAuditSchema,
		Experiment:               "UP-46-UP45-objective-decomposition",
		SourceUP45Seal:           "48fb886f517db0cb21fdbc757cd7bc9500ed9f59",
		OptimizerRun:             false,
		LambdaChanged:            false,
		SelectionUsesHeldOutData: false,
		StatePairFrozenFromUP45:  true,
		SelectedStep:             24,
		ExactFusionStep:          29,
		SelectedExactCapacity:    selectedCapacity,
		FusedExactCapacity:       fusedCapacity,
		SelectedSmooth:           selectedSmooth,
		FusedSmooth:              fusedSmooth,
		SmoothDeltaFusedMinusSelected: ObjectiveComponentDelta{
			PhaseScore:
				fusedSmooth.Signals.NormalizedPhaseScore -
					selectedSmooth.Signals.NormalizedPhaseScore,
			ValueProbability:
				fusedSmooth.Signals.MeanCorrectValueProbability -
					selectedSmooth.Signals.MeanCorrectValueProbability,
			RelationProbability:
				fusedSmooth.Signals.MeanCorrectRelationProbability -
					selectedSmooth.Signals.MeanCorrectRelationProbability,
			HarmonicTaskScore:
				fusedSmooth.Signals.HarmonicTaskScore -
					selectedSmooth.Signals.HarmonicTaskScore,
			SoftCapacity:
				fusedSmooth.SoftCapacity - selectedSmooth.SoftCapacity,
			ResourcePenalty:
				fusedSmooth.ResourcePenalty - selectedSmooth.ResourcePenalty,
			Objective:
				fusedSmooth.Objective - selectedSmooth.Objective,
		},
		SelectedHard: selectedHard,
		FusedHard:    fusedHard,
		HardDeltaFusedMinusSelected: HardCapabilityDelta{
			HeldOutAccuracy:
				fusedHard.Static.HeldOutAccuracy -
					selectedHard.Static.HeldOutAccuracy,
			CommitAccuracy:
				fusedHard.Integration.CommitDecodeAccuracy -
					selectedHard.Integration.CommitDecodeAccuracy,
			FinalTableAccuracy:
				fusedHard.Integration.ExactFinalTableAccuracy -
					selectedHard.Integration.ExactFinalTableAccuracy,
			RelationAccuracy:
				fusedHard.Integration.RelationalQueryAccuracy -
					selectedHard.Integration.RelationalQueryAccuracy,
		},
	}, nil
}
