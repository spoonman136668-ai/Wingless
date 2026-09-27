package unitary

const UP181CStageTwoDelaySchema="wingless.up181c-stage-two-delay-boundary.v1"

type UP181CMetric struct{
	Cadence int `json:"cadence"`
	StageTwoDelay int `json:"stage_two_delay"`
	Arms int `json:"arms"`
	StageOneActions int `json:"stage_one_actions"`
	StageTwoActions int `json:"stage_two_actions"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP181CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP180CSeal string `json:"source_up180c_seal"`
	Cadences []int `json:"cadences"`
	StageOneDelay int `json:"stage_one_delay"`
	StageTwoDelays []int `json:"stage_two_delays"`
	Policy string `json:"policy"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	AdaptiveTargetUsed bool `json:"adaptive_target_used"`
	Metrics []UP181CMetric `json:"metrics"`
}

func RunUP181C()(UP181CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};delays:=[]int{0,1,2}
	res:=UP181CResult{Schema:UP181CStageTwoDelaySchema,Experiment:"UP-181C-stage-two-delay-boundary",SourceUP180CSeal:"176720ea845a6599c42a45747b941a3679badf06",Cadences:cadences,StageOneDelay:0,StageTwoDelays:delays,Policy:"no_refresh",MaxActions:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveDelayUsed:false,AdaptiveTargetUsed:false}
	for _,cad:=range cadences{for _,d2:=range delays{
		m:=UP181CMetric{Cadence:cad,StageTwoDelay:d2};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			treated,a1,a2:=up180cRun(x,e,cad,0,d2);m.StageOneActions+=a1;m.StageTwoActions+=a2
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		total:=m.StageOneActions+m.StageTwoActions
		if total>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(total)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
