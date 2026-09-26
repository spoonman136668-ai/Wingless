package unitary

import (
	"fmt"
	"math"
	"strings"
)

const UPLM1MFourthByteIsolationSchema = "wingless.up-lm1m-fourth-byte-isolation.v1"

type UPLM1MByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Tokens       int     `json:"tokens"`
}

type UPLM1MRoutingMetric struct {
	Arm    string       `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1MFourthByteIsolationResult struct {
	Schema                     string                `json:"schema"`
	Experiment                 string                `json:"experiment"`
	SourceUPLM1LSeal           string                `json:"source_up_lm1l_seal"`
	StateDimension             int                   `json:"state_dimension"`
	ExactRecallCap             int                   `json:"exact_recall_cap"`
	BasePretrainEpochs         int                   `json:"base_pretrain_epochs"`
	JointInterleavedEpochs     int                   `json:"joint_interleaved_epochs"`
	ThreeFamilyAdaptEpochs     int                   `json:"three_family_adapt_epochs"`
	FourthAdaptEpochs          int                   `json:"fourth_adapt_epochs"`
	LearningRate               float64               `json:"learning_rate"`
	RouterRetrainingUsed       bool                  `json:"router_retraining_used"`
	RecurrentParametersTrained bool                  `json:"recurrent_parameters_trained"`
	RecallCapChanged           bool                  `json:"recall_cap_changed"`
	AttentionUsed              bool                  `json:"attention_used"`
	FutureOracleUsed           bool                  `json:"future_oracle_used"`
	ByteMetrics                []UPLM1MByteMetric    `json:"byte_metrics"`
	RoutingMetrics             []UPLM1MRoutingMetric `json:"routing_metrics"`
}

func uplm1mFourthCorpus()(train,held []uplm0fExample){
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				ex:=uplm1jExample(n,v,p,"block")
				if (n+2*v+p)%3!=2 { train=append(train,ex) } else { held=append(held,ex) }
			}
		}
	}
	return
}

func uplm1mPredict(c *uplm0jClassifier,name,verb string,d [64]float64) int {
	h:=uplm0fEncode(name+" "+verb)
	if verb=="states" { h=uplm1lEncode("states_store_orthogonal",name,verb,d) }
	return uplm0jArgmax(c.probs(h))
}

func uplm1mEvaluate(model *uplm0aModel,classifier *uplm0jClassifier,d [64]float64,examples []uplm0fExample,streamSize int,split string) UPLM0JMetric {
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
						trueClass=uplm0nClassForVerb(fields[1])
						if trueClass>=0 {
							routeClass=uplm1mPredict(classifier,fields[0],fields[1],d)
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
				okTarget:=target>=0
				if !okTarget { target=0 }
				p:=model.probs(h)
				pred:=uplm0aArgmax(p)
				prob:=p[target]

				if reportOverridePending&&queryName!="" {
					if value,ok:=recall.values[queryName];ok&&len(value)>0 {
						memByte:=value[0]
						memIndex:=model.index[int(memByte)]
						if memIndex>=0 { pred=memIndex }
						if memByte==targetByte { prob=1 } else { prob=1e-12 }
					}
					reportOverridePending=false
				}

				if prob<1e-12 { prob=1e-12 }
				total++
				nll-=math.Log(prob)
				if okTarget&&pred==target { hits++ }
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok {
					depTotal++
					querySeen[qi]=true
					if okTarget&&pred==target { depHits++;queryCorrect[qi]=true }
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
			for i:=0;i<4;i++ { if !querySeen[i]||!queryCorrect[i] { all=false } }
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
	acc:=func(k int)float64{ if totalBy[k]==0{return 0};return float64(exactBy[k])/float64(totalBy[k]) }
	eventAccuracy:=1.0
	if eventTotal>0 { eventAccuracy=float64(eventHits)/float64(eventTotal) }
	reportAccuracy:=1.0
	if reportTotal>0 { reportAccuracy=float64(reportHits)/float64(reportTotal) }

	return UPLM0JMetric{
		Arm:"learned_corrected_fourth_routing",Split:split,
		Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,
		EventRoutingAccuracy:eventAccuracy,ReportRoutingAccuracy:reportAccuracy,
		MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func uplm1mStartModel()(model *uplm0aModel,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld []uplm0fExample,err error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld=uplm1dCorpus()
	fourthTrain,fourthHeld=uplm1mFourthCorpus()
	if len(baseTrain)!=len(paraTrain)||len(baseTrain)!=len(thirdTrain)||len(baseTrain)!=len(fourthTrain) {
		err=fmt.Errorf("training corpus count mismatch: base=%d para=%d third=%d fourth=%d",len(baseTrain),len(paraTrain),len(thirdTrain),len(fourthTrain))
		return
	}
	model=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(model,baseTrain,paraTrain,20)
	model,_=uplm1dExtendAlphabet(model,thirdTrain,thirdHeld,fourthTrain,fourthHeld)
	uplm1hTrain(model,baseTrain,paraTrain,thirdTrain,"strang_split_base")
	return
}

func uplm1mTrainArm(model *uplm0aModel,arm string,baseTrain,paraTrain,thirdTrain,fourthTrain []uplm0fExample){
	switch arm {
	case "no_fourth_byte_adaptation":
		return
	case "fourth_only_4ep":
		for epoch:=0;epoch<4;epoch++ {
			for _,ex:=range fourthTrain { model.trainSentence(ex.text,0.08) }
		}
	case "four_family_interleaved_4ep":
		for epoch:=0;epoch<4;epoch++ {
			for i:=range baseTrain {
				model.trainSentence(baseTrain[i].text,0.08)
				model.trainSentence(paraTrain[i].text,0.08)
				model.trainSentence(thirdTrain[i].text,0.08)
				model.trainSentence(fourthTrain[i].text,0.08)
			}
		}
	}
}

func RunUPLM1M()(UPLM1MFourthByteIsolationResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1MFourthByteIsolationResult{},err }
	d:=uplm1lStoreDirection()
	classifier:=uplm1lTrainClassifier("states_store_orthogonal",d)

	result:=UPLM1MFourthByteIsolationResult{
		Schema:UPLM1MFourthByteIsolationSchema,Experiment:"UP-LM1M-fourth-byte-isolation",
		SourceUPLM1LSeal:"bab2384f16a154b0af448b96fd096cca3e897a2c",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,ThreeFamilyAdaptEpochs:4,FourthAdaptEpochs:4,
		LearningRate:0.08,RouterRetrainingUsed:false,RecurrentParametersTrained:false,RecallCapChanged:false,AttentionUsed:false,FutureOracleUsed:false,
	}

	arms:=[]string{"no_fourth_byte_adaptation","fourth_only_4ep","four_family_interleaved_4ep"}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(start)
		uplm1mTrainArm(model,arm,baseTrain,paraTrain,thirdTrain,fourthTrain)
		for _,x:=range []struct{name string;data []uplm0fExample}{
			{"base_heldout",baseHeld},
			{"paraphrase_heldout",paraHeld},
			{"third_heldout",thirdHeld},
			{"fourth_heldout",fourthHeld},
		}{
			m:=uplm0oByteEval(model,x.data,0,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1MByteMetric{Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,Tokens:m.Tokens})
		}
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1jHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,UPLM1MRoutingMetric{
					Arm:arm,Metric:uplm1mEvaluate(model,classifier,d,held,stream,"fourth_"+order+"_stream"+itoa(stream)),
				})
			}
		}
	}
	return result,nil
}
