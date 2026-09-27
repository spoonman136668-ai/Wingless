package unitary

const UP206CWarningTraceSchema="wingless.up206c-repeated-switch-warning-trace.v1"

type UP206CPoint struct{
	Schedule string `json:"schedule"`
	CohortIndex int `json:"cohort_index"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	BaselineLossStep int `json:"baseline_loss_step"`
	TreatedLossStep int `json:"treated_loss_step"`
	TriggerOnset int `json:"trigger_onset"`
	ActionWrites []int `json:"action_writes"`
	LastMonitorBeforeLoss int `json:"last_monitor_before_loss"`
	HorizonAtLastMonitor int `json:"horizon_at_last_monitor"`
	LastCriticalBeforeLoss int `json:"last_critical_before_loss"`
	CriticalLeadWrites int `json:"critical_lead_writes"`
	PostTriggerCriticalCount int `json:"post_trigger_critical_count"`
	SurvivesHorizon bool `json:"survives_horizon"`
}
type UP206CSummary struct{
	Schedule string `json:"schedule"`
	Arms int `json:"arms"`
	TreatedFailures int `json:"treated_failures"`
	FailuresCriticalAtFinalMonitor int `json:"failures_critical_at_final_monitor"`
	FailuresWithPostTriggerCritical int `json:"failures_with_post_trigger_critical"`
	MeanCriticalLeadWrites float64 `json:"mean_critical_lead_writes"`
	MinCriticalLeadWrites int `json:"min_critical_lead_writes"`
	MaxCriticalLeadWrites int `json:"max_critical_lead_writes"`
}
type UP206CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP205CSeal string `json:"source_up205c_seal"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	RuleChanged bool `json:"rule_changed"`
	DiagnosticObservationOnly bool `json:"diagnostic_observation_only"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Points []UP206CPoint `json:"points"`
	Summaries []UP206CSummary `json:"summaries"`
}
func up206cTreatedTrace(x0 *up81cAging,endangered int,a,b string)(loss,onset int,writes []int,lastMonitor,lastH,lastCritical,postCritical int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=56;start+=2{
		preH:=up161cAdversarial(&m,endangered)
		if committed&&preH<=2{lastCritical=start;postCritical++}
		if !committed{
			if preH<=2{
				m.query(endangered);writes=append(writes,start);onset=start;committed=true;since=0
			}
		}else if len(writes)<8{
			since++
			if since>=4{m.query(endangered);writes=append(writes,start);since=0}
		}
		lastMonitor=start;lastH=preH
		for j:=0;j<2&&start+j<=56;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1029000+step,step,up204cPolicy(step,a,b)){return step,onset,writes,lastMonitor,lastH,lastCritical,postCritical}
		}
	}
	return 57,onset,writes,lastMonitor,lastH,lastCritical,postCritical
}
func RunUP206C()(UP206CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile_repeated","no_refresh","hostile_shield"},
		{"hostile_no_repeated","hostile_shield","no_refresh"},
		{"alternating_fixed_repeated","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating_repeated","fixed_offset_refresh","alternating_shield"},
	}
	res:=UP206CResult{Schema:UP206CWarningTraceSchema,Experiment:"UP-206C-repeated-switch-warning-trace",SourceUP205CSeal:"5f797af968dca22f1007e5035e64bc868ecd7b3b",Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:56,RuleChanged:false,DiagnosticObservationOnly:true,CounterfactualOnly:true,LiveActivation:false}
	for _,s:=range schedules{
		sm:=UP206CSummary{Schedule:s.name,MinCriticalLeadWrites:1<<30};sumLead,nLead:=0,0
		for ci,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};sm.Arms++
			base:=up204cBaseline(x,e,s.a,s.b)
			loss,onset,writes,lastMon,lastH,lastCrit,postCrit:=up206cTreatedTrace(x,e,s.a,s.b)
			surv:=loss>56;lead:=0
			if !surv{
				sm.TreatedFailures++
				if lastH<=2{sm.FailuresCriticalAtFinalMonitor++}
				if lastCrit>0{
					sm.FailuresWithPostTriggerCritical++;lead=loss-lastCrit;sumLead+=lead;nLead++
					if lead<sm.MinCriticalLeadWrites{sm.MinCriticalLeadWrites=lead};if lead>sm.MaxCriticalLeadWrites{sm.MaxCriticalLeadWrites=lead}
				}
			}
			res.Points=append(res.Points,UP206CPoint{Schedule:s.name,CohortIndex:ci,InitialHand:hand,EndangeredKey:e,BaselineLossStep:base,TreatedLossStep:loss,TriggerOnset:onset,ActionWrites:append([]int(nil),writes...),LastMonitorBeforeLoss:lastMon,HorizonAtLastMonitor:lastH,LastCriticalBeforeLoss:lastCrit,CriticalLeadWrites:lead,PostTriggerCriticalCount:postCrit,SurvivesHorizon:surv})
		}}
		if nLead>0{sm.MeanCriticalLeadWrites=float64(sumLead)/float64(nLead)}
		if sm.MinCriticalLeadWrites==1<<30{sm.MinCriticalLeadWrites=0}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
