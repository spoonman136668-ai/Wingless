package nr1

import (
	"fmt"
	"sort"
)

const LocalityReportSchema = "wingless.nr1.locality-report.v1"

type CoveragePoint struct {
	Target  float64 `json:"target"`
	Experts int     `json:"experts"`
}

type FractionTraffic struct {
	ExpertFraction float64 `json:"expert_fraction"`
	Traffic        float64 `json:"traffic"`
}

type LayerStats struct {
	Layer             int               `json:"layer"`
	Events            int               `json:"events"`
	Selections        int               `json:"selections"`
	UniqueExperts     int               `json:"unique_experts"`
	AdjacentPairs     int               `json:"adjacent_pairs"`
	AdjacentReuseRate float64           `json:"adjacent_reuse_rate"`
	Coverage          []CoveragePoint   `json:"coverage"`
	TopTraffic        []FractionTraffic `json:"top_traffic"`
}

type ScopeHotset struct {
	Scope        string          `json:"scope"`
	Events       int             `json:"events"`
	UniqueExperts int            `json:"unique_experts"`
	Coverage     []CoveragePoint `json:"coverage"`
}

type Transition struct {
	Layer int `json:"layer"`
	From  int `json:"from"`
	To    int `json:"to"`
	Count int `json:"count"`
}

type LocalityReport struct {
	Schema             string            `json:"schema"`
	Events             int               `json:"events"`
	Tokens             int               `json:"tokens"`
	Sessions           int               `json:"sessions"`
	Workloads          int               `json:"workloads"`
	TaskFamilies       int               `json:"task_families"`
	UniqueExpertKeys   int               `json:"unique_expert_keys"`
	AdjacentPairs      int               `json:"adjacent_pairs"`
	AdjacentReuseRate  float64           `json:"adjacent_reuse_rate"`
	GlobalCoverage     []CoveragePoint   `json:"global_coverage"`
	GlobalTopTraffic   []FractionTraffic `json:"global_top_traffic"`
	Layers             []LayerStats      `json:"layers"`
	SessionHotsets     []ScopeHotset     `json:"session_hotsets"`
	TaskHotsets        []ScopeHotset     `json:"task_hotsets"`
	TopTransitions     []Transition      `json:"top_transitions"`
}

type routedKey struct {
	Layer  int
	Expert int
}

type orderedEvent struct {
	Session string
	Layer   int
	Token   int64
	Experts []int
}

func Analyze(events []TraceEvent) (LocalityReport, error) {
	var report LocalityReport
	if len(events) == 0 {
		return report, fmt.Errorf("NR1_ANALYZE_EMPTY")
	}
	report.Schema = LocalityReportSchema
	report.Events = len(events)

	tokens := map[string]struct{}{}
	sessions := map[string]struct{}{}
	workloads := map[string]struct{}{}
	tasks := map[string]struct{}{}
	global := map[routedKey]int{}
	byLayer := map[int]map[int]int{}
	layerEvents := map[int]int{}
	bySession := map[string]map[routedKey]int{}
	sessionEvents := map[string]int{}
	byTask := map[string]map[routedKey]int{}
	taskEvents := map[string]int{}
	ordered := make([]orderedEvent, 0, len(events))

	for i, e := range events {
		if err := e.Validate(); err != nil {
			return report, fmt.Errorf("event %d: %w", i, err)
		}
		tokens[e.SessionID+"\x00"+e.WorkloadID+"\x00"+fmt.Sprint(e.TokenIndex)] = struct{}{}
		sessions[e.SessionID] = struct{}{}
		workloads[e.WorkloadID] = struct{}{}
		tasks[e.TaskFamily] = struct{}{}
		if byLayer[e.Layer] == nil {
			byLayer[e.Layer] = map[int]int{}
		}
		if bySession[e.SessionID] == nil {
			bySession[e.SessionID] = map[routedKey]int{}
		}
		if byTask[e.TaskFamily] == nil {
			byTask[e.TaskFamily] = map[routedKey]int{}
		}
		layerEvents[e.Layer]++
		sessionEvents[e.SessionID]++
		taskEvents[e.TaskFamily]++
		for _, expert := range e.Experts {
			key := routedKey{Layer: e.Layer, Expert: expert}
			global[key]++
			byLayer[e.Layer][expert]++
			bySession[e.SessionID][key]++
			byTask[e.TaskFamily][key]++
		}
		ordered = append(ordered, orderedEvent{Session: e.SessionID, Layer: e.Layer, Token: e.TokenIndex, Experts: append([]int(nil), e.Experts...)})
	}

	report.Tokens = len(tokens)
	report.Sessions = len(sessions)
	report.Workloads = len(workloads)
	report.TaskFamilies = len(tasks)
	report.UniqueExpertKeys = len(global)
	report.GlobalCoverage = coverageRouted(global)
	report.GlobalTopTraffic = topTrafficRouted(global)

	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Session != ordered[j].Session {
			return ordered[i].Session < ordered[j].Session
		}
		if ordered[i].Layer != ordered[j].Layer {
			return ordered[i].Layer < ordered[j].Layer
		}
		return ordered[i].Token < ordered[j].Token
	})

	layerOverlap := map[int]float64{}
	layerPairs := map[int]int{}
	transitions := map[Transition]int{}
	var overlapSum float64
	for i := 1; i < len(ordered); i++ {
		prev, cur := ordered[i-1], ordered[i]
		if prev.Session != cur.Session || prev.Layer != cur.Layer || cur.Token != prev.Token+1 {
			continue
		}
		overlap := selectedOverlap(prev.Experts, cur.Experts)
		den := len(prev.Experts)
		if len(cur.Experts) < den {
			den = len(cur.Experts)
		}
		if den > 0 {
			rate := float64(overlap) / float64(den)
			overlapSum += rate
			layerOverlap[cur.Layer] += rate
			report.AdjacentPairs++
			layerPairs[cur.Layer]++
		}
		for _, from := range prev.Experts {
			for _, to := range cur.Experts {
				key := Transition{Layer: cur.Layer, From: from, To: to}
				transitions[key]++
			}
		}
	}
	if report.AdjacentPairs > 0 {
		report.AdjacentReuseRate = overlapSum / float64(report.AdjacentPairs)
	}

	layerIDs := make([]int, 0, len(byLayer))
	for layer := range byLayer {
		layerIDs = append(layerIDs, layer)
	}
	sort.Ints(layerIDs)
	for _, layer := range layerIDs {
		freq := byLayer[layer]
		stats := LayerStats{
			Layer:         layer,
			Events:        layerEvents[layer],
			UniqueExperts: len(freq),
			Coverage:      coverageInt(freq),
			TopTraffic:    topTrafficInt(freq),
			AdjacentPairs: layerPairs[layer],
		}
		for _, count := range freq {
			stats.Selections += count
		}
		if stats.AdjacentPairs > 0 {
			stats.AdjacentReuseRate = layerOverlap[layer] / float64(stats.AdjacentPairs)
		}
		report.Layers = append(report.Layers, stats)
	}

	report.SessionHotsets = scopeHotsets(bySession, sessionEvents)
	report.TaskHotsets = scopeHotsets(byTask, taskEvents)

	for key, count := range transitions {
		key.Count = count
		report.TopTransitions = append(report.TopTransitions, key)
	}
	sort.Slice(report.TopTransitions, func(i, j int) bool {
		a, b := report.TopTransitions[i], report.TopTransitions[j]
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		if a.Layer != b.Layer {
			return a.Layer < b.Layer
		}
		if a.From != b.From {
			return a.From < b.From
		}
		return a.To < b.To
	})
	if len(report.TopTransitions) > 128 {
		report.TopTransitions = report.TopTransitions[:128]
	}
	return report, nil
}

