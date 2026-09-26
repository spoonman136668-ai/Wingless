package unitary

import "fmt"

const UPLM2TShamSchema="wingless.up-lm2t-sham-cost-control.v1"

type UPLM2TPoint struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	BaselineHarm int `json:"baseline_harm"`
	GuidedHarm int `json:"guided_harm"`
	ShamHarm int `json:"sham_harm"`
	GuidedInterventions int `json:"guided_interventions"`
	ShamInterventions int `json:"sham_interventions"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
}
type UPLM2TResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2SSeal string `json:"source_up_lm2s_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	DeferredLevels []int `json:"deferred_levels"`
	IdentityRotations []int `json:"identity_rotations"`
	ValueShifts []int `json:"value_shifts"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	BaselineHarm int `json:"baseline_harm"`
	GuidedHarm int `json:"guided_harm"`
	ShamHarm int `json:"sham_harm"`
	GuidedInterventions int `json:"guided_interventions"`
	ShamInterventions int `json:"sham_interventions"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
	GuidedPreventionFraction float64 `json:"guided_prevention_fraction"`
	ShamPreventionFraction float64 `json:"sham_prevention_fraction"`
	Points []UPLM2TPoint `json:"points"`
}
func uplm2tTriggers(d,rot,shift int)[]int{
	r,reported:=uplm2sInit(d,rot,shift);out:=[]int{}
	for j:=0;j<12;j++{
		if len(r.order)>=16{victim:=r.order[0];if !reported[victim]{out=append(out,j);reported[victim]=true}}
		r.write(fmt.Sprintf("guide-plan-%d-%d-%d-%d",d,rot,shift,j),"x")
	}
	return out
}
func uplm2tBaseline(d,rot,shift int)int{
	r,reported:=uplm2sInit(d,rot,shift);harm:=0
	for j:=0;j<12;j++{if len(r.order)>=16&&!reported[r.order[0]]{harm++};r.write(fmt.Sprintf("base2-%d-%d-%d-%d",d,rot,shift,j),"x")}
	return harm
}
func uplm2tGuided(d,rot,shift int,triggers map[int]bool)(int,int){
	r,reported:=uplm2sInit(d,rot,shift);harm,acts:=0,0
	for j:=0;j<12;j++{
		if triggers[j]&&len(r.order)>=16{victim:=r.order[0];reported[victim]=true;acts++}
		if len(r.order)>=16&&!reported[r.order[0]]{harm++}
		r.write(fmt.Sprintf("guided2-%d-%d-%d-%d",d,rot,shift,j),"x")
	}
	return harm,acts
}
func uplm2tShamTarget(r *uplm0cRecall,reported map[string]bool,victim string)string{
	for i:=len(r.order)-1;i>=0;i--{n:=r.order[i];if n!=victim&&!reported[n]{return n}}
	return ""
}
func uplm2tSham(d,rot,shift int,triggers map[int]bool)(int,int){
	r,reported:=uplm2sInit(d,rot,shift);harm,acts:=0,0
	for j:=0;j<12;j++{
		if triggers[j]{
			victim:="";if len(r.order)>0{victim=r.order[0]}
			target:=uplm2tShamTarget(r,reported,victim)
			if target!=""{reported[target]=true;acts++}
		}
		if len(r.order)>=16&&!reported[r.order[0]]{harm++}
		r.write(fmt.Sprintf("sham2-%d-%d-%d-%d",d,rot,shift,j),"x")
	}
	return harm,acts
}
func RunUPLM2T()(UPLM2TResult,error){
	levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
	res:=UPLM2TResult{Schema:UPLM2TShamSchema,Experiment:"UP-LM2T-sham-cost-control",SourceUPLM2SSeal:"1b755475d311f70774af1d40dea695a97e9a96d9",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,CounterfactualOnly:true,LiveActivation:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{
		ts:=uplm2tTriggers(d,rot,shift);set:=map[int]bool{};for _,j:=range ts{set[j]=true}
		base:=uplm2tBaseline(d,rot,shift);gh,ga:=uplm2tGuided(d,rot,shift,set);sh,sa:=uplm2tSham(d,rot,shift,set)
		p:=UPLM2TPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,BaselineHarm:base,GuidedHarm:gh,ShamHarm:sh,GuidedInterventions:ga,ShamInterventions:sa,GuidedPrevented:base-gh,ShamPrevented:base-sh}
		res.Points=append(res.Points,p);res.BaselineHarm+=base;res.GuidedHarm+=gh;res.ShamHarm+=sh;res.GuidedInterventions+=ga;res.ShamInterventions+=sa;res.GuidedPrevented+=base-gh;res.ShamPrevented+=base-sh
	}}}
	res.GuidedPreventionFraction=uplm2sRate(res.GuidedPrevented,res.BaselineHarm);res.ShamPreventionFraction=uplm2sRate(res.ShamPrevented,res.BaselineHarm)
	return res,nil
}
