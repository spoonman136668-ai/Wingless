package unitary

import "fmt"

const UPLM3LTemporalResourceSchema="wingless.up-lm3l-temporal-resource-coordinate.v1"

type UPLM3LMetric struct{
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	EarliestFailed int `json:"earliest_failed"`
	EarliestActions int `json:"earliest_actions"`
}
type UPLM3LCoordinateSummary struct{
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3LResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3KSeal string `json:"source_up_lm3k_seal"`
	PressureMode string `json:"pressure_mode"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveGroupingUsed bool `json:"adaptive_grouping_used"`
	AdaptiveResourcesUsed bool `json:"adaptive_resources_used"`
	Metrics []UPLM3LMetric `json:"metrics"`
	CoordinateSummaries []UPLM3LCoordinateSummary `json:"coordinate_summaries"`
}
type uplm3lAcc struct{n,min,max int}
func uplm3lKey(coord string,m UPLM3LMetric)string{
	switch coord{
	case "reachable_budget":return fmt.Sprintf("r=%d",m.ReachableBudget)
	case "reachable_budget+start":return fmt.Sprintf("r=%d|s=%d",m.ReachableBudget,m.ActionStartRound)
	case "reachable_budget+throughput":return fmt.Sprintf("r=%d|t=%d",m.ReachableBudget,m.Throughput)
	default:return fmt.Sprintf("r=%d|s=%d|t=%d",m.ReachableBudget,m.ActionStartRound,m.Throughput)
	}
}
func RunUPLM3L()(UPLM3LResult,error){
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"reachable_budget","reachable_budget+start","reachable_budget+throughput","reachable_budget+start+throughput"}
	res:=UPLM3LResult{Schema:UPLM3LTemporalResourceSchema,Experiment:"UP-LM3L-temporal-resource-coordinate",SourceUPLM3KSeal:"d077e4f8ee66b0d5a0b9c55f7b3108ec0290569d",PressureMode:"static_pre",Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveGroupingUsed:false,AdaptiveResourcesUsed:false}
	for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3LMetric{Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp)}
		for _,rot:=range rots{for _,perm:=range perms{
			_,ef,ea:=uplm3kRun(rot,perm,"earliest_deadline","static_pre",budget,start,tp)
			m.EarliestFailed+=ef;m.EarliestActions+=ea
		}}
		res.Metrics=append(res.Metrics,m)
	}}}
	for _,coord:=range coords{
		groups:=map[string]*uplm3lAcc{}
		for _,m:=range res.Metrics{
			k:=uplm3lKey(coord,m);a:=groups[k]
			if a==nil{a=&uplm3lAcc{min:m.EarliestFailed,max:m.EarliestFailed};groups[k]=a}
			a.n++;if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3LCoordinateSummary{Coordinate:coord,Groups:len(groups)};spreadSum:=0
		for _,a:=range groups{spread:=a.max-a.min;spreadSum+=spread;if spread>0{s.NonzeroSpreadGroups++};if spread>s.MaxFailureSpread{s.MaxFailureSpread=spread}}
		if len(groups)>0{s.MeanFailureSpread=float64(spreadSum)/float64(len(groups))}
		res.CoordinateSummaries=append(res.CoordinateSummaries,s)
	}
	return res,nil
}
