package unitary

import "fmt"

const UPLM4KOutageSchema="wingless.up-lm4k-throughput-outage-transfer.v1"

type UPLM4KSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	OutageRound int `json:"outage_round"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4KResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4JSeal string `json:"source_up_lm4j_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	OutageRounds []int `json:"outage_rounds"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveOutageUsed bool `json:"adaptive_outage_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4KSummary `json:"summaries"`
}
func uplm4kRate(tp,round,outage int)int{
	if round==outage{return 0}
	return tp
}
func uplm4kDynamicReachable(budget,start,tp,outage int)int{
	capacity:=0
	for round:=start;round<6;round++{capacity+=uplm4kRate(tp,round,outage)}
	if budget<capacity{return budget}
	return capacity
}
func uplm4kDynamicEnd(budget,start,tp,outage int)int{
	r:=uplm4kDynamicReachable(budget,start,tp,outage)
	if r<=0{return start-1}
	cum:=0
	for round:=start;round<6;round++{
		cum+=uplm4kRate(tp,round,outage)
		if cum>=r{return round}
	}
	return 5
}
func uplm4kRun(rot int,perm,profile string,budget,start,tp,outage int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			rate:=uplm4kRate(tp,round,outage);used:=map[int]bool{}
			for k:=0;k<rate&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4k-%s-%d-%s-%d-%d-%d-%d",profile,outage,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4kKey(coord string,rs,tp,es,rd,ed int)string{
	if coord=="static_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",rs,tp,es)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d",rd,tp,ed)
}
func RunUPLM4K()(UPLM4KResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	outages:=[]int{2,3,4,5};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"static_r+tp+end","dynamic_r+tp+end"}
	res:=UPLM4KResult{Schema:UPLM4KOutageSchema,Experiment:"UP-LM4K-throughput-outage-transfer",SourceUPLM4JSeal:"232895f37c3a49ecb4b87a45d61ce2c9f9ce734e",DeadlineProfiles:profiles,OutageRounds:outages,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(outages)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveOutageUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{rs,tp,es,rd,ed,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,outage:=range outages{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			points=append(points,point{rs:uplm3fReachable(budget,start,tp),tp:tp,es:uplm4cCoverageEnd(budget,start,tp),rd:uplm4kDynamicReachable(budget,start,tp,outage),ed:uplm4kDynamicEnd(budget,start,tp,outage),out:uplm4kRun(rot,perm,profile,budget,start,tp,outage)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4kKey(coord,p.rs,p.tp,p.es,p.rd,p.ed);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4KSummary{DeadlineProfile:profile,OutageRound:outage,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
