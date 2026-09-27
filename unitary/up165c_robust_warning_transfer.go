package unitary

const UP165CTransferSchema="wingless.up165c-robust-warning-transfer.v1"

type UP165CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	ExpiredWarningIntervals int `json:"expired_warning_intervals"`
	WarningRecall float64 `json:"warning_recall"`
}
type UP165CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP164CSeal string `json:"source_up164c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	WarningModelChangedAcrossPolicies bool `json:"warning_model_changed_across_policies"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP165CMetric `json:"metrics"`
}
func up165cRealStep(m *up81cAging,endangered,key,step int,policy string)bool{
	switch policy{
	case "alternating_shield":
		if step%2==1{
			slot:=up151cPredict(m.hand,m.age)
			if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
		}
	case "fixed_offset_refresh":
		slot:=(m.hand+5)%16
		if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
	case "hostile_shield":
		slot:=up151cPredict(m.hand,m.age)
		if m.entries[slot].used&&m.entries[slot].key!=endangered{m.query(m.entries[slot].key)}
	}
	m.write(key,key%3)
	return m.find(endangered)<0
}
func up165cRun(x0 *up81cAging,endangered,cadence int,policy string)(lost,warning bool,expired int){
	m:=*x0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		intervalLost:=false
		for j:=0;j<cadence&&start+j<=64;j++{
			if up165cRealStep(&m,endangered,950000+start+j,start+j,policy){intervalLost=true;break}
		}
		if intervalLost{return true,warn,expired}
		if warn{expired++}
	}
	return false,false,expired
}
func RunUP165C()(UP165CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4}
	res:=UP165CResult{Schema:UP165CTransferSchema,Experiment:"UP-165C-robust-warning-transfer",SourceUP164CSeal:"e388b14467b099d0fce933570fc2d6c549d6d5c0",ExactRecallCap:16,Policies:policies,Cadences:cadences,ShadowOnly:true,CorrectiveActionUsed:false,WarningModelChangedAcrossPolicies:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{
		m:=UP165CMetric{Policy:policy,Cadence:cad}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			lost,warn,expired:=up165cRun(x,endangered,cad,policy);m.ExpiredWarningIntervals+=expired
			if lost{m.EventualLosses++;if warn{m.WarnedLosses++}else{m.MissedLosses++}}
		}}
		if m.EventualLosses>0{m.WarningRecall=float64(m.WarnedLosses)/float64(m.EventualLosses)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
