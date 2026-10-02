package unitary

import (
	"fmt"
	"math"
)

const UP50AThreeRoleSchema = "wingless.up50a-three-role-transfer.v1"

type UP50AInstruction struct {
	Kind     string `json:"kind"`
	Start    int    `json:"start,omitempty"`
	Position int    `json:"position,omitempty"`
}

type up50aSample struct {
	roles    [3]int
	program  []UP50AInstruction
	target   [3]int
	category string
}

type UP50ACategoryMetric struct {
	Category string  `json:"category"`
	Samples  int     `json:"samples"`
	Accuracy float64 `json:"accuracy"`
}

type UP50APathResult struct {
	Name                 string                  `json:"name"`
	InitialTrainAccuracy float64                 `json:"initial_train_accuracy"`
	TrainAccuracy        float64                 `json:"train_accuracy"`
	HeldOutAccuracy      float64                 `json:"heldout_accuracy"`
	CategoryMetrics      []UP50ACategoryMetric   `json:"category_metrics"`
	InitialLoss          float64                 `json:"initial_loss"`
	FinalLoss            float64                 `json:"final_loss"`
	Parameters           []float64               `json:"parameters"`
	MaxNormDrift         float64                 `json:"max_norm_drift"`
	MaxRoundTripError    float64                 `json:"max_round_trip_error"`
}

type UP50ADiagnosis struct {
	UnitaryTrainAccuracy       float64 `json:"unitary_train_accuracy"`
	UnitaryHeldOutAccuracy     float64 `json:"unitary_heldout_accuracy"`
	UnseenSwap12Accuracy       float64 `json:"unseen_swap12_accuracy"`
	DerivedRole1Accuracy       float64 `json:"derived_role1_accuracy"`
	DerivedRole2Accuracy       float64 `json:"derived_role2_accuracy"`
	LongProgramAccuracy        float64 `json:"long_program_accuracy"`
	ThreeRoleTransferGate      bool    `json:"three_role_transfer_gate"`
	NonUnitaryHeldOutAccuracy  float64 `json:"nonunitary_heldout_accuracy"`
}

type UP50AThreeRoleResult struct {
	Schema                 string          `json:"schema"`
	Experiment             string          `json:"experiment"`
	SourceUP49ASeal        string          `json:"source_up49a_seal"`
	ValuesPerRole          int             `json:"values_per_role"`
	Roles                  int             `json:"roles"`
	JointDimension         int             `json:"joint_dimension"`
	TrainSingleStepOnly    bool            `json:"train_single_step_only"`
	Swap12Trained          bool            `json:"swap12_trained"`
	Role1PrimitiveTrained  bool            `json:"role1_primitive_trained"`
	Role2PrimitiveTrained  bool            `json:"role2_primitive_trained"`
	HeldOutLengths         []int           `json:"heldout_lengths"`
	Steps                  int             `json:"steps"`
	LearningRate           float64         `json:"learning_rate"`
	GradientEpsilon        float64         `json:"gradient_epsilon"`
	Unitary                UP50APathResult `json:"unitary"`
	NonUnitary             UP50APathResult `json:"nonunitary_matched"`
	Diagnosis              UP50ADiagnosis  `json:"diagnosis"`
}

func up50aIndex(roles [3]int) (int, error) {
	for _, value := range roles {
		if value < 0 || value >= 3 {
			return 0, fmt.Errorf("UP50A role value out of range")
		}
	}
	return roles[0]*9 + roles[1]*3 + roles[2], nil
}

func up50aDecode(index int) ([3]int, error) {
	if index < 0 || index >= 27 {
		return [3]int{}, fmt.Errorf("UP50A joint index out of range")
	}
	return [3]int{index / 9, (index / 3) % 3, index % 3}, nil
}

func up50aBasis(roles [3]int) (State, error) {
	index, err := up50aIndex(roles)
	if err != nil {
		return nil, err
	}
	state := make(State, 27)
	state[index] = 1
	return state, nil
}

func up50aApplyLogical(roles [3]int, instruction UP50AInstruction) ([3]int, error) {
	switch instruction.Kind {
	case "value_swap":
		if instruction.Start < 0 || instruction.Start >= 3 {
			return [3]int{}, fmt.Errorf("UP50A value start out of range")
		}
		a := instruction.Start
		b := (a + 1) % 3
		switch roles[0] {
		case a:
			roles[0] = b
		case b:
			roles[0] = a
		}
	case "role_swap":
		if instruction.Position < 0 || instruction.Position > 1 {
			return [3]int{}, fmt.Errorf("UP50A role swap position out of range")
		}
		p := instruction.Position
		roles[p], roles[p+1] = roles[p+1], roles[p]
	default:
		return [3]int{}, fmt.Errorf("UP50A unknown instruction %q", instruction.Kind)
	}
	return roles, nil
}

