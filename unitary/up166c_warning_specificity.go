package unitary

const UP166CSpecificitySchema="wingless.up166c-warning-specificity.v1"

type UP166CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	EventualLosses int `json:"eventual_losses"`
	WarningIntervals int `json:"warning_intervals"`
	WarnedLosses int `json:"warned_losses"`
	MissedLosses int `json:"missed_losses"`
	ExpiredWarningIntervals int `json:"expired_warning_intervals"`
	WarningPrecision float64 `json:"warning_precision"`
	WarningRecall float64 `json:"warning_recall"`
	MeanFirstWarningLeadWrites float64 `json:"mean_first_warning_lead_writes"`
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
	WarningModelChangedAcrossPolicies bool `json:"warning_model_changed_across_policies"`
	FutureRealPolicyScheduleUsedByWarning bool `json:"future_real_policy_schedule_used_by_warning"`
	SafeControlIncluded bool `json:"safe_control_included"`
	Metrics []UP166CMetric `json:"metrics"`
}

func up166cRealStep(m *up81cAging,endangered,key,step int,policy string)bool{
	if policy=="protected_endangered"{
		m.query(endangered)
		m.write(key,key%3)
		return m.find(endangered)<0
	}
	return up165cRealStep(m,endangered,key,step,policy)
}

func up166cRun(x0 *up81cAging,endangered,cadence int,policy string)(lost,warnedAtLoss bool,warningIntervals,expired,leadWrites int){
	m:=*x0
	firstWarningStart:=0
	for start:=1;start<=64;start+=cadence{
		h:=up161cAdversarial(&m,endangered)
		warn:=h<=cadence
		if warn{
			warningIntervals++
			if firstWarningStart==0{firstWarningStart=start}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up166cRealStep(&m,endangered,960000+step,step,policy){
				lost=true
				warnedAtLoss=warn
				if firstWarningStart>0{leadWrites=step-firstWarningStart+1}
				return
			}
		}
		if warn{expired++}
	}
	return
}

func RunUP166C()(UP166CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,4,8,12}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield","protected_endangered"}
	cadences:=[]int{2,4}
	res:=UP166CResult{
		Schema:UP166CSpecificitySchema,
		Experiment:"UP-166C-warning-specificity",
		SourceUP165CSeal:"f93cd9cdc9e4b371dcd46f88ef6aa987ee9cb792",
		ExactRecallCap:16,
		Policies:policies,
		Cadences:cadences,
		ShadowOnly:true,
		CorrectiveActionUsed:false,
		WarningModelChangedAcrossPolicies:false,
		FutureRealPolicyScheduleUsedByWarning:false,
		SafeControlIncluded:true,
	}
	for _,policy:=range policies{
		for _,cad:=range cadences{
			m:=UP166CMetric{Policy:policy,Cadence:cad}
			leadTotal:=0
			for _,c:=range cohorts{
				for _,hand:=range hands{
					x,endangered,ok:=up161cTriggerState(hand,c)
					if !ok{continue}
					m.Arms++
					lost,warned,wi,exp,lead:=up166cRun(x,endangered,cad,policy)
					m.WarningIntervals+=wi
					m.ExpiredWarningIntervals+=exp
					if lost{
						m.EventualLosses++
						if warned{m.WarnedLosses++;leadTotal+=lead}else{m.MissedLosses++}
					}
				}
			}
			if m.WarningIntervals>0{m.WarningPrecision=float64(m.WarnedLosses)/float64(m.WarningIntervals)}
			if m.EventualLosses>0{m.WarningRecall=float64(m.WarnedLosses)/float64(m.EventualLosses)}
			if m.WarnedLosses>0{m.MeanFirstWarningLeadWrites=float64(leadTotal)/float64(m.WarnedLosses)}
			res.Metrics=append(res.Metrics,m)
		}
	}
	return res,nil
}
