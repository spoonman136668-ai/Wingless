package unitary

const UP82ASixRoleCoverageBoundarySchema = "wingless.up82a-six-role-coverage-boundary.v1"

type UP82ASixRoleCoverageBoundaryResult struct {
	Schema          string               `json:"schema"`
	Experiment      string               `json:"experiment"`
	SourceUP81ASeal string               `json:"source_up81a_seal"`
	TrainingLevels  []int                `json:"training_levels"`
	TrainingSteps   int                  `json:"training_steps"`
	LearningRate    float64              `json:"learning_rate"`
	HeldOutStates   int                  `json:"heldout_states"`
	Points          []UP81ASixRolePoint  `json:"points"`
}

func up82aTrainAtLevel(r up81aRoles, level int) bool {
	if up81aPrimary(r) != 0 { return false }
	s := up81aSecondary(r)
	t := up81aTertiary(r)
	switch level {
	case 54:
		return s == 0 && t < 2
	case 81:
		return s == 0
	case 108:
		return s == 0 || (s == 1 && t == 0)
	default:
		return false
	}
}

func up82aRunPoint(level int, name string, dim int, featureFn func(up81aRoles) []float64) (UP81ASixRolePoint,error) {
	all:=up81aAllRoles()
	var trainRoles,heldRoles []up81aRoles
	for _,r:=range all {
		if up82aTrainAtLevel(r,level) { trainRoles=append(trainRoles,r) }
		if up81aPrimary(r)!=0 { heldRoles=append(heldRoles,r) }
	}
	point:=UP81ASixRolePoint{TrainingStates:len(trainRoles),TrainingSteps:1,Representation:name,FeatureDimension:dim,HeldOutStates:len(heldRoles),Gate:true}
	for role:=0;role<6;role++ {
		train:=make([]headSample,0,len(trainRoles));held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles { train=append(train,headSample{features:featureFn(r),target:r[role]}) }
		for _,r:=range heldRoles { held=append(held,headSample{features:featureFn(r),target:r[role]}) }
		head,metric,err:=trainLinearSoftmax(train,3,dim,1,1.0);if err!=nil{return UP81ASixRolePoint{},err}
		_,heldAcc,err:=evaluateHead(head,held);if err!=nil{return UP81ASixRolePoint{},err}
		point.PerRole=append(point.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
		point.MeanTrainAccuracy+=metric.TrainAccuracy;point.MeanHeldOutAccuracy+=heldAcc
		if heldAcc<0.98 { point.Gate=false }
	}
	point.MeanTrainAccuracy/=6;point.MeanHeldOutAccuracy/=6
	return point,nil
}

func RunUP82A()(UP82ASixRoleCoverageBoundaryResult,error){
	levels:=[]int{54,81,108}
	result:=UP82ASixRoleCoverageBoundaryResult{Schema:UP82ASixRoleCoverageBoundarySchema,Experiment:"UP-82A-six-role-coverage-boundary",SourceUP81ASeal:"8acef3abd83d0a40fe8ecbd2dbb5627bb3bace91",TrainingLevels:append([]int(nil),levels...),TrainingSteps:1,LearningRate:1.0,HeldOutStates:486}
	for _,level:=range levels {
		onehot,err:=up82aRunPoint(level,"raw_factorized_one_hot",18,up81aRawOneHot);if err!=nil{return UP82ASixRoleCoverageBoundaryResult{},err};result.Points=append(result.Points,onehot)
		simplex,err:=up82aRunPoint(level,"role_factorized_simplex",12,up81aSimplex);if err!=nil{return UP82ASixRoleCoverageBoundaryResult{},err};result.Points=append(result.Points,simplex)
	}
	return result,nil
}
