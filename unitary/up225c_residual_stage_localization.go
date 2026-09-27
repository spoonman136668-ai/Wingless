package unitary

const UP225CStageSchema="wingless.up225c-residual-stage-localization.v1"

type UP225CArm struct{
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Outcome string `json:"outcome"`
	LossStep int `json:"loss_step"`
	ActionsTaken int `json:"actions_taken"`
	ActionSteps []int `json:"action_steps"`
	AdversarialCriticalAfterAction1 bool `json:"adversarial_critical_after_action1"`
	MinAdversarialHorizon int `json:"min_adversarial_horizon"`
	MinNoQueryHorizon int `json:"min_no_query_horizon"`
	MinEndangeredAge int `json:"min_endangered_age"`
	MinHandDistance int `json:"min_hand_distance"`
	PredictedEndangeredCount int `json:"predicted_endangered_count"`
}
type UP225CStageSummary struct{
	ActionsTaken int `json:"actions_taken"`
	LossArms int `json:"loss_arms"`
	SurvivorArms int `json:"survivor_arms"`
}
type UP225CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP224CSeal string `json:"source_up224c_seal"`
	Horizon int `json:"horizon"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP225CArm `json:"arms"`
	StageSummaries []UP225CStageSummary `json:"stage_summaries"`
	SilentLossArms []UP225CArm `json:"silent_loss_arms"`
	SilentLossSnapshots []UP223CSnapshot `json:"silent_loss_snapshots"`
}
func up225cObserve(a *UP225CArm,s UP223CSnapshot){
	if s.AdversarialHorizon<=2{a.AdversarialCriticalAfterAction1=true}
	if s.AdversarialHorizon<a.MinAdversarialHorizon{a.MinAdversarialHorizon=s.AdversarialHorizon}
	if s.NoQueryHorizon<a.MinNoQueryHorizon{a.MinNoQueryHorizon=s.NoQueryHorizon}
	if s.EndangeredAge>=0&&s.EndangeredAge<a.MinEndangeredAge{a.MinEndangeredAge=s.EndangeredAge}
	if s.HandDistance<a.MinHandDistance{a.MinHandDistance=s.HandDistance}
	if s.PredictedIsEndangered{a.PredictedEndangeredCount++}
}
func up225cRun(x0 *up81cAging,endangered int)(UP225CArm,[]UP223CSnapshot){
	m:=*x0;committed:=false;since:=0;actions:=0
	steps:=make([]int,8);for i:=range steps{steps[i]=-1}
	a:=UP225CArm{EndangeredKey:endangered,Outcome:"survive",LossStep:81,ActionSteps:steps,MinAdversarialHorizon:65,MinNoQueryHorizon:65,MinEndangeredAge:99,MinHandDistance:99}
	trace:=[]UP223CSnapshot{}
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;a.ActionSteps[0]=start;committed=true;since=0
			}
		}else{
			since++
			if actions<7&&since>=4{
				m.query(endangered);actions++;a.ActionSteps[actions-1]=start;since=0
			}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;a.ActionSteps[7]=start;since=0
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if actions>0{
				s:=up223cSnapshot(&m,endangered,step);trace=append(trace,s);up225cObserve(&a,s)
			}
			if up165cRealStep(&m,endangered,1048000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){
				a.Outcome="loss";a.LossStep=step;a.ActionsTaken=actions
				if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
				return a,trace
			}
		}
	}
	a.ActionsTaken=actions
	if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
	return a,trace
}
func RunUP225C()(UP225CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP225CResult{Schema:UP225CStageSchema,Experiment:"UP-225C-residual-stage-localization",SourceUP224CSeal:"6b5bc0d6774374ea79d6594bd548101345af51ca",Horizon:80,InterventionChanged:false,ThresholdFittingUsed:false,AdaptiveFeatureSelectionUsed:false,LiveActivation:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,tr:=up225cRun(x,e);a.Cohort=ci;a.InitialHand=hand;res.Arms=append(res.Arms,a)
		if a.Outcome=="loss"&&!a.AdversarialCriticalAfterAction1{
			res.SilentLossArms=append(res.SilentLossArms,a)
			res.SilentLossSnapshots=append(res.SilentLossSnapshots,tr...)
		}
	}}
	for n:=0;n<=8;n++{
		s:=UP225CStageSummary{ActionsTaken:n}
		for _,a:=range res.Arms{
			if a.ActionsTaken!=n{continue}
			if a.Outcome=="loss"{s.LossArms++}else{s.SurvivorArms++}
		}
		res.StageSummaries=append(res.StageSummaries,s)
	}
	return res,nil
}
