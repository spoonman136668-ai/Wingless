package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

const UP53CPhaseMultiplexSchema = "wingless.up53c-phase-multiplex-memory.v1"

type UP53CMultiplexMetric struct {
	Banks                     int     `json:"banks"`
	Scenarios                 int     `json:"scenarios"`
	PrototypeCountPerEntity   int     `json:"prototype_count_per_entity"`
	CoherenceValueAccuracy    float64 `json:"coherence_value_accuracy"`
	CoherenceExactScenario    float64 `json:"coherence_exact_scenario_accuracy"`
	CoherenceMinimumMargin    float64 `json:"coherence_minimum_margin"`
	MagnitudeValueAccuracy    float64 `json:"magnitude_value_accuracy"`
	MagnitudeExactScenario    float64 `json:"magnitude_exact_scenario_accuracy"`
	MaxNormDrift              float64 `json:"max_norm_drift"`
	Gate                      bool    `json:"gate"`
}

type UP53CPhaseMultiplexResult struct {
	Schema                    string                  `json:"schema"`
	Experiment                string                  `json:"experiment"`
	SourceUP52CSeal           string                  `json:"source_up52c_seal"`
	StateDimension            int                     `json:"state_dimension"`
	Entities                  int                     `json:"entities"`
	ValuesPerEntity           int                     `json:"values_per_entity"`
	MemoryNoise               float64                 `json:"memory_noise"`
	GlobalPhaseNuisance       bool                    `json:"global_phase_nuisance"`
	BankLevels                []int                   `json:"bank_levels"`
	PhaseTags                 []float64               `json:"phase_tags"`
	CoherenceGlobalPhaseInvariant bool                `json:"coherence_global_phase_invariant"`
	MagnitudeControl          bool                    `json:"magnitude_control"`
	Metrics                   []UP53CMultiplexMetric `json:"metrics"`
	MaximumPassingBanks       int                     `json:"maximum_passing_banks"`
	FirstFailingBanks         int                     `json:"first_failing_banks"`
}

type up53cPrototype struct {
	values []int
	state  State
}

func up53cPhaseTags(count int) []float64 {
	out:=make([]float64,count)
	for bank:=0;bank<count;bank++{
		b:=float64(bank)
		out[bank]=0.41*b+0.173*b*b
	}
	return out
}

func up53cLocalEncode(values []int,tags []float64)(State,error){
	if len(values)!=len(tags){return nil,fmt.Errorf("UP53C value/tag count mismatch")}
	state:=make(State,4)
	for bank,value:=range values{
		if value<0||value>=4{return nil,fmt.Errorf("UP53C value out of range")}
		state[value]+=cmplx.Rect(1,tags[bank])
	}
	norm,err:=NormSquared(state)
	if err!=nil{return nil,err}
	if norm<=1e-15{return nil,fmt.Errorf("UP53C phase-code cancellation")}
	return Normalize(state)
}

func up53cScenarioState(scenario,banks int,tags []float64)(State,[][]int,error){
	tables:=make([][]int,banks)
	for bank:=0;bank<banks;bank++{
		tables[bank]=make([]int,4)
		for entity:=0;entity<4;entity++{
			value:=((scenario+1)*(bank*7+3)+entity*11+scenario*scenario*(bank+1)+bank*bank)%4
			tables[bank][entity]=value
		}
	}
	state:=make(State,16)
	for entity:=0;entity<4;entity++{
		values:=make([]int,banks)
		for bank:=0;bank<banks;bank++{values[bank]=tables[bank][entity]}
		local,err:=up53cLocalEncode(values,tags)
		if err!=nil{return nil,nil,err}
		for value:=0;value<4;value++{
			state[entity*4+value]=local[value]
		}
	}
	state,err:=Normalize(state)
	if err!=nil{return nil,nil,err}
	return state,tables,nil
}

