package nr1

import (
	"container/heap"
	"container/list"
	"errors"
	"sort"
)

const ResidencyReportSchema = "wingless.nr1.residency-report.v1"

type ResidencyConfig struct {
	ExpertBytes   uint64   `json:"expert_bytes"`
	ReservedBytes uint64   `json:"reserved_bytes"`
	BudgetsBytes  []uint64 `json:"budgets_bytes"`
}

type ResidencyResult struct {
	Policy            string  `json:"policy"`
	BudgetBytes       uint64  `json:"budget_bytes"`
	ReservedBytes     uint64  `json:"reserved_bytes"`
	ExpertCacheBytes  uint64  `json:"expert_cache_bytes"`
	Slots             int     `json:"slots"`
	Accesses          int64   `json:"accesses"`
	Hits              int64   `json:"hits"`
	Misses            int64   `json:"misses"`
	HitRate           float64 `json:"hit_rate"`
	MissBytes         uint64  `json:"miss_bytes"`
	WarmBytesPerToken float64 `json:"warm_bytes_per_token"`
	Tokens            int     `json:"tokens"`
}

type ResidencyReport struct {
	Schema      string            `json:"schema"`
	ExpertBytes uint64            `json:"expert_bytes"`
	Results     []ResidencyResult `json:"results"`
}

func Simulate(events []TraceEvent, cfg ResidencyConfig) (ResidencyReport, error) {
	report := ResidencyReport{Schema: ResidencyReportSchema, ExpertBytes: cfg.ExpertBytes}
	if len(events) == 0 || cfg.ExpertBytes == 0 || len(cfg.BudgetsBytes) == 0 {
		return report, errors.New("NR1_RESIDENCY_CONFIG_INVALID")
	}
	for _, e := range events {
		if err := e.Validate(); err != nil {
			return report, err
		}
	}
	accesses := flatten(events)
	tokens := distinctTokenCount(events)
	globalRank := rankedKeys(accesses)
	sessionRanks := rankedSessionKeys(events)

	budgets := append([]uint64(nil), cfg.BudgetsBytes...)
	sort.Slice(budgets, func(i, j int) bool { return budgets[i] < budgets[j] })
	for _, budget := range budgets {
		cacheBytes := uint64(0)
		if budget > cfg.ReservedBytes {
			cacheBytes = budget - cfg.ReservedBytes
		}
		slots64 := cacheBytes / cfg.ExpertBytes
		if slots64 > uint64(^uint(0)>>1) {
			return report, errors.New("NR1_RESIDENCY_SLOT_OVERFLOW")
		}
		slots := int(slots64)
		for _, policy := range []string{"lru", "lfu", "slru", "static-global", "session-static"} {
			var hits int64
			switch policy {
			case "lru":
				hits = simulateLRU(accesses, slots)
			case "lfu":
				hits = simulateLFU(accesses, slots)
			case "slru":
				hits = simulateSLRU(accesses, slots)
			case "static-global":
				hits = simulateStatic(accesses, slots, globalRank)
			case "session-static":
				hits = simulateSessionStatic(events, slots, sessionRanks)
			}
			result := buildResidencyResult(policy, budget, cfg.ReservedBytes, cacheBytes, slots, int64(len(accesses)), hits, cfg.ExpertBytes, tokens)
			report.Results = append(report.Results, result)
		}
	}
	return report, nil
}

type access struct {
	Session string
	Key     routedKey
}

func flatten(events []TraceEvent) []access {
	n := 0
	for _, e := range events {
		n += len(e.Experts)
	}
	out := make([]access, 0, n)
	for _, e := range events {
		for _, expert := range e.Experts {
			out = append(out, access{Session: e.SessionID, Key: routedKey{Layer: e.Layer, Expert: expert}})
		}
	}
	return out
}

