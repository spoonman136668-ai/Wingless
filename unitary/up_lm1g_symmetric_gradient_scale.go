package unitary

import "math"

const UPLM1GSymmetricScaleSchema = "wingless.up-lm1g-symmetric-gradient-scale.v1"

type UPLM1GByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1GRoutingMetric struct {
	Arm    string          `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1GSymmetricScaleResult struct {
	Schema                     string                `json:"schema"`
	Experiment                 string                `json:"experiment"`
	SourceUPLM1FSeal           string                `json:"source_up_lm1f_seal"`
	StateDimension             int                   `json:"state_dimension"`
	ExactRecallCap             int                   `json:"exact_recall_cap"`
	BasePretrainEpochs         int                   `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                   `json:"joint_interleaved_epochs"`
	AdaptationEpochs           int                   `json:"adaptation_epochs"`
	LearningRate               float64               `json:"learning_rate"`
	RouterRetrainingUsed       bool                  `json:"router_retraining_used"`
	RecurrentParametersTrained bool                  `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                  `json:"attention_used"`
	FutureOracleUsed           bool                  `json:"future_oracle_used"`
	ByteMetrics                []UPLM1GByteMetric    `json:"byte_metrics"`
	RoutingMetrics             []UPLM1GRoutingMetric `json:"routing_metrics"`
}

func uplm1gTrain(model *uplm0aModel,baseTrain,paraTrain,thirdTrain []uplm0fExample,arm string) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range baseTrain {
			if arm=="cyclic_by_index" {
				switch i%3 {
				case 0:
					model.trainSentence(thirdTrain[i].text,0.08);model.trainSentence(baseTrain[i].text,0.08);model.trainSentence(paraTrain[i].text,0.08)
				case 1:
					model.trainSentence(baseTrain[i].text,0.08);model.trainSentence(paraTrain[i].text,0.08);model.trainSentence(thirdTrain[i].text,0.08)
				case 2:
					model.trainSentence(paraTrain[i].text,0.08);model.trainSentence(thirdTrain[i].text,0.08);model.trainSentence(baseTrain[i].text,0.08)
				}
				continue
			}
			bg:=uplm0xSentenceGradient(model,baseTrain[i].text)
			pg:=uplm0xSentenceGradient(model,paraTrain[i].text)
			tg:=uplm0xSentenceGradient(model,thirdTrain[i].text)
			scale:=0.08/3.0
			switch arm {
			case "symmetric_scale_1oversqrt3":
				scale=0.08/math.Sqrt(3)
			case "symmetric_scale_1":
				scale=0.08
			}
			uplm1fApplyGradient(model,[]uplm0xGradient{bg,pg,tg},scale)
		}
	}
}

func RunUPLM1G()(UPLM1GSymmetricScaleResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) } }
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,_=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)
	classifier:=uplm1cTrainClassifier()

	result:=UPLM1GSymmetricScaleResult{
		Schema:UPLM1GSymmetricScaleSchema,Experiment:"UP-LM1G-symmetric-gradient-scale",
		SourceUPLM1FSeal:"a24fcfb761d68e6ce191bdb21d994efde68fddb9",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	for _,arm:=range []string{"cyclic_by_index","symmetric_scale_1over3","symmetric_scale_1oversqrt3","symmetric_scale_1"} {
		model:=uplm0oCloneModel(start)
		uplm1gTrain(model,baseTrain,paraTrain,thirdTrain,arm)
		for _,x:=range []struct{name string;data []uplm0fExample}{
			{"base_block_stream1",baseHeld},
			{"paraphrase_block_stream1",paraHeld},
			{"third_block_stream1",thirdHeld},
		}{
			m:=uplm1cEval(model,classifier,x.data,1,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1GByteMetric{Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1cHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,
					UPLM1GRoutingMetric{Arm:arm,Metric:uplm1cEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream))},
				)
			}
		}
	}
	return result,nil
}
