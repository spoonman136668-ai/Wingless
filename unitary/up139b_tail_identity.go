package unitary

const UP139BTailIdentitySchema = "wingless.up139b-tail-identity.v1"

type UP139BArmMetric struct {
	Arm                string            `json:"arm"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
	MeanOldRetention   float64           `json:"mean_old_retention"`
	MeanNewAccuracy    float64           `json:"mean_new_accuracy"`
}

type UP139BTailIdentityResult struct {
	Schema              string            `json:"schema"`
	Experiment          string            `json:"experiment"`
	SourceUP138BSeal    string            `json:"source_up138b_seal"`
	StateDimension      int               `json:"state_dimension"`
	GateInputDimension  int               `json:"gate_input_dimension"`
	GroundingEpochs     int               `json:"grounding_epochs"`
	LearningRate        float64           `json:"learning_rate"`
	NewUpdatesPerEpoch  int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch  int               `json:"old_updates_per_epoch"`
	FinalNewUpdates     int               `json:"final_new_updates"`
	ExtraUpdatesUsed    bool              `json:"extra_updates_used"`
	AdaptiveTailChoice  bool              `json:"adaptive_tail_choice"`
	ProjectorRecomputed bool              `json:"projector_recomputed"`
	ArmMetrics          []UP139BArmMetric `json:"arm_metrics"`
}

func up139bRotate(in []up135bExample,shift int) []up135bExample {
	out:=make([]up135bExample,len(in))
	if len(in)==0 { return out }
	shift%=len(in)
	for i:=range in { out[i]=in[(shift+i)%len(in)] }
	return out
}

func up139bWithout(base []up135bExample,start,end int) []up135bExample {
	out:=make([]up135bExample,0,len(base)-(end-start))
	for i,e:=range base {
		if i>=start&&i<end { continue }
		out=append(out,e)
	}
	return out
}

func up139bTrain(g *up129bGate,arm string,o,r [64]float64) {
	base:=up135bNewExamples()
	first4:=append([]up135bExample{},base[:4]...)
	last4:=append([]up135bExample{},base[20:24]...)
	withoutFirst:=up139bWithout(base,0,4)
	withoutLast:=up139bWithout(base,20,24)

	for epoch:=0;epoch<20;epoch++ {
		var pre,tail []up135bExample
		switch arm {
		case "moving_tail_control":
			full:=up139bRotate(base,epoch%24)
			pre=full[:20]
			tail=full[20:24]
		case "fixed_tail_last4":
			pre=up139bRotate(withoutLast,epoch%20)
			tail=last4
		case "fixed_tail_first4":
			pre=up139bRotate(withoutFirst,epoch%20)
			tail=first4
		case "alternating_fixed_tails":
			if epoch%2==0 {
				pre=up139bRotate(withoutLast,epoch%20)
				tail=last4
			}else{
				pre=up139bRotate(withoutFirst,epoch%20)
				tail=first4
			}
		}
		for _,e:=range pre { up135bStep(g,e,o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for _,e:=range tail { up135bStep(g,e,o,r) }
	}
}

func RunUP139B()(UP139BTailIdentityResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	result:=UP139BTailIdentityResult{
		Schema:UP139BTailIdentitySchema,Experiment:"UP-139B-tail-identity",
		SourceUP138BSeal:"661d0328475ccf459a9acb19f8e35f6d79a08338",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		ExtraUpdatesUsed:false,AdaptiveTailChoice:false,ProjectorRecomputed:false,
	}
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	for _,arm:=range []string{"moving_tail_control","fixed_tail_last4","fixed_tail_first4","alternating_fixed_tails"} {
		clone:=*base
		g:=&clone
		up139bTrain(g,arm,o,r)
		p,_:=up130bSplit(g,"primary_new",primary,o,r)
		s,_:=up130bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP139BArmMetric{
			Arm:arm,PrimaryNew:p,SecondaryNew:s,OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
			MeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
			MeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
		})
	}
	return result,nil
}
