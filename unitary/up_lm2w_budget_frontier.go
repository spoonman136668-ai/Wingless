package unitary

import "fmt"

const UPLM2WBudgetSchema="wingless.up-lm2w-budget-frontier.v1"

type UPLM2WPoint struct{
	Budget int `json:"budget"`
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	PendingLayout string `json:"pending_layout"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
}
type UPLM2WBudgetMetric struct{
	Budget int `json:"budget"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
	GuidedPreventedPerAction float64 `json:"guided_prevented_per_action"`
	ShamPreventedPerAction float64 `json:"sham_prevented_per_action"`
}
type UPLM2WResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2VSeal string `json:"source_up_lm2v_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Budgets []int `json:"budgets"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	Points []UPLM2WPoint `json:"points"`
	Metrics []UPLM2WBudgetMetric `json:"metrics"`
}
func uplm2wRun(d,rot,shift int,layout,mode string,budget int)(completed,failed,actions int){
	r,reported:=uplm2vInit(d,rot,shift,layout)
	for j:=0;j<12;j++{
		victim:="";harm:=false
		if len(r.order)>=16{victim=r.order[0];harm=!reported[victim]}
		if harm&&actions<budget&&mode=="guided"{reported[victim]=true;actions++}
		if harm&&actions<budget&&mode=="sham"{if t:=uplm2vShamTarget(r,reported,victim);t!=""{reported[t]=true;actions++}}
		r.write(fmt.Sprintf("w-%s-%s-%d-%d-%d-%d-%d",mode,layout,budget,d,rot,shift,j),"x")
	}
	completed,failed=uplm2vFinish(r,reported,rot);return completed,failed,actions
}
func RunUPLM2W()(UPLM2WResult,error){
	levels:=[]int{4,5,6};rots:=[]int{3,11};shifts:=[]int{0,2};layouts:=[]string{"suffix_reported","alternating_reported"};budgets:=[]int{0,1,2,3}
	res:=UPLM2WResult{Schema:UPLM2WBudgetSchema,Experiment:"UP-LM2W-budget-frontier",SourceUPLM2VSeal:"d639f6e56f3cb923daa3c5d91da173b8a3eea9f0",ExactRecallCap:16,Budgets:budgets,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for _,budget:=range budgets{
		m:=UPLM2WBudgetMetric{Budget:budget}
		for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{for _,layout:=range layouts{
			_,bf,_:=uplm2wRun(d,rot,shift,layout,"baseline",0)
			_,gf,ga:=uplm2wRun(d,rot,shift,layout,"guided",budget)
			_,sf,sa:=uplm2wRun(d,rot,shift,layout,"sham",budget)
			res.Points=append(res.Points,UPLM2WPoint{Budget:budget,DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,PendingLayout:layout,BaselineFailed:bf,GuidedFailed:gf,ShamFailed:sf,GuidedActions:ga,ShamActions:sa})
			m.BaselineFailed+=bf;m.GuidedFailed+=gf;m.ShamFailed+=sf;m.GuidedActions+=ga;m.ShamActions+=sa
		}}}}
		m.GuidedPrevented=m.BaselineFailed-m.GuidedFailed;m.ShamPrevented=m.BaselineFailed-m.ShamFailed
		m.GuidedPreventedPerAction=uplm2sRate(m.GuidedPrevented,m.GuidedActions);m.ShamPreventedPerAction=uplm2sRate(m.ShamPrevented,m.ShamActions)
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
