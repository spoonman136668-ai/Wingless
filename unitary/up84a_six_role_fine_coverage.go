package unitary

const UP84ASixRoleFineCoverageSchema = "wingless.up84a-six-role-fine-coverage.v1"

type UP84ASixRoleFineCoverageResult struct {
	Schema          string              `json:"schema"`
	Experiment      string              `json:"experiment"`
	SourceUP83ASeal string              `json:"source_up83a_seal"`
	TrainingLevels  []int               `json:"training_levels"`
	TrainingSteps   int                 `json:"training_steps"`
	LearningRate    float64             `json:"learning_rate"`
	HeldOutStates   int                 `json:"heldout_states"`
	Points          []UP81ASixRolePoint `json:"points"`
}

func up84aQuaternary(r up81aRoles) int { return r[0] % 3 }

func up84aTrainAtLevel(r up81aRoles, level int) bool {
	if up81aPrimary(r)!=0 || up81aSecondary(r)!=0 || up81aTertiary(r)!=0 { return false }
	q:=up84aQuaternary(r)
	switch level {
	case 9: return q==0
	case 18: return q<2
	case 27: return true
	default: return false
	}
}

func up84aRunPoint(level int,name string,dim int,featureFn func(up81aRoles)[]float64)(UP81ASixRolePoint,error){
	all:=up81aAllRoles();var trainRoles,heldRoles []up81aRoles
	for _,r:=range all{if up84aTrainAtLevel(r,level){trainRoles=append(trainRoles,r)};if up81aPrimary(r)!=0{heldRoles=append(heldRoles,r)}}
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

func RunUP84A()(UP84ASixRoleFineCoverageResult,error){
	levels:=[]int{9,18,27}
	result:=UP84ASixRoleFineCoverageResult{Schema:UP84ASixRoleFineCoverageSchema,Experiment:"UP-84A-six-role-fine-coverage",SourceUP83ASeal:"9167f6909c9102a82fe81753e4c337cc22cebe1f",TrainingLevels:append([]int(nil),levels...),TrainingSteps:1,LearningRate:1.0,HeldOutStates:486}
	for _,level:=range levels{
		onehot,err:=up84aRunPoint(level,"raw_factorized_one_hot",18,up81aRawOneHot);if err!=nil{return UP84ASixRoleFineCoverageResult{},err};result.Points=append(result.Points,onehot)
		simplex,err:=up84aRunPoint(level,"role_factorized_simplex",12,up81aSimplex);if err!=nil{return UP84ASixRoleFineCoverageResult{},err};result.Points=append(result.Points,simplex)
	}
	return result,nil
}
