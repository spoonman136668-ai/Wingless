package unitary

const UP180CStageDelaySchema="wingless.up180c-stage-delay-matrix.v1"

type UP180CMetric struct{
	Cadence int `json:"cadence"`
	StageOneDelay int `json:"stage_one_delay"`
	StageTwoDelay int `json:"stage_two_delay"`
	Arms int `json:"arms"`
	StageOneActions int `json:"stage_one_actions"`
	StageTwoActions int `json:"stage_two_actions"`
	TreatedLosses int `json:"treated_losses"`
	PreventedLosses int `json:"prevented_losses"`
	PreventedPerAction float64 `json:"prevented_per_action"`
	MeanLossDelayWrites float64 `json:"mean_loss_delay_writes"`
}
type UP180CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP179CSeal string `json:"source_up179c_seal"`
	Cadences []int `json:"cadences"`
	DelayPairs [][2]int `json:"delay_pairs"`
	Policy string `json:"policy"`
	MaxActions int `json:"max_actions"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDelayUsed bool `json:"adaptive_delay_used"`
	AdaptiveTargetUsed bool `json:"adaptive_target_used"`
	Metrics []UP180CMetric `json:"metrics"`
}
func up180cRun(x0 *up81cAging,endangered,cadence,d1,d2 int)(lossStep,a1,a2 int){
	m:=*x0
	stage:=1
	pending:=false
	pendingDelay:=0
	for start:=1;start<=64;start+=cadence{
		actedThis:=false
		if pending{
			if pendingDelay>0{pendingDelay--}
			if pendingDelay==0{
				m.query(endangered)
				if stage==1{a1++}else{a2++}
				stage++;pending=false;actedThis=true
			}
		}
		if !actedThis&&!pending&&stage<=2{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{
				delay:=d1
				if stage==2{delay=d2}
				if delay==0{
					m.query(endangered)
					if stage==1{a1++}else{a2++}
					stage++;actedThis=true
				}else{
					pending=true;pendingDelay=delay
				}
			}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1010000+step,step,"no_refresh"){return step,a1,a2}
		}
	}
	return 65,a1,a2
}
func RunUP180C()(UP180CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};pairs:=[][2]int{{0,0},{0,1},{1,0},{1,1}}
	res:=UP180CResult{Schema:UP180CStageDelaySchema,Experiment:"UP-180C-stage-delay-matrix",SourceUP179CSeal:"88fe5debf5424c205e0274bee40a9c4b61ff1eaa",Cadences:cadences,DelayPairs:pairs,Policy:"no_refresh",MaxActions:2,CounterfactualOnly:true,LiveActivation:false,AdaptiveDelayUsed:false,AdaptiveTargetUsed:false}
	for _,cad:=range cadences{for _,p:=range pairs{
		m:=UP180CMetric{Cadence:cad,StageOneDelay:p[0],StageTwoDelay:p[1]};delaySum,delayN:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			treated,a1,a2:=up180cRun(x,e,cad,p[0],p[1]);m.StageOneActions+=a1;m.StageTwoActions+=a2
			if treated<=64{m.TreatedLosses++;delaySum+=treated-base;delayN++}else if base<=64{m.PreventedLosses++}
		}}
		totalActions:=m.StageOneActions+m.StageTwoActions
		if totalActions>0{m.PreventedPerAction=float64(m.PreventedLosses)/float64(totalActions)}
		if delayN>0{m.MeanLossDelayWrites=float64(delaySum)/float64(delayN)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
