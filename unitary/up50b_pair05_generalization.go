package unitary

import "fmt"

const UP50BGeneralizationSchema = "wingless.up50b-pair05-generalization.v1"

type UP50BCaseSpec struct {
	Name        string    `json:"name"`
	MemoryNoise float64   `json:"memory_noise"`
	HeldDepths  []int     `json:"held_depths"`
}

type UP50BCase struct {
	Spec             UP50BCaseSpec       `json:"spec"`
	SelectedHard     MultiplicityDoseArm `json:"selected_hard"`
	FullControl      MultiplicityDoseArm `json:"full_control"`
	HeldDelta        float64             `json:"held_delta"`
	CommitDelta      float64             `json:"commit_delta"`
	HardGate         bool                `json:"hard_gate"`
}

type UP50BGeneralizationResult struct {
	Schema              string       `json:"schema"`
	Experiment          string       `json:"experiment"`
	SourceUP49BSeal     string       `json:"source_up49b_seal"`
	SourceStep          int          `json:"source_step"`
	FrozenPair          []int        `json:"frozen_pair"`
	FrozenOffsets       []float64    `json:"frozen_offsets"`
	SelectionRun        bool         `json:"selection_run"`
	Cases               []UP50BCase `json:"cases"`
	AllCasesPass        bool         `json:"all_cases_pass"`
}

func up50bFusedOffsets() ([]float64,error){
	offsets,err:=up48bFrozenOffsets(42)
	if err!=nil{return nil,err}
	center:=0.5*(offsets[0]+offsets[5])
	offsets[0]=center
	offsets[5]=center
	_,capacity,err:=exactOffsetGroups(offsets)
	if err!=nil{return nil,err}
	if capacity<=compositeChannels{return nil,fmt.Errorf("UP50B frozen pair not fused")}
	return offsets,nil
}

func RunUP50B()(UP50BGeneralizationResult,error){
	offsets,err:=up50bFusedOffsets()
	if err!=nil{return UP50BGeneralizationResult{},err}
	mixer:=fullLatentMixer()
	allTrain:=fullObserverTablePool(true)
	trueHeld:=fullObserverTablePool(false)
	trainDepths:=[]int{0}
	specs:=[]UP50BCaseSpec{
		{Name:"low_noise_standard_depths",MemoryNoise:0.03,HeldDepths:[]int{32,128,512,1024}},
		{Name:"high_noise_standard_depths",MemoryNoise:0.08,HeldDepths:[]int{32,128,512,1024}},
		{Name:"low_noise_shifted_depths",MemoryNoise:0.03,HeldDepths:[]int{64,256,1024,2048}},
		{Name:"high_noise_shifted_depths",MemoryNoise:0.08,HeldDepths:[]int{64,256,1024,2048}},
	}
	result:=UP50BGeneralizationResult{
		Schema:UP50BGeneralizationSchema,
		Experiment:"UP-50B-pair05-generalization",
		SourceUP49BSeal:"3b69129d94c97eccc90af320f99dcac7c2450053",
		SourceStep:42,
		FrozenPair:[]int{0,5},
		FrozenOffsets:append([]float64(nil),offsets...),
		SelectionRun:false,
		AllCasesPass:true,
	}
	for _,spec:=range specs{
		allDepths:=append([]int{0},spec.HeldDepths...)
		selected,err:=evaluateMultiplicityDoseArm(
			"up50b_"+spec.Name+"_selected",
			append([]float64(nil),offsets...),
			mixer,allTrain,trueHeld,
			trainDepths,spec.HeldDepths,allDepths,spec.MemoryNoise,
		)
		if err!=nil{return UP50BGeneralizationResult{},err}
		full,err:=evaluateMultiplicityDoseArm(
			"up50b_"+spec.Name+"_full",
			[]float64{0,0,0,0,0,0},
			mixer,allTrain,trueHeld,
			trainDepths,spec.HeldDepths,allDepths,spec.MemoryNoise,
		)
		if err!=nil{return UP50BGeneralizationResult{},err}
		heldDelta:=full.Static.HeldOutAccuracy-selected.Static.HeldOutAccuracy
		commitDelta:=full.Integration.CommitDecodeAccuracy-selected.Integration.CommitDecodeAccuracy
		gate:=up49bHardGate(selected,full)
		if !gate{result.AllCasesPass=false}
		result.Cases=append(result.Cases,UP50BCase{
			Spec:spec,SelectedHard:selected,FullControl:full,
			HeldDelta:heldDelta,CommitDelta:commitDelta,HardGate:gate,
		})
	}
	return result,nil
}
