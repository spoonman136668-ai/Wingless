package unitary

import "fmt"

const UPLM3OPressureSweepSchema="wingless.up-lm3o-pressure-level-sweep.v1"

type UPLM3OMetric struct{
	PressureLevel int `json:"pressure_level"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3OSummary struct{
	PressureLevel int `json:"pressure_level"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3OResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3NSeal string `json:"source_up_lm3n_seal"`
	PressureLevels []int `json:"pressure_levels"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptivePressureLevelsUsed bool `json:"adaptive_pressure_levels_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3OMetric `json:"metrics"`
	Summaries []UPLM3OSummary `json:"summaries"`
}
func uplm3oRun(rot int,perm,policy string,budget,start,tp,pressure int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3hPrepressure(arms,rot)
	for w:=0;w<pressure;w++{
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3o-pre-%d-%d-%s-%s-%d-%d-%d-%d",pressure,w,policy,perm,budget,start,tp,i),"x")}
	}
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3o-%d-%s-%s-%d-%d-%d-%d-%d",pressure,policy,perm,budget,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func uplm3oKey(coord string,m UPLM3OMetric)string{
	if coord=="reachable_budget"{return fmt.Sprintf("r=%d",m.ReachableBudget)}
	return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)
}
func RunUPLM3O()(UPLM3OResult,error){
	pressures:=[]int{0,1,2,3,4};coords:=[]string{"reachable_budget","reachable_budget+capacity_by_round4"}
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3OResult{Schema:UPLM3OPressureSweepSchema,Experiment:"UP-LM3O-pressure-level-sweep",SourceUPLM3NSeal:"a99020004253810d937b5e95e12c558fb5643824",PressureLevels:pressures,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptivePressureLevelsUsed:false,AdaptiveCoordinateSearchUsed:false}
	for _,pressure:=range pressures{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3OMetric{PressureLevel:pressure,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp),CapacityByRound4:uplm3mCapByRound(budget,start,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3oRun(rot,perm,"earliest_deadline",budget,start,tp,pressure);m.EarliestFailed+=ef}}
		res.Metrics=append(res.Metrics,m)
	}}}}
	type acc struct{min,max int}
	for _,pressure:=range pressures{for _,coord:=range coords{
		g:=map[string]*acc{}
		for _,m:=range res.Metrics{
			if m.PressureLevel!=pressure{continue}
			k:=uplm3oKey(coord,m);a:=g[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a};if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
		}
		s:=UPLM3OSummary{PressureLevel:pressure,Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}}
	return res,nil
}
