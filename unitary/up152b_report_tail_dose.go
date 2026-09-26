package unitary

const UP152BReportDoseSchema="wingless.up152b-report-tail-dose.v1"

type UP152BArm struct{
	Subject string `json:"subject"`
	ReportCount int `json:"report_count"`
	TailVerbs []string `json:"tail_verbs"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanTerminalPrefixDamage float64 `json:"mean_terminal_prefix_damage"`
	MeanTerminalAnchorRecovery float64 `json:"mean_terminal_anchor_recovery"`
	MeanTerminalTailDamage float64 `json:"mean_terminal_tail_damage"`
	MeanTerminalNetEpochChange float64 `json:"mean_terminal_net_epoch_change"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP152BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP151BSeal string `json:"source_up151b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	Subjects int `json:"subjects"`
	ReportDosePoints int `json:"report_dose_points"`
	AdaptiveTailSelection bool `json:"adaptive_tail_selection"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP152BArm `json:"arms"`
}

func up152bVerbClass(verb string)int{
	for _,v:=range up130bStore{if v==verb{return up97bStore}}
	for _,v:=range up130bObserve{if v==verb{return up97bObserve}}
	return up97bReport
}
func up152bTail(subject string,reportCount int)[]up135bExample{
	sets:=[][]string{
		{"caches","files","scans","checks"},
		{"caches","files","scans","relays"},
		{"caches","files","relays","announces"},
		{"caches","relays","announces","cites"},
		{"relays","announces","cites","summarizes"},
	}
	out:=make([]up135bExample,0,4)
	for _,verb:=range sets[reportCount]{out=append(out,up135bExample{name:subject,verb:verb,class:up152bVerbClass(verb)})}
	return out
}
func up152bRun(common *up129bGate,subject string,reportCount int,o,r [64]float64)UP152BArm{
	g:=*common
	tail:=up152bTail(subject,reportCount)
	prefixState:=up143bOldMean(&g,o,r)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,a,t,n:=up151bEpoch(&g,tail,epoch,o,r)
		sp+=p;sa+=a;st+=t;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	verbs:=make([]string,0,4);for _,e:=range tail{verbs=append(verbs,e.verb)}
	return UP152BArm{
		Subject:subject,ReportCount:reportCount,TailVerbs:verbs,PrefixStateOldRetention:prefixState,
		MeanTerminalPrefixDamage:sp/5,MeanTerminalAnchorRecovery:sa/5,MeanTerminalTailDamage:st/5,MeanTerminalNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:pn,FinalSecondaryNew:snM,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP152B()(UP152BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP152BResult{
		Schema:UP152BReportDoseSchema,Experiment:"UP-152B-report-tail-dose",
		SourceUP151BSeal:"069a1fbc5c0db1b16c1a8b51acf004f8f8685229",
		CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		Subjects:2,ReportDosePoints:5,AdaptiveTailSelection:false,ExtraUpdatesUsed:false,
	}
	for _,subject:=range []string{"mia","noah"}{
		for rc:=0;rc<=4;rc++{res.Arms=append(res.Arms,up152bRun(common,subject,rc,o,r))}
	}
	return res,nil
}
