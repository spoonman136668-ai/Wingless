package unitary

const UP157CRedistributionSchema="wingless.up157c-loss-redistribution.v1"

type UP157CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	GuidedOutcome string `json:"guided_outcome"`
	GuidedEndangeredFinal bool `json:"guided_endangered_final"`
	GuidedEndangeredLossStep int `json:"guided_endangered_loss_step"`
	BaselineFinalProtected int `json:"baseline_final_protected"`
	GuidedFinalProtected int `json:"guided_final_protected"`
	ShamFinalProtected int `json:"sham_final_protected"`
	BaselineFinalOriginal int `json:"baseline_final_original"`
	GuidedFinalOriginal int `json:"guided_final_original"`
	ShamFinalOriginal int `json:"sham_final_original"`
	BaselineFinalNonCohortOriginal int `json:"baseline_final_noncohort_original"`
	GuidedFinalNonCohortOriginal int `json:"guided_final_noncohort_original"`
	ShamFinalNonCohortOriginal int `json:"sham_final_noncohort_original"`
}
type UP157CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP156CSeal string `json:"source_up156c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Arms int `json:"arms"`
	PermanentSaves int `json:"permanent_saves"`
	DelayedLosses int `json:"delayed_losses"`
	GuidedFinalProtectedGain int `json:"guided_final_protected_gain"`
	GuidedFinalOriginalGain int `json:"guided_final_original_gain"`
	GuidedFinalNonCohortOriginalDelta int `json:"guided_final_noncohort_original_delta"`
	Points []UP157CPoint `json:"points"`
}
func up157cOriginalCount(x *up81cAging)int{n:=0;for k:=0;k<16;k++{if x.find(k)>=0{n++}};return n}
func up157cRun(hand int,cohort []int,trigger,endangered int,mode string)(protected,original,noncohort int,endangeredFinal bool,lossStep int){
	x:=up156cInit(hand);lossStep=-1;wasPresent:=x.find(endangered)>=0
	for step:=1;step<=24;step++{
		if rk,ok:=up155cRefresh(step);ok{x.query(rk)}
		key:=up156cStreamKey(step)
		if step==trigger{
			if mode=="guided"{x.query(endangered)}
			if mode=="sham"{t:=up156cShamTarget(x,cohort,endangered);if t>=0{x.query(t)}}
		}
		x.write(key,step%3)
		now:=x.find(endangered)>=0
		if wasPresent&&!now&&lossStep<0{lossStep=step}
		wasPresent=now
	}
	protected=up155cPresent(x,cohort);original=up157cOriginalCount(x);noncohort=original-protected;endangeredFinal=x.find(endangered)>=0
	return
}
func RunUP157C()(UP157CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"};cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}};hands:=[]int{0,4,8,12}
	res:=UP157CResult{Schema:UP157CRedistributionSchema,Experiment:"UP-157C-loss-redistribution",SourceUP156CSeal:"527fe0829e1224b02d47e0827d07368ead8ce730",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up156cFindTrigger(hand,c);if trigger<0{continue}
		bp,bo,bn,_,_:=up157cRun(hand,c,trigger,endangered,"baseline")
		gp,go0,gn,gef,gl:=up157cRun(hand,c,trigger,endangered,"guided")
		sp,so,sn,_,_:=up157cRun(hand,c,trigger,endangered,"sham")
		outcome:="delayed_loss";if gef{outcome="permanent_save";res.PermanentSaves++}else{res.DelayedLosses++}
		res.Arms++;res.GuidedFinalProtectedGain+=gp-bp;res.GuidedFinalOriginalGain+=go0-bo;res.GuidedFinalNonCohortOriginalDelta+=gn-bn
		res.Points=append(res.Points,UP157CPoint{Cohort:names[ci],InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,GuidedOutcome:outcome,GuidedEndangeredFinal:gef,GuidedEndangeredLossStep:gl,BaselineFinalProtected:bp,GuidedFinalProtected:gp,ShamFinalProtected:sp,BaselineFinalOriginal:bo,GuidedFinalOriginal:go0,ShamFinalOriginal:so,BaselineFinalNonCohortOriginal:bn,GuidedFinalNonCohortOriginal:gn,ShamFinalNonCohortOriginal:sn})
	}}
	return res,nil
}
