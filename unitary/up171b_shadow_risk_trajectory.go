package unitary

import "math"

const UP171BShadowSchema="wingless.up171b-shadow-risk-trajectory.v1"

type UP171BStep struct{
	Subject string `json:"subject"`
	Path string `json:"path"`
	Step int `json:"step"`
	ProbeIndex int `json:"probe_index"`
	Class string `json:"class"`
	BandMaxRank int `json:"band_max_rank"`
	FlaggedExamples int `json:"flagged_examples"`
	Crossings int `json:"crossings"`
	CoveredCrossings int `json:"covered_crossings"`
	FalseWarningSlots int `json:"false_warning_slots"`
	Coverage float64 `json:"coverage"`
	CrossingDensity float64 `json:"crossing_density"`
	MinAbsoluteMargin float64 `json:"min_absolute_margin"`
	MeanAbsoluteMargin float64 `json:"mean_absolute_margin"`
}
type UP171BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP170BSeal string `json:"source_up170b_seal"`
	ContextPhase int `json:"context_phase"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	StepsPerPath int `json:"steps_per_path"`
	ExamplesPerState int `json:"examples_per_state"`
	BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
	AdaptiveBandUsed bool `json:"adaptive_band_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	UpdatesSuppressed bool `json:"updates_suppressed"`
	Steps []UP171BStep `json:"steps"`
}
func up171bClass(idx int)string{if idx<5{return "STORE"};if idx<10{return "OBSERVE"};return "REPORT"}
func up171bBand(cls string)int{if cls=="REPORT"{return 4};return 10}
func up171bAnalyze(subject,path string,step,idx int,before,after []up165bEval)UP171BStep{
	cls:=up171bClass(idx);band:=up171bBand(cls);cross,covered:=0,0
	minAbs,sumAbs:=math.Inf(1),0.0
	for i:=range before{
		a:=math.Abs(before[i].margin);sumAbs+=a;if a<minAbs{minAbs=a}
		if before[i].correct!=after[i].correct{
			cross++;if up166bRank(before,i)<=band{covered++}
		}
	}
	rate:=func(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
	return UP171BStep{Subject:subject,Path:path,Step:step,ProbeIndex:idx,Class:cls,BandMaxRank:band,FlaggedExamples:band,Crossings:cross,CoveredCrossings:covered,FalseWarningSlots:band-covered,Coverage:rate(covered,cross),CrossingDensity:rate(covered,band),MinAbsoluteMargin:minAbs,MeanAbsoluteMargin:sumAbs/float64(len(before))}
}
func RunUP171B()(UP171BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=17
	paths:=[]struct{name string;indices []int}{
		{"STORE_OBSERVE_REPORT",[]int{0,1,2,3,4,5,6,7,8,9,13,14}},
		{"OBSERVE_STORE_REPORT",[]int{5,6,7,8,9,0,1,2,3,4,13,14}},
	}
	res:=UP171BResult{Schema:UP171BShadowSchema,Experiment:"UP-171B-shadow-risk-trajectory",SourceUP170BSeal:"7447e0cf2a50b0fd76b63dea442c3c0a3b2962e2",ContextPhase:17,Subjects:6,Paths:2,StepsPerPath:12,ExamplesPerState:120,BandsFrozenBeforeRun:true,AdaptiveBandUsed:false,MaintenanceTriggered:false,UpdatesSuppressed:false}
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,p:=range paths{
			g:=*base
			for si,idx:=range p.indices{
				before:=up165bMargins(&g,o,r)
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r)
				after:=up165bMargins(&g,o,r)
				res.Steps=append(res.Steps,up171bAnalyze(subject,p.name,si+1,idx,before,after))
			}
		}
	}
	return res,nil
}
