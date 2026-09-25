package unitary

import (
	"math"
	"sort"
)

const UPLM0BLongDependencyByteSchema = "wingless.up-lm0b-long-dependency-byte.v1"

type UPLM0BMetric struct {
	Arm                         string  `json:"arm"`
	Split                       string  `json:"split"`
	Top1Accuracy                float64 `json:"top1_accuracy"`
	CrossEntropy                float64 `json:"cross_entropy"`
	Perplexity                  float64 `json:"perplexity"`
	DependentFirstByteAccuracy  float64 `json:"dependent_first_byte_accuracy"`
	DependentCrossEntropy       float64 `json:"dependent_cross_entropy"`
	DependentPerplexity         float64 `json:"dependent_perplexity"`
	Tokens                      int     `json:"tokens"`
	DependentTokens             int     `json:"dependent_tokens"`
	RecurrentStateBytes         int     `json:"recurrent_state_bytes"`
}

type UPLM0BLongDependencyByteResult struct {
	Schema              string          `json:"schema"`
	Experiment          string          `json:"experiment"`
	SourceUPLM0ASeal    string          `json:"source_up_lm0a_seal"`
	StateDimension      int             `json:"state_dimension"`
	Epochs              int             `json:"epochs"`
	LearningRate        float64         `json:"learning_rate"`
	TrainSentences      int             `json:"train_sentences"`
	HeldOutSentences    int             `json:"heldout_sentences"`
	AttentionUsed       bool            `json:"attention_used"`
	ExactRecallUsed     bool            `json:"exact_recall_used"`
	PretrainedWeights   bool            `json:"pretrained_weights"`
	Metrics             []UPLM0BMetric  `json:"metrics"`
}

type uplm0bExample struct {
	text      string
	targetPos int
}

func uplm0bCorpus() (train, held []uplm0bExample, alphabet []byte) {
	names:=[]string{"ada","ben","cy","dee","eli","fay"}
	values:=[]string{"amber","cobalt","ivory","jade","mauve","silver"}
	distractors:=[]string{"birds sing.","rain falls slowly.","lamps glow at dusk.","quiet winds cross hills."}
	seen:=map[byte]bool{}
	for ni,name:=range names {
		for vi,value:=range values {
			for di,d:=range distractors {
				prefix:=name+" stores "+value+". "+d+" "+name+" reports "
				text:=prefix+value+".\n"
				ex:=uplm0bExample{text:text,targetPos:len(prefix)}
				if (ni+2*vi+di)%3 != 2 { train=append(train,ex) } else { held=append(held,ex) }
				for i:=0;i<len(text);i++ { seen[text[i]]=true }
			}
		}
	}
	ints:=make([]int,0,len(seen))
	for b:=range seen { ints=append(ints,int(b)) }
	sort.Ints(ints)
	alphabet=make([]byte,len(ints))
	for i,v:=range ints { alphabet[i]=byte(v) }
	return
}

func uplm0bEvaluateModel(m *uplm0aModel, examples []uplm0bExample, stream4 bool) UPLM0BMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	nll,depNLL:=0.0,0.0
	for start:=0; start<len(examples); {
		end:=start+1
		if stream4 {
			end=start+4
			if end>len(examples) { end=len(examples) }
		}
		var h [64]float64
		for e:=start;e<end;e++ {
			ex:=examples[e]
			s:=ex.text
			for t:=0;t<len(s)-1;t++ {
				h=uplm0aStep(h,s[t])
				target:=m.index[int(s[t+1])]
				p:=m.probs(h)
				prob:=p[target]
				if prob<1e-12 { prob=1e-12 }
				total++
				nll-=math.Log(prob)
				if uplm0aArgmax(p)==target { hits++ }
				if t+1==ex.targetPos {
					depTotal++
					depNLL-=math.Log(prob)
					if uplm0aArgmax(p)==target { depHits++ }
				}
			}
			if stream4 { h=uplm0aStep(h,s[len(s)-1]) }
		}
		start=end
	}
	ce:=nll/float64(total)
	depCE:=depNLL/float64(depTotal)
	return UPLM0BMetric{
		Arm:"recurrent64",
		Top1Accuracy:float64(hits)/float64(total),
		CrossEntropy:ce,Perplexity:math.Exp(ce),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		DependentCrossEntropy:depCE,DependentPerplexity:math.Exp(depCE),
		Tokens:total,DependentTokens:depTotal,RecurrentStateBytes:512,
	}
}

func uplm0bBigramEvaluate(b *uplm0aBigram, examples []uplm0bExample, split string) UPLM0BMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	nll,depNLL:=0.0,0.0
	for _,ex:=range examples {
		s:=ex.text
		for t:=0;t<len(s)-1;t++ {
			cur,next:=s[t],s[t+1]
			total++
			if b.best(cur)==next { hits++ }
			prob:=0.0
			if b.totals[int(cur)]>0 {
				prob=float64(b.counts[int(cur)][int(next)])/float64(b.totals[int(cur)])
			}
			if prob<1e-12 { prob=1e-12 }
			nll-=math.Log(prob)
			if t+1==ex.targetPos {
				depTotal++
				depNLL-=math.Log(prob)
				if b.best(cur)==next { depHits++ }
			}
		}
	}
	ce:=nll/float64(total)
	depCE:=depNLL/float64(depTotal)
	return UPLM0BMetric{
		Arm:"bigram_mle",Split:split,
		Top1Accuracy:float64(hits)/float64(total),
		CrossEntropy:ce,Perplexity:math.Exp(ce),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		DependentCrossEntropy:depCE,DependentPerplexity:math.Exp(depCE),
		Tokens:total,DependentTokens:depTotal,RecurrentStateBytes:0,
	}
}

func RunUPLM0B()(UPLM0BLongDependencyByteResult,error){
	train,held,alphabet:=uplm0bCorpus()
	model:=newUPLM0AModel(alphabet)
	const epochs=20
	const lr=0.08
	for epoch:=0;epoch<epochs;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,lr) }
	}
	trainM:=uplm0bEvaluateModel(model,train,false); trainM.Split="train"
	heldM:=uplm0bEvaluateModel(model,held,false); heldM.Split="heldout_recombination"
	streamM:=uplm0bEvaluateModel(model,held,true); streamM.Split="heldout_stream4"

	baseTrain:=make([]uplm0aExample,len(train))
	for i,ex:=range train { baseTrain[i]=uplm0aExample{text:ex.text} }
	bigram:=newUPLM0ABigram(baseTrain,alphabet)
	bTrain:=uplm0bBigramEvaluate(bigram,train,"train")
	bHeld:=uplm0bBigramEvaluate(bigram,held,"heldout_recombination")
	bStream:=uplm0bBigramEvaluate(bigram,held,"heldout_stream4")

	return UPLM0BLongDependencyByteResult{
		Schema:UPLM0BLongDependencyByteSchema,
		Experiment:"UP-LM0B-long-dependency-byte",
		SourceUPLM0ASeal:"ddc471404518ba61ca0c56bd37142d6b29164e35",
		StateDimension:64,Epochs:epochs,LearningRate:lr,
		TrainSentences:len(train),HeldOutSentences:len(held),
		AttentionUsed:false,ExactRecallUsed:false,PretrainedWeights:false,
		Metrics:[]UPLM0BMetric{trainM,heldM,streamM,bTrain,bHeld,bStream},
	},nil
}
