package unitary

import (
	"fmt"
	"math"
)

const UP64AFiveRoleZeroShotSchema = "wingless.up64a-five-role-zero-shot.v1"

type UP64AInstruction struct {
	Kind     string `json:"kind"`
	Start    int    `json:"start,omitempty"`
	Position int    `json:"position,omitempty"`
}

type up64aRoles [5]int

type up64aSample struct {
	Roles    up64aRoles
	Program  []UP64AInstruction
	Target   up64aRoles
	Category string
}

type UP64ACategoryMetric struct {
	Category string  `json:"category"`
	Samples  int     `json:"samples"`
	Accuracy float64 `json:"accuracy"`
}

type UP64AFiveRoleZeroShotResult struct {
	Schema                 string                `json:"schema"`
	Experiment             string                `json:"experiment"`
	SourceUP63ASeal        string                `json:"source_up63a_seal"`
	SourceTrainingTask     string                `json:"source_training_task"`
	SourceTrainingSteps    int                   `json:"source_training_steps"`
	SourceParameters       []float64             `json:"source_parameters"`
	FiveRoleTrainingUsed   bool                  `json:"five_role_training_used"`
	Roles                  int                   `json:"roles"`
	ValuesPerRole          int                   `json:"values_per_role"`
	JointDimension         int                   `json:"joint_dimension"`
	HeldOutLengths         []int                 `json:"heldout_lengths"`
	AggregateAccuracy      float64               `json:"aggregate_accuracy"`
	CategoryMetrics        []UP64ACategoryMetric `json:"category_metrics"`
	MaxNormDrift           float64               `json:"max_norm_drift"`
	MaxRoundTripError      float64               `json:"max_round_trip_error"`
	MatchedControlAccuracy float64               `json:"matched_control_accuracy"`
	Gate                   bool                  `json:"gate"`
}

func up64aAllRoles() []up64aRoles {
	out := make([]up64aRoles, 0, 243)
	for a:=0;a<3;a++ { for b:=0;b<3;b++ { for c:=0;c<3;c++ { for d:=0;d<3;d++ { for e:=0;e<3;e++ {
		out=append(out,up64aRoles{a,b,c,d,e})
	}}}}}
	return out
}

func up64aIndex(r up64aRoles)(int,error){
	idx:=0
	for _,v:=range r {
		if v<0||v>=3{return 0,fmt.Errorf("UP64A role value out of range")}
		idx=idx*3+v
	}
	return idx,nil
}

func up64aDecode(idx int)(up64aRoles,error){
	if idx<0||idx>=243{return up64aRoles{},fmt.Errorf("UP64A joint index out of range")}
	var r up64aRoles
	for i:=4;i>=0;i-- { r[i]=idx%3; idx/=3 }
	return r,nil
}

func up64aBasis(r up64aRoles)(State,error){
	idx,err:=up64aIndex(r);if err!=nil{return nil,err}
	s:=make(State,243);s[idx]=1
	return s,nil
}

func up64aApplyLogical(r up64aRoles,in UP64AInstruction)(up64aRoles,error){
	switch in.Kind {
	case "value_swap":
		if in.Start<0||in.Start>=3{return up64aRoles{},fmt.Errorf("UP64A value start")}
		a:=in.Start;b:=(a+1)%3
		switch r[0] { case a:r[0]=b; case b:r[0]=a }
	case "role_swap":
		if in.Position<0||in.Position>3{return up64aRoles{},fmt.Errorf("UP64A role position")}
		p:=in.Position;r[p],r[p+1]=r[p+1],r[p]
	default:
		return up64aRoles{},fmt.Errorf("UP64A instruction kind")
	}
	return r,nil
}

func up64aTarget(r up64aRoles,p []UP64AInstruction)(up64aRoles,error){
	var err error
	for _,in:=range p { r,err=up64aApplyLogical(r,in);if err!=nil{return up64aRoles{},err} }
	return r,nil
}

func up64aDerivedMutation(role,start int)[]UP64AInstruction{
	var p []UP64AInstruction
	for current:=role-1;current>=0;current-- { p=append(p,UP64AInstruction{Kind:"role_swap",Position:current}) }
	p=append(p,UP64AInstruction{Kind:"value_swap",Start:start})
	for current:=0;current<role;current++ { p=append(p,UP64AInstruction{Kind:"role_swap",Position:current}) }
	return p
}

