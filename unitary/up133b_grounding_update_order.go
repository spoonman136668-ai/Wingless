package unitary

const UP133BUpdateOrderSchema = "wingless.up133b-grounding-update-order.v1"

type UP133BArmMetric struct {
	Arm                  string            `json:"arm"`
	PrimaryNew           UP130BSplitMetric `json:"primary_new"`
	SecondaryNew         UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained   UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained    UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP133BUpdateOrderResult struct {
	Schema                   string            `json:"schema"`
	Experiment               string            `json:"experiment"`
	SourceUP132BSeal         string            `json:"source_up132b_seal"`
	StateDimension           int               `json:"state_dimension"`
	GateInputDimension       int               `json:"gate_input_dimension"`
	GroundingEpochs          int               `json:"grounding_epochs"`
	LearningRate             float64           `json:"learning_rate"`
	NewGroundingSubjects     int               `json:"new_grounding_subjects"`
	NewUpdatesPerEpoch       int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch       int               `json:"old_updates_per_epoch"`
	ProjectorRecomputed      bool              `json:"projector_recomputed"`
	AdaptiveOrderingUsed     bool              `json:"adaptive_ordering_used"`
	ArmMetrics               []UP133BArmMetric `json:"arm_metrics"`
}

func up133bTrainNewBlock(g *up129bGate,o,r [64]float64) {
	for _,name:=range up130bNewNames[:2] {
		for _,v:=range up130bStore { up129bGateStep(g,name,v,up97bStore,o,r) }
		for _,v:=range up130bObserve { up129bGateStep(g,name,v,up97bObserve,o,r) }
		for _,v:=range up130bReport { up129bGateStep(g,name,v,up97bReport,o,r) }
	}
}

func up133bTrainOldBlock(g *up129bGate,epoch int,o,r [64]float64) {
	old:=up132bOldSurfaces()
	for i,s:=range old {
		subjectIndex:=(i+epoch)%2
		up129bGateStep(g,up121bOriginalNames[subjectIndex],s.verb,s.class,o,r)
	}
}

func up133bTrain(g *up129bGate,arm string,o,r [64]float64) {
	for epoch:=0;epoch<20;epoch++ {
		switch arm {
		case "new_then_old":
			up133bTrainNewBlock(g,o,r)
			up133bTrainOldBlock(g,epoch,o,r)
		case "old_then_new":
			up133bTrainOldBlock(g,epoch,o,r)
			up133bTrainNewBlock(g,o,r)
		case "alternating_epoch_order":
			if epoch%2==0 {
				up133bTrainNewBlock(g,o,r)
				up133bTrainOldBlock(g,epoch,o,r)
			}else{
				up133bTrainOldBlock(g,epoch,o,r)
				up133bTrainNewBlock(g,o,r)
			}
		}
	}
}

func RunUP133B()(UP133BUpdateOrderResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP133BUpdateOrderResult{
		Schema:UP133BUpdateOrderSchema,Experiment:"UP-133B-grounding-update-order",
		SourceUP132BSeal:"e55b98cfb0449a090db4ac248896673235f4f9d3",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewGroundingSubjects:2,NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,
		ProjectorRecomputed:false,AdaptiveOrderingUsed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range []string{"new_then_old","old_then_new","alternating_epoch_order"} {
		clone:=*base
		g:=&clone
		up133bTrain(g,arm,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP133BArmMetric{
			Arm:arm,PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
