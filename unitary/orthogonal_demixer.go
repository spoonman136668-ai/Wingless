package unitary

import (
	"fmt"
	"math"
)

const OrthogonalDemixerSchema = "wingless.unitary-orthogonal-demixer.v1"

type GivensAngles [15]float64

type OrthogonalLearningTrace struct {
	Outer             int          `json:"outer"`
	Loss              float64      `json:"loss"`
	MeanRoleAlignment float64      `json:"mean_role_alignment"`
	MinRoleAlignment  float64      `json:"min_role_alignment"`
	Angles            GivensAngles `json:"angles"`
}

type OrthogonalDemixerDiagnosis struct {
	OracleRecoverabilityPass bool    `json:"oracle_recoverability_pass"`
	StructuredLearningPass   bool    `json:"structured_learning_pass"`
	UnseenDepthPass          bool    `json:"unseen_depth_pass"`
	MutableIntegrationPass   bool    `json:"mutable_integration_pass"`
	InitialCapacityAccuracy  float64 `json:"initial_capacity_accuracy"`
	FinalTrainAccuracy       float64 `json:"final_train_accuracy"`
	FinalHeldOutAccuracy     float64 `json:"final_heldout_accuracy"`
	MatchedControlAccuracy   float64 `json:"matched_control_accuracy"`
	OracleHeldOutAccuracy    float64 `json:"oracle_heldout_accuracy"`
	OracleFeatureError       float64 `json:"oracle_feature_error"`
	MeanRoleAlignment        float64 `json:"mean_role_alignment"`
	MinimumRoleAlignment     float64 `json:"minimum_role_alignment"`
	AngleL2Shift             float64 `json:"angle_l2_shift"`
	LearnedVsOracleMatrixL2  float64 `json:"learned_vs_oracle_matrix_l2"`
}

type OrthogonalDemixerProbeResult struct {
	Schema                         string                     `json:"schema"`
	Experiment                     string                     `json:"experiment"`
	CompositeDimension             int                        `json:"composite_dimension"`
	RuntimeStateObjects            int                        `json:"runtime_state_objects"`
	SemanticChannelLocationsKnown  bool                       `json:"semantic_channel_locations_known"`
	KnownMixerExposedToLearner     bool                       `json:"known_mixer_exposed_to_learner"`
	OracleUsedForLearning          bool                       `json:"oracle_used_for_learning"`
	RuntimeMixerInverseApplied     bool                       `json:"runtime_mixer_inverse_applied"`
	LearnedReadoutProjection       bool                       `json:"learned_readout_projection"`
	ReadoutProjectionOrthogonal    bool                       `json:"readout_projection_orthogonal"`
	TrainableDemixerParameters     int                        `json:"trainable_demixer_parameters"`
	RuntimePrototypeLookup         bool                       `json:"runtime_prototype_lookup"`
	ExplicitInverseTransportReadout bool                      `json:"explicit_inverse_transport_readout"`
	ExplicitDepthProvided          bool                       `json:"explicit_depth_provided"`
	GlobalPhaseNuisance            bool                       `json:"memory_only_global_phase_nuisance"`
	MemoryNoiseAmplitude           float64                    `json:"memory_noise_amplitude"`
	LearningTables                 int                        `json:"channel_learning_tables"`
	TrainTables                    int                        `json:"train_tables"`
	HeldOutTables                  int                        `json:"held_out_tables"`
	TrainDepths                    []int                      `json:"train_depths"`
	HeldOutDepths                  []int                      `json:"held_out_depths"`
	InitialAngles                  GivensAngles               `json:"initial_angles"`
	LearnedAngles                  GivensAngles               `json:"learned_angles"`
	InitialLoss                    float64                    `json:"initial_loss"`
	FinalLoss                      float64                    `json:"final_loss"`
	Trace                          []OrthogonalLearningTrace  `json:"trace"`
	OracleUnitary                  BlindChannelStaticResult   `json:"oracle_unitary"`
	Unitary                        BlindChannelStaticResult   `json:"unitary"`
	MatchedControl                 BlindChannelStaticResult   `json:"matched_control"`
	Integration                    BlindChannelIntegration    `json:"unitary_mutable_integration"`
	Diagnosis                      OrthogonalDemixerDiagnosis `json:"diagnosis"`
}

