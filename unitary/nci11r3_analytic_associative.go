package unitary

import("math";"time";"unsafe")

// NCI11-R3: analytic task-gradient training of 576-byte native association
// state, with two alternative fixed-time decay dynamics. State encoder and
// query decoder weights are task-trained; four matrix products are coded.
type nci11R3Model struct {
 Write [12]float64
 Query [12]float64
 Gain float64
 Bias float64
 Decay float64
}
type nci11R3Example struct {
 Item nci11WCase
 Matrix [12][12]float64
}
type nci11R3Candidate struct {
 Mode string `json:"mode"`
 ValidationAccuracy float64 `json:"validation_accuracy"`
 InitialAccuracy float64 `json:"same_architecture_initial_accuracy"`
 ZeroMemoryAccuracy float64 `json:"trained_zero_memory_accuracy"`
 ProgrammedControlAccuracy float64 `json:"frozen_coded_compression_accuracy"`
 TrainingMinusInitial float64 `json:"training_minus_initial"`
 TrainingMinusZero float64 `json:"training_minus_zero_memory"`
 WeightDeltaL2 float64 `json:"weight_l2_change"`
 TrainingSteps int `json:"training_steps"`
 TrainLoss float64 `json:"train_bce"`
 ResidentBytes int `json:"resident_state_bytes"`
 ModelBytes int `json:"model_parameter_bytes_float64"`
 CheckpointPass bool `json:"native_checkpoint_pass"`
 JoinEvaluations int `json:"coded_join_candidate_evaluations"`
 TrainingWallMs int64 `json:"training_wall_ms"`
}
func nci11R3Init(decay float64)nci11R3Model{
 m:=nci11R3Model{Gain:16,Decay:decay}
 for i:=0;i<12;i++{m.Write[i]=1;m.Query[i]=1}
 return m
}
func nci11R3Rename(id,seed,offset int)int{
 return offset+(id*31+seed*19)%10007
}
func nci11R3Data(lo,hi int,families []int,offset int)([]nci11R3Example,error){
 out:=[]nci11R3Example{}
 for seed:=lo;seed<hi;seed++{for _,fam:=range families{
  for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
   item,err:=nci11WFixture(seed,fam,ca,cb);if err!=nil{return nil,err}
   for i,e:=range item.Edges{
    item.Edges[i]=nci02WEdge{
     From:nci11R3Rename(e.From,seed,offset),
     To:nci11R3Rename(e.To,seed,offset)}
   }
   for i,q:=range item.Questions{
    item.Questions[i]=[2]int{
      nci11R3Rename(q[0],seed,offset),
      nci11R3Rename(q[1],seed,offset)}
   }
   // Write a single deterministic source-target rank-1 public observation
   // into a 12x12 float32 resident state; no target question in this stage.
   var m [12][12]float64
   for _,e:=range item.Edges{
    a,b:=nci11WCode(e.From),nci11WCode(e.To)
    for i:=0;i<12;i++{for j:=0;j<12;j++{
      m[i][j]=m[i][j]*.93+a[i]*b[j]/4
    }}
   }
   out=append(out,nci11R3Example{Item:item,Matrix:m})
  }}
 }}
 return out,nil
}
func nci11R3Build(item nci11WCase,decay float64)[12][12]float64{
 var matrix [12][12]float64
 for _,e:=range item.Edges{
  a,b:=nci11WCode(e.From),nci11WCode(e.To)
  for i:=0;i<12;i++{for j:=0;j<12;j++{
    matrix[i][j]=decay*matrix[i][j]+a[i]*b[j]/4
  }}
 }
 return matrix
}
func nci11R3State(e nci11R3Example,m nci11R3Model)nci11WMemory{
 data:=nci11R3Build(e.Item,m.Decay)
 var state nci11WMemory
 for i:=0;i<12;i++{for j:=0;j<12;j++{
   state[i][j]=float32(data[i][j]*m.Write[i])
 }}
 return state
}
func nci11R3Predict(e nci11R3Example,m nci11R3Model,q int,
 calculate bool)(float64,[26]float64){
 base:=nci11R3Build(e.Item,m.Decay)
 var layers [5][12]float64
 layers[0]=nci11WCode(e.Item.Questions[q][0])
 target:=nci11WCode(e.Item.Questions[q][1])
 for hop:=0;hop<4;hop++{
  for i:=0;i<12;i++{for j:=0;j<12;j++{
    layers[hop+1][j]+=layers[hop][i]*base[i][j]*m.Write[i]
  }}
 }
 dot:=0.
 for j:=0;j<12;j++{dot+=layers[4][j]*target[j]*m.Query[j]}
 logit:=m.Gain*dot+m.Bias
 var gradients [26]float64
 if !calculate{return logit,gradients}
 y:=0.;if e.Item.Truth[q]{y=1}
 d:=1/(1+math.Exp(-max(-35.,min(35.,logit))))-y
 gradients[24]=d*dot;gradients[25]=d
 var adj [12]float64
 for j:=0;j<12;j++{
  gradients[12+j]=d*m.Gain*layers[4][j]*target[j]
  adj[j]=d*m.Gain*target[j]*m.Query[j]
 }
 for hop:=3;hop>=0;hop--{
  var prev [12]float64
  for i:=0;i<12;i++{for j:=0;j<12;j++{
    gradients[i]+=adj[j]*layers[hop][i]*base[i][j]
    prev[i]+=adj[j]*base[i][j]*m.Write[i]
  }}
  adj=prev
 }
 return logit,gradients
}
func nci11R3Accuracy(model nci11R3Model,data []nci11R3Example)(float64,float64) {
 correct,zero:=0,0
 for _,ex:=range data{for q:=0;q<2;q++{
  logit,_:=nci11R3Predict(ex,model,q,false)
  if (logit>=0)==ex.Item.Truth[q]{correct++}
  if (model.Bias>=0)==ex.Item.Truth[q]{zero++}
 }}
 return float64(correct)/float64(len(data)*2),float64(zero)/float64(len(data)*2)
}
func nci11R3Train(data []nci11R3Example,decay float64)nci11R3Model {
 model:=nci11R3Init(decay)
 var first,second [26]float64
 for epoch:=0;epoch<80;epoch++{
  var gradients [26]float64
  for _,ex:=range data{for q:=0;q<2;q++{
    _,g:=nci11R3Predict(ex,model,q,true)
    for i:=0;i<26;i++{gradients[i]+=g[i]}
  }}
  factor:=1/float64(len(data)*2)
  for k:=0;k<26;k++{
    grad:=max(-10.,min(10.,gradients[k]*factor))
    first[k]=.9*first[k]+.1*grad
    second[k]=.999*second[k]+.001*grad*grad
    adjusted:=first[k]/(1-math.Pow(.9,float64(epoch+1)))/
      (math.Sqrt(second[k]/(1-math.Pow(.999,float64(epoch+1))))+1e-8)
    lr:=.04;if k==24{lr=.3};if k==25{lr=.08}
    delta:=lr*adjusted
    if k<12{model.Write[k]=max(.05,min(4,model.Write[k]-delta))}
    if k>=12&&k<24{model.Query[k-12]=max(.05,min(4,model.Query[k-12]-delta))}
    if k==24{model.Gain=max(.1,min(300,model.Gain-delta))}
    if k==25{model.Bias=max(-15,min(15,model.Bias-delta))}
  }
 }
 return model
}
func nci11R3CodedAccuracy(data []nci11R3Example)(float64,int){
 trained,_,_:=nci07WTrain()
 join,cut,_:=nci02WTrain()
 correct,calls:=0,0
 for _,example:=range data{
  stream:=[]nci06Fact{}
  for i,e:=range example.Item.Edges{
   stream=append(stream,nci06Fact{Edge:e,Exposure:2,Inserted:i,Last:i})
  }
  compressed,_:=nci10R2Coded(stream,trained)
  state:=nci10WChunksToNCI09(compressed)
  for q,pair:=range example.Item.Questions{
   seen:=nci09WReach(join,cut,state,pair[0],pair[1],4,&calls)
   if seen==example.Item.Truth[q]{correct++}
  }
 }
 return float64(correct)/float64(len(data)*2),calls
}
func nci11R3Trial(decay float64,train,valid []nci11R3Example)(nci11R3Candidate,error){
 initial:=nci11R3Init(decay)
 start:=time.Now()
 model:=nci11R3Train(train,decay)
 acc,zero:=nci11R3Accuracy(model,valid)
 initAcc,_:=nci11R3Accuracy(initial,valid)
 ctl,ops:=nci11R3CodedAccuracy(valid)
 delta:=0.
 for i:=0;i<12;i++{
  delta+=math.Pow(model.Write[i]-initial.Write[i],2)
  delta+=math.Pow(model.Query[i]-initial.Query[i],2)
 }
 delta+=math.Pow(model.Gain-initial.Gain,2)+math.Pow(model.Bias-initial.Bias,2)
 state:=nci11R3State(valid[0],model)
 recovered,ok:=nci11WCheckpoint(state)
 name:="unattenuated_associative"
 if decay<1{name="recurrent_decayed_associative"}
 loss:=0.
 for _,ex:=range train{for q:=0;q<2;q++{
  logit,_:=nci11R3Predict(ex,model,q,false)
  logit=max(-30.,min(30.,logit))
  label:=0.;if ex.Item.Truth[q]{label=1}
  loss+=math.Log1p(math.Exp(logit))-label*logit
 }}
 return nci11R3Candidate{
  Mode:name,ValidationAccuracy:acc,InitialAccuracy:initAcc,
  ZeroMemoryAccuracy:zero,ProgrammedControlAccuracy:ctl,
  TrainingMinusInitial:acc-initAcc,TrainingMinusZero:acc-zero,
  WeightDeltaL2:math.Sqrt(delta),TrainingSteps:80,
  TrainLoss:loss/float64(len(train)*2),
  ResidentBytes:int(unsafe.Sizeof(state)),ModelBytes:26*8,
  CheckpointPass:ok&&recovered==state,JoinEvaluations:ops,
  TrainingWallMs:time.Since(start).Milliseconds(),
 },nil
}
func RunNCI11R3Wingless()(map[string]any,error){
 train,err:=nci11R3Data(37344,37368,[]int{0,1,2},100);if err!=nil{return nil,err}
 valid,err:=nci11R3Data(37424,37436,[]int{4,5,6,7},1400);if err!=nil{return nil,err}
 one,err:=nci11R3Trial(1,train,valid);if err!=nil{return nil,err}
 two,err:=nci11R3Trial(.93,train,valid);if err!=nil{return nil,err}
 selected:=one
 if two.TrainingMinusInitial>one.TrainingMinusInitial{selected=two}
 qualified:=selected.ValidationAccuracy>=.70&&
   selected.TrainingMinusInitial>=.10&&
   selected.TrainingMinusZero>=.10&&
   selected.ResidentBytes==576&&selected.CheckpointPass&&
   selected.WeightDeltaL2>.01
 classification:="NCI11R3_DEVELOPMENT_NO_GENERALIZABLE_STATE_GAIN"
 if qualified{classification="NCI11R3_DEVELOPMENT_NATIVE_LEARNING_QUALIFIED"}
 return map[string]any{
  "schema":"wingless.nci11r3.analytic-gradient-associative-development.v1",
  "stage":"DEVELOPMENT_ONLY",
  "frozen_original_parent":"2261747ae1925d6aae2e200b09b8b589bb034a30",
  "train_seed_bounds":[]int{37344,37368},"validation_seed_bounds":[]int{37424,37436},
  "train_families":[]int{0,1,2},"validation_families":[]int{4,5,6,7},
  "changed_public_entity_vocabulary":true,
  "train_graph_count":len(train),"validation_graph_count":len(valid),
  "new_native_state_bytes":576,"trained_parameters_float64":26,
  "training_method":"analytic reverse-mode gradient through four-hop native matrix composition",
  "trainable_state_encoder":true,"trainable_query_decoder":true,
  "coded_four_hop_composition":true,
  "models":[]nci11R3Candidate{one,two},"selected":selected.Mode,
  "development_qualified":qualified,"verdict":classification,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "nci01_full_acceptance":false,
  "sealed_cohort_reserved":false,"sealed_first_use":false,
 },nil
}
