package unitary

const UP153BAnchorSwapSchema="wingless.up153b-anchor-side-swap.v1"

type UP153BArm struct{
	Subject string `json:"subject"`
	Arm string `json:"arm"`
	PostAnchorBlock string `json:"post_anchor_block"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanPreAnchorEffect float64 `json:"mean_pre_anchor_effect"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanPostAnchorEffect float64 `json:"mean_post_anchor_effect"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP153BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP152BSeal string `json:"source_up152b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	PostAnchorUpdates int `json:"post_anchor_updates"`
	Subjects int `json:"subjects"`
	ArmsPerSubject int `json:"arms_per_subject"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP153BArm `json:"arms"`
}

func up153bBlock(subject string,verbs []string)[]up135bExample{
	out:=make([]up135bExample,0,len(verbs))
	for _,verb:=range verbs{out=append(out,up135bExample{name:subject,verb:verb,class:up152bVerbClass(verb)})}
	return out
}
func up153bInBlock(e up135bExample,block []up135bExample)bool{
	for _,x:=range block{if up151bSame(e,x){return true}}
	return false
}
func up153bEpoch(g *up129bGate,subject,arm string,epoch int,o,r [64]float64)(pre,anchor,post,net float64){
	report:=up153bBlock(subject,[]string{"relays","announces","cites","summarizes"})
	mixed:=up153bBlock(subject,[]string{"caches","files","scans","checks"})
	beforeBlock,afterBlock:=mixed,report
	if arm=="report_before_anchor"{beforeBlock,afterBlock=report,mixed}
	before:=up143bOldMean(g,o,r)
	for _,e:=range up135bNewExamples(){
		if up153bInBlock(e,report)||up153bInBlock(e,mixed){continue}
		up135bStep(g,e,o,r)
	}
	for _,e:=range beforeBlock{up135bStep(g,e,o,r)}
	preState:=up143bOldMean(g,o,r)
	up133bTrainOldBlock(g,epoch,o,r)
	anchorState:=up143bOldMean(g,o,r)
	for _,e:=range afterBlock{up135bStep(g,e,o,r)}
	after:=up143bOldMean(g,o,r)
	return preState-before,anchorState-preState,after-anchorState,after-before
}
func up153bRun(common *up129bGate,subject,arm string,o,r [64]float64)UP153BArm{
	g:=*common
	prefixState:=up143bOldMean(&g,o,r)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,a,t,n:=up153bEpoch(&g,subject,arm,epoch,o,r)
		sp+=p;sa+=a;st+=t;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	post:="report"
	if arm=="report_before_anchor"{post="store_observe"}
	return UP153BArm{
		Subject:subject,Arm:arm,PostAnchorBlock:post,PrefixStateOldRetention:prefixState,
		MeanPreAnchorEffect:sp/5,MeanAnchorRecovery:sa/5,MeanPostAnchorEffect:st/5,MeanNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:pn,FinalSecondaryNew:snM,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP153B()(UP153BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP153BResult{
		Schema:UP153BAnchorSwapSchema,Experiment:"UP-153B-anchor-side-swap",
		SourceUP152BSeal:"5f7ca686414f67cd56f0e126f87df9c90188e2c7",
		CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,PostAnchorUpdates:4,
		Subjects:2,ArmsPerSubject:2,AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false,
	}
	for _,subject:=range []string{"mia","noah"}{
		for _,arm:=range []string{"report_after_anchor","report_before_anchor"}{
			res.Arms=append(res.Arms,up153bRun(common,subject,arm,o,r))
		}
	}
	return res,nil
}
