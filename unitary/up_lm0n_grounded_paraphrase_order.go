package unitary

import (
	"math"
	"strings"
)

const UPLM0NGroundedParaphraseOrderSchema = "wingless.up-lm0n-grounded-paraphrase-order.v1"

type UPLM0NGroundedParaphraseOrderResult struct {
	Schema                          string                   `json:"schema"`
	Experiment                      string                   `json:"experiment"`
	SourceUPLM0MSeal                string                   `json:"source_up_lm0m_seal"`
	StateDimension                  int                      `json:"state_dimension"`
	ExactRecallCap                  int                      `json:"exact_recall_cap"`
	ClassifierEpochs                int                      `json:"classifier_epochs"`
	ClassifierLearningRate          float64                  `json:"classifier_learning_rate"`
	ExplicitTypeAtLearnedInference  bool                     `json:"explicit_type_at_learned_inference"`
	ValueBytesAtClassifierInference bool                    `json:"value_bytes_at_classifier_inference"`
	AttentionUsed                   bool                     `json:"attention_used"`
	FutureOracleUsed                bool                     `json:"future_oracle_used"`
	LanguageRetrainingOnParaphrases bool                     `json:"language_retraining_on_paraphrases"`
	ClassifierMetrics               []UPLM0JClassifierMetric `json:"classifier_metrics"`
	Metrics                         []UPLM0JMetric           `json:"metrics"`
}

func uplm0nAnchor(name, verb string) string { return name+" "+verb }

func uplm0nClassForVerb(verb string) int {
	switch verb {
	case "stores","saves","archives":
		return uplm0jStore
	case "observes","sees","notices":
		return uplm0jObserve
	case "reports","recalls","recounts":
		return uplm0jReport
	default:
		return -1
	}
}

func uplm0nTrainStep(c *uplm0jClassifier, anchor string, target int) {
	h:=uplm0fEncode(anchor)
	p:=c.probs(h)
	for k:=0;k<3;k++ {
		g:=p[k]
		if k==target { g-=1 }
		for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
		c.b[k]-=0.08*g
	}
}

func uplm0nTrainClassifier() *uplm0jClassifier {
	c:=uplm0jTrainClassifier()
	names:=uplm0gNames()
	newVerbs:=[]string{"saves","sees","recalls"}
	baseVerbs:=[]string{"stores","observes","reports"}
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ {
			for class,verb:=range newVerbs {
				uplm0nTrainStep(c,uplm0nAnchor(names[ni],verb),class)
			}
		}
		for class,verb:=range baseVerbs {
			uplm0nTrainStep(c,uplm0nAnchor(names[0],verb),class)
		}
	}
	return c
}

func uplm0nClassifierEval(c *uplm0jClassifier, split,label string, class int, verbs []string, start,end int) UPLM0JClassifierMetric {
	names:=uplm0gNames()
	total,hits,tp,fp,fn:=0,0,0,0,0
	for ni:=start;ni<end;ni++ {
		for _,verb:=range verbs {
			target:=uplm0nClassForVerb(verb)
			pred:=uplm0jPredict(c,uplm0nAnchor(names[ni],verb))
			total++
			if pred==target { hits++ }
			if class>=0 {
				if pred==class&&target==class { tp++ }
				if pred==class&&target!=class { fp++ }
				if pred!=class&&target==class { fn++ }
			}
		}
	}
	precision,recall:=1.0,1.0
	if class>=0 {
		if tp+fp>0 { precision=float64(tp)/float64(tp+fp) }
		if tp+fn>0 { recall=float64(tp)/float64(tp+fn) }
	}
	return UPLM0JClassifierMetric{Split:split,Label:label,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func uplm0nExample(n,v,p int,order string) uplm0fExample {
	names:=uplm0gNames()
	values:=uplm0gValues()
	ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
	initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
	latest:=initial
	uc:=uplm0eUpdateCount(p)
	s:=""
	var targets [4]int

	storeInitial:=func(i int){ s+=ns[i]+" saves "+initial[i]+". " }
	observe:=func(i int){
		obs:=values[(v+i+3)%6]
		s+=ns[(i+p)%4]+" sees "+obs+". "
	}
	update:=func(i int){
		if i>=uc { return }
		latest[i]=values[(v+i+2)%6]
		s+=ns[i]+" saves "+latest[i]+". "
	}
	report:=func(i int,last bool){
		s+=ns[i]+" recalls "
		targets[i]=len(s)
		s+=latest[i]+"."
		if !last { s+=" " }
	}

	switch order {
	case "per_name":
		for i:=0;i<4;i++ { storeInitial(i);observe(i);update(i);report(i,i==3) }
	case "paired_names":
		for pair:=0;pair<4;pair+=2 {
			for i:=pair;i<pair+2;i++ { storeInitial(i) }
			for i:=pair;i<pair+2;i++ { observe(i) }
			for i:=pair;i<pair+2;i++ { update(i) }
			for i:=pair;i<pair+2;i++ { report(i,pair==2&&i==3) }
		}
	case "stores_then_local_reports":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i);update(i);report(i,i==3) }
	case "reverse_report_tail":
		for i:=0;i<4;i++ { storeInitial(i) }
		for i:=0;i<4;i++ { observe(i) }
		for i:=0;i<4;i++ { update(i) }
		for i:=3;i>=0;i-- { report(i,i==0) }
	}
	s+="\n"
	return uplm0fExample{text:s,targetPos:targets,updateCount:uc}
}

