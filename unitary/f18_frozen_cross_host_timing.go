package unitary

import ("fmt";"math";"sort";"time";"runtime")

// F18: F17 policy, core, graph grammar and threshold frozen. New seeds,
// five independent timing passes per case, counterbalanced execution order.
// Do not count replayed F17 source as independent scientific original.
// Every result uses truly online inference from an encoded 576-byte state.
type wf18Row struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 Depth int `json:"depth"`
 Variant int `json:"variant"`
 Truth bool `json:"truth"`
 Risk wf15Choice `json:"risk"`
 Full wf15Choice `json:"full"`
 Static wf15Choice `json:"static"`
 RiskNs [5]int64 `json:"risk_ns"`
 FullNs [5]int64 `json:"full_ns"`
 StaticNs [5]int64 `json:"static_ns"`
}
func wf18Median(values [5]float64)float64 {
 sorted:=values
 sort.Float64s(sorted[:])
 return sorted[2]
}
func RunWFoundationF18()(map[string]any,error){
 beginning:=time.Now()
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 roles:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,
  -3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 frozen:=wf16Weights{W:[8]float64{
  .5233848397033497,.111980090092255,-4.3433500918591745,
  -4.043757386252545,.5349164073718833,1.7913415291476846,
  .5233848395808397,.38694717564180703}}
 rows:=[]wf18Row{}
 for seed:=68600;seed<68606;seed++{for fam:=92;fam<=95;fam++{
  for k:=7;k<=8;k++{for roleShift:=0;roleShift<2;roleShift++{
   g,err:=wf17World(seed,k,fam,roleShift);if err!=nil{return nil,err}
   state,err:=wf14Serialize(g);if err!=nil{return nil,err}
   query:=[2]int{g.Src,g.Dst}
   edgeIndex:=func(a,b int)int{
    a=wf1Rename(a,seed,2);b=wf1Rename(b,seed,2)
    for i,e:=range g.Facts.Edges{if e.From==a&&e.To==b{return i}}
    return -1
   }
   late,early:=edgeIndex(k-1,k),edgeIndex(1,2)
   extra:=map[int][2]int{92:{0,3},93:{3,1},94:{k,0},95:{2,5}}[fam]
   distractor:=edgeIndex(extra[0],extra[1])
   if late<0||early<0||distractor<0{return nil,fmt.Errorf("missing exact arm")}
   wrong:=g.Required;wrong[0],wrong[1]=wrong[1],wrong[0]
   plans:=[]struct{change int;pattern [8]int;depth int;truth bool}{
    {-1,g.Required,k,true},{late,g.Required,k,false},
    {early,g.Required,k,false},{-1,wrong,k,false},
    {-1,g.Required,k-1,false},{distractor,g.Required,k,true}}
   for variant,v:=range plans{
    oracle:=wf13Exact(g,v.pattern,v.depth,nil)
    if v.change>=0{oracle=wf13Exact(g,v.pattern,v.depth,&v.change)}
    if oracle!=v.truth{return nil,fmt.Errorf("INVALID_F18_ORACLE")}
    packed:=state
    if v.change>=0{packed[12+v.change]=float32((g.Roles[v.change]+1)%4)}
    row:=wf18Row{Seed:seed,Family:fam,Depth:v.depth,Variant:variant,Truth:v.truth}
    // warm code paths but discard predictions/timings
    for _,mode:=range []string{"risk","full","static"}{
     _=wf17Online(packed,query,v.pattern,v.depth,roles,entity,core,frozen,mode)
    }
    for repeat:=0;repeat<5;repeat++{
     // counterbalanced order over records AND rounds
     order:=[]string{"risk","full","static"}
     offset:=(seed+fam+variant+repeat)%3
     for pos:=0;pos<3;pos++{
      mode:=order[(pos+offset)%3]
      start:=time.Now()
      prediction:=wf17Online(packed,query,v.pattern,v.depth,roles,entity,core,frozen,mode)
      ns:=time.Since(start).Nanoseconds()
      switch mode{
      case "risk":
       row.RiskNs[repeat]=ns
       if repeat==0{row.Risk=prediction}else if row.Risk!=prediction{return nil,fmt.Errorf("risk prediction drift")}
      case "full":
       row.FullNs[repeat]=ns
       if repeat==0{row.Full=prediction}else if row.Full!=prediction{return nil,fmt.Errorf("full prediction drift")}
      case "static":
       row.StaticNs[repeat]=ns
       if repeat==0{row.Static=prediction}else if row.Static!=prediction{return nil,fmt.Errorf("static prediction drift")}
      }
     }
    }
    rows=append(rows,row)
   }
  }}
 }}
 if len(rows)!=576{return nil,fmt.Errorf("F18 validation count drift")}
 type totals struct{Correct,FP,FN,Steps int}
 stats:=map[string]totals{}
 var riskNs,fullNs,staticNs [5]int64
 for _,r:=range rows{
  for _,kind:=range []string{"risk","full","static"}{
   var value wf15Choice
   switch kind{case "risk":value=r.Risk;case "full":value=r.Full;case "static":value=r.Static}
   a:=stats[kind];if value.Answer==r.Truth{a.Correct++}
   if r.Truth&&!value.Answer{a.FN++}
   if !r.Truth&&value.Answer{a.FP++}
   a.Steps+=value.Steps;stats[kind]=a
  }
  for i:=0;i<5;i++{riskNs[i]+=r.RiskNs[i];fullNs[i]+=r.FullNs[i];staticNs[i]+=r.StaticNs[i]}
 }
 var gains,staticGains [5]float64
 for i:=0;i<5;i++{
  if fullNs[i]<=0{return nil,fmt.Errorf("invalid wall measurement")}
  gains[i]=1-float64(riskNs[i])/float64(fullNs[i])
  staticGains[i]=1-float64(staticNs[i])/float64(fullNs[i])
 }
 ksum:=0;for _,r:=range rows{ksum+=r.Depth}
 scores:=map[string]any{}
 for _,kind:=range []string{"risk","full","static"}{
  s:=stats[kind]
  scores[kind]=map[string]any{
   "correct":s.Correct,"accuracy":float64(s.Correct)/576,
   "false_positive_count":s.FP,"false_negative_count":s.FN,
   "recurrent_steps":s.Steps,"full_steps":ksum,
   "recurrent_savings":1-float64(s.Steps)/float64(ksum)}
 }
 worst:=math.Inf(1)
 for _,v:=range gains{worst=math.Min(worst,v)}
 verdict:="DEVELOPMENT_NO_REPLICATED_FROZEN_WALL_GAIN"
 if stats["full"].Correct!=576{verdict="INVALID_FROZEN_FULL_ORACLE"
 }else if stats["risk"].Correct<548||stats["risk"].FP>19||stats["risk"].FN>10{
  verdict="DEVELOPMENT_FROZEN_POLICY_ROBUSTNESS_FAILED"
 }else if 1-float64(stats["risk"].Steps)/float64(ksum)>=.15&&
   wf18Median(gains)>=.05&&worst>=.05{
  verdict="DEVELOPMENT_REPLICATED_SINGLE_HOST_WALL_GAIN"
 }
 return map[string]any{
  "schema":"wingless.f18.frozen-online-cross-host-repeats.v1",
  "stage":"DEVELOPMENT_ONLY","source_F17":"3ef00f6db0fda03385730640fe5c0354e10307bf",
  "frozen_policy_source_F16":"0fc2d87f9160cbb1a8f7678992642eef461cf104",
  "frozen_weights":frozen.W,"frozen_threshold":.5,"model_retrained":false,
  "host_os":runtime.GOOS,"host_arch":runtime.GOARCH,
  "new_seed_bounds":[]int{68600,68606},"new_grammar_families":[]int{92,93,94,95},
  "reused_grammar_disclosed":true,"graph_states":96,"questions":576,
  "positive_count":192,"negative_count":384,"repeats_per_arm":5,
  "state_bytes":576,"online_future_steps_skipped":true,
  "wall_scope":"time.Now each full online invocation incl risk features, state decode and inference, excludes source construction and frozen parameter load. Per case warm-up then five balanced passes; cross-host total ratio per pass. Host-local measurement only.",
  "risk_wall_ns_by_repeat":riskNs,"full_wall_ns_by_repeat":fullNs,
  "static_wall_ns_by_repeat":staticNs,
  "risk_fraction_saved_by_repeat":gains,"static_fraction_saved_by_repeat":staticGains,
  "risk_median_fraction_saved":wf18Median(gains),
  "risk_worst_fraction_saved":worst,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "F3_independent_source_birth":false,"F3_sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "wall_ms":time.Since(beginning).Milliseconds()},nil
}
