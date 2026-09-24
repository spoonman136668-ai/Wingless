package unitary

import (
	"fmt"
	"math"
)

const UP56ADoubleControlledSchema = "wingless.up56a-double-controlled.v1"

type UP56AInstruction struct {
	Kind          string `json:"kind"`
	Position      int    `json:"position,omitempty"`
	TargetRole    int    `json:"target_role,omitempty"`
	ControlRoleA  int    `json:"control_role_a,omitempty"`
	ControlValueA int    `json:"control_value_a,omitempty"`
	ControlRoleB  int    `json:"control_role_b,omitempty"`
	ControlValueB int    `json:"control_value_b,omitempty"`
	Start         int    `json:"start,omitempty"`
}

type up56aRoles [3]int
type up56aSample struct{ Roles up56aRoles; Program []UP56AInstruction; Target up56aRoles; Category string }

type UP56ACategoryMetric struct{ Category string `json:"category"`; Samples int `json:"samples"`; Accuracy float64 `json:"accuracy"` }
type UP56APathResult struct {
	Name string `json:"name"`; InitialTrainAccuracy float64 `json:"initial_train_accuracy"`; TrainAccuracy float64 `json:"train_accuracy"`;
	HeldOutAccuracy float64 `json:"heldout_accuracy"`; CategoryMetrics []UP56ACategoryMetric `json:"category_metrics"`;
	InitialLoss float64 `json:"initial_loss"`; FinalLoss float64 `json:"final_loss"`; Parameters []float64 `json:"parameters"`;
	MaxNormDrift float64 `json:"max_norm_drift"`; MaxRoundTripError float64 `json:"max_round_trip_error"`
}
type UP56ADoubleControlledResult struct {
	Schema string `json:"schema"`; Experiment string `json:"experiment"`; SourceUP55ASeal string `json:"source_up55a_seal"`;
	TrainSingleStepOnly bool `json:"train_single_step_only"`; TrainedControlValues []int `json:"trained_control_values"`;
	Unitary UP56APathResult `json:"unitary"`; NonUnitary UP56APathResult `json:"nonunitary_matched"`;
	UnseenControlTupleAccuracy float64 `json:"unseen_control_tuple_accuracy"`; MovedTargetAccuracy float64 `json:"moved_target_accuracy"`;
	LongProgramAccuracy float64 `json:"long_program_accuracy"`; DoubleControlledGate bool `json:"double_controlled_gate"`;
	NonUnitaryHeldOutAccuracy float64 `json:"nonunitary_heldout_accuracy"`
}

func up56aAllRoles() []up56aRoles { out:=make([]up56aRoles,0,27); for a:=0;a<3;a++{for b:=0;b<3;b++{for c:=0;c<3;c++{out=append(out,up56aRoles{a,b,c})}}}; return out }
func up56aIndex(r up56aRoles)(int,error){for _,v:=range r{if v<0||v>2{return 0,fmt.Errorf("UP56A role value")}};return r[0]*9+r[1]*3+r[2],nil}
func up56aDecode(i int)(up56aRoles,error){if i<0||i>=27{return up56aRoles{},fmt.Errorf("UP56A index")};return up56aRoles{i/9,(i/3)%3,i%3},nil}
func up56aBasis(r up56aRoles)(State,error){i,e:=up56aIndex(r);if e!=nil{return nil,e};s:=make(State,27);s[i]=1;return s,nil}

