package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const PilotRankProbeSchema = "wingless.unitary-pilot-rank-probe.v1"

type PilotRankCell struct {
	Name                string    `json:"name"`
	Path                string    `json:"path"`
	CodeRank            int       `json:"code_rank"`
	PilotStates         int       `json:"pilot_states"`
	PilotComplexScalars int       `json:"pilot_complex_scalars"`
	FeaturesPerEntity   int       `json:"features_per_entity"`
	MemoryNoise         float64   `json:"memory_noise"`
	TrainAccuracy       float64   `json:"train_accuracy"`
	HeldOutAccuracy     float64   `json:"held_out_accuracy"`
	PerEntityTrain      []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut    []float64 `json:"per_entity_held_out_accuracy"`
}

type PilotRankDiagnosis struct {
	MinimalNoiselessRank         int     `json:"minimal_noiseless_rank"`
	MinimalNoise005Rank          int     `json:"minimal_noise_0_05_rank"`
	Rank2NoiselessAccuracy       float64 `json:"rank2_noiseless_accuracy"`
	Rank2Noise005Accuracy        float64 `json:"rank2_noise_0_05_accuracy"`
	Rank3NoiselessAccuracy       float64 `json:"rank3_noiseless_accuracy"`
	Rank3Noise005Accuracy        float64 `json:"rank3_noise_0_05_accuracy"`
	Rank4NoiselessAccuracy       float64 `json:"rank4_noiseless_accuracy"`
	Rank4Noise005Accuracy        float64 `json:"rank4_noise_0_05_accuracy"`
	IntrinsicRankLossSupported   bool    `json:"intrinsic_rank_loss_supported"`
	MemoryNoiseSensitivity       bool    `json:"memory_noise_sensitivity_supported"`
	FourPilotCandidateSupported  bool    `json:"four_pilot_candidate_supported"`
}

