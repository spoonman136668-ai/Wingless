package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const RidgeAnonymousSchema = "wingless.unitary-ridge-anonymous.v1"

const ridgeAnonymousLambda = 1e-6

type RidgeAnonymousDiagnosis struct {
	OracleWitnessPass          bool    `json:"oracle_witness_pass"`
	ClosedFormLearningPass     bool    `json:"closed_form_learning_pass"`
	UnseenDepthPass            bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass     bool    `json:"mutable_integration_pass"`
	OracleWitnessMaxError      float64 `json:"oracle_witness_max_error"`
	OracleWitnessHeldAccuracy  float64 `json:"oracle_witness_held_accuracy"`
	RidgeTrainAccuracy         float64 `json:"ridge_train_accuracy"`
	RidgeHeldAccuracy          float64 `json:"ridge_held_accuracy"`
	MatchedControlAccuracy     float64 `json:"matched_control_accuracy"`
	MaxUnitaryFeatureDrift     float64 `json:"max_unitary_feature_drift"`
}

type RidgeAnonymousProbeResult struct {
	Schema                          string                    `json:"schema"`
	Experiment                      string                    `json:"experiment"`
	CompositeDimension              int                       `json:"composite_dimension"`
	RuntimeStateObjects             int                       `json:"runtime_state_objects"`
	AnonymousChannels               int                       `json:"anonymous_channels"`
	SemanticChannelLocationsKnown   bool                      `json:"semantic_channel_locations_known"`
	KnownMixerExposedToLearner      bool                      `json:"known_mixer_exposed_to_learner"`
	OracleWitnessUsedForLearning    bool                      `json:"oracle_witness_used_for_learning"`
	RuntimeMixerInverseApplied      bool                      `json:"runtime_mixer_inverse_applied"`
	LearnedDemixerUsed              bool                      `json:"learned_demixer_used"`
	DirectAnonymousReadout          bool                      `json:"direct_anonymous_readout"`
	RuntimePrototypeLookup          bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                      `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided           bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance             bool                      `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude            float64                   `json:"memory_noise_amplitude"`
	QuadraticFeatureDimension       int                       `json:"quadratic_feature_dimension"`
	RidgeLambda                     float64                   `json:"ridge_lambda"`
	RidgeTrainingTrials             int                       `json:"ridge_training_trials_per_table"`
	TrainTables                     int                       `json:"train_tables"`
	HeldOutTables                   int                       `json:"held_out_tables"`
	HeldOutDepths                   []int                     `json:"held_out_depths"`
	OracleWitnessUnitary            AnonymousGramStaticResult `json:"oracle_witness_unitary"`
	RidgeUnitary                    AnonymousGramStaticResult `json:"ridge_unitary"`
	RidgeMatchedControl             AnonymousGramStaticResult `json:"ridge_matched_control"`
	Integration                     AnonymousGramIntegration  `json:"unitary_mutable_integration"`
	Diagnosis                       RidgeAnonymousDiagnosis   `json:"diagnosis"`
}

type ridgeAnonymousRow struct {
	features []float64
	table    memoryTable
}

func buildRidgeAnonymousRows(
	tables []memoryTable,
	memoryNoise float64,
	trials int,
	seedOffset int,
) ([]ridgeAnonymousRow, error) {
	rows := make([]ridgeAnonymousRow, 0, len(tables)*trials)
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for trial := 0; trial < trials; trial++ {
			seed := seedOffset +
				memoryTableIndex(table)*1000 +
				trial*17
			memory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, err
			}
			memory = rotateGlobalPhase(
				memory,
				math.Mod(0.223*float64(seed+1), 2*math.Pi),
			)
			mixed, err := blindMixedComposite(memory)
			if err != nil {
				return nil, err
			}
			features, err := anonymousQuadraticFeatures(mixed)
			if err != nil {
				return nil, err
			}
			rows = append(rows, ridgeAnonymousRow{
				features: features,
				table:    table,
			})
		}
	}
	return rows, nil
}

