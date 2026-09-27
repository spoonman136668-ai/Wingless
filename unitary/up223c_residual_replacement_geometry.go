package unitary

const UP223CGeometrySchema="wingless.up223c-residual-replacement-geometry.v1"

type UP223CSnapshot struct{
	Step int `json:"step"`
	Hand int `json:"hand"`
	EndangeredSlot int `json:"endangered_slot"`
	EndangeredAge int `json:"endangered_age"`
	PredictedSlot int `json:"predicted_slot"`
	PredictedIsEndangered bool `json:"predicted_is_endangered"`
	HandDistance int `json:"hand_distance"`
	Age0 int `json:"age0"`
	Age1 int `json:"age1"`
	Age2 int `json:"age2"`
	Age3 int `json:"age3"`
	AdversarialHorizon int `json:"adversarial_horizon"`
	NoQueryHorizon int `json:"no_query_horizon"`
}
type UP223CArm struct{
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Outcome string `json:"outcome"`
	Action8Step int `json:"action8_step"`
	LossStep int `json:"loss_step"`
	AdversarialCriticalSeen bool `json:"adversarial_critical_seen"`
	MinEndangeredAge int `json:"min_endangered_age"`
	MinHandDistance int `json:"min_hand_distance"`
	PredictedEndangeredCount int `json:"predicted_endangered_count"`
	LastEndangeredAge int `json:"last_endangered_age"`
	LastHandDistance int `json:"last_hand_distance"`
	LastPredictedIsEndangered bool `json:"last_predicted_is_endangered"`
}
type UP223CGroupSummary struct{
	Group string `json:"group"`
	Arms int `json:"arms"`
	MinEndangeredAgeMin int `json:"min_endangered_age_min"`
	MinEndangeredAgeMax int `json:"min_endangered_age_max"`
	MinHandDistanceMin int `json:"min_hand_distance_min"`
	MinHandDistanceMax int `json:"min_hand_distance_max"`
	PredictedEndangeredCountMin int `json:"predicted_endangered_count_min"`
	PredictedEndangeredCountMax int `json:"predicted_endangered_count_max"`
	LastEndangeredAgeMin int `json:"last_endangered_age_min"`
	LastEndangeredAgeMax int `json:"last_endangered_age_max"`
	LastHandDistanceMin int `json:"last_hand_distance_min"`
	LastHandDistanceMax int `json:"last_hand_distance_max"`
}
type UP223CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP222CSeal string `json:"source_up222c_seal"`
	Horizon int `json:"horizon"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP223CArm `json:"arms"`
	GroupSummaries []UP223CGroupSummary `json:"group_summaries"`
	SilentSnapshots []UP223CSnapshot `json:"silent_snapshots"`
}
func up223cSnapshot(m *up81cAging,endangered,step int)UP223CSnapshot{
	idx:=m.find(endangered)
	s:=UP223CSnapshot{Step:step,Hand:m.hand,EndangeredSlot:idx,EndangeredAge:-1}
	if idx>=0{
		s.EndangeredAge=int(m.age[idx]);s.HandDistance=(idx-m.hand+16)%16
	}
	s.PredictedSlot=up151cPredict(m.hand,m.age);s.PredictedIsEndangered=idx>=0&&s.PredictedSlot==idx
	for i:=0;i<16;i++{
		if !m.entries[i].used{continue}
		switch m.age[i]{case 0:s.Age0++;case 1:s.Age1++;case 2:s.Age2++;case 3:s.Age3++}
	}
	s.AdversarialHorizon=up161cAdversarial(m,endangered);s.NoQueryHorizon=up158cHorizon(m,endangered)
	return s
}
func up223cRun(x0 *up81cAging,endangered int)(UP223CArm,[]UP223CSnapshot){
	m:=*x0;committed:=false;since:=0;actions:=0
	a:=UP223CArm{EndangeredKey:endangered,Outcome:"survive",Action8Step:-1,LossStep:81,MinEndangeredAge:99,MinHandDistance:99}
	trace:=[]UP223CSnapshot{}
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;since=0;a.Action8Step=start}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if a.Action8Step>=0{
				s:=up223cSnapshot(&m,endangered,step);trace=append(trace,s)
				if s.AdversarialHorizon<=2{a.AdversarialCriticalSeen=true}
				if s.EndangeredAge>=0&&s.EndangeredAge<a.MinEndangeredAge{a.MinEndangeredAge=s.EndangeredAge}
				if s.HandDistance<a.MinHandDistance{a.MinHandDistance=s.HandDistance}
				if s.PredictedIsEndangered{a.PredictedEndangeredCount++}
				a.LastEndangeredAge=s.EndangeredAge;a.LastHandDistance=s.HandDistance;a.LastPredictedIsEndangered=s.PredictedIsEndangered
			}
			if up165cRealStep(&m,endangered,1046000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){
				a.Outcome="loss";a.LossStep=step
				if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
				return a,trace
			}
		}
	}
	if a.MinEndangeredAge==99{a.MinEndangeredAge=-1};if a.MinHandDistance==99{a.MinHandDistance=-1}
	return a,trace
}
func up223cGroup(a UP223CArm)string{
	if a.Outcome=="loss"&&!a.AdversarialCriticalSeen{return "silent_loss"}
	if a.Outcome=="loss"{return "warned_loss"}
	return "survive"
}
func RunUP223C()(UP223CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP223CResult{Schema:UP223CGeometrySchema,Experiment:"UP-223C-residual-replacement-geometry",SourceUP222CSeal:"ef2fe261ce81565591f981182efd5c8a25575c7f",Horizon:80,InterventionChanged:false,ThresholdFittingUsed:false,AdaptiveFeatureSelectionUsed:false,LiveActivation:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,tr:=up223cRun(x,e);a.Cohort=ci;a.InitialHand=hand;res.Arms=append(res.Arms,a)
		if up223cGroup(a)=="silent_loss"{res.SilentSnapshots=append(res.SilentSnapshots,tr...)}
	}}
	for _,g:=range []string{"silent_loss","warned_loss","survive"}{
		s:=UP223CGroupSummary{Group:g,MinEndangeredAgeMin:999,MinHandDistanceMin:999,PredictedEndangeredCountMin:999,LastEndangeredAgeMin:999,LastHandDistanceMin:999}
		for _,a:=range res.Arms{
			if up223cGroup(a)!=g{continue};s.Arms++
			if a.MinEndangeredAge<s.MinEndangeredAgeMin{s.MinEndangeredAgeMin=a.MinEndangeredAge};if a.MinEndangeredAge>s.MinEndangeredAgeMax{s.MinEndangeredAgeMax=a.MinEndangeredAge}
			if a.MinHandDistance<s.MinHandDistanceMin{s.MinHandDistanceMin=a.MinHandDistance};if a.MinHandDistance>s.MinHandDistanceMax{s.MinHandDistanceMax=a.MinHandDistance}
			if a.PredictedEndangeredCount<s.PredictedEndangeredCountMin{s.PredictedEndangeredCountMin=a.PredictedEndangeredCount};if a.PredictedEndangeredCount>s.PredictedEndangeredCountMax{s.PredictedEndangeredCountMax=a.PredictedEndangeredCount}
			if a.LastEndangeredAge<s.LastEndangeredAgeMin{s.LastEndangeredAgeMin=a.LastEndangeredAge};if a.LastEndangeredAge>s.LastEndangeredAgeMax{s.LastEndangeredAgeMax=a.LastEndangeredAge}
			if a.LastHandDistance<s.LastHandDistanceMin{s.LastHandDistanceMin=a.LastHandDistance};if a.LastHandDistance>s.LastHandDistanceMax{s.LastHandDistanceMax=a.LastHandDistance}
		}
		if s.Arms==0{s.MinEndangeredAgeMin=0;s.MinHandDistanceMin=0;s.PredictedEndangeredCountMin=0;s.LastEndangeredAgeMin=0;s.LastHandDistanceMin=0}
		res.GroupSummaries=append(res.GroupSummaries,s)
	}
	return res,nil
}
