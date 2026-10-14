package unitary

import (
  "encoding/json"
  "fmt"
  "strings"
)

type NCI01WinglessCase struct {
  Seed int `json:"seed"`
  Direct bool `json:"direct"`
  TwoHop bool `json:"two_hop"`
  ThreeHop bool `json:"three_hop"`
  Interrupted bool `json:"interrupted_two_hop"`
  Retained bool `json:"old_knowledge_retained"`
  NewlyIntegrated bool `json:"new_knowledge_integrated"`
  DirectOnly bool `json:"direct_only_two_hop_baseline"`
  AdmissionTP int `json:"admission_tp"`
  AdmissionFP int `json:"admission_fp"`
  AdmissionFN int `json:"admission_fn"`
  Entries int `json:"entries"`
}
type NCI01WinglessResult struct {
  Schema string `json:"schema"`
  Phase string `json:"phase"`
  Verdict string `json:"verdict"`
  Parent string `json:"parent"`
  DevelopmentThreshold float64 `json:"development_threshold"`
  DevelopmentAccuracy float64 `json:"development_accuracy"`
  SeedRange [2]int `json:"seed_range"`
  OriginalLearningExamples int `json:"original_classifier_training_examples"`
  NativeAdmitReadoutWeights int `json:"native_admission_readout_weights"`
  RecurrentStateDimension int `json:"recurrent_state_dimension"`
  NativeRecallSlots int `json:"native_recall_slots"`
  MaxActualEntries int `json:"max_actual_entries"`
  AdmissionCalls int `json:"admission_calls"`
  RelationLookupHops int `json:"relation_lookup_hops"`
  CheckpointRestores int `json:"checkpoint_restores"`
  DirectAccuracy float64 `json:"direct_accuracy"`
  TwoHopAccuracy float64 `json:"two_hop_accuracy"`
  ThreeHopAccuracy float64 `json:"three_hop_accuracy"`
  InterruptedAccuracy float64 `json:"interrupted_accuracy"`
  RetentionAccuracy float64 `json:"retention_accuracy"`
  NewKnowledgeAccuracy float64 `json:"new_knowledge_accuracy"`
  DirectOnlyTwoHopBaseline float64 `json:"direct_only_two_hop_baseline"`
  AdmissionPrecision float64 `json:"admission_precision"`
  AdmissionRecall float64 `json:"admission_recall"`
  NativeCompositionRuleLearned bool `json:"native_composition_rule_learned"`
  RelationalStateLearned bool `json:"relational_state_learned"`
  NCI01FullAcceptance bool `json:"nci01_full_acceptance"`
  ExternalModelCalls int `json:"external_model_calls"`
  AcceptedRefMutation bool `json:"accepted_ref_mutation"`
  Records []NCI01WinglessCase `json:"records"`
}
func nci01WClause(a,b int,store bool)string{
 verb:="observes"
 if store{verb="stores"}
 return fmt.Sprintf("k%02d %s v%02d.",a,verb,b)
}
func nci01WAdmission(model *uplm0fClassifier,cut float64,memory *uplm0cRecall,clause string)(bool,bool){
 name,value,truth,_:=uplm0fClauseKV(clause)
 predicted:=model.prob(uplm0fEncode(clause))>=cut
 if predicted&&name!=""&&value!=""{memory.write(name,value)}
 return predicted,truth
}
func nci01WTraverse(memory *uplm0cRecall,start,hops int)string{
 key:=fmt.Sprintf("k%02d",start)
 for i:=0;i<hops;i++{
    value,ok:=memory.values[key]
    if !ok||!strings.HasPrefix(value,"v"){return ""}
    var next int
    if _,err:=fmt.Sscanf(value,"v%d",&next);err!=nil{return ""}
    if i==hops-1{return value}
    key=fmt.Sprintf("k%02d",next)
 }
 return ""
}
func nci01WCheckpoint(memory *uplm0cRecall)(*uplm0cRecall,error){
 blob,err:=json.Marshal(struct{
  Values map[string]string `json:"values"`
  Order []string `json:"order"`
 }{memory.values,memory.order})
 if err!=nil{return nil,err}
 var state struct{
  Values map[string]string `json:"values"`
  Order []string `json:"order"`
 }
 if err:=json.Unmarshal(blob,&state);err!=nil{return nil,err}
 return &uplm0cRecall{values:state.Values,order:state.Order},nil
}
func RunNCI01Wingless(phase string)(NCI01WinglessResult,error){
 if phase!="development"&&phase!="sealed"{return NCI01WinglessResult{},fmt.Errorf("unknown phase")}
 model:=uplm0fTrainClassifier()
 selected,best:=0.50,-1.0
 for _,cut:=range []float64{0.30,0.40,0.50,0.60,0.70}{
   hit,total:=0,0
   for a:=0;a<16;a++{for _,truth:=range []bool{true,false}{
     clause:=nci01WClause(a,(a+3)%24,truth)
     predicted:=model.prob(uplm0fEncode(clause))>=cut
     if predicted==truth{hit++};total++
   }}
   score:=float64(hit)/float64(total)
   if score>best{selected,best=cut,score}
 }
 result:=NCI01WinglessResult{
   Schema:"wingless.nci01.native-language-memory-integration.v1",
   Phase:phase,Verdict:"DEV_ONLY",
   Parent:"d5e32af7eff2f3e2a5ebd372f07c6072fec09c1b",
   DevelopmentThreshold:selected,DevelopmentAccuracy:best,SeedRange:[2]int{20700,20724},
   OriginalLearningExamples:20*24*32*2,NativeAdmitReadoutWeights:65,
   RecurrentStateDimension:64,NativeRecallSlots:16,
   NativeCompositionRuleLearned:false,RelationalStateLearned:false,
   NCI01FullAcceptance:false,ExternalModelCalls:0,AcceptedRefMutation:false,
 }
 if phase=="development"{return result,nil}
 tp,fp,fn:=0,0,0
 score:=func(f func(NCI01WinglessCase)bool)float64{
   count:=0
   for _,rec:=range result.Records{if f(rec){count++}}
   return float64(count)/float64(len(result.Records))
 }
 for seed:=20700;seed<20724;seed++{
   start:=(seed*7)%8
   var nodes [5]int
   for i:=0;i<5;i++{nodes[i]=24+(start+i)%8}
   memory:=newUPLM0CRecall()
   rec:=NCI01WinglessCase{Seed:seed}
   clauses:=[]string{
     nci01WClause(nodes[0],nodes[1],true),
     nci01WClause(43,29,false),
     nci01WClause(nodes[1],nodes[2],true),
     nci01WClause(44,28,false),
     nci01WClause(nodes[2],nodes[3],true),
   }
   for _,clause:=range clauses{
     pred,truth:=nci01WAdmission(model,selected,memory,clause)
     result.AdmissionCalls++
     if pred&&truth{tp++;rec.AdmissionTP++}
     if pred&&!truth{fp++;rec.AdmissionFP++}
     if !pred&&truth{fn++;rec.AdmissionFN++}
   }
   rec.Direct=nci01WTraverse(memory,nodes[0],1)==fmt.Sprintf("v%02d",nodes[1])
   rec.TwoHop=nci01WTraverse(memory,nodes[0],2)==fmt.Sprintf("v%02d",nodes[2])
   rec.ThreeHop=nci01WTraverse(memory,nodes[0],3)==fmt.Sprintf("v%02d",nodes[3])
   rec.DirectOnly=nci01WTraverse(memory,nodes[0],1)==fmt.Sprintf("v%02d",nodes[2])
   restored,err:=nci01WCheckpoint(memory)
   if err!=nil{return result,err}
   result.CheckpointRestores++
   rec.Interrupted=nci01WTraverse(restored,nodes[0],2)==fmt.Sprintf("v%02d",nodes[2])
   nci01WAdmission(model,selected,restored,nci01WClause(nodes[3],nodes[4],true))
   result.AdmissionCalls++
   rec.Retained=nci01WTraverse(restored,nodes[0],2)==fmt.Sprintf("v%02d",nodes[2])
   rec.NewlyIntegrated=nci01WTraverse(restored,nodes[0],4)==fmt.Sprintf("v%02d",nodes[4])
   rec.Entries=len(restored.order)
   if rec.Entries>result.MaxActualEntries{result.MaxActualEntries=rec.Entries}
   if rec.Entries>16{return result,fmt.Errorf("recall cap exceeded")}
   result.RelationLookupHops+=1+2+3+1+2+2+4
   result.Records=append(result.Records,rec)
 }
 if len(result.Records)!=24{return result,fmt.Errorf("unexpected cohort size")}
 result.DirectAccuracy=score(func(r NCI01WinglessCase)bool{return r.Direct})
 result.TwoHopAccuracy=score(func(r NCI01WinglessCase)bool{return r.TwoHop})
 result.ThreeHopAccuracy=score(func(r NCI01WinglessCase)bool{return r.ThreeHop})
 result.InterruptedAccuracy=score(func(r NCI01WinglessCase)bool{return r.Interrupted})
 result.RetentionAccuracy=score(func(r NCI01WinglessCase)bool{return r.Retained})
 result.NewKnowledgeAccuracy=score(func(r NCI01WinglessCase)bool{return r.NewlyIntegrated})
 result.DirectOnlyTwoHopBaseline=score(func(r NCI01WinglessCase)bool{return r.DirectOnly})
 result.AdmissionPrecision=float64(tp)/float64(max(1,tp+fp))
 result.AdmissionRecall=float64(tp)/float64(max(1,tp+fn))
 result.Verdict="COMPONENT_INTEGRATION_NOT_NATIVE_LEARNED_REASONING"
 if result.TwoHopAccuracy<0.75||result.NewKnowledgeAccuracy<0.75{
   result.Verdict="NCI01_COMPONENT_INTEGRATION_BELOW_FEASIBILITY_GATE"
 }
 return result,nil
}
