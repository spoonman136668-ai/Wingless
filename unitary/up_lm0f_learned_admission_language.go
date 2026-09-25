package unitary

import (
	"math"
	"strings"
)

const UPLM0FLearnedAdmissionSchema = "wingless.up-lm0f-learned-admission-language.v1"

type UPLM0FMetric struct {
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

type UPLM0FLearnedAdmissionResult struct {
	Schema                         string         `json:"schema"`
	Experiment                     string         `json:"experiment"`
	SourceUPLM0ESeal               string         `json:"source_up_lm0e_seal"`
	SourceUP95BSeal                string         `json:"source_up95b_seal"`
	StateDimension                 int            `json:"state_dimension"`
	ExactRecallCap                 int            `json:"exact_recall_cap"`
	ClassifierEpochs               int            `json:"classifier_epochs"`
	ClassifierLearningRate         float64        `json:"classifier_learning_rate"`
	ClassifierThreshold            float64        `json:"classifier_threshold"`
	ExplicitTypeAtLearnedInference bool           `json:"explicit_type_at_learned_inference"`
	AttentionUsed                  bool           `json:"attention_used"`
	FutureOracleUsed               bool           `json:"future_oracle_used"`
	Metrics                        []UPLM0FMetric  `json:"metrics"`
}

type uplm0fExample struct {
	text string
	targetPos [4]int
	updateCount int
}

func uplm0fCorpus()(train,held []uplm0fExample, alphabet []byte){
	names:=[]string{"ada","ben","cy","dee","eli","fay"}
	values:=[]string{"amber","cobalt","ivory","jade","mauve","silver"}
	seen:=map[byte]bool{}
	for n:=0;n<6;n++{
		for v:=0;v<6;v++{
			for p:=0;p<4;p++{
				ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
				initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
				latest:=initial
				s:=""
				for i:=0;i<4;i++{s+=ns[i]+" stores "+initial[i]+". "}
				for i:=0;i<4;i++{
					obs:=values[(v+i+3)%6]
					s+=ns[(i+p)%4]+" observes "+obs+". "
				}
				uc:=uplm0eUpdateCount(p)
				for i:=0;i<uc;i++{
					latest[i]=values[(v+i+2)%6]
					s+=ns[i]+" stores "+latest[i]+". "
				}
				var targets [4]int
				for qi:=0;qi<4;qi++{
					idx:=(qi+p)%4
					s+=ns[idx]+" reports "
					targets[qi]=len(s)
					s+=latest[idx]+"."
					if qi<3{s+=" "}
				}
				s+="\n"
				ex:=uplm0fExample{text:s,targetPos:targets,updateCount:uc}
				if (n+2*v+p)%3!=2{train=append(train,ex)}else{held=append(held,ex)}
				for i:=0;i<len(s);i++{seen[s[i]]=true}
			}
		}
	}
	ints:=make([]int,0,len(seen))
	for b:=range seen{ints=append(ints,int(b))}
	for i:=0;i<len(ints);i++{for j:=i+1;j<len(ints);j++{if ints[j]<ints[i]{ints[i],ints[j]=ints[j],ints[i]}}}
	alphabet=make([]byte,len(ints));for i,v:=range ints{alphabet[i]=byte(v)}
	return
}

func uplm0fEmbed(b byte,i int) float64 {
	x:=uint64(b+1)*0x9e3779b97f4a7c15^uint64(i+1)*0xbf58476d1ce4e5b9^0x510e527fade682d1
	x=sq0Mix64(x)
	if x&1==0{return -0.125}
	return 0.125
}
func uplm0fEncode(s string)[64]float64{
	var h [64]float64
	for j:=0;j<len(s);j++{
		var next [64]float64
		for i:=0;i<64;i++{
			src:=(13*i+7)&63
			sign:=1.0;if i&1==1{sign=-1}
			next[i]=math.Tanh(0.90*sign*h[src]+0.35*uplm0fEmbed(s[j],i))
		}
		h=next
	}
	return h
}
type uplm0fClassifier struct{w [64]float64;b float64}
func uplm0fSigmoid(x float64)float64{if x>=0{z:=math.Exp(-x);return 1/(1+z)};z:=math.Exp(x);return z/(1+z)}
func (c *uplm0fClassifier) prob(h [64]float64)float64{s:=c.b;for i:=0;i<64;i++{s+=c.w[i]*h[i]};return uplm0fSigmoid(s)}
func uplm0fTwo(n int)string{return string([]byte{'0'+byte((n/10)%10),'0'+byte(n%10)})}
func uplm0fTrainClassifier()*uplm0fClassifier{
	c:=&uplm0fClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for key:=0;key<24;key++{
			for value:=0;value<32;value++{
				for _,store:=range []bool{false,true}{
					verb:=" observes ";if store{verb=" stores "}
					s:="k"+uplm0fTwo(key)+verb+"v"+uplm0fTwo(value)+"."
					h:=uplm0fEncode(s);y:=0.0;if store{y=1}
					g:=c.prob(h)-y
					for i:=0;i<64;i++{c.w[i]-=0.08*g*h[i]}
					c.b-=0.08*g
				}
			}
		}
	}
	return c
}

func uplm0fClauseKV(clause string)(name,value string,store,observe bool){
	clause=strings.TrimSpace(clause)
	clause=strings.TrimSuffix(clause,".")
	if x:=strings.SplitN(clause," stores ",2);len(x)==2{return strings.TrimSpace(x[0]),strings.TrimSpace(x[1]),true,false}
	if x:=strings.SplitN(clause," observes ",2);len(x)==2{return strings.TrimSpace(x[0]),strings.TrimSpace(x[1]),false,true}
	return "","",false,false
}

func uplm0fEvaluate(model *uplm0aModel, classifier *uplm0fClassifier, examples []uplm0fExample, stream4 bool, learned bool, split string) UPLM0FMetric {
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
				if strings.HasSuffix(strings.TrimSpace(clause),"reports") && queryName!=""{
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
	arm:="explicit_store_admission";if learned{arm="learned_local_admission"}
	return UPLM0FMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(nll/float64(total)),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		AdmissionPrecision:precision,AdmissionRecall:recallScore,MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0F()(UPLM0FLearnedAdmissionResult,error){
	train,held,alphabet:=uplm0fCorpus()
	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++{for _,ex:=range train{model.trainSentence(ex.text,0.08)}}
	classifier:=uplm0fTrainClassifier()
	return UPLM0FLearnedAdmissionResult{
		Schema:UPLM0FLearnedAdmissionSchema,Experiment:"UP-LM0F-learned-admission-language",
		SourceUPLM0ESeal:"5a4b8d05abe54e1b1c0d8107c4d4f9434cf65801",SourceUP95BSeal:"f15ff278753d3af1786b38f9958e41ded3c00fc5",
		StateDimension:64,ExactRecallCap:16,ClassifierEpochs:20,ClassifierLearningRate:0.08,ClassifierThreshold:0.5,
		ExplicitTypeAtLearnedInference:false,AttentionUsed:false,FutureOracleUsed:false,
		Metrics:[]UPLM0FMetric{
			uplm0fEvaluate(model,classifier,train,false,false,"train"),
			uplm0fEvaluate(model,classifier,held,false,false,"heldout_recombination"),
			uplm0fEvaluate(model,classifier,held,true,false,"heldout_stream4"),
			uplm0fEvaluate(model,classifier,train,false,true,"train"),
			uplm0fEvaluate(model,classifier,held,false,true,"heldout_recombination"),
			uplm0fEvaluate(model,classifier,held,true,true,"heldout_stream4"),
		},
	},nil
}
