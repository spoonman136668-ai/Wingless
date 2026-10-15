package unitary
import("fmt";"math";"time")

// W-F22 integrates a learned evidence-conditioned ALIAS LINK directly into
// the native graph recurrent transition, rather than rewriting input edges.
// The predecessor's 4 learned weights and 576B state are frozen. Coded
// candidate-pair matching and initial evidence-channel semantics remain
// architectural assumptions; this is NOT unrestricted native identity learning.
type wf22Row struct{
 Seed,Family,Depth,Shift,Variant int
 Truth bool
 Evidence [3]float64
 Native,NoAlias,AlwaysAlias,CodedUpper,RewriteUpper wf15Choice
}
func wf22Native(state wf1State,q [2]int,pattern [8]int,k int,
 rm wf13Roles,enc wf7Encoder,core wf5Composer,p wf21Weights,mode string)wf15Choice{
 var src,dst,role [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(state[i])
  src[i],dst[i]=int(bits>>16),int(bits&65535)
  role[i]=int(state[12+i])
 }
 aliasBits:=math.Float32bits(state[24])
 alias,old:=int(aliasBits>>16),int(aliasBits&65535)
 x:=[3]float64{float64(state[25]),float64(state[26]),float64(state[27])}
 confidence:=wf21Sig(p.V[0]+p.V[1]*x[0]+p.V[2]*x[1]+p.V[3]*x[2])
 allow:=false
 switch mode{
 case "native":allow=confidence>=.5
 case "coded_upper":allow=wf21EvidenceTruth(x)
 case "always":allow=true
 }
 var active [12]float64
 for h:=0;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   value:=0.
   if h==0{
    value,_,_=wf7Pair(enc,q[0],src[j])
   }else{
    msg:=0.
    for i:=0;i<12;i++{
     if i==j{continue}
     sim,_,_:=wf7Pair(enc,dst[i],src[j])
     // No canonicalized edge/node values are created. This relation
     // changes the recurrent message coefficient itself. The pair
     // identity check is still CODED; only trust is learned.
     if allow && old!=alias && dst[i]==old && src[j]==alias{
      sim=math.Max(sim,1.)
     }
     msg+=active[i]*sim
    }
    value=wf1Prob(core.W[0]*msg+core.W[1]*active[j]+core.W[2])
   }
   next[j]=value*wf13RoleScore(rm,role[j],pattern[h])
  }
  active=next
 }
 terminal:=0.
 for i:=0;i<12;i++{
  sim,_,_:=wf7Pair(enc,dst[i],q[1])
  terminal+=active[i]*sim
 }
 prob:=wf1Prob(core.W[3]*terminal+core.W[4])
 return wf15Choice{Answer:prob>=.5,Steps:k,EarlyHalt:false}
}
func RunWFoundationF22()(map[string]any,error){
 started:=time.Now()
 gate:=wf21Train()
 original:=[4]float64{1.1898503414192076e-16,2.8522704808696697,2.85227048086967,-2.85227048086967}
 for i,v:=range gate.V{if math.Abs(v-original[i])>1e-8{
  return nil,fmt.Errorf("W21 recovered head drift at %d",i)
 }}
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 risk:=wf16Weights{W:[8]float64{.5233848397033497,.111980090092255,-4.3433500918591745,-4.043757386252545,.5349164073718833,1.7913415291476846,.5233848395808397,.38694717564180703}}
 rows:=[]wf22Row{}
 expectedExtra:=[]bool{true,false,false,true}
 for seed:=83900;seed<83906;seed++{for fam:=108;fam<=111;fam++{
  for depth:=7;depth<=8;depth++{for shift:=0;shift<2;shift++{
   w,main,alt,noise,err:=wf19Graph(seed,depth,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108
   old:=wf1Rename(depth+35,seed,2)
   alias:=wf1Rename(depth+90,seed,2)
   unrelated:=wf1Rename(depth+91,seed,2)
   second:=-1
   for i,e:=range w.Facts.Edges{
    if e.From==old&&e.To==wf1Rename(fork+2,seed,2){second=i}
   }
   if second<0||old==alias||alias==unrelated{return nil,fmt.Errorf("alias fixture conflict")}
   q:=[2]int{w.Src,w.Dst}
   for variant:=0;variant<12;variant++{
    raw:=w
    x:=[3]float64{1,1,-1}
    truth:=true
    if variant<8{
     x=wf21Patterns[variant]
     raw.Roles[main]=(raw.Roles[main]+1)%4
     raw.Facts.Edges[second].From=alias
     truth=wf21EvidenceTruth(x)
    }else{
     truth=expectedExtra[variant-8]
     switch variant{
     case 9:raw.Roles[main]=(raw.Roles[main]+1)%4;raw.Facts.Edges[second].From=unrelated
     case 10:raw.Roles[main]=(raw.Roles[main]+1)%4;raw.Roles[alt]=(raw.Roles[alt]+1)%4
     case 11:raw.Roles[noise]=(raw.Roles[noise]+1)%4
     }
    }
    verify:=raw
    if variant<8&&truth{verify.Facts.Edges[second].From=old}
    if wf13Exact(verify,w.Required,depth,nil)!=truth{
     return nil,fmt.Errorf("INVALID_W22_TYPED_REFERENCE seed=%d fam=%d variant=%d",seed,fam,variant)
    }
    state,err:=wf14Serialize(raw);if err!=nil{return nil,err}
    state[24]=math.Float32frombits(uint32(alias)<<16|uint32(old))
    for i,v:=range x{state[25+i]=float32(v)}
    rows=append(rows,wf22Row{Seed:seed,Family:fam,Depth:depth,Shift:shift,
     Variant:variant,Truth:truth,Evidence:x,
     Native:wf22Native(state,q,w.Required,depth,role,enc,core,gate,"native"),
     NoAlias:wf22Native(state,q,w.Required,depth,role,enc,core,gate,"none"),
     AlwaysAlias:wf22Native(state,q,w.Required,depth,role,enc,core,gate,"always"),
     CodedUpper:wf22Native(state,q,w.Required,depth,role,enc,core,gate,"coded_upper"),
     RewriteUpper:wf21Online(state,q,w.Required,depth,role,enc,core,risk,gate,"coded_upper")})
   }
  }}
 }}
 if len(rows)!=1152{return nil,fmt.Errorf("W22 question count drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"native","no_alias","always_alias","coded_upper","rewrite_upper"}{
  correct,fp,fn:=0,0,0
  by:=map[int][2]int{}
  for _,r:=range rows{
   var v wf15Choice
   switch name{
    case "native":v=r.Native
    case "no_alias":v=r.NoAlias
    case "always_alias":v=r.AlwaysAlias
    case "coded_upper":v=r.CodedUpper
    default:v=r.RewriteUpper
   }
   if v.Answer==r.Truth{correct++}
   if v.Answer&&!r.Truth{fp++}
   if !v.Answer&&r.Truth{fn++}
   s:=by[r.Family];if v.Answer!=r.Truth{s[0]++};s[1]++;by[r.Family]=s
  }
  worst:=0.
  for _,s:=range by{rate:=float64(s[0])/float64(s[1]);worst=math.Max(worst,rate)}
  scores[name]=map[string]float64{"correct":float64(correct),
    "accuracy":float64(correct)/1152,"false_positives":float64(fp),
    "false_negatives":float64(fn),"worst_family_error":worst}
 }
 a:=scores["native"];base:=scores["no_alias"];upper:=scores["rewrite_upper"]
 verdict:="INVALID_W22_CODED_ORACLE_UPPER"
 if upper["correct"]==1152{
  if a["accuracy"]>=.95&&a["worst_family_error"]<=.05&&
   a["accuracy"]-base["accuracy"]>=.05&&
   scores["coded_upper"]["accuracy"]>=.95{
   verdict="DEVELOPMENT_LEARNED_ALIAS_GRAPH_LINK_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_NATIVE_ALIAS_GRAPH_LINK_TRANSFER"}
 }
 return map[string]any{"schema":"wingless.f22.in-recurrence-learned-alias-link.dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_f21_sha":"4ddac012e9cbbb960fdc03a49fe86f4d08607e6b",
  "previous_evidence_head_weights":gate.V,
  "prior_frozen_model_retained":true,"new_head_training_performed":false,
  "resident_state_bytes":576,"coded_alias_endpoint_candidate_detection":true,
  "programmed_edge_endpoint_rewriting_inside_candidate":false,
  "learned_witness_controls_alias_link_during_recurrence":true,
  "structural_families_reused":[]int{108,109,110,111},"new_seeds":[]int{83900,83906},
  "graphs":96,"questions":len(rows),
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "elapsed_ms":time.Since(started).Milliseconds()},nil
}
