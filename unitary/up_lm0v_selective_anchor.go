package unitary

const UPLM0VSelectiveAnchorSchema = "wingless.up-lm0v-selective-anchor.v1"

type UPLM0VByteMetric struct {
	Arm          string  `json:"arm"`
	Lambda       float64 `json:"lambda"`
	AnchorMode   string  `json:"anchor_mode"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0VRoutingMetric struct {
	Arm    string       `json:"arm"`
	Lambda float64      `json:"lambda"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0VSelectiveAnchorResult struct {
	Schema               string                 `json:"schema"`
	Experiment           string                 `json:"experiment"`
	SourceUPLM0USeal     string                 `json:"source_up_lm0u_seal"`
	StateDimension       int                    `json:"state_dimension"`
	ExactRecallCap       int                    `json:"exact_recall_cap"`
	BaseEpochs           int                    `json:"base_epochs"`
	AdaptationEpochs     int                    `json:"adaptation_epochs"`
	LearningRate         float64                `json:"learning_rate"`
	AnchorLambda         float64                `json:"anchor_lambda"`
	SelectiveAnchorBytes []int                  `json:"selective_anchor_bytes"`
	RouterRetrainingUsed bool                   `json:"router_retraining_used"`
	AttentionUsed        bool                   `json:"attention_used"`
	FutureOracleUsed     bool                   `json:"future_oracle_used"`
	ByteMetrics          []UPLM0VByteMetric     `json:"byte_metrics"`
	RoutingMetrics       []UPLM0VRoutingMetric  `json:"routing_metrics"`
}

func uplm0vMask(model *uplm0aModel,mode string) []bool {
	mask:=make([]bool,len(model.alphabet))
	switch mode {
	case "full":
		for i:=range mask { mask[i]=true }
	case "selective":
		for _,b:=range []byte{'t','o','b','p'} {
			idx:=model.index[int(b)]
			if idx>=0 { mask[idx]=true }
		}
	}
	return mask
}

func uplm0vTrainSentence(model,base *uplm0aModel,s string,lr,lambda float64,mask []bool) {
	var h [64]float64
	for t:=0;t<len(s)-1;t++ {
		h=uplm0aStep(h,s[t])
		target:=model.index[int(s[t+1])]
		p:=model.probs(h)
		for c:=range p {
			g:=p[c]
			if c==target { g-=1 }
			anchored:=lambda>0 && c<len(mask) && mask[c]
			for i:=0;i<64;i++ {
				grad:=g*h[i]
				if anchored { grad+=lambda*(model.w[c][i]-base.w[c][i]) }
				model.w[c][i]-=lr*grad
			}
			biasGrad:=g
			if anchored { biasGrad+=lambda*(model.b[c]-base.b[c]) }
			model.b[c]-=lr*biasGrad
		}
	}
}

func uplm0vTrainInterleaved(model,base *uplm0aModel,baseTrain,paraTrain []uplm0fExample,lambda float64,mode string) {
	mask:=uplm0vMask(model,mode)
	for epoch:=0;epoch<4;epoch++ {
		for i:=range paraTrain {
			uplm0vTrainSentence(model,base,paraTrain[i].text,0.08,lambda,mask)
			uplm0vTrainSentence(model,base,baseTrain[i].text,0.08,lambda,mask)
		}
	}
}

func RunUPLM0V()(UPLM0VSelectiveAnchorResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0VSelectiveAnchorResult{
		Schema:UPLM0VSelectiveAnchorSchema,
		Experiment:"UP-LM0V-selective-anchor",
		SourceUPLM0USeal:"3939e1f232c159bad6b0f90297d60ef0b3201f39",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		AnchorLambda:0.0001,SelectiveAnchorBytes:[]int{98,111,112,116},
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	arms:=[]struct{name,mode string; lambda float64}{
		{"interleaved_anchor0","none",0},
		{"interleaved_full_anchor0001","full",0.0001},
		{"interleaved_selective_anchor0001","selective",0.0001},
	}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(baseModel)
		uplm0vTrainInterleaved(model,baseModel,baseTrain,paraTrain,arm.lambda,arm.mode)
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0VByteMetric{Arm:arm.name,Lambda:arm.lambda,AnchorMode:arm.mode,Split:b.Split,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0VByteMetric{Arm:arm.name,Lambda:arm.lambda,AnchorMode:arm.mode,Split:p.Split,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0VRoutingMetric{Arm:arm.name,Lambda:arm.lambda,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0VRoutingMetric{Arm:arm.name,Lambda:arm.lambda,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
