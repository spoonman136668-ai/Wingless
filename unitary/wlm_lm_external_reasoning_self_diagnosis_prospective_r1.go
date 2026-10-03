package unitary

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"sort"
)

type wlmLmExternalReasoningSelfDiagnosisProspectiveR1Preflight struct {
	TargetDomain string `json:"target_domain"`
	TrainingDomains []string `json:"training_domains"`
	TrainingBytes int `json:"training_bytes"`
	TargetEvalBytes int `json:"target_eval_bytes"`
	ZeroSupportClasses []int `json:"zero_support_classes"`
	Risk bool `json:"risk"`
	FreezeSHA256 string `json:"freeze_sha256"`
}

type wlmLmExternalReasoningSelfDiagnosisProspectiveR1Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
	Preflight []wlmLmExternalReasoningSelfDiagnosisProspectiveR1Preflight `json:"preflight"`
}

func wlmLmProspectiveBudgetR1(a,b []byte)([]byte,[]byte){
	const budget=15819
	total:=len(a)+len(b)
	if total<budget || total==0 { return nil,nil }
	na:=budget*len(a)/total
	nb:=budget-na
	if na>len(a) { na=len(a); nb=budget-na }
	if nb>len(b) { nb=len(b); na=budget-nb }
	if na<0 || nb<0 || na>len(a) || nb>len(b) || na+nb!=budget { return nil,nil }
	return a[:na],b[:nb]
}

func wlmLmProspectiveMedianR1(xs []float64) float64 {
	sort.Float64s(xs)
	if len(xs)==0 { return 0 }
	n:=len(xs)
	if n%2==1 { return xs[n/2] }
	return (xs[n/2-1]+xs[n/2])/2
}