func up64aSamples(lengths []int)([]up64aSample,[]up64aSample,error){
	var short,long []up64aSample
	for _,r:=range up64aAllRoles(){
		for start:=0;start<3;start++ {
			p:=[]UP64AInstruction{{Kind:"value_swap",Start:start}}
			t,err:=up64aTarget(r,p);if err!=nil{return nil,nil,err}
			short=append(short,up64aSample{r,p,t,"role0_value"})
		}
		for pos:=0;pos<4;pos++ {
			p:=[]UP64AInstruction{{Kind:"role_swap",Position:pos}}
			t,err:=up64aTarget(r,p);if err!=nil{return nil,nil,err}
			short=append(short,up64aSample{r,p,t,fmt.Sprintf("swap_%d%d",pos,pos+1)})
		}
		for role:=1;role<=4;role++ {
			for start:=0;start<3;start++ {
				p:=up64aDerivedMutation(role,start)
				t,err:=up64aTarget(r,p);if err!=nil{return nil,nil,err}
				short=append(short,up64aSample{r,p,t,fmt.Sprintf("derived_role%d",role)})
			}
		}
	}
	for _,length:=range lengths {
		for _,r:=range up64aAllRoles(){
			for seed:=0;seed<2;seed++ {
				p:=make([]UP64AInstruction,length)
				for i:=range p {
					switch (i+seed+r[0]+2*r[1]+3*r[2]+4*r[3]+5*r[4])%6 {
					case 0:p[i]=UP64AInstruction{Kind:"role_swap",Position:0}
					case 1:p[i]=UP64AInstruction{Kind:"role_swap",Position:1}
					case 2:p[i]=UP64AInstruction{Kind:"role_swap",Position:2}
					case 3:p[i]=UP64AInstruction{Kind:"role_swap",Position:3}
					default:p[i]=UP64AInstruction{Kind:"value_swap",Start:(i*i+seed+r[2]+r[4])%3}
					}
				}
				t,err:=up64aTarget(r,p);if err!=nil{return nil,nil,err}
				long=append(long,up64aSample{r,p,t,fmt.Sprintf("long_program_%d",length)})
			}
		}
	}
	return short,long,nil
}

func up64aCouplings(in UP64AInstruction,sourceParams []float64)([]Coupling,error){
	if len(sourceParams)!=2{return nil,fmt.Errorf("UP64A source parameter count")}
	roleTheta:=sourceParams[0]
	valueTheta:=sourceParams[1]
	switch in.Kind {
	case "value_swap":
		if in.Start<0||in.Start>=3{return nil,fmt.Errorf("UP64A value start")}
		a:=in.Start;b:=(a+1)%3
		var out []Coupling
		for _,r:=range up64aAllRoles(){
			if r[0]!=a{continue}
			s:=r;s[0]=b
			ia,_:=up64aIndex(r);ib,_:=up64aIndex(s)
			out=append(out,Coupling{A:ia,B:ib,Theta:valueTheta})
		}
		return out,nil
	case "role_swap":
		if in.Position<0||in.Position>3{return nil,fmt.Errorf("UP64A role position")}
		var out []Coupling
		p:=in.Position
		for _,r:=range up64aAllRoles(){
			if r[p]>=r[p+1]{continue}
			s:=r;s[p],s[p+1]=s[p+1],s[p]
			ia,_:=up64aIndex(r);ib,_:=up64aIndex(s)
			out=append(out,Coupling{A:ia,B:ib,Theta:roleTheta})
		}
		return out,nil
	default:return nil,fmt.Errorf("UP64A instruction kind")
	}
}

type up64aApply func(State,[]float64,[]UP64AInstruction)(State,error)

func up64aApplyUnitary(initial State,params []float64,p []UP64AInstruction)(State,error){
	out:=append(State(nil),initial...)
	for _,in:=range p { cs,err:=up64aCouplings(in,params);if err!=nil{return nil,err};out,err=Propagate(out,cs);if err!=nil{return nil,err} }
	return out,nil
}

func up64aApplyNonUnitary(initial State,params []float64,p []UP64AInstruction)(State,error){
	out:=append(State(nil),initial...)
	for _,in:=range p {
		cs,err:=up64aCouplings(in,params);if err!=nil{return nil,err}
		for _,c:=range cs { g:=c.Theta;a:=out[c.A];b:=out[c.B];out[c.A]=a-complex(g,0)*b;out[c.B]=complex(g,0)*a+b }
	}
	return out,nil
}

