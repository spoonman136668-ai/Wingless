package unitary

const UP174BRankConditionedSchema="wingless.up174b-rank-conditioned-transition-hazard.v1"

type UP174BMetric struct{
	Scope string `json:"scope"`
	RankStratum string `json:"rank_stratum"`
	TemporalState string `json:"temporal_state"`
	Slots int `json:"slots"`
	Crossings int `json:"crossings"`
	CrossingDensity float64 `json:"crossing_density"`
}
type UP174BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP173BSeal string `json:"source_up173b_seal"`
	ContextPhase int `json:"context_phase"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	StepsPerPath int `json:"steps_per_path"`
	ExamplesPerState int `json:"examples_per_state"`
	RankStrata []string `json:"rank_strata"`
	TemporalStates []string `json:"temporal_states"`
	BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
	AdaptiveBandUsed bool `json:"adaptive_band_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	UpdatesSuppressed bool `json:"updates_suppressed"`
	Metrics []UP174BMetric `json:"metrics"`
}
type up174bCount struct{slots,cross int}
func up174bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func up174bStratum(rank int)string{if rank>=1&&rank<=4{return "rank_1_4"};if rank>=5&&rank<=10{return "rank_5_10"};return ""}
func up174bTemporal(streak int)string{if streak==1{return "current_only"};if streak==2{return "transition_to_persistent"};if streak>=3{return "established_persistent"};return ""}
func RunUP174B()(UP174BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=19
	paths:=[]struct{name string;indices []int}{
		{"STORE_OBSERVE_REPORT",[]int{0,1,2,3,4,5,6,7,8,9,13,14}},
		{"OBSERVE_STORE_REPORT",[]int{5,6,7,8,9,0,1,2,3,4,13,14}},
	}
	scopes:=[]string{"ALL","STORE","OBSERVE","REPORT"}
	strata:=[]string{"rank_1_4","rank_5_10"}
	states:=[]string{"current_only","transition_to_persistent","established_persistent"}
	acc:=map[string]map[string]map[string]*up174bCount{}
	for _,s:=range scopes{acc[s]=map[string]map[string]*up174bCount{};for _,st:=range strata{acc[s][st]=map[string]*up174bCount{};for _,t:=range states{acc[s][st][t]=&up174bCount{}}}}
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,p:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range p.indices{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				ranks:=make([]int,len(before))
				for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r)
				after:=up165bMargins(&g,o,r)
				for i:=range before{
					rs:=up174bStratum(ranks[i]);ts:=up174bTemporal(streak[i]);if rs==""||ts==""{continue}
					changed:=before[i].correct!=after[i].correct
					for _,scope:=range []string{"ALL",cls}{x:=acc[scope][rs][ts];x.slots++;if changed{x.cross++}}
				}
			}
		}
	}
	res:=UP174BResult{Schema:UP174BRankConditionedSchema,Experiment:"UP-174B-rank-conditioned-transition-hazard",SourceUP173BSeal:"d3c7ab99f3facc5b7c84a69db20bf8baefbd40ec",ContextPhase:19,Subjects:6,Paths:2,StepsPerPath:12,ExamplesPerState:120,RankStrata:strata,TemporalStates:states,BandsFrozenBeforeRun:true,AdaptiveBandUsed:false,MaintenanceTriggered:false,UpdatesSuppressed:false}
	for _,scope:=range scopes{for _,rs:=range strata{for _,ts:=range states{x:=acc[scope][rs][ts];res.Metrics=append(res.Metrics,UP174BMetric{Scope:scope,RankStratum:rs,TemporalState:ts,Slots:x.slots,Crossings:x.cross,CrossingDensity:up174bRate(x.cross,x.slots)})}}}
	return res,nil
}
