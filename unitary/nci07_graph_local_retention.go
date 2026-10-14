package unitary
import("fmt";"math")
type nci07WTrial struct{
 Seed int `json:"seed"`; Family int `json:"family"`
 GraphTwo bool `json:"graph_two"`; GraphFour bool `json:"graph_four"`
 GraphRestored bool `json:"graph_restored"`
 PriorTwo bool `json:"prior_two"`; PriorFour bool `json:"prior_four"`
 FIFOFour bool `json:"fifo_four"`; LRUFour bool `json:"lru_four"`
 OracleFour bool `json:"oracle_four"`; GraphEdges int `json:"graph_edges"`
 NewFactsRejected int `json:"new_facts_rejected"`
}
func nci07WFeatures(f nci06Fact,candidates []nci06Fact)[4]float64{
 adjacent:=0
 for _,other:=range candidates{
  if other.Edge==f.Edge{continue}
  if f.Edge.To==other.Edge.From||f.Edge.From==other.Edge.To{adjacent++}
 }
 return [4]float64{1,float64(f.Exposure)/4,float64(f.Inserted)/12,float64(adjacent)/3}
}
func nci07WProb(w [4]float64,x [4]float64)float64{
 z:=0.;for i,k:=range x{z+=w[i]*k}
 if z>35{z=35};if z< -35{z= -35}
 return 1/(1+math.Exp(-z))
}
func nci07WTrain()([4]float64,int,float64){
 type ex struct{x [4]float64;label float64}
 rows:=[]ex{}
 for seed:=26800;seed<26960;seed++{for family:=0;family<4;family++{
  _,facts,chain:=nci06Fixture(seed,family)
  for t:=7;t<12;t++{
   current:=facts[t-7:t+1]
   for _,v:=range current{
    y:=0.
    for _,e:=range chain{if e==v.Edge{y=1}}
    rows=append(rows,ex{nci07WFeatures(v,current),y})
   }
  }
 }}
 w:=[4]float64{}
 for epoch:=0;epoch<45;epoch++{
  for _,v:=range rows{
   delta:=.035*(nci07WProb(w,v.x)-v.label)
   for k:=range w{w[k]-=delta*v.x[k]}
  }
 }
 correct,n:=0,0
 for seed:=26960;seed<26992;seed++{for fam:=0;fam<4;fam++{
  _,facts,chain:=nci06Fixture(seed,fam)
  for t:=7;t<12;t++{for _,v:=range facts[t-7:t+1]{
   truth:=false;for _,e:=range chain{if v.Edge==e{truth=true}}
   if (nci07WProb(w,nci07WFeatures(v,facts[t-7:t+1]))>=.5)==truth{correct++}
   n++
  }}
 }}
 return w,len(rows),float64(correct)/float64(n)
}
func nci07WRetain(stream []nci06Fact,w [4]float64)([]nci02WEdge,int,int){
 kept:=[]nci06Fact{};scoring,rejected:=0,0
 for _,incoming:=range stream{
  if len(kept)<7{kept=append(kept,incoming);continue}
  candidates:=append(append([]nci06Fact(nil),kept...),incoming)
  worst:=0;lowest:=math.Inf(1)
  for idx,f:=range candidates{
   x:=nci07WFeatures(f,candidates)
   value:=nci07WProb(w,x);scoring+=4
   if value<lowest{lowest=value;worst=idx}
  }
  if worst==len(candidates)-1{rejected++;continue}
  kept=append(candidates[:worst],candidates[worst+1:]...)
 }
 out:=[]nci02WEdge{}
 for _,f:=range kept{out=append(out,f.Edge)}
 return out,scoring,rejected
}
func RunNCI07Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("phase")}
 weights,n,dev:=nci07WTrain()
 result:=map[string]any{
  "schema":"wingless.nci07.graph-local-retention-learning.v1","phase":phase,
  "parent":"c3195345f82ca4e0bd23c6015f27b4f7f7ac2d44",
  "trained_retention_weights":weights,"development_examples":n,
  "validation_accuracy":dev,"retention_parameters":4,
  "frozen_previous_retention_parameters":3,
  "memory_cap":7,"sealed_seeds":[]int{27136,27160},
  "graph_local_features_learned":false,
  "graph_retention_choice_learned":true,
  "learned_graph_search":false,"nci01_full_acceptance":false,
  "external_llm_calls":0,"accepted_ref_mutation":false,"verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 oldWeights,_,_:=nci06Train()
 join,cut,_:=nci02WTrain()
 admit:=uplm0fTrainClassifier()
 rows:=[]nci07WTrial{}
 work:=map[string]int{"graph":0,"prior":0,"fifo":0,"lru":0}
 accepted,missed,scoreMult,rejected,restores:=0,0,0,0,0
 counts:=map[string]int{}
 for seed:=27136;seed<27160;seed++{for family:=0;family<4;family++{
  nodes,raw,chain:=nci06Fixture(seed,family)
  stream:=[]nci06Fact{}
  for _,fact:=range raw{
   if admit.prob(uplm0fEncode(nci01WClause(fact.Edge.From,fact.Edge.To,true)))>=.5{
    stream=append(stream,fact);accepted++
   }else{missed++}
  }
  graph,used,dropped:=nci07WRetain(stream,weights)
  scoreMult+=used;rejected+=dropped
  previous,_,_:=nci06Retain(stream,"learned",oldWeights)
  fifo,_,_:=nci06Retain(stream,"fifo",oldWeights)
  lru,_,_:=nci06Retain(stream,"lru",oldWeights)
  if len(graph)!=7||len(previous)!=7||len(fifo)!=7||len(lru)!=7{
   return nil,fmt.Errorf("memory budget or native admission failed")
  }
  a,b,c,d:=0,0,0,0
  row:=nci07WTrial{Seed:seed,Family:family,GraphEdges:len(graph),NewFactsRejected:dropped}
  row.GraphTwo=nci02WReach(join,cut,graph,nodes[0],nodes[2],2,&a)
  row.GraphFour=nci02WReach(join,cut,graph,nodes[0],nodes[4],4,&a)
  row.PriorTwo=nci02WReach(join,cut,previous,nodes[0],nodes[2],2,&b)
  row.PriorFour=nci02WReach(join,cut,previous,nodes[0],nodes[4],4,&b)
  row.FIFOFour=nci02WReach(join,cut,fifo,nodes[0],nodes[4],4,&c)
  row.LRUFour=nci02WReach(join,cut,lru,nodes[0],nodes[4],4,&d)
  copy,err:=nci06Reload(graph);if err!=nil{return nil,err};restores++
  row.GraphRestored=nci02WReach(join,cut,copy,nodes[0],nodes[4],4,&a)
  row.OracleFour=nci04WOracle(chain,nodes[0],nodes[4],4)
  work["graph"]+=a;work["prior"]+=b;work["fifo"]+=c;work["lru"]+=d
  for k,v:=range map[string]bool{
   "graph_two":row.GraphTwo,"graph_four":row.GraphFour,
   "graph_restored":row.GraphRestored,"prior_two":row.PriorTwo,
   "prior_four":row.PriorFour,"fifo_four":row.FIFOFour,
   "lru_four":row.LRUFour,"oracle_four":row.OracleFour}{if v{counts[k]++}}
  rows=append(rows,row)
 }}
 acc:=map[string]float64{}
 for _,k:=range []string{"graph_two","graph_four","graph_restored","prior_two",
  "prior_four","fifo_four","lru_four","oracle_four"}{acc[k]=float64(counts[k])/96}
 fam:=map[string]map[string]float64{}
 for family:=0;family<4;family++{
  m:=map[string]float64{}
  for _,k:=range []string{"graph_four","prior_four","lru_four"}{
   success:=0
   for _,r:=range rows{if r.Family!=family{continue}
    if k=="graph_four"&&r.GraphFour{success++}
    if k=="prior_four"&&r.PriorFour{success++}
    if k=="lru_four"&&r.LRUFour{success++}
   }
   m[k]=float64(success)/24
  }
  fam[fmt.Sprintf("%d",family)]=m
 }
 gainOld:=acc["graph_four"]-acc["prior_four"]
 gainLRU:=acc["graph_four"]-acc["lru_four"]
 noCollapse:=true;for _,f:=range fam{if f["graph_four"]<f["prior_four"]-.20{noCollapse=false}}
 verdict:="NO_GRAPH_LOCAL_RETENTION_TRANSFER"
 if acc["graph_four"]>=.60&&gainOld>=.15&&gainLRU>=.10&&noCollapse{
  verdict="GRAPH_LOCAL_LEARNED_RETENTION_TRANSFER"
 }else if gainOld>=.10{verdict="PARTIAL_GRAPH_LOCAL_RETENTION_GAIN"}
 result["records"]=rows;result["accuracy"]=acc;result["by_family"]=fam
 result["gain_vs_previous_learned"]=gainOld;result["gain_vs_lru"]=gainLRU
 result["join_candidate_evaluations"]=work
 result["score_feature_multiplies"]=scoreMult
 result["new_fact_rejections"]=rejected
 result["native_admission_true_positive"]=accepted
 result["native_admission_false_negative"]=missed
 result["checkpoint_restores"]=restores;result["verdict"]=verdict
 return result,nil
}
