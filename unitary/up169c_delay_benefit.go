package unitary

const UP169CDelayBenefitSchema="wingless.up169c-delay-benefit.v1"

type UP169CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Intervention string `json:"intervention"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	MeanBaselineLossStep float64 `json:"mean_baseline_loss_step"`
	MeanTreatedLossStep float64 `json:"mean_treated_loss_step"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
	MinLossDelayWrites int `json:"min_loss_delay_writes"`
	MaxLossDelayWrites int `json:"max_loss_delay_writes"`
	DelayedArms int `json:"delayed_arms"`
	UnchangedArms int `json:"unchanged_arms"`
	AcceleratedArms int `json:"accelerated_arms"`
}
type UP169CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP168CSeal string `json:"source_up168c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Interventions []string `json:"interventions"`
	Confirmation int `json:"confirmation"`
	MaxActionsPerArm int `json:"max_actions_per_arm"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	RepeatedMaintenanceUsed bool `json:"repeated_maintenance_used"`
	AdaptiveTriggerUsed bool `json:"adaptive_trigger_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP169CMetric `json:"metrics"`
}
func up169cBaselineLossStep(x0 *up81cAging,endangered,cadence int,policy string)int{
	m:=*x0
	for start:=1;start<=64;start+=cadence{
		for j:=0;j<cadence&&start+j<=64;j++{step:=start+j;if up165cRealStep(&m,endangered,990000+step,step,policy){return step}}
	}
	return 65
}
func up169cTreatedLossStep(x0 *up81cAging,endangered,cadence int,policy,intervention string)(int,bool){
	m:=*x0;acted:=false
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered);warn:=h<=cadence
		if !acted&&warn{
			if intervention=="targeted_refresh"{m.query(endangered);acted=true}
			if intervention=="sham_refresh"{acted=up168cShamQuery(&m,endangered)}
		}
		for j:=0;j<cadence&&start+j<=64;j++{step:=start+j;if up165cRealStep(&m,endangered,991000+step,step,policy){return step,acted}}
	}
	return 65,acted
}
func RunUP169C()(UP169CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};interventions:=[]string{"targeted_refresh","sham_refresh"}
	res:=UP169CResult{Schema:UP169CDelayBenefitSchema,Experiment:"UP-169C-delay-benefit",SourceUP168CSeal:"71ba926d3024ffac81b27092a2a30a98aa90fba7",Policies:policies,Cadences:cadences,Interventions:interventions,Confirmation:1,MaxActionsPerArm:1,CounterfactualOnly:true,LiveActivation:false,RepeatedMaintenanceUsed:false,AdaptiveTriggerUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,intervention:=range interventions{
		m:=UP169CMetric{Policy:policy,Cadence:cad,Intervention:intervention,MinLossDelayWrites:1<<30};sumBase,sumTreated,sumDelay:=0,0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,endangered,cad,policy);treated,acted:=up169cTreatedLossStep(x,endangered,cad,policy,intervention)
			if acted{m.ActionsTaken++};d:=treated-base;sumBase+=base;sumTreated+=treated;sumDelay+=d
			if d<m.MinLossDelayWrites{m.MinLossDelayWrites=d};if d>m.MaxLossDelayWrites{m.MaxLossDelayWrites=d}
			if d>0{m.DelayedArms++}else if d<0{m.AcceleratedArms++}else{m.UnchangedArms++}
		}}
		if m.Arms>0{m.MeanBaselineLossStep=float64(sumBase)/float64(m.Arms);m.MeanTreatedLossStep=float64(sumTreated)/float64(m.Arms);m.MeanLossDelayWrites=float64(sumDelay)/float64(m.Arms)}
		if m.MinLossDelayWrites==1<<30{m.MinLossDelayWrites=0}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
