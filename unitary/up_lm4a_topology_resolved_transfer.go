package unitary

import "fmt"

const UPLM4ATopologyResolvedSchema="wingless.up-lm4a-topology-resolved-transfer.v1"

type UPLM4ASummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	PressureLevel int `json:"pressure_level"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
	GlobalOutcomeRange int `json:"global_outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
}
type UPLM4AResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3ZSeal string `json:"source_up_lm3z_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	PressureLevels []int `json:"pressure_levels"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4ASummary `json:"summaries"`
}
func uplm4aKey(coord string,r,c2,c3,c4,c5 int)string{
	switch coord{
	case "r+c3+c4":
		return fmt.Sprintf("r=%d|c3=%d|c4=%d",r,c3,c4)
	case "r+c2+c3+c4":
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d",r,c2,c3,c4)
	default:
		return fmt.Sprintf("r=%d|c2=%d|c3=%d|c4=%d|c5=%d",r,c2,c3,c4,c5)
	}
}
func RunUPLM4A()(UPLM4AResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	pressures:=[]int{0,1,2,3,4,5,6,7,8}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{3,4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"r+c3+c4","r+c2+c3+c4","r+c2+c3+c4+c5"}
	res:=UPLM4AResult{Schema:UPLM4ATopologyResolvedSchema,Experiment:"UP-LM4A-topology-resolved-transfer",SourceUPLM3ZSeal:"10578d7c383bd64cfd75afa9e873db3d6e9f6126",DeadlineProfiles:profiles,PressureLevels:pressures,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:80,TopologyConditionCells:len(profiles)*len(pressures)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,c2,c3,c4,c5,out int}
	type acc struct{min,max int}
	for _,profile:=range profiles{for _,pressure:=range pressures{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,80)
		minOut,maxOut:=0,0;distinct:=map[int]bool{};first:=true
		for _,b:=range budgets{for _,s:=range starts{for _,tp:=range tps{
			p:=point{r:uplm3fReachable(b,s,tp),c2:uplm3mCapByRound(b,s,tp,2),c3:uplm3mCapByRound(b,s,tp,3),c4:uplm3mCapByRound(b,s,tp,4),c5:uplm3mCapByRound(b,s,tp,5),out:uplm3uRun(rot,perm,profile,b,s,tp,pressure)}
			points=append(points,p);distinct[p.out]=true
			if first{minOut,maxOut=p.out,p.out;first=false}else{if p.out<minOut{minOut=p.out};if p.out>maxOut{maxOut=p.out}}
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4aKey(coord,p.r,p.c2,p.c3,p.c4,p.c5);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4ASummary{DeadlineProfile:profile,PressureLevel:pressure,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g),GlobalOutcomeRange:maxOut-minOut,DistinctOutcomes:len(distinct)}
			sum:=0
			for _,a:=range g{d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}
	return res,nil
}
