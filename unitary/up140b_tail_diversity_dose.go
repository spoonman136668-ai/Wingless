package unitary

const UP140BTailDiversitySchema = "wingless.up140b-tail-diversity-dose.v1"

type UP140BArmMetric struct {
	Arm                 string            `json:"arm"`
	DistinctTailSets    int               `json:"distinct_tail_sets"`
	PrimaryNew          UP130BSplitMetric `json:"primary_new"`
	SecondaryNew        UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained  UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained   UP129BSplitMetric `json:"old_unseen_retained"`
	MeanOldRetention    float64           `json:"mean_old_retention"`
	MeanNewAccuracy     float64           `json:"mean_new_accuracy"`
}

type UP140BTailDiversityResult struct {
	Schema              string            `json:"schema"`
	Experiment          string            `json:"experiment"`
	SourceUP139BSeal    string            `json:"source_up139b_seal"`
	StateDimension      int               `json:"state_dimension"`
	GateInputDimension  int               `json:"gate_input_dimension"`
	GroundingEpochs     int               `json:"grounding_epochs"`
	LearningRate        float64           `json:"learning_rate"`
	NewUpdatesPerEpoch  int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch  int               `json:"old_updates_per_epoch"`
	FinalNewUpdates     int               `json:"final_new_updates"`
	ExtraUpdatesUsed    bool              `json:"extra_updates_used"`
	AdaptiveShiftChoice bool              `json:"adaptive_shift_choice"`
	ProjectorRecomputed bool              `json:"projector_recomputed"`
	ArmMetrics          []UP140BArmMetric `json:"arm_metrics"`
}

func up140bShifts(arm string) []int {
	switch arm {
	case "diversity_1":
		return []int{0}
	case "diversity_2":
		return []int{0,12}
	case "diversity_4":
		return []int{0,6,12,18}
	case "diversity_8":
		return []int{0,3,6,9,12,15,18,21}
	case "diversity_20":
		out:=make([]int,20)
		for i:=0;i<20;i++ { out[i]=i }
		return out
	}
	return nil
}

func up140bTrain(g *up129bGate,arm string,o,r [64]float64) {
	base:=up135bNewExamples()
	shifts:=up140bShifts(arm)
	for epoch:=0;epoch<20;epoch++ {
		shift:=shifts[epoch%len(shifts)]
		full:=up139bRotate(base,shift)
		for i:=0;i<20;i++ { up135bStep(g,full[i],o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for i:=20;i<24;i++ { up135bStep(g,full[i],o,r) }
	}
}

func RunUP140B()(UP140BTailDiversityResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP140BTailDiversityResult{
		Schema:UP140BTailDiversitySchema,Experiment:"UP-140B-tail-diversity-dose",
		SourceUP139BSeal:"c9963417bf1ecd9ec4fdc3353dcb577801febd94",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		ExtraUpdatesUsed:false,AdaptiveShiftChoice:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range []string{"diversity_1","diversity_2","diversity_4","diversity_8","diversity_20"} {
		clone:=*base
		g:=&clone
		up140bTrain(g,arm,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP140BArmMetric{
			Arm:arm,DistinctTailSets:len(up140bShifts(arm)),
			PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
			MeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
			MeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
		})
	}
	return result,nil
}
