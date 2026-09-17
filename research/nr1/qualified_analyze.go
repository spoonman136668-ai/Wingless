package nr1

import "sort"

// AnalyzeQualifiedQwen3Coder applies the known logical expert universe of the
// hash-pinned Qwen3-Coder-30B-A3B target. This matters for top-X%-of-experts
// traffic metrics: experts with zero observations still belong to the logical
// population and must remain in the denominator.
func AnalyzeQualifiedQwen3Coder(events []TraceEvent) (LocalityReport, error) {
	report, err := Analyze(events)
	if err != nil {
		return report, err
	}

	global := make(map[routedKey]int)
	byLayer := make(map[int]map[int]int)
	for _, event := range events {
		if byLayer[event.Layer] == nil {
			byLayer[event.Layer] = make(map[int]int)
		}
		for _, expert := range event.Experts {
			key := routedKey{Layer: event.Layer, Expert: expert}
			global[key]++
			byLayer[event.Layer][expert]++
		}
	}

	report.GlobalTopTraffic = topTrafficRoutedUniverse(
		global,
		QualifiedQwen3CoderLayers*QualifiedQwen3CoderExpertsPerLayer,
	)
	for i := range report.Layers {
		layer := report.Layers[i].Layer
		report.Layers[i].TopTraffic = topTrafficIntUniverse(
			byLayer[layer],
			QualifiedQwen3CoderExpertsPerLayer,
		)
	}
	return report, nil
}

func topTrafficIntUniverse(freq map[int]int, universe int) []FractionTraffic {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return topTrafficCountsUniverse(counts, universe)
}

func topTrafficRoutedUniverse(freq map[routedKey]int, universe int) []FractionTraffic {
	counts := make([]int, 0, len(freq))
	for _, n := range freq {
		counts = append(counts, n)
	}
	return topTrafficCountsUniverse(counts, universe)
}

func topTrafficCountsUniverse(counts []int, universe int) []FractionTraffic {
	if universe < len(counts) {
		universe = len(counts)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	total := 0
	for _, n := range counts {
		total += n
	}
	out := make([]FractionTraffic, 0, len(topFractions))
	for _, fraction := range topFractions {
		k := int(float64(universe)*fraction + 0.999999999)
		if k < 1 && universe > 0 {
			k = 1
		}
		if k > universe {
			k = universe
		}
		hits := 0
		observed := k
		if observed > len(counts) {
			observed = len(counts)
		}
		for i := 0; i < observed; i++ {
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
