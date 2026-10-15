package unitary

import("fmt";"math";"sort";"time")

// W-F23 prospective DEVELOPMENT: infer whether two graph nodes denote
// the same entity using a four-weight metric TRAINED on 3-feature known
// node descriptions. Unlike W-F22, there is NO pair certificate or
// source-provided alias pair. Native 576B state contains 12 edges +12
// roles +24 independent (ID,descriptor[3]) records. Matching descriptors
// are supplied by an authored simulator; NOT semantically discovered.
type wf23Metric struct{ W [4]float64 }
type wf23Row struct{
 Seed,Family,Depth,Shift,Variant int
 Truth bool
 Learned,Initialized,ExactID,Unrestricted,CodedDescription wf15Choice
}
type wf23Record struct{ ID int; Attr [3]float64 }
var wf23TrainCodes=[4][3]float64{{0,0,0},{0,1,1},{1,0,1},{1,1,0}}
var wf23TestCodes=[4][3]float64{{0,0,1},{0,1,0},{1,0,0},{1,1,1}}
func wf23MetricProb(model wf23Metric,a,b [3]float64)float64{
 z:=model.W[0]
 for i:=0;i<3;i++{d:=a[i]-b[i];z+=model.W[i+1]*d*d}
 return wf21Sig(z)
}
func wf23Train()wf23Metric{
 var m wf23Metric
 for epoch:=0;epoch<450;epoch++{
  var g [4]float64
  for _,a:=range wf23TrainCodes{for _,b:=range wf23TrainCodes{
   y:=0.;if a==b{y=1}
   p:=wf23MetricProb(m,a,b)
   d:=p-y
   g[0]+=d
   for j:=0;j<3;j++{dist:=a[j]-b[j];g[j+1]+=d*dist*dist}
  }}
  for i:=0;i<4;i++{m.W[i]-=.20*(g[i]/16+.0005*m.W[i])}
 }
 return m
}
func wf23WorldRecords(w wf13World,seed int,variant int,old,alias int)([]wf23Record,error){
 uniq:=map[int]bool{}
 for _,e:=range w.Facts.Edges{uniq[e.From]=true;uniq[e.To]=true}
 ids:=[]int{}
 for id:=range uniq{ids=append(ids,id)}
 sort.Ints(ids)
 if len(ids)>24{return nil,fmt.Errorf("native node descriptor record cap exceeded %d",len(ids))}
 code:=wf23TestCodes[variant%4]
 alt:=wf23TestCodes[(variant%4+1)%4]
 records:=make([]wf23Record,0,len(ids))
 for i,id:=range ids{
  // Non-participating entities get well-separated continuous descriptors.
  // Descriptions represent authored context, NOT discovered ontology.
  v:=float64(i+12)
  attr:=[3]float64{v*7+2,v*11+4,v*13+6}
  if id==old{attr=code}
  if id==alias{
   if variant<4{attr=code}else{attr=alt}
  }
  records=append(records,wf23Record{ID:id,Attr:attr})
 }
 return records,nil
}
func wf23Encode(world wf13World,records []wf23Record)(wf1State,error){
 state,err:=wf14Serialize(world);if err!=nil{return state,err}
 if len(records)>24{return state,fmt.Errorf("descriptor budget exceeded")}
 for i,r:=range records{
  if r.ID<0||r.ID>65535{return state,fmt.Errorf("source ID overflow")}
  offset:=24+i*4
  state[offset]=float32(r.ID)
  for j:=0;j<3;j++{state[offset+j+1]=float32(r.Attr[j])}
 }
 return state,nil
}
func wf23Native(state wf1State,q [2]int,pattern [8]int,k int,
 role wf13Roles,entity wf7Encoder,rec wf5Composer,m wf23Metric,mode string)wf15Choice{
 var src,dst,tags [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(state[i])
  src[i],dst[i]=int(bits>>16),int(bits&65535)
  tags[i]=int(state[12+i])
 }
 attributes:=map[int][3]float64{}
 for i:=0;i<24;i++{
  at:=24+i*4
  id:=int(state[at])
  if id==0{continue}
  attributes[id]=[3]float64{float64(state[at+1]),float64(state[at+2]),float64(state[at+3])}
 }
 var active [12]float64
 for h:=0;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   value:=0.
   if h==0{value,_,_=wf7Pair(entity,q[0],src[j])
   }else{
    msg:=0.
    for i:=0;i<12;i++{
     if i==j{continue}
     sim,_,_:=wf7Pair(entity,dst[i],src[j])
     if dst[i]!=src[j]{
      a,okA:=attributes[dst[i]]
      b,okB:=attributes[src[j]]
      if okA&&okB{
       score:=0.
       switch mode{
       case "learned":score=wf23MetricProb(m,a,b)
       case "untrained":score=wf23MetricProb(wf23Metric{},a,b)
       case "coded_description":if a==b{score=1}
       case "unrestricted":score=1
       }
       if score>=.5{sim=math.Max(sim,score)}
      }
     }
     msg+=active[i]*sim
    }
    value=wf1Prob(rec.W[0]*msg+rec.W[1]*active[j]+rec.W[2])
   }
   next[j]=value*wf13RoleScore(role,tags[j],pattern[h])
  }
  active=next
 }
 endpoint:=0.
 for i:=0;i<12;i++{v,_,_:=wf7Pair(entity,dst[i],q[1]);endpoint+=v*active[i]}
 p:=wf1Prob(rec.W[3]*endpoint+rec.W[4])
 return wf15Choice{Answer:p>=.5,Steps:k,EarlyHalt:false}
}
func RunWFoundationF23()(map[string]any,error){
 started:=time.Now()
 metric:=wf23Train()
 trainedPairs:=0
 for _,a:=range wf23TrainCodes{for _,b:=range wf23TrainCodes{
  prediction:=wf23MetricProb(metric,a,b)>=.5
  if prediction==(a==b){trainedPairs++}
 }}
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 rec:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 rows:=[]wf23Row{}
 for seed:=88300;seed<88306;seed++{for fam:=108;fam<=111;fam++{
  for depth:=7;depth<=8;depth++{for shift:=0;shift<2;shift++{
   orig,main,alt,noise,err:=wf19Graph(seed,depth,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108
   old,alias:=wf1Rename(depth+35,seed,2),wf1Rename(depth+90,seed,2)
   unrelated:=wf1Rename(depth+91,seed,2)
   idx:=-1
   for j,e:=range orig.Facts.Edges{
    if e.From==old&&e.To==wf1Rename(fork+2,seed,2){idx=j}
   }
   if idx<0{return nil,fmt.Errorf("W23 alternate chain edge absent")}
   q:=[2]int{orig.Src,orig.Dst}
   for variant:=0;variant<12;variant++{
    w:=orig
    truth:=true
    if variant<8{
     w.Roles[main]=(w.Roles[main]+1)%4
     w.Facts.Edges[idx].From=alias
     truth=variant<4
    }else{
     switch variant{
      case 8: // intact, positive
      case 9:w.Roles[main]=(w.Roles[main]+1)%4 // alt surviving, positive
      case 10:w.Roles[main]=(w.Roles[main]+1)%4;w.Roles[alt]=(w.Roles[alt]+1)%4;truth=false
      case 11:w.Roles[main]=(w.Roles[main]+1)%4;w.Facts.Edges[idx].From=unrelated;truth=false
     }
    }
    canonical:=w
    if variant<4{canonical.Facts.Edges[idx].From=old}
    if wf13Exact(canonical,orig.Required,depth,nil)!=truth{
     return nil,fmt.Errorf("INVALID_W23_TYPED_SOURCE seed=%d family=%d variant=%d",seed,fam,variant)
    }
    records,err:=wf23WorldRecords(w,seed,variant,old,alias)
    if err!=nil{return nil,err}
    state,err:=wf23Encode(w,records);if err!=nil{return nil,err}
    rows=append(rows,wf23Row{Seed:seed,Family:fam,Depth:depth,Shift:shift,Variant:variant,Truth:truth,
     Learned:wf23Native(state,q,orig.Required,depth,role,entity,rec,metric,"learned"),
     Initialized:wf23Native(state,q,orig.Required,depth,role,entity,rec,metric,"untrained"),
     ExactID:wf23Native(state,q,orig.Required,depth,role,entity,rec,metric,"exact"),
     Unrestricted:wf23Native(state,q,orig.Required,depth,role,entity,rec,metric,"unrestricted"),
     CodedDescription:wf23Native(state,q,orig.Required,depth,role,entity,rec,metric,"coded_description")})
    _=noise
   }
  }}
 }}
 if len(rows)!=1152{return nil,fmt.Errorf("science size invalid")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"learned","untrained","raw_id","unrestricted","coded_description"}{
  corr,fp,fn:=0,0,0
  for _,r:=range rows{
   var v wf15Choice
   switch name{
   case "learned":v=r.Learned
   case "untrained":v=r.Initialized
   case "raw_id":v=r.ExactID
   case "unrestricted":v=r.Unrestricted
   case "coded_description":v=r.CodedDescription
   }
   if v.Answer==r.Truth{corr++}
   if v.Answer&&!r.Truth{fp++}
   if !v.Answer&&r.Truth{fn++}
  }
  scores[name]=map[string]float64{"correct":float64(corr),"accuracy":float64(corr)/1152,
    "false_positives":float64(fp),"false_negatives":float64(fn)}
 }
 a:=scores["learned"];coded:=scores["coded_description"]
 verdict:="INVALID_W23_ENTITY_DESCRIPTION_ORACLE"
 if coded["correct"]==1152 && trainedPairs==16{
  if a["accuracy"]>=.95&&a["false_positives"]<=29&&a["false_negatives"]<=29&&
     a["accuracy"]-scores["raw_id"]["accuracy"]>=.05&&
     a["accuracy"]-scores["untrained"]["accuracy"]>=.05{
   verdict="DEVELOPMENT_LEARNED_NODE_DESCRIPTION_ALIAS_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_LEARNED_NODE_DESCRIPTION_ALIAS_TRANSFER"}
 }
 return map[string]any{"schema":"wingless.f23.learnt-per-node-evidence-alias.dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_F22":"5a785822bf8fe7aa9878ab374dfd799a9d6f3633",
  "training_descriptions":wf23TrainCodes,"unseen_test_descriptions":wf23TestCodes,
  "train_only_metric_parameters":metric.W,"train_pairs_correct":trainedPairs,
  "source_description_features_given":true,"pair_certificate_supplied":false,
  "programmed_lookup_of_node_annotations":true,"programmed_squared_difference_features":true,
  "alias_candidate_inferred_from_learned_node_metric":true,
  "prior_entity_role_recurrent_weights_frozen":true,"new_graph_seeds":[]int{88300,88306},
  "graph_families_reused":[]int{108,109,110,111},"graph_count":96,"questions":1152,
  "actual_native_state_bytes":576,"per_node_descriptor_record_limit":24,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,"elapsed_ms":time.Since(started).Milliseconds()},nil
}
