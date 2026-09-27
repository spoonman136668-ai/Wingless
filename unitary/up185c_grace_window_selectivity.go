package unitary

const UP185CGraceSelectivitySchema="wingless.up185c-grace-window-selectivity.v1"

type UP185CMetric struct{
	Cadence int `json:"cadence"`
	Horizon int `json:"horizon"`
	Grace int `json:"grace"`
	Arms int `json:"arms"`
	BaselineFailuresByHorizon int `json:"baseline_failures_by_horizon"`
	PreventedFailuresByHorizon int `json:"prevented_failures_by_horizon"`
	ApparentUnnecessaryActionArms int `json:"apparent_unnecessary_action_arms"`
	NearFutureActionArms int `json:"near_future_action_arms"`
	TrueUnnecessaryActionArms int `json:"true_unnecessary_action_arms"`
	BaselineSurvivorsBeyondGrace int `json:"baseline_survivors_beyond_grace"`
	PreventedPerTrueUnnecessary float64 `json:"prevented_per_true_unnecessary"`
	TrueUnnecessaryFractionAmongBeyondGraceSurvivors float64 `json:"true_unnecessary_fraction_among_beyond_grace_survivors"`
}
type UP185CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP184CSeal string `json:"source_up184c_seal"`
	Cadences []int `json:"cadences"`
	Horizons []int `json:"horizons"`
	GraceWindows []int `json:"grace_windows"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	AdaptiveGraceSelectionUsed bool `json:"adaptive_grace_selection_used"`
	Metrics []UP185CMetric `json:"metrics"`
}
func RunUP185C()(UP185CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};horizons:=[]int{16,17,18,19,20,21,22,23,24,25,26,27,28,29};graces:=[]int{0,2,4,8}
	res:=UP185CResult{Schema:UP185CGraceSelectivitySchema,Experiment:"UP-185C-grace-window-selectivity",SourceUP184CSeal:"0f9ca8328048efb421b7d2bb2f06ef02b842e47c",Cadences:cadences,Horizons:horizons,GraceWindows:graces,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,AdaptiveGraceSelectionUsed:false}
	for _,cad:=range cadences{for _,horizon:=range horizons{
		type armObs struct{base,treated,actions int}
		obs:=[]armObs{}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			treated,actions:=up182cRun(x,e,cad,horizon)
			obs=append(obs,armObs{base:base,treated:treated,actions:actions})
		}}
		for _,grace:=range graces{
			m:=UP185CMetric{Cadence:cad,Horizon:horizon,Grace:grace,Arms:len(obs)}
			for _,o:=range obs{
				if o.base<=horizon{
					m.BaselineFailuresByHorizon++
					if o.treated>horizon{m.PreventedFailuresByHorizon++}
					continue
				}
				if o.base>horizon+grace{m.BaselineSurvivorsBeyondGrace++}
				if o.actions==0{continue}
				m.ApparentUnnecessaryActionArms++
				if o.base<=horizon+grace{m.NearFutureActionArms++}else{m.TrueUnnecessaryActionArms++}
			}
			if m.TrueUnnecessaryActionArms>0{m.PreventedPerTrueUnnecessary=float64(m.PreventedFailuresByHorizon)/float64(m.TrueUnnecessaryActionArms)}
			if m.BaselineSurvivorsBeyondGrace>0{m.TrueUnnecessaryFractionAmongBeyondGraceSurvivors=float64(m.TrueUnnecessaryActionArms)/float64(m.BaselineSurvivorsBeyondGrace)}
			res.Metrics=append(res.Metrics,m)
		}
	}}
	return res,nil
}
