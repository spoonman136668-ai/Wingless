package unitary

import (
 "encoding/json"
 "fmt"
 "math"
 "strconv"
 "strings"
)

type nci02WEdge struct {From int; To int}
type nci02WJoin struct {W [5]float64; B float64}
type NCI02WTrial struct {
 Seed int `json:"seed"`
 DistinctStructure int `json:"distinct_structure"`
 TwoHop bool `json:"two_hop"`
 ThreeHop bool `json:"three_hop"`
 TwoHopNegative bool `json:"two_hop_negative"`
 ThreeHopNegative bool `json:"three_hop_negative"`
 Interrupted bool `json:"interrupted"`
 Retained bool `json:"retained"`
 NewFourHop bool `json:"new_four_hop"`
 DirectOnly bool `json:"direct_only_baseline"`
 OracleTwoHop bool `json:"oracle_two_hop"`
 AdmissionTP int `json:"admission_tp"`
 AdmissionFP int `json:"admission_fp"`
 AdmissionFN int `json:"admission_fn"`
 Entries int `json:"entries"`
}
type NCI02WResult struct {
 Schema string `json:"schema"`
 Parent string `json:"parent"`
 Phase string `json:"phase"`
 Verdict string `json:"verdict"`
 DevelopmentJoinAccuracy float64 `json:"development_join_accuracy"`
 JoinThreshold float64 `json:"join_threshold"`
 JoinWeights [5]float64 `json:"join_weights"`
 JoinBias float64 `json:"join_bias"`
 HeldoutSeedRange [2]int `json:"heldout_seed_range"`
 TrainingExamples int `json:"training_examples"`
 TrainingUpdates int `json:"training_updates"`
 LearnableJoinParameters int `json:"learnable_join_parameters"`
 NativeAdmissionParameters int `json:"native_admission_parameters"`
 NativeRecallCap int `json:"native_recall_cap"`
 PeakEntries int `json:"peak_entries"`
 JoinCandidateEvaluations int `json:"join_candidate_evaluations"`
 OriginalAdmissionCalls int `json:"original_admission_calls"`
 CheckpointRestores int `json:"checkpoint_restores"`
 Accuracy map[string]float64 `json:"accuracy"`
 AdmissionPrecision float64 `json:"admission_precision"`
 AdmissionRecall float64 `json:"admission_recall"`
 LearnedCompositionPredicate bool `json:"learned_composition_predicate"`
 LearnedRelationalStorage bool `json:"learned_relational_storage"`
 NCI01FullAcceptance bool `json:"nci01_full_acceptance"`
 ExternalModelCalls int `json:"external_model_calls"`
 AcceptedRefMutation bool `json:"accepted_ref_mutation"`
 Records []NCI02WTrial `json:"records"`
}

