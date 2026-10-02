package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const MemoryProbeSchema = "wingless.unitary-memory-probe.v1"

type MemoryPathResult struct {
	Name                      string  `json:"name"`
	CommitDecodeAccuracy      float64 `json:"commit_decode_accuracy"`
	ExactFinalTableAccuracy   float64 `json:"exact_final_table_accuracy"`
	RelationalQueryAccuracy   float64 `json:"relational_query_accuracy"`
	MinDecodeMargin           float64 `json:"min_decode_margin"`
	MaxNormDrift              float64 `json:"max_norm_drift"`
}

type MemoryProbeResult struct {
	Schema                 string           `json:"schema"`
	Experiment             string           `json:"experiment"`
	Entities               int              `json:"entities"`
	ValuesPerEntity        int              `json:"values_per_entity"`
	Dimension              int              `json:"dimension"`
	Scenarios              int              `json:"scenarios"`
	WritesPerScenario      int              `json:"writes_per_scenario"`
	CommitDecodesPerPath   int              `json:"commit_decodes_per_path"`
	TransportDepths        []int            `json:"transport_depths"`
	NoiseAmplitude         float64          `json:"noise_amplitude"`
	WriteBoundary          string           `json:"write_boundary"`
	Unitary                MemoryPathResult `json:"unitary"`
	NonUnitary             MemoryPathResult `json:"non_unitary_matched"`
}

type memoryTable [4]int

type memoryPrototype struct {
	table memoryTable
	state State
}

type memoryScenario struct {
	initial  memoryTable
	writes   []memoryWrite
	queryA   int
	queryB   int
	finalGap int
}

type memoryWrite struct {
	entity int
	value  int
	gap    int
}

func validateMemoryTable(table memoryTable) error {
	for entity, value := range table {
		if value < 0 || value >= 4 {
			return fmt.Errorf("memory value out of range: entity=%d value=%d", entity, value)
		}
	}
	return nil
}

func encodeMemory(table memoryTable) (State, error) {
	if err := validateMemoryTable(table); err != nil {
		return nil, err
	}
	state := make(State, 16)
	for entity, value := range table {
		state[entity*4+value] = 0.5
	}
	return state, nil
}

func allMemoryTables() []memoryTable {
	out := make([]memoryTable, 0, 256)
	for a := 0; a < 4; a++ {
		for b := 0; b < 4; b++ {
			for c := 0; c < 4; c++ {
				for d := 0; d < 4; d++ {
					out = append(out, memoryTable{a, b, c, d})
				}
			}
		}
	}
	return out
}

func memoryNoise(seed, dimension int, amplitude float64) State {
	out := make(State, dimension)
	for i := range out {
		magnitude := amplitude * (1 + 0.25*float64((i+seed)%3))
		phase := 0.41*float64(i) + 0.17*float64(seed)
		out[i] = cmplx.Rect(magnitude, phase)
	}
	return out
}

func perturbMemory(state State, seed int, amplitude float64) (State, error) {
	noise := memoryNoise(seed, len(state), amplitude)
	out := append(State(nil), state...)
	for i := range out {
		out[i] += noise[i]
	}
	return Normalize(out)
}

func buildMemoryPrototypes(block []Coupling, depth int, apply stressApply) ([]memoryPrototype, error) {
	tables := allMemoryTables()
	out := make([]memoryPrototype, 0, len(tables))
	for _, table := range tables {
		canonical, err := encodeMemory(table)
		if err != nil {
			return nil, err
		}
		evolved, err := apply(canonical, block, depth)
		if err != nil {
			return nil, err
		}
		out = append(out, memoryPrototype{table: table, state: evolved})
	}
	return out, nil
}

func decodeMemory(state State, prototypes []memoryPrototype) (memoryTable, float64, error) {
	if len(prototypes) == 0 {
		return memoryTable{}, 0, fmt.Errorf("memory prototype set is empty")
	}
	bestIndex := 0
	bestScore := -1.0
	secondScore := -1.0
	for i, prototype := range prototypes {
		score, err := stressFidelity(prototype.state, state)
		if err != nil {
			return memoryTable{}, 0, err
		}
		if score > bestScore {
			secondScore = bestScore
			bestScore = score
			bestIndex = i
		} else if score > secondScore {
			secondScore = score
		}
	}
	margin := bestScore - secondScore
	if !finite(margin) || margin < 0 {
		return memoryTable{}, 0, fmt.Errorf("memory decode margin is invalid")
	}
	return prototypes[bestIndex].table, margin, nil
}

