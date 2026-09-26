package unitary

const UP167BReportGeneralizationSchema="wingless.up167b-margin-rank-report-generalization.v1"

type UP167BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP166BSeal string `json:"source_up166b_seal"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	ReportSurfaces []int `json:"report_surfaces"`
	IndependentProbeClones bool `json:"independent_probe_clones"`
	RankingSignal string `json:"ranking_signal"`
	PostResultThresholdUsed bool `json:"post_result_threshold_used"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	Steps []UP166BStep `json:"steps"`
	Crossings []UP166BCrossingRank `json:"crossings"`
}
func RunUP167B()(UP167BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	reports:=[]int{10,11,12,13,14}
	res:=UP167BResult{Schema:UP167BReportGeneralizationSchema,Experiment:"UP-167B-margin-rank-report-generalization",SourceUP166BSeal:"46d8e1d00cf6cde7c0eb8cee6f5fb610466ae902",Subjects:6,Paths:2,ReportSurfaces:append([]int(nil),reports...),IndependentProbeClones:true,RankingSignal:"absolute_pre_step_probability_margin",PostResultThresholdUsed:false,DiagnosticUpdatesRetained:false}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r);a0,b0:=up164bPreReportStates(base,o,r)
		for _,report:=range reports{
			a,b:=*a0,*b0
			ab,bb:=up165bMargins(&a,o,r),up165bMargins(&b,o,r)
			up163bApplyBlock(&a,[]int{report},o,r);up163bApplyBlock(&b,[]int{report},o,r)
			aa,ba:=up165bMargins(&a,o,r),up165bMargins(&b,o,r)
			as,ac:=up166bAnalyze(subject,"STORE_OBSERVE",report,ab,aa)
			bs,bc:=up166bAnalyze(subject,"OBSERVE_STORE",report,bb,ba)
			res.Steps=append(res.Steps,as,bs);res.Crossings=append(res.Crossings,ac...);res.Crossings=append(res.Crossings,bc...)
		}
	}
	return res,nil
}
