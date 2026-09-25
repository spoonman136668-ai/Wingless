package unitary

import "math"

const UP97BThreeWayNaturalEventSchema = "wingless.up97b-threeway-natural-event.v1"

const (
	up97bStore = 0
	up97bObserve = 1
	up97bReport = 2
)

var up97bNames=[]string{"ada","ben","cy","dee","eli","fay"}
var up97bValues=[]string{"amber","cobalt","ivory","jade","mauve","silver"}
var up97bVerbs=[]string{"stores","keeps","holds","observes","sees","notes","reports","recalls","tells"}

type UP97BClassMetric struct {
	Split string `json:"split"`
	Label string `json:"label"`
	Accuracy float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Examples int `json:"examples"`
}

type UP97BVerbMetric struct {
	Verb string `json:"verb"`
	Accuracy float64 `json:"accuracy"`
	Examples int `json:"examples"`
}

type UP97BRoutingPoint struct {
	Arm string `json:"arm"`
	TotalWrites int `json:"total_writes"`
	TargetKeys int `json:"target_keys"`
	QueryAccuracy float64 `json:"query_accuracy"`
	ExactTargetSetAccuracy float64 `json:"exact_target_set_accuracy"`
	RecallEntriesUsed int `json:"recall_entries_used"`
}

type UP97BThreeWayNaturalEventResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP96BSeal string `json:"source_up96b_seal"`
	StateDimension int `json:"state_dimension"`
	Epochs int `json:"epochs"`
	LearningRate float64 `json:"learning_rate"`
	ExplicitClassAtLearnedInference bool `json:"explicit_class_at_learned_inference"`
	ClassMetrics []UP97BClassMetric `json:"class_metrics"`
	VerbMetrics []UP97BVerbMetric `json:"verb_metrics"`
	RoutingPoints []UP97BRoutingPoint `json:"routing_points"`
}

type up97bClassifier struct { w [3][64]float64; b [3]float64 }

func up97bVerbClass(vi int) int {
	if vi<3{return up97bStore}
	if vi<6{return up97bObserve}
	return up97bReport
}
func up97bSurface(ni,vi,vali int) string {
	return up97bNames[ni]+" "+up97bVerbs[vi]+" "+up97bValues[vali]+"."
}
func up97bTrainSplit(ni,vali,vi int) bool { return (ni+2*vali+vi)%4!=3 }

func (c *up97bClassifier) probs(h [64]float64) [3]float64 {
	var logits [3]float64
	maxV:=math.Inf(-1)
	for k:=0;k<3;k++{
		s:=c.b[k]
		for i:=0;i<64;i++{s+=c.w[k][i]*h[i]}
		logits[k]=s;if s>maxV{maxV=s}
	}
	sum:=0.0
	for k:=0;k<3;k++{logits[k]=math.Exp(logits[k]-maxV);sum+=logits[k]}
	for k:=0;k<3;k++{logits[k]/=sum}
	return logits
}
func up97bArgmax(p [3]float64) int {best:=0;for k:=1;k<3;k++{if p[k]>p[best]{best=k}};return best}

func up97bTrainClassifier()*up97bClassifier{
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<6;ni++{
			for vali:=0;vali<6;vali++{
				for vi:=0;vi<9;vi++{
					if !up97bTrainSplit(ni,vali,vi){continue}
					h:=up95bEncode(up97bSurface(ni,vi,vali))
					p:=c.probs(h);target:=up97bVerbClass(vi)
					for k:=0;k<3;k++{
						g:=p[k];if k==target{g-=1}
						for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
						c.b[k]-=0.08*g
					}
				}
			}
		}
	}
	return c
}

func up97bClassEval(c *up97bClassifier, split string, class int, label string) UP97BClassMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			for vi:=0;vi<9;vi++{
				isTrain:=up97bTrainSplit(ni,vali,vi)
				if split=="train"&&!isTrain{continue}
				if split=="heldout_recombination"&&isTrain{continue}
				target:=up97bVerbClass(vi)
				pred:=up97bArgmax(c.probs(up95bEncode(up97bSurface(ni,vi,vali))))
				total++;if pred==target{hits++}
				if class>=0{
					if pred==class&&target==class{tp++}
					if pred==class&&target!=class{fp++}
					if pred!=class&&target==class{fn++}
				}
			}
		}
	}
	precision:=1.0;recall:=1.0
	if class>=0{
		if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
		if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	}
	return UP97BClassMetric{Split:split,Label:label,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}
