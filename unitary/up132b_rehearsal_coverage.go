package unitary

const UP132BRehearsalCoverageSchema = "wingless.up132b-rehearsal-coverage.v1"

type UP132BArmMetric struct {
	Arm                   string            `json:"arm"`
	OldRehearsalSubjects  int               `json:"old_rehearsal_subjects"`
	OldUpdatesPerEpoch    int               `json:"old_updates_per_epoch"`
	NewUpdatesPerEpoch    int               `json:"new_updates_per_epoch"`
	PrimaryNew            UP130BSplitMetric `json:"primary_new"`
	SecondaryNew          UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained    UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained     UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP132BRehearsalCoverageResult struct {
	Schema                   string             `json:"schema"`
	Experiment               string             `json:"experiment"`
	SourceUP131BSeal         string             `json:"source_up131b_seal"`
	StateDimension           int                `json:"state_dimension"`
	GateInputDimension       int                `json:"gate_input_dimension"`
	GroundingEpochs          int                `json:"grounding_epochs"`
	LearningRate             float64            `json:"learning_rate"`
	NewGroundingSubjects     int                `json:"new_grounding_subjects"`
	OldUpdatesPerEpoch       int                `json:"old_updates_per_epoch"`
	ProjectorRecomputed      bool               `json:"projector_recomputed"`
	AdaptiveSubjectChoice    bool               `json:"adaptive_subject_choice"`
	ArmMetrics               []UP132BArmMetric  `json:"arm_metrics"`
}

type up132bSurfaceSpec struct {
	verb string
	class int
}

func up132bOldSurfaces() []up132bSurfaceSpec {
	out:=make([]up132bSurfaceSpec,0,15)
	for _,v:=range up128bStore { out=append(out,up132bSurfaceSpec{v,up97bStore}) }
	for _,v:=range up128bObserve { out=append(out,up132bSurfaceSpec{v,up97bObserve}) }
	for _,v:=range up128bReport { out=append(out,up132bSurfaceSpec{v,up97bReport}) }
	return out
}

func up132bTrain(g *up129bGate,oldSubjectCount int,o,r [64]float64) {
	old:=up132bOldSurfaces()
	for epoch:=0;epoch<20;epoch++ {
		for _,name:=range up130bNewNames[:2] {
			for _,v:=range up130bStore { up129bGateStep(g,name,v,up97bStore,o,r) }
			for _,v:=range up130bObserve { up129bGateStep(g,name,v,up97bObserve,o,r) }
			for _,v:=range up130bReport { up129bGateStep(g,name,v,up97bReport,o,r) }
		}
		for i,s:=range old {
			subjectIndex:=0
			if oldSubjectCount>1 { subjectIndex=(i+epoch)%oldSubjectCount }
			up129bGateStep(g,up121bOriginalNames[subjectIndex],s.verb,s.class,o,r)
		}
	}
}

func RunUP132B()(UP132BRehearsalCoverageResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP132BRehearsalCoverageResult{
		Schema:UP132BRehearsalCoverageSchema,Experiment:"UP-132B-rehearsal-coverage",
		SourceUP131BSeal:"72638010e1a03e3b32249f6c84ee05b38b21f0e5",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewGroundingSubjects:2,OldUpdatesPerEpoch:15,ProjectorRecomputed:false,AdaptiveSubjectChoice:false,
	}
	arms:=[]struct{name string;subjects int}{
		{"ada_only",1},{"two_subject_balanced",2},{"four_subject_balanced",4},
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range arms {
		clone:=*base
		g:=&clone
		up132bTrain(g,arm.subjects,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP132BArmMetric{
			Arm:arm.name,OldRehearsalSubjects:arm.subjects,OldUpdatesPerEpoch:15,NewUpdatesPerEpoch:24,
			PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	return result,nil
}