type PilotRankProbeResult struct {
	Schema                  string             `json:"schema"`
	Experiment              string             `json:"experiment"`
	Dimension               int                `json:"dimension"`
	Entities                int                `json:"entities"`
	ValuesPerEntity         int                `json:"values_per_entity"`
	RuntimePrototypeLookup  bool               `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout  bool               `json:"explicit_inverse_readout"`
	ExplicitDepthProvided   bool               `json:"explicit_depth_provided"`
	GlobalPhaseNuisance     bool               `json:"global_phase_nuisance"`
	SelectedTrainTables     int                `json:"selected_train_tables"`
	SelectedHeldOutTables   int                `json:"selected_held_out_tables"`
	MinTrainMarginalCount   int                `json:"min_train_marginal_count"`
	MinHeldOutMarginalCount int                `json:"min_held_out_marginal_count"`
	TrainDepths             []int              `json:"train_depths"`
	HeldOutDepths           []int              `json:"held_out_depths"`
	Cells                   []PilotRankCell    `json:"cells"`
	Diagnosis               PilotRankDiagnosis `json:"diagnosis"`
}

type pilotRankBank struct {
	anchor State
	codes  []State
}

func makePilotRankBank(rank int) (pilotRankBank, error) {
	if rank < 1 || rank > 4 {
		return pilotRankBank{}, fmt.Errorf("pilot code rank=%d want 1..4", rank)
	}
	full, err := makeFrameBank()
	if err != nil {
		return pilotRankBank{}, err
	}

	codes := make([]State, rank)
	for k := 0; k < rank; k++ {
		code := make(State, 16)
		for entity := 0; entity < 4; entity++ {
			coefficient := cmplx.Rect(1, 2*math.Pi*float64(k*entity)/4)
			for i := range code {
				code[i] += coefficient * full.phase[entity][i]
			}
		}
		code, err = Normalize(code)
		if err != nil {
			return pilotRankBank{}, err
		}
		codes[k] = code
	}
	return pilotRankBank{
		anchor: full.anchor,
		codes:  codes,
	}, nil
}

func transportPilotRankBank(
	bank pilotRankBank,
	block []Coupling,
	depth int,
	apply stressApply,
) (pilotRankBank, error) {
	anchor, err := apply(bank.anchor, block, depth)
	if err != nil {
		return pilotRankBank{}, err
	}
	codes := make([]State, len(bank.codes))
	for i, code := range bank.codes {
		codes[i], err = apply(code, block, depth)
		if err != nil {
			return pilotRankBank{}, err
		}
	}
	return pilotRankBank{anchor: anchor, codes: codes}, nil
}

func pilotRankFeatures(memory State, bank pilotRankBank) ([]float64, error) {
	anchor, err := stateInner(bank.anchor, memory)
	if err != nil {
		return nil, err
	}
	if cmplx.Abs(anchor) <= 1e-15 {
		return nil, fmt.Errorf("pilot-rank anchor correlation too small")
	}

	features := make([]float64, 0, 2*len(bank.codes))
	for _, code := range bank.codes {
		correlation, err := stateInner(code, memory)
		if err != nil {
			return nil, err
		}
		ratio := correlation / anchor
		if !finite(real(ratio)) || !finite(imag(ratio)) {
			return nil, fmt.Errorf("pilot-rank feature is not finite")
		}
		features = append(features, real(ratio), imag(ratio))
	}
	return features, nil
}

func buildPilotRankSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	rank int,
	block []Coupling,
	apply stressApply,
	memoryNoise float64,
	seedOffset int,
) ([]headSample, error) {
	baseBank, err := makePilotRankBank(rank)
	if err != nil {
		return nil, err
	}
	transformed := make(map[int]pilotRankBank, len(depths))
	for _, depth := range depths {
		bank, err := transportPilotRankBank(baseBank, block, depth, apply)
		if err != nil {
			return nil, err
		}
		transformed[depth] = bank
	}

	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for depthIndex, depth := range depths {
			seed := seedOffset + memoryTableIndex(table)*10000 + depthIndex*149 + entity*23 + rank*1000000
			state := canonical
			if memoryNoise > 0 {
				state, err = perturbMemory(canonical, seed, memoryNoise)
				if err != nil {
					return nil, err
				}
			} else {
				state = append(State(nil), canonical...)
			}
			state = rotateGlobalPhase(state, math.Mod(0.211*float64(seed+1), 2*math.Pi))
			forward, err := apply(state, block, depth)
			if err != nil {
				return nil, err
			}
			features, err := pilotRankFeatures(forward, transformed[depth])
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

func runPilotRankCell(
	name, path string,
	rank int,
	memoryNoise float64,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
) (PilotRankCell, error) {
	const (
		steps        = 250
		learningRate = 1.0
	)
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)

	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildPilotRankSamples(
			trainTables, trainDepths, entity, rank, block, apply,
			memoryNoise, 0,
		)
		if err != nil {
			return PilotRankCell{}, err
		}
		heldSamples, err := buildPilotRankSamples(
			heldTables, heldDepths, entity, rank, block, apply,
			memoryNoise, 900000,
		)
		if err != nil {
			return PilotRankCell{}, err
		}
		head, _, err := trainLinearSoftmax(
			trainSamples, 4, 2*rank, steps, learningRate,
		)
		if err != nil {
			return PilotRankCell{}, err
		}
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return PilotRankCell{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return PilotRankCell{}, err
		}
		perTrain[entity] = trainAccuracy
		perHeld[entity] = heldAccuracy
	}

	var trainAverage, heldAverage float64
	for i := 0; i < 4; i++ {
		trainAverage += perTrain[i]
		heldAverage += perHeld[i]
	}
	trainAverage /= 4
	heldAverage /= 4

	return PilotRankCell{
		Name:                name,
		Path:                path,
		CodeRank:            rank,
		PilotStates:         rank + 1,
		PilotComplexScalars: (rank + 1) * 16,
		FeaturesPerEntity:   2 * rank,
		MemoryNoise:         memoryNoise,
		TrainAccuracy:       trainAverage,
		HeldOutAccuracy:     heldAverage,
		PerEntityTrain:      perTrain,
		PerEntityHeldOut:    perHeld,
	}, nil
}

func pilotRankCellByName(cells []PilotRankCell, name string) (PilotRankCell, error) {
	for _, cell := range cells {
		if cell.Name == name {
			return cell, nil
		}
	}
	return PilotRankCell{}, fmt.Errorf("pilot-rank cell %q missing", name)
}

func minimalPassingPilotRank(cells []PilotRankCell, noise float64) int {
	minimum := 0
	for _, cell := range cells {
		if cell.Path != "unitary" || math.Abs(cell.MemoryNoise-noise) > 1e-15 {
			continue
		}
		if cell.HeldOutAccuracy < 0.95 {
			continue
		}
		if minimum == 0 || cell.CodeRank < minimum {
			minimum = cell.CodeRank
		}
	}
	return minimum
}

// RunUP11 isolates pilot rank from memory-noise sensitivity.
//
// The four successful UP-9 entity-phase pilots are treated as an orthonormal
// rank-4 reference family. For code rank r=1..4, UP-11 constructs the first r
// equal-energy Fourier mixtures of those pilots, plus the common anchor.
//
// Each entity head receives the complete complex pilot/anchor ratios: 2r real
// features. This avoids the one-scalar projection used by UP-10 and asks a
// cleaner question: how many independent co-evolving code directions are
// required to decode the four-entity memory?
//
// Every rank is tested both noiselessly and with the standing memory-noise
// amplitude 0.05, across the unseen-depth split. A matched non-unitary rank-4
// noisy control is also measured.
func RunUP11() (PilotRankProbeResult, error) {
	const selectedTables = 32
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}

	trainTables, err := selectObserverTables(true, selectedTables)
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	heldTables, err := selectObserverTables(false, selectedTables)
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	minTrain := observerMarginalMinimum(trainTables)
	minHeld := observerMarginalMinimum(heldTables)
	if minTrain <= 0 || minHeld <= 0 {
		return PilotRankProbeResult{}, fmt.Errorf("pilot-rank table split lost a marginal")
	}

	block := stressProgram()
	cells := make([]PilotRankCell, 0, 9)
	for rank := 1; rank <= 4; rank++ {
		for _, noise := range []float64{0, 0.05} {
			name := fmt.Sprintf("unitary_rank%d_noise_%0.2f", rank, noise)
			cell, err := runPilotRankCell(
				name, "unitary", rank, noise,
				trainTables, heldTables,
				trainDepths, heldDepths,
				block, applyStressUnitary,
			)
			if err != nil {
				return PilotRankProbeResult{}, err
			}
			cells = append(cells, cell)
		}
	}
	control, err := runPilotRankCell(
		"control_rank4_noise_0.05",
		"non_unitary_matched",
		4,
		0.05,
		trainTables,
		heldTables,
		trainDepths,
		heldDepths,
		block,
		applyStressNonUnitary,
	)
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	cells = append(cells, control)

	r2zero, err := pilotRankCellByName(cells, "unitary_rank2_noise_0.00")
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	r2noise, err := pilotRankCellByName(cells, "unitary_rank2_noise_0.05")
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	r3zero, err := pilotRankCellByName(cells, "unitary_rank3_noise_0.00")
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	r3noise, err := pilotRankCellByName(cells, "unitary_rank3_noise_0.05")
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	r4zero, err := pilotRankCellByName(cells, "unitary_rank4_noise_0.00")
	if err != nil {
		return PilotRankProbeResult{}, err
	}
	r4noise, err := pilotRankCellByName(cells, "unitary_rank4_noise_0.05")
	if err != nil {
		return PilotRankProbeResult{}, err
	}

	minNoiseless := minimalPassingPilotRank(cells, 0)
	minNoisy := minimalPassingPilotRank(cells, 0.05)
	diagnosis := PilotRankDiagnosis{
		MinimalNoiselessRank:        minNoiseless,
		MinimalNoise005Rank:         minNoisy,
		Rank2NoiselessAccuracy:      r2zero.HeldOutAccuracy,
		Rank2Noise005Accuracy:       r2noise.HeldOutAccuracy,
		Rank3NoiselessAccuracy:      r3zero.HeldOutAccuracy,
		Rank3Noise005Accuracy:       r3noise.HeldOutAccuracy,
		Rank4NoiselessAccuracy:      r4zero.HeldOutAccuracy,
		Rank4Noise005Accuracy:       r4noise.HeldOutAccuracy,
		IntrinsicRankLossSupported:  minNoiseless >= 3,
		MemoryNoiseSensitivity:      r2zero.HeldOutAccuracy-r2noise.HeldOutAccuracy >= 0.10 ||
			r3zero.HeldOutAccuracy-r3noise.HeldOutAccuracy >= 0.10,
		FourPilotCandidateSupported: r3noise.HeldOutAccuracy >= 0.95,
	}

	return PilotRankProbeResult{
		Schema:                  PilotRankProbeSchema,
		Experiment:              "UP-11-pilot-rank-memory-noise-sweep",
		Dimension:               16,
		Entities:                4,
		ValuesPerEntity:         4,
		RuntimePrototypeLookup:  false,
		ExplicitInverseReadout:  false,
		ExplicitDepthProvided:   false,
		GlobalPhaseNuisance:     true,
		SelectedTrainTables:     len(trainTables),
		SelectedHeldOutTables:   len(heldTables),
		MinTrainMarginalCount:   minTrain,
		MinHeldOutMarginalCount: minHeld,
		TrainDepths:             append([]int(nil), trainDepths...),
		HeldOutDepths:           append([]int(nil), heldDepths...),
		Cells:                   cells,
		Diagnosis:               diagnosis,
	}, nil
}
