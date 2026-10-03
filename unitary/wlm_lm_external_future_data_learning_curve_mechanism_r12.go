package unitary

import "math"

type wlmLmFutureDataLearningCurveMechanismR12Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR12Prediction(
	d []byte,
	i int,
	b *[256][256]uint32,
	s map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif,
)(pred uint8,motifPresent bool,motifBest uint8) {
	pred=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(b,d[i-1])
	if i>=4 {
		k:=[4]uint8{d[i-4],d[i-3],d[i-2],d[i-1]}
		if r,ok:=s[k];ok {
			return r.best,true,r.best
		}
	}
	return pred,false,0
}

func RunWlmLmExternalFutureDataLearningCurveMechanismR12(code,structured,prose []byte) interface{} {
	type src struct{ domain string; data []byte; sha string; size int }
	ss:=[]src{
		{"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453},
		{"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365},
		{"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454},
	}
	cats:=[]string{"motif_admission","motif_eviction","motif_successor_refinement","baseline_byte_refinement","prediction_unchanged"}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,
		"non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
		"target_window_bytes_per_arm":1454,"target_adaptation_bytes_per_arm":872,
		"target_evaluation_bytes_per_arm":582,"adaptation_packet_count":8,
		"transition_observation_count":0,"transition_accounting_error_count":0,
		"arm_code_reconstructed_cumulative_gain":0,
		"arm_structured_reconstructed_cumulative_gain":0,
		"arm_technical_prose_reconstructed_cumulative_gain":0,
		"code_excess_gain":109.5,
		"dominant_category_code":0,
		"dominant_category_excess_gain":0,
		"dominant_category_excess_fraction":0,
		"capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
		"counter_overflow_count":0,"invalid_row_count":0,
	}
	for _,c:=range cats {
		m["arm_code_"+c+"_net_gain"]=0
		m["arm_structured_"+c+"_net_gain"]=0
		m["arm_technical_prose_"+c+"_net_gain"]=0
	}
	for _,s:=range ss {
		m["total_source_bytes"]+=float64(len(s.data))
		if len(s.data)!=s.size || wlmLmExternalRawRepPredR1SHA256(s.data)!=s.sha { m["source_identity_mismatch_count"]++ }
	}
	for ti,t:=range ss {
		idx:=make([]int,0,2)
		for i:=range ss { if i!=ti { idx=append(idx,i) } }
		a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].data,ss[idx[1]].data)
		if len(a)+len(b)!=15819 { m["invalid_row_count"]++; continue }
		if len(t.data)<1454 { m["invalid_row_count"]++; continue }
		adapt:=t.data[:872]
		eval:=t.data[872:1454]
		baseRel,baseSel:=wlmLmR10Train256([][]byte{a,b},m)
		if len(baseSel)!=256 { m["invalid_row_count"]++ }
		armPrefix:="arm_"+t.domain+"_"
		armNet:=0.0
		for p:=1;p<=8;p++ {
			adaptRel,adaptSel:=wlmLmR10Train256([][]byte{a,b,adapt[:p*109]},m)
			if len(adaptSel)!=256 { m["invalid_row_count"]++ }
			for i:=1;i<len(eval);i++ {
				basePred,baseMotif,baseBest:=wlmLmR12Prediction(eval,i,&baseRel,baseSel)
				adaptPred,adaptMotif,adaptBest:=wlmLmR12Prediction(eval,i,&adaptRel,adaptSel)
				truth:=eval[i]
				delta:=0.0
				if adaptPred==truth { delta++ }
				if basePred==truth { delta-- }
				category:="prediction_unchanged"
				if basePred!=adaptPred {
					switch {
					case !baseMotif && adaptMotif:
						category="motif_admission"
					case baseMotif && !adaptMotif:
						category="motif_eviction"
					case baseMotif && adaptMotif && baseBest!=adaptBest:
						category="motif_successor_refinement"
					case !baseMotif && !adaptMotif:
						category="baseline_byte_refinement"
					default:
						m["transition_accounting_error_count"]++
					}
				}
				m[armPrefix+category+"_net_gain"]+=delta
				armNet+=delta
				m["transition_observation_count"]++
			}
		}
		m[armPrefix+"reconstructed_cumulative_gain"]=armNet
	}
	if m["arm_code_reconstructed_cumulative_gain"]!=210 { m["transition_accounting_error_count"]++ }
	if m["arm_structured_reconstructed_cumulative_gain"]!=98 { m["transition_accounting_error_count"]++ }
	if m["arm_technical_prose_reconstructed_cumulative_gain"]!=103 { m["transition_accounting_error_count"]++ }
	categoryCodes:=map[string]float64{
		"motif_admission":1,
		"motif_eviction":2,
		"motif_successor_refinement":3,
		"baseline_byte_refinement":4,
	}
	bestExcess:=-math.MaxFloat64
	bestCode:=0.0
	for _,c:=range cats[:4] {
		ex:=m["arm_code_"+c+"_net_gain"]-(m["arm_structured_"+c+"_net_gain"]+m["arm_technical_prose_"+c+"_net_gain"])/2.0
		m[c+"_code_excess_gain"]=ex
		if ex>bestExcess {
			bestExcess=ex
			bestCode=categoryCodes[c]
		}
	}
	m["dominant_category_code"]=bestCode
	m["dominant_category_excess_gain"]=bestExcess
	if m["code_excess_gain"]!=0 { m["dominant_category_excess_fraction"]=bestExcess/m["code_excess_gain"] }
	if m["transition_observation_count"]!=13944 { m["transition_accounting_error_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmFutureDataLearningCurveMechanismR12Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNING-CURVE-MECHANISM-R12",
		Metrics:m,
	}
}
