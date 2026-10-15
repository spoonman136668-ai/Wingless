package unitary
import("fmt";"math";"time")
// F20 tests whether frozen W19 reasoning consumes an explicit alias witness
// written inside its existing 576B resident state. The online model does
// not read reserve words 24/25: falsification of alias-binding integration,
// NOT evidence that aliases are impossible in an expanded architecture.
type wf20Row struct{Seed,Family,Depth,Variant int;Truth,RawOracle bool
 Full,Risk,Canonical wf15Choice;AliasDeclared bool}
func RunWFoundationF20()(map[string]any,error){
 start:=time.Now()
 enc,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 rm:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 risk:=wf16Weights{W:[8]float64{.5233848397033497,.111980090092255,-4.3433500918591745,-4.043757386252545,.5349164073718833,1.7913415291476846,.5233848395808397,.38694717564180703}}
 rows:=[]wf20Row{}
 for seed:=78600;seed<78606;seed++{for fam:=108;fam<=111;fam++{
  for k:=7;k<=8;k++{for shift:=0;shift<2;shift++{
   world,main,alt,noise,err:=wf19Graph(seed,k,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108
   old:=wf1Rename(k+35,seed,2)
   alias:=wf1Rename(k+90,seed,2)
   second:=-1
   for i,e:=range world.Facts.Edges{
    if e.From==old&&e.To==wf1Rename(fork+2,seed,2){second=i}
   }
   if second<0||old==alias{return nil,fmt.Errorf("missing identity bridge")}
   q:=[2]int{world.Src,world.Dst}
   // Each graph creates paired identical edge/role inputs differing
   // ONLY in whether alias metadata is declared true or false.
   for v:=0;v<6;v++{
    test:=world;aliasDeclared:=false
    if v==1||v==2||v==3{test.Roles[main]=(test.Roles[main]+1)%4}
    if v==2||v==3{
     test.Facts.Edges[second].From=alias
     aliasDeclared=v==2
    }
    if v==4{
     test.Roles[main]=(test.Roles[main]+1)%4
     test.Roles[alt]=(test.Roles[alt]+1)%4
    }
    if v==5{test.Roles[noise]=(test.Roles[noise]+1)%4}
    expected:=[]bool{true,true,true,false,false,true}[v]
    rawTruth:=wf13Exact(test,test.Required,k,nil)
    canon:=test
    if aliasDeclared{canon.Facts.Edges[second].From=old}
    declaredTruth:=wf13Exact(canon,canon.Required,k,nil)
    if declaredTruth!=expected{return nil,fmt.Errorf("INVALID_ALIAS_ORACLE seed=%d fam=%d v=%d",seed,fam,v)}
    packed,err:=wf14Serialize(test);if err!=nil{return nil,err}
    // Public alias metadata lives INSIDE (not outside) same 576B state.
    packed[24]=math.Float32frombits(uint32(alias)<<16|uint32(old))
    if aliasDeclared{packed[25]=1}
    original:=wf17Online(packed,q,test.Required,k,rm,enc,core,risk,"full")
    early:=wf17Online(packed,q,test.Required,k,rm,enc,core,risk,"risk")
    normalized,err:=wf14Serialize(canon);if err!=nil{return nil,err}
    fixed:=wf17Online(normalized,q,test.Required,k,rm,enc,core,risk,"full")
    rows=append(rows,wf20Row{Seed:seed,Family:fam,Depth:k,Variant:v,
     Truth:expected,RawOracle:rawTruth,AliasDeclared:aliasDeclared,
     Full:original,Risk:early,Canonical:fixed})
   }
  }}
 }}
 if len(rows)!=576{return nil,fmt.Errorf("F20 cohort mismatch")}
 // Fairness diagnosis: the original model is given alias metadata in its
 // state but cannot use it; paired input edges are bit-identical.
 indistinguishable:=0;positiveAliased:=0;canonicalCorrect:=0
 for i:=0;i<len(rows);i+=6{
  if rows[i+2].Full==rows[i+3].Full&&rows[i+2].Risk==rows[i+3].Risk{indistinguishable++}
  if rows[i+2].Full.Answer{positiveAliased++}
 }
 scores:=map[string]map[string]float64{}
 for _,mode:=range []string{"full","risk","coded_canonicalize_then_frozen"}{
  correct,fp,fn:=0,0,0
  by:=map[int][2]int{}
  for _,r:=range rows{
   var p wf15Choice
   switch mode{case "full":p=r.Full;case "risk":p=r.Risk;default:p=r.Canonical}
   if p.Answer==r.Truth{correct++}
   if p.Answer&&!r.Truth{fp++}
   if !p.Answer&&r.Truth{fn++}
   b:=by[r.Family];if p.Answer!=r.Truth{b[0]++};b[1]++;by[r.Family]=b
   if mode=="coded_canonicalize_then_frozen"&&p.Answer==r.Truth{canonicalCorrect++}
  }
  worst:=0.
  for _,v:=range by{if float64(v[0])/float64(v[1])>worst{worst=float64(v[0])/float64(v[1])}}
  scores[mode]=map[string]float64{"correct":float64(correct),
   "accuracy":float64(correct)/float64(len(rows)),
   "false_positives":float64(fp),"false_negatives":float64(fn),
   "worst_family_error_rate":worst}
 }
 passed:=indistinguishable==96&&canonicalCorrect==576&&
  scores["full"]["accuracy"]<.95
 verdict:="INVALID_ALIAS_CONTROL_OR_IDENTITY_PAIR"
 if indistinguishable==96&&canonicalCorrect==576{
  verdict="DEVELOPMENT_FROZEN_ALIAS_INFORMATION_PATH_NEGATIVE"
  if !passed{verdict="DEVELOPMENT_ALIAS_INFORMATION_PATH_NOT_FALSIFIED"}
 }
 return map[string]any{"schema":"wingless.f20.explicit-alias-info-path-falsifier.dev.v1",
  "stage":"DEVELOPMENT_ONLY","source_F19":"51d82274d97c076edaba0b45cca45e4a05f093c0",
  "model_retrained":false,"native_state_bytes":576,
  "alias_metadata_state_words":[]int{24,25},
  "frozen_model_reads_alias_metadata":false,
  "coded_alias_canonicalizer_is_not_learned":true,
  "families_reused":[]int{108,109,110,111},"new_seed_bounds":[]int{78600,78606},
  "graph_count":96,"question_count":len(rows),
  "alias_positive_queries":96,"alias_negative_queries":96,
  "identical_model_outputs_for_differing_alias_authority_pairs":indistinguishable,
  "positive_aliased_path_predicted_true_by_frozen":positiveAliased,
  "scores":scores,"rows":rows,"registered_verdict":verdict,
  "source_birth_independent":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"external_model_calls":0,
  "elapsed_ms":time.Since(start).Milliseconds()},nil
}