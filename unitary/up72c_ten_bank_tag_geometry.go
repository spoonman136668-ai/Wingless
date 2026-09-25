package unitary

import "math"

const UP72CTenBankTagGeometrySchema = "wingless.up72c-ten-bank-tag-geometry.v1"

type UP72CTagGeometryPoint struct {
	TagFamily              string  `json:"tag_family"`
	ScheduleBase           int     `json:"schedule_base"`
	Banks                  int     `json:"banks"`
	MemoryNoise            float64 `json:"memory_noise"`
	ValueAccuracy          float64 `json:"value_accuracy"`
	ExactScenarioAccuracy  float64 `json:"exact_scenario_accuracy"`
	MinimumMargin          float64 `json:"minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	Gate                   bool    `json:"gate"`
}

type UP72CTenBankTagGeometryResult struct {
	Schema             string                  `json:"schema"`
	Experiment         string                  `json:"experiment"`
	SourceUP70CSeal    string                  `json:"source_up70c_seal"`
	StateDimension     int                     `json:"state_dimension"`
	Banks              int                     `json:"banks"`
	ScheduleBases      []int                   `json:"schedule_bases"`
	NoiseLevels        []float64               `json:"noise_levels"`
	TagFamilies        []string                `json:"tag_families"`
	SelectionPerformed bool                    `json:"selection_performed"`
	Points             []UP72CTagGeometryPoint `json:"points"`
}

func up72cUniformTags(banks int) []float64 {
	tags:=make([]float64,banks)
	for i:=0;i<banks;i++ { tags[i]=2*math.Pi*float64(i)/float64(banks) }
	return tags
}

func up72cEvaluate(family string,tags []float64,noise float64,seedBase int)(UP72CTagGeometryPoint,error){
	const(banks=10;scenarios=256;depth=64)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP72CTagGeometryPoint{},err}
	for i:=range prototypes{prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth);if err!=nil{return UP72CTagGeometryPoint{},err}}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP72CTagGeometryPoint{},err}
		n:=memoryNoise(seedBase+banks*100000+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP72CTagGeometryPoint{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP72CTagGeometryPoint{},err}}
		exact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP72CTagGeometryPoint{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP72CTagGeometryPoint{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{totalValues++;want:=tables[bank][entity];if cv[bank]==want{coherentCorrect++}else{exact=false};if mv[bank]==want{magnitudeCorrect++}}
		}
		if exact{coherentExact++}
	}
	p:=UP72CTagGeometryPoint{TagFamily:family,ScheduleBase:seedBase,Banks:banks,MemoryNoise:noise,ValueAccuracy:float64(coherentCorrect)/float64(totalValues),ExactScenarioAccuracy:float64(coherentExact)/scenarios,MinimumMargin:minMargin,MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues)}
	p.Gate=p.ValueAccuracy>=0.99&&p.ExactScenarioAccuracy>=0.95
	return p,nil
}

func RunUP72C()(UP72CTenBankTagGeometryResult,error){
	schedules:=[]int{115000000,116000000};noises:=[]float64{0,0.001,0.002,0.004};families:=[]string{"golden_rotation","uniform_phase"}
	result:=UP72CTenBankTagGeometryResult{Schema:UP72CTenBankTagGeometrySchema,Experiment:"UP-72C-ten-bank-tag-geometry",SourceUP70CSeal:"52af265b43b4074e62bdc3836ab68c7d90e58452",StateDimension:16,Banks:10,ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),TagFamilies:append([]string(nil),families...),SelectionPerformed:false}
	for _,seedBase:=range schedules{for _,noise:=range noises{
		p,err:=up72cEvaluate("golden_rotation",up57cGoldenTags(10),noise,seedBase);if err!=nil{return UP72CTenBankTagGeometryResult{},err};result.Points=append(result.Points,p)
		p,err=up72cEvaluate("uniform_phase",up72cUniformTags(10),noise,seedBase);if err!=nil{return UP72CTenBankTagGeometryResult{},err};result.Points=append(result.Points,p)
	}}
	return result,nil
}
