package unitary

const UP143BTrajectorySchema = "wingless.up143b-within-epoch-interference-trajectory.v1"

type UP143BEpochMetric struct {
	Family string `json:"family"`
	Arm string `json:"arm"`
	Epoch int `json:"epoch"`
	BeforeEpoch float64 `json:"before_epoch"`
	AfterNewPrefix float64 `json:"after_new_prefix"`
	AfterOldAnchor float64 `json:"after_old_anchor"`
	AfterNewTail float64 `json:"after_new_tail"`
	PrefixDamage float64 `json:"prefix_damage"`
	AnchorRecovery float64 `json:"anchor_recovery"`
	TailDamage float64 `json:"tail_damage"`
	NetEpochChange float64 `json:"net_epoch_change"`
}

type UP143BArmSummary struct {
	Family string `json:"family"`
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

type UP143BTrajectoryResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP142BSeal string `json:"source_up142b_seal"`
	StateDimension int `json:"state_dimension"`
	GateInputDimension int `json:"gate_input_dimension"`
	GroundingEpochs int `json:"grounding_epochs"`
	LearningRate float64 `json:"learning_rate"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ProjectorRecomputed bool `json:"projector_recomputed"`
	EpochMetrics []UP143BEpochMetric `json:"epoch_metrics"`
	ArmSummaries []UP143BArmSummary `json:"arm_summaries"`
}

func up143bOldMean(g *up129bGate,o,r [64]float64) float64 {
	h:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	u:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	return (h.ClassAccuracy+u.ClassAccuracy)/2
}

func up143bExamples(family string) []up135bExample {
	if family=="original" { return up135bNewExamples() }
	return up141bExamples()
}

func up143bShift(arm string,epoch int) int {
	if arm=="diversity_1" { return 0 }
	return epoch
}

func up143bNewEval(g *up129bGate,family string,names []string,o,r [64]float64) UP130BSplitMetric {
	if family=="original" {
		m,_:=up130bSplit(g,"new",names,o,r)
		return m
	}
	return up141bSplit(g,"new",names,o,r)
}

func up143bRun(family,arm string,o,r [64]float64)([]UP143BEpochMetric,UP143BArmSummary) {
	base:=up129bTrainGate(o,r)
	g:=base
	examples:=up143bExamples(family)
	rows:=make([]UP143BEpochMetric,0,20)
	sumPrefix,sumAnchor,sumTail,sumNet:=0.0,0.0,0.0,0.0

	for epoch:=0;epoch<20;epoch++ {
		before:=up143bOldMean(g,o,r)
		full:=up139bRotate(examples,up143bShift(arm,epoch))
		for i:=0;i<20;i++ { up135bStep(g,full[i],o,r) }
		afterPrefix:=up143bOldMean(g,o,r)
		up133bTrainOldBlock(g,epoch,o,r)
		afterAnchor:=up143bOldMean(g,o,r)
		for i:=20;i<24;i++ { up135bStep(g,full[i],o,r) }
		afterTail:=up143bOldMean(g,o,r)

		row:=UP143BEpochMetric{
			Family:family,Arm:arm,Epoch:epoch+1,
			BeforeEpoch:before,AfterNewPrefix:afterPrefix,AfterOldAnchor:afterAnchor,AfterNewTail:afterTail,
			PrefixDamage:afterPrefix-before,AnchorRecovery:afterAnchor-afterPrefix,
			TailDamage:afterTail-afterAnchor,NetEpochChange:afterTail-before,
		}
		rows=append(rows,row)
		sumPrefix+=row.PrefixDamage;sumAnchor+=row.AnchorRecovery;sumTail+=row.TailDamage;sumNet+=row.NetEpochChange
	}

	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	primary:=up143bNewEval(g,family,up130bNewNames[4:6],o,r)
	secondary:=up143bNewEval(g,family,up131bSecondaryNames(),o,r)
	s:=UP143BArmSummary{
		Family:family,Arm:arm,
		MeanPrefixDamage:sumPrefix/20,MeanAnchorRecovery:sumAnchor/20,
		MeanTailDamage:sumTail/20,MeanNetEpochChange:sumNet/20,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,
		FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:primary,FinalSecondaryNew:secondary,
		FinalMeanNewAccuracy:(primary.OverallAccuracy+secondary.OverallAccuracy)/2,
	}
	return rows,s
}

func RunUP143B()(UP143BTrajectoryResult,error) {
	o,r:=up124bCompetitorDirections()
	res:=UP143BTrajectoryResult{
		Schema:UP143BTrajectorySchema,Experiment:"UP-143B-within-epoch-interference-trajectory",
		SourceUP142BSeal:"110d063dee8ecbdadb5996de9afd1fa295f0d2d3",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false,ProjectorRecomputed:false,
	}
	for _,family:=range []string{"original","replication"} {
		for _,arm:=range []string{"diversity_1","diversity_20"} {
			rows,s:=up143bRun(family,arm,o,r)
			res.EpochMetrics=append(res.EpochMetrics,rows...)
			res.ArmSummaries=append(res.ArmSummaries,s)
		}
	}
	return res,nil
}
