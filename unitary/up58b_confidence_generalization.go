package unitary

const UP58BConfidenceGeneralizationSchema = "wingless.up58b-confidence-generalization.v1"

type UP58BConfidenceGeneralizationResult struct {
	Schema            string              `json:"schema"`
	Experiment        string              `json:"experiment"`
	SourceUP57BSeal   string              `json:"source_up57b_seal"`
	ScheduleBases     []int               `json:"schedule_bases"`
	NoiseLevels       []float64           `json:"noise_levels"`
	WriteLevels       []int               `json:"write_levels"`
	Thresholds        []float64           `json:"thresholds"`
	OracleUsed        bool                `json:"oracle_used"`
	TrainingChanged   bool                `json:"training_changed"`
	ThresholdSelected bool                `json:"threshold_selected"`
	Metrics           []UP57BPolicyMetric `json:"metrics"`
}

func RunUP58B()(UP58BConfidenceGeneralizationResult,error){
	schedules:=[]int{70000000,71000000,72000000}
	noises:=[]float64{0.06,0.07,0.075}
	writeLevels:=[]int{32,128}
	thresholds:=[]float64{0.25,0.5,0.75}
	trainDepths:=[]int{0};heldDepths:=[]int{32,128,512,1024};allDepths:=[]int{0,32,128,512,1024}
	mixer:=fullLatentMixer();trainTables:=fullObserverTablePool(true);heldTables:=fullObserverTablePool(false)
	offsets,err:=up50bFusedOffsets();if err!=nil{return UP58BConfidenceGeneralizationResult{},err}
	result:=UP58BConfidenceGeneralizationResult{
		Schema:UP58BConfidenceGeneralizationSchema,
		Experiment:"UP-58B-confidence-generalization",
		SourceUP57BSeal:"35b4000c1b25372419ad1f740e19146694c84952",
		ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),
		WriteLevels:append([]int(nil),writeLevels...),
		Thresholds:append([]float64(nil),thresholds...),
		OracleUsed:false,TrainingChanged:false,ThresholdSelected:false,
	}
	for _,noise:=range noises {
		prepared,err:=up52bPrepareArm("selected",offsets,mixer,trainTables,heldTables,trainDepths,heldDepths,allDepths,noise)
		if err!=nil{return UP58BConfidenceGeneralizationResult{},err}
		for _,writes:=range writeLevels {
			for _,seedBase:=range schedules {
				base,err:=runUP57BPolicy(mixer,prepared,heldTables,heldDepths,noise,writes,seedBase,0,false)
				if err!=nil{return UP58BConfidenceGeneralizationResult{},err}
				result.Metrics=append(result.Metrics,base)
				for _,threshold:=range thresholds {
					m,err:=runUP57BPolicy(mixer,prepared,heldTables,heldDepths,noise,writes,seedBase,threshold,true)
					if err!=nil{return UP58BConfidenceGeneralizationResult{},err}
					result.Metrics=append(result.Metrics,m)
				}
			}
		}
	}
	return result,nil
}
