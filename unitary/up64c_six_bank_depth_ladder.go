package unitary

import "math"

const UP64CSixBankDepthLadderSchema = "wingless.up64c-six-bank-depth-ladder.v1"

type UP64CDepthMetric struct {
	Family                string  `json:"family"`
	ScheduleBase          int     `json:"schedule_base"`
	MemoryNoise           float64 `json:"memory_noise"`
	Depth                 int     `json:"depth"`
	ValueAccuracy         float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin         float64 `json:"minimum_margin"`
	Gate                  bool    `json:"gate"`
}

type UP64CSixBankDepthLadderResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP62CSeal    string             `json:"source_up62c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	ScheduleBases      []int              `json:"schedule_bases"`
	NoiseLevels        []float64          `json:"noise_levels"`
	DepthLevels        []int              `json:"depth_levels"`
	Families           []string           `json:"families"`
	SelectionPerformed bool               `json:"selection_performed"`
	Metrics            []UP64CDepthMetric `json:"metrics"`
}

func up64cEvaluate(tags []float64, family string, noise float64, seedBase, depth int) (UP64CDepthMetric, error) {
	const banks = 6
	const scenarios = 256
	prototypes, err := up53cPrototypeBank(banks, tags)
	if err != nil {
		return UP64CDepthMetric{}, err
	}
	for i := range prototypes {
		prototypes[i].state, err = up53cTransportLocalPrototype(prototypes[i].state, depth)
		if err != nil {
			return UP64CDepthMetric{}, err
		}
	}
	block := up53cTransportBlock()
	var correct, total, exact int
	minMargin := math.Inf(1)
	for scenario := 0; scenario < scenarios; scenario++ {
		state, tables, err := up53cScenarioState(scenario, banks, tags)
		if err != nil {
			return UP64CDepthMetric{}, err
		}
		n := memoryNoise(seedBase+int(math.Round(noise*100000))*1000+scenario, 16, noise)
		for i := range state {
			state[i] += n[i]
		}
		state, err = Normalize(state)
		if err != nil {
			return UP64CDepthMetric{}, err
		}
		state = rotateGlobalPhase(state, math.Mod(0.317*float64(scenario+banks+1), 2*math.Pi))
		for d := 0; d < depth; d++ {
			state, err = Propagate(state, block)
			if err != nil {
				return UP64CDepthMetric{}, err
			}
		}
		scenarioExact := true
		for entity := 0; entity < 4; entity++ {
			local := append(State(nil), state[entity*4:(entity+1)*4]...)
			values, margin, err := up53cDecode(local, prototypes, true)
			if err != nil {
				return UP64CDepthMetric{}, err
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
	m := UP64CDepthMetric{
		Family: family, ScheduleBase: seedBase, MemoryNoise: noise, Depth: depth,
		ValueAccuracy: float64(correct) / float64(total),
		ExactScenarioAccuracy: float64(exact) / float64(scenarios),
		MinimumMargin: minMargin,
	}
	m.Gate = m.ValueAccuracy >= 0.99 && m.ExactScenarioAccuracy >= 0.95
	return m, nil
}

func RunUP64C() (UP64CSixBankDepthLadderResult, error) {
	schedules := []int{91000000, 92000000}
	noises := []float64{0.0075, 0.0100}
	depths := []int{16, 32, 64, 128, 256}
	all := up58cFamilies()
	var families []UP58CTagFamily
	for _, f := range all {
		if f.Name == "golden_rotation" || f.Name == "fixed_irregular" {
			families = append(families, f)
		}
	}
	result := UP64CSixBankDepthLadderResult{
		Schema: UP64CSixBankDepthLadderSchema,
		Experiment: "UP-64C-six-bank-depth-ladder",
		SourceUP62CSeal: "e2da0e6e6ee17b2eae765e07105ea2049e2a99e5",
		Banks: 6,
		StateDimension: 16,
		ScheduleBases: append([]int(nil), schedules...),
		NoiseLevels: append([]float64(nil), noises...),
		DepthLevels: append([]int(nil), depths...),
		SelectionPerformed: false,
	}
	for _, f := range families {
		result.Families = append(result.Families, f.Name)
	}
	for _, seedBase := range schedules {
		for _, f := range families {
			for _, noise := range noises {
				for _, depth := range depths {
					m, err := up64cEvaluate(f.Tags, f.Name, noise, seedBase, depth)
					if err != nil {
						return UP64CSixBankDepthLadderResult{}, err
					}
					result.Metrics = append(result.Metrics, m)
				}
			}
		}
	}
	return result, nil
}
