package unitary

const UP151CForecastSchema="wingless.up151c-state-eviction-forecast.v1"

type UP151CPoint struct{
	InitialHand int `json:"initial_hand"`
	Admission int `json:"admission"`
	ScheduledRefreshKey int `json:"scheduled_refresh_key"`
	RefreshHit bool `json:"refresh_hit"`
	HandBeforeWrite int `json:"hand_before_write"`
	PredictedSlot int `json:"predicted_slot"`
	ActualSlot int `json:"actual_slot"`
	PredictionCorrect bool `json:"prediction_correct"`
	EvictedKey int `json:"evicted_key"`
	OriginalFactsRemaining int `json:"original_facts_remaining"`
}
type UP151CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP150CSeal string `json:"source_up150c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	InitialHands []int `json:"initial_hands"`
	AdmissionsPerArm int `json:"admissions_per_arm"`
	SparseRefreshAdmissions []int `json:"sparse_refresh_admissions"`
	SparseRefreshKeys []int `json:"sparse_refresh_keys"`
	PredictorUsesHand bool `json:"predictor_uses_hand"`
	PredictorUsesAges bool `json:"predictor_uses_ages"`
	PredictorUsesIdentity bool `json:"predictor_uses_identity"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	MemoryPolicyChanged bool `json:"memory_policy_changed"`
	CorrectPredictions int `json:"correct_predictions"`
	TotalPredictions int `json:"total_predictions"`
	PredictionAccuracy float64 `json:"prediction_accuracy"`
	Points []UP151CPoint `json:"points"`
}
func up151cPredict(hand int,age [16]uint8)int{
	h:=hand
	for{
		if age[h]==0{return h}
		age[h]--
		h=(h+1)%16
	}
}
func up151cRefreshKey(admission int)(int,bool){
	switch admission{case 3:return 3,true;case 6:return 7,true;case 9:return 11,true;case 12:return 15,true;case 15:return 0,true}
	return -1,false
}
func up151cRun(initialHand int)([]UP151CPoint,int){
	x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
	out:=[]UP151CPoint{};correct:=0
	for j:=1;j<=16;j++{
		rk,scheduled:=up151cRefreshKey(j);hit:=false
		if scheduled{_,hit=x.query(rk)}
		handBefore:=x.hand;pred:=up151cPredict(x.hand,x.age)
		evictedKey:=-1;if x.entries[pred].used{evictedKey=x.entries[pred].key}
		newKey:=2000+j;x.write(newKey,j%3);actual:=x.find(newKey);ok:=pred==actual;if ok{correct++}
		remaining:=0;for k:=0;k<16;k++{if x.find(k)>=0{remaining++}}
		if !scheduled{rk=-1}
		out=append(out,UP151CPoint{InitialHand:initialHand,Admission:j,ScheduledRefreshKey:rk,RefreshHit:hit,HandBeforeWrite:handBefore,PredictedSlot:pred,ActualSlot:actual,PredictionCorrect:ok,EvictedKey:evictedKey,OriginalFactsRemaining:remaining})
	}
	return out,correct
}
func RunUP151C()(UP151CResult,error){
	hands:=[]int{0,4,8,12};res:=UP151CResult{Schema:UP151CForecastSchema,Experiment:"UP-151C-state-eviction-forecast",SourceUP150CSeal:"e758b9094ce0a113dee157424b3e5e8415c582a6",ExactRecallCap:16,InitialHands:hands,AdmissionsPerArm:16,SparseRefreshAdmissions:[]int{3,6,9,12,15},SparseRefreshKeys:[]int{3,7,11,15,0},PredictorUsesHand:true,PredictorUsesAges:true,PredictorUsesIdentity:false,FutureOracleUsed:false,MemoryPolicyChanged:false}
	for _,h:=range hands{pts,c:=up151cRun(h);res.Points=append(res.Points,pts...);res.CorrectPredictions+=c}
	res.TotalPredictions=len(res.Points);res.PredictionAccuracy=float64(res.CorrectPredictions)/float64(res.TotalPredictions)
	return res,nil
}
