package unitary

import (
	"fmt"
	"math"
)

const UP56BSelectiveRiskSchema = "wingless.up56b-selective-risk.v1"

type UP56BThresholdMetric struct {
	Threshold        float64 `json:"threshold"`
	TotalSamples     int     `json:"total_samples"`
	SelectedSamples  int     `json:"selected_samples"`
	Coverage         float64 `json:"coverage"`
	SelectedAccuracy float64 `json:"selected_accuracy"`
	BaselineAccuracy float64 `json:"baseline_accuracy"`
	RiskReduction    float64 `json:"risk_reduction"`
}
type UP56BPoint struct {
	ScheduleBase      int                    `json:"schedule_base"`
	MemoryNoise       float64                `json:"memory_noise"`
	WritesPerScenario int                    `json:"writes_per_scenario"`
	Metrics           []UP56BThresholdMetric `json:"metrics"`
}
type UP56BSelectiveRiskResult struct {
	Schema                string       `json:"schema"`
	Experiment            string       `json:"experiment"`
	SourceUP55BSeal       string       `json:"source_up55b_seal"`
	ScheduleBases         []int        `json:"schedule_bases"`
	NoiseLevels           []float64    `json:"noise_levels"`
	WriteLevels           []int        `json:"write_levels"`
	Thresholds            []float64    `json:"thresholds"`
	TrainingChanged       bool         `json:"training_changed"`
	CommitBehaviorChanged bool         `json:"commit_behavior_changed"`
	Points                []UP56BPoint `json:"points"`
}

type up56bEvent struct {
	margin  float64
	correct bool
}

func up56bCollectEvents(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
	seedBase int,
) ([]up56bEvent,error) {
	const scenarios=48
	var events []up56bEvent
	for scenarioIndex:=0;scenarioIndex<scenarios;scenarioIndex++ {
		scenario:=makeFullRankScenario(scenarioIndex,heldTables[(scenarioIndex*7)%len(heldTables)],writes,depths)
		pathTable:=scenario.initial
		trueTable:=scenario.initial
		for writeIndex,write:=range scenario.writes {
			canonical,err:=encodeMemory(pathTable);if err!=nil{return nil,err}
			seed:=seedBase+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return nil,err}
			memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
			state,err:=fullLatentEncode(memory,mixer);if err!=nil{return nil,err}
			operator,ok:=prepared.Ops[write.gap];if !ok{return nil,fmt.Errorf("missing UP56B depth=%d",write.gap)}
			state,err=latentMatrixVector(operator,state);if err!=nil{return nil,err}
			decoded,_,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers);if err!=nil{return nil,err}
			events=append(events,up56bEvent{margin:margin,correct:decoded==trueTable})
			pathTable,err=applyMemoryWrite(decoded,write.entity,write.value);if err!=nil{return nil,err}
			trueTable,err=applyMemoryWrite(trueTable,write.entity,write.value);if err!=nil{return nil,err}
		}
	}
	return events,nil
}

func up56bMetrics(events []up56bEvent, thresholds []float64) []UP56BThresholdMetric {
	var baselineCorrect int
	for _,e:=range events{if e.correct{baselineCorrect++}}
	baseline:=float64(baselineCorrect)/float64(len(events))
	out:=make([]UP56BThresholdMetric,0,len(thresholds))
	for _,threshold:=range thresholds {
		var selected,correct int
		for _,e:=range events {
			if e.margin < threshold {continue}
			selected++
			if e.correct{correct++}
		}
		accuracy:=0.0
		if selected>0{accuracy=float64(correct)/float64(selected)}
		coverage:=float64(selected)/float64(len(events))
		out=append(out,UP56BThresholdMetric{
			Threshold:threshold,TotalSamples:len(events),SelectedSamples:selected,Coverage:coverage,
			SelectedAccuracy:accuracy,BaselineAccuracy:baseline,
			RiskReduction:(1-baseline)-(1-accuracy),
		})
	}
	return out
}

func RunUP56B()(UP56BSelectiveRiskResult,error){
	schedules:=[]int{59000000,60000000,64000000}
	noises:=[]float64{0.065,0.07}
	writes:=[]int{32,64}
	thresholds:=[]float64{0,0.05,0.1,0.25,0.5,0.75}
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP56BSelectiveRiskResult{},err}
	result:=UP56BSelectiveRiskResult{
		Schema:UP56BSelectiveRiskSchema,Experiment:"UP-56B-selective-risk",
		SourceUP55BSeal:"7d95194326285268ba65b87a079436145dc1cc6c",
		ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),
		WriteLevels:append([]int(nil),writes...),Thresholds:append([]float64(nil),thresholds...),
		TrainingChanged:false,CommitBehaviorChanged:false,
	}
	for _,noise:=range noises {
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP56BSelectiveRiskResult{},err}
		for _,seedBase:=range schedules {
			for _,w:=range writes {
				events,err:=up56bCollectEvents(mixer,prepared,heldTables,heldDepths,noise,w,seedBase)
				if err!=nil{return UP56BSelectiveRiskResult{},err}
				result.Points=append(result.Points,UP56BPoint{
					ScheduleBase:seedBase,MemoryNoise:noise,WritesPerScenario:w,
					Metrics:up56bMetrics(events,thresholds),
				})
			}
		}
	}
	return result,nil
}
