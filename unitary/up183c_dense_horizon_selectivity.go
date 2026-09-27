package unitary

const UP183CDenseSelectivitySchema="wingless.up183c-dense-horizon-selectivity.v1"

type UP183CMetric struct{
	Cadence int `json:"cadence"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineFailures int `json:"baseline_failures"`
	BaselineSurvivors int `json:"baseline_survivors"`
	InterventionActions int `json:"intervention_actions"`
	ArmsWithAnyAction int `json:"arms_with_any_action"`
	PreventedFailures int `json:"prevented_failures"`
	UnnecessaryActionArms int `json:"unnecessary_action_arms"`
	NecessaryActionArms int `json:"necessary_action_arms"`
	NoActionSurvivors int `json:"no_action_survivors"`
	PreventedPerUnnecessaryActionArm float64 `json:"prevented_per_unnecessary_action_arm"`
	UnnecessaryFractionAmongSurvivors float64 `json:"unnecessary_fraction_among_survivors"`
}
type UP183CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP182CSeal string `json:"source_up182c_seal"`
	Cadences []int `json:"cadences"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveHorizonSelectionUsed bool `json:"adaptive_horizon_selection_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP183CMetric `json:"metrics"`
}
func RunUP183C()(UP183CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};horizons:=[]int{16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31}
	res:=UP183CResult{Schema:UP183CDenseSelectivitySchema,Experiment:"UP-183C-dense-horizon-selectivity",SourceUP182CSeal:"4526fdb17de1bbdd873e79968177f6416a79225e",Cadences:cadences,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,AdaptiveHorizonSelectionUsed:false,WarningThresholdChanged:false}
	for _,cad:=range cadences{for _,horizon:=range horizons{
		m:=UP183CMetric{Cadence:cad,Horizon:horizon}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			bf:=base<=horizon
			if bf{m.BaselineFailures++}else{m.BaselineSurvivors++}
			treated,actions:=up182cRun(x,e,cad,horizon)
			m.InterventionActions+=actions
			if actions>0{
				m.ArmsWithAnyAction++
				if bf{m.NecessaryActionArms++}else{m.UnnecessaryActionArms++}
			}else if !bf{m.NoActionSurvivors++}
			if bf&&treated>horizon{m.PreventedFailures++}
		}}
		if m.UnnecessaryActionArms>0{m.PreventedPerUnnecessaryActionArm=float64(m.PreventedFailures)/float64(m.UnnecessaryActionArms)}
		if m.BaselineSurvivors>0{m.UnnecessaryFractionAmongSurvivors=float64(m.UnnecessaryActionArms)/float64(m.BaselineSurvivors)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
