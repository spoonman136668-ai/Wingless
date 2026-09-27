package unitary

import "fmt"

const UPLM3TDeadlineProfileSchema="wingless.up-lm3t-deadline-profile-transfer.v1"

type UPLM3TMetric struct{
	DeadlineProfile string `json:"deadline_profile"`
	Budget int `json:"budget"`
	ActionStartRound int `json:"action_start_round"`
	Throughput int `json:"throughput"`
	ReachableBudget int `json:"reachable_budget"`
	CapacityByRound3 int `json:"capacity_by_round3"`
	CapacityByRound4 int `json:"capacity_by_round4"`
	EarliestFailed int `json:"earliest_failed"`
}
type UPLM3TCoordinateSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM3TProfileSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	OutcomeRange int `json:"outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
}
type UPLM3TResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3SSeal string `json:"source_up_lm3s_seal"`
	PressureLevel int `json:"pressure_level"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveProfileSelectionUsed bool `json:"adaptive_profile_selection_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	Metrics []UPLM3TMetric `json:"metrics"`
	CoordinateSummaries []UPLM3TCoordinateSummary `json:"coordinate_summaries"`
	ProfileSummaries []UPLM3TProfileSummary `json:"profile_summaries"`
}
func uplm3tTarget(a uplm2xArm,profile string)int{
	if profile=="deferred_only"{return uplm3hDeferredTarget(a)}
	if profile=="layout_only"{return uplm3hLayoutTarget(a)}
	return uplm3hTarget(a)
}
func uplm3tPrepressure(arms []uplm2xArm,rot int,profile string){
	for i:=range arms{
		target:=uplm3tTarget(arms[i],profile)
		for n:=0;n<32;n++{
			_,cd,ok:=uplm2xFirstPending(&arms[i])
			if !ok||cd<=target{break}
			arms[i].r.write(fmt.Sprintf("3t-pre-%s-%d-%d-%d",profile,rot,i,n),"x")
		}
	}
}
func uplm3tRun(rot int,perm,profile string,budget,start,tp int)(completed,failed,actions int){
	arms:=uplm2yPermute(uplm2xArms(rot),perm)
	uplm3tPrepressure(arms,rot,profile)
	for w:=0;w<4;w++{for i:=range arms{arms[i].r.write(fmt.Sprintf("3t-static-%s-%d-%s-%d-%d-%d-%d",profile,w,perm,budget,start,tp,i),"x")}}
	for round:=0;round<6;round++{
		if round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used)
				if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("3t-%s-%s-%d-%d-%d-%d-%d",profile,perm,budget,start,tp,round,i),"x")}
	}
	for i:=range arms{c,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);completed+=c;failed+=f}
	return
}
func uplm3tKey(coord string,m UPLM3TMetric)string{
	if coord=="reachable_budget+capacity_by_round4"{return fmt.Sprintf("r=%d|c4=%d",m.ReachableBudget,m.CapacityByRound4)}
	return fmt.Sprintf("r=%d|c3=%d|c4=%d",m.ReachableBudget,m.CapacityByRound3,m.CapacityByRound4)
}
func RunUPLM3T()(UPLM3TResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	coords:=[]string{"reachable_budget+capacity_by_round4","reachable_budget+capacity_by_round3+capacity_by_round4"}
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3TResult{Schema:UPLM3TDeadlineProfileSchema,Experiment:"UP-LM3T-deadline-profile-transfer",SourceUPLM3SSeal:"26b995d12ffcb92afbcbcf637ac08a4ca846d277",PressureLevel:4,DeadlineProfiles:profiles,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveProfileSelectionUsed:false,AdaptiveCoordinateSearchUsed:false}
	for _,profile:=range profiles{for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
		m:=UPLM3TMetric{DeadlineProfile:profile,Budget:budget,ActionStartRound:start,Throughput:tp,ReachableBudget:uplm3fReachable(budget,start,tp),CapacityByRound3:uplm3mCapByRound(budget,start,tp,3),CapacityByRound4:uplm3mCapByRound(budget,start,tp,4)}
		for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3tRun(rot,perm,profile,budget,start,tp);m.EarliestFailed+=ef}}
		res.Metrics=append(res.Metrics,m)
	}}}}
	type acc struct{min,max int}
	for _,profile:=range profiles{
		vals:=[]int{}
		for _,m:=range res.Metrics{if m.DeadlineProfile==profile{vals=append(vals,m.EarliestFailed)}}
		minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
		for _,v:=range vals{distinct[v]=true;if v<minv{minv=v};if v>maxv{maxv=v}}
		res.ProfileSummaries=append(res.ProfileSummaries,UPLM3TProfileSummary{DeadlineProfile:profile,MinEarliestFailed:minv,MaxEarliestFailed:maxv,OutcomeRange:maxv-minv,DistinctOutcomes:len(distinct)})
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,m:=range res.Metrics{
				if m.DeadlineProfile!=profile{continue}
				k:=uplm3tKey(coord,m);a:=g[k]
				if a==nil{a=&acc{min:m.EarliestFailed,max:m.EarliestFailed};g[k]=a}
				if m.EarliestFailed<a.min{a.min=m.EarliestFailed}
				if m.EarliestFailed>a.max{a.max=m.EarliestFailed}
			}
			s:=UPLM3TCoordinateSummary{DeadlineProfile:profile,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{s.NonzeroSpreadGroups++};if d>s.MaxFailureSpread{s.MaxFailureSpread=d}}
			if len(g)>0{s.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.CoordinateSummaries=append(res.CoordinateSummaries,s)
		}
	}
	return res,nil
}
