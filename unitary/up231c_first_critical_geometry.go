package unitary

import "fmt"

const UP231CGeometrySchema="wingless.up231c-first-critical-geometry.v1"

type UP231CArm struct{
	Schedule string `json:"schedule"`
	PhaseAdvanceWrites int `json:"phase_advance_writes"`
	Cohort int `json:"cohort"`
	InitialHand int `json:"initial_hand"`
	EndangeredKey int `json:"endangered_key"`
	Outcome string `json:"outcome"`
	Action8Step int `json:"action8_step"`
	EventPresent bool `json:"event_present"`
	EventStep int `json:"event_step"`
	LossStep int `json:"loss_step"`
	Signature string `json:"signature,omitempty"`
	Snapshot UP223CSnapshot `json:"snapshot"`
}
type UP231CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP230CSeal string `json:"source_up230c_seal"`
	Horizon int `json:"horizon"`
	MonitoringCadence int `json:"monitoring_cadence"`
	ArmsTotal int `json:"arms_total"`
	Failures int `json:"failures"`
	Survivors int `json:"survivors"`
	EventBearingFailures int `json:"event_bearing_failures"`
	EventBearingSurvivors int `json:"event_bearing_survivors"`
	SilentFailures int `json:"silent_failures"`
	NoEventSurvivors int `json:"no_event_survivors"`
	DistinctFailureSignatures int `json:"distinct_failure_signatures"`
	DistinctSurvivorSignatures int `json:"distinct_survivor_signatures"`
	SharedSignatures int `json:"shared_signatures"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFeatureUsed bool `json:"new_native_feature_used"`
	LiveActivation bool `json:"live_activation"`
	Arms []UP231CArm `json:"arms"`
}
func up231cSignature(s UP223CSnapshot)string{
	return fmt.Sprintf("age=%d|dist=%d|pred=%t|a0=%d|a1=%d|a2=%d|a3=%d|adv=%d|noq=%d",
		s.EndangeredAge,s.HandDistance,s.PredictedIsEndangered,s.Age0,s.Age1,s.Age2,s.Age3,s.AdversarialHorizon,s.NoQueryHorizon)
}
func up231cRun(x0 *up81cAging,endangered,advance int,a,b string)(UP231CArm,error){
	m:=*x0
	arm:=UP231CArm{EndangeredKey:endangered,Outcome:"survive",Action8Step:-1,EventStep:-1,LossStep:81}
	committed:=false;actions:=0
	steps:=make([]int,8);for i:=range steps{steps[i]=-1}
	action4:=-1;due5,due6,due7:=-1,-1,-1
	for start:=1;start<=80;start+=2{
		if !committed{
			if up161cAdversarial(&m,endangered)<=2{
				m.query(endangered);actions++;steps[0]=start;committed=true
			}
		}else{
			if actions>=1&&actions<4{
				last:=steps[actions-1]
				if last>=0&&start>=last+8{
					m.query(endangered);actions++;steps[actions-1]=start
					if actions==4{action4=start;due5=action4+8;due6=action4+16;due7=action4+24}
				}
			}else if actions==4{
				fire:=start>=due5
				if start<due5&&up161cAdversarial(&m,endangered)<=2{fire=true}
				if fire{m.query(endangered);actions++;steps[4]=start}
			}else if actions==5{
				if start>=due6{m.query(endangered);actions++;steps[5]=start}
			}else if actions==6{
				if start>=due7{m.query(endangered);actions++;steps[6]=start}
			}else if actions==7{
				if up161cAdversarial(&m,endangered)<=2{
					m.query(endangered);actions++;steps[7]=start;arm.Action8Step=start
				}
			}else if actions==8&&arm.Action8Step>=0&&start>arm.Action8Step&&!arm.EventPresent{
				if up161cAdversarial(&m,endangered)<=2{
					s:=up223cSnapshot(&m,endangered,start)
					arm.EventPresent=true;arm.EventStep=start;arm.Snapshot=s;arm.Signature=up231cSignature(s)
				}
			}
		}
		for j:=0;j<2&&start+j<=80;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1053000+step,step,up210cPolicy(step,advance,a,b)){
				arm.Outcome="loss";arm.LossStep=step
				return arm,nil
			}
		}
	}
	return arm,nil
}
func RunUP231C()(UP231CResult,error){
	cohorts:=[][]int{{0,5,10,15},{1,6,11,12},{2,7,8,13},{3,4,9,14}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	type sched struct{name,a,b string}
	schedules:=[]sched{
		{"no_hostile","no_refresh","hostile_shield"},
		{"hostile_no","hostile_shield","no_refresh"},
		{"alternating_fixed","alternating_shield","fixed_offset_refresh"},
		{"fixed_alternating","fixed_offset_refresh","alternating_shield"},
	}
	advances:=[]int{8,12}
	res:=UP231CResult{
		Schema:UP231CGeometrySchema,Experiment:"UP-231C-first-critical-geometry",
		SourceUP230CSeal:"869367b4c7b66ced2d757f3d1b0a4be12e0eeedd",
		Horizon:80,MonitoringCadence:2,InterventionChanged:false,ThresholdFittingUsed:false,
		AdaptiveFeatureSelectionUsed:false,NewNativeFeatureUsed:false,LiveActivation:false,
	}
	failS:=map[string]bool{};survS:=map[string]bool{}
	for _,s:=range schedules{for _,advance:=range advances{for ci,c:=range cohorts{for _,hand:=range hands{
		x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
		arm,err:=up231cRun(x,e,advance,s.a,s.b);if err!=nil{return UP231CResult{},err}
		arm.Schedule=s.name;arm.PhaseAdvanceWrites=advance;arm.Cohort=ci;arm.InitialHand=hand
		res.ArmsTotal++
		if arm.Outcome=="loss"{
			res.Failures++
			if arm.EventPresent{res.EventBearingFailures++;failS[arm.Signature]=true}else{res.SilentFailures++}
		}else{
			res.Survivors++
			if arm.EventPresent{res.EventBearingSurvivors++;survS[arm.Signature]=true}else{res.NoEventSurvivors++}
		}
		res.Arms=append(res.Arms,arm)
	}}}}
	res.DistinctFailureSignatures=len(failS);res.DistinctSurvivorSignatures=len(survS)
	for k:=range failS{if survS[k]{res.SharedSignatures++}}
	return res,nil
}
