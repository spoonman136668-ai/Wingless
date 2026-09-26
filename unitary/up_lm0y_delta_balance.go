package unitary

import "math"

const UPLM0YDeltaBalanceSchema = "wingless.up-lm0y-delta-balance.v1"

type UPLM0YByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM0YRoutingMetric struct {
	Arm    string       `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0YArmDiagnostic struct {
	Arm                       string  `json:"arm"`
	MeanRawParaToBaseNormRatio float64 `json:"mean_raw_para_to_base_norm_ratio"`
	PairRowProjections         int     `json:"pair_row_projections"`
	PairsProcessed             int     `json:"pairs_processed"`
}

type UPLM0YDeltaBalanceResult struct {
	Schema               string                  `json:"schema"`
	Experiment           string                  `json:"experiment"`
	SourceUPLM0XSeal     string                  `json:"source_up_lm0x_seal"`
	StateDimension       int                     `json:"state_dimension"`
	ExactRecallCap       int                     `json:"exact_recall_cap"`
	BaseEpochs           int                     `json:"base_epochs"`
	AdaptationEpochs     int                     `json:"adaptation_epochs"`
	LearningRate         float64                 `json:"learning_rate"`
	RouterRetrainingUsed bool                    `json:"router_retraining_used"`
	AttentionUsed        bool                    `json:"attention_used"`
	FutureOracleUsed     bool                    `json:"future_oracle_used"`
	AnchoringUsed        bool                    `json:"anchoring_used"`
	ByteMetrics          []UPLM0YByteMetric      `json:"byte_metrics"`
	RoutingMetrics       []UPLM0YRoutingMetric   `json:"routing_metrics"`
	ArmDiagnostics       []UPLM0YArmDiagnostic   `json:"arm_diagnostics"`
}

func uplm0yPairUpdate(model *uplm0aModel,paraText,baseText,arm string)(float64,int){
	pre:=uplm0oCloneModel(model)
	pm:=uplm0oCloneModel(pre)
	bm:=uplm0oCloneModel(pre)
	pm.trainSentence(paraText,0.08)
	bm.trainSentence(baseText,0.08)

	baseNorm2,paraNorm2:=0.0,0.0
	for c:=range model.w {
		for i:=range model.w[c] {
			bd:=bm.w[c][i]-pre.w[c][i]
			pd:=pm.w[c][i]-pre.w[c][i]
			baseNorm2+=bd*bd
			paraNorm2+=pd*pd
		}
		bd:=bm.b[c]-pre.b[c]
		pd:=pm.b[c]-pre.b[c]
		baseNorm2+=bd*bd
		paraNorm2+=pd*pd
	}

	ratio:=0.0
	scale:=1.0
	if baseNorm2>0 && paraNorm2>0 {
		ratio=math.Sqrt(paraNorm2/baseNorm2)
		if arm!="paired_delta_sum" {
			scale=math.Sqrt(baseNorm2/paraNorm2)
		}
	}

	projections:=0
	for c:=range model.w {
		rowDot,rowBase2:=0.0,0.0
		if arm=="norm_balanced_row_projected" {
			for i:=range model.w[c] {
				bd:=bm.w[c][i]-pre.w[c][i]
				pd:=scale*(pm.w[c][i]-pre.w[c][i])
				rowDot+=bd*pd
				rowBase2+=bd*bd
			}
			bdBias:=bm.b[c]-pre.b[c]
			pdBias:=scale*(pm.b[c]-pre.b[c])
			rowDot+=bdBias*pdBias
			rowBase2+=bdBias*bdBias
			if rowDot<0 && rowBase2>0 { projections++ }
		}

		coeff:=0.0
		if arm=="norm_balanced_row_projected" && rowDot<0 && rowBase2>0 {
			coeff=rowDot/rowBase2
		}
		for i:=range model.w[c] {
			bd:=bm.w[c][i]-pre.w[c][i]
			pd:=scale*(pm.w[c][i]-pre.w[c][i])
			if coeff!=0 { pd-=coeff*bd }
			model.w[c][i]=pre.w[c][i]+bd+pd
		}
		bdBias:=bm.b[c]-pre.b[c]
		pdBias:=scale*(pm.b[c]-pre.b[c])
		if coeff!=0 { pdBias-=coeff*bdBias }
		model.b[c]=pre.b[c]+bdBias+pdBias
	}
	return ratio,projections
}

func uplm0yTrain(model *uplm0aModel,baseTrain,paraTrain []uplm0fExample,arm string)(float64,int,int){
	ratioSum:=0.0
	projections:=0
	pairs:=0
	for epoch:=0;epoch<4;epoch++ {
		for i:=range baseTrain {
			ratio,p:=uplm0yPairUpdate(model,paraTrain[i].text,baseTrain[i].text,arm)
			ratioSum+=ratio
			projections+=p
			pairs++
		}
	}
	mean:=0.0
	if pairs>0 { mean=ratioSum/float64(pairs) }
	return mean,projections,pairs
}

func RunUPLM0Y()(UPLM0YDeltaBalanceResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0YDeltaBalanceResult{
		Schema:UPLM0YDeltaBalanceSchema,
		Experiment:"UP-LM0Y-delta-balance",
		SourceUPLM0XSeal:"8b4bbbd93538f1f001f73f7e888e901769868a7e",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,AnchoringUsed:false,
	}
	for _,arm:=range []string{"paired_delta_sum","norm_balanced_delta","norm_balanced_row_projected"} {
		model:=uplm0oCloneModel(baseModel)
		meanRatio,projections,pairs:=uplm0yTrain(model,baseTrain,paraTrain,arm)
		result.ArmDiagnostics=append(result.ArmDiagnostics,UPLM0YArmDiagnostic{
			Arm:arm,MeanRawParaToBaseNormRatio:meanRatio,PairRowProjections:projections,PairsProcessed:pairs,
		})

		bm:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		pm:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0YByteMetric{Arm:arm,Split:bm.Split,Top1Accuracy:bm.Top1Accuracy,Perplexity:bm.Perplexity,Tokens:bm.Tokens},
			UPLM0YByteMetric{Arm:arm,Split:pm.Split,Top1Accuracy:pm.Top1Accuracy,Perplexity:pm.Perplexity,Tokens:pm.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0YRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0YRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
