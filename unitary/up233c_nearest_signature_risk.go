package unitary

import "fmt"

const UP233CNearestRiskSchema="wingless.up233c-nearest-signature-risk.v1"

type up233cVec [9]int

type UP233CArm struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	Outcome string `json:"outcome"`
	EventPresent bool `json:"event_present"`
	RiskClass string `json:"risk_class"`
	MinFailureDistance int `json:"min_failure_distance"`
	MinSurvivorDistance int `json:"min_survivor_distance"`
}

type UP233CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP232CSeal string `json:"source_up232c_seal"`
	ArmsTotal int `json:"arms_total"`
	Failures int `json:"failures"`
	Survivors int `json:"survivors"`
	EventBearingFailures int `json:"event_bearing_failures"`
	EventBearingSurvivors int `json:"event_bearing_survivors"`
	HighRisk int `json:"high_risk"`
	LowRisk int `json:"low_risk"`
	Unknown int `json:"unknown"`
	NoEvent int `json:"no_event"`
	HighRiskFailures int `json:"high_risk_failures"`
	HighRiskSurvivors int `json:"high_risk_survivors"`
	LowRiskFailures int `json:"low_risk_failures"`
	LowRiskSurvivors int `json:"low_risk_survivors"`
	UnknownFailures int `json:"unknown_failures"`
	UnknownSurvivors int `json:"unknown_survivors"`
	HighRiskFailureRecall float64 `json:"high_risk_failure_recall"`
	HighRiskPrecision float64 `json:"high_risk_precision"`
	MeanFailureMinFailureDistance float64 `json:"mean_failure_min_failure_distance"`
	MeanFailureMinSurvivorDistance float64 `json:"mean_failure_min_survivor_distance"`
	MeanSurvivorMinFailureDistance float64 `json:"mean_survivor_min_failure_distance"`
	MeanSurvivorMinSurvivorDistance float64 `json:"mean_survivor_min_survivor_distance"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	AdaptiveWeightingUsed bool `json:"adaptive_weighting_used"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	InterventionChanged bool `json:"intervention_changed"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP233CArm `json:"arms"`
}

func up233cParse(sig string)(up233cVec,bool){
	var age,dist,a0,a1,a2,a3,adv,noq int
	var pred bool
	n,err:=fmt.Sscanf(sig,"age=%d|dist=%d|pred=%t|a0=%d|a1=%d|a2=%d|a3=%d|adv=%d|noq=%d",&age,&dist,&pred,&a0,&a1,&a2,&a3,&adv,&noq)
	if err!=nil||n!=9{return up233cVec{},false}
	p:=0;if pred{p=1}
	return up233cVec{age,dist,p,a0,a1,a2,a3,adv,noq},true
}
func up233cFromSnapshot(s UP223CSnapshot)up233cVec{
	p:=0;if s.PredictedIsEndangered{p=1}
	return up233cVec{s.EndangeredAge,s.HandDistance,p,s.Age0,s.Age1,s.Age2,s.Age3,s.AdversarialHorizon,s.NoQueryHorizon}
}
func up233cDist(a,b up233cVec)int{d:=0;for i:=0;i<9;i++{if a[i]!=b[i]{d++}};return d}
func up233cExemplars(m map[string]bool)[]up233cVec{
	out:=[]up233cVec{}
	for s:=range m{if v,ok:=up233cParse(s);ok{out=append(out,v)}}
	return out
}
func up233cClass(event bool,s UP223CSnapshot,fail,surv []up233cVec)(string,int,int){
	if !event{return "no_event",-1,-1}
	x:=up233cFromSnapshot(s);fd,sd:=99,99
	for _,v:=range fail{if d:=up233cDist(x,v);d<fd{fd=d}}
	for _,v:=range surv{if d:=up233cDist(x,v);d<sd{sd=d}}
	if fd<sd{return "high_risk",fd,sd}
	if sd<fd{return "low_risk",fd,sd}
	return "unknown",fd,sd
}
func RunUP233C()(UP233CResult,error){
	failEx:=up233cExemplars(up232cFailureSignatures);survEx:=up233cExemplars(up232cSurvivorSignatures)
	cohorts:=[][]int{{0,7,10,13},{1,4,11,14},{2,5,8,15},{3,6,9,12}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{{"no_hostile","no_refresh","hostile_shield"},{"hostile_no","hostile_shield","no_refresh"},{"alternating_fixed","alternating_shield","fixed_offset_refresh"},{"fixed_alternating","fixed_offset_refresh","alternating_shield"}}
	res:=UP233CResult{Schema:UP233CNearestRiskSchema,Experiment:"UP-233C-nearest-signature-risk",SourceUP232CSeal:"5054470dab9f5b157c907da40ec0da9b08a4ff20",HeldoutFittingUsed:false,NewNativeFeatureUsed:false,AdaptiveWeightingUsed:false,ThresholdFittingUsed:false,InterventionChanged:false,LiveActivation:false}
	var ffd,fsd,sfd,ssd int
	for _,s:=range schedules{for _,advance:=range []int{6,16}{for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,err:=up231cRun(x,e,advance,s.a,s.b);if err!=nil{return UP233CResult{},err}
		risk,fd,sd:=up233cClass(a.EventPresent,a.Snapshot,failEx,survEx)
		fail:=a.Outcome=="loss";res.ArmsTotal++
		if fail{res.Failures++}else{res.Survivors++}
		if a.EventPresent{
			if fail{res.EventBearingFailures++;ffd+=fd;fsd+=sd}else{res.EventBearingSurvivors++;sfd+=fd;ssd+=sd}
		}
		switch risk{
		case "high_risk":res.HighRisk++;if fail{res.HighRiskFailures++}else{res.HighRiskSurvivors++}
		case "low_risk":res.LowRisk++;if fail{res.LowRiskFailures++}else{res.LowRiskSurvivors++}
		case "unknown":res.Unknown++;if fail{res.UnknownFailures++}else{res.UnknownSurvivors++}
		default:res.NoEvent++
		}
		res.Arms=append(res.Arms,UP233CArm{Schedule:s.name,PhaseAdvanceWrites:advance,Cohort:ci,InitialHand:hand,Outcome:a.Outcome,EventPresent:a.EventPresent,RiskClass:risk,MinFailureDistance:fd,MinSurvivorDistance:sd})
	}}}}
	if res.Failures>0{res.HighRiskFailureRecall=float64(res.HighRiskFailures)/float64(res.Failures)}
	if res.HighRisk>0{res.HighRiskPrecision=float64(res.HighRiskFailures)/float64(res.HighRisk)}
	if res.EventBearingFailures>0{res.MeanFailureMinFailureDistance=float64(ffd)/float64(res.EventBearingFailures);res.MeanFailureMinSurvivorDistance=float64(fsd)/float64(res.EventBearingFailures)}
	if res.EventBearingSurvivors>0{res.MeanSurvivorMinFailureDistance=float64(sfd)/float64(res.EventBearingSurvivors);res.MeanSurvivorMinSurvivorDistance=float64(ssd)/float64(res.EventBearingSurvivors)}
	return res,nil
}
