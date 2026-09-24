package unitary

import (
	"fmt"
	"math"
)

const UP54AFourRoleSchema = "wingless.up54a-four-role-transfer.v1"

type UP54AInstruction struct {
	Kind     string `json:"kind"`
	Start    int    `json:"start,omitempty"`
	Position int    `json:"position,omitempty"`
}

type up54aRoles [4]int

type up54aSample struct {
	Roles    up54aRoles
	Program  []UP54AInstruction
	Target   up54aRoles
	Category string
}

type UP54ACategoryMetric struct {
	Category string  `json:"category"`
	Samples  int     `json:"samples"`
	Accuracy float64 `json:"accuracy"`
}

type UP54APathResult struct {
	Name                 string                `json:"name"`
	InitialTrainAccuracy float64               `json:"initial_train_accuracy"`
	TrainAccuracy        float64               `json:"train_accuracy"`
	HeldOutAccuracy      float64               `json:"heldout_accuracy"`
	CategoryMetrics      []UP54ACategoryMetric `json:"category_metrics"`
	InitialLoss          float64               `json:"initial_loss"`
	FinalLoss            float64               `json:"final_loss"`
	Parameters           []float64             `json:"parameters"`
	MaxNormDrift         float64               `json:"max_norm_drift"`
	MaxRoundTripError    float64               `json:"max_round_trip_error"`
}

type UP54AFourRoleResult struct {
	Schema                    string          `json:"schema"`
	Experiment                string          `json:"experiment"`
	SourceUP53ASeal           string          `json:"source_up53a_seal"`
	ValuesPerRole             int             `json:"values_per_role"`
	Roles                     int             `json:"roles"`
	JointDimension            int             `json:"joint_dimension"`
	TrainSingleStepOnly       bool            `json:"train_single_step_only"`
	OnlySwap01Trained         bool            `json:"only_swap01_trained"`
	HeldOutLengths            []int           `json:"heldout_lengths"`
	Steps                     int             `json:"steps"`
	LearningRate              float64         `json:"learning_rate"`
	GradientEpsilon           float64         `json:"gradient_epsilon"`
	Unitary                   UP54APathResult `json:"unitary"`
	NonUnitary                UP54APathResult `json:"nonunitary_matched"`
	UnseenSwap12Accuracy      float64         `json:"unseen_swap12_accuracy"`
	UnseenSwap23Accuracy      float64         `json:"unseen_swap23_accuracy"`
	DerivedRole1Accuracy      float64         `json:"derived_role1_accuracy"`
	DerivedRole2Accuracy      float64         `json:"derived_role2_accuracy"`
	DerivedRole3Accuracy      float64         `json:"derived_role3_accuracy"`
	LongProgramAccuracy       float64         `json:"long_program_accuracy"`
	FourRoleTransferGate      bool            `json:"four_role_transfer_gate"`
	NonUnitaryHeldOutAccuracy float64         `json:"nonunitary_heldout_accuracy"`
}

func up54aIndex(roles up54aRoles) (int,error) {
	index:=0
	for _,v:=range roles {
		if v<0 || v>=3 {return 0,fmt.Errorf("UP54A role value out of range")}
		index=index*3+v
	}
	return index,nil
}

func up54aDecode(index int)(up54aRoles,error){
	if index<0 || index>=81{return up54aRoles{},fmt.Errorf("UP54A joint index out of range")}
	var roles up54aRoles
	for i:=3;i>=0;i--{
		roles[i]=index%3
		index/=3
	}
	return roles,nil
}

func up54aBasis(roles up54aRoles)(State,error){
	index,err:=up54aIndex(roles)
	if err!=nil{return nil,err}
	state:=make(State,81)
	state[index]=1
	return state,nil
}

func up54aAllRoles()[]up54aRoles{
	out:=make([]up54aRoles,0,81)
	for a:=0;a<3;a++{
		for b:=0;b<3;b++{
			for c:=0;c<3;c++{
				for d:=0;d<3;d++{
					out=append(out,up54aRoles{a,b,c,d})
				}
			}
		}
	}
	return out
}

func up54aApplyLogical(roles up54aRoles,instruction UP54AInstruction)(up54aRoles,error){
	switch instruction.Kind{
	case "value_swap":
		if instruction.Start<0 || instruction.Start>=3{return up54aRoles{},fmt.Errorf("UP54A value start out of range")}
		a:=instruction.Start
		b:=(a+1)%3
		switch roles[0]{
		case a: roles[0]=b
		case b: roles[0]=a
		}
	case "role_swap":
		if instruction.Position<0 || instruction.Position>2{return up54aRoles{},fmt.Errorf("UP54A role swap position out of range")}
		p:=instruction.Position
		roles[p],roles[p+1]=roles[p+1],roles[p]
	default:
		return up54aRoles{},fmt.Errorf("UP54A unknown instruction %q",instruction.Kind)
	}
	return roles,nil
}

