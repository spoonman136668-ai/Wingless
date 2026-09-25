package unitary

import (
	"math"
	"sort"
)

const UPLM0AByteMicroLanguageSchema = "wingless.up-lm0a-byte-microlanguage.v1"

type UPLM0AMetric struct {
	Arm                 string  `json:"arm"`
	Split               string  `json:"split"`
	Top1Accuracy        float64 `json:"top1_accuracy"`
	CrossEntropy        float64 `json:"cross_entropy"`
	Perplexity          float64 `json:"perplexity"`
	Tokens              int     `json:"tokens"`
	RecurrentStateBytes int     `json:"recurrent_state_bytes"`
}

type UPLM0AByteMicroLanguageResult struct {
	Schema              string         `json:"schema"`
	Experiment          string         `json:"experiment"`
	SourceUPSQ0Seal     string         `json:"source_up_sq0_seal"`
	StateDimension      int            `json:"state_dimension"`
	Epochs              int            `json:"epochs"`
	LearningRate        float64        `json:"learning_rate"`
	TrainSentences      int            `json:"train_sentences"`
	HeldOutSentences    int            `json:"heldout_sentences"`
	AlphabetBytes       []int          `json:"alphabet_bytes"`
	AttentionUsed       bool           `json:"attention_used"`
	ExactRecallUsed     bool           `json:"exact_recall_used"`
	PretrainedWeights   bool           `json:"pretrained_weights"`
	Metrics             []UPLM0AMetric `json:"metrics"`
}

type uplm0aExample struct {
	text string
	n    int
	v    int
	o    int
}

