package unitary

const UP164BReportSensitivitySchema="wingless.up164b-report-path-sensitivity.v1"

type UP164BReportStep struct{
	ReportIndex int `json:"report_index"`
	OldRetentionBefore float64 `json:"old_retention_before"`
	OldRetentionAfter float64 `json:"old_retention_after"`
	MarginalRetentionChange float64 `json:"marginal_retention_change"`
	UpdateNorm float64 `json:"update_norm"`
	CosineUpdateVsOldReference float64 `json:"cosine_update_vs_old_reference"`
}
type UP164BSubject struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	PathA string `json:"path_a"`
	PathB string `json:"path_b"`
	PreReportRetentionDeltaAminusB float64 `json:"pre_report_retention_delta_a_minus_b"`
	AfterReport13RetentionDeltaAminusB float64 `json:"after_report13_retention_delta_a_minus_b"`
	AfterReport14RetentionDeltaAminusB float64 `json:"after_report14_retention_delta_a_minus_b"`
	GateDistancePreReport float64 `json:"gate_distance_pre_report"`
	GateDistanceAfterReport13 float64 `json:"gate_distance_after_report13"`
	GateDistanceAfterReport14 float64 `json:"gate_distance_after_report14"`
	PathAReportSteps []UP164BReportStep `json:"path_a_report_steps"`
	PathBReportSteps []UP164BReportStep `json:"path_b_report_steps"`
}
type UP164BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP163BSeal string `json:"source_up163b_seal"`
	Subjects int `json:"subjects"`
	ReportStepsPerPath int `json:"report_steps_per_path"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	SubjectsData []UP164BSubject `json:"subjects_data"`
}
func up164bPreReportStates(base *up129bGate,o,r [64]float64)(*up129bGate,*up129bGate){
	a:=*base;b:=*base
	up163bApplyBlock(&a,[]int{0,1,2,3,4},o,r)
	up163bApplyBlock(&a,[]int{5,6,7,8,9},o,r)
	up163bApplyBlock(&b,[]int{5,6,7,8,9},o,r)
	up163bApplyBlock(&b,[]int{0,1,2,3,4},o,r)
	return &a,&b
}
func up164bApplyReportStep(g *up129bGate,idx int,oldRef []float64,o,r [64]float64)UP164BReportStep{
	before:=up160bOld(g,o,r)
	start:=*g
	up163bApplyBlock(g,[]int{idx},o,r)
	after:=up160bOld(g,o,r)
	d:=up150bDelta(&start,g)
	return UP164BReportStep{ReportIndex:idx,OldRetentionBefore:before,OldRetentionAfter:after,MarginalRetentionChange:after-before,UpdateNorm:up150bNorm(d),CosineUpdateVsOldReference:up150bCos(d,oldRef)}
}
func RunUP164B()(UP164BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP164BResult{Schema:UP164BReportSensitivitySchema,Experiment:"UP-164B-report-path-sensitivity",SourceUP163BSeal:"76db76eb540b1eb572186699d4a6dd62dc36210f",Subjects:6,ReportStepsPerPath:2,DiagnosticUpdatesRetained:false,AdaptiveOrderingUsed:false}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r)
		oldRef:=up150bOldAnchorDelta(base,o,r)
		a,b:=up164bPreReportStates(base,o,r)
		preDelta:=up160bOld(a,o,r)-up160bOld(b,o,r)
		preDist:=up160bDistance(a,b)
		a13:=up164bApplyReportStep(a,13,oldRef,o,r)
		b13:=up164bApplyReportStep(b,13,oldRef,o,r)
		d13:=up160bOld(a,o,r)-up160bOld(b,o,r);dist13:=up160bDistance(a,b)
		a14:=up164bApplyReportStep(a,14,oldRef,o,r)
		b14:=up164bApplyReportStep(b,14,oldRef,o,r)
		d14:=up160bOld(a,o,r)-up160bOld(b,o,r);dist14:=up160bDistance(a,b)
		res.SubjectsData=append(res.SubjectsData,UP164BSubject{
			Subject:subject,TrainingPair:append([]string(nil),up159bPair(subject)...),PathA:"STORE_OBSERVE",PathB:"OBSERVE_STORE",
			PreReportRetentionDeltaAminusB:preDelta,AfterReport13RetentionDeltaAminusB:d13,AfterReport14RetentionDeltaAminusB:d14,
			GateDistancePreReport:preDist,GateDistanceAfterReport13:dist13,GateDistanceAfterReport14:dist14,
			PathAReportSteps:[]UP164BReportStep{a13,a14},PathBReportSteps:[]UP164BReportStep{b13,b14},
		})
	}
	return res,nil
}