func up56aValidate(in UP56AInstruction) error {
	switch in.Kind {
	case "role_swap":
		if in.Position<0||in.Position>1{return fmt.Errorf("UP56A role position")}
	case "double_controlled_value_swap":
		roles:=[]int{in.TargetRole,in.ControlRoleA,in.ControlRoleB}
		for _,r:=range roles{if r<0||r>2{return fmt.Errorf("UP56A role")}}
		if in.TargetRole==in.ControlRoleA||in.TargetRole==in.ControlRoleB||in.ControlRoleA==in.ControlRoleB{return fmt.Errorf("UP56A roles not distinct")}
		if in.ControlValueA<0||in.ControlValueA>2||in.ControlValueB<0||in.ControlValueB>2||in.Start<0||in.Start>2{return fmt.Errorf("UP56A value")}
	default:return fmt.Errorf("UP56A kind")
	}
	return nil
}
func up56aApplyLogical(r up56aRoles,in UP56AInstruction)(up56aRoles,error){
	if e:=up56aValidate(in);e!=nil{return up56aRoles{},e}
	if in.Kind=="role_swap"{r[in.Position],r[in.Position+1]=r[in.Position+1],r[in.Position];return r,nil}
	if r[in.ControlRoleA]==in.ControlValueA&&r[in.ControlRoleB]==in.ControlValueB{
		a:=in.Start;b:=(a+1)%3;switch r[in.TargetRole]{case a:r[in.TargetRole]=b;case b:r[in.TargetRole]=a}
	}
	return r,nil
}
func up56aTarget(r up56aRoles,p []UP56AInstruction)(up56aRoles,error){var e error;for _,in:=range p{r,e=up56aApplyLogical(r,in);if e!=nil{return up56aRoles{},e}};return r,nil}

