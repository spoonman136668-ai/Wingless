package unitary

import (
	"fmt"
	"math"
)

const UP52BBoundaryWriteSchema = "wingless.up52b-boundary-write-chain.v1"

type up52bPreparedArm struct {
	Static          DiscoveryBreadthArm
	Ops             map[int]latentMatrix
	Observables     []latentMatrix
	Regressors      [4]breadthRegressor
	Classifiers     [4]linearSoftmaxHead
	SelectedIndices []int
}

type UP52BWritePoint struct {
	WritesPerScenario int                   `json:"writes_per_scenario"`
	Selected          DiscoveredIntegration `json:"selected"`
	FullControl       DiscoveredIntegration `json:"full_control"`
	HeldDelta         float64               `json:"held_delta"`
	CommitDelta       float64               `json:"commit_delta"`
	HardGate          bool                  `json:"hard_gate"`
}

type UP52BNoiseSeries struct {
	MemoryNoise        float64             `json:"memory_noise"`
	StaticSelectedHeld float64             `json:"static_selected_held"`
	StaticFullHeld     float64             `json:"static_full_held"`
	Points             []UP52BWritePoint   `json:"points"`
	LastPassingWrites  int                 `json:"last_passing_writes"`
	FirstFailingWrites int                 `json:"first_failing_writes"`
}

type UP52BBoundaryWriteResult struct {
	Schema            string               `json:"schema"`
	Experiment        string               `json:"experiment"`
	SourceUP51BSeal   string               `json:"source_up51b_seal"`
	FrozenPair        []int                `json:"frozen_pair"`
	BoundaryNoises    []float64            `json:"boundary_noises"`
	WriteLevels       []int                `json:"write_levels"`
	HeldDepths        []int                `json:"held_depths"`
	SelectionRun      bool                 `json:"selection_run"`
	Series            []UP52BNoiseSeries   `json:"series"`
}

func up52bPrepareArm(
	name string,
	offsets []float64,
	mixer latentMatrix,
	trainTables, heldTables []memoryTable,
	trainDepths, heldDepths, allDepths []int,
	memoryNoise float64,
) (up52bPreparedArm, error) {
	step, err := fullLatentDoseStep(mixer, offsets)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	ops, err := latentDepthOperators(step, allDepths)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	candidates, _, err := discoverCommutingObservables(
		step, interactionCandidateCount, interactionRounds,
	)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	selectionStates, selectionTables, err := taskSelectedTrainingStates(
		trainTables, mixer, memoryNoise, 2,
	)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	selectedIndices, selectedObservables, err := selectInteractionRelevantObservables(
		selectionStates, selectionTables, candidates, interactionRuntimeCount,
	)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	regressors, classifiers, static, err := trainAndEvaluateMultiplicityArm(
		name,
		"up52b_"+name,
		trainTables, heldTables,
		trainDepths, heldDepths,
		ops, mixer, selectedObservables, memoryNoise,
	)
	if err != nil {
		return up52bPreparedArm{}, err
	}
	return up52bPreparedArm{
		Static:static,
		Ops:ops,
		Observables:selectedObservables,
		Regressors:regressors,
		Classifiers:classifiers,
		SelectedIndices:append([]int(nil),selectedIndices...),
	}, nil
}

