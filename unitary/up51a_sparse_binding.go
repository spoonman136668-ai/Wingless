package unitary

const UP51ASparseBindingSchema = "wingless.up51a-sparse-binding-generalization.v1"

type UP51ADiagnosis struct {
	TrainingBasisStates       int     `json:"training_basis_states"`
	HeldOutBasisStates        int     `json:"heldout_basis_states"`
	TrainingCoverageFraction  float64 `json:"training_coverage_fraction"`
	UnitaryTrainAccuracy      float64 `json:"unitary_train_accuracy"`
	UnseenPrimitiveAccuracy   float64 `json:"unseen_primitive_accuracy"`
	UnseenSwap12Accuracy      float64 `json:"unseen_swap12_accuracy"`
	DerivedRole1Accuracy      float64 `json:"derived_role1_accuracy"`
	DerivedRole2Accuracy      float64 `json:"derived_role2_accuracy"`
	LongProgramAccuracy       float64 `json:"long_program_accuracy"`
	AggregateHeldOutAccuracy  float64 `json:"aggregate_heldout_accuracy"`
	SparseGeneralizationGate  bool    `json:"sparse_generalization_gate"`
	NonUnitaryHeldOutAccuracy float64 `json:"nonunitary_heldout_accuracy"`
}

type UP51ASparseBindingResult struct {
	Schema                 string          `json:"schema"`
	Experiment             string          `json:"experiment"`
	SourceUP50ASeal        string          `json:"source_up50a_seal"`
	ValuesPerRole          int             `json:"values_per_role"`
	Roles                  int             `json:"roles"`
	JointDimension         int             `json:"joint_dimension"`
	TrainingBasisRule      string          `json:"training_basis_rule"`
	TrainSingleStepOnly    bool            `json:"train_single_step_only"`
	HeldOutLengths         []int           `json:"heldout_lengths"`
	Steps                  int             `json:"steps"`
	LearningRate           float64         `json:"learning_rate"`
	GradientEpsilon        float64         `json:"gradient_epsilon"`
	Unitary                UP50APathResult `json:"unitary"`
	NonUnitary             UP50APathResult `json:"nonunitary_matched"`
	Diagnosis              UP51ADiagnosis  `json:"diagnosis"`
}

func up51aTrainingBasis(roles [3]int) bool {
	return (roles[0]+roles[1]+roles[2])%3 == 0
}

func up51aSparseTrainSamples() ([]up50aSample, int, error) {
	var out []up50aSample
	basisCount := 0
	for _, roles := range up50aAllRoles() {
		if !up51aTrainingBasis(roles) {
			continue
		}
		basisCount++
		for start := 0; start < 3; start++ {
			program := []UP50AInstruction{{Kind: "value_swap", Start: start}}
			target, err := up50aTarget(roles, program)
			if err != nil {
				return nil, 0, err
			}
			out = append(out, up50aSample{
				roles: roles, program: program, target: target,
				category: "train_sparse_role0_value",
			})
		}
		program := []UP50AInstruction{{Kind: "role_swap", Position: 0}}
		target, err := up50aTarget(roles, program)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, up50aSample{
			roles: roles, program: program, target: target,
			category: "train_sparse_swap01",
		})
	}
	return out, basisCount, nil
}

func up51aSparseHeldSamples(lengths []int) ([]up50aSample, int, error) {
	var out []up50aSample
	heldBasis := 0

	for _, roles := range up50aAllRoles() {
		if up51aTrainingBasis(roles) {
			continue
		}
		heldBasis++

		for start := 0; start < 3; start++ {
			primitive := []UP50AInstruction{{Kind: "value_swap", Start: start}}
			target, err := up50aTarget(roles, primitive)
			if err != nil {
				return nil, 0, err
			}
			out = append(out, up50aSample{
				roles: roles, program: primitive, target: target,
				category: "unseen_basis_primitive",
			})
		}

		swap12 := []UP50AInstruction{{Kind: "role_swap", Position: 1}}
		target, err := up50aTarget(roles, swap12)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, up50aSample{
			roles: roles, program: swap12, target: target,
			category: "unseen_swap12",
		})

		for start := 0; start < 3; start++ {
			role1 := []UP50AInstruction{
				{Kind: "role_swap", Position: 0},
				{Kind: "value_swap", Start: start},
				{Kind: "role_swap", Position: 0},
			}
			target1, err := up50aTarget(roles, role1)
			if err != nil {
				return nil, 0, err
			}
			out = append(out, up50aSample{
				roles: roles, program: role1, target: target1,
				category: "derived_role1",
			})

			role2 := []UP50AInstruction{
				{Kind: "role_swap", Position: 1},
				{Kind: "role_swap", Position: 0},
				{Kind: "value_swap", Start: start},
				{Kind: "role_swap", Position: 0},
				{Kind: "role_swap", Position: 1},
			}
			target2, err := up50aTarget(roles, role2)
			if err != nil {
				return nil, 0, err
			}
			out = append(out, up50aSample{
				roles: roles, program: role2, target: target2,
				category: "derived_role2",
			})
		}
	}

	for _, length := range lengths {
		for _, roles := range up50aAllRoles() {
			if up51aTrainingBasis(roles) {
				continue
			}
			for seed := 0; seed < 3; seed++ {
				program := make([]UP50AInstruction, length)
				for i := range program {
					switch (i + seed + roles[0] + 2*roles[1] + roles[2]) % 4 {
					case 0:
						program[i] = UP50AInstruction{Kind: "role_swap", Position: 0}
					case 1:
						program[i] = UP50AInstruction{Kind: "role_swap", Position: 1}
					default:
						program[i] = UP50AInstruction{
							Kind: "value_swap",
							Start: (i*i + seed + roles[1]) % 3,
						}
					}
				}
				target, err := up50aTarget(roles, program)
				if err != nil {
					return nil, 0, err
				}
				out = append(out, up50aSample{
					roles: roles, program: program, target: target,
					category: "long_program",
				})
			}
		}
	}
	return out, heldBasis, nil
}

