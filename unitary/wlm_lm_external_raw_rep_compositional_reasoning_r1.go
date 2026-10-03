package unitary

import (
	"math"
	"sort"
)

type wlmLmExternalRawRepCompositionalReasoningR1Source struct {
	data []byte
	sha256 string
	bytes int
}

type wlmLmExternalRawRepCompositionalReasoningR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmLmExternalRawRepCompositionalReasoningR1DecodeSlots(data []byte, values [][4]uint8) ([]int,bool) {
	if len(data)%4!=0 { return nil,false }
	ids:=make(map[[4]uint8]int,len(values))
	for i,key:=range values { ids[key]=i }
	out:=make([]int,0,len(data)/4)
	for off:=0;off<len(data);off+=4 {
		key:=[4]uint8{data[off],data[off+1],data[off+2],data[off+3]}
		id,ok:=ids[key]
		if !ok { return nil,false }
		out=append(out,id)
	}
	return out,true
}

func wlmLmExternalRawRepCompositionalReasoningR1Record(values [][4]uint8,a,b,target int,includeTarget bool) []byte {
	out:=make([]byte,0,12)
	out=append(out,values[a][:]...)
	out=append(out,values[b][:]...)
	if includeTarget { out=append(out,values[target][:]...) }
	return out
}

func RunWlmLmExternalRawRepCompositionalReasoningR1(code, structured, prose []byte) interface{} {
	sources:=[]wlmLmExternalRawRepCompositionalReasoningR1Source{
		{data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	metrics:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":float64(len(sources)),
		"total_source_bytes":0,
		"external_representation_pool_count":0,
		"reasoning_value_count":0,
		"reasoning_value_distinct_count":0,
		"training_record_count":0,
		"heldout_pair_count":0,
		"training_representation_decode_failure_count":0,
		"heldout_representation_decode_failure_count":0,
		"heldout_context_leak_count":0,
		"factorized_bit_coverage_missing_count":0,
		"factorized_heldout_accuracy":0,
		"lookup_heldout_accuracy":0,
		"raw_output_exact_accuracy":0,
		"factorized_counter_count":24,
		"lookup_counter_count":512,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"capacity_growth_event_count":0,
		"invalid_row_count":0,
		"counter_overflow_count":0,
		"counter_overflow_rows":0,
	}

	train:=make([][]byte,0,len(sources))
	for _,src:=range sources {
		metrics["total_source_bytes"]+=float64(len(src.data))
		if len(src.data)!=src.bytes || wlmLmExternalRawRepPredR1SHA256(src.data)!=src.sha256 {
			metrics["source_identity_mismatch_count"]++
		}
		split:=len(src.data)*3/5
		if split<5 { metrics["invalid_row_count"]++ }
		train=append(train,src.data[:split])
	}

	_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel(train,metrics)
	metrics["external_representation_pool_count"]=float64(len(selected))
	rows:=make([]wlmLmRawRepPredFreshHoldoutR1Motif,0,len(selected))
	for _,row:=range selected { rows=append(rows,row) }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].bestCount!=rows[j].bestCount { return rows[i].bestCount>rows[j].bestCount }
		if rows[i].consistency!=rows[j].consistency { return rows[i].consistency>rows[j].consistency }
		return wlmLmRawRepPredFreshHoldoutR1Less(rows[i].key,rows[j].key)
	})
	if len(rows)<8 {
		metrics["invalid_row_count"]++
	}
	limit:=8
	if len(rows)<limit { limit=len(rows) }
	values:=make([][4]uint8,limit)
	for i:=0;i<limit;i++ { values[i]=rows[i].key }
	sort.Slice(values,func(i,j int)bool{
		return wlmLmRawRepPredFreshHoldoutR1Less(values[i],values[j])
	})
	metrics["reasoning_value_count"]=float64(len(values))
	distinct:=make(map[[4]uint8]bool,len(values))
	for _,v:=range values { distinct[v]=true }
	metrics["reasoning_value_distinct_count"]=float64(len(distinct))
	if len(values)!=8 {
		metrics["invalid_row_count"]++
		return wlmLmExternalRawRepCompositionalReasoningR1Result{
			Schema:"wingless.research-scientific-result.v1",
			Experiment:"WLM-LM-EXTERNAL-RAW-REP-COMPOSITIONAL-REASONING-R1",
			Metrics:metrics,
		}
	}

	var factor wlmLmRawRepCompositionalReasoningR1Factor
	var lookup wlmLmRawRepCompositionalReasoningR1Lookup
	trainingPairs:=make([][2]int,0,48)
	seenTrainingPair:=make(map[[2]int]bool)

	for repeat:=0;repeat<32;repeat++ {
		for a:=0;a<8;a++ {
			for b:=0;b<8;b++ {
				if wlmLmRawRepCompositionalReasoningR1Heldout(a,b) { continue }
				t:=a^b
				raw:=wlmLmExternalRawRepCompositionalReasoningR1Record(values,a,b,t,true)
				decoded,ok:=wlmLmExternalRawRepCompositionalReasoningR1DecodeSlots(raw,values)
				if !ok || len(decoded)!=3 {
					metrics["training_representation_decode_failure_count"]++
					continue
				}
				if repeat==0 { trainingPairs=append(trainingPairs,[2]int{decoded[0],decoded[1]}) }
				seenTrainingPair[[2]int{decoded[0],decoded[1]}]=true
				factor.observe(decoded[0],decoded[1],decoded[2],metrics)
				lookup.observe(decoded[0],decoded[1],decoded[2],metrics)
				metrics["training_record_count"]++
			}
		}
	}
	metrics["factorized_bit_coverage_missing_count"]=float64(
		wlmLmRawRepCompositionalReasoningR1CoverageMissing(trainingPairs),
	)

	factorHits:=0
	lookupHits:=0
	rawHits:=0
	validHeldout:=0
	for a:=0;a<8;a++ {
		for b:=0;b<8;b++ {
			if !wlmLmRawRepCompositionalReasoningR1Heldout(a,b) { continue }
			metrics["heldout_pair_count"]++
			t:=a^b
			raw:=wlmLmExternalRawRepCompositionalReasoningR1Record(values,a,b,t,false)
			decoded,ok:=wlmLmExternalRawRepCompositionalReasoningR1DecodeSlots(raw,values)
			if !ok || len(decoded)!=2 {
				metrics["heldout_representation_decode_failure_count"]++
				continue
			}
			if seenTrainingPair[[2]int{decoded[0],decoded[1]}] {
				metrics["heldout_context_leak_count"]++
			}
			fp:=factor.predict(decoded[0],decoded[1])
			lp:=lookup.predict(decoded[0],decoded[1])
			if fp==t { factorHits++ }
			if lp==t { lookupHits++ }
			if fp>=0 && fp<len(values) && values[fp]==values[t] { rawHits++ }
			validHeldout++
		}
	}
	if validHeldout>0 {
		metrics["factorized_heldout_accuracy"]=float64(factorHits)/float64(validHeldout)
		metrics["lookup_heldout_accuracy"]=float64(lookupHits)/float64(validHeldout)
		metrics["raw_output_exact_accuracy"]=float64(rawHits)/float64(validHeldout)
	}
	for _,v:=range metrics {
		if math.IsNaN(v)||math.IsInf(v,0) { metrics["invalid_row_count"]++ }
	}

	return wlmLmExternalRawRepCompositionalReasoningR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-RAW-REP-COMPOSITIONAL-REASONING-R1",
		Metrics:metrics,
	}
}
