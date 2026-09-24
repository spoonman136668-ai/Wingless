package unitary

import (
	"fmt"
	"math"
)

const UP55AControlledSchema = "wingless.up55a-controlled-composition.v1"

type UP55AInstruction struct {
	Kind         string `json:"kind"`
	Position     int    `json:"position,omitempty"`
	TargetRole   int    `json:"target_role,omitempty"`
	ControlRole  int    `json:"control_role,omitempty"`
	ControlValue int    `json:"control_value,omitempty"`
	Start        int    `json:"start,omitempty"`
}

type up55aRoles [3]int

type up55aSample struct {
	Roles    up55aRoles
	Program  []UP55AInstruction
	Target   up55aRoles
	Category string
}

type UP55ACategoryMetric struct {
	Category string  `json:"category"`
	Samples  int     `json:"samples"`
	Accuracy float64 `json:"accuracy"`
}

type UP55APathResult struct {
	Name                 string                 `json:"name"`
	InitialTrainAccuracy float64                `json:"initial_train_accuracy"`
	TrainAccuracy        float64                `json:"train_accuracy"`
	HeldOutAccuracy      float64                `json:"heldout_accuracy"`
	CategoryMetrics      []UP55ACategoryMetric `json:"category_metrics"`
	InitialLoss          float64                `json:"initial_loss"`
	FinalLoss            float64                `json:"final_loss"`
	Parameters           []float64              `json:"parameters"`
	MaxNormDrift         float64                `json:"max_norm_drift"`
	MaxRoundTripError    float64                `json:"max_round_trip_error"`
}

type UP55AControlledResult struct {
	Schema                       string          `json:"schema"`
	Experiment                   string          `json:"experiment"`
	SourceUP54ASeal              string          `json:"source_up54a_seal"`
	Roles                        int             `json:"roles"`
	ValuesPerRole                int             `json:"values_per_role"`
	JointDimension               int             `json:"joint_dimension"`
	TrainSingleStepOnly          bool            `json:"train_single_step_only"`
	TrainedControlValue          int             `json:"trained_control_value"`
	TrainedTargetRole            int             `json:"trained_target_role"`
	TrainedControlRole           int             `json:"trained_control_role"`
	OnlySwap01Trained            bool            `json:"only_swap01_trained"`
	HeldOutLengths               []int           `json:"heldout_lengths"`
	Unitary                      UP55APathResult `json:"unitary"`
	NonUnitary                   UP55APathResult `json:"nonunitary_matched"`
	UnseenControlValueAccuracy   float64         `json:"unseen_control_value_accuracy"`
	MovedControlAccuracy         float64         `json:"moved_control_accuracy"`
	MovedTargetAccuracy          float64         `json:"moved_target_accuracy"`
	LongProgramAccuracy          float64         `json:"long_program_accuracy"`
	ControlledCompositionGate    bool            `json:"controlled_composition_gate"`
	NonUnitaryHeldOutAccuracy    float64         `json:"nonunitary_heldout_accuracy"`
}

func up55aAllRoles() []up55aRoles {
	out:=make([]up55aRoles,0,27)
	for a:=0;a<3;a++{for b:=0;b<3;b++{for c:=0;c<3;c++{out=append(out,up55aRoles{a,b,c})}}}
	return out
}

func up55aIndex(r up55aRoles)(int,error){
	for _,v:=range r{if v<0||v>=3{return 0,fmt.Errorf("UP55A role value out of range")}}
	return r[0]*9+r[1]*3+r[2],nil
}

func up55aDecode(index int)(up55aRoles,error){
	if index<0||index>=27{return up55aRoles{},fmt.Errorf("UP55A index out of range")}
	return up55aRoles{index/9,(index/3)%3,index%3},nil
}

func up55aBasis(r up55aRoles)(State,error){
	i,err:=up55aIndex(r);if err!=nil{return nil,err}
	s:=make(State,27);s[i]=1;return s,nil
}

func up55aValidateInstruction(in UP55AInstruction) error {
	switch in.Kind {
	case "role_swap":
		if in.Position<0||in.Position>1{return fmt.Errorf("UP55A role swap position")}
	case "controlled_value_swap":
		if in.TargetRole<0||in.TargetRole>2||in.ControlRole<0||in.ControlRole>2||in.TargetRole==in.ControlRole{return fmt.Errorf("UP55A controlled roles")}
		if in.ControlValue<0||in.ControlValue>2{return fmt.Errorf("UP55A control value")}
		if in.Start<0||in.Start>2{return fmt.Errorf("UP55A start")}
	default:
		return fmt.Errorf("UP55A unknown instruction %q",in.Kind)
	}
	return nil
}

