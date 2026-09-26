package unitary

const UP147BBlockOrderSchema="wingless.up147b-tail-block-order.v1"

type UP147BBlockMetric struct{
	Arm string `json:"arm"`
	BlockPosition int `json:"block_position"`
	Shift int `json:"shift"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
}
type UP147BArmSummary struct{
	Arm string `json:"arm"`
	Switches int `json:"switches"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP147BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP146BSeal string `json:"source_up146b_seal"`
	Epochs int `json:"epochs"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	UsesPerTailSet int `json:"uses_per_tail_set"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	BlockMetrics []UP147BBlockMetric `json:"block_metrics"`
	ArmSummaries []UP147BArmSummary `json:"arm_summaries"`
}

func up147bOrder(arm string) []int{
	switch arm{
	case "order_0": return []int{0,6,12,18}
	case "order_1": return []int{6,12,18,0}
	case "order_2": return []int{12,18,0,6}
	default: return []int{18,0,6,12}
	}
}
func up147bSchedule(order []int) []int{
	out:=make([]int,0,20)
	for _,shift:=range order{for i:=0;i<5;i++{out=append(out,shift)}}
	return out
}
func up147bRun(arm string,o,r [64]float64)([]UP147BBlockMetric,UP147BArmSummary){
	g:=up129bTrainGate(o,r)
	examples:=up135bNewExamples()
	order:=up147bOrder(arm)
	blocks:=make([]UP147BBlockMetric,0,4)
	epoch:=0
	for blockPos,shift:=range order{
		sp,sa,st,sn:=0.0,0.0,0.0,0.0
		for rep:=0;rep<5;rep++{
			before:=up143bOldMean(g,o,r)
			full:=up139bRotate(examples,shift)
			for i:=0;i<20;i++{up135bStep(g,full[i],o,r)}
			prefix:=up143bOldMean(g,o,r)
			up133bTrainOldBlock(g,epoch,o,r)
			anchor:=up143bOldMean(g,o,r)
			for i:=20;i<24;i++{up135bStep(g,full[i],o,r)}
			tail:=up143bOldMean(g,o,r)
			sp+=prefix-before
			sa+=anchor-prefix
			st+=tail-anchor
			sn+=tail-before
			epoch++
		}
		blocks=append(blocks,UP147BBlockMetric{
			Arm:arm,BlockPosition:blockPos+1,Shift:shift,
			MeanPrefixDamage:sp/5,MeanAnchorRecovery:sa/5,
			MeanTailDamage:st/5,MeanNetEpochChange:sn/5,
		})
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return blocks,UP147BArmSummary{
		Arm:arm,Switches:up145bSwitchCount(up147bSchedule(order)),
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,
		FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:p,FinalSecondaryNew:s,
		FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}
func RunUP147B()(UP147BResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP147BResult{
		Schema:UP147BBlockOrderSchema,Experiment:"UP-147B-tail-block-order",
		SourceUP146BSeal:"fe75f2ca45f4566b71afa14e18b439812514d9e3",
		Epochs:20,DistinctTailSets:4,UsesPerTailSet:5,
		AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false,
	}
	for _,arm:=range []string{"order_0","order_1","order_2","order_3"}{
		blocks,s:=up147bRun(arm,o,r)
		res.BlockMetrics=append(res.BlockMetrics,blocks...)
		res.ArmSummaries=append(res.ArmSummaries,s)
	}
	return res,nil
}
