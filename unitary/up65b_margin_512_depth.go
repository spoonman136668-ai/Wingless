package unitary

const UP65BMargin512DepthSchema = "wingless.up65b-margin-512-depth.v1"

type UP65BMargin512DepthResult struct {
	Schema            string                        `json:"schema"`
	Experiment        string                        `json:"experiment"`
	SourceUP63BSeal   string                        `json:"source_up63b_seal"`
	ScheduleBases     []int                         `json:"schedule_bases"`
	NoiseLevels       []float64                     `json:"noise_levels"`
	DepthLevels       []int                         `json:"depth_levels"`
	WritesPerScenario int                           `json:"writes_per_scenario"`
	BinEdges          []float64                     `json:"bin_edges"`
	OracleUsed        bool                          `json:"oracle_used"`
	PolicyChanged     bool                          `json:"policy_changed"`
	TrainingChanged   bool                          `json:"training_changed"`
	ThresholdSelected bool                          `json:"threshold_selected"`
	Metrics           []UP63BDepthCalibrationMetric `json:"metrics"`
}

func RunUP65B()(UP65BMargin512DepthResult,error){
	schedules:=[]int{98000000,99000000};noises:=[]float64{0.085,0.09};depths:=[]int{32,128,512,1024};const writes=512
	trainDepths:=[]int{0};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP65BMargin512DepthResult{},err}
	result:=UP65BMargin512DepthResult{Schema:UP65BMargin512DepthSchema,Experiment:"UP-65B-margin-512-depth",SourceUP63BSeal:"e0b9d3ed351bdd20ee336655c0daf5dce0e9ea55",ScheduleBases:append([]int(nil),schedules...),NoiseLevels:append([]float64(nil),noises...),DepthLevels:append([]int(nil),depths...),WritesPerScenario:writes,BinEdges:[]float64{0,0.25,0.5,0.75,1},OracleUsed:false,PolicyChanged:false,TrainingChanged:false,ThresholdSelected:false}
	for _,noise:=range noises{
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,depths,allDepths,noise);if err!=nil{return UP65BMargin512DepthResult{},err}
		for _,depth:=range depths{
			fixedDepth:=[]int{depth}
			for _,seedBase:=range schedules{
				m,err:=up61bRunCalibration(mixer,prepared,heldTables,fixedDepth,noise,writes,seedBase);if err!=nil{return UP65BMargin512DepthResult{},err}
				result.Metrics=append(result.Metrics,UP63BDepthCalibrationMetric{ScheduleBase:m.ScheduleBase,MemoryNoise:m.MemoryNoise,WritesPerScenario:m.WritesPerScenario,Depth:depth,Buckets:m.Buckets,MaxNormDrift:m.MaxNormDrift})
			}
		}
	}
	return result,nil
}
