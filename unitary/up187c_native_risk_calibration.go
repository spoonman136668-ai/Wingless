package unitary

const UP187CRiskCalibrationSchema="wingless.up187c-native-risk-calibration.v1"

type UP187CMetric struct{
	Cadence int `json:"cadence"`
	RiskBucket string `json:"risk_bucket"`
	Checkpoints int `json:"checkpoints"`
	FailuresWithinOneCadence int `json:"failures_within_one_cadence"`
	FailuresWithinTwoCadences int `json:"failures_within_two_cadences"`
	OneCadenceFailureRate float64 `json:"one_cadence_failure_rate"`
	TwoCadenceFailureRate float64 `json:"two_cadence_failure_rate"`
	MeanAdversarialShieldHorizon float64 `json:"mean_adversarial_shield_horizon"`
	MeanWritesRemainingToLoss float64 `json:"mean_writes_remaining_to_loss"`
}
type UP187CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP186CSeal string `json:"source_up186c_seal"`
	Cadences []int `json:"cadences"`
	RiskBuckets []string `json:"risk_buckets"`
	Policy string `json:"policy"`
	InterventionUsed bool `json:"intervention_used"`
	AdaptiveBucketingUsed bool `json:"adaptive_bucketing_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP187CMetric `json:"metrics"`
}
func up187cBucket(h,cad int)string{
	if h<=cad{return "critical"}
	if h<=2*cad{return "near"}
	return "far"
}
func RunUP187C()(UP187CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4};buckets:=[]string{"far","near","critical"}
	res:=UP187CResult{Schema:UP187CRiskCalibrationSchema,Experiment:"UP-187C-native-risk-calibration",SourceUP186CSeal:"0565a60b0388bf927f138ee2c467440352d6d6dc",Cadences:cadences,RiskBuckets:buckets,Policy:"no_refresh",InterventionUsed:false,AdaptiveBucketingUsed:false,WarningThresholdChanged:false,LiveActivation:false}
	for _,cad:=range cadences{
		type acc struct{n,one,two int;sumH,sumRemain float64}
		a:=map[string]*acc{"far":{},"near":{},"critical":{}}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			loss:=up169cBaselineLossStep(x,e,cad,"no_refresh")
			m:=*x
			for start:=1;start<=64;start+=cad{
				if loss<start{break}
				h:=up161cAdversarial(&m,e);bucket:=up187cBucket(h,cad);z:=a[bucket]
				remaining:=loss-start+1
				z.n++;z.sumH+=float64(h);z.sumRemain+=float64(remaining)
				if loss<=start+cad-1{z.one++}
				if loss<=start+2*cad-1{z.two++}
				lost:=false
				for j:=0;j<cad&&start+j<=64;j++{
					step:=start+j
					if up165cRealStep(&m,e,990000+step,step,"no_refresh"){lost=true;break}
				}
				if lost{break}
			}
		}}
		for _,bucket:=range buckets{
			z:=a[bucket];m:=UP187CMetric{Cadence:cad,RiskBucket:bucket,Checkpoints:z.n,FailuresWithinOneCadence:z.one,FailuresWithinTwoCadences:z.two}
			if z.n>0{
				m.OneCadenceFailureRate=float64(z.one)/float64(z.n)
				m.TwoCadenceFailureRate=float64(z.two)/float64(z.n)
				m.MeanAdversarialShieldHorizon=z.sumH/float64(z.n)
				m.MeanWritesRemainingToLoss=z.sumRemain/float64(z.n)
			}
			res.Metrics=append(res.Metrics,m)
		}
	}
	return res,nil
}
