package unitary

const UP175BHazardSchema="wingless.up175b-heldout-hazard-calibration.v1"

type UP175BCell struct{
	Class string `json:"class"`
	RankStratum string `json:"rank_stratum"`
	TemporalState string `json:"temporal_state"`
	FrozenProbability float64 `json:"frozen_probability"`
	Slots int `json:"slots"`
	Crossings int `json:"crossings"`
	ObservedDensity float64 `json:"observed_density"`
}
type UP175BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP174BSeal string `json:"source_up174b_seal"`
	ContextPhase int `json:"context_phase"`
	WarningThreshold float64 `json:"warning_threshold"`
	BackgroundProbability float64 `json:"background_probability"`
	TotalSlots int `json:"total_slots"`
	ActualCrossings int `json:"actual_crossings"`
	PredictedCrossingsSum float64 `json:"predicted_crossings_sum"`
	BrierScore float64 `json:"brier_score"`
	TruePositive int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	FalseNegative int `json:"false_negative"`
	TrueNegative int `json:"true_negative"`
	WarningPrecision float64 `json:"warning_precision"`
	WarningRecall float64 `json:"warning_recall"`
	LeadTimeUpdates int `json:"lead_time_updates"`
	PredictorFrozenBeforeRun bool `json:"predictor_frozen_before_run"`
	AdaptiveThresholdUsed bool `json:"adaptive_threshold_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Cells []UP175BCell `json:"cells"`
}
type up175bCount struct{slots,cross int}
func up175bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func up175bRank(rank int)string{if rank<=4{return "rank_1_4"};if rank<=10{return "rank_5_10"};return "background"}
func up175bTemporal(streak int)string{if streak==1{return "current_only"};if streak==2{return "transition_to_persistent"};if streak>=3{return "established_persistent"};return "background"}
func up175bP(cls,rankState,temp string)float64{
	bg:=0.00043793793793793793
	if rankState=="background"||temp=="background"{return bg}
	switch cls{
	case "STORE":
		if rankState=="rank_1_4"{if temp=="current_only"{return 0.4358974358974359};if temp=="transition_to_persistent"{return 0.43478260869565216};return 0.25}
		if temp=="current_only"{return 0.08130081300813008};if temp=="transition_to_persistent"{return 0.09722222222222222};return 0.04242424242424243
	case "OBSERVE":
		if rankState=="rank_1_4"{if temp=="current_only"{return 0.23333333333333334};if temp=="transition_to_persistent"{return 0.4406779661016949};return 0.2251655629139073}
		if temp=="current_only"{return 0.10569105691056911};if temp=="transition_to_persistent"{return 0.15384615384615385};return 0.029069767441860465
	case "REPORT":
		if rankState=="rank_1_4"{if temp=="current_only"{return 0.3333333333333333};if temp=="transition_to_persistent"{return 0.18181818181818182};return 0.14473684210526316}
	}
	return bg
}
func RunUP175B()(UP175BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=20;threshold:=0.20
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	acc:=map[string]*up175bCount{};predSum,brier:=0.0,0.0;total,cross,tp,fp,fn,tn:=0,0,0,0,0,0
	for _,subject:=range up130bNewNames{
		base:=up169bPostReport(common,subject,epoch,o,r)
		for _,path:=range paths{
			g:=*base;streak:=make([]int,120)
			for _,idx:=range path{
				before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
				ranks:=make([]int,len(before))
				for i:=range before{ranks[i]=up166bRank(before,i);if ranks[i]<=band{streak[i]++}else{streak[i]=0}}
				items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
				for i:=range before{
					rs:=up175bRank(ranks[i]);ts:=up175bTemporal(streak[i]);p:=up175bP(cls,rs,ts);changed:=before[i].correct!=after[i].correct
					total++;predSum+=p;y:=0.0;if changed{cross++;y=1};d:=p-y;brier+=d*d
					warn:=p>=threshold
					if warn&&changed{tp++}else if warn&&!changed{fp++}else if !warn&&changed{fn++}else{tn++}
					key:=cls+"|"+rs+"|"+ts
					if acc[key]==nil{acc[key]=&up175bCount{}}
					acc[key].slots++;if changed{acc[key].cross++}
				}
			}
		}
	}
	res:=UP175BResult{Schema:UP175BHazardSchema,Experiment:"UP-175B-heldout-hazard-calibration",SourceUP174BSeal:"5fdc7debe728fe385baa1f4c9e28ab628220da6f",ContextPhase:20,WarningThreshold:threshold,BackgroundProbability:0.00043793793793793793,TotalSlots:total,ActualCrossings:cross,PredictedCrossingsSum:predSum,BrierScore:brier/float64(total),TruePositive:tp,FalsePositive:fp,FalseNegative:fn,TrueNegative:tn,WarningPrecision:up175bRate(tp,tp+fp),WarningRecall:up175bRate(tp,tp+fn),LeadTimeUpdates:1,PredictorFrozenBeforeRun:true,AdaptiveThresholdUsed:false,MaintenanceTriggered:false}
	classes:=[]string{"STORE","OBSERVE","REPORT"};ranks:=[]string{"rank_1_4","rank_5_10","background"};temps:=[]string{"current_only","transition_to_persistent","established_persistent","background"}
	for _,cls:=range classes{for _,rs:=range ranks{for _,ts:=range temps{key:=cls+"|"+rs+"|"+ts;if x:=acc[key];x!=nil{res.Cells=append(res.Cells,UP175BCell{Class:cls,RankStratum:rs,TemporalState:ts,FrozenProbability:up175bP(cls,rs,ts),Slots:x.slots,Crossings:x.cross,ObservedDensity:up175bRate(x.cross,x.slots)})}}}}
	return res,nil
}
