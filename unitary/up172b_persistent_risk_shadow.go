package unitary

const UP172BPersistentSchema="wingless.up172b-persistent-risk-shadow.v1"

type UP172BBinSummary struct{
	Scope string `json:"scope"`
	Bin string `json:"bin"`
	Slots int `json:"slots"`
	Crossings int `json:"crossings"`
	CrossingDensity float64 `json:"crossing_density"`
}
type UP172BWarningSummary struct{
	Scope string `json:"scope"`
	Warning string `json:"warning"`
	FlaggedSlots int `json:"flagged_slots"`
	CoveredCrossings int `json:"covered_crossings"`
	TotalCrossings int `json:"total_crossings"`
	CrossingDensity float64 `json:"crossing_density"`
	CrossingCoverage float64 `json:"crossing_coverage"`
}
type UP172BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP171BSeal string `json:"source_up171b_seal"`
	ContextPhase int `json:"context_phase"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	StepsPerPath int `json:"steps_per_path"`
	ExamplesPerState int `json:"examples_per_state"`
	PersistenceSteps int `json:"persistence_steps"`
	BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
	AdaptiveBandUsed bool `json:"adaptive_band_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	UpdatesSuppressed bool `json:"updates_suppressed"`
	Bins []UP172BBinSummary `json:"bins"`
	Warnings []UP172BWarningSummary `json:"warnings"`
}
type up172bCount struct{slots,cross int}
type up172bWarn struct{flagged,covered,totalCross int}
func up172bBin(streak int)string{
	if streak<=0{return "out_of_band"}
	if streak==1{return "current_only"}
	if streak==2{return "persistent_2"}
	return "persistent_3plus"
}
func up172bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func RunUP172B()(UP172BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=17
	paths:=[]struct{name string;indices []int}{
		{"STORE_OBSERVE_REPORT",[]int{0,1,2,3,4,5,6,7,8,9,13,14}},
		{"OBSERVE_STORE_REPORT",[]int{5,6,7,8,9,0,1,2,3,4,13,14}},
	}
	scopes:=[]string{"ALL","STORE","OBSERVE","REPORT"}
	bins:=map[string]map[string]*up172bCount{}
	warns:=map[string]map[string]*up172bWarn{}
	for _,s:=range scopes{
		bins[s]=map[string]*up172bCount{}
		for _,b:=range []string{"out_of_band","current_only","persistent_2","persistent_3plus"}{bins[s][b]=&up172bCount{}}
		warns[s]=map[string]*up172bWarn{"one_step":{},"persistent_2plus":{}}
	}
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,p:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range p.indices{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				for i:=range before{
					if up166bRank(before,i)<=band{streak[i]++}else{streak[i]=0}
				}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r)
				after:=up165bMargins(&g,o,r)
				for i:=range before{
					b:=up172bBin(streak[i]);changed:=before[i].correct!=after[i].correct
					for _,scope:=range []string{"ALL",cls}{
						bins[scope][b].slots++;if changed{bins[scope][b].cross++}
						warns[scope]["one_step"].totalCross++;warns[scope]["persistent_2plus"].totalCross++
						if streak[i]>=1{warns[scope]["one_step"].flagged++;if changed{warns[scope]["one_step"].covered++}}
						if streak[i]>=2{warns[scope]["persistent_2plus"].flagged++;if changed{warns[scope]["persistent_2plus"].covered++}}
					}
				}
			}
		}
	}
	res:=UP172BResult{Schema:UP172BPersistentSchema,Experiment:"UP-172B-persistent-risk-shadow",SourceUP171BSeal:"4516ed4184145beebee4e368b19220f714187313",ContextPhase:17,Subjects:6,Paths:2,StepsPerPath:12,ExamplesPerState:120,PersistenceSteps:2,BandsFrozenBeforeRun:true,AdaptiveBandUsed:false,MaintenanceTriggered:false,UpdatesSuppressed:false}
	for _,scope:=range scopes{
		for _,b:=range []string{"out_of_band","current_only","persistent_2","persistent_3plus"}{
			x:=bins[scope][b];res.Bins=append(res.Bins,UP172BBinSummary{Scope:scope,Bin:b,Slots:x.slots,Crossings:x.cross,CrossingDensity:up172bRate(x.cross,x.slots)})
		}
		for _,w:=range []string{"one_step","persistent_2plus"}{
			x:=warns[scope][w];res.Warnings=append(res.Warnings,UP172BWarningSummary{Scope:scope,Warning:w,FlaggedSlots:x.flagged,CoveredCrossings:x.covered,TotalCrossings:x.totalCross,CrossingDensity:up172bRate(x.covered,x.flagged),CrossingCoverage:up172bRate(x.covered,x.totalCross)})
		}
	}
	return res,nil
}
