package unitary

import (
	"fmt"
	"math"
)

const CompressedFrameSchema = "wingless.unitary-compressed-frame.v1"

type CompressedFrameCell struct {
	Name                   string    `json:"name"`
	Path                   string    `json:"path"`
	ReferenceNoise         float64   `json:"reference_noise"`
	ReferenceTransportScale float64  `json:"reference_transport_scale"`
	TrainAccuracy          float64   `json:"train_accuracy"`
	HeldOutAccuracy        float64   `json:"held_out_accuracy"`
	PerEntityHeldOut       []float64 `json:"per_entity_held_out_accuracy"`
}

type CompressedFrameIntegration struct {
	Scenarios               int     `json:"scenarios"`
	WritesPerScenario       int     `json:"writes_per_scenario"`
	ReferenceNoise          float64 `json:"reference_noise"`
	CommitDecodeAccuracy    float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy float64 `json:"relational_query_accuracy"`
	MinValueMargin          float64 `json:"min_value_margin"`
	MinRelationMargin       float64 `json:"min_relation_margin"`
	MaxNormDrift            float64 `json:"max_norm_drift"`
}

type CompressedFrameDiagnosis struct {
	CompressionPass       bool    `json:"compression_pass"`
	NoisyFramePass        bool    `json:"noisy_frame_pass"`
	MutableProgramPass    bool    `json:"mutable_program_pass"`
	FullFiveAccuracy      float64 `json:"full_five_accuracy"`
	CompressedCleanAccuracy float64 `json:"compressed_clean_accuracy"`
	ReferenceNoise01Accuracy float64 `json:"reference_noise_0_01_accuracy"`
	ReferenceNoise03Accuracy float64 `json:"reference_noise_0_03_accuracy"`
	Drift1e4Accuracy      float64 `json:"reference_drift_1e_4_accuracy"`
	Drift1e3Accuracy      float64 `json:"reference_drift_1e_3_accuracy"`
	PilotStateReduction   float64 `json:"pilot_state_reduction"`
}

