package unitary

const UP182CHorizonSelectivitySchema="wingless.up182c-horizon-selectivity.v1"

type UP182CMetric struct{
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
}
type UP182CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP181CSeal string `json:"source_up181c_seal"`
	Cadences []int `json:"cadences"`
	Horizons []int `json:"horizons"`
	MaxActions int `json:"max_actions"`
	Policy string `json:"policy"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveHorizonUsed bool `json:"adaptive_horizon_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	Metrics []UP182CMetric `json:"metrics"`
}
func up182cRun(x0 *up81cAging,endangered,cadence,horizon int)(lossStep,actions int){
	m:=*x0
	stage:=1
	pendingSecond:=false
	for start:=1;start<=horizon;start+=cadence{
		actedThis:=false
		if pendingSecond&&stage==2{
			m.query(endangered);actions++;stage++;pendingSecond=false;actedThis=true
		}
		if !actedThis&&stage<=2{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{
				if stage==1{
					m.query(endangered);actions++;stage++;actedThis=true
				}else{
					pendingSecond=true
				}
			}
		}
		for j:=0;j<cadence&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1012000+step,step,"no_refresh"){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP182C()(UP182CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};horizons:=[]int{16,24,32,40,48,56,64}
	res:=UP182CResult{Schema:UP182CHorizonSelectivitySchema,Experiment:"UP-182C-horizon-selectivity",SourceUP181CSeal:"a8c6b5bb2df4761e0fb00f1e0666543304835970",Cadences:cadences,Horizons:horizons,MaxActions:2,Policy:"no_refresh",CounterfactualOnly:true,LiveActivation:false,AdaptiveHorizonUsed:false,WarningThresholdChanged:false}
	for _,cad:=range cadences{for _,horizon:=range horizons{
		m:=UP182CMetric{Cadence:cad,Horizon:horizon}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			baselineFailure:=base<=horizon
			if baselineFailure{m.BaselineFailures++}else{m.BaselineSurvivors++}
			treated,actions:=up182cRun(x,e,cad,horizon)
			m.InterventionActions+=actions
			if actions>0{
				m.ArmsWithAnyAction++
				if baselineFailure{m.NecessaryActionArms++}else{m.UnnecessaryActionArms++}
			}else if !baselineFailure{m.NoActionSurvivors++}
			if baselineFailure&&treated>horizon{m.PreventedFailures++}
		}}
		if m.UnnecessaryActionArms>0{m.PreventedPerUnnecessaryActionArm=float64(m.PreventedFailures)/float64(m.UnnecessaryActionArms)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
