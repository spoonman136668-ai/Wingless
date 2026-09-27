package unitary

const UP168CBoundedCorrectionSchema="wingless.up168c-bounded-correction.v1"

type UP168CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Confirmation int `json:"confirmation"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	HarmedArms int `json:"harmed_arms"`
	WastedActions int `json:"wasted_actions"`
	PreventedPerAction float64 `json:"prevented_per_action"`
}
type UP168CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP167CSeal string `json:"source_up167c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Confirmations []int `json:"confirmations"`
	Interventions []string `json:"interventions"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveTriggerUsed bool `json:"adaptive_trigger_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP168CMetric `json:"metrics"`
}
func up168cShamQuery(m *up81cAging,endangered int)bool{
	slot:=up151cPredict(m.hand,m.age)
	if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key);return true}
	for i:=range m.entries{if m.entries[i].used&&m.entries[i].key!=endangered{m.query(m.entries[i].key);return true}}
	return false
}
func up168cRun(x0 *up81cAging,endangered,cadence,confirmation int,policy,intervention string)(lost,acted bool){
	m:=*x0;streak:=0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		if warn{streak++}else{streak=0}
		if !acted&&warn&&streak>=confirmation{
			if intervention=="targeted_refresh"{m.query(endangered);acted=true}
			if intervention=="sham_refresh"{acted=up168cShamQuery(&m,endangered)}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,980000+step,step,policy){return true,acted}
		}
	}
	return false,acted
}
func RunUP168C()(UP168CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};confirms:=[]int{1,2};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP168CResult{Schema:UP168CBoundedCorrectionSchema,Experiment:"UP-168C-bounded-correction",SourceUP167CSeal:"8acfc2c4e35967728992598246194efa38cd79e2",Policies:policies,Cadences:cadences,Confirmations:confirms,Interventions:interventions,MaxActionsPerArm:1,CounterfactualOnly:true,LiveActivation:false,AdaptiveTriggerUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,confirm:=range confirms{for _,intervention:=range interventions{
		m:=UP168CMetric{Policy:policy,Cadence:cad,Confirmation:confirm,Intervention:intervention}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			baseLost,_,_:=up165cRun(x,endangered,cad,policy)
			treatedLost,acted:=up168cRun(x,endangered,cad,confirm,policy,intervention)
			if baseLost{m.BaselineLosses++}
			if acted{m.ActionsTaken++}
			if treatedLost{m.TreatedLosses++}
			if baseLost&&!treatedLost{m.PreventedLosses++}
			if !baseLost&&treatedLost{m.HarmedArms++}
			if acted&&!(baseLost&&!treatedLost){m.WastedActions++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
