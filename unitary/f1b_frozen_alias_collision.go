package unitary

import("fmt";"math";"time";"unsafe")

// W-F1B: development-only adversarial near-alias test of ORIGINAL frozen F1
// three-parameter controllers. All endpoint equality and four-hop traversal
// remain programmed and are evaluated at the exact same 576 resident bytes.
type wf1bAlias struct {
 Intended int `json:"intended"`
 Alias int `json:"alias"`
 Cosine float64 `json:"cosine"`
}
type wf1bAnswers struct {
 Base [2]bool `json:"base"`
 Stress [2]bool `json:"stress"`
}
type wf1bRecord struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 CutA bool `json:"cut_a"`
 CutB bool `json:"cut_b"`
 Truth [2]bool `json:"truth"`
 Aliases []wf1bAlias `json:"aliases"`
 Models map[string]wf1bAnswers `json:"models"`
}
type wf1bScore struct {
 BaseCorrect int `json:"base_correct"`
 StressCorrect int `json:"stress_correct"`
 BaseWrongPositive int `json:"base_wrong_positive"`
 StressWrongPositive int `json:"stress_wrong_positive"`
}
var wf1bFrozen=map[string][3]float64{
 "exact_binding":{.09502010707601494,0,-1.4396918650883532},
 "identity_similarity":{.10059148842305103,0,-1.4396162587555883},
 "hybrid":{.13090927403123004,.31557067618237516,-1.4679481396881129},
}
func wf1bCatalog()[][12]float64{
 data:=make([][12]float64,60000)
 for id:=range data{data[id]=nci11WCode(id)}
 return data
}
func wf1bClosest(catalog [][12]float64,want int,used map[int]bool)(int,float64){
 best,bestSimilarity:=-1,math.Inf(-1)
 target:=catalog[want]
 for id,emb:=range catalog{
  if id==want||used[id]{continue}
  cos:=0.
  for k,x:=range target{cos+=x*emb[k]}
  if cos>bestSimilarity{best,bestSimilarity=id,cos}
 }
 return best,bestSimilarity
}
func wf1bStress(seed,fam int,ca,cb bool,catalog [][12]float64)(wf1Datum,wf1Datum,[]wf1bAlias,error){
 original,err:=wf1New(seed,fam,ca,cb,1)
 if err!=nil{return original,original,nil,err}
 stress:=original
 chains,_,_:=nci10WFixture(seed,fam%4)
 used:=map[int]bool{}
 for _,e:=range stress.Edges{used[e.From]=true;used[e.To]=true}
 aliases:=[]wf1bAlias{}
 for side,broken:=range []bool{ca,cb}{
  if !broken{continue}
  oldFrom:=wf1Rename(chains[side][2],seed,1)
  oldTo:=wf1Rename(300+seed%13,seed,1)
  if side==1{oldTo=wf1Rename(320+seed%13,seed,1)}
  intended:=wf1Rename(chains[side][3],seed,1)
  chosen,cos:=wf1bClosest(catalog,intended,used)
  if chosen<0{return original,stress,nil,fmt.Errorf("missing near-collision identity")}
  found:=false
  for j,e:=range stress.Edges{
   if e.From==oldFrom&&e.To==oldTo{
    stress.Edges[j]=nci02WEdge{From:oldFrom,To:chosen}
    found=true;break
   }
  }
  if !found{return original,stress,nil,fmt.Errorf("missing counterfactual edge for side%d",side)}
  aliases=append(aliases,wf1bAlias{Intended:intended,Alias:chosen,Cosine:cos})
  used[chosen]=true
 }
 stress.Eq=[12][12]float64{}
 stress.Similar=[12][12]float64{}
 for i,a:=range stress.Edges{
  vecA:=nci11WCode(a.To)
  for j,b:=range stress.Edges{
   if a.To==b.From{stress.Eq[i][j]=1}
   vecB:=nci11WCode(b.From)
   for k:=0;k<12;k++{stress.Similar[i][j]+=vecA[k]*vecB[k]}
  }
 }
 for q,pair:=range stress.Questions{
  if nci04WOracle(stress.Edges[:],pair[0],pair[1],4)!=stress.Labels[q]{
   return original,stress,nil,fmt.Errorf("adversarial alias changed ground-truth")
  }
 }
 return original,stress,aliases,nil
}
func RunWFoundationF1B()(map[string]any,error){
 start:=time.Now()
 catalog:=wf1bCatalog()
 models:=map[string]wf1Model{}
 for name,w:=range wf1bFrozen{
  models[name]=wf1Model{Form:name,W:w}
 }
 records:=[]wf1bRecord{}
 totals:=map[string]wf1bScore{}
 similaritySum,similarityMin,aliasCount:=0.,math.Inf(1),0
 for seed:=39312;seed<39336;seed++{
  for fam:=4;fam<=7;fam++{
   for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
    base,stress,aliases,err:=wf1bStress(seed,fam,ca,cb,catalog)
    if err!=nil{return nil,err}
    for _,x:=range aliases{
     similaritySum+=x.Cosine
     if x.Cosine<similarityMin{similarityMin=x.Cosine}
     aliasCount++
    }
    r:=wf1bRecord{Seed:seed,Family:fam,CutA:ca,CutB:cb,
      Truth:stress.Labels,Aliases:aliases,Models:map[string]wf1bAnswers{}}
    for name,m:=range models{
     oldState:=wf1Memory(base,m)
     newState:=wf1Memory(stress,m)
     values:=wf1bAnswers{}
     for q,pair:=range stress.Questions{
      values.Base[q]=wf1Answer(oldState,pair)>=.5
      values.Stress[q]=wf1Answer(newState,pair)>=.5
     }
     r.Models[name]=values
     c:=totals[name]
     for q,truth:=range stress.Labels{
      if values.Base[q]==truth{c.BaseCorrect++}
      if values.Stress[q]==truth{c.StressCorrect++}
      if !truth{
       if values.Base[q]{c.BaseWrongPositive++}
       if values.Stress[q]{c.StressWrongPositive++}
      }
     }
     totals[name]=c
    }
    records=append(records,r)
   }}
  }
 }
 if len(records)!=384||aliasCount!=384{return nil,fmt.Errorf("F1B expected 384 cases/aliases")}
 questionCount:=len(records)*2
 negativeCount:=0
 for _,r:=range records{for _,truth:=range r.Truth{if !truth{negativeCount++}}}
 if negativeCount!=384{return nil,fmt.Errorf("unexpected class balance")}
 measured:=map[string]map[string]float64{}
 for name,count:=range totals{
  measured[name]=map[string]float64{
   "base_accuracy":float64(count.BaseCorrect)/float64(questionCount),
   "stress_accuracy":float64(count.StressCorrect)/float64(questionCount),
   "false_positive_rate_negative":float64(count.StressWrongPositive)/float64(negativeCount),
   "baseline_false_positive_rate":float64(count.BaseWrongPositive)/float64(negativeCount)}
 }
 eq:=measured["exact_binding"]
 sim:=measured["identity_similarity"]
 useful:=eq["stress_accuracy"]>=.90&&eq["stress_accuracy"]-sim["stress_accuracy"]>=.10&&
   sim["false_positive_rate_negative"]-eq["false_positive_rate_negative"]>=.10
 verdict:="NEAR_COLLISION_DIAGNOSTIC_INCONCLUSIVE"
 if useful{verdict="CODED_EQUALITY_RESISTS_ADVERSARIAL_ID_COLLISIONS"}
 if unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("native 576 memory budget drift")}
 return map[string]any{
  "schema":"wingless.foundation.f1b.adversarial-alias-development.v1",
  "stage":"DEVELOPMENT_ONLY",
  "source_parent":"079b075a2f939caa7546bbbdd37742169c88620a",
  "frozen_training_origin_artifact":11925452674,
  "seed_bounds":[]int{39312,39336},
  "families":[]int{4,5,6,7},
  "original_rows":records,"question_count":questionCount,
  "negative_question_count":negativeCount,
  "frozen_model_weights":wf1bFrozen,
  "candidate_scores":measured,
  "minimum_alias_feature_similarity":similarityMin,
  "mean_alias_feature_similarity":similaritySum/float64(aliasCount),
  "exact_public_graph_oracle_correct":questionCount,
  "registered_verdict":verdict,
  "native_state_bytes":576,"trained_parameters_per_candidate":3,
  "trained_model_parameter_bytes_float64":24,
  "adversarial_search_count":60000,
  "public_similarity_based_alias_construction":true,
  "coded_entity_equality":true,"coded_four_hop_traversal":true,
  "learned_identity_binding":false,
  "original_development_wall_ms":time.Since(start).Milliseconds(),
  "sealed_first_use":false,"sealed_cohort_reserved":false,
  "external_llm_calls":0,"accepted_ref_mutation":false,
  "full_nci01_acceptance":false,
 },nil
}
