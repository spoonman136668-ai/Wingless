package unitary

import "sort"

const UP190BNativeRegimeSchema="wingless.up190b-native-regime-scan.v1"

type UP190BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
	MaxAbsoluteMassRatioError float64 `json:"max_absolute_mass_ratio_error"`
}
type UP190BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP189BSeal string `json:"source_up189b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	ScanPhases []int `json:"scan_phases"`
	ContextThreshold int `json:"context_threshold"`
	ContextAFactor float64 `json:"context_a_factor"`
	ContextBFactor float64 `json:"context_b_factor"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	UniqueNativeCorrectCounts []int `json:"unique_native_correct_counts"`
	Metrics []UP188BMetric `json:"metrics"`
	Summaries []UP190BSummary `json:"summaries"`
}
func RunUP190B()(UP190BResult,error){
	cal:=[]int{26,27,28,29,30}
	scan:=make([]int,0,16);for p:=43;p<=58;p++{scan=append(scan,p)}
	pooled:=up184bBuild("pooled_26_30",cal)
	f:=map[string]float64{"context_A":1.0103904235372925,"context_B":0.8983526543771945}
	res:=UP190BResult{Schema:UP190BNativeRegimeSchema,Experiment:"UP-190B-native-regime-scan",SourceUP189BSeal:"62429a84c67c974773efe5ddb51d7324be880507",CalibrationPhases:cal,ScanPhases:scan,ContextThreshold:621,ContextAFactor:f["context_A"],ContextBFactor:f["context_B"],PhaseInputUsed:false,ParityInputUsed:false,HeldoutFittingUsed:false,MaintenanceTriggered:false}
	seen:=map[int]bool{};sum:=map[string]float64{};maxv:=map[string]float64{}
	for _,phase:=range scan{
		_,nc:=up187bContext(phase);seen[nc]=true
		a:=up188bEval(pooled,f,phase,false);b:=up188bEval(pooled,f,phase,true)
		res.Metrics=append(res.Metrics,a,b)
		for _,m:=range []UP188BMetric{a,b}{sum[m.Model]+=m.AbsoluteMassRatioError;if m.AbsoluteMassRatioError>maxv[m.Model]{maxv[m.Model]=m.AbsoluteMassRatioError}}
	}
	for n:=range seen{res.UniqueNativeCorrectCounts=append(res.UniqueNativeCorrectCounts,n)}
	sort.Ints(res.UniqueNativeCorrectCounts)
	for _,name:=range []string{"pooled_26_30","pooled_plus_native_mass_correction"}{res.Summaries=append(res.Summaries,UP190BSummary{Model:name,MeanAbsoluteMassRatioError:sum[name]/float64(len(scan)),MaxAbsoluteMassRatioError:maxv[name]})}
	return res,nil
}