func solveDenseMultiple(
	matrix [][]float64,
	right [][]float64,
) ([][]float64, error) {
	n := len(matrix)
	if n == 0 || len(right) != n {
		return nil, fmt.Errorf("dense solve dimension mismatch")
	}
	rhs := len(right[0])
	if rhs == 0 {
		return nil, fmt.Errorf("dense solve requires right-hand sides")
	}

	aug := make([][]float64, n)
	for row := 0; row < n; row++ {
		if len(matrix[row]) != n || len(right[row]) != rhs {
			return nil, fmt.Errorf("dense solve row dimension mismatch")
		}
		aug[row] = make([]float64, n+rhs)
		copy(aug[row], matrix[row])
		copy(aug[row][n:], right[row])
	}

	for pivot := 0; pivot < n; pivot++ {
		best := pivot
		bestAbs := math.Abs(aug[pivot][pivot])
		for row := pivot + 1; row < n; row++ {
			value := math.Abs(aug[row][pivot])
			if value > bestAbs {
				best = row
				bestAbs = value
			}
		}
		if !finite(bestAbs) || bestAbs <= 1e-14 {
			return nil, fmt.Errorf(
				"dense solve singular pivot=%d magnitude=%g",
				pivot, bestAbs,
			)
		}
		if best != pivot {
			aug[pivot], aug[best] = aug[best], aug[pivot]
		}

		pivotValue := aug[pivot][pivot]
		for column := pivot; column < n+rhs; column++ {
			aug[pivot][column] /= pivotValue
		}

		for row := 0; row < n; row++ {
			if row == pivot {
				continue
			}
			factor := aug[row][pivot]
			if factor == 0 {
				continue
			}
			for column := pivot; column < n+rhs; column++ {
				aug[row][column] -= factor * aug[pivot][column]
			}
		}
	}

	out := make([][]float64, n)
	for row := 0; row < n; row++ {
		out[row] = append([]float64(nil), aug[row][n:]...)
		for _, value := range out[row] {
			if !finite(value) {
				return nil, fmt.Errorf("dense solve produced non-finite result")
			}
		}
	}
	return out, nil
}

func trainRidgeAnonymousHeads(
	rows []ridgeAnonymousRow,
	lambda float64,
) ([4]linearSoftmaxHead, AnonymousGramStaticResult, error) {
	var heads [4]linearSoftmaxHead
	if len(rows) == 0 || lambda <= 0 {
		return heads, AnonymousGramStaticResult{}, fmt.Errorf(
			"invalid ridge training configuration",
		)
	}

	n := len(rows)
	const outputs = 16

	gram := make([][]float64, n)
	right := make([][]float64, n)
	for i := 0; i < n; i++ {
		if len(rows[i].features) != anonymousGramQuadraticDim {
			return heads, AnonymousGramStaticResult{}, fmt.Errorf(
				"ridge feature dimension mismatch",
			)
		}
		gram[i] = make([]float64, n)
		right[i] = make([]float64, outputs)
		for entity := 0; entity < 4; entity++ {
			right[i][entity*4+rows[i].table[entity]] = 1
		}
	}

	for first := 0; first < n; first++ {
		gram[first][first] = 1 + lambda
		for second := first + 1; second < n; second++ {
			value := 1.0
			for feature := 0; feature < anonymousGramQuadraticDim; feature++ {
				value += rows[first].features[feature] *
					rows[second].features[feature]
			}
			if !finite(value) {
				return heads, AnonymousGramStaticResult{}, fmt.Errorf(
					"ridge gram value is non-finite",
				)
			}
			gram[first][second] = value
			gram[second][first] = value
		}
		if n > 0 {
			for feature := 0; feature < anonymousGramQuadraticDim; feature++ {
				gram[first][first] +=
					rows[first].features[feature] *
						rows[first].features[feature]
			}
		}
	}

	alpha, err := solveDenseMultiple(gram, right)
	if err != nil {
		return heads, AnonymousGramStaticResult{}, err
	}

	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(
			4, anonymousGramQuadraticDim,
		)
	}

	for row := 0; row < n; row++ {
		for entity := 0; entity < 4; entity++ {
			for class := 0; class < 4; class++ {
				coefficient := alpha[row][entity*4+class]
				heads[entity].bias[class] += coefficient
				for feature, value := range rows[row].features {
					heads[entity].weights[class][feature] +=
						coefficient * value
				}
			}
		}
	}

	perTrain := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples := make([]headSample, 0, len(rows))
		for _, row := range rows {
			samples = append(samples, headSample{
				features: row.features,
				target:   row.table[entity],
			})
		}
		_, accuracy, err := evaluateHead(heads[entity], samples)
		if err != nil {
			return heads, AnonymousGramStaticResult{}, err
		}
		perTrain[entity] = accuracy
	}

	var average float64
	for _, accuracy := range perTrain {
		average += accuracy
	}
	average /= 4

	return heads, AnonymousGramStaticResult{
		Name:           "unitary_anonymous_quadratic_ridge",
		Path:           "unitary_closed_form_ridge",
		TrainAccuracy:  average,
		PerEntityTrain: perTrain,
	}, nil
}

