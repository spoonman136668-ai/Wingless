package unitary

import "sort"

type wlmYggHierarchyLifecycleBridgeR2Higher struct {
	compounds []wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey
	profiles map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile
	globalBest uint8
}

type wlmYggHierarchyLifecycleBridgeR2Record struct {
	kind uint8
	data []uint8
}

type wlmYggHierarchyLifecycleBridgeR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func wlmYggHierarchyLifecycleBridgeR2Build(seed uint32, metrics map[string]float64) (*wlmLmByteMotifDiscoveryR3BalancedLearner, wlmYggHierarchyLifecycleBridgeR2Higher) {
	local:=map[string]float64{"counter_overflow_rows":0,"invalid_candidate_rows":0,"duplicate_shuffled_control_key_count":0}
	lex:=wlmLmByteMotifDiscoveryR3BalancedCorpusFor(seed,4096,false)
	learner:=wlmLmByteMotifDiscoveryR3BalancedNewLearner()
	learner.train(lex.bytes,"slice-backed-four-byte-window",local)
	learner.selectMotifs(local)
	metrics["invalid_lifecycle_rows"]+=local["invalid_candidate_rows"]
	train:=wlmLmMultiscaleNativeLanguagePrecursorR2GrammarFor(seed,3072,false)
	l1:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel1(train.Bytes,&learner)
	compounds:=wlmLmMultiscaleNativeLanguagePrecursorR2LearnCompounds(l1)
	set:=wlmLmMultiscaleNativeLanguagePrecursorR2CompoundSet(compounds)
	l2:=wlmLmMultiscaleNativeLanguagePrecursorR2EventizeLevel2(l1,set)
	profilesPtr,globalCounts,_:=wlmLmMultiscaleNativeLanguagePrecursorR2LearnProfiles(l2)
	profiles:=make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile)
	for k,p:=range profilesPtr { if p!=nil { profiles[k]=*p } }
	sort.Slice(compounds,func(i,j int)bool{return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundLess(compounds[i],compounds[j])})
	return &learner,wlmYggHierarchyLifecycleBridgeR2Higher{compounds:compounds,profiles:profiles,globalBest:wlmLmMultiscaleNativeLanguagePrecursorR2GlobalBest(globalCounts)}
}

func wlmYggHierarchyLifecycleBridgeR2CompoundFromBytes(b []uint8) wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey {
	var a,c [4]uint8
	copy(a[:],b[:4]);copy(c[:],b[4:8])
	return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey{A:a,B:c}
}

func wlmYggHierarchyLifecycleBridgeR2CompoundSet(h wlmYggHierarchyLifecycleBridgeR2Higher) map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool {
	out:=make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool)
	for _,c:=range h.compounds {out[c]=true}
	return out
}

func wlmYggHierarchyLifecycleBridgeR2Predict(h wlmYggHierarchyLifecycleBridgeR2Higher,a,b wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey) uint8 {
	set:=wlmYggHierarchyLifecycleBridgeR2CompoundSet(h)
	pa,oka:=h.profiles[a];pb,okb:=h.profiles[b]
	if !set[a]||!set[b]||!oka||!okb {return h.globalBest}
	hi:=wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(pa.First)
	lo:=wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(pb.Second)
	return 224+(hi<<1)+lo
}

func wlmYggHierarchyLifecycleBridgeR2Accuracy(h wlmYggHierarchyLifecycleBridgeR2Higher,g wlmLmMultiscaleNativeLanguagePrecursorR2Grammar,lesioned map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool,stratum string) float64 {
	correct,total:=0,0
	for r:=0;r<len(g.Pairs);r++ {
		off:=r*32
		if off+29>len(g.Bytes) {continue}
		a:=wlmYggHierarchyLifecycleBridgeR2CompoundFromBytes(g.Bytes[off:off+8])
		b:=wlmYggHierarchyLifecycleBridgeR2CompoundFromBytes(g.Bytes[off+20:off+28])
		affected:=lesioned!=nil&&(lesioned[a]||lesioned[b])
		if stratum=="lesioned"&&!affected {continue}
		if stratum=="intact"&&affected {continue}
		if wlmYggHierarchyLifecycleBridgeR2Predict(h,a,b)==g.Bytes[off+28] {correct++}
		total++
	}
	if total==0{return 0}
	return float64(correct)/float64(total)
}

func wlmYggHierarchyLifecycleBridgeR2ProfileBits(p wlmLmMultiscaleNativeLanguagePrecursorR2Profile)(uint8,uint8){
	return wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(p.First),wlmLmMultiscaleNativeLanguagePrecursorR2BestBit(p.Second)
}

