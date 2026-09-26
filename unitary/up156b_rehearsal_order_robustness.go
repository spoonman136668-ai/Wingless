package unitary

const UP156BOrderSchema="wingless.up156b-rehearsal-order-robustness.v1"

type UP156BArm struct{
	Subject string `json:"subject"`
	RehearsalOrder string `json:"rehearsal_order"`
	RehearsalAfter int `json:"rehearsal_after"`
	RehearsalBefore int `json:"rehearsal_before"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportRehearsalRecovery float64 `json:"mean_post_report_rehearsal_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP156BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP155BSeal string `json:"source_up155b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	ReportBlockUpdates int `json:"report_block_updates"`
	Subjects int `json:"subjects"`
	Doses int `json:"doses"`
	Orders int `json:"orders"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptivePlacementUsed bool `json:"adaptive_placement_used"`
	Arms []UP156BArm `json:"arms"`
}

type up156bOldItem struct{surface up132bSurfaceSpec;subjectIndex int;originalIndex int}

func up156bOldItems(epoch int)[]up156bOldItem{
	old:=up132bOldSurfaces()
	out:=make([]up156bOldItem,0,len(old))
	for i,s:=range old{out=append(out,up156bOldItem{surface:s,subjectIndex:(i+epoch)%2,originalIndex:i})}
	return out
}
func up156bOrdered(items []up156bOldItem,order string,epoch int)[]up156bOldItem{
	n:=len(items);out:=make([]up156bOldItem,0,n)
	switch order{
	case "canonical":
		out=append(out,items...)
	case "reverse":
		for i:=n-1;i>=0;i--{out=append(out,items[i])}
	default:
		start:=((epoch-15)*3)%n
		for i:=0;i<n;i++{out=append(out,items[(start+i)%n])}
	}
	return out
}
func up156bApplyOld(g *up129bGate,items []up156bOldItem,start,end int,o,r [64]float64){
	for i:=start;i<end;i++{
		it:=items[i]
		up129bGateStep(g,up121bOriginalNames[it.subjectIndex],it.surface.verb,it.surface.class,o,r)
	}
}
func up156bEpoch(g *up129bGate,subject,order string,afterN,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){if !up153bInBlock(e,reportBlock){up135bStep(g,e,o,r)}}
	afterPrefix:=up143bOldMean(g,o,r)
	items:=up156bOrdered(up156bOldItems(epoch),order,epoch)
	beforeN:=15-afterN
	if beforeN>0{up156bApplyOld(g,items,0,beforeN,o,r)}
	afterPre:=up143bOldMean(g,o,r)
	for _,e:=range reportBlock{up135bStep(g,e,o,r)}
	afterReport:=up143bOldMean(g,o,r)
	if afterN>0{up156bApplyOld(g,items,beforeN,15,o,r)}
	afterAll:=up143bOldMean(g,o,r)
	return afterPrefix-before,afterPre-afterPrefix,afterReport-afterPre,afterAll-afterReport,afterAll-before
}
func up156bRun(common *up129bGate,subject,order string,afterN int,o,r [64]float64)UP156BArm{
	g:=*common
	prefixState:=up143bOldMean(&g,o,r)
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up156bEpoch(&g,subject,order,afterN,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP156BArm{Subject:subject,RehearsalOrder:order,RehearsalAfter:afterN,RehearsalBefore:15-afterN,PrefixStateOldRetention:prefixState,MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,MeanPostReportRehearsalRecovery:sa/5,MeanNetEpochChange:sn/5,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2}
}
func RunUP156B()(UP156BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP156BResult{Schema:UP156BOrderSchema,Experiment:"UP-156B-rehearsal-order-robustness",SourceUP155BSeal:"78f13eaa8774efe6ee17cc8343d4e19e8bff025e",CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,ReportBlockUpdates:4,Subjects:2,Doses:4,Orders:3,ExtraUpdatesUsed:false,AdaptivePlacementUsed:false}
	for _,subject:=range []string{"mia","noah"}{
		for _,order:=range []string{"canonical","reverse","rotating"}{
			for _,afterN:=range []int{0,8,12,15}{res.Arms=append(res.Arms,up156bRun(common,subject,order,afterN,o,r))}
		}
	}
	return res,nil
}
