package unitary

import "fmt"

const UPLM5AOrderedPoolingSchema="wingless.up-lm5a-ordered-path-pooling.v1"

type UPLM5ASummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	PooledPoints int `json:"pooled_points"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM5AResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4ZSeal string `json:"source_up_lm4z_seal"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	ScheduleIdentityInKey bool `json:"schedule_identity_in_key"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptivePoolingUsed bool `json:"adaptive_pooling_used"`
	Summaries []UPLM5ASummary `json:"summaries"`
}
func uplm5aKey(coord string,r,tp,end,pb,pt,pe,preEarly,preLate int)string{
	if coord=="resource-specific-full"{
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d|pe=%d",r,tp,end,pb,pt,pe)
	}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d",r,tp,end,preEarly,pe,preLate)
}
func RunUPLM5A()(UPLM5AResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"resource-specific-full","ordered-path"}
	res:=UPLM5AResult{Schema:UPLM5AOrderedPoolingSchema,Experiment:"UP-LM5A-ordered-path-pooling",SourceUPLM4ZSeal:"c38dc038c49a397906a4cd6482783760dd010d9c",PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,ScheduleIdentityInKey:false,CounterfactualOnly:true,LiveActivation:false,AdaptivePoolingUsed:false}
	type point struct{r,tp,end,pb,pt,pe,preE,preL,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			preE,preL:=pb,pt
			if w.ThroughputRestore<w.BudgetRestore{preE,preL=pt,pb}
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r,tp,end,pb,pt,pe,preE,preL,out})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5aKey(coord,p.r,p.tp,p.end,p.pb,p.pt,p.pe,p.preE,p.preL);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM5ASummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
