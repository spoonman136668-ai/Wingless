package unitary

import "math"

const UPLM0EMutableBindingSchema = "wingless.up-lm0e-mutable-binding-language.v1"

type UPLM0EMetric struct {
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
	MaxRecallEntries           int     `json:"max_recall_entries"`
	ExactRecallBytes           int     `json:"exact_recall_bytes"`
	RecurrentStateBytes        int     `json:"recurrent_state_bytes"`
}

type UPLM0EMutableBindingResult struct {
	Schema            string         `json:"schema"`
	Experiment        string         `json:"experiment"`
	SourceUPLM0DSeal  string         `json:"source_up_lm0d_seal"`
	StateDimension    int            `json:"state_dimension"`
	ExactRecallCap    int            `json:"exact_recall_cap"`
	Epochs            int            `json:"epochs"`
	LearningRate      float64        `json:"learning_rate"`
	TrainParagraphs   int            `json:"train_paragraphs"`
	HeldOutParagraphs int            `json:"heldout_paragraphs"`
	AttentionUsed     bool           `json:"attention_used"`
	FutureOracleUsed  bool           `json:"future_oracle_used"`
	Metrics           []UPLM0EMetric `json:"metrics"`
}

type uplm0eExample struct {
	text        string
	targetPos   [4]int
	updateCount int
}

func uplm0eUpdateCount(d int) int {
	switch d {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	default:
		return 4
	}
}

func uplm0eCorpus()(train,held []uplm0eExample, alphabet []byte){
	names:=[]string{"ada","ben","cy","dee","eli","fay"}
	values:=[]string{"amber","cobalt","ivory","jade","mauve","silver"}
	distractors:=[]string{"birds sing.","rain falls slowly.","lamps glow at dusk.","quiet winds cross hills."}
	seen:=map[byte]bool{}
	for n:=0;n<6;n++{
		for v:=0;v<6;v++{
			for d:=0;d<4;d++{
				ns:=[4]string{names[n],names[(n+1)%6],names[(n+2)%6],names[(n+3)%6]}
				initial:=[4]string{values[v],values[(v+1)%6],values[(v+2)%6],values[(v+3)%6]}
				latest:=initial
				s:=""
				for i:=0;i<4;i++{s+=ns[i]+" stores "+initial[i]+". "}
				s+=distractors[d]+" "
				uc:=uplm0eUpdateCount(d)
				for i:=0;i<uc;i++{
					latest[i]=values[(v+i+2)%6]
					s+=ns[i]+" stores "+latest[i]+". "
				}
				s+=distractors[(d+1)%4]+" "
				var target [4]int
				for qi:=0;qi<4;qi++{
					idx:=(qi+d)%4
					s+=ns[idx]+" reports "
					target[qi]=len(s)
					s+=latest[idx]+"."
					if qi<3{s+=" "}
				}
				s+="\n"
				ex:=uplm0eExample{text:s,targetPos:target,updateCount:uc}
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

func uplm0eEvaluate(model *uplm0aModel, examples []uplm0eExample, stream4 bool, memory bool, split string) UPLM0EMetric {
	hits,total:=0,0
	depHits,depTotal:=0,0
	exactParagraphs,paragraphs:=0,0
	nll:=0.0
	maxEntries:=0
	exactBy:=map[int]int{0:0,1:0,2:0,4:0}
	totalBy:=map[int]int{0:0,1:0,2:0,4:0}

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
					if value,ok:=recall.values[parser.queryName];ok&&len(value)>0{
						memByte:=value[0]
						pred=model.index[int(memByte)]
						if memByte==targetByte{prob=1}else{prob=1e-12}
					}
					parser.queryNext=false
				}
				if prob<1e-12{prob=1e-12}
				total++;nll-=math.Log(prob);if pred==target{hits++}
				if qi,ok:=uplm0dIsTarget(t+1,examples[e].targetPos);ok{
					depTotal++;querySeen[qi]=true
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
	ce:=nll/float64(total)
	arm:="recurrent64";if memory{arm="recurrent64_syntax_recall16"}
	acc:=func(k int)float64{if totalBy[k]==0{return 0};return float64(exactBy[k])/float64(totalBy[k])}
	return UPLM0EMetric{
		Arm:arm,Split:split,Top1Accuracy:float64(hits)/float64(total),Perplexity:math.Exp(ce),
		DependentFirstByteAccuracy:float64(depHits)/float64(depTotal),
		QuerySetExactAccuracy:float64(exactParagraphs)/float64(paragraphs),
		Update0ExactAccuracy:acc(0),Update1ExactAccuracy:acc(1),Update2ExactAccuracy:acc(2),Update4ExactAccuracy:acc(4),
		MaxRecallEntries:maxEntries,ExactRecallBytes:maxEntries*16,RecurrentStateBytes:512,
	}
}

func RunUPLM0E()(UPLM0EMutableBindingResult,error){
	train,held,alphabet:=uplm0eCorpus()
	model:=newUPLM0AModel(alphabet)
	const epochs=20
	const lr=0.08
	for epoch:=0;epoch<epochs;epoch++{for _,ex:=range train{model.trainSentence(ex.text,lr)}}
	return UPLM0EMutableBindingResult{
		Schema:UPLM0EMutableBindingSchema,Experiment:"UP-LM0E-mutable-binding-language",
		SourceUPLM0DSeal:"c44b28d8feb60320f0ed5c15548b26b2da727c1a",
		StateDimension:64,ExactRecallCap:16,Epochs:epochs,LearningRate:lr,
		TrainParagraphs:len(train),HeldOutParagraphs:len(held),AttentionUsed:false,FutureOracleUsed:false,
		Metrics:[]UPLM0EMetric{
			uplm0eEvaluate(model,train,false,false,"train"),
			uplm0eEvaluate(model,held,false,false,"heldout_recombination"),
			uplm0eEvaluate(model,held,true,false,"heldout_stream4"),
			uplm0eEvaluate(model,train,false,true,"train"),
			uplm0eEvaluate(model,held,false,true,"heldout_recombination"),
			uplm0eEvaluate(model,held,true,true,"heldout_stream4"),
		},
	},nil
}
