package unitary

import (
	"fmt"
	"math"
)

const UP49ARoleBindingSchema = "wingless.up49a-role-binding-composition.v1"

type UP49AInstruction struct {
	Kind  string `json:"kind"`
	Start int    `json:"start,omitempty"`
}

type up49aSample struct {
	left, right int
	program     []UP49AInstruction
	targetLeft  int
	targetRight int
	category    string
}

type UP49ACategoryMetric struct {
	Category string  `json:"category"`
	Samples  int     `json:"samples"`
	Accuracy float64 `json:"accuracy"`
}

type UP49APathResult struct {
	Name              string                  `json:"name"`
	InitialTrainAccuracy float64               `json:"initial_train_accuracy"`
	TrainAccuracy     float64                  `json:"train_accuracy"`
	HeldOutAccuracy   float64                  `json:"heldout_accuracy"`
	CategoryMetrics   []UP49ACategoryMetric    `json:"category_metrics"`
	InitialLoss       float64                  `json:"initial_loss"`
	FinalLoss         float64                  `json:"final_loss"`
	Parameters        []float64                `json:"parameters"`
	MaxNormDrift      float64                  `json:"max_norm_drift"`
	MaxRoundTripError float64                  `json:"max_round_trip_error"`
}

type UP49ADiagnosis struct {
	UnitaryTrainAccuracy            float64 `json:"unitary_train_accuracy"`
	UnitaryHeldOutAccuracy          float64 `json:"unitary_heldout_accuracy"`
	ConjugatedRightAccuracy         float64 `json:"conjugated_right_accuracy"`
	LongProgramAccuracy             float64 `json:"long_program_accuracy"`
	RoleBindingGate                 bool    `json:"role_binding_gate"`
	NonUnitaryHeldOutAccuracy       float64 `json:"nonunitary_heldout_accuracy"`
	UnitaryHeldOutAdvantage         float64 `json:"unitary_heldout_advantage"`
}

type UP49ARoleBindingResult struct {
	Schema              string          `json:"schema"`
	Experiment          string          `json:"experiment"`
	SourceUP48ASeal     string          `json:"source_up48a_seal"`
	ValuesPerRole       int             `json:"values_per_role"`
	JointDimension      int             `json:"joint_dimension"`
	TrainSingleStepOnly bool            `json:"train_single_step_only"`
	RightPrimitiveTrained bool          `json:"right_primitive_trained"`
	HeldOutLengths      []int           `json:"heldout_lengths"`
	Steps               int             `json:"steps"`
	LearningRate        float64         `json:"learning_rate"`
	GradientEpsilon     float64         `json:"gradient_epsilon"`
	Unitary             UP49APathResult `json:"unitary"`
	NonUnitary          UP49APathResult `json:"nonunitary_matched"`
	Diagnosis           UP49ADiagnosis  `json:"diagnosis"`
}

func up49aIndex(left, right int) (int, error) {
	if left < 0 || left >= 4 || right < 0 || right >= 4 {
		return 0, fmt.Errorf("UP49A role value out of range")
	}
	return left*4 + right, nil
}

func up49aDecodeIndex(index int) (int, int, error) {
	if index < 0 || index >= 16 {
		return 0, 0, fmt.Errorf("UP49A joint index out of range")
	}
	return index / 4, index % 4, nil
}

func up49aBasis(left, right int) (State, error) {
	index, err := up49aIndex(left, right)
	if err != nil {
		return nil, err
	}
	state := make(State, 16)
	state[index] = 1
	return state, nil
}

func up49aApplyLogical(left, right int, instruction UP49AInstruction) (int, int, error) {
	switch instruction.Kind {
	case "left_swap":
		if instruction.Start < 0 || instruction.Start >= 4 {
			return 0, 0, fmt.Errorf("UP49A left_swap start out of range")
		}
		a := instruction.Start
		b := (a + 1) % 4
		switch left {
		case a:
			left = b
		case b:
			left = a
		}
	case "swap_roles":
		left, right = right, left
	default:
		return 0, 0, fmt.Errorf("UP49A unknown instruction kind %q", instruction.Kind)
	}
	return left, right, nil
}