func applyMemoryWrite(table memoryTable, entity, value int) (memoryTable, error) {
	if err := validateMemoryTable(table); err != nil {
		return memoryTable{}, err
	}
	if entity < 0 || entity >= 4 || value < 0 || value >= 4 {
		return memoryTable{}, fmt.Errorf("memory write out of range")
	}
	out := table
	out[entity] = value
	return out, nil
}

func memoryRelation(table memoryTable, a, b int) (int, error) {
	if err := validateMemoryTable(table); err != nil {
		return 0, err
	}
	if a < 0 || a >= 4 || b < 0 || b >= 4 {
		return 0, fmt.Errorf("memory query entity out of range")
	}
	value := (table[a] - table[b]) % 4
	if value < 0 {
		value += 4
	}
	return value, nil
}

func makeMemoryScenario(seed int) (memoryScenario, error) {
	initial := memoryTable{}
	for entity := 0; entity < 4; entity++ {
		initial[entity] = (seed + entity) % 4
	}

	gaps := []int{32, 128, 512, 8}
	writes := make([]memoryWrite, 0, 12)
	for j := 0; j < 12; j++ {
		entity := (seed*3 + j*2 + j/3) % 4
		value := (seed + j*j + 2*j + entity) % 4
		gap := gaps[(seed+j)%len(gaps)]
		writes = append(writes, memoryWrite{entity: entity, value: value, gap: gap})
	}

	scenario := memoryScenario{
		initial:  initial,
		writes:   writes,
		queryA:   (seed + 1) % 4,
		queryB:   (seed + 3) % 4,
		finalGap: gaps[(seed+12)%len(gaps)],
	}
	return scenario, nil
}

func precomputeMemoryPrototypeSets(block []Coupling, depths []int, apply stressApply) (map[int][]memoryPrototype, error) {
	out := make(map[int][]memoryPrototype, len(depths))
	for _, depth := range depths {
		prototypes, err := buildMemoryPrototypes(block, depth, apply)
		if err != nil {
			return nil, err
		}
		out[depth] = prototypes
	}
	return out, nil
}

func runMemoryPath(name string, scenarios []memoryScenario, prototypeSets map[int][]memoryPrototype, block []Coupling, apply stressApply, noiseAmplitude float64) (MemoryPathResult, error) {
	var commitCorrect, commitTotal int
	var finalExact, relationCorrect int
	minMargin := math.Inf(1)
	var maxNormDrift float64

	for scenarioIndex, scenario := range scenarios {
		pathTable := scenario.initial
		trueTable := scenario.initial

		for writeIndex, write := range scenario.writes {
			canonical, err := encodeMemory(pathTable)
			if err != nil {
				return MemoryPathResult{}, err
			}
			noisy, err := perturbMemory(canonical, scenarioIndex*100+writeIndex, noiseAmplitude)
			if err != nil {
				return MemoryPathResult{}, err
			}
			evolved, err := apply(noisy, block, write.gap)
			if err != nil {
				return MemoryPathResult{}, err
			}
			norm2, err := NormSquared(evolved)
			if err != nil {
				return MemoryPathResult{}, err
			}
			drift := math.Abs(norm2 - 1)
			if drift > maxNormDrift {
				maxNormDrift = drift
			}

			prototypes, ok := prototypeSets[write.gap]
			if !ok {
				return MemoryPathResult{}, fmt.Errorf("missing memory prototypes for depth %d", write.gap)
			}
			decoded, margin, err := decodeMemory(evolved, prototypes)
			if err != nil {
				return MemoryPathResult{}, err
			}
			if margin < minMargin {
				minMargin = margin
			}
			commitTotal++
			if decoded == trueTable {
				commitCorrect++
			}

			// Explicit irreversible commit boundary: observe the current table,
			// overwrite exactly one entity, then re-encode a canonical state.
			// This boundary is identical for unitary and control paths.
			pathTable, err = applyMemoryWrite(decoded, write.entity, write.value)
			if err != nil {
				return MemoryPathResult{}, err
			}
			trueTable, err = applyMemoryWrite(trueTable, write.entity, write.value)
			if err != nil {
				return MemoryPathResult{}, err
			}
		}

		canonical, err := encodeMemory(pathTable)
		if err != nil {
			return MemoryPathResult{}, err
		}
		noisy, err := perturbMemory(canonical, scenarioIndex*100+99, noiseAmplitude)
		if err != nil {
			return MemoryPathResult{}, err
		}
		evolved, err := apply(noisy, block, scenario.finalGap)
		if err != nil {
			return MemoryPathResult{}, err
		}
		norm2, err := NormSquared(evolved)
		if err != nil {
			return MemoryPathResult{}, err
		}
		drift := math.Abs(norm2 - 1)
		if drift > maxNormDrift {
			maxNormDrift = drift
		}

		prototypes, ok := prototypeSets[scenario.finalGap]
		if !ok {
			return MemoryPathResult{}, fmt.Errorf("missing final memory prototypes for depth %d", scenario.finalGap)
		}
		decoded, margin, err := decodeMemory(evolved, prototypes)
		if err != nil {
			return MemoryPathResult{}, err
		}
		if margin < minMargin {
			minMargin = margin
		}

		if decoded == trueTable {
			finalExact++
		}

		gotRelation, err := memoryRelation(decoded, scenario.queryA, scenario.queryB)
		if err != nil {
			return MemoryPathResult{}, err
		}
		wantRelation, err := memoryRelation(trueTable, scenario.queryA, scenario.queryB)
		if err != nil {
			return MemoryPathResult{}, err
		}
		if gotRelation == wantRelation {
			relationCorrect++
		}
	}

	if commitTotal == 0 || len(scenarios) == 0 || !finite(minMargin) {
		return MemoryPathResult{}, fmt.Errorf("memory path produced no evaluable results")
	}
	return MemoryPathResult{
		Name:                    name,
		CommitDecodeAccuracy:    float64(commitCorrect) / float64(commitTotal),
		ExactFinalTableAccuracy: float64(finalExact) / float64(len(scenarios)),
		RelationalQueryAccuracy: float64(relationCorrect) / float64(len(scenarios)),
		MinDecodeMargin:         minMargin,
		MaxNormDrift:            maxNormDrift,
	}, nil
}

