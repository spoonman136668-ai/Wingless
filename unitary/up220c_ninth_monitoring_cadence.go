package unitary

const UP220CNinthCadenceSchema="wingless.up220c-ninth-monitoring-cadence.v1"

type UP220CMetric struct{
	Horizon int `json:"horizon"`
	MonitoringMode string `json:"monitoring_mode"`
	Arms int `json:"arms"`
	Cap8Losses int `json:"cap8_losses"`
	Cap9Losses int `json:"cap9_losses"`
	NinthActions int `json:"ninth_actions"`
	NinthRescues int `json:"ninth_rescues"`
	UnnecessaryNinth int `json:"unnecessary_ninth"`
	NecessaryNinthAttempts int `json:"necessary_ninth_attempts"`
	IneffectiveNinth int `json:"ineffective_ninth"`
	RescuePerUnnecessary float64 `json:"rescue_per_unnecessary"`
}
type UP220CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP219CSeal string `json:"source_up219c_seal"`
	CriticalHorizon int `json:"critical_horizon"`
	MaxActions int `json:"max_actions"`
	Horizons []int `json:"horizons"`
	MonitoringModes []string `json:"monitoring_modes"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveCadenceSelectionUsed bool `json:"adaptive_cadence_selection_used"`
	ScheduleRetuningUsed bool `json:"schedule_retuning_used"`
	Metrics []UP220CMetric `json:"metrics"`
}
func up220cTreated(x0 *up81cAging,endangered,horizon int,mode string)(loss,actions,ninth int){
	m:=*x0;committed:=false;since:=0;eighthBoundary:=-1
	for start:=1;start<=horizon;start+=2{
		eighthThisBoundary:=false
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;since=0;eighthBoundary=start;eighthThisBoundary=true
			}else if mode=="cadence2"&&actions==8&&eighthBoundary>=0&&start>eighthBoundary&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;ninth++;since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if mode=="every_write"&&actions==8&&eighthBoundary>=0&&step>eighthBoundary{
				if !(eighthThisBoundary&&j==0)&&up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++;ninth++;since=0
				}
			}
			if up165cRealStep(&m,endangered,1043000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){return step,actions,ninth}
		}
	}
	return horizon+1,actions,ninth
}
func RunUP220C()(UP220CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	horizons:=[]int{74,76,78,80};modes:=[]string{"cadence2","every_write"}
	res:=UP220CResult{Schema:UP220CNinthCadenceSchema,Experiment:"UP-220C-ninth-monitoring-cadence",SourceUP219CSeal:"c56186d088c0a58ff44d9e1e5df94c5034aea60f",CriticalHorizon:2,MaxActions:9,Horizons:horizons,MonitoringModes:modes,CounterfactualOnly:true,LiveActivation:false,ThresholdFittingUsed:false,AdaptiveCadenceSelectionUsed:false,ScheduleRetuningUsed:false}
	for _,h:=range horizons{for _,mode:=range modes{
		m:=UP220CMetric{Horizon:h,MonitoringMode:mode}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			cap8,_,_:=up218cTreated(x,e,8,h,"alternating_shield","fixed_offset_refresh",8)
			cap9,_,ninth:=up220cTreated(x,e,h,mode)
			l8:=cap8<=h;l9:=cap9<=h
			if l8{m.Cap8Losses++};if l9{m.Cap9Losses++};m.NinthActions+=ninth
			if l8&&!l9{m.NinthRescues++}
			if !l8&&ninth>0{m.UnnecessaryNinth++}
			if l8&&ninth>0{m.NecessaryNinthAttempts++}
			if l8&&l9&&ninth>0{m.IneffectiveNinth++}
		}}
		if m.UnnecessaryNinth>0{m.RescuePerUnnecessary=float64(m.NinthRescues)/float64(m.UnnecessaryNinth)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
