package unitary

const UP134BTerminalAnchorSchema = "wingless.up134b-terminal-anchor-allocation.v1"

type UP134BArmMetric struct {
	Arm                   string            `json:"arm"`
	TerminalOldUpdates    int               `json:"terminal_old_updates"`
	PrimaryNew            UP130BSplitMetric `json:"primary_new"`
	SecondaryNew          UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained    UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained     UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP134BTerminalAnchorResult struct {
	Schema                 string            `json:"schema"`
	Experiment             string            `json:"experiment"`
	SourceUP133BSeal       string            `json:"source_up133b_seal"`
	StateDimension         int               `json:"state_dimension"`
	GateInputDimension     int               `json:"gate_input_dimension"`
	GroundingEpochs        int               `json:"grounding_epochs"`
	LearningRate           float64           `json:"learning_rate"`
	NewUpdatesPerEpoch     int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch     int               `json:"old_updates_per_epoch"`
	AdaptiveSplitUsed      bool              `json:"adaptive_split_used"`
	ProjectorRecomputed    bool              `json:"projector_recomputed"`
	ArmMetrics             []UP134BArmMetric `json:"arm_metrics"`
}

func up134bOldExample(g *up129bGate,epoch,originalIndex int,o,r [64]float64) {
	old:=up132bOldSurfaces()
	s:=old[originalIndex]
	subjectIndex:=(originalIndex+epoch)%2
	up129bGateStep(g,up121bOriginalNames[subjectIndex],s.verb,s.class,o,r)
}

func up134bTrain(g *up129bGate,terminal int,o,r [64]float64) {
	for epoch:=0;epoch<20;epoch++ {
		rot:=make([]int,15)
		for j:=0;j<15;j++ { rot[j]=(epoch+j)%15 }

		prefix:=15-terminal
		for j:=0;j<prefix;j++ { up134bOldExample(g,epoch,rot[j],o,r) }

		up133bTrainNewBlock(g,o,r)

		for j:=prefix;j<15;j++ { up134bOldExample(g,epoch,rot[j],o,r) }
	}
}

func RunUP134B()(UP134BTerminalAnchorResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP134BTerminalAnchorResult{
		Schema:UP134BTerminalAnchorSchema,Experiment:"UP-134B-terminal-anchor-allocation",
		SourceUP133BSeal:"9c724753be853be217da23f2355b2841464d0b40",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,AdaptiveSplitUsed:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,terminal:=range []int{0,5,10,15} {
		clone:=*base
		g:=&clone
		up134bTrain(g,terminal,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP134BArmMetric{
			Arm:"terminal_"+itoa(terminal),TerminalOldUpdates:terminal,
			PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
