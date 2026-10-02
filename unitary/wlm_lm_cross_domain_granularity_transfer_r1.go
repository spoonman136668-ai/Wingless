package unitary

import "sort"

type wlmLmCrossDomainGranularityTransferR1Stats struct {
	total  uint32
	counts map[uint8]uint32
}

func (s *wlmLmCrossDomainGranularityTransferR1Stats) observe(next uint8, metrics map[string]float64) {
	if s.counts == nil {
		s.counts = make(map[uint8]uint32)
	}
	if s.total == ^uint32(0) || s.counts[next] == ^uint32(0) {
		metrics["counter_overflow_rows"]++
		return
	}
	s.total++
	s.counts[next]++
}

func (s *wlmLmCrossDomainGranularityTransferR1Stats) best() (uint8, uint32) {
	var bestByte uint8
	var bestCount uint32
	first := true
	for b, count := range s.counts {
		if first || count > bestCount || (count == bestCount && b < bestByte) {
			bestByte = b
			bestCount = count
			first = false
		}
	}
	return bestByte, bestCount
}

type wlmLmCrossDomainGranularityTransferR1Entry struct {
	key       [4]uint8
	successor uint8
	total     uint32
	bestCount uint32
}

type wlmLmCrossDomainGranularityTransferR1Store struct {
	entries map[[4]uint8]uint8
	order   [][4]uint8
}

func wlmLmCrossDomainGranularityTransferR1NewStore() *wlmLmCrossDomainGranularityTransferR1Store {
	return &wlmLmCrossDomainGranularityTransferR1Store{
		entries: make(map[[4]uint8]uint8),
		order:   make([][4]uint8, 0, 12),
	}
}

func (s *wlmLmCrossDomainGranularityTransferR1Store) add(key [4]uint8, successor uint8, metrics map[string]float64) bool {
	if existing, ok := s.entries[key]; ok {
		if existing != successor {
			metrics["motif_rewrite_count"]++
		}
		return false
	}
	if len(s.order) >= 12 {
		return false
	}
	s.entries[key] = successor
	s.order = append(s.order, key)
	return true
}

func (s *wlmLmCrossDomainGranularityTransferR1Store) predict(key [4]uint8) uint8 {
	if v, ok := s.entries[key]; ok {
		return v
	}
	return 0
}

type wlmLmCrossDomainGranularityTransferR1Corpus struct {
	bytes   []uint8
	targets []int
	starts  []int
}

type wlmLmCrossDomainGranularityTransferR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmLmCrossDomainGranularityTransferR1SharedMotif(id int) [4]uint8 {
	return wlmLmByteMotifDiscoveryR3BalancedMotif(id)
}

func wlmLmCrossDomainGranularityTransferR1UniqueMotif(domain, local int) [4]uint8 {
	return [4]uint8{
		uint8(32 + 16*domain + local),
		uint8(96 + 5*domain + 3*local),
		204,
		51,
	}
}

func wlmLmCrossDomainGranularityTransferR1Successor(domain, motifIndex int) uint8 {
	if motifIndex < 4 {
		return wlmLmByteMotifDiscoveryR3BalancedSuccessor(motifIndex)
	}
	return uint8(128 + 8*domain + 3*(motifIndex-4))
}

func wlmLmCrossDomainGranularityTransferR1Step(state uint32) uint32 {
	return state*1664525 + 1013904223
}

func wlmLmCrossDomainGranularityTransferR1CorpusFor(seed uint32, domain, records, phase int) wlmLmCrossDomainGranularityTransferR1Corpus {
	out := wlmLmCrossDomainGranularityTransferR1Corpus{
		bytes:   make([]uint8, 0, records*8),
		targets: make([]int, 0, records),
		starts:  make([]int, 0, records),
	}
	state := seed ^ uint32(0x6a09e667+domain*0x10001+phase*0x010101)
	offset := (int(seed) + domain + phase) % 6
	fillerBase := [...]uint8{160,176,208,224}[domain]
	for record := 0; record < records; record++ {
		motifIndex := (5*record + offset) % 6
		var key [4]uint8
		if motifIndex < 4 {
			key = wlmLmCrossDomainGranularityTransferR1SharedMotif(motifIndex)
		} else {
			key = wlmLmCrossDomainGranularityTransferR1UniqueMotif(domain, motifIndex-4)
		}
		out.starts = append(out.starts, len(out.bytes))
		out.bytes = append(out.bytes, key[:]...)
		out.targets = append(out.targets, len(out.bytes))
		out.bytes = append(out.bytes, wlmLmCrossDomainGranularityTransferR1Successor(domain, motifIndex))
		for filler := 0; filler < 3; filler++ {
			state = wlmLmCrossDomainGranularityTransferR1Step(state)
			out.bytes = append(out.bytes, fillerBase+uint8((state>>24)&15))
		}
	}
	return out
}

