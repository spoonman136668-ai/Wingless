package unitary

import ("fmt";"math";"time";"unsafe")

// W-F11 is an explicitly query-depth-conditioned learned recurrence.
// Public IDs/embedding and an explicit bounded number of recurrence steps
// are engineered. ONLY five recurrent update/readout weights are learned.
// No coded entity equality or coded graph-OR is used in candidate inference.
// All source graphs are development-only, not authenticated F3 originals.
type wf11Case struct {
 Seed int `json:"seed"`
 ChainDepth int `json:"chain_depth"`
 OrderVariant int `json:"order_variant"`
 RequestedDepths [2]int `json:"requested_depths"`
 Query [2]int `json:"query"`
 Truth [2]bool `json:"truth"`
 Predictions map[string][2]bool `json:"predictions"`
}
func wf11Build(seed,k,variant int)(wf1Datum,[2]int,error){
 var d wf1Datum
 if k<3||k>8 {return d,[2]int{},fmt.Errorf("unsupported path depth %d",k)}
 var edges [12][2]int
 for i:=0;i<k;i++{edges[i]=[2]int{i,i+1}}
 for i:=k;i<12;i++{u:=20+i-k;edges[i]=[2]int{u,u+1}}
 for i:=11;i>0;i--{
  j:=(seed*19+k*37+i*11+variant*53)%(i+1)
  edges[i],edges[j]=edges[j],edges[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range edges {
  a,b:=wf1Rename(e[0],seed,variant+1),wf1Rename(e[1],seed,variant+1)
  if seen[[2]int{a,b}] {return d,[2]int{},fmt.Errorf("duplicate edge")}
  seen[[2]int{a,b}]=true
  d.Edges[i]=nci02WEdge{From:a,To:b}
 }
 src,dst:=wf1Rename(0,seed,variant+1),wf1Rename(k,seed,variant+1)
 if wf10Oracle(d.Edges,src,dst,k-1)||!wf10Oracle(d.Edges,src,dst,k){
  return d,[2]int{},fmt.Errorf("exact-depth reference did not flip")
 }
 d.Questions=[2][2]int{{src,dst},{src,dst}}
 d.Labels=[2]bool{false,true}
 return d,[2]int{src,dst},nil
}
func wf11Answer(mem wf1State,q [2]int,requested int,enc wf7Encoder,m wf5Composer)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{bits:=math.Float32bits(mem[i]);src[i]=int(bits>>16);dst[i]=int(bits&65535)}
 var active [12]float64
 for i:=0;i<12;i++{active[i],_,_=wf7Pair(enc,q[0],src[i])}
 // The *query* controls recurrence count; updates are learned and shared.
 // This is NOT learned halting. The schedule is a programmed inductive bias.
 for hop:=1;hop<requested;hop++{
  var next [12]float64
  for j:=0;j<12;j++{
   incoming:=0.
   for i:=0;i<12;i++{if i==j{continue}
    off:=12+i*11+j;if j>i{off--}
    incoming+=active[i]*float64(mem[off])
   }
   next[j]=wf1Prob(m.W[0]*incoming+m.W[1]*active[j]+m.W[2])
  }
  active=next
 }
 score:=0.
 for i:=0;i<12;i++{
  match,_,_:=wf7Pair(enc,dst[i],q[1])
  score+=active[i]*match
 }
 return wf1Prob(m.W[3]*score+m.W[4])
}
func wf11BCE(batch []wf1Datum,depths []int,enc wf7Encoder,params wf5Composer)float64{
 sum:=0.
 for i,d:=range batch{
  mem:=wf7Memory(d,enc)
  for j,requested:=range []int{depths[i]-1,depths[i]}{
   p:=max(.000001,min(.999999,wf11Answer(mem,d.Questions[j],requested,enc,params)))
   if d.Labels[j]{sum-=math.Log(p)}else{sum-=math.Log(1-p)}
  }
 }
 return sum/float64(len(batch)*2)
}
func wf11Train(enc wf7Encoder)(wf5Composer,int64,float64,error){
 start:=time.Now()
 var data []wf1Datum;var depths []int
 for seed:=51900;seed<51916;seed++{
  for k:=3;k<=6;k++{
   for v:=0;v<2;v++{
    d,_,err:=wf11Build(seed,k,v);if err!=nil{return wf5Composer{},0,0,err}
    data=append(data,d);depths=append(depths,k)
   }
  }
 }
 if len(data)!=128{return wf5Composer{},0,0,fmt.Errorf("train cardinality")}
 m:=wf5Start()
 var mm,vv [5]float64
 for epoch:=1;epoch<=105;epoch++{
  batch:=make([]wf1Datum,0,24); ds:=make([]int,0,24)
  for j:=0;j<24;j++{
   i:=(epoch*17+j*23)%len(data)
   batch=append(batch,data[i]);ds=append(ds,depths[i])
  }
  var grad [5]float64
  for w:=0;w<5;w++{
   plus,minus:=m,m
   plus.W[w]+=.035;minus.W[w]-=.035
   grad[w]=(wf11BCE(batch,ds,enc,plus)-wf11BCE(batch,ds,enc,minus))/.07
  }
  for w:=0;w<5;w++{
   g:=max(-10,min(10,grad[w]))
   mm[w]=.9*mm[w]+.1*g;vv[w]=.999*vv[w]+.001*g*g
   delta:=.11*(mm[w]/(1-math.Pow(.9,float64(epoch))))/
     (math.Sqrt(vv[w]/(1-math.Pow(.999,float64(epoch))))+1e-8)
   m.W[w]=max(-20,min(20,m.W[w]-delta))
  }
 }
 return m,time.Since(start).Milliseconds(),wf11BCE(data,depths,enc,m),nil
}
func RunWFoundationF11()(map[string]any,error){
 begin:=time.Now()
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 candidate,trainMs,loss,err:=wf11Train(enc);if err!=nil{return nil,err}
 untrained:=wf5Start()
 fixedFour:=wf5Composer{W:[5]float64{6.2402561172942175,-.44436093640380137,
  -2.6062431834020403,8.506846585720483,-4.236991668757422}}
 rows:=[]wf11Case{}
 metrics:=map[string][2]int{}
 names:=[]string{"learned_query_conditioned_recurrence","untrained_same_interface",
  "frozen_depth_blind_F10","coded_exact_depth_oracle","always_yes_no_graph"}
 for seed:=52300;seed<52324;seed++{for depth:=7;depth<=8;depth++{for variant:=0;variant<2;variant++{
  d,q,err:=wf11Build(seed,depth,variant);if err!=nil{return nil,err}
  state:=wf7Memory(d,enc)
  row:=wf11Case{Seed:seed,ChainDepth:depth,OrderVariant:variant,
   RequestedDepths:[2]int{depth-1,depth},Query:q,
   Truth:[2]bool{false,true},Predictions:map[string][2]bool{}}
  for _,name:=range names{
   var answers [2]bool
   for j,req:=range row.RequestedDepths{
    switch name{
    case "learned_query_conditioned_recurrence":answers[j]=wf11Answer(state,q,req,enc,candidate)>=.5
    case "untrained_same_interface":answers[j]=wf11Answer(state,q,req,enc,untrained)>=.5
    case "frozen_depth_blind_F10":answers[j]=wf7Answer(state,q,enc,fixedFour)>=.5
    case "coded_exact_depth_oracle":answers[j]=wf10Oracle(d.Edges,q[0],q[1],req)
    case "always_yes_no_graph":answers[j]=true
    }
    m:=metrics[name]
    if answers[j]==row.Truth[j]{m[0]++}
    if answers==row.Truth{m[1]++}
    metrics[name]=m
   }
   row.Predictions[name]=answers
  }
  rows=append(rows,row)
 }}}
 if len(rows)!=96||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("validation shape/memory drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  c:=metrics[name]
  paired:=0;correct:=0;fp:=0
  for _,r:=range rows{
   a:=r.Predictions[name]
   for j,truth:=range r.Truth{
    if a[j]==truth{correct++}
    if !truth&&a[j]{fp++}
   }
   if a==r.Truth{paired++}
  }
  scores[name]=map[string]float64{
   "correct":float64(correct),"accuracy":float64(correct)/192,
   "paired_correct":float64(paired),"paired_accuracy":float64(paired)/96,
   "negative_false_positives":float64(fp),
  }
  _=c
 }
 status:="INVALID_EXACT_DEPTH_REFERENCE"
 if scores["coded_exact_depth_oracle"]["paired_correct"]==96{
  s:=scores["learned_query_conditioned_recurrence"]
  initialized:=scores["untrained_same_interface"]
  passed:=s["paired_accuracy"]>=.85&&s["accuracy"]>=.90&&
    s["negative_false_positives"]<=10&&s["paired_accuracy"]-initialized["paired_accuracy"]>=.15
  if passed{status="DEVELOPMENT_QUERY_DEPTH_TRANSFER_OBSERVED"}else{
   status="DEVELOPMENT_NO_LEARNED_QUERY_DEPTH_TRANSFER"
  }
 }
 return map[string]any{
  "schema":"wingless.f11.query_conditioned_recursion-development.v1",
  "stage":"DEVELOPMENT_ONLY","source_F10_sha":"fab5f094a29203a5ef1d63fb40225a3aa9a7e57f",
  "train_seed_bounds":[]int{51900,51916},"train_depths":[]int{3,4,5,6},
  "train_graph_states":128,"validation_seed_bounds":[]int{52300,52324},
  "validation_depths":[]int{7,8},"validation_graph_states":96,
  "validation_questions":192,"negative_questions":96,"state_bytes":576,
  "query_depth_side_input_bytes":1,"candidate_trainable_params":5,
  "pretrained_identity_projection_params":257,
  "optimizer_steps":105,"training_ms":trainMs,"training_loss":loss,
  "candidate_weights":candidate.W,"scores":scores,"rows":rows,
  "programmed_depth_loop":true,"learned_halting":false,
  "coded_identity_equality_in_candidate":false,
  "coded_OR_reachability_in_candidate":false,
  "registered_verdict":status,"independent_source_birth":false,
  "sealed_first_use":false,"accepted_ref_mutation":false,
  "external_llm_calls":0,"original_wall_ms":time.Since(begin).Milliseconds(),
 },nil
}
