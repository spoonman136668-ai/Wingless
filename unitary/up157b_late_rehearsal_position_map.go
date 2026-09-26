package unitary

const UP157BPositionMapSchema="wingless.up157b-late-rehearsal-position-map.v1"

type UP157BArm struct{
	Subject string `json:"subject"`
	Rotation int `json:"rotation"`
	PreReportOriginalIndices []int `json:"pre_report_original_indices"`
	PostReportOriginalIndices []int `json:"post_report_original_indices"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportRehearsalRecovery float64 `json:"mean_post_report_rehearsal_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP157BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP156BSeal string `json:"source_up156b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	RehearsalAfter int `json:"rehearsal_after"`
	RehearsalBefore int `json:"rehearsal_before"`
	Rotations int `json:"rotations"`
	Subjects int `json:"subjects"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Arms []UP157BArm `json:"arms"`
}

func up157bRotate(items []up156bOldItem,shift int)[]up156bOldItem{
	n:=len(items);out:=make([]up156bOldItem,0,n)
	for i:=0;i<n;i++{out=append(out,items[(shift+i)%n])}
	return out
}
func up157bEpoch(g *up129bGate,subject string,rotation,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){if !up153bInBlock(e,reportBlock){up135bStep(g,e,o,r)}}
	afterPrefix:=up143bOldMean(g,o,r)
	items:=up157bRotate(up156bOldItems(epoch),rotation)
	up156bApplyOld(g,items,0,3,o,r)
	afterPre:=up143bOldMean(g,o,r)
	for _,e:=range reportBlock{up135bStep(g,e,o,r)}
	afterReport:=up143bOldMean(g,o,r)
	up156bApplyOld(g,items,3,15,o,r)
	afterAll:=up143bOldMean(g,o,r)
	return afterPrefix-before,afterPre-afterPrefix,afterReport-afterPre,afterAll-afterReport,afterAll-before
}
func up157bRun(common *up129bGate,subject string,rotation int,o,r [64]float64)UP157BArm{
	g:=*common
	prefixState:=up143bOldMean(&g,o,r)
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up157bEpoch(&g,subject,rotation,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	items:=up157bRotate(up156bOldItems(15),rotation)
	preIdx:=make([]int,0,3);postIdx:=make([]int,0,12)
	for i,it:=range items{if i<3{preIdx=append(preIdx,it.originalIndex)}else{postIdx=append(postIdx,it.originalIndex)}}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP157BArm{
		Subject:subject,Rotation:rotation,PreReportOriginalIndices:preIdx,PostReportOriginalIndices:postIdx,
		PrefixStateOldRetention:prefixState,MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,
		MeanReportEffect:sr/5,MeanPostReportRehearsalRecovery:sa/5,MeanNetEpochChange:sn/5,
		FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP157B()(UP157BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP157BResult{Schema:UP157BPositionMapSchema,Experiment:"UP-157B-late-rehearsal-position-map",SourceUP156BSeal:"c1a5ffb93fcb4664012eeddbf891b157b69b7fbf",CommonPrefixEpochs:15,TerminalEpochs:5,RehearsalAfter:12,RehearsalBefore:3,Rotations:15,Subjects:2,OldUpdatesPerEpoch:15,NewUpdatesPerEpoch:24,ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false}
	for _,subject:=range []string{"mia","noah"}{for rotation:=0;rotation<15;rotation++{res.Arms=append(res.Arms,up157bRun(common,subject,rotation,o,r))}}
	return res,nil
}
