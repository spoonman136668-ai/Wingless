package unitary

import ("fmt";"math";"time";"unsafe")

// W-F10 probes whether a frozen 4-hop reasoning circuit can honor a query
// for another depth, using pairs whose graph and endpoints are identical.
// It is a DEVELOPMENT capability audit, not F3 heldout/first-use clearance.
type wf10Observation struct {
 Seed int `json:"seed"`
 PathDepth int `json:"path_depth"`
 RequestedDepths [2]int `json:"requested_depths"`
 Query [2]int `json:"query"`
 Edges [12][2]int `json:"edges"`
 Truth [2]bool `json:"truth"`
 Models map[string][2]bool `json:"models"`
 FrozenProbability float64 `json:"frozen_probability"`
}
func wf10Oracle(edges [12]nci02WEdge,source,target,depth int)bool {
 reach:=map[int]bool{source:true}
 for step:=0;step<depth;step++{
  next:=map[int]bool{}
  for _,e:=range edges { if reach[e.From]{next[e.To]=true} }
  reach=next
 }
 return reach[target]
}
func wf10Pair(seed,k int)(wf1Datum,[2]int,error){
 var d wf1Datum
 if k<4||k>7{return d,[2]int{},fmt.Errorf("unregistered exact path depth")}
 var raw [12][2]int
 for i:=0;i<k;i++{raw[i]=[2]int{i,i+1}}
 for i:=k;i<12;i++{
  n:=i-k
  if n<7 {raw[i]=[2]int{8+n,8+(n+1)%7}} else {raw[i]=[2]int{8,10}}
 }
 for i:=11;i>0;i--{
  j:=(seed*31+k*17+i*29)%(i+1)
  raw[i],raw[j]=raw[j],raw[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range raw{
  a,b:=wf1Rename(e[0],seed,2),wf1Rename(e[1],seed,2)
  if seen[[2]int{a,b}] {return d,[2]int{},fmt.Errorf("duplicate edge")}
  seen[[2]int{a,b}]=true
  d.Edges[i]=nci02WEdge{From:a,To:b}
 }
 start,end:=wf1Rename(0,seed,2),wf1Rename(k,seed,2)
 if wf10Oracle(d.Edges,start,end,k-1)||!wf10Oracle(d.Edges,start,end,k){
  return d,[2]int{},fmt.Errorf("depth-conditioned oracle does not flip")
 }
 if start==end{return d,[2]int{},fmt.Errorf("source target collision")}
 d.Questions=[2][2]int{{start,end},{start,end}}
 d.Labels=[2]bool{false,true}
 for i,a:=range d.Edges{for j,b:=range d.Edges{if a.To==b.From{d.Eq[i][j]=1}}}
 return d,[2]int{start,end},nil
}
func RunWFoundationF10()(map[string]any,error){
 start:=time.Now()
 frozen,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 composer:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
 -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 old:=wf5Matcher()
 results:=[]wf10Observation{}
 names:=[]string{"frozen_learned_composition_depth_blind","old_F5_depth_blind","coded_depth_aware_oracle","always_yes_shortcut"}
 stats:=map[string][3]int{}
 for seed:=51100;seed<51124;seed++{
  for depth:=4;depth<=7;depth++{
   data,pair,err:=wf10Pair(seed,depth);if err!=nil{return nil,err}
   mem:=wf7Memory(data,frozen)
   previous:=wf4Memory(data,old)
   prob:=wf7Answer(mem,pair,frozen,composer)
   oldProb:=wf5Readout(previous,pair,old,composer)
   row:=wf10Observation{Seed:seed,PathDepth:depth,
    RequestedDepths:[2]int{depth-1,depth},Query:pair,Truth:[2]bool{false,true},
    FrozenProbability:prob,Models:map[string][2]bool{}}
   for i,e:=range data.Edges{row.Edges[i]=[2]int{e.From,e.To}}
   row.Models["frozen_learned_composition_depth_blind"]=[2]bool{prob>=.5,prob>=.5}
   row.Models["old_F5_depth_blind"]=[2]bool{oldProb>=.5,oldProb>=.5}
   row.Models["coded_depth_aware_oracle"]=[2]bool{
    wf10Oracle(data.Edges,pair[0],pair[1],depth-1),
    wf10Oracle(data.Edges,pair[0],pair[1],depth)}
   row.Models["always_yes_shortcut"]=[2]bool{true,true}
   for _,name:=range names {
    prediction:=row.Models[name]
    v:=stats[name]
    if !prediction[0]{v[0]++}
    if prediction[1]{v[1]++}
    if prediction==row.Truth{v[2]++}
    stats[name]=v
   }
   results=append(results,row)
  }
 }
 if len(results)!=96||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("cohort or state changed")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  c:=stats[name]
  scores[name]=map[string]float64{"negative_correct":float64(c[0]),
   "positive_correct":float64(c[1]),"paired_correct":float64(c[2]),
   "accuracy":float64(c[0]+c[1])/192,"paired_accuracy":float64(c[2])/96}
 }
 if math.IsNaN(results[0].FrozenProbability){return nil,fmt.Errorf("nonfinite frozen output")}
 status:="INVALID_DEPTH_ORACLE"
 if scores["coded_depth_aware_oracle"]["paired_correct"]==96{
  status="DEPTH_CONDITIONING_MISSING_FROM_FROZEN_INPUT_INTERFACE"
  if scores["frozen_learned_composition_depth_blind"]["paired_correct"]!=0{
   status="INVALID_DEPTH_INVARIANCE_NOT_PROVEN"
  }
 }
 return map[string]any{
  "schema":"wingless.f10.depth_conditioned_capability_audit.v1",
  "stage":"DEVELOPMENT_ONLY","source_F9_sha":"d79f29806be4746dee219dfabf77a5abc2ecd14e",
  "original_F7_weights_sha256":"c7df241ca99627ea816091a84d321819828328a02a56a507cd7862d8afe67fe4",
  "seeds":[51100,51124],"exact_chain_lengths":[]int{4,5,6,7},
  "graph_states":96,"paired_questions":192,"negative_questions":96,
  "state_bytes":576,"model_retrained":false,"weights_adjusted":false,
  "query_depth_is_present_in_frozen_model_interface":false,
  "same_query_and_edges_per_depth_pair":true,"explicit_depth_is_only_control_input":true,
  "scores":scores,"rows":results,"registered_verdict":status,
  "source_birth_independent":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "elapsed_ms":time.Since(start).Milliseconds(),
 },nil
}
