package unitary

import("fmt";"time";"unsafe")

// W-F9: paired counterfactual graph edge deletion is a causal development
// check of frozen representation/recurrence, not an F3 sealed first use.
type wf9Row struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 Question [2]int `json:"question"`
 OriginalEdges [12][2]int `json:"original_edges"`
 RemovedBridge [2]int `json:"removed_bridge"`
 ReplacementBridge [2]int `json:"replacement_bridge"`
 Truth [2]bool `json:"truth"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf9Motif(f int)([12][2]int,error){
 switch f{
 case 32:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,6},{6,7},{7,8},{9,10},{10,11},{11,12},{12,13}},nil
 case 33:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,5},{5,6},{6,7},{7,8},{9,10},{10,11},{11,12},{12,9}},nil
 case 34:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{3,2},{2,1},{1,0},{5,6},{6,7},{7,8},{8,5},{0,9}},nil
 case 35:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,6},{6,7},{7,8},{8,9},{10,11},{11,12},{12,13}},nil
 default:return [12][2]int{},fmt.Errorf("unregistered F9 topology %d",f)
 }
}
func wf9Pair(seed,fam int)(wf1Datum,wf1Datum,[2]int,error){
 var d wf1Datum
 motif,err:=wf9Motif(fam);if err!=nil{return d,d,[2]int{},err}
 var edges [12][2]int
 for i,p:=range motif{edges[i]=[2]int{wf1Rename(p[0],seed,2),wf1Rename(p[1],seed,2)}}
 for i:=11;i>0;i--{j:=(seed*31+fam*13+i*17)%(i+1);edges[i],edges[j]=edges[j],edges[i]}
 for i,p:=range edges{d.Edges[i]=nci02WEdge{From:p[0],To:p[1]}}
 source,target:=wf1Rename(0,seed,2),wf1Rename(4,seed,2)
 d.Questions=[2][2]int{{source,target},{source,target}}
 d.Labels=[2]bool{true,true}
 altered:=d
 replacement:=wf1Rename(15,seed,2)
 wantedSrc,wantedDst:=wf1Rename(2,seed,2),wf1Rename(3,seed,2)
 count:=0
 for i,e:=range altered.Edges{
  if e.From==wantedSrc&&e.To==wantedDst{
   altered.Edges[i]=nci02WEdge{From:wantedSrc,To:replacement};count++
  }
 }
 if count!=1{return d,altered,[2]int{},fmt.Errorf("bridge absent or duplicate")}
 for i,a:=range d.Edges{for j,b:=range d.Edges{
  if a.To==b.From{d.Eq[i][j]=1}
 }}
 altered.Eq=[12][12]float64{}
 for i,a:=range altered.Edges{for j,b:=range altered.Edges{
  if a.To==b.From{altered.Eq[i][j]=1}
 }}
 if !nci04WOracle(d.Edges[:],source,target,4)||
   nci04WOracle(altered.Edges[:],source,target,4){
  return d,altered,[2]int{},fmt.Errorf("causal cut did not change exact oracle")
 }
 altered.Labels=[2]bool{false,false}
 return d,altered,[2]int{source,target},nil
}
func RunWFoundationF9()(map[string]any,error){
 started:=time.Now()
 encoder,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 old:=wf5Matcher()
 recurrent:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
   -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 coded:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 totals:=map[string][3]int{}
 rows:=[]wf9Row{}
 for seed:=50300;seed<50348;seed++{for fam:=32;fam<=35;fam++{
  base,cut,q,err:=wf9Pair(seed,fam)
  if err!=nil{return nil,err}
  a,b:=wf7Memory(base,encoder),wf7Memory(cut,encoder)
  oldA,oldB:=wf4Memory(base,old),wf4Memory(cut,old)
  codeA,codeB:=wf1Memory(base,coded),wf1Memory(cut,coded)
  row:=wf9Row{Seed:seed,Family:fam,Question:q,Truth:[2]bool{true,false},
   RemovedBridge:[2]int{wf1Rename(2,seed,2),wf1Rename(3,seed,2)},
   ReplacementBridge:[2]int{wf1Rename(2,seed,2),wf1Rename(15,seed,2)},
   Predictions:map[string][2]bool{}}
  for i,e:=range base.Edges{row.OriginalEdges[i]=[2]int{e.From,e.To}}
  for _,name:=range []string{"frozen_F7_learned_identity_and_recurrence",
   "frozen_F7_identity_coded_OR_control","old_F5_learned_identity_recurrence",
   "fixed_exact_coded_control","no_graph_always_yes_negative_control"}{
   var pred [2]bool
   switch name{
   case "frozen_F7_learned_identity_and_recurrence":
    pred=[2]bool{wf7Answer(a,q,encoder,recurrent)>=.5,
                 wf7Answer(b,q,encoder,recurrent)>=.5}
   case "frozen_F7_identity_coded_OR_control":
    pred=[2]bool{wf7CodedTraversal(a,q,encoder)>=.5,
                 wf7CodedTraversal(b,q,encoder)>=.5}
   case "old_F5_learned_identity_recurrence":
    pred=[2]bool{wf5Readout(oldA,q,old,recurrent)>=.5,
                 wf5Readout(oldB,q,old,recurrent)>=.5}
   case "fixed_exact_coded_control":
    pred=[2]bool{wf1Answer(codeA,q)>=.5,wf1Answer(codeB,q)>=.5}
   case "no_graph_always_yes_negative_control":pred=[2]bool{true,true}
   }
   record:=totals[name]
   if pred[0]{record[0]++}
   if !pred[1]{record[1]++}
   if pred==[2]bool{true,false}{record[2]++}
   totals[name]=record
   row.Predictions[name]=pred
  }
  rows=append(rows,row)
 }}
 if len(rows)!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F9 population/state drift")}
 scores:=map[string]map[string]float64{}
 for name,t:=range totals{
  scores[name]=map[string]float64{
   "base_true_positive_rate":float64(t[0])/192,
   "cut_true_negative_rate":float64(t[1])/192,
   "paired_both_correct_rate":float64(t[2])/192,
   "accuracy":float64(t[0]+t[1])/384,
   "correct":float64(t[0]+t[1]),
   "cut_false_positives":float64(192-t[1])}
 }
 verdict:="INVALID_CODED_CAUSAL_CONTROL"
 if scores["fixed_exact_coded_control"]["paired_both_correct_rate"]==1{
  candidate:=scores["frozen_F7_learned_identity_and_recurrence"]
  oldCandidate:=scores["old_F5_learned_identity_recurrence"]
  if candidate["paired_both_correct_rate"]>=.95&&candidate["accuracy"]>=.95&&
     candidate["cut_false_positives"]<=10&&
     candidate["accuracy"]-oldCandidate["accuracy"]>=.10{
   verdict="DEVELOPMENT_FROZEN_LEARNED_CAUSAL_BRIDGE_SENSITIVITY"
  }else{verdict="DEVELOPMENT_NO_FROZEN_CAUSAL_BRIDGE_SENSITIVITY"}
 }
 return map[string]any{"schema":"wingless.foundation.f9.frozen-causal-bridge-dev.v1",
  "stage":"DEVELOPMENT_ONLY","parent_f8_sha":"14e4f3eadaaf4c29bcfefc16e966c21f9ebb5379",
  "projection_source_original_sha256":"799364227498074bb3ef45f9017a67d83a721b71e7e92cc5ad6b01fc8506b140",
  "new_family_ids":[]int{32,33,34,35},"seed_bounds":[]int{50300,50348},
  "paired_graph_states":len(rows),"question_count":384,"positive_count":192,
  "negative_count":192,"native_state_bytes":576,
  "freeze_no_training":true,"source_birth_independent":false,
  "coded_exact_control_must_pass":true,"scores":scores,
  "raw_rows":rows,"registered_verdict":verdict,
  "sealed_first_use":false,"scientifically_accepted":false,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "original_wall_ms":time.Since(started).Milliseconds()},nil
}
