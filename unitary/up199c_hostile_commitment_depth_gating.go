package unitary

const UP199CCommitmentDepthSchema="wingless.up199c-hostile-commitment-depth-gating.v1"

type UP199CMetric struct{
	CommitmentDepth int `json:"commitment_depth"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineFailures int `json:"baseline_failures"`
	TreatedFailures int `json:"treated_failures"`
	PreventedFailures int `json:"prevented_failures"`
	ActionsTaken int `json:"actions_taken"`
	ArmsWithAction int `json:"arms_with_action"`
	AcceleratedFailures int `json:"accelerated_failures"`
	SkippedOpportunities int `json:"skipped_opportunities"`
}
type UP199CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP198CSeal string `json:"source_up198c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	CommitmentDepths []int `json:"commitment_depths"`
	Horizons []int `json:"horizons"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Metrics []UP199CMetric `json:"metrics"`
}
func up199cRun(x0 *up81cAging,endangered,horizon,depth int)(lossStep,actions,skipped int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			h:=up161cAdversarial(&m,endangered)
			if h<=2{m.query(endangered);actions++;committed=true;since=0}
		}else if actions<8{
			since++
			if since>=4{
				if actions<depth || up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++
				}else{skipped++}
				since=0
			}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1020000+step,step,"hostile_shield"){return step,actions,skipped}
		}
	}
	return horizon+1,actions,skipped
}
func RunUP199C()(UP199CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	depths:=[]int{1,2,4,8};horizons:=[]int{16,24,32,40,48,56}
	res:=UP199CResult{Schema:UP199CCommitmentDepthSchema,Experiment:"UP-199C-hostile-commitment-depth-gating",SourceUP198CSeal:"c8eb2c79ca5fe09371799e9c29031e3051b7c11b",Policy:"hostile_shield",Cadence:2,SpacingIntervals:4,ActionCap:8,CommitmentDepths:depths,Horizons:horizons,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	for _,depth:=range depths{for _,horizon:=range horizons{
		m:=UP199CMetric{CommitmentDepth:depth,Horizon:horizon}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,2,"hostile_shield")
			treated,actions,skipped:=up199cRun(x,e,horizon,depth)
			if base<=horizon{m.BaselineFailures++}
			if treated<=horizon{m.TreatedFailures++}
			if base<=horizon&&treated>horizon{m.PreventedFailures++}
			if treated<base{m.AcceleratedFailures++}
			if actions>0{m.ArmsWithAction++}
			m.ActionsTaken+=actions;m.SkippedOpportunities+=skipped
		}}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
