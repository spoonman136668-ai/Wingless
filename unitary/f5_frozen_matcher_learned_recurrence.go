package unitary

import("fmt";"math";"time";"unsafe")

// W-F5: train composition ONLY, freezing the F4 17-parameter matcher.
// Candidate inference uses neither direct equality nor hard-coded path-OR.
// This is new development data; F3 sealed cohort is NOT accessed.
type wf5Composer struct{ W [5]float64 }
type wf5Row struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 CutA bool `json:"cut_a"`
 CutB bool `json:"cut_b"`
 Truth [2]bool `json:"truth"`
 Models map[string][2]bool `json:"models"`
}
func wf5Matcher()wf4Model{
 return wf4Model{
  Bias:4.5293453573490075,
  Weight:[16]float64{
   -1.5947942456474447,.6900551970356293,3.428782411828214,
   1.85858649386416,1.2460604579000123,1.3580545663581163,
   .23929760726945476,.9141765894027869,-.4408711085751792,
   1.993532764022319,4.634090762401231,4.6060390424996385,
   2.431704800149348,3.5912512274531316,4.287418167947507,-1.,
  }}
}
func wf5Start()wf5Composer{return wf5Composer{W:[5]float64{1.5,-1.5,-1,3,-3}}}
func wf5Readout(state wf1State,q [2]int,cmp wf4Model,m wf5Composer)float64{
 var src,dst [12]int
 for i:=0;i<12;i++{
  v:=math.Float32bits(state[i]);src[i]=int(v>>16);dst[i]=int(v&65535)
 }
 var active [12]float64
 for i:=0;i<12;i++{active[i]=wf4Match(cmp,q[0],src[i])}
 for round:=0;round<3;round++{
  var next [12]float64
  for j:=0;j<12;j++{
   sum:=0.
   for i:=0;i<12;i++{
    if i==j{continue}
    off:=12+i*11+j;if j>i{off--}
    sum+=active[i]*float64(state[off])
   }
   next[j]=wf1Prob(m.W[0]*sum+m.W[1]*active[j]+m.W[2])
  }
  active=next
 }
 combined:=0.
 for i:=0;i<12;i++{combined+=active[i]*wf4Match(cmp,dst[i],q[1])}
 return wf1Prob(m.W[3]*combined+m.W[4])
}
func wf5Loss(data []wf1Datum,cmp wf4Model,m wf5Composer)float64{
 loss:=0.
 for _,d:=range data{
  state:=wf4Memory(d,cmp)
  for j,q:=range d.Questions{
   p:=max(.00001,min(.99999,wf5Readout(state,q,cmp,m)))
   if d.Labels[j]{loss-=math.Log(p)}else{loss-=math.Log(1-p)}
  }
 }
 return loss/float64(2*len(data))
}
func wf5Train(data []wf1Datum,cmp wf4Model)(wf5Composer,int64,float64){
 m:=wf5Start()
 var moment,sq [5]float64
 start:=time.Now()
 for epoch:=1;epoch<=140;epoch++{
  batch:=make([]wf1Datum,0,30)
  for j:=0;j<30;j++{batch=append(batch,data[(epoch*29+j*23)%len(data)])}
  var grad [5]float64
  for k:=0;k<5;k++{
   plus,minus:=m,m
   plus.W[k]+=.035;minus.W[k]-=.035
   grad[k]=(wf5Loss(batch,cmp,plus)-wf5Loss(batch,cmp,minus))/.07
  }
  for k:=0;k<5;k++{
   g:=max(-15,min(15,grad[k]))
   moment[k]=.9*moment[k]+.1*g
   sq[k]=.999*sq[k]+.001*g*g
   mm:=moment[k]/(1-math.Pow(.9,float64(epoch)))
   vv:=sq[k]/(1-math.Pow(.999,float64(epoch)))
   m.W[k]-=.12*mm/(math.Sqrt(vv)+1e-8)
   m.W[k]=max(-20,min(20,m.W[k]))
  }
 }
 return m,time.Since(start).Milliseconds(),wf5Loss(data,cmp,m)
}
func RunWFoundationF5()(map[string]any,error){
 started:=time.Now()
 cmp:=wf5Matcher()
 train,err:=wf1Data(44300,44324,[]int{0,1,2},0)
 if err!=nil{return nil,err}
 if len(train)!=288{return nil,fmt.Errorf("training states drift")}
 trained,trainingMs,loss:=wf5Train(train,cmp)
 frozenFourhop:=wf1Model{Form:"exact_binding",W:[3]float64{13,0,-6.5}}
 init:=wf5Start()
 totals:=map[string][2]int{}
 rows:=[]wf5Row{}
 neg:=0
 catalog:=wf1bCatalog()
 for seed:=44600;seed<44612;seed++{for fam:=4;fam<=7;fam++{
  for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
   _,d,_,err:=wf1bStress(seed,fam,ca,cb,catalog)
   if err!=nil{return nil,err}
   learnedState:=wf4Memory(d,cmp)
   fixedState:=wf1Memory(d,frozenFourhop)
   row:=wf5Row{Seed:seed,Family:fam,CutA:ca,CutB:cb,Truth:d.Labels,Models:map[string][2]bool{}}
   for _,name:=range []string{"learned_recurrent_composition","untrained_same_recurrence","frozen_F4_coded_fourhop","coded_exact_upper_control"}{
    var answer [2]bool
    for j,q:=range d.Questions{
     switch name{
     case "learned_recurrent_composition": answer[j]=wf5Readout(learnedState,q,cmp,trained)>=.5
     case "untrained_same_recurrence":answer[j]=wf5Readout(learnedState,q,cmp,init)>=.5
     case "frozen_F4_coded_fourhop":answer[j]=wf4Answer(learnedState,q,cmp)>=.5
     case "coded_exact_upper_control":answer[j]=wf1Answer(fixedState,q)>=.5
     }
     prev:=totals[name]
     if answer[j]==d.Labels[j]{prev[0]++}
     if !d.Labels[j]&&answer[j]{prev[1]++}
     totals[name]=prev
    }
    row.Models[name]=answer
   }
   for _,truth:=range d.Labels{if !truth{neg++}}
   rows=append(rows,row)
  }}
 }}
 if len(rows)!=192||neg!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F5 shape or memory drift")}
 scores:=map[string]map[string]float64{}
 for name,s:=range totals{
  scores[name]=map[string]float64{"accuracy":float64(s[0])/384,
   "correct":float64(s[0]),"false_positives":float64(s[1]),
   "negative_fp_rate":float64(s[1])/192}
 }
 w:=scores["learned_recurrent_composition"];base:=scores["untrained_same_recurrence"]
 upper:=scores["frozen_F4_coded_fourhop"]
 positive:=w["accuracy"]>=.85&&w["negative_fp_rate"]<=.15&&
  w["accuracy"]-base["accuracy"]>=.15&&upper["accuracy"]-w["accuracy"]<=.15
 verdict:="DEVELOPMENT_NO_LEARNED_COMPOSITION_TRANSFER"
 if positive{verdict="DEVELOPMENT_FROZEN_MATCHER_LEARNED_COMPOSITION_OBSERVED"}
 return map[string]any{"schema":"wingless.foundation.f5.frozen-matcher-recurrence-dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_F4_sha":"30577b1e86354e3af3167cd59193819646a3a82e",
  "frozen_F4_matcher":map[string]any{"weights":cmp.Weight,"bias":cmp.Bias},
  "frozen_matcher_parameter_count":17,"learned_composer_parameter_count":5,
  "new_train_bounds":[]int{44300,44324},"new_validation_bounds":[]int{44600,44612},
  "train_graph_states":288,"validation_graph_states":len(rows),"validation_questions":384,
  "negative_questions":neg,"native_state_bytes":576,
  "training_steps":140,"training_wall_ms":trainingMs,"training_loss":loss,
  "trained_composer":trained.W,"scores":scores,"rows":rows,"registered_verdict":verdict,
  "coded_equality_in_candidate":false,"coded_or_traversal_in_candidate":false,
  "engineered_bits":true,"programmed_three_hops":true,
  "sealed_first_use":false,"scientifically_accepted":false,
  "accepted_ref_mutation":false,"external_llm_calls":0,
  "original_development_wall_ms":time.Since(started).Milliseconds()},nil
}
