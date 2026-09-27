package unitary

const UP164CRobustWarningSchema="wingless.up164c-robust-horizon-warning.v1"

type UP164CMetric struct{
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	ExpiredWarningIntervals int `json:"expired_warning_intervals"`
	WarningRecall float64 `json:"warning_recall"`
}
type UP164CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP163CSeal string `json:"source_up163c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Cadences []int `json:"cadences"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	ExactFutureScheduleUsed bool `json:"exact_future_schedule_used"`
	CurrentStateRobustHorizonUsed bool `json:"current_state_robust_horizon_used"`
	Metrics []UP164CMetric `json:"metrics"`
}
func up164cRealStep(m *up81cAging,endangered,key int)bool{
	slot:=up151cPredict(m.hand,m.age)
	if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
	m.write(key,key%3)
	return m.find(endangered)<0
}
func up164cRun(x0 *up81cAging,endangered,cadence int)(warned,missed,expired int){
	m:=*x0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		lost:=false
		for j:=0;j<cadence && start+j<=64;j++{
			if up164cRealStep(&m,endangered,940000+start+j){lost=true;break}
		}
		if lost{
			if warn{warned++}else{missed++}
			return
		}
		if warn{expired++}
	}
	return
}
func RunUP164C()(UP164CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};cadences:=[]int{2,4}
	res:=UP164CResult{Schema:UP164CRobustWarningSchema,Experiment:"UP-164C-robust-horizon-warning",SourceUP163CSeal:"382c766c147ea3a2e91e6874b7e296d8b52b5261",ExactRecallCap:16,Cadences:cadences,ShadowOnly:true,CorrectiveActionUsed:false,ExactFutureScheduleUsed:false,CurrentStateRobustHorizonUsed:true}
	for _,cad:=range cadences{
		m:=UP164CMetric{Cadence:cad}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++;m.EventualLosses++
			w,miss,exp:=up164cRun(x,endangered,cad);m.WarnedLosses+=w;m.MissedLosses+=miss;m.ExpiredWarningIntervals+=exp
		}}
		if m.EventualLosses>0{m.WarningRecall=float64(m.WarnedLosses)/float64(m.EventualLosses)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
