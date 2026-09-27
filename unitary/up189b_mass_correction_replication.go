package unitary

const UP189BReplicationSchema="wingless.up189b-mass-correction-replication.v1"

type UP189BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
	MaxAbsoluteMassRatioError float64 `json:"max_absolute_mass_ratio_error"`
}
type UP189BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP188BSeal string `json:"source_up188b_seal"`
	CalibrationPhases []int `json:"calibration_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ContextAFactor float64 `json:"context_a_factor"`
	ContextBFactor float64 `json:"context_b_factor"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP188BMetric `json:"metrics"`
	Summaries []UP189BSummary `json:"summaries"`
}
func RunUP189B()(UP189BResult,error){
	cal:=[]int{26,27,28,29,30};eval:=[]int{37,38,39,40,41,42}
	pooled:=up184bBuild("pooled_26_30",cal)
	f:=map[string]float64{"context_A":1.0103904235372925,"context_B":0.8983526543771945}
	res:=UP189BResult{Schema:UP189BReplicationSchema,Experiment:"UP-189B-mass-correction-replication",SourceUP188BSeal:"ec79b82c146b0c75634f0abec190434d1f58cebd",CalibrationPhases:cal,EvaluationPhases:eval,ContextAFactor:f["context_A"],ContextBFactor:f["context_B"],PhaseInputUsed:false,ParityInputUsed:false,HeldoutFittingUsed:false,MaintenanceTriggered:false}
	sum:=map[string]float64{};maxv:=map[string]float64{}
	for _,phase:=range eval{
		a:=up188bEval(pooled,f,phase,false);b:=up188bEval(pooled,f,phase,true)
		res.Metrics=append(res.Metrics,a,b)
		for _,m:=range []UP188BMetric{a,b}{sum[m.Model]+=m.AbsoluteMassRatioError;if m.AbsoluteMassRatioError>maxv[m.Model]{maxv[m.Model]=m.AbsoluteMassRatioError}}
	}
	for _,name:=range []string{"pooled_26_30","pooled_plus_native_mass_correction"}{res.Summaries=append(res.Summaries,UP189BSummary{Model:name,MeanAbsoluteMassRatioError:sum[name]/float64(len(eval)),MaxAbsoluteMassRatioError:maxv[name]})}
	return res,nil
}
