package unitary

const UP190COracleCeilingSchema="wingless.up190c-two-action-oracle-ceiling.v1"

type UP190CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	FrozenRulePrevented int `json:"frozen_rule_prevented"`
	OraclePreventable int `json:"oracle_preventable"`
	OraclePreventableFrozenMissed int `json:"oracle_preventable_frozen_missed"`
	NotPreventableWithTwoActions int `json:"not_preventable_with_two_actions"`
	MeanBestOracleLossExtension float64 `json:"mean_best_oracle_loss_extension"`
	MaxBestOracleLossExtension int `json:"max_best_oracle_loss_extension"`
}
type UP190CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP189CSeal string `json:"source_up189c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	MaxActions int `json:"max_actions"`
	ActionGrid string `json:"action_grid"`
	OfflineOracleOnly bool `json:"offline_oracle_only"`
	LiveActivation bool `json:"live_activation"`
	ActionBudgetChanged bool `json:"action_budget_changed"`
	ArbitraryWriteTimingUsed bool `json:"arbitrary_write_timing_used"`
	Metrics []UP190CMetric `json:"metrics"`
}
func up190cScheduledLoss(x0 *up81cAging,endangered,cadence int,policy string,a1,a2 int)(int,int){
	m:=*x0;actions:=0
	for start:=1;start<=64;start+=cadence{
		if start==a1||start==a2{m.query(endangered);actions++}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,990000+step,step,policy){return step,actions}
		}
	}
	return 65,actions
}
func up190cBestTwoAction(x0 *up81cAging,endangered,cadence int,policy string)(best int){
	best=up169cBaselineLossStep(x0,endangered,cadence,policy)
	starts:=[]int{}
	for s:=1;s<=64;s+=cadence{starts=append(starts,s)}
	for _,s:=range starts{
		v,_:=up190cScheduledLoss(x0,endangered,cadence,policy,s,-1)
		if v>best{best=v}
	}
	for i:=0;i<len(starts);i++{
		for j:=i+1;j<len(starts);j++{
			v,_:=up190cScheduledLoss(x0,endangered,cadence,policy,starts[i],starts[j])
			if v>best{best=v}
			if best>64{return best}
		}
	}
	return best
}
func RunUP190C()(UP190CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	cadences:=[]int{2,4}
	res:=UP190CResult{
		Schema:UP190COracleCeilingSchema,Experiment:"UP-190C-two-action-oracle-ceiling",
		SourceUP189CSeal:"5a8f28303e12b2adaee3d588b0551211b87b08f5",
		Policies:policies,Cadences:cadences,MaxActions:2,ActionGrid:"monitoring_interval_starts",
		OfflineOracleOnly:true,LiveActivation:false,ActionBudgetChanged:false,ArbitraryWriteTimingUsed:false,
	}
	for _,policy:=range policies{for _,cad:=range cadences{
		m:=UP190CMetric{Policy:policy,Cadence:cad};sumExt:=0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			base:=up169cBaselineLossStep(x,e,cad,policy)
			frozen,_,_:=up189cRun(x,e,cad,policy)
			best:=up190cBestTwoAction(x,e,cad,policy)
			if base<=64{
				m.BaselineLosses++
				if frozen>64{m.FrozenRulePrevented++}
				if best>64{
					m.OraclePreventable++
					if frozen<=64{m.OraclePreventableFrozenMissed++}
				}else{
					m.NotPreventableWithTwoActions++
				}
			}
			ext:=best-base
			sumExt+=ext
			if ext>m.MaxBestOracleLossExtension{m.MaxBestOracleLossExtension=ext}
		}}
		if m.Arms>0{m.MeanBestOracleLossExtension=float64(sumExt)/float64(m.Arms)}
		res.Metrics=append(res.Metrics,m)
	}}
	return res,nil
}
