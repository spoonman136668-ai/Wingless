package unitary

import "math"

const UP57CGoldenScaleSchema = "wingless.up57c-golden-bank-scale.v1"

type UP57CGoldenScalePoint struct {
	ScheduleBase           int     `json:"schedule_base"`
	Banks                  int     `json:"banks"`
	MemoryNoise            float64 `json:"memory_noise"`
	ValueAccuracy          float64 `json:"value_accuracy"`
	ExactScenarioAccuracy  float64 `json:"exact_scenario_accuracy"`
	MinimumMargin          float64 `json:"minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	Gate                   bool    `json:"gate"`
}
type UP57CGoldenScaleResult struct {
	Schema             string                  `json:"schema"`
	Experiment         string                  `json:"experiment"`
	SourceUP56CSeal    string                  `json:"source_up56c_seal"`
	StateDimension     int                     `json:"state_dimension"`
	ScheduleBases      []int                   `json:"schedule_bases"`
	SelectionPerformed bool                    `json:"selection_performed"`
	Points             []UP57CGoldenScalePoint `json:"points"`
}

func up57cGoldenTags(banks int) []float64 {
	const golden=2.399963229728653
	tags:=make([]float64,banks)
	for i:=0;i<banks;i++{tags[i]=golden*float64(i)}
	return tags
}

func up57cEvaluate(banks int,noise float64,seedBase int)(UP57CGoldenScalePoint,error){
	const(scenarios=256;depth=64)
	tags:=up57cGoldenTags(banks)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP57CGoldenScalePoint{},err}
	for i:=range prototypes{prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth);if err!=nil{return UP57CGoldenScalePoint{},err}}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP57CGoldenScalePoint{},err}
		n:=memoryNoise(seedBase+banks*100000+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP57CGoldenScalePoint{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP57CGoldenScalePoint{},err}}
		exact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP57CGoldenScalePoint{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP57CGoldenScalePoint{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{
				totalValues++;want:=tables[bank][entity]
				if cv[bank]==want{coherentCorrect++}else{exact=false}
				if mv[bank]==want{magnitudeCorrect++}
			}
		}
		if exact{coherentExact++}
	}
	p:=UP57CGoldenScalePoint{
		ScheduleBase:seedBase,Banks:banks,MemoryNoise:noise,
		ValueAccuracy:float64(coherentCorrect)/float64(totalValues),
		ExactScenarioAccuracy:float64(coherentExact)/scenarios,
		MinimumMargin:minMargin,MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues),
	}
	p.Gate=p.ValueAccuracy>=0.99&&p.ExactScenarioAccuracy>=0.95
	return p,nil
}

func RunUP57C()(UP57CGoldenScaleResult,error){
	schedules:=[]int{65000000,66000000}
	type condition struct{banks int;noise float64}
	conditions:=[]condition{
		{4,0.04},
		{5,0},{5,0.005},{5,0.01},{5,0.02},{5,0.03},{5,0.04},
		{6,0},
	}
	result:=UP57CGoldenScaleResult{
		Schema:UP57CGoldenScaleSchema,Experiment:"UP-57C-golden-bank-scale",
		SourceUP56CSeal:"8c61fb93f7764367ef6a17405431756c9537a0f3",
		StateDimension:16,ScheduleBases:append([]int(nil),schedules...),SelectionPerformed:false,
	}
	for _,seedBase:=range schedules{
		for _,c:=range conditions{
			p,err:=up57cEvaluate(c.banks,c.noise,seedBase);if err!=nil{return UP57CGoldenScaleResult{},err}
			result.Points=append(result.Points,p)
		}
	}
	return result,nil
}
