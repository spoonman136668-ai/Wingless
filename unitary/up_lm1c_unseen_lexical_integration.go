package unitary

const UPLM1CUnseenLexicalSchema = "wingless.up-lm1c-unseen-lexical-integration.v1"

type UPLM1CRouterMetric struct {
	Split    string  `json:"split"`
	Accuracy float64 `json:"accuracy"`
	Examples int     `json:"examples"`
}

type UPLM1CUnseenLexicalResult struct {
	Schema                     string               `json:"schema"`
	Experiment                 string               `json:"experiment"`
	SourceUPLM1BSeal           string               `json:"source_up_lm1b_seal"`
	StateDimension             int                  `json:"state_dimension"`
	ExactRecallCap             int                  `json:"exact_recall_cap"`
	BasePretrainEpochs         int                  `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                  `json:"joint_interleaved_epochs"`
	RouterGroundingEpochs      int                  `json:"router_grounding_epochs"`
	LearningRate               float64              `json:"learning_rate"`
	ByteModelThirdFamilyTraining bool               `json:"byte_model_third_family_training"`
	RecurrentParametersTrained bool                 `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                 `json:"attention_used"`
	FutureOracleUsed           bool                 `json:"future_oracle_used"`
	RouterMetrics              []UPLM1CRouterMetric `json:"router_metrics"`
	Metrics                    []UPLM0JMetric       `json:"metrics"`
}

func uplm1cTrainClassifier() *uplm0jClassifier {
	c:=uplm0nTrainClassifier()
	names:=uplm0gNames()
	third:=[]string{"archives","notices","recounts"}
	base:=[]string{"stores","observes","reports"}
	para:=[]string{"saves","sees","recalls"}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			for class,verb:=range third {
				uplm0nTrainStep(c,uplm0nAnchor(names[ni],verb),class)
			}
		}
		for class,verb:=range base {
			uplm0nTrainStep(c,uplm0nAnchor(names[0],verb),class)
		}
		for class,verb:=range para {
			uplm0nTrainStep(c,uplm0nAnchor(names[0],verb),class)
		}
	}
	return c
}

func uplm1cRouterAccuracy(c *uplm0jClassifier, start,end int) UPLM1CRouterMetric {
	names:=uplm0gNames()
	verbs:=[]string{"archives","notices","recounts"}
	hits,total:=0,0
	for ni:=start;ni<end;ni++ {
		for class,verb:=range verbs {
			total++
			if uplm0jPredict(c,uplm0nAnchor(names[ni],verb))==class { hits++ }
		}
	}
	return UPLM1CRouterMetric{
		Split:"heldout_third_family_subjects",
		Accuracy:float64(hits)/float64(total),
		Examples:total,
	}
}

func uplm1cExample(n,v,p int,order string) uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial
	uc:=uplm0eUpdateCount(p)
	s:=""
	var targets [4]int

	storeInitial:=func(i int){ s+=ns[i]+" archives "+initial[i]+". " }
	observe:=func(i int){
		obs:=values[(v+i+3)%6]
		s+=ns[(i+p)%4]+" notices "+obs+". "
	}
	update:=func(i int){
		if i>=uc { return }
		latest[i]=values[(v+i+2)%6]
		s+=ns[i]+" archives "+latest[i]+". "
	}
	report:=func(i int,last bool){
		s+=ns[i]+" recounts "
		targets[i]=len(s)
		s+=latest[i]+"."
		if !last { s+=" " }
	}

	switch order {
	case "block":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i) }
		for i:=0;i<4;i++ { update(i) }
		for qi:=0;qi<4;qi++ {
			idx:=(qi+p)%4
			report(idx,qi==3)
		}
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

func uplm1cHeldout(order string) []uplm0fExample {
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				out=append(out,uplm1cExample(n,v,p,order))
			}
		}
	}
	return out
}

func uplm1cEval(model *uplm0aModel,classifier *uplm0jClassifier,examples []uplm0fExample,stream int,split string) UPLM0JMetric {
	m:=uplm0nEvaluate(model,classifier,examples,stream,true,split)
	b:=uplm0oByteEval(model,examples,0,split)
	m.Top1Accuracy=b.Top1Accuracy
	m.Perplexity=b.Perplexity
	return m
}

func RunUPLM1C()(UPLM1CUnseenLexicalResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(model,baseTrain,paraTrain,20)
	classifier:=uplm1cTrainClassifier()

	result:=UPLM1CUnseenLexicalResult{
		Schema:UPLM1CUnseenLexicalSchema,
		Experiment:"UP-LM1C-unseen-lexical-integration",
		SourceUPLM1BSeal:"33b5d83b477caa5104cdb67ce2ff3a9e1f9526fc",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,RouterGroundingEpochs:20,
		LearningRate:0.08,ByteModelThirdFamilyTraining:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	result.RouterMetrics=append(result.RouterMetrics,uplm1cRouterAccuracy(classifier,4,6))

	for _,stream:=range []int{1,4} {
		result.Metrics=append(result.Metrics,
			uplm1cEval(model,classifier,baseHeld,stream,"base_block_stream"+itoa(stream)),
			uplm1cEval(model,classifier,paraHeld,stream,"paraphrase_block_stream"+itoa(stream)),
		)
	}
	for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
		held:=uplm1cHeldout(order)
		for _,stream:=range []int{1,4} {
			result.Metrics=append(result.Metrics,
				uplm1cEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream)),
			)
		}
	}
	return result,nil
}

func itoa(v int) string {
	if v==1 { return "1" }
	if v==4 { return "4" }
	return "0"
}
