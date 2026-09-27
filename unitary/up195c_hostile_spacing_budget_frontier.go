package unitary

const UP195CSpacingBudgetSchema="wingless.up195c-hostile-spacing-budget-frontier.v1"

type UP195CMetric struct{
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	PreventedLosses int `json:"prevented_losses"`
	ActionsTaken int `json:"actions_taken"`
	AcceleratedLosses int `json:"accelerated_losses"`
	MeanLossStepExtensionAmongFailures float64 `json:"mean_loss_step_extension_among_failures"`
	MaxLossStepExtensionAmongFailures int `json:"max_loss_step_extension_among_failures"`
}
type UP195CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP194CSeal string `json:"source_up194c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Spacings []int `json:"spacings"`
	ActionCaps []int `json:"action_caps"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	AdaptiveSpacingUsed bool `json:"adaptive_spacing_used"`
	AdaptiveCapUsed bool `json:"adaptive_cap_used"`
	Metrics []UP195CMetric `json:"metrics"`
}
func RunUP195C()(UP195CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	spacings:=[]int{2,3,4};caps:=[]int{4,6,8,10,12,16}
	res:=UP195CResult{Schema:UP195CSpacingBudgetSchema,Experiment:"UP-195C-hostile-spacing-budget-frontier",SourceUP194CSeal:"5a25ed974a31ca3d34a4f2b8a55133db800ff361",Policy:"hostile_shield",Cadence:2,Spacings:spacings,ActionCaps:caps,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,AdaptiveSpacingUsed:false,AdaptiveCapUsed:false}
	for _,spacing:=range spacings{for _,cap:=range caps{
		m:=UP195CMetric{SpacingIntervals:spacing,ActionCap:cap};sumExt,n:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,2,"hostile_shield")
			treated,actions:=up194cRun(x,e,2,spacing,cap);m.ActionsTaken+=actions
			if base<=64{m.BaselineLosses++}
			if base<=64&&treated>64{m.PreventedLosses++}
			if treated<base{m.AcceleratedLosses++}
			if treated<=64{
				ext:=treated-base;sumExt+=ext;n++
				if ext>m.MaxLossStepExtensionAmongFailures{m.MaxLossStepExtensionAmongFailures=ext}
			}
		}}
		if n>0{m.MeanLossStepExtensionAmongFailures=float64(sumExt)/float64(n)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
