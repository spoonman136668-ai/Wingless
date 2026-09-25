package unitary

const UP85ASixRoleSelectorReplicationSchema = "wingless.up85a-six-role-selector-replication.v1"

type UP85ASixRoleSelectorReplicationResult struct {
	Schema          string              `json:"schema"`
	Experiment      string              `json:"experiment"`
	SourceUP84ASeal string              `json:"source_up84a_seal"`
	TrainingLevels  []int               `json:"training_levels"`
	TrainingSteps   int                 `json:"training_steps"`
	LearningRate    float64             `json:"learning_rate"`
	HeldOutStates   int                 `json:"heldout_states"`
	Points          []UP81ASixRolePoint `json:"points"`
}

func up85aQuaternaryB(r up81aRoles) int { return r[1] % 3 }

func up85aTrainAtLevel(r up81aRoles,level int) bool {
	if up81aPrimary(r)!=0 || up81aSecondary(r)!=0 || up81aTertiary(r)!=0{return false}
	q:=up85aQuaternaryB(r)
	switch level{case 18:return q<2;case 27:return true;default:return false}
}

func up85aRunPoint(level int,name string,dim int,featureFn func(up81aRoles)[]float64)(UP81ASixRolePoint,error){
	all:=up81aAllRoles();var trainRoles,heldRoles []up81aRoles
	for _,r:=range all{if up85aTrainAtLevel(r,level){trainRoles=append(trainRoles,r)};if up81aPrimary(r)!=0{heldRoles=append(heldRoles,r)}}
	point:=UP81ASixRolePoint{TrainingStates:len(trainRoles),TrainingSteps:1,Representation:name,FeatureDimension:dim,HeldOutStates:len(heldRoles),Gate:true}
	for role:=0;role<6;role++{
		train:=make([]headSample,0,len(trainRoles));held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles{train=append(train,headSample{features:featureFn(r),target:r[role]})}
		for _,r:=range heldRoles{held=append(held,headSample{features:featureFn(r),target:r[role]})}
		head,metric,err:=trainLinearSoftmax(train,3,dim,1,1.0);if err!=nil{return UP81ASixRolePoint{},err}
		_,heldAcc,err:=evaluateHead(head,held);if err!=nil{return UP81ASixRolePoint{},err}
		point.PerRole=append(point.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc});point.MeanTrainAccuracy+=metric.TrainAccuracy;point.MeanHeldOutAccuracy+=heldAcc;if heldAcc<0.98{point.Gate=false}
	}
	point.MeanTrainAccuracy/=6;point.MeanHeldOutAccuracy/=6;return point,nil
}

func RunUP85A()(UP85ASixRoleSelectorReplicationResult,error){
	levels:=[]int{18,27};result:=UP85ASixRoleSelectorReplicationResult{Schema:UP85ASixRoleSelectorReplicationSchema,Experiment:"UP-85A-six-role-selector-replication",SourceUP84ASeal:"fb50c21756f2a1f6a1cb472fd382c74760cb52e0",TrainingLevels:append([]int(nil),levels...),TrainingSteps:1,LearningRate:1.0,HeldOutStates:486}
	for _,level:=range levels{
		onehot,err:=up85aRunPoint(level,"raw_factorized_one_hot",18,up81aRawOneHot);if err!=nil{return UP85ASixRoleSelectorReplicationResult{},err};result.Points=append(result.Points,onehot)
		simplex,err:=up85aRunPoint(level,"role_factorized_simplex",12,up81aSimplex);if err!=nil{return UP85ASixRoleSelectorReplicationResult{},err};result.Points=append(result.Points,simplex)
	}
	return result,nil
}
