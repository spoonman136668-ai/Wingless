package unitary

import (
	"fmt"
	"math"
)

const UP53BOracleResetSchema = "wingless.up53b-oracle-reset-decomposition.v1"

type UP53BPoint struct {
	MemoryNoise              float64               `json:"memory_noise"`
	WritesPerScenario        int                   `json:"writes_per_scenario"`
	ClosedLoop               DiscoveredIntegration `json:"closed_loop"`
	OracleReset              DiscoveredIntegration `json:"oracle_reset"`
	FullControl              DiscoveredIntegration `json:"full_control"`
	ClosedLoopGate           bool                  `json:"closed_loop_gate"`
	OracleMinusClosedCommit  float64               `json:"oracle_minus_closed_commit"`
	OracleMinusClosedFinal   float64               `json:"oracle_minus_closed_final"`
}

type UP53BOracleResetResult struct {
	Schema                 string        `json:"schema"`
	Experiment             string        `json:"experiment"`
	SourceUP52BSeal        string        `json:"source_up52b_seal"`
	FrozenPair             []int         `json:"frozen_pair"`
	BoundaryNoises         []float64     `json:"boundary_noises"`
	WriteLevels            []int         `json:"write_levels"`
	HeldDepths             []int         `json:"held_depths"`
	Points                 []UP53BPoint `json:"points"`
	OracleResetUsedForPath bool          `json:"oracle_reset_used_for_path"`
	TrainingChanged        bool          `json:"training_changed"`
}

func runUP53BIntegration(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
	oracleReset bool,
) (DiscoveredIntegration,error) {
	const scenarios=48
	relationSamples,err:=relationHeadTrainingSamples()
	if err!=nil{return DiscoveredIntegration{},err}
	relationHead,_,err:=trainLinearSoftmax(relationSamples,4,16,600,1.0)
	if err!=nil{return DiscoveredIntegration{},err}

	var commitCorrect,commitTotal,finalCorrect,relationCorrect int
	minValueMargin:=math.Inf(1);minRelationMargin:=math.Inf(1);var maxNormDrift float64

	for scenarioIndex:=0;scenarioIndex<scenarios;scenarioIndex++{
		scenario:=makeFullRankScenario(scenarioIndex,heldTables[(scenarioIndex*7)%len(heldTables)],writes,depths)
		pathTable:=scenario.initial
		trueTable:=scenario.initial

		for writeIndex,write:=range scenario.writes{
			canonical,err:=encodeMemory(pathTable);if err!=nil{return DiscoveredIntegration{},err}
			seed:=34000000+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return DiscoveredIntegration{},err}
			memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
			state,err:=fullLatentEncode(memory,mixer);if err!=nil{return DiscoveredIntegration{},err}
			operator,ok:=prepared.Ops[write.gap];if !ok{return DiscoveredIntegration{},fmt.Errorf("missing UP53B mutable depth=%d",write.gap)}
			state,err=latentMatrixVector(operator,state);if err!=nil{return DiscoveredIntegration{},err}
			norm2,err:=NormSquared(state);if err!=nil{return DiscoveredIntegration{},err}
			if d:=math.Abs(norm2-1);d>maxNormDrift{maxNormDrift=d}
			decoded,_,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers)
			if err!=nil{return DiscoveredIntegration{},err}
			if margin<minValueMargin{minValueMargin=margin}
			commitTotal++;if decoded==trueTable{commitCorrect++}

			nextTrue,err:=applyMemoryWrite(trueTable,write.entity,write.value);if err!=nil{return DiscoveredIntegration{},err}
			if oracleReset{
				pathTable=nextTrue
			}else{
				pathTable,err=applyMemoryWrite(decoded,write.entity,write.value);if err!=nil{return DiscoveredIntegration{},err}
			}
			trueTable=nextTrue
		}

		canonical,err:=encodeMemory(pathTable);if err!=nil{return DiscoveredIntegration{},err}
		seed:=34000000+scenarioIndex*10000+9999
		memory,err:=perturbMemory(canonical,seed,memoryNoise);if err!=nil{return DiscoveredIntegration{},err}
		memory=rotateGlobalPhase(memory,math.Mod(0.271*float64(seed+1),2*math.Pi))
		state,err:=fullLatentEncode(memory,mixer);if err!=nil{return DiscoveredIntegration{},err}
		operator,ok:=prepared.Ops[scenario.finalGap];if !ok{return DiscoveredIntegration{},fmt.Errorf("missing UP53B final depth=%d",scenario.finalGap)}
		state,err=latentMatrixVector(operator,state);if err!=nil{return DiscoveredIntegration{},err}
		norm2,err:=NormSquared(state);if err!=nil{return DiscoveredIntegration{},err}
		if d:=math.Abs(norm2-1);d>maxNormDrift{maxNormDrift=d}
		decoded,distributions,margin,err:=decodeBreadthTable(state,prepared.Observables,prepared.Regressors,prepared.Classifiers)
		if err!=nil{return DiscoveredIntegration{},err}
		if margin<minValueMargin{minValueMargin=margin}
		if decoded==trueTable{finalCorrect++}
		relationInput,err:=relationFeatures(distributions[scenario.queryA],distributions[scenario.queryB]);if err!=nil{return DiscoveredIntegration{},err}
		relationProbabilities,err:=relationHead.probabilities(relationInput);if err!=nil{return DiscoveredIntegration{},err}
		got,rm,err:=classAndMargin(relationProbabilities);if err!=nil{return DiscoveredIntegration{},err}
		if rm<minRelationMargin{minRelationMargin=rm}
		want,err:=memoryRelation(trueTable,scenario.queryA,scenario.queryB);if err!=nil{return DiscoveredIntegration{},err}
		if got==want{relationCorrect++}
	}
	return DiscoveredIntegration{
		Scenarios:scenarios,WritesPerScenario:writes,
		CommitDecodeAccuracy:float64(commitCorrect)/float64(commitTotal),
		ExactFinalTableAccuracy:float64(finalCorrect)/float64(scenarios),
		RelationalQueryAccuracy:float64(relationCorrect)/float64(scenarios),
		MinValueMargin:minValueMargin,MinRelationMargin:minRelationMargin,MaxNormDrift:maxNormDrift,
	},nil
}

