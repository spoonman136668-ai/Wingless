package unitary

import "math"

const UP65ADistributedReadoutSchema = "wingless.up65a-distributed-readout.v1"

type UP65ARoleMetric struct {
	Role            int     `json:"role"`
	TrainAccuracy   float64 `json:"train_accuracy"`
	HeldOutAccuracy float64 `json:"heldout_accuracy"`
}

type UP65AArm struct {
	Name                string            `json:"name"`
	FeatureDimension    int               `json:"feature_dimension"`
	TrainStates         int               `json:"train_states"`
	HeldOutStates       int               `json:"heldout_states"`
	PerRole             []UP65ARoleMetric `json:"per_role"`
	MeanTrainAccuracy   float64           `json:"mean_train_accuracy"`
	MeanHeldOutAccuracy float64           `json:"mean_heldout_accuracy"`
	Gate                bool              `json:"gate"`
}

type UP65ADistributedReadoutResult struct {
	Schema            string   `json:"schema"`
	Experiment        string   `json:"experiment"`
	SourceUP64ASeal   string   `json:"source_up64a_seal"`
	Roles             int      `json:"roles"`
	ValuesPerRole     int      `json:"values_per_role"`
	LogicalStates     int      `json:"logical_states"`
	TrainingSplitRule string   `json:"training_split_rule"`
	ObserverOnly      bool     `json:"observer_only"`
	FactorizedDense   UP65AArm `json:"factorized_dense"`
	EntangledJoint    UP65AArm `json:"entangled_joint"`
}

func up65aTrainState(r up64aRoles) bool {
	sum:=0
	for _,v:=range r{sum+=v}
	return sum%3==0
}

func up65aHadamardFeatures(r up64aRoles) []float64 {
	base:=make([]float64,64)
	for role,value:=range r{
		base[role*3+value]=1/math.Sqrt(5)
	}
	out:=append([]float64(nil),base...)
	for width:=1;width<64;width*=2{
		for start:=0;start<64;start+=2*width{
			for j:=0;j<width;j++{
				a:=out[start+j]
				b:=out[start+j+width]
				out[start+j]=a+b
				out[start+j+width]=a-b
			}
		}
	}
	for i:=range out{out[i]/=8}
	return out
}

func up65aEntangledFeatures(r up64aRoles) []float64 {
	index,_:=up64aIndex(r)
	out:=make([]float64,64)
	scale:=1/math.Sqrt(32)
	for f:=1;f<=32;f++{
		phase:=2*math.Pi*float64(f*index)/243
		out[2*(f-1)]=scale*math.Cos(phase)
		out[2*(f-1)+1]=scale*math.Sin(phase)
	}
	return out
}

func up65aRunArm(name string,featureFn func(up64aRoles)[]float64)(UP65AArm,error){
	all:=up64aAllRoles()
	var trainRoles,heldRoles []up64aRoles
	for _,r:=range all{
		if up65aTrainState(r){trainRoles=append(trainRoles,r)}else{heldRoles=append(heldRoles,r)}
	}
	arm:=UP65AArm{
		Name:name,FeatureDimension:64,TrainStates:len(trainRoles),HeldOutStates:len(heldRoles),Gate:true,
	}
	for role:=0;role<5;role++{
		train:=make([]headSample,0,len(trainRoles))
		held:=make([]headSample,0,len(heldRoles))
		for _,r:=range trainRoles{train=append(train,headSample{features:featureFn(r),target:r[role]})}
		for _,r:=range heldRoles{held=append(held,headSample{features:featureFn(r),target:r[role]})}
		head,metric,err:=trainLinearSoftmax(train,3,64,800,1.0)
		if err!=nil{return UP65AArm{},err}
		_,heldAcc,err:=evaluateHead(head,held)
		if err!=nil{return UP65AArm{},err}
		arm.PerRole=append(arm.PerRole,UP65ARoleMetric{Role:role,TrainAccuracy:metric.TrainAccuracy,HeldOutAccuracy:heldAcc})
		arm.MeanTrainAccuracy+=metric.TrainAccuracy
		arm.MeanHeldOutAccuracy+=heldAcc
		if heldAcc<0.98{arm.Gate=false}
	}
	arm.MeanTrainAccuracy/=5
	arm.MeanHeldOutAccuracy/=5
	return arm,nil
}

func RunUP65A()(UP65ADistributedReadoutResult,error){
	factorized,err:=up65aRunArm("dense_hadamard_factorized",up65aHadamardFeatures)
	if err!=nil{return UP65ADistributedReadoutResult{},err}
	entangled,err:=up65aRunArm("joint_fourier_entangled",up65aEntangledFeatures)
	if err!=nil{return UP65ADistributedReadoutResult{},err}
	return UP65ADistributedReadoutResult{
		Schema:UP65ADistributedReadoutSchema,
		Experiment:"UP-65A-distributed-readout",
		SourceUP64ASeal:"4bc2a1ec7418b1be091fe54d6ad0fc952f46325c",
		Roles:5,ValuesPerRole:3,LogicalStates:243,
		TrainingSplitRule:"sum(role values) mod 3 == 0",
		ObserverOnly:true,
		FactorizedDense:factorized,
		EntangledJoint:entangled,
	},nil
}
