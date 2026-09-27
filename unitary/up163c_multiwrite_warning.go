package unitary

const UP163CMultiwriteSchema="wingless.up163c-multiwrite-warning.v1"

type UP163CMetric struct{
	Cadence int `json:"cadence"`
	LookaheadWrites int `json:"lookahead_writes"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	FalseWarnings int `json:"false_warnings"`
	WarningRecall float64 `json:"warning_recall"`
}
type UP163CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP162CSeal string `json:"source_up162c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	ExternalFutureScheduleUsed bool `json:"external_future_schedule_used"`
	KnownPressurePolicySimulated bool `json:"known_pressure_policy_simulated"`
	Metrics []UP163CMetric `json:"metrics"`
}
func up163cStep(m *up81cAging,endangered,key int)bool{
	slot:=up151cPredict(m.hand,m.age)
	if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
	m.write(key,key%3)
	return m.find(endangered)<0
}
func up163cWouldLoseWithin(x *up81cAging,endangered,startKey,horizon int)bool{
	m:=*x
	for i:=0;i<horizon;i++{if up163cStep(&m,endangered,startKey+i){return true}}
	return false
}
func up163cRun(x0 *up81cAging,endangered,cadence int)(warned,missed,falseWarn int){
	m:=*x0;warningActive:=false
	for n:=1;n<=64;n++{
		if up162cScheduled(n,cadence){
			warningActive=up163cWouldLoseWithin(&m,endangered,930000+n,cadence)
		}
		lost:=up163cStep(&m,endangered,930000+n)
		if lost{
			if warningActive{warned++}else{missed++}
			return
		}
		if warningActive&&up162cScheduled(n+1,cadence){falseWarn++}
	}
	return
}
func RunUP163C()(UP163CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};cadences:=[]int{2,4}
	res:=UP163CResult{Schema:UP163CMultiwriteSchema,Experiment:"UP-163C-multiwrite-warning",SourceUP162CSeal:"181b0bb3dace65aaad51d3926627c70142439bd1",ExactRecallCap:16,ShadowOnly:true,CorrectiveActionUsed:false,ExternalFutureScheduleUsed:false,KnownPressurePolicySimulated:true}
	for _,cad:=range cadences{
		m:=UP163CMetric{Cadence:cad,LookaheadWrites:cad}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++;m.EventualLosses++
			w,miss,fp:=up163cRun(x,endangered,cad);m.WarnedLosses+=w;m.MissedLosses+=miss;m.FalseWarnings+=fp
		}}
		if m.EventualLosses>0{m.WarningRecall=float64(m.WarnedLosses)/float64(m.EventualLosses)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