func up50aTarget(roles [3]int, program []UP50AInstruction) ([3]int, error) {
	var err error
	for _, instruction := range program {
		roles, err = up50aApplyLogical(roles, instruction)
		if err != nil {
			return [3]int{}, err
		}
	}
	return roles, nil
}

func up50aAllRoles() [][3]int {
	out := make([][3]int, 0, 27)
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			for c := 0; c < 3; c++ {
				out = append(out, [3]int{a, b, c})
			}
		}
	}
	return out
}

func up50aTrainSamples() ([]up50aSample, error) {
	var out []up50aSample
	for _, roles := range up50aAllRoles() {
		for start := 0; start < 3; start++ {
			program := []UP50AInstruction{{Kind: "value_swap", Start: start}}
			target, err := up50aTarget(roles, program)
			if err != nil {
				return nil, err
			}
			out = append(out, up50aSample{roles: roles, program: program, target: target, category: "train_role0_value"})
		}
		program := []UP50AInstruction{{Kind: "role_swap", Position: 0}}
		target, err := up50aTarget(roles, program)
		if err != nil {
			return nil, err
		}
		out = append(out, up50aSample{roles: roles, program: program, target: target, category: "train_swap01"})
	}
	return out, nil
}

func up50aHeldSamples(lengths []int) ([]up50aSample, error) {
	var out []up50aSample
	for _, roles := range up50aAllRoles() {
		// Role-swap parameter is never trained at position 1.
		p := []UP50AInstruction{{Kind: "role_swap", Position: 1}}
		t, err := up50aTarget(roles, p)
		if err != nil {
			return nil, err
		}
		out = append(out, up50aSample{roles: roles, program: p, target: t, category: "unseen_swap12"})

		for start := 0; start < 3; start++ {
			// Move role 1 into role 0, apply learned value rule, move it back.
			p1 := []UP50AInstruction{
				{Kind: "role_swap", Position: 0},
				{Kind: "value_swap", Start: start},
				{Kind: "role_swap", Position: 0},
			}
			t1, err := up50aTarget(roles, p1)
			if err != nil {
				return nil, err
			}
			out = append(out, up50aSample{roles: roles, program: p1, target: t1, category: "derived_role1"})

			// Move role 2 across two positions, mutate it, and restore ordering.
			p2 := []UP50AInstruction{
				{Kind: "role_swap", Position: 1},
				{Kind: "role_swap", Position: 0},
				{Kind: "value_swap", Start: start},
				{Kind: "role_swap", Position: 0},
				{Kind: "role_swap", Position: 1},
			}
			t2, err := up50aTarget(roles, p2)
			if err != nil {
				return nil, err
			}
			out = append(out, up50aSample{roles: roles, program: p2, target: t2, category: "derived_role2"})
		}
	}
	for _, length := range lengths {
		for _, roles := range up50aAllRoles() {
			for seed := 0; seed < 3; seed++ {
				program := make([]UP50AInstruction, length)
				for i := range program {
					switch (i + seed + roles[0] + 2*roles[1] + roles[2]) % 4 {
					case 0:
						program[i] = UP50AInstruction{Kind: "role_swap", Position: 0}
					case 1:
						program[i] = UP50AInstruction{Kind: "role_swap", Position: 1}
					default:
						program[i] = UP50AInstruction{Kind: "value_swap", Start: (i*i + seed + roles[1]) % 3}
					}
				}
				target, err := up50aTarget(roles, program)
				if err != nil {
					return nil, err
				}
				out = append(out, up50aSample{roles: roles, program: program, target: target, category: "long_program"})
			}
		}
	}
	return out, nil
}

