package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const CoevolvingFrameSchema = "wingless.unitary-coevolving-frame.v1"

type FrameObserverCell struct {
	Name               string    `json:"name"`
	Path               string    `json:"path"`
	FrameMode          string    `json:"frame_mode"`
	TrainDepths        []int     `json:"train_depths"`
	HeldOutDepths      []int     `json:"held_out_depths"`
	TrainAccuracy      float64   `json:"train_accuracy"`
	HeldOutAccuracy    float64   `json:"held_out_accuracy"`
	PerEntityTrain     []float64 `json:"per_entity_train_accuracy"`
	PerEntityHeldOut   []float64 `json:"per_entity_held_out_accuracy"`
}

type CoevolvingFrameDiagnosis struct {
	BalancedSplitValid        bool    `json:"balanced_split_valid"`
	UnitaryCoevolvingPass     bool    `json:"unitary_coevolving_unseen_depth_pass"`
	UnitaryStaticPass         bool    `json:"unitary_static_unseen_depth_pass"`
	ControlCoevolvingPass     bool    `json:"control_coevolving_unseen_depth_pass"`
	InternalFrameSupported    bool    `json:"internal_frame_supported"`
	UnitaryFrameGain          float64 `json:"unitary_frame_gain"`
	ControlFrameGain          float64 `json:"control_frame_gain"`
}

type CoevolvingFrameProbeResult struct {
	Schema                  string                    `json:"schema"`
	Experiment              string                    `json:"experiment"`
	Dimension               int                       `json:"dimension"`
	Entities                int                       `json:"entities"`
	ValuesPerEntity         int                       `json:"values_per_entity"`
	PilotStates             int                       `json:"pilot_states"`
	PilotComplexScalars     int                       `json:"pilot_complex_scalars"`
	FeaturesPerEntity       int                       `json:"features_per_entity"`
	RuntimePrototypeLookup  bool                      `json:"runtime_prototype_lookup"`
	ExplicitInverseReadout  bool                      `json:"explicit_inverse_readout"`
	ExplicitDepthProvided   bool                      `json:"explicit_depth_provided"`
	GlobalPhaseNuisance     bool                      `json:"global_phase_nuisance"`
	ReferenceNoiseAmplitude float64                   `json:"reference_noise_amplitude"`
	MemoryNoiseAmplitude    float64                   `json:"memory_noise_amplitude"`
	SelectedTrainTables     int                       `json:"selected_train_tables"`
	SelectedHeldOutTables   int                       `json:"selected_held_out_tables"`
	MinTrainMarginalCount   int                       `json:"min_train_marginal_count"`
	MinHeldOutMarginalCount int                       `json:"min_held_out_marginal_count"`
	Cells                   []FrameObserverCell       `json:"cells"`
	Diagnosis               CoevolvingFrameDiagnosis `json:"diagnosis"`
}

type frameBank struct {
	anchor State
	phase  [4]State
}

func makeFrameBank() (frameBank, error) {
	var bank frameBank
	bank.anchor = make(State, 16)
	for i := range bank.anchor {
		bank.anchor[i] = 0.25
	}
	var err error
	bank.anchor, err = Normalize(bank.anchor)
	if err != nil {
		return frameBank{}, err
	}

	for entity := 0; entity < 4; entity++ {
		pilot := make(State, 16)
		for value := 0; value < 4; value++ {
			pilot[entity*4+value] = cmplx.Rect(0.5, math.Pi*float64(value)/2)
		}
		pilot, err = Normalize(pilot)
		if err != nil {
			return frameBank{}, err
		}
		bank.phase[entity] = pilot
	}
	return bank, nil
}

func stateInner(a, b State) (complex128, error) {
	if err := validateState(a); err != nil {
		return 0, err
	}
	if err := validateState(b); err != nil {
		return 0, err
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("inner-product dimension mismatch")
	}
	var out complex128
	for i := range a {
		out += cmplx.Conj(a[i]) * b[i]
	}
	if !finite(real(out)) || !finite(imag(out)) {
		return 0, fmt.Errorf("inner product is not finite")
	}
	return out, nil
}

func frameFeature(memory, anchor, phase State) ([]float64, error) {
	cPhase, err := stateInner(phase, memory)
	if err != nil {
		return nil, err
	}
	cAnchor, err := stateInner(anchor, memory)
	if err != nil {
		return nil, err
	}
	denominator := cmplx.Abs(cPhase) * cmplx.Abs(cAnchor)
	if !finite(denominator) || denominator <= 1e-15 {
		return nil, fmt.Errorf("frame correlation denominator is invalid")
	}
	value := cPhase * cmplx.Conj(cAnchor) / complex(denominator, 0)
	if !finite(real(value)) || !finite(imag(value)) {
		return nil, fmt.Errorf("frame feature is not finite")
	}
	return []float64{real(value), imag(value)}, nil
}