// RunUP5 introduces an explicit irreversible write boundary because true
// overwrite is many-to-one and cannot be represented by a closed unitary map
// without preserving the displaced information elsewhere.
//
// Between write/commit boundaries, the memory is transported through the same
// deep 16-dimensional unitary vs matched non-unitary substrate used by UP-3.
// The experiment asks whether a mutable four-entity table can survive repeated
// overwrite, distractor transport and a final relational query.
func RunUP5() (MemoryProbeResult, error) {
	const (
		scenarioCount   = 32
		writesPerCase   = 12
		noiseAmplitude  = 0.10
	)
	depths := []int{8, 32, 128, 512}
	block := stressProgram()

	scenarios := make([]memoryScenario, 0, scenarioCount)
	for seed := 0; seed < scenarioCount; seed++ {
		scenario, err := makeMemoryScenario(seed)
		if err != nil {
			return MemoryProbeResult{}, err
		}
		if len(scenario.writes) != writesPerCase {
			return MemoryProbeResult{}, fmt.Errorf("unexpected memory write count")
		}
		scenarios = append(scenarios, scenario)
	}

	unitaryPrototypes, err := precomputeMemoryPrototypeSets(block, depths, applyStressUnitary)
	if err != nil {
		return MemoryProbeResult{}, err
	}
	controlPrototypes, err := precomputeMemoryPrototypeSets(block, depths, applyStressNonUnitary)
	if err != nil {
		return MemoryProbeResult{}, err
	}

	unitaryResult, err := runMemoryPath(
		"unitary_transport",
		scenarios,
		unitaryPrototypes,
		block,
		applyStressUnitary,
		noiseAmplitude,
	)
	if err != nil {
		return MemoryProbeResult{}, err
	}
	controlResult, err := runMemoryPath(
		"non_unitary_matched_transport",
		scenarios,
		controlPrototypes,
		block,
		applyStressNonUnitary,
		noiseAmplitude,
	)
	if err != nil {
		return MemoryProbeResult{}, err
	}

	return MemoryProbeResult{
		Schema:               MemoryProbeSchema,
		Experiment:           "UP-5-relational-read-write-memory",
		Entities:             4,
		ValuesPerEntity:      4,
		Dimension:            16,
		Scenarios:            scenarioCount,
		WritesPerScenario:    writesPerCase,
		CommitDecodesPerPath: scenarioCount * writesPerCase,
		TransportDepths:      append([]int(nil), depths...),
		NoiseAmplitude:       noiseAmplitude,
		WriteBoundary:        "observe_decode_overwrite_reencode",
		Unitary:              unitaryResult,
		NonUnitary:           controlResult,
	}, nil
}
