package unitary

import "sort"

type wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey struct {
	A [4]uint8
	B [4]uint8
}

type wlmLmMultiscaleNativeLanguagePrecursorR2Event struct {
	Kind     uint8
	Raw      uint8
	Motif    [4]uint8
	Compound wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey
}

type wlmLmMultiscaleNativeLanguagePrecursorR2Profile struct {
	First  [2]uint32
	Second [2]uint32
}

type wlmLmMultiscaleNativeLanguagePrecursorR2Grammar struct {
	Bytes []uint8
	Pairs [][2]int
}

type wlmLmMultiscaleNativeLanguagePrecursorR2Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmMultiscaleNativeLanguagePrecursorR2Compound(id int) wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey {
	return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{
		A: wlmLmByteMotifDiscoveryR3BalancedMotif(id),
		B: wlmLmByteMotifDiscoveryR3BalancedMotif(id + 4),
	}
}

func wlmLmMultiscaleNativeLanguagePrecursorR2Step(state uint32) uint32 {
	return state*1664525 + 1013904223
}

func wlmLmMultiscaleNativeLanguagePrecursorR2Heldout(seed uint32, a, b int) bool {
	return (a+b+int(seed%4))%4 == 1
}

func wlmLmMultiscaleNativeLanguagePrecursorR2GrammarFor(seed uint32, records int, evaluation bool) wlmLmMultiscaleNativeLanguagePrecursorR2Grammar {
	pairs := make([][2]int, 0, 16)
	for a:=0;a<4;a++ {
		for b:=0;b<4;b++ {
			h:=wlmLmMultiscaleNativeLanguagePrecursorR2Heldout(seed,a,b)
			if evaluation == h {
				pairs=append(pairs,[2]int{a,b})
			}
		}
	}
	out:=wlmLmMultiscaleNativeLanguagePrecursorR2Grammar{
		Bytes:make([]uint8,0,records*32),
		Pairs:make([][2]int,0,records),
	}
	state:=seed^0x3c6ef372
	if evaluation { state=seed^0xa54ff53a }
	offset:=(int(seed)+boolToInt(evaluation))%len(pairs)
	for record:=0;record<records;record++ {
		pair:=pairs[(5*record+offset)%len(pairs)]
		a,b:=pair[0],pair[1]
		ca:=wlmLmMultiscaleNativeLanguagePrecursorR2Compound(a)
		cb:=wlmLmMultiscaleNativeLanguagePrecursorR2Compound(b)
		out.Pairs=append(out.Pairs,pair)
		out.Bytes=append(out.Bytes,ca.A[:]...)
		out.Bytes=append(out.Bytes,ca.B[:]...)
		for i:=0;i<12;i++ {
			state=wlmLmMultiscaleNativeLanguagePrecursorR2Step(state)
			out.Bytes=append(out.Bytes,uint8(192+(state%64)))
		}
		out.Bytes=append(out.Bytes,cb.A[:]...)
		out.Bytes=append(out.Bytes,cb.B[:]...)
		ref:=uint8(224+((a&1)<<1)+(b&1))
		out.Bytes=append(out.Bytes,ref)
		for i:=0;i<3;i++ {
			state=wlmLmMultiscaleNativeLanguagePrecursorR2Step(state)
			out.Bytes=append(out.Bytes,uint8(192+(state%64)))
		}
	}
	return out
}

func boolToInt(v bool) int { if v { return 1 }; return 0 }

func wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel1(
	bytes []uint8,
	learner *wlmLmByteMotifDiscoveryR3BalancedLearner,
) []wlmLmMultiscaleNativeLanguagePrecursorR2Event {
	events:=make([]wlmLmMultiscaleNativeLanguagePrecursorR2Event,0,len(bytes))
	for i:=0;i<len(bytes); {
		if i+4<=len(bytes) {
			key:=[4]uint8{bytes[i],bytes[i+1],bytes[i+2],bytes[i+3]}
			if _,ok:=learner.selectedMap[key];ok {
				events=append(events,wlmLmMultiscaleNativeLanguagePrecursorR2Event{Kind:1,Motif:key})
				i+=4
				continue
			}
		}
		events=append(events,wlmLmMultiscaleNativeLanguagePrecursorR2Event{Kind:0,Raw:bytes[i]})
		i++
	}
	return events
}

