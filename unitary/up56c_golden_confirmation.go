package unitary

import "math"

const UP56CGoldenConfirmationSchema = "wingless.up56c-golden-confirmation.v1"

type UP56CGoldenPoint struct {
	ScheduleBase          int     `json:"schedule_base"`
	MemoryNoise           float64 `json:"memory_noise"`
	ValueAccuracy         float64 `json:"value_accuracy"`
	ExactScenarioAccuracy float64 `json:"exact_scenario_accuracy"`
	MinimumMargin         float64 `json:"minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	Gate                  bool    `json:"gate"`
}
type UP56CGoldenConfirmationResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP55CSeal    string             `json:"source_up55c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	GoldenTags         []float64          `json:"golden_tags"`
	ScheduleBases      []int              `json:"schedule_bases"`
	NoiseLevels        []float64          `json:"noise_levels"`
	SelectionPerformed bool               `json:"selection_performed"`
	Points             []UP56CGoldenPoint `json:"points"`
}

func up56cEvaluate(tags []float64, noise float64, seedBase int) (UP56CGoldenPoint,error) {
	const(banks=4;scenarios=256;depth=64)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP56CGoldenPoint{},err}
	for i:=range prototypes {
		prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth)
		if err!=nil{return UP56CGoldenPoint{},err}
	}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1)
	for scenario:=0;scenario<scenarios;scenario++ {
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP56CGoldenPoint{},err}
		n:=memoryNoise(seedBase+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP56CGoldenPoint{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP56CGoldenPoint{},err}}
		exact:=true
		for entity:=0;entity<4;entity++ {
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP56CGoldenPoint{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP56CGoldenPoint{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++ {
				totalValues++;want:=tables[bank][entity]
				if cv[bank]==want{coherentCorrect++}else{exact=false}
				if mv[bank]==want{magnitudeCorrect++}
			}
		}
		if exact{coherentExact++}
	}
	p:=UP56CGoldenPoint{
		ScheduleBase:seedBase,MemoryNoise:noise,
		ValueAccuracy:float64(coherentCorrect)/float64(totalValues),
		ExactScenarioAccuracy:float64(coherentExact)/scenarios,
		MinimumMargin:minMargin,
		MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues),
	}
	p.Gate=p.ValueAccuracy>=0.99&&p.ExactScenarioAccuracy>=0.95
	return p,nil
}

func RunUP56C()(UP56CGoldenConfirmationResult,error) {
	tags:=[]float64{0,2.399963229728653,4.799926459457306,7.199889689185959}
	schedules:=[]int{61000000,62000000,63000000}
	noises:=[]float64{0.035,0.04,0.0425,0.045,0.0475,0.05}
	result:=UP56CGoldenConfirmationResult{
		Schema:UP56CGoldenConfirmationSchema,
		Experiment:"UP-56C-golden-confirmation",
		SourceUP55CSeal:"986d4f350c985b753f89fa535e4040c54f3e7c21",
		Banks:4,StateDimension:16,GoldenTags:append([]float64(nil),tags...),
		ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),
		SelectionPerformed:false,
	}
	for _,seedBase:=range schedules {
		for _,noise:=range noises {
			p,err:=up56cEvaluate(tags,noise,seedBase);if err!=nil{return UP56CGoldenConfirmationResult{},err}
			result.Points=append(result.Points,p)
		}
	}
	return result,nil
}
