package unitary

const UP136BMatchedRefreshSchema = "wingless.up136b-matched-refresh-confirmation.v1"

type UP136BArmMetric struct {
	Arm                string            `json:"arm"`
	FinalNewUpdates    int               `json:"final_new_updates"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP136BMatchedRefreshResult struct {
	Schema              string            `json:"schema"`
	Experiment          string            `json:"experiment"`
	SourceUP135BSeal    string            `json:"source_up135b_seal"`
	StateDimension      int               `json:"state_dimension"`
	GateInputDimension  int               `json:"gate_input_dimension"`
	GroundingEpochs     int               `json:"grounding_epochs"`
	LearningRate        float64           `json:"learning_rate"`
	NewUpdatesPerEpoch  int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch  int               `json:"old_updates_per_epoch"`
	RotationUsed        bool              `json:"rotation_used"`
	ExtraUpdatesUsed    bool              `json:"extra_updates_used"`
	ProjectorRecomputed bool              `json:"projector_recomputed"`
	ArmMetrics          []UP136BArmMetric `json:"arm_metrics"`
}

func up136bFixedNewExamples() []up135bExample {
	return up135bNewExamples()
}

func up136bTrain(g *up129bGate,refresh bool,o,r [64]float64) {
	newExamples:=up136bFixedNewExamples()
	for epoch:=0;epoch<20;epoch++ {
		if !refresh {
			for _,e:=range newExamples { up135bStep(g,e,o,r) }
			up133bTrainOldBlock(g,epoch,o,r)
			continue
		}
		for i:=0;i<20;i++ { up135bStep(g,newExamples[i],o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for i:=20;i<24;i++ { up135bStep(g,newExamples[i],o,r) }
	}
}

func RunUP136B()(UP136BMatchedRefreshResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP136BMatchedRefreshResult{
		Schema:UP136BMatchedRefreshSchema,Experiment:"UP-136B-matched-refresh-confirmation",
		SourceUP135BSeal:"6b0d70a5552403c4af843bbc5fe59d59829c3c67",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,RotationUsed:false,ExtraUpdatesUsed:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range []struct{name string;refresh bool;count int}{
		{"exact_terminal15_control",false,0},
		{"fixed_tail4_refresh",true,4},
	} {
		clone:=*base
		g:=&clone
		up136bTrain(g,arm.refresh,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP136BArmMetric{
			Arm:arm.name,FinalNewUpdates:arm.count,PrimaryNew:p,SecondaryNew:s,
			OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
