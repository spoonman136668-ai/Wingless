package unitary

const UP147CRecencySchema="wingless.up147c-durable-recency-gradient.v1"

type UP147CPoint struct{
	Cohort string `json:"cohort"`
	CohortStart int `json:"cohort_start"`
	RefreshBeforeAdmission int `json:"refresh_before_admission"`
	SubsequentAdmissions int `json:"subsequent_admissions"`
	PresentBeforeRefresh int `json:"present_before_refresh"`
	SuccessfulRefreshQueries int `json:"successful_refresh_queries"`
	TargetFinalRetained int `json:"target_final_retained"`
	TargetFinalAccuracy float64 `json:"target_final_accuracy"`
	NonTargetDurableFinalRetained int `json:"non_target_durable_final_retained"`
	NewCandidatesFinalRetained int `json:"new_candidates_final_retained"`
	RecallEntriesUsed int `json:"recall_entries_used"`
	FinalHand int `json:"final_hand"`
}
type UP147CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP146CSeal string `json:"source_up146c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Cohorts int `json:"cohorts"`
	RefreshArms int `json:"refresh_arms"`
	AdmissionsPerArm int `json:"admissions_per_arm"`
	RefreshAttemptsPerArm int `json:"refresh_attempts_per_arm"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	AdaptiveQuerySelection bool `json:"adaptive_query_selection"`
	CapacityIncreased bool `json:"capacity_increased"`
	Points []UP147CPoint `json:"points"`
}
func up147cCohortName(start int)string{
	switch start{case 0:return "keys_0_3";case 4:return "keys_4_7";case 8:return "keys_8_11";default:return "keys_12_15"}
}
func up147cRun(start,refreshBefore int)UP147CPoint{
	x:=&up81cAging{};durable:=up145cDurable();cands:=up145cCandidates();truth:=map[int]int{}
	for _,d:=range durable{truth[d.key]=d.class;x.write(d.key,d.class)}
	presentBefore,successful:=0,0
	for ai,c:=range cands{
		admission:=ai+1
		if admission==refreshBefore{
			for k:=start;k<start+4;k++{
				if x.find(k)>=0{presentBefore++}
				if _,ok:=x.query(k);ok{successful++}
			}
		}
		x.write(c.key,c.class)
	}
	targetRet,nonTargetRet,newRet:=0,0,0
	for _,d:=range durable{
		_,ok:=x.query(d.key)
		if d.key>=start&&d.key<start+4{if ok{targetRet++}}else{if ok{nonTargetRet++}}
	}
	for _,c:=range cands{if _,ok:=x.query(c.key);ok{newRet++}}
	return UP147CPoint{Cohort:up147cCohortName(start),CohortStart:start,RefreshBeforeAdmission:refreshBefore,SubsequentAdmissions:5-refreshBefore,PresentBeforeRefresh:presentBefore,SuccessfulRefreshQueries:successful,TargetFinalRetained:targetRet,TargetFinalAccuracy:float64(targetRet)/4.0,NonTargetDurableFinalRetained:nonTargetRet,NewCandidatesFinalRetained:newRet,RecallEntriesUsed:x.count,FinalHand:x.hand}
}
func RunUP147C()(UP147CResult,error){
	res:=UP147CResult{Schema:UP147CRecencySchema,Experiment:"UP-147C-durable-recency-gradient",SourceUP146CSeal:"242980067208103a4757eb194b5585c99cf73e04",ExactRecallCap:16,Cohorts:4,RefreshArms:4,AdmissionsPerArm:4,RefreshAttemptsPerArm:4,MemoryPolicyChanged:false,AdaptiveQuerySelection:false,CapacityIncreased:false}
	for _,start:=range []int{0,4,8,12}{for _,refresh:=range []int{1,2,3,4}{res.Points=append(res.Points,up147cRun(start,refresh))}}
	return res,nil
}
