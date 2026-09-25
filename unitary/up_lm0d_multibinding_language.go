package unitary

import (
	"math"
	"strings"
)

const UPLM0DMultiBindingSchema = "wingless.up-lm0d-multibinding-language.v1"

type UPLM0DMetric struct {
	Arm                        string  `json:"arm"`
	Split                      string  `json:"split"`
	Top1Accuracy               float64 `json:"top1_accuracy"`
	Perplexity                 float64 `json:"perplexity"`
	DependentFirstByteAccuracy float64 `json:"dependent_first_byte_accuracy"`
	QuerySetExactAccuracy      float64 `json:"query_set_exact_accuracy"`
	DependentPerplexity        float64 `json:"dependent_perplexity"`
	MaxRecallEntries           int     `json:"max_recall_entries"`
	ExactRecallBytes           int     `json:"exact_recall_bytes"`
	RecurrentStateBytes        int     `json:"recurrent_state_bytes"`
}

type UPLM0DMultiBindingResult struct {
	Schema            string         `json:"schema"`
	Experiment        string         `json:"experiment"`
	SourceUPLM0CSeal  string         `json:"source_up_lm0c_seal"`
	StateDimension    int            `json:"state_dimension"`
	ExactRecallCap    int            `json:"exact_recall_cap"`
	Epochs            int            `json:"epochs"`
	LearningRate      float64        `json:"learning_rate"`
	TrainParagraphs   int            `json:"train_paragraphs"`
	HeldOutParagraphs int            `json:"heldout_paragraphs"`
	AttentionUsed     bool           `json:"attention_used"`
	FutureOracleUsed  bool           `json:"future_oracle_used"`
	Metrics           []UPLM0DMetric `json:"metrics"`
}

type uplm0dExample struct {
	text string
	targetPos [4]int
}

