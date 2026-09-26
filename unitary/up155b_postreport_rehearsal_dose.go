package unitary

const UP155BPostReportDoseSchema="wingless.up155b-postreport-rehearsal-dose.v1"

type UP155BArm struct{
	Subject string `json:"subject"`
	RehearsalAfter int `json:"rehearsal_after"`
	RehearsalBefore int `json:"rehearsal_before"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportRehearsalRecovery float64 `json:"mean_post_report_rehearsal_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP155BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP154BSeal string `json:"source_up154b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	ReportBlockUpdates int `json:"report_block_updates"`
	Subjects int `json:"subjects"`
	DosePoints int `json:"dose_points"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptivePlacementUsed bool `json:"adaptive_placement_used"`
	Arms []UP155BArm `json:"arms"`
}

func up155bEpoch(g *up129bGate,subject string,afterN,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){
		if up153bInBlock(e,reportBlock){continue}
		up135bStep(g,e,o,r)
	}
	afterPrefix:=up143bOldMean(g,o,r)
	beforeN:=15-afterN
	if beforeN>0{up154bOldSlice(g,epoch,0,beforeN,o,r)}
	afterPre:=up143bOldMean(g,o,r)
	for _,e:=range reportBlock{up135bStep(g,e,o,r)}
	afterReport:=up143bOldMean(g,o,r)
	if afterN>0{up154bOldSlice(g,epoch,beforeN,15,o,r)}
	afterAll:=up143bOldMean(g,o,r)
	return afterPrefix-before,afterPre-afterPrefix,afterReport-afterPre,afterAll-afterReport,afterAll-before
}
func up155bRun(common *up129bGate,subject string,afterN int,o,r [64]float64)UP155BArm{
	g:=*common
	prefixState:=up143bOldMean(&g,o,r)
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up155bEpoch(&g,subject,afterN,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP155BArm{
		Subject:subject,RehearsalAfter:afterN,RehearsalBefore:15-afterN,PrefixStateOldRetention:prefixState,
		MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,MeanPostReportRehearsalRecovery:sa/5,MeanNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:pn,FinalSecondaryNew:snM,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP155B()(UP155BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP155BResult{
		Schema:UP155BPostReportDoseSchema,Experiment:"UP-155B-postreport-rehearsal-dose",
		SourceUP154BSeal:"315f2b88ae66d87c15b2765b4de6fba9943b886b",
		CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,ReportBlockUpdates:4,
		Subjects:2,DosePoints:5,ExtraUpdatesUsed:false,AdaptivePlacementUsed:false,
	}
	for _,subject:=range []string{"mia","noah"}{
		for _,afterN:=range []int{0,4,8,12,15}{
			res.Arms=append(res.Arms,up155bRun(common,subject,afterN,o,r))
		}
	}
	return res,nil
}
