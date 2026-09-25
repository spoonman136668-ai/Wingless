package unitary

const UP81ASixRoleCoverageTransferSchema = "wingless.up81a-six-role-coverage-transfer.v1"

type up81aRoles [6]int

type UP81ASixRolePoint struct {
	TrainingStates      int               `json:"training_states"`
	TrainingSteps       int               `json:"training_steps"`
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP81ASixRoleCoverageTransferResult struct {
	Schema          string               `json:"schema"`
	Experiment      string               `json:"experiment"`
	SourceUP80ASeal string               `json:"source_up80a_seal"`
	Roles           int                  `json:"roles"`
	ValuesPerRole   int                  `json:"values_per_role"`
	JointStates     int                  `json:"joint_states"`
	TrainingSteps   int                  `json:"training_steps"`
	LearningRate    float64              `json:"learning_rate"`
	HeldOutStates   int                  `json:"heldout_states"`
	Points          []UP81ASixRolePoint  `json:"points"`
}

func up81aAllRoles() []up81aRoles {
	out := make([]up81aRoles, 0, 729)
	for a:=0;a<3;a++ { for b:=0;b<3;b++ { for c:=0;c<3;c++ { for d:=0;d<3;d++ { for e:=0;e<3;e++ { for f:=0;f<3;f++ {
		out=append(out,up81aRoles{a,b,c,d,e,f})
	}}}}}}
	return out
}

func up81aPrimary(r up81aRoles) int {
	s:=0
	for _,v:=range r { s+=v }
	return s%3
}

func up81aSecondary(r up81aRoles) int {
	return (r[0]+2*r[1]+r[2]+2*r[3]+r[4]+2*r[5])%3
}

func up81aTertiary(r up81aRoles) int {
	return (r[0]+r[1]+2*r[2]+2*r[3]+r[4])%3
}

func up81aTrain(r up81aRoles) bool {
	if up81aPrimary(r)!=0 { return false }
	s:=up81aSecondary(r)
	t:=up81aTertiary(r)
	return s==0 || (s==1 && t==0)
}

func up81aRawOneHot(r up81aRoles) []float64 {
	out:=make([]float64,18)
	for role,v:=range r { out[role*3+v]=1 }
	return out
}

func up81aSimplex(r up81aRoles) []float64 {
	out:=make([]float64,12)
	table:=[3][2]float64{{1,0},{-0.5,0.8660254037844386},{-0.5,-0.8660254037844386}}
	for role,v:=range r {
		out[role*2]=table[v][0]
		out[role*2+1]=table[v][1]
	}
	return out
}

func up81aRunPoint(name string, dim int, featureFn func(up81aRoles) []float64) (UP81ASixRolePoint,error) {
	all:=up81aAllRoles()
	var trainRoles,heldRoles []up81aRoles
	for _,r:=range all {
		if up81aTrain(r) { trainRoles=append(trainRoles,r) }
		if up81aPrimary(r)!=0 { heldRoles=append(heldRoles,r) }
	}
	point:=UP81ASixRolePoint{
		TrainingStates:len(trainRoles),TrainingSteps:1,Representation:name,
		FeatureDimension:dim,HeldOutStates:len(heldRoles),Gate:true,
	}
	for role:=0;role<6;role++ {
		train:=make([]headSample,0,len(trainRoles))
		held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles { train=append(train,headSample{features:featureFn(r),target:r[role]}) }
		for _,r:=range heldRoles { held=append(held,headSample{features:featureFn(r),target:r[role]}) }
		head,metric,err:=trainLinearSoftmax(train,3,dim,1,1.0)
		if err!=nil{return UP81ASixRolePoint{},err}
		_,heldAcc,err:=evaluateHead(head,held)
		if err!=nil{return UP81ASixRolePoint{},err}
		point.PerRole=append(point.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
		point.MeanTrainAccuracy+=metric.TrainAccuracy
		point.MeanHeldOutAccuracy+=heldAcc
		if heldAcc<0.98 { point.Gate=false }
	}
	point.MeanTrainAccuracy/=6
	point.MeanHeldOutAccuracy/=6
	return point,nil
}

func RunUP81A()(UP81ASixRoleCoverageTransferResult,error){
	result:=UP81ASixRoleCoverageTransferResult{
		Schema:UP81ASixRoleCoverageTransferSchema,
		Experiment:"UP-81A-six-role-coverage-transfer",
		SourceUP80ASeal:"83ed6b29dc439653505e2b0dbd9239f2fa0373fe",
		Roles:6,ValuesPerRole:3,JointStates:729,TrainingSteps:1,LearningRate:1.0,HeldOutStates:486,
	}
	onehot,err:=up81aRunPoint("raw_factorized_one_hot",18,up81aRawOneHot)
	if err!=nil{return UP81ASixRoleCoverageTransferResult{},err}
	result.Points=append(result.Points,onehot)
	simplex,err:=up81aRunPoint("role_factorized_simplex",12,up81aSimplex)
	if err!=nil{return UP81ASixRoleCoverageTransferResult{},err}
	result.Points=append(result.Points,simplex)
	return result,nil
}
