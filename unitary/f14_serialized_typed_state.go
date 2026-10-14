package unitary

import ("fmt";"math";"time";"unsafe")

// W-F14: frozen W-F13 parameters, but candidate now receives only a
// REAL 576-byte serialized state (12 packed edges, 12 role words, 120 zero
// reserve words) plus separately budgeted query. Fresh DEVELOPMENT seeds.
type wf14Row struct{
 Seed int `json:"seed"`
 Family int `json:"family"`
 Depth int `json:"depth"`
 RoleShift int `json:"role_shift"`
 Truth [4]bool `json:"truth"`
 Predictions map[string][4]bool `json:"predictions"`
 StateBytes int `json:"state_bytes"`
}
func wf14FrozenRoles()wf13Roles{
 return wf13Roles{Bias:4.997920317446051,V:[4][4]float64{
 {-1.2516429336771078,-1.2571723585102805,-1.3702159581837117,1.3591861342639062},
 {-1.4097662197605667,-1.1859954550309044,1.5019227713467245,-1.3638159189834762},
 {1.2630898461119602,1.2728972755052257,1.3737904685355704,1.3650768322993967},
 {1.4058446864910592,1.2194933772476702,-1.4878312823167705,-1.3704207571068523}}}
}
func wf14Serialize(world wf13World)(wf1State,error){
 var state wf1State
 if unsafe.Sizeof(state)!=576{return state,fmt.Errorf("wrong byte budget")}
 for i,e:=range world.Facts.Edges{
  if e.From<0||e.From>65535||e.To<0||e.To>65535{return state,fmt.Errorf("ID too wide")}
  if world.Roles[i]<0||world.Roles[i]>3{return state,fmt.Errorf("invalid role")}
  state[i]=math.Float32frombits(uint32(e.From)<<16|uint32(e.To))
  state[12+i]=float32(world.Roles[i])
 }
 for i:=24;i<144;i++{if state[i]!=0{return state,fmt.Errorf("nonzero reserve")}}
 return state,nil
}
func wf14Predict(state wf1State,query [2]int,roleSeq [8]int,k int,changed *int,
 roleEmbedding wf13Roles,entity wf7Encoder,rec wf5Composer,blind bool)float64{
 var source,dest [12]int;var edgeRole [12]int
 for i:=0;i<12;i++{
  x:=math.Float32bits(state[i])
  source[i],dest[i]=int(x>>16),int(x&65535)
  edgeRole[i]=int(state[12+i])
 }
 var active [12]float64
 for i:=0;i<12;i++{
  sim,_,_:=wf7Pair(entity,query[0],source[i])
  r:=edgeRole[i]
  if changed!=nil&&i==*changed{r=(r+1)%4}
  gate:=1.
  if !blind{gate=wf13RoleScore(roleEmbedding,r,roleSeq[0])}
  active[i]=sim*gate
 }
 for h:=1;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   message:=0.
   for i:=0;i<12;i++{if i==j{continue}
    a,_,_:=wf7Pair(entity,dest[i],source[j]);message+=active[i]*a
   }
   r:=edgeRole[j]
   if changed!=nil&&j==*changed{r=(r+1)%4}
   gate:=1.
   if !blind{gate=wf13RoleScore(roleEmbedding,r,roleSeq[h])}
   next[j]=wf1Prob(rec.W[0]*message+rec.W[1]*active[j]+rec.W[2])*gate
  }
  active=next
 }
 score:=0.
 for i:=0;i<12;i++{p,_,_:=wf7Pair(entity,dest[i],query[1]);score+=active[i]*p}
 return wf1Prob(rec.W[3]*score+rec.W[4])
}
func RunWFoundationF14()(map[string]any,error){
 start:=time.Now()
 embedding,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 learned:=wf14FrozenRoles()
 rec:=wf5Composer{W:[5]float64{
  6.470679800286544,-1.3847365385924115,-3.0464843010034035,
  7.51083733991036,-3.4005406394991717}}
 names:=[]string{"frozen_serialized_candidate","raw_world_F13_debug",
   "untrained_serialized_roles","roleblind_serialized","coded_exact_typed_upper"}
 rows:=[]wf14Row{}
 for seed:=57900;seed<57912;seed++{for family:=64;family<=67;family++{
  for depth:=7;depth<=8;depth++{for shift:=0;shift<2;shift++{
   world,err:=wf13Build(seed,depth,family,shift);if err!=nil{return nil,err}
   state,err:=wf14Serialize(world);if err!=nil{return nil,err}
   q:=[2]int{world.Src,world.Dst}
   idx:=-1
   for j,e:=range world.Facts.Edges{
    if e.From==wf1Rename(depth/2,seed,2)&&e.To==wf1Rename(depth/2+1,seed,2){idx=j;break}
   }
   if idx<0{return nil,fmt.Errorf("original bridge disappeared")}
   perm:=world.Required
   if perm[0]==perm[1]{return nil,fmt.Errorf("cannot permute identical roles")}
   perm[0],perm[1]=perm[1],perm[0]
   row:=wf14Row{Seed:seed,Family:family,Depth:depth,RoleShift:shift,
    Truth:[4]bool{true,false,false,false},Predictions:map[string][4]bool{},
    StateBytes:int(unsafe.Sizeof(state))}
   var truth [4]bool
   for i:=0;i<4;i++{
    r,k:=world.Required,depth;var mutation *int
    if i==1{mutation=&idx};if i==2{r=perm};if i==3{k=depth-1}
    truth[i]=wf13Exact(world,r,k,mutation)
   }
   if truth!=row.Truth{return nil,fmt.Errorf("coded typed reference disagrees")}
   for _,name:=range names{
    var pred [4]bool
    for i:=0;i<4;i++{
     r,k:=world.Required,depth;var mutation *int
     if i==1{mutation=&idx};if i==2{r=perm};if i==3{k=depth-1}
     switch name{
     case "frozen_serialized_candidate":
      pred[i]=wf14Predict(state,q,r,k,mutation,learned,embedding,rec,false)>=.5
     case "raw_world_F13_debug":
      pred[i]=wf13Answer(world,r,k,learned,embedding,rec,mutation,false)>=.5
     case "untrained_serialized_roles":
      pred[i]=wf14Predict(state,q,r,k,mutation,wf13Init(),embedding,rec,false)>=.5
     case "roleblind_serialized":
      pred[i]=wf14Predict(state,q,r,k,mutation,wf13Init(),embedding,rec,true)>=.5
     case "coded_exact_typed_upper":pred[i]=truth[i]
     }
    }
    row.Predictions[name]=pred
   }
   if row.Predictions["frozen_serialized_candidate"]!=row.Predictions["raw_world_F13_debug"]{
    return nil,fmt.Errorf("serialized computation diverged from original")}
   rows=append(rows,row)
  }}
 }}
 if len(rows)!=192{return nil,fmt.Errorf("F14 graph count drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  correct,pairs,fp:=0,0,0
  for _,r:=range rows{
   p:=r.Predictions[name]
   if p==r.Truth{pairs++}
   for j,t:=range r.Truth{
    if p[j]==t{correct++}
    if !t&&p[j]{fp++}
   }
  }
  scores[name]=map[string]float64{"correct":float64(correct),
   "accuracy":float64(correct)/768,"fourway_correct":float64(pairs),
   "fourway_accuracy":float64(pairs)/192,
   "false_positives":float64(fp),"negative_fp_rate":float64(fp)/576}
 }
 verdict:="INVALID_SERIALIZED_REFERENCE"
 if scores["coded_exact_typed_upper"]["fourway_correct"]==192{
  a:=scores["frozen_serialized_candidate"];baseline:=scores["untrained_serialized_roles"]
  if a["fourway_accuracy"]>=.80&&a["accuracy"]>=.90&&a["negative_fp_rate"]<=.10&&
     a["fourway_accuracy"]-baseline["fourway_accuracy"]>=.15{
   verdict="DEVELOPMENT_FROZEN_SERIALIZED_ROLE_DEPTH_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_FROZEN_SERIALIZED_ROLE_DEPTH_TRANSFER"}
 }
 return map[string]any{
  "schema":"wingless.f14.actual-serialized-state-dev.v1",
  "stage":"DEVELOPMENT_ONLY","original_F13_sha":"0b9e769b00dbaf57fe1357e17592ad41eed47b29",
  "original_F13_json_sha256":"29af0774cf719b1d103659ac1caaa3d6334718b43ad190888c325af179f1ea1b",
  "frozen_role_vectors":learned.V,"frozen_role_bias":learned.Bias,
  "new_seed_bounds":[]int{57900,57912},"grammar_families_reused":[]int{64,65,66,67},
  "quadruple_graphs":len(rows),"questions":768,"negative_questions":576,
  "resident_state_bytes":576,"packed_endpoint_bytes":48,"packed_role_bytes":48,
  "zero_reserved_bytes":480,"max_role_query_bytes":8,"query_depth_bytes":1,
  "serialized_state_used_by_candidate":true,"model_retrained":false,
  "programmed_depth_loops":true,"learned_halting":false,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_model_calls":0,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "wall_ms":time.Since(start).Milliseconds()},nil
}
