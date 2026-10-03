package unitary

import (
  "math"
  "sort"
)

type wlmLmFutureDataSignalFactorialR11Result struct {
  Schema string `json:"schema"`
  Experiment string `json:"experiment"`
  Metrics map[string]float64 `json:"metrics"`
}

func wlmLmR11Eligible(files [][]byte) int {
  c:=map[[4]uint8]uint32{}
  for _,d:=range files {
    for i:=4;i<len(d);i++ {
      k:=[4]uint8{d[i-4],d[i-3],d[i-2],d[i-1]}
      c[k]++
    }
  }
  n:=0
  for _,v:=range c { if v>=4 { n++ } }
  return n
}

func wlmLmR11Train256(files [][]byte, m map[string]float64) ([256][256]uint32,map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) {
  local:=map[string]float64{"counter_overflow_rows":0}
  baseline,raw:=wlmLmRawRepPredFreshHoldoutR1TrainModel(files,local)
  m["counter_overflow_count"]+=local["counter_overflow_rows"]
  rows:=make([]wlmLmRawRepPredFreshHoldoutR1Motif,0,len(raw))
  for _,r:=range raw { rows=append(rows,r) }
  sort.Slice(rows,func(i,j int)bool{
    if rows[i].bestCount!=rows[j].bestCount { return rows[i].bestCount>rows[j].bestCount }
    if rows[i].consistency!=rows[j].consistency { return rows[i].consistency>rows[j].consistency }
    return wlmLmRawRepPredFreshHoldoutR1Less(rows[i].key,rows[j].key)
  })
  if len(rows)>256 { rows=rows[:256] }
  out:=make(map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif,len(rows))
  for _,r:=range rows { out[r.key]=r }
  return baseline,out
}

func wlmLmR11Hits(d []byte,b *[256][256]uint32,s map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) int {
  hits:=0
  for i:=1;i<len(d);i++ {
    p:=wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(b,d[i-1])
    if i>=4 {
      k:=[4]uint8{d[i-4],d[i-3],d[i-2],d[i-1]}
      if r,ok:=s[k];ok { p=r.best }
    }
    if p==d[i] { hits++ }
  }
  return hits
}

