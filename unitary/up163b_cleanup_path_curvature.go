package unitary

const UP163BPathCurvatureSchema="wingless.up163b-cleanup-path-curvature.v1"

type UP163BCheckpoint struct{
	Checkpoint int `json:"checkpoint"`
	Block string `json:"block"`
	OldRetention float64 `json:"old_retention"`
	DisplacementNorm float64 `json:"displacement_norm"`
	CosineDisplacementVsOldReference float64 `json:"cosine_displacement_vs_old_reference"`
}
type UP163BSubject struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	PathA string `json:"path_a"`
	PathB string `json:"path_b"`
	PathACheckpoints []UP163BCheckpoint `json:"path_a_checkpoints"`
	PathBCheckpoints []UP163BCheckpoint `json:"path_b_checkpoints"`
	MatchedGateDistance []float64 `json:"matched_gate_distance"`
	MatchedRetentionDeltaAminusB []float64 `json:"matched_retention_delta_a_minus_b"`
}
type UP163BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP162BSeal string `json:"source_up162b_seal"`
	Subjects int `json:"subjects"`
	CheckpointsPerPath int `json:"checkpoints_per_path"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	SubjectsData []UP163BSubject `json:"subjects_data"`
}

func up163bApplyBlock(g *up129bGate,indices []int,o,r [64]float64){
	items:=up158bItems(15,indices)
	up156bApplyOld(g,items,0,len(items),o,r)
}
func up163bCheckpoint(base,g *up129bGate,oldRef []float64,checkpoint int,block string,o,r [64]float64)UP163BCheckpoint{
	d:=up150bDelta(base,g)
	return UP163BCheckpoint{
		Checkpoint:checkpoint,Block:block,OldRetention:up160bOld(g,o,r),
		DisplacementNorm:up150bNorm(d),CosineDisplacementVsOldReference:up150bCos(d,oldRef),
	}
}
func up163bTrace(base *up129bGate,blocks []struct{name string;idx []int},oldRef []float64,o,r [64]float64)([]UP163BCheckpoint,[]*up129bGate){
	g:=*base
	points:=[]UP163BCheckpoint{up163bCheckpoint(base,&g,oldRef,0,"post_report",o,r)}
	states:=[]*up129bGate{func()*up129bGate{x:=g;return &x}()}
	for i,b:=range blocks{
		up163bApplyBlock(&g,b.idx,o,r)
		points=append(points,up163bCheckpoint(base,&g,oldRef,i+1,b.name,o,r))
		x:=g;states=append(states,&x)
	}
	return points,states
}
func RunUP163B()(UP163BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	store:=[]int{0,1,2,3,4};observe:=[]int{5,6,7,8,9};report:=[]int{13,14}
	aBlocks:=[]struct{name string;idx []int}{{"STORE",store},{"OBSERVE",observe},{"REPORT",report}}
	bBlocks:=[]struct{name string;idx []int}{{"OBSERVE",observe},{"STORE",store},{"REPORT",report}}
	res:=UP163BResult{Schema:UP163BPathCurvatureSchema,Experiment:"UP-163B-cleanup-path-curvature",SourceUP162BSeal:"a705c3abb67a77c4575f43e2f872424fe881c899",Subjects:6,CheckpointsPerPath:4,DiagnosticUpdatesRetained:false,AdaptiveOrderingUsed:false}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r)
		oldRef:=up150bOldAnchorDelta(base,o,r)
		ap,as:=up163bTrace(base,aBlocks,oldRef,o,r)
		bp,bs:=up163bTrace(base,bBlocks,oldRef,o,r)
		dists:=make([]float64,0,3);deltas:=make([]float64,0,3)
		for i:=1;i<=3;i++{
			dists=append(dists,up160bDistance(as[i],bs[i]))
			deltas=append(deltas,ap[i].OldRetention-bp[i].OldRetention)
		}
		res.SubjectsData=append(res.SubjectsData,UP163BSubject{
			Subject:subject,TrainingPair:append([]string(nil),up159bPair(subject)...),
			PathA:"STORE_OBSERVE_REPORT",PathB:"OBSERVE_STORE_REPORT",
			PathACheckpoints:ap,PathBCheckpoints:bp,
			MatchedGateDistance:dists,MatchedRetentionDeltaAminusB:deltas,
		})
	}
	return res,nil
}