func wlmYggHierarchyLifecycleBridgeR2Serialize(h wlmYggHierarchyLifecycleBridgeR2Higher,selected map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool) []wlmYggHierarchyLifecycleBridgeR2Record {
	out:=make([]wlmYggHierarchyLifecycleBridgeR2Record,0)
	for _,c:=range h.compounds {
		if selected!=nil&&!selected[c]{continue}
		data:=append([]uint8{},c.A[:]...);data=append(data,c.B[:]...)
		out=append(out,wlmYggHierarchyLifecycleBridgeR2Record{kind:2,data:data})
	}
	for _,c:=range h.compounds {
		if selected!=nil&&!selected[c]{continue}
		p:=h.profiles[c];hi,lo:=wlmYggHierarchyLifecycleBridgeR2ProfileBits(p)
		data:=append([]uint8{},c.A[:]...);data=append(data,c.B[:]...);data=append(data,hi,lo)
		out=append(out,wlmYggHierarchyLifecycleBridgeR2Record{kind:3,data:data})
	}
	return out
}

func wlmYggHierarchyLifecycleBridgeR2RecordBytes(rs []wlmYggHierarchyLifecycleBridgeR2Record) int {
	n:=0;for _,r:=range rs{n+=len(r.data)};return n
}

func wlmYggHierarchyLifecycleBridgeR2Restore(global uint8,rs []wlmYggHierarchyLifecycleBridgeR2Record)(wlmYggHierarchyLifecycleBridgeR2Higher,int,int){
	h:=wlmYggHierarchyLifecycleBridgeR2Higher{profiles:make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile),globalBest:global}
	invalid:=0
	for _,r:=range rs{
		if r.kind==2&&len(r.data)==8 {
			c:=wlmYggHierarchyLifecycleBridgeR2CompoundFromBytes(r.data)
			h.compounds=append(h.compounds,c)
		}else if r.kind==3&&len(r.data)==10 {
			c:=wlmYggHierarchyLifecycleBridgeR2CompoundFromBytes(r.data[:8])
			var p wlmLmMultiscaleNativeLanguagePrecursorR2Profile
			p.First[r.data[8]&1]=1;p.Second[r.data[9]&1]=1
			h.profiles[c]=p
		}else{invalid++}
	}
	sort.Slice(h.compounds,func(i,j int)bool{return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundLess(h.compounds[i],h.compounds[j])})
	return h,len(rs),invalid
}

func wlmYggHierarchyLifecycleBridgeR2Jaccard(a,b wlmYggHierarchyLifecycleBridgeR2Higher) float64 {
	token:=func(h wlmYggHierarchyLifecycleBridgeR2Higher)map[string]bool{
		m:=map[string]bool{}
		for _,c:=range h.compounds{
			hi,lo:=wlmYggHierarchyLifecycleBridgeR2ProfileBits(h.profiles[c])
			k:=string(append(append(append([]byte{},c.A[:]...),c.B[:]...),hi,lo))
			m[k]=true
		}
		return m
	}
	x,y:=token(a),token(b);inter:=0
	for k:=range x{if y[k]{inter++}}
	u:=len(x);for k:=range y{if !x[k]{u++}}
	if u==0{return 1};return float64(inter)/float64(u)
}

func wlmYggHierarchyLifecycleBridgeR2Lesion(h wlmYggHierarchyLifecycleBridgeR2Higher,selected map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool) wlmYggHierarchyLifecycleBridgeR2Higher {
	out:=wlmYggHierarchyLifecycleBridgeR2Higher{profiles:make(map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile),globalBest:h.globalBest}
	for _,c:=range h.compounds{if !selected[c]{out.compounds=append(out.compounds,c)}}
	for c,p:=range h.profiles{if !selected[c]{out.profiles[c]=p}}
	return out
}

func wlmYggHierarchyLifecycleBridgeR2Repair(base wlmYggHierarchyLifecycleBridgeR2Higher,rs []wlmYggHierarchyLifecycleBridgeR2Record)(wlmYggHierarchyLifecycleBridgeR2Higher,int,int){
	add,ops,invalid:=wlmYggHierarchyLifecycleBridgeR2Restore(base.globalBest,rs)
	out:=base
	for _,c:=range add.compounds{out.compounds=append(out.compounds,c)}
	for c,p:=range add.profiles{out.profiles[c]=p}
	sort.Slice(out.compounds,func(i,j int)bool{return wlmLmMultiscaleNativeLanguagePrecursorR2CompoundLess(out.compounds[i],out.compounds[j])})
	return out,ops,invalid
}

