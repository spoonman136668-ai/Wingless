package unitary

const UP188CRiskPolicySchema="wingless.up188c-risk-policy-transfer.v1"

type UP188CMetric struct{
	Policy string `json:"policy"`
	Cadence int `json:"cadence"`
	RiskBucket string `json:"risk_bucket"`
	Checkpoints int `json:"checkpoints"`
	FailuresWithinOneCadence int `json:"failures_within_one_cadence"`
	FailuresWithinTwoCadences int `json:"failures_within_two_cadences"`
	OneCadenceFailureRate float64 `json:"one_cadence_failure_rate"`
	TwoCadenceFailureRate float64 `json:"two_cadence_failure_rate"`
	TrajectoryLossesBy64 int `json:"trajectory_losses_by_64"`
}

type UP188CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP187CSeal string `json:"source_up187c_seal"`
	Policies []string `json:"policies"`
	Cadences []int `json:"cadences"`
	RiskBuckets []string `json:"risk_buckets"`
	InterventionUsed bool `json:"intervention_used"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	AdaptiveBucketingUsed bool `json:"adaptive_bucketing_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	FuturePolicyScheduleUsedByWarning bool `json:"future_policy_schedule_used_by_warning"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP188CMetric `json:"metrics"`
}

func RunUP188C()(UP188CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	policies:=[]string{"no_refresh","alternating_shield","fixed_offset_refresh","hostile_shield"}
	cadences:=[]int{2,4}
	buckets:=[]string{"far","near","critical"}
	res:=UP188CResult{
		Schema:UP188CRiskPolicySchema,Experiment:"UP-188C-risk-policy-transfer",
		SourceUP187CSeal:"d575783d52fac813325204d884dfde0573ea648c",
		Policies:policies,Cadences:cadences,RiskBuckets:buckets,
		InterventionUsed:false,AdaptivePolicySelectionUsed:false,AdaptiveBucketingUsed:false,
		WarningThresholdChanged:false,FuturePolicyScheduleUsedByWarning:false,LiveActivation:false,
	}
	for _,policy:=range policies{for _,cad:=range cadences{
		type acc struct{n,one,two,losses int}
		a:=map[string]*acc{"far":{},"near":{},"critical":{}}
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue}
			loss:=up169cBaselineLossStep(x,e,cad,policy)
			if loss<=64{
				for _,bucket:=range buckets{a[bucket].losses++}
			}
			m:=*x
			for start:=1;start<=64;start+=cad{
				if loss<start{break}
				h:=up161cAdversarial(&m,e);bucket:=up187cBucket(h,cad);z:=a[bucket];z.n++
				if loss<=start+cad-1{z.one++}
				if loss<=start+2*cad-1{z.two++}
				lost:=false
				for j:=0;j<cad&&start+j<=64;j++{
					step:=start+j
					if up165cRealStep(&m,e,990000+step,step,policy){lost=true;break}
				}
				if lost{break}
			}
		}}
		for _,bucket:=range buckets{
			z:=a[bucket]
			m:=UP188CMetric{Policy:policy,Cadence:cad,RiskBucket:bucket,Checkpoints:z.n,FailuresWithinOneCadence:z.one,FailuresWithinTwoCadences:z.two,TrajectoryLossesBy64:z.losses}
			if z.n>0{
				m.OneCadenceFailureRate=float64(z.one)/float64(z.n)
				m.TwoCadenceFailureRate=float64(z.two)/float64(z.n)
			}
			res.Metrics=append(res.Metrics,m)
		}
	}}
	return res,nil
}
