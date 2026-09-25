package unitary

import "math"

const UP87ASequentialRoleOverwriteSchema = "wingless.up87a-sequential-role-overwrite.v1"

type UP87ASequentialRoleOverwritePoint struct {
	Representation   string  `json:"representation"`
	FeatureDimension int     `json:"feature_dimension"`
	Length           int     `json:"length"`
	NoiseAmplitude   float64 `json:"noise_amplitude"`
	Scenarios        int     `json:"scenarios"`
	ValueAccuracy    float64 `json:"value_accuracy"`
	ExactStateAccuracy float64 `json:"exact_state_accuracy"`
	Gate             bool    `json:"gate"`
}

type UP87ASequentialRoleOverwriteResult struct {
	Schema          string                            `json:"schema"`
	Experiment      string                            `json:"experiment"`
	SourceUP86ASeal string                            `json:"source_up86a_seal"`
	TrainingStates  int                               `json:"training_states"`
	TrainingSteps   int                               `json:"training_steps"`
	LearningRate    float64                           `json:"learning_rate"`
	Scenarios       int                               `json:"scenarios"`
	Lengths         []int                             `json:"lengths"`
	NoiseLevels     []float64                         `json:"noise_levels"`
	Retraining      bool                              `json:"retraining"`
	Points          []UP87ASequentialRoleOverwritePoint `json:"points"`
}

func up87aScenarioRoles(scenario int) up81aRoles {
	n := (scenario*37 + 11) % 729
	var r up81aRoles
	for i:=5;i>=0;i--{
		r[i]=n%3
		n/=3
	}
	return r
}

func up87aAddNoise(state []float64,scenario,step int,amplitude float64){
	if amplitude==0{return}
	for j:=range state{
		phase:=float64((scenario+1)*(step+1)*(j+1)*13)
		state[j]+=amplitude*math.Sin(phase)
	}
}

func up87aRunRepresentation(name string,dim,width int,featureFn func(up81aRoles)[]float64,lengths []int,noises []float64)([]UP87ASequentialRoleOverwritePoint,error){
	all:=up81aAllRoles()
	var trainRoles []up81aRoles
	for _,r:=range all{
		if up81aPrimary(r)==0&&up81aSecondary(r)==0&&up81aTertiary(r)==0{
			trainRoles=append(trainRoles,r)
		}
	}
	const scenarios=64
	var points []UP87ASequentialRoleOverwritePoint
	for _,length:=range lengths{
		for _,noise:=range noises{
			eventErrors:=make([]bool,scenarios*length)
			valueCorrect:=0
			valueTotal:=0
			for role:=0;role<6;role++{
				train:=make([]headSample,0,len(trainRoles))
				for _,r:=range trainRoles{
					train=append(train,headSample{features:featureFn(r),target:r[role]})
				}
				head,_,err:=trainLinearSoftmax(train,3,dim,1,1.0)
				if err!=nil{return nil,err}
				for scenario:=0;scenario<scenarios;scenario++{
					truth:=up87aScenarioRoles(scenario)
					state:=append([]float64(nil),featureFn(truth)...)
					for step:=0;step<length;step++{
						writeRole:=(scenario*5+step*7)%6
						writeValue:=(scenario+2*step+writeRole)%3
						truth[writeRole]=writeValue
						exact:=featureFn(truth)
						start:=writeRole*width
						copy(state[start:start+width],exact[start:start+width])
						up87aAddNoise(state,scenario,step,noise)
						probabilities,err:=head.probabilities(state)
						if err!=nil{return nil,err}
						got,_,err:=classAndMargin(probabilities)
						if err!=nil{return nil,err}
						valueTotal++
						if got==truth[role]{
							valueCorrect++
						}else{
							eventErrors[scenario*length+step]=true
						}
					}
				}
			}
			exactCorrect:=0
			for _,bad:=range eventErrors{if !bad{exactCorrect++}}
			valueAccuracy:=float64(valueCorrect)/float64(valueTotal)
			exactAccuracy:=float64(exactCorrect)/float64(len(eventErrors))
			points=append(points,UP87ASequentialRoleOverwritePoint{
				Representation:name,FeatureDimension:dim,Length:length,NoiseAmplitude:noise,Scenarios:scenarios,
				ValueAccuracy:valueAccuracy,ExactStateAccuracy:exactAccuracy,
				Gate:valueAccuracy>=0.99&&exactAccuracy>=0.95,
			})
		}
	}
	return points,nil
}

func RunUP87A()(UP87ASequentialRoleOverwriteResult,error){
	lengths:=[]int{16,64,256}
	noises:=[]float64{0,0.0005,0.001,0.002}
	result:=UP87ASequentialRoleOverwriteResult{
		Schema:UP87ASequentialRoleOverwriteSchema,
		Experiment:"UP-87A-sequential-role-overwrite",
		SourceUP86ASeal:"c4b554ccacee9d0a7807603be070e6ed0d6857ae",
		TrainingStates:27,TrainingSteps:1,LearningRate:1.0,Scenarios:64,
		Lengths:append([]int(nil),lengths...),NoiseLevels:append([]float64(nil),noises...),Retraining:false,
	}
	onehot,err:=up87aRunRepresentation("raw_factorized_one_hot",18,3,up81aRawOneHot,lengths,noises)
	if err!=nil{return UP87ASequentialRoleOverwriteResult{},err}
	result.Points=append(result.Points,onehot...)
	simplex,err:=up87aRunRepresentation("role_factorized_simplex",12,2,up81aSimplex,lengths,noises)
	if err!=nil{return UP87ASequentialRoleOverwriteResult{},err}
	result.Points=append(result.Points,simplex...)
	return result,nil
}
