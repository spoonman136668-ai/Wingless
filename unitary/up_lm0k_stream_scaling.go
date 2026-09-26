package unitary

import (
	"math"
	"strings"
)

const UPLM0KStreamScalingSchema = "wingless.up-lm0k-stream-scaling.v1"

type UPLM0KStreamScalingResult struct {
	Schema                          string                    `json:"schema"`
	Experiment                      string                    `json:"experiment"`
	SourceUPLM0JSeal                string                    `json:"source_up_lm0j_seal"`
	StateDimension                  int                       `json:"state_dimension"`
	ExactRecallCap                  int                       `json:"exact_recall_cap"`
	ClassifierEpochs                int                       `json:"classifier_epochs"`
	ClassifierLearningRate          float64                   `json:"classifier_learning_rate"`
	ExplicitTypeAtLearnedInference  bool                      `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool                     `json:"value_bytes_at_classifier_inference"`
	AttentionUsed                   bool                      `json:"attention_used"`
	FutureOracleUsed                bool                      `json:"future_oracle_used"`
	ClassifierMetrics               []UPLM0JClassifierMetric  `json:"classifier_metrics"`
	Metrics                         []UPLM0JMetric            `json:"metrics"`
}

func uplm0kEvaluate(model *uplm0aModel, classifier *uplm0jClassifier, examples []uplm0fExample, streamSize int, learned bool, split string) UPLM0JMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	exactParagraphs,paragraphs:=0,0
	nll:=0.0
	maxEntries:=0
	tp,fp,fn:=0,0,0
	eventHits,eventTotal:=0,0
	reportHits,reportTotal:=0,0
	exactBy:=map[int]int{0:0,1:0,2:0,4:0}
	totalBy:=map[int]int{0:0,1:0,2:0,4:0}

	for start:=0;start<len(examples);{
		end:=start+streamSize
		if end>len(examples){end=len(examples)}
		var h [64]float64
		recall:=newUPLM0CRecall()

		for e:=start;e<end;e++ {
			s:=examples[e].text
			queryCorrect:=[4]bool{}
			querySeen:=[4]bool{}
			clause:=""
			routeKnown:=false
			routeClass:=-1
			trueClass:=-1
			queryName:=""
			reportOverridePending:=false

			for t:=0;t<len(s)-1;t++ {
				b:=s[t]
				h=uplm0aStep(h,b)
				clause+=string(b)
				trim:=strings.TrimSpace(clause)

				if b==' '&&!routeKnown {
					fields:=strings.Fields(trim)
					if len(fields)==2 {
						trueClass=uplm0jTrueClass(fields[1])
						if trueClass>=0 {
							routeClass=trueClass
							if learned { routeClass=uplm0jPredict(classifier,fields[0]+" "+fields[1]) }
							routeKnown=true
							eventTotal++
							if routeClass==trueClass { eventHits++ }
							if trueClass==uplm0jReport {
								reportTotal++
								if routeClass==uplm0jReport { reportHits++ }
							}
							if routeClass==uplm0jReport {
								queryName=fields[0]
								reportOverridePending=true
							}
						}
					}
				}

				targetByte:=s[t+1]
				target:=model.index[int(targetByte)]
				p:=model.probs(h)
				pred:=uplm0aArgmax(p)
				prob:=p[target]

				if reportOverridePending&&queryName!="" {
					if value,ok:=recall.values[queryName];ok&&len(value)>0 {
						memByte:=value[0]
						pred=model.index[int(memByte)]
						if memByte==targetByte { prob=1 } else { prob=1e-12 }
					}
					reportOverridePending=false
				}

				if prob<1e-12 { prob=1e-12 }
				total++
				nll-=math.Log(prob)
				if pred==target { hits++ }
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok {
					depTotal++
					querySeen[qi]=true
					if pred==target { depHits++;queryCorrect[qi]=true }
				}

				if b=='.' {
					clean:=strings.TrimSuffix(strings.TrimSpace(clause),".")
					fields:=strings.Fields(clean)
					if len(fields)>=3&&trueClass>=0 {
						name:=fields[0]
						value:=fields[2]
						predStore:=routeKnown&&routeClass==uplm0jStore
						trueStore:=trueClass==uplm0jStore
						if predStore {
							recall.write(name,value)
							if trueStore { tp++ } else { fp++ }
						} else if trueStore { fn++ }
					}
					clause=""
					routeKnown=false
					routeClass=-1
					trueClass=-1
					queryName=""
					reportOverridePending=false
				}

				if len(recall.order)>maxEntries { maxEntries=len(recall.order) }
			}

			all:=true
			for i:=0;i<4;i++ {
				if !querySeen[i]||!queryCorrect[i] { all=false }
			}
			paragraphs++
			totalBy[examples[e].updateCount]++
			if all { exactParagraphs++;exactBy[examples[e].updateCount]++ }
			if streamSize>1 { h=uplm0aStep(h,s[len(s)-1]) }
		}
		start=end
	}

	precision,recallScore:=1.0,1.0
	if tp+fp>0 { precision=float64(tp)/float64(tp+fp) }
	if tp+fn>0 { recallScore=float64(tp)/float64(tp+fn) }
	acc:=func(k int)float64{
		if totalBy[k]==0 { return 0 }
		return float64(exactBy[k])/float64(totalBy[k])
	}
	arm:="explicit_event_routing"
	if learned { arm="learned_prefix_threeway" }
	eventAccuracy:=1.0
	if eventTotal>0 { eventAccuracy=float64(eventHits)/float64(eventTotal) }
	reportAccuracy:=1.0
	if reportTotal>0 { reportAccuracy=float64(reportHits)/float64(reportTotal) }

	return UPLM0JMetric{
		Arm:arm,Split:split,
		Top1Accuracy:float64(hits)/float64(total),
		Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,
		EventRoutingAccuracy:eventAccuracy,ReportRoutingAccuracy:reportAccuracy,
		MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0K()(UPLM0KStreamScalingResult,error){
	train,held,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,0.08) }
	}
	classifier:=uplm0jTrainClassifier()

	result:=UPLM0KStreamScalingResult{
		Schema:UPLM0KStreamScalingSchema,
		Experiment:"UP-LM0K-stream-scaling",
		SourceUPLM0JSeal:"01af2d8f9a0aa00245a567be70413734cb7daaaa",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,
		ExplicitTypeAtLearnedInference:false,ValueBytesAtClassifierInference:false,AttentionUsed:false,FutureOracleUsed:false,
	}
	for _,split:=range []struct{name string;start,end int}{
		{"train_subjects",0,4},{"heldout_subjects",4,6},
	}{
		result.ClassifierMetrics=append(result.ClassifierMetrics,
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,-1,"all"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jStore,"STORE"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jObserve,"OBSERVE"),
			uplm0jClassifierEval(classifier,split.name,split.start,split.end,uplm0jReport,"REPORT"),
		)
	}
	for _,arm:=range []struct{name string;learned bool}{
		{"explicit_event_routing",false},{"learned_prefix_threeway",true},
	}{
		for _,stream:=range []int{1,4,8,16} {
			split:="heldout_stream"+string(rune('0'+stream))
			if stream>=10 { split="heldout_stream16" }
			result.Metrics=append(result.Metrics,uplm0kEvaluate(model,classifier,held,stream,arm.learned,split))
		}
	}
	return result,nil
}
