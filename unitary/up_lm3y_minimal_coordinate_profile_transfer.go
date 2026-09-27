package unitary

import "fmt"

const UPLM3YMinimalProfileSchema="wingless.up-lm3y-minimal-coordinate-profile-transfer.v1"

type UPLM3YMetric struct{
	DeadlineProfile string `json:"deadline_profile"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound2 int `json:"capacity_by_round2"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	CapacityByRound5 int `json:"capacity_by_round5"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3YCoordSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3YProfileSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	OutcomeRange int `json:"outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
}
type UPLM3YResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3XSeal string `json:"source_up_lm3x_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PressureLevel int `json:"pressure_level"`
	ResourceConfigurationsPerProfile int `json:"resource_configurations_per_profile"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	AdaptiveProfileSelectionUsed bool `json:"adaptive_profile_selection_used"`
	Metrics []UPLM3YMetric `json:"metrics"`
	CoordinateSummaries []UPLM3YCoordSummary `json:"coordinate_summaries"`
	ProfileSummaries []UPLM3YProfileSummary `json:"profile_summaries"`
}
func uplm3yKey(coord string,m UPLM3YMetric)string{
	switch coord{
	case "r+c3+c4":
		return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
	case "r+c2+c3+c4":
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound2,m.CapacityByRound3,m.CapacityByRound4)
	default:
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d|c5=%d",m.ReachableBudget,m.CapacityByRound2,m.CapacityByRound3,m.CapacityByRound4,m.CapacityByRound5)
	}
}
func RunUPLM3Y()(UPLM3YResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"r+c3+c4","r+c2+c3+c4","r+c2+c3+c4+c5"}
	res:=UPLM3YResult{
		Schema:UPLM3YMinimalProfileSchema,Experiment:"UP-LM3Y-minimal-coordinate-profile-transfer",
		SourceUPLM3XSeal:"b0bb842b995260571a01dfcb649e47b1dbbf0a53",
		DeadlineProfiles:profiles,PressureLevel:4,ResourceConfigurationsPerProfile:80,CandidateCoordinates:coords,
		CounterfactualOnly:true,LiveActivation:false,AdaptiveCoordinateSearchUsed:false,AdaptiveProfileSelectionUsed:false,
	}
	for _,profile:=range profiles{
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			m:=UPLM3YMetric{DeadlineProfile:profile,Budget:b,ActionStartRound:s,Throughput:tp,ReachableBudget:uplm3fReachable(b,s,tp),CapacityByRound2:uplm3mCapByRound(b,s,tp,2),CapacityByRound3:uplm3mCapByRound(b,s,tp,3),CapacityByRound4:uplm3mCapByRound(b,s,tp,4),CapacityByRound5:uplm3mCapByRound(b,s,tp,5)}
			for _,rot:=range rots{for _,perm:=range perms{m.EarliestFailed+=uplm3uRun(rot,perm,profile,b,s,tp,4)}}
			res.Metrics=append(res.Metrics,m)
		}}}
	}
	type acc struct{min,max int}
	for _,profile:=range profiles{
		vals:=[]int{}
		for _,m:=range res.Metrics{if m.DeadlineProfile==profile{vals=append(vals,m.EarliestFailed)}}
		minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
		for _,v:=range vals{distinct[v]=true;if v<minv{minv=v};if v>maxv{maxv=v}}
		res.ProfileSummaries=append(res.ProfileSummaries,UPLM3YProfileSummary{DeadlineProfile:profile,OutcomeRange:maxv-minv,DistinctOutcomes:len(distinct)})
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,m:=range res.Metrics{
				if m.DeadlineProfile!=profile{continue}
				k:=uplm3yKey(coord,m);a:=g[k]
				if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
				if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
			}
			s:=UPLM3YCoordSummary{DeadlineProfile:profile,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
			if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.CoordinateSummaries=append(res.CoordinateSummaries,s)
		}
	}
	return res,nil
}
