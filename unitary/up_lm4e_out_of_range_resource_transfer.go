package unitary

import "fmt"

const UPLM4EOutOfRangeSchema="wingless.up-lm4e-out-of-range-resource-transfer.v1"

type UPLM4ESummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4EResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4DSeal string `json:"source_up_lm4d_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PressureLevels []int `json:"pressure_levels"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	Budgets []int `json:"budgets"`
	ActionStartRounds []int `json:"action_start_rounds"`
	Throughputs []int `json:"throughputs"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4ESummary `json:"summaries"`
}
func uplm4eKey(coord string,b,start,tp,r,c2,c3,c4,end int)string{
	if coord=="r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",r,c2,c3,c4)
}
func RunUPLM4E()(UPLM4EResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pressures:=[]int{0,4,8};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{2,3,4,5,6,7,8};starts:=[]int{1,2,3,4,5,6};tps:=[]int{1,2,3,4,5}
	coords:=[]string{"r+tp+end","r+c2+c3+c4"}
	res:=UPLM4EResult{Schema:UPLM4EOutOfRangeSchema,Experiment:"UP-LM4E-out-of-range-resource-transfer",SourceUPLM4DSeal:"707f192a9cd6838a467161a1b2967d75562fa7e0",DeadlineProfiles:profiles,PressureLevels:pressures,Rotations:rots,Permutations:perms,Budgets:budgets,ActionStartRounds:starts,Throughputs:tps,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(pressures)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{b,start,tp,r,c2,c3,c4,end,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,pressure:=range pressures{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			points=append(points,point{b:b,start:s,tp:tp,r:uplm3fReachable(b,s,tp),c2:uplm3mCapByRound(b,s,tp,2),c3:uplm3mCapByRound(b,s,tp,3),c4:uplm3mCapByRound(b,s,tp,4),end:uplm4cCoverageEnd(b,s,tp),out:uplm3uRun(rot,perm,profile,b,s,tp,pressure)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4eKey(coord,p.b,p.start,p.tp,p.r,p.c2,p.c3,p.c4,p.end);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4ESummary{DeadlineProfile:profile,PressureLevel:pressure,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
