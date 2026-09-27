package unitary

const UP213CEighthTimingSchema="wingless.up213c-eighth-action-timing-window.v1"

type UP213CMetric struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	Mode string `json:"mode"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	AcceleratedLosses int `json:"accelerated_losses"`
	ActionsTaken int `json:"actions_taken"`
}
type UP213CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP212CSeal string `json:"source_up212c_seal"`
	Horizons []int `json:"horizons"`
	Modes []string `json:"modes"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveTimingUsed bool `json:"adaptive_timing_used"`
	Metrics []UP213CMetric `json:"metrics"`
}
func up213cTreated(x0 *up81cAging,endangered,advance,horizon int,a,b,mode string)(loss,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&mode!="none"{
				need:=4
				if mode=="early"{need=3}else if mode=="late"{need=5}
				if since>=need{m.query(endangered);actions++;since=0}
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1038000+step,step,up210cPolicy(step,advance,a,b)){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP213C()(UP213CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{0,4};horizons:=[]int{66,68,70,72};modes:=[]string{"none","early","nominal","late"}
	res:=UP213CResult{Schema:UP213CEighthTimingSchema,Experiment:"UP-213C-eighth-action-timing-window",SourceUP212CSeal:"c4f349bf8341b725d435c3109c24c81a00feb438",Horizons:horizons,Modes:modes,CounterfactualOnly:true,LiveActivation:false,AdaptiveTimingUsed:false}
	for _,s:=range schedules{for _,advance:=range advances{for _,h:=range horizons{for _,mode:=range modes{
		m:=UP213CMetric{Schedule:s.name,PhaseAdvanceWrites:advance,Horizon:h,Mode:mode}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			base:=up211cBaseline(x,e,advance,h,s.a,s.b)
			treated,actions:=up213cTreated(x,e,advance,h,s.a,s.b,mode)
			if base<=h{m.BaselineLosses++};if treated<=h{m.TreatedLosses++};if base<=h&&treated>h{m.PreventedLosses++};if treated<base{m.AcceleratedLosses++};m.ActionsTaken+=actions
		}}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