func up53cPrototypeBank(banks int,tags []float64)([]up53cPrototype,error){
	count:=1
	for i:=0;i<banks;i++{count*=4}
	out:=make([]up53cPrototype,0,count)
	for code:=0;code<count;code++{
		x:=code
		values:=make([]int,banks)
		for bank:=banks-1;bank>=0;bank--{
			values[bank]=x%4
			x/=4
		}
		state,err:=up53cLocalEncode(values,tags)
		if err!=nil{return nil,err}
		out=append(out,up53cPrototype{values:append([]int(nil),values...),state:state})
	}
	return out,nil
}

func up53cFidelity(a,b State)(float64,error){
	if len(a)!=len(b)||len(a)==0{return 0,fmt.Errorf("UP53C fidelity dimension mismatch")}
	var inner complex128
	var an,bn float64
	for i:=range a{
		inner+=cmplx.Conj(a[i])*b[i]
		an+=cmplx.Abs(a[i])*cmplx.Abs(a[i])
		bn+=cmplx.Abs(b[i])*cmplx.Abs(b[i])
	}
	if an<=0||bn<=0{return 0,fmt.Errorf("UP53C zero fidelity norm")}
	return cmplx.Abs(inner)/math.Sqrt(an*bn),nil
}

func up53cMagnitudeCosine(a,b State)(float64,error){
	if len(a)!=len(b)||len(a)==0{return 0,fmt.Errorf("UP53C magnitude dimension mismatch")}
	var dot,an,bn float64
	for i:=range a{
		av:=cmplx.Abs(a[i]);bv:=cmplx.Abs(b[i])
		dot+=av*bv;an+=av*av;bn+=bv*bv
	}
	if an<=0||bn<=0{return 0,fmt.Errorf("UP53C zero magnitude norm")}
	return dot/math.Sqrt(an*bn),nil
}

func up53cDecode(local State,prototypes []up53cPrototype,coherence bool)([]int,float64,error){
	best:=-1.0
	second:=-1.0
	bestIndex:=-1
	for i,p:=range prototypes{
		var score float64
		var err error
		if coherence{score,err=up53cFidelity(local,p.state)}else{score,err=up53cMagnitudeCosine(local,p.state)}
		if err!=nil{return nil,0,err}
		if score>best{
			second=best;best=score;bestIndex=i
		}else if score>second{
			second=score
		}
	}
	if bestIndex<0{return nil,0,fmt.Errorf("UP53C no prototype")}
	margin:=best-second
	return append([]int(nil),prototypes[bestIndex].values...),margin,nil
}

func up53cTransportBlock()[]Coupling{
	var out []Coupling
	for entity:=0;entity<4;entity++{
		base:=entity*4
		out=append(out,
			Coupling{A:base+0,B:base+1,Theta:0.173},
			Coupling{A:base+1,B:base+2,Theta:-0.119},
			Coupling{A:base+2,B:base+3,Theta:0.137},
			Coupling{A:base+0,B:base+3,Theta:-0.091},
		)
	}
	return out
}

func up53cTransportLocalPrototype(p State,depth int)(State,error){
	state:=make(State,16)
	for i:=0;i<4;i++{state[i]=p[i]}
	block:=up53cTransportBlock()
	var err error
	for d:=0;d<depth;d++{
		state,err=Propagate(state,block)
		if err!=nil{return nil,err}
	}
	local:=append(State(nil),state[:4]...)
	return Normalize(local)
}

