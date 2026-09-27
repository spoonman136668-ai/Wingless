package unitary

import "fmt"

const UPLM3NDoubleStaticSchema="wingless.up-lm3n-double-static-pressure.v1"

type UPLM3NMetric struct{
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3NSummary struct{
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3NResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3MSeal string `json:"source_up_lm3m_seal"`
	StaticExtraWritesPerArm int `json:"static_extra_writes_per_arm"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	AdaptivePressureUsed bool `json:"adaptive_pressure_used"`
	Metrics []UPLM3NMetric `json:"metrics"`
	Summaries []UPLM3NSummary `json:"summaries"`
}
func uplm3nRun(rot int,perm,policy string,budget,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3hPrepressure(arms,rot)
	for i:=range arms{
		arms[i].r.write(fmt.Sprintf("3n-pre1-%s-%s-%d-%d-%d-%d",policy,perm,budget,start,tp,i),"x")
		arms[i].r.write(fmt.Sprintf("3n-pre2-%s-%s-%d-%d-%d-%d",policy,perm,budget,start,tp,i),"x")
	}
	for round:=0;round<6;round++{
		if policy!="baseline"&&round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,policy,used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3n-%s-%s-%d-%d-%d-%d-%d",policy,perm,budget,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func uplm3nKey(coord string,m UPLM3NMetric)string{
	if coord=="reachable_budget"{return fmt.Sprintf("r=%d",m.ReachableBudget)}
	return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)
}
func RunUPLM3N()(UPLM3NResult,error){
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	coords:=[]string{"reachable_budget","reachable_budget+capacity_by_round4"}
	res:=UPLM3NResult{Schema:UPLM3NDoubleStaticSchema,Experiment:"UP-LM3N-double-static-pressure",SourceUPLM3MSeal:"d067be9ceea97f83c61bec184767a56494e8fa33",StaticExtraWritesPerArm:2,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCoordinateSearchUsed:false,AdaptivePressureUsed:false}
	for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3NMetric{Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp),CapacityByRound4:uplm3mCapByRound(budget,start,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3nRun(rot,perm,"earliest_deadline",budget,start,tp);m.EarliestFailed+=ef}}
		res.Metrics=append(res.Metrics,m)
	}}}
	type acc struct{min,max int}
	for _,coord:=range coords{
		g:=map[string]*acc{}
		for _,m:=range res.Metrics{k:=uplm3nKey(coord,m);a:=g[k];if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a};if m.EarliestFailed<a.min{a.min=m.EarliestFailed};if m.EarliestFailed>a.max{a.max=m.EarliestFailed}}
		s:=UPLM3NSummary{Coordinate:coord,Groups:len(g)};sum:=0
		for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
		if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
