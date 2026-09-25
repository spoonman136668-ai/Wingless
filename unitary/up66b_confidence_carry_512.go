package unitary

const UP66BConfidenceCarry512Schema = "wingless.up66b-confidence-carry-512.v1"

type UP66BDepthPolicyMetric struct {
	Depth  int               `json:"depth"`
	Metric UP57BPolicyMetric `json:"metric"`
}

type UP66BConfidenceCarry512Result struct {
	Schema            string                   `json:"schema"`
	Experiment        string                   `json:"experiment"`
	SourceUP65BSeal   string                   `json:"source_up65b_seal"`
	ScheduleBases     []int                    `json:"schedule_bases"`
	NoiseLevels       []float64                `json:"noise_levels"`
	DepthLevels       []int                    `json:"depth_levels"`
	WritesPerScenario int                      `json:"writes_per_scenario"`
	Thresholds        []float64                `json:"thresholds"`
	OracleUsed        bool                     `json:"oracle_used"`
	TrainingChanged   bool                     `json:"training_changed"`
	ThresholdSelected bool                     `json:"threshold_selected"`
	Metrics           []UP66BDepthPolicyMetric `json:"metrics"`
}

func RunUP66B()(UP66BConfidenceCarry512Result,error){
	schedules:=[]int{100000000,101000000}
	noises:=[]float64{0.085,0.09}
	depths:=[]int{32,128,512,1024}
	thresholds:=[]float64{0.25,0.5,0.75}
	const writes=512
	trainDepths:=[]int{0}
	allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer()
	trainTables:=fullObserverTablePool(true)
	heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets()
	if err!=nil{return UP66BConfidenceCarry512Result{},err}
	result:=UP66BConfidenceCarry512Result{
		Schema:UP66BConfidenceCarry512Schema,Experiment:"UP-66B-confidence-carry-512",
		SourceUP65BSeal:"957cef220597126b2ca9b99872f0a3b1b5a38a5b",
		ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),
		DepthLevels:append([]int(nil),depths...),WritesPerScenario:writes,
		Thresholds:append([]float64(nil),thresholds...),OracleUsed:false,TrainingChanged:false,ThresholdSelected:false,
	}
	for _,noise:=range noises{
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,depths,allDepths,noise)
		if err!=nil{return UP66BConfidenceCarry512Result{},err}
		for _,depth:=range depths{
			fixedDepth:=[]int{depth}
			for _,seedBase:=range schedules{
				base,err:=runUP57BPolicy(mixer,prepared,heldTables,fixedDepth,noise,writes,seedBase,0,false)
				if err!=nil{return UP66BConfidenceCarry512Result{},err}
				result.Metrics=append(result.Metrics,UP66BDepthPolicyMetric{Depth:depth,Metric:base})
				for _,threshold:=range thresholds{
					m,err:=runUP57BPolicy(mixer,prepared,heldTables,fixedDepth,noise,writes,seedBase,threshold,true)
					if err!=nil{return UP66BConfidenceCarry512Result{},err}
					result.Metrics=append(result.Metrics,UP66BDepthPolicyMetric{Depth:depth,Metric:m})
				}
			}
		}
	}
	return result,nil
}