func up55aApplyLogical(r up55aRoles,in UP55AInstruction)(up55aRoles,error){
	if err:=up55aValidateInstruction(in);err!=nil{return up55aRoles{},err}
	switch in.Kind{
	case "role_swap":
		r[in.Position],r[in.Position+1]=r[in.Position+1],r[in.Position]
	case "controlled_value_swap":
		if r[in.ControlRole]==in.ControlValue{
			a:=in.Start;b:=(a+1)%3
			switch r[in.TargetRole]{case a:r[in.TargetRole]=b;case b:r[in.TargetRole]=a}
		}
	}
	return r,nil
}

func up55aTarget(r up55aRoles,p []UP55AInstruction)(up55aRoles,error){
	var err error
	for _,in:=range p{r,err=up55aApplyLogical(r,in);if err!=nil{return up55aRoles{},err}}
	return r,nil
}

func up55aTrainSamples()([]up55aSample,error){
	var out []up55aSample
	for _,r:=range up55aAllRoles(){
		for start:=0;start<3;start++{
			p:=[]UP55AInstruction{{Kind:"controlled_value_swap",TargetRole:0,ControlRole:1,ControlValue:0,Start:start}}
			t,err:=up55aTarget(r,p);if err!=nil{return nil,err}
			out=append(out,up55aSample{Roles:r,Program:p,Target:t,Category:"train_control0_target0_control1"})
		}
		p:=[]UP55AInstruction{{Kind:"role_swap",Position:0}}
		t,err:=up55aTarget(r,p);if err!=nil{return nil,err}
		out=append(out,up55aSample{Roles:r,Program:p,Target:t,Category:"train_swap01"})
	}
	return out,nil
}

func up55aHeldSamples(lengths []int)([]up55aSample,error){
	var out []up55aSample
	for _,r:=range up55aAllRoles(){
		for _,controlValue:=range []int{1,2}{
			for start:=0;start<3;start++{
				p:=[]UP55AInstruction{{Kind:"controlled_value_swap",TargetRole:0,ControlRole:1,ControlValue:controlValue,Start:start}}
				t,err:=up55aTarget(r,p);if err!=nil{return nil,err}
				out=append(out,up55aSample{Roles:r,Program:p,Target:t,Category:"unseen_control_value"})
			}
		}
		for controlValue:=0;controlValue<3;controlValue++{
			for start:=0;start<3;start++{
				p:=[]UP55AInstruction{{Kind:"controlled_value_swap",TargetRole:0,ControlRole:2,ControlValue:controlValue,Start:start}}
				t,err:=up55aTarget(r,p);if err!=nil{return nil,err}
				out=append(out,up55aSample{Roles:r,Program:p,Target:t,Category:"moved_control"})
				p2:=[]UP55AInstruction{{Kind:"controlled_value_swap",TargetRole:1,ControlRole:0,ControlValue:controlValue,Start:start}}
				t2,err:=up55aTarget(r,p2);if err!=nil{return nil,err}
				out=append(out,up55aSample{Roles:r,Program:p2,Target:t2,Category:"moved_target"})
			}
		}
	}
	for _,length:=range lengths{
		for _,r:=range up55aAllRoles(){
			for seed:=0;seed<3;seed++{
				p:=make([]UP55AInstruction,length)
				for i:=range p{
					switch (i+seed+r[0]+2*r[1]+r[2])%5{
					case 0:
						p[i]=UP55AInstruction{Kind:"role_swap",Position:0}
					case 1:
						p[i]=UP55AInstruction{Kind:"role_swap",Position:1}
					default:
						target:=(i+seed)%3
						control:=(target+1+(i%2))%3
						if control==target{control=(control+1)%3}
						p[i]=UP55AInstruction{
							Kind:"controlled_value_swap",
							TargetRole:target,
							ControlRole:control,
							ControlValue:(i+seed+r[1])%3,
							Start:(i*i+seed+r[2])%3,
						}
					}
				}
				t,err:=up55aTarget(r,p);if err!=nil{return nil,err}
				out=append(out,up55aSample{Roles:r,Program:p,Target:t,Category:"long_program"})
			}
		}
	}
	return out,nil
}

