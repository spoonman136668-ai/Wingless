package unitary

import (
	"fmt"
	"math"
)

const DiscoveryBreadthSchema = "wingless.unitary-discovery-breadth.v1"

const (
	breadthBaselineCount = 32
	breadthPrimaryCount  = 64
	breadthRounds        = discoveredProjectionRounds
	breadthRidgeLambda   = discoveredRidgeLambda
)

type breadthRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type DiscoveryBreadthArm struct {
	Name                     string    `json:"name"`
	Path                     string    `json:"path"`
	ObservableCount          int       `json:"observable_count"`
	RawFeatureDimension      int       `json:"raw_feature_dimension"`
	QuadraticFeatureDimension int      `json:"quadratic_feature_dimension"`
	TrainAccuracy            float64   `json:"train_accuracy"`
	HeldOutAccuracy          float64   `json:"held_out_accuracy"`
	PerEntityTrain           []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut         []float64 `json:"per_entity_held_out_accuracy"`
	MeanPhaseCosine          float64   `json:"mean_phase_cosine"`
	MaxNormDrift             float64   `json:"max_norm_drift"`
}

type DiscoveryBreadthDiagnosis struct {
	BaselineReproductionPass bool    `json:"baseline_reproduction_pass"`
	DiscoveryCommutatorPass  bool    `json:"discovery_commutator_pass"`
	FeatureInvariancePass    bool    `json:"feature_invariance_pass"`
	PhaseCodeLearningPass    bool    `json:"phase_code_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	BaselineHeldOutDelta     float64 `json:"baseline_heldout_delta"`
	MaxCommutatorEntryError  float64 `json:"max_commutator_entry_error"`
	MaxFeatureDrift          float64 `json:"max_feature_drift"`
	MinPreNormalizationFrobenius float64 `json:"min_pre_normalization_frobenius"`
	PrimaryTrainAccuracy     float64 `json:"primary_train_accuracy"`
	PrimaryHeldOutAccuracy   float64 `json:"primary_heldout_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine      float64 `json:"mean_held_phase_cosine"`
}

type DiscoveryBreadthProbeResult struct {
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
	ProjectionRounds                 int                       `json:"projection_rounds"`
	EffectiveOrbitSize               int                       `json:"effective_orbit_size"`
	RidgeLambda                      float64                   `json:"ridge_lambda"`
	TrainingDepths                   []int                     `json:"training_depths"`
	HeldOutDepths                    []int                     `json:"held_out_depths"`
	TrainTables                      int                       `json:"train_tables"`
	HeldOutTables                    int                       `json:"held_out_tables"`
	Baseline32                       DiscoveryBreadthArm        `json:"baseline_32"`
	Primary64                        DiscoveryBreadthArm        `json:"primary_64"`
	MatchedControl64                 DiscoveryBreadthArm        `json:"matched_control_64"`
	Integration64                    DiscoveredIntegration     `json:"unitary_mutable_integration_64"`
	Diagnosis                        DiscoveryBreadthDiagnosis  `json:"diagnosis"`
}

type breadthRow struct {
	features []float64
	table    memoryTable
}

func breadthRawFeatureDim(count int) int {
	return count * 2
}

func breadthQuadraticFeatureDim(count int) int {
	raw := breadthRawFeatureDim(count)
	return raw + (raw*(raw+1))/2
}

func breadthRawFeatures(
	state State,
	observables []latentMatrix,
) ([]float64, error) {
	if len(observables) < 1 {
		return nil, fmt.Errorf("breadth requires observables")
	}
	out := make([]float64, 0, breadthRawFeatureDim(len(observables)))
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
			return nil, fmt.Errorf("breadth raw feature is non-finite")
		}
	}
	return out, nil
}

func breadthQuadraticFeatures(
	state State,
	observables []latentMatrix,
) ([]float64, error) {
	raw, err := breadthRawFeatures(state, observables)
	if err != nil {
		return nil, err
	}
	expected := breadthQuadraticFeatureDim(len(observables))
	out := make([]float64, 0, expected)
	out = append(out, raw...)
	for first := 0; first < len(raw); first++ {
		for second := first; second < len(raw); second++ {
			value := raw[first] * raw[second]
			if !finite(value) {
				return nil, fmt.Errorf("breadth quadratic feature is non-finite")
			}
			out = append(out, value)
		}
	}
	if len(out) != expected {
		return nil, fmt.Errorf(
			"breadth feature dimension=%d want=%d",
			len(out), expected,
		)
	}
	return out, nil
}