func up50aCouplings(instruction UP50AInstruction, params []float64) ([]Coupling, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("UP50A parameter count=%d want=2", len(params))
	}
	switch instruction.Kind {
	case "value_swap":
		if instruction.Start < 0 || instruction.Start >= 3 {
			return nil, fmt.Errorf("UP50A value start out of range")
		}
		a := instruction.Start
		b := (a + 1) % 3
		out := make([]Coupling, 0, 9)
		for r1 := 0; r1 < 3; r1++ {
			for r2 := 0; r2 < 3; r2++ {
				ia, _ := up50aIndex([3]int{a, r1, r2})
				ib, _ := up50aIndex([3]int{b, r1, r2})
				out = append(out, Coupling{A: ia, B: ib, Theta: params[0]})
			}
		}
		return out, nil
	case "role_swap":
		if instruction.Position < 0 || instruction.Position > 1 {
			return nil, fmt.Errorf("UP50A role swap position out of range")
		}
		out := make([]Coupling, 0, 9)
		for other := 0; other < 3; other++ {
			for a := 0; a < 3; a++ {
				for b := a + 1; b < 3; b++ {
					var ra, rb [3]int
					if instruction.Position == 0 {
						ra = [3]int{a, b, other}
						rb = [3]int{b, a, other}
					} else {
						ra = [3]int{other, a, b}
						rb = [3]int{other, b, a}
					}
					ia, _ := up50aIndex(ra)
					ib, _ := up50aIndex(rb)
					out = append(out, Coupling{A: ia, B: ib, Theta: params[1]})
				}
			}
		}
		return out, nil
	default:
		return nil, fmt.Errorf("UP50A unknown instruction %q", instruction.Kind)
	}
}

type up50aApply func(State, []float64, []UP50AInstruction) (State, error)

