package unitary

const UPLM1HSymmetricSplitSchema = "wingless.up-lm1h-symmetric-operator-split.v1"

type UPLM1HMetric struct {
	Arm             string  `json:"arm"`
	Split           string  `json:"split"`
	Top1Accuracy    float64 `json:"top1_accuracy"`
	Perplexity      float64 `json:"perplexity"`
}

type UPLM1HArmSummary struct {
	Arm                   string  `json:"arm"`
	MinHeldoutAccuracy    float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy   float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
}

type UPLM1HSymmetricSplitResult struct {
	Schema                     string               `json:"schema"`
	Experiment                 string               `json:"experiment"`
	SourceUPLM1GSeal           string               `json:"source_up_lm1g_seal"`
	StateDimension             int                  `json:"state_dimension"`
	BasePretrainEpochs         int                  `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                  `json:"joint_interleaved_epochs"`
	AdaptationEpochs           int                  `json:"adaptation_epochs"`
	LearningRate               float64              `json:"learning_rate"`
	HalfStepLearningRate       float64              `json:"half_step_learning_rate"`
	RecurrentParametersTrained bool                 `json:"recurrent_parameters_trained"`
	RouterUsed                 bool                 `json:"router_used"`
	ExactRecallUsed            bool                 `json:"exact_recall_used"`
	AttentionUsed              bool                 `json:"attention_used"`
	Metrics                    []UPLM1HMetric       `json:"metrics"`
	Summaries                  []UPLM1HArmSummary   `json:"summaries"`
}

func uplm1hTrain(model *uplm0aModel,baseTrain,paraTrain,thirdTrain []uplm0fExample,arm string) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range baseTrain {
			switch arm {
			case "cyclic_by_index":
				switch i%3 {
				case 0:
					model.trainSentence(thirdTrain[i].text,0.08)
					model.trainSentence(baseTrain[i].text,0.08)
					model.trainSentence(paraTrain[i].text,0.08)
				case 1:
					model.trainSentence(baseTrain[i].text,0.08)
					model.trainSentence(paraTrain[i].text,0.08)
					model.trainSentence(thirdTrain[i].text,0.08)
				case 2:
					model.trainSentence(paraTrain[i].text,0.08)
					model.trainSentence(thirdTrain[i].text,0.08)
					model.trainSentence(baseTrain[i].text,0.08)
				}
			case "strang_split_base":
				model.trainSentence(baseTrain[i].text,0.04)
				model.trainSentence(paraTrain[i].text,0.04)
				model.trainSentence(thirdTrain[i].text,0.08)
				model.trainSentence(paraTrain[i].text,0.04)
				model.trainSentence(baseTrain[i].text,0.04)
			case "strang_split_para":
				model.trainSentence(paraTrain[i].text,0.04)
				model.trainSentence(thirdTrain[i].text,0.04)
				model.trainSentence(baseTrain[i].text,0.08)
				model.trainSentence(thirdTrain[i].text,0.04)
				model.trainSentence(paraTrain[i].text,0.04)
			case "strang_split_third":
				model.trainSentence(thirdTrain[i].text,0.04)
				model.trainSentence(baseTrain[i].text,0.04)
				model.trainSentence(paraTrain[i].text,0.08)
				model.trainSentence(baseTrain[i].text,0.04)
				model.trainSentence(thirdTrain[i].text,0.04)
			}
		}
	}
}

func uplm1hSummary(arm string,metrics []UPLM1HMetric) UPLM1HArmSummary {
	vals:=[]float64{}
	for _,m:=range metrics {
		if m.Arm==arm && (m.Split=="base_heldout"||m.Split=="paraphrase_heldout"||m.Split=="third_heldout") {
			vals=append(vals,m.Top1Accuracy)
		}
	}
	min,max,sum:=1.0,0.0,0.0
	for _,v:=range vals {
		if v<min { min=v }
		if v>max { max=v }
		sum+=v
	}
	mean:=0.0
	if len(vals)>0 { mean=sum/float64(len(vals)) }
	return UPLM1HArmSummary{Arm:arm,MinHeldoutAccuracy:min,MeanHeldoutAccuracy:mean,HeldoutAccuracySpread:max-min}
}

func RunUPLM1H()(UPLM1HSymmetricSplitResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) } }
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,_=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)

	result:=UPLM1HSymmetricSplitResult{
		Schema:UPLM1HSymmetricSplitSchema,Experiment:"UP-LM1H-symmetric-operator-split",
		SourceUPLM1GSeal:"4d7ad1c003f23c67749c229ac766c2a7353ac41c",
		StateDimension:64,BasePretrainEpochs:20,JointInterleavedEpochs:20,AdaptationEpochs:4,
		LearningRate:0.08,HalfStepLearningRate:0.04,
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,AttentionUsed:false,
	}
	arms:=[]string{"cyclic_by_index","strang_split_base","strang_split_para","strang_split_third"}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(start)
		uplm1hTrain(model,baseTrain,paraTrain,thirdTrain,arm)
		for _,x:=range []struct{name string;train,held []uplm0fExample}{
			{"base",baseTrain,baseHeld},
			{"paraphrase",paraTrain,paraHeld},
			{"third",thirdTrain,thirdHeld},
		}{
			mt:=uplm0oByteEval(model,x.train,0,x.name+"_train")
			mh:=uplm0oByteEval(model,x.held,0,x.name+"_heldout")
			result.Metrics=append(result.Metrics,
				UPLM1HMetric{Arm:arm,Split:x.name+"_train",Top1Accuracy:mt.Top1Accuracy,Perplexity:mt.Perplexity},
				UPLM1HMetric{Arm:arm,Split:x.name+"_heldout",Top1Accuracy:mh.Top1Accuracy,Perplexity:mh.Perplexity},
			)
		}
	}
	for _,arm:=range arms { result.Summaries=append(result.Summaries,uplm1hSummary(arm,result.Metrics)) }
	return result,nil
}
