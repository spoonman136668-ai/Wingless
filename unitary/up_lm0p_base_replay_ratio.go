package unitary

const UPLM0PBaseReplaySchema = "wingless.up-lm0p-base-replay-ratio.v1"

type UPLM0PByteMetric struct {
	BaseReplayPasses int     `json:"base_replay_passes"`
	Split            string  `json:"split"`
	Top1Accuracy     float64 `json:"top1_accuracy"`
	Perplexity       float64 `json:"perplexity"`
	Tokens           int     `json:"tokens"`
}

type UPLM0PRoutingMetric struct {
	BaseReplayPasses int          `json:"base_replay_passes"`
	Metric           UPLM0JMetric `json:"metric"`
}

type UPLM0PBaseReplayResult struct {
	Schema                          string                  `json:"schema"`
	Experiment                      string                  `json:"experiment"`
	SourceUPLM0OSeal                string                  `json:"source_up_lm0o_seal"`
	StateDimension                  int                     `json:"state_dimension"`
	ExactRecallCap                  int                     `json:"exact_recall_cap"`
	BaseEpochs                      int                     `json:"base_epochs"`
	AdaptationEpochs                int                     `json:"adaptation_epochs"`
	LearningRate                    float64                 `json:"learning_rate"`
	RouterRetrainingUsed            bool                    `json:"router_retraining_used"`
	AttentionUsed                   bool                    `json:"attention_used"`
	FutureOracleUsed                bool                    `json:"future_oracle_used"`
	ByteMetrics                     []UPLM0PByteMetric      `json:"byte_metrics"`
	RoutingMetrics                  []UPLM0PRoutingMetric   `json:"routing_metrics"`
}

func RunUPLM0P()(UPLM0PBaseReplayResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0PBaseReplayResult{
		Schema:UPLM0PBaseReplaySchema,
		Experiment:"UP-LM0P-base-replay-ratio",
		SourceUPLM0OSeal:"d2c6fd9ed91b607cbb86b4dc3615cbff5342b6c8",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	for _,replay:=range []int{0,1,2,4} {
		model:=uplm0oCloneModel(baseModel)
		for epoch:=0;epoch<4;epoch++ {
			for _,ex:=range paraTrain { model.trainSentence(ex.text,0.08) }
			for pass:=0;pass<replay;pass++ {
				for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
			}
		}
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0PByteMetric{BaseReplayPasses:replay,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0PByteMetric{BaseReplayPasses:replay,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0PRoutingMetric{BaseReplayPasses:replay,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0PRoutingMetric{BaseReplayPasses:replay,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
