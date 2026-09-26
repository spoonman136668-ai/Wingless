package unitary

const UP135BFinalRefreshSchema = "wingless.up135b-final-new-refresh.v1"

type UP135BArmMetric struct {
	Arm                string            `json:"arm"`
	FinalNewUpdates    int               `json:"final_new_updates"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP135BFinalRefreshResult struct {
	Schema                string            `json:"schema"`
	Experiment            string            `json:"experiment"`
	SourceUP134BSeal      string            `json:"source_up134b_seal"`
	StateDimension        int               `json:"state_dimension"`
	GateInputDimension    int               `json:"gate_input_dimension"`
	GroundingEpochs       int               `json:"grounding_epochs"`
	LearningRate          float64           `json:"learning_rate"`
	NewUpdatesPerEpoch    int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch    int               `json:"old_updates_per_epoch"`
	ExtraUpdatesUsed      bool              `json:"extra_updates_used"`
	AdaptiveSplitUsed     bool              `json:"adaptive_split_used"`
	ProjectorRecomputed   bool              `json:"projector_recomputed"`
	ArmMetrics            []UP135BArmMetric `json:"arm_metrics"`
}

type up135bExample struct {
	name string
	verb string
	class int
}

func up135bNewExamples() []up135bExample {
	out:=make([]up135bExample,0,24)
	for _,name:=range up130bNewNames[:2] {
		for _,v:=range up130bStore { out=append(out,up135bExample{name,v,up97bStore}) }
		for _,v:=range up130bObserve { out=append(out,up135bExample{name,v,up97bObserve}) }
		for _,v:=range up130bReport { out=append(out,up135bExample{name,v,up97bReport}) }
	}
	return out
}

func up135bStep(g *up129bGate,e up135bExample,o,r [64]float64) {
	up129bGateStep(g,e.name,e.verb,e.class,o,r)
}

func up135bTrain(g *up129bGate,refresh int,o,r [64]float64) {
	baseNew:=up135bNewExamples()
	for epoch:=0;epoch<20;epoch++ {
		rot:=make([]up135bExample,24)
		for i:=0;i<24;i++ { rot[i]=baseNew[(epoch+i)%24] }
		prefix:=24-refresh
		for i:=0;i<prefix;i++ { up135bStep(g,rot[i],o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for i:=prefix;i<24;i++ { up135bStep(g,rot[i],o,r) }
	}
}

func RunUP135B()(UP135BFinalRefreshResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP135BFinalRefreshResult{
		Schema:UP135BFinalRefreshSchema,Experiment:"UP-135B-final-new-refresh",
		SourceUP134BSeal:"28d28abf2039726a84bff6b7197ad3472207b541",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,ExtraUpdatesUsed:false,AdaptiveSplitUsed:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,refresh:=range []int{0,4,8,12} {
		clone:=*base
		g:=&clone
		up135bTrain(g,refresh,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP135BArmMetric{
			Arm:"refresh_"+itoa(refresh),FinalNewUpdates:refresh,
			PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
