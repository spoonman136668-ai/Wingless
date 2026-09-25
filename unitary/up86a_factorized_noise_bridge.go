package unitary

import "math"

const UP86AFactorizedNoiseBridgeSchema = "wingless.up86a-factorized-noise-bridge.v1"

type UP86ANoisePoint struct {
	Representation      string            `json:"representation"`
	FeatureDimension    int               `json:"feature_dimension"`
	NoiseAmplitude      float64           `json:"noise_amplitude"`
	TrainingStates      int               `json:"training_states"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP86AFactorizedNoiseBridgeResult struct {
	Schema          string           `json:"schema"`
	Experiment      string           `json:"experiment"`
	SourceUP85ASeal string           `json:"source_up85a_seal"`
	TrainingStates  int              `json:"training_states"`
	HeldOutStates   int              `json:"heldout_states"`
	TrainingSteps   int              `json:"training_steps"`
	LearningRate    float64          `json:"learning_rate"`
	NoiseLevels     []float64        `json:"noise_levels"`
	Retraining      bool             `json:"retraining"`
	Points          []UP86ANoisePoint `json:"points"`
}

func up86aPerturb(features []float64, sampleIndex int, amplitude float64) []float64 {
	out:=append([]float64(nil),features...)
	if amplitude==0{return out}
	for j:=range out{
		phase:=float64((sampleIndex+1)*(j+1)*17)
		out[j]+=amplitude*math.Sin(phase)
	}
	return out
}

func up86aRunRepresentation(name string,dim int,featureFn func(up81aRoles)[]float64,noises []float64)([]UP86ANoisePoint,error){
	all:=up81aAllRoles()
	var trainRoles,heldRoles []up81aRoles
	for _,r:=range all{
		if up81aPrimary(r)==0&&up81aSecondary(r)==0&&up81aTertiary(r)==0{trainRoles=append(trainRoles,r)}
		if up81aPrimary(r)!=0{heldRoles=append(heldRoles,r)}
	}
	points:=make([]UP86ANoisePoint,len(noises))
	for i,noise:=range noises{
		points[i]=UP86ANoisePoint{
			Representation:name,FeatureDimension:dim,NoiseAmplitude:noise,
			TrainingStates:len(trainRoles),HeldOutStates:len(heldRoles),Gate:true,
		}
	}
	for role:=0;role<6;role++{
		train:=make([]headSample,0,len(trainRoles))
		for _,r:=range trainRoles{train=append(train,headSample{features:featureFn(r),target:r[role]})}
		head,metric,err:=trainLinearSoftmax(train,3,dim,1,1.0)
		if err!=nil{return nil,err}
		for ni,noise:=range noises{
			held:=make([]headSample,0,len(heldRoles))
			for si,r:=range heldRoles{
				held=append(held,headSample{features:up86aPerturb(featureFn(r),si,noise),target:r[role]})
			}
			_,heldAcc,err:=evaluateHead(head,held)
			if err!=nil{return nil,err}
			points[ni].PerRole=append(points[ni].PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
			points[ni].MeanTrainAccuracy+=metric.TrainAccuracy
			points[ni].MeanHeldOutAccuracy+=heldAcc
			if heldAcc<0.98{points[ni].Gate=false}
		}
	}
	for i:=range points{
		points[i].MeanTrainAccuracy/=6
		points[i].MeanHeldOutAccuracy/=6
	}
	return points,nil
}

func RunUP86A()(UP86AFactorizedNoiseBridgeResult,error){
	noises:=[]float64{0,0.01,0.025,0.05,0.10}
	result:=UP86AFactorizedNoiseBridgeResult{
		Schema:UP86AFactorizedNoiseBridgeSchema,
		Experiment:"UP-86A-factorized-noise-bridge",
		SourceUP85ASeal:"02d0e76c50e756cab685aec2a2b51f99254663d3",
		TrainingStates:27,HeldOutStates:486,TrainingSteps:1,LearningRate:1.0,
		NoiseLevels:append([]float64(nil),noises...),Retraining:false,
	}
	onehot,err:=up86aRunRepresentation("raw_factorized_one_hot",18,up81aRawOneHot,noises)
	if err!=nil{return UP86AFactorizedNoiseBridgeResult{},err}
	result.Points=append(result.Points,onehot...)
	simplex,err:=up86aRunRepresentation("role_factorized_simplex",12,up81aSimplex,noises)
	if err!=nil{return UP86AFactorizedNoiseBridgeResult{},err}
	result.Points=append(result.Points,simplex...)
	return result,nil
}
