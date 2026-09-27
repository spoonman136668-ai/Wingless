package unitary

import "fmt"

const UPLM4LTwoOutageSchema="wingless.up-lm4l-two-round-outage-transfer.v1"

type UPLM4LWindow struct{
	StartRound int `json:"start_round"`
	EndRound int `json:"end_round"`
}
type UPLM4LSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	OutageStart int `json:"outage_start"`
	OutageEnd int `json:"outage_end"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4LResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4KSeal string `json:"source_up_lm4k_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	OutageWindows []UPLM4LWindow `json:"outage_windows"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveOutageUsed bool `json:"adaptive_outage_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4LSummary `json:"summaries"`
}
func uplm4lRate(tp,round,start,end int)int{
	if round>=start&&round<=end{return 0}
	return tp
}
func uplm4lDynamicReachable(budget,startRound,tp,outStart,outEnd int)int{
	capacity:=0
	for round:=startRound;round<6;round++{capacity+=uplm4lRate(tp,round,outStart,outEnd)}
	if budget<capacity{return budget}
	return capacity
}
func uplm4lDynamicEnd(budget,startRound,tp,outStart,outEnd int)int{
	r:=uplm4lDynamicReachable(budget,startRound,tp,outStart,outEnd)
	if r<=0{return startRound-1}
	cum:=0
	for round:=startRound;round<6;round++{
		cum+=uplm4lRate(tp,round,outStart,outEnd)
		if cum>=r{return round}
	}
	return 5
}
func uplm4lRun(rot int,perm,profile string,budget,startRound,tp,outStart,outEnd int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=startRound{
			rate:=uplm4lRate(tp,round,outStart,outEnd);used:=map[int]bool{}
			for k:=0;k<rate&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4l-%s-%d-%d-%s-%d-%d-%d-%d",profile,outStart,outEnd,perm,budget,startRound,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4lKey(coord string,rs,tp,es,rd,ed int)string{
	if coord=="static_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",rs,tp,es)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d",rd,tp,ed)
}
func RunUPLM4L()(UPLM4LResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4LWindow{{2,3},{3,4},{4,5}}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"static_r+tp+end","dynamic_r+tp+end"}
	res:=UPLM4LResult{Schema:UPLM4LTwoOutageSchema,Experiment:"UP-LM4L-two-round-outage-transfer",SourceUPLM4KSeal:"3ae95dcc94af0a6399fd5a933de277b12a169dc1",DeadlineProfiles:profiles,OutageWindows:windows,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(windows)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveOutageUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{rs,tp,es,rd,ed,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,w:=range windows{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			points=append(points,point{rs:uplm3fReachable(budget,start,tp),tp:tp,es:uplm4cCoverageEnd(budget,start,tp),rd:uplm4lDynamicReachable(budget,start,tp,w.StartRound,w.EndRound),ed:uplm4lDynamicEnd(budget,start,tp,w.StartRound,w.EndRound),out:uplm4lRun(rot,perm,profile,budget,start,tp,w.StartRound,w.EndRound)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4lKey(coord,p.rs,p.tp,p.es,p.rd,p.ed);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4LSummary{DeadlineProfile:profile,OutageStart:w.StartRound,OutageEnd:w.EndRound,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
