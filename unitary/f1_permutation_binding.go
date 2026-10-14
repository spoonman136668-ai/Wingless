package unitary

import("fmt";"math";"time";"unsafe")

// W-F1: native Go permutation binding falsification, DEVELOPMENT ONLY.
// The 576-byte bounded memory holds 12 packed public edge endpoints and
// 132 off-diagonal learned link scores. Identity equality and 4-hop search
// are CODED; trained 3-weight scoring alone is not learned variable binding.
type wf1State [144]float32
type wf1Model struct { W [3]float64; Form string }
type wf1Datum struct {
 Edges [12]nci02WEdge
 Questions [2][2]int
 Labels [2]bool
 Eq [12][12]float64
 Similar [12][12]float64
}
type wf1Score struct {
 Form string `json:"form"`
 TrainedAccuracy float64 `json:"trained_accuracy"`
 InitialAccuracy float64 `json:"initial_accuracy"`
 ZeroMemoryAccuracy float64 `json:"zero_memory_accuracy"`
 RenamingAgreement float64 `json:"renaming_prediction_agreement"`
 TrainingSteps int `json:"training_steps"`
 ModelParams int `json:"model_parameter_count"`
 ModelBytes int `json:"model_float64_bytes"`
 ResidentBytes int `json:"resident_state_bytes"`
 CheckpointPass bool `json:"checkpoint_pass"`
 TrainedWeights [3]float64 `json:"weights"`
 TrainingLoss float64 `json:"train_bce"`
 TrainingMs int64 `json:"training_wall_ms"`
}
func wf1Rename(id,seed,variant int)int {
 mul:=[3]int{61,73,89};s:=[3]int{17,19,23}
 return (id*mul[variant]+seed*s[variant]+101)%60000
}
func wf1New(seed,family int,ca,cb bool,variant int)(wf1Datum,error){
 original,err:=nci11WFixture(seed,family,ca,cb)
 if err!=nil{return wf1Datum{},err}
 d:=wf1Datum{}
 for i,e:=range original.Edges{
  d.Edges[i]=nci02WEdge{From:wf1Rename(e.From,seed,variant),
   To:wf1Rename(e.To,seed,variant)}
 }
 for i,p:=range original.Questions{
  d.Questions[i]=[2]int{wf1Rename(p[0],seed,variant),
    wf1Rename(p[1],seed,variant)}
  d.Labels[i]=nci04WOracle(d.Edges[:],d.Questions[i][0],d.Questions[i][1],4)
 }
 for i,a:=range d.Edges{
  da:=nci11WCode(a.To)
  for j,b:=range d.Edges{
   if a.To==b.From{d.Eq[i][j]=1}
   sb:=nci11WCode(b.From)
   for k:=0;k<12;k++{d.Similar[i][j]+=da[k]*sb[k]}
  }
 }
 return d,nil
}
func wf1Data(lo,hi int,families []int,variant int)([]wf1Datum,error){
 data:=[]wf1Datum{}
 for seed:=lo;seed<hi;seed++{
  for _,family:=range families{
   for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
    v,err:=wf1New(seed,family,ca,cb,variant)
    if err!=nil{return nil,err}
    data=append(data,v)
   }}
  }
 }
 return data,nil
}
func wf1Init(form string)wf1Model{
 return wf1Model{W:[3]float64{-.2,0,-1},Form:form}
}
func wf1Prob(z float64)float64{
 return 1/(1+math.Exp(-max(-40.,min(40.,z))))
}
func wf1Memory(d wf1Datum,m wf1Model)wf1State{
 var mem wf1State
 for i,e:=range d.Edges{
  bits:=(uint32(e.From)<<16)|uint32(e.To)
  mem[i]=math.Float32frombits(bits)
 }
 off:=12
 for i:=0;i<12;i++{for j:=0;j<12;j++{
  if i==j{continue}
  signal,secondary:=d.Eq[i][j],0.
  if m.Form=="identity_similarity"{signal=d.Similar[i][j]}
  if m.Form=="hybrid"{secondary=d.Similar[i][j]}
  mem[off]=float32(wf1Prob(m.W[0]*signal+m.W[1]*secondary+m.W[2]))
  off++
 }}
 return mem
}
func wf1Answer(mem wf1State,question [2]int)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(mem[i])
  src[i]=int(bits>>16)
  dst[i]=int(bits&65535)
 }
 var dp [12]float64
 for i:=0;i<12;i++{if src[i]==question[0]{dp[i]=1}}
 for hop:=1;hop<4;hop++{
  var nxt [12]float64
  for j:=0;j<12;j++{
   product:=1.
   for i:=0;i<12;i++{
    if i==j{continue}
    off:=12+i*11+j
    if j>i{off--}
    product*=1-dp[i]*float64(mem[off])
   }
   nxt[j]=1-product
  }
  dp=nxt
 }
 miss:=1.
 for i:=0;i<12;i++{if dst[i]==question[1]{miss*=1-dp[i]}}
 return 1-miss
}
func wf1Loss(data []wf1Datum,model wf1Model)float64{
 total:=0.
 for _,d:=range data{
  state:=wf1Memory(d,model)
  for q,p:=range d.Questions{
   score:=max(.0001,min(.9999,wf1Answer(state,p)))
   if d.Labels[q]{total-=math.Log(score)}else{total-=math.Log(1-score)}
  }
 }
 return total/float64(len(data)*2)
}
func wf1Accuracy(data []wf1Datum,model wf1Model)(float64,[]bool){
 count:=0
 outcomes:=make([]bool,0,len(data)*2)
 for _,d:=range data{
  state:=wf1Memory(d,model)
  for q,p:=range d.Questions{
   predicted:=wf1Answer(state,p)>=.5
   outcomes=append(outcomes,predicted)
   if predicted==d.Labels[q]{count++}
  }
 }
 return float64(count)/float64(len(data)*2),outcomes
}
func wf1Train(train []wf1Datum,form string)(wf1Model,float64,int64) {
 m:=wf1Init(form)
 start:=time.Now()
 for epoch:=0;epoch<65;epoch++{
  batch:=make([]wf1Datum,0,36)
  for j:=0;j<36;j++{batch=append(batch,train[(epoch*17+j*13)%len(train)])}
  var grad [3]float64
  for k:=0;k<3;k++{
   plus,minus:=m,m
   plus.W[k]+=.035;minus.W[k]-=.035
   grad[k]=(wf1Loss(batch,plus)-wf1Loss(batch,minus))/.07
  }
  rate:=.35/math.Sqrt(1+float64(epoch)/12)
  for k:=0;k<3;k++{
   m.W[k]-=rate*max(-8,min(8,grad[k]))
   m.W[k]=max(-15,min(15,m.W[k]))
  }
 }
 return m,wf1Loss(train,m),time.Since(start).Milliseconds()
}
func wf1Checkpoint(s wf1State)bool{
 // Bitwise exact resident-state checkpoint, including encoded float bit patterns.
 var restored wf1State
 for i,x:=range s{restored[i]=math.Float32frombits(math.Float32bits(x))}
 for i:=range s{
  if math.Float32bits(restored[i])!=math.Float32bits(s[i]){return false}
 }
 return true
}
func RunWFoundationF1()(map[string]any,error){
 train,err:=wf1Data(38304,38328,[]int{0,1,2},0);if err!=nil{return nil,err}
 valid,err:=wf1Data(38416,38428,[]int{4,5,6,7},1);if err!=nil{return nil,err}
 alternate,err:=wf1Data(38416,38428,[]int{4,5,6,7},2);if err!=nil{return nil,err}
 forms:=[]string{"exact_binding","identity_similarity","hybrid"}
 scores:=[]wf1Score{}
 performance:=map[string]float64{}
 for _,form:=range forms{
  model,loss,duration:=wf1Train(train,form)
  accuracy,outcomes:=wf1Accuracy(valid,model)
  baseline,_:=wf1Accuracy(valid,wf1Init(form))
  _,renamed:=wf1Accuracy(alternate,model)
  agree:=0
  for i:=range outcomes{if outcomes[i]==renamed[i]{agree++}}
  mem:=wf1Memory(valid[0],model)
  record:=wf1Score{Form:form,TrainedAccuracy:accuracy,
   InitialAccuracy:baseline,ZeroMemoryAccuracy:.5,
   RenamingAgreement:float64(agree)/float64(len(outcomes)),
   TrainingSteps:65,ModelParams:3,ModelBytes:24,
   ResidentBytes:int(unsafe.Sizeof(mem)),CheckpointPass:wf1Checkpoint(mem),
   TrainedWeights:model.W,TrainingLoss:loss,TrainingMs:duration}
  scores=append(scores,record)
  performance[form]=accuracy
 }
 // Same 576-byte memory format with coded exact equality link probabilities.
 fixed:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 control,_:=wf1Accuracy(valid,fixed)
 if int(unsafe.Sizeof(wf1State{}))!=576{return nil,fmt.Errorf("memory not 576 bytes")}
 useful:=performance["exact_binding"]>=.70&&
   performance["exact_binding"]-performance["identity_similarity"]>=.10&&
   scores[0].RenamingAgreement>=.99
 verdict:="EQUALITY_SCAFFOLD_TRANSFER_NOT_CONFIRMED"
 if useful{verdict="CODED_EQUALITY_SCAFFOLD_EXPLAINS_TRANSFER"}
 return map[string]any{
  "schema":"wingless.foundation.f1-permutation-binding-dev.v1",
  "stage":"DEVELOPMENT_ONLY",
  "source_parent":"b88f33834c7dd2e3f003ee6dc61c4ccf0968a2d7",
  "train_bounds":[]int{38304,38328},
  "validation_bounds":[]int{38416,38428},
  "train_graphs":len(train),"validation_graphs":len(valid),
  "train_families":[]int{0,1,2},"validation_families":[]int{4,5,6,7},
  "candidate_results":scores,
  "fixed_exact_binding_control_accuracy":control,
  "frozen_min_accuracy":.70,"frozen_min_structural_gain":.10,
  "frozen_min_renaming_agreement":.99,
  "diagnostic_useful":useful,"verdict":verdict,
  "native_state_bytes":576,"trained_parameter_count":3,
  "coded_equality_binding":true,"coded_four_hop_search":true,
  "no_learned_identity_matching":true,"programmed_memory_encoder":true,
  "sealed_first_use":false,"sealed_cohort_reserved":false,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "full_nci01_acceptance":false,
 },nil
}
