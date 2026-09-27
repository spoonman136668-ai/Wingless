package unitary

const UP219CNinthThresholdSchema="wingless.up219c-ninth-critical-vs-near.v1"

type UP219CMetric struct{
	Horizon int `json:"horizon"`
	NinthTrigger string `json:"ninth_trigger"`
	Arms int `json:"arms"`
	Cap8Losses int `json:"cap8_losses"`
	Cap9Losses int `json:"cap9_losses"`
	NinthActions int `json:"ninth_actions"`
	NinthRescues int `json:"ninth_rescues"`
	UnnecessaryNinth int `json:"unnecessary_ninth"`
	NecessaryNinthAttempts int `json:"necessary_ninth_attempts"`
	IneffectiveNinth int `json:"ineffective_ninth"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP219CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP218CSeal string `json:"source_up218c_seal"`
	CriticalHorizon int `json:"critical_horizon"`
	NearHorizon int `json:"near_horizon"`
	MaxActions int `json:"max_actions"`
	Horizons []int `json:"horizons"`
	Triggers []string `json:"triggers"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ScheduleRetuningUsed bool `json:"schedule_retuning_used"`
	Metrics []UP219CMetric `json:"metrics"`
}
func up219cTreated(x0 *up81cAging,endangered,horizon,ninthThreshold int)(loss,actions,ninth int){
	m:=*x0;committed:=false;since:=0;eighthBoundary:=-1
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;since=0;eighthBoundary=start
			}else if actions==8&&eighthBoundary>=0&&start>eighthBoundary&&up161cAdversarial(&m,endangered)<=ninthThreshold{
				m.query(endangered);actions++;ninth++;since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1042000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){return step,actions,ninth}
		}
	}
	return horizon+1,actions,ninth
}
func RunUP219C()(UP219CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	horizons:=[]int{74,76,78,80};triggers:=[]string{"critical","near"}
	res:=UP219CResult{Schema:UP219CNinthThresholdSchema,Experiment:"UP-219C-ninth-critical-vs-near",SourceUP218CSeal:"0c3385e8c39eba77a10478cf5f4f9f448036687e",CriticalHorizon:2,NearHorizon:4,MaxActions:9,Horizons:horizons,Triggers:triggers,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,ScheduleRetuningUsed:false}
	for _,h:=range horizons{for _,trigger:=range triggers{
		threshold:=2;if trigger=="near"{threshold=4}
		m:=UP219CMetric{Horizon:h,NinthTrigger:trigger}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap8,_,_:=up219cTreated(x,e,h,0)
			cap9,_,ninth:=up219cTreated(x,e,h,threshold)
			l8:=cap8<=h;l9:=cap9<=h
			if l8{m.Cap8Losses++};if l9{m.Cap9Losses++};m.NinthActions+=ninth
			if l8&&!l9{m.NinthRescues++}
			if !l8&&ninth>0{m.UnnecessaryNinth++}
			if l8&&ninth>0{m.NecessaryNinthAttempts++}
			if l8&&l9&&ninth>0{m.IneffectiveNinth++}
		}}
		if m.UnnecessaryNinth>0{m.RescuePerUnnecessary=float64(m.NinthRescues)/float64(m.UnnecessaryNinth)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
