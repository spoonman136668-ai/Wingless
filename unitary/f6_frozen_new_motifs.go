package unitary

import("fmt";"time";"unsafe")

// W-F6: independently new STRUCTURAL GRAMMAR families 8-11 in development.
// Frozen learned F4 matcher + frozen learned F5 recurrent composition; no fit
// or threshold selection. These are author-created DEVELOPMENT fixtures and
// NOT independent third-party source-birth or sealed first-use evidence.
type wf6Graph struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 Edges [12][2]int `json:"edges"`
 Questions [2][2]int `json:"questions"`
 Truth [2]bool `json:"truth"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf6Topology(family int)([12][2]int,error){
 switch family{
 case 8: return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,6},{6,7},{7,8},{9,10},{10,11},{11,12},{12,13}},nil
 case 9: return [12][2]int{{0,1},{1,2},{2,3},{3,4},{1,5},{2,6},{5,7},{6,7},{7,8},{9,0},{8,9},{10,11}},nil
 case 10:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{4,0},{2,0},{4,5},{5,6},{6,4},{7,8},{8,9},{9,7}},nil
 case 11:return [12][2]int{{0,1},{1,2},{2,3},{3,4},{0,5},{5,2},{3,1},{2,6},{6,7},{7,4},{7,8},{8,3}},nil
 default:return [12][2]int{},fmt.Errorf("unregistered graph grammar %d",family)
 }
}
func wf6New(seed,family int)(wf1Datum,error){
 topology,err:=wf6Topology(family)
 if err!=nil{return wf1Datum{},err}
 var d wf1Datum
 var mapped [12][2]int
 for i,e:=range topology{
  mapped[i]=[2]int{wf1Rename(e[0],seed,2),wf1Rename(e[1],seed,2)}
 }
 // All original motif edges are preserved, but position is permuted.
 for i:=len(mapped)-1;i>0;i--{
  j:=(seed*17+i*29+family*3)%(i+1)
  mapped[i],mapped[j]=mapped[j],mapped[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range mapped{
  if seen[e]{return wf1Datum{},fmt.Errorf("duplicate graph edge")}
  seen[e]=true
  d.Edges[i]=nci02WEdge{From:e[0],To:e[1]}
 }
 source,target:=wf1Rename(0,seed,2),wf1Rename(4,seed,2)
 if !nci04WOracle(d.Edges[:],source,target,4){
  return wf1Datum{},fmt.Errorf("registered four-edge positive unreachable")
 }
 var possible []int
 used:=map[int]bool{}
 for _,e:=range d.Edges{used[e.From]=true;used[e.To]=true}
 for node:=1;node<=13;node++{
  id:=wf1Rename(node,seed,2)
  if id==target||!used[id]||nci04WOracle(d.Edges[:],source,id,4){continue}
  possible=append(possible,id)
 }
 if len(possible)==0{return wf1Datum{},fmt.Errorf("no structurally present hard negative")}
 negative:=possible[(seed+family)%len(possible)]
 if (seed+family)%2==0{
  d.Questions=[2][2]int{{source,target},{source,negative}}
  d.Labels=[2]bool{true,false}
 }else{
  d.Questions=[2][2]int{{source,negative},{source,target}}
  d.Labels=[2]bool{false,true}
 }
 for j,q:=range d.Questions{
  if nci04WOracle(d.Edges[:],q[0],q[1],4)!=d.Labels[j]{
   return wf1Datum{},fmt.Errorf("reference oracle disagreement")
  }
 }
 for i,a:=range d.Edges{for j,b:=range d.Edges{if a.To==b.From{d.Eq[i][j]=1}}}
 return d,nil
}
func RunWFoundationF6()(map[string]any,error){
 start:=time.Now()
 match:=wf5Matcher()
 frozen:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
  -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 init:=wf5Start()
 coded:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 names:=[]string{"frozen_F5_learned_composition","untrained_composition_negative_control",
  "frozen_F4_coded_traversal","fixed_coded_exact_control"}
 totals:=map[string][2]int{}
 rows:=[]wf6Graph{}
 negatives:=0
 for seed:=45400;seed<45448;seed++{for family:=8;family<=11;family++{
  d,err:=wf6New(seed,family)
  if err!=nil{return nil,err}
  learnedState:=wf4Memory(d,match)
  codedState:=wf1Memory(d,coded)
  row:=wf6Graph{Seed:seed,Family:family,Truth:d.Labels,
    Questions:d.Questions,Predictions:map[string][2]bool{}}
  for i,e:=range d.Edges{row.Edges[i]=[2]int{e.From,e.To}}
  for _,name:=range names{
   var answers [2]bool
   for j,q:=range d.Questions{
    switch name{
    case "frozen_F5_learned_composition":
     answers[j]=wf5Readout(learnedState,q,match,frozen)>=.5
    case "untrained_composition_negative_control":
     answers[j]=wf5Readout(learnedState,q,match,init)>=.5
    case "frozen_F4_coded_traversal":
     answers[j]=wf4Answer(learnedState,q,match)>=.5
    case "fixed_coded_exact_control":
     answers[j]=wf1Answer(codedState,q)>=.5
    }
    count:=totals[name]
    if answers[j]==d.Labels[j]{count[0]++}
    if !d.Labels[j]&&answers[j]{count[1]++}
    totals[name]=count
   }
   row.Predictions[name]=answers
  }
  for _,truth:=range d.Labels{if !truth{negatives++}}
  rows=append(rows,row)
 }}
 if len(rows)!=192||negatives!=192||unsafe.Sizeof(wf1State{})!=576{
  return nil,fmt.Errorf("F6 prospective development count drift")
 }
 scores:=map[string]map[string]float64{}
 for name,count:=range totals{
  scores[name]=map[string]float64{"correct":float64(count[0]),"accuracy":float64(count[0])/384,
   "false_positives":float64(count[1]),"negative_fp_rate":float64(count[1])/192}
 }
 candidate:=scores["frozen_F5_learned_composition"]
 code:=scores["fixed_coded_exact_control"]
 pass:=candidate["accuracy"]>=.85&&candidate["negative_fp_rate"]<=.15&&
       code["accuracy"]-candidate["accuracy"]<=.15
 verdict:="DEVELOPMENT_FROZEN_NATIVE_COMPOSITION_GRAMMAR_SHIFT_FAILURE"
 if pass{verdict="DEVELOPMENT_FROZEN_NATIVE_COMPOSITION_GRAMMAR_SHIFT_SUPPORTED"}
 return map[string]any{"schema":"wingless.foundation.f6.new-motifs-frozen-development.v1",
  "stage":"DEVELOPMENT_ONLY","frozen_F5_source_sha":"3ee67c1e7323b84dfbf52b671796ff4382e080f3",
  "frozen_F5_source_artifact_id":11931186944,
  "new_development_seed_bounds":[]int{45400,45448},
  "new_grammar_families":[]int{8,9,10,11},
  "graph_states":len(rows),"question_count":384,"negative_questions":negatives,
  "resident_state_bytes":576,
  "frozen_F4_matcher_parameter_count":17,"frozen_F5_composer_parameter_count":5,
  "trained_during_F6":false,"model_retrained":false,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "author_generated_new_topologies":true,"independently_authenticated_source_birth":false,
  "scientific_sealed_transfer":false,
  "coded_equality_in_candidate":false,"coded_OR_traversal_in_candidate":false,
  "programmed_message_rounds":3,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "development_wall_ms":time.Since(start).Milliseconds()},nil
}
