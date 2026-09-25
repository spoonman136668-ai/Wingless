package unitary

const UP96BSurfaceGeneralizationSchema = "wingless.up96b-surface-generalization.v1"

var up96bStoreVerbs=[]string{"stores","keeps","holds"}
var up96bObserveVerbs=[]string{"observes","sees","notes"}

type UP96BClassifierMetric struct {
	Split     string  `json:"split"`
	Verb      string  `json:"verb"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	Examples  int     `json:"examples"`
}

type UP96BRoutingPoint struct {
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

type UP96BSurfaceGeneralizationResult struct {
	Schema                         string                  `json:"schema"`
	Experiment                     string                  `json:"experiment"`
	SourceUP95BSeal                string                  `json:"source_up95b_seal"`
	StateDimension                 int                     `json:"state_dimension"`
	ClassifierEpochs               int                     `json:"classifier_epochs"`
	LearningRate                   float64                 `json:"learning_rate"`
	DecisionThreshold              float64                 `json:"decision_threshold"`
	ExactRecallCap                 int                     `json:"exact_recall_cap"`
	ExplicitClassAtLearnedInference bool                   `json:"explicit_class_at_learned_inference"`
	ClassifierMetrics              []UP96BClassifierMetric `json:"classifier_metrics"`
	RoutingPoints                  []UP96BRoutingPoint     `json:"routing_points"`
}

func up96bSurface(key,value,verbIndex int) string {
	verb:=""
	if verbIndex<3{verb=up96bStoreVerbs[verbIndex]}else{verb=up96bObserveVerbs[verbIndex-3]}
	return "k"+up95bTwo(key&31)+" "+verb+" v"+up95bTwo(value&31)+"."
}
func up96bIsStoreVerb(verbIndex int) bool { return verbIndex<3 }
func up96bIsTrain(key,value,verbIndex int) bool { return (key+2*value+verbIndex)%4!=3 }

func up96bTrainClassifier() *up95bClassifier {
	c:=&up95bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for key:=0;key<32;key++{
			for value:=0;value<32;value++{
				for vi:=0;vi<6;vi++{
					if !up96bIsTrain(key,value,vi){continue}
					h:=up95bEncode(up96bSurface(key,value,vi))
					y:=0.0;if up96bIsStoreVerb(vi){y=1}
					g:=c.prob(h)-y
					for i:=0;i<64;i++{c.w[i]-=0.08*g*h[i]}
					c.b-=0.08*g
				}
			}
		}
	}
	return c
}

func up96bEval(c *up95bClassifier,split,verb string,verbIndex int) UP96BClassifierMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	for key:=0;key<32;key++{
		for value:=0;value<32;value++{
			start,end:=0,6
			if verbIndex>=0{start=verbIndex;end=verbIndex+1}
			for vi:=start;vi<end;vi++{
				isTrain:=up96bIsTrain(key,value,vi)
				if split=="train" && !isTrain{continue}
				if split=="heldout_recombination" && isTrain{continue}
				want:=up96bIsStoreVerb(vi)
				pred:=c.prob(up95bEncode(up96bSurface(key,value,vi)))>=0.5
				total++;if pred==want{hits++}
				if pred&&want{tp++};if pred&&!want{fp++};if !pred&&want{fn++}
			}
		}
	}
	precision:=1.0;if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
	recall:=1.0;if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	return UP96BClassifierMetric{Split:split,Verb:verb,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func up96bRoute(c *up95bClassifier,arm string,totalWrites,targetKeys int,seedBases []int) UP96BRoutingPoint {
	qHits,qTotal,exactHits,episodes:=0,0,0,0
	maxRecall:=0
	admissions,truePos,storeEvents,falsePos:=0,0,0,0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1201+totalWrites*7+targetKeys,ep)
			rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed)
			truth:=make(map[int]int,targetKeys)
			write:=func(store bool,key,value int){
				m.write(key,value,32)
				if store{storeEvents++}
				verbIndex:=0
				if store{verbIndex=(key+value+ep)%3}else{verbIndex=3+(key+value+ep)%3}
				admit:=store
				if arm=="learned_surface_classifier"{
					admit=c.prob(up95bEncode(up96bSurface(key,value,verbIndex)))>=0.5
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
	return UP96BRoutingPoint{
		Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,QueryAccuracy:float64(qHits)/float64(qTotal),
		ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),RecallEntriesUsed:maxRecall,
		AdmissionPrecision:precision,AdmissionRecall:recall,FalsePositiveAdmissions:falsePos,
	}
}

func RunUP96B()(UP96BSurfaceGeneralizationResult,error){
	c:=up96bTrainClassifier()
	result:=UP96BSurfaceGeneralizationResult{
		Schema:UP96BSurfaceGeneralizationSchema,Experiment:"UP-96B-surface-generalization",
		SourceUP95BSeal:"f15ff278753d3af1786b38f9958e41ded3c00fc5",StateDimension:64,
		ClassifierEpochs:20,LearningRate:0.08,DecisionThreshold:0.5,ExactRecallCap:16,ExplicitClassAtLearnedInference:false,
	}
	result.ClassifierMetrics=append(result.ClassifierMetrics,
		up96bEval(c,"train","all",-1),
		up96bEval(c,"heldout_recombination","all",-1),
	)
	for vi:=0;vi<6;vi++{
		verb:=up96bStoreVerbs[0]
		if vi<3{verb=up96bStoreVerbs[vi]}else{verb=up96bObserveVerbs[vi-3]}
		result.ClassifierMetrics=append(result.ClassifierMetrics,up96bEval(c,"heldout_recombination",verb,vi))
	}
	seedBases:=[]int{143000000,144000000}
	for _,arm:=range []string{"explicit_event_class","learned_surface_classifier"}{
		for _,writes:=range []int{32,64,128,256}{
			for _,targets:=range []int{4,8,16}{
				result.RoutingPoints=append(result.RoutingPoints,up96bRoute(c,arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
