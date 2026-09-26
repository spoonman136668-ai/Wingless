package unitary

const UPLM2NLeadSchema="wingless.up-lm2n-recall-loss-lead-time.v1"

type UPLM2NPoint struct{
	IdentityRotation int `json:"identity_rotation"`
	DeferredReportOrder string `json:"deferred_report_order"`
	ValueShift int `json:"value_shift"`
	LostPendingPosition int `json:"lost_pending_position"`
	LostPendingName string `json:"lost_pending_name"`
	LossEventOrdinal int `json:"loss_event_ordinal"`
	ReportEventOrdinal int `json:"report_event_ordinal"`
	LeadEvents int `json:"lead_events"`
	RecallEntriesAtLoss int `json:"recall_entries_at_loss"`
}
type UPLM2NResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2MSeal string `json:"source_up_lm2m_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	EntityCount int `json:"entity_count"`
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotations []int `json:"identity_rotations"`
	DeferredReportOrders []string `json:"deferred_report_orders"`
	ValueShifts []int `json:"value_shifts"`
	ModelTrainingUsed bool `json:"model_training_used"`
	InterventionTriggered bool `json:"intervention_triggered"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	Points []UPLM2NPoint `json:"points"`
}
func uplm2nName(pos,rot int)string{names:=uplm2iNames();return names[(pos+rot)%24]}
func uplm2nValue(pos,rot,shift int)string{v:=uplm0gValues();return v[((pos+rot)+shift)%len(v)]}
func uplm2nRun(rot int,order string,shift int)UPLM2NPoint{
	recall:=newUPLM0CRecall();event:=0
	pending:=map[string]int{};present:=map[string]bool{};lossEvent:=map[string]int{}
	for i:=0;i<12;i++{event++;n:=uplm2nName(i,rot);recall.write(n,uplm2nValue(i,rot,shift));if i>=7{pending[n]=i;present[n]=true}}
	for i:=0;i<12;i++{event++} // first-chunk OBSERVE
	for i:=0;i<7;i++{event++}  // early first-chunk REPORT
	for i:=12;i<24;i++{
		event++;recall.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))
		for n:=range pending{
			_,ok:=recall.values[n]
			if present[n]&&!ok&&lossEvent[n]==0{lossEvent[n]=event}
			present[n]=ok
		}
	}
	for i:=12;i<24;i++{event++} // second-chunk OBSERVE
	orderPos:=[]int{7,8,9,10,11};if order=="reverse"{orderPos=[]int{11,10,9,8,7}}
	for _,pos:=range orderPos{
		event++;n:=uplm2nName(pos,rot)
		if le:=lossEvent[n];le>0{return UPLM2NPoint{IdentityRotation:rot,DeferredReportOrder:order,ValueShift:shift,LostPendingPosition:pos,LostPendingName:n,LossEventOrdinal:le,ReportEventOrdinal:event,LeadEvents:event-le,RecallEntriesAtLoss:16}}
	}
	return UPLM2NPoint{IdentityRotation:rot,DeferredReportOrder:order,ValueShift:shift,LostPendingPosition:-1,LossEventOrdinal:-1,ReportEventOrdinal:-1,LeadEvents:-1,RecallEntriesAtLoss:len(recall.order)}
}
func RunUPLM2N()(UPLM2NResult,error){
	rots:=[]int{0,7};orders:=[]string{"forward","reverse"};shifts:=[]int{0,1,2,3}
	res:=UPLM2NResult{Schema:UPLM2NLeadSchema,Experiment:"UP-LM2N-recall-loss-lead-time",SourceUPLM2MSeal:"e4482856eb3f80456c88997ac6a80496d99986a0",ExactRecallCap:16,EntityCount:24,DeferredFirstChunkReports:5,IdentityRotations:rots,DeferredReportOrders:orders,ValueShifts:shifts,ModelTrainingUsed:false,InterventionTriggered:false,FutureOracleUsed:false}
	for _,rot:=range rots{for _,ord:=range orders{for _,shift:=range shifts{res.Points=append(res.Points,uplm2nRun(rot,ord,shift))}}}
	return res,nil
}
