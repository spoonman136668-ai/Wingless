package unitary

const UPLM1JFourthFamilySchema = "wingless.up-lm1j-fourth-family-boundary.v1"

type UPLM1JRouterMetric struct {
	Split    string  `json:"split"`
	Accuracy float64 `json:"accuracy"`
	Examples int     `json:"examples"`
}

type UPLM1JFourthFamilyResult struct {
	Schema                     string               `json:"schema"`
	Experiment                 string               `json:"experiment"`
	SourceUPLM1ISeal           string               `json:"source_up_lm1i_seal"`
	StateDimension             int                  `json:"state_dimension"`
	ExactRecallCap             int                  `json:"exact_recall_cap"`
	BasePretrainEpochs         int                  `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                  `json:"joint_interleaved_epochs"`
	ThreeFamilyAdaptEpochs     int                  `json:"three_family_adapt_epochs"`
	FourthRouterEpochs         int                  `json:"fourth_router_epochs"`
	LearningRate               float64              `json:"learning_rate"`
	HalfStepLearningRate       float64              `json:"half_step_learning_rate"`
	FourthFamilyByteTraining   bool                 `json:"fourth_family_byte_training"`
	RecurrentParametersTrained bool                 `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                 `json:"attention_used"`
	FutureOracleUsed           bool                 `json:"future_oracle_used"`
	RouterMetrics              []UPLM1JRouterMetric `json:"router_metrics"`
	Metrics                    []UPLM0JMetric       `json:"metrics"`
}

var uplm1jFourthVerbs=[]string{"retains","inspects","states"}

func uplm1jTrainClassifier()*uplm0jClassifier{
	c:=uplm1cTrainClassifier()
	names:=uplm0gNames()
	prior:=[][]string{
		{"stores","observes","reports"},
		{"saves","sees","recalls"},
		{"archives","notices","recounts"},
	}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			for class,verb:=range uplm1jFourthVerbs {
				uplm0nTrainStep(c,uplm0nAnchor(names[ni],verb),class)
			}
		}
		for _,family:=range prior {
			for class,verb:=range family {
				uplm0nTrainStep(c,uplm0nAnchor(names[0],verb),class)
			}
		}
	}
	return c
}

func uplm1jRouterMetric(c *uplm0jClassifier) UPLM1JRouterMetric {
	names:=uplm0gNames()
	hits,total:=0,0
	for ni:=4;ni<6;ni++ {
		for class,verb:=range uplm1jFourthVerbs {
			total++
			if uplm0jPredict(c,uplm0nAnchor(names[ni],verb))==class { hits++ }
		}
	}
	return UPLM1JRouterMetric{Split:"heldout_fourth_family_subjects",Accuracy:float64(hits)/float64(total),Examples:total}
}

func uplm1jExample(n,v,p int,order string) uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial
	uc:=uplm0eUpdateCount(p)
	s:=""
	var targets [4]int

	storeInitial:=func(i int){ s+=ns[i]+" retains "+initial[i]+". " }
	observe:=func(i int){
		obs:=values[(v+i+3)%6]
		s+=ns[(i+p)%4]+" inspects "+obs+". "
	}
	update:=func(i int){
		if i>=uc { return }
		latest[i]=values[(v+i+2)%6]
		s+=ns[i]+" retains "+latest[i]+". "
	}
	report:=func(i int,last bool){
		s+=ns[i]+" states "
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

func uplm1jHeldout(order string) []uplm0fExample {
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				out=append(out,uplm1jExample(n,v,p,order))
			}
		}
	}
	return out
}

func RunUPLM1J()(UPLM1JFourthFamilyResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(model,baseTrain,paraTrain,20)
	model,_=uplm1dExtendAlphabet(model,thirdTrain,thirdHeld)
	uplm1hTrain(model,baseTrain,paraTrain,thirdTrain,"strang_split_base")

	classifier:=uplm1jTrainClassifier()
	result:=UPLM1JFourthFamilyResult{
		Schema:UPLM1JFourthFamilySchema,Experiment:"UP-LM1J-fourth-family-boundary",
		SourceUPLM1ISeal:"7eedf2833867fa1b154a5113c07f08cd77ae170b",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,ThreeFamilyAdaptEpochs:4,FourthRouterEpochs:20,
		LearningRate:0.08,HalfStepLearningRate:0.04,FourthFamilyByteTraining:false,
		RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	result.RouterMetrics=append(result.RouterMetrics,uplm1jRouterMetric(classifier))

	for _,stream:=range []int{1,4} {
		result.Metrics=append(result.Metrics,
			uplm1cEval(model,classifier,baseHeld,stream,"base_block_stream"+itoa(stream)),
			uplm1cEval(model,classifier,paraHeld,stream,"paraphrase_block_stream"+itoa(stream)),
			uplm1cEval(model,classifier,thirdHeld,stream,"third_block_stream"+itoa(stream)),
		)
	}
	for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
		held:=uplm1jHeldout(order)
		for _,stream:=range []int{1,4} {
			result.Metrics=append(result.Metrics,
				uplm1cEval(model,classifier,held,stream,"fourth_"+order+"_stream"+itoa(stream)),
			)
		}
	}
	return result,nil
}