func channelPairs() [][2]int {
	pairs := make([][2]int, 0, 15)
	for first := 0; first < 6; first++ {
		for second := first + 1; second < 6; second++ {
			pairs = append(pairs, [2]int{first, second})
		}
	}
	return pairs
}

func mixerConstructionAngles() GivensAngles {
	var out GivensAngles
	for i := range out {
		out[i] = 0.43 + 0.071*float64(i)
	}
	return out
}

func transposeChannelMatrix(matrix ChannelMatrix) ChannelMatrix {
	var out ChannelMatrix
	for row := 0; row < 6; row++ {
		for column := 0; column < 6; column++ {
			out[row][column] = matrix[column][row]
		}
	}
	return out
}

// orthogonalDemixer parameterizes a complete six-dimensional orthogonal
// readout projection with 15 Givens angles.
//
// Rotations are applied in reverse pair order. Therefore the exact inverse of
// denseChannelMixer is representable by the negative construction angles,
// without supplying those target values to learning.
func orthogonalDemixer(angles GivensAngles) ChannelMatrix {
	out := identityChannelMatrix()
	pairs := channelPairs()

	for index := len(pairs) - 1; index >= 0; index-- {
		first := pairs[index][0]
		second := pairs[index][1]
		theta := angles[index]
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
	}
	return out
}

func oracleInverseAngles() GivensAngles {
	base := mixerConstructionAngles()
	var out GivensAngles
	for i := range out {
		out[i] = -base[i]
	}
	return out
}

func anglesL2(a, b GivensAngles) float64 {
	var total float64
	for i := range a {
		delta := a[i] - b[i]
		total += delta * delta
	}
	return math.Sqrt(total)
}

func channelMatrixL2(a, b ChannelMatrix) float64 {
	var total float64
	for row := 0; row < 6; row++ {
		for column := 0; column < 6; column++ {
			delta := a[row][column] - b[row][column]
			total += delta * delta
		}
	}
	return math.Sqrt(total)
}

func oracleFeatureRecoveryError() (float64, error) {
	table := memoryTable{3, 0, 2, 1}
	memory, err := encodeMemory(table)
	if err != nil {
		return 0, err
	}
	memory, err = perturbMemory(memory, 181818, 0.05)
	if err != nil {
		return 0, err
	}
	memory = rotateGlobalPhase(memory, 1.319)

	bank, err := qualifiedLearnedFrameBank()
	if err != nil {
		return 0, err
	}
	packed, err := packCompositeState(memory, bank)
	if err != nil {
		return 0, err
	}
	mixed, err := applyChannelMatrix(packed, denseChannelMixer())
	if err != nil {
		return 0, err
	}
	recovered, err := applyChannelMatrix(
		mixed, transposeChannelMatrix(denseChannelMixer()),
	)
	if err != nil {
		return 0, err
	}

	maximum := 0.0
	for entity := 0; entity < 4; entity++ {
		want, err := compositeEntityFeature(packed, entity)
		if err != nil {
			return 0, err
		}
		got, err := compositeEntityFeature(recovered, entity)
		if err != nil {
			return 0, err
		}
		for i := range want {
			delta := math.Abs(want[i] - got[i])
			if delta > maximum {
				maximum = delta
			}
		}
	}
	return maximum, nil
}

func orthogonalHeadsLoss(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	angles GivensAngles,
	memoryNoise float64,
	trials int,
) (float64, error) {
	return blindHeadsLoss(
		heads,
		tables,
		orthogonalDemixer(angles),
		memoryNoise,
		trials,
	)
}

