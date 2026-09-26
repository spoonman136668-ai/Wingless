package unitary

import "math"

const UPLM0OParaphraseAdaptationSchema = "wingless.up-lm0o-paraphrase-byte-adaptation.v1"

type UPLM0OByteMetric struct {
	AdaptationEpochs int     `json:"adaptation_epochs"`
	Split            string  `json:"split"`
	Top1Accuracy     float64 `json:"top1_accuracy"`
	Perplexity       float64 `json:"perplexity"`
	Tokens           int     `json:"tokens"`
}

type UPLM0ORoutingMetric struct {
	AdaptationEpochs int          `json:"adaptation_epochs"`
	Metric           UPLM0JMetric `json:"metric"`
}

type UPLM0OParaphraseAdaptationResult struct {
	Schema                          string                   `json:"schema"`
	Experiment                      string                   `json:"experiment"`
	SourceUPLM0NSeal                string                   `json:"source_up_lm0n_seal"`
	StateDimension                  int                      `json:"state_dimension"`
	ExactRecallCap                  int                      `json:"exact_recall_cap"`
	BaseEpochs                      int                      `json:"base_epochs"`
	LearningRate                    float64                  `json:"learning_rate"`
	RouterRetrainingUsed            bool                     `json:"router_retraining_used"`
	AttentionUsed                   bool                     `json:"attention_used"`
	FutureOracleUsed                bool                     `json:"future_oracle_used"`
	BaseReplayPassesPerAdaptEpoch   int                      `json:"base_replay_passes_per_adapt_epoch"`
	ByteMetrics                     []UPLM0OByteMetric       `json:"byte_metrics"`
	RoutingMetrics                  []UPLM0ORoutingMetric    `json:"routing_metrics"`
}

func uplm0oParaphraseCorpus()(train,held []uplm0fExample){
	names:=uplm0gNames()
	values:=uplm0gValues()
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
				initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
				latest:=initial
				s:=""
				for i:=0;i<4;i++ { s+=ns[i]+" saves "+initial[i]+". " }
				for i:=0;i<4;i++ {
					obs:=values[(v+i+3)%6]
					s+=ns[(i+p)%4]+" sees "+obs+". "
				}
				uc:=uplm0eUpdateCount(p)
				for i:=0;i<uc;i++ {
					latest[i]=values[(v+i+2)%6]
					s+=ns[i]+" saves "+latest[i]+". "
				}
				var targets [4]int
				for qi:=0;qi<4;qi++ {
					idx:=(qi+p)%4
					s+=ns[idx]+" recalls "
					targets[qi]=len(s)
					s+=latest[idx]+"."
					if qi<3 { s+=" " }
				}
				s+="\n"
				ex:=uplm0fExample{text:s,targetPos:targets,updateCount:uc}
				if (n+2*v+p)%3!=2 { train=append(train,ex) } else { held=append(held,ex) }
			}
		}
	}
	return
}

func uplm0oCloneModel(src *uplm0aModel)*uplm0aModel {
	m:=&uplm0aModel{
		alphabet:append([]byte(nil),src.alphabet...),
		index:src.index,
		b:append([]float64(nil),src.b...),
	}
	m.w=make([][]float64,len(src.w))
	for i:=range src.w { m.w[i]=append([]float64(nil),src.w[i]...) }
	return m
}

func uplm0oByteEval(model *uplm0aModel,examples []uplm0fExample,adapt int,split string) UPLM0OByteMetric {
	hits,total:=0,0
	nll:=0.0
	for _,ex:=range examples {
		var h [64]float64
		s:=ex.text
		for t:=0;t<len(s)-1;t++ {
			h=uplm0aStep(h,s[t])
			target:=model.index[int(s[t+1])]
			p:=model.probs(h)
			pred:=uplm0aArgmax(p)
			prob:=1e-12
			if target>=0 {
				prob=p[target]
				if pred==target { hits++ }
			}
			if prob<1e-12 { prob=1e-12 }
			nll-=math.Log(prob)
			total++
		}
	}
	return UPLM0OByteMetric{
		AdaptationEpochs:adapt,Split:split,
		Top1Accuracy:float64(hits)/float64(total),
		Perplexity:math.Exp(nll/float64(total)),
		Tokens:total,
	}
}

func RunUPLM0O()(UPLM0OParaphraseAdaptationResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	baseModel:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { baseModel.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0OParaphraseAdaptationResult{
		Schema:UPLM0OParaphraseAdaptationSchema,
		Experiment:"UP-LM0O-paraphrase-byte-adaptation",
		SourceUPLM0NSeal:"7ca4aefc2d11c6e2bc94267d8d11b90e3770c8e0",
		StateDimension:64,ExactRecallCap:16,BaseEpochs:20,LearningRate:0.08,
		RouterRetrainingUsed:false,AttentionUsed:false,FutureOracleUsed:false,
		BaseReplayPassesPerAdaptEpoch:1,
	}

	for _,adapt:=range []int{0,1,2,4} {
		model:=uplm0oCloneModel(baseModel)
		for epoch:=0;epoch<adapt;epoch++ {
			for _,ex:=range paraTrain { model.trainSentence(ex.text,0.08) }
			for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
		}
		result.ByteMetrics=append(result.ByteMetrics,
			uplm0oByteEval(model,baseHeld,adapt,"base_heldout"),
			uplm0oByteEval(model,paraHeld,adapt,"paraphrase_block_heldout"),
		)
		result.RoutingMetrics=append(result.RoutingMetrics,
			UPLM0ORoutingMetric{AdaptationEpochs:adapt,Metric:uplm0nEvaluate(model,classifier,paraHeld,1,true,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM0ORoutingMetric{AdaptationEpochs:adapt,Metric:uplm0nEvaluate(model,classifier,uplm0nHeldout(order),1,true,order+"_stream1")},
			)
		}
	}
	return result,nil
}
