package unitary

const UP186CWarningLeadSchema="wingless.up186c-warning-lead-calibration.v1"

type UP186CMetric struct{
	Cadence int `json:"cadence"`
	Arms int `json:"arms"`
	WarnedBeforeLoss int `json:"warned_before_loss"`
	LostBy64 int `json:"lost_by_64"`
	MinLeadWrites int `json:"min_lead_writes"`
	MaxLeadWrites int `json:"max_lead_writes"`
	MeanLeadWrites float64 `json:"mean_lead_writes"`
	LeadWithinCadence int `json:"lead_within_cadence"`
	LeadWithinTwoCadences int `json:"lead_within_two_cadences"`
	WarningsAfterLoss int `json:"warnings_after_loss"`
}
type UP186CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP185CSeal string `json:"source_up185c_seal"`
	Cadences []int `json:"cadences"`
	Policy string `json:"policy"`
	InterventionUsed bool `json:"intervention_used"`
	WarningThresholdChanged bool `json:"warning_threshold_changed"`
	LiveActivation bool `json:"live_activation"`
	Metrics []UP186CMetric `json:"metrics"`
}
func up186cObserve(x0 *up81cAging,endangered,cadence int)(warnStep,lossStep int){
	m:=*x0
	warnStep=0;lossStep=65
	for start:=1;start<=64;start+=cadence{
		if warnStep==0{
			if up161cAdversarial(&m,endangered)<=cadence{warnStep=start}
		}
		for j:=0;j<cadence&&start+j<=64;j++{
			step:=start+j
			if up165cRealStep(&m,endangered,1013000+step,step,"no_refresh"){lossStep=step;return}
		}
	}
	return
}
func RunUP186C()(UP186CResult,error){
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	cadences:=[]int{2,4}
	res:=UP186CResult{Schema:UP186CWarningLeadSchema,Experiment:"UP-186C-warning-lead-calibration",SourceUP185CSeal:"97fd5882ceb29b815fdca2a5480822a70df9875c",Cadences:cadences,Policy:"no_refresh",InterventionUsed:false,WarningThresholdChanged:false,LiveActivation:false}
	for _,cad:=range cadences{
		m:=UP186CMetric{Cadence:cad}
		sum,n:=0,0
		for _,c:=range cohorts{for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,c);if !ok{continue};m.Arms++
			w,l:=up186cObserve(x,e,cad)
			if l<=64{m.LostBy64++}
			if w>0&&l<=64&&w<=l{
				lead:=l-w;m.WarnedBeforeLoss++;sum+=lead;n++
				if m.MinLeadWrites==0||lead<m.MinLeadWrites{m.MinLeadWrites=lead}
				if lead>m.MaxLeadWrites{m.MaxLeadWrites=lead}
				if lead<=cad{m.LeadWithinCadence++}
				if lead<=2*cad{m.LeadWithinTwoCadences++}
			}else if w>l{m.WarningsAfterLoss++}
		}}
		if n>0{m.MeanLeadWrites=float64(sum)/float64(n)}
		res.Metrics=append(res.Metrics,m)
	}
	return res,nil
}
