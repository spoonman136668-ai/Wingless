package unitary

import("fmt";"strings")

const UP160CGeometrySchema="wingless.up160c-durability-geometry.v1"

type UP160CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	UniqueWriteHorizon int `json:"unique_write_horizon"`
	CanonicalFingerprint string `json:"canonical_fingerprint"`
	ActualDurable bool `json:"actual_durable"`
}
type UP160CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP159CSeal string `json:"source_up159c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	SemanticIdentityUsedInFingerprint bool `json:"semantic_identity_used_in_fingerprint"`
	Horizon16Cases int `json:"horizon16_cases"`
	Horizon16UniqueFingerprints int `json:"horizon16_unique_fingerprints"`
	Horizon16MixedOutcomeFingerprints int `json:"horizon16_mixed_outcome_fingerprints"`
	Points []UP160CPoint `json:"points"`
}
func up160cFingerprint(x *up81cAging,endangered int)string{
	slot:=x.find(endangered);offset:=-1;if slot>=0{offset=(slot-x.hand+16)%16}
	parts:=make([]string,0,16)
	for i:=0;i<16;i++{parts=append(parts,fmt.Sprintf("%d",x.age[(x.hand+i)%16]))}
	return fmt.Sprintf("offset=%d|ages=%s",offset,strings.Join(parts,","))
}
func up160cArm(hand,trigger,endangered int)(int,string,bool){
	x:=up156cInit(hand);horizon:=-1;fp:=""
	for step:=1;step<=24;step++{
		if rk,ok:=up159cRefresh(step);ok{x.query(rk)}
		key:=up159cStreamKey(step)
		if step==trigger{x.query(endangered);horizon=up158cHorizon(x,endangered);fp=up160cFingerprint(x,endangered)}
		x.write(key,step%3)
	}
	return horizon,fp,x.find(endangered)>=0
}
func RunUP160C()(UP160CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"}
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};outcomes:=map[string]int{}
	res:=UP160CResult{Schema:UP160CGeometrySchema,Experiment:"UP-160C-durability-geometry",SourceUP159CSeal:"3f621ebede6fc7b4aa09bcb0eeb5b06bde295943",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,SemanticIdentityUsedInFingerprint:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up159cFindTrigger(hand,c);if trigger<0{continue}
		h,fp,durable:=up160cArm(hand,trigger,endangered)
		res.Points=append(res.Points,UP160CPoint{Cohort:names[ci],InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,UniqueWriteHorizon:h,CanonicalFingerprint:fp,ActualDurable:durable})
		if h==16{res.Horizon16Cases++;if durable{outcomes[fp]|=1}else{outcomes[fp]|=2}}
	}}
	res.Horizon16UniqueFingerprints=len(outcomes)
	for _,mask:=range outcomes{if mask==3{res.Horizon16MixedOutcomeFingerprints++}}
	return res,nil
}
