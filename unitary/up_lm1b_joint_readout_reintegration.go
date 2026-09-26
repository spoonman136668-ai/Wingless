package unitary

const UPLM1BJointReadoutReintegrationSchema = "wingless.up-lm1b-joint-readout-reintegration.v1"

type UPLM1BByteMetric struct {
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM1BJointReadoutReintegrationResult struct {
	Schema                          string             `json:"schema"`
	Experiment                      string             `json:"experiment"`
	SourceUPLM1ASeal                string             `json:"source_up_lm1a_seal"`
	StateDimension                  int                `json:"state_dimension"`
	ExactRecallCap                  int                `json:"exact_recall_cap"`
	BasePretrainEpochs              int                `json:"base_pretrain_epochs"`
	JointInterleavedEpochs          int                `json:"joint_interleaved_epochs"`
	LearningRate                    float64            `json:"learning_rate"`
	RecurrentParametersTrained      bool               `json:"recurrent_parameters_trained"`
	RouterRetrainingUsed            bool               `json:"router_retraining_used"`
	AttentionUsed                   bool               `json:"attention_used"`
	FutureOracleUsed                bool               `json:"future_oracle_used"`
	ByteMetrics                     []UPLM1BByteMetric `json:"byte_metrics"`
	IntegratedMetrics               []UPLM0JMetric     `json:"integrated_metrics"`
}

func uplm1bByteMetric(split string,m UPLM0OByteMetric) UPLM1BByteMetric {
	return UPLM1BByteMetric{Split:split,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,Tokens:m.Tokens}
}

func RunUPLM1B()(UPLM1BJointReadoutReintegrationResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()

	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(model,baseTrain,paraTrain,20)
	classifier:=uplm0nTrainClassifier()

	result:=UPLM1BJointReadoutReintegrationResult{
		Schema:UPLM1BJointReadoutReintegrationSchema,
		Experiment:"UP-LM1B-joint-readout-reintegration",
		SourceUPLM1ASeal:"7764d7d75d1d4e9d2c200c68a0de26d425647d70",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,LearningRate:0.08,
		RecurrentParametersTrained:false,RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}

	bm:=uplm0oByteEval(model,baseHeld,0,"base_heldout")
	pm:=uplm0oByteEval(model,paraHeld,0,"paraphrase_heldout")
	result.ByteMetrics=append(result.ByteMetrics,
		uplm1bByteMetric("base_heldout",bm),
		uplm1bByteMetric("paraphrase_heldout",pm),
	)

	result.IntegratedMetrics=append(result.IntegratedMetrics,
		uplm0nEvaluate(model,classifier,baseHeld,1,true,"base_block_stream1"),
		uplm0nEvaluate(model,classifier,baseHeld,4,true,"base_block_stream4"),
		uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1"),
		uplm0nEvaluate(model,classifier,paraHeld,4,true,"paraphrase_block_stream4"),
	)
	for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
		held:=uplm0nHeldout(order)
		result.IntegratedMetrics=append(result.IntegratedMetrics,
			uplm0nEvaluate(model,classifier,held,1,true,order+"_stream1"),
			uplm0nEvaluate(model,classifier,held,4,true,order+"_stream4"),
		)
	}
	return result,nil
}
