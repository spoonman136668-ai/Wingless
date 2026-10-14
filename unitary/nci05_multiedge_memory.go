package unitary

import (
 "encoding/json"
 "fmt"
)

type nci05WGraph struct {
 Edges []nci02WEdge `json:"edges"`
}
func (g *nci05WGraph) Add(a,b int) error {
 for _,e:=range g.Edges{if e.From==a&&e.To==b{return nil}}
 if len(g.Edges)>=7{return fmt.Errorf("NCI05 seven-edge experimental cap exceeded")}
 g.Edges=append(g.Edges,nci02WEdge{From:a,To:b})
 return nil
}
func nci05WRestore(g *nci05WGraph)(*nci05WGraph,error){
 payload,e:=json.Marshal(g);if e!=nil{return nil,e}
 var restored nci05WGraph
 if e:=json.Unmarshal(payload,&restored);e!=nil{return nil,e}
 return &restored,nil
}
type NCI05WTrial struct{
 Seed int `json:"seed"`
 Motif string `json:"motif"`
 MultiTwo bool `json:"multiedge_two"`
 MultiThree bool `json:"multiedge_three"`
 MultiTwoNegative bool `json:"multiedge_two_negative"`
 MultiThreeNegative bool `json:"multiedge_three_negative"`
 MultiInterrupted bool `json:"multiedge_interrupted"`
 MultiRetained bool `json:"multiedge_retained"`
 MultiNewFour bool `json:"multiedge_new_four"`
 BaselineTwo bool `json:"baseline_two"`
 BaselineThree bool `json:"baseline_three"`
 BaselineRetained bool `json:"baseline_retained"`
 OracleTwo bool `json:"oracle_two"`
 OracleThree bool `json:"oracle_three"`
 DirectOnly bool `json:"direct_only"`
 AdmissionTP int `json:"admission_tp"`
 AdmissionFP int `json:"admission_fp"`
 AdmissionFN int `json:"admission_fn"`
 ActiveEdges int `json:"active_edges"`
 BaselineEntries int `json:"baseline_entries"`
}
func nci05WFixture(seed,kind int)([5]int,[]nci02WEdge){
 var nodes [5]int
 step:=[4]int{1,3,5,7}[kind]
 shift:=(seed*7+kind*3)%16
 for i:=0;i<5;i++{nodes[i]=64+(shift+step*i)%16}
 extra:=[6]int{}
 for i:=0;i<6;i++{extra[i]=80+3*i+(seed+i+kind)%3}
 edges:=[]nci02WEdge{
  {From:nodes[0],To:nodes[1]},
  {From:nodes[1],To:nodes[2]},
  {From:nodes[2],To:nodes[3]},
 }
 switch kind {
 case 0: edges=append(edges,nci02WEdge{From:nodes[0],To:extra[0]},
         nci02WEdge{From:extra[2],To:extra[3]},nci02WEdge{From:extra[4],To:extra[5]})
 case 1: edges=append(edges,nci02WEdge{From:nodes[1],To:extra[0]},
         nci02WEdge{From:extra[2],To:extra[3]},nci02WEdge{From:extra[4],To:extra[5]})
 case 2: edges=append(edges,nci02WEdge{From:nodes[0],To:extra[0]},
         nci02WEdge{From:nodes[1],To:extra[1]},nci02WEdge{From:extra[4],To:extra[5]})
 default:edges=append(edges,nci02WEdge{From:extra[0],To:extra[1]},
         nci02WEdge{From:extra[2],To:extra[3]},nci02WEdge{From:extra[4],To:extra[5]})
 }
 return nodes,edges
}
func RunNCI05Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("bad phase")}
 learned,cut,dev:=nci02WTrain()
 result:=map[string]any{
  "schema":"wingless.nci05.native-multiedge-capacity-integration.v1",
  "phase":phase,"source_parent":"e2931315bb498c30d9313dd7a8ccbdea74aea2e9",
  "join_dev_accuracy":dev,"join_threshold":cut,
  "native_join_weights":6,"native_admission_weights":65,
  "experimental_edge_cap":7,"native_baseline_slots":16,
  "sealed_seed_range":[]int{24320,24344},
  "motifs":[]string{"root_fork","internal_fork","dual_fork","nonfork_control"},
  "original_storage_learned":false,"multiedge_storage_learned":false,
  "learned_graph_search":false,"nci01_full_acceptance":false,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 admit:=uplm0fTrainClassifier()
 names:=[]string{"root_fork","internal_fork","dual_fork","nonfork_control"}
 trials:=make([]NCI05WTrial,0,96)
 metrics:=map[string]map[string]int{}
 multiOps,baselineOps,checkpoints,admissionCalls,peak:=0,0,0,0,0
 for seed:=24320;seed<24344;seed++{
  for family:=0;family<4;family++{
   nodes,edges:=nci05WFixture(seed,family)
   raw:=&nci05WGraph{}
   baseline:=newUPLM0CRecall()
   row:=NCI05WTrial{Seed:seed,Motif:names[family]}
   observe:=func(g *nci05WGraph,store *uplm0cRecall,a,b int,storeTruth bool)error{
    clause:=nci01WClause(a,b,storeTruth)
    predicted:=admit.prob(uplm0fEncode(clause))>=.5
    admissionCalls++
    if predicted&&storeTruth{row.AdmissionTP++}
    if predicted&&!storeTruth{row.AdmissionFP++}
    if !predicted&&storeTruth{row.AdmissionFN++}
    if predicted{
      if e:=g.Add(a,b);e!=nil{return e}
      store.write(fmt.Sprintf("k%02d",a),fmt.Sprintf("v%02d",b))
    }
    return nil
   }
   for _,e:=range edges{if err:=observe(raw,baseline,e.From,e.To,true);err!=nil{return nil,err}}
   if err:=observe(raw,baseline,96,97,false);err!=nil{return nil,err}
   orig:=nci02WEdges(baseline)
   reach:=func(list []nci02WEdge,hops int,dest int,isMulti bool)bool{
     if isMulti{return nci02WReach(learned,cut,list,nodes[0],dest,hops,&multiOps)}
     return nci02WReach(learned,cut,list,nodes[0],dest,hops,&baselineOps)
   }
   row.MultiTwo=reach(raw.Edges,2,nodes[2],true)
   row.MultiThree=reach(raw.Edges,3,nodes[3],true)
   row.MultiTwoNegative=!reach(raw.Edges,2,nodes[3],true)
   row.MultiThreeNegative=!reach(raw.Edges,3,nodes[4],true)
   row.BaselineTwo=reach(orig,2,nodes[2],false)
   row.BaselineThree=reach(orig,3,nodes[3],false)
   row.DirectOnly=nci02WContains(orig,nodes[0],nodes[2])
   row.OracleTwo=nci04WOracle(edges,nodes[0],nodes[2],2)
   row.OracleThree=nci04WOracle(edges,nodes[0],nodes[3],3)
   resumed,err:=nci05WRestore(raw);if err!=nil{return nil,err}
   resumeBaseline,err:=nci01WCheckpoint(baseline);if err!=nil{return nil,err}
   checkpoints++
   row.MultiInterrupted=reach(resumed.Edges,2,nodes[2],true)
   if err:=observe(resumed,resumeBaseline,nodes[3],nodes[4],true);err!=nil{return nil,err}
   row.MultiRetained=reach(resumed.Edges,2,nodes[2],true)
   row.MultiNewFour=reach(resumed.Edges,4,nodes[4],true)
   row.BaselineRetained=reach(nci02WEdges(resumeBaseline),2,nodes[2],false)
   row.ActiveEdges=len(resumed.Edges)
   row.BaselineEntries=len(resumeBaseline.order)
   if row.ActiveEdges>peak{peak=row.ActiveEdges}
   if row.ActiveEdges>7||row.BaselineEntries>7{return nil,fmt.Errorf("bounded memory failure")}
   bucket:=metrics[row.Motif];if bucket==nil{bucket=map[string]int{};metrics[row.Motif]=bucket}
   for _,entry:=range []struct{name string;ok bool}{
     {"multiedge_two",row.MultiTwo},{"multiedge_three",row.MultiThree},
     {"multiedge_two_negative",row.MultiTwoNegative},
     {"multiedge_three_negative",row.MultiThreeNegative},
     {"multiedge_interrupted",row.MultiInterrupted},
     {"multiedge_retained",row.MultiRetained},
     {"multiedge_new_four",row.MultiNewFour},
     {"baseline_two",row.BaselineTwo},{"baseline_three",row.BaselineThree},
     {"baseline_retained",row.BaselineRetained},
     {"oracle_two",row.OracleTwo},{"oracle_three",row.OracleThree},
     {"direct_only",row.DirectOnly},
   }{if entry.ok{bucket[entry.name]++}}
   trials=append(trials,row)
  }
 }
 if len(trials)!=96{return nil,fmt.Errorf("trial count")}
 acc:=map[string]map[string]float64{}
 for _,motif:=range names{
  acc[motif]=map[string]float64{}
  for _,key:=range []string{"multiedge_two","multiedge_three","multiedge_two_negative",
   "multiedge_three_negative","multiedge_interrupted","multiedge_retained",
   "multiedge_new_four","baseline_two","baseline_three","baseline_retained",
   "oracle_two","oracle_three","direct_only"} {
    acc[motif][key]=float64(metrics[motif][key])/24
  }
 }
 success:=true
 for _,motif:=range names{
  for _,task:=range []string{"multiedge_two","multiedge_three","multiedge_two_negative",
   "multiedge_three_negative","multiedge_interrupted","multiedge_retained","multiedge_new_four"}{
    if acc[motif][task]<.90{success=false}
  }
  if acc[motif]["oracle_two"]<.99||acc[motif]["oracle_three"]<.99{success=false}
 }
 originalBroken:=acc["root_fork"]["baseline_two"]<.5&&
   acc["internal_fork"]["baseline_two"]<.5&&acc["dual_fork"]["baseline_two"]<.5
 verdict:="BOUNDED_MULTIEDGE_STORAGE_DID_NOT_RESTORE_COMPOSITION"
 if success&&originalBroken {
   verdict="BOUNDED_MULTIEDGE_STORAGE_RESTORES_CODED_COMPOSITION"
 } else if success {
   verdict="MULTIEDGE_COMPONENT_TRANSFER_WITHOUT_BASELINE_CONTRAST"
 }
 result["records"]=trials
 result["accuracy_by_motif"]=acc
 result["multi_join_candidate_evaluations"]=multiOps
 result["baseline_join_candidate_evaluations"]=baselineOps
 result["checkpoint_restores"]=checkpoints
 result["admission_calls"]=admissionCalls
 result["peak_multiedge_entries"]=peak
 result["verdict"]=verdict
 return result,nil
}
