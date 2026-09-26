package unitary

import "math"

const UP165BMarginSchema="wingless.up165b-report-margin-crossings.v1"

type UP165BCrossing struct{
	Path string `json:"path"`
	ReportIndex int `json:"report_index"`
	Split string `json:"split"`
	Name string `json:"name"`
	Verb string `json:"verb"`
	TargetClass int `json:"target_class"`
	BeforeMargin float64 `json:"before_margin"`
	AfterMargin float64 `json:"after_margin"`
	Direction string `json:"direction"`
}
type UP165BStep struct{
	Subject string `json:"subject"`
	ReportIndex int `json:"report_index"`
	PathAOldBefore float64 `json:"path_a_old_before"`
	PathAOldAfter float64 `json:"path_a_old_after"`
	PathBOldBefore float64 `json:"path_b_old_before"`
	PathBOldAfter float64 `json:"path_b_old_after"`
	PathACorrectToIncorrect int `json:"path_a_correct_to_incorrect"`
	PathAIncorrectToCorrect int `json:"path_a_incorrect_to_correct"`
	PathBCorrectToIncorrect int `json:"path_b_correct_to_incorrect"`
	PathBIncorrectToCorrect int `json:"path_b_incorrect_to_correct"`
	PathAMeanMarginBefore float64 `json:"path_a_mean_margin_before"`
	PathAMeanMarginAfter float64 `json:"path_a_mean_margin_after"`
	PathBMeanMarginBefore float64 `json:"path_b_mean_margin_before"`
	PathBMeanMarginAfter float64 `json:"path_b_mean_margin_after"`
	PathAMinMarginBefore float64 `json:"path_a_min_margin_before"`
	PathAMinMarginAfter float64 `json:"path_a_min_margin_after"`
	PathBMinMarginBefore float64 `json:"path_b_min_margin_before"`
	PathBMinMarginAfter float64 `json:"path_b_min_margin_after"`
	DifferentCorrectnessAfter int `json:"different_correctness_after"`
}
type UP165BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP164BSeal string `json:"source_up164b_seal"`
	Subjects int `json:"subjects"`
	ReportSteps int `json:"report_steps"`
	DecisionThreshold string `json:"decision_threshold"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Steps []UP165BStep `json:"steps"`
	Crossings []UP165BCrossing `json:"crossings"`
}
type up165bEval struct{split,name,verb string;class int;margin float64;correct bool}
func up165bMargins(g *up129bGate,o,r [64]float64)[]up165bEval{
	out:=[]up165bEval{}
	sets:=[]struct{name string;names []string}{{"heldout",up121bOriginalNames[4:6]},{"unseen",up121bUnseenNames}}
	for _,set:=range sets{for _,name:=range set.names{for _,s:=range up132bOldSurfaces(){
		raw,proj:=up129bViews(name,s.verb,o,r);p:=up129bProbs(g,raw,proj)
		other:=math.Inf(-1);for k:=0;k<3;k++{if k!=s.class&&p[k]>other{other=p[k]}}
		m:=p[s.class]-other
		out=append(out,up165bEval{split:set.name,name:name,verb:s.verb,class:s.class,margin:m,correct:m>=0})
	}}}
	return out
}
func up165bStats(before,after []up165bEval)(cti,itc int,meanB,meanA,minB,minA float64){
	minB,minA=math.Inf(1),math.Inf(1)
	for i:=range before{
		if before[i].correct&&!after[i].correct{cti++}
		if !before[i].correct&&after[i].correct{itc++}
		meanB+=before[i].margin;meanA+=after[i].margin
		if before[i].margin<minB{minB=before[i].margin};if after[i].margin<minA{minA=after[i].margin}
	}
	n:=float64(len(before));meanB/=n;meanA/=n;return
}
func up165bCross(path string,report int,before,after []up165bEval)[]UP165BCrossing{
	out:=[]UP165BCrossing{}
	for i:=range before{
		dir:=""
		if before[i].correct&&!after[i].correct{dir="correct_to_incorrect"}
		if !before[i].correct&&after[i].correct{dir="incorrect_to_correct"}
		if dir!=""{out=append(out,UP165BCrossing{Path:path,ReportIndex:report,Split:before[i].split,Name:before[i].name,Verb:before[i].verb,TargetClass:before[i].class,BeforeMargin:before[i].margin,AfterMargin:after[i].margin,Direction:dir})}
	}
	return out
}
func RunUP165B()(UP165BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP165BResult{Schema:UP165BMarginSchema,Experiment:"UP-165B-report-margin-crossings",SourceUP164BSeal:"1d2c2ba6e1db1c60c39843633d9047de37e37674",Subjects:6,ReportSteps:2,DecisionThreshold:"probability_margin_zero",DiagnosticUpdatesRetained:false,AdaptiveOrderingUsed:false}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r);a,b:=up164bPreReportStates(base,o,r)
		for _,report:=range []int{13,14}{
			ab:=up165bMargins(a,o,r);bb:=up165bMargins(b,o,r)
			aOldBefore,bOldBefore:=up160bOld(a,o,r),up160bOld(b,o,r)
			up163bApplyBlock(a,[]int{report},o,r);up163bApplyBlock(b,[]int{report},o,r)
			aa:=up165bMargins(a,o,r);ba:=up165bMargins(b,o,r)
			acti,aitc,amb,ama,aminb,amina:=up165bStats(ab,aa)
			bcti,bitc,bmb,bma,bminb,bmina:=up165bStats(bb,ba)
			diff:=0;for i:=range aa{if aa[i].correct!=ba[i].correct{diff++}}
			res.Steps=append(res.Steps,UP165BStep{Subject:subject,ReportIndex:report,PathAOldBefore:aOldBefore,PathAOldAfter:up160bOld(a,o,r),PathBOldBefore:bOldBefore,PathBOldAfter:up160bOld(b,o,r),PathACorrectToIncorrect:acti,PathAIncorrectToCorrect:aitc,PathBCorrectToIncorrect:bcti,PathBIncorrectToCorrect:bitc,PathAMeanMarginBefore:amb,PathAMeanMarginAfter:ama,PathBMeanMarginBefore:bmb,PathBMeanMarginAfter:bma,PathAMinMarginBefore:aminb,PathAMinMarginAfter:amina,PathBMinMarginBefore:bminb,PathBMinMarginAfter:bmina,DifferentCorrectnessAfter:diff})
			res.Crossings=append(res.Crossings,up165bCross("STORE_OBSERVE",report,ab,aa)...)
			res.Crossings=append(res.Crossings,up165bCross("OBSERVE_STORE",report,bb,ba)...)
		}
	}
	return res,nil
}
