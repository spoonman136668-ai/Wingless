package unitary

const UPLM0LInterleavedClausesSchema = "wingless.up-lm0l-interleaved-clauses.v1"

type UPLM0LInterleavedClausesResult struct {
	Schema                          string                   `json:"schema"`
	Experiment                      string                   `json:"experiment"`
	SourceUPLM0KSeal                string                   `json:"source_up_lm0k_seal"`
	StateDimension                  int                      `json:"state_dimension"`
	ExactRecallCap                  int                      `json:"exact_recall_cap"`
	ClassifierEpochs                int                      `json:"classifier_epochs"`
	ClassifierLearningRate          float64                  `json:"classifier_learning_rate"`
	ExplicitTypeAtLearnedInference  bool                     `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool                    `json:"value_bytes_at_classifier_inference"`
	AttentionUsed                   bool                     `json:"attention_used"`
	FutureOracleUsed                bool                     `json:"future_oracle_used"`
	InterleavedRetrainingUsed       bool                     `json:"interleaved_retraining_used"`
	ClassifierMetrics               []UPLM0JClassifierMetric `json:"classifier_metrics"`
	Metrics                         []UPLM0JMetric           `json:"metrics"`
}

func uplm0lInterleavedHeldout() []uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
				initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
				latest:=initial
				uc:=uplm0eUpdateCount(p)
				s:=""
				var targets [4]int
				for i:=0;i<4;i++ {
					s+=ns[i]+" stores "+initial[i]+". "
					obs:=values[(v+i+3)%6]
					s+=ns[(i+p)%4]+" observes "+obs+". "
					if i<uc {
						latest[i]=values[(v+i+2)%6]
						s+=ns[i]+" stores "+latest[i]+". "
					}
					s+=ns[i]+" reports "
					targets[i]=len(s)
					s+=latest[i]+"."
					if i<3 { s+=" " }
				}
				s+="\n"
				out=append(out,uplm0fExample{text:s,targetPos:targets,updateCount:uc})
			}
		}
	}
	return out
}

func RunUPLM0L()(UPLM0LInterleavedClausesResult,error){
	train,_,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0jTrainClassifier()
	held:=uplm0lInterleavedHeldout()

	result:=UPLM0LInterleavedClausesResult{
		Schema:UPLM0LInterleavedClausesSchema,
		Experiment:"UP-LM0L-interleaved-clauses",
		SourceUPLM0KSeal:"73584e206a72ae8d3b3c7ebd2c24596f92ee6662",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,
		ExplicitTypeAtLearnedInference:false,ValueBytesAtClassifierInference:false,
		AttentionUsed:false,FutureOracleUsed:false,InterleavedRetrainingUsed:false,
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
	result.Metrics=[]UPLM0JMetric{
		uplm0kEvaluate(model,classifier,held,1,false,"interleaved_stream1"),
		uplm0kEvaluate(model,classifier,held,4,false,"interleaved_stream4"),
		uplm0kEvaluate(model,classifier,held,1,true,"interleaved_stream1"),
		uplm0kEvaluate(model,classifier,held,4,true,"interleaved_stream4"),
	}
	return result,nil
}
