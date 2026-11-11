package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const DiscoveredCommutantSchema = "wingless.unitary-discovered-commutant.v1"

const (
	discoveredObservableCount    = 32
	discoveredProjectionRounds   = 20
	discoveredEffectiveOrbitSize = 1 << discoveredProjectionRounds
	discoveredRawFeatureDim      = discoveredObservableCount * 2
	discoveredQuadraticFeatureDim = discoveredRawFeatureDim +
		(discoveredRawFeatureDim*(discoveredRawFeatureDim+1))/2
	discoveredRidgeLambda = 1e-6
)

type discoveredRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type DiscoveredStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MeanPhaseCosine  float64   `json:"mean_phase_cosine"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type DiscoveredIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type DiscoveredDiagnosis struct {
	DiscoveryCommutatorPass   bool    `json:"discovery_commutator_pass"`
	FeatureInvariancePass     bool    `json:"feature_invariance_pass"`
	PhaseCodeLearningPass     bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass           bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass    bool    `json:"mutable_integration_pass"`
	MaxCommutatorEntryError   float64 `json:"max_commutator_entry_error"`
	MaxFeatureDrift           float64 `json:"max_feature_drift"`
	MinPreNormalizationFrobenius float64 `json:"min_pre_normalization_frobenius"`
	TrainAccuracy             float64 `json:"train_accuracy"`
	HeldOutAccuracy           float64 `json:"held_out_accuracy"`
	MatchedControlAccuracy    float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine       float64 `json:"mean_held_phase_cosine"`
}

type DiscoveredCommutantProbeResult struct {
	Schema                           string                    `json:"schema"`
	Experiment                       string                    `json:"experiment"`
	LatentDimension                  int                       `json:"latent_dimension"`
	RuntimeStateObjects              int                       `json:"runtime_state_objects"`
	VisibleChannelBlocks             bool                      `json:"visible_channel_blocks"`
	FullCoordinateMixing             bool                      `json:"full_coordinate_mixing"`
	ObservableDiscoveryFromTransport bool                      `json:"observable_discovery_from_transport"`
	DiscoveryUsesHiddenMultiplicity  bool                      `json:"discovery_uses_hidden_multiplicity"`
	DiscoveryUsesHiddenMixer         bool                      `json:"discovery_uses_hidden_mixer"`
	TransportAdjointUsedInDiscovery  bool                      `json:"transport_adjoint_used_in_discovery"`
	RuntimeAdjointApplied            bool                      `json:"runtime_adjoint_applied"`
	DiscoverySeedCount               int                       `json:"discovery_seed_count"`
	DiscoveryProjectionRounds        int                       `json:"discovery_projection_rounds"`
	DiscoveryEffectiveOrbitSize      int                       `json:"discovery_effective_orbit_size"`
	KnownFullMixerExposedToLearner   bool                      `json:"known_full_mixer_exposed_to_learner"`
	RuntimeUnmixApplied              bool                      `json:"runtime_unmix_applied"`
	LearnedDemixerUsed               bool                      `json:"learned_demixer_used"`
	PhaseAlphabetSupervision         bool                      `json:"phase_alphabet_supervision"`
	PhaseCodeDimension               int                       `json:"phase_code_dimension"`
	RawFeatureDimension              int                       `json:"raw_feature_dimension"`
	QuadraticFeatureDimension        int                       `json:"quadratic_feature_dimension"`
	RuntimePrototypeLookup           bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout  bool                      `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided            bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance              bool                      `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude             float64                   `json:"memory_noise_amplitude"`
	TrainingDepths                   []int                     `json:"training_depths"`
	HeldOutDepths                    []int                     `json:"held_out_depths"`
	TrainTables                      int                       `json:"train_tables"`
	HeldOutTables                    int                       `json:"held_out_tables"`
	Unitary                          DiscoveredStaticResult    `json:"unitary"`
	MatchedControl                   DiscoveredStaticResult    `json:"matched_control"`
	Integration                      DiscoveredIntegration     `json:"unitary_mutable_integration"`
	Diagnosis                        DiscoveredDiagnosis       `json:"diagnosis"`
}

type discoveredRow struct {
	features []float64
	table    memoryTable
}

func cloneLatentMatrix(matrix latentMatrix) latentMatrix {
	out := make(latentMatrix, len(matrix))
	for row := range matrix {
		out[row] = append([]complex128(nil), matrix[row]...)
	}
	return out
}

func deterministicDiscoveryVector(seed, variant int) (State, error) {
	out := make(State, fullLatentDimension)
	for coordinate := 0; coordinate < fullLatentDimension; coordinate++ {
		x := float64(coordinate + 1)
		s := float64(seed + 1)
		v := float64(variant + 1)
		amplitude := 0.85 + 0.15*math.Sin(
			0.173*s*x + 0.037*v*x*x,
		)
		phase :=
			0.113*s*x +
				0.017*(s+2*v)*x*x +
				0.00071*(s+v)*x*x*x
		out[coordinate] = cmplx.Rect(amplitude, phase)
	}
	return Normalize(out)
}

func discoverySeedMatrix(seed int) (latentMatrix, error) {
	left, err := deterministicDiscoveryVector(seed, 0)
	if err != nil {
		return nil, err
	}
	right, err := deterministicDiscoveryVector(seed, 1)
	if err != nil {
		return nil, err
	}
	out := make(latentMatrix, fullLatentDimension)
	for row := 0; row < fullLatentDimension; row++ {
		out[row] = make([]complex128, fullLatentDimension)
		for column := 0; column < fullLatentDimension; column++ {
			out[row][column] = left[row] * cmplx.Conj(right[column])
		}
	}
	return out, nil
}

func matrixFrobeniusNorm(matrix latentMatrix) (float64, error) {
	if len(matrix) != fullLatentDimension {
		return 0, fmt.Errorf("frobenius matrix dimension mismatch")
	}
	var sum float64
	for row := 0; row < fullLatentDimension; row++ {
		if len(matrix[row]) != fullLatentDimension {
			return 0, fmt.Errorf("frobenius row dimension mismatch")
		}
		for column := 0; column < fullLatentDimension; column++ {
			magnitude := cmplx.Abs(matrix[row][column])
			sum += magnitude * magnitude
		}
	}
	if !finite(sum) || sum <= 0 {
		return 0, fmt.Errorf("invalid frobenius norm square=%g", sum)
	}
	return math.Sqrt(sum), nil
}

func averageLatentMatrices(a, b latentMatrix) (latentMatrix, error) {
	if len(a) != fullLatentDimension || len(b) != fullLatentDimension {
		return nil, fmt.Errorf("average matrix dimension mismatch")
	}
	out := make(latentMatrix, fullLatentDimension)
	for row := 0; row < fullLatentDimension; row++ {
		if len(a[row]) != fullLatentDimension || len(b[row]) != fullLatentDimension {
			return nil, fmt.Errorf("average matrix row dimension mismatch")
		}
		out[row] = make([]complex128, fullLatentDimension)
		for column := 0; column < fullLatentDimension; column++ {
			out[row][column] = 0.5 * (a[row][column] + b[row][column])
		}
	}
	return out, nil
}

func normalizeLatentMatrix(matrix latentMatrix) (latentMatrix, float64, error) {
	norm, err := matrixFrobeniusNorm(matrix)
	if err != nil {
		return nil, 0, err
	}
	out := cloneLatentMatrix(matrix)
	for row := range out {
		for column := range out[row] {
			out[row][column] /= complex(norm, 0)
		}
	}
	return out, norm, nil
}

func discoveryConjugationPowers(
	step latentMatrix,
	rounds int,
) ([]latentMatrix, []latentMatrix, error) {
	if rounds < 1 {
		return nil, nil, fmt.Errorf("discovery rounds must be positive")
	}
	powers := make([]latentMatrix, rounds)
	adjoints := make([]latentMatrix, rounds)
	current := cloneLatentMatrix(step)
	for round := 0; round < rounds; round++ {
		powers[round] = current
		adjoint, err := latentAdjoint(current)
		if err != nil {
			return nil, nil, err
		}
		adjoints[round] = adjoint
		if round+1 < rounds {
			next, err := latentMatrixMultiply(current, current)
			if err != nil {
				return nil, nil, err
			}
			current = next
		}
	}
	return powers, adjoints, nil
}

func discoverCommutingObservables(
	step latentMatrix,
	count, rounds int,
) ([]latentMatrix, float64, error) {
	powers, adjoints, err := discoveryConjugationPowers(step, rounds)
	if err != nil {
		return nil, 0, err
	}
	out := make([]latentMatrix, 0, count)
	minPreNorm := math.Inf(1)

	for seedIndex := 0; seedIndex < count; seedIndex++ {
		current, err := discoverySeedMatrix(seedIndex)
		if err != nil {
			return nil, 0, err
		}
		for round := 0; round < rounds; round++ {
			left, err := latentMatrixMultiply(powers[round], current)
			if err != nil {
				return nil, 0, err
			}
			rotated, err := latentMatrixMultiply(left, adjoints[round])
			if err != nil {
				return nil, 0, err
			}
			current, err = averageLatentMatrices(current, rotated)
			if err != nil {
				return nil, 0, err
			}
		}
		normalized, preNorm, err := normalizeLatentMatrix(current)
		if err != nil {
			return nil, 0, err
		}
		if preNorm < minPreNorm {
			minPreNorm = preNorm
		}
		out = append(out, normalized)
	}
	return out, minPreNorm, nil
}

func discoveredRawFeatures(
	state State,
	observables []latentMatrix,
) ([]float64, error) {
	if len(observables) != discoveredObservableCount {
		return nil, fmt.Errorf(
			"discovered observable count=%d want=%d",
			len(observables), discoveredObservableCount,
		)
	}
	out := make([]float64, 0, discoveredRawFeatureDim)
	for _, observable := range observables {
		applied, err := latentMatrixVector(observable, state)
		if err != nil {
			return nil, err
		}
		value, err := stateInner(state, applied)
		if err != nil {
			return nil, err
		}
		out = append(out, real(value), imag(value))
	}
	for _, value := range out {
		if !finite(value) {
			return nil, fmt.Errorf("discovered raw feature is non-finite")
		}
	}
	return out, nil
}

func discoveredQuadraticFeatures(
	state State,
	observables []latentMatrix,
) ([]float64, error) {
	raw, err := discoveredRawFeatures(state, observables)
	if err != nil {
		return nil, err
	}
	out := make([]float64, 0, discoveredQuadraticFeatureDim)
	out = append(out, raw...)
	for first := 0; first < len(raw); first++ {
		for second := first; second < len(raw); second++ {
			value := raw[first] * raw[second]
			if !finite(value) {
				return nil, fmt.Errorf("discovered quadratic feature is non-finite")
			}
			out = append(out, value)
		}
	}
	if len(out) != discoveredQuadraticFeatureDim {
		return nil, fmt.Errorf(
			"discovered quadratic feature dimension=%d want=%d",
			len(out), discoveredQuadraticFeatureDim,
		)
	}
	return out, nil
}

func discoveredFeatureDrift(
	mixer latentMatrix,
	observables []latentMatrix,
	depthOperator latentMatrix,
) (float64, error) {
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil { return 0, err }
	memory, err = perturbMemory(memory, 252525, 0.05)
	if err != nil { return 0, err }
	memory = rotateGlobalPhase(memory, 1.731)
	state, err := fullLatentEncode(memory, mixer)
	if err != nil { return 0, err }
	before, err := discoveredQuadraticFeatures(state, observables)
	if err != nil { return 0, err }
	evolved, err := latentMatrixVector(depthOperator, state)
	if err != nil { return 0, err }
	after, err := discoveredQuadraticFeatures(evolved, observables)
	if err != nil { return 0, err }
	maximum := 0.0
	for i := range before {
		delta := math.Abs(before[i] - after[i])
		if delta > maximum { maximum = delta }
	}
	return maximum, nil
}

func buildDiscoveredRows(
	tables []memoryTable,
	depths []int,
	operators map[int]latentMatrix,
	mixer latentMatrix,
	observables []latentMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]discoveredRow, float64, error) {
	rows := make([]discoveredRow, 0, len(tables)*len(depths)*trials)
	var maxNormDrift float64
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil { return nil, 0, err }
		for depthIndex, depth := range depths {
			operator, ok := operators[depth]
			if !ok { return nil, 0, fmt.Errorf("missing discovered depth=%d", depth) }
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil { return nil, 0, err }
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.271*float64(seed+1), 2*math.Pi),
				)
				state, err := fullLatentEncode(memory, mixer)
				if err != nil { return nil, 0, err }
				state, err = latentMatrixVector(operator, state)
				if err != nil { return nil, 0, err }
				norm2, err := NormSquared(state)
				if err != nil { return nil, 0, err }
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift { maxNormDrift = drift }
				features, err := discoveredQuadraticFeatures(state, observables)
				if err != nil { return nil, 0, err }
				rows = append(rows, discoveredRow{
					features: features,
					table: table,
				})
			}
		}
	}
	return rows, maxNormDrift, nil
}

func trainDiscoveredRegressors(
	rows []discoveredRow,
	lambda float64,
) ([4]discoveredRegressor, error) {
	var regressors [4]discoveredRegressor
	if len(rows) == 0 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid discovered ridge configuration")
	}
	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)
	for row := 0; row < n; row++ {
		if len(rows[row].features) != discoveredQuadraticFeatureDim {
			return regressors, fmt.Errorf("discovered feature dimension mismatch")
		}
		gram[row] = make([]float64, n)
		right[row] = make([]float64, outputs)
		for entity := 0; entity < 4; entity++ {
			target, err := learnedPhaseTarget(rows[row].table[entity])
			if err != nil { return regressors, err }
			right[row][entity*2] = target[0]
			right[row][entity*2+1] = target[1]
		}
	}
	for first := 0; first < n; first++ {
		value := 1.0 + lambda
		for feature := 0; feature < discoveredQuadraticFeatureDim; feature++ {
			value += rows[first].features[feature] * rows[first].features[feature]
		}
		gram[first][first] = value
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < discoveredQuadraticFeatureDim; feature++ {
				value += rows[first].features[feature] * rows[second].features[feature]
			}
			if !finite(value) {
				return regressors, fmt.Errorf("discovered gram value is non-finite")
			}
			gram[first][second] = value
			gram[second][first] = value
		}
	}
	alpha, err := solveDenseMultiple(gram, right)
	if err != nil { return regressors, err }

	for entity := 0; entity < 4; entity++ {
		for output := 0; output < 2; output++ {
			regressors[entity].weights[output] =
				make([]float64, discoveredQuadraticFeatureDim)
		}
	}
	for row := 0; row < n; row++ {
		for entity := 0; entity < 4; entity++ {
			for output := 0; output < 2; output++ {
				coefficient := alpha[row][entity*2+output]
				regressors[entity].bias[output] += coefficient
				for feature, value := range rows[row].features {
					regressors[entity].weights[output][feature] +=
						coefficient * value
				}
			}
		}
	}
	return regressors, nil
}

func (reg discoveredRegressor) predict(features []float64) ([]float64, error) {
	if len(features) != discoveredQuadraticFeatureDim {
		return nil, fmt.Errorf("discovered prediction feature mismatch")
	}
	out := []float64{reg.bias[0], reg.bias[1]}
	for output := 0; output < 2; output++ {
		if len(reg.weights[output]) != discoveredQuadraticFeatureDim {
			return nil, fmt.Errorf("discovered regressor weight mismatch")
		}
		for i, value := range features {
			out[output] += reg.weights[output][i] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("discovered prediction non-finite")
		}
	}
	return out, nil
}

func trainDiscoveredModel(rows []discoveredRow) (
	[4]discoveredRegressor,
	[4]linearSoftmaxHead,
	DiscoveredStaticResult,
	error,
) {
	var classifiers [4]linearSoftmaxHead
	regressors, err := trainDiscoveredRegressors(rows, discoveredRidgeLambda)
	if err != nil {
		return regressors, classifiers, DiscoveredStaticResult{}, err
	}
	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return regressors, classifiers, DiscoveredStaticResult{}, err
			}
			samples = append(samples, headSample{
				features: phase,
				target: row.table[entity],
			})
		}
		head, _, err := trainLinearSoftmax(samples, 4, 2, 1200, 1.0)
		if err != nil {
			return regressors, classifiers, DiscoveredStaticResult{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, DiscoveredStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}
	average := 0.0
	for _, accuracy := range perTrain { average += accuracy }
	average /= 4
	return regressors, classifiers, DiscoveredStaticResult{
		Name: "unitary_dynamics_discovered_commutant_phase_code",
		Path: "unitary_discovered_commutant",
		TrainAccuracy: average,
		PerEntityTrain: perTrain,
	}, nil
}

func evaluateDiscoveredModel(
	base DiscoveredStaticResult,
	regressors [4]discoveredRegressor,
	classifiers [4]linearSoftmaxHead,
	rows []discoveredRow,
	maxNormDrift float64,
) (DiscoveredStaticResult, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var cosineTotal float64
	var cosineCount int
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil { return DiscoveredStaticResult{}, err }
			target, err := learnedPhaseTarget(row.table[entity])
			if err != nil { return DiscoveredStaticResult{}, err }
			norm := math.Hypot(phase[0], phase[1])
			if norm > 1e-15 {
				cosine := (phase[0]*target[0]+phase[1]*target[1])/norm
				if cosine > 1 { cosine = 1 }
				if cosine < -1 { cosine = -1 }
				cosineTotal += cosine
				cosineCount++
			}
			samples = append(samples, headSample{
				features: phase,
				target: row.table[entity],
			})
		}
		_, accuracy, err := evaluateHead(classifiers[entity], samples)
		if err != nil { return DiscoveredStaticResult{}, err }
		result.PerEntityHeldOut[entity] = accuracy
	}
	for _, accuracy := range result.PerEntityHeldOut {
		result.HeldOutAccuracy += accuracy
	}
	result.HeldOutAccuracy /= 4
	result.MaxNormDrift = maxNormDrift
	if cosineCount > 0 {
		result.MeanPhaseCosine = cosineTotal / float64(cosineCount)
	}
	return result, nil
}

func decodeDiscoveredTable(
	state State,
	observables []latentMatrix,
	regressors [4]discoveredRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	features, err := discoveredQuadraticFeatures(state, observables)
	if err != nil { return table, distributions, 0, err }
	for entity := 0; entity < 4; entity++ {
		phase, err := regressors[entity].predict(features)
		if err != nil { return table, distributions, 0, err }
		probabilities, err := classifiers[entity].probabilities(phase)
		if err != nil { return table, distributions, 0, err }
		value, margin, err := classAndMargin(probabilities)
		if err != nil { return table, distributions, 0, err }
		table[entity] = value
		distributions[entity] = probabilities
		if margin < minMargin { minMargin = margin }
	}
	return table, distributions, minMargin, nil
}

func runDiscoveredIntegration(
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]discoveredRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (DiscoveredIntegration, error) {
	const scenarios, writes = 48, 16
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil { return DiscoveredIntegration{}, err }
	relationHead, _, err := trainLinearSoftmax(relationSamples, 4, 16, 600, 1.0)
	if err != nil { return DiscoveredIntegration{}, err }

	var commitCorrect, commitTotal, finalCorrect, relationCorrect int
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
			if err != nil { return DiscoveredIntegration{}, err }
			seed := 33000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil { return DiscoveredIntegration{}, err }
			memory = rotateGlobalPhase(memory, math.Mod(0.271*float64(seed+1), 2*math.Pi))
			state, err := fullLatentEncode(memory, mixer)
			if err != nil { return DiscoveredIntegration{}, err }
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return DiscoveredIntegration{}, fmt.Errorf("missing mutable discovered depth=%d", write.gap)
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil { return DiscoveredIntegration{}, err }
			norm2, err := NormSquared(state)
			if err != nil { return DiscoveredIntegration{}, err }
			drift := math.Abs(norm2-1)
			if drift > maxNormDrift { maxNormDrift = drift }

			decoded, _, margin, err := decodeDiscoveredTable(
				state, observables, regressors, classifiers,
			)
			if err != nil { return DiscoveredIntegration{}, err }
			if margin < minValueMargin { minValueMargin = margin }
			commitTotal++
			if decoded == trueTable { commitCorrect++ }
			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil { return DiscoveredIntegration{}, err }
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil { return DiscoveredIntegration{}, err }
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil { return DiscoveredIntegration{}, err }
		seed := 33000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil { return DiscoveredIntegration{}, err }
		memory = rotateGlobalPhase(memory, math.Mod(0.271*float64(seed+1), 2*math.Pi))
		state, err := fullLatentEncode(memory, mixer)
		if err != nil { return DiscoveredIntegration{}, err }
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return DiscoveredIntegration{}, fmt.Errorf("missing final discovered depth=%d", scenario.finalGap)
		}
		state, err = latentMatrixVector(operator, state)
		if err != nil { return DiscoveredIntegration{}, err }
		norm2, err := NormSquared(state)
		if err != nil { return DiscoveredIntegration{}, err }
		drift := math.Abs(norm2-1)
		if drift > maxNormDrift { maxNormDrift = drift }

		decoded, distributions, margin, err := decodeDiscoveredTable(
			state, observables, regressors, classifiers,
		)
		if err != nil { return DiscoveredIntegration{}, err }
		if margin < minValueMargin { minValueMargin = margin }
		if decoded == trueTable { finalCorrect++ }

		relationInput, err := relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err != nil { return DiscoveredIntegration{}, err }
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil { return DiscoveredIntegration{}, err }
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil { return DiscoveredIntegration{}, err }
		if relationMargin < minRelationMargin { minRelationMargin = relationMargin }
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil { return DiscoveredIntegration{}, err }
		if gotRelation == wantRelation { relationCorrect++ }
	}

	return DiscoveredIntegration{
		Scenarios: scenarios,
		WritesPerScenario: writes,
		CommitDecodeAccuracy: float64(commitCorrect)/float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect)/float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect)/float64(scenarios),
		MinValueMargin: minValueMargin,
		MinRelationMargin: minRelationMargin,
		MaxNormDrift: maxNormDrift,
	}, nil
}

func RunUP25() (DiscoveredCommutantProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	unitaryStep, err := conjugatedLatentStep(mixer, applyStressUnitary)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	controlStep, err := conjugatedLatentStep(mixer, applyStressNonUnitary)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	observables, minPreNorm, err := discoverCommutingObservables(
		unitaryStep,
		discoveredObservableCount,
		discoveredProjectionRounds,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	commutatorError, err := maxWeylCommutatorEntry(observables, unitaryStep)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	featureDrift, err := discoveredFeatureDrift(mixer, observables, unitaryOps[1024])
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	trainRows, trainNormDrift, err := buildDiscoveredRows(
		trainTables, trainDepths, unitaryOps,
		mixer, observables, memoryNoise, 2, 0,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	regressors, classifiers, base, err := trainDiscoveredModel(trainRows)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	base.MaxNormDrift = trainNormDrift

	heldRows, heldDrift, err := buildDiscoveredRows(
		heldTables, heldDepths, unitaryOps,
		mixer, observables, memoryNoise, 2, 7000000,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	unitaryResult, err := evaluateDiscoveredModel(
		base, regressors, classifiers, heldRows, heldDrift,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	controlRows, controlDrift, err := buildDiscoveredRows(
		heldTables, heldDepths, controlOps,
		mixer, observables, memoryNoise, 2, 7000000,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }
	controlResult, err := evaluateDiscoveredModel(
		DiscoveredStaticResult{
			Name: "non_unitary_discovered_commutant_same_model",
			Path: "non_unitary_same_model",
			TrainAccuracy: base.TrainAccuracy,
			PerEntityTrain: append([]float64(nil), base.PerEntityTrain...),
		},
		regressors, classifiers, controlRows, controlDrift,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	integration, err := runDiscoveredIntegration(
		mixer, observables, unitaryOps,
		regressors, classifiers,
		heldTables, heldDepths, memoryNoise,
	)
	if err != nil { return DiscoveredCommutantProbeResult{}, err }

	commutatorPass := commutatorError <= 1e-5
	invariancePass := featureDrift <= 5e-3
	learningPass := unitaryResult.TrainAccuracy >= 0.99
	unseenPass :=
		unitaryResult.HeldOutAccuracy >= 0.99 &&
		unitaryResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
		integration.ExactFinalTableAccuracy >= 0.95 &&
		integration.RelationalQueryAccuracy >= 0.95 &&
		integration.MaxNormDrift <= 1e-10

	return DiscoveredCommutantProbeResult{
		Schema: DiscoveredCommutantSchema,
		Experiment: "UP-25-dynamics-discovered-commuting-observables",
		LatentDimension: fullLatentDimension,
		RuntimeStateObjects: 1,
		VisibleChannelBlocks: false,
		FullCoordinateMixing: true,
		ObservableDiscoveryFromTransport: true,
		DiscoveryUsesHiddenMultiplicity: false,
		DiscoveryUsesHiddenMixer: false,
		TransportAdjointUsedInDiscovery: true,
		RuntimeAdjointApplied: false,
		DiscoverySeedCount: discoveredObservableCount,
		DiscoveryProjectionRounds: discoveredProjectionRounds,
		DiscoveryEffectiveOrbitSize: discoveredEffectiveOrbitSize,
		KnownFullMixerExposedToLearner: false,
		RuntimeUnmixApplied: false,
		LearnedDemixerUsed: false,
		PhaseAlphabetSupervision: true,
		PhaseCodeDimension: 2,
		RawFeatureDimension: discoveredRawFeatureDim,
		QuadraticFeatureDimension: discoveredQuadraticFeatureDim,
		RuntimePrototypeLookup: false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided: false,
		GlobalPhaseNuisance: true,
		MemoryNoiseAmplitude: memoryNoise,
		TrainingDepths: trainDepths,
		HeldOutDepths: heldDepths,
		TrainTables: len(trainTables),
		HeldOutTables: len(heldTables),
		Unitary: unitaryResult,
		MatchedControl: controlResult,
		Integration: integration,
		Diagnosis: DiscoveredDiagnosis{
			DiscoveryCommutatorPass: commutatorPass,
			FeatureInvariancePass: invariancePass,
			PhaseCodeLearningPass: learningPass,
			UnseenDepthPass: unseenPass,
			MutableIntegrationPass: mutablePass,
			MaxCommutatorEntryError: commutatorError,
			MaxFeatureDrift: featureDrift,
			MinPreNormalizationFrobenius: minPreNorm,
			TrainAccuracy: unitaryResult.TrainAccuracy,
			HeldOutAccuracy: unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy: controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine: unitaryResult.MeanPhaseCosine,
		},
	}, nil
}
