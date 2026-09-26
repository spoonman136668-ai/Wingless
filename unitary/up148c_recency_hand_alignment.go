package unitary

const UP148CAlignmentSchema="wingless.up148c-recency-hand-alignment.v1"

type UP148CPoint struct{
	Rotation int `json:"rotation"`
	TargetInitialSlots []int `json:"target_initial_slots"`
	InitialHand int `json:"initial_hand"`
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
type UP148CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP147CSeal string `json:"source_up147c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	TargetFacts int `json:"target_facts"`
	Rotations []int `json:"rotations"`
	RefreshArms int `json:"refresh_arms"`
	AdmissionsPerArm int `json:"admissions_per_arm"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	AdaptiveRefreshUsed bool `json:"adaptive_refresh_used"`
	CapacityIncreased bool `json:"capacity_increased"`
	Points []UP148CPoint `json:"points"`
}
func up148cDurable(rotation int)[]up145cLexeme{
	d:=up145cDurable();out:=make([]up145cLexeme,0,len(d))
	for i:=0;i<len(d);i++{out=append(out,d[(rotation+i)%len(d)])}
	return out
}
func up148cRun(rotation,refreshBefore int)UP148CPoint{
	x:=&up81cAging{};durable:=up148cDurable(rotation);cands:=up145cCandidates()
	for _,d:=range durable{x.write(d.key,d.class)}
	slots:=make([]int,0,4);for k:=0;k<4;k++{slots=append(slots,x.find(k))}
	initialHand:=x.hand;presentBefore,successful:=0,0
	for ai,c:=range cands{
		admission:=ai+1
		if admission==refreshBefore{
			for k:=0;k<4;k++{
				if x.find(k)>=0{presentBefore++}
				if _,ok:=x.query(k);ok{successful++}
			}
		}
		x.write(c.key,c.class)
	}
	targetRet,nonTargetRet,newRet:=0,0,0
	for k:=0;k<16;k++{
		_,ok:=x.query(k)
		if k<4{if ok{targetRet++}}else{if ok{nonTargetRet++}}
	}
	for _,c:=range cands{if _,ok:=x.query(c.key);ok{newRet++}}
	return UP148CPoint{Rotation:rotation,TargetInitialSlots:slots,InitialHand:initialHand,RefreshBeforeAdmission:refreshBefore,SubsequentAdmissions:5-refreshBefore,PresentBeforeRefresh:presentBefore,SuccessfulRefreshQueries:successful,TargetFinalRetained:targetRet,TargetFinalAccuracy:float64(targetRet)/4.0,NonTargetDurableFinalRetained:nonTargetRet,NewCandidatesFinalRetained:newRet,RecallEntriesUsed:x.count,FinalHand:x.hand}
}
func RunUP148C()(UP148CResult,error){
	rots:=[]int{0,4,8,12}
	res:=UP148CResult{Schema:UP148CAlignmentSchema,Experiment:"UP-148C-recency-hand-alignment",SourceUP147CSeal:"ac6f4f4e0aff80939ced6251a2b6461ad67ff2c7",ExactRecallCap:16,TargetFacts:4,Rotations:append([]int(nil),rots...),RefreshArms:4,AdmissionsPerArm:4,MemoryPolicyChanged:false,AdaptiveRefreshUsed:false,CapacityIncreased:false}
	for _,rot:=range rots{for _,refresh:=range []int{1,2,3,4}{res.Points=append(res.Points,up148cRun(rot,refresh))}}
	return res,nil
}
