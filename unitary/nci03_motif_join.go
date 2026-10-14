package unitary

import (
 "fmt"
 "math"
 "strconv"
 "strings"
)

type nci03WEdge struct {From int; To int}
type nci03WJoin struct {W [5]float64; B float64}
type NCI03WTrial struct {
 Seed int `json:"seed"`
 DistinctStructure int `json:"distinct_structure"`
 TopologyClass int `json:"topology_class"`
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
type NCI03WResult struct {
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
 MotifAccuracy map[string]map[string]float64 `json:"motif_accuracy"`
 AdmissionPrecision float64 `json:"admission_precision"`
 AdmissionRecall float64 `json:"admission_recall"`
 LearnedCompositionPredicate bool `json:"learned_composition_predicate"`
 LearnedRelationalStorage bool `json:"learned_relational_storage"`
 NCI01FullAcceptance bool `json:"nci01_full_acceptance"`
 ExternalModelCalls int `json:"external_model_calls"`
 AcceptedRefMutation bool `json:"accepted_ref_mutation"`
 Records []NCI03WTrial `json:"records"`
}

func nci03WFeatures(a,b nci03WEdge) [5]float64 {
 eq:=func(x,y int)float64 {if x==y{return 1};return 0}
 return [5]float64{eq(a.To,b.From),eq(a.From,b.From),eq(a.To,b.To),
  eq(a.From,b.To),float64((a.To-b.From+256)%7)/7.0}
}
func nci03WScore(model nci03WJoin,a,b nci03WEdge)float64{
 f:=nci03WFeatures(a,b)
 z:=model.B
 for k,v:=range f {z+=model.W[k]*v}
 if z>=0 {return 1/(1+math.Exp(-z))}
 e:=math.Exp(z);return e/(1+e)
}
func nci03WTrain()(nci03WJoin,float64,float64){
 model:=nci03WJoin{}
 // Names only 0..23 during training. Heldout identities 32..63.
 for epoch:=0;epoch<65;epoch++ {
  for i:=0;i<384;i++ {
   a:=nci03WEdge{From:(i*5+3)%24,To:(i*7+9)%24}
   candidates:=[2]nci03WEdge{
    {From:a.To,To:(i*11+13)%24},
    {From:(a.To+1+i%7)%24,To:(i*11+13)%24},
   }
   for j,b:=range candidates {
    target:=1.0;if j==1{target=0}
    p:=nci03WScore(model,a,b);gradient:=p-target
    features:=nci03WFeatures(a,b)
    for k,x:=range features{model.W[k]-=0.10*gradient*x}
    model.B-=0.10*gradient
   }
  }
 }
 threshold,best:=0.5,-1.0
 for _,cut:=range []float64{.4,.5,.6,.7}{
  correct,total:=0,0
  for i:=0;i<160;i++ {
   a:=nci03WEdge{From:24+i%8,To:24+(i*3+1)%8}
   good:=nci03WEdge{From:a.To,To:24+(i*5+3)%8}
   bad:=nci03WEdge{From:24+(a.To-23+i%6)%8,To:good.To}
   if bad.From==a.To {bad.From=24+(bad.From-23)%8}
   for j,b:=range []nci03WEdge{good,bad} {
    truth:=j==0
    if (nci03WScore(model,a,b)>=cut)==truth{correct++}
    total++
   }
  }
  acc:=float64(correct)/float64(total)
  if acc>best{best,threshold=acc,cut}
 }
 return model,threshold,best
}
func nci03WEdges(memory *uplm0cRecall) []nci03WEdge{
 result:=make([]nci03WEdge,0,len(memory.order))
 for _,key:=range memory.order{
  v,ok:=memory.values[key];if !ok{continue}
  if !strings.HasPrefix(key,"k")||!strings.HasPrefix(v,"v"){continue}
  from,e1:=strconv.Atoi(strings.TrimPrefix(key,"k"))
  to,e2:=strconv.Atoi(strings.TrimPrefix(v,"v"))
  if e1==nil&&e2==nil{result=append(result,nci03WEdge{from,to})}
 }
 return result
}
func nci03WContains(edges []nci03WEdge,from,to int)bool{
 for _,e:=range edges{if e.From==from&&e.To==to{return true}}
 return false
}
func nci03WReach(model nci03WJoin,cut float64,edges []nci03WEdge,from,to,hops int,counter *int)bool{
 if hops<1||hops>4{return false}
 frontier:=edges
 if hops==1{return nci03WContains(frontier,from,to)}
 for d:=2;d<=hops;d++{
  seen:=map[nci03WEdge]bool{}
  next:=make([]nci03WEdge,0)
  for _,a:=range frontier{for _,b:=range edges{
   *counter++
   if nci03WScore(model,a,b)<cut{continue}
   joined:=nci03WEdge{From:a.From,To:b.To}
   if !seen[joined]{seen[joined]=true;next=append(next,joined)}
  }}
  if len(next)>256{return false}
  frontier=next
 }
 return nci03WContains(frontier,from,to)
}
func RunNCI03Wingless(phase string)(NCI03WResult,error){
 if phase!="development"&&phase!="sealed"{return NCI03WResult{},fmt.Errorf("invalid phase")}
 join,cut,dev:=nci03WTrain()
 result:=NCI03WResult{
  Schema:"wingless.nci03.motif-generalization.v1",
  Parent:"bb3f873924ffebc5180c3ebe2c563dbe8eaf2f77",Phase:phase,
  Verdict:"DEVELOPMENT_ONLY",DevelopmentJoinAccuracy:dev,
  JoinThreshold:cut,JoinWeights:join.W,JoinBias:join.B,
  HeldoutSeedRange:[2]int{22976,23000},TrainingExamples:384*2,TrainingUpdates:65*384*2,
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
 for seed:=22976;seed<23000;seed++{
  offset:=(seed*11)%32
  var nodes [5]int
  for k:=range nodes{nodes[k]=32+(offset+k*7)%32}
  d1,d2:=32+(offset+3)%32,32+(offset+11)%32
  backbone:=[]nci03WEdge{{nodes[0],nodes[1]},{nodes[1],nodes[2]},{nodes[2],nodes[3]}}
  motifs:=[4][2]nci03WEdge{
   {{d1,d2},{d2,nodes[1]}},
   {{d1,d2},{d2,d1}},
   {{d1,nodes[2]},{d2,d1}},
   {{d1,nodes[3]},{d2,nodes[0]}},
  }
  variant:=seed%4
  relations:=append(backbone,motifs[variant][:]...)
  store:=newUPLM0CRecall()
  trial:=NCI03WTrial{Seed:seed,DistinctStructure:offset,TopologyClass:variant}
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
  edges:=nci03WEdges(store)
  trial.TwoHop=nci03WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  trial.ThreeHop=nci03WReach(join,cut,edges,nodes[0],nodes[3],3,&result.JoinCandidateEvaluations)
  trial.TwoHopNegative=!nci03WReach(join,cut,edges,nodes[0],nodes[3],2,&result.JoinCandidateEvaluations)
  trial.ThreeHopNegative=!nci03WReach(join,cut,edges,nodes[0],nodes[4],3,&result.JoinCandidateEvaluations)
  trial.DirectOnly=nci03WContains(edges,nodes[0],nodes[2])
  trial.OracleTwoHop=nci03WContains(edges,nodes[0],nodes[1])&&nci03WContains(edges,nodes[1],nodes[2])
  restored,err:=nci01WCheckpoint(store)
  if err!=nil{return result,err}
  result.CheckpointRestores++
  edges=nci03WEdges(restored)
  trial.Interrupted=nci03WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  p,_:=nci01WAdmission(admit,.5,restored,nci01WClause(nodes[3],nodes[4],true))
  result.OriginalAdmissionCalls++;if p{tp++}else{fn++}
  edges=nci03WEdges(restored)
  trial.Retained=nci03WReach(join,cut,edges,nodes[0],nodes[2],2,&result.JoinCandidateEvaluations)
  trial.NewFourHop=nci03WReach(join,cut,edges,nodes[0],nodes[4],4,&result.JoinCandidateEvaluations)
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
 result.MotifAccuracy=map[string]map[string]float64{}
 motifPassed:=true
 for m:=0;m<4;m++{
  local:=[]NCI03WTrial{}
  for _,row:=range result.Records{if row.TopologyClass==m{local=append(local,row)}}
  if len(local)!=6{return result,fmt.Errorf("motif cardinality mismatch")}
  rates:=map[string]float64{}
  for _,k:=range []string{"two_hop","three_hop","two_hop_negative","three_hop_negative","interrupted","new_four_hop"}{
   hits:=0
   for _,trial:=range local{
    switch k{
    case "two_hop":if trial.TwoHop{hits++}
    case "three_hop":if trial.ThreeHop{hits++}
    case "two_hop_negative":if trial.TwoHopNegative{hits++}
    case "three_hop_negative":if trial.ThreeHopNegative{hits++}
    case "interrupted":if trial.Interrupted{hits++}
    case "new_four_hop":if trial.NewFourHop{hits++}
    }
   }
   rates[k]=float64(hits)/float64(len(local))
   if hits<5{motifPassed=false}
  }
  result.MotifAccuracy[fmt.Sprintf("%d",m)]=rates
 }
 result.Verdict="NATIVE_JOIN_MOTIF_TRANSFER_FALSIFIED"
 if motifPassed{result.Verdict="NATIVE_JOIN_MOTIF_TRANSFER_BOUNDED"}
 return result,nil
}
