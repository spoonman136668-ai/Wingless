package unitary

const UPLM1IOperatorSplitReintegrationSchema = "wingless.up-lm1i-operator-split-reintegration.v1"

type UPLM1IByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1IRoutingMetric struct {
	Arm    string          `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1IOperatorSplitReintegrationResult struct {
	Schema                     string                 `json:"schema"`
	Experiment                 string                 `json:"experiment"`
	SourceUPLM1HSeal           string                 `json:"source_up_lm1h_seal"`
	StateDimension             int                    `json:"state_dimension"`
	ExactRecallCap             int                    `json:"exact_recall_cap"`
	BasePretrainEpochs         int                    `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                    `json:"joint_interleaved_epochs"`
	AdaptationEpochs           int                    `json:"adaptation_epochs"`
	LearningRate               float64                `json:"learning_rate"`
	HalfStepLearningRate       float64                `json:"half_step_learning_rate"`
	RouterRetrainingUsed       bool                   `json:"router_retraining_used"`
	RecurrentParametersTrained bool                   `json:"recurrent_parameters_trained"`
	AttentionUsed              bool                   `json:"attention_used"`
	FutureOracleUsed           bool                   `json:"future_oracle_used"`
	ByteMetrics                []UPLM1IByteMetric     `json:"byte_metrics"`
	RoutingMetrics             []UPLM1IRoutingMetric  `json:"routing_metrics"`
}

func RunUPLM1I()(UPLM1IOperatorSplitReintegrationResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,_=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)

	classifier:=uplm1cTrainClassifier()

	result:=UPLM1IOperatorSplitReintegrationResult{
		Schema:UPLM1IOperatorSplitReintegrationSchema,
		Experiment:"UP-LM1I-operator-split-reintegration",
		SourceUPLM1HSeal:"5fd53050638d4d636132392647b56933778d82cf",
		StateDimension:64,ExactRecallCap:16,
		BasePretrainEpochs:20,JointInterleavedEpochs:20,AdaptationEpochs:4,
		LearningRate:0.08,HalfStepLearningRate:0.04,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
	}

	for _,arm:=range []string{"cyclic_by_index","strang_split_base"} {
		model:=uplm0oCloneModel(start)
		uplm1hTrain(model,baseTrain,paraTrain,thirdTrain,arm)

		for _,x:=range []struct{name string;data []uplm0fExample}{
			{"base_heldout",baseHeld},
			{"paraphrase_heldout",paraHeld},
			{"third_heldout",thirdHeld},
		}{
			m:=uplm0oByteEval(model,x.data,0,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1IByteMetric{
				Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,
			})
		}

		for _,stream:=range []int{1,4} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM1IRoutingMetric{Arm:arm,Metric:uplm1cEval(model,classifier,baseHeld,stream,"base_block_stream"+itoa(stream))},
				UPLM1IRoutingMetric{Arm:arm,Metric:uplm1cEval(model,classifier,paraHeld,stream,"paraphrase_block_stream"+itoa(stream))},
			)
		}
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1cHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,
					UPLM1IRoutingMetric{Arm:arm,Metric:uplm1cEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream))},
				)
			}
		}
	}
	return result,nil
}
