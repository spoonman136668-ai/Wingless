package unitary

import "fmt"

const UPLM4JThroughputDropSchema="wingless.up-lm4j-throughput-drop-transfer.v1"

type UPLM4JSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	DropRound int `json:"drop_round"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4JResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4ISeal string `json:"source_up_lm4i_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	DropRounds []int `json:"drop_rounds"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	ActionStarts []int `json:"action_starts"`
	InitialThroughputs []int `json:"initial_throughputs"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveDropUsed bool `json:"adaptive_drop_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4JSummary `json:"summaries"`
}
func uplm4jRate(tp,round,drop int)int{
	if round>=drop{return tp-1}
	return tp
}
func uplm4jDynamicReachable(budget,start,tp,drop int)int{
	capacity:=0
	for round:=start;round<6;round++{capacity+=uplm4jRate(tp,round,drop)}
	if budget<capacity{return budget}
	return capacity
}
func uplm4jDynamicEnd(budget,start,tp,drop int)int{
	r:=uplm4jDynamicReachable(budget,start,tp,drop)
	if r<=0{return start-1}
	cum:=0
	for round:=start;round<6;round++{
		cum+=uplm4jRate(tp,round,drop)
		if cum>=r{return round}
	}
	return 5
}
func uplm4jRun(rot int,perm,profile string,budget,start,tp,drop int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			rate:=uplm4jRate(tp,round,drop);used:=map[int]bool{}
			for k:=0;k<rate&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4j-%s-%d-%s-%d-%d-%d-%d",profile,drop,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4jKey(coord string,rs,tp,es,rd,post,ed int)string{
	if coord=="static_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",rs,tp,es)}
	return fmt.Sprintf("r=%d|tp0=%d|tp1=%d|end=%d",rd,tp,post,ed)
}
func RunUPLM4J()(UPLM4JResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	drops:=[]int{2,3,4,5};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4}
	coords:=[]string{"static_r+tp+end","dynamic_r+tp_before+tp_after+end"}
	res:=UPLM4JResult{Schema:UPLM4JThroughputDropSchema,Experiment:"UP-LM4J-throughput-drop-transfer",SourceUPLM4ISeal:"43523c3750ccd0a99df30564c0505b589e61467e",DeadlineProfiles:profiles,DropRounds:drops,Rotations:rots,Permutations:perms,Budgets:budgets,ActionStarts:starts,InitialThroughputs:tps,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(drops)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveDropUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{rs,tp,es,rd,post,ed,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,drop:=range drops{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			points=append(points,point{
				rs:uplm3fReachable(budget,start,tp),tp:tp,es:uplm4cCoverageEnd(budget,start,tp),
				rd:uplm4jDynamicReachable(budget,start,tp,drop),post:tp-1,ed:uplm4jDynamicEnd(budget,start,tp,drop),
				out:uplm4jRun(rot,perm,profile,budget,start,tp,drop),
			})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4jKey(coord,p.rs,p.tp,p.es,p.rd,p.post,p.ed);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4JSummary{DeadlineProfile:profile,DropRound:drop,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{
				if a.n>1{sm.MultiMemberGroups++}
				d:=a.max-a.min;sum+=d
				if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}
			}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
