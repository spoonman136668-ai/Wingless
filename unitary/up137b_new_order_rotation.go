package unitary

const UP137BNewOrderRotationSchema = "wingless.up137b-new-order-rotation.v1"

type UP137BArmMetric struct {
	Arm                string            `json:"arm"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP137BNewOrderRotationResult struct {
	Schema               string            `json:"schema"`
	Experiment           string            `json:"experiment"`
	SourceUP136BSeal     string            `json:"source_up136b_seal"`
	StateDimension       int               `json:"state_dimension"`
	GateInputDimension   int               `json:"gate_input_dimension"`
	GroundingEpochs      int               `json:"grounding_epochs"`
	LearningRate         float64           `json:"learning_rate"`
	NewUpdatesPerEpoch   int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch   int               `json:"old_updates_per_epoch"`
	FinalNewRefreshUsed  bool              `json:"final_new_refresh_used"`
	ExtraUpdatesUsed     bool              `json:"extra_updates_used"`
	ProjectorRecomputed  bool              `json:"projector_recomputed"`
	ArmMetrics           []UP137BArmMetric `json:"arm_metrics"`
}

func up137bOrder(arm string,epoch int,base []up135bExample) []up135bExample {
	out:=make([]up135bExample,len(base))
	switch arm {
	case "fixed_new_order":
		copy(out,base)
	case "rotate_new_by_epoch":
		for i:=range base { out[i]=base[(epoch+i)%len(base)] }
	case "reverse_rotate_new_by_epoch":
		rev:=make([]up135bExample,len(base))
		for i:=range base { rev[i]=base[len(base)-1-i] }
		for i:=range rev { out[i]=rev[(epoch+i)%len(rev)] }
	}
	return out
}

func up137bTrain(g *up129bGate,arm string,o,r [64]float64) {
	base:=up135bNewExamples()
	for epoch:=0;epoch<20;epoch++ {
		for _,e:=range up137bOrder(arm,epoch,base) { up135bStep(g,e,o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
	}
}

func RunUP137B()(UP137BNewOrderRotationResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP137BNewOrderRotationResult{
		Schema:UP137BNewOrderRotationSchema,Experiment:"UP-137B-new-order-rotation",
		SourceUP136BSeal:"e86911ab08a6daef957c4c7873e60aa21e1e3695",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,
		FinalNewRefreshUsed:false,ExtraUpdatesUsed:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range []string{"fixed_new_order","rotate_new_by_epoch","reverse_rotate_new_by_epoch"} {
		clone:=*base
		g:=&clone
		up137bTrain(g,arm,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP137BArmMetric{
			Arm:arm,PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
