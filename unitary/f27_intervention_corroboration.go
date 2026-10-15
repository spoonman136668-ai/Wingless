package unitary
import("fmt";"math";"time")

// W-F27: select candidate identity links by scanning node evidence, with NO
// given alias-pair tuples at inference. A four-weight contextual similarity
// metric is trained only on unrelated TRAIN entity pairs, then frozen.
// All possible graph endpoint pairs are enumerated by CODED scaffolding;
// source-supplied entity attributes and a synthetic ontology are NOT learned.
type wf27Row struct{Seed,Family,Depth,Shift,Variant int;Truth bool
 Learned,NoLinks,FixedDistance,StrictDistance,StaticOnly,CausalFixed,CodedUpper wf15Choice
 CandidatePairs int}
type wf27Metric struct{W [5]float64}
func wf27Mix(x uint32)uint32{
 x^=x>>16;x*=0x7feb352d;x^=x>>15;x*=0x846ca68b;x^=x>>16
 return x
}
func wf27Proto(seed,identity int)[3]float64{
 var z [3]float64
 for j:=0;j<3;j++{
  h:=wf27Mix(uint32(identity*919+seed*101+j*487+104729))
  z[j]=float64(h%10001)/5000-1
 }
 return z
}
func wf27Observed(seed,id,semantic,salt int)[3]float64{
 x:=wf27Proto(seed,semantic)
 for i:=0;i<3;i++{
  noise:=wf27Mix(uint32(seed*727+id*193+salt*41+i*677))
  x[i]+=(float64(noise%1001)/500-1)*.015
 }
 return x
}
// Two distinct simulated environmental intervention responses accompany an uninformative static fingerprint.
 // Train and test use disjoint entity identities and distinct source seeds.
 // For negative identity pairs the impostor shares a static prototype but
 // agrees in first intervention but violates response to second. Source provides event readings;
 // the learned model is NOT told whether two nodes are the same.
