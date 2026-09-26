package unitary

import "strings"

const UPLM0SLexicalLossMaskSchema = "wingless.up-lm0s-lexical-loss-mask.v1"

type UPLM0SByteMetric struct {
	Arm                    string  `json:"arm"`
	Split                  string  `json:"split"`
	UpdatedTargetsPerEpoch int     `json:"updated_targets_per_epoch"`
	Top1Accuracy           float64 `json:"top1_accuracy"`
	Perplexity             float64 `json:"perplexity"`
	Tokens                 int     `json:"tokens"`
}

type UPLM0SRoutingMetric struct {
	Arm    string          `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM0SLexicalLossMaskResult struct {
	Schema                  string                 `json:"schema"`
	Experiment              string                 `json:"experiment"`
	SourceUPLM0RSeal        string                 `json:"source_up_lm0r_seal"`
	StateDimension          int                    `json:"state_dimension"`
	ExactRecallCap          int                    `json:"exact_recall_cap"`
	BaseEpochs              int                    `json:"base_epochs"`
	AdaptationEpochs        int                    `json:"adaptation_epochs"`
	LearningRate            float64                `json:"learning_rate"`
	BaseReplayUsed          bool                   `json:"base_replay_used"`
	RouterRetrainingUsed    bool                   `json:"router_retraining_used"`
	AttentionUsed           bool                   `json:"attention_used"`
	FutureOracleUsed        bool                   `json:"future_oracle_used"`
	ByteMetrics             []UPLM0SByteMetric     `json:"byte_metrics"`
	RoutingMetrics          []UPLM0SRoutingMetric  `json:"routing_metrics"`
}

func uplm0sVerbTargetMask(s string) []bool {
	mask:=make([]bool,len(s))
	for _,verb:=range []string{"saves","sees","recalls"} {
		needle:=" "+verb+" "
		from:=0
		for from<len(s) {
			k:=strings.Index(s[from:],needle)
			if k<0 { break }
			wordStart:=from+k+1
			for j:=0;j<len(verb);j++ { mask[wordStart+j]=true }
			from=wordStart+len(verb)
		}
	}
	return mask
}

func uplm0sTrainMasked(model *uplm0aModel,s string,lr float64) int {
	mask:=uplm0sVerbTargetMask(s)
	var h [64]float64
	updates:=0
	for t:=0;t<len(s)-1;t++ {
		h=uplm0aStep(h,s[t])
		targetPos:=t+1
		if !mask[targetPos] { continue }
		target:=model.index[int(s[targetPos])]
		p:=model.probs(h)
		for c:=range p {
			g:=p[c]
			if c==target { g-=1 }
			for i:=0;i<64;i++ { model.w[c][i]-=lr*g*h[i] }
			model.b[c]-=lr*g
		}
		updates++
	}
	return updates
}

func uplm0sFullTargets(examples []uplm0fExample) int {
	n:=0
	for _,ex:=range examples { n+=len(ex.text)-1 }
	return n
}

func uplm0sMaskedTargets(examples []uplm0fExample) int {
	n:=0
	for _,ex:=range examples {
		mask:=uplm0sVerbTargetMask(ex.text)
		for i:=1;i<len(mask);i++ { if mask[i]{n++} }
	}
	return n
}

func RunUPLM0S()(UPLM0SLexicalLossMaskResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0SLexicalLossMaskResult{
		Schema:UPLM0SLexicalLossMaskSchema,Experiment:"UP-LM0S-lexical-loss-mask",
		SourceUPLM0RSeal:"d75ed3c7f6d3d0447b475477bab3496295cd5188",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,AdaptationEpochs:4,LearningRate:0.08,
		BaseReplayUsed:false,RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	fullTargets:=uplm0sFullTargets(paraTrain)
	maskedTargets:=uplm0sMaskedTargets(paraTrain)
	for _,arm:=range []string{"full_sentence","verb_targets_only"} {
		model:=uplm0oCloneModel(baseModel)
		for epoch:=0;epoch<4;epoch++ {
			for _,ex:=range paraTrain {
				if arm=="full_sentence" { model.trainSentence(ex.text,0.08) } else { uplm0sTrainMasked(model,ex.text,0.08) }
			}
		}
		updates:=fullTargets
		if arm=="verb_targets_only" { updates=maskedTargets }
		b:=uplm0oByteEval(model,baseHeld,4,"base_heldout")
		p:=uplm0oByteEval(model,paraHeld,4,"paraphrase_block_heldout")
		result.ByteMetrics=append(result.ByteMetrics,
			UPLM0SByteMetric{Arm:arm,Split:b.Split,UpdatedTargetsPerEpoch:updates,Top1Accuracy:b.Top1Accuracy,Perplexity:b.Perplexity,Tokens:b.Tokens},
			UPLM0SByteMetric{Arm:arm,Split:p.Split,UpdatedTargetsPerEpoch:updates,Top1Accuracy:p.Top1Accuracy,Perplexity:p.Perplexity,Tokens:p.Tokens},
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0SRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0SRoutingMetric{Arm:arm,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
