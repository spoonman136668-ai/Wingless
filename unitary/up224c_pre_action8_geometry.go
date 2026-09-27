package unitary

const UP224CPreAction8Schema="wingless.up224c-pre-action8-geometry.v1"

type UP224CArm struct{
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Group string `json:"group"`
	Action7Step int `json:"action7_step"`
	Action8Step int `json:"action8_step"`
	LossStep int `json:"loss_step"`
	Observations int `json:"observations"`
	MinAdversarialHorizon int `json:"min_adversarial_horizon"`
	MinNoQueryHorizon int `json:"min_no_query_horizon"`
	MinEndangeredAge int `json:"min_endangered_age"`
	MinHandDistance int `json:"min_hand_distance"`
	PredictedEndangeredCount int `json:"predicted_endangered_count"`
}
type UP224CGroupSummary struct{
	Group string `json:"group"`
	Arms int `json:"arms"`
	ObservationMin int `json:"observation_min"`
	ObservationMax int `json:"observation_max"`
	MinAdversarialHorizonMin int `json:"min_adversarial_horizon_min"`
	MinAdversarialHorizonMax int `json:"min_adversarial_horizon_max"`
	MinNoQueryHorizonMin int `json:"min_no_query_horizon_min"`
	MinNoQueryHorizonMax int `json:"min_no_query_horizon_max"`
	MinEndangeredAgeMin int `json:"min_endangered_age_min"`
	MinEndangeredAgeMax int `json:"min_endangered_age_max"`
	MinHandDistanceMin int `json:"min_hand_distance_min"`
	MinHandDistanceMax int `json:"min_hand_distance_max"`
	PredictedEndangeredCountMin int `json:"predicted_endangered_count_min"`
	PredictedEndangeredCountMax int `json:"predicted_endangered_count_max"`
}
type UP224CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP223CSeal string `json:"source_up223c_seal"`
	Horizon int `json:"horizon"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP224CArm `json:"arms"`
	GroupSummaries []UP224CGroupSummary `json:"group_summaries"`
	PreAction8LossSnapshots []UP223CSnapshot `json:"pre_action8_loss_snapshots"`
}
func up224cObserve(a *UP224CArm,s UP223CSnapshot){
	a.Observations++
	if s.AdversarialHorizon<a.MinAdversarialHorizon{a.MinAdversarialHorizon=s.AdversarialHorizon}
	if s.NoQueryHorizon<a.MinNoQueryHorizon{a.MinNoQueryHorizon=s.NoQueryHorizon}
	if s.EndangeredAge>=0&&s.EndangeredAge<a.MinEndangeredAge{a.MinEndangeredAge=s.EndangeredAge}
	if s.HandDistance<a.MinHandDistance{a.MinHandDistance=s.HandDistance}
	if s.PredictedIsEndangered{a.PredictedEndangeredCount++}
}
func up224cRun(x0 *up81cAging,endangered int)(UP224CArm,[]UP223CSnapshot){
	m:=*x0;committed:=false;since:=0;actions:=0
	a:=UP224CArm{EndangeredKey:endangered,Group:"no_action8_survive",Action7Step:-1,Action8Step:-1,LossStep:81,MinAdversarialHorizon:65,MinNoQueryHorizon:65,MinEndangeredAge:99,MinHandDistance:99}
	trace:=[]UP223CSnapshot{}
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{
				m.query(endangered);actions++;since=0
				if actions==7{a.Action7Step=start}
			}
			if actions==7{
				s:=up223cSnapshot(&m,endangered,start);trace=append(trace,s);up224cObserve(&a,s)
				if s.AdversarialHorizon<=2{
					m.query(endangered);actions++;since=0;a.Action8Step=start;a.Group="action8_reached"
				}
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if a.Action7Step>=0&&a.Action8Step<0&&j==1{
				s:=up223cSnapshot(&m,endangered,step);trace=append(trace,s);up224cObserve(&a,s)
			}
			if up165cRealStep(&m,endangered,1047000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){
				a.LossStep=step
				if a.Action8Step<0{a.Group="pre_action8_loss"}
				if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
				return a,trace
			}
		}
	}
	if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
	return a,trace
}
func RunUP224C()(UP224CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP224CResult{Schema:UP224CPreAction8Schema,Experiment:"UP-224C-pre-action8-geometry",SourceUP223CSeal:"bd93d6f27aa8c2899e646d70e426135c13804d4f",Horizon:80,InterventionChanged:false,ThresholdFittingUsed:false,AdaptiveFeatureSelectionUsed:false,LiveActivation:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,tr:=up224cRun(x,e);a.Cohort=ci;a.InitialHand=hand;res.Arms=append(res.Arms,a)
		if a.Group=="pre_action8_loss"{res.PreAction8LossSnapshots=append(res.PreAction8LossSnapshots,tr...)}
	}}
	for _,g:=range []string{"pre_action8_loss","action8_reached","no_action8_survive"}{
		s:=UP224CGroupSummary{Group:g,ObservationMin:999,MinAdversarialHorizonMin:999,MinNoQueryHorizonMin:999,MinEndangeredAgeMin:999,MinHandDistanceMin:999,PredictedEndangeredCountMin:999}
		for _,a:=range res.Arms{
			if a.Group!=g{continue};s.Arms++
			if a.Observations<s.ObservationMin{s.ObservationMin=a.Observations};if a.Observations>s.ObservationMax{s.ObservationMax=a.Observations}
			if a.MinAdversarialHorizon<s.MinAdversarialHorizonMin{s.MinAdversarialHorizonMin=a.MinAdversarialHorizon};if a.MinAdversarialHorizon>s.MinAdversarialHorizonMax{s.MinAdversarialHorizonMax=a.MinAdversarialHorizon}
			if a.MinNoQueryHorizon<s.MinNoQueryHorizonMin{s.MinNoQueryHorizonMin=a.MinNoQueryHorizon};if a.MinNoQueryHorizon>s.MinNoQueryHorizonMax{s.MinNoQueryHorizonMax=a.MinNoQueryHorizon}
			if a.MinEndangeredAge<s.MinEndangeredAgeMin{s.MinEndangeredAgeMin=a.MinEndangeredAge};if a.MinEndangeredAge>s.MinEndangeredAgeMax{s.MinEndangeredAgeMax=a.MinEndangeredAge}
			if a.MinHandDistance<s.MinHandDistanceMin{s.MinHandDistanceMin=a.MinHandDistance};if a.MinHandDistance>s.MinHandDistanceMax{s.MinHandDistanceMax=a.MinHandDistance}
			if a.PredictedEndangeredCount<s.PredictedEndangeredCountMin{s.PredictedEndangeredCountMin=a.PredictedEndangeredCount};if a.PredictedEndangeredCount>s.PredictedEndangeredCountMax{s.PredictedEndangeredCountMax=a.PredictedEndangeredCount}
		}
		if s.Arms==0{s.ObservationMin=0;s.MinAdversarialHorizonMin=0;s.MinNoQueryHorizonMin=0;s.MinEndangeredAgeMin=0;s.MinHandDistanceMin=0;s.PredictedEndangeredCountMin=0}
		res.GroupSummaries=append(res.GroupSummaries,s)
	}
	return res,nil
}
