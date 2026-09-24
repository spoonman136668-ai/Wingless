package unitary

import "math"

const UP67CBoundaryDepthReplicationSchema = "wingless.up67c-boundary-depth-replication.v1"

type UP67CDepthMetric struct {
	Family                string  `json:"family"`
	ScheduleBase          int     `json:"schedule_base"`
	MemoryNoise           float64 `json:"memory_noise"`
	Depth                 int     `json:"depth"`
	ValueAccuracy         float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin         float64 `json:"minimum_margin"`
	Gate                  bool    `json:"gate"`
}

type UP67CBoundaryDepthReplicationResult struct {
	Schema               string             `json:"schema"`
	Experiment           string             `json:"experiment"`
	SourceUP65CSeal      string             `json:"source_up65c_seal"`
	Banks                int                `json:"banks"`
	StateDimension       int                `json:"state_dimension"`
	ScheduleBases        []int              `json:"schedule_bases"`
	GoldenNoiseLevels    []float64          `json:"golden_noise_levels"`
	IrregularNoiseLevels []float64          `json:"irregular_noise_levels"`
	DepthLevels          []int              `json:"depth_levels"`
	SelectionPerformed   bool               `json:"selection_performed"`
	Metrics              []UP67CDepthMetric `json:"metrics"`
}

func up67cEvaluate(tags []float64, family string, noise float64, seedBase, depth int) (UP67CDepthMetric, error) {
	const banks = 6
	const scenarios = 256
	prototypes, err := up53cPrototypeBank(banks, tags)
	if err != nil {
		return UP67CDepthMetric{}, err
	}
	for i := range prototypes {
		prototypes[i].state, err = up53cTransportLocalPrototype(prototypes[i].state, depth)
		if err != nil {
			return UP67CDepthMetric{}, err
		}
	}
	block := up53cTransportBlock()
	var correct, total, exact int
	minMargin := math.Inf(1)
	for scenario := 0; scenario < scenarios; scenario++ {
		state, tables, err := up53cScenarioState(scenario, banks, tags)
		if err != nil {
			return UP67CDepthMetric{}, err
		}
		n := memoryNoise(seedBase+int(math.Round(noise*100000))*1000+scenario, 16, noise)
		for i := range state {
			state[i] += n[i]
		}
		state, err = Normalize(state)
		if err != nil {
			return UP67CDepthMetric{}, err
		}
		state = rotateGlobalPhase(state, math.Mod(0.317*float64(scenario+banks+1), 2*math.Pi))
		for d := 0; d < depth; d++ {
			state, err = Propagate(state, block)
			if err != nil {
				return UP67CDepthMetric{}, err
			}
		}
		scenarioExact := true
		for entity := 0; entity < 4; entity++ {
			local := append(State(nil), state[entity*4:(entity+1)*4]...)
			values, margin, err := up53cDecode(local, prototypes, true)
			if err != nil {
				return UP67CDepthMetric{}, err
			}
			if margin < minMargin {
				minMargin = margin
			}
			for bank := 0; bank < banks; bank++ {
				total++
				if values[bank] == tables[bank][entity] {
					correct++
				} else {
					scenarioExact = false
				}
			}
		}
		if scenarioExact {
			exact++
		}
	}
	m := UP67CDepthMetric{
		Family: family,
		ScheduleBase: seedBase,
		MemoryNoise: noise,
		Depth: depth,
		ValueAccuracy: float64(correct) / float64(total),
		ExactScenarioAccuracy: float64(exact) / float64(scenarios),
		MinimumMargin: minMargin,
	}
	m.Gate = m.ValueAccuracy >= 0.99 && m.ExactScenarioAccuracy >= 0.95
	return m, nil
}

func RunUP67C() (UP67CBoundaryDepthReplicationResult, error) {
	schedules := []int{105000000, 106000000}
	goldenNoises := []float64{0.0106, 0.0108, 0.0110}
	irregularNoises := []float64{0.0085, 0.0086, 0.0087}
	depths := []int{32, 64, 128}

	var golden, irregular UP58CTagFamily
	for _, f := range up58cFamilies() {
		switch f.Name {
		case "golden_rotation":
			golden = f
		case "fixed_irregular":
			irregular = f
		}
	}

	result := UP67CBoundaryDepthReplicationResult{
		Schema: UP67CBoundaryDepthReplicationSchema,
		Experiment: "UP-67C-boundary-depth-replication",
		SourceUP65CSeal: "1895beda79f067ee1bf6646138445eed828fd3d8",
		Banks: 6,
		StateDimension: 16,
		ScheduleBases: append([]int(nil), schedules...),
		GoldenNoiseLevels: append([]float64(nil), goldenNoises...),
		IrregularNoiseLevels: append([]float64(nil), irregularNoises...),
		DepthLevels: append([]int(nil), depths...),
		SelectionPerformed: false,
	}

	for _, seedBase := range schedules {
		for _, depth := range depths {
			for _, noise := range goldenNoises {
				m, err := up67cEvaluate(golden.Tags, golden.Name, noise, seedBase, depth)
				if err != nil {
					return UP67CBoundaryDepthReplicationResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
			for _, noise := range irregularNoises {
				m, err := up67cEvaluate(irregular.Tags, irregular.Name, noise, seedBase, depth)
				if err != nil {
					return UP67CBoundaryDepthReplicationResult{}, err
				}
				result.Metrics = append(result.Metrics, m)
			}
		}
	}
	return result, nil
}
