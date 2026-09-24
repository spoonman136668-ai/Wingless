package unitary

import (
	"math"
)

const UP54CFourBankNoiseSchema = "wingless.up54c-four-bank-noise-ladder.v1"

type UP54CNoisePoint struct {
	Banks                  int     `json:"banks"`
	MemoryNoise            float64 `json:"memory_noise"`
	CoherenceValueAccuracy float64 `json:"coherence_value_accuracy"`
	CoherenceExactScenario float64 `json:"coherence_exact_scenario_accuracy"`
	CoherenceMinimumMargin float64 `json:"coherence_minimum_margin"`
	MagnitudeValueAccuracy float64 `json:"magnitude_value_accuracy"`
	MaxNormDrift           float64 `json:"max_norm_drift"`
	Gate                   bool    `json:"gate"`
}

type UP54CBankBoundary struct {
	Banks              int     `json:"banks"`
	MaximumPassingNoise float64 `json:"maximum_passing_noise"`
	FirstFailingNoise   float64 `json:"first_failing_noise"`
}

type UP54CFourBankNoiseResult struct {
	Schema          string              `json:"schema"`
	Experiment      string              `json:"experiment"`
	SourceUP53CSeal string              `json:"source_up53c_seal"`
	StateDimension  int                 `json:"state_dimension"`
	BankLevels      []int               `json:"bank_levels"`
	NoiseLevels     []float64           `json:"noise_levels"`
	PhaseTagsFrozen bool                `json:"phase_tags_frozen"`
	TransportDepth  int                 `json:"transport_depth"`
	Scenarios       int                 `json:"scenarios"`
	Points          []UP54CNoisePoint   `json:"points"`
	Boundaries      []UP54CBankBoundary `json:"boundaries"`
}

func up54cEvaluate(banks int,noise float64)(UP54CNoisePoint,error){
	const(scenarios=256;depth=64)
	tags:=up53cPhaseTags(banks)
	prototypes,err:=up53cPrototypeBank(banks,tags);if err!=nil{return UP54CNoisePoint{},err}
	for i:=range prototypes{prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth);if err!=nil{return UP54CNoisePoint{},err}}
	block:=up53cTransportBlock()
	var coherentCorrect,magnitudeCorrect,totalValues,coherentExact int
	minMargin:=math.Inf(1);var maxNormDrift float64
	for scenario:=0;scenario<scenarios;scenario++{
		state,tables,err:=up53cScenarioState(scenario,banks,tags);if err!=nil{return UP54CNoisePoint{},err}
		n:=memoryNoise(54000000+banks*100000+int(math.Round(noise*100000))*1000+scenario,16,noise)
		for i:=range state{state[i]+=n[i]}
		state,err=Normalize(state);if err!=nil{return UP54CNoisePoint{},err}
		state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
		for d:=0;d<depth;d++{state,err=Propagate(state,block);if err!=nil{return UP54CNoisePoint{},err}}
		norm2,err:=NormSquared(state);if err!=nil{return UP54CNoisePoint{},err};if drift:=math.Abs(norm2-1);drift>maxNormDrift{maxNormDrift=drift}
		exact:=true
		for entity:=0;entity<4;entity++{
			local:=append(State(nil),state[entity*4:(entity+1)*4]...)
			cv,margin,err:=up53cDecode(local,prototypes,true);if err!=nil{return UP54CNoisePoint{},err}
			mv,_,err:=up53cDecode(local,prototypes,false);if err!=nil{return UP54CNoisePoint{},err}
			if margin<minMargin{minMargin=margin}
			for bank:=0;bank<banks;bank++{
				totalValues++;want:=tables[bank][entity]
				if cv[bank]==want{coherentCorrect++}else{exact=false}
				if mv[bank]==want{magnitudeCorrect++}
			}
		}
		if exact{coherentExact++}
	}
	p:=UP54CNoisePoint{
		Banks:banks,MemoryNoise:noise,
		CoherenceValueAccuracy:float64(coherentCorrect)/float64(totalValues),
		CoherenceExactScenario:float64(coherentExact)/scenarios,
		CoherenceMinimumMargin:minMargin,
		MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues),
		MaxNormDrift:maxNormDrift,
	}
	p.Gate=p.CoherenceValueAccuracy>=0.99&&p.CoherenceExactScenario>=0.95&&p.MaxNormDrift<=1e-9
	return p,nil
}

func RunUP54C()(UP54CFourBankNoiseResult,error){
	banks:=[]int{3,4}
	noises:=[]float64{0,0.005,0.01,0.02,0.03,0.04,0.05}
	result:=UP54CFourBankNoiseResult{
		Schema:UP54CFourBankNoiseSchema,Experiment:"UP-54C-four-bank-noise-ladder",
		SourceUP53CSeal:"31b5136faf6a3575728a91d18886fff8f6cce122",
		StateDimension:16,BankLevels:append([]int(nil),banks...),NoiseLevels:append([]float64(nil),noises...),
		PhaseTagsFrozen:true,TransportDepth:64,Scenarios:256,
	}
	for _,b:=range banks{
		boundary:=UP54CBankBoundary{Banks:b}
		for _,noise:=range noises{
			p,err:=up54cEvaluate(b,noise);if err!=nil{return UP54CFourBankNoiseResult{},err}
			result.Points=append(result.Points,p)
			if p.Gate{boundary.MaximumPassingNoise=noise}else if boundary.FirstFailingNoise==0{boundary.FirstFailingNoise=noise}
		}
		result.Boundaries=append(result.Boundaries,boundary)
	}
	return result,nil
}
