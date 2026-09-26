package unitary

import (
	"fmt"
	"math"
)

const LearnedReadoutProbeSchema = "wingless.unitary-learned-readout-probe.v1"

type LearnedHeadMetric struct {
	Name          string  `json:"name"`
	TrainAccuracy float64 `json:"train_accuracy"`
	TrainLoss     float64 `json:"train_loss"`
}

type LearnedReadoutPathResult struct {
	Name                    string  `json:"name"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MaxRoundTripError       float64 `json:"max_round_trip_error"`
	MaxForwardNormDrift     float64 `json:"max_forward_norm_drift"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
}

type LearnedReadoutProbeResult struct {
	Schema                 string                    `json:"schema"`
	Experiment             string                    `json:"experiment"`
	Dimension              int                       `json:"dimension"`
	Entities               int                       `json:"entities"`
	ValuesPerEntity        int                       `json:"values_per_entity"`
	Scenarios              int                       `json:"scenarios"`
	WritesPerScenario      int                       `json:"writes_per_scenario"`
	TransportDepths        []int                     `json:"transport_depths"`
	NoiseAmplitude         float64                   `json:"noise_amplitude"`
	RuntimePrototypeLookup bool                      `json:"runtime_prototype_lookup"`
	ValueDecoder           LearnedHeadMetric         `json:"value_decoder"`
	RelationHead           LearnedHeadMetric         `json:"relation_head"`
	HeldOutEntityPairs     [][2]int                  `json:"held_out_entity_pairs"`
	Unitary                LearnedReadoutPathResult  `json:"unitary"`
	NonUnitary             LearnedReadoutPathResult  `json:"non_unitary_matched"`
}

type linearSoftmaxHead struct {
	weights [][]float64
	bias    []float64
}

type headSample struct {
	features []float64
	target   int
}

func newLinearSoftmaxHead(classes, features int) linearSoftmaxHead {
	weights := make([][]float64, classes)
	for i := range weights {
		weights[i] = make([]float64, features)
	}
	return linearSoftmaxHead{
		weights: weights,
		bias:    make([]float64, classes),
	}
}

func softmax(logits []float64) ([]float64, error) {
	if len(logits) == 0 {
		return nil, fmt.Errorf("softmax requires logits")
	}
	maxLogit := logits[0]
	for _, value := range logits[1:] {
		if value > maxLogit {
			maxLogit = value
		}
	}
	out := make([]float64, len(logits))
	var total float64
	for i, value := range logits {
		out[i] = math.Exp(value - maxLogit)
		total += out[i]
	}
	if !finite(total) || total <= 0 {
		return nil, fmt.Errorf("softmax normalization is invalid")
	}
	for i := range out {
		out[i] /= total
	}
	return out, nil
}

func (head linearSoftmaxHead) probabilities(features []float64) ([]float64, error) {
	if len(head.weights) == 0 || len(head.weights[0]) != len(features) {
		return nil, fmt.Errorf("readout feature dimension mismatch")
	}
	logits := make([]float64, len(head.weights))
	for class := range head.weights {
		value := head.bias[class]
		for j, feature := range features {
			value += head.weights[class][j] * feature
		}
		logits[class] = value
	}
	return softmax(logits)
}

func classAndMargin(probabilities []float64) (int, float64, error) {
	if len(probabilities) < 2 {
		return 0, 0, fmt.Errorf("classification requires at least two classes")
	}
	best := 0
	second := 1
	if probabilities[second] > probabilities[best] {
		best, second = second, best
	}
	for i := 2; i < len(probabilities); i++ {
		if probabilities[i] > probabilities[best] {
			second = best
			best = i
		} else if probabilities[i] > probabilities[second] {
			second = i
		}
	}
	margin := probabilities[best] - probabilities[second]
	if !finite(margin) {
		return 0, 0, fmt.Errorf("classification margin is not finite")
	}
	return best, margin, nil
}

