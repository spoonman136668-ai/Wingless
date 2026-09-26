package unitary

const UP158BCleanupOrderSchema="wingless.up158b-cleanup-class-order.v1"

type UP158BArm struct{
	Subject string `json:"subject"`
	CleanupOrder string `json:"cleanup_order"`
	PreReportOriginalIndices []int `json:"pre_report_original_indices"`
	PostReportOriginalIndices []int `json:"post_report_original_indices"`
	MeanNewPrefixEffect float64 `json:"mean_new_prefix_effect"`
	MeanPreReportRehearsalRecovery float64 `json:"mean_pre_report_rehearsal_recovery"`
	MeanReportEffect float64 `json:"mean_report_effect"`
	MeanPostReportCleanupRecovery float64 `json:"mean_post_report_cleanup_recovery"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP158BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP157BSeal string `json:"source_up157b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	RehearsalBefore int `json:"rehearsal_before"`
	RehearsalAfter int `json:"rehearsal_after"`
	CleanupStore int `json:"cleanup_store"`
	CleanupObserve int `json:"cleanup_observe"`
	CleanupReport int `json:"cleanup_report"`
	Orders int `json:"orders"`
	Subjects int `json:"subjects"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Arms []UP158BArm `json:"arms"`
}
func up158bIndicesForOrder(order string)[]int{
	store:=[]int{0,1,2,3,4}
	observe:=[]int{5,6,7,8,9}
	report:=[]int{13,14}
	var groups [][]int
	switch order{
	case "STORE_OBSERVE_REPORT":groups=[][]int{store,observe,report}
	case "STORE_REPORT_OBSERVE":groups=[][]int{store,report,observe}
	case "OBSERVE_STORE_REPORT":groups=[][]int{observe,store,report}
	case "OBSERVE_REPORT_STORE":groups=[][]int{observe,report,store}
	case "REPORT_STORE_OBSERVE":groups=[][]int{report,store,observe}
	default:groups=[][]int{report,observe,store}
	}
	out:=make([]int,0,12)
	for _,g:=range groups{out=append(out,g...)}
	return out
}
func up158bItems(epoch int,indices []int)[]up156bOldItem{
	all:=up156bOldItems(epoch);out:=make([]up156bOldItem,0,len(indices))
	for _,i:=range indices{out=append(out,all[i])}
	return out
}
func up158bEpoch(g *up129bGate,subject,order string,epoch int,o,r [64]float64)(prefix,pre,report,post,net float64){
	reportBlock:=up154bReportBlock(subject)
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){if !up153bInBlock(e,reportBlock){up135bStep(g,e,o,r)}}
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
func up158bRun(common *up129bGate,subject,order string,o,r [64]float64)UP158BArm{
	g:=*common
	sp,sb,sr,sa,sn:=0.0,0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,b,rr,a,n:=up158bEpoch(&g,subject,order,epoch,o,r)
		sp+=p;sb+=b;sr+=rr;sa+=a;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP158BArm{
		Subject:subject,CleanupOrder:order,PreReportOriginalIndices:[]int{10,11,12},
		PostReportOriginalIndices:append([]int(nil),up158bIndicesForOrder(order)...),
		MeanNewPrefixEffect:sp/5,MeanPreReportRehearsalRecovery:sb/5,MeanReportEffect:sr/5,
		MeanPostReportCleanupRecovery:sa/5,MeanNetEpochChange:sn/5,
		FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP158B()(UP158BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP158BResult{Schema:UP158BCleanupOrderSchema,Experiment:"UP-158B-cleanup-class-order",SourceUP157BSeal:"b969ff7c6b4c7ac8ba3eb83b073216780cabc392",CommonPrefixEpochs:15,TerminalEpochs:5,RehearsalBefore:3,RehearsalAfter:12,CleanupStore:5,CleanupObserve:5,CleanupReport:2,Orders:6,Subjects:2,ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false}
	orders:=[]string{"STORE_OBSERVE_REPORT","STORE_REPORT_OBSERVE","OBSERVE_STORE_REPORT","OBSERVE_REPORT_STORE","REPORT_STORE_OBSERVE","REPORT_OBSERVE_STORE"}
	for _,subject:=range []string{"mia","noah"}{for _,order:=range orders{res.Arms=append(res.Arms,up158bRun(common,subject,order,o,r))}}
	return res,nil
}
