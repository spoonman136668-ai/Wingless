package unitary

const UP145BSwitchSchema="wingless.up145b-tail-switch-frequency.v1"

type UP145BSummary struct{
	Arm string `json:"arm"`
	Switches int `json:"switches"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP145BSwitchResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP144BSeal string `json:"source_up144b_seal"`
	Epochs int `json:"epochs"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	UsesPerTailSet int `json:"uses_per_tail_set"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Summaries []UP145BSummary `json:"summaries"`
}
func up145bSchedule(arm string) []int{
	switch arm{
	case "low_switch": return []int{0,0,0,0,0,6,6,6,6,6,12,12,12,12,12,18,18,18,18,18}
	case "medium_switch": return []int{0,0,6,6,12,12,18,18,0,0,6,6,12,12,18,18,0,6,12,18}
	default: return []int{0,6,12,18,0,6,12,18,0,6,12,18,0,6,12,18,0,6,12,18}
	}
}
func up145bSwitchCount(s []int) int{n:=0;for i:=1;i<len(s);i++{if s[i]!=s[i-1]{n++}};return n}
func up145bRun(arm string,o,r [64]float64) UP145BSummary{
	g:=up129bTrainGate(o,r);examples:=up135bNewExamples();sched:=up145bSchedule(arm)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for epoch,shift:=range sched{
		before:=up143bOldMean(g,o,r);full:=up139bRotate(examples,shift)
		for i:=0;i<20;i++{up135bStep(g,full[i],o,r)}
		prefix:=up143bOldMean(g,o,r);up133bTrainOldBlock(g,epoch,o,r);anchor:=up143bOldMean(g,o,r)
		for i:=20;i<24;i++{up135bStep(g,full[i],o,r)}
		tail:=up143bOldMean(g,o,r)
		sp+=prefix-before;sa+=anchor-prefix;st+=tail-anchor;sn+=tail-before
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP145BSummary{
		Arm:arm,Switches:up145bSwitchCount(sched),MeanPrefixDamage:sp/20,MeanAnchorRecovery:sa/20,
		MeanTailDamage:st/20,MeanNetEpochChange:sn/20,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}
func RunUP145B()(UP145BSwitchResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP145BSwitchResult{Schema:UP145BSwitchSchema,Experiment:"UP-145B-tail-switch-frequency",SourceUP144BSeal:"bdd4b4e30dcee2403d611e50f5da3b3af061df39",Epochs:20,DistinctTailSets:4,UsesPerTailSet:5,AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false}
	for _,arm:=range []string{"low_switch","medium_switch","high_switch"}{res.Summaries=append(res.Summaries,up145bRun(arm,o,r))}
	return res,nil
}
