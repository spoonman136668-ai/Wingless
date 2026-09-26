package unitary

const UP148BCrossoverSchema="wingless.up148b-terminal-tail-crossover.v1"

type UP148BBlockMetric struct{
	Arm string `json:"arm"`
	BlockPosition int `json:"block_position"`
	Shift int `json:"shift"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
}
type UP148BArmSummary struct{
	Arm string `json:"arm"`
	Pair string `json:"pair"`
	FinalShift int `json:"final_shift"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP148BPairDelta struct{
	Pair string `json:"pair"`
	SafeArm string `json:"safe_arm"`
	DamagingArm string `json:"damaging_arm"`
	OldRetentionDelta float64 `json:"old_retention_delta"`
	NewAccuracyDelta float64 `json:"new_accuracy_delta"`
}
type UP148BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP147BSeal string `json:"source_up147b_seal"`
	Epochs int `json:"epochs"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	UsesPerTailSet int `json:"uses_per_tail_set"`
	Switches int `json:"switches"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	ProjectorRecomputed bool `json:"projector_recomputed"`
	BlockMetrics []UP148BBlockMetric `json:"block_metrics"`
	ArmSummaries []UP148BArmSummary `json:"arm_summaries"`
	PairDeltas []UP148BPairDelta `json:"pair_deltas"`
}
func up148bSchedule(arm string)[4]int{
	switch arm{
	case "safe_final_18":return [4]int{0,6,12,18}
	case "damaging_final_12":return [4]int{0,6,18,12}
	case "safe_final_6":return [4]int{12,18,0,6}
	default:return [4]int{12,18,6,0}
	}
}
func up148bPair(arm string)string{if arm=="safe_final_18"||arm=="damaging_final_12"{return "pair_a"};return "pair_b"}
func up148bRun(arm string,o,r [64]float64)([]UP148BBlockMetric,UP148BArmSummary){
	g:=up129bTrainGate(o,r);examples:=up135bNewExamples();sched:=up148bSchedule(arm);epoch:=0
	blocks:=make([]UP148BBlockMetric,0,4)
	for bp,shift:=range sched{
		sp,sa,st,sn:=0.0,0.0,0.0,0.0
		for rep:=0;rep<5;rep++{
			before:=up143bOldMean(g,o,r);full:=up139bRotate(examples,shift)
			for i:=0;i<20;i++{up135bStep(g,full[i],o,r)}
			prefix:=up143bOldMean(g,o,r);up133bTrainOldBlock(g,epoch,o,r);anchor:=up143bOldMean(g,o,r)
			for i:=20;i<24;i++{up135bStep(g,full[i],o,r)}
			tail:=up143bOldMean(g,o,r)
			sp+=prefix-before;sa+=anchor-prefix;st+=tail-anchor;sn+=tail-before;epoch++
		}
		blocks=append(blocks,UP148BBlockMetric{Arm:arm,BlockPosition:bp+1,Shift:shift,MeanPrefixDamage:sp/5,MeanAnchorRecovery:sa/5,MeanTailDamage:st/5,MeanNetEpochChange:sn/5})
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return blocks,UP148BArmSummary{Arm:arm,Pair:up148bPair(arm),FinalShift:sched[3],FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2}
}
func RunUP148B()(UP148BResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP148BResult{Schema:UP148BCrossoverSchema,Experiment:"UP-148B-terminal-tail-crossover",SourceUP147BSeal:"8a6c31cc4f3c3c439d8d297d953b847cfe8b53ae",Epochs:20,DistinctTailSets:4,UsesPerTailSet:5,Switches:3,AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false,ProjectorRecomputed:false}
	arms:=[]string{"safe_final_18","damaging_final_12","safe_final_6","damaging_final_0"}
	for _,arm:=range arms{b,s:=up148bRun(arm,o,r);res.BlockMetrics=append(res.BlockMetrics,b...);res.ArmSummaries=append(res.ArmSummaries,s)}
	find:=func(name string)UP148BArmSummary{for _,s:=range res.ArmSummaries{if s.Arm==name{return s}};return UP148BArmSummary{}}
	for _,p:=range []struct{pair,safe,bad string}{{"pair_a","safe_final_18","damaging_final_12"},{"pair_b","safe_final_6","damaging_final_0"}}{
		a,b:=find(p.safe),find(p.bad);res.PairDeltas=append(res.PairDeltas,UP148BPairDelta{Pair:p.pair,SafeArm:p.safe,DamagingArm:p.bad,OldRetentionDelta:a.FinalMeanOldRetention-b.FinalMeanOldRetention,NewAccuracyDelta:a.FinalMeanNewAccuracy-b.FinalMeanNewAccuracy})
	}
	return res,nil
}
