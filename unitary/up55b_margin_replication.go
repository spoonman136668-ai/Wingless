package unitary

import (
	"fmt"
	"math"
)

const UP55BMarginReplicationSchema = "wingless.up55b-margin-replication.v1"

type UP55BReplicationPoint struct {
	ScheduleBase       int              `json:"schedule_base"`
	MemoryNoise        float64          `json:"memory_noise"`
	WritesPerScenario  int              `json:"writes_per_scenario"`
	CommitAccuracy     float64          `json:"commit_accuracy"`
	Bins               []UP54BMarginBin `json:"bins"`
	MeanMarginCorrect  float64          `json:"mean_margin_correct"`
	MeanMarginIncorrect float64         `json:"mean_margin_incorrect"`
	LowMarginAccuracy  float64          `json:"low_margin_accuracy"`
	HighMarginAccuracy float64          `json:"high_margin_accuracy"`
	HighExceedsLow     bool             `json:"high_exceeds_low"`
}
type UP55BMarginReplicationResult struct {
	Schema               string                  `json:"schema"`
	Experiment           string                  `json:"experiment"`
	SourceUP54BSeal      string                  `json:"source_up54b_seal"`
	ScheduleBases        []int                   `json:"schedule_bases"`
	NoiseLevels          []float64               `json:"noise_levels"`
	WriteLevels          []int                   `json:"write_levels"`
	BinEdges             []float64               `json:"bin_edges"`
	TrainingChanged      bool                    `json:"training_changed"`
	CommitBehaviorChanged bool                   `json:"commit_behavior_changed"`
	Points               []UP55BReplicationPoint `json:"points"`
	AllSchedulesSeparated bool                   `json:"all_schedules_separated"`
}

func up55bCollect(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
	edges []float64,
	seedBase int,
) (UP55BReplicationPoint, error) {
	const scenarios = 48
	bins := make([]UP54BMarginBin, len(edges)-1)
	for i := range bins { bins[i].Lower=edges[i]; bins[i].Upper=edges[i+1] }

	var total, correct int
	var correctMargin, incorrectMargin float64
	var correctN, incorrectN int

	for scenarioIndex:=0; scenarioIndex<scenarios; scenarioIndex++ {
		scenario:=makeFullRankScenario(scenarioIndex,heldTables[(scenarioIndex*7)%len(heldTables)],writes,depths)
		pathTable:=scenario.initial
		trueTable:=scenario.initial
		for writeIndex,write:=range scenario.writes {
			canonical,err:=encodeMemory(pathTable);if err!=nil{return UP55BReplicationPoint{},err}
			seed:=seedBase+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return UP55BReplicationPoint{},err}
			memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
			state,err:=fullLatentEncode(memory,mixer);if err!=nil{return UP55BReplicationPoint{},err}
			operator,ok:=prepared.Ops[write.gap];if !ok{return UP55BReplicationPoint{},fmt.Errorf("missing UP55B depth=%d",write.gap)}
			state,err=latentMatrixVector(operator,state);if err!=nil{return UP55BReplicationPoint{},err}
			decoded,_,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers);if err!=nil{return UP55BReplicationPoint{},err}
			isCorrect:=decoded==trueTable
			total++
			idx:=up54bBinIndex(margin,edges);bins[idx].Samples++
			if isCorrect {correct++;bins[idx].Correct++;correctMargin+=margin;correctN++} else {incorrectMargin+=margin;incorrectN++}
			pathTable,err=applyMemoryWrite(decoded,write.entity,write.value);if err!=nil{return UP55BReplicationPoint{},err}
			trueTable,err=applyMemoryWrite(trueTable,write.entity,write.value);if err!=nil{return UP55BReplicationPoint{},err}
		}
	}
	for i:=range bins{if bins[i].Samples>0{bins[i].Accuracy=float64(bins[i].Correct)/float64(bins[i].Samples)}}

	var lowN,lowCorrect,highN,highCorrect int
	for _,b:=range bins {
		if b.Upper <= 0.1 { lowN+=b.Samples; lowCorrect+=b.Correct }
		if b.Lower >= 0.25 { highN+=b.Samples; highCorrect+=b.Correct }
	}
	p:=UP55BReplicationPoint{
		ScheduleBase:seedBase,MemoryNoise:memoryNoise,WritesPerScenario:writes,
		CommitAccuracy:float64(correct)/float64(total),Bins:bins,
	}
	if correctN>0{p.MeanMarginCorrect=correctMargin/float64(correctN)}
	if incorrectN>0{p.MeanMarginIncorrect=incorrectMargin/float64(incorrectN)}
	if lowN>0{p.LowMarginAccuracy=float64(lowCorrect)/float64(lowN)}
	if highN>0{p.HighMarginAccuracy=float64(highCorrect)/float64(highN)}
	p.HighExceedsLow=p.HighMarginAccuracy>p.LowMarginAccuracy
	return p,nil
}

func RunUP55B()(UP55BMarginReplicationResult,error){
	schedules:=[]int{56000000,57000000,58000000}
	noises:=[]float64{0.065,0.07}
	writes:=[]int{32,64}
	edges:=[]float64{0,0.001,0.01,0.05,0.1,0.25,1.0000001}
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP55BMarginReplicationResult{},err}
	result:=UP55BMarginReplicationResult{
		Schema:UP55BMarginReplicationSchema,Experiment:"UP-55B-margin-replication",
		SourceUP54BSeal:"d5964e1a7d07d747c764220a48d1101b49cec814",
		ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),
		WriteLevels:append([]int(nil),writes...),BinEdges:append([]float64(nil),edges...),
		TrainingChanged:false,CommitBehaviorChanged:false,AllSchedulesSeparated:true,
	}
	for _,noise:=range noises {
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP55BMarginReplicationResult{},err}
		for _,seedBase:=range schedules {
			for _,w:=range writes {
				p,err:=up55bCollect(mixer,prepared,heldTables,heldDepths,noise,w,edges,seedBase)
				if err!=nil{return UP55BMarginReplicationResult{},err}
				if !p.HighExceedsLow { result.AllSchedulesSeparated=false }
				result.Points=append(result.Points,p)
			}
		}
	}
	return result,nil
}
