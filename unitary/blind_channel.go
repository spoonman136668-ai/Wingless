package unitary

import (
	"fmt"
	"math"
)

const BlindChannelSchema = "wingless.unitary-blind-channel.v1"

type ChannelMatrix [6][6]float64

type BlindChannelTrace struct {
	Outer             int           `json:"outer"`
	Loss              float64       `json:"loss"`
	MeanRoleAlignment float64       `json:"mean_role_alignment"`
	MinRoleAlignment  float64       `json:"min_role_alignment"`
	Demixer           ChannelMatrix `json:"demixer"`
}

type BlindChannelStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type BlindChannelIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type BlindChannelDiagnosis struct {
	ChannelLearningPass      bool    `json:"channel_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	InitialCapacityAccuracy  float64 `json:"initial_capacity_accuracy"`
	FinalTrainAccuracy       float64 `json:"final_train_accuracy"`
	FinalHeldOutAccuracy     float64 `json:"final_heldout_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	MeanRoleAlignment        float64 `json:"mean_role_alignment"`
	MinimumRoleAlignment     float64 `json:"minimum_role_alignment"`
	DemixerL2Shift           float64 `json:"demixer_l2_shift"`
}

type BlindChannelProbeResult struct {
	Schema                       string                  `json:"schema"`
	Experiment                   string                  `json:"experiment"`
	BaseDimension                int                     `json:"base_dimension"`
	CompositeDimension           int                     `json:"composite_dimension"`
	AnonymousChannels            int                     `json:"anonymous_channels"`
	RuntimeStateObjects          int                     `json:"runtime_state_objects"`
	SemanticChannelLocationsKnown bool                   `json:"semantic_channel_locations_known"`
	KnownMixerExposedToLearner   bool                    `json:"known_mixer_exposed_to_learner"`
	RuntimeMixerInverseApplied   bool                    `json:"runtime_mixer_inverse_applied"`
	LearnedReadoutProjection     bool                    `json:"learned_readout_projection"`
	RuntimePrototypeLookup       bool                    `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout       bool                    `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided        bool                    `json:"explicit_depth_provided"`
	GlobalPhaseNuisance          bool                    `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude         float64                 `json:"memory_noise_amplitude"`
	LearningTables               int                     `json:"channel_learning_tables"`
	TrainTables                  int                     `json:"train_tables"`
	HeldOutTables                int                     `json:"held_out_tables"`
	TrainDepths                  []int                   `json:"train_depths"`
	HeldOutDepths                []int                   `json:"held_out_depths"`
	FixedMixer                   ChannelMatrix           `json:"fixed_anonymous_channel_mixer"`
	InitialDemixer               ChannelMatrix           `json:"initial_demixer"`
	LearnedDemixer               ChannelMatrix           `json:"learned_demixer"`
	InitialLoss                  float64                 `json:"initial_loss"`
	FinalLoss                    float64                 `json:"final_loss"`
	Trace                        []BlindChannelTrace     `json:"trace"`
	Unitary                      BlindChannelStaticResult `json:"unitary"`
	MatchedControl               BlindChannelStaticResult `json:"matched_control"`
	Integration                  BlindChannelIntegration `json:"unitary_mutable_integration"`
	Diagnosis                    BlindChannelDiagnosis   `json:"diagnosis"`
}

func identityChannelMatrix() ChannelMatrix {
	var out ChannelMatrix
	for i := 0; i < 6; i++ {
		out[i][i] = 1
	}
	return out
}

// denseChannelMixer is a deterministic dense real orthogonal matrix.
//
// It is built by applying one nontrivial Givens rotation for every channel
// pair. The resulting transform is strongly scrambled: every output channel
// contains a nonzero contribution from every historical semantic channel.
// The learner is never given this matrix.
func denseChannelMixer() ChannelMatrix {
	out := identityChannelMatrix()
	index := 0
	for first := 0; first < 6; first++ {
		for second := first + 1; second < 6; second++ {
			theta := 0.43 + 0.071*float64(index)
			c := math.Cos(theta)
			s := math.Sin(theta)

			beforeFirst := out[first]
			beforeSecond := out[second]
			for column := 0; column < 6; column++ {
				out[first][column] =
					c*beforeFirst[column] -
						s*beforeSecond[column]
				out[second][column] =
					s*beforeFirst[column] +
						c*beforeSecond[column]
			}
			index++
		}
	}
	return out
}

func channelMatrixProduct(a, b ChannelMatrix) ChannelMatrix {
	var out ChannelMatrix
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			for k := 0; k < 6; k++ {
				out[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return out
}

func normalizeChannelRow(row *[6]float64) error {
	if row == nil {
		return fmt.Errorf("channel row is nil")
	}
	var norm2 float64
	for _, value := range row {
		if !finite(value) {
			return fmt.Errorf("channel row contains non-finite value")
		}
		norm2 += value * value
	}
	if norm2 <= 1e-18 {
		return fmt.Errorf("channel row norm too small")
	}
	scale := 1 / math.Sqrt(norm2)
	for i := range row {
		row[i] *= scale
	}
	return nil
}

func normalizedChannelMatrix(matrix ChannelMatrix) (ChannelMatrix, error) {
	out := matrix
	for row := 0; row < 6; row++ {
		if err := normalizeChannelRow(&out[row]); err != nil {
			return ChannelMatrix{}, err
		}
	}
	return out, nil
}

func applyChannelMatrix(state State, matrix ChannelMatrix) (State, error) {
	if len(state) != compositeDimension {
		return nil, fmt.Errorf(
			"channel-matrix state dimension=%d want=%d",
			len(state), compositeDimension,
		)
	}
	out := make(State, compositeDimension)
	for coordinate := 0; coordinate < compositeChannelDim; coordinate++ {
		for row := 0; row < 6; row++ {
			var value complex128
			for column := 0; column < 6; column++ {
				value += complex(matrix[row][column], 0) *
					state[column*compositeChannelDim+coordinate]
			}
			out[row*compositeChannelDim+coordinate] = value
		}
	}
	return out, nil
}

func mixerOrthogonalityError(matrix ChannelMatrix) float64 {
	var maximum float64
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			var dot float64
			for k := 0; k < 6; k++ {
				dot += matrix[i][k] * matrix[j][k]
			}
			want := 0.0
			if i == j {
				want = 1
			}
			delta := math.Abs(dot - want)
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum
}

func roleAlignments(demixer, mixer ChannelMatrix) []float64 {
	composed := channelMatrixProduct(demixer, mixer)
	out := make([]float64, 6)
	for row := 0; row < 6; row++ {
		var total float64
		for column := 0; column < 6; column++ {
			total += composed[row][column] * composed[row][column]
		}
		if total > 0 {
			out[row] = composed[row][row] * composed[row][row] / total
		}
	}
	return out
}

func blindMeanMin(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	minimum := values[0]
	var total float64
	for _, value := range values {
		total += value
		if value < minimum {
			minimum = value
		}
	}
	return total / float64(len(values)), minimum
}

func channelMatrixL2Shift(a, b ChannelMatrix) float64 {
	var total float64
	for row := 0; row < 6; row++ {
		for column := 0; column < 6; column++ {
			delta := a[row][column] - b[row][column]
			total += delta * delta
		}
	}
	return math.Sqrt(total)
}

func blindMixedComposite(memory State) (State, error) {
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return nil, err
	}
	composite, err := packCompositeState(memory, bank)
	if err != nil {
		return nil, err
	}
	return applyChannelMatrix(composite, denseChannelMixer())
}

func blindFeature(state State, demixer ChannelMatrix, entity int) ([]float64, error) {
	recovered, err := applyChannelMatrix(state, demixer)
	if err != nil {
		return nil, err
	}
	return compositeEntityFeature(recovered, entity)
}

func buildBlindCanonicalSamples(
	tables []memoryTable,
	entity int,
	demixer ChannelMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, error) {
	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for trial := 0; trial < trials; trial++ {
			seed := seedOffset +
				memoryTableIndex(table)*100 +
				entity*17 +
				trial
			memory, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return nil, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.203*float64(seed+1), 2*math.Pi),
			)
			mixed, err := blindMixedComposite(memory)
			if err != nil {
				return nil, err
			}
			features, err := blindFeature(
				mixed, demixer, entity,
			)
			if err != nil {
				return nil, err
			}
			samples = append(samples, headSample{
				features: features,
				target:   table[entity],
			})
		}
	}
	return samples, nil
}

func fitBlindHeads(
	tables []memoryTable,
	demixer ChannelMatrix,
	memoryNoise float64,
	trials int,
	steps int,
) ([4]linearSoftmaxHead, float64, error) {
	var heads [4]linearSoftmaxHead
	var average float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildBlindCanonicalSamples(
			tables, entity, demixer,
			memoryNoise, trials, 0,
		)
		if err != nil {
			return heads, 0, err
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, 0, err
		}
		heads[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return heads, 0, err
		}
		average += accuracy
	}
	return heads, average / 4, nil
}

func blindHeadsLoss(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	demixer ChannelMatrix,
	memoryNoise float64,
	trials int,
) (float64, error) {
	var total float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildBlindCanonicalSamples(
			tables, entity, demixer,
			memoryNoise, trials, 0,
		)
		if err != nil {
			return 0, err
		}
		loss, _, err := evaluateHead(
			heads[entity], samples,
		)
		if err != nil {
			return 0, err
		}
		total += loss
	}
	return total / 4, nil
}

func learnBlindDemixer(
	tables []memoryTable,
	memoryNoise float64,
) (
	initial ChannelMatrix,
	learned ChannelMatrix,
	initialLoss float64,
	finalLoss float64,
	trace []BlindChannelTrace,
	err error,
) {
	const (
		trials             = 1
		outerSteps         = 70
		headStepsPerOuter  = 10
		headLearningRate   = 1.0
		demixLearningRate  = 0.22
		demixEpsilon       = 1e-4
	)
	demixer := identityChannelMatrix()
	initial = demixer

	var heads [4]linearSoftmaxHead
	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(4, 2)
	}

	for entity := 0; entity < 4; entity++ {
		samples, e := buildBlindCanonicalSamples(
			tables, entity, demixer,
			memoryNoise, trials, 0,
		)
		if e != nil {
			err = e
			return
		}
		for step := 0; step < headStepsPerOuter; step++ {
			if e := updateSoftmaxHead(
				&heads[entity], samples, headLearningRate,
			); e != nil {
				err = e
				return
			}
		}
	}
	initialLoss, err = blindHeadsLoss(
		heads, tables, demixer,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}

	for outer := 1; outer <= outerSteps; outer++ {
		for entity := 0; entity < 4; entity++ {
			samples, e := buildBlindCanonicalSamples(
				tables, entity, demixer,
				memoryNoise, trials, 0,
			)
			if e != nil {
				err = e
				return
			}
			for step := 0; step < headStepsPerOuter; step++ {
				if e := updateSoftmaxHead(
					&heads[entity],
					samples,
					headLearningRate,
				); e != nil {
					err = e
					return
				}
			}
		}

		var gradient ChannelMatrix
		for row := 0; row < 6; row++ {
			for column := 0; column < 6; column++ {
				plus := demixer
				minus := demixer
				plus[row][column] += demixEpsilon
				minus[row][column] -= demixEpsilon
				plusLoss, e := blindHeadsLoss(
					heads, tables, plus,
					memoryNoise, trials,
				)
				if e != nil {
					err = e
					return
				}
				minusLoss, e := blindHeadsLoss(
					heads, tables, minus,
					memoryNoise, trials,
				)
				if e != nil {
					err = e
					return
				}
				gradient[row][column] =
					(plusLoss - minusLoss) /
						(2 * demixEpsilon)
				if !finite(gradient[row][column]) {
					err = fmt.Errorf(
						"demixer gradient row=%d col=%d non-finite",
						row, column,
					)
					return
				}
			}
		}

		for row := 0; row < 6; row++ {
			for column := 0; column < 6; column++ {
				demixer[row][column] -=
					demixLearningRate *
						gradient[row][column]
			}
			if e := normalizeChannelRow(&demixer[row]); e != nil {
				err = e
				return
			}
		}

		if outer == 1 || outer == 10 ||
			outer == 20 || outer == 40 ||
			outer == 55 || outer == 70 {
			loss, e := blindHeadsLoss(
				heads, tables, demixer,
				memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			alignments := roleAlignments(
				demixer, denseChannelMixer(),
			)
			mean, minimum := blindMeanMin(alignments)
			trace = append(trace, BlindChannelTrace{
				Outer:             outer,
				Loss:              loss,
				MeanRoleAlignment: mean,
				MinRoleAlignment:  minimum,
				Demixer:           demixer,
			})
		}
	}

	finalLoss, err = blindHeadsLoss(
		heads, tables, demixer,
		memoryNoise, trials,
	)
	if err != nil {
		return
	}
	learned = demixer
	return
}

func buildBlindTransportSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	demixer ChannelMatrix,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]headSample, float64, error) {
	var samples []headSample
	var maxNormDrift float64
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					entity*101 +
					trial*17
				memory, err := perturbMemory(
					canonical, seed, memoryNoise,
				)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.203*float64(seed+1), 2*math.Pi),
				)
				mixed, err := blindMixedComposite(memory)
				if err != nil {
					return nil, 0, err
				}
				mixed, err = evolveCompositeState(
					mixed, block, depth, apply,
				)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(mixed)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}
				features, err := blindFeature(
					mixed, demixer, entity,
				)
				if err != nil {
					return nil, 0, err
				}
				samples = append(samples, headSample{
					features: features,
					target:   table[entity],
				})
			}
		}
	}
	return samples, maxNormDrift, nil
}

func trainBlindTransportHeads(
	name, path string,
	demixer ChannelMatrix,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	steps int,
) ([4]linearSoftmaxHead, BlindChannelStaticResult, error) {
	var heads [4]linearSoftmaxHead
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)
	var maxNormDrift float64

	for entity := 0; entity < 4; entity++ {
		trainSamples, trainDrift, err := buildBlindTransportSamples(
			trainTables, trainDepths, entity,
			demixer, block, apply,
			memoryNoise, 4, 0,
		)
		if err != nil {
			return heads, BlindChannelStaticResult{}, err
		}
		heldSamples, heldDrift, err := buildBlindTransportSamples(
			heldTables, heldDepths, entity,
			demixer, block, apply,
			memoryNoise, 2, 7000000,
		)
		if err != nil {
			return heads, BlindChannelStaticResult{}, err
		}
		if trainDrift > maxNormDrift {
			maxNormDrift = trainDrift
		}
		if heldDrift > maxNormDrift {
			maxNormDrift = heldDrift
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, steps, 1.0,
		)
		if err != nil {
			return heads, BlindChannelStaticResult{}, err
		}
		heads[entity] = head
		_, trainAccuracy, err := evaluateHead(
			head, trainSamples,
		)
		if err != nil {
			return heads, BlindChannelStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(
			head, heldSamples,
		)
		if err != nil {
			return heads, BlindChannelStaticResult{}, err
		}
		perTrain[entity] = trainAccuracy
		perHeld[entity] = heldAccuracy
	}

	var trainAverage, heldAverage float64
	for entity := 0; entity < 4; entity++ {
		trainAverage += perTrain[entity]
		heldAverage += perHeld[entity]
	}
	trainAverage /= 4
	heldAverage /= 4

	return heads, BlindChannelStaticResult{
		Name:             name,
		Path:             path,
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
		MaxNormDrift:     maxNormDrift,
	}, nil
}

func decodeBlindTable(
	state State,
	demixer ChannelMatrix,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := blindFeature(
			state, demixer, entity,
		)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err :=
			heads[entity].probabilities(features)
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

func runBlindIntegration(
	demixer ChannelMatrix,
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise float64,
) (BlindChannelIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return BlindChannelIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return BlindChannelIntegration{}, err
	}

	var commitCorrect, commitTotal int
	var finalCorrect, relationCorrect int
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)
	var maxNormDrift float64

	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeFullRankScenario(
			scenarioIndex,
			heldTables[(scenarioIndex*7)%len(heldTables)],
			writes,
			depths,
		)
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			seed := 19000000 +
				scenarioIndex*10000 +
				writeIndex*31
			memory, err := perturbMemory(
				canonical, seed, memoryNoise,
			)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.203*float64(seed+1), 2*math.Pi),
			)
			mixed, err := blindMixedComposite(memory)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			mixed, err = evolveCompositeState(
				mixed, block, write.gap,
				applyStressUnitary,
			)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			norm2, err := NormSquared(mixed)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeBlindTable(
				mixed, demixer, heads,
			)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			if margin < minValueMargin {
				minValueMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}
			pathTable, err = applyMemoryWrite(
				decoded, write.entity, write.value,
			)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return BlindChannelIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		seed := 19000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(
			canonical, seed, memoryNoise,
		)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.203*float64(seed+1), 2*math.Pi),
		)
		mixed, err := blindMixedComposite(memory)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		mixed, err = evolveCompositeState(
			mixed, block, scenario.finalGap,
			applyStressUnitary,
		)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		norm2, err := NormSquared(mixed)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		decoded, distributions, margin, err := decodeBlindTable(
			mixed, demixer, heads,
		)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		if margin < minValueMargin {
			minValueMargin = margin
		}
		if decoded == trueTable {
			finalCorrect++
		}

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return BlindChannelIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return BlindChannelIntegration{
		Scenarios:               scenarios,
		WritesPerScenario:       writes,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect) / float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(scenarios),
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
		MaxNormDrift:            maxNormDrift,
	}, nil
}

// RunUP17 hides the semantic channel locations behind a dense orthogonal
// channel mixer.
//
// The encoder knows the fixed mixer because it creates the anonymous composite
// state. The task learner does not receive the mixer, its inverse, or target
// demixing coefficients. A trainable 6x6 readout projection begins as identity
// and is updated only through classification loss.
//
// The learned projection is a readout operation only. The recurrent runtime
// remains one mixed 96-D state throughout transport.
func RunUP17() (BlindChannelProbeResult, error) {
	const memoryNoise = 0.05
	learningTables, err := selectObserverTables(true, 32)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}

	mixer := denseChannelMixer()
	if mixerOrthogonalityError(mixer) > 1e-12 {
		return BlindChannelProbeResult{}, fmt.Errorf(
			"dense channel mixer is not orthogonal",
		)
	}

	initial := identityChannelMatrix()
	_, initialCapacity, err := fitBlindHeads(
		learningTables, initial,
		memoryNoise, 2, 900,
	)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}

	initialDemixer, learnedDemixer, initialLoss, finalLoss, trace, err :=
		learnBlindDemixer(
			learningTables, memoryNoise,
		)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}

	alignments := roleAlignments(
		learnedDemixer, mixer,
	)
	meanAlignment, minimumAlignment :=
		blindMeanMin(alignments)
	demixerShift := channelMatrixL2Shift(
		initialDemixer, learnedDemixer,
	)

	block := stressProgram()
	unitaryHeads, unitaryResult, err :=
		trainBlindTransportHeads(
			"unitary_blind_anonymous_channels",
			"unitary",
			learnedDemixer,
			trainTables,
			heldTables,
			trainDepths,
			heldDepths,
			block,
			applyStressUnitary,
			memoryNoise,
			1200,
		)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}

	_, controlResult, err :=
		trainBlindTransportHeads(
			"non_unitary_blind_anonymous_channels",
			"non_unitary_matched",
			learnedDemixer,
			trainTables,
			heldTables,
			trainDepths,
			heldDepths,
			block,
			applyStressNonUnitary,
			memoryNoise,
			1200,
		)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}

	integration, err := runBlindIntegration(
		learnedDemixer,
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return BlindChannelProbeResult{}, err
	}

	channelPass :=
		initialCapacity < 0.70 &&
			unitaryResult.TrainAccuracy >= 0.99 &&
			finalLoss < initialLoss &&
			demixerShift >= 0.50
	unseenPass :=
		unitaryResult.HeldOutAccuracy >= 0.99 &&
			unitaryResult.MaxNormDrift <= 1e-12
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-12

	return BlindChannelProbeResult{
		Schema:                        BlindChannelSchema,
		Experiment:                    "UP-17-blind-anonymous-channel-recovery",
		BaseDimension:                 compositeChannelDim,
		CompositeDimension:            compositeDimension,
		AnonymousChannels:             6,
		RuntimeStateObjects:           1,
		SemanticChannelLocationsKnown: false,
		KnownMixerExposedToLearner:    false,
		RuntimeMixerInverseApplied:    false,
		LearnedReadoutProjection:      true,
		RuntimePrototypeLookup:        false,
		ExplicitInverseReadout:        false,
		ExplicitDepthProvided:         false,
		GlobalPhaseNuisance:           true,
		MemoryNoiseAmplitude:          memoryNoise,
		LearningTables:                len(learningTables),
		TrainTables:                   len(trainTables),
		HeldOutTables:                 len(heldTables),
		TrainDepths:                   append([]int(nil), trainDepths...),
		HeldOutDepths:                 append([]int(nil), heldDepths...),
		FixedMixer:                    mixer,
		InitialDemixer:                initialDemixer,
		LearnedDemixer:                learnedDemixer,
		InitialLoss:                   initialLoss,
		FinalLoss:                     finalLoss,
		Trace:                         trace,
		Unitary:                       unitaryResult,
		MatchedControl:                controlResult,
		Integration:                   integration,
		Diagnosis: BlindChannelDiagnosis{
			ChannelLearningPass:     channelPass,
			UnseenDepthPass:         unseenPass,
			MutableIntegrationPass:  mutablePass,
			InitialCapacityAccuracy: initialCapacity,
			FinalTrainAccuracy:      unitaryResult.TrainAccuracy,
			FinalHeldOutAccuracy:    unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:  controlResult.HeldOutAccuracy,
			MeanRoleAlignment:       meanAlignment,
			MinimumRoleAlignment:    minimumAlignment,
			DemixerL2Shift:          demixerShift,
		},
	}, nil
}