func semanticGramLinearCoefficients(
	firstSemantic, secondSemantic int,
) ([]complex128, error) {
	if firstSemantic < 0 || firstSemantic >= compositeChannels ||
		secondSemantic < 0 || secondSemantic >= compositeChannels {
		return nil, fmt.Errorf("semantic coefficient channel out of range")
	}

	mixer := denseChannelMixer()
	scale := float64(compositeChannels)
	coefficients := make([]complex128, 0, anonymousGramLinearDim)

	for mixed := 0; mixed < compositeChannels; mixed++ {
		coefficients = append(
			coefficients,
			complex(
				mixer[mixed][firstSemantic]*
					mixer[mixed][secondSemantic]/scale,
				0,
			),
		)
	}

	for left := 0; left < compositeChannels; left++ {
		for right := left + 1; right < compositeChannels; right++ {
			a := mixer[left][firstSemantic] *
				mixer[right][secondSemantic]
			b := mixer[right][firstSemantic] *
				mixer[left][secondSemantic]
			coefficients = append(
				coefficients,
				complex((a+b)/scale, 0),
				complex(0, (a-b)/scale),
			)
		}
	}

	if len(coefficients) != anonymousGramLinearDim {
		return nil, fmt.Errorf(
			"semantic coefficient dimension=%d want=%d",
			len(coefficients), anonymousGramLinearDim,
		)
	}
	return coefficients, nil
}

func oracleNumeratorWeights(entity int) ([]complex128, error) {
	if entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("oracle entity out of range")
	}

	phase, err := semanticGramLinearCoefficients(2+entity, 0)
	if err != nil {
		return nil, err
	}
	anchor, err := semanticGramLinearCoefficients(1, 0)
	if err != nil {
		return nil, err
	}

	weights := make([]complex128, anonymousGramQuadraticDim)
	index := anonymousGramLinearDim
	for first := 0; first < anonymousGramLinearDim; first++ {
		for second := first; second < anonymousGramLinearDim; second++ {
			value := phase[first] * cmplx.Conj(anchor[second])
			if first != second {
				value += phase[second] * cmplx.Conj(anchor[first])
			}
			weights[index] = value
			index++
		}
	}
	if index != anonymousGramQuadraticDim {
		return nil, fmt.Errorf("oracle quadratic weight dimension mismatch")
	}
	return weights, nil
}

func oracleNumeratorFromAnonymous(
	state State,
	entity int,
) (complex128, error) {
	features, err := anonymousQuadraticFeatures(state)
	if err != nil {
		return 0, err
	}
	weights, err := oracleNumeratorWeights(entity)
	if err != nil {
		return 0, err
	}
	var out complex128
	for i, feature := range features {
		out += weights[i] * complex(feature, 0)
	}
	if !finite(real(out)) || !finite(imag(out)) {
		return 0, fmt.Errorf("oracle numerator is non-finite")
	}
	return out, nil
}