func up56aTrainSamples()([]up56aSample,error){
	var out []up56aSample
	for _,r:=range up56aAllRoles(){
		for start:=0;start<3;start++{
			p:=[]UP56AInstruction{{Kind:"double_controlled_value_swap",TargetRole:0,ControlRoleA:1,ControlValueA:0,ControlRoleB:2,ControlValueB:0,Start:start}}
			t,e:=up56aTarget(r,p);if e!=nil{return nil,e};out=append(out,up56aSample{r,p,t,"train_double_control_00"})
		}
		p:=[]UP56AInstruction{{Kind:"role_swap",Position:0}};t,e:=up56aTarget(r,p);if e!=nil{return nil,e};out=append(out,up56aSample{r,p,t,"train_swap01"})
	}
	return out,nil
}
func up56aHeldSamples(lengths []int)([]up56aSample,error){
	var out []up56aSample
	for _,r:=range up56aAllRoles(){
		for ca:=0;ca<3;ca++{for cb:=0;cb<3;cb++{if ca==0&&cb==0{continue};for start:=0;start<3;start++{
			p:=[]UP56AInstruction{{Kind:"double_controlled_value_swap",TargetRole:0,ControlRoleA:1,ControlValueA:ca,ControlRoleB:2,ControlValueB:cb,Start:start}}
			t,e:=up56aTarget(r,p);if e!=nil{return nil,e};out=append(out,up56aSample{r,p,t,"unseen_control_tuple"})
		}}}
		for target:=1;target<=2;target++{
			ca:=(target+1)%3;cb:=(target+2)%3
			for va:=0;va<3;va++{for vb:=0;vb<3;vb++{for start:=0;start<3;start++{
				p:=[]UP56AInstruction{{Kind:"double_controlled_value_swap",TargetRole:target,ControlRoleA:ca,ControlValueA:va,ControlRoleB:cb,ControlValueB:vb,Start:start}}
				t,e:=up56aTarget(r,p);if e!=nil{return nil,e};out=append(out,up56aSample{r,p,t,"moved_target"})
			}}}
		}
	}
	for _,length:=range lengths{for _,r:=range up56aAllRoles(){for seed:=0;seed<3;seed++{
		p:=make([]UP56AInstruction,length)
		for i:=range p{
			if (i+seed+r[0]+r[1])%4==0{p[i]=UP56AInstruction{Kind:"role_swap",Position:(i+seed)%2};continue}
			target:=(i+seed)%3;ca:=(target+1)%3;cb:=(target+2)%3
			p[i]=UP56AInstruction{Kind:"double_controlled_value_swap",TargetRole:target,ControlRoleA:ca,ControlValueA:(i+r[1])%3,ControlRoleB:cb,ControlValueB:(i*i+seed+r[2])%3,Start:(i+seed+r[0])%3}
		}
		t,e:=up56aTarget(r,p);if e!=nil{return nil,e};out=append(out,up56aSample{r,p,t,"long_program"})
	}}}
	return out,nil
}
func up56aCouplings(in UP56AInstruction,params []float64)([]Coupling,error){
	if len(params)!=2{return nil,fmt.Errorf("UP56A params")};if e:=up56aValidate(in);e!=nil{return nil,e};var out []Coupling
	if in.Kind=="role_swap"{for _,r:=range up56aAllRoles(){p:=in.Position;if r[p]>=r[p+1]{continue};s:=r;s[p],s[p+1]=s[p+1],s[p];a,_:=up56aIndex(r);b,_:=up56aIndex(s);out=append(out,Coupling{A:a,B:b,Theta:params[0]})};return out,nil}
	aVal:=in.Start;bVal:=(aVal+1)%3
	for _,r:=range up56aAllRoles(){if r[in.ControlRoleA]!=in.ControlValueA||r[in.ControlRoleB]!=in.ControlValueB||r[in.TargetRole]!=aVal{continue};s:=r;s[in.TargetRole]=bVal;a,_:=up56aIndex(r);b,_:=up56aIndex(s);out=append(out,Coupling{A:a,B:b,Theta:params[1]})}
	return out,nil
}
type up56aApply func(State,[]float64,[]UP56AInstruction)(State,error)
func up56aApplyUnitary(initial State,params []float64,p []UP56AInstruction)(State,error){out:=append(State(nil),initial...);for _,in:=range p{cs,e:=up56aCouplings(in,params);if e!=nil{return nil,e};out,e=Propagate(out,cs);if e!=nil{return nil,e}};return out,nil}
func up56aApplyNonUnitary(initial State,params []float64,p []UP56AInstruction)(State,error){if e:=validateState(initial);e!=nil{return nil,e};out:=append(State(nil),initial...);for _,in:=range p{cs,e:=up56aCouplings(in,params);if e!=nil{return nil,e};for _,c:=range cs{g:=c.Theta;a:=out[c.A];b:=out[c.B];out[c.A]=a-complex(g,0)*b;out[c.B]=complex(g,0)*a+b}};return out,nil}
func up56aInverse(state State,params []float64,p []UP56AInstruction)(State,error){out:=append(State(nil),state...);for i:=len(p)-1;i>=0;i--{cs,e:=up56aCouplings(p[i],params);if e!=nil{return nil,e};for j:=range cs{cs[j].Theta=-cs[j].Theta};out,e=Propagate(out,cs);if e!=nil{return nil,e}};return out,nil}

