package unitary
import "math"
type wlmLmFutureDataMechanismGuidedAllocationR13Result struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;Metrics map[string]float64 `json:"metrics"`}
func RunWlmLmExternalFutureDataMechanismGuidedAllocationR13(code,structured,prose []byte) interface{}{
 type src struct{d string;b []byte;h string;n,e,g int}
 ss:=[]src{{"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453,582,872},{"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365,581,436},{"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454,581,436}}
 m:=map[string]float64{"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,"target_window_bytes_per_arm":1454,"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,"total_adaptation_budget_equal":0,"total_adaptation_budget_guided":0,"allocation_target_label_use_count":0,"allocation_post_result_choice_count":0,"packet0_aggregate_exact_hits":0,"equal_aggregate_exact_hits":0,"guided_aggregate_exact_hits":0,"guided_aggregate_exact_hit_delta_vs_equal":0,"guided_code_exact_hit_delta_vs_equal":0,"guided_domain_below_packet0_count":0,"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0}
 for _,s:=range ss{m["total_source_bytes"]+=float64(len(s.b));m["total_adaptation_budget_equal"]+=float64(s.e);m["total_adaptation_budget_guided"]+=float64(s.g);if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h{m["source_identity_mismatch_count"]++}}
 for ti,t:=range ss{
  idx:=[]int{};for i:=range ss{if i!=ti{idx=append(idx,i)}}
  a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b);if len(a)+len(b)!=15819||len(t.b)<1454{m["invalid_row_count"]++;continue}
  eval,res:=t.b[:582],t.b[582:1454]
  br,bs:=wlmLmR10Train256([][]byte{a,b},m);er,es:=wlmLmR10Train256([][]byte{a,b,res[:t.e]},m);gr,gs:=wlmLmR10Train256([][]byte{a,b,res[:t.g]},m)
  if len(bs)!=256||len(es)!=256||len(gs)!=256{m["invalid_row_count"]++}
  p0:=wlmLmR10Hits(eval,&br,bs);eh:=wlmLmR10Hits(eval,&er,es);gh:=wlmLmR10Hits(eval,&gr,gs);p:="arm_"+t.d+"_"
  m[p+"packet0_exact_hits"]=float64(p0);m[p+"equal_exact_hits"]=float64(eh);m[p+"guided_exact_hits"]=float64(gh);m[p+"guided_delta_vs_equal"]=float64(gh-eh);m[p+"guided_delta_vs_packet0"]=float64(gh-p0)
  m["packet0_aggregate_exact_hits"]+=float64(p0);m["equal_aggregate_exact_hits"]+=float64(eh);m["guided_aggregate_exact_hits"]+=float64(gh)
  if gh<p0{m["guided_domain_below_packet0_count"]++};if t.d=="code"{m["guided_code_exact_hit_delta_vs_equal"]=float64(gh-eh)}
 }
 m["guided_aggregate_exact_hit_delta_vs_equal"]=m["guided_aggregate_exact_hits"]-m["equal_aggregate_exact_hits"]
 if m["total_adaptation_budget_equal"]!=1744||m["total_adaptation_budget_guided"]!=1744{m["invalid_row_count"]++}
 for _,v:=range m{if math.IsNaN(v)||math.IsInf(v,0){m["invalid_row_count"]++}}
 return wlmLmFutureDataMechanismGuidedAllocationR13Result{"wingless.research-scientific-result.v1","WLM-LM-EXTERNAL-FUTURE-DATA-MECHANISM-GUIDED-ALLOCATION-R13",m}
}
