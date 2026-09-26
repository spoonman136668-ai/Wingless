package unitary

const UP141BReplicationSchema = "wingless.up141b-tail-diversity-replication.v1"

var up141bStore=[]string{"stows","shelves","tucks","deposits"}
var up141bObserve=[]string{"audits","reviews","tracks","samples"}
var up141bReport=[]string{"briefs","mentions","narrates","explains"}

type UP141BArmMetric struct {
	Arm                string            `json:"arm"`
	DistinctTailSets   int               `json:"distinct_tail_sets"`
	PrimaryNew         UP130BSplitMetric `json:"primary_new"`
	SecondaryNew       UP130BSplitMetric `json:"secondary_new"`
	OldHeldoutRetained UP129BSplitMetric `json:"old_heldout_retained"`
	OldUnseenRetained  UP129BSplitMetric `json:"old_unseen_retained"`
	MeanOldRetention   float64           `json:"mean_old_retention"`
	MeanNewAccuracy    float64           `json:"mean_new_accuracy"`
}

type UP141BReplicationResult struct {
	Schema                 string            `json:"schema"`
	Experiment             string            `json:"experiment"`
	SourceUP140BSeal       string            `json:"source_up140b_seal"`
	StateDimension         int               `json:"state_dimension"`
	GateInputDimension     int               `json:"gate_input_dimension"`
	GroundingEpochs        int               `json:"grounding_epochs"`
	LearningRate           float64           `json:"learning_rate"`
	NewUpdatesPerEpoch     int               `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch     int               `json:"old_updates_per_epoch"`
	FinalNewUpdates        int               `json:"final_new_updates"`
	ReplicationSurfaceCount int              `json:"replication_surface_count"`
	ExtraUpdatesUsed       bool              `json:"extra_updates_used"`
	AdaptiveShiftChoice    bool              `json:"adaptive_shift_choice"`
	ProjectorRecomputed    bool              `json:"projector_recomputed"`
	ZeroShotPrimary        UP130BSplitMetric `json:"zero_shot_primary"`
	ZeroShotSecondary      UP130BSplitMetric `json:"zero_shot_secondary"`
	ArmMetrics             []UP141BArmMetric `json:"arm_metrics"`
}

func up141bSurfaces() []struct{verb string;class int} {
	out:=make([]struct{verb string;class int},0,12)
	for _,v:=range up141bStore { out=append(out,struct{verb string;class int}{v,up97bStore}) }
	for _,v:=range up141bObserve { out=append(out,struct{verb string;class int}{v,up97bObserve}) }
	for _,v:=range up141bReport { out=append(out,struct{verb string;class int}{v,up97bReport}) }
	return out
}

func up141bExamples() []up135bExample {
	out:=make([]up135bExample,0,24)
	for _,name:=range up130bNewNames[:2] {
		for _,v:=range up141bStore { out=append(out,up135bExample{name,v,up97bStore}) }
		for _,v:=range up141bObserve { out=append(out,up135bExample{name,v,up97bObserve}) }
		for _,v:=range up141bReport { out=append(out,up135bExample{name,v,up97bReport}) }
	}
	return out
}

func up141bSplit(g *up129bGate,split string,names []string,o,r [64]float64) UP130BSplitMetric {
	hits,total,tp,fp,fn:=0,0,0,0,0
	classHits:=[3]int{}
	classTotal:=[3]int{}
	worst:=1.0
	for _,s:=range up141bSurfaces() {
		m:=up130bSurface(g,split,names,s.verb,s.class,o,r)
		if m.Accuracy<worst { worst=m.Accuracy }
		for _,name:=range names {
			k:=up129bGateClass(g,name,s.verb,o,r)
			total++;classTotal[s.class]++
			if k==s.class { hits++;classHits[s.class]++ }
			if k==up97bStore&&s.class==up97bStore { tp++ }
			if k==up97bStore&&s.class!=up97bStore { fp++ }
			if k!=up97bStore&&s.class==up97bStore { fn++ }
		}
	}
	prec,rec:=1.0,1.0
	if tp+fp>0 { prec=float64(tp)/float64(tp+fp) }
	if tp+fn>0 { rec=float64(tp)/float64(tp+fn) }
	return UP130BSplitMetric{
		Split:split,OverallAccuracy:float64(hits)/float64(total),
		StoreAccuracy:float64(classHits[up97bStore])/float64(classTotal[up97bStore]),
		ObserveAccuracy:float64(classHits[up97bObserve])/float64(classTotal[up97bObserve]),
		ReportAccuracy:float64(classHits[up97bReport])/float64(classTotal[up97bReport]),
		StorePrecision:prec,StoreRecall:rec,WorstSurfaceAcc:worst,Examples:total,
	}
}

func up141bShifts(arm string) []int {
	switch arm {
	case "diversity_1": return []int{0}
	case "diversity_8": return []int{0,3,6,9,12,15,18,21}
	case "diversity_20":
		out:=make([]int,20)
		for i:=0;i<20;i++ { out[i]=i }
		return out
	}
	return nil
}

func up141bTrain(g *up129bGate,arm string,o,r [64]float64) {
	base:=up141bExamples()
	shifts:=up141bShifts(arm)
	for epoch:=0;epoch<20;epoch++ {
		full:=up139bRotate(base,shifts[epoch%len(shifts)])
		for i:=0;i<20;i++ { up135bStep(g,full[i],o,r) }
		up133bTrainOldBlock(g,epoch,o,r)
		for i:=20;i<24;i++ { up135bStep(g,full[i],o,r) }
	}
}

func RunUP141B()(UP141BReplicationResult,error) {
	o,r:=up124bCompetitorDirections()
	base:=up129bTrainGate(o,r)
	primary:=up130bNewNames[4:6]
	secondary:=up131bSecondaryNames()
	result:=UP141BReplicationResult{
		Schema:UP141BReplicationSchema,Experiment:"UP-141B-tail-diversity-replication",
		SourceUP140BSeal:"3dfc55991f8c2835d1e93114ef7faeddea99e74d",
		StateDimension:64,GateInputDimension:128,GroundingEpochs:20,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,ReplicationSurfaceCount:12,
		ExtraUpdatesUsed:false,AdaptiveShiftChoice:false,ProjectorRecomputed:false,
		ZeroShotPrimary:up141bSplit(base,"zero_shot_primary",primary,o,r),
		ZeroShotSecondary:up141bSplit(base,"zero_shot_secondary",secondary,o,r),
	}
	for _,arm:=range []string{"diversity_1","diversity_8","diversity_20"} {
		clone:=*base
		g:=&clone
		up141bTrain(g,arm,o,r)
		p:=up141bSplit(g,"primary_new",primary,o,r)
		s:=up141bSplit(g,"secondary_new",secondary,o,r)
		oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
		oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
		result.ArmMetrics=append(result.ArmMetrics,UP141BArmMetric{
			Arm:arm,DistinctTailSets:len(up141bShifts(arm)),PrimaryNew:p,SecondaryNew:s,
			OldHeldoutRetained:oldH,OldUnseenRetained:oldU,
			MeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
			MeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
		})
	}
	return result,nil
}
