package unitary

import (
	"crypto/sha1"
	"fmt"
	"math"
	"os/exec"
	"sort"
)

type wlmLmRawRepPredFreshHoldoutR1File struct {
	path string
	sha  string
	size int
}

type wlmLmRawRepPredFreshHoldoutR1Stats struct {
	total  uint32
	counts [256]uint32
}

type wlmLmRawRepPredFreshHoldoutR1Motif struct {
	key         [4]uint8
	stats       *wlmLmRawRepPredFreshHoldoutR1Stats
	best        uint8
	bestCount   uint32
	consistency float64
}

type wlmLmRawRepPredFreshHoldoutR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

var wlmLmRawRepPredFreshHoldoutR1Train = []wlmLmRawRepPredFreshHoldoutR1File{
	{"README.md", "50f1ab5ae468551ee48dd1143e534a0ebe01d6a8", 0},
	{"docs/architecture.md", "dec5b6284ec69c1c94a4cf5b005022981df0e900", 0},
	{"docs/internal-stage-2.md", "d71357155dc1b141369dfd61d971b0eb1d00c9c9", 0},
	{"docs/inference-backends.md", "4e5b3c2c268e1ded129020fd22e8fbe07725ec30", 0},
	{"unitary/wlm_lm_byte_motif_discovery_r3_balanced.go", "f9ce41417b5000e3c2419b1b5bf0ce2fc4c41899", 0},
	{"unitary/wlm_lm_cross_domain_granularity_transfer_r1.go", "07b3a6fe2dec98b09c91f770f8d5648defc49ae7", 0},
	{"unitary/adaptive_continuous_fusion.go", "60ae9ce88419cf2541205ad41d81b8c0d3cbe114", 0},
	{"unitary/anonymous_gram.go", "8a85e5e09a0ffa4a0e8e849df1a63e78ecfd2a1e", 0},
	{"unitary/memory.go", "24b132e0bf81e373d7b82b70719b9ee9419d0054", 0},
	{"unitary/training.go", "f39dcd4cda2f2b78b97df77fc7098199afa55acb", 0},
}

var wlmLmRawRepPredFreshHoldoutR1Eval = []wlmLmRawRepPredFreshHoldoutR1File{
	{"unitary/full_latent.go", "5d9400dcc08c5d10cc34303876d88ab1fb785cc8", 30831},
	{"unitary/real_orthogonal_equivalence.go", "46e946ace9d91347f86db6ad63233df14a5c7c1e", 27322},
	{"unitary/learned_anchor.go", "4a47552fea9927a01ca895674bfbbb89bc9b03cd", 27201},
	{"unitary/discovery_breadth.go", "fd97feba07d4007bab43aeb6663232c42edb6cf1", 26934},
}