func learnOrthogonalDemixer(
	tables []memoryTable,
	memoryNoise float64,
) (
	initial GivensAngles,
	learned GivensAngles,
	initialLoss float64,
	finalLoss float64,
	trace []OrthogonalLearningTrace,
	err error,
) {
	const (
		trials            = 2
		outerSteps        = 100
		headStepsPerOuter = 12
		headLearningRate  = 1.0
		angleLearningRate = 0.18
		angleEpsilon      = 1e-4
	)

	angles := GivensAngles{}
	initial = angles

	var heads [4]linearSoftmaxHead
	for entity := 0; entity < 4; entity++ {
		heads[entity] = newLinearSoftmaxHead(4, 2)
	}

	updateHeads := func() error {
		demixer := orthogonalDemixer(angles)
		for entity := 0; entity < 4; entity++ {
			samples, e := buildBlindCanonicalSamples(
				tables,
				entity,
				demixer,
				memoryNoise,
				trials,
				0,
			)
			if e != nil {
				return e
			}
			for step := 0; step < headStepsPerOuter; step++ {
				if e := updateSoftmaxHead(
					&heads[entity],
					samples,
					headLearningRate,
				); e != nil {
					return e
				}
			}
		}
		return nil
	}

	if err = updateHeads(); err != nil {
		return
	}
	initialLoss, err = orthogonalHeadsLoss(
		heads, tables, angles, memoryNoise, trials,
	)
	if err != nil {
		return
	}

	for outer := 1; outer <= outerSteps; outer++ {
		if err = updateHeads(); err != nil {
			return
		}

		var gradient GivensAngles
		for index := range angles {
			plus := angles
			minus := angles
			plus[index] += angleEpsilon
			minus[index] -= angleEpsilon

			plusLoss, e := orthogonalHeadsLoss(
				heads, tables, plus, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			minusLoss, e := orthogonalHeadsLoss(
				heads, tables, minus, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			gradient[index] =
				(plusLoss - minusLoss) / (2 * angleEpsilon)
			if !finite(gradient[index]) {
				err = fmt.Errorf(
					"orthogonal angle gradient index=%d non-finite",
					index,
				)
				return
			}
		}

		for index := range angles {
			angles[index] -= angleLearningRate * gradient[index]
			angles[index] = math.Remainder(angles[index], 2*math.Pi)
		}

		if outer == 1 || outer == 10 || outer == 20 ||
			outer == 40 || outer == 60 || outer == 80 ||
			outer == 100 {
			loss, e := orthogonalHeadsLoss(
				heads, tables, angles, memoryNoise, trials,
			)
			if e != nil {
				err = e
				return
			}
			alignments := roleAlignments(
				orthogonalDemixer(angles),
				denseChannelMixer(),
			)
			mean, minimum := blindMeanMin(alignments)
			trace = append(trace, OrthogonalLearningTrace{
				Outer:             outer,
				Loss:              loss,
				MeanRoleAlignment: mean,
				MinRoleAlignment:  minimum,
				Angles:            angles,
			})
		}
	}

	finalLoss, err = orthogonalHeadsLoss(
		heads, tables, angles, memoryNoise, trials,
	)
	if err != nil {
		return
	}
	learned = angles
	return
}

// RunUP18 distinguishes two hypotheses left open by the UP-17 negative:
//
//  1. the anonymous mixer destroyed task-relevant information;
//  2. information survived, but UP-17's unconstrained row-normalized demixer
//     was a poor optimization geometry.
//
// The oracle transpose is used only as a recoverability control. It is never
// supplied to the task learner or runtime learned path.
func RunUP18() (OrthogonalDemixerProbeResult, error) {
	const memoryNoise = 0.05

	learningTables, err := selectObserverTables(true, 64)
	if err != nil {
		return OrthogonalDemixerProbeResult{}, err
	}
	trainTables := fullObserverTablePool(true)
	heldTables := fullObserverTablePool(false)
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	block := stressProgram()

	oracleMatrix := transposeChannelMatrix(denseChannelMixer())
	oracleFeatureError, err := oracleFeatureRecoveryError()
	if err != nil {
		return OrthogonalDemixerProbeResult{}, err
	}

	_, oracleStatic, err := trainBlindTransportHeads(
		"unitary_oracle_anonymous_channel_recovery",
		"unitary_oracle_control",
		oracleMatrix,
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
		return OrthogonalDemixerProbeResult{}, err
	}

	initialMatrix := orthogonalDemixer(GivensAngles{})
	_, initialCapacity, err := fitBlindHeads(
		learningTables,
		initialMatrix,
		memoryNoise,
		2,
		900,
	)
	if err != nil {
		return OrthogonalDemixerProbeResult{}, err
	}

	initialAngles, learnedAngles, initialLoss, finalLoss, trace, err :=
		learnOrthogonalDemixer(learningTables, memoryNoise)
	if err != nil {
		return OrthogonalDemixerProbeResult{}, err
	}

	learnedMatrix := orthogonalDemixer(learnedAngles)
	alignments := roleAlignments(learnedMatrix, denseChannelMixer())
	meanAlignment, minimumAlignment := blindMeanMin(alignments)

	unitaryHeads, unitaryStatic, err := trainBlindTransportHeads(
		"unitary_structured_orthogonal_demixer",
		"unitary",
		learnedMatrix,
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
		return OrthogonalDemixerProbeResult{}, err
	}

	_, controlStatic, err := trainBlindTransportHeads(
		"non_unitary_structured_orthogonal_demixer",
		"non_unitary_matched",
		learnedMatrix,
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
		return OrthogonalDemixerProbeResult{}, err
	}

	integration, err := runBlindIntegration(
		learnedMatrix,
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
	)
	if err != nil {
		return OrthogonalDemixerProbeResult{}, err
	}

	oraclePass :=
		oracleFeatureError <= 1e-12 &&
			oracleStatic.HeldOutAccuracy >= 0.99 &&
			oracleStatic.MaxNormDrift <= 1e-12

	learningPass :=
		initialCapacity < 0.70 &&
			finalLoss < initialLoss &&
			anglesL2(initialAngles, learnedAngles) >= 0.50 &&
			unitaryStatic.TrainAccuracy >= 0.99

	unseenPass :=
		unitaryStatic.HeldOutAccuracy >= 0.99 &&
			unitaryStatic.MaxNormDrift <= 1e-12

	mutablePass :=
		integration.CommitDecodeAccuracy >= 0.99 &&
			integration.ExactFinalTableAccuracy >= 0.95 &&
			integration.RelationalQueryAccuracy >= 0.95 &&
			integration.MaxNormDrift <= 1e-12

	return OrthogonalDemixerProbeResult{
		Schema:                          OrthogonalDemixerSchema,
		Experiment:                      "UP-18-orthogonal-anonymous-channel-demixer",
		CompositeDimension:              compositeDimension,
		RuntimeStateObjects:             1,
		SemanticChannelLocationsKnown:   false,
		KnownMixerExposedToLearner:      false,
		OracleUsedForLearning:           false,
		RuntimeMixerInverseApplied:      false,
		LearnedReadoutProjection:        true,
		ReadoutProjectionOrthogonal:     true,
		TrainableDemixerParameters:      15,
		RuntimePrototypeLookup:          false,
		ExplicitInverseTransportReadout: false,
		ExplicitDepthProvided:           false,
		GlobalPhaseNuisance:             true,
		MemoryNoiseAmplitude:            memoryNoise,
		LearningTables:                  len(learningTables),
		TrainTables:                     len(trainTables),
		HeldOutTables:                   len(heldTables),
		TrainDepths:                     append([]int(nil), trainDepths...),
		HeldOutDepths:                   append([]int(nil), heldDepths...),
		InitialAngles:                   initialAngles,
		LearnedAngles:                   learnedAngles,
		InitialLoss:                     initialLoss,
		FinalLoss:                       finalLoss,
		Trace:                           trace,
		OracleUnitary:                   oracleStatic,
		Unitary:                         unitaryStatic,
		MatchedControl:                  controlStatic,
		Integration:                     integration,
		Diagnosis: OrthogonalDemixerDiagnosis{
			OracleRecoverabilityPass: oraclePass,
			StructuredLearningPass:   learningPass,
			UnseenDepthPass:          unseenPass,
			MutableIntegrationPass:   mutablePass,
			InitialCapacityAccuracy:  initialCapacity,
			FinalTrainAccuracy:       unitaryStatic.TrainAccuracy,
			FinalHeldOutAccuracy:     unitaryStatic.HeldOutAccuracy,
			MatchedControlAccuracy:   controlStatic.HeldOutAccuracy,
			OracleHeldOutAccuracy:    oracleStatic.HeldOutAccuracy,
			OracleFeatureError:       oracleFeatureError,
			MeanRoleAlignment:        meanAlignment,
			MinimumRoleAlignment:     minimumAlignment,
			AngleL2Shift:             anglesL2(initialAngles, learnedAngles),
			LearnedVsOracleMatrixL2:  channelMatrixL2(learnedMatrix, oracleMatrix),
		},
	}, nil
}