func RunUP51A() (UP51ASparseBindingResult, error) {
	const (
		steps = 140
		lr    = 0.10
		eps   = 1e-6
	)
	lengths := []int{12, 36}

	train, trainBasis, err := up51aSparseTrainSamples()
	if err != nil {
		return UP51ASparseBindingResult{}, err
	}
	held, heldBasis, err := up51aSparseHeldSamples(lengths)
	if err != nil {
		return UP51ASparseBindingResult{}, err
	}

	unitaryResult, err := up50aTrain(
		"unitary_sparse_basis",
		train, held,
		up50aApplyUnitary, true,
		steps, lr, eps,
	)
	if err != nil {
		return UP51ASparseBindingResult{}, err
	}
	nonUnitaryResult, err := up50aTrain(
		"nonunitary_sparse_basis",
		train, held,
		up50aApplyNonUnitary, false,
		steps, lr, eps,
	)
	if err != nil {
		return UP51ASparseBindingResult{}, err
	}

	primitive := up50aCategory(unitaryResult.CategoryMetrics, "unseen_basis_primitive")
	swap12 := up50aCategory(unitaryResult.CategoryMetrics, "unseen_swap12")
	role1 := up50aCategory(unitaryResult.CategoryMetrics, "derived_role1")
	role2 := up50aCategory(unitaryResult.CategoryMetrics, "derived_role2")
	long := up50aCategory(unitaryResult.CategoryMetrics, "long_program")

	gate :=
		unitaryResult.TrainAccuracy >= 0.99 &&
			unitaryResult.HeldOutAccuracy >= 0.98 &&
			primitive >= 0.98 &&
			swap12 >= 0.98 &&
			role1 >= 0.98 &&
			role2 >= 0.98 &&
			long >= 0.98 &&
			unitaryResult.MaxNormDrift <= 1e-9 &&
			unitaryResult.MaxRoundTripError <= 1e-8

	return UP51ASparseBindingResult{
		Schema:              UP51ASparseBindingSchema,
		Experiment:          "UP-51A-sparse-binding-generalization",
		SourceUP50ASeal:     "8aa5ed497de5a152c634211117e9001ed445b092",
		ValuesPerRole:       3,
		Roles:               3,
		JointDimension:      27,
		TrainingBasisRule:   "(role0+role1+role2) mod 3 == 0",
		TrainSingleStepOnly: true,
		HeldOutLengths:      append([]int(nil), lengths...),
		Steps:               steps,
		LearningRate:        lr,
		GradientEpsilon:     eps,
		Unitary:             unitaryResult,
		NonUnitary:          nonUnitaryResult,
		Diagnosis: UP51ADiagnosis{
			TrainingBasisStates:       trainBasis,
			HeldOutBasisStates:        heldBasis,
			TrainingCoverageFraction:  float64(trainBasis) / 27,
			UnitaryTrainAccuracy:      unitaryResult.TrainAccuracy,
			UnseenPrimitiveAccuracy:   primitive,
			UnseenSwap12Accuracy:      swap12,
			DerivedRole1Accuracy:      role1,
			DerivedRole2Accuracy:      role2,
			LongProgramAccuracy:       long,
			AggregateHeldOutAccuracy:  unitaryResult.HeldOutAccuracy,
			SparseGeneralizationGate:  gate,
			NonUnitaryHeldOutAccuracy: nonUnitaryResult.HeldOutAccuracy,
		},
	}, nil
}
