package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const FullLatentSchema = "wingless.unitary-full-latent-mix.v1"

const (
	fullLatentDimension    = compositeDimension
	fullLatentFeatureDim   = fullLatentDimension * fullLatentDimension
	fullLatentRidgeLambda  = 1e-5
)

type latentMatrix [][]complex128

type fullLatentRegressor struct {
	weights [2][]float64
	bias    [2]float64
}

type FullLatentStaticResult struct {
	Name             string    `json:"name"`
	Path             string    `json:"path"`
	TrainAccuracy    float64   `json:"train_accuracy"`
	HeldOutAccuracy  float64   `json:"held_out_accuracy"`
	PerEntityTrain   []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut []float64 `json:"per_entity_held_out_accuracy"`
	MaxNormDrift     float64   `json:"max_norm_drift"`
}

type FullLatentIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type FullLatentDiagnosis struct {
	ConjugationEquivalencePass bool    `json:"conjugation_equivalence_pass"`
	FullLatentLearningPass     bool    `json:"full_latent_learning_pass"`
	UnseenDepthPass            bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass     bool    `json:"mutable_integration_pass"`
	MaxOracleRecoveryError     float64 `json:"max_oracle_recovery_error"`
	TrainAccuracy              float64 `json:"train_accuracy"`
	HeldOutAccuracy            float64 `json:"held_out_accuracy"`
	MatchedControlAccuracy     float64 `json:"matched_control_accuracy"`
	MeanHeldPhaseCosine        float64 `json:"mean_held_phase_cosine"`
	MinimumMixerParticipation  float64 `json:"minimum_mixer_participation_ratio"`
}

