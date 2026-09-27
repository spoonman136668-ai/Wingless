package unitary

const UP173CFixedScheduleSchema="wingless.up173c-fixed-schedule-sweep.v1"

type UP173CSchedule struct{
	Name string `json:"name"`
	WriteStarts []int `json:"write_starts"`
}
type UP173CMetric struct{
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	Losses int `json:"losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP173CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP172CSeal string `json:"source_up172c_seal"`
	InitialHands []int `json:"initial_hands"`
	Cadence int `json:"cadence"`
	Policy string `json:"policy"`
	Schedules []UP173CSchedule `json:"schedules"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveScheduleSearchUsed bool `json:"adaptive_schedule_search_used"`
	AdaptiveActionBudgetUsed bool `json:"adaptive_action_budget_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP173CMetric `json:"metrics"`
	BestFixedSchedule string `json:"best_fixed_schedule"`
	BestFixedPreventedLosses int `json:"best_fixed_prevented_losses"`
	WarningPreventedLosses int `json:"warning_prevented_losses"`
}
func up173cFixed(x0 *up81cAging,endangered int,starts []int)(lossStep,actions int){
	m:=*x0
	for start:=1;start<=64;start+=2{
		for _,s:=range starts{if start==s&&actions<2{m.query(endangered);actions++;break}}
		for j:=0;j<2&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,996000+step,step,"no_refresh"){return step,actions}
		}
	}
	return 65,actions
}
func RunUP173C()(UP173CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	schedules:=[]UP173CSchedule{
		{Name:"fixed_01_03",WriteStarts:[]int{1,3}},
		{Name:"fixed_09_11",WriteStarts:[]int{9,11}},
		{Name:"fixed_17_19",WriteStarts:[]int{17,19}},
		{Name:"fixed_25_27",WriteStarts:[]int{25,27}},
		{Name:"fixed_33_35",WriteStarts:[]int{33,35}},
		{Name:"fixed_41_43",WriteStarts:[]int{41,43}},
		{Name:"fixed_49_51",WriteStarts:[]int{49,51}},
		{Name:"fixed_57_59",WriteStarts:[]int{57,59}},
	}
	res:=UP173CResult{Schema:UP173CFixedScheduleSchema,Experiment:"UP-173C-fixed-schedule-sweep",SourceUP172CSeal:"0b8e5a7a6cf10d3b88301fc7e667b8ab2ef1c41f",InitialHands:hands,Cadence:2,Policy:"no_refresh",Schedules:schedules,MaxActionsPerArm:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveScheduleSearchUsed:false,AdaptiveActionBudgetUsed:false,WarningThresholdChanged:false}
	type arm struct{x *up81cAging;endangered int}
	arms:=[]arm{}
	for _,c:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,c);if ok{arms=append(arms,arm{x:x,endangered:e})}}}
	base:=make([]int,len(arms))
	for i,a:=range arms{base[i]=up169cBaselineLossStep(a.x,a.endangered,2,"no_refresh")}
	eval:=func(name string,fn func(arm)(int,int))UP173CMetric{
		m:=UP173CMetric{Intervention:name,Arms:len(arms)};delaySum,delayN:=0,0
		for i,a:=range arms{
			treated,actions:=fn(a);m.ActionsTaken+=actions
			if treated<=64{m.Losses++;delaySum+=treated-base[i];delayN++}else if base[i]<=64{m.PreventedLosses++}
		}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		return m
	}
	w:=eval("warning_targeted",func(a arm)(int,int){return up170cRun(a.x,a.endangered,2,"no_refresh","targeted_refresh")})
	res.WarningPreventedLosses=w.PreventedLosses;res.Metrics=append(res.Metrics,w)
	bestName:="";best:=-1
	for _,s:=range schedules{
		ss:=s
		m:=eval(ss.Name,func(a arm)(int,int){return up173cFixed(a.x,a.endangered,ss.WriteStarts)})
		res.Metrics=append(res.Metrics,m)
		if m.PreventedLosses>best{best=m.PreventedLosses;bestName=ss.Name}
	}
	n:=eval("no_action",func(a arm)(int,int){return up169cBaselineLossStep(a.x,a.endangered,2,"no_refresh"),0})
	res.Metrics=append(res.Metrics,n);res.BestFixedSchedule=bestName;res.BestFixedPreventedLosses=best
	return res,nil
}
