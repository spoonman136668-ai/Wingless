package unitary

const UP200COracleGapSchema="wingless.up200c-hostile-trigger-oracle-gap.v1"

type UP200CPoint struct{
	CohortIndex int `json:"cohort_index"`
	InitialHand int `json:"initial_hand"`
	BaselineLossStep int `json:"baseline_loss_step"`
	InternalOnset int `json:"internal_onset"`
	InternalActions int `json:"internal_actions"`
	InternalSurvives bool `json:"internal_survives"`
	OracleLatestSafeOnset int `json:"oracle_latest_safe_onset"`
	OracleActions int `json:"oracle_actions"`
	OnsetHeadroom int `json:"onset_headroom"`
}
type UP200CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP199CSeal string `json:"source_up199c_seal"`
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	SpacingIntervals int `json:"spacing_intervals"`
	ActionCap int `json:"action_cap"`
	Horizon int `json:"horizon"`
	Arms int `json:"arms"`
	BaselineLosses int `json:"baseline_losses"`
	InternalPrevented int `json:"internal_prevented"`
	InternalActions int `json:"internal_actions"`
	OracleSafeArms int `json:"oracle_safe_arms"`
	OracleActions int `json:"oracle_actions"`
	MeanInternalOnset float64 `json:"mean_internal_onset"`
	MeanOracleLatestSafeOnset float64 `json:"mean_oracle_latest_safe_onset"`
	MeanOnsetHeadroom float64 `json:"mean_onset_headroom"`
	OracleActionSavings int `json:"oracle_action_savings"`
	InternalEqualsLatestSafe int `json:"internal_equals_latest_safe"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	OracleFutureOutcomeAccess bool `json:"oracle_future_outcome_access"`
	OracleFeedsCandidateRule bool `json:"oracle_feeds_candidate_rule"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Points []UP200CPoint `json:"points"`
}
func up200cInternal(x0 *up81cAging,endangered,horizon int)(lossStep,actions,onset int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;onset=start;committed=true;since=0
			}
		}else if actions<8{
			since++
			if since>=4{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1021000+step,step,"hostile_shield"){return step,actions,onset}
		}
	}
	return horizon+1,actions,onset
}
func up200cForced(x0 *up81cAging,endangered,horizon,onset int)(lossStep,actions int){
	m:=*x0;committed:=false;since:=0
	for start:=1;start<=horizon;start+=2{
		if !committed&&start==onset{
			m.query(endangered);actions++;committed=true;since=0
		}else if committed&&actions<8{
			since++
			if since>=4{m.query(endangered);actions++;since=0}
		}
		for j:=0;j<2&&start+j<=horizon;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1022000+step,step,"hostile_shield"){return step,actions}
		}
	}
	return horizon+1,actions
}
func RunUP200C()(UP200CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	const horizon=56
	res:=UP200CResult{Schema:UP200COracleGapSchema,Experiment:"UP-200C-hostile-trigger-oracle-gap",SourceUP199CSeal:"ca99fbda42a5e5b3e4b8d5bd9287c7329617387c",Policy:"hostile_shield",Cadence:2,SpacingIntervals:4,ActionCap:8,Horizon:horizon,CounterfactualOnly:true,LiveActivation:false,OracleFutureOutcomeAccess:true,OracleFeedsCandidateRule:false,AdaptiveSelectionUsed:false}
	sumInternal,sumOracle,sumHead,nOracle:=0,0,0,0
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};res.Arms++
		base:=up169cBaselineLossStep(x,e,2,"hostile_shield");if base<=horizon{res.BaselineLosses++}
		iloss,ia,ion:=up200cInternal(x,e,horizon)
		surv:=iloss>horizon
		if base<=horizon&&surv{res.InternalPrevented++}
		res.InternalActions+=ia;sumInternal+=ion
		latest:=0;oa:=0
		for onset:=55;onset>=1;onset-=2{
			loss,actions:=up200cForced(x,e,horizon,onset)
			if loss>horizon{latest=onset;oa=actions;break}
		}
		head:=0
		if latest>0{
			res.OracleSafeArms++;res.OracleActions+=oa;nOracle++;sumOracle+=latest;head=latest-ion;sumHead+=head
			if head==0{res.InternalEqualsLatestSafe++}
		}
		res.Points=append(res.Points,UP200CPoint{CohortIndex:ci,InitialHand:hand,BaselineLossStep:base,InternalOnset:ion,InternalActions:ia,InternalSurvives:surv,OracleLatestSafeOnset:latest,OracleActions:oa,OnsetHeadroom:head})
	}}
	if res.Arms>0{res.MeanInternalOnset=float64(sumInternal)/float64(res.Arms)}
	if nOracle>0{res.MeanOracleLatestSafeOnset=float64(sumOracle)/float64(nOracle);res.MeanOnsetHeadroom=float64(sumHead)/float64(nOracle)}
	res.OracleActionSavings=res.InternalActions-res.OracleActions
	return res,nil
}
