package unitary

import "fmt"

const UPLM1PFifthFamilySchema = "wingless.up-lm1p-fifth-family-scale.v1"

type UPLM1PByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1PSummary struct {
	Arm                   string  `json:"arm"`
	MinHeldoutAccuracy    float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy   float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
}

type UPLM1PRouterMetric struct {
	Family        string  `json:"family"`
	Split         string  `json:"split"`
	Accuracy      float64 `json:"accuracy"`
	StoreRecall   float64 `json:"store_recall"`
	ObserveRecall float64 `json:"observe_recall"`
	ReportRecall  float64 `json:"report_recall"`
	Examples      int     `json:"examples"`
}

type UPLM1PIntegratedMetric struct {
	Arm    string       `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1PFifthFamilyResult struct {
	Schema                     string                  `json:"schema"`
	Experiment                 string                  `json:"experiment"`
	SourceUPLM1OSeal           string                  `json:"source_up_lm1o_seal"`
	StateDimension             int                     `json:"state_dimension"`
	ExactRecallCap             int                     `json:"exact_recall_cap"`
	FifthRouterEpochs          int                     `json:"fifth_router_epochs"`
	AdaptationEpochs           int                     `json:"adaptation_epochs"`
	FullStepLearningRate       float64                 `json:"full_step_learning_rate"`
	HalfStepLearningRate       float64                 `json:"half_step_learning_rate"`
	RecurrentParametersTrained bool                    `json:"recurrent_parameters_trained"`
	RecallCapChanged           bool                    `json:"recall_cap_changed"`
	AttentionUsed              bool                    `json:"attention_used"`
	FutureOracleUsed           bool                    `json:"future_oracle_used"`
	FifthCorrectionUsed        bool                    `json:"fifth_correction_used"`
	ByteMetrics                []UPLM1PByteMetric      `json:"byte_metrics"`
	Summaries                  []UPLM1PSummary         `json:"summaries"`
	RouterMetrics              []UPLM1PRouterMetric    `json:"router_metrics"`
	IntegratedMetrics          []UPLM1PIntegratedMetric `json:"integrated_metrics"`
}

var uplm1pFifthVerbs=[]string{"banks","surveys","declares"}

func uplm1pExample(n,v,p int,order string) uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial
	uc:=uplm0eUpdateCount(p)
	s:=""
	var targets [4]int

	storeInitial:=func(i int){ s+=ns[i]+" banks "+initial[i]+". " }
	observe:=func(i int){
		obs:=values[(v+i+3)%6]
		s+=ns[(i+p)%4]+" surveys "+obs+". "
	}
	update:=func(i int){
		if i>=uc { return }
		latest[i]=values[(v+i+2)%6]
		s+=ns[i]+" banks "+latest[i]+". "
	}
	report:=func(i int,last bool){
		s+=ns[i]+" declares "
		targets[i]=len(s)
		s+=latest[i]+"."
		if !last { s+=" " }
	}

	switch order {
	case "block":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i) }
		for i:=0;i<4;i++ { update(i) }
		for qi:=0;qi<4;qi++ { idx:=(qi+p)%4;report(idx,qi==3) }
	case "per_name":
		for i:=0;i<4;i++ { storeInitial(i);observe(i);update(i);report(i,i==3) }
	case "paired_names":
		for pair:=0;pair<4;pair+=2 {
			for i:=pair;i<pair+2;i++ { storeInitial(i) }
			for i:=pair;i<pair+2;i++ { observe(i) }
			for i:=pair;i<pair+2;i++ { update(i) }
			for i:=pair;i<pair+2;i++ { report(i,pair==2&&i==3) }
		}
	case "stores_then_local_reports":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i);update(i);report(i,i==3) }
	case "reverse_report_tail":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i) }
		for i:=0;i<4;i++ { update(i) }
		for i:=3;i>=0;i-- { report(i,i==0) }
	}
	s+="\n"
	return uplm0fExample{text:s,targetPos:targets,updateCount:uc}
}

func uplm1pCorpus()(train,held []uplm0fExample) {
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				ex:=uplm1pExample(n,v,p,"block")
				if (n+2*v+p)%3!=2 { train=append(train,ex) } else { held=append(held,ex) }
			}
		}
	}
	return
}

func uplm1pHeldout(order string) []uplm0fExample {
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				out=append(out,uplm1pExample(n,v,p,order))
			}
		}
	}
	return out
}

func uplm1pTrainClassifier(d [64]float64)*uplm0jClassifier {
	c:=uplm1lTrainClassifier("states_store_orthogonal",d)
	names:=uplm0gNames()
	prior:=[][]string{
		{"stores","observes","reports"},
		{"saves","sees","recalls"},
		{"archives","notices","recounts"},
	}
	fourth:=[]string{"retains","inspects","states"}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			for class,verb:=range uplm1pFifthVerbs {
				uplm0nTrainStep(c,uplm0nAnchor(names[ni],verb),class)
			}
		}
		for _,family:=range prior {
			for class,verb:=range family {
				uplm0nTrainStep(c,uplm0nAnchor(names[0],verb),class)
			}
		}
		for class,verb:=range fourth {
			uplm1lTrainStep(c,"states_store_orthogonal",names[0],verb,class,d)
		}
	}
	return c
}

func uplm1pRouterEval(c *uplm0jClassifier,d [64]float64,family,split string,names,verbs []string) UPLM1PRouterMetric {
	hits,total:=0,0
	classHits:=[3]int{}
	classTotal:=[3]int{}
	for _,name:=range names {
		for class,verb:=range verbs {
			pred:=uplm1mPredict(c,name,verb,d)
			total++;classTotal[class]++
			if pred==class { hits++;classHits[class]++ }
		}
	}
	return UPLM1PRouterMetric{
		Family:family,Split:split,Accuracy:float64(hits)/float64(total),
		StoreRecall:float64(classHits[0])/float64(classTotal[0]),
		ObserveRecall:float64(classHits[1])/float64(classTotal[1]),
		ReportRecall:float64(classHits[2])/float64(classTotal[2]),Examples:total,
	}
}

func uplm1pTrainByteArm(model *uplm0aModel,arm string,families [5][]uplm0fExample) {
	if arm=="no_fifth_adaptation" { return }
	for epoch:=0;epoch<4;epoch++ {
		for i:=range families[0] {
			switch arm {
			case "five_family_cyclic":
				for f:=0;f<5;f++ { model.trainSentence(families[f][i].text,0.08) }
			case "five_family_rotating_palindromic":
				center:=i%5
				others:=[4]int{(center+1)%5,(center+2)%5,(center+3)%5,(center+4)%5}
				for j:=0;j<4;j++ { model.trainSentence(families[others[j]][i].text,0.04) }
				model.trainSentence(families[center][i].text,0.08)
				for j:=3;j>=0;j-- { model.trainSentence(families[others[j]][i].text,0.04) }
			}
		}
	}
}

func uplm1pSummary(arm string,metrics []UPLM1PByteMetric) UPLM1PSummary {
	vals:=[]float64{}
	for _,m:=range metrics {
		if m.Arm==arm { vals=append(vals,m.Top1Accuracy) }
	}
	min,max,sum:=1.0,0.0,0.0
	for _,v:=range vals {
		if v<min { min=v }
		if v>max { max=v }
		sum+=v
	}
	return UPLM1PSummary{Arm:arm,MinHeldoutAccuracy:min,MeanHeldoutAccuracy:sum/float64(len(vals)),HeldoutAccuracySpread:max-min}
}

func RunUPLM1P()(UPLM1PFifthFamilyResult,error) {
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1PFifthFamilyResult{},err }

	fourFamilies:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",fourFamilies)

	fifthTrain,fifthHeld:=uplm1pCorpus()
	if len(baseTrain)!=len(fifthTrain) {
		return UPLM1PFifthFamilyResult{},fmt.Errorf("fifth-family training count mismatch: base=%d fifth=%d",len(baseTrain),len(fifthTrain))
	}
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	d:=uplm1lStoreDirection()
	classifier:=uplm1pTrainClassifier(d)

	result:=UPLM1PFifthFamilyResult{
		Schema:UPLM1PFifthFamilySchema,Experiment:"UP-LM1P-fifth-family-scale",
		SourceUPLM1OSeal:"fdfc618ce6f6a323683714870ae1ff6b6eeba7c8",
		StateDimension:64,ExactRecallCap:16,FifthRouterEpochs:20,AdaptationEpochs:4,
		FullStepLearningRate:0.08,HalfStepLearningRate:0.04,
		RecurrentParametersTrained:false,RecallCapChanged:false,AttentionUsed:false,FutureOracleUsed:false,FifthCorrectionUsed:false,
	}

	names:=uplm0gNames()
	result.RouterMetrics=append(result.RouterMetrics,
		uplm1pRouterEval(classifier,d,"base","heldout",names[4:6],[]string{"stores","observes","reports"}),
		uplm1pRouterEval(classifier,d,"paraphrase","heldout",names[4:6],[]string{"saves","sees","recalls"}),
		uplm1pRouterEval(classifier,d,"third","heldout",names[4:6],[]string{"archives","notices","recounts"}),
		uplm1pRouterEval(classifier,d,"fourth","heldout",names[4:6],[]string{"retains","inspects","states"}),
		uplm1pRouterEval(classifier,d,"fifth","train",names[:4],uplm1pFifthVerbs),
		uplm1pRouterEval(classifier,d,"fifth","heldout",names[4:6],uplm1pFifthVerbs),
		uplm1pRouterEval(classifier,d,"fifth","unseen",uplm1kUnseenNames,uplm1pFifthVerbs),
	)

	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	heldSets:=[]struct{name string;data []uplm0fExample}{
		{"base_heldout",baseHeld},{"paraphrase_heldout",paraHeld},{"third_heldout",thirdHeld},{"fourth_heldout",fourthHeld},{"fifth_heldout",fifthHeld},
	}
	arms:=[]string{"no_fifth_adaptation","five_family_cyclic","five_family_rotating_palindromic"}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(start)
		uplm1pTrainByteArm(model,arm,families)
		for _,x:=range heldSets {
			m:=uplm0oByteEval(model,x.data,0,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1PByteMetric{Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		result.Summaries=append(result.Summaries,uplm1pSummary(arm,result.ByteMetrics))

		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1pHeldout(order)
			for _,stream:=range []int{1,4} {
				result.IntegratedMetrics=append(result.IntegratedMetrics,UPLM1PIntegratedMetric{
					Arm:arm,
					Metric:uplm1mEvaluate(model,classifier,d,held,stream,"fifth_"+order+"_stream"+itoa(stream)),
				})
			}
		}
	}
	return result,nil
}
