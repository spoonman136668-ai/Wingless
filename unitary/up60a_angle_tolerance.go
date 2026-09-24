package unitary

import "math"

const UP60AToleranceSchema = "wingless.up60a-double-angle-tolerance.v1"

type UP60ATolerancePoint struct {
	DoubleControlDelta float64 `json:"double_control_delta"`
	Length int `json:"length"`
	Accuracy float64 `json:"accuracy"`
	MaxNormDrift float64 `json:"max_norm_drift"`
	MaxRoundTripError float64 `json:"max_round_trip_error"`
	Gate bool `json:"gate"`
}
type UP60AToleranceResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP59ASeal string `json:"source_up59a_seal"`
	TrainingChanged bool `json:"training_changed"`
	RoleSwapAngle float64 `json:"role_swap_angle"`
	ExactDoubleAngle float64 `json:"exact_double_angle"`
	Deltas []float64 `json:"deltas"`
	Lengths []int `json:"lengths"`
	Points []UP60ATolerancePoint `json:"points"`
}

func RunUP60A()(UP60AToleranceResult,error){
	train,err:=up56aTrainSamples();if err!=nil{return UP60AToleranceResult{},err}
	refHeld,err:=up56aHeldSamples([]int{16});if err!=nil{return UP60AToleranceResult{},err}
	learned,err:=up56aTrain("unitary_tolerance_source",train,refHeld,up56aApplyUnitary,true);if err!=nil{return UP60AToleranceResult{},err}
	deltas:=[]float64{0,-0.001,-0.0025,-0.005,-0.01,-0.02,-0.04,-0.08,-0.12,-0.16}
	lengths:=[]int{24,48,128}
	result:=UP60AToleranceResult{
		Schema:UP60AToleranceSchema,Experiment:"UP-60A-double-angle-tolerance",
		SourceUP59ASeal:"1bc1ea41cb01f97f37c516f4745d77b31f112bf8",
		TrainingChanged:false,RoleSwapAngle:learned.Parameters[0],ExactDoubleAngle:math.Pi/2,
		Deltas:append([]float64(nil),deltas...),Lengths:append([]int(nil),lengths...),
	}
	for _,delta:=range deltas{
		params:=[]float64{learned.Parameters[0],math.Pi/2+delta}
		for _,length:=range lengths{
			held,err:=up56aHeldSamples([]int{length});if err!=nil{return UP60AToleranceResult{},err}
			long:=up57aLongOnly(held)
			_,acc,norm,rt,_,err:=up56aEvaluate(params,long,up56aApplyUnitary,true);if err!=nil{return UP60AToleranceResult{},err}
			result.Points=append(result.Points,UP60ATolerancePoint{DoubleControlDelta:delta,Length:length,Accuracy:acc,MaxNormDrift:norm,MaxRoundTripError:rt,Gate:acc>=0.98&&norm<=1e-9&&rt<=1e-8})
		}
	}
	return result,nil
}