func up49aTarget(left, right int, program []UP49AInstruction) (int, int, error) {
	var err error
	for _, instruction := range program {
		left, right, err = up49aApplyLogical(left, right, instruction)
		if err != nil {
			return 0, 0, err
		}
	}
	return left, right, nil
}

func up49aTrainSamples() ([]up49aSample, error) {
	var out []up49aSample
	for left := 0; left < 4; left++ {
		for right := 0; right < 4; right++ {
			for start := 0; start < 4; start++ {
				program := []UP49AInstruction{{Kind: "left_swap", Start: start}}
				tl, tr, err := up49aTarget(left, right, program)
				if err != nil {
					return nil, err
				}
				out = append(out, up49aSample{
					left: left, right: right, program: program,
					targetLeft: tl, targetRight: tr, category: "train_left_primitive",
				})
			}
			program := []UP49AInstruction{{Kind: "swap_roles"}}
			tl, tr, err := up49aTarget(left, right, program)
			if err != nil {
				return nil, err
			}
			out = append(out, up49aSample{
				left: left, right: right, program: program,
				targetLeft: tl, targetRight: tr, category: "train_role_swap",
			})
		}
	}
	return out, nil
}

func up49aHeldSamples(lengths []int) ([]up49aSample, error) {
	var out []up49aSample

	// Unseen right-role primitive emerges only through conjugation:
	// swap_roles -> left_swap -> swap_roles.
	for left := 0; left < 4; left++ {
		for right := 0; right < 4; right++ {
			for start := 0; start < 4; start++ {
				program := []UP49AInstruction{
					{Kind: "swap_roles"},
					{Kind: "left_swap", Start: start},
					{Kind: "swap_roles"},
				}
				tl, tr, err := up49aTarget(left, right, program)
				if err != nil {
					return nil, err
				}
				out = append(out, up49aSample{
					left: left, right: right, program: program,
					targetLeft: tl, targetRight: tr, category: "conjugated_right",
				})
			}
		}
	}

	for _, length := range lengths {
		for left := 0; left < 4; left++ {
			for right := 0; right < 4; right++ {
				for seed := 0; seed < 4; seed++ {
					program := make([]UP49AInstruction, length)
					for i := range program {
						if (i+seed+left+right)%3 == 0 {
							program[i] = UP49AInstruction{Kind: "swap_roles"}
						} else {
							program[i] = UP49AInstruction{
								Kind:  "left_swap",
								Start: (i*i + 3*i + seed + 2*left + right) % 4,
							}
						}
					}
					tl, tr, err := up49aTarget(left, right, program)
					if err != nil {
						return nil, err
					}
					out = append(out, up49aSample{
						left: left, right: right, program: program,
						targetLeft: tl, targetRight: tr, category: "long_program",
					})
				}
			}
		}
	}
	return out, nil
}

func up49aCouplings(instruction UP49AInstruction, params []float64) ([]Coupling, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("UP49A parameter count=%d want=2", len(params))
	}
	switch instruction.Kind {
	case "left_swap":
		if instruction.Start < 0 || instruction.Start >= 4 {
			return nil, fmt.Errorf("UP49A left start out of range")
		}
		a := instruction.Start
		b := (a + 1) % 4
		couplings := make([]Coupling, 0, 4)
		for right := 0; right < 4; right++ {
			ia, _ := up49aIndex(a, right)
			ib, _ := up49aIndex(b, right)
			couplings = append(couplings, Coupling{A: ia, B: ib, Theta: params[0]})
		}
		return couplings, nil
	case "swap_roles":
		couplings := make([]Coupling, 0, 6)
		for left := 0; left < 4; left++ {
			for right := left + 1; right < 4; right++ {
				ia, _ := up49aIndex(left, right)
				ib, _ := up49aIndex(right, left)
				couplings = append(couplings, Coupling{A: ia, B: ib, Theta: params[1]})
			}
		}
		return couplings, nil
	default:
		return nil, fmt.Errorf("UP49A unknown instruction %q", instruction.Kind)
	}
}

type up49aApply func(State, []float64, []UP49AInstruction) (State, error)

