package unitary

const UP178CWarningStageSchema="wingless.up178c-warning-stage-ablation.v1"

type UP178CMetric struct{
	Cadence int `json:"cadence"`
	Arm string `json:"arm"`
	Arms int `json:"arms"`
	ActionsTaken int `json:"actions_taken"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP178CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP177CSeal string `json:"source_up177c_seal"`
	Cadences []int `json:"cadences"`
	Arms []string `json:"arms"`
	DelayIntervals int `json:"delay_intervals"`
	Policy string `json:"policy"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveStageUsed bool `json:"adaptive_stage_used"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	Metrics []UP178CMetric `json:"metrics"`
}
func up178cSecondOnly(x0 *up81cAging,endangered,cadence int)(lossStep,actions int){
	m:=*x0;warningCount:=0;pending:=false
	for start:=1;start<=64;start+=cadence{
		if pending&&actions==0{m.query(endangered);actions++;pending=false}
		if actions==0&&!pending{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{warningCount++;if warningCount==2{pending=true}}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1006000+step,step,"no_refresh"){return step,actions}
		}
	}
	return 65,actions
}
func RunUP178C()(UP178CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};cadences:=[]int{2,4};arms:=[]string{"first_only","second_warning_only","both","sham_both"}
	res:=UP178CResult{Schema:UP178CWarningStageSchema,Experiment:"UP-178C-warning-stage-ablation",SourceUP177CSeal:"5dd0dd36b223926006a798baeb069c9660ff529f",Cadences:cadences,Arms:arms,DelayIntervals:1,Policy:"no_refresh",CounterfactualOnly:true,LiveActivation:false,AdaptiveStageUsed:false,AdaptiveDelayUsed:false}
	for _,cad:=range cadences{for _,armName:=range arms{
		m:=UP178CMetric{Cadence:cad,Arm:armName};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			var treated,actions int
			switch armName{
			case "first_only": treated,actions=up177cRun(x,e,cad,1,"targeted_refresh")
			case "second_warning_only": treated,actions=up178cSecondOnly(x,e,cad)
			case "both": treated,actions=up177cRun(x,e,cad,2,"targeted_refresh")
			default: treated,actions=up177cRun(x,e,cad,2,"sham_refresh")
			}
			m.ActionsTaken+=actions
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if m.ActionsTaken>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(m.ActionsTaken)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
