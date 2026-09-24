package unitary

import "math"

const UP58AExactAngleSchema="wingless.up58a-exact-angle-control.v1"

type UP58AAnglePoint struct{
 Length int `json:"length"`
 LearnedAccuracy float64 `json:"learned_accuracy"`
 ExactAccuracy float64 `json:"exact_accuracy"`
 LearnedNormDrift float64 `json:"learned_norm_drift"`
 ExactNormDrift float64 `json:"exact_norm_drift"`
}
type UP58AExactAngleResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP57ASeal string `json:"source_up57a_seal"`
 TrainingChanged bool `json:"training_changed"`
 LearnedParameters []float64 `json:"learned_parameters"`
 ExactParameters []float64 `json:"exact_parameters"`
 Lengths []int `json:"lengths"`
 Points []UP58AAnglePoint `json:"points"`
 ExactRestoresAll bool `json:"exact_restores_all"`
}
func RunUP58A()(UP58AExactAngleResult,error){
 lengths:=[]int{16,24,32,48,64,96,128}
 train,err:=up56aTrainSamples();if err!=nil{return UP58AExactAngleResult{},err}
 refHeld,err:=up56aHeldSamples([]int{16});if err!=nil{return UP58AExactAngleResult{},err}
 learned,err:=up56aTrain("unitary_exact_angle_source",train,refHeld,up56aApplyUnitary,true);if err!=nil{return UP58AExactAngleResult{},err}
 exact:=[]float64{math.Pi/2,math.Pi/2}
	result:=UP58AExactAngleResult{
		Schema:UP58AExactAngleSchema,Experiment:"UP-58A-exact-angle-control",
		SourceUP57ASeal:"cba417b76999151cc7e70a2da1d5afc577b26dfc",TrainingChanged:false,
		LearnedParameters:append([]float64(nil),learned.Parameters...),ExactParameters:append([]float64(nil),exact...),
		Lengths:append([]int(nil),lengths...),ExactRestoresAll:true,
	}
	for _,length:=range lengths{
		held,err:=up56aHeldSamples([]int{length});if err!=nil{return UP58AExactAngleResult{},err}
		long:=up57aLongOnly(held)
		_,la,ln,_,_,err:=up56aEvaluate(learned.Parameters,long,up56aApplyUnitary,true);if err!=nil{return UP58AExactAngleResult{},err}
		_,ea,en,_,_,err:=up56aEvaluate(exact,long,up56aApplyUnitary,true);if err!=nil{return UP58AExactAngleResult{},err}
		if ea<0.98{result.ExactRestoresAll=false}
		result.Points=append(result.Points,UP58AAnglePoint{Length:length,LearnedAccuracy:la,ExactAccuracy:ea,LearnedNormDrift:ln,ExactNormDrift:en})
	}
	return result,nil
}
