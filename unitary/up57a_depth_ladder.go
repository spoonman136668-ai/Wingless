package unitary

const UP57ADepthSchema = "wingless.up57a-double-control-depth-ladder.v1"

type UP57ADepthPoint struct {
	Length int `json:"length"`
	UnitaryAccuracy float64 `json:"unitary_accuracy"`
	NonUnitaryAccuracy float64 `json:"nonunitary_accuracy"`
	MaxNormDrift float64 `json:"max_norm_drift"`
	MaxRoundTripError float64 `json:"max_round_trip_error"`
	Gate bool `json:"gate"`
}
type UP57ADepthResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP56ASeal string `json:"source_up56a_seal"`
	TrainingChanged bool `json:"training_changed"`
	Lengths []int `json:"lengths"`
	UnitaryParameters []float64 `json:"unitary_parameters"`
	NonUnitaryParameters []float64 `json:"nonunitary_parameters"`
	Points []UP57ADepthPoint `json:"points"`
	LastPassingLength int `json:"last_passing_length"`
	FirstFailingLength int `json:"first_failing_length"`
}

func up57aLongOnly(samples []up56aSample) []up56aSample {
	var out []up56aSample
	for _,s:=range samples{if s.Category=="long_program"{out=append(out,s)}}
	return out
}

func RunUP57A()(UP57ADepthResult,error){
	lengths:=[]int{2,4,8,12,16,24,32,48,64,96,128}
	train,err:=up56aTrainSamples();if err!=nil{return UP57ADepthResult{},err}
	referenceHeld,err:=up56aHeldSamples([]int{16});if err!=nil{return UP57ADepthResult{},err}
	u,err:=up56aTrain("unitary_depth_source",train,referenceHeld,up56aApplyUnitary,true);if err!=nil{return UP57ADepthResult{},err}
	n,err:=up56aTrain("nonunitary_depth_source",train,referenceHeld,up56aApplyNonUnitary,false);if err!=nil{return UP57ADepthResult{},err}
	result:=UP57ADepthResult{
		Schema:UP57ADepthSchema,Experiment:"UP-57A-double-control-depth-ladder",
		SourceUP56ASeal:"fd1731dba6b0f8c73ef5173143c52de1ad435a98",
		TrainingChanged:false,Lengths:append([]int(nil),lengths...),
		UnitaryParameters:append([]float64(nil),u.Parameters...),NonUnitaryParameters:append([]float64(nil),n.Parameters...),
	}
	for _,length:=range lengths{
		held,err:=up56aHeldSamples([]int{length});if err!=nil{return UP57ADepthResult{},err}
		long:=up57aLongOnly(held)
		_,ua,un,urt,_,err:=up56aEvaluate(u.Parameters,long,up56aApplyUnitary,true);if err!=nil{return UP57ADepthResult{},err}
		_,na,_,_,_,err:=up56aEvaluate(n.Parameters,long,up56aApplyNonUnitary,false);if err!=nil{return UP57ADepthResult{},err}
		gate:=ua>=0.98&&un<=1e-9&&urt<=1e-8
		result.Points=append(result.Points,UP57ADepthPoint{Length:length,UnitaryAccuracy:ua,NonUnitaryAccuracy:na,MaxNormDrift:un,MaxRoundTripError:urt,Gate:gate})
		if gate{result.LastPassingLength=length}else if result.FirstFailingLength==0{result.FirstFailingLength=length}
	}
	return result,nil
}
