package unitary

import("fmt";"math";"time")

// W-F15: train a small stop/continue policy from FROZEN W-F14 prefix-state
// statistics. Stopper inference sees only past/current state and public
// iteration budget; it never uses future oracle/rollout or family identity.
// The candidate core's native memory is an ACTUAL 576-byte serialized state.
// This is a separate author-generated DEVELOPMENT experiment, not F3.
type wf15Record struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 Depth int `json:"depth"`
 Variant int `json:"variant"`
 Label bool `json:"label"`
 Features [][6]float64 `json:"prefix_features"`
 PrefixAnswers []bool `json:"prefix_answers"`
 Decisions map[string]wf15Choice `json:"decisions"`
}
type wf15Choice struct{
 Answer bool `json:"answer"`
 Steps int `json:"steps"`
 EarlyHalt bool `json:"early_halt"`
}
type wf15Policy struct{ W [6]float64 }
func wf15Build(seed,k,family,variant int)(wf13World,error){
 var w wf13World
 if k<4||k>8||family<72||family>79{return w,fmt.Errorf("F15 unregistered family/depth")}
 var edges [12][2]int
 for i:=0;i<k;i++{edges[i]=[2]int{i,i+1}}
 switch family%4{
 case 0:edges[k]=[2]int{0,3}
 case 1:edges[k]=[2]int{k-1,2}
 case 2:edges[k]=[2]int{k,0}
 case 3:edges[k]=[2]int{2,0}
 }
 for i:=k+1;i<12;i++{
  v:=28+i-k;edges[i]=[2]int{v,v+1}
 }
 var roles [12]int
 for i:=range roles{roles[i]=(seed*3+family*5+i*3+variant)%4}
 for i:=11;i>0;i--{
  j:=(seed*17+family*11+k*29+i*31+variant*7)%(i+1)
  edges[i],edges[j]=edges[j],edges[i]
  roles[i],roles[j]=roles[j],roles[i]
 }
 seen:=map[[2]int]bool{}
 for i,p:=range edges{
  a,b:=wf1Rename(p[0],seed,2),wf1Rename(p[1],seed,2)
  if seen[[2]int{a,b}]{return w,fmt.Errorf("duplicate graph edge")}
  seen[[2]int{a,b}]=true
  w.Facts.Edges[i]=nci02WEdge{From:a,To:b}
  w.Roles[i]=roles[i]
 }
 w.Src,w.Dst=wf1Rename(0,seed,2),wf1Rename(k,seed,2)
 w.K=k
 for i:=0;i<k;i++{
  a,b:=wf1Rename(i,seed,2),wf1Rename(i+1,seed,2)
  match:=false
  for j,e:=range w.Facts.Edges{
   if e.From==a&&e.To==b{w.Required[i]=w.Roles[j];match=true;break}
  }
  if !match{return w,fmt.Errorf("missing edge after permute")}
 }
 return w,nil
}
func wf15Trace(state wf1State,query [2]int,sequence [8]int,
 k int,changed *int,roles wf13Roles,entity wf7Encoder,core wf5Composer)([][6]float64,[]bool){
 var src,dst,ty [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(state[i])
  src[i],dst[i]=int(bits>>16),int(bits&65535)
  ty[i]=int(state[12+i])
 }
 var active [12]float64
 features:=make([][6]float64,0,k)
 answers:=make([]bool,0,k)
 previous:=0.
 for step:=0;step<k;step++{
  var next [12]float64
  for j:=0;j<12;j++{
   val:=0.
   if step==0{val,_,_=wf7Pair(entity,query[0],src[j])}else{
    sum:=0.
    for i:=0;i<12;i++{
     if i==j{continue}
     match,_,_:=wf7Pair(entity,dst[i],src[j])
     sum+=active[i]*match
    }
    val=wf1Prob(core.W[0]*sum+core.W[1]*active[j]+core.W[2])
   }
   r:=ty[j];if changed!=nil&&j==*changed{r=(r+1)%4}
   val*=wf13RoleScore(roles,r,sequence[step])
   next[j]=val
  }
  active=next
  maxActive,sumActive:=0.,0.
  terminal:=0.
  for i:=0;i<12;i++{
   maxActive=math.Max(maxActive,active[i])
   sumActive+=active[i]
   m,_,_:=wf7Pair(entity,dst[i],query[1])
   terminal+=active[i]*m
  }
  prob:=wf1Prob(core.W[3]*terminal+core.W[4])
  features=append(features,[6]float64{
   1,float64(step+1)/float64(k),maxActive,sumActive/12,prob,prob-previous})
  answers=append(answers,prob>=.5)
  previous=prob
 }
 return features,answers
}
func wf15Examples(seeds []int,fams []int,depths []int,
 roleModel wf13Roles,entity wf7Encoder,core wf5Composer)([]wf15Record,error){
 out:=[]wf15Record{}
 for _,seed:=range seeds{for _,fam:=range fams{for _,k:=range depths{for variant:=0;variant<2;variant++{
  w,err:=wf15Build(seed,k,fam,variant);if err!=nil{return nil,err}
  packed,err:=wf14Serialize(w);if err!=nil{return nil,err}
  q:=[2]int{w.Src,w.Dst}
  idx:=-1
  // Rotate mutation depth to force both easy/late negatives.
  mutationStep:=1+(seed+fam+variant)%(k-2)
  for j,e:=range w.Facts.Edges{
   if e.From==wf1Rename(mutationStep,seed,2)&&e.To==wf1Rename(mutationStep+1,seed,2){idx=j;break}
  }
  if idx<0{return nil,fmt.Errorf("missing mutation edge")}
  wrong:=w.Required
  wrong[0],wrong[1]=wrong[1],wrong[0]
  specs:=[]struct{seq [8]int;cut *int;truth bool}{
   {w.Required,nil,true},{w.Required,&idx,false},{wrong,nil,false}}
  for vi,spec:=range specs{
   exact:=wf13Exact(w,spec.seq,k,spec.cut)
   if exact!=spec.truth{return nil,fmt.Errorf("typed reference invalid seed %d fam %d",seed,fam)}
   feat,preds:=wf15Trace(packed,q,spec.seq,k,spec.cut,roleModel,entity,core)
   out=append(out,wf15Record{Seed:seed,Family:fam,Depth:k,Variant:vi,
    Label:spec.truth,Features:feat,PrefixAnswers:preds,Decisions:map[string]wf15Choice{}})
  }
 }}}}
 return out,nil
}
func wf15Train(rows []wf15Record)(wf15Policy,int64,float64){
 start:=time.Now()
 model:=wf15Policy{}
 var moments,squares [6]float64
 crossEntropy:=0.
 // The TRAIN-ONLY target is a hindsight-safe early stop: current answer
 // equals actual truth, and every later prefix answer stays correct.
 type sample struct{X [6]float64;Y float64}
 samples:=[]sample{}
 for _,r:=range rows{
  for t:=0;t<r.Depth-1;t++{
   safe:=r.PrefixAnswers[t]==r.Label
   for _,v:=range r.PrefixAnswers[t:] {if v!=r.Label{safe=false}}
   y:=0.;if safe{y=1}
   samples=append(samples,sample{r.Features[t],y})
  }
 }
 for epoch:=1;epoch<=220;epoch++{
  var grad [6]float64;loss:=0.
  for _,s:=range samples{
   z:=0.;for j,v:=range s.X{z+=model.W[j]*v}
   p:=max(.0000001,min(.9999999,wf1Prob(z)))
   delta:=p-s.Y
   for j,v:=range s.X{grad[j]+=delta*v}
   loss-=s.Y*math.Log(p)+(1-s.Y)*math.Log(1-p)
  }
  crossEntropy=loss/float64(len(samples))
  for j:=0;j<6;j++{
   g:=grad[j]/float64(len(samples))
   moments[j]=.9*moments[j]+.1*g
   squares[j]=.999*squares[j]+.001*g*g
   delta:=.07*(moments[j]/(1-math.Pow(.9,float64(epoch))))/
    (math.Sqrt(squares[j]/(1-math.Pow(.999,float64(epoch))))+1e-8)
   model.W[j]-=delta
  }
 }
 return model,time.Since(start).Milliseconds(),crossEntropy
}
func wf15Decide(r wf15Record,m wf15Policy,mode string,threshold float64)wf15Choice{
 for t:=0;t<r.Depth;t++{
  halt:=false
  switch mode{
  case "fixed_full":halt=false
  case "fixed_one":halt=true
  case "simple_threshold":halt=r.Features[t][4]<=threshold
  case "learned":
   score:=0.;for j,v:=range r.Features[t]{score+=v*m.W[j]}
   halt=wf1Prob(score)>=.5
  case "hindsight_oracle":
   halt=r.PrefixAnswers[t]==r.Label
   for _,v:=range r.PrefixAnswers[t:] {if v!=r.Label{halt=false}}
  }
  if halt||t==r.Depth-1{
   return wf15Choice{Answer:r.PrefixAnswers[t],Steps:t+1,EarlyHalt:halt&&t<r.Depth-1}
  }
 }
 panic("unreachable")
}
type wf15Totals struct{Correct,FalsePos,Steps,Full int}
func wf15Score(rows []wf15Record,policy wf15Policy,mode string,threshold float64)map[string]float64{
 var score wf15Totals
 negatives:=0
 for _,r:=range rows{
  c:=wf15Decide(r,policy,mode,threshold)
  if c.Answer==r.Label{score.Correct++}
  if !r.Label{negatives++;if c.Answer{score.FalsePos++}}
  score.Steps+=c.Steps;score.Full+=r.Depth
 }
 n:=float64(len(rows));s:=float64(score.Steps)/float64(score.Full)
 return map[string]float64{"correct":float64(score.Correct),"accuracy":float64(score.Correct)/n,
   "false_positives":float64(score.FalsePos),"negative_false_positive_rate":float64(score.FalsePos)/float64(negatives),
   "steps":float64(score.Steps),"full_steps":float64(score.Full),
   "mean_steps":float64(score.Steps)/n,"compute_saved_fraction":1-s,
   "utility":float64(score.Correct)/n-.15*s}
}
func RunWFoundationF15()(map[string]any,error){
 started:=time.Now()
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 roles:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,
    -3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 trainSeeds:=[]int{};for s:=59300;s<59312;s++{trainSeeds=append(trainSeeds,s)}
 validSeeds:=[]int{};for s:=59700;s<59712;s++{validSeeds=append(validSeeds,s)}
 train,err:=wf15Examples(trainSeeds,[]int{72,73,74,75},[]int{4,5,6},roles,entity,core)
 if err!=nil{return nil,err}
 if len(train)!=864{return nil,fmt.Errorf("train case drift")}
 policy,ms,loss:=wf15Train(train)
 thresholds:=[]float64{.000001,.00001,.0001,.001,.01,.05,.1,.25,.49}
 selected:=thresholds[0];best:=math.Inf(-1)
 for _,v:=range thresholds{
  metric:=wf15Score(train,policy,"simple_threshold",v)
  if metric["utility"]>best{best=metric["utility"];selected=v}
 }
 trainPolicy:=wf15Score(train,policy,"learned",0)
 trainHeuristic:=wf15Score(train,policy,"simple_threshold",selected)
 valid,err:=wf15Examples(validSeeds,[]int{76,77,78,79},[]int{7,8},roles,entity,core)
 if err!=nil{return nil,err}
 if len(valid)!=576{return nil,fmt.Errorf("validation case drift")}
 modes:=[]string{"learned_stop_policy","fixed_complete_budget",
   "train_selected_confidence_threshold","fixed_one_step_shortcut",
   "hindsight_oracle_upper"}
 scores:=map[string]map[string]float64{}
 for _,name:=range modes{
  mode:=map[string]string{"learned_stop_policy":"learned",
   "fixed_complete_budget":"fixed_full",
   "train_selected_confidence_threshold":"simple_threshold",
   "fixed_one_step_shortcut":"fixed_one",
   "hindsight_oracle_upper":"hindsight_oracle"}[name]
  scores[name]=wf15Score(valid,policy,mode,selected)
  for i:=range valid{valid[i].Decisions[name]=wf15Decide(valid[i],policy,mode,selected)}
 }
 got:=scores["learned_stop_policy"];fixed:=scores["fixed_complete_budget"]
 simple:=scores["train_selected_confidence_threshold"]
 success:=got["accuracy"]>=.95&&got["negative_false_positive_rate"]<=.05&&
   got["compute_saved_fraction"]>=.15&&
   got["accuracy"]>=fixed["accuracy"]-.03&&
   got["utility"]>=simple["utility"]+.01
 verdict:="DEVELOPMENT_NO_LEARNED_STOPPING_TRANSFER"
 if success{verdict="DEVELOPMENT_LEARNED_NATIVE_STOPPING_TRANSFER"}
 return map[string]any{"schema":"wingless.f15.learned-prefix-stop-compute-dev.v1",
  "stage":"DEVELOPMENT_ONLY","F14_sha":"c37ef4f1213029a0c84cabf1ff099085df01b372",
  "train_families":[]int{72,73,74,75},"valid_families":[]int{76,77,78,79},
  "train_depths":[]int{4,5,6},"valid_depths":[]int{7,8},
  "train_case_count":len(train),"valid_case_count":len(valid),
  "native_state_bytes":576,"policy_parameter_count":6,
  "policy_train_epochs":220,"policy_training_ms":ms,"train_loss":loss,
  "policy_weights":policy.W,"train_selected_threshold":selected,
  "train_policy_score":trainPolicy,"train_heuristic_score":trainHeuristic,
  "scores":scores,"rows":valid,"registered_verdict":verdict,
  "policy_features_prefix_only":true,"training_targets_use_training_only_future":true,
  "candidate_sees_future_oracle":false,"bounded_forced_stop_at_query_end":true,
  "learned_early_stop_decisions":true,
  "author_generated_development":true,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_model_calls":0,
  "wall_ms":time.Since(started).Milliseconds()},nil
}