func wf27TrainingPair(i int)([3]float64,[3]float64,bool){
 seed:=5100+(i%41);sem:=2100+i
 a:=wf27Observed(seed,sem,sem,i)
 b:=wf27Observed(seed,sem+5000,sem,i+300)
 positive:=i%2==0
 if !positive{
   b[2]+=.12
 }
 return a,b,positive
}
func wf27Features(a,b [3]float64)[5]float64{
 var f [5]float64;f[0]=1
 for j:=0;j<3;j++{d:=a[j]-b[j];f[j+1]=100*d*d}
 contrast:=(a[2]-a[1])-(b[2]-b[1]);f[4]=100*contrast*contrast
 return f
}
func wf27Confidence(p wf27Metric,a,b [3]float64)float64{
 x:=wf27Features(a,b);z:=0.
 for i:=range x{z+=p.W[i]*x[i]}
 return wf21Sig(z)
}
func wf27Train()wf27Metric{
 var p wf27Metric
 for epoch:=0;epoch<560;epoch++{
  var g [5]float64
  for i:=0;i<512;i++{
   a,b,pos:=wf27TrainingPair(i)
   x:=wf27Features(a,b)
   z:=0.;for j:=range x{z+=p.W[j]*x[j]}
   y:=0.;if pos{y=1}
   d:=wf21Sig(z)-y
   for j:=range x{g[j]+=d*x[j]}
  }
  for i:=range p.W{
   grad:=g[i]/512
   if grad>10{grad=10};if grad< -10{grad= -10}
   p.W[i]-=.035*(grad+.0002*p.W[i])
  }
 }
 return p
}
func wf27Recurrent(state wf1State,q [2]int,pattern [8]int,k int,
 rm wf13Roles,enc wf7Encoder,core wf5Composer,p wf27Metric,mode string)(wf15Choice,int){
 var src,dst,tags [12]int
 for i:=0;i<12;i++{
  b:=math.Float32bits(state[i])
  src[i],dst[i]=int(b>>16),int(b&65535);tags[i]=int(state[12+i])
 }
 var ids [24]int
 var fingerprints [24][3]float64
 n:=0
 for at:=24;at<120;at+=4{
  id:=int(state[at]);if id==0{continue}
  ids[n]=id
  for j:=0;j<3;j++{fingerprints[n][j]=float64(state[at+j+1])}
  n++
 }
 locate:=func(id int)[3]float64{
  for i:=0;i<n;i++{if ids[i]==id{return fingerprints[i]}}
  return [3]float64{1e6,1e6,1e6}
 }
 var a,b [12][3]float64
 for i:=0;i<12;i++{a[i]=locate(src[i]);b[i]=locate(dst[i])}
 var allowed [12][12]bool
 joins:=0
 for i:=0;i<12;i++{for j:=0;j<12;j++{
  if i==j||dst[i]==src[j]{continue}
  ok:=false
  if mode=="learned"{ok=wf27Confidence(p,b[i],a[j])>=.5}
  if mode=="fixed_distance"||mode=="strict_distance"||mode=="static_only"||mode=="causal_fixed"{
   dist:=0.;for m:=0;m<3;m++{d:=b[i][m]-a[j][m];dist+=d*d}
   cutoff:=.045;if mode=="strict_distance"{cutoff=.0015};if mode=="static_only"{d:=b[i][0]-a[j][0];dist=d*d;cutoff=.0004};if mode=="causal_fixed"{d:=b[i][2]-a[j][2];dist=d*d;cutoff=.009};ok=dist<cutoff
  }
  if ok{allowed[i][j]=true;joins++}
 }}
 var active [12]float64
 for h:=0;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   v:=0.
   if h==0{v,_,_=wf7Pair(enc,q[0],src[j])}else{
    msg:=0.
    for i:=0;i<12;i++{
     if i==j{continue}
     sim,_,_:=wf7Pair(enc,dst[i],src[j])
     if allowed[i][j]{sim=math.Max(sim,1)}
     msg+=active[i]*sim
    }
    v=wf1Prob(core.W[0]*msg+core.W[1]*active[j]+core.W[2])
   }
   next[j]=v*wf13RoleScore(rm,tags[j],pattern[h])
  }
  active=next
 }
 score:=0.
 for i:=0;i<12;i++{sim,_,_:=wf7Pair(enc,dst[i],q[1]);score+=active[i]*sim}
 return wf15Choice{Answer:wf1Prob(core.W[3]*score+core.W[4])>=.5,Steps:k},joins
}
func RunWFoundationF27()(map[string]any,error){
 start:=time.Now();p:=wf27Train()
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 rows:=[]wf27Row{}
 for seed:=116200;seed<116206;seed++{for fam:=108;fam<=111;fam++{
  for depth:=7;depth<=8;depth++{for shift:=0;shift<2;shift++{
   w,main,alt,noise,err:=wf19Graph(seed,depth,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108;original:=wf1Rename(depth+35,seed,2)
   alias:=wf1Rename(depth+90,seed,2)
   unrelated:=wf1Rename(depth+91,seed,2)
   second:=-1
   for j,e:=range w.Facts.Edges{
    if e.From==original&&e.To==wf1Rename(fork+2,seed,2){second=j}
   }
   if second<0{return nil,fmt.Errorf("missing original bridge")}
   q:=[2]int{w.Src,w.Dst}
   for variant:=0;variant<8;variant++{
    graph:=w
    truth:=[]bool{true,true,false,false,true,true,false,true}[variant]
    if variant==1||variant==2||variant==3||variant==4||variant==6{
     graph.Roles[main]=(graph.Roles[main]+1)%4
    }
    if variant==1||variant==3||variant==4{
     graph.Facts.Edges[second].From=alias
    }
    if variant==2||variant==6{
     graph.Facts.Edges[second].From=unrelated
    }
    if variant==5{graph.Roles[noise]=(graph.Roles[noise]+1)%4}
    if variant==7{graph.Roles[main]=(graph.Roles[main]+1)%4}
    if variant==3{graph.Roles[alt]=(graph.Roles[alt]+1)%4}
    if variant==6{graph.Roles[alt]=(graph.Roles[alt]+1)%4}
    reference:=graph
    if (variant==1||variant==4)&&truth{
     reference.Facts.Edges[second].From=original
    }
    if wf13Exact(reference,w.Required,depth,nil)!=truth{
     return nil,fmt.Errorf("INVALID_F27_ORACLE %d/%d/%d",seed,fam,variant)
    }
    state,err:=wf14Serialize(graph);if err!=nil{return nil,err}
    distinct:=map[int]bool{}
    for _,e:=range graph.Facts.Edges{distinct[e.From]=true;distinct[e.To]=true}
    ids:=[]int{}
    for id:=range distinct{ids=append(ids,id)}
    // Sorted to remove Go map nondeterminism in original input construction.
    for i:=0;i<len(ids);i++{for j:=i+1;j<len(ids);j++{if ids[j]<ids[i]{ids[i],ids[j]=ids[j],ids[i]}}}
    if len(ids)>24{return nil,fmt.Errorf("state budget exceeded %d",len(ids))}
    for j,id:=range ids{
     semantic:=id
     if id==alias||id==unrelated{semantic=original}
     // Distinct node profiles supplied by source; the learner is NOT told
     // that alias uses the same latent prototype as original.
     embedding:=wf27Observed(seed,id,semantic,variant)
     if id==unrelated{embedding[2]+=.12}
    if variant==4&&id==alias{embedding[2]+=.08}
     for c:=0;c<3;c++{state[24+j*4+c+1]=float32(embedding[c])}
     state[24+j*4]=float32(id)
    }
    learned,pairs:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"learned")
    none,_:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"none")
    fixed,_:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"fixed_distance")
    strict,_:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"strict_distance")
    static,_:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"static_only")
    fixedCausal,_:=wf27Recurrent(state,q,w.Required,depth,role,entity,core,p,"causal_fixed")
    // Given *semantic oracle*: labels checked from programmed graph
    // canonicalization, not injected into native model features.
    upper:=wf15Choice{Answer:wf13Exact(reference,w.Required,depth,nil),Steps:depth}
    rows=append(rows,wf27Row{Seed:seed,Family:fam,Depth:depth,Shift:shift,
      Variant:variant,Truth:truth,Learned:learned,NoLinks:none,
      FixedDistance:fixed,StrictDistance:strict,StaticOnly:static,CausalFixed:fixedCausal,CodedUpper:upper,CandidatePairs:pairs})
   }
  }}
 }}
 if len(rows)!=768{return nil,fmt.Errorf("W24 cohort count drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"learned","none","fixed_distance","strict_distance","static_only","causal_fixed","coded_oracle"}{
  correct,fp,fn:=0,0,0
  for _,r:=range rows{
   var x wf15Choice
   switch name{case "learned":x=r.Learned;case "none":x=r.NoLinks
   case "fixed_distance":x=r.FixedDistance;case "strict_distance":x=r.StrictDistance;case "static_only":x=r.StaticOnly;case "causal_fixed":x=r.CausalFixed;default:x=r.CodedUpper}
   if x.Answer==r.Truth{correct++}
   if x.Answer&&!r.Truth{fp++};if !x.Answer&&r.Truth{fn++}
  }
  scores[name]=map[string]float64{"correct":float64(correct),
    "accuracy":float64(correct)/768,"fp":float64(fp),"fn":float64(fn)}
 }
 l:=scores["learned"]
 positive:="DEVELOPMENT_CAUSAL_CORROBORATION_NO_GAIN"
 if scores["coded_oracle"]["correct"]!=768{positive="INVALID_F27_SOURCE_ORACLE"
 }else if l["accuracy"]>=.95&&l["accuracy"]-scores["none"]["accuracy"]>=.1&&
  l["accuracy"]-math.Max(scores["causal_fixed"]["accuracy"],math.Max(scores["static_only"]["accuracy"],math.Max(scores["fixed_distance"]["accuracy"],scores["strict_distance"]["accuracy"])))>=.03{
  positive="DEVELOPMENT_CAUSAL_CORROBORATION_LEARNED_GAIN"
 }
 return map[string]any{"schema":"wingless.f27.intervention-response-corroboration.dev.v1",
  "stage":"DEVELOPMENT_ONLY","parent_F26":"c215e14985e5bb971058096dc48e8f6dd20281c7",
  "training_parameters":p.W,"training_only_positive_negative_pairs":512,
  "actual_query_candidate_pairs_supplied":false,
  "all_graph_edge_endpoint_pairs_enumerated_in_code":true,
  "node_profiles_source_generated":true,
  "feature_similarity_learned":true,"trained_temporal_coherence_feature":true,"impostors_share_initial_context_but_diverge_under_second_intervention":true,"hard_near_negative_training":true,"near_negative_validation_source_generated":true,"fixed_strict_distance_cutoff":.0015,"native_state_bytes":576,
  "test_seeds":[]int{116200,116206},"motif_families_reused":[]int{108,109,110,111},
  "graphs":96,"questions":len(rows),"fixed_causal_distance_cutoff":.009,
  "scores":scores,"rows":rows,"registered_verdict":positive,
  "independently_born_source":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "elapsed_ms":time.Since(start).Milliseconds()},nil
}
