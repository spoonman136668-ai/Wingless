package unitary

import "fmt"

const UPLM2UCostSchema="wingless.up-lm2u-downstream-cost.v1"

type UPLM2UPoint struct{
	DeferredFirstChunkReports int `json:"deferred_first_chunk_reports"`
	IdentityRotation int `json:"identity_rotation"`
	ValueShift int `json:"value_shift"`
	BaselineCompleted int `json:"baseline_completed"`
	GuidedCompleted int `json:"guided_completed"`
	ShamCompleted int `json:"sham_completed"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
}
type UPLM2UResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM2TSeal string `json:"source_up_lm2t_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	CapacityChanged bool `json:"capacity_changed"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	TotalArms int `json:"total_arms"`
	BaselineCompleted int `json:"baseline_completed"`
	GuidedCompleted int `json:"guided_completed"`
	ShamCompleted int `json:"sham_completed"`
	BaselineFailed int `json:"baseline_failed"`
	GuidedFailed int `json:"guided_failed"`
	ShamFailed int `json:"sham_failed"`
	GuidedActions int `json:"guided_actions"`
	ShamActions int `json:"sham_actions"`
	Points []UPLM2UPoint `json:"points"`
}
func uplm2uFinish(r *uplm0cRecall,reported map[string]bool,rot int)(completed,failed int){
	for i:=0;i<12;i++{
		n:=uplm2nName(i,rot)
		if reported[n]{completed++;continue}
		if _,ok:=r.values[n];ok{completed++}else{failed++}
	}
	return
}
func uplm2uRun(d,rot,shift int,mode string,triggers map[int]bool)(completed,failed,actions int){
	r,reported:=uplm2sInit(d,rot,shift)
	for j:=0;j<12;j++{
		if triggers[j]{
			if mode=="guided"&&len(r.order)>=16{
				v:=r.order[0]
				if !reported[v]{reported[v]=true;actions++}
			}
			if mode=="sham"{
				v:="";if len(r.order)>0{v=r.order[0]}
				t:=uplm2tShamTarget(r,reported,v)
				if t!=""{reported[t]=true;actions++}
			}
		}
		r.write(fmt.Sprintf("cost-%s-%d-%d-%d-%d",mode,d,rot,shift,j),"x")
	}
	completed,failed=uplm2uFinish(r,reported,rot)
	return completed,failed,actions
}
func RunUPLM2U()(UPLM2UResult,error){
	levels:=[]int{4,5,6};rots:=[]int{0,7};shifts:=[]int{0,1,2,3}
	res:=UPLM2UResult{Schema:UPLM2UCostSchema,Experiment:"UP-LM2U-downstream-cost",SourceUPLM2TSeal:"73b5d5a6d7025a4e97c4eed06b50d9dc502d3417",ExactRecallCap:16,CounterfactualOnly:true,LiveActivation:false,CapacityChanged:false,ExtraTrainingUsed:false}
	for _,d:=range levels{for _,rot:=range rots{for _,shift:=range shifts{
		ts:=uplm2tTriggers(d,rot,shift);set:=map[int]bool{};for _,j:=range ts{set[j]=true}
		bc,bf,_:=uplm2uRun(d,rot,shift,"baseline",map[int]bool{})
		gc,gf,ga:=uplm2uRun(d,rot,shift,"guided",set)
		sc,sf,sa:=uplm2uRun(d,rot,shift,"sham",set)
		p:=UPLM2UPoint{DeferredFirstChunkReports:d,IdentityRotation:rot,ValueShift:shift,BaselineCompleted:bc,GuidedCompleted:gc,ShamCompleted:sc,BaselineFailed:bf,GuidedFailed:gf,ShamFailed:sf,GuidedActions:ga,ShamActions:sa}
		res.Points=append(res.Points,p);res.TotalArms++;res.BaselineCompleted+=bc;res.GuidedCompleted+=gc;res.ShamCompleted+=sc;res.BaselineFailed+=bf;res.GuidedFailed+=gf;res.ShamFailed+=sf;res.GuidedActions+=ga;res.ShamActions+=sa
	}}}
	return res,nil
}