func wlmLmRawRepPredFreshHoldoutR1BlobHash(data []byte) string {
	h := sha1.New()
	_, _ = h.Write([]byte(fmt.Sprintf("blob %d\x00", len(data))))
	_, _ = h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func wlmLmRawRepPredFreshHoldoutR1Load(f wlmLmRawRepPredFreshHoldoutR1File) ([]byte, bool) {
	data, err := exec.Command("git", "cat-file", "blob", f.sha).Output()
	if err != nil {
		return nil, false
	}
	if wlmLmRawRepPredFreshHoldoutR1BlobHash(data) != f.sha {
		return data, false
	}
	if f.size > 0 && len(data) != f.size {
		return data, false
	}
	return data, true
}

func wlmLmRawRepPredFreshHoldoutR1Best(counts *[256]uint32) (uint8, uint32, uint64) {
	best := uint8(0)
	bestCount := counts[0]
	total := uint64(counts[0])
	for i := 1; i < 256; i++ {
		total += uint64(counts[i])
		if counts[i] > bestCount {
			best = uint8(i)
			bestCount = counts[i]
		}
	}
	return best, bestCount, total
}

func wlmLmRawRepPredFreshHoldoutR1Less(a, b [4]uint8) bool {
	for i := 0; i < 4; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func wlmLmRawRepPredFreshHoldoutR1TrainModel(files [][]byte, metrics map[string]float64) ([256][256]uint32, map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) {
	var baseline [256][256]uint32
	candidates := make(map[[4]uint8]*wlmLmRawRepPredFreshHoldoutR1Stats)
	for _, data := range files {
		for i := 1; i < len(data); i++ {
			if baseline[data[i-1]][data[i]] == ^uint32(0) {
				metrics["counter_overflow_rows"]++
				continue
			}
			baseline[data[i-1]][data[i]]++
		}
		for i := 4; i < len(data); i++ {
			key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
			s := candidates[key]
			if s == nil {
				s = &wlmLmRawRepPredFreshHoldoutR1Stats{}
				candidates[key] = s
			}
			if s.total == ^uint32(0) || s.counts[data[i]] == ^uint32(0) {
				metrics["counter_overflow_rows"]++
				continue
			}
			s.total++
			s.counts[data[i]]++
		}
	}
	rows := make([]wlmLmRawRepPredFreshHoldoutR1Motif, 0)
	for key, s := range candidates {
		if s.total < 4 {
			continue
		}
		best, bestCount, _ := wlmLmRawRepPredFreshHoldoutR1Best(&s.counts)
		rows = append(rows, wlmLmRawRepPredFreshHoldoutR1Motif{
			key: key, stats: s, best: best, bestCount: bestCount,
			consistency: float64(bestCount) / float64(s.total),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].bestCount != rows[j].bestCount {
			return rows[i].bestCount > rows[j].bestCount
		}
		if rows[i].consistency != rows[j].consistency {
			return rows[i].consistency > rows[j].consistency
		}
		return wlmLmRawRepPredFreshHoldoutR1Less(rows[i].key, rows[j].key)
	})
	if len(rows) > 512 {
		rows = rows[:512]
	}
	selected := make(map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif, len(rows))
	for _, row := range rows {
		selected[row.key] = row
	}
	return baseline, selected
}

func wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(baseline *[256][256]uint32, prev uint8) uint8 {
	best := uint8(0)
	bestCount := baseline[prev][0]
	for i := 1; i < 256; i++ {
		if baseline[prev][i] > bestCount {
			best = uint8(i)
			bestCount = baseline[prev][i]
		}
	}
	return best
}

func wlmLmRawRepPredFreshHoldoutR1BaselineProb(baseline *[256][256]uint32, prev, target uint8) float64 {
	var total uint64
	for i := 0; i < 256; i++ {
		total += uint64(baseline[prev][i])
	}
	return (float64(baseline[prev][target]) + 0.5) / (float64(total) + 128.0)
}

func wlmLmRawRepPredFreshHoldoutR1MotifProb(counts *[256]uint32, target uint8) float64 {
	_, _, total := wlmLmRawRepPredFreshHoldoutR1Best(counts)
	return (float64(counts[target]) + 0.5) / (float64(total) + 128.0)
}

func wlmLmRawRepPredFreshHoldoutR1EventReduction(data []byte, selected map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) float64 {
	if len(data) == 0 {
		return 0
	}
	events := 0
	for i := 0; i < len(data); {
		if i+4 <= len(data) {
			key := [4]uint8{data[i], data[i+1], data[i+2], data[i+3]}
			if _, ok := selected[key]; ok {
				events++
				i += 4
				continue
			}
		}
		events++
		i++
	}
	return 1.0 - float64(events)/float64(len(data))
}

func wlmLmRawRepPredFreshHoldoutR1Evaluate(data []byte, baseline *[256][256]uint32, selected map[[4]uint8]wlmLmRawRepPredFreshHoldoutR1Motif) (float64, float64, float64, float64) {
	if len(data) < 5 {
		return 0, 0, 0, 0
	}
	covered := 0
	baseCoveredCorrect := 0
	repCoveredCorrect := 0
	var baseBits, repBits float64
	total := 0
	for i := 1; i < len(data); i++ {
		target := data[i]
		basePred := wlmLmRawRepPredFreshHoldoutR1BaselinePrediction(baseline, data[i-1])
		baseP := wlmLmRawRepPredFreshHoldoutR1BaselineProb(baseline, data[i-1], target)
		repP := baseP
		if i >= 4 {
			key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
			if motif, ok := selected[key]; ok {
				covered++
				if basePred == target {
					baseCoveredCorrect++
				}
				if motif.best == target {
					repCoveredCorrect++
				}
				repP = wlmLmRawRepPredFreshHoldoutR1MotifProb(&motif.stats.counts, target)
			}
		}
		baseBits += -math.Log2(baseP)
		repBits += -math.Log2(repP)
		total++
	}
	eligible := len(data) - 4
	if eligible < 1 {
		eligible = 1
	}
	coverage := float64(covered) / float64(eligible)
	gain := 0.0
	if covered > 0 {
		gain = float64(repCoveredCorrect-baseCoveredCorrect) / float64(covered)
	}
	reduction := wlmLmRawRepPredFreshHoldoutR1EventReduction(data, selected)
	bpbDelta := repBits/float64(total) - baseBits/float64(total)
	return coverage, gain, reduction, bpbDelta
}

func wlmLmRawRepPredFreshHoldoutR1Min(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

func wlmLmRawRepPredFreshHoldoutR1Max(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

// RunWlmLmRawRepPredFreshHoldoutR1 tests the frozen raw-input -> representation -> prediction bridge.
func RunWlmLmRawRepPredFreshHoldoutR1() interface{} {
	metrics := map[string]float64{
		"training_file_identity_mismatch_count":                 0,
		"evaluation_file_identity_mismatch_count":               0,
		"train_eval_blob_overlap_count":                         0,
		"valid_evaluation_file_count":                           0,
		"selected_motif_count":                                  0,
		"minimum_eval_representation_coverage_fraction":         1,
		"minimum_covered_top1_accuracy_gain":                    1,
		"minimum_effective_event_reduction_fraction":            1,
		"maximum_representation_minus_baseline_bits_per_byte":  -100,
		"files_with_positive_covered_top1_gain":                 0,
		"files_with_nonworse_bits_per_byte":                     0,
		"capacity_growth_event_count":                           0,
		"tokenizer_use_count":                                   0,
		"external_model_call_count":                             0,
		"invalid_byte_rows":                                     0,
		"counter_overflow_rows":                                 0,
	}
	trainSet := make(map[string]bool)
	for _, f := range wlmLmRawRepPredFreshHoldoutR1Train {
		trainSet[f.sha] = true
	}
	for _, f := range wlmLmRawRepPredFreshHoldoutR1Eval {
		if trainSet[f.sha] {
			metrics["train_eval_blob_overlap_count"]++
		}
	}
	train := make([][]byte, 0, len(wlmLmRawRepPredFreshHoldoutR1Train))
	for _, f := range wlmLmRawRepPredFreshHoldoutR1Train {
		data, ok := wlmLmRawRepPredFreshHoldoutR1Load(f)
		if !ok {
			metrics["training_file_identity_mismatch_count"]++
		}
		train = append(train, data)
	}
	baseline, selected := wlmLmRawRepPredFreshHoldoutR1TrainModel(train, metrics)
	metrics["selected_motif_count"] = float64(len(selected))
	for _, f := range wlmLmRawRepPredFreshHoldoutR1Eval {
		data, ok := wlmLmRawRepPredFreshHoldoutR1Load(f)
		if !ok {
			metrics["evaluation_file_identity_mismatch_count"]++
			continue
		}
		coverage, gain, reduction, bpbDelta := wlmLmRawRepPredFreshHoldoutR1Evaluate(data, &baseline, selected)
		metrics["minimum_eval_representation_coverage_fraction"] = wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_eval_representation_coverage_fraction"], coverage)
		metrics["minimum_covered_top1_accuracy_gain"] = wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_covered_top1_accuracy_gain"], gain)
		metrics["minimum_effective_event_reduction_fraction"] = wlmLmRawRepPredFreshHoldoutR1Min(metrics["minimum_effective_event_reduction_fraction"], reduction)
		metrics["maximum_representation_minus_baseline_bits_per_byte"] = wlmLmRawRepPredFreshHoldoutR1Max(metrics["maximum_representation_minus_baseline_bits_per_byte"], bpbDelta)
		if gain > 0 {
			metrics["files_with_positive_covered_top1_gain"]++
		}
		if bpbDelta <= 0 {
			metrics["files_with_nonworse_bits_per_byte"]++
		}
		metrics["valid_evaluation_file_count"]++
	}
	return wlmLmRawRepPredFreshHoldoutR1Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R1",
		Metrics: metrics,
	}
}
