package unitary

const UP169BCalibrationSchema="wingless.up169b-margin-rank-class-calibration.v1"

type UP169BClassMetric struct{
	Class string `json:"class"`
	PreregisteredBandMaxRank int `json:"preregistered_band_max_rank"`
	ProbeStates int `json:"probe_states"`
	TotalCrossings int `json:"total_crossings"`
	MaxCrossingRank int `json:"max_crossing_rank"`
	ClassBandCoveredCrossings int `json:"class_band_covered_crossings"`
	ClassBandCoverage float64 `json:"class_band_coverage"`
	Universal4CoveredCrossings int `json:"universal4_covered_crossings"`
	Universal4Coverage float64 `json:"universal4_coverage"`
	ClassBandFlaggedSlots int `json:"class_band_flagged_slots"`
	ClassBandCrossingDensity float64 `json:"class_band_crossing_density"`
	Universal4FlaggedSlots int `json:"universal4_flagged_slots"`
	Universal4CrossingDensity float64 `json:"universal4_crossing_density"`
}
type UP169BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP168BSeal string `json:"source_up168b_seal"`
	CalibrationPhase int `json:"calibration_phase"`
	ParentPhase int `json:"parent_phase"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	ProbeSurfaces int `json:"probe_surfaces"`
	ExamplesPerState int `json:"examples_per_state"`
	BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
	AdaptiveBandUsed bool `json:"adaptive_band_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	ClassMetrics []UP169BClassMetric `json:"class_metrics"`
}
func up169bClass(idx int)string{if idx<5{return "STORE"};if idx<10{return "OBSERVE"};return "REPORT"}
func up169bBand(class string)int{if class=="REPORT"{return 4};return 10}
func up169bPostReport(common *up129bGate,subject string,epoch int,o,r [64]float64)*up129bGate{
	g:=*common;pair:=up159bPair(subject);reportBlock:=up154bReportBlock(subject)
	count:=0
	for _,e:=range up159bExamples(pair){if up153bInBlock(e,reportBlock){continue};up135bStep(&g,e,o,r);count++}
	if count!=20{panic("UP169B_PREFIX_COUNT")}
	pre:=up158bItems(epoch,[]int{10,11,12});up156bApplyOld(&g,pre,0,len(pre),o,r)
	for _,e:=range reportBlock{up135bStep(&g,e,o,r)}
	return &g
}
func up169bPaths(base *up129bGate,epoch int,o,r [64]float64)(*up129bGate,*up129bGate){
	a,b:=*base,*base
	s:=up158bItems(epoch,[]int{0,1,2,3,4});ob:=up158bItems(epoch,[]int{5,6,7,8,9})
	up156bApplyOld(&a,s,0,len(s),o,r);up156bApplyOld(&a,ob,0,len(ob),o,r)
	up156bApplyOld(&b,ob,0,len(ob),o,r);up156bApplyOld(&b,s,0,len(s),o,r)
	return &a,&b
}
func RunUP169B()(UP169BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=16
	type acc struct{states,cross,maxRank,classCovered,u4Covered int}
	a:=map[string]*acc{"STORE":{},"OBSERVE":{},"REPORT":{}}
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r);pa,pb:=up169bPaths(base,epoch,o,r)
		for idx:=0;idx<15;idx++{
			class:=up169bClass(idx);band:=up169bBand(class)
			for _,p0:=range []*up129bGate{pa,pb}{
				p:=*p0;before:=up165bMargins(&p,o,r);items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&p,items,0,1,o,r);after:=up165bMargins(&p,o,r)
				x:=a[class];x.states++
				for i:=range before{
					if before[i].correct==after[i].correct{continue}
					x.cross++;rank:=up166bRank(before,i);if rank>x.maxRank{x.maxRank=rank};if rank<=band{x.classCovered++};if rank<=4{x.u4Covered++}
				}
			}
		}
	}
	res:=UP169BResult{Schema:UP169BCalibrationSchema,Experiment:"UP-169B-margin-rank-class-calibration",SourceUP168BSeal:"5f08682281da1418ae2dc2812ae13d1667d5852f",CalibrationPhase:16,ParentPhase:15,Subjects:6,Paths:2,ProbeSurfaces:15,ExamplesPerState:120,BandsFrozenBeforeRun:true,AdaptiveBandUsed:false,MaintenanceTriggered:false}
	for _,class:=range []string{"STORE","OBSERVE","REPORT"}{
		x:=a[class];band:=up169bBand(class);rate:=func(n,d int)float64{if d==0{return 0};return float64(n)/float64(d)}
		flagged:=x.states*band;u4:=x.states*4
		res.ClassMetrics=append(res.ClassMetrics,UP169BClassMetric{Class:class,PreregisteredBandMaxRank:band,ProbeStates:x.states,TotalCrossings:x.cross,MaxCrossingRank:x.maxRank,ClassBandCoveredCrossings:x.classCovered,ClassBandCoverage:rate(x.classCovered,x.cross),Universal4CoveredCrossings:x.u4Covered,Universal4Coverage:rate(x.u4Covered,x.cross),ClassBandFlaggedSlots:flagged,ClassBandCrossingDensity:rate(x.classCovered,flagged),Universal4FlaggedSlots:u4,Universal4CrossingDensity:rate(x.u4Covered,u4)})
	}
	return res,nil
}
