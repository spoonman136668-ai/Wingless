package unitary

import (
	"math"
	"strings"
)

const UPLM0GDomainAdmissionSchema = "wingless.up-lm0g-domain-trained-admission.v1"

type UPLM0GClassifierMetric struct {
	Split     string  `json:"split"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UPLM0GMetric struct {
	Arm                        string  `json:"arm"`
	Split                      string  `json:"split"`
	Top1Accuracy               float64 `json:"top1_accuracy"`
	Perplexity                 float64 `json:"perplexity"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	QuerySetExactAccuracy      float64 `json:"query_set_exact_accuracy"`
	Update0ExactAccuracy       float64 `json:"update0_exact_accuracy"`
	Update1ExactAccuracy       float64 `json:"update1_exact_accuracy"`
	Update2ExactAccuracy       float64 `json:"update2_exact_accuracy"`
	Update4ExactAccuracy       float64 `json:"update4_exact_accuracy"`
	AdmissionPrecision         float64 `json:"admission_precision"`
	AdmissionRecall            float64 `json:"admission_recall"`
	MaxRecallEntries           int     `json:"max_recall_entries"`
	ExactRecallBytes           int     `json:"exact_recall_bytes"`
	RecurrentStateBytes        int     `json:"recurrent_state_bytes"`
}

type UPLM0GDomainAdmissionResult struct {
	Schema                         string                  `json:"schema"`
	Experiment                     string                  `json:"experiment"`
	SourceUPLM0FSeal               string                  `json:"source_up_lm0f_seal"`
	SourceUP96BSeal                string                  `json:"source_up96b_seal"`
	StateDimension                 int                     `json:"state_dimension"`
	ExactRecallCap                 int                     `json:"exact_recall_cap"`
	ClassifierEpochs               int                     `json:"classifier_epochs"`
	ClassifierLearningRate         float64                 `json:"classifier_learning_rate"`
	ClassifierThreshold            float64                 `json:"classifier_threshold"`
	ExplicitTypeAtLearnedInference bool                    `json:"explicit_type_at_learned_inference"`
	AttentionUsed                  bool                    `json:"attention_used"`
	FutureOracleUsed               bool                    `json:"future_oracle_used"`
	ClassifierMetrics              []UPLM0GClassifierMetric `json:"classifier_metrics"`
	Metrics                        []UPLM0GMetric           `json:"metrics"`
}

type uplm0gClassifier struct{ w [64]float64; b float64 }

func (c *uplm0gClassifier) prob(h [64]float64) float64 {
	s:=c.b
	for i:=0;i<64;i++{s+=c.w[i]*h[i]}
	if s>=0 { z:=math.Exp(-s); return 1/(1+z) }
	z:=math.Exp(s); return z/(1+z)
}

func uplm0gNames() []string { return []string{"ada","ben","cy","dee","eli","fay"} }
func uplm0gValues() []string { return []string{"amber","cobalt","ivory","jade","mauve","silver"} }
func uplm0gClause(name,value string, store bool) string {
	if store { return name+" stores "+value+"." }
	return name+" observes "+value+"."
}
func uplm0gTrainSplit(ni,vi,et int) bool { return (ni+2*vi+et)%3!=2 }

func uplm0gTrainClassifier() *uplm0gClassifier {
	c:=&uplm0gClassifier{}
	names,values:=uplm0gNames(),uplm0gValues()
	for epoch:=0;epoch<20;epoch++{
		for ni,name:=range names{
			for vi,value:=range values{
				for et:=0;et<2;et++{
					if !uplm0gTrainSplit(ni,vi,et){continue}
					store:=et==0
					h:=uplm0fEncode(uplm0gClause(name,value,store))
					y:=0.0;if store{y=1}
					g:=c.prob(h)-y
					for i:=0;i<64;i++{c.w[i]-=0.08*g*h[i]}
					c.b-=0.08*g
				}
			}
		}
	}
	return c
}

