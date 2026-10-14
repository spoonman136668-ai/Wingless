package unitary

import ("fmt";"sort";"math")

// NCI10-R2 freezes an outcome-trained linear merge controller.
// Its five numeric parameters and threshold are selected solely with
// development delayed-question utility. Span memory/retrieval stay programmed.
type nci10R2Candidate struct {
 Weights nci10WModel
 Cut float64
 TrainScore float64
 TrainWork int
 Index int
}
func nci10R2Evaluate(weights nci10WModel,threshold float64,
 lo,hi int,retention [4]float64, join nci02WJoin,joinCut float64)(float64,int){
 score,work:=0,0
 for seed:=lo;seed<hi;seed++{
  for family:=0;family<4;family++{
   chains,raw,_:=nci10WFixture(seed,family)
   stream:=[]nci06Fact{}
   for _,f:=range raw{stream=append(stream,f.f)}
   chunks,_,_,_,gate:=nci10WLearned(stream,weights,threshold,retention)
   work+=gate
   state:=nci10WChunksToNCI09(chunks)
   workLocal:=0
   for k:=0;k<2;k++{
    if nci09WReach(join,joinCut,state,chains[k][0],chains[k][4],4,&workLocal){score++}
   }
   work+=workLocal
  }
 }
 return float64(score)/float64((hi-lo)*8),work
}
func nci10WTrainByFutureUtility()(nci10WModel,float64,int,int,float64){
 retention,_,_:=nci07WTrain()
 join,cut,_:=nci02WTrain()
 candidateWeights:=[][5]float64{
  {0,0,0,0,0},{2,1,0,0,-1},{-2,-1,0,0,1},
  {0,0,-2,2,0},{0,0,2,-2,0},{0,0,0,0,-1},
  {0,0,0,0,1},{1,0,-1,1,-.5},{-1,0,1,-1,.5},
 }
 policies:=[]nci10R2Candidate{}
 for k,v:=range candidateWeights{
  m:=nci10WModel{W:[4]float64{v[0],v[1],v[2],v[3]},B:v[4]}
  for _,threshold:=range []float64{.4,.6}{
   score,work:=nci10R2Evaluate(m,threshold,32400,32408,retention,join,cut)
   policies=append(policies,nci10R2Candidate{m,threshold,score,work,k})
  }
 }
 sort.SliceStable(policies,func(i,j int)bool{
  a,b:=policies[i],policies[j]
  if a.TrainScore!=b.TrainScore{return a.TrainScore>b.TrainScore}
  if a.TrainWork!=b.TrainWork{return a.TrainWork<b.TrainWork}
  if a.Index!=b.Index{return a.Index<b.Index}
  return a.Cut<b.Cut
 })
 best,validation,bestWork:=policies[0],-1.,int(^uint(0)>>1)
 for _,p:=range policies[:4]{
  acc,work:=nci10R2Evaluate(p.Weights,p.Cut,32432,32440,retention,join,cut)
  if acc>validation||(acc==validation&&work<bestWork){
   best,validation,bestWork=p,acc,work
  }
 }
 return best.Weights,best.Cut,len(policies)*8*8,4*8*8,validation
}