type FullLatentProbeResult struct {
	Schema                          string                `json:"schema"`
	Experiment                      string                `json:"experiment"`
	LatentDimension                 int                   `json:"latent_dimension"`
	RuntimeStateObjects             int                   `json:"runtime_state_objects"`
	VisibleChannelBlocks            bool                  `json:"visible_channel_blocks"`
	FullCoordinateMixing            bool                  `json:"full_coordinate_mixing"`
	ConjugatedTransport             bool                  `json:"conjugated_transport"`
	KnownFullMixerExposedToLearner  bool                  `json:"known_full_mixer_exposed_to_learner"`
	OracleUnmixUsedForLearning      bool                  `json:"oracle_unmix_used_for_learning"`
	RuntimeUnmixApplied             bool                  `json:"runtime_unmix_applied"`
	LearnedDemixerUsed              bool                  `json:"learned_demixer_used"`
	PhaseAlphabetSupervision        bool                  `json:"phase_alphabet_supervision"`
	PhaseCodeDimension              int                   `json:"phase_code_dimension"`
	FullHermitianFeatureDimension   int                   `json:"full_hermitian_feature_dimension"`
	RuntimePrototypeLookup          bool                  `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                  `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided           bool                  `json:"explicit_depth_provided"`
	GlobalPhaseNuisance             bool                  `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude            float64               `json:"memory_noise_amplitude"`
	RidgeLambda                     float64               `json:"ridge_lambda"`
	TrainingDepths                  []int                 `json:"training_depths"`
	HeldOutDepths                   []int                 `json:"held_out_depths"`
	TrainTables                     int                   `json:"train_tables"`
	HeldOutTables                   int                   `json:"held_out_tables"`
	Unitary                         FullLatentStaticResult `json:"unitary"`
	MatchedControl                  FullLatentStaticResult `json:"matched_control"`
	Integration                     FullLatentIntegration  `json:"unitary_mutable_integration"`
	Diagnosis                       FullLatentDiagnosis    `json:"diagnosis"`
}

type fullLatentRow struct {
	features []float64
	table    memoryTable
}

func identityLatentMatrix(dimension int) latentMatrix {
	out := make(latentMatrix, dimension)
	for row := 0; row < dimension; row++ {
		out[row] = make([]complex128, dimension)
		out[row][row] = 1
	}
	return out
}

func fullLatentMixer() latentMatrix {
	out := make(latentMatrix, fullLatentDimension)
	scale := 1 / math.Sqrt(float64(fullLatentDimension))
	for row := 0; row < fullLatentDimension; row++ {
		out[row] = make([]complex128, fullLatentDimension)
		for column := 0; column < fullLatentDimension; column++ {
			phase := 2 * math.Pi *
				float64(row*column) /
				float64(fullLatentDimension)
			out[row][column] = cmplx.Rect(scale, phase)
		}
	}
	return out
}

func latentAdjoint(matrix latentMatrix) (latentMatrix, error) {
	n := len(matrix)
	if n == 0 {
		return nil, fmt.Errorf("latent adjoint requires matrix")
	}
	out := make(latentMatrix, n)
	for row := 0; row < n; row++ {
		if len(matrix[row]) != n {
			return nil, fmt.Errorf("latent matrix must be square")
		}
		out[row] = make([]complex128, n)
		for column := 0; column < n; column++ {
			out[row][column] = cmplx.Conj(matrix[column][row])
		}
	}
	return out, nil
}

func latentMatrixMultiply(a, b latentMatrix) (latentMatrix, error) {
	n := len(a)
	if n == 0 || len(b) != n {
		return nil, fmt.Errorf("latent matrix multiply dimension mismatch")
	}
	out := make(latentMatrix, n)
	for row := 0; row < n; row++ {
		if len(a[row]) != n || len(b[row]) != n {
			return nil, fmt.Errorf("latent matrices must be square")
		}
		out[row] = make([]complex128, n)
	}
	for row := 0; row < n; row++ {
		for middle := 0; middle < n; middle++ {
			av := a[row][middle]
			if av == 0 {
				continue
			}
			for column := 0; column < n; column++ {
				out[row][column] += av * b[middle][column]
			}
		}
	}
	return out, nil
}

func latentMatrixVector(matrix latentMatrix, state State) (State, error) {
	n := len(matrix)
	if n == 0 || len(state) != n {
		return nil, fmt.Errorf("latent matrix-vector dimension mismatch")
	}
	out := make(State, n)
	for row := 0; row < n; row++ {
		if len(matrix[row]) != n {
			return nil, fmt.Errorf("latent matrix must be square")
		}
		var value complex128
		for column := 0; column < n; column++ {
			value += matrix[row][column] * state[column]
		}
		if !finite(real(value)) || !finite(imag(value)) {
			return nil, fmt.Errorf("latent matrix-vector result is non-finite")
		}
		out[row] = value
	}
	return out, nil
}

func latentMatrixPower(matrix latentMatrix, exponent int) (latentMatrix, error) {
	if exponent < 0 {
		return nil, fmt.Errorf("latent matrix power exponent must be nonnegative")
	}
	n := len(matrix)
	result := identityLatentMatrix(n)
	if exponent == 0 {
		return result, nil
	}
	base := matrix
	for exponent > 0 {
		if exponent&1 == 1 {
			next, err := latentMatrixMultiply(result, base)
			if err != nil {
				return nil, err
			}
			result = next
		}
		exponent >>= 1
		if exponent > 0 {
			next, err := latentMatrixMultiply(base, base)
			if err != nil {
				return nil, err
			}
			base = next
		}
	}
	return result, nil
}

func compositeOneStepMatrix(apply stressApply) (latentMatrix, error) {
	block := stressProgram()
	out := make(latentMatrix, fullLatentDimension)
	for row := range out {
		out[row] = make([]complex128, fullLatentDimension)
	}
	for column := 0; column < fullLatentDimension; column++ {
		basis := make(State, fullLatentDimension)
		basis[column] = 1
		evolved, err := evolveCompositeState(basis, block, 1, apply)
		if err != nil {
			return nil, err
		}
		for row := 0; row < fullLatentDimension; row++ {
			out[row][column] = evolved[row]
		}
	}
	return out, nil
}

func conjugatedLatentStep(
	mixer latentMatrix,
	apply stressApply,
) (latentMatrix, error) {
	base, err := compositeOneStepMatrix(apply)
	if err != nil {
		return nil, err
	}
	adjoint, err := latentAdjoint(mixer)
	if err != nil {
		return nil, err
	}
	left, err := latentMatrixMultiply(mixer, base)
	if err != nil {
		return nil, err
	}
	return latentMatrixMultiply(left, adjoint)
}

func latentDepthOperators(
	step latentMatrix,
	depths []int,
) (map[int]latentMatrix, error) {
	out := make(map[int]latentMatrix, len(depths))
	for _, depth := range depths {
		if _, exists := out[depth]; exists {
			continue
		}
		powered, err := latentMatrixPower(step, depth)
		if err != nil {
			return nil, err
		}
		out[depth] = powered
	}
	return out, nil
}

func fullLatentEncode(memory State, mixer latentMatrix) (State, error) {
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return nil, err
	}
	packed, err := packCompositeState(memory, bank)
	if err != nil {
		return nil, err
	}
	return latentMatrixVector(mixer, packed)
}

func fullLatentHermitianFeatures(state State) ([]float64, error) {
	if len(state) != fullLatentDimension {
		return nil, fmt.Errorf(
			"full latent feature dimension=%d want=%d",
			len(state), fullLatentDimension,
		)
	}
	scale := float64(fullLatentDimension)
	out := make([]float64, 0, fullLatentFeatureDim)

	for coordinate := 0; coordinate < fullLatentDimension; coordinate++ {
		value := cmplx.Abs(state[coordinate])
		out = append(out, scale*value*value)
	}

	for first := 0; first < fullLatentDimension; first++ {
		for second := first + 1; second < fullLatentDimension; second++ {
			value := cmplx.Conj(state[first]) * state[second]
			out = append(
				out,
				scale*real(value),
				scale*imag(value),
			)
		}
	}

	if len(out) != fullLatentFeatureDim {
		return nil, fmt.Errorf(
			"full latent Hermitian feature dimension=%d want=%d",
			len(out), fullLatentFeatureDim,
		)
	}
	for _, value := range out {
		if !finite(value) {
			return nil, fmt.Errorf("full latent feature is non-finite")
		}
	}
	return out, nil
}

func fullLatentMixerParticipation(mixer latentMatrix) (float64, error) {
	if len(mixer) != fullLatentDimension {
		return 0, fmt.Errorf("full latent mixer dimension mismatch")
	}
	minimum := math.Inf(1)
	for row := 0; row < fullLatentDimension; row++ {
		if len(mixer[row]) != fullLatentDimension {
			return 0, fmt.Errorf("full latent mixer row dimension mismatch")
		}
		var fourth float64
		for column := 0; column < fullLatentDimension; column++ {
			magnitude := cmplx.Abs(mixer[row][column])
			fourth += magnitude * magnitude * magnitude * magnitude
		}
		if fourth <= 0 {
			return 0, fmt.Errorf("full latent mixer row has zero mass")
		}
		participation := 1 / fourth
		if participation < minimum {
			minimum = participation
		}
	}
	return minimum, nil
}

func fullLatentRecoveryError(
	mixer latentMatrix,
	unitaryOps map[int]latentMatrix,
	depths []int,
) (float64, error) {
	adjoint, err := latentAdjoint(mixer)
	if err != nil {
		return 0, err
	}
	table := memoryTable{3, 1, 0, 2}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 220022, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.337)
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return 0, err
	}
	packed, err := packCompositeState(memory, bank)
	if err != nil {
		return 0, err
	}
	latent, err := latentMatrixVector(mixer, packed)
	if err != nil {
		return 0, err
	}

	maximum := 0.0
	for _, depth := range depths {
		operator, ok := unitaryOps[depth]
		if !ok {
			return 0, fmt.Errorf("missing latent depth operator %d", depth)
		}
		evolvedLatent, err := latentMatrixVector(operator, latent)
		if err != nil {
			return 0, err
		}
		recovered, err := latentMatrixVector(adjoint, evolvedLatent)
		if err != nil {
			return 0, err
		}
		want, err := evolveCompositeState(
			packed, stressProgram(), depth, applyStressUnitary,
		)
		if err != nil {
			return 0, err
		}
		distance, err := L2Distance(recovered, want)
		if err != nil {
			return 0, err
		}
		if distance > maximum {
			maximum = distance
		}
	}
	return maximum, nil
}

func buildFullLatentRows(
	tables []memoryTable,
	depths []int,
	operators map[int]latentMatrix,
	mixer latentMatrix,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]fullLatentRow, float64, error) {
	rows := make([]fullLatentRow, 0, len(tables)*len(depths)*trials)
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
					"missing full latent operator depth=%d", depth,
				)
			}
			for trial := 0; trial < trials; trial++ {
				seed := seedOffset +
					memoryTableIndex(table)*100000 +
					depthIndex*1000 +
					trial*17
				memory, err := perturbMemory(
					canonical, seed, memoryNoise,
				)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.239*float64(seed+1), 2*math.Pi),
				)
				latent, err := fullLatentEncode(memory, mixer)
				if err != nil {
					return nil, 0, err
				}
				latent, err = latentMatrixVector(operator, latent)
				if err != nil {
					return nil, 0, err
				}
				norm2, err := NormSquared(latent)
				if err != nil {
					return nil, 0, err
				}
				drift := math.Abs(norm2 - 1)
				if drift > maxNormDrift {
					maxNormDrift = drift
				}
				features, err := fullLatentHermitianFeatures(latent)
				if err != nil {
					return nil, 0, err
				}
				rows = append(rows, fullLatentRow{
					features: features,
					table:    table,
				})
			}
		}
	}

	return rows, maxNormDrift, nil
}

func trainFullLatentRegressors(
	rows []fullLatentRow,
	lambda float64,
) ([4]fullLatentRegressor, error) {
	var regressors [4]fullLatentRegressor
	if len(rows) == 0 || lambda <= 0 {
		return regressors, fmt.Errorf("invalid full latent ridge configuration")
	}

	n := len(rows)
	const outputs = 8
	gram := make([][]float64, n)
	right := make([][]float64, n)

	for row := 0; row < n; row++ {
		if len(rows[row].features) != fullLatentFeatureDim {
			return regressors, fmt.Errorf(
				"full latent training feature dimension mismatch",
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
		for feature := 0; feature < fullLatentFeatureDim; feature++ {
			value += rows[first].features[feature] *
				rows[first].features[feature]
		}
		gram[first][first] = value
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < fullLatentFeatureDim; feature++ {
				value += rows[first].features[feature] *
					rows[second].features[feature]
			}
			if !finite(value) {
				return regressors, fmt.Errorf(
					"full latent Gram value is non-finite",
				)
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
				make([]float64, fullLatentFeatureDim)
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

func (regressor fullLatentRegressor) predict(
	features []float64,
) ([]float64, error) {
	if len(features) != fullLatentFeatureDim {
		return nil, fmt.Errorf("full latent prediction feature mismatch")
	}
	out := []float64{regressor.bias[0], regressor.bias[1]}
	for output := 0; output < 2; output++ {
		if len(regressor.weights[output]) != fullLatentFeatureDim {
			return nil, fmt.Errorf("full latent regressor weight mismatch")
		}
		for feature, value := range features {
			out[output] += regressor.weights[output][feature] * value
		}
		if !finite(out[output]) {
			return nil, fmt.Errorf("full latent regressor output non-finite")
		}
	}
	return out, nil
}

func trainFullLatentModel(
	rows []fullLatentRow,
	lambda float64,
) (
	[4]fullLatentRegressor,
	[4]linearSoftmaxHead,
	FullLatentStaticResult,
	error,
) {
	var classifiers [4]linearSoftmaxHead
	regressors, err := trainFullLatentRegressors(rows, lambda)
	if err != nil {
		return regressors, classifiers, FullLatentStaticResult{}, err
	}

	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phaseFeatures, err := regressors[entity].predict(row.features)
			if err != nil {
				return regressors, classifiers, FullLatentStaticResult{}, err
			}
			samples = append(samples, headSample{
				features: phaseFeatures,
				target:   row.table[entity],
			})
		}
		head, _, err := trainLinearSoftmax(
			samples, 4, 2, 1200, 1.0,
		)
		if err != nil {
			return regressors, classifiers, FullLatentStaticResult{}, err
		}
		classifiers[entity] = head
		_, accuracy, err := evaluateHead(head, samples)
		if err != nil {
			return regressors, classifiers, FullLatentStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}

	var average float64
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return regressors, classifiers, FullLatentStaticResult{
		Name:           "unitary_full_latent_phase_code",
		Path:           "unitary_full_96d_latent",
		TrainAccuracy:  average,
		PerEntityTrain: perTrain,
	}, nil
}

func evaluateFullLatentModel(
	base FullLatentStaticResult,
	regressors [4]fullLatentRegressor,
	classifiers [4]linearSoftmaxHead,
	rows []fullLatentRow,
	maxNormDrift float64,
) (FullLatentStaticResult, float64, error) {
	result := base
	result.PerEntityHeldOut = make([]float64, 4)
	var cosineTotal float64
	var cosineCount int

	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			phaseFeatures, err := regressors[entity].predict(row.features)
			if err != nil {
				return FullLatentStaticResult{}, 0, err
			}
			target, err := learnedPhaseTarget(row.table[entity])
			if err != nil {
				return FullLatentStaticResult{}, 0, err
			}
			predNorm := math.Hypot(phaseFeatures[0], phaseFeatures[1])
			if predNorm > 1e-15 {
				cosine := (phaseFeatures[0]*target[0] +
					phaseFeatures[1]*target[1]) / predNorm
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
				features: phaseFeatures,
				target:   row.table[entity],
			})
		}
		_, accuracy, err := evaluateHead(
			classifiers[entity], samples,
		)
		if err != nil {
			return FullLatentStaticResult{}, 0, err
		}
		result.PerEntityHeldOut[entity] = accuracy
	}

	for _, accuracy := range result.PerEntityHeldOut {
		result.HeldOutAccuracy += accuracy
	}
	result.HeldOutAccuracy /= 4
	result.MaxNormDrift = maxNormDrift

	meanCosine := 0.0
	if cosineCount > 0 {
		meanCosine = cosineTotal / float64(cosineCount)
	}
	return result, meanCosine, nil
}

func decodeFullLatentTable(
	state State,
	regressors [4]fullLatentRegressor,
	classifiers [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)

	features, err := fullLatentHermitianFeatures(state)
	if err != nil {
		return memoryTable{}, distributions, 0, err
	}

	for entity := 0; entity < 4; entity++ {
		phaseFeatures, err := regressors[entity].predict(features)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err :=
			classifiers[entity].probabilities(phaseFeatures)
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

func runFullLatentIntegration(
	mixer latentMatrix,
	unitaryOps map[int]latentMatrix,
	regressors [4]fullLatentRegressor,
	classifiers [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
) (FullLatentIntegration, error) {
	const (
		scenarios = 48
		writes    = 16
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return FullLatentIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return FullLatentIntegration{}, err
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
			operator, ok := unitaryOps[write.gap]
			if !ok {
				return FullLatentIntegration{}, fmt.Errorf(
					"missing mutable latent operator depth=%d", write.gap,
				)
			}
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return FullLatentIntegration{}, err
			}
			seed := 27000000 + scenarioIndex*10000 + writeIndex*31
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return FullLatentIntegration{}, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.241*float64(seed+1), 2*math.Pi),
			)
			latent, err := fullLatentEncode(memory, mixer)
			if err != nil {
				return FullLatentIntegration{}, err
			}
			latent, err = latentMatrixVector(operator, latent)
			if err != nil {
				return FullLatentIntegration{}, err
			}
			norm2, err := NormSquared(latent)
			if err != nil {
				return FullLatentIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			decoded, _, margin, err := decodeFullLatentTable(
				latent, regressors, classifiers,
			)
			if err != nil {
				return FullLatentIntegration{}, err
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
				return FullLatentIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(
				trueTable, write.entity, write.value,
			)
			if err != nil {
				return FullLatentIntegration{}, err
			}
		}

		operator, ok := unitaryOps[scenario.finalGap]
		if !ok {
			return FullLatentIntegration{}, fmt.Errorf(
				"missing final latent operator depth=%d", scenario.finalGap,
			)
		}
		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		seed := 27000000 + scenarioIndex*10000 + 9999
		memory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.241*float64(seed+1), 2*math.Pi),
		)
		latent, err := fullLatentEncode(memory, mixer)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		latent, err = latentMatrixVector(operator, latent)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		norm2, err := NormSquared(latent)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		decoded, distributions, margin, err := decodeFullLatentTable(
			latent, regressors, classifiers,
		)
		if err != nil {
			return FullLatentIntegration{}, err
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
			return FullLatentIntegration{}, err
		}
		relationProbabilities, err :=
			relationHead.probabilities(relationInput)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		gotRelation, relationMargin, err :=
			classAndMargin(relationProbabilities)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(
			trueTable, scenario.queryA, scenario.queryB,
		)
		if err != nil {
			return FullLatentIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return FullLatentIntegration{
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

func RunUP22() (FullLatentProbeResult, error) {
	const memoryNoise = 0.05
	trainDepths := []int{0, 216}
	heldDepths := []int{32, 128, 512, 1024}
	allDepths := []int{0, 32, 128, 216, 512, 1024}

	mixer := fullLatentMixer()
	participation, err := fullLatentMixerParticipation(mixer)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	unitaryStep, err := conjugatedLatentStep(
		mixer, applyStressUnitary,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}
	controlStep, err := conjugatedLatentStep(
		mixer, applyStressNonUnitary,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	unitaryOps, err := latentDepthOperators(unitaryStep, allDepths)
	if err != nil {
		return FullLatentProbeResult{}, err
	}
	controlOps, err := latentDepthOperators(controlStep, heldDepths)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	recoveryError, err := fullLatentRecoveryError(
		mixer,
		unitaryOps,
		[]int{0, 32, 216, 1024},
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)

	trainRows, trainDrift, err := buildFullLatentRows(
		trainTables,
		trainDepths,
		unitaryOps,
		mixer,
		memoryNoise,
		1,
		0,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	regressors, classifiers, base, err :=
		trainFullLatentModel(trainRows, fullLatentRidgeLambda)
	if err != nil {
		return FullLatentProbeResult{}, err
	}
	base.MaxNormDrift = trainDrift

	heldRows, heldDrift, err := buildFullLatentRows(
		heldTables,
		heldDepths,
		unitaryOps,
		mixer,
		memoryNoise,
		2,
		7000000,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}
	unitaryResult, meanCosine, err := evaluateFullLatentModel(
		base,
		regressors,
		classifiers,
		heldRows,
		heldDrift,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	controlRows, controlDrift, err := buildFullLatentRows(
		heldTables,
		heldDepths,
		controlOps,
		mixer,
		memoryNoise,
		2,
		7000000,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}
	controlResult, _, err := evaluateFullLatentModel(
		FullLatentStaticResult{
			Name:           "non_unitary_full_latent_phase_code",
			Path:           "non_unitary_full_96d_latent_same_model",
			TrainAccuracy:  base.TrainAccuracy,
			PerEntityTrain: append([]float64(nil), base.PerEntityTrain...),
		},
		regressors,
		classifiers,
		controlRows,
		controlDrift,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	integration, err := runFullLatentIntegration(
		mixer,
		unitaryOps,
		regressors,
		classifiers,
		heldTables,
		heldDepths,
		memoryNoise,
	)
	if err != nil {
		return FullLatentProbeResult{}, err
	}

	equivalencePass := recoveryError <= 1e-10
	learningPass := unitaryResult.TrainAccuracy >= 0.99
	unseenPass :=
		unitaryResult.HeldOutAccuracy >= 0.99 &&
			unitaryResult.MaxNormDrift <= 1e-10
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-10

	return FullLatentProbeResult{
		Schema:                          FullLatentSchema,
		Experiment:                      "UP-22-full-96d-latent-mixing",
		LatentDimension:                 fullLatentDimension,
		RuntimeStateObjects:             1,
		VisibleChannelBlocks:            false,
		FullCoordinateMixing:            true,
		ConjugatedTransport:             true,
		KnownFullMixerExposedToLearner:  false,
		OracleUnmixUsedForLearning:      false,
		RuntimeUnmixApplied:             false,
		LearnedDemixerUsed:              false,
		PhaseAlphabetSupervision:        true,
		PhaseCodeDimension:              2,
		FullHermitianFeatureDimension:   fullLatentFeatureDim,
		RuntimePrototypeLookup:          false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided:           false,
		GlobalPhaseNuisance:             true,
		MemoryNoiseAmplitude:            memoryNoise,
		RidgeLambda:                     fullLatentRidgeLambda,
		TrainingDepths:                  append([]int(nil), trainDepths...),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		TrainTables:                     len(trainTables),
		HeldOutTables:                   len(heldTables),
		Unitary:                         unitaryResult,
		MatchedControl:                  controlResult,
		Integration:                     integration,
		Diagnosis: FullLatentDiagnosis{
			ConjugationEquivalencePass: equivalencePass,
			FullLatentLearningPass:     learningPass,
			UnseenDepthPass:            unseenPass,
			MutableIntegrationPass:     mutablePass,
			MaxOracleRecoveryError:     recoveryError,
			TrainAccuracy:              unitaryResult.TrainAccuracy,
			HeldOutAccuracy:            unitaryResult.HeldOutAccuracy,
			MatchedControlAccuracy:     controlResult.HeldOutAccuracy,
			MeanHeldPhaseCosine:        meanCosine,
			MinimumMixerParticipation:  participation,
		},
	}, nil
}
