package unitary

const UPLM1FSymmetricGradientSchema = "wingless.up-lm1f-symmetric-triplet-gradient.v1"

type UPLM1FByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1FRoutingMetric struct {
	Arm    string          `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1FSymmetricGradientResult struct {
	Schema                     string                  `json:"schema"`
	Experiment                 string                  `json:"experiment"`
	SourceUPLM1ESeal           string                  `json:"source_up_lm1e_seal"`
	StateDimension             int                     `json:"state_dimension"`
	ExactRecallCap             int                     `json:"exact_recall_cap"`
	BasePretrainEpochs         int                     `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                     `json:"joint_interleaved_epochs"`
	AdaptationEpochs           int                     `json:"adaptation_epochs"`
	LearningRate               float64                 `json:"learning_rate"`
	RouterRetrainingUsed       bool                    `json:"router_retraining_used"`
	RecurrentParametersTrained bool                    `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                    `json:"attention_used"`
	FutureOracleUsed           bool                    `json:"future_oracle_used"`
	ByteMetrics                []UPLM1FByteMetric      `json:"byte_metrics"`
	RoutingMetrics             []UPLM1FRoutingMetric   `json:"routing_metrics"`
}

func uplm1fApplyGradient(model *uplm0aModel,grads []uplm0xGradient,scale float64) {
	for c:=range model.w {
		for j:=0;j<64;j++ {
			sum:=0.0
			for _,g:=range grads { sum+=g.w[c][j] }
			model.w[c][j]-=scale*sum
		}
		sumB:=0.0
		for _,g:=range grads { sumB+=g.b[c] }
		model.b[c]-=scale*sumB
	}
}

func uplm1fTrain(model *uplm0aModel,baseTrain,paraTrain,thirdTrain []uplm0fExample,arm string) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range baseTrain {
			switch arm {
			case "cyclic_by_index":
				switch i%3 {
				case 0:
					model.trainSentence(thirdTrain[i].text,0.08);model.trainSentence(baseTrain[i].text,0.08);model.trainSentence(paraTrain[i].text,0.08)
				case 1:
					model.trainSentence(baseTrain[i].text,0.08);model.trainSentence(paraTrain[i].text,0.08);model.trainSentence(thirdTrain[i].text,0.08)
				case 2:
					model.trainSentence(paraTrain[i].text,0.08);model.trainSentence(thirdTrain[i].text,0.08);model.trainSentence(baseTrain[i].text,0.08)
				}
			case "symmetric_triplet_average":
				bg:=uplm0xSentenceGradient(model,baseTrain[i].text)
				pg:=uplm0xSentenceGradient(model,paraTrain[i].text)
				tg:=uplm0xSentenceGradient(model,thirdTrain[i].text)
				uplm1fApplyGradient(model,[]uplm0xGradient{bg,pg,tg},0.08/3.0)
			case "symmetric_triplet_sum":
				bg:=uplm0xSentenceGradient(model,baseTrain[i].text)
				pg:=uplm0xSentenceGradient(model,paraTrain[i].text)
				tg:=uplm0xSentenceGradient(model,thirdTrain[i].text)
				uplm1fApplyGradient(model,[]uplm0xGradient{bg,pg,tg},0.08/3.0)
			}
		}
	}
}

func RunUPLM1F()(UPLM1FSymmetricGradientResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) } }
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,_=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)
	classifier:=uplm1cTrainClassifier()

	result:=UPLM1FSymmetricGradientResult{
		Schema:UPLM1FSymmetricGradientSchema,Experiment:"UP-LM1F-symmetric-triplet-gradient",
		SourceUPLM1ESeal:"1968a301cad7ad318713143037ede9113fa43a0a",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
	}

	for _,arm:=range []string{"cyclic_by_index","symmetric_triplet_average","symmetric_triplet_sum"} {
		model:=uplm0oCloneModel(start)
		uplm1fTrain(model,baseTrain,paraTrain,thirdTrain,arm)

		for _,x:=range []struct{name string;data []uplm0fExample}{
			{"base_block_stream1",baseHeld},
			{"paraphrase_block_stream1",paraHeld},
			{"third_block_stream1",thirdHeld},
		}{
			m:=uplm1cEval(model,classifier,x.data,1,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1FByteMetric{Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1cHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,
					UPLM1FRoutingMetric{Arm:arm,Metric:uplm1cEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream))},
				)
			}
		}
	}
	return result,nil
}
