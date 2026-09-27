package unitary

const UP162CRecheckSchema="wingless.up162c-periodic-recheck.v1"

type UP162CMetric struct{
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	FalseWarnings int `json:"false_warnings"`
	WarningRecall float64 `json:"warning_recall"`
}
type UP162CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP161CSeal string `json:"source_up161c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Cadences []int `json:"cadences"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	AdaptiveCadenceUsed bool `json:"adaptive_cadence_used"`
	Metrics []UP162CMetric `json:"metrics"`
}
func up162cScheduled(step,cadence int)bool{return step==1||(step-1)%cadence==0}
func up162cRun(x0 *up81cAging,endangered,cadence int)(warned,missed,falseWarn int){
	m:=*x0;warnedThisLoss:=false
	for n:=1;n<=64;n++{
		slot:=up151cPredict(m.hand,m.age)
		if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
		slot=up151cPredict(m.hand,m.age)
		harm:=m.entries[slot].used&&m.entries[slot].key==endangered
		if up162cScheduled(n,cadence){
			warn:=harm
			if warn&&!harm{falseWarn++}
			if warn&&harm{warnedThisLoss=true}
		}
		m.write(920000+n,n%3)
		if m.find(endangered)<0{
			if warnedThisLoss{warned++}else{missed++}
			return
		}
		warnedThisLoss=false
	}
	return
}
func RunUP162C()(UP162CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};cadences:=[]int{1,2,4}
	res:=UP162CResult{Schema:UP162CRecheckSchema,Experiment:"UP-162C-periodic-recheck",SourceUP161CSeal:"b3ae72f2192ef21892278bbba86afe45b6439f04",ExactRecallCap:16,Cadences:cadences,ShadowOnly:true,CorrectiveActionUsed:false,FutureScheduleUsed:false,AdaptiveCadenceUsed:false}
	for _,cad:=range cadences{
		m:=UP162CMetric{Cadence:cad}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++;m.EventualLosses++
			w,miss,fp:=up162cRun(x,endangered,cad);m.WarnedLosses+=w;m.MissedLosses+=miss;m.FalseWarnings+=fp
		}}
		if m.EventualLosses>0{m.WarningRecall=float64(m.WarnedLosses)/float64(m.EventualLosses)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