func rotateGlobalPhase(state State, phase float64) State {
	out := append(State(nil), state...)
	g := cmplx.Rect(1, phase)
	for i := range out {
		out[i] *= g
	}
	return out
}

func transportFrameBank(bank frameBank, block []Coupling, depth int, apply stressApply, coevolve bool) (frameBank, error) {
	if !coevolve {
		return bank, nil
	}
	anchor, err := apply(bank.anchor, block, depth)
	if err != nil {
		return frameBank{}, err
	}
	var phases [4]State
	for entity := 0; entity < 4; entity++ {
		phases[entity], err = apply(bank.phase[entity], block, depth)
		if err != nil {
			return frameBank{}, err
		}
	}
	return frameBank{anchor: anchor, phase: phases}, nil
}

func buildFrameSamples(
	tables []memoryTable,
	depths []int,
	entity int,
	block []Coupling,
	apply stressApply,
	coev bool,
	memoryNoise float64,
	seedOffset int,
) ([]headSample, error) {
	bank, err := makeFrameBank()
	if err != nil {
		return nil, err
	}
	transformedBanks := make(map[int]frameBank, len(depths))
	for _, depth := range depths {
		transformed, err := transportFrameBank(bank, block, depth, apply, coev)
		if err != nil {
			return nil, err
		}
		transformedBanks[depth] = transformed
	}

	var samples []headSample
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		for depthIndex, depth := range depths {
			seed := seedOffset + memoryTableIndex(table)*10000 + depthIndex*131 + entity*17
			noisy, err := perturbMemory(canonical, seed, memoryNoise)
			if err != nil {
				return nil, err
			}
			// Deliberate nuisance: the memory receives a deterministic global
			// phase that is not applied to the frame bank. The anchor/phase
			// correlation product must cancel it.
			globalPhase := math.Mod(0.173*float64(seed+1), 2*math.Pi)
			noisy = rotateGlobalPhase(noisy, globalPhase)

			forward, err := apply(noisy, block, depth)
			if err != nil {
				return nil, err
			}
			frame := transformedBanks[depth]
			features, err := frameFeature(forward, frame.anchor, frame.phase[entity])
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

func runFrameCell(
	name, path, mode string,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths []int,
	block []Coupling,
	apply stressApply,
	coev bool,
	memoryNoise float64,
) (FrameObserverCell, error) {
	const (
		steps        = 200
		learningRate = 1.0
	)
	perTrain := make([]float64, 4)
	perHeld := make([]float64, 4)
	for entity := 0; entity < 4; entity++ {
		trainSamples, err := buildFrameSamples(
			trainTables, trainDepths, entity, block, apply, coev,
			memoryNoise, 0,
		)
		if err != nil {
			return FrameObserverCell{}, err
		}
		heldSamples, err := buildFrameSamples(
			heldTables, heldDepths, entity, block, apply, coev,
			memoryNoise, 800000,
		)
		if err != nil {
			return FrameObserverCell{}, err
		}
		head, _, err := trainLinearSoftmax(trainSamples, 4, 2, steps, learningRate)
		if err != nil {
			return FrameObserverCell{}, err
		}
		_, trainAccuracy, err := evaluateHead(head, trainSamples)
		if err != nil {
			return FrameObserverCell{}, err
		}
		_, heldAccuracy, err := evaluateHead(head, heldSamples)
		if err != nil {
			return FrameObserverCell{}, err
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

	return FrameObserverCell{
		Name:             name,
		Path:             path,
		FrameMode:        mode,
		TrainDepths:      append([]int(nil), trainDepths...),
		HeldOutDepths:    append([]int(nil), heldDepths...),
		TrainAccuracy:    trainAverage,
		HeldOutAccuracy:  heldAverage,
		PerEntityTrain:   perTrain,
		PerEntityHeldOut: perHeld,
	}, nil
}

func findFrameCell(cells []FrameObserverCell, name string) (FrameObserverCell, error) {
	for _, cell := range cells {
		if cell.Name == name {
			return cell, nil
		}
	}
	return FrameObserverCell{}, fmt.Errorf("frame cell %q missing", name)
}

// RunUP9 tests whether a compact co-evolving internal reference frame resolves
// the depth/frame ambiguity isolated by UP-8 without supplying depth or an
// explicit inverse.
//
// Five pilot states are used:
//   - one global anchor reference;
//   - one four-phase reference for each of four entities.
//
// Each entity observer receives only two real numbers derived from the
// normalized phase of:
//     <phase_ref|memory> * conj(<anchor_ref|memory>)
//
// For a unitary transport, co-evolving memory/reference inner products are
// invariant with depth. The product also cancels a memory-only global phase.
// A static-reference ablation tests whether co-evolution itself is necessary.
func RunUP9() (CoevolvingFrameProbeResult, error) {
	const (
		selectedTables = 64
		memoryNoise    = 0.05
	)
	fixedDepth := []int{128}
	trainDepths := []int{8, 24, 72, 216, 432, 648}
	heldDepths := []int{32, 128, 512, 1024}

	trainTables, err := selectObserverTables(true, selectedTables)
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}
	heldTables, err := selectObserverTables(false, selectedTables)
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}
	minTrain := observerMarginalMinimum(trainTables)
	minHeld := observerMarginalMinimum(heldTables)
	if minTrain <= 0 || minHeld <= 0 {
		return CoevolvingFrameProbeResult{}, fmt.Errorf("frame experiment lost a table marginal")
	}

	block := stressProgram()
	specs := []struct {
		name, path, mode string
		trainDepths      []int
		heldDepths       []int
		apply            stressApply
		coev             bool
	}{
		{"unitary_static_fixed", "unitary", "static", fixedDepth, fixedDepth, applyStressUnitary, false},
		{"unitary_coevolving_fixed", "unitary", "coevolving", fixedDepth, fixedDepth, applyStressUnitary, true},
		{"unitary_static_unseen_depth", "unitary", "static", trainDepths, heldDepths, applyStressUnitary, false},
		{"unitary_coevolving_unseen_depth", "unitary", "coevolving", trainDepths, heldDepths, applyStressUnitary, true},
		{"control_static_unseen_depth", "non_unitary_matched", "static", trainDepths, heldDepths, applyStressNonUnitary, false},
		{"control_coevolving_unseen_depth", "non_unitary_matched", "coevolving", trainDepths, heldDepths, applyStressNonUnitary, true},
	}

	cells := make([]FrameObserverCell, 0, len(specs))
	for _, spec := range specs {
		cell, err := runFrameCell(
			spec.name,
			spec.path,
			spec.mode,
			trainTables,
			heldTables,
			spec.trainDepths,
			spec.heldDepths,
			block,
			spec.apply,
			spec.coev,
			memoryNoise,
		)
		if err != nil {
			return CoevolvingFrameProbeResult{}, err
		}
		cells = append(cells, cell)
	}

	unitaryStatic, err := findFrameCell(cells, "unitary_static_unseen_depth")
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}
	unitaryCoev, err := findFrameCell(cells, "unitary_coevolving_unseen_depth")
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}
	controlStatic, err := findFrameCell(cells, "control_static_unseen_depth")
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}
	controlCoev, err := findFrameCell(cells, "control_coevolving_unseen_depth")
	if err != nil {
		return CoevolvingFrameProbeResult{}, err
	}

	unitaryGain := unitaryCoev.HeldOutAccuracy - unitaryStatic.HeldOutAccuracy
	controlGain := controlCoev.HeldOutAccuracy - controlStatic.HeldOutAccuracy

	diagnosis := CoevolvingFrameDiagnosis{
		BalancedSplitValid:     minTrain > 0 && minHeld > 0,
		UnitaryCoevolvingPass:  unitaryCoev.HeldOutAccuracy >= 0.95,
		UnitaryStaticPass:      unitaryStatic.HeldOutAccuracy >= 0.95,
		ControlCoevolvingPass:  controlCoev.HeldOutAccuracy >= 0.95,
		InternalFrameSupported: unitaryCoev.HeldOutAccuracy >= 0.95 && unitaryGain >= 0.15,
		UnitaryFrameGain:       unitaryGain,
		ControlFrameGain:       controlGain,
	}

	return CoevolvingFrameProbeResult{
		Schema:                  CoevolvingFrameSchema,
		Experiment:              "UP-9-coevolving-internal-frame",
		Dimension:               16,
		Entities:                4,
		ValuesPerEntity:         4,
		PilotStates:             5,
		PilotComplexScalars:     80,
		FeaturesPerEntity:       2,
		RuntimePrototypeLookup:  false,
		ExplicitInverseReadout:  false,
		ExplicitDepthProvided:   false,
		GlobalPhaseNuisance:     true,
		ReferenceNoiseAmplitude: 0,
		MemoryNoiseAmplitude:    memoryNoise,
		SelectedTrainTables:     len(trainTables),
		SelectedHeldOutTables:   len(heldTables),
		MinTrainMarginalCount:   minTrain,
		MinHeldOutMarginalCount: minHeld,
		Cells:                   cells,
		Diagnosis:               diagnosis,
	}, nil
}
