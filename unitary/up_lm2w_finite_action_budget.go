package unitary

import "fmt"

const UPLM2WBudgetSchema="wingless.up-lm2w-finite-action-budget.v1"

type UPLM2WPoint struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	PendingLayout string `json:"pending_layout"`
	ActionBudget int `json:"action_budget"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
	GuidedResidualFailures int `json:"guided_residual_failures"`
	PreventionPerGuidedAction float64 `json:"prevention_per_guided_action"`
	GuidedUnnecessaryActions int `json:"guided_unnecessary_actions"`
}
type UPLM2WResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2VSeal string `json:"source_up_lm2v_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	ActionBudgets []int `json:"action_budgets"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	TotalArms int `json:"total_arms"`
	TotalGuidedActions int `json:"total_guided_actions"`
	TotalGuidedPrevented int `json:"total_guided_prevented"`
	TotalGuidedResidualFailures int `json:"total_guided_residual_failures"`
	TotalGuidedUnnecessaryActions int `json:"total_guided_unnecessary_actions"`
	Points []UPLM2WPoint `json:"points"`
}
func uplm2wRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func uplm2wRun(d,rot,shift int,layout,mode string,budget int)(completed,failed,actions,unnecessary int){
	r,reported:=uplm2vInit(d,rot,shift,layout)
	for j:=0;j<12;j++{
		victim:="";harm:=false
		if len(r.order)>=16{victim=r.order[0];harm=!reported[victim]}
		if harm&&actions<budget&&mode=="guided"{reported[victim]=true;actions++}
		if harm&&actions<budget&&mode=="sham"{if t:=uplm2vShamTarget(r,reported,victim);t!=""{reported[t]=true;actions++}}
		if mode=="guided"&&actions>budget{unnecessary++}
		r.write(fmt.Sprintf("w-%s-%s-%d-%d-%d-%d-%d",mode,layout,d,rot,shift,budget,j),"x")
	}
	completed,failed=uplm2vFinish(r,reported,rot)
	return
}
func RunUPLM2W()(UPLM2WResult,error){
	levels:=[]int{4,5,6};rots:=[]int{3,11};shifts:=[]int{0,2};layouts:=[]string{"suffix_reported","alternating_reported"};budgets:=[]int{1,2,3}
	res:=UPLM2WResult{Schema:UPLM2WBudgetSchema,Experiment:"UP-LM2W-finite-action-budget",SourceUPLM2VSeal:"d639f6e56f3cb923daa3c5d91da173b8a3eea9f0",ExactRecallCap:16,ActionBudgets:budgets,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{for _,layout:=range layouts{for _,budget:=range budgets{
		_,bf,_,_:=uplm2wRun(d,rot,shift,layout,"baseline",0)
		_,gf,ga,gu:=uplm2wRun(d,rot,shift,layout,"guided",budget)
		_,sf,sa,_:=uplm2wRun(d,rot,shift,layout,"sham",budget)
		gp:=bf-gf;sp:=bf-sf
		res.Points=append(res.Points,UPLM2WPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,PendingLayout:layout,ActionBudget:budget,BaselineFailed:bf,GuidedFailed:gf,ShamFailed:sf,GuidedActions:ga,ShamActions:sa,GuidedPrevented:gp,ShamPrevented:sp,GuidedResidualFailures:gf,PreventionPerGuidedAction:uplm2wRate(gp,ga),GuidedUnnecessaryActions:gu})
		res.TotalArms++;res.TotalGuidedActions+=ga;res.TotalGuidedPrevented+=gp;res.TotalGuidedResidualFailures+=gf;res.TotalGuidedUnnecessaryActions+=gu
	}}}}}
	return res,nil
}