func directSemanticNumerator(
	memory State,
	entity int,
) (complex128, error) {
	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return 0, err
	}
	packed, err := packCompositeState(memory, bank)
	if err != nil {
		return 0, err
	}
	memoryChannel, err := compositeChannel(packed, 0)
	if err != nil {
		return 0, err
	}
	anchorChannel, err := compositeChannel(packed, 1)
	if err != nil {
		return 0, err
	}
	phaseChannel, err := compositeChannel(packed, 2+entity)
	if err != nil {
		return 0, err
	}
	phaseCorrelation, err := stateInner(phaseChannel, memoryChannel)
	if err != nil {
		return 0, err
	}
	anchorCorrelation, err := stateInner(anchorChannel, memoryChannel)
	if err != nil {
		return 0, err
	}
	return phaseCorrelation * cmplx.Conj(anchorCorrelation), nil
}

func oracleNumeratorWitnessError() (float64, error) {
	tables := allMemoryTables()
	maximum := 0.0
	for tableIndex, table := range tables {
		if tableIndex%17 != 0 {
			continue
		}
		canonical, err := encodeMemory(table)
		if err != nil {
			return 0, err
		}
		memory, err := perturbMemory(
			canonical, 23000000+tableIndex*31, 0.05,
		)
		if err != nil {
			return 0, err
		}
		memory = rotateGlobalPhase(
			memory,
			math.Mod(0.223*float64(tableIndex+1), 2*math.Pi),
		)
		mixed, err := blindMixedComposite(memory)
		if err != nil {
			return 0, err
		}
		for entity := 0; entity < 4; entity++ {
			got, err := oracleNumeratorFromAnonymous(mixed, entity)
			if err != nil {
				return 0, err
			}
			want, err := directSemanticNumerator(memory, entity)
			if err != nil {
				return 0, err
			}
			delta := cmplx.Abs(got - want)
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum, nil
}

func buildOracleNumeratorSamples(
	tables []memoryTable,
	depths []int,
	entity int,
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
					trial*17
				memory, err := perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, 0, err
				}
				memory = rotateGlobalPhase(
					memory,
					math.Mod(0.223*float64(seed+1), 2*math.Pi),
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
				numerator, err := oracleNumeratorFromAnonymous(
					mixed, entity,
				)
				if err != nil {
					return nil, 0, err
				}
				features := []float64{real(numerator), imag(numerator)}
				samples = append(samples, headSample{
					features: features,
					target:   table[entity],
				})
			}
		}
	}
	return samples, maxNormDrift, nil
}

func oracleWitnessStatic(
	trainTables, heldTables []memoryTable,
	heldDepths []int,
	block []Coupling,
	memoryNoise float64,
) (AnonymousGramStaticResult, error) {
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)
	var maxNormDrift float64

	for entity := 0; entity < 4; entity++ {
		trainSamples, _, err := buildOracleNumeratorSamples(
			trainTables,
			[]int{0},
			entity,
			block,
			applyStressUnitary,
			memoryNoise,
			2,
			0,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2, 1200, 1.0,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		heldSamples, drift, err := buildOracleNumeratorSamples(
			heldTables,
			heldDepths,
			entity,
			block,
			applyStressUnitary,
			memoryNoise,
			2,
			7000000,
		)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return AnonymousGramStaticResult{}, err
		}
		perTrain[entity] = trainAccuracy
		perHeld[entity] = heldAccuracy
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
	}

	var trainAverage, heldAverage float64
	for entity := 0; entity < 4; entity++ {
		trainAverage += perTrain[entity]
		heldAverage += perHeld[entity]
	}
	trainAverage /= 4
	heldAverage /= 4

	return AnonymousGramStaticResult{
		Name:             "unitary_oracle_quadratic_witness",
		Path:             "oracle_diagnostic_only",
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
		MaxNormDrift:     maxNormDrift,
	}, nil
}

