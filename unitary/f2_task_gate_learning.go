package unitary

import ("fmt";"math";"time";"unsafe")

// W-F2: learns only scalar calibration over PROGRAMMED equality features.
// Query/graph labels are used for TRAINING on disjoint development data, never
// as input to memory at validation time; no learned entity binding is claimed.
type wf2Metric struct {
 Name string `json:"name"`
 Accuracy float64 `json:"accuracy"`
 FalsePositiveRate float64 `json:"false_positive_rate"`
 TruePositiveRate float64 `json:"true_positive_rate"`
 Weights [3]float64 `json:"weights"`
 TrainMs int64 `json:"training_wall_ms"`
 TrainingMethod string `json:"training_method"`
}
type wf2Answer struct {
 Truth [2]bool `json:"truth"`
 Models map[string][2]bool `json:"models"`
}

func wf2WeightedLoss(data []wf1Datum,m wf1Model,negativeWeight float64)float64{
 cost,total:=0.,0.
 for _,d:=range data{
  state:=wf1Memory(d,m)
  for q,pair:=range d.Questions{
   prob:=max(.00001,min(.99999,wf1Answer(state,pair)))
   if d.Labels[q]{cost-=math.Log(prob);total++}else{
    cost-=negativeWeight*math.Log(1-prob);total+=negativeWeight
   }
  }
 }
 return cost/total
}
func wf2Train(train []wf1Datum,method string)(wf1Model,int64){
 m:=wf1Init("exact_binding")
 negativeWeight:=1.
 if method=="adam_negative_cost"{negativeWeight=2.}
 if method=="adam_structural_warmstart"{m.W=[3]float64{3,0,-2}}
 var moments,variance [3]float64
 started:=time.Now()
 for epoch:=1;epoch<=135;epoch++{
  batch:=make([]wf1Datum,0,40)
  for i:=0;i<40;i++{batch=append(batch,train[(epoch*29+i*17)%len(train)])}
  var grads [3]float64
  for k:=0;k<3;k++{
   plus,minus:=m,m
   plus.W[k]+=.045;minus.W[k]-=.045
   grads[k]=(wf2WeightedLoss(batch,plus,negativeWeight)-wf2WeightedLoss(batch,minus,negativeWeight))/.09
  }
  for k:=0;k<3;k++{
   g:=max(-15.,min(15.,grads[k]))
   moments[k]=.9*moments[k]+.1*g
   variance[k]=.999*variance[k]+.001*g*g
   correctedM:=moments[k]/(1-math.Pow(.9,float64(epoch)))
   correctedV:=variance[k]/(1-math.Pow(.999,float64(epoch)))
   m.W[k]-=.65*correctedM/(math.Sqrt(correctedV)+1e-8)
   m.W[k]=max(-20,min(20,m.W[k]))
  }
 }
 return m,time.Since(started).Milliseconds()
}
func RunWFoundationF2()(map[string]any,error){
 started:=time.Now()
 train,err:=wf1Data(40500,40524,[]int{0,1,2},0);if err!=nil{return nil,err}
 valid,err:=wf1Data(40700,40712,[]int{4,5,6,7},1);if err!=nil{return nil,err}
 if len(train)!=288||len(valid)!=192{return nil,fmt.Errorf("cohort cardinality changed")}
 models:=map[string]wf1Model{}
 durations:=map[string]int64{}
 original,_,trainMs:=wf1Train(train,"exact_binding")
 models["original_f1_optimizer"]=original;durations["original_f1_optimizer"]=trainMs
 for _,method:=range []string{"adam_task_bce","adam_negative_cost","adam_structural_warmstart"}{
  model,ms:=wf2Train(train,method);models[method]=model;durations[method]=ms
 }
 models["frozen_f1b_parent"]=wf1Model{Form:"exact_binding",W:wf1bFrozen["exact_binding"]}
 models["fixed_coded_sharp_control"]=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 answers:=make([]wf2Answer,0,len(valid))
 hits:=map[string]int{};posHits:=map[string]int{};falsePos:=map[string]int{}
 positive,negative:=0,0
 for _,d:=range valid{
  row:=wf2Answer{Truth:d.Labels,Models:map[string][2]bool{}}
  for _,truth:=range d.Labels{if truth{positive++}else{negative++}}
  for name,m:=range models{
   state:=wf1Memory(d,m)
   var pred [2]bool
   for q,pair:=range d.Questions{
    pred[q]=wf1Answer(state,pair)>=.5
    if pred[q]==d.Labels[q]{hits[name]++}
    if d.Labels[q]&&pred[q]{posHits[name]++}
    if !d.Labels[q]&&pred[q]{falsePos[name]++}
   }
   row.Models[name]=pred
  }
  answers=append(answers,row)
 }
 if positive!=192||negative!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("balance or state budget drift")}
 scores:=map[string]wf2Metric{}
 for name,m:=range models{
  scores[name]=wf2Metric{Name:name,Accuracy:float64(hits[name])/384,
   FalsePositiveRate:float64(falsePos[name])/192,
   TruePositiveRate:float64(posHits[name])/192,
   Weights:m.W,TrainMs:durations[name],TrainingMethod:name}
 }
 base:=scores["original_f1_optimizer"].Accuracy
 qualified:=false
 for _,method:=range []string{"adam_task_bce","adam_negative_cost","adam_structural_warmstart"}{
  s:=scores[method]
  if s.Accuracy>=.95&&s.FalsePositiveRate<=.05&&s.Accuracy-base>=.10{qualified=true}
 }
 verdict:="DEVELOPMENT_NO_TASK_LEARNED_SHARPNESS_GAIN"
 if qualified{verdict="DEVELOPMENT_TASK_TRAINED_SHARPNESS_OBSERVED"}
 return map[string]any{
  "schema":"wingless.foundation.f2.task-gate-learning-dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_parent":"4a8da70dc242940b42b69dd9d71b0e4f9fc7212c",
  "train_seed_bounds":[]int{40500,40524},"validation_seed_bounds":[]int{40700,40712},
  "train_graphs":len(train),"validation_graphs":len(valid),"validation_question_count":384,
  "validation_negative_count":negative,"validation_positive_count":positive,
  "candidate_scores":scores,"original_validation_rows":answers,
  "registered_verdict":verdict,"native_state_bytes":576,"trainable_parameter_bytes":24,
  "coded_entity_equality":true,"coded_four_hop_traversal":true,
  "learned_entity_binding":false,"no_sealed_first_use":true,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "development_wall_ms":time.Since(started).Milliseconds(),
 },nil
}