func RunUP53C()(UP53CPhaseMultiplexResult,error){
	const(
		scenarios=256
		noise=0.05
		depth=64
	)
	levels:=[]int{1,2,3,4,5,6}
	result:=UP53CPhaseMultiplexResult{
		Schema:UP53CPhaseMultiplexSchema,
		Experiment:"UP-53C-phase-multiplex-memory",
		SourceUP52CSeal:"9a1df4dd9a9556905e7030db633d3064068663a8",
		StateDimension:16,
		Entities:4,
		ValuesPerEntity:4,
		MemoryNoise:noise,
		GlobalPhaseNuisance:true,
		BankLevels:append([]int(nil),levels...),
		PhaseTags:up53cPhaseTags(6),
		CoherenceGlobalPhaseInvariant:true,
		MagnitudeControl:true,
	}
	block:=up53cTransportBlock()
	for _,banks:=range levels{
		tags:=up53cPhaseTags(banks)
		prototypes,err:=up53cPrototypeBank(banks,tags)
		if err!=nil{return UP53CPhaseMultiplexResult{},err}
		for i:=range prototypes{
			prototypes[i].state,err=up53cTransportLocalPrototype(prototypes[i].state,depth)
			if err!=nil{return UP53CPhaseMultiplexResult{},err}
		}

		var coherenceCorrect,magnitudeCorrect,totalValues int
		var coherenceExact,magnitudeExact int
		minMargin:=math.Inf(1)
		var maxNormDrift float64

		for scenario:=0;scenario<scenarios;scenario++{
			state,tables,err:=up53cScenarioState(scenario,banks,tags)
			if err!=nil{return UP53CPhaseMultiplexResult{},err}
			noiseState:=memoryNoise(53000000+banks*100000+scenario,16,noise)
			for i:=range state{state[i]+=noiseState[i]}
			state,err=Normalize(state)
			if err!=nil{return UP53CPhaseMultiplexResult{},err}
			state=rotateGlobalPhase(state,math.Mod(0.317*float64(scenario+banks+1),2*math.Pi))
			for d:=0;d<depth;d++{
				state,err=Propagate(state,block)
				if err!=nil{return UP53CPhaseMultiplexResult{},err}
			}
			norm2,err:=NormSquared(state)
			if err!=nil{return UP53CPhaseMultiplexResult{},err}
			if drift:=math.Abs(norm2-1);drift>maxNormDrift{maxNormDrift=drift}

			coherentScenario:=true
			magnitudeScenario:=true
			for entity:=0;entity<4;entity++{
				local:=append(State(nil),state[entity*4:(entity+1)*4]...)
				coherentValues,margin,err:=up53cDecode(local,prototypes,true)
				if err!=nil{return UP53CPhaseMultiplexResult{},err}
				magnitudeValues,_,err:=up53cDecode(local,prototypes,false)
				if err!=nil{return UP53CPhaseMultiplexResult{},err}
				if margin<minMargin{minMargin=margin}
				for bank:=0;bank<banks;bank++{
					totalValues++
					want:=tables[bank][entity]
					if coherentValues[bank]==want{coherenceCorrect++}else{coherentScenario=false}
					if magnitudeValues[bank]==want{magnitudeCorrect++}else{magnitudeScenario=false}
				}
			}
			if coherentScenario{coherenceExact++}
			if magnitudeScenario{magnitudeExact++}
		}
		metric:=UP53CMultiplexMetric{
			Banks:banks,
			Scenarios:scenarios,
			PrototypeCountPerEntity:len(prototypes),
			CoherenceValueAccuracy:float64(coherenceCorrect)/float64(totalValues),
			CoherenceExactScenario:float64(coherenceExact)/scenarios,
			CoherenceMinimumMargin:minMargin,
			MagnitudeValueAccuracy:float64(magnitudeCorrect)/float64(totalValues),
			MagnitudeExactScenario:float64(magnitudeExact)/scenarios,
			MaxNormDrift:maxNormDrift,
		}
		metric.Gate=metric.CoherenceValueAccuracy>=0.99&&
			metric.CoherenceExactScenario>=0.95&&
			metric.MaxNormDrift<=1e-9
		result.Metrics=append(result.Metrics,metric)
		if metric.Gate{
			result.MaximumPassingBanks=banks
		}else if result.FirstFailingBanks==0{
			result.FirstFailingBanks=banks
		}
	}
	return result,nil
}
