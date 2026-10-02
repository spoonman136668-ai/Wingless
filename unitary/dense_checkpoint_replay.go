package unitary

import "fmt"

const DenseCheckpointReplaySchema =
	"wingless.dense-checkpoint-replay.v1"

type frozenUP45ReplayState struct {
	Step    int
	Offsets []float64
}

var up45FrozenDenseTrajectory = []frozenUP45ReplayState{
	{Step: 0, Offsets: []float64{-0.07319250547113999, -0.043915503282683996, -0.014638501094227999, 0.014638501094227999, 0.043915503282683996, 0.07319250547113999}},
	{Step: 1, Offsets: []float64{-0.07316644270252409, -0.04397276358297841, -0.014596839559772846, 0.014596839559772846, 0.04397276358297841, 0.07316644270252409}},
	{Step: 2, Offsets: []float64{-0.07297525927302635, -0.043782230292719305, -0.014976624090971526, 0.014216404889335541, 0.044162339893537295, 0.07335536887384435}},
	{Step: 3, Offsets: []float64{-0.07300814019784398, -0.04373947625419586, -0.014937360558193459, 0.014175116381827644, 0.04412059834237876, 0.07338926228602688}},
	{Step: 4, Offsets: []float64{-0.0729131972920033, -0.0437352140652298, -0.015023612327382065, 0.013998239721812272, 0.044247900368014695, 0.07342588359478819}},
	{Step: 5, Offsets: []float64{-0.07589790066159101, -0.04024016437472746, -0.015241498534647551, 0.013457334392994673, 0.04769149497338544, 0.07023073420458589}},
	{Step: 6, Offsets: []float64{-0.07854727091820522, -0.03630079972357977, -0.015883044325533387, 0.012638951613074505, 0.05105594860678387, 0.06703621474746}},
	{Step: 7, Offsets: []float64{-0.08034033672506748, -0.03286901312400459, -0.017204938231679616, 0.012467270824376403, 0.05359456403121099, 0.06435245322516428}},
	{Step: 8, Offsets: []float64{-0.07964934726360334, -0.03400477621970551, -0.014906337735706557, 0.008556236628617934, 0.0586343358108007, 0.061369888779596794}},
	{Step: 9, Offsets: []float64{-0.07910105538896371, -0.03565909754277756, -0.01146971427175334, 0.005427182742503438, 0.06302801829538054, 0.05777466616561062}},
	{Step: 10, Offsets: []float64{-0.0783070560381359, -0.03711013842772236, -0.007880348916601426, 0.00217228489808608, 0.06720942389587246, 0.053915834588501146}},
	{Step: 11, Offsets: []float64{-0.07697739916566095, -0.03844017923842922, -0.004471782900601095, -0.001358222209758832, 0.07121866726097625, 0.05002891625347383}},
	{Step: 12, Offsets: []float64{-0.07313926995523916, -0.04265171080118584, -0.0008017537650376425, -0.00561286579673649, 0.07399130880108405, 0.04821429151711507}},
	{Step: 13, Offsets: []float64{-0.06889138001302549, -0.046621766989668134, 0.002762441104174088, -0.009889989648494907, 0.07639641364681934, 0.04624428190019509}},
	{Step: 14, Offsets: []float64{-0.064823219625267, -0.04987262590487743, 0.006539830797834645, -0.013740391692078968, 0.07897230766235513, 0.04292409876203362}},
	{Step: 15, Offsets: []float64{-0.06482247736387664, -0.04949426386551825, 0.010272027524973654, -0.017391909341715805, 0.07755494575279342, 0.043881677293343614}},
	{Step: 16, Offsets: []float64{-0.06436329123108397, -0.04868741582335552, 0.012714425473410538, -0.021202996561164242, 0.07641860929906982, 0.04512066884312338}},
	{Step: 17, Offsets: []float64{-0.063657966485046, -0.047935143702770126, 0.014971729882445228, -0.024875601196743348, 0.07534607297464248, 0.04615090852747177}},
	{Step: 18, Offsets: []float64{-0.06270106936015475, -0.04721227967145074, 0.016026681032886448, -0.028085543130790803, 0.07431491210408045, 0.04765729902542938}},
	{Step: 19, Offsets: []float64{-0.06340586648323157, -0.04379903543552296, 0.015350709470188616, -0.030641664169700728, 0.07510735257170795, 0.0473885040465587}},
	{Step: 20, Offsets: []float64{-0.06401425341725661, -0.040492652447057326, 0.014707266295163744, -0.03303393367291187, 0.07575088310011475, 0.047082690141947305}},
	{Step: 21, Offsets: []float64{-0.06331008895069434, -0.039576128112466295, 0.012193769635826416, -0.03391955031399446, 0.0757645243973178, 0.04884747334401088}},
	{Step: 22, Offsets: []float64{-0.06201665153290661, -0.03908025745592321, 0.010220945615943012, -0.035171756955477744, 0.07608399594480736, 0.04996372438355718}},
	{Step: 23, Offsets: []float64{-0.06058048802159582, -0.038438972259469094, 0.00807705892065937, -0.03656765889886631, 0.07641618313573768, 0.05109387712353417}},
	{Step: 24, Offsets: []float64{-0.05841298929708432, -0.03887820934987345, 0.0051658188506906524, -0.03702321187042985, 0.07759973733795059, 0.05154885432874637}},
	{Step: 25, Offsets: []float64{-0.057287138597534884, -0.03928995671546388, 0.004765552979263114, -0.03734055059955403, 0.07902118896859409, 0.05013090396469559}},
	{Step: 26, Offsets: []float64{-0.05473205208659455, -0.04109552588777735, 0.005748753575759921, -0.03905009430893662, 0.07906382778649103, 0.050065090921057556}},
	{Step: 27, Offsets: []float64{-0.051904783729462715, -0.04206384824927616, 0.006355402130227749, -0.04154178965560916, 0.07923877062684413, 0.04991624887727617}},
	{Step: 28, Offsets: []float64{-0.04892314510252909, -0.04348107195043551, 0.006635381476917126, -0.04354107846428387, 0.07906921245674402, 0.05024070158358733}},
	{Step: 29, Offsets: []float64{-0.045664698831175514, -0.04501676812740249, 0.007218659931181362, -0.045664698831175514, 0.07894943527059196, 0.05017807058798021}},
	{Step: 30, Offsets: []float64{-0.04273134838009583, -0.047200255976499926, 0.008758072772756354, -0.04713778924438781, 0.078570846182747, 0.049740474645480225}},
	{Step: 31, Offsets: []float64{-0.041188603831002535, -0.047149435725171066, 0.012280174465935613, -0.05076268104998453, 0.07588127325658892, 0.05093927288363361}},
	{Step: 32, Offsets: []float64{-0.03954045057871559, -0.046911738206264465, 0.015110809330505721, -0.05416925161574409, 0.07267342905140177, 0.05283720201881665}},
	{Step: 33, Offsets: []float64{-0.039324496292631426, -0.045510301919876976, 0.01670259395503804, -0.05647280908206149, 0.0709926426115726, 0.053612370727959255}},
	{Step: 34, Offsets: []float64{-0.04035851600548692, -0.0452606908150856, 0.018957161162040034, -0.0568353162336186, 0.06907103456074841, 0.05442632733140268}},
	{Step: 35, Offsets: []float64{-0.041344381126643445, -0.045000278866654135, 0.021167200756903155, -0.057144680567527854, 0.06710662298915257, 0.055215516814769716}},
	{Step: 36, Offsets: []float64{-0.0424554016262346, -0.04461753293402405, 0.023233483297371453, -0.05732849915801507, 0.06503520881449358, 0.05613274160640869}},
	{Step: 37, Offsets: []float64{-0.04305753204830245, -0.04373508304830271, 0.0244752142950225, -0.0580617740617826, 0.06314871334032307, 0.05723046152304218}},
	{Step: 38, Offsets: []float64{-0.042767703668147986, -0.04421317811195958, 0.024721661849880075, -0.05798382268979346, 0.06273298608981424, 0.0575100565302067}},
	{Step: 39, Offsets: []float64{-0.042474459297132415, -0.04467438692167782, 0.02469891178729828, -0.057838308239667284, 0.062400518311379466, 0.05788772435979977}},
	{Step: 40, Offsets: []float64{-0.04196399772726684, -0.04526324941777318, 0.02482906952890605, -0.05781350546286627, 0.061831046595615016, 0.058380636483385207}},
	{Step: 41, Offsets: []float64{-0.04152345551224157, -0.0459059965003096, 0.024985874603958756, -0.057677693131135854, 0.06125898725132394, 0.05886228328840431}},
	{Step: 42, Offsets: []float64{-0.041204150839203126, -0.04634799444590295, 0.025271067413355233, -0.05764500456897697, 0.06079229983650971, 0.05913378260421811}},
	{Step: 43, Offsets: []float64{-0.04036910636432661, -0.047079530797945475, 0.025902255902208456, -0.05787982233631974, 0.06059153334394607, 0.058834670252437285}},
	{Step: 44, Offsets: []float64{-0.03904394676074399, -0.04727531084824827, 0.025659717665893388, -0.058748327947075366, 0.06062821778584114, 0.0587796501043331}},
	{Step: 45, Offsets: []float64{-0.03953580506334601, -0.044906442413640196, 0.027240274102503917, -0.06087115707805494, 0.05794817579417078, 0.060124954658366446}},
	{Step: 46, Offsets: []float64{-0.039925806761896436, -0.04251483188223937, 0.028875799605078144, -0.06293140587678792, 0.0552019907939392, 0.06129425412190637}},
	{Step: 47, Offsets: []float64{-0.04025404506589853, -0.04011771324879275, 0.030356604775231524, -0.06484592174255784, 0.05257612666087719, 0.062284948621140424}},
	{Step: 48, Offsets: []float64{-0.04050481968956773, -0.03772334675340676, 0.03174652633686573, -0.0666483907461437, 0.04995983458779618, 0.06317019626445626}},
}