func uplm0dCorpus()(train,held []uplm0dExample, alphabet []byte){
	names:=[]string{"ada","ben","cy","dee","eli","fay"}
	values:=[]string{"amber","cobalt","ivory","jade","mauve","silver"}
	distractors:=[]string{"birds sing.","rain falls slowly.","lamps glow at dusk.","quiet winds cross hills."}
	seen:=map[byte]bool{}
	for n:=0;n<6;n++{
		for v:=0;v<6;v++{
			for d:=0;d<4;d++{
				ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
				vs:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
				s:=""
				for i:=0;i<4;i++{s+=ns[i]+" stores "+vs[i]+". "}
				s+=distractors[d]+" "
				var target [4]int
				for qi:=0;qi<4;qi++{
					idx:=(qi+d)%4
					s+=ns[idx]+" reports "
					target[qi]=len(s)
					s+=vs[idx]+"."
					if qi<3{s+=" "}
				}
				s+="\n"
				ex:=uplm0dExample{text:s,targetPos:target}
				if (n+2*v+d)%3!=2{train=append(train,ex)}else{held=append(held,ex)}
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

func uplm0dTokenBefore(s,suffix string) string {
	if !strings.HasSuffix(s,suffix){return ""}
	prefix:=s[:len(s)-len(suffix)]
	start:=strings.LastIndex(prefix,". ")
	if start>=0{prefix=prefix[start+2:]}
	return strings.TrimSpace(prefix)
}

type uplm0dParser struct {
	seen []byte
	storeName string
	storeValue []byte
	collectStore bool
	queryName string
	queryNext bool
}

func (p *uplm0dParser) observe(b byte, recall *uplm0cRecall){
	if p.collectStore {
		if b=='.' {
			recall.write(p.storeName,string(p.storeValue))
			p.collectStore=false
			p.storeValue=p.storeValue[:0]
			p.storeName=""
		} else {
			p.storeValue=append(p.storeValue,b)
		}
	}
	p.seen=append(p.seen,b)
	s:=string(p.seen)
	if strings.HasSuffix(s," stores "){
		p.storeName=uplm0dTokenBefore(s," stores ")
		p.collectStore=true
		p.storeValue=p.storeValue[:0]
	}
	if strings.HasSuffix(s," reports "){
		p.queryName=uplm0dTokenBefore(s," reports ")
		p.queryNext=true
	}
}

func uplm0dIsTarget(pos int, targets [4]int)(int,bool){
	for i,t:=range targets{if pos==t{return i,true}}
	return 0,false
}

func uplm0dEvaluate(model *uplm0aModel, examples []uplm0dExample, stream4 bool, memory bool, split string) UPLM0DMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	exactParagraphs,paragraphs:=0,0
	nll,depNLL:=0.0,0.0
	maxEntries:=0

	for start:=0;start<len(examples);{
		end:=start+1
		if stream4{end=start+4;if end>len(examples){end=len(examples)}}
		var h [64]float64
		recall:=newUPLM0CRecall()
		for e:=start;e<end;e++{
			var parser uplm0dParser
			queryCorrect:=[4]bool{}
			querySeen:=[4]bool{}
			s:=examples[e].text
			for t:=0;t<len(s)-1;t++{
				h=uplm0aStep(h,s[t])
				if memory{parser.observe(s[t],recall)}
				targetByte:=s[t+1]
				target:=model.index[int(targetByte)]
				p:=model.probs(h)
				pred:=uplm0aArgmax(p)
				prob:=p[target]
				if memory && parser.queryNext{
					if value,ok:=recall.values[parser.queryName];ok && len(value)>0{
						memByte:=value[0]
						pred=model.index[int(memByte)]
						if memByte==targetByte{prob=1}else{prob=1e-12}
					}
					parser.queryNext=false
				}
				if prob<1e-12{prob=1e-12}
				total++;nll-=math.Log(prob);if pred==target{hits++}
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok{
					depTotal++;depNLL-=math.Log(prob);querySeen[qi]=true
					if pred==target{depHits++;queryCorrect[qi]=true}
				}
				if len(recall.order)>maxEntries{maxEntries=len(recall.order)}
			}
			all:=true
			for i:=0;i<4;i++{if !querySeen[i]||!queryCorrect[i]{all=false}}
			paragraphs++;if all{exactParagraphs++}
			if stream4{h=uplm0aStep(h,s[len(s)-1])}
		}
		start=end
	}
	ce:=nll/float64(total);depCE:=depNLL/float64(depTotal)
	arm:="recurrent64";if memory{arm="recurrent64_syntax_recall16"}
	return UPLM0DMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(ce),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		DependentPerplexity:math.Exp(depCE),MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,
		RecurrentStateBytes:512,
	}
}

func RunUPLM0D()(UPLM0DMultiBindingResult,error){
	train,held,alphabet:=uplm0dCorpus()
	model:=newUPLM0AModel(alphabet)
	const epochs=20
	const lr=0.08
	for epoch:=0;epoch<epochs;epoch++{for _,ex:=range train{model.trainSentence(ex.text,lr)}}
	return UPLM0DMultiBindingResult{
		Schema:UPLM0DMultiBindingSchema,Experiment:"UP-LM0D-multibinding-language",
		SourceUPLM0CSeal:"6323443c61da2841f9ad6889d13abacd23590e61",
		StateDimension:64,ExactRecallCap:16,Epochs:epochs,LearningRate:lr,
		TrainParagraphs:len(train),HeldOutParagraphs:len(held),AttentionUsed:false,FutureOracleUsed:false,
		Metrics:[]UPLM0DMetric{
			uplm0dEvaluate(model,train,false,false,"train"),
			uplm0dEvaluate(model,held,false,false,"heldout_recombination"),
			uplm0dEvaluate(model,held,true,false,"heldout_stream4"),
			uplm0dEvaluate(model,train,false,true,"train"),
			uplm0dEvaluate(model,held,false,true,"heldout_recombination"),
			uplm0dEvaluate(model,held,true,true,"heldout_stream4"),
		},
	},nil
}
