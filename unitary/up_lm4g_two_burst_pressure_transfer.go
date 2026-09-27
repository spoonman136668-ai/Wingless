package unitary

import "fmt"

const UPLM4GTwoBurstSchema="wingless.up-lm4g-two-burst-pressure-transfer.v1"

type UPLM4GPair struct{
	RoundA int `json:"round_a"`
	RoundB int `json:"round_b"`
}
type UPLM4GSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BurstRoundA int `json:"burst_round_a"`
	BurstRoundB int `json:"burst_round_b"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4GResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4FSeal string `json:"source_up_lm4f_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	BurstPairs []UPLM4GPair `json:"burst_pairs"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveBurstPairUsed bool `json:"adaptive_burst_pair_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4GSummary `json:"summaries"`
}
func uplm4gRun(rot int,perm,profile string,budget,start,tp,a,b int)int{
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
			arms[i].r.write(fmt.Sprintf("4g-%s-%d-%d-%s-%d-%d-%d-%d",profile,a,b,perm,budget,start,tp,round),"x")
			if round==a||round==b{
				arms[i].r.write(fmt.Sprintf("4g-burst-%s-%d-%d-%s-%d-%d-%d-%d-%d",profile,a,b,perm,budget,start,tp,round,i),"x")
			}
		}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func RunUPLM4G()(UPLM4GResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pairs:=[]UPLM4GPair{{0,1},{0,2},{0,3},{0,4},{0,5},{1,2},{1,3},{1,4},{1,5},{2,3},{2,4},{2,5},{3,4},{3,5},{4,5}}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"r+tp+end","r+c2+c3+c4"}
	res:=UPLM4GResult{Schema:UPLM4GTwoBurstSchema,Experiment:"UP-LM4G-two-burst-pressure-transfer",SourceUPLM4FSeal:"0f9c211e92c1eb784c4d93be06cef2bf89e718f1",DeadlineProfiles:profiles,BurstPairs:pairs,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:80,TopologyConditionCells:len(profiles)*len(pairs)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveBurstPairUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{r,tp,end,c2,c3,c4,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,pair:=range pairs{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			points=append(points,point{
				r:uplm3fReachable(budget,start,tp),tp:tp,end:uplm4cCoverageEnd(budget,start,tp),
				c2:uplm3mCapByRound(budget,start,tp,2),c3:uplm3mCapByRound(budget,start,tp,3),c4:uplm3mCapByRound(budget,start,tp,4),
				out:uplm4gRun(rot,perm,profile,budget,start,tp,pair.RoundA,pair.RoundB),
			})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4fKey(coord,p.r,p.tp,p.end,p.c2,p.c3,p.c4);x:=g[k]
				if x==nil{x=&acc{min:p.out,max:p.out};g[k]=x}
				x.n++;if p.out<x.min{x.min=p.out};if p.out>x.max{x.max=p.out}
			}
			sm:=UPLM4GSummary{DeadlineProfile:profile,BurstRoundA:pair.RoundA,BurstRoundB:pair.RoundB,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,x:=range g{
				if x.n>1{sm.MultiMemberGroups++}
				d:=x.max-x.min;sum+=d
				if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}
			}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