type DenseReplayPoint struct {
	Step                int       `json:"step"`
	Offsets             []float64 `json:"offsets"`
	ExactCapacity       int       `json:"exact_capacity"`
	ExactFusion         bool      `json:"exact_fusion"`
	PhaseScore          float64   `json:"phase_score"`
	ValueProbability    float64   `json:"value_probability"`
	RelationProbability float64   `json:"relation_probability"`
	HarmonicTaskScore   float64   `json:"harmonic_task_score"`
	SoftCapacity        float64   `json:"soft_capacity"`
	ResourcePenalty     float64   `json:"resource_penalty"`
	Objective           float64   `json:"objective"`
}

type DenseCheckpointReplayDiagnosis struct {
	DenseSelectedStep              int     `json:"dense_selected_step"`
	DenseSelectedExactCapacity     int     `json:"dense_selected_exact_capacity"`
	DenseSelectedExactFusion       bool    `json:"dense_selected_exact_fusion"`
	DenseSelectedObjective         float64 `json:"dense_selected_objective"`
	SparseUP45SelectedStep         int     `json:"sparse_up45_selected_step"`
	SparseUP45SelectedObjective    float64 `json:"sparse_up45_selected_objective"`
	ExactFusionStep29Objective     float64 `json:"exact_fusion_step29_objective"`
	DenseMinusSparseObjective      float64 `json:"dense_minus_sparse_objective"`
	SelectedHeldOutAccuracy        float64 `json:"selected_heldout_accuracy"`
	SelectedCommitAccuracy         float64 `json:"selected_commit_accuracy"`
	SelectedFinalAccuracy          float64 `json:"selected_final_accuracy"`
	SelectedRelationAccuracy       float64 `json:"selected_relation_accuracy"`
	HeldOutRetentionDelta          float64 `json:"heldout_retention_delta"`
	CommitRetentionDelta           float64 `json:"commit_retention_delta"`
	CapabilityGates                bool    `json:"capability_gates"`
	DenseScheduleRecoveredBetter   bool    `json:"dense_schedule_recovered_better"`
	SelectedExactFusionAndImproved bool    `json:"selected_exact_fusion_and_improved"`
}