func RunWlmLmExternalFutureDataLearningEfficiencySignalFactorialR11(code,structured,prose []byte) interface{} {
  type src struct{ domain string; data []byte; sha string; size int }
  ss:=[]src{
    {"code",code,"7a95f1c506c9ac4b2277df5f2bdd9d61cc67b520c45021a5a961939770221ef6",41453},
    {"structured",structured,"4c5cbe6cbcd28af73761091367b20e07d0403847e236c06c31fc27061bd81192",14365},
    {"technical_prose",prose,"8247b7c5de1e74854aac1a08aa5894444d1d33b4045c70d5cc3367ad0e25c3f3",1454},
  }
  m:=map[string]float64{
    "source_identity_mismatch_count":0,"source_count":3,"total_source_bytes":0,"arm_count":3,
    "non_target_training_byte_budget_per_arm":15819,"motif_capacity":256,
    "target_window_bytes_per_arm":1454,"target_adaptation_bytes_per_arm":872,
    "target_evaluation_bytes_per_arm":582,"adaptation_packet_count":8,
    "signal_target_byte_use_count":0,"signal_target_label_use_count":0,"signal_post_result_choice_count":0,
    "tested_signal_count":3,
    "H_pairwise_agreement_count":0,"H_pairwise_reversal_count":0,"H_pairwise_tie_count":0,
    "D_pairwise_agreement_count":0,"D_pairwise_reversal_count":0,"D_pairwise_tie_count":0,
    "M_pairwise_agreement_count":0,"M_pairwise_reversal_count":0,"M_pairwise_tie_count":0,
    "best_signal_code":0,"best_signal_pairwise_agreement_count":0,"best_signal_pairwise_reversal_count":0,
    "capacity_growth_event_count":0,"tokenizer_use_count":0,"external_model_call_count":0,
    "counter_overflow_count":0,"invalid_row_count":0,
  }
  for _,s:=range ss {
    m["total_source_bytes"]+=float64(len(s.data))
    if len(s.data)!=s.size||wlmLmExternalRawRepPredR1SHA256(s.data)!=s.sha { m["source_identity_mismatch_count"]++ }
  }
  head:=make([]float64,3)
  deficit:=make([]float64,3)
  marginMedian:=make([]float64,3)
  cumulative:=make([]float64,3)
  finalGain:=make([]float64,3)
  for ti,t:=range ss {
    idx:=[]int{}
    for i:=range ss { if i!=ti { idx=append(idx,i) } }
    a,b:=wlmLmTransferCapacityBudgetR2(ss[idx[0]].data,ss[idx[1]].data)
    if len(a)+len(b)!=15819 { m["invalid_row_count"]++; continue }
    eligible:=wlmLmR11Eligible([][]byte{a,b})
    head[ti]=float64(eligible-256)
    pfx:="arm_"+t.domain+"_"
    m[pfx+"eligible_motif_count"]=float64(eligible)
    m[pfx+"representation_headroom"]=head[ti]

    // Freeze D and M from non-target training state before target bytes are opened.
    trainLocal:=map[string]float64{"counter_overflow_rows":0}
    baseTrain,rawSel:=wlmLmRawRepPredFreshHoldoutR1TrainModel([][]byte{a,b},trainLocal)
    m["counter_overflow_count"]+=trainLocal["counter_overflow_rows"]
    rows:=make([]wlmLmRawRepPredFreshHoldoutR1Motif,0,len(rawSel))
    for _,r:=range rawSel { rows=append(rows,r) }
    sort.Slice(rows,func(i,j int)bool{
      if rows[i].stats.total!=rows[j].stats.total { return rows[i].stats.total>rows[j].stats.total }
      return wlmLmRawRepPredFreshHoldoutR1Less(rows[i].key,rows[j].key)
    })
    if len(rows)<256 { m["invalid_row_count"]++ }
    if len(rows)>256 { rows=rows[:256] }
    classMass:=make([]float64,8)
    for i,r:=range rows { classMass[i/32]+=float64(r.stats.total) }
    minMass,maxMass:=classMass[0],classMass[0]
    for _,v:=range classMass[1:] { if v<minMass{minMass=v}; if v>maxMass{maxMass=v} }
    if maxMass<=0 { m["invalid_row_count"]++ } else { deficit[ti]=1-minMass/maxMass }
    m[pfx+"class_support_deficit"]=deficit[ti]

    margins:=make([]float64,0,256)
    for prev:=0;prev<256;prev++ {
      var total uint64
      var top,runner uint32
      for to:=0;to<256;to++ {
        v:=baseTrain[prev][to]; total+=uint64(v)
        if v>top { runner=top; top=v } else if v>runner { runner=v }
      }
      if total>=4 { margins=append(margins,(float64(top)-float64(runner))/float64(total)) }
    }
    sort.Float64s(margins)
    if len(margins)==0 { m["invalid_row_count"]++ } else if len(margins)%2==1 {
      marginMedian[ti]=margins[len(margins)/2]
    } else {
      n:=len(margins); marginMedian[ti]=(margins[n/2-1]+margins[n/2])/2
    }
    m[pfx+"relation_margin_median"]=marginMedian[ti]
    if len(t.data)<1454 { m["invalid_row_count"]++; continue }
    adapt:=t.data[:872]; eval:=t.data[872:1454]
    base,sel:=wlmLmR11Train256([][]byte{a,b},m)
    if len(sel)!=256 { m["invalid_row_count"]++ }
    h0:=wlmLmR11Hits(eval,&base,sel)
    m[pfx+"packet0_exact_hit_count"]=float64(h0)
    sumGain:=0.0
    h8:=h0
    for p:=1;p<=8;p++ {
      bb,ss2:=wlmLmR11Train256([][]byte{a,b,adapt[:p*109]},m)
      if len(ss2)!=256 { m["invalid_row_count"]++ }
      h:=wlmLmR11Hits(eval,&bb,ss2)
      m[pfx+"packet_"+string(rune('0'+p))+"_exact_hit_count"]=float64(h)
      gain:=float64(h-h0)
      sumGain+=gain
      if p==8 { h8=h }
    }
    cumulative[ti]=sumGain
    finalGain[ti]=float64(h8-h0)
    m[pfx+"cumulative_adaptation_gain"]=cumulative[ti]
    m[pfx+"final_gain"]=finalGain[ti]
  }
  type sig struct{name string; values []float64; higherPredictsHigher bool; code float64}
  sigs:=[]sig{
    {"H",head,true,1},
    {"D",deficit,false,2},
    {"M",marginMedian,true,3},
  }
  for _,sg:=range sigs {
    agree,rev,tie:=0.0,0.0,0.0
    for i:=0;i<3;i++ { for j:=i+1;j<3;j++ {
      if sg.values[i]==sg.values[j] || cumulative[i]==cumulative[j] { tie++; continue }
      same:=(sg.values[i]<sg.values[j] && cumulative[i]<cumulative[j]) || (sg.values[i]>sg.values[j] && cumulative[i]>cumulative[j])
      ok:=same
      if !sg.higherPredictsHigher { ok=!same }
      if ok { agree++ } else { rev++ }
    }}
    m[sg.name+"_pairwise_agreement_count"]=agree
    m[sg.name+"_pairwise_reversal_count"]=rev
    m[sg.name+"_pairwise_tie_count"]=tie
    if agree>m["best_signal_pairwise_agreement_count"] ||
      (agree==m["best_signal_pairwise_agreement_count"] && rev<m["best_signal_pairwise_reversal_count"]) ||
      (agree==m["best_signal_pairwise_agreement_count"] && rev==m["best_signal_pairwise_reversal_count"] && (m["best_signal_code"]==0 || sg.code<m["best_signal_code"])) {
      m["best_signal_code"]=sg.code
      m["best_signal_pairwise_agreement_count"]=agree
      m["best_signal_pairwise_reversal_count"]=rev
    }
  }
  for _,v:=range m { if math.IsNaN(v)||math.IsInf(v,0) { m["invalid_row_count"]++ } }
  return wlmLmFutureDataSignalFactorialR11Result{"wingless.research-scientific-result.v1","WLM-LM-EXTERNAL-FUTURE-DATA-LEARNING-EFFICIENCY-SIGNAL-FACTORIAL-R11",m}
}
