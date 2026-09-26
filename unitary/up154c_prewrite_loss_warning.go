package unitary

const UP154CLossWarningSchema="wingless.up154c-prewrite-loss-warning.v1"

type UP154CPoint struct{
	InitialHand int `json:"initial_hand"`
	Step int `json:"step"`
	IncomingKey int `json:"incoming_key"`
	IncomingWasPresent bool `json:"incoming_was_present"`
	ScheduledRefreshKey int `json:"scheduled_refresh_key"`
	RefreshHit bool `json:"refresh_hit"`
	HandBeforeWrite int `json:"hand_before_write"`
	PredictedSlot int `json:"predicted_slot"`
	PredictedDurableLoss bool `json:"predicted_durable_loss"`
	ActualDurableLoss bool `json:"actual_durable_loss"`
	OriginalBefore int `json:"original_before"`
	OriginalAfter int `json:"original_after"`
}
type UP154CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP153CSeal string `json:"source_up153c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	InitialHands []int `json:"initial_hands"`
	StepsPerArm int `json:"steps_per_arm"`
	NovelWritesPerArm int `json:"novel_writes_per_arm"`
	RepeatWritesPerArm int `json:"repeat_writes_per_arm"`
	TruePositive int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	FalseNegative int `json:"false_negative"`
	TrueNegative int `json:"true_negative"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Accuracy float64 `json:"accuracy"`
	PredictorUsesSemanticClass bool `json:"predictor_uses_semantic_class"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	InterventionTriggered bool `json:"intervention_triggered"`
	Points []UP154CPoint `json:"points"`
}
func up154cRefresh(step int)(int,bool){switch step{case 4:return 3,true;case 8:return 7,true;case 12:return 11,true;case 16:return 15,true;case 20:return 0,true;case 24:return 4,true};return -1,false}
func up154cOriginalCount(x *up81cAging)int{n:=0;for k:=0;k<16;k++{if x.find(k)>=0{n++}};return n}
func up154cRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func up154cRun(initialHand int)([]UP154CPoint,int,int,int,int){
	x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
	out:=[]UP154CPoint{};tp,fp,fn,tn:=0,0,0,0;novelCount:=0;lastNovel:=-1
	for step:=1;step<=24;step++{
		rk,scheduled:=up154cRefresh(step);hit:=false;if scheduled{_,hit=x.query(rk)}
		key:=0
		if step%3==0&&lastNovel>=0{key=lastNovel}else{novelCount++;key=5000+novelCount;lastNovel=key}
		present:=x.find(key)>=0;predSlot:=-1;warn:=false
		if !present{
			predSlot=up151cPredict(x.hand,x.age)
			if x.entries[predSlot].used{old:=x.entries[predSlot].key;warn=old>=0&&old<16}
		}
		before:=up154cOriginalCount(x);hb:=x.hand;x.write(key,step%3);after:=up154cOriginalCount(x);actual:=after<before
		if warn&&actual{tp++}else if warn&&!actual{fp++}else if !warn&&actual{fn++}else{tn++}
		if !scheduled{rk=-1}
		out=append(out,UP154CPoint{InitialHand:initialHand,Step:step,IncomingKey:key,IncomingWasPresent:present,ScheduledRefreshKey:rk,RefreshHit:hit,HandBeforeWrite:hb,PredictedSlot:predSlot,PredictedDurableLoss:warn,ActualDurableLoss:actual,OriginalBefore:before,OriginalAfter:after})
	}
	return out,tp,fp,fn,tn
}
func RunUP154C()(UP154CResult,error){
	hands:=[]int{0,4,8,12};res:=UP154CResult{Schema:UP154CLossWarningSchema,Experiment:"UP-154C-prewrite-loss-warning",SourceUP153CSeal:"da9e256704ddafc9cb20c990f347e9bb59cf2181",ExactRecallCap:16,InitialHands:hands,StepsPerArm:24,NovelWritesPerArm:16,RepeatWritesPerArm:8,PredictorUsesSemanticClass:false,FutureOracleUsed:false,InterventionTriggered:false}
	for _,h:=range hands{pts,tp,fp,fn,tn:=up154cRun(h);res.Points=append(res.Points,pts...);res.TruePositive+=tp;res.FalsePositive+=fp;res.FalseNegative+=fn;res.TrueNegative+=tn}
	res.Precision=up154cRate(res.TruePositive,res.TruePositive+res.FalsePositive);res.Recall=up154cRate(res.TruePositive,res.TruePositive+res.FalseNegative);res.Accuracy=up154cRate(res.TruePositive+res.TrueNegative,len(res.Points))
	return res,nil
}
