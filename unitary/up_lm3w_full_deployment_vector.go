package unitary

import "fmt"

const UPLM3WFullVectorSchema="wingless.up-lm3w-full-deployment-vector.v1"

type UPLM3WMetric struct{
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
type UPLM3WSummary struct{
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3WResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3VSeal string `json:"source_up_lm3v_seal"`
	ResourceConfigurations int `json:"resource_configurations"`
	GlobalOutcomeRange int `json:"global_outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3WMetric `json:"metrics"`
	Summaries []UPLM3WSummary `json:"summaries"`
}
func uplm3wKey(coord string,m UPLM3WMetric)string{
	if coord=="reachable_budget+capacity_by_round3+capacity_by_round4"{
		return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
	}
	return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d|c5=%d",m.ReachableBudget,m.CapacityByRound2,m.CapacityByRound3,m.CapacityByRound4,m.CapacityByRound5)
}
func RunUPLM3W()(UPLM3WResult,error){
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"reachable_budget+capacity_by_round3+capacity_by_round4","reachable_budget+capacity_by_round2+capacity_by_round3+capacity_by_round4+capacity_by_round5"}
	res:=UPLM3WResult{Schema:UPLM3WFullVectorSchema,Experiment:"UP-LM3W-full-deployment-vector",SourceUPLM3VSeal:"66a6e40d82e43973f7e715ef5a4ccee4fa7cf886",ResourceConfigurations:80,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCoordinateSearchUsed:false}
	vals:=[]int{}
	for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
		m:=UPLM3WMetric{Budget:b,ActionStartRound:s,Throughput:tp,ReachableBudget:uplm3fReachable(b,s,tp),CapacityByRound2:uplm3mCapByRound(b,s,tp,2),CapacityByRound3:uplm3mCapByRound(b,s,tp,3),CapacityByRound4:uplm3mCapByRound(b,s,tp,4),CapacityByRound5:uplm3mCapByRound(b,s,tp,5)}
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
			k:=uplm3wKey(coord,m);a:=g[k]
			if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
			if m.EarliestFailed<a.min{a.min=m.EarliestFailed}
			if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3WSummary{Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
