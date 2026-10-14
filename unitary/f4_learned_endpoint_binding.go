package unitary

import("fmt";"math";"time";"unsafe")

// W-F4: task-trained bit-feature matchers for BOTH source/terminal query anchors
// and edge transition binding. No direct entity equality is used by candidate.
// Four-hop reachability aggregation and public 16-bit features remain CODED.
// This is a NEW DEVELOPMENT experiment, not the pending F3 first-use cohort.
type wf4Model struct{ Weight [16]float64; Bias float64 }
type wf4Case struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 CutA bool `json:"cut_a"`
 CutB bool `json:"cut_b"`
 Truth [2]bool `json:"truth"`
 Models map[string][2]bool `json:"models"`
}
func wf4Initial()wf4Model{
 m:=wf4Model{}
 for i:=range m.Weight{m.Weight[i]=-1}
 return m
}
func wf4Softplus(x float64)float64{return math.Log1p(math.Exp(max(-30,min(30,x))))}
func wf4Match(m wf4Model,a,b int)float64{
 mismatch:=(a^b)&65535
 z:=m.Bias
 for bit:=0;bit<16;bit++{if mismatch&(1<<bit)!=0{z-=wf4Softplus(m.Weight[bit])}}
 return wf1Prob(z)
}
func wf4Memory(d wf1Datum,m wf4Model)wf1State{
 var state wf1State
 for i,e:=range d.Edges{
  state[i]=math.Float32frombits((uint32(e.From)<<16)|uint32(e.To))
 }
 off:=12
 for i,a:=range d.Edges{for j,b:=range d.Edges{
  if i==j{continue}
  state[off]=float32(wf4Match(m,a.To,b.From));off++
 }}
 return state
}
func wf4Answer(mem wf1State,q [2]int,m wf4Model)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(mem[i])
  src[i]=int(bits>>16);dst[i]=int(bits&65535)
 }
 var active [12]float64
 for i:=0;i<12;i++{active[i]=wf4Match(m,q[0],src[i])}
 for hop:=1;hop<4;hop++{
  var next [12]float64
  for j:=0;j<12;j++{
   notReached:=1.
   for i:=0;i<12;i++{
    if i==j{continue}
    off:=12+i*11+j;if j>i{off--}
    notReached*=1-active[i]*float64(mem[off])
   }
   next[j]=1-notReached
  }
  active=next
 }
 miss:=1.
 for i:=0;i<12;i++{miss*=1-active[i]*wf4Match(m,dst[i],q[1])}
 return 1-miss
}
func wf4Loss(batch []wf1Datum,m wf4Model)float64{
 total:=0.
 for _,d:=range batch{
  mem:=wf4Memory(d,m)
  for j,q:=range d.Questions{
   pred:=max(.00001,min(.99999,wf4Answer(mem,q,m)))
   if d.Labels[j]{total-=math.Log(pred)}else{total-=math.Log(1-pred)}
  }
 }
 return total/float64(2*len(batch))
}
func wf4Step(m wf4Model,dim int,diff float64)wf4Model{
 if dim==16{m.Bias+=diff}else{m.Weight[dim]+=diff}
 return m
}
func wf4Train(train []wf1Datum)(wf4Model,int64,float64){
 model:=wf4Initial()
 var moments,squares [17]float64
 start:=time.Now()
 for epoch:=1;epoch<=90;epoch++{
  batch:=make([]wf1Datum,0,24)
  for i:=0;i<24;i++{batch=append(batch,train[(epoch*23+i*19)%len(train)])}
  var grad [17]float64
  for k:=0;k<17;k++{
   plus,minus:=wf4Step(model,k,.045),wf4Step(model,k,-.045)
   grad[k]=(wf4Loss(batch,plus)-wf4Loss(batch,minus))/.09
  }
  for k:=0;k<17;k++{
   g:=max(-15,min(15,grad[k]))
   moments[k]=.9*moments[k]+.1*g
   squares[k]=.999*squares[k]+.001*g*g
   mm:=moments[k]/(1-math.Pow(.9,float64(epoch)))
   vv:=squares[k]/(1-math.Pow(.999,float64(epoch)))
   rate:=.35
   delta:=rate*mm/(math.Sqrt(vv)+1e-8)
   if k==16{model.Bias=max(-15,min(15,model.Bias-delta))}
   if k<16{model.Weight[k]=max(-15,min(15,model.Weight[k]-delta))}
  }
 }
 return model,time.Since(start).Milliseconds(),wf4Loss(train,model)
}
func RunWFoundationF4()(map[string]any,error){
 totalStart:=time.Now()
 train,err:=wf1Data(42900,42924,[]int{0,1,2},0)
 if err!=nil{return nil,err}
 if len(train)!=288{return nil,fmt.Errorf("training cardinality drift %d",len(train))}
 model,trainMs,loss:=wf4Train(train)
 initial:=wf4Initial()
 fixed:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 catalog:=wf1bCatalog()
 names:=[]string{"learned_anchor_and_transition","initialized_same_arch","fixed_coded_equality_oracle_control"}
 type count struct{correct,fp int}
 sums:=map[string]count{}
 rows:=make([]wf4Case,0,192)
 negative:=0
 for seed:=43100;seed<43112;seed++{for fam:=4;fam<=7;fam++{
  for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
   _,d,_,err:=wf1bStress(seed,fam,ca,cb,catalog)
   if err!=nil{return nil,err}
   row:=wf4Case{Seed:seed,Family:fam,CutA:ca,CutB:cb,Truth:d.Labels,Models:map[string][2]bool{}}
   state:=wf4Memory(d,model)
   state0:=wf4Memory(d,initial)
   fixedState:=wf1Memory(d,fixed)
   for _,name:=range names{
    var answers [2]bool
    for j,q:=range d.Questions{
     switch name{
      case "learned_anchor_and_transition":answers[j]=wf4Answer(state,q,model)>=.5
      case "initialized_same_arch":answers[j]=wf4Answer(state0,q,initial)>=.5
      case "fixed_coded_equality_oracle_control":answers[j]=wf1Answer(fixedState,q)>=.5
     }
     c:=sums[name]
     if answers[j]==d.Labels[j]{c.correct++}
     if !d.Labels[j]&&answers[j]{c.fp++}
     sums[name]=c
    }
    row.Models[name]=answers
   }
   for _,truth:=range d.Labels{if !truth{negative++}}
   rows=append(rows,row)
  }}
 }}
 if len(rows)!=192||negative!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F4 cohort/state-budget drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  c:=sums[name]
  scores[name]=map[string]float64{
   "correct":float64(c.correct),"false_positives":float64(c.fp),
   "accuracy":float64(c.correct)/384,"negative_false_positive_rate":float64(c.fp)/192}
 }
 trained:=scores["learned_anchor_and_transition"]
 baseline:=scores["initialized_same_arch"]
 coded:=scores["fixed_coded_equality_oracle_control"]
 pass:=trained["accuracy"]>=.85&&trained["negative_false_positive_rate"]<=.15&&
  trained["accuracy"]-baseline["accuracy"]>=.15&&coded["accuracy"]-trained["accuracy"]<=.15
 verdict:="DEVELOPMENT_NO_LEARNED_QUERY_BINDING_GAIN"
 if pass{verdict="DEVELOPMENT_LEARNED_ENDPOINT_BINDING_OBSERVED"}
 return map[string]any{
  "schema":"wingless.foundation.f4.learned-query-binding-dev.v1",
  "stage":"DEVELOPMENT_ONLY",
  "source_f3_sha":"d3f6fb15dd8610833c1712d030acbba5bfe942b8",
  "train_seed_bounds":[]int{42900,42924},"validation_seed_bounds":[]int{43100,43112},
  "train_graph_states":len(train),"validation_graph_states":len(rows),
  "validation_questions":384,"negative_questions":negative,
  "native_state_bytes":576,"trainable_parameter_count":17,
  "train_steps":90,"train_loss":loss,"training_wall_ms":trainMs,
  "development_wall_ms":time.Since(totalStart).Milliseconds(),
  "selected_model":map[string]any{"weights":model.Weight,"bias":model.Bias},
  "candidate_scores":scores,"original_validation_rows":rows,
  "registered_verdict":verdict,
  "coded_entity_equality_in_candidate":false,
  "coded_query_anchor_equality_in_candidate":false,
  "coded_four_hop_or_traversal_in_candidate":true,
  "engineered_public_id_bit_features":true,
  "no_future_question_during_memory_creation":true,
  "sealed_first_use":false,"accepted_ref_mutation":false,
  "external_llm_calls":0,
 },nil
}
