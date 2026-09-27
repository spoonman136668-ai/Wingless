package unitary

const UP212CDenseBoundarySchema="wingless.up212c-dense-tail-cap-boundary.v1"

type UP212CMetric struct{
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
type UP212CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP211CSeal string `json:"source_up211c_seal"`
	PhaseAdvances []int `json:"phase_advances"`
	Horizons []int `json:"horizons"`
	OpportunityCaps []int `json:"opportunity_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP212CMetric `json:"metrics"`
}
func RunUP212C()(UP212CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{0,4};horizons:=[]int{64,66,68,70,72};caps:=[]int{7,8}
	res:=UP212CResult{Schema:UP212CDenseBoundarySchema,Experiment:"UP-212C-dense-tail-cap-boundary",SourceUP211CSeal:"25e0427860cd8532401f40cecf3f6f4559da3778",PhaseAdvances:advances,Horizons:horizons,OpportunityCaps:caps,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,horizon:=range horizons{for _,cap:=range caps{
		m:=UP212CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:horizon,OpportunityCap:cap};sumExt,nFail:=0,0
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
