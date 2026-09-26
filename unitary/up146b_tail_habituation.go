package unitary

const UP146BHabituationSchema="wingless.up146b-tail-habituation.v1"

type UP146BEpochMetric struct{
	Epoch int `json:"epoch"`
	Shift int `json:"shift"`
	Repetition int `json:"repetition"`
	PrefixDamage float64 `json:"prefix_damage"`
	AnchorRecovery float64 `json:"anchor_recovery"`
	TailDamage float64 `json:"tail_damage"`
	NetEpochChange float64 `json:"net_epoch_change"`
}
type UP146BPositionSummary struct{
	Repetition int `json:"repetition"`
	Samples int `json:"samples"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
}
type UP146BHabituationResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP145BSeal string `json:"source_up145b_seal"`
	Epochs int `json:"epochs"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	UsesPerTailSet int `json:"uses_per_tail_set"`
	Switches int `json:"switches"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	EpochMetrics []UP146BEpochMetric `json:"epoch_metrics"`
	PositionSummaries []UP146BPositionSummary `json:"position_summaries"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}

func RunUP146B()(UP146BHabituationResult,error){
	o,r:=up124bCompetitorDirections()
	g:=up129bTrainGate(o,r)
	examples:=up135bNewExamples()
	sched:=up145bSchedule("low_switch")
	res:=UP146BHabituationResult{
		Schema:UP146BHabituationSchema,Experiment:"UP-146B-tail-habituation",
		SourceUP145BSeal:"4bf88b261a8179a6c484fc190e445c2892d210fc",
		Epochs:20,DistinctTailSets:4,UsesPerTailSet:5,Switches:up145bSwitchCount(sched),
		AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false,
	}
	var sp,sa,st,sn [5]float64
	var counts [5]int
	for epoch,shift:=range sched{
		before:=up143bOldMean(g,o,r)
		full:=up139bRotate(examples,shift)
		for i:=0;i<20;i++{up135bStep(g,full[i],o,r)}
		prefix:=up143bOldMean(g,o,r)
		up133bTrainOldBlock(g,epoch,o,r)
		anchor:=up143bOldMean(g,o,r)
		for i:=20;i<24;i++{up135bStep(g,full[i],o,r)}
		tail:=up143bOldMean(g,o,r)
		rep:=epoch%5
		m:=UP146BEpochMetric{
			Epoch:epoch+1,Shift:shift,Repetition:rep+1,
			PrefixDamage:prefix-before,AnchorRecovery:anchor-prefix,
			TailDamage:tail-anchor,NetEpochChange:tail-before,
		}
		res.EpochMetrics=append(res.EpochMetrics,m)
		sp[rep]+=m.PrefixDamage;sa[rep]+=m.AnchorRecovery;st[rep]+=m.TailDamage;sn[rep]+=m.NetEpochChange;counts[rep]++
	}
	for i:=0;i<5;i++{
		n:=float64(counts[i])
		res.PositionSummaries=append(res.PositionSummaries,UP146BPositionSummary{
			Repetition:i+1,Samples:counts[i],MeanPrefixDamage:sp[i]/n,MeanAnchorRecovery:sa[i]/n,
			MeanTailDamage:st[i]/n,MeanNetEpochChange:sn[i]/n,
		})
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	res.FinalMeanOldRetention=(oldH.ClassAccuracy+oldU.ClassAccuracy)/2
	res.FinalMeanNewAccuracy=(p.OverallAccuracy+s.OverallAccuracy)/2
	return res,nil
}
