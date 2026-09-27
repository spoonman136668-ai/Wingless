package unitary

import "fmt"

const UPLM4STripleSchema="wingless.up-lm4s-triple-revoke-restore-transfer.v1"

type UPLM4STripleSchedule struct{Cut1,Restore1,Cut2,Restore2,Cut3,Restore3 int}
type UPLM4SSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	ScheduleIndex int `json:"schedule_index"`
	RevokeAmount int `json:"revoke_amount"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4SResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4RSeal string `json:"source_up_lm4r_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	Summaries []UPLM4SSummary `json:"summaries"`
}
func uplm4sCap(budget,round,revoke int,s UPLM4STripleSchedule)int{
	revoked:=(round>=s.Cut1&&round<s.Restore1)||(round>=s.Cut2&&round<s.Restore2)||(round>=s.Cut3&&round<s.Restore3)
	if !revoked{return budget};x:=budget-revoke;if x<0{return 0};return x
}
func uplm4sSchedule(budget,start,tp,revoke int,s UPLM4STripleSchedule)(r,end,p1,p2,p3 int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==s.Restore1{p1=actions};if round==s.Restore2{p2=actions};if round==s.Restore3{p3=actions}
		if round<start{continue}
		cap:=uplm4sCap(budget,round,revoke,s);before:=actions
		for k:=0;k<tp&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,p1,p2,p3
}
func uplm4sRun(rot int,perm,profile string,budget,start,tp,revoke int,s UPLM4STripleSchedule)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile);actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4sCap(budget,round,revoke,s);used:=map[int]bool{}
			for k:=0;k<tp&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4s-%s-%d-%d-%d-%d-%d-%d-%d-%s-%d-%d-%d-%d",profile,s.Cut1,s.Restore1,s.Cut2,s.Restore2,s.Cut3,s.Restore3,revoke,perm,budget,start,tp,round),"x")}
	}
	failed:=0;for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f};return failed
}
func uplm4sKey(coord string,r,tp,end,p1,p2,p3 int)string{
	if coord=="final"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	if coord=="final+latestrestore"{return fmt.Sprintf("r=%d|tp=%d|end=%d|p3=%d",r,tp,end,p3)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|p1=%d|p2=%d|p3=%d",r,tp,end,p1,p2,p3)
}
func RunUPLM4S()(UPLM4SResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4STripleSchedule{{0,1,2,3,4,5},{0,2,2,3,4,5},{0,1,2,4,4,5}}
	revokes:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"final","final+latestrestore","final+allrestores"}
	res:=UPLM4SResult{Schema:UPLM4STripleSchema,Experiment:"UP-LM4S-triple-revoke-restore-transfer",SourceUPLM4RSeal:"e3df4b476ef273a679da519c0d7632ccda09e4d3",TopologyConditionCells:108,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCheckpointSearchUsed:false}
	type point struct{r,tp,end,p1,p2,p3,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for si,s:=range schedules{for _,rev:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,p1,p2,p3:=uplm4sSchedule(b,st,tp,rev,s)
			pts=append(pts,point{r,tp,end,p1,p2,p3,uplm4sRun(rot,perm,profile,b,st,tp,rev,s)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{k:=uplm4sKey(coord,p.r,p.tp,p.end,p.p1,p.p2,p.p3);a:=g[k];if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a};a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}}
			sm:=UPLM4SSummary{DeadlineProfile:profile,ScheduleIndex:si,RevokeAmount:rev,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))};res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