func wlmLmMultiscaleNativeLanguagePrecursorR2CompoundLess(a,b wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey) bool {
	for i:=0;i<4;i++ { if a.A[i]!=b.A[i] { return a.A[i]<b.A[i] } }
	for i:=0;i<4;i++ { if a.B[i]!=b.B[i] { return a.B[i]<b.B[i] } }
	return false
}

func wlmLmMultiscaleNativeLanguagePrecursorR2LearnCompounds(events []wlmLmMultiscaleNativeLanguagePrecursorR2Event) []wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey {
	counts:=make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]uint32)
	for i:=0;i+1<len(events);i++ {
		if events[i].Kind==1 && events[i+1].Kind==1 {
			key:=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{A:events[i].Motif,B:events[i+1].Motif}
			counts[key]++
		}
	}
	type row struct{ key wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey; count uint32 }
	rows:=make([]row,0)
	for key,count:=range counts {
		if count>=64 { rows=append(rows,row{key:key,count:count}) }
	}
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].count!=rows[j].count { return rows[i].count>rows[j].count }
		return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundLess(rows[i].key,rows[j].key)
	})
	if len(rows)>4 { rows=rows[:4] }
	out:=make([]wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey,len(rows))
	for i,row:=range rows { out[i]=row.key }
	return out
}

func wlmLmMultiscaleNativeLanguagePrecursorR2CompoundSet(compounds []wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey) map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool {
	out:=make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool)
	for _,key:=range compounds { out[key]=true }
	return out
}

func wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel2(
	level1 []wlmLmMultiscaleNativeLanguagePrecursorR2Event,
	compounds map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool,
) []wlmLmMultiscaleNativeLanguagePrecursorR2Event {
	out:=make([]wlmLmMultiscaleNativeLanguagePrecursorR2Event,0,len(level1))
	for i:=0;i<len(level1); {
		if i+1<len(level1) && level1[i].Kind==1 && level1[i+1].Kind==1 {
			key:=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{A:level1[i].Motif,B:level1[i+1].Motif}
			if compounds[key] {
				out=append(out,wlmLmMultiscaleNativeLanguagePrecursorR2Event{Kind:2,Compound:key})
				i+=2
				continue
			}
		}
		out=append(out,level1[i])
		i++
	}
	return out
}

func wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(counts [2]uint32) uint8 {
	if counts[1]>counts[0] { return 1 }
	return 0
}

func wlmLmMultiscaleNativeLanguagePrecursorR2LearnProfiles(
	events []wlmLmMultiscaleNativeLanguagePrecursorR2Event,
) (map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]*wlmLmMultiscaleNativeLanguagePrecursorR2Profile,map[uint8]uint32,int) {
	profiles:=make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]*wlmLmMultiscaleNativeLanguagePrecursorR2Profile)
	global:=make(map[uint8]uint32)
	patterns:=0
	for i:=0;i<len(events);i++ {
		if events[i].Kind!=2 { continue }
		j:=i+1
		rawCount:=0
		for j<len(events) && events[j].Kind==0 {
			rawCount++; j++
		}
		if rawCount!=12 || j>=len(events) || events[j].Kind!=2 { continue }
		k:=j+1
		if k>=len(events) || events[k].Kind!=0 { continue }
		target:=events[k].Raw
		if target<224 || target>227 { continue }
		class:=target-224
		first:=profiles[events[i].Compound]
		if first==nil { first=&wlmLmMultiscaleNativeLanguagePrecursorR2Profile{};profiles[events[i].Compound]=first }
		second:=profiles[events[j].Compound]
		if second==nil { second=&wlmLmMultiscaleNativeLanguagePrecursorR2Profile{};profiles[events[j].Compound]=second }
		first.First[(class>>1)&1]++
		second.Second[class&1]++
		global[target]++
		patterns++
	}
	return profiles,global,patterns
}

func wlmLmMultiscaleNativeLanguagePrecursorR2GlobalBest(counts map[uint8]uint32) uint8 {
	best:=uint8(0); bestCount:=uint32(0); first:=true
	for value,count:=range counts {
		if first || count>bestCount || (count==bestCount && value<best) {
			best=value;bestCount=count;first=false
		}
	}
	return best
}

