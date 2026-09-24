package unitary

const UP52ACoverageLadderSchema = "wingless.up52a-basis-coverage-ladder.v1"

type UP52ACoveragePoint struct {
	TrainingBasisStates      int             `json:"training_basis_states"`
	HeldoutBasisStates       int             `json:"heldout_basis_states"`
	TrainingCoverageFraction float64         `json:"training_coverage_fraction"`
	Unitary                  UP50APathResult `json:"unitary"`
	NonUnitary               UP50APathResult `json:"nonunitary_matched"`
	UnseenPrimitiveAccuracy  float64         `json:"unseen_primitive_accuracy"`
	UnseenSwap12Accuracy     float64         `json:"unseen_swap12_accuracy"`
	DerivedRole1Accuracy     float64         `json:"derived_role1_accuracy"`
	DerivedRole2Accuracy     float64         `json:"derived_role2_accuracy"`
	LongProgramAccuracy      float64         `json:"long_program_accuracy"`
	Gate                     bool            `json:"gate"`
}

type UP52ACoverageLadderResult struct {
	Schema            string               `json:"schema"`
	Experiment        string               `json:"experiment"`
	SourceUP51ASeal   string               `json:"source_up51a_seal"`
	TrainingBasisOrder [][3]int             `json:"training_basis_order"`
	CoverageLevels    []int                `json:"coverage_levels"`
	TrainSingleStepOnly bool               `json:"train_single_step_only"`
	HeldOutLengths    []int                `json:"heldout_lengths"`
	Steps             int                  `json:"steps"`
	LearningRate      float64              `json:"learning_rate"`
	GradientEpsilon   float64              `json:"gradient_epsilon"`
	Points            []UP52ACoveragePoint `json:"points"`
	MinimumPassingBasis int                `json:"minimum_passing_basis"`
}

func up52aBasisOrder() [][3]int {
	return [][3]int{
		{0,0,0},
		{1,1,1},
		{2,2,2},
		{0,1,2},
		{1,2,0},
		{2,0,1},
		{0,2,1},
		{1,0,2},
		{2,1,0},
	}
}

func up52aKey(r [3]int) int {
	return r[0]*9 + r[1]*3 + r[2]
}

func up52aSamples(trainingBasis int, lengths []int) ([]up50aSample, []up50aSample, int, error) {
	order := up52aBasisOrder()
	selected := map[int]bool{}
	for i:=0; i<trainingBasis; i++ {
		selected[up52aKey(order[i])] = true
	}

	var train []up50aSample
	var held []up50aSample
	for _, roles := range up50aAllRoles() {
		isTrain := selected[up52aKey(roles)]
		if isTrain {
			for start:=0; start<3; start++ {
				program:=[]UP50AInstruction{{Kind:"value_swap",Start:start}}
				target,err:=up50aTarget(roles,program)
				if err!=nil{return nil,nil,0,err}
				train=append(train,up50aSample{roles:roles,program:program,target:target,category:"train_sparse_role0_value"})
			}
			program:=[]UP50AInstruction{{Kind:"role_swap",Position:0}}
			target,err:=up50aTarget(roles,program)
			if err!=nil{return nil,nil,0,err}
			train=append(train,up50aSample{roles:roles,program:program,target:target,category:"train_sparse_swap01"})
			continue
		}

		for start:=0; start<3; start++ {
			primitive:=[]UP50AInstruction{{Kind:"value_swap",Start:start}}
			target,err:=up50aTarget(roles,primitive)
			if err!=nil{return nil,nil,0,err}
			held=append(held,up50aSample{roles:roles,program:primitive,target:target,category:"unseen_basis_primitive"})
		}

		swap12:=[]UP50AInstruction{{Kind:"role_swap",Position:1}}
		target,err:=up50aTarget(roles,swap12)
		if err!=nil{return nil,nil,0,err}
		held=append(held,up50aSample{roles:roles,program:swap12,target:target,category:"unseen_swap12"})

		for start:=0; start<3; start++ {
			role1:=[]UP50AInstruction{
				{Kind:"role_swap",Position:0},
				{Kind:"value_swap",Start:start},
				{Kind:"role_swap",Position:0},
			}
			target1,err:=up50aTarget(roles,role1)
			if err!=nil{return nil,nil,0,err}
			held=append(held,up50aSample{roles:roles,program:role1,target:target1,category:"derived_role1"})

			role2:=[]UP50AInstruction{
				{Kind:"role_swap",Position:1},
				{Kind:"role_swap",Position:0},
				{Kind:"value_swap",Start:start},
				{Kind:"role_swap",Position:0},
				{Kind:"role_swap",Position:1},
			}
			target2,err:=up50aTarget(roles,role2)
			if err!=nil{return nil,nil,0,err}
			held=append(held,up50aSample{roles:roles,program:role2,target:target2,category:"derived_role2"})
		}
	}

	for _,length:=range lengths {
		for _,roles:=range up50aAllRoles() {
			if selected[up52aKey(roles)] { continue }
			for seed:=0; seed<3; seed++ {
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
				if err!=nil{return nil,nil,0,err}
				held=append(held,up50aSample{roles:roles,program:program,target:target,category:"long_program"})
			}
		}
	}
	return train,held,27-trainingBasis,nil
}

func RunUP52A() (UP52ACoverageLadderResult,error) {
	const (
		steps=140
		lr=0.10
		eps=1e-6
	)
	lengths:=[]int{12,36}
	levels:=[]int{9,6,3,1}
	result:=UP52ACoverageLadderResult{
		Schema:UP52ACoverageLadderSchema,
		Experiment:"UP-52A-basis-coverage-ladder",
		SourceUP51ASeal:"d7eef7212ad0e309246fcfc85bfd5c255e97e034",
		TrainingBasisOrder:up52aBasisOrder(),
		CoverageLevels:append([]int(nil),levels...),
		TrainSingleStepOnly:true,
		HeldOutLengths:append([]int(nil),lengths...),
		Steps:steps,
		LearningRate:lr,
		GradientEpsilon:eps,
	}
	for _,level:=range levels {
		train,held,heldBasis,err:=up52aSamples(level,lengths)
		if err!=nil{return UP52ACoverageLadderResult{},err}
		unitaryResult,err:=up50aTrain("unitary_coverage",train,held,up50aApplyUnitary,true,steps,lr,eps)
		if err!=nil{return UP52ACoverageLadderResult{},err}
		nonunitaryResult,err:=up50aTrain("nonunitary_coverage",train,held,up50aApplyNonUnitary,false,steps,lr,eps)
		if err!=nil{return UP52ACoverageLadderResult{},err}

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

		result.Points=append(result.Points,UP52ACoveragePoint{
			TrainingBasisStates:level,
			HeldoutBasisStates:heldBasis,
			TrainingCoverageFraction:float64(level)/27,
			Unitary:unitaryResult,
			NonUnitary:nonunitaryResult,
			UnseenPrimitiveAccuracy:primitive,
			UnseenSwap12Accuracy:swap12,
			DerivedRole1Accuracy:role1,
			DerivedRole2Accuracy:role2,
			LongProgramAccuracy:long,
			Gate:gate,
		})
		if gate && (result.MinimumPassingBasis==0 || level<result.MinimumPassingBasis) {
			result.MinimumPassingBasis=level
		}
	}
	return result,nil
}
