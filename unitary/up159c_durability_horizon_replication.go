package unitary

const UP159CReplicationSchema="wingless.up159c-durability-horizon-replication.v1"

type UP159CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	UniqueWriteHorizon int `json:"unique_write_horizon"`
	PredictedDurable bool `json:"predicted_durable"`
	ActualDurable bool `json:"actual_durable"`
	ActualLossStep int `json:"actual_loss_step"`
}
type UP159CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP158CSeal string `json:"source_up158c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	DurabilityThreshold int `json:"durability_threshold"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	SemanticClassUsed bool `json:"semantic_class_used"`
	ThresholdAdaptive bool `json:"threshold_adaptive"`
	Arms int `json:"arms"`
	PermanentSaves int `json:"permanent_saves"`
	DelayedLosses int `json:"delayed_losses"`
	TruePositive int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	FalseNegative int `json:"false_negative"`
	TrueNegative int `json:"true_negative"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Accuracy float64 `json:"accuracy"`
	HorizonAUROC float64 `json:"horizon_auroc"`
	Points []UP159CPoint `json:"points"`
}
func up159cStreamKey(step int)int{
	novel,last:=0,-1;key:=-1
	for s:=1;s<=step;s++{
		if s%4==0&&last>=0{key=last}else{novel++;last=8000+novel;key=last}
	}
	return key
}
func up159cRefresh(step int)(int,bool){
	switch step{case 2:return 5,true;case 6:return 10,true;case 11:return 15,true;case 15:return 4,true;case 20:return 9,true}
	return -1,false
}
func up159cFindTrigger(hand int,cohort []int)(int,int){
	x:=up156cInit(hand)
	for step:=1;step<=24;step++{
		if rk,ok:=up159cRefresh(step);ok{x.query(rk)}
		key:=up159cStreamKey(step)
		if x.find(key)<0{
			slot:=up151cPredict(x.hand,x.age)
			if x.entries[slot].used&&up155cHas(cohort,x.entries[slot].key){return step,x.entries[slot].key}
		}
		x.write(key,step%3)
	}
	return -1,-1
}
func up159cArm(hand int,trigger,endangered int)(int,bool,int){
	x:=up156cInit(hand);horizon:=-1;loss:=-1
	for step:=1;step<=24;step++{
		if rk,ok:=up159cRefresh(step);ok{x.query(rk)}
		key:=up159cStreamKey(step)
		if step==trigger{x.query(endangered);horizon=up158cHorizon(x,endangered)}
		x.write(key,step%3)
		if loss<0&&x.find(endangered)<0{loss=step}
	}
	return horizon,x.find(endangered)>=0,loss
}
func RunUP159C()(UP159CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"}
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};threshold:=21;perm:=[]int{};delay:=[]int{}
	res:=UP159CResult{Schema:UP159CReplicationSchema,Experiment:"UP-159C-durability-horizon-replication",SourceUP158CSeal:"60d470b7d11aae27c98a3d73d62bea2984b8e5e6",ExactRecallCap:16,DurabilityThreshold:threshold,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,SemanticClassUsed:false,ThresholdAdaptive:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up159cFindTrigger(hand,c);if trigger<0{continue}
		h,durable,loss:=up159cArm(hand,trigger,endangered);pred:=h>=threshold
		res.Arms++;if durable{res.PermanentSaves++;perm=append(perm,h)}else{res.DelayedLosses++;delay=append(delay,h)}
		if pred&&durable{res.TruePositive++}else if pred&&!durable{res.FalsePositive++}else if !pred&&durable{res.FalseNegative++}else{res.TrueNegative++}
		res.Points=append(res.Points,UP159CPoint{Cohort:names[ci],InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,UniqueWriteHorizon:h,PredictedDurable:pred,ActualDurable:durable,ActualLossStep:loss})
	}}
	res.Precision=up158cRate(res.TruePositive,res.TruePositive+res.FalsePositive)
	res.Recall=up158cRate(res.TruePositive,res.TruePositive+res.FalseNegative)
	res.Accuracy=up158cRate(res.TruePositive+res.TrueNegative,res.Arms)
	res.HorizonAUROC=up158cAUROC(perm,delay)
	return res,nil
}
