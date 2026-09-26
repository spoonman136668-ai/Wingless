package unitary

const UP160CGeometrySchema="wingless.up160c-horizon16-geometry.v1"

type UP160CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	TriggerStep int `json:"trigger_step"`
	EndangeredKey int `json:"endangered_key"`
	UniqueWriteHorizon int `json:"unique_write_horizon"`
	CurrentHand int `json:"current_hand"`
	EndangeredSlot int `json:"endangered_slot"`
	CyclicDistance int `json:"cyclic_distance"`
	EndangeredAge int `json:"endangered_age"`
	ZeroAgeSlots int `json:"zero_age_slots"`
	NonzeroAgeSlots int `json:"nonzero_age_slots"`
	AgeSum int `json:"age_sum"`
	MaxAge int `json:"max_age"`
	AgeGEEndangered int `json:"age_ge_endangered"`
	ActualDurable bool `json:"actual_durable"`
	ActualLossStep int `json:"actual_loss_step"`
}
type UP160CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP159CSeal string `json:"source_up159c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	PrimaryHorizon int `json:"primary_horizon"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	SemanticClassUsed bool `json:"semantic_class_used"`
	InterventionPolicyChanged bool `json:"intervention_policy_changed"`
	PrimaryCases int `json:"primary_cases"`
	DurableCases int `json:"durable_cases"`
	DelayedLossCases int `json:"delayed_loss_cases"`
	Points []UP160CPoint `json:"points"`
}
func up160cSnapshot(x *up81cAging,key int)(slot,dist,eage,zeros,nonzero,sum,max,ge int){
	slot=x.find(key);if slot<0{return}
	dist=(slot-x.hand+16)%16;eage=int(x.age[slot])
	for i:=0;i<16;i++{a:=int(x.age[i]);sum+=a;if a==0{zeros++}else{nonzero++};if a>max{max=a};if a>=eage{ge++}}
	return
}
func up160cArm(hand int,trigger,endangered int)(UP160CPoint,bool){
	x:=up156cInit(hand);var p UP160CPoint;captured:=false;loss:=-1
	for step:=1;step<=24;step++{
		if rk,ok:=up159cRefresh(step);ok{x.query(rk)}
		key:=up159cStreamKey(step)
		if step==trigger{
			x.query(endangered);h:=up158cHorizon(x,endangered)
			if h==16{
				slot,dist,eage,z,nz,sum,max,ge:=up160cSnapshot(x,endangered)
				p=UP160CPoint{InitialHand:hand,TriggerStep:trigger,EndangeredKey:endangered,UniqueWriteHorizon:h,CurrentHand:x.hand,EndangeredSlot:slot,CyclicDistance:dist,EndangeredAge:eage,ZeroAgeSlots:z,NonzeroAgeSlots:nz,AgeSum:sum,MaxAge:max,AgeGEEndangered:ge};captured=true
			}
		}
		x.write(key,step%3)
		if loss<0&&x.find(endangered)<0{loss=step}
	}
	if captured{p.ActualDurable=x.find(endangered)>=0;p.ActualLossStep=loss}
	return p,captured
}
func RunUP160C()(UP160CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"};cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}};hands:=[]int{0,4,8,12}
	res:=UP160CResult{Schema:UP160CGeometrySchema,Experiment:"UP-160C-horizon16-geometry",SourceUP159CSeal:"3f621ebede6fc7b4aa09bcb0eeb5b06bde295943",ExactRecallCap:16,PrimaryHorizon:16,FutureScheduleUsed:false,SemanticClassUsed:false,InterventionPolicyChanged:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		trigger,endangered:=up159cFindTrigger(hand,c);if trigger<0{continue}
		p,ok:=up160cArm(hand,trigger,endangered);if !ok{continue};p.Cohort=names[ci];res.Points=append(res.Points,p);res.PrimaryCases++;if p.ActualDurable{res.DurableCases++}else{res.DelayedLossCases++}
	}}
	return res,nil
}
