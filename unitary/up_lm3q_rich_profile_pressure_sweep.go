package unitary

import "fmt"

const UPLM3QRichProfileSchema="wingless.up-lm3q-rich-profile-pressure-sweep.v1"

type UPLM3QMetric struct{
	PressureLevel int `json:"pressure_level"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3QSummary struct{
	PressureLevel int `json:"pressure_level"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3QResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3PSeal string `json:"source_up_lm3p_seal"`
	PressureLevels []int `json:"pressure_levels"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptivePressureLevelsUsed bool `json:"adaptive_pressure_levels_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3QMetric `json:"metrics"`
	Summaries []UPLM3QSummary `json:"summaries"`
}
func uplm3qKey(coord string,m UPLM3QMetric)string{
	if coord=="reachable_budget+capacity_by_round4"{
		return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)
	}
	return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
}
func RunUPLM3Q()(UPLM3QResult,error){
	pressures:=[]int{4,5,6,7,8}
	coords:=[]string{"reachable_budget+capacity_by_round4","reachable_budget+capacity_by_round3+capacity_by_round4"}
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3QResult{Schema:UPLM3QRichProfileSchema,Experiment:"UP-LM3Q-rich-profile-pressure-sweep",SourceUPLM3PSeal:"9edc932ad23072bc0cd5bf99586a50718f2a9e26",PressureLevels:pressures,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptivePressureLevelsUsed:false,AdaptiveCoordinateSearchUsed:false}
	for _,pressure:=range pressures{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3QMetric{PressureLevel:pressure,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp),CapacityByRound3:uplm3mCapByRound(budget,start,tp,3),CapacityByRound4:uplm3mCapByRound(budget,start,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3oRun(rot,perm,"earliest_deadline",budget,start,tp,pressure);m.EarliestFailed+=ef}}
		res.Metrics=append(res.Metrics,m)
	}}}}
	type acc struct{min,max int}
	for _,pressure:=range pressures{for _,coord:=range coords{
		g:=map[string]*acc{}
		for _,m:=range res.Metrics{
			if m.PressureLevel!=pressure{continue}
			k:=uplm3qKey(coord,m);a:=g[k]
			if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
			if m.EarliestFailed<a.min{a.min=m.EarliestFailed}
			if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3QSummary{PressureLevel:pressure,Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}}
	return res,nil
}
