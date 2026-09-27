package unitary

const UP185BParitySchema="wingless.up185b-parity-diagnostic.v1"

type UP185BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP184BSeal string `json:"source_up184b_seal"`
	EvaluationPhases []int `json:"evaluation_phases"`
	ParityIsDiagnosticOnly bool `json:"parity_is_diagnostic_only"`
	NativeDeploymentFeature bool `json:"native_deployment_feature"`
	EvaluationFittingUsed bool `json:"evaluation_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP184BMetric `json:"metrics"`
	Summaries []UP184BModelSummary `json:"summaries"`
}
func RunUP185B()(UP185BResult,error){
	pooled:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	even:=up184bBuild("parity_diagnostic",[]int{26,28,30})
	odd:=up184bBuild("parity_diagnostic",[]int{27,29})
	phases:=[]int{31,32,33}
	res:=UP185BResult{Schema:UP185BParitySchema,Experiment:"UP-185B-parity-diagnostic",SourceUP184BSeal:"d0f944790268afd4e7a6d15380710e123e818587",EvaluationPhases:phases,ParityIsDiagnosticOnly:true,NativeDeploymentFeature:false,EvaluationFittingUsed:false,MaintenanceTriggered:false}
	sums:=map[string]float64{}
	for _,phase:=range phases{
		x:=up184bEval(pooled,phase);res.Metrics=append(res.Metrics,x);sums["pooled_26_30"]+=x.AbsoluteMassRatioError
		m:=odd;if phase%2==0{m=even}
		y:=up184bEval(m,phase);y.Model="parity_diagnostic";res.Metrics=append(res.Metrics,y);sums["parity_diagnostic"]+=y.AbsoluteMassRatioError
	}
	for _,name:=range []string{"pooled_26_30","parity_diagnostic"}{res.Summaries=append(res.Summaries,UP184BModelSummary{Model:name,MeanAbsoluteMassRatioError:sums[name]/3.0})}
	return res,nil
}
