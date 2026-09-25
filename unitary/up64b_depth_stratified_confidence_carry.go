package unitary

const UP64BDepthStratifiedConfidenceCarrySchema = "wingless.up64b-depth-stratified-confidence-carry.v1"

type UP64BDepthPolicyMetric struct {
	Depth  int               `json:"depth"`
	Metric UP57BPolicyMetric `json:"metric"`
}

type UP64BDepthStratifiedConfidenceCarryResult struct {
	Schema            string                   `json:"schema"`
	Experiment        string                   `json:"experiment"`
	SourceUP63BSeal   string                   `json:"source_up63b_seal"`
	ScheduleBases     []int                    `json:"schedule_bases"`
	NoiseLevels       []float64                `json:"noise_levels"`
	DepthLevels       []int                    `json:"depth_levels"`
	WritesPerScenario int                      `json:"writes_per_scenario"`
	Thresholds        []float64                `json:"thresholds"`
	OracleUsed        bool                     `json:"oracle_used"`
	TrainingChanged   bool                     `json:"training_changed"`
	ThresholdSelected bool                     `json:"threshold_selected"`
	Metrics           []UP64BDepthPolicyMetric `json:"metrics"`
}

func RunUP64B()(UP64BDepthStratifiedConfidenceCarryResult,error){
	schedules:=[]int{96000000,97000000};noises:=[]float64{0.085,0.09};depths:=[]int{32,128,512,1024};thresholds:=[]float64{0.25,0.5,0.75};const writes=256
	trainDepths:=[]int{0};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP64BDepthStratifiedConfidenceCarryResult{},err}
	result:=UP64BDepthStratifiedConfidenceCarryResult{Schema:UP64BDepthStratifiedConfidenceCarrySchema,Experiment:"UP-64B-depth-stratified-confidence-carry",SourceUP63BSeal:"e0b9d3ed351bdd20ee336655c0daf5dce0e9ea55",ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),DepthLevels:append([]int(nil),depths...),WritesPerScenario:writes,Thresholds:append([]float64(nil),thresholds...),OracleUsed:false,TrainingChanged:false,ThresholdSelected:false}
	for _,noise:=range noises {
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,depths,allDepths,noise);if err!=nil{return UP64BDepthStratifiedConfidenceCarryResult{},err}
		for _,depth:=range depths {
			fixedDepth:=[]int{depth}
			for _,seedBase:=range schedules {
				base,err:=runUP57BPolicy(mixer,prepared,heldTables,fixedDepth,noise,writes,seedBase,0,false);if err!=nil{return UP64BDepthStratifiedConfidenceCarryResult{},err}
				result.Metrics=append(result.Metrics,UP64BDepthPolicyMetric{Depth:depth,Metric:base})
				for _,threshold:=range thresholds {
					m,err:=runUP57BPolicy(mixer,prepared,heldTables,fixedDepth,noise,writes,seedBase,threshold,true);if err!=nil{return UP64BDepthStratifiedConfidenceCarryResult{},err}
					result.Metrics=append(result.Metrics,UP64BDepthPolicyMetric{Depth:depth,Metric:m})
				}
			}
		}
	}
	return result,nil
}
