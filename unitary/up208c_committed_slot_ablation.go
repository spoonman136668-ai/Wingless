package unitary

const UP208CSlotAblationSchema="wingless.up208c-committed-slot-ablation.v1"

type UP208CMetric struct{
	Schedule string `json:"schedule"`
	OmittedOrdinal int `json:"omitted_ordinal"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
}
type UP208CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP207CSeal string `json:"source_up207c_seal"`
	PhaseShiftWrites int `json:"phase_shift_writes"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	MaxScheduledOpportunities int `json:"max_scheduled_opportunities"`
	Horizon int `json:"horizon"`
	OmittedOrdinals []int `json:"omitted_ordinals"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ReplacementActionUsed bool `json:"replacement_action_used"`
	AdaptiveOmissionUsed bool `json:"adaptive_omission_used"`
	Metrics []UP208CMetric `json:"metrics"`
}
func up208cTreated(x0 *up81cAging,endangered int,a,b string,omit int)(loss,actions int){
	m:=*x0;committed:=false;since:=0;opportunities:=0
	for start:=1;start<=56;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				opportunities=1
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if opportunities<8{
			since++
			if since>=4{
				opportunities++
				if opportunities!=omit{m.query(endangered);actions++}
				since=0
			}
		}
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1032000+step,step,up207cPolicy(step,a,b)){return step,actions}
		}
	}
	return 57,actions
}
func RunUP208C()(UP208CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_shifted","no_refresh","hostile_shield"},
		{"hostile_no_shifted","hostile_shield","no_refresh"},
		{"alternating_fixed_shifted","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_shifted","fixed_offset_refresh","alternating_shield"},
	}
	omits:=[]int{0,2,3,4,5,6,7,8}
	res:=UP208CResult{Schema:UP208CSlotAblationSchema,Experiment:"UP-208C-committed-slot-ablation",SourceUP207CSeal:"9023c20d428e3f742b2f9840bf67e52164ec0f3b",PhaseShiftWrites:4,Cadence:2,SpacingIntervals:4,MaxScheduledOpportunities:8,Horizon:56,OmittedOrdinals:omits,CounterfactualOnly:true,LiveActivation:false,ReplacementActionUsed:false,AdaptiveOmissionUsed:false}
	for _,s:=range schedules{for _,omit:=range omits{
		m:=UP208CMetric{Schedule:s.name,OmittedOrdinal:omit};sumExt,nFail:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up207cBaseline(x,e,s.a,s.b)
			treated,actions:=up208cTreated(x,e,s.a,s.b,omit)
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
