package unitary

import "fmt"

const UPLM5FRestoreIdentitySchema="wingless.up-lm5f-restore-time-resource-identity.v1"

type UPLM5FSummary struct{
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
type UPLM5FResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5ESeal string `json:"source_up_lm5e_seal"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	PointsPerPooledCell int `json:"points_per_pooled_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CutRoundLabelsInKeys bool `json:"cut_round_labels_in_keys"`
	FullScheduleIdentityInKeys bool `json:"full_schedule_identity_in_keys"`
	AdaptiveIdentityEncodingUsed bool `json:"adaptive_identity_encoding_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5FSummary `json:"summaries"`
}
func uplm5fKey(coord string,r,tp,end,pe,preE,preL,earlyR,lateR int,earlyResource string)string{
	base:=fmt.Sprintf("r=%d|tp=%d|end=%d|preE=%d|postE=%d|preL=%d|earlyR=%d|lateR=%d",r,tp,end,preE,pe,preL,earlyR,lateR)
	if coord=="timing-aware"{return base}
	return base+"|earlyResource="+earlyResource
}
func RunUPLM5F()(UPLM5FResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:1,ThroughputRestore:3},
		{BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:2,ThroughputRestore:4},
		{BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:3,ThroughputRestore:5},
		{BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:2},
	}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"timing-aware","timing-aware+early-resource"}
	res:=UPLM5FResult{
		Schema:UPLM5FRestoreIdentitySchema,Experiment:"UP-LM5F-restore-time-resource-identity",
		SourceUPLM5ESeal:"a41f2bb21374392e21fdc2c1aa37f0346307a243",
		PooledConditionCells:72,PointsPerPooledCell:384,CandidateCoordinates:coords,
		CutRoundLabelsInKeys:false,FullScheduleIdentityInKeys:false,AdaptiveIdentityEncodingUsed:false,
		CounterfactualOnly:true,LiveActivation:false,
	}
	type point struct{r,tp,end,pe,preE,preL,earlyR,lateR,out int;earlyResource string}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pe:=uplm4ySchedule(b,st,tp,w,br,tr)
			preE,preL:=pb,pt;earlyR,lateR:=w.BudgetRestore,w.ThroughputRestore;earlyResource:="budget_first"
			if w.ThroughputRestore<w.BudgetRestore{
				preE,preL=pt,pb;earlyR,lateR=w.ThroughputRestore,w.BudgetRestore;earlyResource="throughput_first"
			}
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)
			pts=append(pts,point{r:r,tp:tp,end:end,pe:pe,preE:preE,preL:preL,earlyR:earlyR,lateR:lateR,out:out,earlyResource:earlyResource})
		}}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm5fKey(coord,p.r,p.tp,p.end,p.pe,p.preE,p.preL,p.earlyR,p.lateR,p.earlyResource)
				a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM5FSummary{DeadlineProfile:profile,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,PooledPoints:len(pts),Groups:len(g)}
			sum:=0
			for _,a:=range g{
				if a.n>1{sm.MultiMemberGroups++}
				d:=a.max-a.min;sum+=d
				if d>0{sm.NonzeroSpreadGroups++}
				if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}
			}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
