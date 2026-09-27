package unitary

const UP169CCorrectionDelaySchema="wingless.up169c-correction-delay.v1"

type UP169CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Confirmation int `json:"confirmation"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	BaselineLosses int `json:"baseline_losses"`
	TreatedLosses int `json:"treated_losses"`
	TotalDelayWrites int `json:"total_delay_writes"`
	MeanDelayWrites float64 `json:"mean_delay_writes"`
	ArmsDelayed int `json:"arms_delayed"`
	ArmsAdvanced int `json:"arms_advanced"`
	MaxDelayWrites int `json:"max_delay_writes"`
	NoPositiveDelayActions int `json:"no_positive_delay_actions"`
}
type UP169CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP168CSeal string `json:"source_up168c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Confirmations []int `json:"confirmations"`
	Interventions []string `json:"interventions"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	EvaluationWrites int `json:"evaluation_writes"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	RepeatedMaintenanceUsed bool `json:"repeated_maintenance_used"`
	AdaptiveTriggerUsed bool `json:"adaptive_trigger_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP169CMetric `json:"metrics"`
}

func up169cRun(x0 *up81cAging,endangered,cadence,confirmation int,policy,intervention string)(lossStep int,acted bool){
	m:=*x0;streak:=0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered);warn:=h<=cadence
		if warn{streak++}else{streak=0}
		if intervention!="none"&&!acted&&warn&&streak>=confirmation{
			if intervention=="targeted_refresh"{m.query(endangered);acted=true}
			if intervention=="sham_refresh"{acted=up168cShamQuery(&m,endangered)}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,990000+step,step,policy){return step,acted}
		}
	}
	return 65,acted
}
func RunUP169C()(UP169CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};confirms:=[]int{1,2};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP169CResult{Schema:UP169CCorrectionDelaySchema,Experiment:"UP-169C-correction-delay",SourceUP168CSeal:"71ba926d3024ffac81b27092a2a30a98aa90fba7",Policies:policies,Cadences:cadences,Confirmations:confirms,Interventions:interventions,MaxActionsPerArm:1,EvaluationWrites:64,CounterfactualOnly:true,LiveActivation:false,RepeatedMaintenanceUsed:false,AdaptiveTriggerUsed:false,WarningThresholdChanged:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,confirm:=range confirms{for _,intervention:=range interventions{
		m:=UP169CMetric{Policy:policy,Cadence:cad,Confirmation:confirm,Intervention:intervention}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			baseStep,_:=up169cRun(x,endangered,cad,confirm,policy,"none")
			treatedStep,acted:=up169cRun(x,endangered,cad,confirm,policy,intervention)
			if baseStep<=64{m.BaselineLosses++};if treatedStep<=64{m.TreatedLosses++};if acted{m.ActionsTaken++}
			delay:=treatedStep-baseStep;m.TotalDelayWrites+=delay
			if delay>0{m.ArmsDelayed++;if delay>m.MaxDelayWrites{m.MaxDelayWrites=delay}}
			if delay<0{m.ArmsAdvanced++}
			if acted&&delay<=0{m.NoPositiveDelayActions++}
		}}
		if m.Arms>0{m.MeanDelayWrites=float64(m.TotalDelayWrites)/float64(m.Arms)}
		res.Metrics=append(res.Metrics,m)
	}}}}
	return res,nil
}
