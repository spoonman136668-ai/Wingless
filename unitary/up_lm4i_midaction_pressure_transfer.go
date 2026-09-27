package unitary

import "fmt"

const UPLM4IMidActionSchema="wingless.up-lm4i-midaction-pressure-transfer.v1"

type UPLM4ISummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BurstRound int `json:"burst_round"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4IResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4HSeal string `json:"source_up_lm4h_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	BurstRounds []int `json:"burst_rounds"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveBurstUsed bool `json:"adaptive_burst_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4ISummary `json:"summaries"`
}
func uplm4iRun(rot int,perm,profile string,budget,start,tp,burst int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		used:=map[int]bool{}
		if round>=start && actions<budget{
			i:=uplm3cChoose(arms,"earliest_deadline",used)
			if i>=0{
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}
			}
		}
		if round==burst{
			for i:=range arms{arms[i].r.write(fmt.Sprintf("4i-midburst-%s-%d-%s-%d-%d-%d-%d",profile,burst,perm,budget,start,tp,i),"x")}
		}
		if round>=start{
			for k:=1;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4i-%s-%d-%s-%d-%d-%d-%d",profile,burst,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func RunUPLM4I()(UPLM4IResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	bursts:=[]int{0,1,2,3,4,5};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"r+tp+end","r+c2+c3+c4"}
	res:=UPLM4IResult{Schema:UPLM4IMidActionSchema,Experiment:"UP-LM4I-midaction-pressure-transfer",SourceUPLM4HSeal:"83e3c6dfa0f15e3c9b465ee1f3f6729ba2aee0ee",DeadlineProfiles:profiles,BurstRounds:bursts,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:80,TopologyConditionCells:len(profiles)*len(bursts)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveBurstUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{r,tp,end,c2,c3,c4,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,burst:=range bursts{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			points=append(points,point{r:uplm3fReachable(budget,start,tp),tp:tp,end:uplm4cCoverageEnd(budget,start,tp),c2:uplm3mCapByRound(budget,start,tp,2),c3:uplm3mCapByRound(budget,start,tp,3),c4:uplm3mCapByRound(budget,start,tp,4),out:uplm4iRun(rot,perm,profile,budget,start,tp,burst)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4fKey(coord,p.r,p.tp,p.end,p.c2,p.c3,p.c4);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4ISummary{DeadlineProfile:profile,BurstRound:burst,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
