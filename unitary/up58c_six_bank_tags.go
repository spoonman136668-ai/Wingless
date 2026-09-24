package unitary

import "math"

const UP58CSixBankTagSchema="wingless.up58c-six-bank-tag-ablation.v1"

type UP58CTagFamily struct {
	Name string `json:"name"`
	Tags []float64 `json:"tags"`
}
type UP58CTagMetric struct {
	Family string `json:"family"`
	MemoryNoise float64 `json:"memory_noise"`
	ValueAccuracy float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin float64 `json:"minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	Gate bool `json:"gate"`
}
type UP58CSixBankTagResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP57CSeal string `json:"source_up57c_seal"`
	Banks int `json:"banks"`
	StateDimension int `json:"state_dimension"`
	NoiseLevels []float64 `json:"noise_levels"`
	Families []UP58CTagFamily `json:"families"`
	Metrics []UP58CTagMetric `json:"metrics"`
	SelectionPerformed bool `json:"selection_performed"`
}

func up58cFamilies()[]UP58CTagFamily{
	const golden=2.399963229728653
	goldenTags:=make([]float64,6)
	uniform:=make([]float64,6)
	quadratic:=make([]float64,6)
	for i:=0;i<6;i++{
		goldenTags[i]=golden*float64(i)
		uniform[i]=2*math.Pi*float64(i)/6
		b:=float64(i);quadratic[i]=0.41*b+0.173*b*b
	}
	return []UP58CTagFamily{
		{Name:"golden_rotation",Tags:goldenTags},
		{Name:"uniform_hexagon",Tags:uniform},
		{Name:"quadratic_extension",Tags:quadratic},
		{Name:"fixed_irregular",Tags:[]float64{0,0.73,1.91,3.37,5.11,7.03}},
	}
}

func up58cEvaluate(tags []float64,noise float64)(UP58CTagMetric,error){
	const(banks=6;scenarios=256;depth=64)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP58CTagMetric{},err}
	for i:=range prototypes{prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth);if err!=nil{return UP58CTagMetric{},err}}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP58CTagMetric{},err}
		n:=memoryNoise(71000000+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP58CTagMetric{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP58CTagMetric{},err}}
		exact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP58CTagMetric{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP58CTagMetric{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{
				totalValues++;want:=tables[bank][entity]
				if cv[bank]==want{coherentCorrect++}else{exact=false}
				if mv[bank]==want{magnitudeCorrect++}
			}
		}
		if exact{coherentExact++}
	}
	m:=UP58CTagMetric{
		MemoryNoise:noise,
		ValueAccuracy:float64(coherentCorrect)/float64(totalValues),
		ExactScenarioAccuracy:float64(coherentExact)/scenarios,
		MinimumMargin:minMargin,
		MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues),
	}
	m.Gate=m.ValueAccuracy>=0.99&&m.ExactScenarioAccuracy>=0.95
	return m,nil
}

func RunUP58C()(UP58CSixBankTagResult,error){
	noises:=[]float64{0,0.0025}
	families:=up58cFamilies()
	result:=UP58CSixBankTagResult{
		Schema:UP58CSixBankTagSchema,Experiment:"UP-58C-six-bank-tag-ablation",
		SourceUP57CSeal:"010642f03fd990bffcc276ca116fb971eed8b667",
		Banks:6,StateDimension:16,NoiseLevels:append([]float64(nil),noises...),
		Families:families,SelectionPerformed:false,
	}
	for _,f:=range families{
		for _,noise:=range noises{
			m,err:=up58cEvaluate(f.Tags,noise);if err!=nil{return UP58CSixBankTagResult{},err}
			m.Family=f.Name
			result.Metrics=append(result.Metrics,m)
		}
	}
	return result,nil
}
