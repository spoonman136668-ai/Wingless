package unitary

import (
	"fmt"
	"math"
)

const UP54BMarginSchema = "wingless.up54b-margin-calibration.v1"

type UP54BMarginBin struct {
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
	Samples int `json:"samples"`
	Correct int `json:"correct"`
	Accuracy float64 `json:"accuracy"`
}
type UP54BMarginPoint struct {
	MemoryNoise float64 `json:"memory_noise"`
	WritesPerScenario int `json:"writes_per_scenario"`
	CommitAccuracy float64 `json:"commit_accuracy"`
	Bins []UP54BMarginBin `json:"bins"`
	MeanMarginCorrect float64 `json:"mean_margin_correct"`
	MeanMarginIncorrect float64 `json:"mean_margin_incorrect"`
}
type UP54BMarginResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP53BSeal string `json:"source_up53b_seal"`
	FrozenPair []int `json:"frozen_pair"`
	NoiseLevels []float64 `json:"noise_levels"`
	WriteLevels []int `json:"write_levels"`
	BinEdges []float64 `json:"bin_edges"`
	TrainingChanged bool `json:"training_changed"`
	CommitBehaviorChanged bool `json:"commit_behavior_changed"`
	Points []UP54BMarginPoint `json:"points"`
}

func up54bBinIndex(margin float64,edges []float64) int {
	for i:=0;i<len(edges)-1;i++{if margin>=edges[i]&&margin<edges[i+1]{return i}}
	return len(edges)-2
}

func up54bCollect(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
	edges []float64,
)(UP54BMarginPoint,error){
	const scenarios=48
	bins:=make([]UP54BMarginBin,len(edges)-1)
	for i:=range bins{bins[i].Lower=edges[i];bins[i].Upper=edges[i+1]}
	var total,correct int
	var correctMargin,incorrectMargin float64
	var correctN,incorrectN int
	for scenarioIndex:=0;scenarioIndex<scenarios;scenarioIndex++{
		scenario:=makeFullRankScenario(scenarioIndex,heldTables[(scenarioIndex*7)%len(heldTables)],writes,depths)
		pathTable:=scenario.initial;trueTable:=scenario.initial
		for writeIndex,write:=range scenario.writes{
			canonical,err:=encodeMemory(pathTable);if err!=nil{return UP54BMarginPoint{},err}
			seed:=34000000+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return UP54BMarginPoint{},err}
			memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
			state,err:=fullLatentEncode(memory,mixer);if err!=nil{return UP54BMarginPoint{},err}
			operator,ok:=prepared.Ops[write.gap];if !ok{return UP54BMarginPoint{},fmt.Errorf("missing UP54B depth=%d",write.gap)}
			state,err=latentMatrixVector(operator,state);if err!=nil{return UP54BMarginPoint{},err}
			decoded,_,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers);if err!=nil{return UP54BMarginPoint{},err}
			isCorrect:=decoded==trueTable
			total++;idx:=up54bBinIndex(margin,edges);bins[idx].Samples++
			if isCorrect{correct++;bins[idx].Correct++;correctMargin+=margin;correctN++}else{incorrectMargin+=margin;incorrectN++}
			pathTable,err=applyMemoryWrite(decoded,write.entity,write.value);if err!=nil{return UP54BMarginPoint{},err}
			trueTable,err=applyMemoryWrite(trueTable,write.entity,write.value);if err!=nil{return UP54BMarginPoint{},err}
		}
	}
	for i:=range bins{if bins[i].Samples>0{bins[i].Accuracy=float64(bins[i].Correct)/float64(bins[i].Samples)}}
	p:=UP54BMarginPoint{MemoryNoise:memoryNoise,WritesPerScenario:writes,CommitAccuracy:float64(correct)/float64(total),Bins:bins}
	if correctN>0{p.MeanMarginCorrect=correctMargin/float64(correctN)}
	if incorrectN>0{p.MeanMarginIncorrect=incorrectMargin/float64(incorrectN)}
	return p,nil
}

func RunUP54B()(UP54BMarginResult,error){
	noises:=[]float64{0.065,0.07};writes:=[]int{16,32,64};edges:=[]float64{0,0.001,0.01,0.05,0.1,0.25,1.0000001}
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false);offsets,err:=up50bFusedOffsets();if err!=nil{return UP54BMarginResult{},err}
	result:=UP54BMarginResult{Schema:UP54BMarginSchema,Experiment:"UP-54B-margin-calibration",SourceUP53BSeal:"04f996031423faa596b07b150d8e30fb4b02b52b",FrozenPair:[]int{0,5},NoiseLevels:append([]float64(nil),noises...),WriteLevels:append([]int(nil),writes...),BinEdges:append([]float64(nil),edges...),TrainingChanged:false,CommitBehaviorChanged:false}
	for _,noise:=range noises{
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise);if err!=nil{return UP54BMarginResult{},err}
		for _,w:=range writes{p,err:=up54bCollect(mixer,prepared,heldTables,heldDepths,noise,w,edges);if err!=nil{return UP54BMarginResult{},err};result.Points=append(result.Points,p)}
	}
	return result,nil
}
