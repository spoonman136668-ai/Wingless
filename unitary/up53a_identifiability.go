package unitary

const UP53AIdentifiabilitySchema = "wingless.up53a-three-basis-identifiability.v1"

type UP53ATriadSpec struct {
	Name  string   `json:"name"`
	Basis [][3]int `json:"basis"`
}

type UP53ATriadResult struct {
	Spec                    UP53ATriadSpec   `json:"spec"`
	Unitary                 UP50APathResult  `json:"unitary"`
	NonUnitary              UP50APathResult  `json:"nonunitary_matched"`
	UnseenPrimitiveAccuracy float64          `json:"unseen_primitive_accuracy"`
	UnseenSwap12Accuracy    float64          `json:"unseen_swap12_accuracy"`
	DerivedRole1Accuracy    float64          `json:"derived_role1_accuracy"`
	DerivedRole2Accuracy    float64          `json:"derived_role2_accuracy"`
	LongProgramAccuracy     float64          `json:"long_program_accuracy"`
	Gate                    bool             `json:"gate"`
}

type UP53AIdentifiabilityResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP52ASeal      string             `json:"source_up52a_seal"`
	TrainingBasisStates  int                `json:"training_basis_states"`
	TrainSingleStepOnly  bool               `json:"train_single_step_only"`
	HeldOutLengths       []int              `json:"heldout_lengths"`
	Steps                int                `json:"steps"`
	LearningRate         float64            `json:"learning_rate"`
	GradientEpsilon      float64            `json:"gradient_epsilon"`
	Triads               []UP53ATriadResult `json:"triads"`
	InformativeTriadsPass bool              `json:"informative_triads_pass"`
	InvariantTriadPass    bool              `json:"invariant_triad_pass"`
}

func up53aTriads() []UP53ATriadSpec {
	return []UP53ATriadSpec{
		{
			Name:"invariant_diagonal",
			Basis:[][3]int{{0,0,0},{1,1,1},{2,2,2}},
		},
		{
			Name:"informative_cycle_a",
			Basis:[][3]int{{0,1,2},{1,2,0},{2,0,1}},
		},
		{
			Name:"informative_cycle_b",
			Basis:[][3]int{{0,2,1},{1,0,2},{2,1,0}},
		},
	}
}

func up53aSamples(basis [][3]int, lengths []int) ([]up50aSample, []up50aSample, error) {
	selected:=map[int]bool{}
	for _,roles:=range basis {
		selected[up52aKey(roles)]=true
	}

	var train []up50aSample
	var held []up50aSample
	for _,roles:=range up50aAllRoles() {
		if selected[up52aKey(roles)] {
			for start:=0;start<3;start++ {
				program:=[]UP50AInstruction{{Kind:"value_swap",Start:start}}
				target,err:=up50aTarget(roles,program)
				if err!=nil{return nil,nil,err}
				train=append(train,up50aSample{roles:roles,program:program,target:target,category:"train_role0_value"})
			}
			program:=[]UP50AInstruction{{Kind:"role_swap",Position:0}}
			target,err:=up50aTarget(roles,program)
			if err!=nil{return nil,nil,err}
			train=append(train,up50aSample{roles:roles,program:program,target:target,category:"train_swap01"})
			continue
		}

		for start:=0;start<3;start++ {
			primitive:=[]UP50AInstruction{{Kind:"value_swap",Start:start}}
			target,err:=up50aTarget(roles,primitive)
			if err!=nil{return nil,nil,err}
			held=append(held,up50aSample{roles:roles,program:primitive,target:target,category:"unseen_basis_primitive"})
		}
		swap12:=[]UP50AInstruction{{Kind:"role_swap",Position:1}}
		target,err:=up50aTarget(roles,swap12)
		if err!=nil{return nil,nil,err}
		held=append(held,up50aSample{roles:roles,program:swap12,target:target,category:"unseen_swap12"})

		for start:=0;start<3;start++ {
			role1:=[]UP50AInstruction{
				{Kind:"role_swap",Position:0},
				{Kind:"value_swap",Start:start},
				{Kind:"role_swap",Position:0},
			}
			target1,err:=up50aTarget(roles,role1)
			if err!=nil{return nil,nil,err}
			held=append(held,up50aSample{roles:roles,program:role1,target:target1,category:"derived_role1"})

			role2:=[]UP50AInstruction{
				{Kind:"role_swap",Position:1},
				{Kind:"role_swap",Position:0},
				{Kind:"value_swap",Start:start},
				{Kind:"role_swap",Position:0},
				{Kind:"role_swap",Position:1},
			}
			target2,err:=up50aTarget(roles,role2)
			if err!=nil{return nil,nil,err}
			held=append(held,up50aSample{roles:roles,program:role2,target:target2,category:"derived_role2"})
		}
	}

	for _,length:=range lengths {
		for _,roles:=range up50aAllRoles() {
			if selected[up52aKey(roles)] {continue}
			for seed:=0;seed<3;seed++ {
				program:=make([]UP50AInstruction,length)
				for i:=range program {
					switch (i+seed+roles[0]+2*roles[1]+roles[2])%4 {
					case 0:
						program[i]=UP50AInstruction{Kind:"role_swap",Position:0}
					case 1:
						program[i]=UP50AInstruction{Kind:"role_swap",Position:1}
					default:
						program[i]=UP50AInstruction{Kind:"value_swap",Start:(i*i+seed+roles[1])%3}
					}
				}
				target,err:=up50aTarget(roles,program)
				if err!=nil{return nil,nil,err}
				held=append(held,up50aSample{roles:roles,program:program,target:target,category:"long_program"})
			}
		}
	}
	return train,held,nil
}

