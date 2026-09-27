package unitary

const UP196CHeldoutCohortSchema="wingless.up196c-hostile-heldout-cohort-transfer.v1"

type UP196CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP195CSeal string `json:"source_up195c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	PrimarySpacing int `json:"primary_spacing"`
	PrimaryCap int `json:"primary_cap"`
	DenseSpacing int `json:"dense_spacing"`
	DenseCap int `json:"dense_cap"`
	CohortTopology string `json:"cohort_topology"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	PrimaryPrevented int `json:"primary_prevented"`
	PrimaryActions int `json:"primary_actions"`
	PrimaryAccelerated int `json:"primary_accelerated"`
	PrimaryMeanLossExtensionAmongFailures float64 `json:"primary_mean_loss_extension_among_failures"`
	PrimaryMaxLossExtensionAmongFailures int `json:"primary_max_loss_extension_among_failures"`
	DensePrevented int `json:"dense_prevented"`
	DenseActions int `json:"dense_actions"`
	DenseAccelerated int `json:"dense_accelerated"`
	DenseMeanLossExtensionAmongFailures float64 `json:"dense_mean_loss_extension_among_failures"`
	DenseMaxLossExtensionAmongFailures int `json:"dense_max_loss_extension_among_failures"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	CohortSpecificTuningUsed bool `json:"cohort_specific_tuning_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
}
func RunUP196C()(UP196CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP196CResult{
		Schema:UP196CHeldoutCohortSchema,Experiment:"UP-196C-hostile-heldout-cohort-transfer",SourceUP195CSeal:"9fa909c1cc8f7503daea221e395e9d563b2a11a1",
		Policy:"hostile_shield",Cadence:2,PrimarySpacing:4,PrimaryCap:8,DenseSpacing:1,DenseCap:16,CohortTopology:"contiguous_quartets",
		CounterfactualOnly:true,LiveActivation:false,CohortSpecificTuningUsed:false,WarningThresholdChanged:false,
	}
	sumP,nP,sumD,nD:=0,0,0,0
	for _,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};res.Arms++
		base:=up169cBaselineLossStep(x,e,2,"hostile_shield")
		primary,pa:=up194cRun(x,e,2,4,8)
		dense,da:=up194cRun(x,e,2,1,16)
		res.PrimaryActions+=pa;res.DenseActions+=da
		if base<=64{res.BaselineLosses++}
		if base<=64&&primary>64{res.PrimaryPrevented++}
		if base<=64&&dense>64{res.DensePrevented++}
		if primary<base{res.PrimaryAccelerated++}
		if dense<base{res.DenseAccelerated++}
		if primary<=64{
			ext:=primary-base;sumP+=ext;nP++
			if ext>res.PrimaryMaxLossExtensionAmongFailures{res.PrimaryMaxLossExtensionAmongFailures=ext}
		}
		if dense<=64{
			ext:=dense-base;sumD+=ext;nD++
			if ext>res.DenseMaxLossExtensionAmongFailures{res.DenseMaxLossExtensionAmongFailures=ext}
		}
	}}
	if nP>0{res.PrimaryMeanLossExtensionAmongFailures=float64(sumP)/float64(nP)}
	if nD>0{res.DenseMeanLossExtensionAmongFailures=float64(sumD)/float64(nD)}
	return res,nil
}
