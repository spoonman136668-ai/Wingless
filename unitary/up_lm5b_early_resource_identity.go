package unitary

import "fmt"

const UPLM5BEarlyIdentitySchema="wingless.up-lm5b-early-resource-identity.v1"

type UPLM5BSummary struct{
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
type UPLM5BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5ASeal string `json:"source_up_lm5a_seal"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	RestoreRoundLabelsInKey bool `json:"restore_round_labels_in_key"`
	FullScheduleIdentityInKey bool `json:"full_schedule_identity_in_key"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveIdentityEncodingUsed bool `json:"adaptive_identity_encoding_used"`
	Summaries []UPLM5BSummary `json:"summaries"`
}
func uplm5bKey(coord string,r,tp,end,pe,preEarly,preLate int,earlyResource string)string{
	if coord=="ordered-path"{
		return fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d",r,tp,end,preEarly,pe,preLate)
	}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d|early=%s",r,tp,end,preEarly,pe,preLate,earlyResource)
}
func RunUPLM5B()(UPLM5BResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"ordered-path","ordered-path+early-resource"}
	res:=UPLM5BResult{Schema:UPLM5BEarlyIdentitySchema,Experiment:"UP-LM5B-early-resource-identity",SourceUPLM5ASeal:"610210b0fe8c4ac1e6af2ef1f2e160376c629107",PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,RestoreRoundLabelsInKey:false,FullScheduleIdentityInKey:false,CounterfactualOnly:true,LiveActivation:false,AdaptiveIdentityEncodingUsed:false}
	type point struct{r,tp,end,pe,preE,preL,out int;early string}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			preE,preL:=pb,pt
			early:="budget_first"
			if w.ThroughputRestore<w.BudgetRestore{preE,preL=pt,pb;early="throughput_first"}
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r:r,tp:tp,end:end,pe:pe,preE:preE,preL:preL,out:out,early:early})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5bKey(coord,p.r,p.tp,p.end,p.pe,p.preE,p.preL,p.early);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM5BSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
