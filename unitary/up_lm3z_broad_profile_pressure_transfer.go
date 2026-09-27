package unitary

import "fmt"

const UPLM3ZBroadProfilePressureSchema="wingless.up-lm3z-broad-profile-pressure-transfer.v1"

type UPLM3ZMetric struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
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
type UPLM3ZCoordSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3ZProfileSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	OutcomeRange int `json:"outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
}
type UPLM3ZResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3YSeal string `json:"source_up_lm3y_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PressureLevels []int `json:"pressure_levels"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UPLM3ZMetric `json:"metrics"`
	CoordinateSummaries []UPLM3ZCoordSummary `json:"coordinate_summaries"`
	ProfileSummaries []UPLM3ZProfileSummary `json:"profile_summaries"`
}
func uplm3zKey(coord string,m UPLM3ZMetric)string{
	switch coord{
	case "r+c3+c4":
		return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
	case "r+c2+c3+c4":
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound2,m.CapacityByRound3,m.CapacityByRound4)
	default:
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d|c5=%d",m.ReachableBudget,m.CapacityByRound2,m.CapacityByRound3,m.CapacityByRound4,m.CapacityByRound5)
	}
}
func RunUPLM3Z()(UPLM3ZResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pressures:=[]int{0,1,2,3,4,5,6,7,8}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"r+c3+c4","r+c2+c3+c4","r+c2+c3+c4+c5"}
	res:=UPLM3ZResult{Schema:UPLM3ZBroadProfilePressureSchema,Experiment:"UP-LM3Z-broad-profile-pressure-transfer",SourceUPLM3YSeal:"9cef7e13cd1caa3faa6a475be3dbf7a432503fc3",DeadlineProfiles:profiles,PressureLevels:pressures,ResourceConfigurationsPerCell:80,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,profile:=range profiles{for _,pressure:=range pressures{
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			m:=UPLM3ZMetric{DeadlineProfile:profile,PressureLevel:pressure,Budget:b,ActionStartRound:s,Throughput:tp,ReachableBudget:uplm3fReachable(b,s,tp),CapacityByRound2:uplm3mCapByRound(b,s,tp,2),CapacityByRound3:uplm3mCapByRound(b,s,tp,3),CapacityByRound4:uplm3mCapByRound(b,s,tp,4),CapacityByRound5:uplm3mCapByRound(b,s,tp,5)}
			for _,rot:=range rots{for _,perm:=range perms{m.EarliestFailed+=uplm3uRun(rot,perm,profile,b,s,tp,pressure)}}
			res.Metrics=append(res.Metrics,m)
		}}}
	}}
	type acc struct{min,max int}
	for _,profile:=range profiles{for _,pressure:=range pressures{
		vals:=[]int{}
		for _,m:=range res.Metrics{if m.DeadlineProfile==profile&&m.PressureLevel==pressure{vals=append(vals,m.EarliestFailed)}}
		minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
		for _,v:=range vals{distinct[v]=true;if v<minv{minv=v};if v>maxv{maxv=v}}
		res.ProfileSummaries=append(res.ProfileSummaries,UPLM3ZProfileSummary{DeadlineProfile:profile,PressureLevel:pressure,OutcomeRange:maxv-minv,DistinctOutcomes:len(distinct)})
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,m:=range res.Metrics{
				if m.DeadlineProfile!=profile||m.PressureLevel!=pressure{continue}
				k:=uplm3zKey(coord,m);a:=g[k]
				if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
				if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
			}
			s:=UPLM3ZCoordSummary{DeadlineProfile:profile,PressureLevel:pressure,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
			if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.CoordinateSummaries=append(res.CoordinateSummaries,s)
		}
	}}
	return res,nil
}
