package unitary

import "fmt"

const UPLM3VResourceGridSchema="wingless.up-lm3v-resource-grid-transfer.v1"

type UPLM3VMetric struct{
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3VSummary struct{
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3VResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3USeal string `json:"source_up_lm3u_seal"`
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	ResourceConfigurations int `json:"resource_configurations"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	GlobalOutcomeRange int `json:"global_outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveResourceSelectionUsed bool `json:"adaptive_resource_selection_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3VMetric `json:"metrics"`
	Summaries []UPLM3VSummary `json:"summaries"`
}
func uplm3vKey(coord string,m UPLM3VMetric)string{
	if coord=="reachable_budget+capacity_by_round4"{return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)}
	return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
}
func RunUPLM3V()(UPLM3VResult,error){
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"reachable_budget+capacity_by_round4","reachable_budget+capacity_by_round3+capacity_by_round4"}
	res:=UPLM3VResult{Schema:UPLM3VResourceGridSchema,Experiment:"UP-LM3V-resource-grid-transfer",SourceUPLM3USeal:"e037a27f4419c0ac8ac43c4b0ff3fee2ec998cfe",DeadlineProfile:"hybrid_min",PressureLevel:4,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,ResourceConfigurations:len(budgets)*len(starts)*len(tps),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveResourceSelectionUsed:false,AdaptiveCoordinateSearchUsed:false}
	vals:=[]int{}
	for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
		m:=UPLM3VMetric{Budget:b,ActionStartRound:s,Throughput:tp,ReachableBudget:uplm3fReachable(b,s,tp),CapacityByRound3:uplm3mCapByRound(b,s,tp,3),CapacityByRound4:uplm3mCapByRound(b,s,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{m.EarliestFailed+=uplm3uRun(rot,perm,"hybrid_min",b,s,tp,4)}}
		res.Metrics=append(res.Metrics,m);vals=append(vals,m.EarliestFailed)
	}}}
	minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
	for _,v:=range vals{distinct[v]=true;if v<minv{minv=v};if v>maxv{maxv=v}}
	res.GlobalOutcomeRange=maxv-minv;res.DistinctOutcomes=len(distinct)
	type acc struct{min,max int}
	for _,coord:=range coords{
		g:=map[string]*acc{}
		for _,m:=range res.Metrics{
			k:=uplm3vKey(coord,m);a:=g[k]
			if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
			if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3VSummary{Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