func RunUP53B()(UP53BOracleResetResult,error){
	noises:=[]float64{0.065,0.07}
	writeLevels:=[]int{16,32,64}
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	selectedOffsets,err:=up50bFusedOffsets();if err!=nil{return UP53BOracleResetResult{},err}
	fullOffsets:=[]float64{0,0,0,0,0,0}
	result:=UP53BOracleResetResult{
		Schema:UP53BOracleResetSchema,Experiment:"UP-53B-oracle-reset-decomposition",
		SourceUP52BSeal:"a18fafcbd02bfc7ee5f29eaef43a657f10146985",
		FrozenPair:[]int{0,5},BoundaryNoises:append([]float64(nil),noises...),WriteLevels:append([]int(nil),writeLevels...),HeldDepths:append([]int(nil),heldDepths...),
		OracleResetUsedForPath:true,TrainingChanged:false,
	}
	for _,noise:=range noises{
		selected,err:=up52bPrepareArm("selected",selectedOffsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP53BOracleResetResult{},err}
		full,err:=up52bPrepareArm("full",fullOffsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP53BOracleResetResult{},err}
		for _,writes:=range writeLevels{
			closed,err:=runUP53BIntegration(mixer,selected,heldTables,heldDepths,noise,writes,false);if err!=nil{return UP53BOracleResetResult{},err}
			oracle,err:=runUP53BIntegration(mixer,selected,heldTables,heldDepths,noise,writes,true);if err!=nil{return UP53BOracleResetResult{},err}
			fullClosed,err:=runUP53BIntegration(mixer,full,heldTables,heldDepths,noise,writes,false);if err!=nil{return UP53BOracleResetResult{},err}
			closedArm:=MultiplicityDoseArm{Name:"closed",Static:selected.Static,Integration:closed}
			fullArm:=MultiplicityDoseArm{Name:"full",Static:full.Static,Integration:fullClosed}
			result.Points=append(result.Points,UP53BPoint{
				MemoryNoise:noise,WritesPerScenario:writes,ClosedLoop:closed,OracleReset:oracle,FullControl:fullClosed,
				ClosedLoopGate:up49bHardGate(closedArm,fullArm),
				OracleMinusClosedCommit:oracle.CommitDecodeAccuracy-closed.CommitDecodeAccuracy,
				OracleMinusClosedFinal:oracle.ExactFinalTableAccuracy-closed.ExactFinalTableAccuracy,
			})
		}
	}
	return result,nil
}
