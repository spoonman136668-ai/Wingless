package unitary

const UPLM0MOrderPermutationsSchema = "wingless.up-lm0m-order-permutations.v1"

type UPLM0MOrderPermutationsResult struct {
	Schema                          string                   `json:"schema"`
	Experiment                      string                   `json:"experiment"`
	SourceUPLM0LSeal                string                   `json:"source_up_lm0l_seal"`
	StateDimension                  int                      `json:"state_dimension"`
	ExactRecallCap                  int                      `json:"exact_recall_cap"`
	ClassifierEpochs                int                      `json:"classifier_epochs"`
	ClassifierLearningRate          float64                  `json:"classifier_learning_rate"`
	ExplicitTypeAtLearnedInference  bool                     `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool                    `json:"value_bytes_at_classifier_inference"`
	AttentionUsed                   bool                     `json:"attention_used"`
	FutureOracleUsed                bool                     `json:"future_oracle_used"`
	PermutationRetrainingUsed       bool                     `json:"permutation_retraining_used"`
	ClassifierMetrics               []UPLM0JClassifierMetric `json:"classifier_metrics"`
	Metrics                         []UPLM0JMetric           `json:"metrics"`
}

func uplm0mExample(n,v,p int,order string) uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial
	uc:=uplm0eUpdateCount(p)
	s:=""
	var targets [4]int

	storeInitial:=func(i int){ s+=ns[i]+" stores "+initial[i]+". " }
	observe:=func(i int){
		obs:=values[(v+i+3)%6]
		s+=ns[(i+p)%4]+" observes "+obs+". "
	}
	update:=func(i int){
		if i>=uc { return }
		latest[i]=values[(v+i+2)%6]
		s+=ns[i]+" stores "+latest[i]+". "
	}
	report:=func(i int,last bool){
		s+=ns[i]+" reports "
		targets[i]=len(s)
		s+=latest[i]+"."
		if !last { s+=" " }
	}

	switch order {
	case "per_name":
		for i:=0;i<4;i++ {
			storeInitial(i); observe(i); update(i); report(i,i==3)
		}
	case "paired_names":
		for pair:=0;pair<4;pair+=2 {
			for i:=pair;i<pair+2;i++ { storeInitial(i) }
			for i:=pair;i<pair+2;i++ { observe(i) }
			for i:=pair;i<pair+2;i++ { update(i) }
			for i:=pair;i<pair+2;i++ { report(i,pair==2&&i==3) }
		}
	case "stores_then_local_reports":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ {
			observe(i); update(i); report(i,i==3)
		}
	case "reverse_report_tail":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i) }
		for i:=0;i<4;i++ { update(i) }
		for i:=3;i>=0;i-- { report(i,i==0) }
	}
	s+="\n"
	return uplm0fExample{text:s,targetPos:targets,updateCount:uc}
}

func uplm0mHeldout(order string) []uplm0fExample {
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				out=append(out,uplm0mExample(n,v,p,order))
			}
		}
	}
	return out
}

func RunUPLM0M()(UPLM0MOrderPermutationsResult,error){
	train,_,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0jTrainClassifier()

	result:=UPLM0MOrderPermutationsResult{
		Schema:UPLM0MOrderPermutationsSchema,
		Experiment:"UP-LM0M-order-permutations",
		SourceUPLM0LSeal:"2e75e78dc180e6b976572e63f6cb193552e4e6b4",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,
		ExplicitTypeAtLearnedInference:false,ValueBytesAtClassifierInference:false,
		AttentionUsed:false,FutureOracleUsed:false,PermutationRetrainingUsed:false,
	}
	for _,split:=range []struct{name string;start,end int}{
		{"train_subjects",0,4},{"heldout_subjects",4,6},
	}{
		result.ClassifierMetrics=append(result.ClassifierMetrics,
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,-1,"all"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jStore,"STORE"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jObserve,"OBSERVE"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jReport,"REPORT"),
		)
	}
	for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
		held:=uplm0mHeldout(order)
		for _,arm:=range []struct{name string;learned bool}{
			{"explicit_event_routing",false},{"learned_prefix_threeway",true},
		}{
			result.Metrics=append(result.Metrics,
				uplm0kEvaluate(model,classifier,held,1,arm.learned,order+"_stream1"),
				uplm0kEvaluate(model,classifier,held,4,arm.learned,order+"_stream4"),
			)
		}
	}
	return result,nil
}
