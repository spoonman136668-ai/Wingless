package unitary

import "math"

const UP95BLocalEventClassifierSchema = "wingless.up95b-local-event-classifier.v1"

type UP95BClassifierMetric struct {
	Split     string  `json:"split"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP95BRoutingPoint struct {
	Arm                     string  `json:"arm"`
	TotalWrites             int     `json:"total_writes"`
	TargetKeys              int     `json:"target_keys"`
	QueryAccuracy           float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy  float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	AdmissionPrecision      float64 `json:"admission_precision"`
	AdmissionRecall         float64 `json:"admission_recall"`
	FalsePositiveAdmissions int     `json:"false_positive_admissions"`
}

type UP95BLocalEventClassifierResult struct {
	Schema              string                  `json:"schema"`
	Experiment          string                  `json:"experiment"`
	SourceUP94BSeal     string                  `json:"source_up94b_seal"`
	StateDimension      int                     `json:"state_dimension"`
	ClassifierEpochs    int                     `json:"classifier_epochs"`
	LearningRate        float64                 `json:"learning_rate"`
	DecisionThreshold   float64                 `json:"decision_threshold"`
	ExactRecallCap      int                     `json:"exact_recall_cap"`
	RecurrentStateBytes int                     `json:"recurrent_state_bytes"`
	ExplicitTypeAtLearnedInference bool         `json:"explicit_type_at_learned_inference"`
	ClassifierMetrics   []UP95BClassifierMetric `json:"classifier_metrics"`
	RoutingPoints       []UP95BRoutingPoint     `json:"routing_points"`
}

func up95bEmbed(b byte,i int) float64 {
	x:=uint64(b+1)*0x9e3779b97f4a7c15^uint64(i+1)*0xbf58476d1ce4e5b9^0x510e527fade682d1
	x=sq0Mix64(x)
	if x&1==0{return -0.125}
	return 0.125
}

func up95bStep(h [64]float64,b byte)[64]float64{
	var next [64]float64
	for i:=0;i<64;i++{
		src:=(13*i+7)&63
		sign:=1.0;if i&1==1{sign=-1}
		next[i]=math.Tanh(0.90*sign*h[src]+0.35*up95bEmbed(b,i))
	}
	return next
}

func up95bEncode(s string)[64]float64{
	var h [64]float64
	for i:=0;i<len(s);i++{h=up95bStep(h,s[i])}
	return h
}

func up95bTwo(n int) string {
	return string([]byte{'0'+byte((n/10)%10),'0'+byte(n%10)})
}

func up95bSurface(store bool,key,value int) string {
	verb:=" observes "
	if store{verb=" stores "}
	return "k"+up95bTwo(key&31)+verb+"v"+up95bTwo(value&31)+"."
}

type up95bClassifier struct {
	w [64]float64
	b float64
}

func up95bSigmoid(x float64) float64 {
	if x>=0 {z:=math.Exp(-x);return 1/(1+z)}
	z:=math.Exp(x);return z/(1+z)
}

func (c *up95bClassifier) prob(h [64]float64) float64 {
	s:=c.b
	for i:=0;i<64;i++{s+=c.w[i]*h[i]}
	return up95bSigmoid(s)
}

func (c *up95bClassifier) train(){
	const lr=0.08
	for epoch:=0;epoch<20;epoch++{
		for key:=0;key<24;key++{
			for value:=0;value<32;value++{
				for _,store:=range []bool{false,true}{
					h:=up95bEncode(up95bSurface(store,key,value))
					y:=0.0;if store{y=1}
					g:=c.prob(h)-y
					for i:=0;i<64;i++{c.w[i]-=lr*g*h[i]}
					c.b-=lr*g
				}
			}
		}
	}
}

func up95bClassifierEval(c *up95bClassifier,keys []int,split string) UP95BClassifierMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	for _,key:=range keys{
		for value:=0;value<32;value++{
			for _,store:=range []bool{false,true}{
				pred:=c.prob(up95bEncode(up95bSurface(store,key,value)))>=0.5
				total++;if pred==store{hits++}
				if pred&&store{tp++};if pred&&!store{fp++};if !pred&&store{fn++}
			}
		}
	}
	precision:=1.0;if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	recall:=1.0;if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	return UP95BClassifierMetric{Split:split,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func up95bRouting(c *up95bClassifier,arm string,totalWrites,targetKeys int,seedBases []int) UP95BRoutingPoint {
	qHits,qTotal,exactHits,episodes:=0,0,0,0
	maxRecall:=0
	admissions,truePos,storeEvents,falsePos:=0,0,0,0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1001+totalWrites*7+targetKeys,ep)
			rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed)
			truth:=make(map[int]int,targetKeys)

			write:=func(store bool,key,value int){
				m.write(key,value,32)
				if store{storeEvents++}
				admit:=false
				switch arm{
				case "fifo_all_writes":
					admit=true
				case "explicit_store_type":
					admit=store
				case "learned_local_classifier":
					admit=c.prob(up95bEncode(up95bSurface(store,key,value)))>=0.5
				}
				if admit{
					m.recallWrite(key,value);admissions++
					if store{truePos++}else{falsePos++}
				}
			}

			for k:=0;k<targetKeys;k++{v:=rng.intn(32);truth[k]=v;write(true,k,v)}
			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2;post:=distractors-pre
			for d:=0;d<pre;d++{write(false,1000+ep*1000+d,rng.intn(32))}
			for k:=0;k<targetKeys;k++{old:=truth[k];v:=(old+1+rng.intn(31))%32;truth[k]=v;write(true,k,v)}
			for d:=0;d<post;d++{write(false,200000+ep*1000+d,rng.intn(32))}

			exact:=true
			for k:=0;k<targetKeys;k++{got:=up91bQuery(m,k,32);qTotal++;if got==truth[k]{qHits++}else{exact=false}}
			episodes++;if exact{exactHits++};if m.recallUsed>maxRecall{maxRecall=m.recallUsed}
		}
	}
	precision:=1.0;if admissions>0{precision=float64(truePos)/float64(admissions)}
	recall:=0.0;if storeEvents>0{recall=float64(truePos)/float64(storeEvents)}
	return UP95BRoutingPoint{
		Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,
		QueryAccuracy:float64(qHits)/float64(qTotal),ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),
		RecallEntriesUsed:maxRecall,AdmissionPrecision:precision,AdmissionRecall:recall,FalsePositiveAdmissions:falsePos,
	}
}

func RunUP95B()(UP95BLocalEventClassifierResult,error){
	var c up95bClassifier;c.train()
	trainKeys:=make([]int,24);for i:=range trainKeys{trainKeys[i]=i}
	heldKeys:=make([]int,8);for i:=range heldKeys{heldKeys[i]=i+24}
	seedBases:=[]int{139000000,140000000}
	result:=UP95BLocalEventClassifierResult{
		Schema:UP95BLocalEventClassifierSchema,Experiment:"UP-95B-local-event-classifier",
		SourceUP94BSeal:"83fc5d7dac7b4e6bebbec8c93aee0e6d73a111b6",
		StateDimension:64,ClassifierEpochs:20,LearningRate:0.08,DecisionThreshold:0.5,
		ExactRecallCap:16,RecurrentStateBytes:512,ExplicitTypeAtLearnedInference:false,
		ClassifierMetrics:[]UP95BClassifierMetric{
			up95bClassifierEval(&c,trainKeys,"train_keys"),
			up95bClassifierEval(&c,heldKeys,"heldout_keys"),
		},
	}
	for _,arm:=range []string{"fifo_all_writes","explicit_store_type","learned_local_classifier"}{
		for _,writes:=range []int{32,64,128,256}{
			for _,targets:=range []int{4,8,16}{
				result.RoutingPoints=append(result.RoutingPoints,up95bRouting(&c,arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
