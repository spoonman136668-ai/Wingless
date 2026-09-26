package unitary

const UPLM0TReadoutAnchorSchema = "wingless.up-lm0t-readout-anchor.v1"

type UPLM0TByteMetric struct {
	Lambda       float64 `json:"lambda"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0TRoutingMetric struct {
	Lambda float64       `json:"lambda"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0TReadoutAnchorResult struct {
	Schema               string                  `json:"schema"`
	Experiment           string                  `json:"experiment"`
	SourceUPLM0SSeal     string                  `json:"source_up_lm0s_seal"`
	StateDimension       int                     `json:"state_dimension"`
	ExactRecallCap       int                     `json:"exact_recall_cap"`
	BaseEpochs           int                     `json:"base_epochs"`
	AdaptationEpochs     int                     `json:"adaptation_epochs"`
	LearningRate         float64                 `json:"learning_rate"`
	BaseReplayUsed       bool                    `json:"base_replay_used"`
	RouterRetrainingUsed bool                    `json:"router_retraining_used"`
	AttentionUsed        bool                    `json:"attention_used"`
	FutureOracleUsed     bool                    `json:"future_oracle_used"`
	ByteMetrics          []UPLM0TByteMetric      `json:"byte_metrics"`
	RoutingMetrics       []UPLM0TRoutingMetric   `json:"routing_metrics"`
}

func uplm0tTrainAnchored(model,base *uplm0aModel,s string,lr,lambda float64) {
	var h [64]float64
	for t:=0;t<len(s)-1;t++ {
		h=uplm0aStep(h,s[t])
		target:=model.index[int(s[t+1])]
		p:=model.probs(h)
		for c:=range p {
			g:=p[c]
			if c==target { g-=1 }
			for i:=0;i<64;i++ {
				dataGrad:=g*h[i]
				anchorGrad:=lambda*(model.w[c][i]-base.w[c][i])
				model.w[c][i]-=lr*(dataGrad+anchorGrad)
			}
			biasGrad:=g+lambda*(model.b[c]-base.b[c])
			model.b[c]-=lr*biasGrad
		}
	}
}

func RunUPLM0T()(UPLM0TReadoutAnchorResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0TReadoutAnchorResult{
		Schema:UPLM0TReadoutAnchorSchema,Experiment:"UP-LM0T-readout-anchor",
		SourceUPLM0SSeal:"21c33f040b55e43fdd956c9c229d6afd4e08c419",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		BaseReplayUsed:false,RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	for _,lambda:=range []float64{0,0.0001,0.001,0.01} {
		model:=uplm0oCloneModel(baseModel)
		for epoch:=0;epoch<4;epoch++ {
			for _,ex:=range paraTrain { uplm0tTrainAnchored(model,baseModel,ex.text,0.08,lambda) }
		}
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0TByteMetric{Lambda:lambda,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0TByteMetric{Lambda:lambda,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0TRoutingMetric{Lambda:lambda,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0TRoutingMetric{Lambda:lambda,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
