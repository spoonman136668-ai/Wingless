package unitary

const UP149CHandSchema="wingless.up149c-hand-rotation-causality.v1"

type UP149CPoint struct{
	InitialHand int `json:"initial_hand"`
	TargetSlots []int `json:"target_slots"`
	RefreshBeforeAdmission int `json:"refresh_before_admission"`
	PresentBeforeRefresh int `json:"present_before_refresh"`
	SuccessfulRefreshQueries int `json:"successful_refresh_queries"`
	TargetFinalRetained int `json:"target_final_retained"`
	TargetFinalAccuracy float64 `json:"target_final_accuracy"`
	NonTargetDurableFinalRetained int `json:"non_target_durable_final_retained"`
	NewCandidatesFinalRetained int `json:"new_candidates_final_retained"`
	FinalHand int `json:"final_hand"`
	RecallEntriesUsed int `json:"recall_entries_used"`
}
type UP149CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP148CSeal string `json:"source_up148c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	InitialHands []int `json:"initial_hands"`
	RefreshArms int `json:"refresh_arms"`
	AdmissionsPerArm int `json:"admissions_per_arm"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	AdaptiveRefreshUsed bool `json:"adaptive_refresh_used"`
	CapacityIncreased bool `json:"capacity_increased"`
	Points []UP149CPoint `json:"points"`
}
func up149cRun(initialHand,refreshBefore int)UP149CPoint{
	x:=&up81cAging{};durable:=up145cDurable();cands:=up145cCandidates()
	for _,d:=range durable{x.write(d.key,d.class)}
	x.hand=initialHand
	slots:=[]int{x.find(0),x.find(1),x.find(2),x.find(3)}
	present,successful:=0,0
	for ai,c:=range cands{
		admission:=ai+1
		if admission==refreshBefore{
			for k:=0;k<4;k++{if x.find(k)>=0{present++};if _,ok:=x.query(k);ok{successful++}}
		}
		x.write(c.key,c.class)
	}
	target,other,newRet:=0,0,0
	for k:=0;k<16;k++{_,ok:=x.query(k);if k<4{if ok{target++}}else{if ok{other++}}}
	for _,c:=range cands{if _,ok:=x.query(c.key);ok{newRet++}}
	return UP149CPoint{InitialHand:initialHand,TargetSlots:slots,RefreshBeforeAdmission:refreshBefore,PresentBeforeRefresh:present,SuccessfulRefreshQueries:successful,TargetFinalRetained:target,TargetFinalAccuracy:float64(target)/4.0,NonTargetDurableFinalRetained:other,NewCandidatesFinalRetained:newRet,FinalHand:x.hand,RecallEntriesUsed:x.count}
}
func RunUP149C()(UP149CResult,error){
	hands:=[]int{0,4,8,12}
	res:=UP149CResult{Schema:UP149CHandSchema,Experiment:"UP-149C-hand-rotation-causality",SourceUP148CSeal:"84bfd1dde558aa5efdecb13f4e9f17e290c48c08",ExactRecallCap:16,InitialHands:append([]int(nil),hands...),RefreshArms:4,AdmissionsPerArm:4,MemoryPolicyChanged:false,AdaptiveRefreshUsed:false,CapacityIncreased:false}
	for _,h:=range hands{for _,refresh:=range []int{1,2,3,4}{res.Points=append(res.Points,up149cRun(h,refresh))}}
	return res,nil
}