func up50aApplyUnitary(initial State, params []float64, program []UP50AInstruction) (State, error) {
	out := append(State(nil), initial...)
	for _, instruction := range program {
		couplings, err := up50aCouplings(instruction, params)
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

func up50aApplyNonUnitary(initial State, params []float64, program []UP50AInstruction) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	out := append(State(nil), initial...)
	for _, instruction := range program {
		couplings, err := up50aCouplings(instruction, params)
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

func up50aInverseUnitary(state State, params []float64, program []UP50AInstruction) (State, error) {
	out := append(State(nil), state...)
	for i := len(program)-1; i >= 0; i-- {
		couplings, err := up50aCouplings(program[i], params)
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

func up50aEvaluate(params []float64, samples []up50aSample, apply up50aApply, roundTrip bool) (
	loss, accuracy, maxNormDrift, maxRoundTrip float64,
	categories []UP50ACategoryMetric,
	err error,
) {
	type counts struct{ samples, hits int }
	per := map[string]*counts{}
	order := []string{}
	hits := 0
	for _, sample := range samples {
		initial, e := up50aBasis(sample.roles)
		if e != nil { err=e; return }
		evolved, e := apply(initial, params, sample.program)
		if e != nil { err=e; return }
		probs, e := Probabilities(evolved)
		if e != nil { err=e; return }
		targetIndex, e := up50aIndex(sample.target)
		if e != nil { err=e; return }
		p := probs[targetIndex]
		if p < 1e-15 { p = 1e-15 }
		loss += -math.Log(p)
		observedIndex, e := ArgMax(evolved)
		if e != nil { err=e; return }
		observed, e := up50aDecode(observedIndex)
		if e != nil { err=e; return }
		if _, ok := per[sample.category]; !ok {
			per[sample.category] = &counts{}
			order = append(order, sample.category)
		}
		per[sample.category].samples++
		if observed == sample.target {
			hits++
			per[sample.category].hits++
		}
		norm2, e := NormSquared(evolved)
		if e != nil { err=e; return }
		drift := math.Abs(norm2-1)
		if drift > maxNormDrift { maxNormDrift = drift }
		if roundTrip {
			recovered, e := up50aInverseUnitary(evolved, params, sample.program)
			if e != nil { err=e; return }
			distance, e := L2Distance(initial, recovered)
			if e != nil { err=e; return }
			if distance > maxRoundTrip { maxRoundTrip = distance }
		}
	}
	if len(samples)==0 { err=fmt.Errorf("UP50A empty samples"); return }
	loss /= float64(len(samples))
	accuracy = float64(hits)/float64(len(samples))
	for _, category := range order {
		c := per[category]
		categories = append(categories, UP50ACategoryMetric{
			Category: category, Samples: c.samples, Accuracy: float64(c.hits)/float64(c.samples),
		})
	}
	return
}

func up50aLoss(params []float64, samples []up50aSample, apply up50aApply) (float64,error) {
	loss,_,_,_,_,err := up50aEvaluate(params,samples,apply,false)
	return loss,err
}

func up50aTrain(name string, train,held []up50aSample, apply up50aApply, roundTrip bool, steps int, lr,eps float64) (UP50APathResult,error) {
	params:=[]float64{0.1,0.1}
	initialLoss,initialAcc,_,_,_,err:=up50aEvaluate(params,train,apply,false)
	if err!=nil{return UP50APathResult{},err}
	for step:=0;step<steps;step++{
		grad:=make([]float64,2)
		for j:=0;j<2;j++{
			plus:=append([]float64(nil),params...)
			minus:=append([]float64(nil),params...)
			plus[j]+=eps; minus[j]-=eps
			lp,e:=up50aLoss(plus,train,apply); if e!=nil{return UP50APathResult{},e}
			lm,e:=up50aLoss(minus,train,apply); if e!=nil{return UP50APathResult{},e}
			grad[j]=(lp-lm)/(2*eps)
			if !finite(grad[j]){return UP50APathResult{},fmt.Errorf("UP50A %s gradient %d nonfinite",name,j)}
		}
		for j:=0;j<2;j++{params[j]-=lr*grad[j]}
	}
	finalLoss,trainAcc,trainNorm,trainRT,_,err:=up50aEvaluate(params,train,apply,roundTrip)
	if err!=nil{return UP50APathResult{},err}
	_,heldAcc,heldNorm,heldRT,categories,err:=up50aEvaluate(params,held,apply,roundTrip)
	if err!=nil{return UP50APathResult{},err}
	return UP50APathResult{
		Name:name,InitialTrainAccuracy:initialAcc,TrainAccuracy:trainAcc,HeldOutAccuracy:heldAcc,
		CategoryMetrics:categories,InitialLoss:initialLoss,FinalLoss:finalLoss,Parameters:append([]float64(nil),params...),
		MaxNormDrift:math.Max(trainNorm,heldNorm),MaxRoundTripError:math.Max(trainRT,heldRT),
	},nil
}

func up50aCategory(metrics []UP50ACategoryMetric,name string) float64 {
	for _,m:=range metrics{if m.Category==name{return m.Accuracy}}
	return 0
}

func RunUP50A() (UP50AThreeRoleResult,error) {
	const(steps=140; lr=0.10; eps=1e-6)
	lengths:=[]int{12,36}
	train,err:=up50aTrainSamples(); if err!=nil{return UP50AThreeRoleResult{},err}
	held,err:=up50aHeldSamples(lengths); if err!=nil{return UP50AThreeRoleResult{},err}
	u,err:=up50aTrain("unitary",train,held,up50aApplyUnitary,true,steps,lr,eps)
	if err!=nil{return UP50AThreeRoleResult{},err}
	n,err:=up50aTrain("nonunitary_matched",train,held,up50aApplyNonUnitary,false,steps,lr,eps)
	if err!=nil{return UP50AThreeRoleResult{},err}
	swap12:=up50aCategory(u.CategoryMetrics,"unseen_swap12")
	role1:=up50aCategory(u.CategoryMetrics,"derived_role1")
	role2:=up50aCategory(u.CategoryMetrics,"derived_role2")
	long:=up50aCategory(u.CategoryMetrics,"long_program")
	gate:=u.TrainAccuracy>=0.99&&u.HeldOutAccuracy>=0.98&&swap12>=0.98&&role1>=0.98&&role2>=0.98&&long>=0.98&&u.MaxNormDrift<=1e-9&&u.MaxRoundTripError<=1e-8
	return UP50AThreeRoleResult{
		Schema:UP50AThreeRoleSchema,Experiment:"UP-50A-three-role-transfer",
		SourceUP49ASeal:"d5a709e13ead3a7c2fef4df913a1b3e3acd3139b",
		ValuesPerRole:3,Roles:3,JointDimension:27,TrainSingleStepOnly:true,
		Swap12Trained:false,Role1PrimitiveTrained:false,Role2PrimitiveTrained:false,
		HeldOutLengths:append([]int(nil),lengths...),Steps:steps,LearningRate:lr,GradientEpsilon:eps,
		Unitary:u,NonUnitary:n,
		Diagnosis:UP50ADiagnosis{
			UnitaryTrainAccuracy:u.TrainAccuracy,UnitaryHeldOutAccuracy:u.HeldOutAccuracy,
			UnseenSwap12Accuracy:swap12,DerivedRole1Accuracy:role1,DerivedRole2Accuracy:role2,
			LongProgramAccuracy:long,ThreeRoleTransferGate:gate,NonUnitaryHeldOutAccuracy:n.HeldOutAccuracy,
		},
	},nil
}