func breadthFeatureDrift(
	mixer latentMatrix,
	observables []latentMatrix,
	depthOperator latentMatrix,
) (float64, error) {
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 262626, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.819)
	state, err := fullLatentEncode(memory, mixer)
	if err != nil {
		return 0, err
	}
	before, err := breadthQuadraticFeatures(state, observables)
	if err != nil {
		return 0, err
	}
	evolved, err := latentMatrixVector(depthOperator, state)
	if err != nil {
		return 0, err
	}
	after, err := breadthQuadraticFeatures(evolved, observables)
	if err != nil {
		return 0, err
	}
	maximum := 0.0
	for i := range before {
		delta := math.Abs(before[i] - after[i])
		if delta > maximum {
			maximum = delta
		}
	}
	return maximum, nil
}

func buildBreadthRows(
	tables []memoryTable,
	depths []int,
	operators map[int]latentMatrix,
	mixer latentMatrix,
	observables []latentMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]breadthRow, float64, error) {
	rows := make([]breadthRow, 0, len(tables)*len(depths)*trials)
	var maxNormDrift float64
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, 0, err
		}
		for depthIndex, depth := range depths {
			operator, ok := operators[depth]
			if !ok {
				return nil, 0, fmt.Errorf(
					"missing breadth operator depth=%d", depth,
				)
			}
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.271*float64(seed+1), 2*math.Pi),
				)
				state, err := fullLatentEncode(memory, mixer)
				if err != nil {
					return nil, 0, err
				}
				state, err = latentMatrixVector(operator, state)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(state)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}
				features, err := breadthQuadraticFeatures(
					state, observables,
				)
				if err != nil {
					return nil, 0, err
				}
				rows = append(rows, breadthRow{
					features: features,
					table:    table,
				})
			}
		}
	}
	return rows, maxNormDrift, nil
}

func trainBreadthRegressors(
	rows []breadthRow,
	featureDim int,
	lambda float64,
) ([4]breadthRegressor, error) {
	var regressors [4]breadthRegressor
	if len(rows) == 0 || featureDim < 1 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid breadth ridge configuration")
	}
	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)

	for row := 0; row < n; row++ {
		if len(rows[row].features) != featureDim {
			return regressors, fmt.Errorf(
				"breadth training feature dimension=%d want=%d",
				len(rows[row].features), featureDim,
			)
		}
		gram[row] = make([]float64, n)
		right[row] = make([]float64, outputs)
		for entity := 0; entity < 4; entity++ {
			target, err := learnedPhaseTarget(rows[row].table[entity])
			if err != nil {
				return regressors, err
			}
			right[row][entity*2] = target[0]
			right[row][entity*2+1] = target[1]
		}
	}

	for first := 0; first < n; first++ {
		value := 1.0 + lambda
		for feature := 0; feature < featureDim; feature++ {
			value += rows[first].features[feature] *
				rows[first].features[feature]
		}
		gram[first][first] = value
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < featureDim; feature++ {
				value += rows[first].features[feature] *
					rows[second].features[feature]
			}
			if !finite(value) {
				return regressors, fmt.Errorf("breadth gram value non-finite")
			}
			gram[first][second] = value
			gram[second][first] = value
		}
	}

	alpha, err := solveDenseMultiple(gram, right)
	if err != nil {
		return regressors, err
	}

	for entity := 0; entity < 4; entity++ {
		for output := 0; output < 2; output++ {
			regressors[entity].weights[output] =
				make([]float64, featureDim)
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

func (reg breadthRegressor) predict(features []float64) ([]float64, error) {
	featureDim := len(reg.weights[0])
	if featureDim == 0 || len(features) != featureDim ||
		len(reg.weights[1]) != featureDim {
		return nil, fmt.Errorf("breadth prediction feature mismatch")
	}
	out := []float64{reg.bias[0], reg.bias[1]}
	for output := 0; output < 2; output++ {
		for i, value := range features {
			out[output] += reg.weights[output][i] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("breadth prediction non-finite")
		}
	}
	return out, nil
}

func trainBreadthModel(
	rows []breadthRow,
	observableCount int,
	name, path string,
) (
	[4]breadthRegressor,
	[4]linearSoftmaxHead,
	DiscoveryBreadthArm,
	error,
) {
	var classifiers [4]linearSoftmaxHead
	featureDim := breadthQuadraticFeatureDim(observableCount)
	regressors, err := trainBreadthRegressors(
		rows, featureDim, breadthRidgeLambda,
	)
	if err != nil {
		return regressors, classifiers, DiscoveryBreadthArm{}, err
	}

	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return regressors, classifiers, DiscoveryBreadthArm{}, err
			}
			samples = append(samples, headSample{
				features: phase,
				target:   row.table[entity],
			})
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, 2, 1200, 1.0,
		)
		if err != nil {
			return regressors, classifiers, DiscoveryBreadthArm{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, DiscoveryBreadthArm{}, err
		}
		perTrain[entity] = accuracy
	}

	average := 0.0
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return regressors, classifiers, DiscoveryBreadthArm{
		ObservableCount:           observableCount,
		RawFeatureDimension:       breadthRawFeatureDim(observableCount),
		QuadraticFeatureDimension: featureDim,
		TrainAccuracy:             average,
		PerEntityTrain:            perTrain,
		Name:                       name,
		Path:                       path,
	}, nil
}

