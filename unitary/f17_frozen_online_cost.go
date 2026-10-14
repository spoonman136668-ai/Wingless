package unitary

import("fmt";"math";"time")

// W-F17 genuine online frozen-policy test: no full trace is precomputed.
// Candidate sees only 576B serialized graph, bounded role/depth query,
// and current recurrent prefix. On stop, later core steps are NOT executed.
// F16 weights/calibration threshold pinned from artifact 11935598224.
type wf17Sample struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 Depth int `json:"depth"`
 Variant int `json:"variant"`
 Truth bool `json:"truth"`
 Frozen wf15Choice `json:"frozen_online"`
 Full wf15Choice `json:"fixed_full_online"`
 Static wf15Choice `json:"static_one_hop_control"`
 FrozenNs int64 `json:"frozen_ns"`
 FullNs int64 `json:"full_ns"`
 StaticNs int64 `json:"static_ns"`
}
func wf17World(seed,k,fam,variant int)(wf13World,error){
 var w wf13World
 if k<7||k>8||fam<92||fam>95{return w,fmt.Errorf("F17 invalid graph spec")}
 var raw [12][2]int
 for i:=0;i<k;i++{raw[i]=[2]int{i,i+1}}
 switch fam{
 case 92:raw[k]=[2]int{0,3}
 case 93:raw[k]=[2]int{3,1}
 case 94:raw[k]=[2]int{k,0}
 case 95:raw[k]=[2]int{2,5}
 }
 for i:=k+1;i<12;i++{x:=30+i-k;raw[i]=[2]int{x,x+1}}
 var roles [12]int
 for i:=range roles{roles[i]=(seed*3+fam*5+i*3+variant)%4}
 for i:=11;i>0;i--{
  j:=(seed*17+fam*19+k*23+variant*11+i*29)%(i+1)
  raw[i],raw[j]=raw[j],raw[i];roles[i],roles[j]=roles[j],roles[i]
 }
 seen:=map[[2]int]bool{}
 for i,p:=range raw{
  a,b:=wf1Rename(p[0],seed,2),wf1Rename(p[1],seed,2)
  if seen[[2]int{a,b}]{return w,fmt.Errorf("duplicate edge")}
  seen[[2]int{a,b}]=true
  w.Facts.Edges[i]=nci02WEdge{From:a,To:b};w.Roles[i]=roles[i]
 }
 w.Src,w.Dst=wf1Rename(0,seed,2),wf1Rename(k,seed,2);w.K=k
 for i:=0;i<k;i++{
  a,b:=wf1Rename(i,seed,2),wf1Rename(i+1,seed,2)
  found:=false
  for j,e:=range w.Facts.Edges{
   if e.From==a&&e.To==b{w.Required[i]=w.Roles[j];found=true;break}
  }
  if !found{return w,fmt.Errorf("missing true chain")}
 }
 return w,nil
}
func wf17Online(state wf1State,q [2]int,pattern [8]int,k int,
 rm wf13Roles,enc wf7Encoder,core wf5Composer,policy wf16Weights,mode string)wf15Choice{
 var src,dst,role [12]int
 for i:=0;i<12;i++{
  bit:=math.Float32bits(state[i]);src[i],dst[i]=int(bit>>16),int(bit&65535)
  role[i]=int(state[12+i])
 }
 if mode=="static"{
  found:=false
  for i:=0;i<12;i++{
   if src[i]==q[0]&&role[i]==pattern[0]{found=true}
  }
  if !found{return wf15Choice{Answer:false,Steps:0,EarlyHalt:true}}
 }
 var active [12]float64
 prev:=0.
 for h:=0;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   val:=0.
   if h==0{val,_,_=wf7Pair(enc,q[0],src[j])}else{
    msg:=0.
    for i:=0;i<12;i++{if i==j{continue}
     sim,_,_:=wf7Pair(enc,dst[i],src[j]);msg+=active[i]*sim
    }
    val=wf1Prob(core.W[0]*msg+core.W[1]*active[j]+core.W[2])
   }
   next[j]=val*wf13RoleScore(rm,role[j],pattern[h])
  }
  active=next
  maxActive,sumActive,terminal:=0.,0.,0.
  for i:=0;i<12;i++{
   maxActive=math.Max(maxActive,active[i]);sumActive+=active[i]
   sim,_,_:=wf7Pair(enc,dst[i],q[1]);terminal+=active[i]*sim
  }
  prob:=wf1Prob(core.W[3]*terminal+core.W[4])
  halt:=false
  if mode=="risk"&&prob<.5{
   maxRole,meanRole:=0.,0.
   if h+1<k{
    for i:=0;i<12;i++{
     val:=wf13RoleScore(rm,role[i],pattern[h+1])
     maxRole=math.Max(maxRole,val);meanRole+=val/12
    }
   }
   x:=[8]float64{1,float64(h+1)/float64(k),maxActive,sumActive/12,
     prob,prob-prev,maxRole,meanRole}
   score:=0.;for j,v:=range x{score+=policy.W[j]*v}
   halt=wf1Prob(score)>=.5
  }
  if halt||h==k-1{
   return wf15Choice{Answer:prob>=.5,Steps:h+1,EarlyHalt:halt&&h<k-1}
  }
  prev=prob
 }
 panic("invalid online inference")
}
func wf17Metrics(rows []wf17Sample,name string)map[string]float64{
 correct,pos,neg,fp,fn,steps,full,nanos:=0,0,0,0,0,0,0,int64(0)
 byFamily:=map[int][3]int{}
 for _,r:=range rows{
  var value wf15Choice;var elapsed int64
  switch name{
  case "frozen":value,elapsed=r.Frozen,r.FrozenNs
  case "full":value,elapsed=r.Full,r.FullNs
  case "static":value,elapsed=r.Static,r.StaticNs
  }
  if value.Answer==r.Truth{correct++}
  if r.Truth{pos++;if !value.Answer{fn++}}else{neg++;if value.Answer{fp++}}
  steps+=value.Steps;full+=r.Depth;nanos+=elapsed
  by:=byFamily[r.Family]
  if r.Truth&&!value.Answer{by[0]++};if !r.Truth&&value.Answer{by[1]++}
  by[2]++;byFamily[r.Family]=by
 }
 worstFN,worstFP:=0.,0.
 for _,s:=range byFamily{
  // Exactly 2 positive and 4 negatives per graph, hence denominator.
  pcount:=s[2]/3;ncount:=s[2]-pcount
  worstFN=math.Max(worstFN,float64(s[0])/float64(pcount))
  worstFP=math.Max(worstFP,float64(s[1])/float64(ncount))
 }
 return map[string]float64{"correct":float64(correct),
  "accuracy":float64(correct)/float64(len(rows)),
  "false_negatives":float64(fn),"false_negative_rate":float64(fn)/float64(pos),
  "false_positives":float64(fp),"false_positive_rate":float64(fp)/float64(neg),
  "worst_family_FNR":worstFN,"worst_family_FPR":worstFP,
  "actual_recurrent_steps":float64(steps),"full_recurrent_steps":float64(full),
  "recurrent_step_saving":1-float64(steps)/float64(full),
  "online_total_wall_ns":float64(nanos),
  "wall_ns_per_answer":float64(nanos)/float64(len(rows))}
}
func RunWFoundationF17()(map[string]any,error){
 t0:=time.Now()
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 rm:=wf14FrozenRoles()
 rec:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,
  -3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 policy:=wf16Weights{W:[8]float64{
  .5233848397033497,.111980090092255,-4.3433500918591745,
  -4.043757386252545,.5349164073718833,1.7913415291476846,
  .5233848395808397,.38694717564180703}}
 rows:=[]wf17Sample{}
 graphCount:=0
 for seed:=63200;seed<63212;seed++{for family:=92;family<=95;family++{
  for k:=7;k<=8;k++{for shift:=0;shift<2;shift++{
   w,err:=wf17World(seed,k,family,shift);if err!=nil{return nil,err}
   base,err:=wf14Serialize(w);if err!=nil{return nil,err}
   q:=[2]int{w.Src,w.Dst}
   late:=-1;early:=-1;noise:=-1
   for i,e:=range w.Facts.Edges{
    if e.From==wf1Rename(k-1,seed,2)&&e.To==wf1Rename(k,seed,2){late=i}
    if e.From==wf1Rename(1,seed,2)&&e.To==wf1Rename(2,seed,2){early=i}
    if e.From==wf1Rename(map[int]int{92:0,93:3,94:k,95:2}[family],seed,2)&&
       e.To==wf1Rename(map[int]int{92:3,93:1,94:0,95:5}[family],seed,2){noise=i}
   }
   if late<0||early<0||noise<0{return nil,fmt.Errorf("bridge IDs not found")}
   wrong:=w.Required;wrong[0],wrong[1]=wrong[1],wrong[0]
   variants:=[]struct{change int;seq [8]int;depth int;expected bool}{
    {-1,w.Required,k,true},{late,w.Required,k,false},
    {early,w.Required,k,false},{-1,wrong,k,false},
    {-1,w.Required,k-1,false},{noise,w.Required,k,true}}
   for variant,v:=range variants{
    check:=wf13Exact(w,v.seq,v.depth,nil)
    if v.change>=0{check=wf13Exact(w,v.seq,v.depth,&v.change)}
    if check!=v.expected{
     return nil,fmt.Errorf("coded reference mismatch fam=%d seed=%d type=%d",family,seed,variant)
    }
    packed:=base
    if v.change>=0{packed[12+v.change]=float32((w.Roles[v.change]+1)%4)}
    // Equal decoding and decision path for BOTH frozen online and full.
    sample:=wf17Sample{Seed:seed,Family:family,Depth:v.depth,Variant:variant,Truth:v.expected}
    // Warm both arms, discarding JIT/cache warm-up measurements.
    _=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"risk")
    _=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"full")
    _=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"static")
    // Alternate order over graph/variant; no end-to-end precomputation.
    if (seed+family+variant)%2==0{
     start:=time.Now();sample.Frozen=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"risk");sample.FrozenNs=time.Since(start).Nanoseconds()
     start=time.Now();sample.Full=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"full");sample.FullNs=time.Since(start).Nanoseconds()
    }else{
     start:=time.Now();sample.Full=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"full");sample.FullNs=time.Since(start).Nanoseconds()
     start=time.Now();sample.Frozen=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"risk");sample.FrozenNs=time.Since(start).Nanoseconds()
    }
    start:=time.Now();sample.Static=wf17Online(packed,q,v.seq,v.depth,rm,enc,rec,policy,"static");sample.StaticNs=time.Since(start).Nanoseconds()
    rows=append(rows,sample)
   }
   graphCount++
  }}
 }}
 if graphCount!=192||len(rows)!=1152{return nil,fmt.Errorf("F17 size mismatch")}
 scores:=map[string]map[string]float64{
  "frozen_online":wf17Metrics(rows,"frozen"),
  "fixed_full_online":wf17Metrics(rows,"full"),
  "coded_first_hop_rule_control":wf17Metrics(rows,"static")}
 a,b:=scores["frozen_online"],scores["fixed_full_online"]
 a["observed_online_wall_fraction_saved_vs_full"]=1-a["online_total_wall_ns"]/b["online_total_wall_ns"]
 c:=scores["coded_first_hop_rule_control"]
 c["observed_online_wall_fraction_saved_vs_full"]=1-c["online_total_wall_ns"]/b["online_total_wall_ns"]
 valid:=b["accuracy"]>=.95
 verdict:="INVALID_FULL_REASONING_REFERENCE"
 if valid{
  robust:=a["accuracy"]>=.95&&a["false_negative_rate"]<=.05&&
   a["false_positive_rate"]<=.05&&a["worst_family_FNR"]<=.05&&
   a["worst_family_FPR"]<=.05&&a["accuracy"]>=b["accuracy"]-.03
  efficient:=a["recurrent_step_saving"]>=.15&&
   a["observed_online_wall_fraction_saved_vs_full"]>=.05
  verdict="DEVELOPMENT_FROZEN_ROBUST_AND_WALL_EFFICIENT"
  if !robust{verdict="DEVELOPMENT_NO_FROZEN_POLICY_ROBUSTNESS"}else if !efficient{
   verdict="DEVELOPMENT_ROBUST_BUT_NO_PROVEN_ONLINE_WALL_EFFICIENCY"}
  }
 return map[string]any{"schema":"wingless.f17.frozen-online-new-grammars.v1",
  "stage":"DEVELOPMENT_ONLY","frozen_F16_sha":"0fc2d87f9160cbb1a8f7678992642eef461cf104",
  "F16_original_json_sha256":"5f9607a8bdefee8f27ef39ea1b6ab00ef018a2a1e6ec719b1fa944ff840f6bba",
  "frozen_policy_weights":policy.W,"frozen_decision_threshold":.5,
  "source_birth_witnessed":false,"sealed_first_use":false,"accepted_ref_mutation":false,
  "new_families":[]int{92,93,94,95},"seed_bounds":[]int{63200,63212},
  "graph_count":graphCount,"question_count":len(rows),"positive_count":384,
  "negative_count":768,"native_state_bytes":576,"model_retrained":false,
  "query_depth_is_programmed_maximum":true,"online_halt_skips_unneeded_core_operations":true,
  "no_future_prefix_oracle_for_policy":true,"wall_measurement":"time.Now around genuine online inference (includes policy features & inference, excludes dataset construction and one-time model load); per-invocation warmup and alternating order, one timed pass",
  "wall_timing_is_machine_local_not_generalizable":true,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "external_llm_calls":0,"wall_original_ms":time.Since(t0).Milliseconds()},nil
}
