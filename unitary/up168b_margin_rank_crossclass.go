package unitary

import "math"

const UP168BCrossClassSchema="wingless.up168b-margin-rank-crossclass.v1"

type UP168BCrossing struct{
	Subject string `json:"subject"`
	Path string `json:"path"`
	ProbeIndex int `json:"probe_index"`
	ProbeClass string `json:"probe_class"`
	Split string `json:"split"`
	Name string `json:"name"`
	Verb string `json:"verb"`
	TargetClass int `json:"target_class"`
	BeforeMargin float64 `json:"before_margin"`
	AfterMargin float64 `json:"after_margin"`
	AbsoluteMarginBefore float64 `json:"absolute_margin_before"`
	VulnerabilityRank int `json:"vulnerability_rank"`
	RankFraction float64 `json:"rank_fraction"`
	Direction string `json:"direction"`
}
type UP168BStep struct{
	Subject string `json:"subject"`
	Path string `json:"path"`
	ProbeIndex int `json:"probe_index"`
	ProbeClass string `json:"probe_class"`
	Examples int `json:"examples"`
	CrossingCount int `json:"crossing_count"`
	MeanCrossingAbsoluteMargin float64 `json:"mean_crossing_absolute_margin"`
	MeanNonCrossingAbsoluteMargin float64 `json:"mean_non_crossing_absolute_margin"`
	MeanCrossingRank float64 `json:"mean_crossing_rank"`
	MaxCrossingRank int `json:"max_crossing_rank"`
	MeanCrossingRankFraction float64 `json:"mean_crossing_rank_fraction"`
	MaxCrossingRankFraction float64 `json:"max_crossing_rank_fraction"`
}
type UP168BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP167BSeal string `json:"source_up167b_seal"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	ProbeIndices []int `json:"probe_indices"`
	IndependentProbeClones bool `json:"independent_probe_clones"`
	RankingSignal string `json:"ranking_signal"`
	PostResultThresholdUsed bool `json:"post_result_threshold_used"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	Steps []UP168BStep `json:"steps"`
	Crossings []UP168BCrossing `json:"crossings"`
}
func up168bProbeClass(idx int)string{if idx<5{return "STORE"};return "OBSERVE"}
func up168bAnalyze(subject,path string,idx int,before,after []up165bEval)(UP168BStep,[]UP168BCrossing){
	out:=[]UP168BCrossing{};sumC,sumN,sumR,sumF:=0.0,0.0,0.0,0.0;nonN,maxR:=0,0
	for i:=range before{
		changed:=before[i].correct!=after[i].correct
		if !changed{sumN+=math.Abs(before[i].margin);nonN++;continue}
		rank:=up166bRank(before,i);if rank>maxR{maxR=rank}
		dir:="incorrect_to_correct";if before[i].correct{dir="correct_to_incorrect"}
		abs:=math.Abs(before[i].margin);frac:=float64(rank)/float64(len(before))
		sumC+=abs;sumR+=float64(rank);sumF+=frac
		out=append(out,UP168BCrossing{Subject:subject,Path:path,ProbeIndex:idx,ProbeClass:up168bProbeClass(idx),Split:before[i].split,Name:before[i].name,Verb:before[i].verb,TargetClass:before[i].class,BeforeMargin:before[i].margin,AfterMargin:after[i].margin,AbsoluteMarginBefore:abs,VulnerabilityRank:rank,RankFraction:frac,Direction:dir})
	}
	s:=UP168BStep{Subject:subject,Path:path,ProbeIndex:idx,ProbeClass:up168bProbeClass(idx),Examples:len(before),CrossingCount:len(out),MaxCrossingRank:maxR}
	if len(out)>0{n:=float64(len(out));s.MeanCrossingAbsoluteMargin=sumC/n;s.MeanCrossingRank=sumR/n;s.MeanCrossingRankFraction=sumF/n;s.MaxCrossingRankFraction=float64(maxR)/float64(len(before))}
	if nonN>0{s.MeanNonCrossingAbsoluteMargin=sumN/float64(nonN)}
	return s,out
}
func RunUP168B()(UP168BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	probes:=[]int{0,1,2,3,4,5,6,7,8,9}
	res:=UP168BResult{Schema:UP168BCrossClassSchema,Experiment:"UP-168B-margin-rank-crossclass",SourceUP167BSeal:"246eb162750f0e72122936502ee00d99b8a49c8a",Subjects:6,Paths:2,ProbeIndices:append([]int(nil),probes...),IndependentProbeClones:true,RankingSignal:"absolute_pre_step_probability_margin",PostResultThresholdUsed:false,DiagnosticUpdatesRetained:false}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r);a0,b0:=up164bPreReportStates(base,o,r)
		for _,idx:=range probes{
			a,b:=*a0,*b0
			ab,bb:=up165bMargins(&a,o,r),up165bMargins(&b,o,r)
			up163bApplyBlock(&a,[]int{idx},o,r);up163bApplyBlock(&b,[]int{idx},o,r)
			aa,ba:=up165bMargins(&a,o,r),up165bMargins(&b,o,r)
			as,ac:=up168bAnalyze(subject,"STORE_OBSERVE",idx,ab,aa)
			bs,bc:=up168bAnalyze(subject,"OBSERVE_STORE",idx,bb,ba)
			res.Steps=append(res.Steps,as,bs);res.Crossings=append(res.Crossings,ac...);res.Crossings=append(res.Crossings,bc...)
		}
	}
	return res,nil
}
