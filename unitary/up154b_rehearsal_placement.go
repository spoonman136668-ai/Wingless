package unitary

const UP154BRehearsalPlacementSchema="wingless.up154b-rehearsal-placement.v1"

type UP154BArm struct{
	Subject string `json:"subject"`
	Arm string `json:"arm"`
	RehearsalBefore int `json:"rehearsal_before"`
	RehearsalAfter int `json:"rehearsal_after"`
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
type UP154BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP153BSeal string `json:"source_up153b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	ReportBlockUpdates int `json:"report_block_updates"`
	Subjects int `json:"subjects"`
	ArmsPerSubject int `json:"arms_per_subject"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptivePlacementUsed bool `json:"adaptive_placement_used"`
	Arms []UP154BArm `json:"arms"`
}

func up154bOldSlice(g *up129bGate,epoch,start,end int,o,r [64]float64){
	old:=up132bOldSurfaces()
	for i:=start;i<end;i++{
		s:=old[i]
		subjectIndex:=(i+epoch)%2
		up129bGateStep(g,up121bOriginalNames[subjectIndex],s.verb,s.class,o,r)
	}
}
func up154bReportBlock(subject string)[]up135bExample{
	return up153bBlock(subject,[]string{"relays","announces","cites","summarizes"})
}
func up154bEpoch(g *up129bGate,subject,arm string,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){
		if up153bInBlock(e,reportBlock){continue}
		up135bStep(g,e,o,r)
	}
	afterPrefix:=up143bOldMean(g,o,r)
	beforeN,afterN:=15,0
	switch arm{
	case "split_8_7": beforeN,afterN=8,7
	case "all_after": beforeN,afterN=0,15
	}
	if beforeN>0{up154bOldSlice(g,epoch,0,beforeN,o,r)}
	afterPre:=up143bOldMean(g,o,r)
	for _,e:=range reportBlock{up135bStep(g,e,o,r)}
	afterReport:=up143bOldMean(g,o,r)
	if afterN>0{up154bOldSlice(g,epoch,beforeN,15,o,r)}
	afterAll:=up143bOldMean(g,o,r)
	return afterPrefix-before,afterPre-afterPrefix,afterReport-afterPre,afterAll-afterReport,afterAll-before
}
func up154bRun(common *up129bGate,subject,arm string,o,r [64]float64)UP154BArm{
	g:=*common
	prefixState:=up143bOldMean(&g,o,r)
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	beforeN,afterN:=15,0
	switch arm{
	case "split_8_7": beforeN,afterN=8,7
	case "all_after": beforeN,afterN=0,15
	}
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up154bEpoch(&g,subject,arm,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP154BArm{
		Subject:subject,Arm:arm,RehearsalBefore:beforeN,RehearsalAfter:afterN,PrefixStateOldRetention:prefixState,
		MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,MeanPostReportRehearsalRecovery:sa/5,MeanNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:pn,FinalSecondaryNew:snM,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP154B()(UP154BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP154BResult{
		Schema:UP154BRehearsalPlacementSchema,Experiment:"UP-154B-rehearsal-placement",
		SourceUP153BSeal:"00e708d49751c637e3bb2f82018a8f314bdf1ac6",
		CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,ReportBlockUpdates:4,
		Subjects:2,ArmsPerSubject:3,ExtraUpdatesUsed:false,AdaptivePlacementUsed:false,
	}
	for _,subject:=range []string{"mia","noah"}{
		for _,arm:=range []string{"all_before","split_8_7","all_after"}{
			res.Arms=append(res.Arms,up154bRun(common,subject,arm,o,r))
		}
	}
	return res,nil
}