func RunUP53A()(UP53AIdentifiabilityResult,error){
	const(
		steps=140
		lr=0.10
		eps=1e-6
	)
	lengths:=[]int{12,36}
	result:=UP53AIdentifiabilityResult{
		Schema:UP53AIdentifiabilitySchema,
		Experiment:"UP-53A-three-basis-identifiability",
		SourceUP52ASeal:"c2cbcfa72c1169e52b3873032a731058beef2dfb",
		TrainingBasisStates:3,
		TrainSingleStepOnly:true,
		HeldOutLengths:append([]int(nil),lengths...),
		Steps:steps,
		LearningRate:lr,
		GradientEpsilon:eps,
		InformativeTriadsPass:true,
	}
	for _,spec:=range up53aTriads(){
		train,held,err:=up53aSamples(spec.Basis,lengths)
		if err!=nil{return UP53AIdentifiabilityResult{},err}
		unitaryResult,err:=up50aTrain("unitary_"+spec.Name,train,held,up50aApplyUnitary,true,steps,lr,eps)
		if err!=nil{return UP53AIdentifiabilityResult{},err}
		nonunitaryResult,err:=up50aTrain("nonunitary_"+spec.Name,train,held,up50aApplyNonUnitary,false,steps,lr,eps)
		if err!=nil{return UP53AIdentifiabilityResult{},err}
		primitive:=up50aCategory(unitaryResult.CategoryMetrics,"unseen_basis_primitive")
		swap12:=up50aCategory(unitaryResult.CategoryMetrics,"unseen_swap12")
		role1:=up50aCategory(unitaryResult.CategoryMetrics,"derived_role1")
		role2:=up50aCategory(unitaryResult.CategoryMetrics,"derived_role2")
		long:=up50aCategory(unitaryResult.CategoryMetrics,"long_program")
		gate:=unitaryResult.TrainAccuracy>=0.99 &&
			unitaryResult.HeldOutAccuracy>=0.98 &&
			primitive>=0.98 && swap12>=0.98 && role1>=0.98 && role2>=0.98 && long>=0.98 &&
			unitaryResult.MaxNormDrift<=1e-9 &&
			unitaryResult.MaxRoundTripError<=1e-8

		result.Triads=append(result.Triads,UP53ATriadResult{
			Spec:spec,
			Unitary:unitaryResult,
			NonUnitary:nonunitaryResult,
			UnseenPrimitiveAccuracy:primitive,
			UnseenSwap12Accuracy:swap12,
			DerivedRole1Accuracy:role1,
			DerivedRole2Accuracy:role2,
			LongProgramAccuracy:long,
			Gate:gate,
		})
		if spec.Name=="invariant_diagonal" {
			result.InvariantTriadPass=gate
		} else if !gate {
			result.InformativeTriadsPass=false
		}
	}
	return result,nil
}