func up54aTarget(roles up54aRoles,program []UP54AInstruction)(up54aRoles,error){
	var err error
	for _,instruction:=range program{
		roles,err=up54aApplyLogical(roles,instruction)
		if err!=nil{return up54aRoles{},err}
	}
	return roles,nil
}

func up54aTrainSamples()([]up54aSample,error){
	var out []up54aSample
	for _,roles:=range up54aAllRoles(){
		for start:=0;start<3;start++{
			p:=[]UP54AInstruction{{Kind:"value_swap",Start:start}}
			t,err:=up54aTarget(roles,p);if err!=nil{return nil,err}
			out=append(out,up54aSample{Roles:roles,Program:p,Target:t,Category:"train_role0_value"})
		}
		p:=[]UP54AInstruction{{Kind:"role_swap",Position:0}}
		t,err:=up54aTarget(roles,p);if err!=nil{return nil,err}
		out=append(out,up54aSample{Roles:roles,Program:p,Target:t,Category:"train_swap01"})
	}
	return out,nil
}

func up54aDerivedMutation(position,start int)[]UP54AInstruction{
	var p []UP54AInstruction
	for current:=position-1;current>=0;current--{
		p=append(p,UP54AInstruction{Kind:"role_swap",Position:current})
	}
	p=append(p,UP54AInstruction{Kind:"value_swap",Start:start})
	for current:=0;current<position;current++{
		p=append(p,UP54AInstruction{Kind:"role_swap",Position:current})
	}
	return p
}

func up54aHeldSamples(lengths []int)([]up54aSample,error){
	var out []up54aSample
	for _,roles:=range up54aAllRoles(){
		for _,position:=range []int{1,2}{
			p:=[]UP54AInstruction{{Kind:"role_swap",Position:position}}
			t,err:=up54aTarget(roles,p);if err!=nil{return nil,err}
			category:="unseen_swap12"
			if position==2{category="unseen_swap23"}
			out=append(out,up54aSample{Roles:roles,Program:p,Target:t,Category:category})
		}
		for role:=1;role<=3;role++{
			for start:=0;start<3;start++{
				p:=up54aDerivedMutation(role,start)
				t,err:=up54aTarget(roles,p);if err!=nil{return nil,err}
				category:=fmt.Sprintf("derived_role%d",role)
				out=append(out,up54aSample{Roles:roles,Program:p,Target:t,Category:category})
			}
		}
	}
	for _,length:=range lengths{
		for _,roles:=range up54aAllRoles(){
			for seed:=0;seed<3;seed++{
				p:=make([]UP54AInstruction,length)
				for i:=range p{
					switch (i+seed+roles[0]+2*roles[1]+3*roles[2]+roles[3])%5{
					case 0:
						p[i]=UP54AInstruction{Kind:"role_swap",Position:0}
					case 1:
						p[i]=UP54AInstruction{Kind:"role_swap",Position:1}
					case 2:
						p[i]=UP54AInstruction{Kind:"role_swap",Position:2}
					default:
						p[i]=UP54AInstruction{Kind:"value_swap",Start:(i*i+seed+roles[2])%3}
					}
				}
				t,err:=up54aTarget(roles,p);if err!=nil{return nil,err}
				out=append(out,up54aSample{Roles:roles,Program:p,Target:t,Category:"long_program"})
			}
		}
	}
	return out,nil
}

