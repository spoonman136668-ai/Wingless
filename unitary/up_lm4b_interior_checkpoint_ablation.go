package unitary

import "fmt"

const UPLM4BInteriorSchema="wingless.up-lm4b-interior-checkpoint-ablation.v1"

type UPLM4BSummary struct{
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
type UPLM4BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4ASeal string `json:"source_up_lm4a_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4BSummary `json:"summaries"`
}
func uplm4bKey(coord string,r,c2,c3,c4 int)string{
	switch coord{
	case "r+c2+c3":
		return fmt.Sprintf("r=%d|c2=%d|c3=%d",r,c2,c3)
	case "r+c2+c4":
		return fmt.Sprintf("r=%d|c2=%d|c4=%d",r,c2,c4)
	default:
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",r,c2,c3,c4)
	}
}
func RunUPLM4B()(UPLM4BResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pressures:=[]int{0,1,2,3,4,5,6,7,8}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"r+c2+c3","r+c2+c4","r+c2+c3+c4"}
	res:=UPLM4BResult{Schema:UPLM4BInteriorSchema,Experiment:"UP-LM4B-interior-checkpoint-ablation",SourceUPLM4ASeal:"ddeb4b19bbb9a904832242ad7c8a976187787829",TopologyConditionCells:162,ResourceConfigurationsPerCell:80,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,c2,c3,c4,out int}
	type acc struct{min,max int}
	for _,profile:=range profiles{for _,pressure:=range pressures{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			points=append(points,point{r:uplm3fReachable(b,s,tp),c2:uplm3mCapByRound(b,s,tp,2),c3:uplm3mCapByRound(b,s,tp,3),c4:uplm3mCapByRound(b,s,tp,4),out:uplm3uRun(rot,perm,profile,b,s,tp,pressure)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4bKey(coord,p.r,p.c2,p.c3,p.c4);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4BSummary{DeadlineProfile:profile,PressureLevel:pressure,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
