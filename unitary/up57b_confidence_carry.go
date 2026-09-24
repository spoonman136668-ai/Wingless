package unitary

import (
	"fmt"
	"math"
)

const UP57BConfidenceCarrySchema = "wingless.up57b-confidence-carry.v1"

type UP57BPolicyMetric struct {
	ScheduleBase        int     `json:"schedule_base"`
	MemoryNoise         float64 `json:"memory_noise"`
	WritesPerScenario   int     `json:"writes_per_scenario"`
	Policy              string  `json:"policy"`
	MarginThreshold     float64 `json:"margin_threshold"`
	LowConfidenceEvents int     `json:"low_confidence_events"`
	CommitAccuracy      float64 `json:"commit_accuracy"`
	FinalAccuracy       float64 `json:"final_accuracy"`
	RelationAccuracy    float64 `json:"relation_accuracy"`
	MaxNormDrift        float64 `json:"max_norm_drift"`
}
type UP57BConfidenceCarryResult struct {
	Schema              string              `json:"schema"`
	Experiment          string              `json:"experiment"`
	SourceUP56BSeal     string              `json:"source_up56b_seal"`
	FrozenPair          []int               `json:"frozen_pair"`
	ScheduleBases       []int               `json:"schedule_bases"`
	NoiseLevels         []float64           `json:"noise_levels"`
	WritesPerScenario   int                 `json:"writes_per_scenario"`
	Thresholds          []float64           `json:"thresholds"`
	OracleUsed          bool                `json:"oracle_used"`
	TrainingChanged     bool                `json:"training_changed"`
	Metrics             []UP57BPolicyMetric `json:"metrics"`
}

func runUP57BPolicy(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
	seedBase int,
	threshold float64,
	useGate bool,
) (UP57BPolicyMetric,error) {
	const scenarios=48
	relationSamples,err:=relationHeadTrainingSamples()
	if err!=nil{return UP57BPolicyMetric{},err}
	relationHead,_,err:=trainLinearSoftmax(relationSamples,4,16,600,1.0)
	if err!=nil{return UP57BPolicyMetric{},err}

	var commitCorrect,commitTotal,finalCorrect,relationCorrect,lowConfidence int
	var maxNormDrift float64
	for scenarioIndex:=0;scenarioIndex<scenarios;scenarioIndex++ {
		scenario:=makeFullRankScenario(scenarioIndex,heldTables[(scenarioIndex*7)%len(heldTables)],writes,depths)
		pathTable:=scenario.initial
		trueTable:=scenario.initial

		for writeIndex,write:=range scenario.writes {
			canonical,err:=encodeMemory(pathTable);if err!=nil{return UP57BPolicyMetric{},err}
			seed:=seedBase+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return UP57BPolicyMetric{},err}
			memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
			state,err:=fullLatentEncode(memory,mixer);if err!=nil{return UP57BPolicyMetric{},err}
			operator,ok:=prepared.Ops[write.gap];if !ok{return UP57BPolicyMetric{},fmt.Errorf("missing UP57B depth=%d",write.gap)}
			state,err=latentMatrixVector(operator,state);if err!=nil{return UP57BPolicyMetric{},err}
			norm2,err:=NormSquared(state);if err!=nil{return UP57BPolicyMetric{},err}
			if d:=math.Abs(norm2-1);d>maxNormDrift{maxNormDrift=d}
			decoded,_,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers)
			if err!=nil{return UP57BPolicyMetric{},err}
			commitTotal++
			if decoded==trueTable{commitCorrect++}
			trueTable,err=applyMemoryWrite(trueTable,write.entity,write.value);if err!=nil{return UP57BPolicyMetric{},err}
			if useGate && margin<threshold {
				lowConfidence++
				pathTable,err=applyMemoryWrite(pathTable,write.entity,write.value)
			}else{
				pathTable,err=applyMemoryWrite(decoded,write.entity,write.value)
			}
			if err!=nil{return UP57BPolicyMetric{},err}
		}

		canonical,err:=encodeMemory(pathTable);if err!=nil{return UP57BPolicyMetric{},err}
		seed:=seedBase+scenarioIndex*10000+9999
		memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return UP57BPolicyMetric{},err}
		memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
		state,err:=fullLatentEncode(memory,mixer);if err!=nil{return UP57BPolicyMetric{},err}
		operator,ok:=prepared.Ops[scenario.finalGap];if !ok{return UP57BPolicyMetric{},fmt.Errorf("missing UP57B final depth=%d",scenario.finalGap)}
		state,err=latentMatrixVector(operator,state);if err!=nil{return UP57BPolicyMetric{},err}
		decoded,distributions,_,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers);if err!=nil{return UP57BPolicyMetric{},err}
		if decoded==trueTable{finalCorrect++}
		relationInput,err:=relationFeatures(distributions[scenario.queryA],distributions[scenario.queryB]);if err!=nil{return UP57BPolicyMetric{},err}
		relationProbabilities,err:=relationHead.probabilities(relationInput);if err!=nil{return UP57BPolicyMetric{},err}
		got,_,err:=classAndMargin(relationProbabilities);if err!=nil{return UP57BPolicyMetric{},err}
		want,err:=memoryRelation(trueTable,scenario.queryA,scenario.queryB);if err!=nil{return UP57BPolicyMetric{},err}
		if got==want{relationCorrect++}
	}
	policy:="baseline_closed_loop"
	if useGate{policy="confidence_carry"}
	return UP57BPolicyMetric{
		ScheduleBase:seedBase,MemoryNoise:memoryNoise,WritesPerScenario:writes,
		Policy:policy,MarginThreshold:threshold,LowConfidenceEvents:lowConfidence,
		CommitAccuracy:float64(commitCorrect)/float64(commitTotal),
		FinalAccuracy:float64(finalCorrect)/float64(scenarios),
		RelationAccuracy:float64(relationCorrect)/float64(scenarios),
		MaxNormDrift:maxNormDrift,
	},nil
}

func RunUP57B()(UP57BConfidenceCarryResult,error){
	schedules:=[]int{67000000,68000000,69000000}
	noises:=[]float64{0.065,0.07}
	thresholds:=[]float64{0.25,0.5,0.75}
	const writes=64
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP57BConfidenceCarryResult{},err}
	result:=UP57BConfidenceCarryResult{
		Schema:UP57BConfidenceCarrySchema,Experiment:"UP-57B-confidence-carry",
		SourceUP56BSeal:"4e56b9fb943e82b72b8c5172c8a211f81b30665d",
		FrozenPair:[]int{0,5},ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),WritesPerScenario:writes,
		Thresholds:append([]float64(nil),thresholds...),OracleUsed:false,TrainingChanged:false,
	}
	for _,noise:=range noises {
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP57BConfidenceCarryResult{},err}
		for _,seedBase:=range schedules {
			base,err:=runUP57BPolicy(mixer,prepared,heldTables,heldDepths,noise,writes,seedBase,0,false)
			if err!=nil{return UP57BConfidenceCarryResult{},err}
			result.Metrics=append(result.Metrics,base)
			for _,threshold:=range thresholds {
				m,err:=runUP57BPolicy(mixer,prepared,heldTables,heldDepths,noise,writes,seedBase,threshold,true)
				if err!=nil{return UP57BConfidenceCarryResult{},err}
				result.Metrics=append(result.Metrics,m)
			}
		}
	}
	return result,nil
}
