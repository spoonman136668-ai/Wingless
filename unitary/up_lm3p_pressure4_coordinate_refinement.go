package unitary

import "fmt"

const UPLM3PPressure4Schema="wingless.up-lm3p-pressure4-coordinate-refinement.v1"

type UPLM3PMetric struct{
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3PSummary struct{
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3PResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3OSeal string `json:"source_up_lm3o_seal"`
	PressureLevel int `json:"pressure_level"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3PMetric `json:"metrics"`
	Summaries []UPLM3PSummary `json:"summaries"`
}
func uplm3pKey(coord string,m UPLM3PMetric)string{
	switch coord{
	case "reachable_budget+capacity_by_round4":
		return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)
	case "reachable_budget+capacity_by_round3+capacity_by_round4":
		return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
	default:
		return fmt.Sprintf("r=%d|s=%d|t=%d",m.ReachableBudget,m.ActionStartRound,m.Throughput)
	}
}
func RunUPLM3P()(UPLM3PResult,error){
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"reachable_budget+capacity_by_round4","reachable_budget+capacity_by_round3+capacity_by_round4","reachable_budget+action_start_round+throughput"}
	res:=UPLM3PResult{Schema:UPLM3PPressure4Schema,Experiment:"UP-LM3P-pressure4-coordinate-refinement",SourceUPLM3OSeal:"fe054ae3e21b001c722caaf13a694559104e2e3a",PressureLevel:4,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCoordinateSearchUsed:false}
	for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3PMetric{Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp),CapacityByRound3:uplm3mCapByRound(budget,start,tp,3),CapacityByRound4:uplm3mCapByRound(budget,start,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3oRun(rot,perm,"earliest_deadline",budget,start,tp,4);m.EarliestFailed+=ef}}
		res.Metrics=append(res.Metrics,m)
	}}}
	type acc struct{min,max int}
	for _,coord:=range coords{
		g:=map[string]*acc{}
		for _,m:=range res.Metrics{
			k:=uplm3pKey(coord,m);a:=g[k]
			if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
			if m.EarliestFailed<a.min{a.min=m.EarliestFailed}
			if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3PSummary{Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