func RunWlmLmExternalReasoningSelfDiagnosisProspectiveR1(code, structured, prose []byte) interface{} {
	sources:=[]wlmLmExternalMotifRelationDistributionalR1Source{
		{domain:"code",data:code,sha256:"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",bytes:41453},
		{domain:"structured",data:structured,sha256:"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",bytes:14365},
		{domain:"technical_prose",data:prose,sha256:"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",bytes:1454},
	}
	m:=map[string]float64{
		"source_identity_mismatch_count":0,
		"source_count":3,
		"total_source_bytes":0,
		"arm_count":3,
		"training_byte_budget_per_arm":15819,
		"target_eval_byte_budget_per_arm":1454,
		"preflight_target_byte_use_count":0,
		"preflight_target_label_use_count":0,
		"decision_freeze_count":0,
		"diagnostic_serialization_error_count":0,
		"valid_fixed_class_arm_count":0,
		"prospective_risk_arm_count":0,
		"prospective_confirmed_arm_count":0,
		"prospective_false_positive_arm_count":0,
		"relation_capacity_growth_event_count":0,
		"adaptive_readout_growth_event_count":0,
		"tokenizer_use_count":0,
		"external_model_call_count":0,
		"counter_overflow_count":0,
		"invalid_row_count":0,
		"maximum_probability_mass_error":0,
		"minimum_repaired_retained_mass":1,
	}
	for _,s:=range sources {
		m["total_source_bytes"]+=float64(len(s.data))
		if len(s.data)!=s.bytes || wlmLmExternalRawRepPredR1SHA256(s.data)!=s.sha256 { m["source_identity_mismatch_count"]++ }
	}
	preflights:=make([]wlmLmExternalReasoningSelfDiagnosisProspectiveR1Preflight,0,3)
	for targetIdx,targetSource:=range sources {
		trainIdx:=make([]int,0,2)
		for i:=range sources { if i!=targetIdx { trainIdx=append(trainIdx,i) } }
		a,b:=wlmLmProspectiveBudgetR1(sources[trainIdx[0]].data,sources[trainIdx[1]].data)
		arm:="arm_"+targetSource.domain+"_"
		m[arm+"training_bytes"]=float64(len(a)+len(b))
		if len(a)+len(b)!=15819 {
			m["invalid_row_count"]++
			continue
		}
		local:=map[string]float64{"counter_overflow_rows":0}
		_,selected:=wlmLmRawRepPredFreshHoldoutR1TrainModel([][]byte{a,b},local)
		m["counter_overflow_count"]+=local["counter_overflow_rows"]
		m[arm+"selected_motif_count"]=float64(len(selected))
		if len(selected)!=512 {
			m["invalid_row_count"]++
			continue
		}
		keys:=wlmLmExternalMotifRelationHoldoutR1SortedKeys(selected)
		if len(keys)!=512 {
			m["invalid_row_count"]++
			continue
		}
		index:=make(map[[4]uint8]int,512)
		for i,k:=range keys { index[k]=i }

		streams:=[][]int{
			wlmLmExternalMotifRelationHoldoutR1Decode(a,index),
			wlmLmExternalMotifRelationHoldoutR1Decode(b,index),
		}
		var counts [512]uint64
		var rel [512][512]uint32
		var totals [512]uint64
		var uni [512]uint32
		var unitotal uint64
		pairs:=make(map[wlmLmExternalContextStateHoldoutR1Pair]uint32)
		for _,st:=range streams {
			for _,x:=range st {
				if x<0 || x>=512 { m["invalid_row_count"]++; continue }
				counts[x]++
			}
			for j:=0;j+1<len(st);j++ {
				x,y:=st[j],st[j+1]
				if x<0||x>=512||y<0||y>=512 { m["invalid_row_count"]++; continue }
				if rel[x][y]==^uint32(0) || uni[y]==^uint32(0) { m["counter_overflow_count"]++; continue }
				rel[x][y]++; totals[x]++; uni[y]++; unitotal++
			}
			for j:=1;j<len(st);j++ {
				pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
				if pairs[pair]==^uint32(0) { m["counter_overflow_count"]++; continue }
				pairs[pair]++
			}
		}
		if unitotal==0 { m["invalid_row_count"]++; continue }

		ids:=make([]int,512); for i:=range ids { ids[i]=i }
		sort.Slice(ids,func(i,j int)bool{
			if counts[ids[i]]!=counts[ids[j]] { return counts[ids[i]]>counts[ids[j]] }
			return ids[i]<ids[j]
		})
		var classes [512]int
		var classSizes [8]int
		for rank,id:=range ids { c:=rank/64; classes[id]=c; classSizes[c]++ }
		fixed:=true
		for _,n:=range classSizes { if n!=64 { fixed=false } }
		if !fixed { m["invalid_row_count"]++; continue }
		m["valid_fixed_class_arm_count"]++

		rows:=make([]wlmLmExternalContextStateHoldoutR1PairRow,0,len(pairs))
		for pair,c:=range pairs { rows=append(rows,wlmLmExternalContextStateHoldoutR1PairRow{pair:pair,count:c}) }
		sort.Slice(rows,func(i,j int)bool{
			if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
			if rows[i].pair.prev!=rows[j].pair.prev { return rows[i].pair.prev<rows[j].pair.prev }
			return rows[i].pair.curr<rows[j].pair.curr
		})
		n:=511; if len(rows)<n { n=len(rows) }
		exactID:=make(map[wlmLmExternalContextStateHoldoutR1Pair]int,n)
		for i:=0;i<n;i++ { exactID[rows[i].pair]=i }
		var exact [512][512]uint32
		var exactTotals [512]uint64
		for _,st:=range streams {
			for j:=1;j+1<len(st);j++ {
				pair:=wlmLmExternalContextStateHoldoutR1Pair{prev:st[j-1],curr:st[j]}
				id:=511; if v,ok:=exactID[pair];ok { id=v }
				t:=st[j+1]
				if id<0||id>=512||t<0||t>=512 { m["invalid_row_count"]++; continue }
				if exact[id][t]==^uint32(0) { m["counter_overflow_count"]++; continue }
				exact[id][t]++; exactTotals[id]++
			}
		}
		var eligible [512]bool
		var marginsByRow [512]float64
		rowMargins:=make([]float64,0,512)
		for i:=0;i<512;i++ {
			var total uint64; var top,runner uint32
			for _,v:=range rel[i] {
				total+=uint64(v)
				if v>top { runner=top; top=v } else if v>runner { runner=v }
			}
			if total>=4 {
				eligible[i]=true
				marginsByRow[i]=(float64(top)-float64(runner))/float64(total)
				rowMargins=append(rowMargins,marginsByRow[i])
			}
		}
		if len(rowMargins)==0 { m["invalid_row_count"]++; continue }
		median:=wlmLmProspectiveMedianR1(rowMargins)
		classCounts:=make([]int,8)
		gateMargin:=func(dist *[512]float64) float64 {
			top:=wlmLmExternalReasoningReadoutRefinementR1Top1(dist)
			cid:=classes[top]
			mass,runner:=0.0,0.0
			for k,v:=range dist {
				if classes[k]!=cid { continue }
				mass+=v
				if k!=top && v>runner { runner=v }
			}
			if mass<=0 { m["invalid_row_count"]++; return 0 }
			return (dist[top]-runner)/mass
		}
		for _,st:=range streams {
			for j:=0;j+4<len(st);j++ {
				rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(st[j],st[j+1],3,2048,true,&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal)
				if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
				if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
				_ = gateMargin(&rd)
				top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
				classCounts[classes[top]]++
			}
		}
		zero:=make([]int,0)
		zeroSet:=make(map[int]bool)
		for i,n:=range classCounts {
			m[arm+"class_"+string(rune('0'+i))+"_training_margin_sample_count"]=float64(n)
			if n==0 { zero=append(zero,i); zeroSet[i]=true }
		}
		risk:=len(zero)>0
		if risk { m["prospective_risk_arm_count"]++ }
		freeze:=wlmLmExternalReasoningSelfDiagnosisProspectiveR1Preflight{
			TargetDomain:targetSource.domain,
			TrainingDomains:[]string{sources[trainIdx[0]].domain,sources[trainIdx[1]].domain},
			TrainingBytes:len(a)+len(b),
			TargetEvalBytes:1454,
			ZeroSupportClasses:zero,
			Risk:risk,
		}
		raw,err:=json.Marshal(freeze)
		if err!=nil {
			m["diagnostic_serialization_error_count"]++
			m["invalid_row_count"]++
			continue
		}
		sum:=sha256.Sum256(raw)
		freeze.FreezeSHA256=hex.EncodeToString(sum[:])
		preflights=append(preflights,freeze)
		m["decision_freeze_count"]++

		if len(targetSource.data)<1454 { m["invalid_row_count"]++; continue }
		target:=targetSource.data[:1454]
		st:=wlmLmExternalMotifRelationHoldoutR1Decode(target,index)
		unsupported:=0
		unsupportedHits:=0
		totalQueries:=0
		for j:=0;j+4<len(st);j++ {
			truth:=st[j+4]
			rd,re,rr:=wlmLmExternalRepairedRelationReasoningBridgeR1Terminal(st[j],st[j+1],3,2048,true,&rel,&totals,&exact,&exactTotals,exactID,pairs,&eligible,&marginsByRow,median,&uni,unitotal)
			if re>m["maximum_probability_mass_error"] { m["maximum_probability_mass_error"]=re }
			if rr<m["minimum_repaired_retained_mass"] { m["minimum_repaired_retained_mass"]=rr }
			top:=wlmLmExternalReasoningReadoutRefinementR1Top1(&rd)
			totalQueries++
			if zeroSet[classes[top]] {
				unsupported++
				if top==truth { unsupportedHits++ }
			}
		}
		m[arm+"target_query_count"]=float64(totalQueries)
		m[arm+"unsupported_target_query_count"]=float64(unsupported)
		m[arm+"unsupported_target_exact_hit_count"]=float64(unsupportedHits)
		if risk && unsupported>0 {
			m["prospective_confirmed_arm_count"]++
		} else if risk && unsupported==0 {
			m["prospective_false_positive_arm_count"]++
		}
	}
	if len(preflights)!=3 { m["invalid_row_count"]++ }
	if m["maximum_probability_mass_error"]>1e-9 { m["invalid_row_count"]++ }
	for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
	return wlmLmExternalReasoningSelfDiagnosisProspectiveR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-EXTERNAL-REASONING-SELF-DIAGNOSIS-PROSPECTIVE-R1",
		Metrics:m,
		Preflight:preflights,
	}
}
