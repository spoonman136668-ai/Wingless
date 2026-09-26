package unitary

const UP144BSpacingSchema = "wingless.up144b-tail-recurrence-spacing.v1"

type UP144BEpochMetric struct {
	Arm string `json:"arm"`
	Epoch int `json:"epoch"`
	Shift int `json:"shift"`
	BeforeEpoch float64 `json:"before_epoch"`
	AfterNewPrefix float64 `json:"after_new_prefix"`
	AfterOldAnchor float64 `json:"after_old_anchor"`
	AfterNewTail float64 `json:"after_new_tail"`
	PrefixDamage float64 `json:"prefix_damage"`
	AnchorRecovery float64 `json:"anchor_recovery"`
	TailDamage float64 `json:"tail_damage"`
	NetEpochChange float64 `json:"net_epoch_change"`
}
type UP144BArmSummary struct {
	Arm string `json:"arm"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP144BSpacingResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP143BSeal string `json:"source_up143b_seal"`
	Epochs int `json:"epochs"`
	LearningRate float64 `json:"learning_rate"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	UsesPerTailSet int `json:"uses_per_tail_set"`
	AdaptiveShiftChoice bool `json:"adaptive_shift_choice"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	EpochMetrics []UP144BEpochMetric `json:"epoch_metrics"`
	ArmSummaries []UP144BArmSummary `json:"arm_summaries"`
}

func up144bShift(arm string,epoch int) int {
	shifts:=[4]int{0,6,12,18}
	if arm=="interleaved" { return shifts[epoch%4] }
	return shifts[epoch/5]
}
func up144bRun(arm string,o,r [64]float64)([]UP144BEpochMetric,UP144BArmSummary) {
	g:=up129bTrainGate(o,r)
	examples:=up135bNewExamples()
	rows:=make([]UP144BEpochMetric,0,20)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for epoch:=0;epoch<20;epoch++ {
		before:=up143bOldMean(g,o,r)
		shift:=up144bShift(arm,epoch)
		full:=up139bRotate(examples,shift)
		for i:=0;i<20;i++ { up135bStep(g,full[i],o,r) }
		prefix:=up143bOldMean(g,o,r)
		up133bTrainOldBlock(g,epoch,o,r)
		anchor:=up143bOldMean(g,o,r)
		for i:=20;i<24;i++ { up135bStep(g,full[i],o,r) }
		tail:=up143bOldMean(g,o,r)
		m:=UP144BEpochMetric{
			Arm:arm,Epoch:epoch+1,Shift:shift,BeforeEpoch:before,AfterNewPrefix:prefix,
			AfterOldAnchor:anchor,AfterNewTail:tail,PrefixDamage:prefix-before,
			AnchorRecovery:anchor-prefix,TailDamage:tail-anchor,NetEpochChange:tail-before,
		}
		rows=append(rows,m);sp+=m.PrefixDamage;sa+=m.AnchorRecovery;st+=m.TailDamage;sn+=m.NetEpochChange
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return rows,UP144BArmSummary{
		Arm:arm,MeanPrefixDamage:sp/20,MeanAnchorRecovery:sa/20,MeanTailDamage:st/20,MeanNetEpochChange:sn/20,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:p,FinalSecondaryNew:s,FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}
func RunUP144B()(UP144BSpacingResult,error) {
	o,r:=up124bCompetitorDirections()
	res:=UP144BSpacingResult{
		Schema:UP144BSpacingSchema,Experiment:"UP-144B-tail-recurrence-spacing",
		SourceUP143BSeal:"e52125cd6b57be5564ec797327ff5819d822cb04",
		Epochs:20,LearningRate:0.08,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		DistinctTailSets:4,UsesPerTailSet:5,AdaptiveShiftChoice:false,ExtraUpdatesUsed:false,
	}
	for _,arm:=range []string{"interleaved","blocked"} {
		rows,s:=up144bRun(arm,o,r);res.EpochMetrics=append(res.EpochMetrics,rows...);res.ArmSummaries=append(res.ArmSummaries,s)
	}
	return res,nil
}
