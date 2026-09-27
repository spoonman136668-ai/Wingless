package unitary

const UP209CTailCapSchema="wingless.up209c-committed-tail-cap-transfer.v1"

type UP209CMetric struct{
	Schedule string `json:"schedule"`
	OpportunityCap int `json:"opportunity_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP209CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP208CSeal string `json:"source_up208c_seal"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	Horizon int `json:"horizon"`
	OpportunityCaps []int `json:"opportunity_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCapSelectionUsed bool `json:"adaptive_cap_selection_used"`
	Metrics []UP209CMetric `json:"metrics"`
}
func up209cTreated(x0 *up81cAging,endangered int,a,b string,cap int)(loss,actions int){
	m:=*x0;committed:=false;since:=0;opportunities:=0
	for start:=1;start<=56;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				opportunities=1
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if opportunities<cap{
			since++
			if since>=4{
				opportunities++
				m.query(endangered);actions++;since=0
			}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1033000+step,step,up207cPolicy(step,a,b)){return step,actions}
		}
	}
	return 57,actions
}
func RunUP209C()(UP209CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_shifted","no_refresh","hostile_shield"},
		{"hostile_no_shifted","hostile_shield","no_refresh"},
		{"alternating_fixed_shifted","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_shifted","fixed_offset_refresh","alternating_shield"},
	}
	caps:=[]int{6,7,8}
	res:=UP209CResult{Schema:UP209CTailCapSchema,Experiment:"UP-209C-committed-tail-cap-transfer",SourceUP208CSeal:"8f852d58ab5393cbc44dcdfa0f26a068cc9c17af",Cadence:2,SpacingIntervals:4,Horizon:56,OpportunityCaps:caps,CounterfactualOnly:true,LiveActivation:false,AdaptiveCapSelectionUsed:false}
	for _,s:=range schedules{for _,cap:=range caps{
		m:=UP209CMetric{Schedule:s.name,OpportunityCap:cap};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up207cBaseline(x,e,s.a,s.b)
			treated,actions:=up209cTreated(x,e,s.a,s.b,cap)
			if base<=56{m.BaselineLosses++}
			if treated<=56{m.TreatedLosses++;sumExt+=treated-base;nFail++}
			if base<=56&&treated>56{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if actions>0{m.ArmsWithAction++}
			m.ActionsTaken+=actions
		}}
		if nFail>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(nFail)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
