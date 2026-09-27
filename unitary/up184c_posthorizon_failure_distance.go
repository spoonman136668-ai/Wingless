package unitary

const UP184CPostHorizonSchema="wingless.up184c-posthorizon-failure-distance.v1"

type UP184CMetric struct{
	Cadence int `json:"cadence"`
	Horizon int `json:"horizon"`
	BaselineSurvivors int `json:"baseline_survivors"`
	UnnecessaryActionArms int `json:"unnecessary_action_arms"`
	FailWithin1 int `json:"fail_within_1"`
	FailWithin2 int `json:"fail_within_2"`
	FailWithin4 int `json:"fail_within_4"`
	FailWithin8 int `json:"fail_within_8"`
	FailWithin16 int `json:"fail_within_16"`
	SurviveBeyond64 int `json:"survive_beyond_64"`
	ObservedPostHorizonFailures int `json:"observed_post_horizon_failures"`
	MinPostHorizonDistance int `json:"min_post_horizon_distance"`
	MaxPostHorizonDistance int `json:"max_post_horizon_distance"`
	MeanPostHorizonDistance float64 `json:"mean_post_horizon_distance"`
}
type UP184CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP183CSeal string `json:"source_up183c_seal"`
	Cadences []int `json:"cadences"`
	Horizons []int `json:"horizons"`
	BaselineFollowupLimit int `json:"baseline_followup_limit"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	AdaptiveHorizonSelectionUsed bool `json:"adaptive_horizon_selection_used"`
	Metrics []UP184CMetric `json:"metrics"`
}
func RunUP184C()(UP184CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};horizons:=[]int{16,17,18,19,20,21,22,23,24,25,26,27,28,29}
	res:=UP184CResult{Schema:UP184CPostHorizonSchema,Experiment:"UP-184C-posthorizon-failure-distance",SourceUP183CSeal:"86cf3e5da3b7e3ab2bdcfabfff333d3f63d3e17c",Cadences:cadences,Horizons:horizons,BaselineFollowupLimit:64,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,AdaptiveHorizonSelectionUsed:false}
	for _,cad:=range cadences{for _,horizon:=range horizons{
		m:=UP184CMetric{Cadence:cad,Horizon:horizon}
		sumDist:=0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			if base<=horizon{continue}
			m.BaselineSurvivors++
			_,actions:=up182cRun(x,e,cad,horizon)
			if actions==0{continue}
			m.UnnecessaryActionArms++
			if base>64{
				m.SurviveBeyond64++
				continue
			}
			d:=base-horizon
			m.ObservedPostHorizonFailures++
			sumDist+=d
			if m.MinPostHorizonDistance==0||d<m.MinPostHorizonDistance{m.MinPostHorizonDistance=d}
			if d>m.MaxPostHorizonDistance{m.MaxPostHorizonDistance=d}
			if d<=1{m.FailWithin1++}
			if d<=2{m.FailWithin2++}
			if d<=4{m.FailWithin4++}
			if d<=8{m.FailWithin8++}
			if d<=16{m.FailWithin16++}
		}}
		if m.ObservedPostHorizonFailures>0{m.MeanPostHorizonDistance=float64(sumDist)/float64(m.ObservedPostHorizonFailures)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