func distinctTokenCount(events []TraceEvent) int {
	seen := map[string]struct{}{}
	for _, e := range events {
		seen[e.SessionID+"\x00"+e.WorkloadID+"\x00"+itoa64(e.TokenIndex)] = struct{}{}
	}
	return len(seen)
}

func itoa64(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b [32]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func buildResidencyResult(policy string, budget, reserved, cacheBytes uint64, slots int, accesses, hits int64, expertBytes uint64, tokens int) ResidencyResult {
	misses := accesses - hits
	result := ResidencyResult{Policy: policy, BudgetBytes: budget, ReservedBytes: reserved, ExpertCacheBytes: cacheBytes, Slots: slots, Accesses: accesses, Hits: hits, Misses: misses, Tokens: tokens}
	if accesses > 0 {
		result.HitRate = float64(hits) / float64(accesses)
	}
	if misses > 0 && expertBytes > 0 {
		if uint64(misses) > ^uint64(0)/expertBytes {
			result.MissBytes = ^uint64(0)
		} else {
			result.MissBytes = uint64(misses) * expertBytes
		}
	}
	if tokens > 0 {
		result.WarmBytesPerToken = float64(result.MissBytes) / float64(tokens)
	}
	return result
}

func simulateLRU(accesses []access, slots int) int64 {
	if slots <= 0 {
		return 0
	}
	order := list.New()
	resident := map[routedKey]*list.Element{}
	var hits int64
	for _, a := range accesses {
		if el := resident[a.Key]; el != nil {
			hits++
			order.MoveToFront(el)
			continue
		}
		if order.Len() >= slots {
			el := order.Back()
			delete(resident, el.Value.(routedKey))
			order.Remove(el)
		}
		resident[a.Key] = order.PushFront(a.Key)
	}
	return hits
}

type lfuState struct {
	freq    int64
	last    int64
	version int64
}

type lfuItem struct {
	key     routedKey
	freq    int64
	last    int64
	version int64
}

type lfuHeap []lfuItem

func (h lfuHeap) Len() int { return len(h) }
func (h lfuHeap) Less(i, j int) bool {
	if h[i].freq != h[j].freq {
		return h[i].freq < h[j].freq
	}
	return h[i].last < h[j].last
}
func (h lfuHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *lfuHeap) Push(x any)   { *h = append(*h, x.(lfuItem)) }
func (h *lfuHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func simulateLFU(accesses []access, slots int) int64 {
	if slots <= 0 {
		return 0
	}
	states := map[routedKey]lfuState{}
	h := &lfuHeap{}
	heap.Init(h)
	var hits int64
	var clock int64
	for _, a := range accesses {
		clock++
		if st, ok := states[a.Key]; ok {
			hits++
			st.freq++
			st.last = clock
			st.version++
			states[a.Key] = st
			heap.Push(h, lfuItem{key: a.Key, freq: st.freq, last: st.last, version: st.version})
			continue
		}
		for len(states) >= slots {
			item := heap.Pop(h).(lfuItem)
			st, ok := states[item.key]
			if !ok || st.version != item.version || st.freq != item.freq || st.last != item.last {
				continue
			}
			delete(states, item.key)
			break
		}
		st := lfuState{freq: 1, last: clock, version: 1}
		states[a.Key] = st
		heap.Push(h, lfuItem{key: a.Key, freq: st.freq, last: st.last, version: st.version})
	}
	return hits
}

// simulateSLRU is a bounded recency/frequency hybrid. New entries enter a
// probationary LRU; a reuse promotes them into a protected LRU. The protected
// region is capped at half the cache and demotes its LRU entry on overflow.
func simulateSLRU(accesses []access, slots int) int64 {
	if slots <= 0 {
		return 0
	}
	protectedCap := slots / 2
	if protectedCap < 1 {
		protectedCap = 1
	}
	probationCap := slots - protectedCap
	if probationCap < 0 {
		probationCap = 0
	}
	protected := list.New()
	probation := list.New()
	type loc struct {
		protected bool
		el        *list.Element
	}
	resident := map[routedKey]loc{}
	var hits int64

	trimProbation := func() {
		for probation.Len() > probationCap {
			el := probation.Back()
			key := el.Value.(routedKey)
			probation.Remove(el)
			delete(resident, key)
		}
	}
	for _, a := range accesses {
		if where, ok := resident[a.Key]; ok {
			hits++
			if where.protected {
				protected.MoveToFront(where.el)
				continue
			}
			probation.Remove(where.el)
			el := protected.PushFront(a.Key)
			resident[a.Key] = loc{protected: true, el: el}
			if protected.Len() > protectedCap {
				demote := protected.Back()
				key := demote.Value.(routedKey)
				protected.Remove(demote)
				pel := probation.PushFront(key)
				resident[key] = loc{el: pel}
				trimProbation()
			}
			continue
		}
		if probationCap == 0 {
			if protected.Len() >= slots {
				el := protected.Back()
				delete(resident, el.Value.(routedKey))
				protected.Remove(el)
			}
			el := protected.PushFront(a.Key)
			resident[a.Key] = loc{protected: true, el: el}
			continue
		}
		el := probation.PushFront(a.Key)
		resident[a.Key] = loc{el: el}
		trimProbation()
	}
	return hits
}

func rankedKeys(accesses []access) []routedKey {
	freq := map[routedKey]int64{}
	for _, a := range accesses {
		freq[a.Key]++
	}
	keys := make([]routedKey, 0, len(freq))
	for key := range freq {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if freq[keys[i]] != freq[keys[j]] {
			return freq[keys[i]] > freq[keys[j]]
		}
		if keys[i].Layer != keys[j].Layer {
			return keys[i].Layer < keys[j].Layer
		}
		return keys[i].Expert < keys[j].Expert
	})
	return keys
}

func simulateStatic(accesses []access, slots int, rank []routedKey) int64 {
	if slots <= 0 {
		return 0
	}
	if slots > len(rank) {
		slots = len(rank)
	}
	resident := make(map[routedKey]struct{}, slots)
	for _, key := range rank[:slots] {
		resident[key] = struct{}{}
	}
	var hits int64
	for _, a := range accesses {
		if _, ok := resident[a.Key]; ok {
			hits++
		}
	}
	return hits
}

func rankedSessionKeys(events []TraceEvent) map[string][]routedKey {
	freq := map[string]map[routedKey]int64{}
	for _, e := range events {
		if freq[e.SessionID] == nil {
			freq[e.SessionID] = map[routedKey]int64{}
		}
		for _, expert := range e.Experts {
			freq[e.SessionID][routedKey{Layer: e.Layer, Expert: expert}]++
		}
	}
	out := map[string][]routedKey{}
	for session, counts := range freq {
		keys := make([]routedKey, 0, len(counts))
		for key := range counts {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			if counts[keys[i]] != counts[keys[j]] {
				return counts[keys[i]] > counts[keys[j]]
			}
			if keys[i].Layer != keys[j].Layer {
				return keys[i].Layer < keys[j].Layer
			}
			return keys[i].Expert < keys[j].Expert
		})
		out[session] = keys
	}
	return out
}

func simulateSessionStatic(events []TraceEvent, slots int, ranks map[string][]routedKey) int64 {
	if slots <= 0 {
		return 0
	}
	resident := map[string]map[routedKey]struct{}{}
	for session, rank := range ranks {
		n := slots
		if n > len(rank) {
			n = len(rank)
		}
		set := make(map[routedKey]struct{}, n)
		for _, key := range rank[:n] {
			set[key] = struct{}{}
		}
		resident[session] = set
	}
	var hits int64
	for _, e := range events {
		set := resident[e.SessionID]
		for _, expert := range e.Experts {
			if _, ok := set[routedKey{Layer: e.Layer, Expert: expert}]; ok {
				hits++
			}
		}
	}
	return hits
}