func RunUP20() (RidgeAnonymousProbeResult, error) {
	const (
		memoryNoise = 0.05
		trainTrials = 4
		heldTrials  = 2
	)
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	heldDepths := []int{32, 128, 512, 1024}
	block := stressProgram()

	witnessError, err := oracleNumeratorWitnessError()
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}
	witnessStatic, err := oracleWitnessStatic(
		trainTables,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	rows, err := buildRidgeAnonymousRows(
		trainTables, memoryNoise, trainTrials, 0,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}
	ridgeHeads, ridgeBase, err := trainRidgeAnonymousHeads(
		rows, ridgeAnonymousLambda,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	ridgeUnitary, err := evaluateAnonymousHeads(
		ridgeBase,
		ridgeHeads,
		heldTables,
		heldDepths,
		true,
		block,
		applyStressUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	ridgeControl, err := evaluateAnonymousHeads(
		AnonymousGramStaticResult{
			Name:           "non_unitary_anonymous_quadratic_ridge",
			Path:           "non_unitary_matched_same_heads",
			TrainAccuracy:  ridgeBase.TrainAccuracy,
			PerEntityTrain: append([]float64(nil), ridgeBase.PerEntityTrain...),
		},
		ridgeHeads,
		heldTables,
		heldDepths,
		true,
		block,
		applyStressNonUnitary,
		memoryNoise,
		heldTrials,
		7000000,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	featureDrift, err := anonymousFeatureDrift(
		applyStressUnitary, 1024,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	integration, err := runAnonymousGramIntegration(
		ridgeHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return RidgeAnonymousProbeResult{}, err
	}

	witnessPass :=
		witnessError <= 1e-12 &&
			witnessStatic.HeldOutAccuracy >= 0.99
	learningPass := ridgeUnitary.TrainAccuracy >= 0.99
	unseenPass :=
		ridgeUnitary.HeldOutAccuracy >= 0.99 &&
			ridgeUnitary.MaxNormDrift <= 1e-12
	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-12

	return RidgeAnonymousProbeResult{
		Schema:                          RidgeAnonymousSchema,
		Experiment:                      "UP-20-closed-form-anonymous-quadratic-readout",
		CompositeDimension:              compositeDimension,
		RuntimeStateObjects:             1,
		AnonymousChannels:               compositeChannels,
		SemanticChannelLocationsKnown:   false,
		KnownMixerExposedToLearner:      false,
		OracleWitnessUsedForLearning:    false,
		RuntimeMixerInverseApplied:      false,
		LearnedDemixerUsed:              false,
		DirectAnonymousReadout:          true,
		RuntimePrototypeLookup:          false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided:           false,
		GlobalPhaseNuisance:             true,
		MemoryNoiseAmplitude:            memoryNoise,
		QuadraticFeatureDimension:       anonymousGramQuadraticDim,
		RidgeLambda:                     ridgeAnonymousLambda,
		RidgeTrainingTrials:             trainTrials,
		TrainTables:                     len(trainTables),
		HeldOutTables:                   len(heldTables),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		OracleWitnessUnitary:            witnessStatic,
		RidgeUnitary:                    ridgeUnitary,
		RidgeMatchedControl:             ridgeControl,
		Integration:                     integration,
		Diagnosis: RidgeAnonymousDiagnosis{
			OracleWitnessPass:         witnessPass,
			ClosedFormLearningPass:    learningPass,
			UnseenDepthPass:           unseenPass,
			MutableIntegrationPass:    mutablePass,
			OracleWitnessMaxError:     witnessError,
			OracleWitnessHeldAccuracy: witnessStatic.HeldOutAccuracy,
			RidgeTrainAccuracy:        ridgeUnitary.TrainAccuracy,
			RidgeHeldAccuracy:         ridgeUnitary.HeldOutAccuracy,
			MatchedControlAccuracy:    ridgeControl.HeldOutAccuracy,
			MaxUnitaryFeatureDrift:    featureDrift,
		},
	}, nil
}
