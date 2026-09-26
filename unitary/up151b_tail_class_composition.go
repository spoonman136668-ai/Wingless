package unitary

const UP151BCompositionSchema="wingless.up151b-tail-class-composition.v1"

type UP151BArm struct{
	Subject string `json:"subject"`
	Composition string `json:"composition"`
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
type UP151BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP150BSeal string `json:"source_up150b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	Subjects int `json:"subjects"`
	Compositions int `json:"compositions"`
	AdaptiveTailSelection bool `json:"adaptive_tail_selection"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP151BArm `json:"arms"`
}

func up151bTail(subject,composition string)[]up135bExample{
	var verbs []string
	switch composition{
	case "report_all": verbs=[]string{"relays","announces","cites","summarizes"}
	case "store_all": verbs=[]string{"lodges","stashes","caches","files"}
	case "observe_all": verbs=[]string{"scans","checks","views","monitors"}
	case "mixed_seen": verbs=[]string{"caches","files","scans","checks"}
	default: verbs=[]string{"lodges","stashes","views","monitors"}
	}
	out:=make([]up135bExample,0,4)
	for _,verb:=range verbs{
		class:=up97bReport
		for _,v:=range up130bStore{if v==verb{class=up97bStore}}
		for _,v:=range up130bObserve{if v==verb{class=up97bObserve}}
		out=append(out,up135bExample{name:subject,verb:verb,class:class})
	}
	return out
}
func up151bSame(a,b up135bExample)bool{return a.name==b.name&&a.verb==b.verb}
func up151bEpoch(g *up129bGate,tail []up135bExample,epoch int,o,r [64]float64)(p,a,t,n float64){
	before:=up143bOldMean(g,o,r)
	all:=up135bNewExamples()
	for _,e:=range all{
		inTail:=false
		for _,x:=range tail{if up151bSame(e,x){inTail=true;break}}
		if !inTail{up135bStep(g,e,o,r)}
	}
	prefix:=up143bOldMean(g,o,r)
	up133bTrainOldBlock(g,epoch,o,r)
	anchor:=up143bOldMean(g,o,r)
	for _,e:=range tail{up135bStep(g,e,o,r)}
	after:=up143bOldMean(g,o,r)
	return prefix-before,anchor-prefix,after-anchor,after-before
}
func up151bRun(common *up129bGate,subject,composition string,o,r [64]float64)UP151BArm{
	g:=*common
	tail:=up151bTail(subject,composition)
	prefixState:=up143bOldMean(&g,o,r)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for epoch:=15;epoch<20;epoch++{
		p,a,t,n:=up151bEpoch(&g,tail,epoch,o,r);sp+=p;sa+=a;st+=t;sn+=n
	}
	oldH:=up129bSplit(&g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(&g,"old_unseen",up121bUnseenNames,o,r)
	pn,_:=up130bSplit(&g,"primary_new",up130bNewNames[4:6],o,r)
	snM,_:=up130bSplit(&g,"secondary_new",up131bSecondaryNames(),o,r)
	verbs:=make([]string,0,4);for _,e:=range tail{verbs=append(verbs,e.verb)}
	return UP151BArm{
		Subject:subject,Composition:composition,TailVerbs:verbs,PrefixStateOldRetention:prefixState,
		MeanTerminalPrefixDamage:sp/5,MeanTerminalAnchorRecovery:sa/5,MeanTerminalTailDamage:st/5,MeanTerminalNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:pn,FinalSecondaryNew:snM,FinalMeanNewAccuracy:(pn.OverallAccuracy+snM.OverallAccuracy)/2,
	}
}
func RunUP151B()(UP151BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP151BResult{
		Schema:UP151BCompositionSchema,Experiment:"UP-151B-tail-class-composition",
		SourceUP150BSeal:"fffb8ea47a2cd58cf40c8eebb63b226cdadaa015",
		CommonPrefixEpochs:15,TerminalEpochs:5,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		Subjects:2,Compositions:5,AdaptiveTailSelection:false,ExtraUpdatesUsed:false,
	}
	for _,subject:=range []string{"mia","noah"}{
		for _,composition:=range []string{"report_all","store_all","observe_all","mixed_seen","mixed_alternate"}{
			res.Arms=append(res.Arms,up151bRun(common,subject,composition,o,r))
		}
	}
	return res,nil
}