func up54aCouplings(instruction UP54AInstruction,params []float64)([]Coupling,error){
	if len(params)!=2{return nil,fmt.Errorf("UP54A parameter count=%d want=2",len(params))}
	switch instruction.Kind{
	case "value_swap":
		if instruction.Start<0 || instruction.Start>=3{return nil,fmt.Errorf("UP54A value start out of range")}
		a:=instruction.Start;b:=(a+1)%3
		var out []Coupling
		for _,roles:=range up54aAllRoles(){
			if roles[0]!=a{continue}
			rb:=roles;rb[0]=b
			ia,_:=up54aIndex(roles);ib,_:=up54aIndex(rb)
			out=append(out,Coupling{A:ia,B:ib,Theta:params[0]})
		}
		return out,nil
	case "role_swap":
		if instruction.Position<0 || instruction.Position>2{return nil,fmt.Errorf("UP54A role swap position out of range")}
		var out []Coupling
		p:=instruction.Position
		for _,roles:=range up54aAllRoles(){
			if roles[p]>=roles[p+1]{continue}
			swapped:=roles
			swapped[p],swapped[p+1]=swapped[p+1],swapped[p]
			ia,_:=up54aIndex(roles);ib,_:=up54aIndex(swapped)
			out=append(out,Coupling{A:ia,B:ib,Theta:params[1]})
		}
		return out,nil
	default:
		return nil,fmt.Errorf("UP54A unknown instruction %q",instruction.Kind)
	}
}

type up54aApply func(State,[]float64,[]UP54AInstruction)(State,error)

func up54aApplyUnitary(initial State,params []float64,program []UP54AInstruction)(State,error){
	out:=append(State(nil),initial...)
	for _,instruction:=range program{
		couplings,err:=up54aCouplings(instruction,params);if err!=nil{return nil,err}
		out,err=Propagate(out,couplings);if err!=nil{return nil,err}
	}
	return out,nil
}

func up54aApplyNonUnitary(initial State,params []float64,program []UP54AInstruction)(State,error){
	if err:=validateState(initial);err!=nil{return nil,err}
	out:=append(State(nil),initial...)
	for _,instruction:=range program{
		couplings,err:=up54aCouplings(instruction,params);if err!=nil{return nil,err}
		for _,c:=range couplings{
			g:=c.Theta;a:=out[c.A];b:=out[c.B]
			out[c.A]=a-complex(g,0)*b
			out[c.B]=complex(g,0)*a+b
		}
	}
	return out,nil
}

func up54aInverse(state State,params []float64,program []UP54AInstruction)(State,error){
	out:=append(State(nil),state...)
	for i:=len(program)-1;i>=0;i--{
		couplings,err:=up54aCouplings(program[i],params);if err!=nil{return nil,err}
		for j:=range couplings{couplings[j].Theta=-couplings[j].Theta}
		out,err=Propagate(out,couplings);if err!=nil{return nil,err}
	}
	return out,nil
}

func up54aEvaluate(params []float64,samples []up54aSample,apply up54aApply,roundTrip bool)(loss,accuracy,maxNorm,maxRT float64,categories []UP54ACategoryMetric,err error){
	type counts struct{samples,hits int}
	per:=map[string]*counts{};order:=[]string{};hits:=0
	for _,sample:=range samples{
		initial,e:=up54aBasis(sample.Roles);if e!=nil{err=e;return}
		evolved,e:=apply(initial,params,sample.Program);if e!=nil{err=e;return}
		probs,e:=Probabilities(evolved);if e!=nil{err=e;return}
		targetIndex,e:=up54aIndex(sample.Target);if e!=nil{err=e;return}
		p:=probs[targetIndex];if p<1e-15{p=1e-15};loss+=-math.Log(p)
		observedIndex,e:=ArgMax(evolved);if e!=nil{err=e;return}
		observed,e:=up54aDecode(observedIndex);if e!=nil{err=e;return}
		if _,ok:=per[sample.Category];!ok{per[sample.Category]=&counts{};order=append(order,sample.Category)}
		per[sample.Category].samples++
		if observed==sample.Target{hits++;per[sample.Category].hits++}
		norm2,e:=NormSquared(evolved);if e!=nil{err=e;return}
		if drift:=math.Abs(norm2-1);drift>maxNorm{maxNorm=drift}
		if roundTrip{
			recovered,e:=up54aInverse(evolved,params,sample.Program);if e!=nil{err=e;return}
			distance,e:=L2Distance(initial,recovered);if e!=nil{err=e;return}
			if distance>maxRT{maxRT=distance}
		}
	}
	if len(samples)==0{err=fmt.Errorf("UP54A empty samples");return}
	loss/=float64(len(samples));accuracy=float64(hits)/float64(len(samples))
	for _,category:=range order{
		c:=per[category]
		categories=append(categories,UP54ACategoryMetric{Category:category,Samples:c.samples,Accuracy:float64(c.hits)/float64(c.samples)})
	}
	return
}

func up54aLoss(params []float64,samples []up54aSample,apply up54aApply)(float64,error){
	loss,_,_,_,_,err:=up54aEvaluate(params,samples,apply,false);return loss,err
}