type DenseCheckpointReplayResult struct {
	Schema                   string                         `json:"schema"`
	Experiment               string                         `json:"experiment"`
	SourceUP45Seal           string                         `json:"source_up45_seal"`
	SourceUP46Seal           string                         `json:"source_up46_seal"`
	TrajectoryFrozenFromUP45 bool                           `json:"trajectory_frozen_from_up45"`
	OptimizerRun             bool                           `json:"optimizer_run"`
	LambdaChanged            bool                           `json:"lambda_changed"`
	SelectionUsesHeldOutData bool                           `json:"selection_uses_heldout_data"`
	ReferenceEvaluatorUsed   bool                           `json:"reference_evaluator_used"`
	CandidateCount           int                            `json:"candidate_count"`
	Points                   []DenseReplayPoint             `json:"points"`
	SelectedPoint            DenseReplayPoint               `json:"selected_point"`
	SelectedHard             MultiplicityDoseArm            `json:"selected_hard"`
	FullCapacityControl      MultiplicityDoseArm            `json:"full_capacity_control"`
	Diagnosis                DenseCheckpointReplayDiagnosis `json:"diagnosis"`
}

func denseReplayPoint(
	state frozenUP45ReplayState,
	evaluation RichDirectEvaluation,
) (DenseReplayPoint, error) {
	_, capacity, err := exactOffsetGroups(state.Offsets)
	if err != nil {
		return DenseReplayPoint{}, err
	}
	return DenseReplayPoint{
		Step:                state.Step,
		Offsets:             append([]float64(nil), state.Offsets...),
		ExactCapacity:       capacity,
		ExactFusion:         capacity > compositeChannels,
		PhaseScore:          evaluation.Signals.NormalizedPhaseScore,
		ValueProbability:    evaluation.Signals.MeanCorrectValueProbability,
		RelationProbability: evaluation.Signals.MeanCorrectRelationProbability,
		HarmonicTaskScore:   evaluation.Signals.HarmonicTaskScore,
		SoftCapacity:        evaluation.SoftCapacity,
		ResourcePenalty:     evaluation.ResourcePenalty,
		Objective:           evaluation.Objective,
	}, nil
}

