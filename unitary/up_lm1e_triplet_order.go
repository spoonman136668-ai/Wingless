package unitary

const UPLM1ETripletOrderSchema = "wingless.up-lm1e-triplet-order.v1"

type UPLM1EByteMetric struct {
	Schedule     string  `json:"schedule"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1ERoutingMetric struct {
	Schedule string          `json:"schedule"`
	Metric   UPLM0JMetric `json:"metric"`
}

type UPLM1ETripletOrderResult struct {
	Schema                     string                 `json:"schema"`
	Experiment                 string                 `json:"experiment"`
	SourceUPLM1DSeal           string                 `json:"source_up_lm1d_seal"`
	StateDimension             int                    `json:"state_dimension"`
	ExactRecallCap             int                    `json:"exact_recall_cap"`
	BasePretrainEpochs         int                    `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                    `json:"joint_interleaved_epochs"`
	AdaptationEpochs           int                    `json:"adaptation_epochs"`
	LearningRate               float64                `json:"learning_rate"`
	RouterRetrainingUsed       bool                   `json:"router_retraining_used"`
	RecurrentParametersTrained bool                   `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                   `json:"attention_used"`
	FutureOracleUsed           bool                   `json:"future_oracle_used"`
	NewOutputBytes             []int                  `json:"new_output_bytes"`
	ByteMetrics                []UPLM1EByteMetric     `json:"byte_metrics"`
	RoutingMetrics             []UPLM1ERoutingMetric  `json:"routing_metrics"`
}

func uplm1eTrain(model *uplm0aModel,baseTrain,paraTrain,thirdTrain []uplm0fExample,schedule string) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range baseTrain {
			train:=func(which int){
				switch which {
				case 0: model.trainSentence(thirdTrain[i].text,0.08)
				case 1: model.trainSentence(baseTrain[i].text,0.08)
				case 2: model.trainSentence(paraTrain[i].text,0.08)
				}
			}
			switch schedule {
			case "third_base_para":
				train(0);train(1);train(2)
			case "base_para_third":
				train(1);train(2);train(0)
			case "cyclic_by_index":
				switch i%3 {
				case 0: train(0);train(1);train(2)
				case 1: train(1);train(2);train(0)
				case 2: train(2);train(0);train(1)
				}
			case "cyclic_by_epoch":
				switch epoch%3 {
				case 0: train(0);train(1);train(2)
				case 1: train(1);train(2);train(0)
				case 2: train(2);train(0);train(1)
				}
			}
		}
	}
}

func RunUPLM1E()(UPLM1ETripletOrderResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) } }
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,added:=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)
	classifier:=uplm1cTrainClassifier()

	result:=UPLM1ETripletOrderResult{
		Schema:UPLM1ETripletOrderSchema,Experiment:"UP-LM1E-triplet-order",
		SourceUPLM1DSeal:"76c653cc9a20208d42041804b54400e946c3a52e",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
		NewOutputBytes:append([]int(nil),added...),
	}
	for _,schedule:=range []string{"third_base_para","base_para_third","cyclic_by_index","cyclic_by_epoch"} {
		model:=uplm0oCloneModel(start)
		uplm1eTrain(model,baseTrain,paraTrain,thirdTrain,schedule)

		for _,x:=range []struct{name string;data []uplm0fExample}{
			{"base_block_stream1",baseHeld},
			{"paraphrase_block_stream1",paraHeld},
			{"third_block_stream1",thirdHeld},
		}{
			m:=uplm1cEval(model,classifier,x.data,1,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1EByteMetric{Schedule:schedule,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity})
		}
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1cHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,
					UPLM1ERoutingMetric{Schedule:schedule,Metric:uplm1cEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream))},
				)
			}
		}
	}
	return result,nil
}
