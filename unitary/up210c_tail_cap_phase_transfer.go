package unitary

const UP210CPhaseTransferSchema="wingless.up210c-tail-cap-phase-transfer.v1"

type UP210CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	OpportunityCap int `json:"opportunity_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP210CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP209CSeal string `json:"source_up209c_seal"`
	PhaseAdvances []int `json:"phase_advances"`
	OpportunityCaps []int `json:"opportunity_caps"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	Horizon int `json:"horizon"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP210CMetric `json:"metrics"`
}
func up210cPolicy(step,advance int,a,b string)string{
	first:=8-advance
	if step<=first{return a}
	block:=(step-first-1)/8
	if block%2==0{return b}
	return a
}
func up210cBaseline(x0 *up81cAging,endangered,advance int,a,b string)int{
	m:=*x0
	for start:=1;start<=56;start+=2{
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1034000+step,step,up210cPolicy(step,advance,a,b)){return step}
		}
	}
	return 57
}
func up210cTreated(x0 *up81cAging,endangered,advance int,a,b string,cap int)(loss,actions int){
	m:=*x0;committed:=false;since:=0;opportunities:=0
	for start:=1;start<=56;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				opportunities=1;m.query(endangered);actions++;committed=true;since=0
			}
		}else if opportunities<cap{
			since++
			if since>=4{opportunities++;m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1035000+step,step,up210cPolicy(step,advance,a,b)){return step,actions}
		}
	}
	return 57,actions
}
func RunUP210C()(UP210CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{0,2,4,6};caps:=[]int{7,8}
	res:=UP210CResult{Schema:UP210CPhaseTransferSchema,Experiment:"UP-210C-tail-cap-phase-transfer",SourceUP209CSeal:"f30b098ea5914c3edd2f513633b2084d3afe1a70",PhaseAdvances:advances,OpportunityCaps:caps,Cadence:2,SpacingIntervals:4,Horizon:56,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,cap:=range caps{
		m:=UP210CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,OpportunityCap:cap};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up210cBaseline(x,e,advance,s.a,s.b)
			treated,actions:=up210cTreated(x,e,advance,s.a,s.b,cap)
			if base<=56{m.BaselineLosses++}
			if treated<=56{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=56&&treated>56{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			m.ActionsTaken+=actions
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