func RunUP47() (DenseCheckpointReplayResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	if len(up45FrozenDenseTrajectory) != 49 {
		return DenseCheckpointReplayResult{}, fmt.Errorf(
			"frozen trajectory count=%d want=49",
			len(up45FrozenDenseTrajectory),
		)
	}

	mixer := fullLatentMixer()
	allTrain := fullObserverTablePool(true)
	trueHeld := fullObserverTablePool(false)
	fitTables, validationTables :=
		splitTaskAllocationTrainingPool(allTrain)

	points := make([]DenseReplayPoint, 0, len(up45FrozenDenseTrajectory))
	var selectedPoint DenseReplayPoint
	var selectedEvaluation RichDirectEvaluation
	selectedSet := false

	for index, state := range up45FrozenDenseTrajectory {
		if state.Step != index {
			return DenseCheckpointReplayResult{}, fmt.Errorf(
				"frozen trajectory step[%d]=%d",
				index, state.Step,
			)
		}
		evaluation, err := evaluateSqrtFreeRunningDirectOffsets(
			fmt.Sprintf("up47_dense_step_%02d", state.Step),
			append([]float64(nil), state.Offsets...),
			mixer,
			fitTables, validationTables,
			trainDepths, heldDepths, allDepths,
			memoryNoise,
		)
		if err != nil {
			return DenseCheckpointReplayResult{}, err
		}
		if !richEvaluationStructurallyValid(evaluation) {
			return DenseCheckpointReplayResult{}, fmt.Errorf(
				"dense replay evaluation invalid step=%d",
				state.Step,
			)
		}
		point, err := denseReplayPoint(state, evaluation)
		if err != nil {
			return DenseCheckpointReplayResult{}, err
		}
		points = append(points, point)

		if !selectedSet ||
			evaluation.Objective > selectedEvaluation.Objective ||
			(evaluation.Objective == selectedEvaluation.Objective &&
				evaluation.SoftCapacity < selectedEvaluation.SoftCapacity) {
			selectedSet = true
			selectedPoint = point
			selectedEvaluation = evaluation
		}
	}
	if !selectedSet {
		return DenseCheckpointReplayResult{}, fmt.Errorf(
			"dense replay selected no point",
		)
	}

	// True held-out data is touched only after the training-side dense
	// objective has frozen the selected replay state.
	selectedHard, err := evaluateMultiplicityDoseArm(
		fmt.Sprintf(
			"up47_dense_selected_step_%02d",
			selectedPoint.Step,
		),
		append([]float64(nil), selectedPoint.Offsets...),
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return DenseCheckpointReplayResult{}, err
	}
	fullControl, err := evaluateMultiplicityDoseArm(
		"up47_full_capacity_control",
		[]float64{0, 0, 0, 0, 0, 0},
		mixer,
		allTrain,
		trueHeld,
		trainDepths, heldDepths, allDepths,
		memoryNoise,
	)
	if err != nil {
		return DenseCheckpointReplayResult{}, err
	}

	heldDelta := fullControl.Static.HeldOutAccuracy -
		selectedHard.Static.HeldOutAccuracy
	commitDelta := fullControl.Integration.CommitDecodeAccuracy -
		selectedHard.Integration.CommitDecodeAccuracy
	capability :=
		selectedHard.Static.HeldOutAccuracy >= taskAllocationMinHeld &&
			selectedHard.Integration.CommitDecodeAccuracy >= taskAllocationMinCommit &&
			selectedHard.Integration.ExactFinalTableAccuracy >= 0.90 &&
			selectedHard.Integration.RelationalQueryAccuracy >= 0.95 &&
			heldDelta <= 0.02 &&
			commitDelta <= 0.05

	const sparseStep = 24
	const sparseObjective = 0.5581278352889602
	const fusedStep29Objective = 0.5591129650448196

	return DenseCheckpointReplayResult{
		Schema:                   DenseCheckpointReplaySchema,
		Experiment:               "UP-47-frozen-UP45-dense-checkpoint-replay",
		SourceUP45Seal:           "48fb886f517db0cb21fdbc757cd7bc9500ed9f59",
		SourceUP46Seal:           "3a4a5c209818a803b5451d2ca4adf74c62d16680",
		TrajectoryFrozenFromUP45: true,
		OptimizerRun:             false,
		LambdaChanged:            false,
		SelectionUsesHeldOutData: false,
		ReferenceEvaluatorUsed:   true,
		CandidateCount:           len(points),
		Points:                   points,
		SelectedPoint:            selectedPoint,
		SelectedHard:             selectedHard,
		FullCapacityControl:      fullControl,
		Diagnosis: DenseCheckpointReplayDiagnosis{
			DenseSelectedStep:              selectedPoint.Step,
			DenseSelectedExactCapacity:     selectedPoint.ExactCapacity,
			DenseSelectedExactFusion:       selectedPoint.ExactFusion,
			DenseSelectedObjective:         selectedPoint.Objective,
			SparseUP45SelectedStep:         sparseStep,
			SparseUP45SelectedObjective:    sparseObjective,
			ExactFusionStep29Objective:     fusedStep29Objective,
			DenseMinusSparseObjective:      selectedPoint.Objective - sparseObjective,
			SelectedHeldOutAccuracy:        selectedHard.Static.HeldOutAccuracy,
			SelectedCommitAccuracy:         selectedHard.Integration.CommitDecodeAccuracy,
			SelectedFinalAccuracy:          selectedHard.Integration.ExactFinalTableAccuracy,
			SelectedRelationAccuracy:       selectedHard.Integration.RelationalQueryAccuracy,
			HeldOutRetentionDelta:          heldDelta,
			CommitRetentionDelta:           commitDelta,
			CapabilityGates:                capability,
			DenseScheduleRecoveredBetter:   selectedPoint.Objective > sparseObjective,
			SelectedExactFusionAndImproved: selectedPoint.ExactFusion && selectedPoint.Objective > sparseObjective,
		},
	}, nil
}