func selectedOverlap(a, b []int) int {
	set := make(map[int]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}
	n := 0
	for _, v := range b {
		if _, ok := set[v]; ok {
			n++
		}
	}
	return n
}

func scopeHotsets(freq map[string]map[routedKey]int, events map[string]int) []ScopeHotset {
	ids := make([]string, 0, len(freq))
	for id := range freq {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]ScopeHotset, 0, len(ids))
	for _, id := range ids {
		out = append(out, ScopeHotset{Scope: id, Events: events[id], UniqueExperts: len(freq[id]), Coverage: coverageRouted(freq[id])})
	}
	return out
}

var coverageTargets = []float64{0.80, 0.90, 0.95, 0.99}
var topFractions = []float64{0.05, 0.10, 0.20, 0.50}

func coverageInt(freq map[int]int) []CoveragePoint {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return coverageCounts(counts)
}

func coverageRouted(freq map[routedKey]int) []CoveragePoint {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return coverageCounts(counts)
}

func coverageCounts(counts []int) []CoveragePoint {
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	total := 0
	for _, n := range counts {
		total += n
	}
	out := make([]CoveragePoint, 0, len(coverageTargets))
	for _, target := range coverageTargets {
		running, needed := 0, 0
		for i, n := range counts {
			running += n
			needed = i + 1
			if total > 0 && float64(running)/float64(total) >= target {
				break
			}
		}
		out = append(out, CoveragePoint{Target: target, Experts: needed})
	}
	return out
}

func topTrafficInt(freq map[int]int) []FractionTraffic {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return topTrafficCounts(counts)
}

func topTrafficRouted(freq map[routedKey]int) []FractionTraffic {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return topTrafficCounts(counts)
}

func topTrafficCounts(counts []int) []FractionTraffic {
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	total := 0
	for _, n := range counts {
		total += n
	}
	out := make([]FractionTraffic, 0, len(topFractions))
	for _, fraction := range topFractions {
		k := int(float64(len(counts))*fraction + 0.999999999)
		if k < 1 && len(counts) > 0 {
			k = 1
		}
		if k > len(counts) {
			k = len(counts)
		}
		hits := 0
		for i := 0; i < k; i++ {
			hits += counts[i]
		}
		traffic := 0.0
		if total > 0 {
			traffic = float64(hits) / float64(total)
		}
		out = append(out, FractionTraffic{ExpertFraction: fraction, Traffic: traffic})
	}
	return out
}
