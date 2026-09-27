package unitary

const UPLM5HCollisionSchema="wingless.up-lm5h-residual-collision-trace.v1"

type UPLM5HMember struct{
	Budget int `json:"budget"`
	ActionStart int `json:"action_start"`
	Throughput int `json:"throughput"`
	BudgetCut int `json:"budget_cut"`
	BudgetRestore int `json:"budget_restore"`
	ThroughputCut int `json:"throughput_cut"`
	ThroughputRestore int `json:"throughput_restore"`
	EarlyResource string `json:"early_resource"`
	Reachable int `json:"reachable"`
	CoverageEnd int `json:"coverage_end"`
	PreEarly int `json:"pre_early"`
	PostEarly int `json:"post_early"`
	PreLate int `json:"pre_late"`
	PostLate int `json:"post_late"`
	EarlyRestoreRound int `json:"early_restore_round"`
	LateRestoreRound int `json:"late_restore_round"`
	Outcome int `json:"outcome"`
}
type UPLM5HCollision struct{
	DeadlineProfile string `json:"deadline_profile"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	CoordinateKey string `json:"coordinate_key"`
	MinOutcome int `json:"min_outcome"`
	MaxOutcome int `json:"max_outcome"`
	Members []UPLM5HMember `json:"members"`
}
type UPLM5HSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	PooledPoints int `json:"pooled_points"`
	CollisionGroups int `json:"collision_groups"`
	CollisionMembers int `json:"collision_members"`
	MaxFailureSpread int `json:"max_failure_spread"`
}
type UPLM5HResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5GSeal string `json:"source_up_lm5g_seal"`
	ResidualConditionCells int `json:"residual_condition_cells"`
	PointsPerConditionCell int `json:"points_per_condition_cell"`
	FrozenCoordinate string `json:"frozen_coordinate"`
	CoordinateModified bool `json:"coordinate_modified"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5HSummary `json:"summaries"`
	Collisions []UPLM5HCollision `json:"collisions"`
}
func RunUPLM5H()(UPLM5HResult,error){
	profiles:=[]string{"layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:1,ThroughputRestore:3},
		{BudgetCut:1,BudgetRestore:3,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:2,ThroughputRestore:4},
		{BudgetCut:2,BudgetRestore:4,ThroughputCut:0,ThroughputRestore:2},
		{BudgetCut:0,BudgetRestore:2,ThroughputCut:3,ThroughputRestore:5},
		{BudgetCut:3,BudgetRestore:5,ThroughputCut:0,ThroughputRestore:2},
	}
	treds:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	res:=UPLM5HResult{
		Schema:UPLM5HCollisionSchema,Experiment:"UP-LM5H-residual-collision-trace",
		SourceUPLM5GSeal:"ef2c4e50db44007c089c14ae928799c8fdbf1a3f",
		ResidualConditionCells:24,PointsPerConditionCell:384,
		FrozenCoordinate:"timing-aware+early-resource+post-late",
		CoordinateModified:false,AdaptiveFeatureSelectionUsed:false,CounterfactualOnly:true,LiveActivation:false,
	}
	type point struct{key string;member UPLM5HMember}
	type group struct{min,max int;members []UPLM5HMember}
	for _,profile:=range profiles{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		pts:=make([]point,0,384)
		for _,w:=range windows{for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,preE,postE,preL,postL,earlyR,lateR,earlyResource:=uplm5gSchedule(b,st,tp,w,2,tr)
			out:=uplm4wRun(rot,perm,profile,b,st,tp,w,2,tr)
			key:=uplm5gKey("timing-aware+early-resource+post-late",r,tp,end,preE,postE,preL,postL,earlyR,lateR,earlyResource)
			m:=UPLM5HMember{Budget:b,ActionStart:st,Throughput:tp,BudgetCut:w.BudgetCut,BudgetRestore:w.BudgetRestore,ThroughputCut:w.ThroughputCut,ThroughputRestore:w.ThroughputRestore,EarlyResource:earlyResource,Reachable:r,CoverageEnd:end,PreEarly:preE,PostEarly:postE,PreLate:preL,PostLate:postL,EarlyRestoreRound:earlyR,LateRestoreRound:lateR,Outcome:out}
			pts=append(pts,point{key:key,member:m})
		}}}}
		groups:=map[string]*group{}
		for _,p:=range pts{
			g:=groups[p.key]
			if g==nil{g=&group{min:p.member.Outcome,max:p.member.Outcome};groups[p.key]=g}
			if p.member.Outcome<g.min{g.min=p.member.Outcome};if p.member.Outcome>g.max{g.max=p.member.Outcome}
			g.members=append(g.members,p.member)
		}
		s:=UPLM5HSummary{DeadlineProfile:profile,ThroughputReduction:tr,Rotation:rot,Permutation:perm,PooledPoints:len(pts)}
		for key,g:=range groups{
			if g.max==g.min{continue}
			s.CollisionGroups++;s.CollisionMembers+=len(g.members)
			if d:=g.max-g.min;d>s.MaxFailureSpread{s.MaxFailureSpread=d}
			res.Collisions=append(res.Collisions,UPLM5HCollision{DeadlineProfile:profile,ThroughputReduction:tr,Rotation:rot,Permutation:perm,CoordinateKey:key,MinOutcome:g.min,MaxOutcome:g.max,Members:g.members})
		}
		res.Summaries=append(res.Summaries,s)
	}}}}
	return res,nil
}