func evaluateHead(head linearSoftmaxHead, samples []headSample) (float64, float64, error) {
	if len(samples) == 0 {
		return 0, 0, fmt.Errorf("head evaluation requires samples")
	}
	var loss float64
	var correct int
	for _, sample := range samples {
		probabilities, err := head.probabilities(sample.features)
		if err != nil {
			return 0, 0, err
		}
		if sample.target < 0 || sample.target >= len(probabilities) {
			return 0, 0, fmt.Errorf("head target out of range")
		}
		p := probabilities[sample.target]
		if p < 1e-15 {
			p = 1e-15
		}
		loss += -math.Log(p)
		predicted, _, err := classAndMargin(probabilities)
		if err != nil {
			return 0, 0, err
		}
		if predicted == sample.target {
			correct++
		}
	}
	return loss / float64(len(samples)), float64(correct) / float64(len(samples)), nil
}

func trainLinearSoftmax(samples []headSample, classes, features, steps int, learningRate float64) (linearSoftmaxHead, LearnedHeadMetric, error) {
	if len(samples) == 0 || classes < 2 || features < 1 || steps < 1 || learningRate <= 0 {
		return linearSoftmaxHead{}, LearnedHeadMetric{}, fmt.Errorf("invalid learned-head training configuration")
	}
	head := newLinearSoftmaxHead(classes, features)

	for step := 0; step < steps; step++ {
		gradWeights := make([][]float64, classes)
		for i := range gradWeights {
			gradWeights[i] = make([]float64, features)
		}
		gradBias := make([]float64, classes)

		for _, sample := range samples {
			if len(sample.features) != features || sample.target < 0 || sample.target >= classes {
				return linearSoftmaxHead{}, LearnedHeadMetric{}, fmt.Errorf("invalid learned-head sample")
			}
			probabilities, err := head.probabilities(sample.features)
			if err != nil {
				return linearSoftmaxHead{}, LearnedHeadMetric{}, err
			}
			for class := 0; class < classes; class++ {
				delta := probabilities[class]
				if class == sample.target {
					delta -= 1
				}
				gradBias[class] += delta
				for j, feature := range sample.features {
					gradWeights[class][j] += delta * feature
				}
			}
		}

		scale := learningRate / float64(len(samples))
		for class := 0; class < classes; class++ {
			head.bias[class] -= scale * gradBias[class]
			for j := 0; j < features; j++ {
				head.weights[class][j] -= scale * gradWeights[class][j]
			}
		}
	}

	loss, accuracy, err := evaluateHead(head, samples)
	if err != nil {
		return linearSoftmaxHead{}, LearnedHeadMetric{}, err
	}
	return head, LearnedHeadMetric{
		TrainAccuracy: accuracy,
		TrainLoss:     loss,
	}, nil
}

func entityMagnitudeFeatures(state State, entity int) ([]float64, error) {
	if len(state) != 16 || entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("invalid entity feature request")
	}
	features := make([]float64, 4)
	var total float64
	for value := 0; value < 4; value++ {
		index := entity*4 + value
		amplitude := cmplxAbsSquared(state[index])
		features[value] = amplitude
		total += amplitude
	}
	if !finite(total) || total <= 0 {
		return nil, fmt.Errorf("entity feature mass is invalid")
	}
	for i := range features {
		features[i] /= total
	}
	return features, nil
}

func cmplxAbsSquared(value complex128) float64 {
	return real(value)*real(value) + imag(value)*imag(value)
}

func valueHeadTrainingSamples() ([]headSample, error) {
	var samples []headSample
	for value := 0; value < 4; value++ {
		for trial := 0; trial < 32; trial++ {
			features := make([]float64, 4)
			var total float64
			for j := 0; j < 4; j++ {
				base := 0.02 + 0.002*float64((trial+j*3)%7)
				if j == value {
					base += 1
				}
				features[j] = base
				total += base
			}
			for j := range features {
				features[j] /= total
			}
			samples = append(samples, headSample{features: features, target: value})
		}
	}
	return samples, nil
}

func relationFeatures(a, b []float64) ([]float64, error) {
	if len(a) != 4 || len(b) != 4 {
		return nil, fmt.Errorf("relation features require two four-value distributions")
	}
	features := make([]float64, 16)
	for av := 0; av < 4; av++ {
		for bv := 0; bv < 4; bv++ {
			features[av*4+bv] = a[av] * b[bv]
		}
	}
	return features, nil
}