func nci02WFeatures(a,b nci02WEdge) [5]float64 {
 eq:=func(x,y int)float64 {if x==y{return 1};return 0}
 return [5]float64{eq(a.To,b.From),eq(a.From,b.From),eq(a.To,b.To),
  eq(a.From,b.To),float64((a.To-b.From+256)%7)/7.0}
}
func nci02WScore(model nci02WJoin,a,b nci02WEdge)float64{
 f:=nci02WFeatures(a,b)
 z:=model.B
 for k,v:=range f {z+=model.W[k]*v}
 if z>=0 {return 1/(1+math.Exp(-z))}
 e:=math.Exp(z);return e/(1+e)
}
func nci02WTrain()(nci02WJoin,float64,float64){
 model:=nci02WJoin{}
 // Names only 0..23 during training. Heldout identities 32..63.
 for epoch:=0;epoch<65;epoch++ {
  for i:=0;i<384;i++ {
   a:=nci02WEdge{From:(i*5+3)%24,To:(i*7+9)%24}
   candidates:=[2]nci02WEdge{
    {From:a.To,To:(i*11+13)%24},
    {From:(a.To+1+i%7)%24,To:(i*11+13)%24},
   }
   for j,b:=range candidates {
    target:=1.0;if j==1{target=0}
    p:=nci02WScore(model,a,b);gradient:=p-target
    features:=nci02WFeatures(a,b)
    for k,x:=range features{model.W[k]-=0.10*gradient*x}
    model.B-=0.10*gradient
   }
  }
 }
 threshold,best:=0.5,-1.0
 for _,cut:=range []float64{.4,.5,.6,.7}{
  correct,total:=0,0
  for i:=0;i<160;i++ {
   a:=nci02WEdge{From:24+i%8,To:24+(i*3+1)%8}
   good:=nci02WEdge{From:a.To,To:24+(i*5+3)%8}
   bad:=nci02WEdge{From:24+(a.To-23+i%6)%8,To:good.To}
   if bad.From==a.To {bad.From=24+(bad.From-23)%8}
   for j,b:=range []nci02WEdge{good,bad} {
    truth:=j==0
    if (nci02WScore(model,a,b)>=cut)==truth{correct++}
    total++
   }
  }
  acc:=float64(correct)/float64(total)
  if acc>best{best,threshold=acc,cut}
 }
 return model,threshold,best
}
func nci02WEdges(memory *uplm0cRecall) []nci02WEdge{
 result:=make([]nci02WEdge,0,len(memory.order))
 for _,key:=range memory.order{
  v,ok:=memory.values[key];if !ok{continue}
  if !strings.HasPrefix(key,"k")||!strings.HasPrefix(v,"v"){continue}
  from,e1:=strconv.Atoi(strings.TrimPrefix(key,"k"))
  to,e2:=strconv.Atoi(strings.TrimPrefix(v,"v"))
  if e1==nil&&e2==nil{result=append(result,nci02WEdge{from,to})}
 }
 return result
}
func nci02WContains(edges []nci02WEdge,from,to int)bool{
 for _,e:=range edges{if e.From==from&&e.To==to{return true}}
 return false
}
func nci02WReach(model nci02WJoin,cut float64,edges []nci02WEdge,from,to,hops int,counter *int)bool{
 if hops<1||hops>4{return false}
 frontier:=edges
 if hops==1{return nci02WContains(frontier,from,to)}
 for d:=2;d<=hops;d++{
  seen:=map[nci02WEdge]bool{}
  next:=make([]nci02WEdge,0)
  for _,a:=range frontier{for _,b:=range edges{
   *counter++
   if nci02WScore(model,a,b)<cut{continue}
   joined:=nci02WEdge{From:a.From,To:b.To}
   if !seen[joined]{seen[joined]=true;next=append(next,joined)}
  }}
  if len(next)>256{return false}
  frontier=next
 }
 return nci02WContains(frontier,from,to)
}
func RunNCI02Wingless(phase string)(NCI02WResult,error){
 if phase!="development"&&phase!="sealed"{return NCI02WResult{},fmt.Errorf("invalid phase")}
 join,cut,dev:=nci02WTrain()
 result:=NCI02WResult{
  Schema:"wingless.nci02.learned-composition-predicate.v1",
  Parent:"4cac8fa96a9c1aaad854ae62fbf4a22777d6ba09",Phase:phase,
  Verdict:"DEVELOPMENT_ONLY",DevelopmentJoinAccuracy:dev,
  JoinThreshold:cut,JoinWeights:join.W,JoinBias:join.B,
  HeldoutSeedRange:[2]int{22112,22136},TrainingExamples:384*2,TrainingUpdates:65*384*2,
  LearnableJoinParameters:6,NativeAdmissionParameters:65,NativeRecallCap:16,
  LearnedCompositionPredicate:true,LearnedRelationalStorage:false,
  NCI01FullAcceptance:false,ExternalModelCalls:0,AcceptedRefMutation:false,
 }
 if phase=="development"{return result,nil}
 admit:=uplm0fTrainClassifier()
 tp,fp,fn:=0,0,0
 keys:=[]string{"two_hop","three_hop","two_hop_negative","three_hop_negative",
  "interrupted","retained","new_four_hop","direct_only_baseline","oracle_two_hop"}
 counts:=map[string]int{};for _,key:=range keys{counts[key]=0}
 for seed:=22112;seed<22136;seed++{
  offset:=(seed*11)%32
  var nodes [5]int
  for k:=range nodes{nodes[k]=32+(offset+k*7)%32}
  d1,d2:=32+(offset+3)%32,32+(offset+11)%32
  relations:=[]nci02WEdge{{nodes[0],nodes[1]},{d1,d2},
   {nodes[1],nodes[2]},{d2,nodes[1]},{nodes[2],nodes[3]}}
  store:=newUPLM0CRecall()
  trial:=NCI02WTrial{Seed:seed,DistinctStructure:offset}
  for i:=0;i<len(relations);i++{
   edge:=relations[(i+seed%len(relations))%len(relations)]
   statement:=nci01WClause(edge.From,edge.To,true)
   pred,truth:=nci01WAdmission(admit,.5,store,statement)
   result.OriginalAdmissionCalls++
   if pred&&truth{tp++;trial.AdmissionTP++}
   if pred&&!truth{fp++;trial.AdmissionFP++}
   if !pred&&truth{fn++;trial.AdmissionFN++}
  }
  noise:=nci01WClause(80+seed%5,84+seed%5,false)
  predicted,truth:=nci01WAdmission(admit,.5,store,noise)
  result.OriginalAdmissionCalls++
  if predicted&&!truth{fp++;trial.AdmissionFP++}
  if !predicted&&truth{fn++;trial.AdmissionFN++}
  edges:=nci02WEdges(store)
  trial.TwoHop=nci02WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  trial.ThreeHop=nci02WReach(join,cut,edges,nodes[0],nodes[3],3,&result.JoinCandidateEvaluations)
  trial.TwoHopNegative=!nci02WReach(join,cut,edges,nodes[0],nodes[3],2,&result.JoinCandidateEvaluations)
  trial.ThreeHopNegative=!nci02WReach(join,cut,edges,nodes[0],nodes[4],3,&result.JoinCandidateEvaluations)
  trial.DirectOnly=nci02WContains(edges,nodes[0],nodes[2])
  trial.OracleTwoHop=nci02WContains(edges,nodes[0],nodes[1])&&nci02WContains(edges,nodes[1],nodes[2])
  restored,err:=nci01WCheckpoint(store)
  if err!=nil{return result,err}
  result.CheckpointRestores++
  edges=nci02WEdges(restored)
  trial.Interrupted=nci02WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  p,_:=nci01WAdmission(admit,.5,restored,nci01WClause(nodes[3],nodes[4],true))
  result.OriginalAdmissionCalls++;if p{tp++}else{fn++}
  edges=nci02WEdges(restored)
  trial.Retained=nci02WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  trial.NewFourHop=nci02WReach(join,cut,edges,nodes[0],nodes[4],4,&result.JoinCandidateEvaluations)
  trial.Entries=len(restored.order)
  if trial.Entries>result.PeakEntries{result.PeakEntries=trial.Entries}
  if trial.Entries>result.NativeRecallCap{return result,fmt.Errorf("native memory cap")}
  for _,kv:=range []struct{key string;value bool}{
   {"two_hop",trial.TwoHop},{"three_hop",trial.ThreeHop},
   {"two_hop_negative",trial.TwoHopNegative},{"three_hop_negative",trial.ThreeHopNegative},
   {"interrupted",trial.Interrupted},{"retained",trial.Retained},
   {"new_four_hop",trial.NewFourHop},{"direct_only_baseline",trial.DirectOnly},
   {"oracle_two_hop",trial.OracleTwoHop}} {if kv.value{counts[kv.key]++}}
  result.Records=append(result.Records,trial)
 }
 result.Accuracy=map[string]float64{}
 for _,key:=range keys{result.Accuracy[key]=float64(counts[key])/float64(len(result.Records))}
 result.AdmissionPrecision=float64(tp)/float64(max(1,tp+fp))
 result.AdmissionRecall=float64(tp)/float64(max(1,tp+fn))
 result.Verdict="LEARNED_PREDICATE_ONLY_NOT_NATIVE_RELATIONAL_STATE"
 if result.Accuracy["two_hop"]<.75||result.Accuracy["three_hop"]<.75||result.Accuracy["new_four_hop"]<.75||result.Accuracy["two_hop_negative"]<.75{result.Verdict="LEARNED_JOIN_INTEGRATION_BELOW_GATE"}
 return result,nil
}