func up97bVerbEval(c *up97bClassifier, vi int) UP97BVerbMetric {
	total,hits:=0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			if up97bTrainSplit(ni,vali,vi){continue}
			target:=up97bVerbClass(vi)
			pred:=up97bArgmax(c.probs(up95bEncode(up97bSurface(ni,vi,vali))))
			total++;if pred==target{hits++}
		}
	}
	return UP97BVerbMetric{Verb:up97bVerbs[vi],Accuracy:float64(hits)/float64(total),Examples:total}
}

func up97bRoute(c *up97bClassifier, arm string,totalWrites,targetKeys int,seedBases []int) UP97BRoutingPoint {
	qHits,qTotal,exactHits,episodes:=0,0,0,0
	maxRecall:=0
	for _,base:=range seedBases{
		for ep:=0;ep<64;ep++{
			seed:=sq0Seed(base,1401+totalWrites*7+targetKeys,ep);rng:=newSQ0RNG(seed)
			m:=newSQ0Machine("transport_gated_correction",seed);truth:=make(map[int]int,targetKeys)
			event:=func(class,key,value int){
				ni:=key%6;if ni<0{ni=-ni}
				vali:=value%6;if vali<0{vali=-vali}
				vi:=0
				switch class{
				case up97bStore:vi=(key+value+ep)%3
				case up97bObserve:vi=3+(key+value+ep)%3
				case up97bReport:vi=6+(key+value+ep)%3
				}
				predClass:=class
				if arm=="learned_threeway"{
					predClass=up97bArgmax(c.probs(up95bEncode(up97bSurface(ni,vi,vali))))
				}
				switch predClass{
				case up97bStore:
					m.write(key,value,32);m.recallWrite(key,value)
				case up97bObserve:
					m.write(key,value,32)
				case up97bReport:
					// report carries no write side effect
				}
			}
			for k:=0;k<targetKeys;k++{v:=rng.intn(32);truth[k]=v;event(up97bStore,k,v)}
			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2;post:=distractors-pre
			for d:=0;d<pre;d++{event(up97bObserve,1000+ep*1000+d,rng.intn(32))}
			for k:=0;k<targetKeys;k++{old:=truth[k];v:=(old+1+rng.intn(31))%32;truth[k]=v;event(up97bStore,k,v)}
			for d:=0;d<post;d++{event(up97bObserve,200000+ep*1000+d,rng.intn(32))}
			exact:=true
			for k:=0;k<targetKeys;k++{
				// Classify a report surface using dummy value 0, then query only if classified REPORT.
				ni:=k%6;vi:=6+(k+ep)%3;predClass:=up97bReport
				if arm=="learned_threeway"{predClass=up97bArgmax(c.probs(up95bEncode(up97bSurface(ni,vi,0))))}
				got,ok:=0,false
				if predClass==up97bReport{got,ok=m.recallQuery(k)}
				qTotal++;if ok&&got==truth[k]{qHits++}else{exact=false}
			}
			episodes++;if exact{exactHits++};if m.recallUsed>maxRecall{maxRecall=m.recallUsed}
		}
	}
	return UP97BRoutingPoint{Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,QueryAccuracy:float64(qHits)/float64(qTotal),ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),RecallEntriesUsed:maxRecall}
}

func RunUP97B()(UP97BThreeWayNaturalEventResult,error){
	c:=up97bTrainClassifier()
	result:=UP97BThreeWayNaturalEventResult{
		Schema:UP97BThreeWayNaturalEventSchema,Experiment:"UP-97B-threeway-natural-event",
		SourceUP96BSeal:"87ca9afbb789733ca50cf39ee4ed730938cac9fd",
		StateDimension:64,Epochs:20,LearningRate:0.08,ExplicitClassAtLearnedInference:false,
		ClassMetrics:[]UP97BClassMetric{
			up97bClassEval(c,"train",-1,"all"),
			up97bClassEval(c,"heldout_recombination",-1,"all"),
			up97bClassEval(c,"heldout_recombination",up97bStore,"STORE"),
			up97bClassEval(c,"heldout_recombination",up97bObserve,"OBSERVE"),
			up97bClassEval(c,"heldout_recombination",up97bReport,"REPORT"),
		},
	}
	for vi:=0;vi<9;vi++{result.VerbMetrics=append(result.VerbMetrics,up97bVerbEval(c,vi))}
	seedBases:=[]int{147000000,148000000}
	for _,arm:=range []string{"explicit_event_class","learned_threeway"}{
		for _,writes:=range []int{32,64,128,256}{
			for _,targets:=range []int{4,8,16}{
				result.RoutingPoints=append(result.RoutingPoints,up97bRoute(c,arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
