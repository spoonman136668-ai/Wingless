package unitary

const UP142BInteractionSchema="wingless.up142b-family-diversity-interaction.v1"

type UP142BCell struct {
	Family string `json:"family"`
	Arm string `json:"arm"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	PrimaryNew UP130BSplitMetric `json:"primary_new"`
	SecondaryNew UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained UP129BSplitMetric `json:"old_unseen_retained"`
	MeanOldRetention float64 `json:"mean_old_retention"`
	MeanNewAccuracy float64 `json:"mean_new_accuracy"`
}
type UP142BInteractionResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP141BSeal string `json:"source_up141b_seal"`
	StateDimension int `json:"state_dimension"`
	GateInputDimension int `json:"gate_input_dimension"`
	GroundingEpochs int `json:"grounding_epochs"`
	LearningRate float64 `json:"learning_rate"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveFamilyChoice bool `json:"adaptive_family_choice"`
	ProjectorRecomputed bool `json:"projector_recomputed"`
	Cells []UP142BCell `json:"cells"`
	OriginalOldRetentionDelta20Vs1 float64 `json:"original_old_retention_delta20_vs_1"`
	ReplicationOldRetentionDelta20Vs1 float64 `json:"replication_old_retention_delta20_vs_1"`
	OriginalNewAccuracyDelta20Vs1 float64 `json:"original_new_accuracy_delta20_vs_1"`
	ReplicationNewAccuracyDelta20Vs1 float64 `json:"replication_new_accuracy_delta20_vs_1"`
	OldRetentionInteractionDelta float64 `json:"old_retention_interaction_delta"`
}

func up142bCell(family,arm string,o,r [64]float64) UP142BCell {
	base:=up129bTrainGate(o,r)
	g:=base
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	var p,s UP130BSplitMetric
	if family=="original" {
		up140bTrain(g,arm,o,r)
		p,_=up130bSplit(g,"primary_new",primary,o,r)
		s,_=up130bSplit(g,"secondary_new",secondary,o,r)
	}else{
		up141bTrain(g,arm,o,r)
		p=up141bSplit(g,"primary_new",primary,o,r)
		s=up141bSplit(g,"secondary_new",secondary,o,r)
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	return UP142BCell{
		Family:family,Arm:arm,DistinctTailSets:len(up141bShifts(arm)),
		PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		MeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		MeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}

func RunUP142B()(UP142BInteractionResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP142BInteractionResult{
		Schema:UP142BInteractionSchema,Experiment:"UP-142B-family-diversity-interaction",
		SourceUP141BSeal:"7040fd48fd135ab0d0e68966582b877f8d8d79f1",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		ExtraUpdatesUsed:false,AdaptiveFamilyChoice:false,ProjectorRecomputed:false,
	}
	var o1,o20,r1,r20 UP142BCell
	for _,fam:=range []string{"original","replication"}{
		for _,arm:=range []string{"diversity_1","diversity_8","diversity_20"}{
			c:=up142bCell(fam,arm,o,r)
			res.Cells=append(res.Cells,c)
			if fam=="original"&&arm=="diversity_1"{o1=c}
			if fam=="original"&&arm=="diversity_20"{o20=c}
			if fam=="replication"&&arm=="diversity_1"{r1=c}
			if fam=="replication"&&arm=="diversity_20"{r20=c}
		}
	}
	res.OriginalOldRetentionDelta20Vs1=o20.MeanOldRetention-o1.MeanOldRetention
	res.ReplicationOldRetentionDelta20Vs1=r20.MeanOldRetention-r1.MeanOldRetention
	res.OriginalNewAccuracyDelta20Vs1=o20.MeanNewAccuracy-o1.MeanNewAccuracy
	res.ReplicationNewAccuracyDelta20Vs1=r20.MeanNewAccuracy-r1.MeanNewAccuracy
	res.OldRetentionInteractionDelta=res.OriginalOldRetentionDelta20Vs1-res.ReplicationOldRetentionDelta20Vs1
	return res,nil
}
