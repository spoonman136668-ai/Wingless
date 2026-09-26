package unitary

const UP156CRefreshCFSchema="wingless.up156c-bounded-refresh-counterfactual.v1"

type UP156CPoint struct{
	Cohort string `json:"cohort"`
	Keys []int `json:"keys"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	BaselineImmediateLoss bool `json:"baseline_immediate_loss"`
	GuidedImmediateLoss bool `json:"guided_immediate_loss"`
	ShamImmediateLoss bool `json:"sham_immediate_loss"`
	GuidedActionTarget int `json:"guided_action_target"`
	ShamActionTarget int `json:"sham_action_target"`
	BaselineFinalProtected int `json:"baseline_final_protected"`
	GuidedFinalProtected int `json:"guided_final_protected"`
	ShamFinalProtected int `json:"sham_final_protected"`
}
type UP156CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP155CSeal string `json:"source_up155c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	Arms int `json:"arms"`
	BaselineImmediateLosses int `json:"baseline_immediate_losses"`
	GuidedImmediateLosses int `json:"guided_immediate_losses"`
	ShamImmediateLosses int `json:"sham_immediate_losses"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
	GuidedFinalRetentionGain int `json:"guided_final_retention_gain"`
	ShamFinalRetentionGain int `json:"sham_final_retention_gain"`
	Points []UP156CPoint `json:"points"`
}
func up156cStreamKey(step int)int{
	novel,last:=0,-1;key:=-1
	for s:=1;s<=step;s++{
		if s%3==0&&last>=0{key=last}else{novel++;last=7000+novel;key=last}
	}
	return key
}
func up156cInit(hand int)*up81cAging{x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=hand;return x}
func up156cFindTrigger(hand int,cohort []int)(int,int){
	x:=up156cInit(hand)
	for step:=1;step<=24;step++{
		if rk,ok:=up155cRefresh(step);ok{x.query(rk)}
		key:=up156cStreamKey(step)
		if x.find(key)<0{
			slot:=up151cPredict(x.hand,x.age)
			if x.entries[slot].used&&up155cHas(cohort,x.entries[slot].key){return step,x.entries[slot].key}
		}
		x.write(key,step%3)
	}
	return -1,-1
}
func up156cShamTarget(x *up81cAging,cohort []int,endangered int)int{
	best:=-1
	for i:=0;i<16;i++{if x.entries[i].used{k:=x.entries[i].key;if k!=endangered&&!up155cHas(cohort,k)&&(best<0||k<best){best=k}}}
	return best
}
func up156cRunMode(hand int,cohort []int,trigger,endangered int,mode string)(bool,int,int){
	x:=up156cInit(hand);immediate:=false;actionTarget:=-1
	for step:=1;step<=24;step++{
		if rk,ok:=up155cRefresh(step);ok{x.query(rk)}
		key:=up156cStreamKey(step)
		before:=up155cPresent(x,cohort)
		if step==trigger{
			if mode=="guided"{actionTarget=endangered;x.query(endangered)}
			if mode=="sham"{actionTarget=up156cShamTarget(x,cohort,endangered);if actionTarget>=0{x.query(actionTarget)}}
		}
		x.write(key,step%3)
		after:=up155cPresent(x,cohort)
		if step==trigger{immediate=after<before}
	}
	return immediate,up155cPresent(x,cohort),actionTarget
}
func RunUP156C()(UP156CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"};cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}};hands:=[]int{0,4,8,12}
	res:=UP156CResult{Schema:UP156CRefreshCFSchema,Experiment:"UP-156C-bounded-refresh-counterfactual",SourceUP155CSeal:"4ba5db866befcbe662167e0d32cad10a2b317cf5",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up156cFindTrigger(hand,c);if trigger<0{continue}
		bl,bf,_:=up156cRunMode(hand,c,trigger,endangered,"baseline")
		gl,gf,gt:=up156cRunMode(hand,c,trigger,endangered,"guided")
		sl,sf,st:=up156cRunMode(hand,c,trigger,endangered,"sham")
		res.Arms++;if bl{res.BaselineImmediateLosses++};if gl{res.GuidedImmediateLosses++};if sl{res.ShamImmediateLosses++}
		if bl&&!gl{res.GuidedPrevented++};if bl&&!sl{res.ShamPrevented++};res.GuidedActions++;res.ShamActions++;res.GuidedFinalRetentionGain+=gf-bf;res.ShamFinalRetentionGain+=sf-bf
		res.Points=append(res.Points,UP156CPoint{Cohort:names[ci],Keys:append([]int(nil),c...),InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,BaselineImmediateLoss:bl,GuidedImmediateLoss:gl,ShamImmediateLoss:sl,GuidedActionTarget:gt,ShamActionTarget:st,BaselineFinalProtected:bf,GuidedFinalProtected:gf,ShamFinalProtected:sf})
	}}
	return res,nil
}