func relationHeadTrainingSamples() ([]headSample, error) {
	var samples []headSample
	for av := 0; av < 4; av++ {
		for bv := 0; bv < 4; bv++ {
			a := make([]float64, 4)
			b := make([]float64, 4)
			a[av] = 1
			b[bv] = 1
			features, err := relationFeatures(a, b)
			if err != nil {
				return nil, err
			}
			target := (av - bv) % 4
			if target < 0 {
				target += 4
			}
			samples = append(samples, headSample{features: features, target: target})
		}
	}
	return samples, nil
}

func decodeTableWithHead(state State, head linearSoftmaxHead) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := entityMagnitudeFeatures(state, entity)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err := head.probabilities(features)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		value, margin, err := classAndMargin(probabilities)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		table[entity] = value
		distributions[entity] = probabilities
		if margin < minMargin {
			minMargin = margin
		}
	}
	return table, distributions, minMargin, nil
}

func inverseStressUnitary(initial State, block []Coupling, depth int) (State, error) {
	out := append(State(nil), initial...)
	inverse := Invert(block)
	var err error
	for i := 0; i < depth; i++ {
		out, err = Propagate(out, inverse)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func inverseStressNonUnitary(initial State, block []Coupling, depth int) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	out := append(State(nil), initial...)
	for step := 0; step < depth; step++ {
		for i := len(block) - 1; i >= 0; i-- {
			coupling := block[i]
			g := coupling.Theta
			denominator := 1 + g*g
			a := out[coupling.A]
			b := out[coupling.B]
			out[coupling.A] = (a + complex(g, 0)*b) / complex(denominator, 0)
			out[coupling.B] = (-complex(g, 0)*a + b) / complex(denominator, 0)
		}
	}
	return out, nil
}

type learnedInverse func(State, []Coupling, int) (State, error)

func makeReadoutScenario(seed int, writes int, depths []int, heldPairs [][2]int) (memoryScenario, error) {
	initial := memoryTable{}
	for entity := 0; entity < 4; entity++ {
		initial[entity] = (seed*seed + entity*2 + seed) % 4
	}
	ops := make([]memoryWrite, 0, writes)
	for j := 0; j < writes; j++ {
		entity := (seed + 3*j + j/2) % 4
		value := (seed*3 + j*j + entity + 1) % 4
		gap := depths[(seed+j*3)%len(depths)]
		ops = append(ops, memoryWrite{entity: entity, value: value, gap: gap})
	}
	pair := heldPairs[seed%len(heldPairs)]
	return memoryScenario{
		initial:  initial,
		writes:   ops,
		queryA:   pair[0],
		queryB:   pair[1],
		finalGap: depths[(seed+writes)%len(depths)],
	}, nil
}

func runLearnedReadoutPath(
	name string,
	scenarios []memoryScenario,
	block []Coupling,
	apply stressApply,
	inverse learnedInverse,
	valueHead linearSoftmaxHead,
	relationHead linearSoftmaxHead,
	noiseAmplitude float64,
) (LearnedReadoutPathResult, error) {
	var commitCorrect, commitTotal int
	var finalExact, relationCorrect int
	var maxRoundTrip, maxForwardNorm float64
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)

	for scenarioIndex, scenario := range scenarios {
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			noisy, err := perturbMemory(canonical, scenarioIndex*1000+writeIndex, noiseAmplitude)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			forward, err := apply(noisy, block, write.gap)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxForwardNorm {
				maxForwardNorm = drift
			}

			recovered, err := inverse(forward, block, write.gap)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			roundTrip, err := L2Distance(noisy, recovered)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			if roundTrip > maxRoundTrip {
				maxRoundTrip = roundTrip
			}

			decoded, _, margin, err := decodeTableWithHead(recovered, valueHead)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			if margin < minValueMargin {
				minValueMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}

			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return LearnedReadoutPathResult{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		noisy, err := perturbMemory(canonical, scenarioIndex*1000+999, noiseAmplitude)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		forward, err := apply(noisy, block, scenario.finalGap)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxForwardNorm {
			maxForwardNorm = drift
		}
		recovered, err := inverse(forward, block, scenario.finalGap)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		roundTrip, err := L2Distance(noisy, recovered)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		if roundTrip > maxRoundTrip {
			maxRoundTrip = roundTrip
		}

		decoded, distributions, margin, err := decodeTableWithHead(recovered, valueHead)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		if margin < minValueMargin {
			minValueMargin = margin
		}
		if decoded == trueTable {
			finalExact++
		}

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil {
			return LearnedReadoutPathResult{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	if commitTotal == 0 || len(scenarios) == 0 || !finite(minValueMargin) || !finite(minRelationMargin) {
		return LearnedReadoutPathResult{}, fmt.Errorf("learned readout path produced no evaluable results")
	}
	return LearnedReadoutPathResult{
		Name:                    name,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalExact) / float64(len(scenarios)),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(len(scenarios)),
		MaxRoundTripError:       maxRoundTrip,
		MaxForwardNormDrift:     maxForwardNorm,
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
	}, nil
}

// RunUP6 removes the 256-state runtime prototype lookup and deterministic
// relation arithmetic used by UP-5.
//
// A learned shared value decoder observes each entity after transport is
// inverted back into the canonical frame. A learned shared relation head maps
// the two decoded value distributions to one of four modular-relation classes.
// The relation head is trained without entity IDs and is evaluated only on
// ordered entity pairs reserved from its integration training examples.
//
// The transport comparison remains unitary versus matched non-unitary.
func RunUP6() (LearnedReadoutProbeResult, error) {
	const (
		scenarios         = 48
		writesPerScenario = 16
		noiseAmplitude    = 0.10
	)
	depths := []int{32, 128, 512, 1024}
	heldPairs := [][2]int{{0, 3}, {3, 0}, {2, 0}, {3, 1}}
	block := stressProgram()

	valueSamples, err := valueHeadTrainingSamples()
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}
	valueHead, valueMetric, err := trainLinearSoftmax(valueSamples, 4, 4, 250, 1.0)
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}
	valueMetric.Name = "shared_value_decoder"

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}
	relationHead, relationMetric, err := trainLinearSoftmax(relationSamples, 4, 16, 400, 1.0)
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}
	relationMetric.Name = "shared_relation_head"

	workload := make([]memoryScenario, 0, scenarios)
	for seed := 0; seed < scenarios; seed++ {
		scenario, err := makeReadoutScenario(seed, writesPerScenario, depths, heldPairs)
		if err != nil {
			return LearnedReadoutProbeResult{}, err
		}
		workload = append(workload, scenario)
	}

	unitaryResult, err := runLearnedReadoutPath(
		"unitary_transport",
		workload,
		block,
		applyStressUnitary,
		inverseStressUnitary,
		valueHead,
		relationHead,
		noiseAmplitude,
	)
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}
	controlResult, err := runLearnedReadoutPath(
		"non_unitary_matched_transport",
		workload,
		block,
		applyStressNonUnitary,
		inverseStressNonUnitary,
		valueHead,
		relationHead,
		noiseAmplitude,
	)
	if err != nil {
		return LearnedReadoutProbeResult{}, err
	}

	return LearnedReadoutProbeResult{
		Schema:                 LearnedReadoutProbeSchema,
		Experiment:             "UP-6-learned-readout-query",
		Dimension:              16,
		Entities:               4,
		ValuesPerEntity:        4,
		Scenarios:              scenarios,
		WritesPerScenario:      writesPerScenario,
		TransportDepths:        append([]int(nil), depths...),
		NoiseAmplitude:         noiseAmplitude,
		RuntimePrototypeLookup: false,
		ValueDecoder:           valueMetric,
		RelationHead:           relationMetric,
		HeldOutEntityPairs:     append([][2]int(nil), heldPairs...),
		Unitary:                unitaryResult,
		NonUnitary:             controlResult,
	}, nil
}