func wlmLmCrossDomainGranularityTransferR1KeyLess(a, b [4]uint8) bool {
	for i := 0; i < 4; i++ {
		if a[i] < b[i] { return true }
		if a[i] > b[i] { return false }
	}
	return false
}

func wlmLmCrossDomainGranularityTransferR1Mine(
	corpus wlmLmCrossDomainGranularityTransferR1Corpus,
	existing *wlmLmCrossDomainGranularityTransferR1Store,
	metrics map[string]float64,
) []wlmLmCrossDomainGranularityTransferR1Entry {
	stats := make(map[[4]uint8]*wlmLmCrossDomainGranularityTransferR1Stats)
	for _, start := range corpus.starts {
		if start+4 >= len(corpus.bytes) {
			metrics["invalid_stream_rows"]++
			continue
		}
		key := [4]uint8{corpus.bytes[start], corpus.bytes[start+1], corpus.bytes[start+2], corpus.bytes[start+3]}
		next := corpus.bytes[start+4]
		s := stats[key]
		if s == nil {
			s = &wlmLmCrossDomainGranularityTransferR1Stats{}
			stats[key] = s
		}
		s.observe(next, metrics)
	}
	candidates := make([]wlmLmCrossDomainGranularityTransferR1Entry, 0)
	for key, s := range stats {
		if _, already := existing.entries[key]; already || s.total < 64 {
			continue
		}
		best, bestCount := s.best()
		if float64(bestCount)/float64(s.total) < 0.90 {
			continue
		}
		candidates = append(candidates, wlmLmCrossDomainGranularityTransferR1Entry{
			key:key, successor:best, total:s.total, bestCount:bestCount,
		})
	}
	sort.Slice(candidates, func(i,j int) bool {
		if candidates[i].bestCount != candidates[j].bestCount {
			return candidates[i].bestCount > candidates[j].bestCount
		}
		return wlmLmCrossDomainGranularityTransferR1KeyLess(candidates[i].key,candidates[j].key)
	})
	return candidates
}

func wlmLmCrossDomainGranularityTransferR1LearnCorpus(
	store *wlmLmCrossDomainGranularityTransferR1Store,
	corpus wlmLmCrossDomainGranularityTransferR1Corpus,
	metrics map[string]float64,
) int {
	added := 0
	for _, candidate := range wlmLmCrossDomainGranularityTransferR1Mine(corpus,store,metrics) {
		if len(store.order) >= 12 { break }
		if store.add(candidate.key,candidate.successor,metrics) {
			added++
		}
	}
	return added
}

func wlmLmCrossDomainGranularityTransferR1Accuracy(
	store *wlmLmCrossDomainGranularityTransferR1Store,
	corpus wlmLmCrossDomainGranularityTransferR1Corpus,
	metrics map[string]float64,
) float64 {
	correct := 0
	total := 0
	for _, pos := range corpus.targets {
		if pos < 4 || pos >= len(corpus.bytes) {
			metrics["invalid_stream_rows"]++
			continue
		}
		key := [4]uint8{corpus.bytes[pos-4],corpus.bytes[pos-3],corpus.bytes[pos-2],corpus.bytes[pos-1]}
		if store.predict(key) == corpus.bytes[pos] { correct++ }
		total++
	}
	if total == 0 { return 0 }
	return float64(correct)/float64(total)
}

func wlmLmCrossDomainGranularityTransferR1EventReduction(
	store *wlmLmCrossDomainGranularityTransferR1Store,
	corpus wlmLmCrossDomainGranularityTransferR1Corpus,
) float64 {
	events := 0
	for i:=0;i<len(corpus.bytes); {
		if i+4 <= len(corpus.bytes) {
			key := [4]uint8{corpus.bytes[i],corpus.bytes[i+1],corpus.bytes[i+2],corpus.bytes[i+3]}
			if _,ok:=store.entries[key]; ok {
				events++
				i+=4
				continue
			}
		}
		events++
		i++
	}
	return 1-float64(events)/float64(len(corpus.bytes))
}

func wlmLmCrossDomainGranularityTransferR1SharedReuse(store *wlmLmCrossDomainGranularityTransferR1Store) int {
	count:=0
	for id:=0;id<4;id++ {
		if _,ok:=store.entries[wlmLmCrossDomainGranularityTransferR1SharedMotif(id)];ok { count++ }
	}
	return count
}

func wlmLmCrossDomainGranularityTransferR1Min(a,b float64) float64 { if b<a{return b};return a }
func wlmLmCrossDomainGranularityTransferR1Max(a,b float64) float64 { if b>a{return b};return a }