func evaluateBreadthModel(
	base DiscoveryBreadthArm,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	rows []breadthRow,
	maxNormDrift float64,
) (DiscoveryBreadthArm, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var cosineTotal float64
	var cosineCount int

	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phase, err := regressors[entity].predict(row.features)
			if err != nil {
				return DiscoveryBreadthArm{}, err
			}
			target, err := learnedPhaseTarget(row.table[entity])
			if err != nil {
				return DiscoveryBreadthArm{}, err
			}
			norm := math.Hypot(phase[0], phase[1])
			if norm > 1e-15 {
				cosine := (phase[0]*target[0] +
					phase[1]*target[1]) / norm
				if cosine > 1 {
					cosine = 1
				}
				if cosine < -1 {
					cosine = -1
				}
				cosineTotal += cosine
				cosineCount++
			}
			samples = append(samples, headSample{
				features: phase,
				target:   row.table[entity],
			})
		}
		_, accuracy, err := evaluateHead(
			classifiers[entity], samples,
		)
		if err != nil {
			return DiscoveryBreadthArm{}, err
		}
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

func decodeBreadthTable(
	state State,
	observables []latentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)

	features, err := breadthQuadraticFeatures(state, observables)
	if err != nil {
		return table, distributions, 0, err
	}

	for entity := 0; entity < 4; entity++ {
		phase, err := regressors[entity].predict(features)
		if err != nil {
			return table, distributions, 0, err
		}
		probabilities, err := classifiers[entity].probabilities(phase)
		if err != nil {
			return table, distributions, 0, err
		}
		value, margin, err := classAndMargin(probabilities)
		if err != nil {
			return table, distributions, 0, err
		}
		table[entity] = value
		distributions[entity] = probabilities
		if margin < minMargin {
			minMargin = margin
		}
	}
	return table, distributions, minMargin, nil
}

