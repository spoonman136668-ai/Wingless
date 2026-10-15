package unitary
import("fmt";"time")

// W-F19 tests independently redundant typed witnesses (not shortcuts).
// It freezes W-F14's native learned role representation and W-F16's stopping
// head, supplies NEW graph structures, and never touches sealed F3 sources.
type wf19Row struct{
 Seed,Family,Depth,Shift,Variant int
 Truth bool
 Full,Risk,Static wf15Choice
}
func wf19Graph(seed,k,fam,shift int)(wf13World,int,int,int,error){
 var w wf13World
 if fam<108||fam>111||k<7||k>8{return w,0,0,0,fmt.Errorf("unregistered graph")}
 fork:=1+fam-108
 var edges [12][2]int
 var roles [12]int
 for i:=0;i<k;i++{edges[i]=[2]int{i,i+1};roles[i]=(seed*3+fam*5+shift+i*3)%4}
 other:=k+35
 edges[k]=[2]int{fork,other};edges[k+1]=[2]int{other,fork+2}
 roles[k]=roles[fork];roles[k+1]=roles[fork+1]
 for i:=k+2;i<12;i++{edges[i]=[2]int{50+i,51+i};roles[i]=(seed*7+fam+i)%4}
 for i:=11;i>0;i--{
  j:=(seed*17+fam*13+k*23+shift*11+i*29)%(i+1)
  edges[i],edges[j]=edges[j],edges[i]
  roles[i],roles[j]=roles[j],roles[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range edges{
  a,b:=wf1Rename(e[0],seed,2),wf1Rename(e[1],seed,2)
  if seen[[2]int{a,b}]{return w,0,0,0,fmt.Errorf("duplicate edge")}
  seen[[2]int{a,b}]=true
  w.Facts.Edges[i]=nci02WEdge{From:a,To:b};w.Roles[i]=roles[i]
 }
 w.Src,w.Dst=wf1Rename(0,seed,2),wf1Rename(k,seed,2);w.K=k
 find:=func(a,b int)int{
  a=wf1Rename(a,seed,2);b=wf1Rename(b,seed,2)
  for i,e:=range w.Facts.Edges{if e.From==a&&e.To==b{return i}}
  return -1
 }
 for i:=0;i<k;i++{
  ix:=find(i,i+1);if ix<0{return w,0,0,0,fmt.Errorf("missing chain")}
  w.Required[i]=w.Roles[ix]
 }
 main,alt,noise:=find(fork,fork+1),find(fork,other),find(50+k+2,51+k+2)
 if main<0||alt<0||noise<0{return w,0,0,0,fmt.Errorf("missing fork/noise")}
 if w.Roles[main]!=w.Roles[alt]{return w,0,0,0,fmt.Errorf("wrong role on alternative")}
 return w,main,alt,noise,nil
}
func RunWFoundationF19()(map[string]any,error){
 started:=time.Now()
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 rm:=wf14FrozenRoles()
 rec:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 risk:=wf16Weights{W:[8]float64{.5233848397033497,.111980090092255,-4.3433500918591745,-4.043757386252545,.5349164073718833,1.7913415291476846,.5233848395808397,.38694717564180703}}
 rows:=[]wf19Row{}
 expected:=[]bool{true,true,true,false,false,false,true}
 for seed:=73100;seed<73106;seed++{for fam:=108;fam<=111;fam++{
  for k:=7;k<=8;k++{for shift:=0;shift<2;shift++{
   world,main,alt,noise,err:=wf19Graph(seed,k,fam,shift);if err!=nil{return nil,err}
   wrong:=world.Required
   if wrong[0]==wrong[1]{return nil,fmt.Errorf("indistinguishable role swap")}
   wrong[0],wrong[1]=wrong[1],wrong[0]
   for variant,truth:=range expected{
    modified:=world
    queryRoles:=world.Required
    depth:=k
    mutate:=func(i int){modified.Roles[i]=(modified.Roles[i]+1)%4}
    switch variant{
    case 1:mutate(main)
    case 2:mutate(alt)
    case 3:mutate(main);mutate(alt)
    case 4:queryRoles=wrong
    case 5:depth=k-1
    case 6:mutate(noise)
    }
    if got:=wf13Exact(modified,queryRoles,depth,nil);got!=truth{
     return nil,fmt.Errorf("INVALID_DOUBLE_WITNESS_ORACLE seed %d family %d variant %d",seed,fam,variant)
    }
    state,err:=wf14Serialize(modified);if err!=nil{return nil,err}
    q:=[2]int{world.Src,world.Dst}
    rows=append(rows,wf19Row{Seed:seed,Family:fam,Depth:depth,Shift:shift,Variant:variant,
      Truth:truth,Full:wf17Online(state,q,queryRoles,depth,rm,enc,rec,risk,"full"),
      Risk:wf17Online(state,q,queryRoles,depth,rm,enc,rec,risk,"risk"),
      Static:wf17Online(state,q,queryRoles,depth,rm,enc,rec,risk,"static")})
   }
  }}
 }}
 if len(rows)!=672{return nil,fmt.Errorf("science cohort drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"full","risk","static"}{
  correct,fp,fn,steps,total:=0,0,0,0,0
  worst:=0.
  for fam:=108;fam<=111;fam++{
   n,bad:=0,0
   for _,r:=range rows{if r.Family!=fam{continue}
    var p wf15Choice
    switch name{case "full":p=r.Full;case "risk":p=r.Risk;case "static":p=r.Static}
    n++;total+=r.Depth;steps+=p.Steps
    if p.Answer==r.Truth{correct++}else{bad++}
    if p.Answer&&!r.Truth{fp++}
    if !p.Answer&&r.Truth{fn++}
   }
   if float64(bad)/float64(n)>worst{worst=float64(bad)/float64(n)}
  }
  scores[name]=map[string]float64{"correct":float64(correct),"accuracy":float64(correct)/672,
   "false_positives":float64(fp),"false_negatives":float64(fn),
   "worst_family_error_rate":worst,"steps":float64(steps),
   "full_steps":float64(total),"recurrent_steps_saved":1-float64(steps)/float64(total)}
 }
 surviveMain,surviveAlt:=0,0
 for _,r:=range rows{if r.Variant==1&&r.Full.Answer{surviveMain++}
  if r.Variant==2&&r.Full.Answer{surviveAlt++}}
 full:=scores["full"];policy:=scores["risk"]
 fullGood:=full["accuracy"]>=.95&&full["worst_family_error_rate"]<=.05&&surviveMain>=91&&surviveAlt>=91
 riskGood:=policy["accuracy"]>=.95&&policy["worst_family_error_rate"]<=.05
 verdict:="DEVELOPMENT_FROZEN_TWO_WITNESS_TRANSFER"
 if !fullGood{verdict="DEVELOPMENT_NO_FROZEN_TWO_WITNESS_TRANSFER"
 }else if !riskGood{verdict="DEVELOPMENT_TWO_WITNESS_FULL_REASONING_PASS_RISK_STOP_FAILURE"}
 return map[string]any{
  "schema":"wingless.f19.redundant-causal-witness-development.v1",
  "stage":"DEVELOPMENT_ONLY","parent_F18_sha":"5ebcc66b63c4de08b1819d82690f6315a3b61a67",
  "frozen_role_core_source":"0fc2d87f9160cbb1a8f7678992642eef461cf104",
  "new_families":[]int{108,109,110,111},"new_seed_bounds":[]int{73100,73106},
  "graph_states":96,"cases_per_graph":7,"question_count":len(rows),
  "positive_count":384,"negative_count":288,
  "native_state_bytes":576,"model_retrained":false,
  "survive_main_path_mutation":surviveMain,"survive_alternative_mutation":surviveAlt,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "source_birth_independent":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "wall_ms":time.Since(started).Milliseconds()},nil
}
