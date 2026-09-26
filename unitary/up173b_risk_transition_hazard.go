package unitary

const UP173BTransitionSchema="wingless.up173b-risk-transition-hazard.v1"

type UP173BMetric struct{
	Scope string `json:"scope"`
	State string `json:"state"`
	Slots int `json:"slots"`
	Crossings int `json:"crossings"`
	CrossingDensity float64 `json:"crossing_density"`
}
type UP173BSummary struct{
	Scope string `json:"scope"`
	TransitionDensity float64 `json:"transition_density"`
	EstablishedDensity float64 `json:"established_density"`
	TransitionVsEstablishedRatio float64 `json:"transition_vs_established_ratio"`
}
type UP173BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP172BSeal string `json:"source_up172b_seal"`
	ContextPhase int `json:"context_phase"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	StepsPerPath int `json:"steps_per_path"`
	ExamplesPerState int `json:"examples_per_state"`
	BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	UpdatesSuppressed bool `json:"updates_suppressed"`
	Metrics []UP173BMetric `json:"metrics"`
	Summaries []UP173BSummary `json:"summaries"`
}
type up173bCount struct{slots,cross int}
func up173bState(streak int)string{if streak<=0{return "out_of_band"};if streak==1{return "current_only"};if streak==2{return "transition_to_persistent"};return "established_persistent"}
func up173bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func RunUP173B()(UP173BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=18
	paths:=[]struct{name string;indices []int}{
		{"STORE_OBSERVE_REPORT",[]int{0,1,2,3,4,5,6,7,8,9,13,14}},
		{"OBSERVE_STORE_REPORT",[]int{5,6,7,8,9,0,1,2,3,4,13,14}},
	}
	scopes:=[]string{"ALL","STORE","OBSERVE","REPORT"}
	states:=[]string{"out_of_band","current_only","transition_to_persistent","established_persistent"}
	acc:=map[string]map[string]*up173bCount{}
	for _,s:=range scopes{acc[s]=map[string]*up173bCount{};for _,st:=range states{acc[s][st]=&up173bCount{}}}
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,p:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range p.indices{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				for i:=range before{if up166bRank(before,i)<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r)
				after:=up165bMargins(&g,o,r)
				for i:=range before{
					st:=up173bState(streak[i]);changed:=before[i].correct!=after[i].correct
					for _,scope:=range []string{"ALL",cls}{acc[scope][st].slots++;if changed{acc[scope][st].cross++}}
				}
			}
		}
	}
	res:=UP173BResult{Schema:UP173BTransitionSchema,Experiment:"UP-173B-risk-transition-hazard",SourceUP172BSeal:"c41a28b9fe3c42f6e43aeaaeb0a6d334c96b672a",ContextPhase:18,Subjects:6,Paths:2,StepsPerPath:12,ExamplesPerState:120,BandsFrozenBeforeRun:true,MaintenanceTriggered:false,UpdatesSuppressed:false}
	for _,scope:=range scopes{
		for _,st:=range states{x:=acc[scope][st];res.Metrics=append(res.Metrics,UP173BMetric{Scope:scope,State:st,Slots:x.slots,Crossings:x.cross,CrossingDensity:up173bRate(x.cross,x.slots)})}
		td:=up173bRate(acc[scope]["transition_to_persistent"].cross,acc[scope]["transition_to_persistent"].slots)
		ed:=up173bRate(acc[scope]["established_persistent"].cross,acc[scope]["established_persistent"].slots)
		ratio:=0.0;if ed>0{ratio=td/ed}
		res.Summaries=append(res.Summaries,UP173BSummary{Scope:scope,TransitionDensity:td,EstablishedDensity:ed,TransitionVsEstablishedRatio:ratio})
	}
	return res,nil
}
