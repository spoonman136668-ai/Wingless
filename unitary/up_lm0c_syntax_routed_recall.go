package unitary

import (
	"math"
	"strings"
)

const UPLM0CSyntaxRecallSchema = "wingless.up-lm0c-syntax-routed-recall.v1"

type UPLM0CMetric struct {
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
	MaxRecallEntries            int     `json:"max_recall_entries"`
	ExactRecallBytes            int     `json:"exact_recall_bytes"`
}

type UPLM0CSyntaxRecallResult struct {
	Schema              string         `json:"schema"`
	Experiment          string         `json:"experiment"`
	SourceUPLM0BSeal    string         `json:"source_up_lm0b_seal"`
	SourceUP91BSeal     string         `json:"source_up91b_seal"`
	StateDimension      int            `json:"state_dimension"`
	ExactRecallCap      int            `json:"exact_recall_cap"`
	Epochs              int            `json:"epochs"`
	LearningRate        float64        `json:"learning_rate"`
	AttentionUsed       bool           `json:"attention_used"`
	FutureOracleUsed    bool           `json:"future_oracle_used"`
	Metrics             []UPLM0CMetric `json:"metrics"`
}

type uplm0cRecall struct {
	values map[string]string
	order  []string
}

func newUPLM0CRecall() *uplm0cRecall {
	return &uplm0cRecall{values: map[string]string{}}
}

func (r *uplm0cRecall) write(key, value string) {
	if _, ok := r.values[key]; ok {
		r.values[key] = value
		return
	}
	if len(r.order) >= 16 {
		old := r.order[0]
		r.order = r.order[1:]
		delete(r.values, old)
	}
	r.order = append(r.order, key)
	r.values[key] = value
}

type uplm0cParser struct {
	seen          []byte
	name          string
	collectValue  bool
	valueBuf      []byte
	queryNext     bool
}

func (p *uplm0cParser) observe(b byte, recall *uplm0cRecall) {
	if p.collectValue {
		if b == '.' {
			recall.write(p.name, string(p.valueBuf))
			p.collectValue = false
			p.valueBuf = p.valueBuf[:0]
		} else {
			p.valueBuf = append(p.valueBuf, b)
		}
	}

	p.seen = append(p.seen, b)
	s := string(p.seen)
	const stores = " stores "
	if p.name == "" && strings.HasSuffix(s, stores) {
		p.name = s[:len(s)-len(stores)]
		p.collectValue = true
		p.valueBuf = p.valueBuf[:0]
	}
	if p.name != "" && strings.HasSuffix(s, p.name+" reports ") {
		p.queryNext = true
	}
}

func uplm0cFromBaseline(m UPLM0BMetric) UPLM0CMetric {
	return UPLM0CMetric{
		Arm:"recurrent64",Split:m.Split,
		Top1Accuracy:m.Top1Accuracy,CrossEntropy:m.CrossEntropy,Perplexity:m.Perplexity,
		DependentFirstByteAccuracy:m.DependentFirstByteAccuracy,
		DependentCrossEntropy:m.DependentCrossEntropy,DependentPerplexity:m.DependentPerplexity,
		Tokens:m.Tokens,DependentTokens:m.DependentTokens,RecurrentStateBytes:512,
	}
}

func uplm0cEvaluateMemory(model *uplm0aModel, examples []uplm0bExample, stream4 bool, split string) UPLM0CMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	nll,depNLL:=0.0,0.0
	maxEntries:=0

	for start:=0;start<len(examples); {
		end:=start+1
		if stream4 {
			end=start+4
			if end>len(examples) { end=len(examples) }
		}
		var h [64]float64
		recall:=newUPLM0CRecall()
		for e:=start;e<end;e++ {
			var parser uplm0cParser
			s:=examples[e].text
			for t:=0;t<len(s)-1;t++ {
				h=uplm0aStep(h,s[t])
				parser.observe(s[t],recall)
				targetByte:=s[t+1]
				target:=model.index[int(targetByte)]
				p:=model.probs(h)
				pred:=uplm0aArgmax(p)
				prob:=p[target]

				if parser.queryNext {
					if value,ok:=recall.values[parser.name];ok && len(value)>0 {
						memByte:=value[0]
						pred=model.index[int(memByte)]
						if memByte==targetByte { prob=1 } else { prob=1e-12 }
					}
					parser.queryNext=false
				}

				if prob<1e-12 { prob=1e-12 }
				total++
				nll-=math.Log(prob)
				if pred==target { hits++ }

				if t+1==examples[e].targetPos {
					depTotal++
					depNLL-=math.Log(prob)
					if pred==target { depHits++ }
				}
				if len(recall.order)>maxEntries { maxEntries=len(recall.order) }
			}
			if stream4 { h=uplm0aStep(h,s[len(s)-1]) }
		}
		start=end
	}
	ce:=nll/float64(total)
	depCE:=depNLL/float64(depTotal)
	return UPLM0CMetric{
		Arm:"recurrent64_syntax_recall16",Split:split,
		Top1Accuracy:float64(hits)/float64(total),CrossEntropy:ce,Perplexity:math.Exp(ce),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		DependentCrossEntropy:depCE,DependentPerplexity:math.Exp(depCE),
		Tokens:total,DependentTokens:depTotal,RecurrentStateBytes:512,
		MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,
	}
}

func RunUPLM0C()(UPLM0CSyntaxRecallResult,error){
	train,held,alphabet:=uplm0bCorpus()
	model:=newUPLM0AModel(alphabet)
	const epochs=20
	const lr=0.08
	for epoch:=0;epoch<epochs;epoch++ {
		for _,ex:=range train { model.trainSentence(ex.text,lr) }
	}

	baseTrain:=uplm0bEvaluateModel(model,train,false);baseTrain.Split="train"
	baseHeld:=uplm0bEvaluateModel(model,held,false);baseHeld.Split="heldout_recombination"
	baseStream:=uplm0bEvaluateModel(model,held,true);baseStream.Split="heldout_stream4"

	return UPLM0CSyntaxRecallResult{
		Schema:UPLM0CSyntaxRecallSchema,
		Experiment:"UP-LM0C-syntax-routed-recall",
		SourceUPLM0BSeal:"413285fd458b87db03dde2667c59d73148c65b4a",
		SourceUP91BSeal:"9c2b5c8195bba9589afb488236c5ce6ec9271939",
		StateDimension:64,ExactRecallCap:16,Epochs:epochs,LearningRate:lr,
		AttentionUsed:false,FutureOracleUsed:false,
		Metrics:[]UPLM0CMetric{
			uplm0cFromBaseline(baseTrain),
			uplm0cFromBaseline(baseHeld),
			uplm0cFromBaseline(baseStream),
			uplm0cEvaluateMemory(model,train,false,"train"),
			uplm0cEvaluateMemory(model,held,false,"heldout_recombination"),
			uplm0cEvaluateMemory(model,held,true,"heldout_stream4"),
		},
	},nil
}