func uplm0nHeldout(order string) []uplm0fExample {
	out:=[]uplm0fExample{}
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				if (n+2*v+p)%3!=2 { continue }
				out=append(out,uplm0nExample(n,v,p,order))
			}
		}
	}
	return out
}

func uplm0nEvaluate(model *uplm0aModel,classifier *uplm0jClassifier,examples []uplm0fExample,streamSize int,learned bool,split string) UPLM0JMetric {
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
	arm:="explicit_event_routing"
	if learned { arm="learned_grounded_paraphrase_routing" }
	eventAccuracy:=1.0
	if eventTotal>0 { eventAccuracy=float64(eventHits)/float64(eventTotal) }
	reportAccuracy:=1.0
	if reportTotal>0 { reportAccuracy=float64(reportHits)/float64(reportTotal) }

	return UPLM0JMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,EventRoutingAccuracy:eventAccuracy,ReportRoutingAccuracy:reportAccuracy,
		MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0N()(UPLM0NGroundedParaphraseOrderResult,error){
	train,_,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ { for _,ex:=range train { model.trainSentence(ex.text,0.08) } }
	classifier:=uplm0nTrainClassifier()

	result:=UPLM0NGroundedParaphraseOrderResult{
		Schema:UPLM0NGroundedParaphraseOrderSchema,
		Experiment:"UP-LM0N-grounded-paraphrase-order",
		SourceUPLM0MSeal:"e4cee8ac2b372c89fbbc68b5e06c8747a6bfa3eb",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,
		ExplicitTypeAtLearnedInference:false,ValueBytesAtClassifierInference:false,
		AttentionUsed:false,FutureOracleUsed:false,LanguageRetrainingOnParaphrases:false,
	}
	paraphrase:=[]string{"saves","sees","recalls"}
	base:=[]string{"stores","observes","reports"}
	result.ClassifierMetrics=[]UPLM0JClassifierMetric{
		uplm0nClassifierEval(classifier,"heldout_paraphrase_subjects","all",-1,paraphrase,4,6),
		uplm0nClassifierEval(classifier,"heldout_paraphrase_subjects","STORE",uplm0jStore,paraphrase,4,6),
		uplm0nClassifierEval(classifier,"heldout_paraphrase_subjects","OBSERVE",uplm0jObserve,paraphrase,4,6),
		uplm0nClassifierEval(classifier,"heldout_paraphrase_subjects","REPORT",uplm0jReport,paraphrase,4,6),
		uplm0nClassifierEval(classifier,"retained_base_subjects","all",-1,base,0,6),
		uplm0nClassifierEval(classifier,"retained_base_subjects","STORE",uplm0jStore,base,0,6),
		uplm0nClassifierEval(classifier,"retained_base_subjects","OBSERVE",uplm0jObserve,base,0,6),
		uplm0nClassifierEval(classifier,"retained_base_subjects","REPORT",uplm0jReport,base,0,6),
	}
	for _,order:=range []string{"per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
		held:=uplm0nHeldout(order)
		for _,arm:=range []struct{learned bool}{ {false},{true} } {
			result.Metrics=append(result.Metrics,
				uplm0nEvaluate(model,classifier,held,1,arm.learned,order+"_stream1"),
				uplm0nEvaluate(model,classifier,held,4,arm.learned,order+"_stream4"),
			)
		}
	}
	return result,nil
}
