package unitary

const UP131BFewShotSchema = "wingless.up131b-fewshot-lexical-grounding.v1"

type UP131BArmMetric struct {
	Arm                    string            `json:"arm"`
	GroundingSubjects      int               `json:"grounding_subjects"`
	PrimaryNew             UP130BSplitMetric `json:"primary_new"`
	SecondaryNew           UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained     UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained      UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP131BSurfaceMetric struct {
	Arm     string              `json:"arm"`
	Metric  UP130BSurfaceMetric `json:"metric"`
}

type UP131BFewShotResult struct {
	Schema                   string                 `json:"schema"`
	Experiment               string                 `json:"experiment"`
	SourceUP130BSeal         string                 `json:"source_up130b_seal"`
	StateDimension           int                    `json:"state_dimension"`
	GateInputDimension       int                    `json:"gate_input_dimension"`
	OriginalGateEpochs       int                    `json:"original_gate_epochs"`
	GroundingEpochs          int                    `json:"grounding_epochs"`
	LearningRate             float64                `json:"learning_rate"`
	OriginalRehearsalSubject string                 `json:"original_rehearsal_subject"`
	ProjectorRecomputed      bool                   `json:"projector_recomputed"`
	ThresholdChanged         bool                   `json:"threshold_changed"`
	ArmMetrics               []UP131BArmMetric      `json:"arm_metrics"`
	SurfaceMetrics           []UP131BSurfaceMetric  `json:"surface_metrics"`
}

func up131bTrainNew(g *up129bGate,subjects []string,o,r [64]float64){
	for epoch:=0;epoch<20;epoch++{
		for _,name:=range subjects{
			for _,v:=range up130bStore{up129bGateStep(g,name,v,up97bStore,o,r)}
			for _,v:=range up130bObserve{up129bGateStep(g,name,v,up97bObserve,o,r)}
			for _,v:=range up130bReport{up129bGateStep(g,name,v,up97bReport,o,r)}
		}
		for i:=0;i<5;i++{
			up129bGateStep(g,up121bOriginalNames[0],up128bStore[i],up97bStore,o,r)
			up129bGateStep(g,up121bOriginalNames[0],up128bObserve[i],up97bObserve,o,r)
			up129bGateStep(g,up121bOriginalNames[0],up128bReport[i],up97bReport,o,r)
		}
	}
}

func up131bSecondaryNames() []string{
	out:=append([]string{},up121bOriginalNames[4:6]...)
	out=append(out,up121bUnseenNames...)
	return out
}

func RunUP131B()(UP131BFewShotResult,error){
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP131BFewShotResult{
		Schema:UP131BFewShotSchema,Experiment:"UP-131B-fewshot-lexical-grounding",
		SourceUP130BSeal:"9d58620a0871e67ecfcec6f2f2f59c1d3e0fcb80",
		StateDimension:64,GateInputDimension:128,OriginalGateEpochs:20,GroundingEpochs:20,LearningRate:0.08,
		OriginalRehearsalSubject:"ada",ProjectorRecomputed:false,ThresholdChanged:false,
	}
	arms:=[]struct{name string;n int}{
		{"zero_shot",0},{"one_subject",1},{"two_subjects",2},{"four_subjects",4},
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range arms{
		clone:=*base
		g:=&clone
		if arm.n>0{up131bTrainNew(g,up130bNewNames[:arm.n],o,r)}
		p,surf:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP131BArmMetric{
			Arm:arm.name,GroundingSubjects:arm.n,PrimaryNew:p,SecondaryNew:s,
			OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
		for _,m:=range surf{result.SurfaceMetrics=append(result.SurfaceMetrics,UP131BSurfaceMetric{Arm:arm.name,Metric:m})}
	}
	return result,nil
}