func wlmLmMultiscaleNativeLanguagePrecursorR2PredictReference(
	profiles map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]*wlmLmMultiscaleNativeLanguagePrecursorR2Profile,
	a,b wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey,
) uint8 {
	pa,pb:=profiles[a],profiles[b]
	if pa==nil || pb==nil { return 0 }
	high:=wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(pa.First)
	low:=wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(pb.Second)
	return 224+(high<<1)+low
}

func wlmLmMultiscaleNativeLanguagePrecursorR2ReferenceAccuracy(
	events []wlmLmMultiscaleNativeLanguagePrecursorR2Event,
	profiles map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]*wlmLmMultiscaleNativeLanguagePrecursorR2Profile,
	globalBest uint8,
	shuffled bool,
) (float64,float64,int) {
	hierCorrect:=0
	controlCorrect:=0
	total:=0
	for i:=0;i<len(events);i++ {
		if events[i].Kind!=2 { continue }
		j:=i+1; rawCount:=0
		for j<len(events)&&events[j].Kind==0 { rawCount++;j++ }
		if rawCount!=12 || j>=len(events) || events[j].Kind!=2 { continue }
		k:=j+1
		if k>=len(events)||events[k].Kind!=0 { continue }
		target:=events[k].Raw
		if target<224||target>227 { continue }
		a,b:=events[i].Compound,events[j].Compound
		if shuffled {
			a=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{A:a.B,B:a.A}
			b=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{A:b.B,B:b.A}
		}
		if wlmLmMultiscaleNativeLanguagePrecursorR2PredictReference(profiles,a,b)==target { hierCorrect++ }
		if globalBest==target { controlCorrect++ }
		total++
	}
	if total==0 { return 0,0,0 }
	return float64(hierCorrect)/float64(total),float64(controlCorrect)/float64(total),total
}

func wlmLmMultiscaleNativeLanguagePrecursorR2MotifRecall(learner *wlmLmByteMotifDiscoveryR3BalancedLearner) float64 {
	count:=0
	for id:=0;id<8;id++ {
		if _,ok:=learner.selectedMap[wlmLmByteMotifDiscoveryR3BalancedMotif(id)];ok { count++ }
	}
	return float64(count)/8.0
}

func wlmLmMultiscaleNativeLanguagePrecursorR2CompoundRecall(compounds []wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey) float64 {
	set:=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundSet(compounds)
	count:=0
	for id:=0;id<4;id++ { if set[wlmLmMultiscaleNativeLanguagePrecursorR2Compound(id)] { count++ } }
	return float64(count)/4.0
}

func wlmLmMultiscaleNativeLanguagePrecursorR2LocalAccuracy(seed uint32,learner *wlmLmByteMotifDiscoveryR3BalancedLearner) float64 {
	eval:=wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed,1024,true)
	correct:=0
	total:=0
	for _,pos:=range eval.targetPositions {
		if pos<4||pos>=len(eval.bytes){continue}
		key:=[4]uint8{eval.bytes[pos-4],eval.bytes[pos-3],eval.bytes[pos-2],eval.bytes[pos-1]}
		if learner.predict(key,"motif")==eval.bytes[pos] { correct++ }
		total++
	}
	if total==0{return 0}
	return float64(correct)/float64(total)
}

func wlmLmMultiscaleNativeLanguagePrecursorR2Min(a,b float64)float64{if b<a{return b};return a}
func wlmLmMultiscaleNativeLanguagePrecursorR2Max(a,b float64)float64{if b>a{return b};return a}