// Same NCI09 first-compatible span2 merge and graph-local overflow scoring,
// stored in identical 8-byte packed NCI10 chunks as the learned controller.
func nci10R2Coded(stream []nci06Fact,retention [4]float64)([]nci10WChunk,int){
 stored:=[]nci10WChunk{}
 merges:=0
 for _,f:=range stream{
  c:=nci10WChunk{From:uint16(f.Edge.From),To:uint16(f.Edge.To),
    Span:1,Exposure:uint8(f.Exposure),Time:uint8(f.Inserted)}
  pair:=-1
  for i,old:=range stored{
   if old.Span==1&&(old.To==c.From||c.To==old.From){pair=i;break}
  }
  if pair>=0{
   old:=stored[pair]
   next:=nci10WChunk{Span:2,Exposure:max(old.Exposure,c.Exposure),
     Time:min(old.Time,c.Time)}
   if old.To==c.From{next.From=old.From;next.To=c.To
   }else{next.From=c.From;next.To=old.To}
   stored[pair]=next;merges++
  }else{stored=append(stored,c)}
  if len(stored)>7{
   candidates:=[]nci06Fact{}
   for _,v:=range stored{
    candidates=append(candidates,nci06Fact{
      Edge:nci02WEdge{From:int(v.From),To:int(v.To)},
      Exposure:int(v.Exposure),Inserted:int(v.Time)})
   }
   worst,minScore:=0,math.Inf(1)
   for i,item:=range candidates{
    score:=nci07WProb(retention,nci07WFeatures(item,candidates))
    if score<minScore{worst,minScore=i,score}
   }
   stored=append(stored[:worst],stored[worst+1:]...)
  }
 }
 return stored,merges
}
func RunNCI10R2Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("invalid phase")}
 model,cut,trainN,validN,dev:=nci10WTrainByFutureUtility()
 result:=map[string]any{
  "schema":"wingless.nci10r2.future-query-utility-merge-selection.v1",
  "phase":phase,"source_parent":"db2e93d7f7e38d0d0e9cc50462033313fd1b7d1c",
  "training_examples":trainN,"validation_examples":validN,
  "development_accuracy":dev,"threshold":cut,
  "merge_weights":model.W,"merge_bias":model.B,
  "merge_gate_parameters":5,"relation_features_hand_coded":true,
  "compressed_representation_learned":false,"learned_graph_search":false,
  "chunk_record_bytes":unsafe.Sizeof(nci10WChunk{}),
  "baseline_edge_record_bytes":unsafe.Sizeof(nci02WEdge{}),
  "experimental_slot_cap":7,"baseline_slot_cap":7,
  "sealed_seed_range":[]int{32640,32664},
  "families":[]string{"development_like","reversed_exposure","uniform_exposure","alternating_decoys"},
  "training_objective":"development future two-question accuracy; no future query at inference",
  "utility_candidate_policies":18,"utility_validation_finalists":4,
  "compression_representation_learned":false,
  "nci01_full_acceptance":false,"external_llm_calls":0,
  "accepted_ref_mutation":false,"verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{
  _,raw,_:=nci10WFixture(30970,0)
  s:=[]nci06Fact{};for _,x:=range raw{s=append(s,x.f)}
  ret,_,_,_,_:=nci10WLearned(s,model,cut,[4]float64{})
  if len(ret)>7{return nil,fmt.Errorf("native dev cap")}
  return result,nil
 }
 retention,_,_:=nci07WTrain()
 admit:=uplm0fTrainClassifier()
 join,joinCut,_:=nci02WTrain()
 cases:=[]nci10WCase{}
 counts:=map[string]int{}
 byFamily:=map[string]map[string]int{}
 cost:=map[string]int{"learned":0,"coded":0,"graph":0}
 merged:=map[string]int{"learned":0,"coded":0}
 multiply,mergeLogits,rejected,admitted,missed,restores:=0,0,0,0,0,0
 for seed:=32640;seed<32664;seed++{
  for family:=0;family<4;family++{
   chains,raw,full:=nci10WFixture(seed,family)
   stream:=[]nci06Fact{}
   for _,fact:=range raw{
    if admit.prob(uplm0fEncode(nci01WClause(fact.f.Edge.From,fact.f.Edge.To,true)))>=.5{
     stream=append(stream,fact.f);admitted++
    }else{missed++}
   }
   learned,lm,reject,scored,logits:=nci10WLearned(stream,model,cut,retention)
   codedPacked,cm:=nci10R2Coded(stream,retention)
   coded:=nci10WChunksToNCI09(codedPacked)
   graph,_,_:=nci07WRetain(stream,retention)
   merged["learned"]+=lm;merged["coded"]+=cm
   multiply+=scored;mergeLogits+=logits;rejected+=reject
   if len(learned)>7||len(coded)>7||len(graph)!=7{return nil,fmt.Errorf("bad original occupancy")}
   saved,err:=nci10WCheckpoint(learned);if err!=nil{return nil,err}
   restores++
   row:=nci10WCase{Seed:seed,Family:family,LearnedChunks:len(saved),
    CodedChunks:len(coded),LearnedMerges:lm,CodedMerges:cm}
   for _,v:=range saved{row.LearnedCoverage+=int(v.Span)}
   lc,cc,gc:=0,0,0
   row.LearnedA=nci09WReach(join,joinCut,nci10WChunksToNCI09(saved),chains[0][0],chains[0][4],4,&lc)
   row.LearnedB=nci09WReach(join,joinCut,nci10WChunksToNCI09(saved),chains[1][0],chains[1][4],4,&lc)
   row.CodedA=nci09WReach(join,joinCut,coded,chains[0][0],chains[0][4],4,&cc)
   row.CodedB=nci09WReach(join,joinCut,coded,chains[1][0],chains[1][4],4,&cc)
   row.GraphA=nci02WReach(join,joinCut,graph,chains[0][0],chains[0][4],4,&gc)
   row.GraphB=nci02WReach(join,joinCut,graph,chains[1][0],chains[1][4],4,&gc)
   cost["learned"]+=lc;cost["coded"]+=cc;cost["graph"]+=gc
   row.OracleA=nci04WOracle(full,chains[0][0],chains[0][4],4)
   row.OracleB=nci04WOracle(full,chains[1][0],chains[1][4],4)
   key:=fmt.Sprintf("%d",family)
   if byFamily[key]==nil{byFamily[key]=map[string]int{}}
   for _,kv:=range []struct{k string;b bool}{
    {"learned",row.LearnedA},{"learned",row.LearnedB},
    {"coded",row.CodedA},{"coded",row.CodedB},
    {"graph",row.GraphA},{"graph",row.GraphB},
    {"oracle",row.OracleA},{"oracle",row.OracleB},
   }{
    if kv.b{counts[kv.k]++;byFamily[key][kv.k]++}
   }
   cases=append(cases,row)
  }
 }
 if len(cases)!=96{return nil,fmt.Errorf("sealed cohort wrong")}
 accuracy:=map[string]float64{}
 for _,key:=range []string{"learned","coded","graph","oracle"}{accuracy[key]=float64(counts[key])/192}
 familyAccuracy:=map[string]map[string]float64{}
 for family,summary:=range byFamily{
  familyAccuracy[family]=map[string]float64{}
  for _,name:=range []string{"learned","coded","graph","oracle"}{
    familyAccuracy[family][name]=float64(summary[name])/48
  }
 }
 delta:=accuracy["learned"]-accuracy["coded"]
 novelSafe:=true
 for _,fam:=range []string{"1","2","3"}{
  if familyAccuracy[fam]["learned"]<familyAccuracy[fam]["coded"]-.10{novelSafe=false}
 }
 verdict:="R2_NO_QUERY_UTILITY_TRANSFER"
 if accuracy["oracle"]<1{verdict="R2_NOVEL_GRAPH_ORACLE_INVALID"
 }else if accuracy["learned"]>=.70&&delta>=.10&&novelSafe{
  verdict="R2_LEARNED_QUERY_UTILITY_TRANSFER"
 }else if delta>=.05{verdict="R2_PARTIAL_QUERY_UTILITY_GAIN"}
 result["verdict"]=verdict;result["records"]=cases
 result["paired_accuracy"]=accuracy;result["family_accuracy"]=familyAccuracy
 result["learned_minus_coded"]=delta;result["join_candidate_evaluations"]=cost
 result["merges"]=merged;result["gate_inference_evaluations"]=mergeLogits
 result["graph_feature_multiplies"]=multiply;result["incoming_rejections"]=rejected
 result["admission_tp"]=admitted;result["admission_fn"]=missed
 result["checkpoint_restores"]=restores
 return result,nil
}
