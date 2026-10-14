package unitary

import("fmt";"math";"time")

// W-F16 learns a prefix-safe negative stop with calibration on a SEPARATE
// disjoint seed split. A conservative one-sided Wilson bound on false
// negatives must satisfy the risk cap before any early stopping is admitted.
// Features: F15 prefix stats + current graph's NEXT role presence summary.
// No oracle/future inference. All metrics report recurrence-step proxies
// rather than runtime savings; feature computation is separately disclosed.
type wf16Entry struct{
 wf15Record
 Signals [][8]float64 `json:"prefix_signals"`
 Decisions map[string]wf15Choice `json:"decisions"`
}
type wf16Weights struct{W [8]float64}
func wf16Build(seed,k,fam,v int)(wf13World,error){
 var w wf13World
 if k<4||k>8||fam<80||fam>91{return w,fmt.Errorf("unregistered F16 graph")}
 var edges [12][2]int
 for i:=0;i<k;i++{edges[i]=[2]int{i,i+1}}
 switch fam%4{
 case 0:edges[k]=[2]int{0,3}
 case 1:edges[k]=[2]int{k-1,2}
 case 2:edges[k]=[2]int{k,0}
 case 3:edges[k]=[2]int{2,0}
 }
 for i:=k+1;i<12;i++{u:=28+i-k;edges[i]=[2]int{u,u+1}}
 var roles [12]int
 for i:=0;i<12;i++{roles[i]=(seed*3+fam*5+i*3+v)%4}
 for i:=11;i>0;i--{
  j:=(seed*17+fam*11+k*29+i*31+v*7)%(i+1)
  edges[i],edges[j]=edges[j],edges[i]
  roles[i],roles[j]=roles[j],roles[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range edges{
  a,b:=wf1Rename(e[0],seed,2),wf1Rename(e[1],seed,2)
  if seen[[2]int{a,b}]{return w,fmt.Errorf("duplicate ID edge")}
  seen[[2]int{a,b}]=true
  w.Facts.Edges[i]=nci02WEdge{From:a,To:b}
  w.Roles[i]=roles[i]
 }
 w.Src,w.Dst=wf1Rename(0,seed,2),wf1Rename(k,seed,2)
 w.K=k
 for i:=0;i<k;i++{
  a,b:=wf1Rename(i,seed,2),wf1Rename(i+1,seed,2)
  good:=false
  for j,e:=range w.Facts.Edges{
   if e.From==a&&e.To==b{w.Required[i]=w.Roles[j];good=true;break}
  }
  if !good{return w,fmt.Errorf("missing chain")}
 }
 return w,nil
}
func wf16Data(seeds []int,fams []int,depths []int,
 rmod wf13Roles,enc wf7Encoder,rec wf5Composer)([]wf16Entry,error){
 rows:=[]wf16Entry{}
 for _,seed:=range seeds{for _,fam:=range fams{for _,k:=range depths{for v:=0;v<2;v++{
  w,err:=wf16Build(seed,k,fam,v);if err!=nil{return nil,err}
  state,err:=wf14Serialize(w);if err!=nil{return nil,err}
  q:=[2]int{w.Src,w.Dst}
  targetDepth:=1+(seed+fam+v)%(k-2)
  cut:=-1
  for i,e:=range w.Facts.Edges{
   if e.From==wf1Rename(targetDepth,seed,2)&&e.To==wf1Rename(targetDepth+1,seed,2){cut=i;break}
  }
  if cut<0{return nil,fmt.Errorf("cut index unavailable")}
  wrong:=w.Required;wrong[0],wrong[1]=wrong[1],wrong[0]
  spec:=[]struct{order [8]int;change *int;truth bool}{
   {w.Required,nil,true},{w.Required,&cut,false},{wrong,nil,false}}
  for variant,s:=range spec{
   if wf13Exact(w,s.order,k,s.change)!=s.truth{
    return nil,fmt.Errorf("invalid typed oracle label")}
   base,pred:=wf15Trace(state,q,s.order,k,s.change,rmod,enc,rec)
   signals:=make([][8]float64,0,k)
   for i,b:=range base{
    nextMax,nextMean:=0.,0.
    if i+1<k{
     for e:=0;e<12;e++{
      r:=int(state[12+e])
      if s.change!=nil&&e==*s.change{r=(r+1)%4}
      sim:=wf13RoleScore(rmod,r,s.order[i+1])
      nextMax=math.Max(nextMax,sim);nextMean+=sim/12
     }
    }
    signals=append(signals,[8]float64{b[0],b[1],b[2],b[3],b[4],b[5],
      nextMax,nextMean})
   }
   rows=append(rows,wf16Entry{wf15Record:wf15Record{
    Seed:seed,Family:fam,Depth:k,Variant:variant,Label:s.truth,
    Features:base,PrefixAnswers:pred},Signals:signals,
    Decisions:map[string]wf15Choice{}})
  }
 }}}}
 return rows,nil
}
func wf16Train(rows []wf16Entry)(wf16Weights,float64,int64){
 start:=time.Now()
 var model wf16Weights
 var moment,variance [8]float64
 type sample struct{X [8]float64;Y float64}
 samples:=[]sample{}
 for _,r:=range rows{
  for i:=0;i<r.Depth-1;i++{
   safe:=!r.Label&&!r.PrefixAnswers[i]
   for _,later:=range r.PrefixAnswers[i:]{if later{safe=false}}
   y:=0.;if safe{y=1}
   samples=append(samples,sample{r.Signals[i],y})
  }
 }
 loss:=0.
 for epoch:=1;epoch<=220;epoch++{
  var grad [8]float64;loss=0.
  for _,s:=range samples{
   z:=0.;for j,v:=range s.X{z+=model.W[j]*v}
   p:=max(.0000001,min(.9999999,wf1Prob(z)))
   weight:=22.;if s.Y==1{weight=1}
   loss-=weight*(s.Y*math.Log(p)+(1-s.Y)*math.Log(1-p))
   for j,v:=range s.X{grad[j]+=weight*(p-s.Y)*v}
  }
  loss/=float64(len(samples))
  for j:=0;j<8;j++{
   g:=grad[j]/float64(len(samples))
   moment[j]=.9*moment[j]+.1*g
   variance[j]=.999*variance[j]+.001*g*g
   mm:=moment[j]/(1-math.Pow(.9,float64(epoch)))
   vv:=variance[j]/(1-math.Pow(.999,float64(epoch)))
   model.W[j]-=.06*mm/(math.Sqrt(vv)+1e-8)
  }
 }
 return model,loss,time.Since(start).Milliseconds()
}
func wf16Decision(r wf16Entry,weights wf16Weights,mode string,
 threshold float64)wf15Choice{
 for i:=0;i<r.Depth;i++{
  negative:=!r.PrefixAnswers[i]
  stop:=false
  switch mode{
  case "full":stop=false
  case "aggressive":stop=negative
  case "confidence":stop=negative&&r.Signals[i][4]<threshold
  case "learned":
   value:=0.;for j,f:=range r.Signals[i]{value+=weights.W[j]*f}
   stop=negative&&wf1Prob(value)>=threshold
  case "oracle":
   stop=negative&&!r.Label
   for _,future:=range r.PrefixAnswers[i:]{if future{stop=false}}
  }
  if stop||i==r.Depth-1{
   return wf15Choice{Answer:r.PrefixAnswers[i],Steps:i+1,
    EarlyHalt:stop&&i<r.Depth-1}
  }
 }
 panic("invalid recursion stop")
}
func wf16Metric(rows []wf16Entry,weights wf16Weights,mode string,v float64)map[string]float64{
 correct,fp,fn,steps,full,pos,neg:=0,0,0,0,0,0,0
 for _,r:=range rows{
  dec:=wf16Decision(r,weights,mode,v)
  if dec.Answer==r.Label{correct++}
  if r.Label{pos++;if !dec.Answer{fn++}}else{neg++;if dec.Answer{fp++}}
  steps+=dec.Steps;full+=r.Depth
 }
 accuracy:=float64(correct)/float64(len(rows))
 ratio:=float64(steps)/float64(full)
 return map[string]float64{"correct":float64(correct),"accuracy":accuracy,
  "false_negatives":float64(fn),"false_negative_rate":float64(fn)/float64(pos),
  "false_positives":float64(fp),"false_positive_rate":float64(fp)/float64(neg),
  "steps":float64(steps),"full_steps":float64(full),
  "recurrent_step_saving":1-ratio,"utility_proxy":accuracy-.15*ratio}
}
func wf16WilsonUpper(count,total int)float64{
 z:=1.645
 n:=float64(total);p:=float64(count)/n
 return (p+z*z/(2*n)+z*math.Sqrt((p*(1-p)+z*z/(4*n))/n))/(1+z*z/n)
}
func wf16RiskPass(m map[string]float64,positives int)bool{
 if m["accuracy"]<.995{return false}
 return wf16WilsonUpper(int(m["false_negatives"]),positives)<=.05
}
func RunWFoundationF16()(map[string]any,error){
 started:=time.Now()
 encoder,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 rmod:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,
  -3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 mk:=func(a,b int)[]int{r:=[]int{};for i:=a;i<b;i++{r=append(r,i)};return r}
 fit,err:=wf16Data(mk(61300,61312),[]int{80,81,82,83},[]int{4,5,6},rmod,encoder,core)
 if err!=nil{return nil,err}
 if len(fit)!=864{return nil,fmt.Errorf("train size drift")}
 weights,loss,trainMs:=wf16Train(fit)
 calibrate,err:=wf16Data(mk(61800,61806),[]int{80,81,82,83},[]int{4,5,6},rmod,encoder,core)
 if err!=nil{return nil,err}
 if len(calibrate)!=432{return nil,fmt.Errorf("calibration shape drift")}
 riskGrid:=[]float64{.5,.6,.7,.8,.9,.95,.98,.99,.995,.999,1.000001}
 simpleGrid:=[]float64{0,.0000001,.000001,.00001,.0001,.001,.01,.05,.1}
 pick:=func(grid []float64,mode string)float64{
  chosen:=grid[len(grid)-1];best:=-1.
  for _,threshold:=range grid{
   m:=wf16Metric(calibrate,weights,mode,threshold)
   if wf16RiskPass(m,144)&&m["recurrent_step_saving"]>best{
    best=m["recurrent_step_saving"];chosen=threshold
   }
  }
  return chosen
 }
 riskThreshold:=pick(riskGrid,"learned")
 simpleThreshold:=pick(simpleGrid,"confidence")
 val,err:=wf16Data(mk(62300,62312),[]int{84,85,86,87},[]int{7,8},rmod,encoder,core)
 if err!=nil{return nil,err}
 if len(val)!=576{return nil,fmt.Errorf("validation size drift")}
 names:=[]string{"risk_calibrated_learned_stop","calibration_safe_confidence",
  "full_fixed_budget","F15_aggressive_negative","hindsight_oracle_upper"}
 modes:=map[string]string{"risk_calibrated_learned_stop":"learned",
   "calibration_safe_confidence":"confidence","full_fixed_budget":"full",
   "F15_aggressive_negative":"aggressive","hindsight_oracle_upper":"oracle"}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  threshold:=riskThreshold;if name=="calibration_safe_confidence"{threshold=simpleThreshold}
  scores[name]=wf16Metric(val,weights,modes[name],threshold)
  for i:=range val{val[i].Decisions[name]=wf16Decision(val[i],weights,modes[name],threshold)}
 }
 actual:=scores["risk_calibrated_learned_stop"]
 full:=scores["full_fixed_budget"]
 confidence:=scores["calibration_safe_confidence"]
 success:=actual["accuracy"]>=.95&&actual["false_negative_rate"]<=.05&&
  actual["false_positive_rate"]<=.05&&actual["recurrent_step_saving"]>=.15&&
  actual["accuracy"]>=full["accuracy"]-.03&&
  actual["utility_proxy"]>=confidence["utility_proxy"]+.01
 verdict:="DEVELOPMENT_NO_CALIBRATED_SAFE_STOPPING_TRANSFER"
 if success{verdict="DEVELOPMENT_CALIBRATED_RISK_AWARE_STOPPING_TRANSFER"}
 return map[string]any{"schema":"wingless.f16.calibrated-risk-continuation-dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_F15_sha":"b0f0e049b72df11de47c767f32f36ca887a37472",
  "train_families":[]int{80,81,82,83},"validation_families":[]int{84,85,86,87},
  "train_seed_bounds":[]int{61300,61312},"calibration_seed_bounds":[]int{61800,61806},
  "validation_seed_bounds":[]int{62300,62312},
  "train_questions":len(fit),"calibration_questions":len(calibrate),
  "validation_questions":len(val),"resident_state_bytes":576,
  "new_trainable_weights":8,"training_steps":220,"training_loss":loss,
  "training_ms":trainMs,"risk_weights":weights.W,
  "calibration_one_sided_wilson_z":1.645,
  "calibration_risk_threshold":riskThreshold,
  "calibration_confidence_threshold":simpleThreshold,
  "calibration_learned_score":wf16Metric(calibrate,weights,"learned",riskThreshold),
  "calibration_simple_score":wf16Metric(calibrate,weights,"confidence",simpleThreshold),
  "scores":scores,"rows":val,"registered_verdict":verdict,
  "stop_policy_only_negative":true,
  "calibration_seeds_disjoint_from_fitting_and_validation":true,
  "training_labels_use_training_only_future":true,
  "validation_future_not_visible_to_policy":true,
  "query_depth_hard_limit_programmed":true,
  "reported_savings_recurrence_steps_only":true,
  "feature_cost_excluded_from_savings":true,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_model_calls":0,
  "original_ms":time.Since(started).Milliseconds()},nil
}
