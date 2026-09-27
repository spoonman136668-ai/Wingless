package unitary

const UP161CHorizonSchema="wingless.up161c-bounded-protection-horizon.v1"

type UP161CPoint struct{
	Cohort string `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	NoQueryHorizon int `json:"no_query_horizon"`
	AdversarialShieldHorizon int `json:"adversarial_shield_horizon"`
	TestedRobustHorizon int `json:"tested_robust_horizon"`
}
type UP161CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP160CSeal string `json:"source_up160c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	Arms int `json:"arms"`
	MinRobustHorizon int `json:"min_robust_horizon"`
	MeanRobustHorizon float64 `json:"mean_robust_horizon"`
	MaxRobustHorizon int `json:"max_robust_horizon"`
	Points []UP161CPoint `json:"points"`
}
func up161cTriggerState(hand int,cohort []int)(*up81cAging,int,bool){
	x:=up156cInit(hand)
	for step:=1;step<=24;step++{
		if rk,ok:=up159cRefresh(step);ok{x.query(rk)}
		key:=up159cStreamKey(step)
		if x.find(key)<0{
			slot:=up151cPredict(x.hand,x.age)
			if x.entries[slot].used&&up155cHas(cohort,x.entries[slot].key){
				endangered:=x.entries[slot].key;x.query(endangered);return x,endangered,true
			}
		}
		x.write(key,step%3)
	}
	return nil,-1,false
}
func up161cAdversarial(x *up81cAging,endangered int)int{
	m:=*x
	for n:=1;n<=64;n++{
		slot:=up151cPredict(m.hand,m.age)
		if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
		m.write(910000+n,n%3)
		if m.find(endangered)<0{return n}
	}
	return 65
}
func RunUP161C()(UP161CResult,error){
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"}
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};vals:=[]int{}
	res:=UP161CResult{Schema:UP161CHorizonSchema,Experiment:"UP-161C-bounded-protection-horizon",SourceUP160CSeal:"c70c1f53358d88dd984080724562242920c19ff1",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue}
		noq:=up158cHorizon(x,endangered);adv:=up161cAdversarial(x,endangered);rob:=noq;if adv<rob{rob=adv}
		vals=append(vals,rob);res.Arms++
		res.Points=append(res.Points,UP161CPoint{Cohort:names[ci],InitialHand:hand,EndangeredKey:endangered,NoQueryHorizon:noq,AdversarialShieldHorizon:adv,TestedRobustHorizon:rob})
	}}
	if len(vals)>0{res.MinRobustHorizon,res.MaxRobustHorizon=up158cMinMax(vals);res.MeanRobustHorizon=up158cMean(vals)}
	return res,nil
}
