package unitary

const UP221CDualHorizonSchema="wingless.up221c-residual-dual-horizon-diagnostic.v1"

type UP221CArm struct{
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Action8Step int `json:"action8_step"`
	LossStep int `json:"loss_step"`
	Cap8Loss bool `json:"cap8_loss"`
	MinAdversarialHorizon int `json:"min_adversarial_horizon"`
	MinNoQueryHorizon int `json:"min_no_query_horizon"`
	AdversarialCriticalSeen bool `json:"adversarial_critical_seen"`
	NoQueryCriticalSeen bool `json:"no_query_critical_seen"`
	NoQueryNearSeen bool `json:"no_query_near_seen"`
}
type UP221CSummary struct{
	Outcome string `json:"outcome"`
	Arms int `json:"arms"`
	AdversarialCritical int `json:"adversarial_critical"`
	AdversarialCriticalRate float64 `json:"adversarial_critical_rate"`
	NoQueryCritical int `json:"no_query_critical"`
	NoQueryCriticalRate float64 `json:"no_query_critical_rate"`
	NoQueryNear int `json:"no_query_near"`
	NoQueryNearRate float64 `json:"no_query_near_rate"`
}
type UP221CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP220CSeal string `json:"source_up220c_seal"`
	PolicyA string `json:"policy_a"`
	PolicyB string `json:"policy_b"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Horizon int `json:"horizon"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveSignalSelectionUsed bool `json:"adaptive_signal_selection_used"`
	LiveActivation bool `json:"live_activation"`
	SilentCap8Failures int `json:"silent_cap8_failures"`
	SilentCaughtByNoQueryCritical int `json:"silent_caught_by_no_query_critical"`
	SilentCaughtByNoQueryNear int `json:"silent_caught_by_no_query_near"`
	Arms []UP221CArm `json:"arms"`
	Summaries []UP221CSummary `json:"summaries"`
}
func up221cArm(x0 *up81cAging,endangered int)(UP221CArm,error){
	m:=*x0;committed:=false;since:=0;actions:=0
	a:=UP221CArm{EndangeredKey:endangered,Action8Step:-1,LossStep:81,MinAdversarialHorizon:65,MinNoQueryHorizon:65}
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{m.query(endangered);actions++;committed=true;since=0}
		}else{
			since++
			if actions<7&&since>=4{m.query(endangered);actions++;since=0}
			if actions==7&&up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;since=0;a.Action8Step=start
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if a.Action8Step>=0{
				adv:=up161cAdversarial(&m,endangered);nq:=up158cHorizon(&m,endangered)
				if adv<a.MinAdversarialHorizon{a.MinAdversarialHorizon=adv}
				if nq<a.MinNoQueryHorizon{a.MinNoQueryHorizon=nq}
				if adv<=2{a.AdversarialCriticalSeen=true}
				if nq<=2{a.NoQueryCriticalSeen=true}
				if nq<=4{a.NoQueryNearSeen=true}
			}
			if up165cRealStep(&m,endangered,1044000+step,step,up210cPolicy(step,8,"alternating_shield","fixed_offset_refresh")){
				a.LossStep=step;a.Cap8Loss=true;return a,nil
			}
		}
	}
	return a,nil
}
func RunUP221C()(UP221CResult,error){
	cohorts:=[][]int{{0,1,2,3},{4,5,6,7},{8,9,10,11},{12,13,14,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	res:=UP221CResult{Schema:UP221CDualHorizonSchema,Experiment:"UP-221C-residual-dual-horizon-diagnostic",SourceUP220CSeal:"b06af67a6c5ef511c168d03dfdefd87f687000c5",PolicyA:"alternating_shield",PolicyB:"fixed_offset_refresh",PhaseAdvanceWrites:8,Horizon:80,InterventionChanged:false,ThresholdFittingUsed:false,AdaptiveSignalSelectionUsed:false,LiveActivation:false}
	for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,err:=up221cArm(x,e);if err!=nil{return UP221CResult{},err};a.Cohort=ci;a.InitialHand=hand
		res.Arms=append(res.Arms,a)
		if a.Cap8Loss&&!a.AdversarialCriticalSeen{
			res.SilentCap8Failures++
			if a.NoQueryCriticalSeen{res.SilentCaughtByNoQueryCritical++}
			if a.NoQueryNearSeen{res.SilentCaughtByNoQueryNear++}
		}
	}}
	for _,wantLoss:=range []bool{true,false}{
		out:="survive";if wantLoss{out="loss"}
		s:=UP221CSummary{Outcome:out}
		for _,a:=range res.Arms{
			if a.Cap8Loss!=wantLoss{continue}
			s.Arms++
			if a.AdversarialCriticalSeen{s.AdversarialCritical++}
			if a.NoQueryCriticalSeen{s.NoQueryCritical++}
			if a.NoQueryNearSeen{s.NoQueryNear++}
		}
		if s.Arms>0{
			s.AdversarialCriticalRate=float64(s.AdversarialCritical)/float64(s.Arms)
			s.NoQueryCriticalRate=float64(s.NoQueryCritical)/float64(s.Arms)
			s.NoQueryNearRate=float64(s.NoQueryNear)/float64(s.Arms)
		}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
