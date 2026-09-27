package unitary

const UP198CFreshGateSchema="wingless.up198c-hostile-fresh-warning-gating.v1"

type UP198CMetric struct{
	Horizon int `json:"horizon"`
	Schedule string `json:"schedule"`
	Arms int `json:"arms"`
	BaselineFailures int `json:"baseline_failures"`
	TreatedFailures int `json:"treated_failures"`
	PreventedFailures int `json:"prevented_failures"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	AcceleratedFailures int `json:"accelerated_failures"`
	SkippedOpportunities int `json:"skipped_opportunities"`
}
type UP198CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP197CSeal string `json:"source_up197c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizons []int `json:"horizons"`
	Schedules []string `json:"schedules"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP198CMetric `json:"metrics"`
}
func up198cRun(x0 *up81cAging,endangered,horizon int,gated bool)(lossStep,actions,skipped int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			h:=up161cAdversarial(&m,endangered)
			if h<=2{m.query(endangered);actions++;committed=true;since=0}
		}else if actions<8{
			since++
			if since>=4{
				if !gated || up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++
				}else{skipped++}
				since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1019000+step,step,"hostile_shield"){return step,actions,skipped}
		}
	}
	return horizon+1,actions,skipped
}
func RunUP198C()(UP198CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	horizons:=[]int{16,24,32,40,48,56};schedules:=[]string{"committed","fresh_warning_gated"}
	res:=UP198CResult{Schema:UP198CFreshGateSchema,Experiment:"UP-198C-hostile-fresh-warning-gating",SourceUP197CSeal:"d5dcbca3fde66f0079b73ecd9831267b18ec6d1e",Policy:"hostile_shield",Cadence:2,SpacingIntervals:4,ActionCap:8,Horizons:horizons,Schedules:schedules,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,horizon:=range horizons{
		for _,sched:=range schedules{
			m:=UP198CMetric{Horizon:horizon,Schedule:sched}
			gated:=sched=="fresh_warning_gated"
			for _,c:=range cohorts{for _,hand:=range hands{
				x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
				base:=up169cBaselineLossStep(x,e,2,"hostile_shield")
				treated,actions,skipped:=up198cRun(x,e,horizon,gated)
				if base<=horizon{m.BaselineFailures++}
				if treated<=horizon{m.TreatedFailures++}
				if base<=horizon&&treated>horizon{m.PreventedFailures++}
				if treated<base{m.AcceleratedFailures++}
				if actions>0{m.ArmsWithAction++}
				m.ActionsTaken+=actions;m.SkippedOpportunities+=skipped
			}}
			res.Metrics=append(res.Metrics,m)
		}
	}
	return res,nil
}
