package unitary
import "math"
type wlmLmFutureDataAllocationAttributionR14Result struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;Metrics map[string]float64 `json:"metrics"`}
func RunWlmLmExternalFutureDataAllocationAttributionR14(code,structured,prose []byte) interface{}{
 type src struct{d string;b []byte;h string;n,e,g int}
 ss:=[]src{{"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453,582,872},{"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365,581,436},{"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454,581,436}}
 cats:=[]string{"motif_admission_in_guided","motif_eviction_in_guided","motif_successor_refinement","baseline_byte_refinement","prediction_unchanged"}
 m:=map[string]float64{"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,"target_window_bytes_per_arm":1454,"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,"transition_observation_count":0,"transition_accounting_error_count":0,"arm_code_reconstructed_delta":0,"arm_structured_reconstructed_delta":0,"arm_technical_prose_reconstructed_delta":0,"aggregate_reconstructed_delta":0,"dominant_prose_loss_category_code":0,"dominant_prose_loss_abs":0,"dominant_prose_loss_fraction":0,"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0}
 for _,c:=range cats{m["arm_code_"+c+"_net_delta"]=0;m["arm_structured_"+c+"_net_delta"]=0;m["arm_technical_prose_"+c+"_net_delta"]=0}
 for _,s:=range ss{m["total_source_bytes"]+=float64(len(s.b));if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h{m["source_identity_mismatch_count"]++}}
 for ti,t:=range ss{
  idx:=[]int{};for i:=range ss{if i!=ti{idx=append(idx,i)}}
  a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b);if len(a)+len(b)!=15819||len(t.b)<1454{m["invalid_row_count"]++;continue}
  eval,res:=t.b[:582],t.b[582:1454]
  er,es:=wlmLmR10Train256([][]byte{a,b,res[:t.e]},m);gr,gs:=wlmLmR10Train256([][]byte{a,b,res[:t.g]},m)
  if len(es)!=256||len(gs)!=256{m["invalid_row_count"]++}
  p:="arm_"+t.d+"_";arm:=0.0
  for i:=1;i<len(eval);i++{
   ep,em,eb:=wlmLmR12Prediction(eval,i,&er,es);gp,gm,gb:=wlmLmR12Prediction(eval,i,&gr,gs);truth:=eval[i]
   delta:=0.0;if gp==truth{delta++};if ep==truth{delta--}
   cat:="prediction_unchanged"
   if ep!=gp{
    switch{
    case !em&&gm:cat="motif_admission_in_guided"
    case em&&!gm:cat="motif_eviction_in_guided"
    case em&&gm&&eb!=gb:cat="motif_successor_refinement"
    case !em&&!gm:cat="baseline_byte_refinement"
    default:m["transition_accounting_error_count"]++
    }
   }
   m[p+cat+"_net_delta"]+=delta;arm+=delta;m["transition_observation_count"]++
  }
  m[p+"reconstructed_delta"]=arm;m["aggregate_reconstructed_delta"]+=arm
 }
 if m["arm_code_reconstructed_delta"]!=5{m["transition_accounting_error_count"]++}
 if m["arm_structured_reconstructed_delta"]!=0{m["transition_accounting_error_count"]++}
 if m["arm_technical_prose_reconstructed_delta"]!=-10{m["transition_accounting_error_count"]++}
 if m["aggregate_reconstructed_delta"]!=-5{m["transition_accounting_error_count"]++}
 codes:=map[string]float64{"motif_admission_in_guided":1,"motif_eviction_in_guided":2,"motif_successor_refinement":3,"baseline_byte_refinement":4}
 best:=0.0;codev:=0.0
 for _,c:=range cats[:4]{v:=m["arm_technical_prose_"+c+"_net_delta"];if v<best{best=v;codev=codes[c]}}
 m["dominant_prose_loss_category_code"]=codev;m["dominant_prose_loss_abs"]=-best;if best<0{m["dominant_prose_loss_fraction"]=(-best)/10.0}
 if m["transition_observation_count"]!=1743{m["transition_accounting_error_count"]++}
 for _,v:=range m{if math.IsNaN(v)||math.IsInf(v,0){m["invalid_row_count"]++}}
 return wlmLmFutureDataAllocationAttributionR14Result{"wingless.research-scientific-result.v1","WLM-LM-EXTERNAL-FUTURE-DATA-ALLOCATION-ATTRIBUTION-R14",m}
}
