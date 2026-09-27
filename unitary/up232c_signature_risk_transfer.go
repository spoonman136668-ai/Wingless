package unitary

const UP232CSignatureRiskSchema="wingless.up232c-signature-risk-transfer.v1"

type UP232CArm struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Outcome string `json:"outcome"`
	EventPresent bool `json:"event_present"`
	Signature string `json:"signature,omitempty"`
	RiskClass string `json:"risk_class"`
}
type UP232CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP231CSeal string `json:"source_up231c_seal"`
	Horizon int `json:"horizon"`
	MonitoringCadence int `json:"monitoring_cadence"`
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
	ExactSignatureCoverage float64 `json:"exact_signature_coverage"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	SignatureMutationUsed bool `json:"signature_mutation_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	InterventionChanged bool `json:"intervention_changed"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP232CArm `json:"arms"`
}

var up232cFailureSignatures=map[string]bool{
	"age=0|dist=11|pred=false|a0=3|a1=8|a2=5|a3=0|adv=1|noq=2":true,
	"age=0|dist=14|pred=false|a0=4|a1=7|a2=4|a3=1|adv=2|noq=4":true,
	"age=0|dist=1|pred=false|a0=9|a1=5|a2=2|a3=0|adv=1|noq=2":true,
	"age=0|dist=2|pred=false|a0=9|a1=5|a2=2|a3=0|adv=1|noq=2":true,
	"age=0|dist=3|pred=false|a0=7|a1=7|a2=2|a3=0|adv=1|noq=2":true,
	"age=0|dist=4|pred=false|a0=6|a1=8|a2=2|a3=0|adv=1|noq=2":true,
	"age=0|dist=5|pred=false|a0=5|a1=7|a2=4|a3=0|adv=1|noq=2":true,
	"age=0|dist=5|pred=false|a0=8|a1=7|a2=1|a3=0|adv=2|noq=3":true,
	"age=0|dist=5|pred=false|a0=9|a1=6|a2=1|a3=0|adv=2|noq=4":true,
	"age=0|dist=7|pred=false|a0=7|a1=8|a2=1|a3=0|adv=2|noq=3":true,
	"age=0|dist=9|pred=false|a0=4|a1=10|a2=2|a3=0|adv=2|noq=3":true,
	"age=0|dist=9|pred=false|a0=5|a1=8|a2=3|a3=0|adv=1|noq=2":true,
}
var up232cSurvivorSignatures=map[string]bool{
	"age=0|dist=5|pred=false|a0=8|a1=7|a2=1|a3=0|adv=2|noq=4":true,
	"age=0|dist=6|pred=false|a0=7|a1=6|a2=1|a3=2|adv=2|noq=4":true,
	"age=0|dist=6|pred=false|a0=8|a1=7|a2=1|a3=0|adv=2|noq=4":true,
	"age=0|dist=7|pred=false|a0=8|a1=6|a2=2|a3=0|adv=2|noq=4":true,
	"age=0|dist=8|pred=false|a0=7|a1=8|a2=1|a3=0|adv=2|noq=4":true,
	"age=1|dist=2|pred=false|a0=3|a1=7|a2=4|a3=2|adv=2|noq=4":true,
}
func up232cClass(event bool,sig string)string{
	if !event{return "no_event"}
	if up232cFailureSignatures[sig]{return "high_risk"}
	if up232cSurvivorSignatures[sig]{return "low_risk"}
	return "unknown"
}
func RunUP232C()(UP232CResult,error){
	cohorts:=[][]int{{0,6,9,15},{1,7,8,14},{2,4,11,13},{3,5,10,12}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{10,14}
	res:=UP232CResult{
		Schema:UP232CSignatureRiskSchema,Experiment:"UP-232C-signature-risk-transfer",
		SourceUP231CSeal:"3b6edd53b9cf2bfcb70f413fcba4e0e82cd8deb3",
		Horizon:80,MonitoringCadence:2,HeldoutFittingUsed:false,SignatureMutationUsed:false,
		AdaptiveFeatureSelectionUsed:false,ThresholdFittingUsed:false,InterventionChanged:false,LiveActivation:false,
	}
	for _,s:=range schedules{for _,advance:=range advances{for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		a,err:=up231cRun(x,e,advance,s.a,s.b);if err!=nil{return UP232CResult{},err}
		risk:=up232cClass(a.EventPresent,a.Signature)
		out:=UP232CArm{Schedule:s.name,PhaseAdvanceWrites:advance,Cohort:ci,InitialHand:hand,EndangeredKey:e,Outcome:a.Outcome,EventPresent:a.EventPresent,Signature:a.Signature,RiskClass:risk}
		res.ArmsTotal++
		fail:=a.Outcome=="loss"
		if fail{res.Failures++}else{res.Survivors++}
		if a.EventPresent{if fail{res.EventBearingFailures++}else{res.EventBearingSurvivors++}}
		switch risk{
		case "high_risk":
			res.HighRisk++;if fail{res.HighRiskFailures++}else{res.HighRiskSurvivors++}
		case "low_risk":
			res.LowRisk++;if fail{res.LowRiskFailures++}else{res.LowRiskSurvivors++}
		case "unknown":
			res.Unknown++;if fail{res.UnknownFailures++}else{res.UnknownSurvivors++}
		default:
			res.NoEvent++
		}
		res.Arms=append(res.Arms,out)
	}}}}
	if res.Failures>0{res.HighRiskFailureRecall=float64(res.HighRiskFailures)/float64(res.Failures)}
	if res.HighRisk>0{res.HighRiskPrecision=float64(res.HighRiskFailures)/float64(res.HighRisk)}
	eventTotal:=res.EventBearingFailures+res.EventBearingSurvivors
	if eventTotal>0{res.ExactSignatureCoverage=float64(res.HighRisk+res.LowRisk)/float64(eventTotal)}
	return res,nil
}
