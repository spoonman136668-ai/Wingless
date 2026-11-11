package unitary

import (
	"fmt"
	"math"
	"sort"
)

const proximalFusionLambda =
	directOffsetLearningRate * directOffsetResourcePrice

type proximalSortedPoint struct {
	value float64
	index int
}

type proximalBlock struct {
	start  int
	end    int
	weight float64
	mean   float64
}

type ProximalFusionDiagnostics struct {
	Lambda                    float64 `json:"lambda"`
	InputPairPenalty          float64 `json:"input_pair_penalty"`
	OutputPairPenalty         float64 `json:"output_pair_penalty"`
	ProxObjectiveBefore       float64 `json:"prox_objective_before"`
	ProxObjectiveAfter        float64 `json:"prox_objective_after"`
	ExactCapacity             int     `json:"exact_capacity"`
	ExactFusionOccurred       bool    `json:"exact_fusion_occurred"`
	MaximumMonotonicViolation float64 `json:"maximum_monotonic_violation"`
	MaximumBlockMeanResidual  float64 `json:"maximum_block_mean_residual"`
}

func completeGraphPairPenalty(values []float64) float64 {
	var total float64
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			total += math.Abs(values[i] - values[j])
		}
	}
	return total
}

func completeGraphProxObjective(
	x, y []float64,
	lambda float64,
) (float64, error) {
	if len(x) != len(y) || len(x) == 0 {
		return 0, fmt.Errorf("prox objective dimension mismatch")
	}
	if !finite(lambda) || lambda < 0 {
		return 0, fmt.Errorf("prox lambda invalid")
	}
	var squared float64
	for i := range x {
		if !finite(x[i]) || !finite(y[i]) {
			return 0, fmt.Errorf("prox objective non-finite input")
		}
		delta := x[i] - y[i]
		squared += delta * delta
	}
	return 0.5*squared +
		lambda*completeGraphPairPenalty(x), nil
}

func exactOffsetGroups(offsets []float64) ([][]int, int, error) {
	if len(offsets) != compositeChannels {
		return nil, 0, fmt.Errorf(
			"exact offset group dimension=%d want=%d",
			len(offsets), compositeChannels,
		)
	}
	groupMap := make(map[uint64][]int)
	keys := make([]uint64, 0, len(offsets))
	for index, value := range offsets {
		if !finite(value) {
			return nil, 0, fmt.Errorf(
				"exact offset group non-finite index=%d",
				index,
			)
		}
		key := math.Float64bits(value)
		if _, ok := groupMap[key]; !ok {
			keys = append(keys, key)
		}
		groupMap[key] = append(groupMap[key], index)
	}
	sort.Slice(keys, func(i, j int) bool {
		first := groupMap[keys[i]][0]
		second := groupMap[keys[j]][0]
		return first < second
	})
	groups := make([][]int, 0, len(keys))
	capacity := 0
	for _, key := range keys {
		group := append([]int(nil), groupMap[key]...)
		sort.Ints(group)
		groups = append(groups, group)
		capacity += len(group) * len(group)
	}
	return groups, capacity, nil
}

