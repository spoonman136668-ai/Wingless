package unitary

import ("encoding/json";"fmt";"math";"unsafe")

// W-NCI10: train only a supervised merge-selection gate.
// The path representation, adjacency prefilter and query traversal remain code.
// This is not a learned state encoder or general reasoning engine.
type nci10WChunk struct {
 From uint16
 To uint16
 Span uint8
 Exposure uint8
 Time uint8
}
type nci10WDatum struct {
 f nci06Fact
 main bool // development supervision only; never an inference input
}
type nci10WModel struct { W [4]float64; B float64 }
type nci10WCase struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 LearnedA bool `json:"learned_a"`
 LearnedB bool `json:"learned_b"`
 CodedA bool `json:"coded_a"`
 CodedB bool `json:"coded_b"`
 GraphA bool `json:"graph_a"`
 GraphB bool `json:"graph_b"`
 OracleA bool `json:"oracle_a"`
 OracleB bool `json:"oracle_b"`
 LearnedChunks int `json:"learned_chunks"`
 CodedChunks int `json:"coded_chunks"`
 LearnedCoverage int `json:"learned_coverage"`
 LearnedMerges int `json:"learned_merges"`
 CodedMerges int `json:"coded_merges"`
}
func nci10WFixture(seed,fam int)([2][5]int,[]nci10WDatum,[]nci02WEdge) {
 var chains [2][5]int
 for k:=0;k<2;k++{for i:=0;i<5;i++{
  if k==0{chains[k][i]=64+(seed*3+5*i)%16
  }else{chains[k][i]=96+(seed*7+3*i)%16}
 }}
 edges:=make([]nci02WEdge,0,12)
 for c:=0;c<2;c++{for i:=0;i<4;i++{
  edges=append(edges,nci02WEdge{From:chains[c][i],To:chains[c][i+1]})
 }}
 edges=append(edges,nci02WEdge{From:chains[0][1],To:146+seed%3},
  nci02WEdge{From:chains[1][2],To:152+seed%3},
  nci02WEdge{From:160+seed%3,To:chains[0][2]},
  nci02WEdge{From:168+seed%3,To:chains[1][1]})
 orders:=[4][12]int{
  {0,4,8,1,5,9,2,6,10,3,7,11},
  {8,0,9,4,1,10,5,2,11,6,3,7},
  {0,1,8,4,5,9,2,6,10,3,7,11},
  {8,4,0,9,5,1,10,6,2,11,7,3},
 }
 facts:=make([]nci10WDatum,0,12)
 for t,index:=range orders[fam]{
  target:=index<8
  exposure:=1;if target{exposure=4}
  if fam==1{if target{exposure=1}else{exposure=4}}
  if fam==2{exposure=2}
  if fam==3{if index%2==0{exposure=4}else{exposure=1}}
  facts=append(facts,nci10WDatum{
   f:nci06Fact{Edge:edges[index],Exposure:exposure,Inserted:t,Last:t},
   main:target})
 }
 return chains,facts,edges
}
func nci10WFeatures(old,new nci10WChunk,considered []nci10WChunk)[4]float64{
 degree:=0
 for _,q:=range considered{
  if q.Span!=1{continue}
  if q.From==old.From||q.From==old.To||q.From==new.From||q.From==new.To||
    q.To==old.From||q.To==old.To||q.To==new.From||q.To==new.To{degree++}
 }
 return [4]float64{
  float64(min(old.Exposure,new.Exposure))/4,
  float64(int(old.Exposure)+int(new.Exposure))/8,
  math.Abs(float64(int(old.Time)-int(new.Time)))/12,
  float64(min(degree,4))/4,
 }
}
func nci10WProb(w nci10WModel,x [4]float64)float64{
 z:=w.B;for i,v:=range x{z+=w.W[i]*v}
 return 1/(1+math.Exp(-max(-60.,min(60.,z))))
}
type nci10WExample struct{x [4]float64;y float64}
func nci10WExamples(lo,hi int)[]nci10WExample{
 out:=[]nci10WExample{}
 for seed:=lo;seed<hi;seed++{
  for _,fam:=range []int{0,3}{
   _,facts,_:=nci10WFixture(seed,fam)
   observed:=[]nci10WChunk{}
   targets:=[]bool{}
   for _,f:=range facts{
    c:=nci10WChunk{From:uint16(f.f.Edge.From),To:uint16(f.f.Edge.To),
       Span:1,Exposure:uint8(f.f.Exposure),Time:uint8(f.f.Inserted)}
    for i,old:=range observed{
     if old.To==c.From||c.To==old.From{
      label:=0.;if targets[i]&&f.main{label=1}
      out=append(out,nci10WExample{nci10WFeatures(old,c,append(append([]nci10WChunk(nil),observed...),c)),label})
     }
    }
    observed=append(observed,c);targets=append(targets,f.main)
   }
  }
 }
 return out
}
func nci10WTrain()(nci10WModel,float64,int,int,float64){
 train:=nci10WExamples(30976,31024)
 valid:=nci10WExamples(31072,31088)
 if len(train)==0||len(valid)==0{panic("NCI10 dev labels missing")}
 model:=nci10WModel{}
 for epoch:=0;epoch<160;epoch++{
  for _,v:=range train{
   difference:=nci10WProb(model,v.x)-v.y
   for i,x:=range v.x{model.W[i]-=.08*difference*x}
   model.B-=.08*difference
  }
 }
 bestAcc,bestCut:=-1.,.5
 for _,cut:=range []float64{.35,.45,.50,.55,.65,.75}{
  count:=0
  for _,item:=range valid{
   if (nci10WProb(model,item.x)>=cut)==(item.y>0){count++}
  }
  acc:=float64(count)/float64(len(valid))
  if acc>bestAcc{bestAcc,bestCut=acc,cut}
 }
 return model,bestCut,len(train),len(valid),bestAcc
}
func nci10WLearned(stream []nci06Fact, model nci10WModel,cut float64,
 retention [4]float64)([]nci10WChunk,int,int,int,int){
 stored:=[]nci10WChunk{}
 merges,rejections,scoreMults,gateCalls:=0,0,0,0
 for _,f:=range stream{
  edge:=f.Edge
  item:=nci10WChunk{From:uint16(edge.From),To:uint16(edge.To),
    Span:1,Exposure:uint8(f.Exposure),Time:uint8(f.Inserted)}
  selected,best:=-1,cut
  for i,old:=range stored{
   if old.Span!=1||!(old.To==item.From||item.To==old.From){continue}
   candidates:=append(append([]nci10WChunk(nil),stored...),item)
   score:=nci10WProb(model,nci10WFeatures(old,item,candidates))
   gateCalls++
   if score>=best{selected,best=i,score}
  }
  if selected>=0{
   old:=stored[selected]
   next:=nci10WChunk{Span:2,Exposure:max(old.Exposure,item.Exposure),
     Time:min(old.Time,item.Time)}
   if old.To==item.From{next.From=old.From;next.To=item.To
   }else{next.From=item.From;next.To=old.To}
   stored[selected]=next;merges++
  }else{stored=append(stored,item)}
  if len(stored)>7{
   considered:=[]nci06Fact{}
   for _,q:=range stored{
    considered=append(considered,nci06Fact{
      Edge:nci02WEdge{From:int(q.From),To:int(q.To)},
      Exposure:int(q.Exposure),Inserted:int(q.Time)})
   }
   worst,minScore:=0,math.Inf(1)
   for i,q:=range considered{
    score:=nci07WProb(retention,nci07WFeatures(q,considered))
    scoreMults+=4
    if score<minScore{worst,minScore=i,score}
   }
   if worst==len(stored)-1{rejections++}
   stored=append(stored[:worst],stored[worst+1:]...)
  }
 }
 return stored,merges,rejections,scoreMults,gateCalls
}
func nci10WChunksToNCI09(c []nci10WChunk)[]nci09WChunk{
 out:=[]nci09WChunk{}
 for _,v:=range c{
  out=append(out,nci09WChunk{From:int(v.From),To:int(v.To),
   Span:int(v.Span),Exposure:int(v.Exposure),Inserted:int(v.Time)})
 }
 return out
}
func nci10WCheckpoint(c []nci10WChunk)([]nci10WChunk,error){
 raw,err:=json.Marshal(c);if err!=nil{return nil,err}
 var o []nci10WChunk
 if err=json.Unmarshal(raw,&o);err!=nil{return nil,err}
 return o,nil
}
func RunNCI10Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("invalid phase")}
 model,cut,trainN,validN,dev:=nci10WTrain()
 result:=map[string]any{
  "schema":"wingless.nci10.learned-merge-choice.v1",
  "phase":phase,"source_parent":"5e6bd3416d275c731c2a14060fcef09c99e7bb22",
  "training_examples":trainN,"validation_examples":validN,
  "development_accuracy":dev,"threshold":cut,
  "merge_weights":model.W,"merge_bias":model.B,
  "merge_gate_parameters":5,"relation_features_hand_coded":true,
  "compressed_representation_learned":false,"learned_graph_search":false,
  "chunk_record_bytes":unsafe.Sizeof(nci10WChunk{}),
  "baseline_edge_record_bytes":unsafe.Sizeof(nci02WEdge{}),
  "experimental_slot_cap":7,"baseline_slot_cap":7,
  "sealed_seed_range":[]int{31424,31448},
  "families":[]string{"development_like","reversed_exposure","uniform_exposure","alternating_decoys"},
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
 for seed:=31424;seed<31448;seed++{
  for family:=0;family<4;family++{
   chains,raw,full:=nci10WFixture(seed,family)
   stream:=[]nci06Fact{}
   for _,fact:=range raw{
    if admit.prob(uplm0fEncode(nci01WClause(fact.f.Edge.From,fact.f.Edge.To,true)))>=.5{
     stream=append(stream,fact.f);admitted++
    }else{missed++}
   }
   learned,lm,reject,scored,logits:=nci10WLearned(stream,model,cut,retention)
   coded,cm,_,_:=nci09WCompress(stream,retention)
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
 verdict:="NO_LEARNED_MERGE_TRANSFER"
 if accuracy["oracle"]<1{verdict="NOVEL_GRAPH_ORACLE_INVALID"
 }else if accuracy["learned"]>=.70&&delta>=.10&&novelSafe{
  verdict="LEARNED_MERGE_GATE_TRANSFER"
 }else if delta>=.05{verdict="PARTIAL_LEARNED_MERGE_GAIN"}
 result["verdict"]=verdict;result["records"]=cases
 result["paired_accuracy"]=accuracy;result["family_accuracy"]=familyAccuracy
 result["learned_minus_coded"]=delta;result["join_candidate_evaluations"]=cost
 result["merges"]=merged;result["gate_inference_evaluations"]=mergeLogits
 result["graph_feature_multiplies"]=multiply;result["incoming_rejections"]=rejected
 result["admission_tp"]=admitted;result["admission_fn"]=missed
 result["checkpoint_restores"]=restores
 return result,nil
}
