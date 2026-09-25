package unitary

const UP67BMixedDepthConfidenceCarrySchema = "wingless.up67b-mixed-depth-confidence-carry.v1"

type UP67BMixedDepthConfidenceCarryResult struct {
	Schema            string              `json:"schema"`
	Experiment        string              `json:"experiment"`
	SourceUP66BSeal   string              `json:"source_up66b_seal"`
	ScheduleBases     []int               `json:"schedule_bases"`
	NoiseLevels       []float64           `json:"noise_levels"`
	DepthPattern      []int               `json:"depth_pattern"`
	WritesPerScenario int                 `json:"writes_per_scenario"`
	Thresholds        []float64           `json:"thresholds"`
	OracleUsed        bool                `json:"oracle_used"`
	TrainingChanged   bool                `json:"training_changed"`
	ThresholdSelected bool                `json:"threshold_selected"`
	Metrics           []UP57BPolicyMetric `json:"metrics"`
}

func RunUP67B()(UP67BMixedDepthConfidenceCarryResult,error){
	schedules:=[]int{102000000,103000000}
	noises:=[]float64{0.085,0.09}
	depthPattern:=[]int{32,128,512,1024}
	thresholds:=[]float64{0.25,0.5,0.75}
	const writes=512
	trainDepths:=[]int{0}
	allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer()
	trainTables:=fullObserverTablePool(true)
	heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets()
	if err!=nil{return UP67BMixedDepthConfidenceCarryResult{},err}
	result:=UP67BMixedDepthConfidenceCarryResult{
		Schema:UP67BMixedDepthConfidenceCarrySchema,
		Experiment:"UP-67B-mixed-depth-confidence-carry",
		SourceUP66BSeal:"4d4c814b04f69bcf5128f0d1a12729f62ae8bd5e",
		ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),
		DepthPattern:append([]int(nil),depthPattern...),
		WritesPerScenario:writes,
		Thresholds:append([]float64(nil),thresholds...),
		OracleUsed:false,TrainingChanged:false,ThresholdSelected:false,
	}
	for _,noise:=range noises{
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,depthPattern,allDepths,noise)
		if err!=nil{return UP67BMixedDepthConfidenceCarryResult{},err}
		for _,seedBase:=range schedules{
			base,err:=runUP57BPolicy(mixer,prepared,heldTables,depthPattern,noise,writes,seedBase,0,false)
			if err!=nil{return UP67BMixedDepthConfidenceCarryResult{},err}
			result.Metrics=append(result.Metrics,base)
			for _,threshold:=range thresholds{
				m,err:=runUP57BPolicy(mixer,prepared,heldTables,depthPattern,noise,writes,seedBase,threshold,true)
				if err!=nil{return UP67BMixedDepthConfidenceCarryResult{},err}
				result.Metrics=append(result.Metrics,m)
			}
		}
	}
	return result,nil
}
