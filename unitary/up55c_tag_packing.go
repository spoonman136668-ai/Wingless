package unitary

import "math"

const UP55CTagPackingSchema="wingless.up55c-phase-tag-packing.v1"

type UP55CTagFamily struct{
	Name string `json:"name"`
	Tags []float64 `json:"tags"`
}
type UP55CTagMetric struct{
	Family string `json:"family"`
	Noise float64 `json:"noise"`
	ValueAccuracy float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin float64 `json:"minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	Gate bool `json:"gate"`
}
type UP55CTagPackingResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP54CSeal string `json:"source_up54c_seal"`
	Banks int `json:"banks"`
	StateDimension int `json:"state_dimension"`
	NoiseLevels []float64 `json:"noise_levels"`
	Families []UP55CTagFamily `json:"families"`
	Metrics []UP55CTagMetric `json:"metrics"`
	SelectionPerformed bool `json:"selection_performed"`
}

func up55cFamilies()[]UP55CTagFamily{
	return []UP55CTagFamily{
		{Name:"quadratic_current",Tags:[]float64{0,0.583,1.512,2.787}},
		{Name:"golden_rotation",Tags:[]float64{0,2.399963229728653,4.799926459457306,7.199889689185959}},
		{Name:"irrational_spread_a",Tags:[]float64{0,0.79,2.11,4.43}},
		{Name:"irrational_spread_b",Tags:[]float64{0,0.91,2.37,5.08}},
	}
}

func up55cEvaluate(tags []float64,noise float64)(UP55CTagMetric,error){
	const(banks=4;scenarios=256;depth=64)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP55CTagMetric{},err}
	for i:=range prototypes{prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth);if err!=nil{return UP55CTagMetric{},err}}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP55CTagMetric{},err}
		n:=memoryNoise(55000000+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP55CTagMetric{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP55CTagMetric{},err}}
		exact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP55CTagMetric{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP55CTagMetric{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{totalValues++;want:=tables[bank][entity];if cv[bank]==want{coherentCorrect++}else{exact=false};if mv[bank]==want{magnitudeCorrect++}}
		}
		if exact{coherentExact++}
	}
	m:=UP55CTagMetric{Noise:noise,ValueAccuracy:float64(coherentCorrect)/float64(totalValues),ExactScenarioAccuracy:float64(coherentExact)/scenarios,MinimumMargin:minMargin,MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues)}
	m.Gate=m.ValueAccuracy>=0.99&&m.ExactScenarioAccuracy>=0.95
	return m,nil
}

func RunUP55C()(UP55CTagPackingResult,error){
	noises:=[]float64{0.04,0.05};families:=up55cFamilies()
	result:=UP55CTagPackingResult{Schema:UP55CTagPackingSchema,Experiment:"UP-55C-phase-tag-packing",SourceUP54CSeal:"84ea3654f26c21dfbbda292bc7a9f67ad3d16f69",Banks:4,StateDimension:16,NoiseLevels:append([]float64(nil),noises...),Families:families,SelectionPerformed:false}
	for _,f:=range families{for _,noise:=range noises{m,err:=up55cEvaluate(f.Tags,noise);if err!=nil{return UP55CTagPackingResult{},err};m.Family=f.Name;result.Metrics=append(result.Metrics,m)}}
	return result,nil
}