func up55aCouplings(in UP55AInstruction,params []float64)([]Coupling,error){
	if len(params)!=2{return nil,fmt.Errorf("UP55A params")}
	if err:=up55aValidateInstruction(in);err!=nil{return nil,err}
	var out []Coupling
	switch in.Kind{
	case "role_swap":
		for _,r:=range up55aAllRoles(){
			p:=in.Position
			if r[p]>=r[p+1]{continue}
			s:=r;s[p],s[p+1]=s[p+1],s[p]
			a,_:=up55aIndex(r);b,_:=up55aIndex(s)
			out=append(out,Coupling{A:a,B:b,Theta:params[0]})
		}
	case "controlled_value_swap":
		aVal:=in.Start;bVal:=(aVal+1)%3
		for _,r:=range up55aAllRoles(){
			if r[in.ControlRole]!=in.ControlValue||r[in.TargetRole]!=aVal{continue}
			s:=r;s[in.TargetRole]=bVal
			a,_:=up55aIndex(r);b,_:=up55aIndex(s)
			out=append(out,Coupling{A:a,B:b,Theta:params[1]})
		}
	}
	return out,nil
}

type up55aApply func(State,[]float64,[]UP55AInstruction)(State,error)

func up55aApplyUnitary(initial State,params []float64,p []UP55AInstruction)(State,error){
	out:=append(State(nil),initial...)
	for _,in:=range p{c,err:=up55aCouplings(in,params);if err!=nil{return nil,err};out,err=Propagate(out,c);if err!=nil{return nil,err}}
	return out,nil
}

func up55aApplyNonUnitary(initial State,params []float64,p []UP55AInstruction)(State,error){
	if err:=validateState(initial);err!=nil{return nil,err}
	out:=append(State(nil),initial...)
	for _,in:=range p{
		cs,err:=up55aCouplings(in,params);if err!=nil{return nil,err}
		for _,c:=range cs{g:=c.Theta;a:=out[c.A];b:=out[c.B];out[c.A]=a-complex(g,0)*b;out[c.B]=complex(g,0)*a+b}
	}
	return out,nil
}

func up55aInverse(state State,params []float64,p []UP55AInstruction)(State,error){
	out:=append(State(nil),state...)
	for i:=len(p)-1;i>=0;i--{
		cs,err:=up55aCouplings(p[i],params);if err!=nil{return nil,err}
		for j:=range cs{cs[j].Theta=-cs[j].Theta}
		out,err=Propagate(out,cs);if err!=nil{return nil,err}
	}
	return out,nil
}

func up55aEvaluate(params []float64,samples []up55aSample,apply up55aApply,roundTrip bool)(loss,accuracy,maxNorm,maxRT float64,cats []UP55ACategoryMetric,err error){
	type counts struct{n,h int}; per:=map[string]*counts{};order:=[]string{};hits:=0
	for _,s:=range samples{
		initial,e:=up55aBasis(s.Roles);if e!=nil{err=e;return}
		evolved,e:=apply(initial,params,s.Program);if e!=nil{err=e;return}
		probs,e:=Probabilities(evolved);if e!=nil{err=e;return}
		ti,e:=up55aIndex(s.Target);if e!=nil{err=e;return}
		p:=probs[ti];if p<1e-15{p=1e-15};loss+=-math.Log(p)
		oi,e:=ArgMax(evolved);if e!=nil{err=e;return}
		obs,e:=up55aDecode(oi);if e!=nil{err=e;return}
		if _,ok:=per[s.Category];!ok{per[s.Category]=&counts{};order=append(order,s.Category)}
		per[s.Category].n++;if obs==s.Target{hits++;per[s.Category].h++}
		n2,e:=NormSquared(evolved);if e!=nil{err=e;return};if d:=math.Abs(n2-1);d>maxNorm{maxNorm=d}
		if roundTrip{recovered,e:=up55aInverse(evolved,params,s.Program);if e!=nil{err=e;return};d,e:=L2Distance(initial,recovered);if e!=nil{err=e;return};if d>maxRT{maxRT=d}}
	}
	if len(samples)==0{err=fmt.Errorf("UP55A empty samples");return}
	loss/=float64(len(samples));accuracy=float64(hits)/float64(len(samples))
	for _,name:=range order{c:=per[name];cats=append(cats,UP55ACategoryMetric{Category:name,Samples:c.n,Accuracy:float64(c.h)/float64(c.n)})}
	return
}