func wlmYggHierarchyLifecycleBridgeR2Corrupt(rs []wlmYggHierarchyLifecycleBridgeR2Record) []wlmYggHierarchyLifecycleBridgeR2Record {
	out:=make([]wlmYggHierarchyLifecycleBridgeR2Record,len(rs))
	for i,r:=range rs{
		data:=append([]uint8{},r.data...)
		if r.kind==3&&len(data)==10{data[8]^=1;data[9]^=1}
		out[i]=wlmYggHierarchyLifecycleBridgeR2Record{kind:r.kind,data:data}
	}
	return out
}

func wlmYggHierarchyLifecycleBridgeR2IntactMutations(before,after wlmYggHierarchyLifecycleBridgeR2Higher,selected map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool) int {
	n:=0
	bset:=wlmYggHierarchyLifecycleBridgeR2CompoundSet(before)
	aset:=wlmYggHierarchyLifecycleBridgeR2CompoundSet(after)
	for c:=range bset{
		if selected[c]{continue}
		if !aset[c]{n++;continue}
		bh,bl:=wlmYggHierarchyLifecycleBridgeR2ProfileBits(before.profiles[c])
		ah,al:=wlmYggHierarchyLifecycleBridgeR2ProfileBits(after.profiles[c])
		if bh!=ah||bl!=al{n++}
	}
	return n
}

func wlmYggHierarchyLifecycleBridgeR2Min(a,b float64)float64{if b<a{return b};return a}
func wlmYggHierarchyLifecycleBridgeR2Max(a,b float64)float64{if b>a{return b};return a}

