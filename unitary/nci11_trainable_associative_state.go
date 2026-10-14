package unitary

import ("fmt";"math";"math/rand";"encoding/binary";"unsafe")

// W-NCI11: native 576-byte trainable associative state. 48 learned diagonal
// encoding/readout coefficients, output gain and bias, trained with development
// future-question loss via deterministic simultaneous perturbation SGD.
// 12x12 float32 state is compressed and later decoded by a learned projection.
// Outer-product updates and four matrix products remain coded mathematical
// operations, not learned semantic language parsing or arbitrary graph search.
const nci11D=12
type nci11WMemory [nci11D][nci11D]float32
type nci11WModel struct{
 Scale [4][nci11D]float64
 Gain float64
 Bias float64
 Leaky bool
}
type nci11WCase struct{
 Edges [12]nci02WEdge
 Questions [2][2]int
 Truth [2]bool
 Seed int
 Family int
 CutA bool
 CutB bool
}
type nci11WRecord struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 CutA bool `json:"cut_a"`
 CutB bool `json:"cut_b"`
 Truth [2]bool `json:"truth"`
 Learned [2]bool `json:"learned"`
 Untrained [2]bool `json:"untrained"`
 ZeroMemory [2]bool `json:"zero_memory"`
 ExactOracle [2]bool `json:"exact_oracle"`
 CodedCompression [2]bool `json:"coded_compression"`
 NativeStateBytes int `json:"native_state_bytes"`
}
func nci11WCode(id int)[nci11D]float64{
 var x [nci11D]float64
 sum:=0.
 for k:=0;k<nci11D;k++{
  f:=float64(k+1)
  x[k]=math.Sin(float64(id)*f*.119)+math.Cos(float64(id)*f*f*.017)
  sum+=x[k]*x[k]
 }
 for k:=range x{x[k]/=math.Sqrt(sum)}
 return x
}
func nci11WVector(id int,w [nci11D]float64)[nci11D]float64{
 v:=nci11WCode(id)
 sum:=0.
 for k:=range v{v[k]*=w[k];sum+=v[k]*v[k]}
 norm:=math.Sqrt(max(1e-12,sum))
 for k:=range v{v[k]/=norm}
 return v
}
func nci11WBuild(edges [12]nci02WEdge,m nci11WModel)nci11WMemory{
 var state nci11WMemory
 for _,e:=range edges{
  src:=nci11WVector(e.From,m.Scale[0])
  dst:=nci11WVector(e.To,m.Scale[1])
  for i:=0;i<nci11D;i++{for j:=0;j<nci11D;j++{
   next:=float64(state[i][j])
   if m.Leaky{next*=.92}
   state[i][j]=float32(math.Tanh(next+src[i]*dst[j]/4))
  }}
 }
 return state
}
func nci11WLogit(state nci11WMemory,m nci11WModel,a,b int)float64{
 start:=nci11WVector(a,m.Scale[2])
 end:=nci11WVector(b,m.Scale[3])
 for hop:=0;hop<4;hop++{
  var next [nci11D]float64
  for i:=0;i<nci11D;i++{for j:=0;j<nci11D;j++{
   next[j]+=start[i]*float64(state[i][j])
  }}
  start=next
 }
 dot:=0.
 for k:=range start{dot+=start[k]*end[k]}
 return m.Gain*dot+m.Bias
}
func nci11WInit(leaky bool)nci11WModel{
 m:=nci11WModel{Gain:10,Leaky:leaky}
 for c:=0;c<4;c++{for i:=0;i<nci11D;i++{m.Scale[c][i]=1}}
 return m
}
func nci11WFixture(seed,family int,ca,cb bool)(nci11WCase,error){
 base:=family%4
 chains,raw,_:=nci10WFixture(seed,base)
 if len(raw)!=12{return nci11WCase{},fmt.Errorf("bad experiment fixture")}
 r:=nci11WCase{Seed:seed,Family:family,CutA:ca,CutB:cb}
 r.Questions=[2][2]int{{chains[0][0],chains[0][4]},{chains[1][0],chains[1][4]}}
 r.Truth=[2]bool{!ca,!cb}
 copies:=0
 decoyOrdinal:=0
 for i,record:=range raw{
  e:=record.f.Edge
  if ca&&e.From==chains[0][2]&&e.To==chains[0][3]{
   e.To=300+seed%13;copies++
  }
  if cb&&e.From==chains[1][2]&&e.To==chains[1][3]{
   e.To=320+seed%13;copies++
  }
  if family>=4&&!record.main{
   // Public graph motifs absent from development: cycles, fanouts and
   // incoming distractors; preserve 12 observed relations and order.
   j:=decoyOrdinal
   decoyOrdinal++
   switch family{
   case 4:
    targets:=[4]nci02WEdge{{From:chains[0][1],To:200+seed%7},
      {From:200+seed%7,To:chains[0][1]},
      {From:chains[1][2],To:210+seed%7},
      {From:210+seed%7,To:chains[1][2]}}
    e=targets[j]
   case 5:
    targets:=[4]nci02WEdge{{From:220+seed%7,To:chains[0][1]},
      {From:chains[0][1],To:230+seed%7},
      {From:240+seed%7,To:chains[1][2]},
      {From:chains[1][2],To:250+seed%7}}
    e=targets[j]
   case 6:
    targets:=[4]nci02WEdge{{From:chains[0][2],To:260+seed%7},
      {From:chains[0][2],To:270+seed%7},
      {From:chains[1][1],To:280+seed%7},
      {From:chains[1][1],To:290+seed%7}}
    e=targets[j]
   case 7:
    targets:=[4]nci02WEdge{{From:200+seed%7,To:chains[0][1]},
      {From:210+seed%7,To:chains[0][2]},
      {From:220+seed%7,To:chains[1][1]},
      {From:230+seed%7,To:chains[1][2]}}
    e=targets[j]
   }
  }
  r.Edges[i]=e
 }
 if copies!=nci11WBoolInt(ca)+nci11WBoolInt(cb){
  return nci11WCase{},fmt.Errorf("counterfactual target mutation missing")
 }
 return r,nil
}
func nci11WBoolInt(x bool)int{if x{return 1};return 0}
func nci11WSet(start,end int, families []int)([]nci11WCase,error){
 cases:=[]nci11WCase{}
 for seed:=start;seed<end;seed++{for _,fam:=range families{
  for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
   item,err:=nci11WFixture(seed,fam,ca,cb);if err!=nil{return nil,err}
   cases=append(cases,item)
  }}
 }}
 return cases,nil
}
func nci11WLoss(m nci11WModel,items []nci11WCase)float64{
 total:=0.
 for _,item:=range items{
  memory:=nci11WBuild(item.Edges,m)
  for q:=0;q<2;q++{
   logit:=nci11WLogit(memory,m,item.Questions[q][0],item.Questions[q][1])
   p:=max(-25.,min(25.,logit))
   label:=0.;if item.Truth[q]{label=1}
   total+=math.Log1p(math.Exp(p))-label*p
  }
 }
 return total/float64(2*len(items))
}
func nci11WAccuracy(m nci11WModel,items []nci11WCase)(float64,float64){
 ok,ablated:=0,0
 for _,item:=range items{
  memory:=nci11WBuild(item.Edges,m)
  for q:=0;q<2;q++{
   predicted:=nci11WLogit(memory,m,item.Questions[q][0],item.Questions[q][1])>=0
   if predicted==item.Truth[q]{ok++}
   absent:=(m.Bias>=0)
   if absent==item.Truth[q]{ablated++}
  }
 }
 return float64(ok)/float64(2*len(items)),float64(ablated)/float64(2*len(items))
}
func nci11WLearn(train []nci11WCase,leaky bool)nci11WModel{
 m:=nci11WInit(leaky)
 random:=rand.New(rand.NewSource(118811))
 for t:=0;t<105;t++{
  c:=.08
  rate:=.11/math.Sqrt(float64(t+1))
  var perturb [4][nci11D]float64
  plus,minus:=m,m
  for k:=0;k<4;k++{for j:=0;j<nci11D;j++{
   d:=-1.;if random.Intn(2)==1{d=1}
   perturb[k][j]=d
   plus.Scale[k][j]+=c*d
   minus.Scale[k][j]-=c*d
  }}
  lp,lm:=nci11WLoss(plus,train),nci11WLoss(minus,train)
  gradient:=(lp-lm)/(2*c)
  for k:=0;k<4;k++{for j:=0;j<nci11D;j++{
   m.Scale[k][j]-=rate*gradient*perturb[k][j]
   m.Scale[k][j]=max(.10,min(3.,m.Scale[k][j]))
  }}
  // Independently train query calibration on future-question cross-entropy.
  gradGain,gradBias:=0.,0.
  for _,example:=range train{
   matrix:=nci11WBuild(example.Edges,m)
   for q:=0;q<2;q++{
    logit:=nci11WLogit(matrix,m,example.Questions[q][0],example.Questions[q][1])
    sig:=1/(1+math.Exp(-max(-25.,min(25.,logit))))
    y:=0.;if example.Truth[q]{y=1}
    gradBias+=sig-y
    if math.Abs(m.Gain)>1e-5{gradGain+=(sig-y)*(logit-m.Bias)/m.Gain}
   }
  }
  n:=float64(len(train)*2)
  m.Bias-=.3*gradBias/n
  m.Gain=max(.1,min(100.,m.Gain-.3*gradGain/n))
 }
 return m
}
func nci11WCheckpoint(memory nci11WMemory)(nci11WMemory,bool){
 bytes:=make([]byte,576)
 for i:=0;i<nci11D;i++{for j:=0;j<nci11D;j++{
  binary.LittleEndian.PutUint32(bytes[4*(i*nci11D+j):],math.Float32bits(memory[i][j]))
 }}
 var restored nci11WMemory
 for i:=0;i<nci11D;i++{for j:=0;j<nci11D;j++{
  restored[i][j]=math.Float32frombits(binary.LittleEndian.Uint32(bytes[4*(i*nci11D+j):]))
 }}
 return restored,restored==memory
}
func RunNCI11Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("bad phase")}
 train,err:=nci11WSet(33920,33944,[]int{0,1});if err!=nil{return nil,err}
 valid,err:=nci11WSet(34000,34008,[]int{2,3});if err!=nil{return nil,err}
 learned:=[]nci11WModel{nci11WLearn(train,false),nci11WLearn(train,true)}
 selected,val,ablated:=0,-1.,0.
 for i,model:=range learned{
  acc,zero:=nci11WAccuracy(model,valid)
  if acc>val{selected,val,ablated=i,acc,zero}
 }
 m:=learned[selected]
 initial:=nci11WInit(m.Leaky)
 d:=0.
 for k:=0;k<4;k++{for j:=0;j<nci11D;j++{delta:=m.Scale[k][j]-initial.Scale[k][j];d+=delta*delta}}
 result:=map[string]any{
  "schema":"wingless.nci11.learned-associative-relation-state.v1",
  "phase":phase,"parent":"256281ee6873ebeb0890624c605b5b274785087c",
  "training_seed_bounds":[]int{33920,33944},
  "validation_seed_bounds":[]int{34000,34008},
  "sealed_seed_bounds":[]int{34112,34136},
  "train_families":[]int{0,1},"validate_families":[]int{2,3},
  "sealed_novel_families":[]int{4,5,6,7},
  "development_validation_accuracy":val,"zero_memory_validation_accuracy":ablated,
  "model_modes":[]string{"additive","leaky_recurrent"},
  "selected_model_index":selected,"learned_projection_delta_l2":math.Sqrt(d),
  "native_state_bytes":unsafe.Sizeof(nci11WMemory{}),
  "learned_parameter_count":50,"learned_parameters_float64_bytes":400,
  "state_encoding_weights_trainable":true,"query_projection_weights_trainable":true,
  "public_symbol_hash_fixed":true,"four_matrix_products_coded":true,
  "supervised_future_query_training_examples":len(train)*2,
  "training_full_passes_per_model":105,
  "training_loss_evaluations_per_model":210,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "nci01_full_acceptance":false,"verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 sealed,err:=nci11WSet(34112,34136,[]int{4,5,6,7});if err!=nil{return nil,err}
 baseRetention,_,_:=nci07WTrain()
 join,cut,_:=nci02WTrain()
 untrained:=nci11WInit(m.Leaky)
 rows:=[]nci11WRecord{}
 correct,untrainedCorrect,zeroCorrect,oracleCorrect,codedCorrect:=0,0,0,0,0
 grouped:=map[string]int{}
 restoreCount:=0
 for _,item:=range sealed{
  state:=nci11WBuild(item.Edges,m)
  recovered,ok:=nci11WCheckpoint(state)
  if !ok{return nil,fmt.Errorf("checkpoint integrity mismatch")}
  restoreCount++
  uState:=nci11WBuild(item.Edges,untrained)
  stream:=[]nci06Fact{}
  for i,e:=range item.Edges{
   stream=append(stream,nci06Fact{Edge:e,Exposure:2,Inserted:i,Last:i})
  }
  codedPacked,_:=nci10R2Coded(stream,baseRetention)
  coded:=nci10WChunksToNCI09(codedPacked)
  record:=nci11WRecord{
   Seed:item.Seed,Family:item.Family,CutA:item.CutA,CutB:item.CutB,
   Truth:item.Truth,NativeStateBytes:576,
  }
  cost:=0
  for i,q:=range item.Questions{
   record.Learned[i]=nci11WLogit(recovered,m,q[0],q[1])>=0
   record.Untrained[i]=nci11WLogit(uState,untrained,q[0],q[1])>=0
   record.ZeroMemory[i]=m.Bias>=0
   record.ExactOracle[i]=nci04WOracle(item.Edges[:],q[0],q[1],4)
   record.CodedCompression[i]=nci09WReach(join,cut,coded,q[0],q[1],4,&cost)
   if record.Learned[i]==item.Truth[i]{correct++;grouped[fmt.Sprintf("%d",item.Family)]++}
   if record.Untrained[i]==item.Truth[i]{untrainedCorrect++}
   if record.ZeroMemory[i]==item.Truth[i]{zeroCorrect++}
   if record.ExactOracle[i]==item.Truth[i]{oracleCorrect++}
   if record.CodedCompression[i]==item.Truth[i]{codedCorrect++}
  }
  rows=append(rows,record)
 }
 total:=float64(len(rows)*2)
 pct:=float64(correct)/total
 familyAccuracy:=map[string]float64{}
 for _,f:=range []int{4,5,6,7}{key:=fmt.Sprintf("%d",f);familyAccuracy[key]=float64(grouped[key])/float64(len(rows)/4*2)}
 verdict:="NCI11_NO_ROBUST_LATENT_TRANSFER"
 if oracleCorrect!=len(rows)*2{verdict="NCI11_INVALID_ORACLE"
 }else if pct>=.70 && pct-float64(zeroCorrect)/total>=.10{
  safe:=true
  for _,v:=range familyAccuracy{if v<.60{safe=false}}
  if safe{verdict="NCI11_BOUNDED_LEARNED_LATENT_TRANSFER"}
 }
 result["verdict"]=verdict;result["records"]=rows
 result["sealed_questions"]=len(rows)*2;result["sealed_accuracy"]=pct
 result["untrained_latent_accuracy"]=float64(untrainedCorrect)/total
 result["zero_memory_accuracy"]=float64(zeroCorrect)/total
 result["coded_compression_accuracy"]=float64(codedCorrect)/total
 result["oracle_accuracy"]=float64(oracleCorrect)/total
 result["by_new_structure"]=familyAccuracy
 result["checkpoint_restores"]=restoreCount
 return result,nil
}