func runUP52BIntegration(
	mixer latentMatrix,
	prepared up52bPreparedArm,
	heldTables []memoryTable,
	depths []int,
	memoryNoise float64,
	writes int,
) (DiscoveredIntegration, error) {
	const scenarios = 48

	relationSamples, err := relationHeadTrainingSamples()
	if err != nil {
		return DiscoveredIntegration{}, err
	}
	relationHead, _, err := trainLinearSoftmax(
		relationSamples, 4, 16, 600, 1.0,
	)
	if err != nil {
		return DiscoveredIntegration{}, err
	}

	var commitCorrect, commitTotal, finalCorrect, relationCorrect int
	minValueMargin := math.Inf(1)
	minRelationMargin := math.Inf(1)
	var maxNormDrift float64

	for scenarioIndex:=0; scenarioIndex<scenarios; scenarioIndex++ {
		scenario:=makeFullRankScenario(
			scenarioIndex,
			heldTables[(scenarioIndex*7)%len(heldTables)],
			writes,
			depths,
		)
		pathTable:=scenario.initial
		trueTable:=scenario.initial

		for writeIndex,write:=range scenario.writes {
			canonical,err:=encodeMemory(pathTable)
			if err!=nil{return DiscoveredIntegration{},err}
			seed:=34000000+scenarioIndex*10000+writeIndex*31
			memory,err:=perturbMemory(canonical,seed,memoryNoise)
			if err!=nil{return DiscoveredIntegration{},err}
			memory=rotateGlobalPhase(
				memory,
				math.Mod(0.271*float64(seed+1),2*math.Pi),
			)
			state,err:=fullLatentEncode(memory,mixer)
			if err!=nil{return DiscoveredIntegration{},err}
			operator,ok:=prepared.Ops[write.gap]
			if !ok{
				return DiscoveredIntegration{},fmt.Errorf("missing UP52B mutable depth=%d",write.gap)
			}
			state,err=latentMatrixVector(operator,state)
			if err!=nil{return DiscoveredIntegration{},err}
			norm2,err:=NormSquared(state)
			if err!=nil{return DiscoveredIntegration{},err}
			if drift:=math.Abs(norm2-1);drift>maxNormDrift{maxNormDrift=drift}

			decoded,_,margin,err:=decodeBreadthTable(
				state,prepared.Observables,prepared.Regressors,prepared.Classifiers,
			)
			if err!=nil{return DiscoveredIntegration{},err}
			if margin<minValueMargin{minValueMargin=margin}
			commitTotal++
			if decoded==trueTable{commitCorrect++}
			pathTable,err=applyMemoryWrite(decoded,write.entity,write.value)
			if err!=nil{return DiscoveredIntegration{},err}
			trueTable,err=applyMemoryWrite(trueTable,write.entity,write.value)
			if err!=nil{return DiscoveredIntegration{},err}
		}

		canonical,err:=encodeMemory(pathTable)
		if err!=nil{return DiscoveredIntegration{},err}
		seed:=34000000+scenarioIndex*10000+9999
		memory,err:=perturbMemory(canonical,seed,memoryNoise)
		if err!=nil{return DiscoveredIntegration{},err}
		memory=rotateGlobalPhase(
			memory,
			math.Mod(0.271*float64(seed+1),2*math.Pi),
		)
		state,err:=fullLatentEncode(memory,mixer)
		if err!=nil{return DiscoveredIntegration{},err}
		operator,ok:=prepared.Ops[scenario.finalGap]
		if !ok{return DiscoveredIntegration{},fmt.Errorf("missing UP52B final depth=%d",scenario.finalGap)}
		state,err=latentMatrixVector(operator,state)
		if err!=nil{return DiscoveredIntegration{},err}
		norm2,err:=NormSquared(state)
		if err!=nil{return DiscoveredIntegration{},err}
		if drift:=math.Abs(norm2-1);drift>maxNormDrift{maxNormDrift=drift}

		decoded,distributions,margin,err:=decodeBreadthTable(
			state,prepared.Observables,prepared.Regressors,prepared.Classifiers,
		)
		if err!=nil{return DiscoveredIntegration{},err}
		if margin<minValueMargin{minValueMargin=margin}
		if decoded==trueTable{finalCorrect++}

		relationInput,err:=relationFeatures(
			distributions[scenario.queryA],
			distributions[scenario.queryB],
		)
		if err!=nil{return DiscoveredIntegration{},err}
		relationProbabilities,err:=relationHead.probabilities(relationInput)
		if err!=nil{return DiscoveredIntegration{},err}
		gotRelation,relationMargin,err:=classAndMargin(relationProbabilities)
		if err!=nil{return DiscoveredIntegration{},err}
		if relationMargin<minRelationMargin{minRelationMargin=relationMargin}
		wantRelation,err:=memoryRelation(trueTable,scenario.queryA,scenario.queryB)
		if err!=nil{return DiscoveredIntegration{},err}
		if gotRelation==wantRelation{relationCorrect++}
	}

	return DiscoveredIntegration{
		Scenarios:scenarios,
		WritesPerScenario:writes,
		CommitDecodeAccuracy:float64(commitCorrect)/float64(commitTotal),
		ExactFinalTableAccuracy:float64(finalCorrect)/float64(scenarios),
		RelationalQueryAccuracy:float64(relationCorrect)/float64(scenarios),
		MinValueMargin:minValueMargin,
		MinRelationMargin:minRelationMargin,
		MaxNormDrift:maxNormDrift,
	},nil
}

