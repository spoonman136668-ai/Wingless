package unitary

import "math"

const UP59CSixBankNoiseSchema = "wingless.up59c-six-bank-noise-ladder.v1"

type UP59CNoiseMetric struct {
	Family                string  `json:"family"`
	ScheduleBase          int     `json:"schedule_base"`
	MemoryNoise           float64 `json:"memory_noise"`
	ValueAccuracy         float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin         float64 `json:"minimum_margin"`
	Gate                  bool    `json:"gate"`
}

type UP59CSixBankNoiseResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP58CSeal    string             `json:"source_up58c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	ScheduleBases      []int              `json:"schedule_bases"`
	NoiseLevels        []float64          `json:"noise_levels"`
	Families           []string           `json:"families"`
	SelectionPerformed bool               `json:"selection_performed"`
	Metrics            []UP59CNoiseMetric `json:"metrics"`
}

func up59cEvaluate(tags []float64,family string,noise float64,seedBase int)(UP59CNoiseMetric,error){
	const(banks=6;scenarios=256;depth=64)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP59CNoiseMetric{},err}
	for i:=range prototypes{
		prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth)
		if err!=nil{return UP59CNoiseMetric{},err}
	}
	block:=up53cTransportBlock()
	var correct,total,exact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP59CNoiseMetric{},err}
		n:=memoryNoise(seedBase+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP59CNoiseMetric{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP59CNoiseMetric{},err}}
		scenarioExact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			values,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP59CNoiseMetric{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{
				total++
				if values[bank]==tables[bank][entity]{correct++}else{scenarioExact=false}
			}
		}
		if scenarioExact{exact++}
	}
	m:=UP59CNoiseMetric{
		Family:family,ScheduleBase:seedBase,MemoryNoise:noise,
		ValueAccuracy:float64(correct)/float64(total),
		ExactScenarioAccuracy:float64(exact)/float64(scenarios),
		MinimumMargin:minMargin,
	}
	m.Gate=m.ValueAccuracy>=0.99&&m.ExactScenarioAccuracy>=0.95
	return m,nil
}

func RunUP59C()(UP59CSixBankNoiseResult,error){
	schedules:=[]int{73000000,74000000}
	noises:=[]float64{0.0025,0.0035,0.0045,0.0055,0.0075}
	all:=up58cFamilies()
	var families []UP58CTagFamily
	for _,f:=range all{
		if f.Name=="uniform_hexagon"{continue}
		families=append(families,f)
	}
	result:=UP59CSixBankNoiseResult{
		Schema:UP59CSixBankNoiseSchema,
		Experiment:"UP-59C-six-bank-noise-ladder",
		SourceUP58CSeal:"98999e0eda92f6a17e5ad9ac9db6d4d5da533c58",
		Banks:6,StateDimension:16,
		ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),
		SelectionPerformed:false,
	}
	for _,f:=range families{result.Families=append(result.Families,f.Name)}
	for _,seedBase:=range schedules{
		for _,f:=range families{
			for _,noise:=range noises{
				m,err:=up59cEvaluate(f.Tags,f.Name,noise,seedBase)
				if err!=nil{return UP59CSixBankNoiseResult{},err}
				result.Metrics=append(result.Metrics,m)
			}
		}
	}
	return result,nil
}
