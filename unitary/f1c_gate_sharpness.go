package unitary

import ("fmt"; "math"; "time"; "unsafe")

// F1C only diagnoses the ORIGINAL F1 gate family; it does not retrain F1/F1B
// or use any sealed data. Source identities and scientific boundaries stay fixed.
type wf1cAnswers struct {
 Base [2]bool `json:"base"`
 Stress [2]bool `json:"stress"`
}
type wf1cRow struct {
 Seed int `json:"seed"`
 Family int `json:"family"`
 CutA bool `json:"cut_a"`
 CutB bool `json:"cut_b"`
 Truth [2]bool `json:"truth"`
 Models map[string]wf1cAnswers `json:"models"`
}
type wf1cCount struct { BaseCorrect,StressCorrect,BaseFP,StressFP,Changed int }

func RunWFoundationF1C()(map[string]any,error){
 start:=time.Now()
 catalog:=wf1bCatalog()
 models:=map[string]wf1Model{
  "frozen_learned_gate":{Form:"exact_binding",W:wf1bFrozen["exact_binding"]},
  "preexisting_fixed_sharp_control":{Form:"exact_binding",W:[3]float64{13,0,-6.5}},
  "coded_near_binary_gate":{Form:"exact_binding",W:[3]float64{20,0,-15}},
  "zero_link_negative_control":{Form:"exact_binding",W:[3]float64{0,0,-30}},
 }
 totals:=map[string]wf1cCount{}
 rows:=make([]wf1cRow,0,384)
 neg,aliasCount:=0,0
 for seed:=39600;seed<39624;seed++{
  for fam:=4;fam<=7;fam++{
   for _,ca:=range []bool{false,true}{for _,cb:=range []bool{false,true}{
    base,stress,aliases,err:=wf1bStress(seed,fam,ca,cb,catalog)
    if err!=nil{return nil,err}
    aliasCount+=len(aliases)
    row:=wf1cRow{Seed:seed,Family:fam,CutA:ca,CutB:cb,Truth:stress.Labels,Models:map[string]wf1cAnswers{}}
    for _,truth:=range stress.Labels{if !truth{neg++}}
    for name,m:=range models{
     before:=wf1Memory(base,m)
     after:=wf1Memory(stress,m)
     var ans wf1cAnswers
     for q,pair:=range stress.Questions{
      ans.Base[q]=wf1Answer(before,pair)>=.5
      ans.Stress[q]=wf1Answer(after,pair)>=.5
     }
     row.Models[name]=ans
     c:=totals[name]
     for q,truth:=range stress.Labels{
      if ans.Base[q]==truth{c.BaseCorrect++}
      if ans.Stress[q]==truth{c.StressCorrect++}
      if !truth&&ans.Base[q]{c.BaseFP++}
      if !truth&&ans.Stress[q]{c.StressFP++}
      if ans.Base[q]!=ans.Stress[q]{c.Changed++}
     }
     totals[name]=c
    }
    rows=append(rows,row)
   }}
  }
 }
 if len(rows)!=384||neg!=384||aliasCount!=384{return nil,fmt.Errorf("cohort/cardinality drift rows=%d negative=%d aliases=%d",len(rows),neg,aliasCount)}
 if unsafe.Sizeof(wf1State{})!=576{return nil,fmt.Errorf("native state-size drift")}
 summary:=map[string]map[string]float64{}
 for name,c:=range totals{
  summary[name]=map[string]float64{
   "base_accuracy":float64(c.BaseCorrect)/768,
   "stress_accuracy":float64(c.StressCorrect)/768,
   "base_false_positive_rate":float64(c.BaseFP)/384,
   "stress_false_positive_rate":float64(c.StressFP)/384,
   "prediction_flip_count":float64(c.Changed),
  }
 }
 old,sharp:=summary["frozen_learned_gate"],summary["preexisting_fixed_sharp_control"]
 verdict:="GATE_SHARPNESS_CAUSAL_DIAGNOSTIC_UNRESOLVED"
 if sharp["stress_accuracy"]>=.95&&sharp["stress_false_positive_rate"]<=.05&&
   sharp["stress_accuracy"]-old["stress_accuracy"]>=.10&&
   old["stress_false_positive_rate"]-sharp["stress_false_positive_rate"]>=.20 {
  verdict="CODED_GATE_SHARPNESS_EXPLAINS_FROZEN_FALSE_POSITIVES"
 }
 weakMatch:=1/(1+math.Exp(-(wf1bFrozen["exact_binding"][0]+wf1bFrozen["exact_binding"][2])))
 weakNonmatch:=1/(1+math.Exp(-wf1bFrozen["exact_binding"][2]))
 return map[string]any{
  "schema":"wingless.foundation.f1c.gate-sharpness-development.v1",
  "stage":"DEVELOPMENT_ONLY","new_development_seed_bounds":[]int{39600,39624},
  "families":[]int{4,5,6,7},"rows":rows,
  "candidate_scores":summary,"registered_verdict":verdict,
  "frozen_learned_exact_gate":wf1bFrozen["exact_binding"],
  "frozen_match_probability":weakMatch,"frozen_nonmatch_probability":weakNonmatch,
  "source_f1b_commit":"418f822029e6e1528d2f50ae932fb90db6fe6414",
  "source_f1b_artifact_id":11925139848,
  "reference_fixed_control_origin":"unitary/f1_permutation_binding.go (W-F1)",
  "new_graph_states":len(rows),"question_count":768,"negative_count":neg,"aliases":aliasCount,
  "native_state_bytes":576,"coded_equality":true,"coded_four_hop_traversal":true,
  "learned_identity_binding":false,"frozen_no_retraining":true,
  "sealed_first_use":false,"accepted_ref_mutation":false,
  "external_llm_calls":0,"development_wall_ms":time.Since(start).Milliseconds(),
 },nil
}
