package unitary

const UP167CWarningConfirmationSchema="wingless.up167c-warning-confirmation.v1"

type UP167CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Confirmation int `json:"confirmation"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	EligibleLosses int `json:"eligible_losses"`
	MissedLosses int `json:"missed_losses"`
	EligibleIntervals int `json:"eligible_intervals"`
	ExpiredEligibleIntervals int `json:"expired_eligible_intervals"`
	TriggerPrecision float64 `json:"trigger_precision"`
	TriggerRecall float64 `json:"trigger_recall"`
	MeanFirstTriggerLeadWrites float64 `json:"mean_first_trigger_lead_writes"`
}
type UP167CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP166CSeal string `json:"source_up166c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	Confirmations []int `json:"confirmations"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	AdaptiveConfirmationUsed bool `json:"adaptive_confirmation_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP167CMetric `json:"metrics"`
}
func up167cRun(x0 *up81cAging,endangered,cadence,confirmation int,policy string)(lost,eligibleAtLoss bool,eligibleIntervals,expired,lead int){
	m:=*x0;streak:=0;firstEligibleStart:=0;lossStep:=0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		if warn{streak++}else{streak=0}
		eligible:=warn&&streak>=confirmation
		if eligible{eligibleIntervals++;if firstEligibleStart==0{firstEligibleStart=start}}
		intervalLost:=false
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,970000+step,step,policy){intervalLost=true;lossStep=step;break}
		}
		if intervalLost{lost=true;eligibleAtLoss=eligible;break}
		if eligible{expired++}
	}
	if lost&&firstEligibleStart>0{lead=lossStep-firstEligibleStart+1}
	return
}
func RunUP167C()(UP167CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4};confirms:=[]int{1,2,3}
	res:=UP167CResult{Schema:UP167CWarningConfirmationSchema,Experiment:"UP-167C-warning-confirmation",SourceUP166CSeal:"5e4879b1653490ebae62421f84620ce8e1c6a667",Policies:policies,Cadences:cadences,Confirmations:confirms,ShadowOnly:true,CorrectiveActionUsed:false,AdaptiveConfirmationUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{for _,confirm:=range confirms{
		m:=UP167CMetric{Policy:policy,Cadence:cad,Confirmation:confirm};leadSum:=0;leadN:=0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			lost,eligible,ei,exp,lead:=up167cRun(x,endangered,cad,confirm,policy)
			m.EligibleIntervals+=ei;m.ExpiredEligibleIntervals+=exp
			if lost{m.EventualLosses++;if eligible{m.EligibleLosses++}else{m.MissedLosses++};if lead>0{leadSum+=lead;leadN++}}
		}}
		if m.EligibleIntervals>0{m.TriggerPrecision=float64(m.EligibleLosses)/float64(m.EligibleIntervals)}
		if m.EventualLosses>0{m.TriggerRecall=float64(m.EligibleLosses)/float64(m.EventualLosses)}
		if leadN>0{m.MeanFirstTriggerLeadWrites=float64(leadSum)/float64(leadN)}
		res.Metrics=append(res.Metrics,m)
	}}}
	return res,nil
}
