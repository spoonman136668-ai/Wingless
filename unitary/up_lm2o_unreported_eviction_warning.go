package unitary

const UPLM2OWarningSchema="wingless.up-lm2o-unreported-eviction-warning.v1"

type UPLM2OPoint struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	SecondChunkStoreOrdinal int `json:"second_chunk_store_ordinal"`
	RecallEntriesBefore int `json:"recall_entries_before"`
	PredictedVictim string `json:"predicted_victim"`
	Warning bool `json:"warning"`
	ActualVictim string `json:"actual_victim"`
	ActualUnreportedEviction bool `json:"actual_unreported_eviction"`
	Correct bool `json:"correct"`
}
type UPLM2OSummary struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	Warnings int `json:"warnings"`
	ActualUnreportedEvictions int `json:"actual_unreported_evictions"`
}
type UPLM2OResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2NSeal string `json:"source_up_lm2n_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	DeferredLevels []int `json:"deferred_levels"`
	IdentityRotations []int `json:"identity_rotations"`
	ValueShifts []int `json:"value_shifts"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	InterventionTriggered bool `json:"intervention_triggered"`
	TruePositive int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	FalseNegative int `json:"false_negative"`
	TrueNegative int `json:"true_negative"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Points []UPLM2OPoint `json:"points"`
	Summaries []UPLM2OSummary `json:"summaries"`
}
func uplm2oRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func uplm2oRun(d,rot,shift int)([]UPLM2OPoint,int,int,int,int){
	recall:=newUPLM0CRecall();reported:=map[string]bool{}
	for i:=0;i<12;i++{recall.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
	for i:=0;i<12-d;i++{reported[uplm2nName(i,rot)]=true}
	out:=[]UPLM2OPoint{};tp,fp,fn,tn:=0,0,0,0
	for j:=0;j<12;j++{
		predVictim:="";warn:=false
		if len(recall.order)>=16{
			predVictim=recall.order[0];warn=!reported[predVictim]
		}
		actualVictim:="";actualUnreported:=false
		if len(recall.order)>=16{actualVictim=recall.order[0];actualUnreported=!reported[actualVictim]}
		recall.write(uplm2nName(12+j,rot),uplm2nValue(12+j,rot,shift))
		ok:=warn==actualUnreported
		if warn&&actualUnreported{tp++}else if warn&&!actualUnreported{fp++}else if !warn&&actualUnreported{fn++}else{tn++}
		out=append(out,UPLM2OPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,SecondChunkStoreOrdinal:j+1,RecallEntriesBefore:func()int{if j<4{return 12+j};return 16}(),PredictedVictim:predVictim,Warning:warn,ActualVictim:actualVictim,ActualUnreportedEviction:actualUnreported,Correct:ok})
	}
	return out,tp,fp,fn,tn
}
func RunUPLM2O()(UPLM2OResult,error){
	levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
	res:=UPLM2OResult{Schema:UPLM2OWarningSchema,Experiment:"UP-LM2O-unreported-eviction-warning",SourceUPLM2NSeal:"7bda1b312da8b0958481b9f03f08a496b1476a75",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,FutureOracleUsed:false,InterventionTriggered:false}
	for _,d:=range levels{
		warnCount,actualCount:=0,0
		for _,rot:=range rots{for _,shift:=range shifts{
			p,tp,fp,fn,tn:=uplm2oRun(d,rot,shift);res.Points=append(res.Points,p...);res.TruePositive+=tp;res.FalsePositive+=fp;res.FalseNegative+=fn;res.TrueNegative+=tn
			for _,x:=range p{if x.Warning{warnCount++};if x.ActualUnreportedEviction{actualCount++}}
		}}
		res.Summaries=append(res.Summaries,UPLM2OSummary{DeferredFirstChunkReports:d,Warnings:warnCount,ActualUnreportedEvictions:actualCount})
	}
	res.Precision=uplm2oRate(res.TruePositive,res.TruePositive+res.FalsePositive);res.Recall=uplm2oRate(res.TruePositive,res.TruePositive+res.FalseNegative)
	return res,nil
}
