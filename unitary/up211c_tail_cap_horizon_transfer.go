package unitary

const UP211CHorizonTransferSchema="wingless.up211c-tail-cap-horizon-transfer.v1"

type UP211CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	OpportunityCap int `json:"opportunity_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP211CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP210CSeal string `json:"source_up210c_seal"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	OpportunityCaps []int `json:"opportunity_caps"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP211CMetric `json:"metrics"`
}
func up211cBaseline(x0 *up81cAging,endangered,advance,horizon int,a,b string)int{
	m:=*x0
	for start:=1;start<=horizon;start+=2{
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1036000+step,step,up210cPolicy(step,advance,a,b)){return step}
		}
	}
	return horizon+1
}
func up211cTreated(x0 *up81cAging,endangered,advance,horizon int,a,b string,cap int)(loss,actions int){
	m:=*x0;committed:=false;since:=0;opportunities:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				opportunities=1;m.query(endangered);actions++;committed=true;since=0
			}
		}else if opportunities<cap{
			since++
			if since>=4{opportunities++;m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1037000+step,step,up210cPolicy(step,advance,a,b)){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP211C()(UP211CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{0,4};horizons:=[]int{56,64,72};caps:=[]int{7,8}
	res:=UP211CResult{Schema:UP211CHorizonTransferSchema,Experiment:"UP-211C-tail-cap-horizon-transfer",SourceUP210CSeal:"9b6227c8ddf0d5ef28e248173b929b1af8b95da1",PhaseAdvances:advances,Horizons:horizons,OpportunityCaps:caps,Cadence:2,SpacingIntervals:4,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,horizon:=range horizons{for _,cap:=range caps{
		m:=UP211CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:horizon,OpportunityCap:cap};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up211cBaseline(x,e,advance,horizon,s.a,s.b)
			treated,actions:=up211cTreated(x,e,advance,horizon,s.a,s.b,cap)
			if base<=horizon{m.BaselineLosses++}
			if treated<=horizon{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=horizon&&treated>horizon{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			m.ActionsTaken+=actions
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