// completeGraphFusionProx solves
//
//   min_x 0.5 ||x-y||^2 + lambda * sum_{i<j}|x_i-x_j|
//
// exactly for scalar coordinates. Sorting reduces the complete-graph
// absolute-difference penalty to a linear term under the monotone order.
// The remaining problem is isotonic least squares and is solved by PAVA.
func completeGraphFusionProx(
	y []float64,
	lambda float64,
) ([]float64, ProximalFusionDiagnostics, error) {
	if len(y) != compositeChannels {
		return nil, ProximalFusionDiagnostics{}, fmt.Errorf(
			"prox dimension=%d want=%d",
			len(y), compositeChannels,
		)
	}
	if !finite(lambda) || lambda < 0 {
		return nil, ProximalFusionDiagnostics{}, fmt.Errorf(
			"prox lambda invalid",
		)
	}

	sorted := make([]proximalSortedPoint, len(y))
	for index, value := range y {
		if !finite(value) {
			return nil, ProximalFusionDiagnostics{}, fmt.Errorf(
				"prox input non-finite index=%d",
				index,
			)
		}
		sorted[index] = proximalSortedPoint{
			value: value,
			index: index,
		}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].value == sorted[j].value {
			return sorted[i].index < sorted[j].index
		}
		return sorted[i].value < sorted[j].value
	})

	n := len(sorted)
	z := make([]float64, n)
	for i := range sorted {
		coefficient := float64(2*i - n + 1)
		z[i] = sorted[i].value - lambda*coefficient
	}

	blocks := make([]proximalBlock, 0, n)
	for i, value := range z {
		blocks = append(blocks, proximalBlock{
			start: i,
			end: i,
			weight: 1,
			mean: value,
		})
		for len(blocks) >= 2 {
			a := len(blocks) - 2
			b := len(blocks) - 1
			if blocks[a].mean <= blocks[b].mean {
				break
			}
			weight := blocks[a].weight + blocks[b].weight
			mean :=
				(blocks[a].weight*blocks[a].mean +
					blocks[b].weight*blocks[b].mean) /
					weight
			blocks[a] = proximalBlock{
				start: blocks[a].start,
				end: blocks[b].end,
				weight: weight,
				mean: mean,
			}
			blocks = blocks[:b]
		}
	}

	sortedOutput := make([]float64, n)
	maxBlockResidual := 0.0
	for _, block := range blocks {
		var blockSum float64
		for i := block.start; i <= block.end; i++ {
			sortedOutput[i] = block.mean
			blockSum += z[i] - block.mean
		}
		if residual := math.Abs(blockSum); residual > maxBlockResidual {
			maxBlockResidual = residual
		}
	}

	maxViolation := 0.0
	for i := 1; i < n; i++ {
		if violation := sortedOutput[i-1] - sortedOutput[i]; violation > maxViolation {
			maxViolation = violation
		}
	}

	raw := make([]float64, n)
	for sortedIndex, point := range sorted {
		raw[point.index] = sortedOutput[sortedIndex]
	}

	before, err := completeGraphProxObjective(y, y, lambda)
	if err != nil {
		return nil, ProximalFusionDiagnostics{}, err
	}
	after, err := completeGraphProxObjective(raw, y, lambda)
	if err != nil {
		return nil, ProximalFusionDiagnostics{}, err
	}

	normalized, err := normalizeContinuousOffsets(raw)
	if err != nil {
		return nil, ProximalFusionDiagnostics{}, err
	}
	_, capacity, err := exactOffsetGroups(normalized)
	if err != nil {
		return nil, ProximalFusionDiagnostics{}, err
	}

	return normalized, ProximalFusionDiagnostics{
		Lambda: proximalFusionLambda,
		InputPairPenalty: completeGraphPairPenalty(y),
		OutputPairPenalty: completeGraphPairPenalty(raw),
		ProxObjectiveBefore: before,
		ProxObjectiveAfter: after,
		ExactCapacity: capacity,
		ExactFusionOccurred: capacity > compositeChannels,
		MaximumMonotonicViolation: maxViolation,
		MaximumBlockMeanResidual: maxBlockResidual,
	}, nil
}

func applyProximalFusionGradient(
	offsets, gradient []float64,
) ([]float64, []float64, ProximalFusionDiagnostics, error) {
	proposal, update, err :=
		applyCalibrationGradientNoSticky(offsets, gradient)
	if err != nil {
		return nil, nil, ProximalFusionDiagnostics{}, err
	}
	fused, diagnostics, err :=
		completeGraphFusionProx(
			proposal, proximalFusionLambda,
		)
	if err != nil {
		return nil, nil, ProximalFusionDiagnostics{}, err
	}
	return fused, update, diagnostics, nil
}