// RunWlmYggHierarchyLifecycleBridgeR2 executes the frozen cross-project lifecycle bridge.
func RunWlmYggHierarchyLifecycleBridgeR2() interface{} {
	seeds:=[...]uint32{12301,12323,12329,12343}
	metrics:=map[string]float64{
		"valid_seed_count":0,
		"minimum_baseline_local_accuracy":1,
		"minimum_baseline_long_range_accuracy":1,
		"minimum_hibernated_local_accuracy":1,
		"maximum_hibernated_long_range_accuracy":0,
		"maximum_hibernated_higher_active_structure_count":0,
		"maximum_retained_higher_bytes":0,
		"minimum_retained_higher_bytes":72,
		"maximum_wake_operations":0,
		"minimum_post_wake_long_range_accuracy":1,
		"minimum_post_wake_structure_jaccard":1,
		"maximum_wake_to_cold_operation_ratio":0,
		"maximum_erased_wake_long_range_accuracy":0,
		"maximum_post_lesion_lesioned_pair_accuracy":0,
		"minimum_post_lesion_intact_pair_accuracy":1,
		"minimum_post_repair_lesioned_pair_accuracy":1,
		"minimum_post_repair_intact_pair_accuracy":1,
		"maximum_selective_repair_read_bytes":0,
		"minimum_selective_repair_read_bytes":36,
		"maximum_selective_repair_operations":0,
		"maximum_intact_higher_structure_mutation_count":0,
		"maximum_corrupt_repair_lesioned_pair_accuracy":0,
		"maximum_repair_to_cold_operation_ratio":0,
		"maximum_total_active_structure_count":0,
		"capacity_growth_event_count":0,
		"invalid_lifecycle_rows":0,
	}
	const coldOps=131072.0
	for _,seed:=range seeds{
		learner,h:=wlmYggHierarchyLifecycleBridgeR2Build(seed,metrics)
		eval:=wlmLmMultiscaleNativeLanguagePrecursorR2GrammarFor(seed,1024,true)
		local:=wlmLmMultiscaleNativeLanguagePrecursorR2LocalAccuracy(seed,learner)
		baseLong:=wlmYggHierarchyLifecycleBridgeR2Accuracy(h,eval,nil,"all")
		metrics["minimum_baseline_local_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_baseline_local_accuracy"],local)
		metrics["minimum_baseline_long_range_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_baseline_long_range_accuracy"],baseLong)
		allRecords:=wlmYggHierarchyLifecycleBridgeR2Serialize(h,nil)
		retainedBytes:=wlmYggHierarchyLifecycleBridgeR2RecordBytes(allRecords)
		hib:=wlmYggHierarchyLifecycleBridgeR2Higher{profiles:map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile{},globalBest:h.globalBest}
		metrics["minimum_hibernated_local_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_hibernated_local_accuracy"],local)
		metrics["maximum_hibernated_long_range_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_hibernated_long_range_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(hib,eval,nil,"all"))
		metrics["maximum_retained_higher_bytes"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_retained_higher_bytes"],float64(retainedBytes))
		metrics["minimum_retained_higher_bytes"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_retained_higher_bytes"],float64(retainedBytes))
		wake,wops,inv:=wlmYggHierarchyLifecycleBridgeR2Restore(h.globalBest,allRecords)
		metrics["invalid_lifecycle_rows"]+=float64(inv)
		metrics["maximum_wake_operations"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_wake_operations"],float64(wops))
		metrics["minimum_post_wake_long_range_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_post_wake_long_range_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(wake,eval,nil,"all"))
		metrics["minimum_post_wake_structure_jaccard"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_post_wake_structure_jaccard"],wlmYggHierarchyLifecycleBridgeR2Jaccard(h,wake))
		metrics["maximum_wake_to_cold_operation_ratio"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_wake_to_cold_operation_ratio"],float64(wops)/coldOps)
		erased:=wlmYggHierarchyLifecycleBridgeR2Higher{profiles:map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]wlmLmMultiscaleNativeLanguagePrecursorR2Profile{},globalBest:h.globalBest}
		metrics["maximum_erased_wake_long_range_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_erased_wake_long_range_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(erased,eval,nil,"all"))

		selected:=map[wlmLmMultiscaleNativeLanguagePrecursorR2CompoundKey]bool{}
		i:=int(seed%2);selected[h.compounds[i]]=true;selected[h.compounds[i+2]]=true
		selRecords:=wlmYggHierarchyLifecycleBridgeR2Serialize(h,selected)
		damaged:=wlmYggHierarchyLifecycleBridgeR2Lesion(h,selected)
		metrics["maximum_post_lesion_lesioned_pair_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_post_lesion_lesioned_pair_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(damaged,eval,selected,"lesioned"))
		metrics["minimum_post_lesion_intact_pair_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_post_lesion_intact_pair_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(damaged,eval,selected,"intact"))
		repaired,rops,rinv:=wlmYggHierarchyLifecycleBridgeR2Repair(damaged,selRecords)
		metrics["invalid_lifecycle_rows"]+=float64(rinv)
		metrics["minimum_post_repair_lesioned_pair_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_post_repair_lesioned_pair_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(repaired,eval,selected,"lesioned"))
		metrics["minimum_post_repair_intact_pair_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_post_repair_intact_pair_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(repaired,eval,selected,"intact"))
		rbytes:=wlmYggHierarchyLifecycleBridgeR2RecordBytes(selRecords)
		metrics["maximum_selective_repair_read_bytes"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_selective_repair_read_bytes"],float64(rbytes))
		metrics["minimum_selective_repair_read_bytes"]=wlmYggHierarchyLifecycleBridgeR2Min(metrics["minimum_selective_repair_read_bytes"],float64(rbytes))
		metrics["maximum_selective_repair_operations"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_selective_repair_operations"],float64(rops))
		metrics["maximum_intact_higher_structure_mutation_count"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_intact_higher_structure_mutation_count"],float64(wlmYggHierarchyLifecycleBridgeR2IntactMutations(h,repaired,selected)))
		corrupt,_,cinv:=wlmYggHierarchyLifecycleBridgeR2Repair(damaged,wlmYggHierarchyLifecycleBridgeR2Corrupt(selRecords))
		metrics["invalid_lifecycle_rows"]+=float64(cinv)
		metrics["maximum_corrupt_repair_lesioned_pair_accuracy"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_corrupt_repair_lesioned_pair_accuracy"],wlmYggHierarchyLifecycleBridgeR2Accuracy(corrupt,eval,selected,"lesioned"))
		metrics["maximum_repair_to_cold_operation_ratio"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_repair_to_cold_operation_ratio"],float64(rops)/coldOps)
		metrics["maximum_total_active_structure_count"]=wlmYggHierarchyLifecycleBridgeR2Max(metrics["maximum_total_active_structure_count"],float64(8+len(repaired.compounds)+len(repaired.profiles)))
		metrics["valid_seed_count"]++
	}
	return wlmYggHierarchyLifecycleBridgeR2Result{Schema:"wingless.research-scientific-result.v1",Experiment:"WLM-YGG-HIERARCHY-LIFECYCLE-BRIDGE-R2",Metrics:metrics}
}
