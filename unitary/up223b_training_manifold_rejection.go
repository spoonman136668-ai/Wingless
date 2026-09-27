package unitary

import "math"

const UP223BManifoldSchema="wingless.up223b-training-manifold-rejection.v1"

type UP223BDecision struct{
 Family string `json:"family"`; Phase int `json:"phase"`; RoutedRegime string `json:"routed_regime"`
 Margin float64 `json:"margin"`; CentroidDistance float64 `json:"centroid_distance"`; ExemplarDistance float64 `json:"exemplar_distance"`
 BaselineRejected bool `json:"baseline_rejected"`; AugmentedRejected bool `json:"augmented_rejected"`
 Harmful bool `json:"harmful"`; RoutedMAE float64 `json:"routed_mae"`; StaticMAE float64 `json:"static_mae"`
}
type UP223BSummary struct{
 Gate string `json:"gate"`; Scope string `json:"scope"`; Family string `json:"family,omitempty"`
 Decisions int `json:"decisions"`; Rejected int `json:"rejected"`; Harmful int `json:"harmful"`; Benign int `json:"benign"`
 HarmfulRejected int `json:"harmful_rejected"`; HarmfulAccepted int `json:"harmful_accepted"`; BenignRejected int `json:"benign_rejected"`
 HarmfulRejectionRate float64 `json:"harmful_rejection_rate"`; BenignRejectionRate float64 `json:"benign_rejection_rate"`
}
type UP223BResult struct{
 Schema string `json:"schema"`; Experiment string `json:"experiment"`; SourceUP222BSeal string `json:"source_up222b_seal"`
 TrainingCorrectStates int `json:"training_correct_states"`
 MarginThreshold float64 `json:"margin_threshold"`; CentroidDistanceThreshold float64 `json:"centroid_distance_threshold"`; ManifoldDistanceThreshold float64 `json:"manifold_distance_threshold"`
 EvaluationPhases []int `json:"evaluation_phases"`
 UnseenFamilyFittingUsed bool `json:"unseen_family_fitting_used"`; NewNativeFeatureUsed bool `json:"new_native_feature_used"`
 EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`; AdaptiveGateSelectionUsed bool `json:"adaptive_gate_selection_used"`
 NonlinearModelUsed bool `json:"nonlinear_model_used"`; PhaseInputUsed bool `json:"phase_input_used"`; ParityInputUsed bool `json:"parity_input_used"`
 LiveActivation bool `json:"live_activation"`
 Decisions []UP223BDecision `json:"decisions"`; Summaries []UP223BSummary `json:"summaries"`
}
type up223bKnown struct{name string;canonical []int;centroid [4]float64;coef [5]float64}
func up223bVDist(a,b [4]float64)float64{v:=0.0;for i:=0;i<4;i++{d:=a[i]-b[i];v+=d*d};return math.Sqrt(v)}
func RunUP223B()(UP223BResult,error){
 template:=[]int{55,56,57};train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
 eval:=[]int{352,353,354,355,356,357,358,359,360,361,362,363,364,365,366,367,368,369,370,371,372,373,374,375,376,377,378,379,380,381,382,383}
 known:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
 unseen:=[]UP197BComposition{{Name:"mixed4b",Canonical:[]int{2,7,3,8}},{Name:"observe3",Canonical:[]int{5,6,8}},{Name:"store3",Canonical:[]int{0,2,3}},{Name:"cross4",Canonical:[]int{0,5,8,13}}}
 model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
 type raw struct{name string;x [4]float64}; raws:=[]raw{}
 for _,f:=range known{for _,ph:=range train{raws=append(raws,raw{f.Name,up193bFeature(ph,f.Canonical)})}}
 var mean,std [4]float64
 for _,r:=range raws{for j:=0;j<4;j++{mean[j]+=r.x[j]}};for j:=0;j<4;j++{mean[j]/=float64(len(raws))}
 for _,r:=range raws{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
 for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(raws)));if std[j]==0{std[j]=1}}
 zfun:=func(x [4]float64)(z [4]float64){for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]};return}
 kms:=[]up223bKnown{}
 for _,f:=range known{
  km:=up223bKnown{name:f.Name,canonical:f.Canonical};n:=0
  for _,r:=range raws{if r.name==f.Name{z:=zfun(r.x);for j:=0;j<4;j++{km.centroid[j]+=z[j]};n++}}
  for j:=0;j<4;j++{km.centroid[j]/=float64(n)}
  perms:=up196bPermutations(f.Canonical);means:=make([]float64,len(perms));canon:=-1
  for i,p:=range perms{if up201bKey(p)==up201bKey(f.Canonical){canon=i};s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor};means[i]=s/float64(len(template))}
  var a [5][5]float64;var b [5]float64
  for _,ph:=range train{x:=up193bFeature(ph,f.Canonical);z:=zfun(x);v:=[5]float64{1,z[0],z[1],z[2],z[3]};pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph);y:=pt.RequiredMassFactor-means[canon];for i:=0;i<5;i++{b[i]+=v[i]*y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}}
  for i:=1;i<5;i++{a[i][i]+=1e-6};c,ok:=up193bSolve(a,b);if !ok{return UP223BResult{},nil};km.coef=c;kms=append(kms,km)
 }
 marginT:=math.Inf(1);centT:=0.0;correct:=[][4]float64{}
 for fi,f:=range known{for _,ph:=range train{
  z:=zfun(up193bFeature(ph,f.Canonical));bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
  for gi,g:=range kms{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
  if bestIdx==fi{correct=append(correct,z);m:=second-best;if m<marginT{marginT=m};if best>centT{centT=best}}
 }}
 manifoldT:=0.0
 for i,z:=range correct{near:=math.Inf(1);for j,q:=range correct{if i==j{continue};d:=up223bVDist(z,q);if d<near{near=d}};if !math.IsInf(near,1)&&near>manifoldT{manifoldT=near}}
 res:=UP223BResult{Schema:UP223BManifoldSchema,Experiment:"UP-223B-training-manifold-rejection",SourceUP222BSeal:"c1c98b5fa5ffbc5760139ecf234406b2e874d53c",TrainingCorrectStates:len(correct),MarginThreshold:marginT,CentroidDistanceThreshold:centT,ManifoldDistanceThreshold:manifoldT,EvaluationPhases:eval,UnseenFamilyFittingUsed:false,NewNativeFeatureUsed:false,EvaluationDerivedThresholdUsed:false,AdaptiveGateSelectionUsed:false,NonlinearModelUsed:false,PhaseInputUsed:false,ParityInputUsed:false,LiveActivation:false}
 for _,f:=range unseen{
  perms:=up196bPermutations(f.Canonical);means:=make([]float64,len(perms));canon:=-1
  for i,p:=range perms{if up201bKey(p)==up201bKey(f.Canonical){canon=i};s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor};means[i]=s/float64(len(template))}
  for _,ph:=range eval{
   x:=up193bFeature(ph,f.Canonical);z:=zfun(x);v:=[5]float64{1,z[0],z[1],z[2],z[3]}
   bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1);for gi,g:=range kms{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
   margin:=second-best;ex:=math.Inf(1);for _,q:=range correct{if d:=up223bVDist(z,q);d<ex{ex=d}}
   shift:=0.0;for i:=0;i<5;i++{shift+=kms[bestIdx].coef[i]*v[i]}
   canonPt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_eval",SurfaceIndices:f.Canonical},ph);_ = canonPt
   sumR,sumS:=0.0,0.0;n:=0
   for i,p:=range perms{if i==canon{continue};pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_eval",SurfaceIndices:p},ph);sp:=means[i];sumR+=math.Abs(sp+shift-pt.RequiredMassFactor);sumS+=math.Abs(sp-pt.RequiredMassFactor);n++}
   rmae,smae:=sumR/float64(n),sumS/float64(n);base:=margin<marginT||best>centT;aug:=base||ex>manifoldT
   res.Decisions=append(res.Decisions,UP223BDecision{Family:f.Name,Phase:ph,RoutedRegime:kms[bestIdx].name,Margin:margin,CentroidDistance:best,ExemplarDistance:ex,BaselineRejected:base,AugmentedRejected:aug,Harmful:rmae>smae,RoutedMAE:rmae,StaticMAE:smae})
  }
 }
 for _,gate:=range []string{"baseline","augmented"}{for _,fam:=range []string{"mixed4b","observe3","store3","cross4",""}{
  sc:="family";if fam==""{sc="overall"};s:=UP223BSummary{Gate:gate,Scope:sc,Family:fam}
  for _,d:=range res.Decisions{if fam!=""&&d.Family!=fam{continue};rej:=d.BaselineRejected;if gate=="augmented"{rej=d.AugmentedRejected};s.Decisions++;if rej{s.Rejected++};if d.Harmful{s.Harmful++;if rej{s.HarmfulRejected++}else{s.HarmfulAccepted++}}else{s.Benign++;if rej{s.BenignRejected++}}}
  if s.Harmful>0{s.HarmfulRejectionRate=float64(s.HarmfulRejected)/float64(s.Harmful)};if s.Benign>0{s.BenignRejectionRate=float64(s.BenignRejected)/float64(s.Benign)}
  res.Summaries=append(res.Summaries,s)
 }}
 return res,nil
}
