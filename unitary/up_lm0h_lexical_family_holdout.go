package unitary

const UPLM0HLexicalFamilyHoldoutSchema = "wingless.up-lm0h-lexical-family-holdout.v1"

type UPLM0HLexicalFamilyHoldoutResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM0GSeal string `json:"source_up_lm0g_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	ClassifierEpochs int `json:"classifier_epochs"`
	ClassifierLearningRate float64 `json:"classifier_learning_rate"`
	ClassifierThreshold float64 `json:"classifier_threshold"`
	ExplicitTypeAtLearnedInference bool `json:"explicit_type_at_learned_inference"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	ClassifierMetrics []UPLM0GClassifierMetric `json:"classifier_metrics"`
	Metrics []UPLM0GMetric `json:"metrics"`
}

func uplm0hTrainClassifier() *uplm0gClassifier {
	c:=&uplm0gClassifier{}
	names,values:=uplm0gNames(),uplm0gValues()
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{
			for vi:=0;vi<4;vi++{
				for et:=0;et<2;et++{
					store:=et==0
					h:=uplm0fEncode(uplm0gClause(names[ni],values[vi],store))
					y:=0.0
					if store{y=1}
					g:=c.prob(h)-y
					for i:=0;i<64;i++{c.w[i]-=0.08*g*h[i]}
					c.b-=0.08*g
				}
			}
		}
	}
	return c
}

func uplm0hClassifierEval(c *uplm0gClassifier,split string,nameStart,nameEnd,valueStart,valueEnd int) UPLM0GClassifierMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	names,values:=uplm0gNames(),uplm0gValues()
	for ni:=nameStart;ni<nameEnd;ni++{
		for vi:=valueStart;vi<valueEnd;vi++{
			for et:=0;et<2;et++{
				store:=et==0
				pred:=c.prob(uplm0fEncode(uplm0gClause(names[ni],values[vi],store)))>=0.5
				total++
				if pred==store{hits++}
				if pred&&store{tp++}
				if pred&&!store{fp++}
				if !pred&&store{fn++}
			}
		}
	}
	precision,recall:=1.0,1.0
	if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	return UPLM0GClassifierMetric{
		Split:split,Accuracy:float64(hits)/float64(total),
		Precision:precision,Recall:recall,Examples:total,
	}
}

func RunUPLM0H()(UPLM0HLexicalFamilyHoldoutResult,error){
	train,held,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++{
		for _,ex:=range train{model.trainSentence(ex.text,0.08)}
	}
	classifier:=uplm0hTrainClassifier()
	return UPLM0HLexicalFamilyHoldoutResult{
		Schema:UPLM0HLexicalFamilyHoldoutSchema,
		Experiment:"UP-LM0H-lexical-family-holdout",
		SourceUPLM0GSeal:"e2b6dd94977bfde44c80f91480a308aae5d943ee",
		StateDimension:64,ExactRecallCap:16,
		ClassifierEpochs:20,ClassifierLearningRate:0.08,ClassifierThreshold:0.5,
		ExplicitTypeAtLearnedInference:false,AttentionUsed:false,FutureOracleUsed:false,
		ClassifierMetrics:[]UPLM0GClassifierMetric{
			uplm0hClassifierEval(classifier,"train_domain",0,4,0,4),
			uplm0hClassifierEval(classifier,"heldout_subject_family",4,6,0,4),
			uplm0hClassifierEval(classifier,"heldout_value_family",0,4,4,6),
			uplm0hClassifierEval(classifier,"heldout_joint_family",4,6,4,6),
		},
		Metrics:[]UPLM0GMetric{
			uplm0gEvaluate(model,classifier,held,false,false,"heldout_recombination"),
			uplm0gEvaluate(model,classifier,held,true,false,"heldout_stream4"),
			uplm0gEvaluate(model,classifier,held,false,true,"heldout_recombination"),
			uplm0gEvaluate(model,classifier,held,true,true,"heldout_stream4"),
		},
	},nil
}