func up55aLoss(params []float64,s []up55aSample,a up55aApply)(float64,error){l,_,_,_,_,e:=up55aEvaluate(params,s,a,false);return l,e}

func up55aTrain(name string,train,held []up55aSample,apply up55aApply,roundTrip bool)(UP55APathResult,error){
	const(steps=160;lr=0.10;eps=1e-6)
	params:=[]float64{0.1,0.1}
	initialLoss,initialAcc,_,_,_,err:=up55aEvaluate(params,train,apply,false);if err!=nil{return UP55APathResult{},err}
	for step:=0;step<steps;step++{
		grad:=make([]float64,2)
		for j:=0;j<2;j++{
			plus:=append([]float64(nil),params...);minus:=append([]float64(nil),params...)
			plus[j]+=eps;minus[j]-=eps
			lp,e:=up55aLoss(plus,train,apply);if e!=nil{return UP55APathResult{},e}
			lm,e:=up55aLoss(minus,train,apply);if e!=nil{return UP55APathResult{},e}
			grad[j]=(lp-lm)/(2*eps);if !finite(grad[j]){return UP55APathResult{},fmt.Errorf("UP55A nonfinite gradient")}
		}
		for j:=range params{params[j]-=lr*grad[j]}
	}
	finalLoss,trainAcc,trainNorm,trainRT,_,err:=up55aEvaluate(params,train,apply,roundTrip);if err!=nil{return UP55APathResult{},err}
	_,heldAcc,heldNorm,heldRT,cats,err:=up55aEvaluate(params,held,apply,roundTrip);if err!=nil{return UP55APathResult{},err}
	return UP55APathResult{Name:name,InitialTrainAccuracy:initialAcc,TrainAccuracy:trainAcc,HeldOutAccuracy:heldAcc,CategoryMetrics:cats,InitialLoss:initialLoss,FinalLoss:finalLoss,Parameters:append([]float64(nil),params...),MaxNormDrift:math.Max(trainNorm,heldNorm),MaxRoundTripError:math.Max(trainRT,heldRT)},nil
}

func up55aCategory(metrics []UP55ACategoryMetric,name string)float64{for _,m:=range metrics{if m.Category==name{return m.Accuracy}};return 0}

func RunUP55A()(UP55AControlledResult,error){
	lengths:=[]int{16,48}
	train,err:=up55aTrainSamples();if err!=nil{return UP55AControlledResult{},err}
	held,err:=up55aHeldSamples(lengths);if err!=nil{return UP55AControlledResult{},err}
	u,err:=up55aTrain("unitary",train,held,up55aApplyUnitary,true);if err!=nil{return UP55AControlledResult{},err}
	n,err:=up55aTrain("nonunitary_matched",train,held,up55aApplyNonUnitary,false);if err!=nil{return UP55AControlledResult{},err}
	uc:=up55aCategory(u.CategoryMetrics,"unseen_control_value")
	mc:=up55aCategory(u.CategoryMetrics,"moved_control")
	mt:=up55aCategory(u.CategoryMetrics,"moved_target")
	long:=up55aCategory(u.CategoryMetrics,"long_program")
	gate:=u.TrainAccuracy>=0.99&&u.HeldOutAccuracy>=0.98&&uc>=0.98&&mc>=0.98&&mt>=0.98&&long>=0.98&&u.MaxNormDrift<=1e-9&&u.MaxRoundTripError<=1e-8
	return UP55AControlledResult{
		Schema:UP55AControlledSchema,Experiment:"UP-55A-controlled-composition",SourceUP54ASeal:"bddeaf6fc74cd4bfc509b09bdb948b7600a41088",
		Roles:3,ValuesPerRole:3,JointDimension:27,TrainSingleStepOnly:true,TrainedControlValue:0,TrainedTargetRole:0,TrainedControlRole:1,OnlySwap01Trained:true,HeldOutLengths:append([]int(nil),lengths...),
		Unitary:u,NonUnitary:n,UnseenControlValueAccuracy:uc,MovedControlAccuracy:mc,MovedTargetAccuracy:mt,LongProgramAccuracy:long,ControlledCompositionGate:gate,NonUnitaryHeldOutAccuracy:n.HeldOutAccuracy,
	},nil
}