func up64aInverse(state State,params []float64,p []UP64AInstruction)(State,error){
	out:=append(State(nil),state...)
	for i:=len(p)-1;i>=0;i-- {
		cs,err:=up64aCouplings(p[i],params);if err!=nil{return nil,err}
		for j:=range cs { cs[j].Theta=-cs[j].Theta }
		out,err=Propagate(out,cs);if err!=nil{return nil,err}
	}
	return out,nil
}

func up64aEvaluate(params []float64,samples []up64aSample,apply up64aApply,roundTrip bool)(accuracy,maxNorm,maxRT float64,cats []UP64ACategoryMetric,err error){
	type count struct{n,h int}
	per:=map[string]*count{};order:=[]string{};hits:=0
	for _,sample:=range samples {
		initial,e:=up64aBasis(sample.Roles);if e!=nil{err=e;return}
		evolved,e:=apply(initial,params,sample.Program);if e!=nil{err=e;return}
		oi,e:=ArgMax(evolved);if e!=nil{err=e;return}
		obs,e:=up64aDecode(oi);if e!=nil{err=e;return}
		if _,ok:=per[sample.Category];!ok{per[sample.Category]=&count{};order=append(order,sample.Category)}
		per[sample.Category].n++
		if obs==sample.Target{hits++;per[sample.Category].h++}
		n2,e:=NormSquared(evolved);if e!=nil{err=e;return}
		if d:=math.Abs(n2-1);d>maxNorm{maxNorm=d}
		if roundTrip {
			rec,e:=up64aInverse(evolved,params,sample.Program);if e!=nil{err=e;return}
			d,e:=L2Distance(initial,rec);if e!=nil{err=e;return}
			if d>maxRT{maxRT=d}
		}
	}
	if len(samples)==0{err=fmt.Errorf("UP64A empty samples");return}
	accuracy=float64(hits)/float64(len(samples))
	for _,name:=range order { c:=per[name];cats=append(cats,UP64ACategoryMetric{name,c.n,float64(c.h)/float64(c.n)}) }
	return
}

func RunUP64A()(UP64AFiveRoleZeroShotResult,error){
	const sourceSteps=1440
	sourceTrain,err:=up56aTrainSamples();if err!=nil{return UP64AFiveRoleZeroShotResult{},err}
	sourceParams,err:=up61aTrainBudget(sourceTrain,sourceSteps);if err!=nil{return UP64AFiveRoleZeroShotResult{},err}
	lengths:=[]int{16,48}
	short,long,err:=up64aSamples(lengths);if err!=nil{return UP64AFiveRoleZeroShotResult{},err}
	all:=append(append([]up64aSample(nil),short...),long...)
	acc,norm,rt,cats,err:=up64aEvaluate(sourceParams,all,up64aApplyUnitary,true);if err!=nil{return UP64AFiveRoleZeroShotResult{},err}
	control,_,_,_,err:=up64aEvaluate(sourceParams,short,up64aApplyNonUnitary,false);if err!=nil{return UP64AFiveRoleZeroShotResult{},err}
	gate:=acc>=0.98&&norm<=1e-9&&rt<=1e-8
	for _,m:=range cats { if m.Accuracy<0.98 { gate=false } }
	return UP64AFiveRoleZeroShotResult{
		Schema:UP64AFiveRoleZeroShotSchema,
		Experiment:"UP-64A-five-role-zero-shot",
		SourceUP63ASeal:"5f05e9300bcf8ce6319fd01d3f5a0499759061ae",
		SourceTrainingTask:"UP-56A-three-role-double-controlled single-step primitives",
		SourceTrainingSteps:sourceSteps,
		SourceParameters:append([]float64(nil),sourceParams...),
		FiveRoleTrainingUsed:false,
		Roles:5,ValuesPerRole:3,JointDimension:243,
		HeldOutLengths:append([]int(nil),lengths...),
		AggregateAccuracy:acc,CategoryMetrics:cats,
		MaxNormDrift:norm,MaxRoundTripError:rt,
		MatchedControlAccuracy:control,Gate:gate,
	},nil
}
