package unitary

const UP166CWarningPrecisionSchema="wingless.up166c-warning-precision.v1"

type UP166CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	TotalWarningIntervals int `json:"total_warning_intervals"`
	ExpiredWarningIntervals int `json:"expired_warning_intervals"`
	IntervalPrecision float64 `json:"interval_precision"`
	WarningIntervalsPerLoss float64 `json:"warning_intervals_per_loss"`
	MeanFirstWarningLeadWrites float64 `json:"mean_first_warning_lead_writes"`
	MinFirstWarningLeadWrites int `json:"min_first_warning_lead_writes"`
	MaxFirstWarningLeadWrites int `json:"max_first_warning_lead_writes"`
}
type UP166CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP165CSeal string `json:"source_up165c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	ShadowOnly bool `json:"shadow_only"`
	CorrectiveActionUsed bool `json:"corrective_action_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	Metrics []UP166CMetric `json:"metrics"`
}

func up166cRun(x0 *up81cAging,endangered,cadence int,policy string)(lost,warnAtLoss bool,totalWarn,expired,firstLead int){
	m:=*x0;firstWarnStart:=0;lossStep:=0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		if warn{totalWarn++;if firstWarnStart==0{firstWarnStart=start}}
		intervalLost:=false
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,960000+step,step,policy){intervalLost=true;lossStep=step;break}
		}
		if intervalLost{
			if warn{warnAtLoss=true}
			lost=true
			break
		}
		if warn{expired++}
	}
	if lost&&firstWarnStart>0{firstLead=lossStep-firstWarnStart+1}
	return
}
func RunUP166C()(UP166CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12};policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"};cadences:=[]int{2,4}
	res:=UP166CResult{Schema:UP166CWarningPrecisionSchema,Experiment:"UP-166C-warning-precision",SourceUP165CSeal:"f93cd9cdc9e4b371dcd46f88ef6aa987ee9cb792",ExactRecallCap:16,Policies:policies,Cadences:cadences,ShadowOnly:true,CorrectiveActionUsed:false,WarningThresholdChanged:false,FutureRealPolicyScheduleUsedByWarning:false}
	for _,policy:=range policies{for _,cad:=range cadences{
		m:=UP166CMetric{Policy:policy,Cadence:cad,MinFirstWarningLeadWrites:1<<30}
		leadSum:=0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,endangered,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			lost,warned,tw,exp,lead:=up166cRun(x,endangered,cad,policy)
			m.TotalWarningIntervals+=tw;m.ExpiredWarningIntervals+=exp
			if lost{
				m.EventualLosses++
				if warned{m.WarnedLosses++}else{m.MissedLosses++}
				if lead>0{leadSum+=lead;if lead<m.MinFirstWarningLeadWrites{m.MinFirstWarningLeadWrites=lead};if lead>m.MaxFirstWarningLeadWrites{m.MaxFirstWarningLeadWrites=lead}}
			}
		}}
		if m.TotalWarningIntervals>0{m.IntervalPrecision=float64(m.WarnedLosses)/float64(m.TotalWarningIntervals)}
		if m.EventualLosses>0{
			m.WarningIntervalsPerLoss=float64(m.TotalWarningIntervals)/float64(m.EventualLosses)
			m.MeanFirstWarningLeadWrites=float64(leadSum)/float64(m.EventualLosses)
		}
		if m.MinFirstWarningLeadWrites==1<<30{m.MinFirstWarningLeadWrites=0}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