func up54aTrain(name string,train,held []up54aSample,apply up54aApply,roundTrip bool,steps int,lr,eps float64)(UP54APathResult,error){
	params:=[]float64{0.1,0.1}
	initialLoss,initialAcc,_,_,_,err:=up54aEvaluate(params,train,apply,false);if err!=nil{return UP54APathResult{},err}
	for step:=0;step<steps;step++{
		grad:=make([]float64,2)
		for j:=0;j<2;j++{
			plus:=append([]float64(nil),params...);minus:=append([]float64(nil),params...)
			plus[j]+=eps;minus[j]-=eps
			lp,e:=up54aLoss(plus,train,apply);if e!=nil{return UP54APathResult{},e}
			lm,e:=up54aLoss(minus,train,apply);if e!=nil{return UP54APathResult{},e}
			grad[j]=(lp-lm)/(2*eps)
			if !finite(grad[j]){return UP54APathResult{},fmt.Errorf("UP54A %s gradient %d nonfinite",name,j)}
		}
		for j:=0;j<2;j++{params[j]-=lr*grad[j]}
	}
	finalLoss,trainAcc,trainNorm,trainRT,_,err:=up54aEvaluate(params,train,apply,roundTrip);if err!=nil{return UP54APathResult{},err}
	_,heldAcc,heldNorm,heldRT,categories,err:=up54aEvaluate(params,held,apply,roundTrip);if err!=nil{return UP54APathResult{},err}
	return UP54APathResult{
		Name:name,InitialTrainAccuracy:initialAcc,TrainAccuracy:trainAcc,HeldOutAccuracy:heldAcc,
		CategoryMetrics:categories,InitialLoss:initialLoss,FinalLoss:finalLoss,Parameters:append([]float64(nil),params...),
		MaxNormDrift:math.Max(trainNorm,heldNorm),MaxRoundTripError:math.Max(trainRT,heldRT),
	},nil
}

func up54aCategory(metrics []UP54ACategoryMetric,name string)float64{
	for _,m:=range metrics{if m.Category==name{return m.Accuracy}}
	return 0
}

func RunUP54A()(UP54AFourRoleResult,error){
	const(steps=140;lr=0.10;eps=1e-6)
	lengths:=[]int{16,48}
	train,err:=up54aTrainSamples();if err!=nil{return UP54AFourRoleResult{},err}
	held,err:=up54aHeldSamples(lengths);if err!=nil{return UP54AFourRoleResult{},err}
	u,err:=up54aTrain("unitary",train,held,up54aApplyUnitary,true,steps,lr,eps);if err!=nil{return UP54AFourRoleResult{},err}
	n,err:=up54aTrain("nonunitary_matched",train,held,up54aApplyNonUnitary,false,steps,lr,eps);if err!=nil{return UP54AFourRoleResult{},err}
	swap12:=up54aCategory(u.CategoryMetrics,"unseen_swap12")
	swap23:=up54aCategory(u.CategoryMetrics,"unseen_swap23")
	role1:=up54aCategory(u.CategoryMetrics,"derived_role1")
	role2:=up54aCategory(u.CategoryMetrics,"derived_role2")
	role3:=up54aCategory(u.CategoryMetrics,"derived_role3")
	long:=up54aCategory(u.CategoryMetrics,"long_program")
	gate:=u.TrainAccuracy>=0.99&&u.HeldOutAccuracy>=0.98&&swap12>=0.98&&swap23>=0.98&&role1>=0.98&&role2>=0.98&&role3>=0.98&&long>=0.98&&u.MaxNormDrift<=1e-9&&u.MaxRoundTripError<=1e-8
	return UP54AFourRoleResult{
		Schema:UP54AFourRoleSchema,Experiment:"UP-54A-four-role-transfer",
		SourceUP53ASeal:"b02660bf643d98f70468ed6296e3547be4ea422c",
		ValuesPerRole:3,Roles:4,JointDimension:81,TrainSingleStepOnly:true,OnlySwap01Trained:true,
		HeldOutLengths:append([]int(nil),lengths...),Steps:steps,LearningRate:lr,GradientEpsilon:eps,
		Unitary:u,NonUnitary:n,
		UnseenSwap12Accuracy:swap12,UnseenSwap23Accuracy:swap23,
		DerivedRole1Accuracy:role1,DerivedRole2Accuracy:role2,DerivedRole3Accuracy:role3,
		LongProgramAccuracy:long,FourRoleTransferGate:gate,NonUnitaryHeldOutAccuracy:n.HeldOutAccuracy,
	},nil
}
