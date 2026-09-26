package unitary

const UP155CCohortSchema="wingless.up155c-protected-cohort-warning.v1"

type UP155CCohortMetric struct{
	Cohort string `json:"cohort"`
	Keys []int `json:"keys"`
	TruePositive int `json:"true_positive"`
	FalsePositive int `json:"false_positive"`
	FalseNegative int `json:"false_negative"`
	TrueNegative int `json:"true_negative"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Accuracy float64 `json:"accuracy"`
}
type UP155CResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP154CSeal string `json:"source_up154c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	InitialHands []int `json:"initial_hands"`
	StepsPerArm int `json:"steps_per_arm"`
	Cohorts int `json:"cohorts"`
	PredictorUsesSemanticClass bool `json:"predictor_uses_semantic_class"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	InterventionTriggered bool `json:"intervention_triggered"`
	Metrics []UP155CCohortMetric `json:"metrics"`
}
type up155cAgg struct{tp,fp,fn,tn int}
func up155cRefresh(step int)(int,bool){switch step{case 5:return 2,true;case 10:return 6,true;case 15:return 10,true;case 20:return 14,true};return -1,false}
func up155cHas(keys []int,k int)bool{for _,x:=range keys{if x==k{return true}};return false}
func up155cPresent(x *up81cAging,keys []int)int{n:=0;for _,k:=range keys{if x.find(k)>=0{n++}};return n}
func RunUP155C()(UP155CResult,error){
	hands:=[]int{0,4,8,12}
	names:=[]string{"cohort_A","cohort_B","cohort_C","cohort_D"}
	cohorts:=[][]int{{0,4,8,12},{1,5,9,13},{2,6,10,14},{3,7,11,15}}
	aggs:=make([]up155cAgg,len(cohorts))
	for _,initialHand:=range hands{
		x:=&up81cAging{};for _,d:=range up145cDurable(){x.write(d.key,d.class)};x.hand=initialHand
		novelCount:=0;lastNovel:=-1
		for step:=1;step<=24;step++{
			if rk,ok:=up155cRefresh(step);ok{x.query(rk)}
			key:=0
			if step%3==0&&lastNovel>=0{key=lastNovel}else{novelCount++;key=6000+novelCount;lastNovel=key}
			present:=x.find(key)>=0;predSlot:=-1;predOld:=-1
			if !present{predSlot=up151cPredict(x.hand,x.age);if x.entries[predSlot].used{predOld=x.entries[predSlot].key}}
			before:=make([]int,len(cohorts));for i,c:=range cohorts{before[i]=up155cPresent(x,c)}
			warns:=make([]bool,len(cohorts));for i,c:=range cohorts{warns[i]=!present&&predOld>=0&&up155cHas(c,predOld)}
			x.write(key,step%3)
			for i,c:=range cohorts{
				after:=up155cPresent(x,c);actual:=after<before[i];a:=&aggs[i]
				if warns[i]&&actual{a.tp++}else if warns[i]&&!actual{a.fp++}else if !warns[i]&&actual{a.fn++}else{a.tn++}
			}
		}
	}
	res:=UP155CResult{Schema:UP155CCohortSchema,Experiment:"UP-155C-protected-cohort-warning",SourceUP154CSeal:"6dcfd9ecaf4d58f989943517a14cb545f409b086",ExactRecallCap:16,InitialHands:hands,StepsPerArm:24,Cohorts:4,PredictorUsesSemanticClass:false,FutureOracleUsed:false,InterventionTriggered:false}
	for i,a:=range aggs{total:=a.tp+a.fp+a.fn+a.tn;res.Metrics=append(res.Metrics,UP155CCohortMetric{Cohort:names[i],Keys:cohorts[i],TruePositive:a.tp,FalsePositive:a.fp,FalseNegative:a.fn,TrueNegative:a.tn,Precision:up154cRate(a.tp,a.tp+a.fp),Recall:up154cRate(a.tp,a.tp+a.fn),Accuracy:up154cRate(a.tp+a.tn,total)})}
	return res,nil
}
