package unitary

import "fmt"

const UPLM3UProfilePressureSchema="wingless.up-lm3u-profile-pressure-map.v1"

type UPLM3UMetric struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3UCoordSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3UProfileSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	OutcomeRange int `json:"outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
}
type UPLM3UResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3TSeal string `json:"source_up_lm3t_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PressureLevels []int `json:"pressure_levels"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UPLM3UMetric `json:"metrics"`
	CoordinateSummaries []UPLM3UCoordSummary `json:"coordinate_summaries"`
	ProfileSummaries []UPLM3UProfileSummary `json:"profile_summaries"`
}
func uplm3uRun(rot int,perm,profile string,budget,start,tp,pressure int)(failed int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	for w:=0;w<pressure;w++{for i:=range arms{arms[i].r.write(fmt.Sprintf("3u-pre-%s-%d-%d",profile,w,i),"x")}}
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3u-%s-%d-%d",profile,round,i),"x")}
	}
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return
}
func RunUPLM3U()(UPLM3UResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"};pressures:=[]int{0,1,2,3,4,5,6,7,8}
	coords:=[]string{"reachable_budget+capacity_by_round4","reachable_budget+capacity_by_round3+capacity_by_round4"}
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3UResult{Schema:UPLM3UProfilePressureSchema,Experiment:"UP-LM3U-profile-pressure-map",SourceUPLM3TSeal:"50df407515a97ac6f8141a233a6381dc074e8dab",DeadlineProfiles:profiles,PressureLevels:pressures,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,profile:=range profiles{for _,pressure:=range pressures{for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
		m:=UPLM3UMetric{DeadlineProfile:profile,PressureLevel:pressure,Budget:b,ActionStartRound:s,Throughput:tp,ReachableBudget:uplm3fReachable(b,s,tp),CapacityByRound3:uplm3mCapByRound(b,s,tp,3),CapacityByRound4:uplm3mCapByRound(b,s,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{m.EarliestFailed+=uplm3uRun(rot,perm,profile,b,s,tp,pressure)}}
		res.Metrics=append(res.Metrics,m)
	}}}}}
	type acc struct{min,max int}
	for _,profile:=range profiles{for _,pressure:=range pressures{
		vals:=[]int{}
		for _,m:=range res.Metrics{if m.DeadlineProfile==profile&&m.PressureLevel==pressure{vals=append(vals,m.EarliestFailed)}}
		minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
		for _,v:=range vals{distinct[v]=true;if v<minv{minv=v};if v>maxv{maxv=v}}
		res.ProfileSummaries=append(res.ProfileSummaries,UPLM3UProfileSummary{DeadlineProfile:profile,PressureLevel:pressure,OutcomeRange:maxv-minv,DistinctOutcomes:len(distinct)})
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,m:=range res.Metrics{
				if m.DeadlineProfile!=profile||m.PressureLevel!=pressure{continue}
				k:=uplm3tKey(coord,UPLM3TMetric{ReachableBudget:m.ReachableBudget,CapacityByRound3:m.CapacityByRound3,CapacityByRound4:m.CapacityByRound4})
				a:=g[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
				if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
			}
			s:=UPLM3UCoordSummary{DeadlineProfile:profile,PressureLevel:pressure,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
			if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.CoordinateSummaries=append(res.CoordinateSummaries,s)
		}
	}}
	return res,nil
}
