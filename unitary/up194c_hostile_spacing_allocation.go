package unitary

const UP194CHostileSpacingSchema="wingless.up194c-hostile-spacing-allocation.v1"

type UP194CMetric struct{
	Cadence int `json:"cadence"`
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
type UP194CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP193CSeal string `json:"source_up193c_seal"`
	Policy string `json:"policy"`
	ActionCap int `json:"action_cap"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	AdaptiveSpacingUsed bool `json:"adaptive_spacing_used"`
	Metrics []UP194CMetric `json:"metrics"`
}
func up194cRun(x0 *up81cAging,endangered,cadence,spacing,cap int)(lossStep,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=64;start+=cadence{
		if !committed{
			h:=up161cAdversarial(&m,endangered)
			if h<=cadence{
				m.query(endangered);actions++;committed=true;since=0
			}
		}else if actions<cap{
			since++
			if since>=spacing{
				m.query(endangered);actions++;since=0
			}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1017000+step,step,"hostile_shield"){return step,actions}
		}
	}
	return 65,actions
}
func RunUP194C()(UP194CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{cad,spacing int}
	schedules:=[]sched{{2,1},{2,2},{2,3},{2,4},{4,1}}
	res:=UP194CResult{Schema:UP194CHostileSpacingSchema,Experiment:"UP-194C-hostile-spacing-allocation",SourceUP193CSeal:"8fe1c3cdb7a645463f2654a985c0147fee152a71",Policy:"hostile_shield",ActionCap:16,CounterfactualOnly:true,LiveActivation:false,WarningThresholdChanged:false,AdaptiveSpacingUsed:false}
	for _,sc:=range schedules{
		m:=UP194CMetric{Cadence:sc.cad,SpacingIntervals:sc.spacing,ActionCap:16};sumExt,n:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,sc.cad,"hostile_shield")
			treated,actions:=up194cRun(x,e,sc.cad,sc.spacing,16);m.ActionsTaken+=actions
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
	}
	return res,nil
}
