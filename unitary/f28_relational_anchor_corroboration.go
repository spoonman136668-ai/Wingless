package unitary
import("fmt";"math";"time")

// W-F28: select candidate identity links by scanning node evidence, with NO
// given alias-pair tuples at inference. A four-weight contextual similarity
// metric is trained only on unrelated TRAIN entity pairs, then frozen.
// All possible graph endpoint pairs are enumerated by CODED scaffolding;
// source-supplied entity attributes and a synthetic ontology are NOT learned.
type wf28Row struct{Seed,Family,Depth,Shift,Variant int;Truth bool
 Learned,NoLinks,FixedDistance,StrictDistance,StaticOnly,CausalFixed,Consensus,CodedUpper wf15Choice
 CandidatePairs int}
type wf28Metric struct{W [8]float64}
func wf28Mix(x uint32)uint32{
 x^=x>>16;x*=0x7feb352d;x^=x>>15;x*=0x846ca68b;x^=x>>16
 return x
}
// Each of three independent synthetic witness channels represents a
// node's observed relationship to a different stable external anchor.
// Truth and identity labels are used ONLY to generate these observations,
// never by the learned inference policy.
func wf28Witness(seed,id,semantic,salt int)[3]float64{
 var x [3]float64
 for j:=0;j<3;j++{
  h:=wf28Mix(uint32((semantic+421*j)*919+seed*101+j*977+104729))
  proto:=float64(h%10001)/5000-1
  noise:=wf28Mix(uint32(seed*727+id*193+salt*41+j*677))
  x[j]=proto+(float64(noise%1001)/500-1)*.012
 }
 return x
}
// Training observes independent node pairs, with positives that may have
// one wrong witness and negatives that copy exactly ONE witness.
func wf28TrainPair(i int)([3]float64,[3]float64,bool){
 seed:=7000+i%53;semantic:=3000+i
 a:=wf28Witness(seed,semantic,semantic,i)
 b:=wf28Witness(seed,semantic+9000,semantic,i+100)
 positive:=i%2==0
 witness:=(i/2)%3
 if positive{
  if i%8==0{b[witness]+=.4}
 }else{
  b=wf28Witness(seed,semantic+9001,semantic+5000,i+100)
  b[witness]=a[witness]+.008
 }
 return a,b,positive
}
func wf28Features(a,b [3]float64)[8]float64{
 var d [3]float64
 for i:=range d{d[i]=math.Abs(a[i]-b[i])}
 return [8]float64{1,d[0],d[1],d[2],math.Min(d[0],d[1]),
  math.Min(d[1],d[2]),math.Min(d[0],d[2]),math.Max(d[0],math.Max(d[1],d[2]))}
}
func wf28Confidence(p wf28Metric,a,b [3]float64)float64{
 f:=wf28Features(a,b);z:=0.
 for i,v:=range f{z+=p.W[i]*v}
 return wf21Sig(z)
}
func wf28Train()wf28Metric{
 var p wf28Metric
 for epoch:=0;epoch<650;epoch++{
  var g [8]float64
  for i:=0;i<768;i++{
   a,b,positive:=wf28TrainPair(i)
   x:=wf28Features(a,b)
   z:=0.;for k,v:=range x{z+=p.W[k]*v}
   y:=0.;if positive{y=1}
   e:=wf21Sig(z)-y
   for k,v:=range x{g[k]+=e*v}
  }
  for k:=range p.W{
   grad:=g[k]/768
   if grad>10{grad=10};if grad< -10{grad= -10}
   p.W[k]-=.11*(grad+.0004*p.W[k])
  }
 }
 return p
}
func wf28Recurrent(state wf1State,q [2]int,pattern [8]int,k int,
 rm wf13Roles,enc wf7Encoder,core wf5Composer,p wf28Metric,mode string)(wf15Choice,int){
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
  if mode=="learned"{ok=wf28Confidence(p,b[i],a[j])>=.5}
  if mode=="fixed_distance"||mode=="strict_distance"||mode=="static_only"||mode=="causal_fixed"||mode=="consensus"{
   dist:=0.;for m:=0;m<3;m++{d:=b[i][m]-a[j][m];dist+=d*d}
   cutoff:=.045;if mode=="strict_distance"{cutoff=.0015};if mode=="static_only"{d:=b[i][0]-a[j][0];dist=d*d;cutoff=.0004};if mode=="causal_fixed"{d:=b[i][2]-a[j][2];dist=d*d;cutoff=.009}
   if mode=="consensus"{
    matches:=0;for z:=0;z<3;z++{if math.Abs(b[i][z]-a[j][z])<.075{matches++}}
    ok=matches>=2
   }else{ok=dist<cutoff}
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
func RunWFoundationF28()(map[string]any,error){
 start:=time.Now();p:=wf28Train()
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 rows:=[]wf28Row{}
 for seed:=142200;seed<142206;seed++{for fam:=108;fam<=111;fam++{
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
     return nil,fmt.Errorf("INVALID_F28_ORACLE %d/%d/%d",seed,fam,variant)
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
     if id==alias{semantic=original}
     // Distinct node profiles supplied by source; the learner is NOT told
     // that alias uses the same latent prototype as original.
     embedding:=wf28Witness(seed,id,semantic,variant)
     if id==unrelated{sample:=wf28Witness(seed,original,original,variant);embedding[0]=sample[0]+.008}
    if variant==4&&id==alias{embedding[0]+=.4}
     for c:=0;c<3;c++{state[24+j*4+c+1]=float32(embedding[c])}
     state[24+j*4]=float32(id)
    }
    learned,pairs:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"learned")
    none,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"none")
    fixed,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"fixed_distance")
    strict,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"strict_distance")
    static,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"static_only")
    consensus,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"consensus")
    fixedCausal,_:=wf28Recurrent(state,q,w.Required,depth,role,entity,core,p,"causal_fixed")
    // Given *semantic oracle*: labels checked from programmed graph
    // canonicalization, not injected into native model features.
    upper:=wf15Choice{Answer:wf13Exact(reference,w.Required,depth,nil),Steps:depth}
    rows=append(rows,wf28Row{Seed:seed,Family:fam,Depth:depth,Shift:shift,
      Variant:variant,Truth:truth,Learned:learned,NoLinks:none,
      FixedDistance:fixed,StrictDistance:strict,StaticOnly:static,CausalFixed:fixedCausal,Consensus:consensus,CodedUpper:upper,CandidatePairs:pairs})
   }
  }}
 }}
 if len(rows)!=768{return nil,fmt.Errorf("W24 cohort count drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"learned","none","fixed_distance","strict_distance","static_only","causal_fixed","consensus","coded_oracle"}{
  correct,fp,fn:=0,0,0
  for _,r:=range rows{
   var x wf15Choice
   switch name{case "learned":x=r.Learned;case "none":x=r.NoLinks
   case "fixed_distance":x=r.FixedDistance;case "strict_distance":x=r.StrictDistance;case "static_only":x=r.StaticOnly;case "causal_fixed":x=r.CausalFixed;case "consensus":x=r.Consensus;default:x=r.CodedUpper}
   if x.Answer==r.Truth{correct++}
   if x.Answer&&!r.Truth{fp++};if !x.Answer&&r.Truth{fn++}
  }
  scores[name]=map[string]float64{"correct":float64(correct),
    "accuracy":float64(correct)/768,"fp":float64(fp),"fn":float64(fn)}
 }
 l:=scores["learned"]
 positive:="DEVELOPMENT_RELATIONAL_CONSENSUS_NO_QUALIFIED_GAIN"
 if scores["coded_oracle"]["correct"]!=768{positive="INVALID_F28_SOURCE_ORACLE"
 }else if l["accuracy"]>=.95&&l["accuracy"]-scores["none"]["accuracy"]>=.1&&
  l["accuracy"]-math.Max(scores["consensus"]["accuracy"],math.Max(scores["causal_fixed"]["accuracy"],math.Max(scores["static_only"]["accuracy"],math.Max(scores["fixed_distance"]["accuracy"],scores["strict_distance"]["accuracy"]))))>=.03{
  positive="DEVELOPMENT_RELATIONAL_CONSENSUS_LEARNED_GAIN"
 }
 return map[string]any{"schema":"wingless.f28.independent-anchor-relationship-corroboration.dev.v1",
  "stage":"DEVELOPMENT_ONLY","parent_F27":"86af22a4e9958826227b11749a6901c9ae84fe8a",
  "training_parameters":p.W,"training_only_positive_negative_pairs":768,
  "actual_query_candidate_pairs_supplied":false,
  "all_graph_edge_endpoint_pairs_enumerated_in_code":true,
  "node_profiles_source_generated":true,
  "feature_similarity_learned":true,"trained_temporal_coherence_feature":true,"impostors_share_one_independent_anchor_relation":true,"source_relational_witness_count":3,"learned_relational_agreement_weights":8,"hard_near_negative_training":true,"near_negative_validation_source_generated":true,"fixed_strict_distance_cutoff":.0015,"native_state_bytes":576,
  "test_seeds":[]int{142200,142206},"motif_families_reused":[]int{108,109,110,111},
  "graphs":96,"questions":len(rows),"fixed_causal_distance_cutoff":.009,"consensus_source_match_threshold":.075,
  "scores":scores,"rows":rows,"registered_verdict":positive,
  "independently_born_source":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "elapsed_ms":time.Since(start).Milliseconds()},nil
}