func up49aApplyUnitary(initial State, params []float64, program []UP49AInstruction) (State, error) {
	out := append(State(nil), initial...)
	for _, instruction := range program {
		couplings, err := up49aCouplings(instruction, params)
		if err != nil {
			return nil, err
		}
		out, err = Propagate(out, couplings)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func up49aApplyNonUnitary(initial State, params []float64, program []UP49AInstruction) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	out := append(State(nil), initial...)
	for _, instruction := range program {
		couplings, err := up49aCouplings(instruction, params)
		if err != nil {
			return nil, err
		}
		for _, coupling := range couplings {
			g := coupling.Theta
			a := out[coupling.A]
			b := out[coupling.B]
			out[coupling.A] = a - complex(g, 0)*b
			out[coupling.B] = complex(g, 0)*a + b
		}
	}
	return out, nil
}

func up49aInverseUnitary(state State, params []float64, program []UP49AInstruction) (State, error) {
	out := append(State(nil), state...)
	for i := len(program) - 1; i >= 0; i-- {
		couplings, err := up49aCouplings(program[i], params)
		if err != nil {
			return nil, err
		}
		for j := range couplings {
			couplings[j].Theta = -couplings[j].Theta
		}
		out, err = Propagate(out, couplings)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func up49aEvaluate(params []float64, samples []up49aSample, apply up49aApply, roundTrip bool) (
	loss, accuracy, maxNormDrift, maxRoundTrip float64,
	categories []UP49ACategoryMetric,
	err error,
) {
	type counts struct{ samples, hits int }
	perCategory := map[string]*counts{}
	order := []string{}
	hits := 0

	for _, sample := range samples {
		initial, e := up49aBasis(sample.left, sample.right)
		if e != nil {
			err = e
			return
		}
		evolved, e := apply(initial, params, sample.program)
		if e != nil {
			err = e
			return
		}
		probabilities, e := Probabilities(evolved)
		if e != nil {
			err = e
			return
		}
		targetIndex, e := up49aIndex(sample.targetLeft, sample.targetRight)
		if e != nil {
			err = e
			return
		}
		p := probabilities[targetIndex]
		if p < 1e-15 {
			p = 1e-15
		}
		loss += -math.Log(p)

		observedIndex, e := ArgMax(evolved)
		if e != nil {
			err = e
			return
		}
		observedLeft, observedRight, e := up49aDecodeIndex(observedIndex)
		if e != nil {
			err = e
			return
		}
		if _, ok := perCategory[sample.category]; !ok {
			perCategory[sample.category] = &counts{}
			order = append(order, sample.category)
		}
		perCategory[sample.category].samples++
		if observedLeft == sample.targetLeft && observedRight == sample.targetRight {
			hits++
			perCategory[sample.category].hits++
		}

		norm2, e := NormSquared(evolved)
		if e != nil {
			err = e
			return
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		if roundTrip {
			recovered, e := up49aInverseUnitary(evolved, params, sample.program)
			if e != nil {
				err = e
				return
			}
			distance, e := L2Distance(initial, recovered)
			if e != nil {
				err = e
				return
			}
			if distance > maxRoundTrip {
				maxRoundTrip = distance
			}
		}
	}
	if len(samples) == 0 {
		err = fmt.Errorf("UP49A empty sample set")
		return
	}
	loss /= float64(len(samples))
	accuracy = float64(hits) / float64(len(samples))
	for _, category := range order {
		c := perCategory[category]
		categories = append(categories, UP49ACategoryMetric{
			Category: category,
			Samples: c.samples,
			Accuracy: float64(c.hits) / float64(c.samples),
		})
	}
	return
}

func up49aLoss(params []float64, samples []up49aSample, apply up49aApply) (float64, error) {
	loss, _, _, _, _, err := up49aEvaluate(params, samples, apply, false)
	return loss, err
}

func up49aTrain(name string, train, held []up49aSample, apply up49aApply, roundTrip bool, steps int, lr, eps float64) (UP49APathResult, error) {
	params := []float64{0.1, 0.1}
	initialLoss, initialAccuracy, _, _, _, err := up49aEvaluate(params, train, apply, false)
	if err != nil {
		return UP49APathResult{}, err
	}
	for step := 0; step < steps; step++ {
		gradient := make([]float64, len(params))
		for j := range params {
			plus := append([]float64(nil), params...)
			minus := append([]float64(nil), params...)
			plus[j] += eps
			minus[j] -= eps
			lp, e := up49aLoss(plus, train, apply)
			if e != nil {
				return UP49APathResult{}, e
			}
			lm, e := up49aLoss(minus, train, apply)
			if e != nil {
				return UP49APathResult{}, e
			}
			gradient[j] = (lp - lm) / (2 * eps)
			if !finite(gradient[j]) {
				return UP49APathResult{}, fmt.Errorf("UP49A %s gradient %d not finite", name, j)
			}
		}
		for j := range params {
			params[j] -= lr * gradient[j]
		}
	}
	finalLoss, trainAccuracy, trainNorm, trainRound, _, err := up49aEvaluate(params, train, apply, roundTrip)
	if err != nil {
		return UP49APathResult{}, err
	}
	_, heldAccuracy, heldNorm, heldRound, categories, err := up49aEvaluate(params, held, apply, roundTrip)
	if err != nil {
		return UP49APathResult{}, err
	}
	return UP49APathResult{
		Name: name,
		InitialTrainAccuracy: initialAccuracy,
		TrainAccuracy: trainAccuracy,
		HeldOutAccuracy: heldAccuracy,
		CategoryMetrics: categories,
		InitialLoss: initialLoss,
		FinalLoss: finalLoss,
		Parameters: append([]float64(nil), params...),
		MaxNormDrift: math.Max(trainNorm, heldNorm),
		MaxRoundTripError: math.Max(trainRound, heldRound),
	}, nil
}

func up49aCategory(metrics []UP49ACategoryMetric, name string) float64 {
	for _, metric := range metrics {
		if metric.Category == name {
			return metric.Accuracy
		}
	}
	return 0
}

func RunUP49A() (UP49ARoleBindingResult, error) {
	const (
		steps = 120
		lr = 0.12
		eps = 1e-6
	)
	lengths := []int{9, 27, 81}
	train, err := up49aTrainSamples()
	if err != nil {
		return UP49ARoleBindingResult{}, err
	}
	held, err := up49aHeldSamples(lengths)
	if err != nil {
		return UP49ARoleBindingResult{}, err
	}
	unitaryResult, err := up49aTrain("unitary", train, held, up49aApplyUnitary, true, steps, lr, eps)
	if err != nil {
		return UP49ARoleBindingResult{}, err
	}
	nonUnitaryResult, err := up49aTrain("nonunitary_matched", train, held, up49aApplyNonUnitary, false, steps, lr, eps)
	if err != nil {
		return UP49ARoleBindingResult{}, err
	}
	rightAccuracy := up49aCategory(unitaryResult.CategoryMetrics, "conjugated_right")
	longAccuracy := up49aCategory(unitaryResult.CategoryMetrics, "long_program")
	gate := unitaryResult.TrainAccuracy >= 0.99 &&
		unitaryResult.HeldOutAccuracy >= 0.98 &&
		rightAccuracy >= 0.98 &&
		longAccuracy >= 0.98 &&
		unitaryResult.MaxNormDrift <= 1e-9 &&
		unitaryResult.MaxRoundTripError <= 1e-8

	return UP49ARoleBindingResult{
		Schema: UP49ARoleBindingSchema,
		Experiment: "UP-49A-role-binding-composition",
		SourceUP48ASeal: "9d76277d68b392597ac3a4e6cad2352a3c15b849",
		ValuesPerRole: 4,
		JointDimension: 16,
		TrainSingleStepOnly: true,
		RightPrimitiveTrained: false,
		HeldOutLengths: append([]int(nil), lengths...),
		Steps: steps,
		LearningRate: lr,
		GradientEpsilon: eps,
		Unitary: unitaryResult,
		NonUnitary: nonUnitaryResult,
		Diagnosis: UP49ADiagnosis{
			UnitaryTrainAccuracy: unitaryResult.TrainAccuracy,
			UnitaryHeldOutAccuracy: unitaryResult.HeldOutAccuracy,
			ConjugatedRightAccuracy: rightAccuracy,
			LongProgramAccuracy: longAccuracy,
			RoleBindingGate: gate,
			NonUnitaryHeldOutAccuracy: nonUnitaryResult.HeldOutAccuracy,
			UnitaryHeldOutAdvantage: unitaryResult.HeldOutAccuracy - nonUnitaryResult.HeldOutAccuracy,
		},
	}, nil
}