func uplm0aCorpus() (train, held []uplm0aExample, alphabet []byte) {
	names := []string{"ada", "ben", "cy", "dee", "eli", "fay"}
	verbs := []string{"sees", "likes", "finds", "moves"}
	objects := []string{"red fox", "blue cat", "green owl", "small dog", "bright sun", "calm sea"}
	seen := map[byte]bool{}
	for ni, name := range names {
		for vi, verb := range verbs {
			for oi, object := range objects {
				s := name + " " + verb + " " + object + ".\n"
				ex := uplm0aExample{text:s,n:ni,v:vi,o:oi}
				if (ni+2*vi+oi)%3 != 2 {
					train = append(train, ex)
				} else {
					held = append(held, ex)
				}
				for i:=0;i<len(s);i++ { seen[s[i]]=true }
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

func uplm0aEmbedding(b byte, i int) float64 {
	x:=uint64(b+1)*0x9e3779b97f4a7c15 ^ uint64(i+1)*0xbf58476d1ce4e5b9 ^ 0x6a09e667f3bcc909
	x=sq0Mix64(x)
	if x&1==0 { return -1.0/8.0 }
	return 1.0/8.0
}

func uplm0aStep(h [64]float64, b byte) [64]float64 {
	var next [64]float64
	for i:=0;i<64;i++ {
		src:=(13*i+7)&63
		sign:=1.0
		if i&1==1 { sign=-1.0 }
		next[i]=math.Tanh(0.90*sign*h[src] + 0.35*uplm0aEmbedding(b,i))
	}
	return next
}

type uplm0aModel struct {
	alphabet []byte
	index    [256]int
	w        [][]float64
	b        []float64
}

func newUPLM0AModel(alphabet []byte) *uplm0aModel {
	m:=&uplm0aModel{alphabet:append([]byte(nil),alphabet...)}
	for i:=0;i<256;i++ { m.index[i]=-1 }
	for i,b:=range alphabet { m.index[int(b)]=i }
	m.w=make([][]float64,len(alphabet))
	for i:=range m.w { m.w[i]=make([]float64,64) }
	m.b=make([]float64,len(alphabet))
	return m
}

func (m *uplm0aModel) probs(h [64]float64) []float64 {
	logits:=make([]float64,len(m.alphabet))
	maxLogit:=math.Inf(-1)
	for c:=range logits {
		s:=m.b[c]
		for i:=0;i<64;i++ { s+=m.w[c][i]*h[i] }
		logits[c]=s
		if s>maxLogit { maxLogit=s }
	}
	sum:=0.0
	for i:=range logits {
		logits[i]=math.Exp(logits[i]-maxLogit)
		sum+=logits[i]
	}
	for i:=range logits { logits[i]/=sum }
	return logits
}

func uplm0aArgmax(p []float64) int {
	best:=0
	for i:=1;i<len(p);i++ {
		if p[i]>p[best] { best=i }
	}
	return best
}

func (m *uplm0aModel) trainSentence(s string, lr float64) {
	var h [64]float64
	for t:=0;t<len(s)-1;t++ {
		h=uplm0aStep(h,s[t])
		target:=m.index[int(s[t+1])]
		p:=m.probs(h)
		for c:=range p {
			g:=p[c]
			if c==target { g-=1 }
			for i:=0;i<64;i++ { m.w[c][i]-=lr*g*h[i] }
			m.b[c]-=lr*g
		}
	}
}

func (m *uplm0aModel) evaluateSentences(examples []uplm0aExample) UPLM0AMetric {
	hits,total:=0,0
	nll:=0.0
	for _,ex:=range examples {
		var h [64]float64
		s:=ex.text
		for t:=0;t<len(s)-1;t++ {
			h=uplm0aStep(h,s[t])
			target:=m.index[int(s[t+1])]
			p:=m.probs(h)
			total++
			if uplm0aArgmax(p)==target { hits++ }
			prob:=p[target]
			if prob<1e-12 { prob=1e-12 }
			nll-=math.Log(prob)
		}
	}
	ce:=nll/float64(total)
	return UPLM0AMetric{Arm:"recurrent64",Top1Accuracy:float64(hits)/float64(total),CrossEntropy:ce,Perplexity:math.Exp(ce),Tokens:total,RecurrentStateBytes:512}
}

func (m *uplm0aModel) evaluateStreams(examples []uplm0aExample) UPLM0AMetric {
	hits,total:=0,0
	nll:=0.0
	for start:=0;start<len(examples);start+=4 {
		end:=start+4
		if end>len(examples) { end=len(examples) }
		var h [64]float64
		for e:=start;e<end;e++ {
			s:=examples[e].text
			for t:=0;t<len(s)-1;t++ {
				h=uplm0aStep(h,s[t])
				target:=m.index[int(s[t+1])]
				p:=m.probs(h)
				total++
				if uplm0aArgmax(p)==target { hits++ }
				prob:=p[target]
				if prob<1e-12 { prob=1e-12 }
				nll-=math.Log(prob)
			}
			h=uplm0aStep(h,s[len(s)-1])
		}
	}
	ce:=nll/float64(total)
	return UPLM0AMetric{Arm:"recurrent64",Top1Accuracy:float64(hits)/float64(total),CrossEntropy:ce,Perplexity:math.Exp(ce),Tokens:total,RecurrentStateBytes:512}
}

type uplm0aBigram struct {
	counts [256][256]int
	totals [256]int
	alphabet []byte
}

func newUPLM0ABigram(train []uplm0aExample, alphabet []byte) *uplm0aBigram {
	b:=&uplm0aBigram{alphabet:append([]byte(nil),alphabet...)}
	for _,ex:=range train {
		s:=ex.text
		for t:=0;t<len(s)-1;t++ {
			b.counts[int(s[t])][int(s[t+1])]++
			b.totals[int(s[t])]++
		}
	}
	return b
}

func (b *uplm0aBigram) best(current byte) byte {
	best:=b.alphabet[0]
	bestCount:=-1
	for _,candidate:=range b.alphabet {
		c:=b.counts[int(current)][int(candidate)]
		if c>bestCount || (c==bestCount && candidate<best) {
			best=candidate
			bestCount=c
		}
	}
	return best
}

func (b *uplm0aBigram) evaluate(examples []uplm0aExample, split string) UPLM0AMetric {
	hits,total:=0,0
	nll:=0.0
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
		}
	}
	ce:=nll/float64(total)
	return UPLM0AMetric{Arm:"bigram_mle",Split:split,Top1Accuracy:float64(hits)/float64(total),CrossEntropy:ce,Perplexity:math.Exp(ce),Tokens:total,RecurrentStateBytes:0}
}

func RunUPLM0A()(UPLM0AByteMicroLanguageResult,error){
	train,held,alphabet:=uplm0aCorpus()
	model:=newUPLM0AModel(alphabet)
	const epochs=20
	const lr=0.08
	for epoch:=0;epoch<epochs;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,lr) }
	}
	trainM:=model.evaluateSentences(train); trainM.Split="train"
	heldM:=model.evaluateSentences(held); heldM.Split="heldout_recombination"
	streamM:=model.evaluateStreams(held); streamM.Split="heldout_stream4"
	bigram:=newUPLM0ABigram(train,alphabet)
	bTrain:=bigram.evaluate(train,"train")
	bHeld:=bigram.evaluate(held,"heldout_recombination")
	bStream:=bigram.evaluate(held,"heldout_stream4")
	ints:=make([]int,len(alphabet))
	for i,b:=range alphabet { ints[i]=int(b) }
	return UPLM0AByteMicroLanguageResult{
		Schema:UPLM0AByteMicroLanguageSchema,
		Experiment:"UP-LM0A-byte-microlanguage",
		SourceUPSQ0Seal:"8de9e11f26cc44c75f51d26cb606f97e35816479",
		StateDimension:64,Epochs:epochs,LearningRate:lr,
		TrainSentences:len(train),HeldOutSentences:len(held),
		AlphabetBytes:ints,AttentionUsed:false,ExactRecallUsed:false,PretrainedWeights:false,
		Metrics:[]UPLM0AMetric{trainM,heldM,streamM,bTrain,bHeld,bStream},
	},nil
}
