package unitary

import "fmt"

const UPLM4FDynamicSchema="wingless.up-lm4f-dynamic-pressure-transfer.v1"

type UPLM4FSummary struct{
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
type UPLM4FResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4ESeal string `json:"source_up_lm4e_seal"`
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
	Summaries []UPLM4FSummary `json:"summaries"`
}
func uplm4fRun(rot int,perm,profile string,budget,start,tp,burst int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			used:=map[int]bool{}
			for k:=0;k<tp&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{
			arms[i].r.write(fmt.Sprintf("4f-%s-%d-%s-%d-%d-%d-%d",profile,burst,perm,budget,start,tp,round),"x")
			if round==burst{arms[i].r.write(fmt.Sprintf("4f-burst-%s-%d-%s-%d-%d-%d-%d",profile,burst,perm,budget,start,tp,i),"x")}
		}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4fKey(coord string,r,tp,end,c2,c3,c4 int)string{
	if coord=="r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",r,c2,c3,c4)
}
func RunUPLM4F()(UPLM4FResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	bursts:=[]int{0,1,2,3,4,5};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"r+tp+end","r+c2+c3+c4"}
	res:=UPLM4FResult{Schema:UPLM4FDynamicSchema,Experiment:"UP-LM4F-dynamic-pressure-transfer",SourceUPLM4ESeal:"81ca89b4c8eb786a513391e6dca377e17ecc1c3e",DeadlineProfiles:profiles,BurstRounds:bursts,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:80,TopologyConditionCells:len(profiles)*len(bursts)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveBurstUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{r,tp,end,c2,c3,c4,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,burst:=range bursts{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			points=append(points,point{r:uplm3fReachable(b,s,tp),tp:tp,end:uplm4cCoverageEnd(b,s,tp),c2:uplm3mCapByRound(b,s,tp,2),c3:uplm3mCapByRound(b,s,tp,3),c4:uplm3mCapByRound(b,s,tp,4),out:uplm4fRun(rot,perm,profile,b,s,tp,burst)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4fKey(coord,p.r,p.tp,p.end,p.c2,p.c3,p.c4);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4FSummary{DeadlineProfile:profile,BurstRound:burst,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
