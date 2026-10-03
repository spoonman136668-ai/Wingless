package unitary
import "math"
type wlmLmFutureDataAllocationMarginalEfficiencyR15Result struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;Metrics map[string]float64 `json:"metrics"`}
func RunWlmLmExternalFutureDataAllocationMarginalEfficiencyR15(code,structured,prose []byte) interface{}{
 type src struct{d string;b []byte;h string;n int;budgets []int}
 ss:=[]src{
  {"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453,[]int{0,436,582,727,872}},
  {"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365,[]int{0,291,436,581,726}},
  {"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454,[]int{0,291,436,581,726}},
 }
 m:=map[string]float64{"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,"target_evaluation_bytes_per_arm":582,"target_adaptation_reservoir_bytes_per_arm":872,"budget_point_count_total":15,"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,"counter_overflow_count":0,"invalid_row_count":0}
 for _,s:=range ss{m["total_source_bytes"]+=float64(len(s.b));if len(s.b)!=s.n||wlmLmExternalRawRepPredR1SHA256(s.b)!=s.h{m["source_identity_mismatch_count"]++}}
 for ti,t:=range ss{
  idx:=[]int{};for i:=range ss{if i!=ti{idx=append(idx,i)}}
  a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].b,ss[idx[1]].b);if len(a)+len(b)!=15819||len(t.b)<1454{m["invalid_row_count"]++;continue}
  eval,res:=t.b[:582],t.b[582:1454];hits:=make([]int,len(t.budgets))
  for bi,bud:=range t.budgets{
   train:=[][]byte{a,b};if bud>0{train=append(train,res[:bud])}
   rel,sel:=wlmLmR10Train256(train,m);if len(sel)!=256{m["invalid_row_count"]++}
   h:=wlmLmR10Hits(eval,&rel,sel);hits[bi]=h;m["arm_"+t.d+"_budget_"+itoaR15(bud)+"_exact_hits"]=float64(h)
  }
  for i:=1;i<len(t.budgets);i++{
   lo,hi:=t.budgets[i-1],t.budgets[i];gain:=hits[i]-hits[i-1];eff:=float64(gain)/float64(hi-lo)
   k:="arm_"+t.d+"_interval_"+itoaR15(lo)+"_"+itoaR15(hi)
   m[k+"_gain"]=float64(gain);m[k+"_efficiency"]=eff
  }
 }
 pe:=m["arm_technical_prose_interval_436_581_efficiency"];ce:=float64(m["arm_code_budget_872_exact_hits"]-m["arm_code_budget_582_exact_hits"])/290.0
 m["prose_eff_436_581"]=pe;m["code_eff_582_872"]=ce;m["prose_gain_436_581"]=m["arm_technical_prose_interval_436_581_gain"]
 m["code_eff_582_727"]=m["arm_code_interval_582_727_efficiency"];m["code_eff_727_872"]=m["arm_code_interval_727_872_efficiency"]
 for _,v:=range m{if math.IsNaN(v)||math.IsInf(v,0){m["invalid_row_count"]++}}
 return wlmLmFutureDataAllocationMarginalEfficiencyR15Result{"wingless.research-scientific-result.v1","WLM-LM-EXTERNAL-FUTURE-DATA-ALLOCATION-MARGINAL-EFFICIENCY-R15",m}
}
func itoaR15(n int)string{if n==0{return"0"};b:=[20]byte{};i:=len(b);for n>0{i--;b[i]=byte('0'+n%10);n/=10};return string(b[i:])}