type CompressedFrameProbeResult struct {
	Schema                   string                     `json:"schema"`
	Experiment               string                     `json:"experiment"`
	Dimension                int                        `json:"dimension"`
	Entities                 int                        `json:"entities"`
	ValuesPerEntity          int                        `json:"values_per_entity"`
	FullPilotStates          int                        `json:"full_pilot_states"`
	FullPilotComplexScalars  int                        `json:"full_pilot_complex_scalars"`
	CompressedPilotStates    int                        `json:"compressed_pilot_states"`
	CompressedComplexScalars int                        `json:"compressed_pilot_complex_scalars"`
	FeaturesPerEntity        int                        `json:"features_per_entity"`
	RuntimePrototypeLookup   bool                       `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout   bool                       `json:"explicit_inverse_readout"`
	ExplicitDepthProvided    bool                       `json:"explicit_depth_provided"`
	SelectedTrainTables      int                        `json:"selected_train_tables"`
	SelectedHeldOutTables    int                        `json:"selected_held_out_tables"`
	MinTrainMarginalCount    int                        `json:"min_train_marginal_count"`
	MinHeldOutMarginalCount  int                        `json:"min_held_out_marginal_count"`
	TrainDepths              []int                      `json:"train_depths"`
	HeldOutDepths            []int                      `json:"held_out_depths"`
	Cells                    []CompressedFrameCell      `json:"cells"`
	Integration              CompressedFrameIntegration `json:"integration"`
	Diagnosis                CompressedFrameDiagnosis   `json:"diagnosis"`
}

type compressedFrameBank struct {
	anchor State
	codeA  State
	codeB  State
}

func makeCompressedFrameBank() (compressedFrameBank, error) {
	full, err := makeFrameBank()
	if err != nil {
		return compressedFrameBank{}, err
	}
	levels := []float64{-3, -1, 1, 3}
	codeA := make(State, 16)
	codeB := make(State, 16)
	for value, level := range levels {
		codeA[value] = complex(level, 0)
		codeA[4+value] = complex(0, level)
		codeB[8+value] = complex(level, 0)
		codeB[12+value] = complex(0, level)
	}
	codeA, err = Normalize(codeA)
	if err != nil {
		return compressedFrameBank{}, err
	}
	codeB, err = Normalize(codeB)
	if err != nil {
		return compressedFrameBank{}, err
	}
	return compressedFrameBank{
		anchor: full.anchor,
		codeA:  codeA,
		codeB:  codeB,
	}, nil
}

func compressedFrameFeature(memory State, bank compressedFrameBank, entity int) ([]float64, error) {
	if entity < 0 || entity >= 4 {
		return nil, fmt.Errorf("compressed-frame entity out of range")
	}
	anchor, err := stateInner(bank.anchor, memory)
	if err != nil {
		return nil, err
	}
	if cmplxAbsSquared(anchor) <= 1e-18 {
		return nil, fmt.Errorf("compressed-frame anchor correlation too small")
	}
	var code complex128
	if entity < 2 {
		code, err = stateInner(bank.codeA, memory)
	} else {
		code, err = stateInner(bank.codeB, memory)
	}
	if err != nil {
		return nil, err
	}
	ratio := code / anchor
	if !finite(real(ratio)) || !finite(imag(ratio)) {
		return nil, fmt.Errorf("compressed-frame ratio not finite")
	}
	switch entity {
	case 0, 2:
		return []float64{real(ratio)}, nil
	case 1, 3:
		return []float64{-imag(ratio)}, nil
	default:
		panic("unreachable")
	}
}

func scaleStressBlock(block []Coupling, factor float64) ([]Coupling, error) {
	if !finite(factor) || factor <= 0 {
		return nil, fmt.Errorf("invalid stress-block scale")
	}
	out := make([]Coupling, len(block))
	for i, coupling := range block {
		out[i] = coupling
		out[i].Theta *= factor
	}
	return out, nil
}

func perturbCompressedFrame(bank compressedFrameBank, seed int, amplitude float64) (compressedFrameBank, error) {
	if amplitude == 0 {
		return bank, nil
	}
	anchor, err := perturbMemory(bank.anchor, seed+11, amplitude)
	if err != nil {
		return compressedFrameBank{}, err
	}
	codeA, err := perturbMemory(bank.codeA, seed+23, amplitude)
	if err != nil {
		return compressedFrameBank{}, err
	}
	codeB, err := perturbMemory(bank.codeB, seed+37, amplitude)
	if err != nil {
		return compressedFrameBank{}, err
	}
	return compressedFrameBank{anchor: anchor, codeA: codeA, codeB: codeB}, nil
}

func transportCompressedFrame(
	bank compressedFrameBank,
	block []Coupling,
	depth int,
	apply stressApply,
	referenceScale float64,
) (compressedFrameBank, error) {
	referenceBlock := block
	var err error
	if referenceScale != 1 {
		referenceBlock, err = scaleStressBlock(block, referenceScale)
		if err != nil {
			return compressedFrameBank{}, err
		}
	}
	anchor, err := apply(bank.anchor, referenceBlock, depth)
	if err != nil {
		return compressedFrameBank{}, err
	}
	codeA, err := apply(bank.codeA, referenceBlock, depth)
	if err != nil {
		return compressedFrameBank{}, err
	}
	codeB, err := apply(bank.codeB, referenceBlock, depth)
	if err != nil {
		return compressedFrameBank{}, err
	}
	return compressedFrameBank{anchor: anchor, codeA: codeA, codeB: codeB}, nil
}

func buildCompressedSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	block []Coupling,
	apply stressApply,
	memoryNoise, referenceNoise, referenceScale float64,
	seedOffset int,
) ([]headSample, error) {
	baseBank, err := makeCompressedFrameBank()
	if err != nil {
		return nil, err
	}
	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for depthIndex, depth := range depths {
			seed := seedOffset + memoryTableIndex(table)*10000 + depthIndex*137 + entity*19
			noisyMemory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, err
			}
			globalPhase := math.Mod(0.191*float64(seed+1), 2*math.Pi)
			noisyMemory = rotateGlobalPhase(noisyMemory, globalPhase)

			forward, err := apply(noisyMemory, block, depth)
			if err != nil {
				return nil, err
			}

			frame, err := perturbCompressedFrame(baseBank, seed+500000, referenceNoise)
			if err != nil {
				return nil, err
			}
			frame, err = transportCompressedFrame(frame, block, depth, apply, referenceScale)
			if err != nil {
				return nil, err
			}

			features, err := compressedFrameFeature(forward, frame, entity)
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

func trainCompressedHeads(
	tables []memoryTable,
	depths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
) ([4]linearSoftmaxHead, float64, error) {
	var heads [4]linearSoftmaxHead
	var average float64
	for entity := 0; entity < 4; entity++ {
		samples, err := buildCompressedSamples(
			tables, depths, entity, block, apply,
			memoryNoise, 0, 1, 0,
		)
		if err != nil {
			return heads, 0, err
		}
		head, metric, err := trainLinearSoftmax(samples, 4, 1, 300, 1.0)
		if err != nil {
			return heads, 0, err
		}
		heads[entity] = head
		average += metric.TrainAccuracy
	}
	return heads, average / 4, nil
}

func evaluateCompressedHeads(
	heads [4]linearSoftmaxHead,
	tables []memoryTable,
	depths []int,
	block []Coupling,
	apply stressApply,
	memoryNoise, referenceNoise, referenceScale float64,
) (float64, []float64, error) {
	perEntity := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		samples, err := buildCompressedSamples(
			tables, depths, entity, block, apply,
			memoryNoise, referenceNoise, referenceScale, 800000,
		)
		if err != nil {
			return 0, nil, err
		}
		_, accuracy, err := evaluateHead(heads[entity], samples)
		if err != nil {
			return 0, nil, err
		}
		perEntity[entity] = accuracy
	}
	var average float64
	for _, accuracy := range perEntity {
		average += accuracy
	}
	return average / 4, perEntity, nil
}

func decodeCompressedTable(
	memory State,
	bank compressedFrameBank,
	heads [4]linearSoftmaxHead,
) (memoryTable, [4][]float64, float64, error) {
	var table memoryTable
	var distributions [4][]float64
	minMargin := math.Inf(1)
	for entity := 0; entity < 4; entity++ {
		features, err := compressedFrameFeature(memory, bank, entity)
		if err != nil {
			return memoryTable{}, distributions, 0, err
		}
		probabilities, err := heads[entity].probabilities(features)
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

func makeCompressedScenario(seed int, initial memoryTable, writes int, depths []int) memoryScenario {
	ops := make([]memoryWrite, 0, writes)
	for j := 0; j < writes; j++ {
		entity := (seed*3 + j*2 + j/3) % 4
		value := (seed + j*j + 2*j + entity + 1) % 4
		gap := depths[(seed+j)%len(depths)]
		ops = append(ops, memoryWrite{entity: entity, value: value, gap: gap})
	}
	return memoryScenario{
		initial:  initial,
		writes:   ops,
		queryA:   (seed + 1) % 4,
		queryB:   (seed + 3) % 4,
		finalGap: depths[(seed+writes)%len(depths)],
	}
}

func runCompressedIntegration(
	heads [4]linearSoftmaxHead,
	heldTables []memoryTable,
	depths []int,
	block []Coupling,
	memoryNoise, referenceNoise float64,
) (CompressedFrameIntegration, error) {
	const (
		scenarios = 32
		writes    = 12
	)
	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return CompressedFrameIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(relationSamples, 4, 16, 400, 1.0)
	if err != nil {
		return CompressedFrameIntegration{}, err
	}

	var commitCorrect, commitTotal, finalCorrect, relationCorrect int
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)
	var maxNormDrift float64
	baseBank, err := makeCompressedFrameBank()
	if err != nil {
		return CompressedFrameIntegration{}, err
	}

	for scenarioIndex := 0; scenarioIndex < scenarios; scenarioIndex++ {
		scenario := makeCompressedScenario(
			scenarioIndex,
			heldTables[scenarioIndex%len(heldTables)],
			writes,
			depths,
		)
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			seed := 1200000 + scenarioIndex*1000 + writeIndex
			noisyMemory, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			noisyMemory = rotateGlobalPhase(noisyMemory, math.Mod(0.191*float64(seed+1), 2*math.Pi))
			forward, err := applyStressUnitary(noisyMemory, block, write.gap)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			norm2, err := NormSquared(forward)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			frame, err := perturbCompressedFrame(baseBank, seed+500000, referenceNoise)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			frame, err = transportCompressedFrame(frame, block, write.gap, applyStressUnitary, 1)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
			decoded, _, margin, err := decodeCompressedTable(forward, frame, heads)
			if err != nil {
				return CompressedFrameIntegration{}, err
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
				return CompressedFrameIntegration{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return CompressedFrameIntegration{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		seed := 1200000 + scenarioIndex*1000 + 999
		noisyMemory, err := perturbMemory(canonical, seed, memoryNoise)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		noisyMemory = rotateGlobalPhase(noisyMemory, math.Mod(0.191*float64(seed+1), 2*math.Pi))
		forward, err := applyStressUnitary(noisyMemory, block, scenario.finalGap)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		norm2, err := NormSquared(forward)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}
		frame, err := perturbCompressedFrame(baseBank, seed+500000, referenceNoise)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		frame, err = transportCompressedFrame(frame, block, scenario.finalGap, applyStressUnitary, 1)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		decoded, distributions, margin, err := decodeCompressedTable(forward, frame, heads)
		if err != nil {
			return CompressedFrameIntegration{}, err
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
			return CompressedFrameIntegration{}, err
		}
		relationProbabilities, err := relationHead.probabilities(relationInput)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		gotRelation, relationMargin, err := classAndMargin(relationProbabilities)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		if relationMargin < minRelationMargin {
			minRelationMargin = relationMargin
		}
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil {
			return CompressedFrameIntegration{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	return CompressedFrameIntegration{
		Scenarios:               scenarios,
		WritesPerScenario:       writes,
		ReferenceNoise:          referenceNoise,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalCorrect) / float64(scenarios),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(scenarios),
		MinValueMargin:          minValueMargin,
		MinRelationMargin:       minRelationMargin,
		MaxNormDrift:            maxNormDrift,
	}, nil
}

// RunUP10 compresses the five-pilot UP-9 frame to three co-evolving states:
//
//   - one global anchor;
//   - one multiplexed phase-code pilot for entities 0/1;
//   - one multiplexed phase-code pilot for entities 2/3.
//
// The runtime pilot-state footprint therefore falls from 80 to 48 complex
// scalars (40%). Each entity readout uses one real scalar extracted from a
// global-phase-invariant pilot/anchor ratio.
//
// The experiment then measures clean unseen-depth readout, independent pilot
// corruption, reference-transport mismatch, a matched non-unitary control, and
// a full mutable read/write integration using the compressed frame.
func RunUP10() (CompressedFrameProbeResult, error) {
	const (
		selectedTables = 32
		memoryNoise    = 0.05
	)
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}
	trainTables, err := selectObserverTables(true, selectedTables)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	heldTables, err := selectObserverTables(false, selectedTables)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	minTrain := observerMarginalMinimum(trainTables)
	minHeld := observerMarginalMinimum(heldTables)
	if minTrain <= 0 || minHeld <= 0 {
		return CompressedFrameProbeResult{}, fmt.Errorf("compressed-frame table selection lost a marginal")
	}

	block := stressProgram()

	fullCell, err := runFrameCell(
		"unitary_full5_clean",
		"unitary",
		"coevolving",
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressUnitary,
		true,
		memoryNoise,
	)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}

	unitaryHeads, unitaryTrainAccuracy, err := trainCompressedHeads(
		trainTables, trainDepths, block, applyStressUnitary, memoryNoise,
	)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	controlHeads, controlTrainAccuracy, err := trainCompressedHeads(
		trainTables, trainDepths, block, applyStressNonUnitary, memoryNoise,
	)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}

	type cellSpec struct {
		name          string
		path          string
		heads         [4]linearSoftmaxHead
		trainAccuracy float64
		apply         stressApply
		refNoise      float64
		refScale      float64
	}
	specs := []cellSpec{
		{"unitary_compressed3_clean", "unitary", unitaryHeads, unitaryTrainAccuracy, applyStressUnitary, 0, 1},
		{"unitary_compressed3_refnoise_0_01", "unitary", unitaryHeads, unitaryTrainAccuracy, applyStressUnitary, 0.01, 1},
		{"unitary_compressed3_refnoise_0_03", "unitary", unitaryHeads, unitaryTrainAccuracy, applyStressUnitary, 0.03, 1},
		{"unitary_compressed3_drift_1e_4", "unitary", unitaryHeads, unitaryTrainAccuracy, applyStressUnitary, 0, 1.0001},
		{"unitary_compressed3_drift_1e_3", "unitary", unitaryHeads, unitaryTrainAccuracy, applyStressUnitary, 0, 1.001},
		{"control_compressed3_clean", "non_unitary_matched", controlHeads, controlTrainAccuracy, applyStressNonUnitary, 0, 1},
	}

	cells := make([]CompressedFrameCell, 0, len(specs)+1)
	cells = append(cells, CompressedFrameCell{
		Name:                    fullCell.Name,
		Path:                    fullCell.Path,
		ReferenceNoise:          0,
		ReferenceTransportScale: 1,
		TrainAccuracy:           fullCell.TrainAccuracy,
		HeldOutAccuracy:         fullCell.HeldOutAccuracy,
		PerEntityHeldOut:        append([]float64(nil), fullCell.PerEntityHeldOut...),
	})

	for _, spec := range specs {
		accuracy, perEntity, err := evaluateCompressedHeads(
			spec.heads,
			heldTables,
			heldDepths,
			block,
			spec.apply,
			memoryNoise,
			spec.refNoise,
			spec.refScale,
		)
		if err != nil {
			return CompressedFrameProbeResult{}, err
		}
		cells = append(cells, CompressedFrameCell{
			Name:                    spec.name,
			Path:                    spec.path,
			ReferenceNoise:          spec.refNoise,
			ReferenceTransportScale: spec.refScale,
			TrainAccuracy:           spec.trainAccuracy,
			HeldOutAccuracy:         accuracy,
			PerEntityHeldOut:        perEntity,
		})
	}

	integration, err := runCompressedIntegration(
		unitaryHeads,
		heldTables,
		heldDepths,
		block,
		memoryNoise,
		0.01,
	)
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}

	find := func(name string) (CompressedFrameCell, error) {
		for _, cell := range cells {
			if cell.Name == name {
				return cell, nil
			}
		}
		return CompressedFrameCell{}, fmt.Errorf("compressed-frame cell %q missing", name)
	}
	clean, err := find("unitary_compressed3_clean")
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	noise01, err := find("unitary_compressed3_refnoise_0_01")
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	noise03, err := find("unitary_compressed3_refnoise_0_03")
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	drift1e4, err := find("unitary_compressed3_drift_1e_4")
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}
	drift1e3, err := find("unitary_compressed3_drift_1e_3")
	if err != nil {
		return CompressedFrameProbeResult{}, err
	}

	diagnosis := CompressedFrameDiagnosis{
		CompressionPass:          clean.HeldOutAccuracy >= 0.95 && fullCell.HeldOutAccuracy >= 0.95,
		NoisyFramePass:           noise01.HeldOutAccuracy >= 0.95,
		MutableProgramPass:       integration.CommitDecodeAccuracy >= 0.95 &&
			integration.ExactFinalTableAccuracy >= 0.90 &&
			integration.RelationalQueryAccuracy >= 0.90,
		FullFiveAccuracy:         fullCell.HeldOutAccuracy,
		CompressedCleanAccuracy:  clean.HeldOutAccuracy,
		ReferenceNoise01Accuracy: noise01.HeldOutAccuracy,
		ReferenceNoise03Accuracy: noise03.HeldOutAccuracy,
		Drift1e4Accuracy:         drift1e4.HeldOutAccuracy,
		Drift1e3Accuracy:         drift1e3.HeldOutAccuracy,
		PilotStateReduction:      1 - 3.0/5.0,
	}

	return CompressedFrameProbeResult{
		Schema:                   CompressedFrameSchema,
		Experiment:               "UP-10-compressed-noisy-internal-frame",
		Dimension:                16,
		Entities:                 4,
		ValuesPerEntity:          4,
		FullPilotStates:          5,
		FullPilotComplexScalars:  80,
		CompressedPilotStates:    3,
		CompressedComplexScalars: 48,
		FeaturesPerEntity:        1,
		RuntimePrototypeLookup:   false,
		ExplicitInverseReadout:   false,
		ExplicitDepthProvided:    false,
		SelectedTrainTables:      len(trainTables),
		SelectedHeldOutTables:    len(heldTables),
		MinTrainMarginalCount:    minTrain,
		MinHeldOutMarginalCount:  minHeld,
		TrainDepths:              append([]int(nil), trainDepths...),
		HeldOutDepths:            append([]int(nil), heldDepths...),
		Cells:                    cells,
		Integration:              integration,
		Diagnosis:                diagnosis,
	}, nil
}
