package unitary

const UP179CStageTargetSchema="wingless.up179c-stage-target-ablation.v1"

type UP179CMetric struct{
	Cadence int `json:"cadence"`
	Arm string `json:"arm"`
	Arms int `json:"arms"`
	StageOneActions int `json:"stage_one_actions"`
	StageTwoActions int `json:"stage_two_actions"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP179CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP178CSeal string `json:"source_up178c_seal"`
	Cadences []int `json:"cadences"`
	TargetArms []string `json:"target_arms"`
	DelayIntervals int `json:"delay_intervals"`
	Policy string `json:"policy"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveTargetUsed bool `json:"adaptive_target_used"`
	AdaptiveTimingUsed bool `json:"adaptive_timing_used"`
	Metrics []UP179CMetric `json:"metrics"`
}
func up179cTargeted(name string,stage int)bool{
	if name=="targeted_targeted"{return true}
	if name=="targeted_sham"{return stage==1}
	if name=="sham_targeted"{return stage==2}
	return false
}
func up179cAct(m *up81cAging,endangered int,targeted bool)bool{
	if targeted{m.query(endangered);return true}
	return up168cShamQuery(m,endangered)
}
func up179cRun(x0 *up81cAging,endangered,cadence int,arm string)(lossStep,stage1,stage2 int){
	m:=*x0;stage:=1;pending:=false
	for start:=1;start<=64;start+=cadence{
		actedThis:=false
		if pending&&stage<=2{
			ok:=up179cAct(&m,endangered,up179cTargeted(arm,stage))
			if ok{if stage==1{stage1++}else{stage2++};stage++;actedThis=true}
			pending=false
		}
		if !actedThis&&!pending&&stage<=2{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{pending=true}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1008000+step,step,"no_refresh"){return step,stage1,stage2}
		}
	}
	return 65,stage1,stage2
}
func RunUP179C()(UP179CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};arms:=[]string{"targeted_targeted","targeted_sham","sham_targeted","sham_sham"}
	res:=UP179CResult{Schema:UP179CStageTargetSchema,Experiment:"UP-179C-stage-target-ablation",SourceUP178CSeal:"1f456fdca38cba2f5a4535af5241105304b57da0",Cadences:cadences,TargetArms:arms,DelayIntervals:1,Policy:"no_refresh",MaxActions:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveTargetUsed:false,AdaptiveTimingUsed:false}
	for _,cad:=range cadences{for _,arm:=range arms{
		m:=UP179CMetric{Cadence:cad,Arm:arm};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			treated,a1,a2:=up179cRun(x,e,cad,arm);m.StageOneActions+=a1;m.StageTwoActions+=a2
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
