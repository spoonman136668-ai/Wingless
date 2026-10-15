package unitary
import("fmt";"math";"time")
// W-F21: learn authorization from imperfect contextual testimony BEFORE new
// native alias relation test. The structural merge is still CODED; only the
// decision to trust supplied in-state alias evidence is LEARNED. It does
// not yet qualify as autonomous learned identity composition.
type wf21Row struct{
 Seed,Family,Depth,Shift,Variant int
 Truth bool
 Evidence [3]float64
 Learned wf15Choice
 RawFrozen wf15Choice
 Witness0 wf15Choice
 Majority wf15Choice
 OracleCanonical wf15Choice
}
type wf21Weights struct{V [4]float64}
var wf21Patterns=[8][3]float64{
 {-1,-1,-1},{-1,-1,1},{-1,1,-1},{-1,1,1},
 {1,-1,-1},{1,-1,1},{1,1,-1},{1,1,1}}
func wf21Sig(z float64)float64{
 if z<0{e:=math.Exp(z);return e/(1+e)}
 return 1/(1+math.Exp(-z))
}
func wf21EvidenceTruth(x [3]float64)bool{return x[0]+x[1]-x[2]>0}
// Only the TRAIN-only historical synthetics contribute gradients. The
// validation source is not opened until this single policy has frozen.
func wf21Train()wf21Weights{
 var model wf21Weights
 for epoch:=0;epoch<500;epoch++{
  var gradient [4]float64
  for seed:=0;seed<32;seed++{
   for _,x:=range wf21Patterns{
    target:=0.;if wf21EvidenceTruth(x){target=1}
    score:=model.V[0]+model.V[1]*x[0]+model.V[2]*x[1]+model.V[3]*x[2]
    delta:=wf21Sig(score)-target
    gradient[0]+=delta;for i:=0;i<3;i++{gradient[i+1]+=delta*x[i]}
   }
  }
  for i:=0;i<4;i++{model.V[i]-=.15*(gradient[i]/256+.0005*model.V[i])}
 }
 return model
}
// The input to this bridge consists of its 576B native state, fixed learned
// policy weights and the typed query. No graph record/truth is consulted.
func wf21Online(state wf1State,q [2]int,seq [8]int,k int,
 roles wf13Roles,entity wf7Encoder,rec wf5Composer,risk wf16Weights,
 learned wf21Weights,kind string)wf15Choice{
 bits:=math.Float32bits(state[24])
 alias,original:=int(bits>>16),int(bits&65535)
 x:=[3]float64{float64(state[25]),float64(state[26]),float64(state[27])}
 allow:=false
 switch kind{
 case "learned":allow=wf21Sig(learned.V[0]+learned.V[1]*x[0]+learned.V[2]*x[1]+learned.V[3]*x[2])>=.5
 case "first_witness":allow=x[0]>0
 case "unweighted_majority":allow=x[0]+x[1]+x[2]>0
 case "raw":allow=false
 case "coded_upper":allow=wf21EvidenceTruth(x)
 }
 packed:=state
 if allow&&alias!=0&&alias!=original{
  for j:=0;j<12;j++{
   value:=math.Float32bits(packed[j])
   src,dst:=int(value>>16),int(value&65535)
   if src==alias{packed[j]=math.Float32frombits(uint32(original)<<16|uint32(dst))}
  }
 }
 return wf17Online(packed,q,seq,k,roles,entity,rec,risk,"full")
}
func RunWFoundationF21()(map[string]any,error){
 began:=time.Now()
 learned:=wf21Train()
 allPatterns:=0
 for _,x:=range wf21Patterns{
  prediction:=wf21Sig(learned.V[0]+learned.V[1]*x[0]+learned.V[2]*x[1]+learned.V[3]*x[2])>=.5
  if prediction==wf21EvidenceTruth(x){allPatterns++}
 }
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 risk:=wf16Weights{W:[8]float64{.5233848397033497,.111980090092255,-4.3433500918591745,-4.043757386252545,.5349164073718833,1.7913415291476846,.5233848395808397,.38694717564180703}}
 rows:=[]wf21Row{}
 for seed:=80800;seed<80806;seed++{for fam:=108;fam<=111;fam++{
  for k:=7;k<=8;k++{for shift:=0;shift<2;shift++{
   w,main,alt,noise,err:=wf19Graph(seed,k,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108
   old:=wf1Rename(k+35,seed,2)
   alias:=wf1Rename(k+90,seed,2)
   unrelated:=wf1Rename(k+91,seed,2)
   second:=-1
   for j,e:=range w.Facts.Edges{
    if e.From==old&&e.To==wf1Rename(fork+2,seed,2){second=j}
   }
   if second<0||old==alias||alias==unrelated{return nil,fmt.Errorf("missing heldout bridge")}
   q:=[2]int{w.Src,w.Dst}
   for variant:=0;variant<12;variant++{
    raw:=w
    evidence:=wf21Patterns[variant%8]
    if variant<8{
     raw.Roles[main]=(raw.Roles[main]+1)%4
     raw.Facts.Edges[second].From=alias
    }else{
     evidence=[3]float64{1,1,-1}
     switch variant{
     case 8: // untouched both path witnesses: true
     case 9:raw.Roles[main]=(raw.Roles[main]+1)%4 // second witness true
     case 10:raw.Roles[main]=(raw.Roles[main]+1)%4;raw.Roles[alt]=(raw.Roles[alt]+1)%4
     case 11:raw.Roles[main]=(raw.Roles[main]+1)%4;raw.Facts.Edges[second].From=unrelated
     }
    }
    truth:=true
    if variant<8{truth=wf21EvidenceTruth(evidence)}
    if variant==10||variant==11{truth=false}
    validWorld:=raw
    if variant<8&&truth{validWorld.Facts.Edges[second].From=old}
    oracle:=wf13Exact(validWorld,w.Required,k,nil)
    if oracle!=truth{return nil,fmt.Errorf("INVALID_F21_SOURCE_ORACLE seed %d family %d variant %d",seed,fam,variant)}
    state,err:=wf14Serialize(raw);if err!=nil{return nil,err}
    state[24]=math.Float32frombits(uint32(alias)<<16|uint32(old))
    for j,v:=range evidence{state[25+j]=float32(v)}
    row:=wf21Row{Seed:seed,Family:fam,Depth:k,Shift:shift,Variant:variant,
      Truth:truth,Evidence:evidence,
      Learned:wf21Online(state,q,w.Required,k,role,enc,core,risk,learned,"learned"),
      RawFrozen:wf21Online(state,q,w.Required,k,role,enc,core,risk,learned,"raw"),
      Witness0:wf21Online(state,q,w.Required,k,role,enc,core,risk,learned,"first_witness"),
      Majority:wf21Online(state,q,w.Required,k,role,enc,core,risk,learned,"unweighted_majority"),
      OracleCanonical:wf21Online(state,q,w.Required,k,role,enc,core,risk,learned,"coded_upper")}
    rows=append(rows,row)
    _=noise
   }
  }}
 }}
 if len(rows)!=1152{return nil,fmt.Errorf("F21 source count drift")}
 scores:=map[string]map[string]float64{}
 for _,mode:=range []string{"learned","raw","first_witness","majority","coded_upper"}{
  correct,fp,fn:=0,0,0
  for _,r:=range rows{
   var p wf15Choice
   switch mode{
    case "learned":p=r.Learned
    case "raw":p=r.RawFrozen
    case "first_witness":p=r.Witness0
    case "majority":p=r.Majority
    default:p=r.OracleCanonical
   }
   if p.Answer==r.Truth{correct++}
   if p.Answer&&!r.Truth{fp++}
   if !p.Answer&&r.Truth{fn++}
  }
  scores[mode]=map[string]float64{"correct":float64(correct),
    "accuracy":float64(correct)/float64(len(rows)),
    "false_positives":float64(fp),"false_negatives":float64(fn)}
 }
 a:=scores["learned"];ctl:=scores["first_witness"];coded:=scores["coded_upper"]
 verdict:="INVALID_F21_CODED_UPPER_OR_SOURCE"
 if coded["correct"]==1152{
  success:=allPatterns==8&&a["accuracy"]>=.95&&a["false_negatives"]<=29&&
    a["false_positives"]<=29&&a["accuracy"]-ctl["accuracy"]>=.05&&
    a["accuracy"]-scores["majority"]["accuracy"]>=.05
  if success{verdict="DEVELOPMENT_LEARNED_ALIAS_AUTHORIZATION_TRANSFER"
  }else{verdict="DEVELOPMENT_NO_LEARNED_ALIAS_AUTHORIZATION_TRANSFER"}
 }
 return map[string]any{
  "schema":"wingless.f21.learned-evidence-alias-authorization.v1",
  "stage":"DEVELOPMENT_ONLY","frozen_F20_source":"7298eab8cdc8c50a626a58cb9bb111c94a89c9db",
  "training_only_evidence_patterns":8,"training_only_examples":256,
  "policy_weights":learned.V,"policy_threshold":.5,
  "evidence_policy_correct_patterns_train":allPatterns,
  "frozen_entity_role_core":true,"alias_structural_canonicalization_programmed":true,
  "learned_component":"evidence-conditioned authorization decision only",
  "new_seeds":[]int{80800,80806},"graph_families_reused":[]int{108,109,110,111},
  "graphs":96,"questions":1152,"native_state_bytes":576,
  "independent_source_birth":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_paid_api_calls":0,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "elapsed_ms":time.Since(began).Milliseconds()},nil
}
