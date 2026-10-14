package unitary
import ("encoding/json";"fmt";"math";"unsafe")

// NCI09 deliberately tests coded online two-hop path chunks, not learned
// compression. This makes native relational content explicit under a bounded
// seven-chunk record budget. Original learned admission/join remain frozen.
type nci09WChunk struct{
 From int `json:"from"`
 To int `json:"to"`
 Span int `json:"span"`
 Exposure int `json:"exposure"`
 Inserted int `json:"inserted"`
}
type nci09WTrial struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 CompressedA bool `json:"compressed_a"`
 CompressedB bool `json:"compressed_b"`
 GraphA bool `json:"graph_a"`
 GraphB bool `json:"graph_b"`
 OracleA bool `json:"oracle_a"`
 OracleB bool `json:"oracle_b"`
 ChunkCount int `json:"chunk_count"`
 CompressedCoverage int `json:"compressed_coverage"`
 MergedCount int `json:"merged_count"`
 GraphEdges int `json:"graph_edges"`
}
func nci09WCompress(facts []nci06Fact,weights [4]float64)([]nci09WChunk,int,int,int){
 kept:=[]nci09WChunk{}
 merged,rejected,featureMultiplies:=0,0,0
 for _,fact:=range facts{
  incoming:=nci09WChunk{From:fact.Edge.From,To:fact.Edge.To,Span:1,
   Exposure:fact.Exposure,Inserted:fact.Inserted}
  pair:=-1
  for i,existing:=range kept{
   if existing.Span==1 && (existing.To==incoming.From||incoming.To==existing.From){
    pair=i;break
   }
  }
  if pair>=0{
   old:=kept[pair]
   result:=nci09WChunk{Span:2,
     Exposure:max(old.Exposure,incoming.Exposure),
     Inserted:min(old.Inserted,incoming.Inserted)}
   if old.To==incoming.From{
    result.From=old.From;result.To=incoming.To
   }else{result.From=incoming.From;result.To=old.To}
   kept[pair]=result;merged++
  }else{kept=append(kept,incoming)}
  if len(kept)>7{
   candidates:=make([]nci06Fact,0,len(kept))
   for _,chunk:=range kept{
    candidates=append(candidates,nci06Fact{Edge:nci02WEdge{From:chunk.From,To:chunk.To},
        Exposure:chunk.Exposure,Inserted:chunk.Inserted})
   }
   worst:=0;minScore:=math.Inf(1)
   for i,item:=range candidates{
    score:=nci07WProb(weights,nci07WFeatures(item,candidates))
    featureMultiplies+=4
    if score<minScore{minScore=score;worst=i}
   }
   if worst==len(kept)-1{rejected++}
   kept=append(kept[:worst],kept[worst+1:]...)
  }
 }
 return kept,merged,rejected,featureMultiplies
}
func nci09WRestore(items []nci09WChunk)([]nci09WChunk,error){
 bytes,err:=json.Marshal(items);if err!=nil{return nil,err}
 var next []nci09WChunk
 if err=json.Unmarshal(bytes,&next);err!=nil{return nil,err}
 return next,nil
}
func nci09WReach(join nci02WJoin,cut float64,items []nci09WChunk,start,dest,length int,work *int)bool{
 // Existing trained join evaluates adjacent stored macro boundaries; hop spans
 // and exact boundary matching remain coded, not learned.
 var seek func(node,used int,last nci02WEdge,hasPrev bool)bool
 seek=func(node,used int,last nci02WEdge,hasPrev bool)bool{
  for _,c:=range items{
   if c.From!=node||c.Span<1||c.Span>2||used+c.Span>length{continue}
   edge:=nci02WEdge{From:c.From,To:c.To}
   if hasPrev{
    *work++
    if nci02WScore(join,last,edge)<cut{continue}
   }
   total:=used+c.Span
   if total==length&&c.To==dest{return true}
   if total<length&&seek(c.To,total,edge,true){return true}
  }
  return false
 }
 return seek(start,0,nci02WEdge{},false)
}
func RunNCI09Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("unknown phase")}
 retention,n,dev:=nci07WTrain()
 result:=map[string]any{
  "schema":"wingless.nci09.two-link-chunk-composition.v1",
  "phase":phase,"source_parent":"a42e5eecfc8742c38430ca629744d500fa707ba0",
  "sealed_seed_range":[]int{30112,30136},
  "retention_learned_weights":4,"learned_join_weights":6,
  "dev_examples":n,"dev_accuracy":dev,
  "memory_slots":7,"chunk_size_bytes":unsafe.Sizeof(nci09WChunk{}),
  "original_edge_size_bytes":unsafe.Sizeof(nci02WEdge{}),
  "compression_rule_learned":false,"span_search_learned":false,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "nci01_full_acceptance":false,"verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{
  tiny:=[]nci06Fact{
   {Edge:nci02WEdge{From:1,To:2},Exposure:1,Inserted:0},
   {Edge:nci02WEdge{From:2,To:3},Exposure:1,Inserted:1},
  }
  chunks,count,_,_:=nci09WCompress(tiny,retention)
  if count!=1||len(chunks)!=1||chunks[0].Span!=2{return nil,fmt.Errorf("invalid dev compression")}
  return result,nil
 }
 admit:=uplm0fTrainClassifier()
 join,cut,_:=nci02WTrain()
 trials:=[]nci09WTrial{}
 summary:=map[string]int{}
 families:=map[string]map[string]int{}
 ops:=map[string]int{"compressed":0,"graph":0}
 merges,discarded,scoreMuls,admitted,missed,restores:=0,0,0,0,0,0
 for seed:=30112;seed<30136;seed++{
  for fam:=0;fam<4;fam++{
   chains,raw,oracle:=nci08Fixture(seed,fam)
   stream:=[]nci06Fact{}
   for _,f:=range raw{
    if admit.prob(uplm0fEncode(nci01WClause(f.Edge.From,f.Edge.To,true)))>=.5{
     stream=append(stream,f);admitted++
    }else{missed++}
   }
   chunks,count,denied,used:=nci09WCompress(stream,retention)
   graph,_,_:=nci07WRetain(stream,retention)
   merges+=count;discarded+=denied;scoreMuls+=used
   if len(chunks)>7||len(graph)!=7{return nil,fmt.Errorf("NCI09 native occupancy violation")}
   copy,err:=nci09WRestore(chunks);if err!=nil{return nil,err}
   restores++
   row:=nci09WTrial{Seed:seed,Family:fam,ChunkCount:len(copy),GraphEdges:len(graph),MergedCount:count}
   for _,chunk:=range copy{row.CompressedCoverage+=chunk.Span}
   compressedWork,graphWork:=0,0
   row.CompressedA=nci09WReach(join,cut,copy,chains[0][0],chains[0][4],4,&compressedWork)
   row.CompressedB=nci09WReach(join,cut,copy,chains[1][0],chains[1][4],4,&compressedWork)
   row.GraphA=nci02WReach(join,cut,graph,chains[0][0],chains[0][4],4,&graphWork)
   row.GraphB=nci02WReach(join,cut,graph,chains[1][0],chains[1][4],4,&graphWork)
   ops["compressed"]+=compressedWork;ops["graph"]+=graphWork
   row.OracleA=nci04WOracle(oracle,chains[0][0],chains[0][4],4)
   row.OracleB=nci04WOracle(oracle,chains[1][0],chains[1][4],4)
   key:=fmt.Sprintf("%d",fam)
   if families[key]==nil{families[key]=map[string]int{}}
   for _,metric:=range []struct{name string;ok bool}{
    {"compressed",row.CompressedA},{"compressed",row.CompressedB},
    {"graph",row.GraphA},{"graph",row.GraphB},
    {"oracle",row.OracleA},{"oracle",row.OracleB},
   } {
    if metric.ok{summary[metric.name]++;families[key][metric.name]++}
   }
   trials=append(trials,row)
  }
 }
 if len(trials)!=96{return nil,fmt.Errorf("original record count")}
 score:=map[string]float64{}
 for _,name:=range []string{"compressed","graph","oracle"}{score[name]=float64(summary[name])/192}
 byFamily:=map[string]map[string]float64{}
 for family,stats:=range families{
  byFamily[family]=map[string]float64{}
  for _,name:=range []string{"compressed","graph","oracle"}{
   byFamily[family][name]=float64(stats[name])/48
  }
 }
 verdict:="NO_TWO_LINK_CHUNK_COMPOSITION_GAIN"
 if score["oracle"]==1&&score["compressed"]>=.75&&
   score["compressed"]-score["graph"]>=.15{
  verdict="CODED_TWO_LINK_COMPRESSION_RESTORES_BOUNDED_RECALL"
 }else if score["compressed"]-score["graph"]>=.05{
  verdict="PARTIAL_TWO_LINK_CHUNK_RECALL_GAIN"
 }
 result["verdict"]=verdict;result["records"]=trials
 result["paired_accuracy"]=score;result["per_family"]=byFamily
 result["join_candidate_evaluations"]=ops;result["compression_merges"]=merges
 result["incoming_rejections"]=discarded
 result["graph_scorer_feature_multiplies"]=scoreMuls
 result["admission_true_positive"]=admitted;result["admission_false_negative"]=missed
 result["checkpoint_restores"]=restores
 return result,nil
}