// RunWlmLmCrossDomainGranularityTransferR1 evaluates the frozen WBG-4 domain-transfer matrix.
func RunWlmLmCrossDomainGranularityTransferR1() interface{} {
	seeds := [...]uint32{12007,12011,12037,12041}
	metrics := map[string]float64{
		"valid_seed_count":0,
		"minimum_zero_shot_shared_transfer_accuracy":1,
		"minimum_adapted_target_accuracy":1,
		"minimum_fresh_target_accuracy":1,
		"maximum_transfer_added_motif_count":0,
		"minimum_fresh_added_motif_count":12,
		"maximum_transfer_to_fresh_added_motif_ratio":0,
		"maximum_source_accuracy_loss_after_target_adaptation":0,
		"minimum_post_adaptation_source_accuracy":1,
		"minimum_heldout_event_reduction_fraction":1,
		"minimum_shared_motif_reuse_count":4,
		"maximum_global_store_size":0,
		"minimum_global_store_size_after_source_training":12,
		"maximum_global_store_size_after_source_training":0,
		"motif_eviction_count":0,
		"motif_rewrite_count":0,
		"invalid_stream_rows":0,
		"counter_overflow_rows":0,
	}
	for _,seed:=range seeds {
		store:=wlmLmCrossDomainGranularityTransferR1NewStore()
		for domain:=0;domain<3;domain++ {
			train:=wlmLmCrossDomainGranularityTransferR1CorpusFor(seed,domain,3072,0)
			wlmLmCrossDomainGranularityTransferR1LearnCorpus(store,train,metrics)
		}
		sourceSize:=len(store.order)
		metrics["minimum_global_store_size_after_source_training"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_global_store_size_after_source_training"],float64(sourceSize))
		metrics["maximum_global_store_size_after_source_training"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_global_store_size_after_source_training"],float64(sourceSize))
		metrics["maximum_global_store_size"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_global_store_size"],float64(sourceSize))
		metrics["minimum_shared_motif_reuse_count"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_shared_motif_reuse_count"],float64(wlmLmCrossDomainGranularityTransferR1SharedReuse(store)))

		preSource:=make([]float64,3)
		for domain:=0;domain<3;domain++ {
			eval:=wlmLmCrossDomainGranularityTransferR1CorpusFor(seed,domain,768,1)
			preSource[domain]=wlmLmCrossDomainGranularityTransferR1Accuracy(store,eval,metrics)
		}
		targetEval:=wlmLmCrossDomainGranularityTransferR1CorpusFor(seed,3,1536,1)
		zeroShot:=wlmLmCrossDomainGranularityTransferR1Accuracy(store,targetEval,metrics)
		metrics["minimum_zero_shot_shared_transfer_accuracy"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_zero_shot_shared_transfer_accuracy"],zeroShot)

		targetAdapt:=wlmLmCrossDomainGranularityTransferR1CorpusFor(seed,3,768,2)
		before:=len(store.order)
		wlmLmCrossDomainGranularityTransferR1LearnCorpus(store,targetAdapt,metrics)
		transferAdded:=len(store.order)-before
		metrics["maximum_transfer_added_motif_count"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_transfer_added_motif_count"],float64(transferAdded))
		metrics["maximum_global_store_size"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_global_store_size"],float64(len(store.order)))
		adapted:=wlmLmCrossDomainGranularityTransferR1Accuracy(store,targetEval,metrics)
		metrics["minimum_adapted_target_accuracy"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_adapted_target_accuracy"],adapted)
		metrics["minimum_heldout_event_reduction_fraction"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_heldout_event_reduction_fraction"],wlmLmCrossDomainGranularityTransferR1EventReduction(store,targetEval))

		fresh:=wlmLmCrossDomainGranularityTransferR1NewStore()
		freshAdded:=wlmLmCrossDomainGranularityTransferR1LearnCorpus(fresh,targetAdapt,metrics)
		metrics["minimum_fresh_added_motif_count"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_fresh_added_motif_count"],float64(freshAdded))
		freshAcc:=wlmLmCrossDomainGranularityTransferR1Accuracy(fresh,targetEval,metrics)
		metrics["minimum_fresh_target_accuracy"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_fresh_target_accuracy"],freshAcc)
		if freshAdded>0 {
			ratio:=float64(transferAdded)/float64(freshAdded)
			metrics["maximum_transfer_to_fresh_added_motif_ratio"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_transfer_to_fresh_added_motif_ratio"],ratio)
		}

		for domain:=0;domain<3;domain++ {
			eval:=wlmLmCrossDomainGranularityTransferR1CorpusFor(seed,domain,768,1)
			post:=wlmLmCrossDomainGranularityTransferR1Accuracy(store,eval,metrics)
			metrics["minimum_post_adaptation_source_accuracy"]=wlmLmCrossDomainGranularityTransferR1Min(metrics["minimum_post_adaptation_source_accuracy"],post)
			loss:=preSource[domain]-post
			if loss<0 { loss=0 }
			metrics["maximum_source_accuracy_loss_after_target_adaptation"]=wlmLmCrossDomainGranularityTransferR1Max(metrics["maximum_source_accuracy_loss_after_target_adaptation"],loss)
		}
		metrics["valid_seed_count"]++
	}
	return wlmLmCrossDomainGranularityTransferR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-CROSS-DOMAIN-GRANULARITY-TRANSFER-R1",
		Metrics:metrics,
	}
}
