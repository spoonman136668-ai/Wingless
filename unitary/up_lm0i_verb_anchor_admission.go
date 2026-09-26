package unitary

import (
	"math"
	"strings"
)

const UPLM0IVerbAnchorSchema = "wingless.up-lm0i-verb-anchor-admission.v1"

type UPLM0IVerbAnchorResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM0HSeal string `json:"source_up_lm0h_seal"`
	StateDimension int `json:"state_dimension"`
	ExactRecallCap int `json:"exact_recall_cap"`
	ClassifierEpochs int `json:"classifier_epochs"`
	ClassifierLearningRate float64 `json:"classifier_learning_rate"`
	ClassifierThreshold float64 `json:"classifier_threshold"`
	ExplicitTypeAtLearnedInference bool `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool `json:"value_bytes_at_classifier_inference"`
	AttentionUsed bool `json:"attention_used"`
	FutureOracleUsed bool `json:"future_oracle_used"`
	ClassifierMetrics []UPLM0GClassifierMetric `json:"classifier_metrics"`
	Metrics []UPLM0GMetric `json:"metrics"`
}

func uplm0iAnchor(name string,store bool) string {
	if store{return name+" stores"}
	return name+" observes"
}

func uplm0iTrainClassifier()*uplm0gClassifier{
	c:=&uplm0gClassifier{}
	names:=uplm0gNames()
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<4;ni++{
			for et:=0;et<2;et++{
				store:=et==0
				h:=uplm0fEncode(uplm0iAnchor(names[ni],store))
				y:=0.0
				if store{y=1}
				g:=c.prob(h)-y
				for i:=0;i<64;i++{c.w[i]-=0.08*g*h[i]}
				c.b-=0.08*g
			}
		}
	}
	return c
}

func uplm0iClassifierEval(c *uplm0gClassifier,split string,start,end int) UPLM0GClassifierMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	names:=uplm0gNames()
	for ni:=start;ni<end;ni++{
		for et:=0;et<2;et++{
			store:=et==0
			pred:=c.prob(uplm0fEncode(uplm0iAnchor(names[ni],store)))>=0.5
			total++
			if pred==store{hits++}
			if pred&&store{tp++}
			if pred&&!store{fp++}
			if !pred&&store{fn++}
		}
	}
	precision,recall:=1.0,1.0
	if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	return UPLM0GClassifierMetric{Split:split,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func uplm0iEvaluate(model *uplm0aModel,classifier *uplm0gClassifier,examples []uplm0fExample,stream4,learned bool,split string) UPLM0GMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	exactParagraphs,paragraphs:=0,0
	nll:=0.0
	maxEntries:=0
	tp,fp,fn:=0,0,0
	exactBy:=map[int]int{0:0,1:0,2:0,4:0}
	totalBy:=map[int]int{0:0,1:0,2:0,4:0}

	for start:=0;start<len(examples);{
		end:=start+1
		if stream4{end=start+4;if end>len(examples){end=len(examples)}}
		var h [64]float64
		recall:=newUPLM0CRecall()
		for e:=start;e<end;e++{
			s:=examples[e].text
			queryCorrect:=[4]bool{}
			querySeen:=[4]bool{}
			clause:=""
			queryName:=""
			admissionKnown:=false
			admission:=false
			for t:=0;t<len(s)-1;t++{
				b:=s[t]
				h=uplm0aStep(h,b)
				clause+=string(b)
				trim:=strings.TrimSpace(clause)

				if learned&&!admissionKnown{
					fields:=strings.Fields(trim)
					if len(fields)>=2&&(fields[1]=="stores"||fields[1]=="observes"){
						anchor:=fields[0]+" "+fields[1]
						admission=classifier.prob(uplm0fEncode(anchor))>=0.5
						admissionKnown=true
					}
				}

				if strings.HasSuffix(trim," reports"){
					parts:=strings.Split(trim," ")
					if len(parts)>=2{queryName=parts[len(parts)-2]}
				}

				if b=='.'{
					name,value,isStore,isObserve:=uplm0fClauseKV(clause)
					if isStore||isObserve{
						admit:=isStore
						if learned{admit=admissionKnown&&admission}
						if admit{
							recall.write(name,value)
							if isStore{tp++}else{fp++}
						}else if isStore{fn++}
					}
					clause=""
					admissionKnown=false
					admission=false
				}

				targetByte:=s[t+1]
				target:=model.index[int(targetByte)]
				p:=model.probs(h)
				pred:=uplm0aArgmax(p)
				prob:=p[target]
				if strings.HasSuffix(strings.TrimSpace(clause),"reports")&&queryName!=""{
					if value,ok:=recall.values[queryName];ok&&len(value)>0{
						memByte:=value[0]
						pred=model.index[int(memByte)]
						if memByte==targetByte{prob=1}else{prob=1e-12}
					}
				}
				if prob<1e-12{prob=1e-12}
				total++
				nll-=math.Log(prob)
				if pred==target{hits++}
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok{
					depTotal++
					querySeen[qi]=true
					if pred==target{depHits++;queryCorrect[qi]=true}
				}
				if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
			}
			all:=true
			for i:=0;i<4;i++{if !querySeen[i]||!queryCorrect[i]{all=false}}
			paragraphs++
			totalBy[examples[e].updateCount]++
			if all{exactParagraphs++;exactBy[examples[e].updateCount]++}
			if stream4{h=uplm0aStep(h,s[len(s)-1])}
		}
		start=end
	}
	precision,recallScore:=1.0,1.0
	if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	if tp+fn>0{recallScore=float64(tp)/float64(tp+fn)}
	acc:=func(k int)float64{if totalBy[k]==0{return 0};return float64(exactBy[k])/float64(totalBy[k])}
	arm:="explicit_store_admission"
	if learned{arm="verb_anchor_classifier"}
	return UPLM0GMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0I()(UPLM0IVerbAnchorResult,error){
	train,held,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++{for _,ex:=range train{model.trainSentence(ex.text,0.08)}}
	classifier:=uplm0iTrainClassifier()
	return UPLM0IVerbAnchorResult{
		Schema:UPLM0IVerbAnchorSchema,Experiment:"UP-LM0I-verb-anchor-admission",
		SourceUPLM0HSeal:"d9b8b8db0042cd4b55eb05afc1fb4b52d1ef5c67",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,ClassifierThreshold:0.5,
		ExplicitTypeAtLearnedInference:false,ValueBytesAtClassifierInference:false,AttentionUsed:false,FutureOracleUsed:false,
		ClassifierMetrics:[]UPLM0GClassifierMetric{
			uplm0iClassifierEval(classifier,"train_subjects",0,4),
			uplm0iClassifierEval(classifier,"heldout_subjects",4,6),
		},
		Metrics:[]UPLM0GMetric{
			uplm0iEvaluate(model,classifier,held,false,false,"heldout_recombination"),
			uplm0iEvaluate(model,classifier,held,true,false,"heldout_stream4"),
			uplm0iEvaluate(model,classifier,held,false,true,"heldout_recombination"),
			uplm0iEvaluate(model,classifier,held,true,true,"heldout_stream4"),
		},
	},nil
}
