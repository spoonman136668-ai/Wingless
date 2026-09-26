package unitary

const UP159BSubjectGeneralizationSchema="wingless.up159b-cleanup-order-subject-generalization.v1"

type UP159BArm struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	HeldoutNewNames []string `json:"heldout_new_names"`
	CleanupOrder string `json:"cleanup_order"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportCleanupRecovery float64 `json:"mean_post_report_cleanup_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalHeldoutNewAccuracy float64 `json:"final_heldout_new_accuracy"`
}
type UP159BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP158BSeal string `json:"source_up158b_seal"`
	Subjects int `json:"subjects"`
	Orders int `json:"orders"`
	Arms int `json:"arms"`
	TrainingPairs int `json:"training_pairs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	RehearsalBefore int `json:"rehearsal_before"`
	RehearsalAfter int `json:"rehearsal_after"`
	CleanupStore int `json:"cleanup_store"`
	CleanupObserve int `json:"cleanup_observe"`
	CleanupReport int `json:"cleanup_report"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Metrics []UP159BArm `json:"metrics"`
}

func up159bPair(subject string)[]string{
	for i:=0;i<len(up130bNewNames);i+=2{
		if up130bNewNames[i]==subject||up130bNewNames[i+1]==subject{return []string{up130bNewNames[i],up130bNewNames[i+1]}}
	}
	panic("UP159B_SUBJECT_NOT_IN_FROZEN_PAIRS")
}
func up159bHeldout(pair []string)[]string{
	out:=make([]string,0,4)
	for _,n:=range up130bNewNames{
		if n!=pair[0]&&n!=pair[1]{out=append(out,n)}
	}
	return out
}
func up159bExamples(pair []string)[]up135bExample{
	out:=make([]up135bExample,0,24)
	for _,name:=range pair{
		for _,v:=range up130bStore{out=append(out,up135bExample{name:name,verb:v,class:up97bStore})}
		for _,v:=range up130bObserve{out=append(out,up135bExample{name:name,verb:v,class:up97bObserve})}
		for _,v:=range up130bReport{out=append(out,up135bExample{name:name,verb:v,class:up97bReport})}
	}
	return out
}
func up159bEpoch(g *up129bGate,pair []string,subject,order string,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	count:=0
	for _,e:=range up159bExamples(pair){
		if up153bInBlock(e,reportBlock){continue}
		up135bStep(g,e,o,r);count++
	}
	if count!=20{panic("UP159B_PREFIX_COUNT")}
	afterPrefix:=up143bOldMean(g,o,r)
	preItems:=up158bItems(epoch,[]int{10,11,12})
	up156bApplyOld(g,preItems,0,len(preItems),o,r)
	afterPre:=up143bOldMean(g,o,r)
	for _,e:=range reportBlock{up135bStep(g,e,o,r)}
	afterReport:=up143bOldMean(g,o,r)
	postItems:=up158bItems(epoch,up158bIndicesForOrder(order))
	up156bApplyOld(g,postItems,0,len(postItems),o,r)
	afterAll:=up143bOldMean(g,o,r)
	return afterPrefix-before,afterPre-afterPrefix,afterReport-afterPre,afterAll-afterReport,afterAll-before
}
func up159bRun(common *up129bGate,subject,order string,o,r [64]float64)UP159BArm{
	pair:=up159bPair(subject);held:=up159bHeldout(pair)
	g:=*common
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up159bEpoch(&g,pair,subject,order,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	hn,_:=up130bSplit(&g,"heldout_new",held,o,r)
	return UP159BArm{
		Subject:subject,TrainingPair:append([]string(nil),pair...),HeldoutNewNames:append([]string(nil),held...),CleanupOrder:order,
		MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,
		MeanPostReportCleanupRecovery:sa/5,MeanNetEpochChange:sn/5,
		FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalHeldoutNewAccuracy:hn.OverallAccuracy,
	}
}
func RunUP159B()(UP159BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	orders:=[]string{"STORE_OBSERVE_REPORT","STORE_REPORT_OBSERVE","OBSERVE_STORE_REPORT","OBSERVE_REPORT_STORE","REPORT_STORE_OBSERVE","REPORT_OBSERVE_STORE"}
	res:=UP159BResult{
		Schema:UP159BSubjectGeneralizationSchema,Experiment:"UP-159B-cleanup-order-subject-generalization",
		SourceUP158BSeal:"50d4146aa51fca6289a8271ca1e92b9a8ce801bb",
		Subjects:6,Orders:6,Arms:36,TrainingPairs:3,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,
		RehearsalBefore:3,RehearsalAfter:12,CleanupStore:5,CleanupObserve:5,CleanupReport:2,
		ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false,
	}
	for _,subject:=range up130bNewNames{
		for _,order:=range orders{res.Metrics=append(res.Metrics,up159bRun(common,subject,order,o,r))}
	}
	return res,nil
}