// RunWlmLmMultiscaleNativeLanguagePrecursorR2 evaluates the frozen WBG-5 language-like hierarchy.
func RunWlmLmMultiscaleNativeLanguagePrecursorR2() interface{} {
	seeds:=[...]uint32{12203,12211,12227,12239}
	metrics:=map[string]float64{
		"valid_seed_count":0,
		"minimum_level1_motif_recall":1,
		"minimum_level2_compound_recall":1,
		"minimum_local_successor_accuracy":1,
		"minimum_long_range_reference_accuracy":1,
		"maximum_motif_only_long_range_reference_accuracy":0,
		"maximum_shuffled_compound_long_range_accuracy":0,
		"minimum_hierarchical_long_range_gain":1,
		"minimum_effective_event_reduction_fraction":1,
		"maximum_total_learned_structure_count":0,
		"minimum_reference_profile_coverage":4,
		"maximum_heldout_pair_leak_count":0,
		"oracle_boundary_use_count":0,
		"capacity_growth_event_count":0,
		"invalid_stream_rows":0,
		"counter_overflow_rows":0,
	}
	for _,seed:=range seeds {
		local:=map[string]float64{"counter_overflow_rows":0}
		lex:=wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed,4096,false)
		learner:=wlmLmByteMotifDiscoveryR3BalancedNewLearner()
		learner.train(lex.bytes,"slice-backed-four-byte-window",local)
		learner.selectMotifs(local)
		metrics["counter_overflow_rows"]+=local["counter_overflow_rows"]
		metrics["minimum_level1_motif_recall"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_level1_motif_recall"],wlmLmMultiscaleNativeLanguagePrecursorR2MotifRecall(&learner))
		metrics["minimum_local_successor_accuracy"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_local_successor_accuracy"],wlmLmMultiscaleNativeLanguagePrecursorR2LocalAccuracy(seed,&learner))

		train:=wlmLmMultiscaleNativeLanguagePrecursorR2GrammarFor(seed,3072,false)
		for _,pair:=range train.Pairs {
			if wlmLmMultiscaleNativeLanguagePrecursorR2Heldout(seed,pair[0],pair[1]) {
				metrics["maximum_heldout_pair_leak_count"]++
			}
		}
		l1Train:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel1(train.Bytes,&learner)
		compounds:=wlmLmMultiscaleNativeLanguagePrecursorR2LearnCompounds(l1Train)
		metrics["minimum_level2_compound_recall"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_level2_compound_recall"],wlmLmMultiscaleNativeLanguagePrecursorR2CompoundRecall(compounds))
		compoundSet:=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundSet(compounds)
		l2Train:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel2(l1Train,compoundSet)
		profiles,globalCounts,_:=wlmLmMultiscaleNativeLanguagePrecursorR2LearnProfiles(l2Train)
		globalBest:=wlmLmMultiscaleNativeLanguagePrecursorR2GlobalBest(globalCounts)
		coverage:=float64(len(profiles))
		metrics["minimum_reference_profile_coverage"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_reference_profile_coverage"],coverage)

		eval:=wlmLmMultiscaleNativeLanguagePrecursorR2GrammarFor(seed,1024,true)
		l1Eval:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel1(eval.Bytes,&learner)
		l2Eval:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel2(l1Eval,compoundSet)
		hier,control,total:=wlmLmMultiscaleNativeLanguagePrecursorR2ReferenceAccuracy(l2Eval,profiles,globalBest,false)
		shuffled,_,_:=wlmLmMultiscaleNativeLanguagePrecursorR2ReferenceAccuracy(l2Eval,profiles,globalBest,true)
		if total!=1024 { metrics["invalid_stream_rows"]+=float64(absInt(total-1024)) }
		metrics["minimum_long_range_reference_accuracy"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_long_range_reference_accuracy"],hier)
		metrics["maximum_motif_only_long_range_reference_accuracy"]=wlmLmMultiscaleNativeLanguagePrecursorR2Max(metrics["maximum_motif_only_long_range_reference_accuracy"],control)
		metrics["maximum_shuffled_compound_long_range_accuracy"]=wlmLmMultiscaleNativeLanguagePrecursorR2Max(metrics["maximum_shuffled_compound_long_range_accuracy"],shuffled)
		metrics["minimum_hierarchical_long_range_gain"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_hierarchical_long_range_gain"],hier-control)
		reduction:=1-float64(len(l2Eval))/float64(len(eval.Bytes))
		metrics["minimum_effective_event_reduction_fraction"]=wlmLmMultiscaleNativeLanguagePrecursorR2Min(metrics["minimum_effective_event_reduction_fraction"],reduction)
		structures:=float64(len(learner.selected)+len(compounds)+len(profiles))
		metrics["maximum_total_learned_structure_count"]=wlmLmMultiscaleNativeLanguagePrecursorR2Max(metrics["maximum_total_learned_structure_count"],structures)
		metrics["valid_seed_count"]++
	}
	return wlmLmMultiscaleNativeLanguagePrecursorR2Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-LM-MULTISCALE-NATIVE-LANGUAGE-PRECURSOR-R2",Metrics:metrics}
}

func absInt(v int) int { if v<0{return -v};return v }
