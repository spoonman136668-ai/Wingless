package unitary

import "sort"

const UP158CDurabilitySchema="wingless.up158c-save-durability-signal.v1"

type UP158CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	UniqueWriteHorizon int `json:"unique_write_horizon"`
	Outcome string `json:"outcome"`
	ActualLossStep int `json:"actual_loss_step"`
}
type UP158CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP157CSeal string `json:"source_up157c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	SemanticClassUsed bool `json:"semantic_class_used"`
	Arms int `json:"arms"`
	PermanentSaves int `json:"permanent_saves"`
	DelayedLosses int `json:"delayed_losses"`
	PermanentMeanHorizon float64 `json:"permanent_mean_horizon"`
	PermanentMinHorizon int `json:"permanent_min_horizon"`
	PermanentMaxHorizon int `json:"permanent_max_horizon"`
	DelayedMeanHorizon float64 `json:"delayed_mean_horizon"`
	DelayedMinHorizon int `json:"delayed_min_horizon"`
	DelayedMaxHorizon int `json:"delayed_max_horizon"`
	HorizonAUROC float64 `json:"horizon_auroc"`
	Points []UP158CPoint `json:"points"`
}
func up158cHorizon(x *up81cAging,key int)int{
	m:=*x
	for n:=1;n<=64;n++{
		m.write(900000+n,n%3)
		if m.find(key)<0{return n}
	}
	return 65
}
func up158cArm(hand int,cohort []int,trigger,endangered int)(int,bool,int){
	x:=up156cInit(hand);lossStep:=-1;horizon:=-1
	for step:=1;step<=24;step++{
		if rk,ok:=up155cRefresh(step);ok{x.query(rk)}
		key:=up156cStreamKey(step)
		if step==trigger{x.query(endangered);horizon=up158cHorizon(x,endangered)}
		x.write(key,step%3)
		if lossStep<0&&x.find(endangered)<0{lossStep=step}
	}
	return horizon,x.find(endangered)>=0,lossStep
}
func up158cMean(xs []int)float64{if len(xs)==0{return 0};s:=0;for _,x:=range xs{s+=x};return float64(s)/float64(len(xs))}
func up158cMinMax(xs []int)(int,int){if len(xs)==0{return 0,0};a:=append([]int(nil),xs...);sort.Ints(a);return a[0],a[len(a)-1]}
func up158cAUROC(pos,neg []int)float64{
	if len(pos)==0||len(neg)==0{return 0}
	score:=0.0
	for _,p:=range pos{for _,n:=range neg{if p>n{score+=1}else if p==n{score+=0.5}}}
	return score/float64(len(pos)*len(neg))
}
func RunUP158C()(UP158CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"}
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};perm:=[]int{};delay:=[]int{}
	res:=UP158CResult{Schema:UP158CDurabilitySchema,Experiment:"UP-158C-save-durability-signal",SourceUP157CSeal:"e528d4509aeb4de120b343328e4f90b75fb09ca5",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,SemanticClassUsed:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up156cFindTrigger(hand,c);if trigger<0{continue}
		h,final,loss:=up158cArm(hand,c,trigger,endangered);outcome:="delayed_loss"
		if final{outcome="permanent_save";res.PermanentSaves++;perm=append(perm,h)}else{res.DelayedLosses++;delay=append(delay,h)}
		res.Arms++
		res.Points=append(res.Points,UP158CPoint{Cohort:names[ci],InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,UniqueWriteHorizon:h,Outcome:outcome,ActualLossStep:loss})
	}}
	res.PermanentMeanHorizon=up158cMean(perm);res.DelayedMeanHorizon=up158cMean(delay)
	res.PermanentMinHorizon,res.PermanentMaxHorizon=up158cMinMax(perm)
	res.DelayedMinHorizon,res.DelayedMaxHorizon=up158cMinMax(delay)
	res.HorizonAUROC=up158cAUROC(perm,delay)
	return res,nil
}
