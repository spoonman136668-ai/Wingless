package unitary

const UP193CHostileCeilingSchema="wingless.up193c-hostile-commitment-ceiling.v1"

type UP193CMetric struct{
	Cadence int `json:"cadence"`
	ActionCap int `json:"action_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	FreshGatedPrevented int `json:"fresh_gated_prevented"`
	FreshGatedActions int `json:"fresh_gated_actions"`
	CommitmentPrevented int `json:"commitment_prevented"`
	CommitmentActions int `json:"commitment_actions"`
	CommitmentAccelerated int `json:"commitment_accelerated"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
	MaxLossStepExtensionAmongFailures int `json:"max_loss_step_extension_among_failures"`
}
type UP193CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP192CSeal string `json:"source_up192c_seal"`
	Policy string `json:"policy"`
	Cadences []int `json:"cadences"`
	ActionCaps []int `json:"action_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	NewActionTypeUsed bool `json:"new_action_type_used"`
	AdaptiveCapUsed bool `json:"adaptive_cap_used"`
	Metrics []UP193CMetric `json:"metrics"`
}
func RunUP193C()(UP193CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};caps:=[]int{8,12,16}
	res:=UP193CResult{Schema:UP193CHostileCeilingSchema,Experiment:"UP-193C-hostile-commitment-ceiling",SourceUP192CSeal:"ec29832fca7bce53d7f95c0471cec373a3d358c3",Policy:"hostile_shield",Cadences:cadences,ActionCaps:caps,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,NewActionTypeUsed:false,AdaptiveCapUsed:false}
	for _,cad:=range cadences{for _,cap:=range caps{
		m:=UP193CMetric{Cadence:cad,ActionCap:cap};sumExt,n:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,"hostile_shield")
			fresh,fa:=up191cRun(x,e,cad,"hostile_shield",cap)
			commit,ca:=up192cRun(x,e,cad,"hostile_shield",cap)
			if base<=64{m.BaselineLosses++}
			if base<=64&&fresh>64{m.FreshGatedPrevented++}
			if base<=64&&commit>64{m.CommitmentPrevented++}
			m.FreshGatedActions+=fa;m.CommitmentActions+=ca
			if commit<base{m.CommitmentAccelerated++}
			if commit<=64{
				ext:=commit-base;sumExt+=ext;n++
				if ext>m.MaxLossStepExtensionAmongFailures{m.MaxLossStepExtensionAmongFailures=ext}
			}
		}}
		if n>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(n)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
