package unitary

const UP83ASixRoleLowCoverageSchema = "wingless.up83a-six-role-low-coverage.v1"

type UP83ASixRoleLowCoverageResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP82ASeal string `json:"source_up82a_seal"`
	TrainingLevels []int `json:"training_levels"`
	TrainingSteps int `json:"training_steps"`
	LearningRate float64 `json:"learning_rate"`
	HeldOutStates int `json:"heldout_states"`
	Points []UP81ASixRolePoint `json:"points"`
}

func up83aTrainAtLevel(r up81aRoles, level int) bool {
	if up81aPrimary(r)!=0{return false}
	s:=up81aSecondary(r);t:=up81aTertiary(r)
	switch level {
	case 27:return s==0&&t==0
	case 54:return s==0&&t<2
	default:return false
	}
}

func up83aRunPoint(level int,name string,dim int,featureFn func(up81aRoles)[]float64)(UP81ASixRolePoint,error){
	all:=up81aAllRoles();var trainRoles,heldRoles []up81aRoles
	for _,r:=range all{if up83aTrainAtLevel(r,level){trainRoles=append(trainRoles,r)};if up81aPrimary(r)!=0{heldRoles=append(heldRoles,r)}}
	point:=UP81ASixRolePoint{TrainingStates:len(trainRoles),TrainingSteps:1,Representation:name,FeatureDimension:dim,HeldOutStates:len(heldRoles),Gate:true}
	for role:=0;role<6;role++{
		train:=make([]headSample,0,len(trainRoles));held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles{train=append(train,headSample{features:featureFn(r),target:r[role]})}
		for _,r:=range heldRoles{held=append(held,headSample{features:featureFn(r),target:r[role]})}
		head,metric,err:=trainLinearSoftmax(train,3,dim,1,1.0);if err!=nil{return UP81ASixRolePoint{},err}
		_,heldAcc,err:=evaluateHead(head,held);if err!=nil{return UP81ASixRolePoint{},err}
		point.PerRole=append(point.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
		point.MeanTrainAccuracy+=metric.TrainAccuracy;point.MeanHeldOutAccuracy+=heldAcc;if heldAcc<0.98{point.Gate=false}
	}
	point.MeanTrainAccuracy/=6;point.MeanHeldOutAccuracy/=6;return point,nil
}

func RunUP83A()(UP83ASixRoleLowCoverageResult,error){
	levels:=[]int{27,54};result:=UP83ASixRoleLowCoverageResult{Schema:UP83ASixRoleLowCoverageSchema,Experiment:"UP-83A-six-role-low-coverage",SourceUP82ASeal:"fa2f6e41aa2871643f23556bde2ac4dfc875e5a0",TrainingLevels:append([]int(nil),levels...),TrainingSteps:1,LearningRate:1.0,HeldOutStates:486}
	for _,level:=range levels{
		onehot,err:=up83aRunPoint(level,"raw_factorized_one_hot",18,up81aRawOneHot);if err!=nil{return UP83ASixRoleLowCoverageResult{},err};result.Points=append(result.Points,onehot)
		simplex,err:=up83aRunPoint(level,"role_factorized_simplex",12,up81aSimplex);if err!=nil{return UP83ASixRoleLowCoverageResult{},err};result.Points=append(result.Points,simplex)
	}
	return result,nil
}
