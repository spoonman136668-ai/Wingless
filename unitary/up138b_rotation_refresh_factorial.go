package unitary

const UP138BFactorialSchema = "wingless.up138b-rotation-refresh-factorial.v1"

type UP138BArmMetric struct {
	Arm                string            `json:"arm"`
	Rotated            bool              `json:"rotated"`
	FinalNewUpdates    int               `json:"final_new_updates"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
}

type UP138BContrast struct {
	Metric                 string  `json:"metric"`
	RotationMainEffect     float64 `json:"rotation_main_effect"`
	RefreshMainEffect      float64 `json:"refresh_main_effect"`
	RotationRefreshInteraction float64 `json:"rotation_refresh_interaction"`
}

type UP138BFactorialResult struct {
	Schema              string           `json:"schema"`
	Experiment          string           `json:"experiment"`
	SourceUP137BSeal    string           `json:"source_up137b_seal"`
	StateDimension      int              `json:"state_dimension"`
	GateInputDimension  int              `json:"gate_input_dimension"`
	GroundingEpochs     int              `json:"grounding_epochs"`
	LearningRate        float64          `json:"learning_rate"`
	NewUpdatesPerEpoch  int              `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch  int              `json:"old_updates_per_epoch"`
	ExtraUpdatesUsed    bool             `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool            `json:"adaptive_ordering_used"`
	ProjectorRecomputed bool             `json:"projector_recomputed"`
	ArmMetrics          []UP138BArmMetric `json:"arm_metrics"`
	Contrasts           []UP138BContrast  `json:"contrasts"`
}

func up138bOrder(rotated bool,epoch int,base []up135bExample) []up135bExample {
	out:=make([]up135bExample,len(base))
	if !rotated {
		copy(out,base)
		return out
	}
	for i:=range base { out[i]=base[(epoch+i)%len(base)] }
	return out
}

func up138bTrain(g *up129bGate,rotated bool,refresh int,o,r [64]float64) {
	base:=up135bNewExamples()
	for epoch:=0;epoch<20;epoch++ {
		order:=up138bOrder(rotated,epoch,base)
		prefix:=24-refresh
		for i:=0;i<prefix;i++ { up135bStep(g,order[i],o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for i:=prefix;i<24;i++ { up135bStep(g,order[i],o,r) }
	}
}

func up138bValue(m UP138BArmMetric,metric string) float64 {
	switch metric {
	case "primary_new": return m.PrimaryNew.OverallAccuracy
	case "secondary_new": return m.SecondaryNew.OverallAccuracy
	case "old_heldout": return m.OldHeldoutRetained.ClassAccuracy
	case "old_unseen": return m.OldUnseenRetained.ClassAccuracy
	}
	return 0
}

func up138bContrast(metric string,arms []UP138BArmMetric) UP138BContrast {
	var f0,f4,r0,r4 float64
	for _,a:=range arms {
		v:=up138bValue(a,metric)
		switch {
		case !a.Rotated&&a.FinalNewUpdates==0: f0=v
		case !a.Rotated&&a.FinalNewUpdates==4: f4=v
		case a.Rotated&&a.FinalNewUpdates==0: r0=v
		case a.Rotated&&a.FinalNewUpdates==4: r4=v
		}
	}
	return UP138BContrast{
		Metric:metric,
		RotationMainEffect:((r0+r4)-(f0+f4))/2,
		RefreshMainEffect:((f4+r4)-(f0+r0))/2,
		RotationRefreshInteraction:(r4-r0)-(f4-f0),
	}
}

func RunUP138B()(UP138BFactorialResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP138BFactorialResult{
		Schema:UP138BFactorialSchema,Experiment:"UP-138B-rotation-refresh-factorial",
		SourceUP137BSeal:"34eb0e5cf6432290e5cbf0a01c91db8e3f48c41c",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,
		ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	arms:=[]struct{name string;rotated bool;refresh int}{
		{"fixed_refresh0",false,0},
		{"fixed_refresh4",false,4},
		{"rotated_refresh0",true,0},
		{"rotated_refresh4",true,4},
	}
	for _,a:=range arms {
		clone:=*base
		g:=&clone
		up138bTrain(g,a.rotated,a.refresh,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP138BArmMetric{
			Arm:a.name,Rotated:a.rotated,FinalNewUpdates:a.refresh,
			PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
		})
	}
	for _,metric:=range []string{"primary_new","secondary_new","old_heldout","old_unseen"} {
		result.Contrasts=append(result.Contrasts,up138bContrast(metric,result.ArmMetrics))
	}
	return result,nil
}
