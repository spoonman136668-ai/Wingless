package unitary
import("fmt")
type nci08Case struct{
 Seed int `json:"seed"`; OrderFamily int `json:"order_family"`
 GraphA bool `json:"graph_a"`;GraphB bool `json:"graph_b"`
 PreviousA bool `json:"previous_a"`;PreviousB bool `json:"previous_b"`
 FIFOA bool `json:"fifo_a"`;FIFOB bool `json:"fifo_b"`
 LRUA bool `json:"lru_a"`;LRUB bool `json:"lru_b"`
 OracleA bool `json:"oracle_a"`;OracleB bool `json:"oracle_b"`
 NativeFacts int `json:"native_facts"`;GraphEdges int `json:"graph_edges"`
}
func nci08Fixture(seed,fam int)([2][5]int,[]nci06Fact,[]nci02WEdge){
 var chains [2][5]int
 for i:=0;i<5;i++{
  chains[0][i]=64+((seed*3+i*5)%16)
  chains[1][i]=96+((seed*7+i*3)%16)
 }
 edges:=[]nci02WEdge{}
 for c:=0;c<2;c++{for i:=0;i<4;i++{
  edges=append(edges,nci02WEdge{From:chains[c][i],To:chains[c][i+1]})
 }}
 for i:=0;i<4;i++{edges=append(edges,nci02WEdge{
  From:128+i*4+(seed%3),To:150+i*4+(seed%3)})}
 orders:=[4][12]int{
  {0,4,1,5,2,6,3,7,8,9,10,11},
  {0,1,2,3,4,5,6,7,8,9,10,11},
  {4,5,6,7,0,1,2,3,8,9,10,11},
  {8,0,4,9,1,5,10,2,6,11,3,7}}
 facts:=[]nci06Fact{}
 for t,idx:=range orders[fam]{
  repetitions:=1+(seed+t)%4
  facts=append(facts,nci06Fact{Edge:edges[idx],Exposure:repetitions,Inserted:t,Last:t})
 }
 return chains,facts,edges
}
func RunNCI08Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("phase")}
 graphWeight,nGraph,devGraph:=nci07WTrain()
 oldWeight,nOld,devOld:=nci06Train()
 result:=map[string]any{
 "schema":"wingless.nci08.counterfactual-dual-chain-memory-bound.v1",
 "phase":phase,"source_parent":"3dd82ec6fcd10e134b788ca0f6058284c4f668f1",
 "frozen_graph_parameters":4,"frozen_previous_parameters":3,
 "graph_dev_examples":nGraph,"prior_dev_examples":nOld,
 "graph_dev_accuracy":devGraph,"prior_dev_accuracy":devOld,
 "sealed_seed_range":[]int{28000,28024},
 "paired_chain_queries_per_state":2,"edge_memory_limit":7,
 "full_oracle_minimum_edges":8,"theoretical_paired_accuracy_ceiling":.5,
 "new_model_training":false,"learned_relational_storage":false,"learned_graph_search":false,
 "nci01_full_acceptance":false,"external_llm_calls":0,"accepted_ref_mutation":false,
 "verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 admit:=uplm0fTrainClassifier()
 join,cut,_:=nci02WTrain()
 counts:=map[string]int{}
 work:=map[string]int{"graph":0,"previous":0,"fifo":0,"lru":0}
 trials:=[]nci08Case{}
 admitted,missed,scoring,rejections:=0,0,0,0
 for seed:=28000;seed<28024;seed++{for family:=0;family<4;family++{
  chains,raw,oracleEdges:=nci08Fixture(seed,family)
  stream:=[]nci06Fact{}
  for _,f:=range raw{
   if admit.prob(uplm0fEncode(nci01WClause(f.Edge.From,f.Edge.To,true)))>=.5{
    stream=append(stream,f);admitted++
   }else{missed++}
  }
  graph,used,rejected:=nci07WRetain(stream,graphWeight)
  scoring+=used;rejections+=rejected
  older,_,_:=nci06Retain(stream,"learned",oldWeight)
  fifo,_,_:=nci06Retain(stream,"fifo",oldWeight)
  lru,_,_:=nci06Retain(stream,"lru",oldWeight)
  if len(graph)!=7||len(older)!=7||len(fifo)!=7||len(lru)!=7{
   return nil,fmt.Errorf("native admission or capacity changed")
  }
  row:=nci08Case{Seed:seed,OrderFamily:family,NativeFacts:len(stream),GraphEdges:len(graph)}
  evaluate:=func(edges []nci02WEdge,mode string)(bool,bool){
   counter:=0
   a:=nci02WReach(join,cut,edges,chains[0][0],chains[0][4],4,&counter)
   b:=nci02WReach(join,cut,edges,chains[1][0],chains[1][4],4,&counter)
   work[mode]+=counter;return a,b
  }
  row.GraphA,row.GraphB=evaluate(graph,"graph")
  row.PreviousA,row.PreviousB=evaluate(older,"previous")
  row.FIFOA,row.FIFOB=evaluate(fifo,"fifo")
  row.LRUA,row.LRUB=evaluate(lru,"lru")
  row.OracleA=nci04WOracle(oracleEdges,chains[0][0],chains[0][4],4)
  row.OracleB=nci04WOracle(oracleEdges,chains[1][0],chains[1][4],4)
  if (row.GraphA&&row.GraphB)||(row.PreviousA&&row.PreviousB)||
   (row.FIFOA&&row.FIFOB)||(row.LRUA&&row.LRUB){
   return nil,fmt.Errorf("paired-chain storage bound violation")
  }
  for k,v:=range map[string]bool{
   "graph_a":row.GraphA,"graph_b":row.GraphB,"previous_a":row.PreviousA,
   "previous_b":row.PreviousB,"fifo_a":row.FIFOA,"fifo_b":row.FIFOB,
   "lru_a":row.LRUA,"lru_b":row.LRUB,"oracle_a":row.OracleA,"oracle_b":row.OracleB}{
    if v{counts[k]++}
  }
  trials=append(trials,row)
 }}
 if len(trials)!=96{return nil,fmt.Errorf("wrong cohort")}
 scores:=map[string]float64{}
 for _,k:=range []string{"graph_a","graph_b","previous_a","previous_b",
  "fifo_a","fifo_b","lru_a","lru_b","oracle_a","oracle_b"}{
  scores[k]=float64(counts[k])/96
 }
 paired:=map[string]float64{}
 for _,prefix:=range []string{"graph","previous","fifo","lru","oracle"}{
  paired[prefix]=(scores[prefix+"_a"]+scores[prefix+"_b"])/2
 }
 verdict:="DUAL_CHAIN_MEMORY_BOUND_NOT_QUALIFIED"
 if scores["oracle_a"]==1&&scores["oracle_b"]==1{
  verdict="CAPACITY_SYMMETRY_LIMIT_CONFIRMED"
 }
 for _,arm:=range []string{"graph","previous","fifo","lru"}{
  if paired[arm]>.5+1e-9{return nil,fmt.Errorf("physical edge cap violated")}
 }
 family:=map[string]map[string]float64{}
 for f:=0;f<4;f++{
  grouped:=map[string]float64{}
  for _,arm:=range []string{"graph","previous","fifo","lru","oracle"}{
   hits:=0
   for _,row:=range trials{if row.OrderFamily!=f{continue}
    switch arm {
     case "graph":if row.GraphA{hits++};if row.GraphB{hits++}
     case "previous":if row.PreviousA{hits++};if row.PreviousB{hits++}
     case "fifo":if row.FIFOA{hits++};if row.FIFOB{hits++}
     case "lru":if row.LRUA{hits++};if row.LRUB{hits++}
     case "oracle":if row.OracleA{hits++};if row.OracleB{hits++}
    }
   }
   grouped[arm]=float64(hits)/48
  }
  family[fmt.Sprintf("%d",f)]=grouped
 }
 result["records"]=trials;result["query_accuracy"]=scores
 result["paired_query_accuracy"]=paired
 result["paired_by_family"]=family
 result["join_candidate_evaluations"]=work
 result["graph_score_feature_multiplies"]=scoring
 result["incoming_rejections"]=rejections
 result["native_admission_true_positive"]=admitted
 result["native_admission_false_negative"]=missed
 result["verdict"]=verdict
 return result,nil
}