func uplm0gClassifierEval(c *uplm0gClassifier, split string) UPLM0GClassifierMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	names,values:=uplm0gNames(),uplm0gValues()
	for ni,name:=range names{
		for vi,value:=range values{
			for et:=0;et<2;et++{
				isTrain:=uplm0gTrainSplit(ni,vi,et)
				if split=="train" && !isTrain {continue}
				if split=="heldout_recombination" && isTrain {continue}
				store:=et==0
				pred:=c.prob(uplm0fEncode(uplm0gClause(name,value,store)))>=0.5
				total++;if pred==store{hits++}
				if pred&&store{tp++};if pred&&!store{fp++};if !pred&&store{fn++}
			}
		}
	}
	precision:=1.0;if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	recall:=1.0;if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	return UPLM0GClassifierMetric{Split:split,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func uplm0gEvaluate(model *uplm0aModel, classifier *uplm0gClassifier, examples []uplm0fExample, stream4 bool, learned bool, split string) UPLM0GMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	exactParagraphs,paragraphs:=0,0
	nll:=0.0
	maxEntries:=0
	tp,fp,fn:=0,0,0
	exactBy:=map[int]int{0:0,1:0,2:0,4:0};totalBy:=map[int]int{0:0,1:0,2:0,4:0}

	for start:=0;start<len(examples);{
		end:=start+1;if stream4{end=start+4;if end>len(examples){end=len(examples)}}
		var h [64]float64
		recall:=newUPLM0CRecall()
		for e:=start;e<end;e++{
			s:=examples[e].text
			queryCorrect:=[4]bool{};querySeen:=[4]bool{}
			clause:=""
			queryName:=""
			for t:=0;t<len(s)-1;t++{
				b:=s[t]
				h=uplm0aStep(h,b)
				clause+=string(b)
				trim:=strings.TrimSpace(clause)
				if strings.HasSuffix(trim," reports"){
					parts:=strings.Split(trim," ")
					if len(parts)>=2{queryName=parts[len(parts)-2]}
				}
				if b=='.'{
					name,value,isStore,isObserve:=uplm0fClauseKV(clause)
					if isStore||isObserve{
						admit:=isStore
						if learned{admit=classifier.prob(uplm0fEncode(strings.TrimSpace(clause)))>=0.5}
						if admit{
							recall.write(name,value)
							if isStore{tp++}else{fp++}
						}else if isStore{fn++}
					}
					clause=""
				}
				targetByte:=s[t+1]
				target:=model.index[int(targetByte)]
				p:=model.probs(h);pred:=uplm0aArgmax(p);prob:=p[target]
				if strings.HasSuffix(strings.TrimSpace(clause),"reports")&&queryName!=""{
					if value,ok:=recall.values[queryName];ok&&len(value)>0{
						memByte:=value[0];pred=model.index[int(memByte)]
						if memByte==targetByte{prob=1}else{prob=1e-12}
					}
				}
				if prob<1e-12{prob=1e-12}
				total++;nll-=math.Log(prob);if pred==target{hits++}
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok{
					depTotal++;querySeen[qi]=true
					if pred==target{depHits++;queryCorrect[qi]=true}
				}
				if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
			}
			all:=true;for i:=0;i<4;i++{if !querySeen[i]||!queryCorrect[i]{all=false}}
			paragraphs++;totalBy[examples[e].updateCount]++;if all{exactParagraphs++;exactBy[examples[e].updateCount]++}
			if stream4{h=uplm0aStep(h,s[len(s)-1])}
		}
		start=end
	}
	precision:=1.0;if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	recallScore:=1.0;if tp+fn>0{recallScore=float64(tp)/float64(tp+fn)}
	acc:=func(k int)float64{if totalBy[k]==0{return 0};return float64(exactBy[k])/float64(totalBy[k])}
	arm:="explicit_store_admission";if learned{arm="natural_domain_classifier"}
	return UPLM0GMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0G()(UPLM0GDomainAdmissionResult,error){
	train,held,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++{for _,ex:=range train{model.trainSentence(ex.text,0.08)}}
	classifier:=uplm0gTrainClassifier()
	return UPLM0GDomainAdmissionResult{
		Schema:UPLM0GDomainAdmissionSchema,Experiment:"UP-LM0G-domain-trained-admission",
		SourceUPLM0FSeal:"c135cb029d072be3a78b8cdc3929094aaea6ebf6",SourceUP96BSeal:"87ca9afbb789733ca50cf39ee4ed730938cac9fd",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,ClassifierThreshold:0.5,
		ExplicitTypeAtLearnedInference:false,AttentionUsed:false,FutureOracleUsed:false,
		ClassifierMetrics:[]UPLM0GClassifierMetric{
			uplm0gClassifierEval(classifier,"train"),
			uplm0gClassifierEval(classifier,"heldout_recombination"),
		},
		Metrics:[]UPLM0GMetric{
			uplm0gEvaluate(model,classifier,train,false,false,"train"),
			uplm0gEvaluate(model,classifier,held,false,false,"heldout_recombination"),
			uplm0gEvaluate(model,classifier,held,true,false,"heldout_stream4"),
			uplm0gEvaluate(model,classifier,train,false,true,"train"),
			uplm0gEvaluate(model,classifier,held,false,true,"heldout_recombination"),
			uplm0gEvaluate(model,classifier,held,true,true,"heldout_stream4"),
		},
	},nil
}
