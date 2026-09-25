package unitary

const UP88ADelayedRoleRecallSchema = "wingless.up88a-delayed-role-recall.v1"

type UP88ADelayedRoleRecallPoint struct {
	Representation       string  `json:"representation"`
	FeatureDimension     int     `json:"feature_dimension"`
	Delay                int     `json:"delay"`
	NoiseAmplitude       float64 `json:"noise_amplitude"`
	Scenarios            int     `json:"scenarios"`
	TargetRecallAccuracy float64 `json:"target_recall_accuracy"`
	ExactStateAccuracy   float64 `json:"exact_state_accuracy"`
	Gate                 bool    `json:"gate"`
}

type UP88ADelayedRoleRecallResult struct {
	Schema          string                         `json:"schema"`
	Experiment      string                         `json:"experiment"`
	SourceUP87ASeal string                         `json:"source_up87a_seal"`
	TrainingStates  int                            `json:"training_states"`
	TrainingSteps   int                            `json:"training_steps"`
	LearningRate    float64                        `json:"learning_rate"`
	Scenarios       int                            `json:"scenarios"`
	Delays          []int                          `json:"delays"`
	NoiseLevels     []float64                      `json:"noise_levels"`
	Retraining      bool                           `json:"retraining"`
	Points          []UP88ADelayedRoleRecallPoint  `json:"points"`
}

func up88aTrainHeads(dim int, featureFn func(up81aRoles)[]float64)([]linearSoftmax,error){
	all:=up81aAllRoles()
	var trainRoles []up81aRoles
	for _,r:=range all{
		if up81aPrimary(r)==0&&up81aSecondary(r)==0&&up81aTertiary(r)==0{
			trainRoles=append(trainRoles,r)
		}
	}
	heads:=make([]linearSoftmax,0,6)
	for role:=0;role<6;role++{
		train:=make([]headSample,0,len(trainRoles))
		for _,r:=range trainRoles{
			train=append(train,headSample{features:featureFn(r),target:r[role]})
		}
		head,_,err:=trainLinearSoftmax(train,3,dim,1,1.0)
		if err!=nil{return nil,err}
		heads=append(heads,head)
	}
	return heads,nil
}

func up88aDecode(heads []linearSoftmax,state []float64)(up81aRoles,error){
	var out up81aRoles
	for role:=0;role<6;role++{
		p,err:=heads[role].probabilities(state)
		if err!=nil{return up81aRoles{},err}
		got,_,err:=classAndMargin(p)
		if err!=nil{return up81aRoles{},err}
		out[role]=got
	}
	return out,nil
}

func up88aRunRepresentation(name string,dim,width int,featureFn func(up81aRoles)[]float64,delays []int,noises []float64)([]UP88ADelayedRoleRecallPoint,error){
	heads,err:=up88aTrainHeads(dim,featureFn)
	if err!=nil{return nil,err}
	const scenarios=64
	var points []UP88ADelayedRoleRecallPoint
	for _,delay:=range delays{
		for _,noise:=range noises{
			targetHits:=0
			exactHits:=0
			for scenario:=0;scenario<scenarios;scenario++{
				truth:=up87aScenarioRoles(scenario)
				targetRole:=scenario%6
				targetValue:=(scenario*2+1)%3
				truth[targetRole]=targetValue
				state:=append([]float64(nil),featureFn(truth)...)
				for step:=0;step<delay;step++{
					candidate:=(scenario+step*3)%5
					writeRole:=candidate
					if writeRole>=targetRole{writeRole++}
					writeValue:=(scenario+2*step+writeRole+1)%3
					truth[writeRole]=writeValue
					exact:=featureFn(truth)
					start:=writeRole*width
					copy(state[start:start+width],exact[start:start+width])
					up87aAddNoise(state,scenario,step,noise)
				}
				decoded,err:=up88aDecode(heads,state)
				if err!=nil{return nil,err}
				if decoded[targetRole]==truth[targetRole]{targetHits++}
				if decoded==truth{exactHits++}
			}
			targetAcc:=float64(targetHits)/scenarios
			exactAcc:=float64(exactHits)/scenarios
			points=append(points,UP88ADelayedRoleRecallPoint{
				Representation:name,FeatureDimension:dim,Delay:delay,NoiseAmplitude:noise,Scenarios:scenarios,
				TargetRecallAccuracy:targetAcc,ExactStateAccuracy:exactAcc,
				Gate:targetAcc>=0.99&&exactAcc>=0.95,
			})
		}
	}
	return points,nil
}

func RunUP88A()(UP88ADelayedRoleRecallResult,error){
	delays:=[]int{16,64,256,1024}
	noises:=[]float64{0,0.0005,0.001,0.002}
	result:=UP88ADelayedRoleRecallResult{
		Schema:UP88ADelayedRoleRecallSchema,
		Experiment:"UP-88A-delayed-role-recall",
		SourceUP87ASeal:"eda467a422930058f7ef380f02aac7a240240560",
		TrainingStates:27,TrainingSteps:1,LearningRate:1.0,Scenarios:64,
		Delays:append([]int(nil),delays...),NoiseLevels:append([]float64(nil),noises...),Retraining:false,
	}
	onehot,err:=up88aRunRepresentation("raw_factorized_one_hot",18,3,up81aRawOneHot,delays,noises)
	if err!=nil{return UP88ADelayedRoleRecallResult{},err}
	result.Points=append(result.Points,onehot...)
	simplex,err:=up88aRunRepresentation("role_factorized_simplex",12,2,up81aSimplex,delays,noises)
	if err!=nil{return UP88ADelayedRoleRecallResult{},err}
	result.Points=append(result.Points,simplex...)
	return result,nil
}