func RunUP52B()(UP52BBoundaryWriteResult,error){
	noises:=[]float64{0.065,0.07}
	writeLevels:=[]int{1,2,4,8,16,32,64}
	trainDepths:=[]int{0}
	heldDepths:=[]int{32,128,512,1024}
	allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer()
	trainTables:=fullObserverTablePool(true)
	heldTables:=fullObserverTablePool(false)
	selectedOffsets,err:=up50bFusedOffsets()
	if err!=nil{return UP52BBoundaryWriteResult{},err}
	fullOffsets:=[]float64{0,0,0,0,0,0}

	result:=UP52BBoundaryWriteResult{
		Schema:UP52BBoundaryWriteSchema,
		Experiment:"UP-52B-boundary-write-chain",
		SourceUP51BSeal:"3452be7012cb385dfc32941f89a5d1d8c35b4517",
		FrozenPair:[]int{0,5},
		BoundaryNoises:append([]float64(nil),noises...),
		WriteLevels:append([]int(nil),writeLevels...),
		HeldDepths:append([]int(nil),heldDepths...),
		SelectionRun:false,
	}

	for _,noise:=range noises{
		selectedPrepared,err:=up52bPrepareArm(
			"selected",selectedOffsets,mixer,
			trainTables,heldTables,trainDepths,heldDepths,allDepths,noise,
		)
		if err!=nil{return UP52BBoundaryWriteResult{},err}
		fullPrepared,err:=up52bPrepareArm(
			"full",fullOffsets,mixer,
			trainTables,heldTables,trainDepths,heldDepths,allDepths,noise,
		)
		if err!=nil{return UP52BBoundaryWriteResult{},err}

		series:=UP52BNoiseSeries{
			MemoryNoise:noise,
			StaticSelectedHeld:selectedPrepared.Static.HeldOutAccuracy,
			StaticFullHeld:fullPrepared.Static.HeldOutAccuracy,
		}
		for _,writes:=range writeLevels{
			selectedIntegration,err:=runUP52BIntegration(
				mixer,selectedPrepared,heldTables,heldDepths,noise,writes,
			)
			if err!=nil{return UP52BBoundaryWriteResult{},err}
			fullIntegration,err:=runUP52BIntegration(
				mixer,fullPrepared,heldTables,heldDepths,noise,writes,
			)
			if err!=nil{return UP52BBoundaryWriteResult{},err}

			selectedArm:=MultiplicityDoseArm{
				Name:"up52b_selected",
				Static:selectedPrepared.Static,
				Integration:selectedIntegration,
			}
			fullArm:=MultiplicityDoseArm{
				Name:"up52b_full",
				Static:fullPrepared.Static,
				Integration:fullIntegration,
			}
			gate:=up49bHardGate(selectedArm,fullArm)
			point:=UP52BWritePoint{
				WritesPerScenario:writes,
				Selected:selectedIntegration,
				FullControl:fullIntegration,
				HeldDelta:fullPrepared.Static.HeldOutAccuracy-selectedPrepared.Static.HeldOutAccuracy,
				CommitDelta:fullIntegration.CommitDecodeAccuracy-selectedIntegration.CommitDecodeAccuracy,
				HardGate:gate,
			}
			series.Points=append(series.Points,point)
			if gate{
				series.LastPassingWrites=writes
			}else if series.FirstFailingWrites==0{
				series.FirstFailingWrites=writes
			}
		}
		result.Series=append(result.Series,series)
	}
	return result,nil
}