func up56aEvaluate(params []float64,samples []up56aSample,apply up56aApply,rt bool)(loss,acc,maxNorm,maxRT float64,cats []UP56ACategoryMetric,err error){
	type c struct{n,h int};per:=map[string]*c{};order:=[]string{};hits:=0
	for _,s:=range samples{initial,e:=up56aBasis(s.Roles);if e!=nil{err=e;return};ev,e:=apply(initial,params,s.Program);if e!=nil{err=e;return};probs,e:=Probabilities(ev);if e!=nil{err=e;return};ti,_:=up56aIndex(s.Target);p:=probs[ti];if p<1e-15{p=1e-15};loss+=-math.Log(p);oi,e:=ArgMax(ev);if e!=nil{err=e;return};obs,_:=up56aDecode(oi);if _,ok:=per[s.Category];!ok{per[s.Category]=&c{};order=append(order,s.Category)};per[s.Category].n++;if obs==s.Target{hits++;per[s.Category].h++};n2,e:=NormSquared(ev);if e!=nil{err=e;return};if d:=math.Abs(n2-1);d>maxNorm{maxNorm=d};if rt{rec,e:=up56aInverse(ev,params,s.Program);if e!=nil{err=e;return};d,e:=L2Distance(initial,rec);if e!=nil{err=e;return};if d>maxRT{maxRT=d}}}
	if len(samples)==0{err=fmt.Errorf("UP56A empty");return};loss/=float64(len(samples));acc=float64(hits)/float64(len(samples));for _,name:=range order{x:=per[name];cats=append(cats,UP56ACategoryMetric{name,x.n,float64(x.h)/float64(x.n)})};return
}
func up56aLoss(p []float64,s []up56aSample,a up56aApply)(float64,error){l,_,_,_,_,e:=up56aEvaluate(p,s,a,false);return l,e}
func up56aTrain(name string,train,held []up56aSample,apply up56aApply,rt bool)(UP56APathResult,error){
	const(steps=180;lr=0.10;eps=1e-6);params:=[]float64{0.1,0.1};il,ia,_,_,_,e:=up56aEvaluate(params,train,apply,false);if e!=nil{return UP56APathResult{},e}
	for step:=0;step<steps;step++{g:=make([]float64,2);for j:=0;j<2;j++{plus:=append([]float64(nil),params...);minus:=append([]float64(nil),params...);plus[j]+=eps;minus[j]-=eps;lp,e:=up56aLoss(plus,train,apply);if e!=nil{return UP56APathResult{},e};lm,e:=up56aLoss(minus,train,apply);if e!=nil{return UP56APathResult{},e};g[j]=(lp-lm)/(2*eps);if !finite(g[j]){return UP56APathResult{},fmt.Errorf("UP56A gradient")}};for j:=range params{params[j]-=lr*g[j]}}
	fl,ta,tn,tr,_,e:=up56aEvaluate(params,train,apply,rt);if e!=nil{return UP56APathResult{},e};_,ha,hn,hr,cats,e:=up56aEvaluate(params,held,apply,rt);if e!=nil{return UP56APathResult{},e};return UP56APathResult{name,ia,ta,ha,cats,il,fl,append([]float64(nil),params...),math.Max(tn,hn),math.Max(tr,hr)},nil
}
func up56aCat(ms []UP56ACategoryMetric,n string)float64{for _,m:=range ms{if m.Category==n{return m.Accuracy}};return 0}
func RunUP56A()(UP56ADoubleControlledResult,error){
	train,e:=up56aTrainSamples();if e!=nil{return UP56ADoubleControlledResult{},e};held,e:=up56aHeldSamples([]int{16,48});if e!=nil{return UP56ADoubleControlledResult{},e}
	u,e:=up56aTrain("unitary",train,held,up56aApplyUnitary,true);if e!=nil{return UP56ADoubleControlledResult{},e};n,e:=up56aTrain("nonunitary_matched",train,held,up56aApplyNonUnitary,false);if e!=nil{return UP56ADoubleControlledResult{},e}
	uc:=up56aCat(u.CategoryMetrics,"unseen_control_tuple");mt:=up56aCat(u.CategoryMetrics,"moved_target");long:=up56aCat(u.CategoryMetrics,"long_program")
	gate:=u.TrainAccuracy>=0.99&&u.HeldOutAccuracy>=0.98&&uc>=0.98&&mt>=0.98&&long>=0.98&&u.MaxNormDrift<=1e-9&&u.MaxRoundTripError<=1e-8
	return UP56ADoubleControlledResult{Schema:UP56ADoubleControlledSchema,Experiment:"UP-56A-double-controlled",SourceUP55ASeal:"6a0f346422c273998953b8a84793614456369dc3",TrainSingleStepOnly:true,TrainedControlValues:[]int{0,0},Unitary:u,NonUnitary:n,UnseenControlTupleAccuracy:uc,MovedTargetAccuracy:mt,LongProgramAccuracy:long,DoubleControlledGate:gate,NonUnitaryHeldOutAccuracy:n.HeldOutAccuracy},nil
}
