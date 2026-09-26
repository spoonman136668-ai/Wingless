package unitary

import "fmt"

const UPLM2VGeneralizationSchema="wingless.up-lm2v-response-generalization.v1"

type UPLM2VPoint struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	PendingLayout string `json:"pending_layout"`
	BaselineCompleted int `json:"baseline_completed"`
	GuidedCompleted int `json:"guided_completed"`
	ShamCompleted int `json:"sham_completed"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
}
type UPLM2VResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2USeal string `json:"source_up_lm2u_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	DeferredLevels []int `json:"deferred_levels"`
	IdentityRotations []int `json:"identity_rotations"`
	ValueShifts []int `json:"value_shifts"`
	PendingLayouts []string `json:"pending_layouts"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	FutureScheduleUsed bool `json:"future_schedule_used"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	TotalArms int `json:"total_arms"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
	GuidedPrevented int `json:"guided_prevented"`
	ShamPrevented int `json:"sham_prevented"`
	Points []UPLM2VPoint `json:"points"`
}
func uplm2vInit(d,rot,shift int,layout string)(*uplm0cRecall,map[string]bool){
	r:=newUPLM0CRecall();reported:=map[string]bool{}
	for i:=0;i<12;i++{r.write(uplm2nName(i,rot),uplm2nValue(i,rot,shift))}
	n:=12-d
	if layout=="suffix_reported"{
		for i:=12-n;i<12;i++{reported[uplm2nName(i,rot)]=true}
	}else{
		perm:=[]int{1,3,5,7,9,11,0,2,4,6,8,10}
		for i:=0;i<n;i++{reported[uplm2nName(perm[i],rot)]=true}
	}
	return r,reported
}
func uplm2vFinish(r *uplm0cRecall,reported map[string]bool,rot int)(completed,failed int){
	for i:=0;i<12;i++{n:=uplm2nName(i,rot);if reported[n]{completed++}else if _,ok:=r.values[n];ok{completed++}else{failed++}}
	return
}
func uplm2vShamTarget(r *uplm0cRecall,reported map[string]bool,victim string)string{
	for i:=len(r.order)-1;i>=0;i--{n:=r.order[i];if n!=victim&&!reported[n]{return n}}
	return ""
}
func uplm2vRun(d,rot,shift int,layout,mode string)(completed,failed,actions int){
	r,reported:=uplm2vInit(d,rot,shift,layout)
	for j:=0;j<12;j++{
		victim:="";harm:=false
		if len(r.order)>=16{victim=r.order[0];harm=!reported[victim]}
		if harm&&mode=="guided"{reported[victim]=true;actions++}
		if harm&&mode=="sham"{if t:=uplm2vShamTarget(r,reported,victim);t!=""{reported[t]=true;actions++}}
		r.write(fmt.Sprintf("v-%s-%s-%d-%d-%d-%d",mode,layout,d,rot,shift,j),"x")
	}
	completed,failed=uplm2vFinish(r,reported,rot)
	return completed,failed,actions
}
func RunUPLM2V()(UPLM2VResult,error){
	levels:=[]int{4,5,6};rots:=[]int{3,11};shifts:=[]int{0,2};layouts:=[]string{"suffix_reported","alternating_reported"}
	res:=UPLM2VResult{Schema:UPLM2VGeneralizationSchema,Experiment:"UP-LM2V-response-generalization",SourceUPLM2USeal:"e5a45d44e7be97dabfc21700bdcc95c567dc57e1",ExactRecallCap:16,DeferredLevels:levels,IdentityRotations:rots,ValueShifts:shifts,PendingLayouts:layouts,CounterfactualOnly:true,LiveActivation:false,FutureScheduleUsed:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{for _,layout:=range layouts{
		bc,bf,_:=uplm2vRun(d,rot,shift,layout,"baseline");gc,gf,ga:=uplm2vRun(d,rot,shift,layout,"guided");sc,sf,sa:=uplm2vRun(d,rot,shift,layout,"sham")
		res.Points=append(res.Points,UPLM2VPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,PendingLayout:layout,BaselineCompleted:bc,GuidedCompleted:gc,ShamCompleted:sc,BaselineFailed:bf,GuidedFailed:gf,ShamFailed:sf,GuidedActions:ga,ShamActions:sa})
		res.TotalArms++;res.BaselineFailed+=bf;res.GuidedFailed+=gf;res.ShamFailed+=sf;res.GuidedActions+=ga;res.ShamActions+=sa;res.GuidedPrevented+=bf-gf;res.ShamPrevented+=bf-sf
	}}}}
	return res,nil
}
