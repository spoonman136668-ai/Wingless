package unitary

import (
 "fmt"
 "math"
 "time"
 "unsafe"
)

// NCI11-R2 is a DEVELOPMENT-ONLY falsification of trained native recurrent
// state updates, not a new hidden-source experiment. The 12x12 float32 state
// is 576 physical bytes. Learned scalar gates modulate *formation of state*;
// all token hashes, outer products, query traversal are still programmed.
type nci11R2WModel struct {
 Gate [4]float64
 Gain float64
 Bias float64
 AdaptiveForget bool
}
type nci11R2WResult struct {
 Name string `json:"architecture"`
 InitialAccuracy float64 `json:"same_architecture_initial_accuracy"`
 TrainedAccuracy float64 `json:"validation_accuracy"`
 ZeroMemoryAccuracy float64 `json:"trained_zero_memory_accuracy"`
 GainAboveInit float64 `json:"gain_above_initial"`
 GainAboveZero float64 `json:"gain_above_zero"`
 WeightDeltaL2 float64 `json:"trained_weight_delta_l2"`
 TrainingLoss float64 `json:"training_loss"`
 Iterations int `json:"training_steps"`
 GateWeights [4]float64 `json:"gate_weights"`
 TrainingMilliseconds int64 `json:"training_wall_ms"`
}
func nci11R2WInit(adaptive bool)nci11R2WModel {
 return nci11R2WModel{Gate:[4]float64{.1,1.0,-.3,.2},
   Gain:10,Bias:0,AdaptiveForget:adaptive}
}
func nci11R2WSigmoid(z float64)float64{
 z=max(-40.,min(40.,z))
 return 1/(1+math.Exp(-z))
}
func nci11R2WBuild(item nci11WCase,m nci11R2WModel)nci11WMemory {
 var state nci11WMemory
 for t,e:=range item.Edges {
  a,b:=nci11WCode(e.From),nci11WCode(e.To)
  dot,meanAbs:=0.,0.
  for i:=0;i<12;i++{dot+=a[i]*b[i]}
  for i:=0;i<12;i++{for j:=0;j<12;j++{meanAbs+=math.Abs(float64(state[i][j]))}}
  meanAbs/=144
  z:=m.Gate[0]+m.Gate[1]*dot+m.Gate[2]*float64(t)/12+m.Gate[3]*meanAbs
  write:=nci11R2WSigmoid(z)
  keep:=.91
  if m.AdaptiveForget{keep=.72+.24*write}
  for i:=0;i<12;i++{for j:=0;j<12;j++{
   state[i][j]=float32(math.Tanh(keep*float64(state[i][j])+write*a[i]*b[j]/4))
  }}
 }
 return state
}
func nci11R2WLogit(state nci11WMemory,m nci11R2WModel,a,b int)float64{
 n:=nci11WInit(m.AdaptiveForget)
 n.Gain=m.Gain
 n.Bias=m.Bias
 return nci11WLogit(state,n,a,b)
}
func nci11R2WLoss(m nci11R2WModel,cases []nci11WCase)float64 {
 total:=0.
 for _,item:=range cases{
  state:=nci11R2WBuild(item,m)
  for q,question:=range item.Questions{
   z:=max(-25.,min(25.,nci11R2WLogit(state,m,question[0],question[1])))
   label:=0.;if item.Truth[q]{label=1}
   total+=math.Log1p(math.Exp(z))-label*z
  }
 }
 return total/float64(len(cases)*2)
}
func nci11R2WAccuracy(m nci11R2WModel,cases []nci11WCase)(float64,float64){
 correct,zero:=0,0
 for _,item:=range cases{
  state:=nci11R2WBuild(item,m)
  for q,question:=range item.Questions{
   if (nci11R2WLogit(state,m,question[0],question[1])>=0)==item.Truth[q]{correct++}
   if (m.Bias>=0)==item.Truth[q]{zero++}
  }
 }
 return float64(correct)/float64(len(cases)*2),
    float64(zero)/float64(len(cases)*2)
}
func nci11R2WGet(m nci11R2WModel,k int)float64 {
 switch k{
 case 0,1,2,3:return m.Gate[k]
 case 4:return m.Gain
 default:return m.Bias
 }
}
func nci11R2WSet(m *nci11R2WModel,k int,v float64){
 switch k{
 case 0,1,2,3:m.Gate[k]=v
 case 4:m.Gain=v
 default:m.Bias=v
 }
}
func nci11R2WTrain(cases,valid []nci11WCase,adaptive bool)nci11R2WResult{
 m:=nci11R2WInit(adaptive)
 initial:=m
 initialAccuracy,_:=nci11R2WAccuracy(initial,valid)
 start:=time.Now()
 const steps=38
 const batchSize=24
 for step:=0;step<steps;step++{
  subset:=make([]nci11WCase,0,batchSize)
  for j:=0;j<batchSize;j++{
   ix:=(step*19+j*7)%len(cases)
   subset=append(subset,cases[ix])
  }
  for k:=0;k<6;k++{
   base:=nci11R2WGet(m,k)
   epsilon:=.03
   if k==4{epsilon=.06}
   plus,minus:=m,m
   nci11R2WSet(&plus,k,base+epsilon)
   nci11R2WSet(&minus,k,base-epsilon)
   gradient:=(nci11R2WLoss(plus,subset)-nci11R2WLoss(minus,subset))/(2*epsilon)
   gradient=max(-5.,min(5.,gradient))
   rate:=.12/math.Sqrt(1+float64(step)/4)
   if k==4{rate=.4/math.Sqrt(1+float64(step)/4)}
   updated:=base-rate*gradient
   if k<4{updated=max(-6.,min(6.,updated))}
   if k==4{updated=max(.1,min(100.,updated))}
   if k==5{updated=max(-10.,min(10.,updated))}
   nci11R2WSet(&m,k,updated)
  }
 }
 trainedAcc,zeroAcc:=nci11R2WAccuracy(m,valid)
 delta:=0.
 for k:=0;k<6;k++{
  v:=nci11R2WGet(m,k)-nci11R2WGet(initial,k)
  delta+=v*v
 }
 architecture:="fixed_decay_write_gate"
 if adaptive{architecture="predictive_write_and_forget"}
 return nci11R2WResult{
  Name:architecture,InitialAccuracy:initialAccuracy,
  TrainedAccuracy:trainedAcc,ZeroMemoryAccuracy:zeroAcc,
  GainAboveInit:trainedAcc-initialAccuracy,
  GainAboveZero:trainedAcc-zeroAcc,
  WeightDeltaL2:math.Sqrt(delta),
  TrainingLoss:nci11R2WLoss(m,cases),Iterations:steps,
  GateWeights:m.Gate,TrainingMilliseconds:time.Since(start).Milliseconds(),
 }
}
func RunNCI11R2Wingless()(map[string]any,error){
 train,err:=nci11WSet(35328,35352,[]int{0,1});if err!=nil{return nil,err}
 valid,err:=nci11WSet(35392,35404,[]int{2,3});if err!=nil{return nil,err}
 cases:=[]nci11R2WResult{
  nci11R2WTrain(train,valid,false),
  nci11R2WTrain(train,valid,true),
 }
 best:=cases[0]
 for _,row:=range cases[1:]{
  score:=row.GainAboveInit+row.GainAboveZero
  chosen:=best.GainAboveInit+best.GainAboveZero
  if score>chosen||(score==chosen&&row.TrainedAccuracy>best.TrainedAccuracy){
   best=row
  }
 }
 qualified:=best.TrainedAccuracy>=.70&&best.GainAboveInit>=.10&&
  best.GainAboveZero>=.10&&best.WeightDeltaL2>.01
 if unsafe.Sizeof(nci11WMemory{})!=576{
  return nil,fmt.Errorf("native resident state no longer 576 bytes")
 }
 verdict:="NCI11R2_DEV_NO_RELIABLE_LEARNED_MEMORY_GAIN"
 if qualified{verdict="NCI11R2_DEV_NATIVE_LEARNING_QUALIFIED"}
 return map[string]any{
  "schema":"wingless.nci11r2.dev-learned-recurrent-state.v1",
  "mode":"DEVELOPMENT_ONLY",
  "source_parent":"b820ea1ff76f17e350a914177ea823f6984a8e09",
  "training_seeds":[]int{35328,35352},"validation_seeds":[]int{35392,35404},
  "original_first_use_executed":false,"sealed_cohort_reserved":false,
  "training_graphs":len(train),"validation_graphs":len(valid),
  "models":cases,"selected_architecture":best.Name,
  "native_resident_state_bytes":unsafe.Sizeof(nci11WMemory{}),
  "trained_gate_parameters":6,"model_parameter_bytes_float64":48,
  "public_feature_hash_fixed":true,
  "outer_product_and_four_step_composition_programmed":true,
  "training_method":"coordinatewise finite-difference supervised future-query BCE",
  "development_qualified":qualified,
  "qualification_min_accuracy":.70,
  "qualification_min_gain_vs_initialized":.10,
  "qualification_min_gain_vs_zero_memory":.10,
  "verdict":verdict,"external_llm_calls":0,
  "accepted_ref_mutation":false,"nci01_full_acceptance":false,
 },nil
}
