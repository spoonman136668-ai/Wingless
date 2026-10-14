package unitary

import("fmt";"time";"unsafe")

// W-F12 holds the W-F11 recurrence frozen and challenges it with four novel
// branching/cyclic graph grammars. At least 3/4 grammars have a shorter
// shortcut path: answering "reachable within N" is provably insufficient for
// the negative exact-depth query. DEV ONLY and no F3 source consumption.
type wf12Row struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 ChainDepth int `json:"chain_depth"`
 RequestedDepths [2]int `json:"requested_depths"`
 Query [2]int `json:"query"`
 Edges [12][2]int `json:"edges"`
 Truth [2]bool `json:"truth"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf12Build(seed,k,family int)(wf1Datum,[2]int,error){
 var d wf1Datum
 if k<7||k>8||family<60||family>63{return d,[2]int{},fmt.Errorf("unregistered W-F12 cohort")}
 var e [12][2]int
 for i:=0;i<k;i++{e[i]=[2]int{i,i+1}}
 switch family{
 case 60:e[k]=[2]int{0,3}
 case 61:e[k]=[2]int{1,4}
 case 62:e[k]=[2]int{2,5}
 case 63:e[k]=[2]int{5,2}
 }
 for i:=k+1;i<12;i++{
  v:=20+i-k-1;e[i]=[2]int{v,v+1}
 }
 for i:=11;i>0;i--{j:=(seed*31+family*11+k*17+i*23)%(i+1);e[i],e[j]=e[j],e[i]}
 seen:=map[[2]int]bool{}
 for i,p:=range e{
  a,b:=wf1Rename(p[0],seed,2),wf1Rename(p[1],seed,2)
  if seen[[2]int{a,b}]{return d,[2]int{},fmt.Errorf("duplicate after renaming")}
  seen[[2]int{a,b}]=true
  d.Edges[i]=nci02WEdge{From:a,To:b}
 }
 q:=[2]int{wf1Rename(0,seed,2),wf1Rename(k,seed,2)}
 if wf10Oracle(d.Edges,q[0],q[1],k-1)||!wf10Oracle(d.Edges,q[0],q[1],k){
  return d,q,fmt.Errorf("exact-depth contrast changed")
 }
 if family!=63&&!wf12Within(d.Edges,q[0],q[1],k-1){
  return d,q,fmt.Errorf("nontrivial shortcut control lost")
 }
 return d,q,nil
}
func wf12Within(edges [12]nci02WEdge,src,dst,depth int)bool{
 if src==dst{return true}
 for k:=1;k<=depth;k++{if wf10Oracle(edges,src,dst,k){return true}}
 return false
}
func RunWFoundationF12()(map[string]any,error){
 started:=time.Now()
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 frozen:=wf5Composer{W:[5]float64{
  6.470679800286544,-1.3847365385924115,-3.0464843010034035,
  7.51083733991036,-3.4005406394991717}}
 original:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
  -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 names:=[]string{"frozen_WF11_depth_conditioned","frozen_F10_depth_blind",
  "coded_exact_depth","coded_reachable_within_depth","always_yes_shortcut"}
 rows:=[]wf12Row{}
 for seed:=53800;seed<53824;seed++{for fam:=60;fam<=63;fam++{for k:=7;k<=8;k++{
  d,q,err:=wf12Build(seed,k,fam);if err!=nil{return nil,err}
  state:=wf7Memory(d,enc)
  row:=wf12Row{Seed:seed,Family:fam,ChainDepth:k,
   RequestedDepths:[2]int{k-1,k},Query:q,
   Truth:[2]bool{false,true},Predictions:map[string][2]bool{}}
  for i,e:=range d.Edges{row.Edges[i]=[2]int{e.From,e.To}}
  for _,name:=range names{
   var prediction [2]bool
   for j,depth:=range row.RequestedDepths{
    switch name{
    case "frozen_WF11_depth_conditioned":prediction[j]=wf11Answer(state,q,depth,enc,frozen)>=.5
    case "frozen_F10_depth_blind":prediction[j]=wf7Answer(state,q,enc,original)>=.5
    case "coded_exact_depth":prediction[j]=wf10Oracle(d.Edges,q[0],q[1],depth)
    case "coded_reachable_within_depth":prediction[j]=wf12Within(d.Edges,q[0],q[1],depth)
    case "always_yes_shortcut":prediction[j]=true
    }
   }
   row.Predictions[name]=prediction
  }
  rows=append(rows,row)
 }}}
 if len(rows)!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F12 cohort/state size changed")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  correct,paired,fp:=0,0,0
  for _,r:=range rows{
   p:=r.Predictions[name]
   for j,y:=range r.Truth{if p[j]==y{correct++};if !y&&p[j]{fp++}}
   if p==r.Truth{paired++}
  }
  scores[name]=map[string]float64{"correct":float64(correct),
   "accuracy":float64(correct)/384,"paired_correct":float64(paired),
   "paired_accuracy":float64(paired)/192,"false_positives":float64(fp)}
 }
 verdict:="INVALID_DEPTH_REFERENCE_OR_SHORTCUT_CONTROL"
 exact:=scores["coded_exact_depth"]
 within:=scores["coded_reachable_within_depth"]
 if exact["paired_correct"]==192&&within["paired_correct"]<=48 {
  score:=scores["frozen_WF11_depth_conditioned"]
  if score["paired_accuracy"]>=.85&&score["accuracy"]>=.90&&score["false_positives"]<=19{
   verdict="DEVELOPMENT_FROZEN_VARIABLE_DEPTH_BRANCHING_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_FROZEN_VARIABLE_DEPTH_BRANCHING_TRANSFER"}
 }
 return map[string]any{
 "schema":"wingless.f12.frozen-branching-depth-dev.v1","stage":"DEVELOPMENT_ONLY",
 "source_F11_sha":"b82091fc9247502f10b3c3ecbe9f76ae689b5784",
 "source_F11_raw_sha256":"59859a3ec8211b3e80e160260acea06c9c954d5268338845965b6a96d4f43933",
 "frozen_F11_weights":frozen.W,"frozen_F11_model_retrained":false,
 "seed_bounds":[]int{53800,53824},"new_families":[]int{60,61,62,63},
 "chain_depths":[]int{7,8},"graph_states":len(rows),"questions":384,
 "negative_questions":192,"native_state_bytes":576,
 "query_depth_side_input_bytes":1,"programmed_recursion_length":true,
 "frozen_identity_projection":true,"scores":scores,"rows":rows,
 "registered_verdict":verdict,"source_birth_independent":false,
 "sealed_first_use":false,"accepted_ref_mutation":false,"external_llm_calls":0,
 "original_wall_ms":time.Since(started).Milliseconds()},nil
}
