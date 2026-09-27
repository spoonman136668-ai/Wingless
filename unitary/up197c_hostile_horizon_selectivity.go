package unitary

const UP197CHorizonSelectivitySchema="wingless.up197c-hostile-horizon-selectivity.v1"

type UP197CMetric struct{
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineFailures int `json:"baseline_failures"`
	BaselineSurvivors int `json:"baseline_survivors"`
	TreatedFailures int `json:"treated_failures"`
	PreventedFailures int `json:"prevented_failures"`
	ActionsTaken int `json:"actions_taken"`
	ApparentUnnecessaryActionArms int `json:"apparent_unnecessary_action_arms"`
	NearFutureActionArms int `json:"near_future_action_arms"`
	TrueUnnecessaryActionArms int `json:"true_unnecessary_action_arms"`
	BaselineSurvivorsBeyondGrace int `json:"baseline_survivors_beyond_grace"`
	PreventedPerTrueUnnecessary float64 `json:"prevented_per_true_unnecessary"`
	TrueUnnecessaryFractionAmongBeyondGraceSurvivors float64 `json:"true_unnecessary_fraction_among_beyond_grace_survivors"`
}
type UP197CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP196CSeal string `json:"source_up196c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	GraceWrites int `json:"grace_writes"`
	Horizons []int `json:"horizons"`
	CohortTopology string `json:"cohort_topology"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP197CMetric `json:"metrics"`
}
func up197cRun(x0 *up81cAging,endangered,horizon int)(lossStep,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			h:=up161cAdversarial(&m,endangered)
			if h<=2{m.query(endangered);actions++;committed=true;since=0}
		}else if actions<8{
			since++
			if since>=4{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1018000+step,step,"hostile_shield"){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP197C()(UP197CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	horizons:=[]int{8,12,16,20,24,28,32,36,40,44,48,52,56}
	res:=UP197CResult{Schema:UP197CHorizonSelectivitySchema,Experiment:"UP-197C-hostile-horizon-selectivity",SourceUP196CSeal:"dd33632854139ef0ff360ecbcbc0802e3d059977",Policy:"hostile_shield",Cadence:2,SpacingIntervals:4,ActionCap:8,GraceWrites:8,Horizons:horizons,CohortTopology:"contiguous_quartets",CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,horizon:=range horizons{
		m:=UP197CMetric{Horizon:horizon}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,2,"hostile_shield")
			treated,actions:=up197cRun(x,e,horizon);m.ActionsTaken+=actions
			if base<=horizon{m.BaselineFailures++}else{m.BaselineSurvivors++}
			if treated<=horizon{m.TreatedFailures++}
			if base<=horizon&&treated>horizon{m.PreventedFailures++}
			if base>horizon+8{m.BaselineSurvivorsBeyondGrace++}
			if base>horizon&&actions>0{
				m.ApparentUnnecessaryActionArms++
				if base<=horizon+8{m.NearFutureActionArms++}else{m.TrueUnnecessaryActionArms++}
			}
		}}
		if m.TrueUnnecessaryActionArms>0{m.PreventedPerTrueUnnecessary=float64(m.PreventedFailures)/float64(m.TrueUnnecessaryActionArms)}
		if m.BaselineSurvivorsBeyondGrace>0{m.TrueUnnecessaryFractionAmongBeyondGraceSurvivors=float64(m.TrueUnnecessaryActionArms)/float64(m.BaselineSurvivorsBeyondGrace)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