func runBreadthIntegration(
	mixer latentMatrix,
	observables []latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]breadthRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (DiscoveredIntegration, error) {
	const scenarios, writes = 48, 16

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return DiscoveredIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return DiscoveredIntegration{}, err
	}

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
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			seed := 34000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1), 2*math.Pi),
			)
			state, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return DiscoveredIntegration{}, fmt.Errorf(
					"missing breadth mutable depth=%d", write.gap,
				)
			}
			state, err = latentMatrixVector(operator, state)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			norm2, err := NormSquared(state)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeBreadthTable(
				state, observables, regressors, classifiers,
			)
			if err != nil {
				return DiscoveredIntegration{}, err
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
				return DiscoveredIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return DiscoveredIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		seed := 34000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1), 2*math.Pi),
		)
		state, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return DiscoveredIntegration{}, fmt.Errorf(
				"missing breadth final depth=%d", scenario.finalGap,
			)
		}
		state, err = latentMatrixVector(operator, state)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		norm2, err := NormSquared(state)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		decoded, distributions, margin, err := decodeBreadthTable(
			state, observables, regressors, classifiers,
		)
		if err != nil {
			return DiscoveredIntegration{}, err
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
			return DiscoveredIntegration{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return DiscoveredIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return DiscoveredIntegration{
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

func RunUP26() (DiscoveryBreadthProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 512, 1024}

	mixer := fullLatentMixer()
	unitaryStep, err := conjugatedLatentStep(
		mixer, applyStressUnitary,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	controlStep, err := conjugatedLatentStep(
		mixer, applyStressNonUnitary,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	allObservables, minPreNorm, err := discoverCommutingObservables(
		unitaryStep,
		breadthPrimaryCount,
		breadthRounds,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	baselineObservables := allObservables[:breadthBaselineCount]
	primaryObservables := allObservables

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	commutatorError, err := maxWeylCommutatorEntry(
		primaryObservables, unitaryStep,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	featureDrift, err := breadthFeatureDrift(
		mixer, primaryObservables, unitaryOps[1024],
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	base32Rows, base32TrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, baselineObservables, memoryNoise, 2, 0,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	base32Regs, base32Heads, base32Base, err := trainBreadthModel(
		base32Rows,
		breadthBaselineCount,
		"unitary_discovered_breadth_32",
		"unitary_discovered_32_reproduction",
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	base32Base.MaxNormDrift = base32TrainDrift

	base32HeldRows, base32HeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, baselineObservables, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	base32Result, err := evaluateBreadthModel(
		base32Base,
		base32Regs,
		base32Heads,
		base32HeldRows,
		base32HeldDrift,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	primaryRows, primaryTrainDrift, err := buildBreadthRows(
		trainTables, trainDepths, unitaryOps,
		mixer, primaryObservables, memoryNoise, 2, 0,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	primaryRegs, primaryHeads, primaryBase, err := trainBreadthModel(
		primaryRows,
		breadthPrimaryCount,
		"unitary_discovered_breadth_64",
		"unitary_discovered_64_primary",
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	primaryBase.MaxNormDrift = primaryTrainDrift

	primaryHeldRows, primaryHeldDrift, err := buildBreadthRows(
		heldTables, heldDepths, unitaryOps,
		mixer, primaryObservables, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	primaryResult, err := evaluateBreadthModel(
		primaryBase,
		primaryRegs,
		primaryHeads,
		primaryHeldRows,
		primaryHeldDrift,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	controlRows, controlDrift, err := buildBreadthRows(
		heldTables, heldDepths, controlOps,
		mixer, primaryObservables, memoryNoise, 2, 7000000,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}
	controlResult, err := evaluateBreadthModel(
		DiscoveryBreadthArm{
			ObservableCount:           breadthPrimaryCount,
			RawFeatureDimension:       breadthRawFeatureDim(breadthPrimaryCount),
			QuadraticFeatureDimension: breadthQuadraticFeatureDim(breadthPrimaryCount),
			TrainAccuracy:             primaryBase.TrainAccuracy,
			PerEntityTrain:            append([]float64(nil), primaryBase.PerEntityTrain...),
		},
		primaryRegs,
		primaryHeads,
		controlRows,
		controlDrift,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	integration, err := runBreadthIntegration(
		mixer,
		primaryObservables,
		unitaryOps,
		primaryRegs,
		primaryHeads,
		heldTables,
		heldDepths,
		memoryNoise,
	)
	if err != nil {
		return DiscoveryBreadthProbeResult{}, err
	}

	baselineDelta := math.Abs(
		base32Result.HeldOutAccuracy - 0.9248046875,
	)
	baselineReproductionPass := baselineDelta <= 1e-12
	commutatorPass := commutatorError <= 1e-5
	invariancePass := featureDrift <= 5e-3
	learningPass := primaryResult.TrainAccuracy >= 0.99
	unseenPass :=
		primaryResult.HeldOutAccuracy >= 0.99 &&
			primaryResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-10

	return DiscoveryBreadthProbeResult{
		Schema:                           DiscoveryBreadthSchema,
		Experiment:                       "UP-26-dynamics-discovered-observable-breadth",
		LatentDimension:                  fullLatentDimension,
		RuntimeStateObjects:              1,
		VisibleChannelBlocks:             false,
		FullCoordinateMixing:             true,
		ObservableDiscoveryFromTransport: true,
		DiscoveryUsesHiddenMultiplicity:  false,
		DiscoveryUsesHiddenMixer:         false,
		TransportAdjointUsedInDiscovery:  true,
		RuntimeAdjointApplied:            false,
		ProjectionRounds:                 breadthRounds,
		EffectiveOrbitSize:               1 << breadthRounds,
		RidgeLambda:                      breadthRidgeLambda,
		TrainingDepths:                   append([]int(nil), trainDepths...),
		HeldOutDepths:                    append([]int(nil), heldDepths...),
		TrainTables:                      len(trainTables),
		HeldOutTables:                    len(heldTables),
		Baseline32:                       base32Result,
		Primary64:                        primaryResult,
		MatchedControl64:                 controlResult,
		Integration64:                    integration,
		Diagnosis: DiscoveryBreadthDiagnosis{
			BaselineReproductionPass: baselineReproductionPass,
			DiscoveryCommutatorPass:  commutatorPass,
			FeatureInvariancePass:    invariancePass,
			PhaseCodeLearningPass:    learningPass,
			UnseenDepthPass:          unseenPass,
			MutableIntegrationPass:   mutablePass,
			BaselineHeldOutDelta:     baselineDelta,
			MaxCommutatorEntryError:  commutatorError,
			MaxFeatureDrift:          featureDrift,
			MinPreNormalizationFrobenius: minPreNorm,
			PrimaryTrainAccuracy:     primaryResult.TrainAccuracy,
			PrimaryHeldOutAccuracy:   primaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:   controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine:      primaryResult.MeanPhaseCosine,
		},
	}, nil
}
