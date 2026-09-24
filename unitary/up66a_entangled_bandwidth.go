package unitary

import "math"

const UP66AEntangledBandwidthSchema = "wingless.up66a-entangled-bandwidth.v1"

type UP66ABandwidthPoint struct {
	FrequencyPairs       int               `json:"frequency_pairs"`
	FeatureDimension     int               `json:"feature_dimension"`
	TrainStates          int               `json:"train_states"`
	HeldOutStates        int               `json:"heldout_states"`
	PerRole              []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy    float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy  float64           `json:"mean_heldout_accuracy"`
	Gate                 bool              `json:"gate"`
}

type UP66AEntangledBandwidthResult struct {
	Schema          string                 `json:"schema"`
	Experiment      string                 `json:"experiment"`
	SourceUP65ASeal string                 `json:"source_up65a_seal"`
	TrainingSplitRule string               `json:"training_split_rule"`
	FrequencyPairs  []int                  `json:"frequency_pairs"`
	ObserverClassChanged bool              `json:"observer_class_changed"`
	Points          []UP66ABandwidthPoint  `json:"points"`
}

func up66aEntangledFeatures(r up64aRoles,pairs int) []float64 {
	index,_:=up64aIndex(r)
	out:=make([]float64,2*pairs)
	scale:=1/math.Sqrt(float64(pairs))
	for f:=1;f<=pairs;f++ {
		phase:=2*math.Pi*float64(f*index)/243
		out[2*(f-1)]=scale*math.Cos(phase)
		out[2*(f-1)+1]=scale*math.Sin(phase)
	}
	return out
}

func up66aRunPoint(pairs int)(UP66ABandwidthPoint,error){
	all:=up64aAllRoles()
	var trainRoles,heldRoles []up64aRoles
	for _,r:=range all{
		if up65aTrainState(r){trainRoles=append(trainRoles,r)}else{heldRoles=append(heldRoles,r)}
	}
	point:=UP66ABandwidthPoint{
		FrequencyPairs:pairs,FeatureDimension:2*pairs,
		TrainStates:len(trainRoles),HeldOutStates:len(heldRoles),Gate:true,
	}
	for role:=0;role<5;role++{
		train:=make([]headSample,0,len(trainRoles))
		held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles{train=append(train,headSample{features:up66aEntangledFeatures(r,pairs),target:r[role]})}
		for _,r:=range heldRoles{held=append(held,headSample{features:up66aEntangledFeatures(r,pairs),target:r[role]})}
		head,metric,err:=trainLinearSoftmax(train,3,2*pairs,800,1.0)
		if err!=nil{return UP66ABandwidthPoint{},err}
		_,heldAcc,err:=evaluateHead(head,held)
		if err!=nil{return UP66ABandwidthPoint{},err}
		point.PerRole=append(point.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
		point.MeanTrainAccuracy+=metric.TrainAccuracy
		point.MeanHeldOutAccuracy+=heldAcc
		if heldAcc<0.98{point.Gate=false}
	}
	point.MeanTrainAccuracy/=5
	point.MeanHeldOutAccuracy/=5
	return point,nil
}

func RunUP66A()(UP66AEntangledBandwidthResult,error){
	levels:=[]int{16,32,64,121}
	result:=UP66AEntangledBandwidthResult{
		Schema:UP66AEntangledBandwidthSchema,
		Experiment:"UP-66A-entangled-bandwidth",
		SourceUP65ASeal:"bd76bbe20261cccb039c319ae36b0bc2ee34c6b9",
		TrainingSplitRule:"sum(role values) mod 3 == 0",
		FrequencyPairs:append([]int(nil),levels...),
		ObserverClassChanged:false,
	}
	for _,pairs:=range levels{
		p,err:=up66aRunPoint(pairs)
		if err!=nil{return UP66AEntangledBandwidthResult{},err}
		result.Points=append(result.Points,p)
	}
	return result,nil
}
