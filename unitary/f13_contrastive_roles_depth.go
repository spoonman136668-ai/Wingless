package unitary

import("fmt";"math";"time";"unsafe")

// W-F13: separately learned Euclidean 4D relation-role symbols via pairwise
// contrastive classification; no coded relation equality in inference.
// Compose them with the F11 frozen task-trained shared recurrence, conditioned
// on novel seven/eight-symbol ordered role queries. All new graph contexts
// DEVELOPMENT ONLY, no F3 independent original.
type wf13Roles struct{ V [4][4]float64; Bias float64 }
type wf13World struct{ Facts wf1Datum; Roles [12]int; Src,Dst int; Required [8]int; K int }
type wf13Row struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 Chain int `json:"chain"`
 RoleShift int `json:"role_shift"`
 Query [2]int `json:"query"`
 Truth [4]bool `json:"truth"`
 Predictions map[string][4]bool `json:"predictions"`
}
func wf13Init()wf13Roles{
 m:=wf13Roles{}
 for a:=0;a<4;a++{for k:=0;k<4;k++{
  m.V[a][k]=.19*math.Sin(float64((a+2)*(k+3)*11))
 }}
 return m
}
func wf13RoleScore(m wf13Roles,a,b int)float64{
 dist:=0.
 for i:=0;i<4;i++{d:=m.V[a][i]-m.V[b][i];dist+=d*d}
 return wf1Prob(m.Bias-dist)
}
func wf13Train()(wf13Roles,float64,int64){
 start:=time.Now()
 m:=wf13Init()
 var mom,variance [17]float64
 get:=func(m wf13Roles,k int)float64{
  if k==16{return m.Bias}
  return m.V[k/4][k%4]
 }
 set:=func(m *wf13Roles,k int,v float64){
  if k==16{m.Bias=v}else{m.V[k/4][k%4]=v}
 }
 loss:=func(m wf13Roles)float64{
  total:=0.
  for a:=0;a<4;a++{for b:=0;b<4;b++{
   p:=max(.0000001,min(.9999999,wf13RoleScore(m,a,b)))
   if a==b{total-=3*math.Log(p)}else{total-=math.Log(1-p)}
  }}
  return total/24
 }
 for epoch:=1;epoch<=210;epoch++{
  var grad [17]float64
  for j:=0;j<17;j++{
   plus,minus:=m,m
   set(&plus,j,get(m,j)+.0005)
   set(&minus,j,get(m,j)-.0005)
   grad[j]=(loss(plus)-loss(minus))/.001
  }
  for j:=0;j<17;j++{
   g:=max(-10,min(10,grad[j]))
   mom[j]=.9*mom[j]+.1*g
   variance[j]=.999*variance[j]+.001*g*g
   step:=.075*(mom[j]/(1-math.Pow(.9,float64(epoch))))/
     (math.Sqrt(variance[j]/(1-math.Pow(.999,float64(epoch))))+1e-8)
   set(&m,j,max(-15,min(15,get(m,j)-step)))
  }
 }
 return m,loss(m),time.Since(start).Milliseconds()
}
func wf13Build(seed,k,fam,shift int)(wf13World,error){
 var world wf13World
 if fam<64||fam>67||k<7||k>8{return world,fmt.Errorf("unknown F13 family/depth")}
 var edges [12][2]int
 for i:=0;i<k;i++{edges[i]=[2]int{i,i+1}}
 switch fam{
 case 64:edges[k]=[2]int{0,3}
 case 65:edges[k]=[2]int{2,5}
 case 66:edges[k]=[2]int{k-1,2}
 case 67:edges[k]=[2]int{k,1}
 }
 for i:=k+1;i<12;i++{u:=20+i-k;edges[i]=[2]int{u,u+1}}
 var roles [12]int
 for i:=0;i<12;i++{roles[i]=(i*3+seed+fam+shift)%4}
 for i:=11;i>0;i--{
  j:=(seed*29+k*11+fam*17+shift*7+i*23)%(i+1)
  edges[i],edges[j]=edges[j],edges[i]
  roles[i],roles[j]=roles[j],roles[i]
 }
 seen:=map[[2]int]bool{}
 for i,e:=range edges{
  a,b:=wf1Rename(e[0],seed,2),wf1Rename(e[1],seed,2)
  if seen[[2]int{a,b}]{return world,fmt.Errorf("duplicate ID binding")}
  seen[[2]int{a,b}]=true
  world.Facts.Edges[i]=nci02WEdge{From:a,To:b}
  world.Roles[i]=roles[i]
 }
 world.Src,world.Dst=wf1Rename(0,seed,2),wf1Rename(k,seed,2)
 world.K=k
 for i:=0;i<k;i++{
  a,b:=wf1Rename(i,seed,2),wf1Rename(i+1,seed,2)
  found:=false
  for j,e:=range world.Facts.Edges{
   if e.From==a&&e.To==b{world.Required[i]=world.Roles[j];found=true;break}
  }
  if !found{return world,fmt.Errorf("missing chain")}
 }
 return world,nil
}
func wf13Exact(world wf13World,roleCodes [8]int,k int,changed *int)bool{
 if k<1||k>8{return false}
 active:=map[int]bool{world.Src:true}
 for depth:=0;depth<k;depth++{
  next:=map[int]bool{}
  for i,e:=range world.Facts.Edges{
   typ:=world.Roles[i]
   if changed!=nil&&i==*changed{typ=(typ+1)%4}
   if active[e.From]&&typ==roleCodes[depth]{next[e.To]=true}
  }
  active=next
 }
 return active[world.Dst]
}
func wf13Answer(world wf13World,roleCodes [8]int,k int,roleModel wf13Roles,encoder wf7Encoder,rec wf5Composer,change *int,roleblind bool)float64{
 var active [12]float64
 for i,e:=range world.Facts.Edges{
  start,_,_:=wf7Pair(encoder,world.Src,e.From)
  typ:=world.Roles[i]
  if change!=nil&&i==*change{typ=(typ+1)%4}
  gate:=1.
  if !roleblind{gate=wf13RoleScore(roleModel,typ,roleCodes[0])}
  active[i]=start*gate
 }
 for hop:=1;hop<k;hop++{
  var next [12]float64
  for j,to:=range world.Facts.Edges{
   incoming:=0.
   for i,from:=range world.Facts.Edges{
    if i==j{continue}
    edge,_,_:=wf7Pair(encoder,from.To,to.From)
    incoming+=active[i]*edge
   }
   gate:=1.
   if !roleblind{
    code:=world.Roles[j]
    if change!=nil&&j==*change{code=(code+1)%4}
    gate=wf13RoleScore(roleModel,code,roleCodes[hop])
   }
   next[j]=wf1Prob(rec.W[0]*incoming+rec.W[1]*active[j]+rec.W[2])*gate
  }
  active=next
 }
 score:=0.
 for i,e:=range world.Facts.Edges{
  end,_,_:=wf7Pair(encoder,e.To,world.Dst)
  score+=active[i]*end
 }
 return wf1Prob(rec.W[3]*score+rec.W[4])
}
func RunWFoundationF13()(map[string]any,error){
 start:=time.Now()
 encoder,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 roleModel,loss,trainMs:=wf13Train()
 rec:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,
   7.51083733991036,-3.4005406394991717}}
 init:=wf13Init()
 names:=[]string{"contrastive_roles_frozen_depth_recurrence",
    "untrained_role_embedding_same_recurrence","role_blind_recurrence",
    "coded_exact_typed_role_depth","always_yes_shortcut"}
 rows:=[]wf13Row{}
 for seed:=56300;seed<56312;seed++{for family:=64;family<=67;family++{
  for depth:=7;depth<=8;depth++{for shift:=0;shift<2;shift++{
   world,err:=wf13Build(seed,depth,family,shift);if err!=nil{return nil,err}
   required:=world.Required
   changedIndex:=-1
   for i,e:=range world.Facts.Edges{
    if e.From==wf1Rename(depth/2,seed,2)&&e.To==wf1Rename(depth/2+1,seed,2){changedIndex=i;break}
   }
   if changedIndex<0{return nil,fmt.Errorf("required bridge absent")}
   permuted:=required
   a,b:=0,1
   if permuted[a]==permuted[b]{return nil,fmt.Errorf("identical roles for swap")}
   permuted[a],permuted[b]=permuted[b],permuted[a]
   row:=wf13Row{Seed:seed,Family:family,Chain:depth,RoleShift:shift,
     Query:[2]int{world.Src,world.Dst},Truth:[4]bool{true,false,false,false},
     Predictions:map[string][4]bool{}}
   oracle:=[4]bool{
    wf13Exact(world,required,depth,nil),
    wf13Exact(world,required,depth,&changedIndex),
    wf13Exact(world,permuted,depth,nil),
    wf13Exact(world,required,depth-1,nil)}
   if oracle!=row.Truth{return nil,fmt.Errorf("typed exact oracle case conflict")}
   for _,name:=range names{
    var pred [4]bool
    for i:=0;i<4;i++{
     p:=required;k:=depth;var mod *int
     if i==1{mod=&changedIndex}
     if i==2{p=permuted}
     if i==3{k=depth-1}
     switch name{
     case "contrastive_roles_frozen_depth_recurrence":
      pred[i]=wf13Answer(world,p,k,roleModel,encoder,rec,mod,false)>=.5
     case "untrained_role_embedding_same_recurrence":
      pred[i]=wf13Answer(world,p,k,init,encoder,rec,mod,false)>=.5
     case "role_blind_recurrence":
      pred[i]=wf13Answer(world,p,k,init,encoder,rec,mod,true)>=.5
     case "coded_exact_typed_role_depth":pred[i]=oracle[i]
     case "always_yes_shortcut":pred[i]=true
     }
    }
    row.Predictions[name]=pred
   }
   rows=append(rows,row)
  }}
 }}
 if len(rows)!=192||unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("F13 corpus/byte budget drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range names{
  correct,triple,fp:=0,0,0
  for _,r:=range rows{
   prediction:=r.Predictions[name]
   if prediction==r.Truth{triple++}
   for j,y:=range r.Truth{
    if prediction[j]==y{correct++}
    if !y&&prediction[j]{fp++}
   }
  }
  scores[name]=map[string]float64{"correct":float64(correct),"accuracy":float64(correct)/768,
   "all_four_correct":float64(triple),"all_four_rate":float64(triple)/192,
   "false_positives":float64(fp),"negative_fp_rate":float64(fp)/576}
 }
 status:="INVALID_TYPED_DEPTH_REFERENCE"
 if scores["coded_exact_typed_role_depth"]["all_four_correct"]==192{
  a:=scores["contrastive_roles_frozen_depth_recurrence"]
  baseline:=scores["untrained_role_embedding_same_recurrence"]
  if a["all_four_rate"]>=.80&&a["accuracy"]>=.90&&a["negative_fp_rate"]<=.10&&
      a["all_four_rate"]-baseline["all_four_rate"]>=.15{
   status="DEVELOPMENT_CONTRASTIVE_ROLE_DEPTH_COMPOSITION_TRANSFER"
  }else{status="DEVELOPMENT_NO_CONTRASTIVE_ROLE_DEPTH_TRANSFER"}
 }
 return map[string]any{"schema":"wingless.f13.contrastive-role-depth-composition-development.v1",
  "stage":"DEVELOPMENT_ONLY","parent_F12_sha":"6b8c3e45840d8606b35fa32af9b49c87382aa6b3",
  "role_training_only":"4-symbol positive/negative identity classification, no graph question training",
  "role_parameter_count":17,"pretrained_entity_projection_parameters":257,
  "frozen_recurrence_parameters":5,"role_optimization_steps":210,"role_training_ms":trainMs,
  "role_training_loss":loss,"role_embedding":roleModel.V,"role_bias":roleModel.Bias,
  "new_seed_bounds":[]int{56300,56312},"new_families":[]int{64,65,66,67},
  "heldout_depths":[]int{7,8},"quadruple_graphs":len(rows),
  "questions":768,"negative_questions":576,"native_state_bytes":576,
  "resident_role_bytes":48,"resident_entity_words_bytes":48,
  "query_role_bytes_max":8,"depth_query_bytes":1,
  "recurrence_iteration_programmed":true,"learned_halting":false,
  "candidate_coded_role_equality":false,"candidate_coded_entity_equality":false,
  "scores":scores,"rows":rows,"registered_verdict":status,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_model_calls":0,
  "wall_ms":time.Since(start).Milliseconds()},nil
}
