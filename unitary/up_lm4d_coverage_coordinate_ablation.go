package unitary

import "fmt"

const UPLM4DCoverageAblationSchema="wingless.up-lm4d-coverage-coordinate-ablation.v1"

type UPLM4DSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4DResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4CSeal string `json:"source_up_lm4c_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4DSummary `json:"summaries"`
}
func uplm4dKey(coord string,r,c2,tp,end int)string{
	switch coord{
	case "c2+tp+end":
		return fmt.Sprintf("c2=%d|tp=%d|end=%d",c2,tp,end)
	case "r+tp+end":
		return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)
	case "r+c2+end":
		return fmt.Sprintf("r=%d|c2=%d|end=%d",r,c2,end)
	case "r+c2+tp":
		return fmt.Sprintf("r=%d|c2=%d|tp=%d",r,c2,tp)
	default:
		return fmt.Sprintf("r=%d|c2=%d|tp=%d|end=%d",r,c2,tp,end)
	}
}
func RunUPLM4D()(UPLM4DResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pressures:=[]int{0,1,2,3,4,5,6,7,8}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"c2+tp+end","r+tp+end","r+c2+end","r+c2+tp","r+c2+tp+end"}
	res:=UPLM4DResult{Schema:UPLM4DCoverageAblationSchema,Experiment:"UP-LM4D-coverage-coordinate-ablation",SourceUPLM4CSeal:"a74d3337d57aa7da9fc36dc8e6ce7958b90c68c8",TopologyConditionCells:162,ResourceConfigurationsPerCell:80,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,c2,tp,end,out int}
	type acc struct{min,max int}
	for _,profile:=range profiles{for _,pressure:=range pressures{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			points=append(points,point{r:uplm3fReachable(b,s,tp),c2:uplm3mCapByRound(b,s,tp,2),tp:tp,end:uplm4cCoverageEnd(b,s,tp),out:uplm3uRun(rot,perm,profile,b,s,tp,pressure)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4dKey(coord,p.r,p.c2,p.tp,p.end);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4DSummary{DeadlineProfile:profile,PressureLevel:pressure,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
