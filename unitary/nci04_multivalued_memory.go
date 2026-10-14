package unitary
import ("fmt")

// NCI04 falsifies multivalued association handling. This does not change the
// frozen NCI02 learned join, lexical admission, or exact native recall.
type nci04WCase struct {
 Seed int `json:"seed"`
 Motif string `json:"motif"`
 OriginalTwoHop bool `json:"original_two_hop"`
 OriginalThreeHop bool `json:"original_three_hop"`
 AlternateDirect bool `json:"alternate_direct"`
 InterruptedOriginal bool `json:"interrupted_original"`
 OldAfterNew bool `json:"old_after_new"`
 NewFourHop bool `json:"new_four_hop"`
 OracleOriginalTwoHop bool `json:"oracle_original_two_hop"`
 OracleOriginalThreeHop bool `json:"oracle_original_three_hop"`
 DirectOnlyBaseline bool `json:"direct_only_baseline"`
 ActiveEntries int `json:"active_entries"`
}
func nci04WOracle(edges []nci02WEdge,from,to,depth int)bool {
 frontier:=map[int]bool{from:true}
 for hop:=0;hop<depth;hop++{
  next:=map[int]bool{}
  for _,e:=range edges{if frontier[e.From]{next[e.To]=true}}
  frontier=next
 }
 return frontier[to]
}
func RunNCI04Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("bad phase")}
 join,cut,dev:=nci02WTrain()
 result:=map[string]any{
  "schema":"wingless.nci04.multivalued-native-memory-falsification.v1",
  "phase":phase,"source_parent":"4ec09c56e6964b9bb65d2c4c85388f026d0755f3",
  "development_join_accuracy":dev,"join_params":6,
  "admission_params":65,"native_recall_cap":16,"experimental_active_cap":7,
  "sealed_seed_range":[]int{23744,23768},"case_motifs":[]string{
     "root_fork","internal_fork","no_fork_control"},
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "learned_graph_search":false,"learned_multivalued_storage":false,
  "nci01_full_acceptance":false,"verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 admit:=uplm0fTrainClassifier()
 cases:=make([]nci04WCase,0,72)
 summary:=map[string]map[string]int{}
 ops,admissionCalls,restores,peak:=0,0,0,0
 for seed:=23744;seed<23768;seed++{
  for family:=0;family<3;family++{
   shift:=(seed*5+family*3)%16
   nodes:=[5]int{}
   for j:=0;j<5;j++{nodes[j]=64+(shift+3*j)%16}
   extra,noiseA,noiseB:=80+(seed+family)%5,86+(seed+family)%5,92+(seed+family)%5
   edges:=[]nci02WEdge{{nodes[0],nodes[1]},{nodes[1],nodes[2]},
     {nodes[2],nodes[3]}}
   alternateFrom:=nodes[0]
   if family==1{alternateFrom=nodes[1]}
   if family==2{alternateFrom=97}
   edges=append(edges,nci02WEdge{alternateFrom,extra},nci02WEdge{noiseA,noiseB})
   oracleEdges:=append([]nci02WEdge(nil),edges...)
   memory:=newUPLM0CRecall()
   for _,e:=range edges {
      nci01WAdmission(admit,.5,memory,nci01WClause(e.From,e.To,true))
      admissionCalls++
   }
   nci01WAdmission(admit,.5,memory,nci01WClause(98,99,false))
   admissionCalls++
   active:=nci02WEdges(memory)
   motif:=[3]string{"root_fork","internal_fork","no_fork_control"}[family]
   rec:=nci04WCase{Seed:seed,Motif:motif}
   rec.OriginalTwoHop=nci02WReach(join,cut,active,nodes[0],nodes[2],2,&ops)
   rec.OriginalThreeHop=nci02WReach(join,cut,active,nodes[0],nodes[3],3,&ops)
   rec.AlternateDirect=nci02WContains(active,alternateFrom,extra)
   rec.OracleOriginalTwoHop=nci04WOracle(oracleEdges,nodes[0],nodes[2],2)
   rec.OracleOriginalThreeHop=nci04WOracle(oracleEdges,nodes[0],nodes[3],3)
   rec.DirectOnlyBaseline=nci02WContains(active,nodes[0],nodes[2])
   restored,err:=nci01WCheckpoint(memory)
   if err!=nil{return nil,err}
   restores++
   active=nci02WEdges(restored)
   rec.InterruptedOriginal=nci02WReach(join,cut,active,nodes[0],nodes[2],2,&ops)
   nci01WAdmission(admit,.5,restored,nci01WClause(nodes[3],nodes[4],true))
   admissionCalls++
   active=nci02WEdges(restored)
   rec.OldAfterNew=nci02WReach(join,cut,active,nodes[0],nodes[2],2,&ops)
   rec.NewFourHop=nci02WReach(join,cut,active,nodes[0],nodes[4],4,&ops)
   rec.ActiveEntries=len(restored.order)
   if rec.ActiveEntries>peak{peak=rec.ActiveEntries}
   if rec.ActiveEntries>7{return nil,fmt.Errorf("experimental active budget exceeded")}
   fields:=map[string]bool{
      "original_two_hop":rec.OriginalTwoHop,
      "original_three_hop":rec.OriginalThreeHop,
      "alternate_direct":rec.AlternateDirect,
      "interrupted_original":rec.InterruptedOriginal,
      "old_after_new":rec.OldAfterNew,
      "new_four_hop":rec.NewFourHop,
      "oracle_original_two_hop":rec.OracleOriginalTwoHop,
      "oracle_original_three_hop":rec.OracleOriginalThreeHop,
      "direct_only_baseline":rec.DirectOnlyBaseline,
   }
   bucket:=summary[motif]
   if bucket==nil{bucket=map[string]int{};summary[motif]=bucket}
   for k,v:=range fields{if v{bucket[k]++}}
   cases=append(cases,rec)
  }
 }
 if len(cases)!=72{return nil,fmt.Errorf("wrong sealed trial count")}
 accuracy:=map[string]map[string]float64{}
 for motif,fields:=range summary{
  accuracy[motif]=map[string]float64{}
  for _,key:=range []string{"original_two_hop","original_three_hop",
    "alternate_direct","interrupted_original","old_after_new","new_four_hop",
    "oracle_original_two_hop","oracle_original_three_hop","direct_only_baseline"}{
      accuracy[motif][key]=float64(fields[key])/24
  }
 }
 control:=accuracy["no_fork_control"]["original_two_hop"]
 forkA:=accuracy["root_fork"]["original_two_hop"]
 forkB:=accuracy["internal_fork"]["original_two_hop"]
 verdict:="MULTIVALUED_MEMORY_LIMITATION_NOT_ISOLATED"
 if control>=.85 && forkA<.50 && forkB<.50 {
  verdict="FUNCTIONAL_MAP_FORK_RETENTION_FAILURE"
 } else if control>=.85 && forkA>=.85 && forkB>=.85 {
  verdict="MULTIVALUED_TRANSFER_DEMONSTRATED"
 }
 result["verdict"]=verdict
 result["records"]=cases
 result["accuracy_by_motif"]=accuracy
 result["peak_entries"]=peak
 result["checkpoint_restores"]=restores
 result["join_candidate_evaluations"]=ops
 result["admission_calls"]=admissionCalls
 return result,nil
}
