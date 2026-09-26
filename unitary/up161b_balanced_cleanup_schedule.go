package unitary

const UP161BBalancedCleanupSchema="wingless.up161b-balanced-cleanup-schedule.v1"

type UP161BArm struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	Schedule string `json:"schedule"`
	EpochOrders []string `json:"epoch_orders"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportCleanupRecovery float64 `json:"mean_post_report_cleanup_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalHeldoutNewAccuracy float64 `json:"final_heldout_new_accuracy"`
}
type UP161BScheduleSummary struct{
	Schedule string `json:"schedule"`
	MeanOldRetention float64 `json:"mean_old_retention"`
	MinOldRetention float64 `json:"min_old_retention"`
	MaxOldRetention float64 `json:"max_old_retention"`
	SubjectRetentionSpread float64 `json:"subject_retention_spread"`
}
type UP161BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP160BSeal string `json:"source_up160b_seal"`
	Subjects int `json:"subjects"`
	Schedules int `json:"schedules"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP161BArm `json:"arms"`
	Summaries []UP161BScheduleSummary `json:"summaries"`
}
func up161bOrders(schedule string)[]string{
	switch schedule{
	case "balanced_forward":
		return []string{"STORE_OBSERVE_REPORT","OBSERVE_REPORT_STORE","REPORT_STORE_OBSERVE","STORE_OBSERVE_REPORT","OBSERVE_REPORT_STORE"}
	case "balanced_reverse":
		return []string{"STORE_REPORT_OBSERVE","REPORT_OBSERVE_STORE","OBSERVE_STORE_REPORT","STORE_REPORT_OBSERVE","REPORT_OBSERVE_STORE"}
	default:
		return []string{"STORE_OBSERVE_REPORT","STORE_OBSERVE_REPORT","STORE_OBSERVE_REPORT","STORE_OBSERVE_REPORT","STORE_OBSERVE_REPORT"}
	}
}
func up161bRun(common *up129bGate,subject,schedule string,o,r [64]float64)UP161BArm{
	g:=*common;pair:=up159bPair(subject);held:=up159bHeldout(pair);orders:=up161bOrders(schedule)
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for i,epoch:=range []int{15,16,17,18,19}{
		p,b,rr,a,n:=up159bEpoch(&g,pair,subject,orders[i],epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	hn,_:=up130bSplit(&g,"heldout_new",held,o,r)
	return UP161BArm{Subject:subject,TrainingPair:append([]string(nil),pair...),Schedule:schedule,EpochOrders:append([]string(nil),orders...),MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,MeanPostReportCleanupRecovery:sa/5,MeanNetEpochChange:sn/5,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,FinalHeldoutNewAccuracy:hn.OverallAccuracy}
}
func up161bSummary(schedule string,arms []UP161BArm)UP161BScheduleSummary{
	sum:=0.0;min,max:=1.0,0.0;n:=0
	for _,a:=range arms{
		if a.Schedule!=schedule{continue};n++;sum+=a.FinalMeanOldRetention
		if a.FinalMeanOldRetention<min{min=a.FinalMeanOldRetention}
		if a.FinalMeanOldRetention>max{max=a.FinalMeanOldRetention}
	}
	return UP161BScheduleSummary{Schedule:schedule,MeanOldRetention:sum/float64(n),MinOldRetention:min,MaxOldRetention:max,SubjectRetentionSpread:max-min}
}
func RunUP161B()(UP161BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP161BResult{Schema:UP161BBalancedCleanupSchema,Experiment:"UP-161B-balanced-cleanup-schedule",SourceUP160BSeal:"d4a12a5ecd2d5f24671d1af3a2a7267732f10fdd",Subjects:6,Schedules:3,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false}
	schedules:=[]string{"canonical_fixed","balanced_forward","balanced_reverse"}
	for _,subject:=range up130bNewNames{for _,schedule:=range schedules{res.Arms=append(res.Arms,up161bRun(common,subject,schedule,o,r))}}
	for _,schedule:=range schedules{res.Summaries=append(res.Summaries,up161bSummary(schedule,res.Arms))}
	return res,nil
}
