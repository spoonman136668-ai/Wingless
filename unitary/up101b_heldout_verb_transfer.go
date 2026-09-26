package unitary

const UP101BHeldoutVerbTransferSchema = "wingless.up101b-heldout-verb-transfer.v1"

type UP101BHeldoutVerbTransferResult struct {
	Schema                          string               `json:"schema"`
	Experiment                      string               `json:"experiment"`
	SourceUP100BSeal                string               `json:"source_up100b_seal"`
	StateDimension                  int                  `json:"state_dimension"`
	Epochs                          int                  `json:"epochs"`
	LearningRate                    float64              `json:"learning_rate"`
	ExactRecallCap                  int                  `json:"exact_recall_cap"`
	ExplicitClassAtLearnedInference bool                 `json:"explicit_class_at_learned_inference"`
	ClassMetrics                    []UP99BClassMetric    `json:"class_metrics"`
	VerbMetrics                     []UP99BVerbMetric     `json:"verb_metrics"`
	Confusions                      []UP99BConfusion      `json:"confusions"`
	RoutingPoints                   []UP100BRoutingPoint  `json:"routing_points"`
}

var up101bTrainVerbs = []int{0,1,3,4,6,7}
var up101bHeldoutVerbs = []int{2,5,8}

func up101bContains(xs []int, v int) bool {
	for _, x := range xs {
		if x == v { return true }
	}
	return false
}

func up101bTrain() *up97bClassifier {
	c := &up97bClassifier{}
	for epoch := 0; epoch < 20; epoch++ {
		for ni := 0; ni < 6; ni++ {
			for _, vi := range up101bTrainVerbs {
				h := up99bPrefixEncoder(ni, vi, 0)
				p := c.probs(h)
				target := up97bVerbClass(vi)
				for k := 0; k < 3; k++ {
					g := p[k]
					if k == target { g -= 1 }
					for i := 0; i < 64; i++ {
						c.w[k][i] -= 0.08 * g * h[i]
					}
					c.b[k] -= 0.08 * g
				}
			}
		}
	}
	return c
}

func up101bEval(c *up97bClassifier, split, label string, class int, verbs []int) UP99BClassMetric {
	total, hits, tp, fp, fn := 0, 0, 0, 0, 0
	for ni := 0; ni < 6; ni++ {
		for _, vi := range verbs {
			target := up97bVerbClass(vi)
			pred := up97bArgmax(c.probs(up99bPrefixEncoder(ni, vi, 0)))
			total++
			if pred == target { hits++ }
			if class >= 0 {
				if pred == class && target == class { tp++ }
				if pred == class && target != class { fp++ }
				if pred != class && target == class { fn++ }
			}
		}
	}
	precision, recall := 1.0, 1.0
	if class >= 0 {
		if tp+fp > 0 { precision = float64(tp)/float64(tp+fp) }
		if tp+fn > 0 { recall = float64(tp)/float64(tp+fn) }
	}
	return UP99BClassMetric{Arm:"prefix_through_verb",Split:split,Label:label,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func up101bVerbEval(c *up97bClassifier, vi int) UP99BVerbMetric {
	total, hits := 0, 0
	for ni := 0; ni < 6; ni++ {
		target := up97bVerbClass(vi)
		pred := up97bArgmax(c.probs(up99bPrefixEncoder(ni, vi, 0)))
		total++
		if pred == target { hits++ }
	}
	return UP99BVerbMetric{Arm:"prefix_through_verb",Verb:up97bVerbs[vi],Accuracy:float64(hits)/float64(total),Examples:total}
}

func up101bHeldoutConfusion(c *up97bClassifier) UP99BConfusion {
	var counts [3][3]int
	for ni := 0; ni < 6; ni++ {
		for _, vi := range up101bHeldoutVerbs {
			target := up97bVerbClass(vi)
			pred := up97bArgmax(c.probs(up99bPrefixEncoder(ni, vi, 0)))
			counts[target][pred]++
		}
	}
	return UP99BConfusion{Arm:"heldout_verbs",Counts:counts}
}

func up101bPredictHeldout(c *up97bClassifier, class, ni int) int {
	vi := 2
	if class == up97bObserve { vi = 5 }
	if class == up97bReport { vi = 8 }
	return up97bArgmax(c.probs(up99bPrefixEncoder(ni, vi, 0)))
}

func up101bRoute(c *up97bClassifier, arm string, totalWrites, targetKeys int, seedBases []int) UP100BRoutingPoint {
	qHits, qTotal, exactHits, episodes := 0, 0, 0, 0
	maxRecall := 0
	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 2001+totalWrites*7+targetKeys, ep)
			rng := newSQ0RNG(seed)
			m := newSQ0Machine("transport_gated_correction", seed)
			truth := make(map[int]int, targetKeys)
			event := func(class, key, value int) {
				ni := key % 6
				if ni < 0 { ni = -ni }
				predClass := class
				if arm == "learned_heldout_verb" {
					predClass = up101bPredictHeldout(c, class, ni)
				}
				switch predClass {
				case up97bStore:
					m.write(key,value,32)
					m.recallWrite(key,value)
				case up97bObserve:
					m.write(key,value,32)
				case up97bReport:
				}
			}
			for k:=0;k<targetKeys;k++ {
				v:=rng.intn(32); truth[k]=v; event(up97bStore,k,v)
			}
			distractors:=totalWrites-2*targetKeys
			pre:=distractors/2; post:=distractors-pre
			for d:=0;d<pre;d++ { event(up97bObserve,1000+ep*1000+d,rng.intn(32)) }
			for k:=0;k<targetKeys;k++ {
				old:=truth[k]; v:=(old+1+rng.intn(31))%32; truth[k]=v; event(up97bStore,k,v)
			}
			for d:=0;d<post;d++ { event(up97bObserve,200000+ep*1000+d,rng.intn(32)) }
			exact:=true
			for k:=0;k<targetKeys;k++ {
				predClass:=up97bReport
				if arm=="learned_heldout_verb" { predClass=up101bPredictHeldout(c,up97bReport,k%6) }
				got,ok:=0,false
				if predClass==up97bReport { got,ok=m.recallQuery(k) }
				qTotal++
				if ok&&got==truth[k] { qHits++ } else { exact=false }
			}
			episodes++
			if exact { exactHits++ }
			if m.recallUsed>maxRecall { maxRecall=m.recallUsed }
		}
	}
	return UP100BRoutingPoint{Arm:arm,TotalWrites:totalWrites,TargetKeys:targetKeys,QueryAccuracy:float64(qHits)/float64(qTotal),ExactTargetSetAccuracy:float64(exactHits)/float64(episodes),RecallEntriesUsed:maxRecall}
}

func RunUP101B()(UP101BHeldoutVerbTransferResult,error){
	c:=up101bTrain()
	result:=UP101BHeldoutVerbTransferResult{
		Schema:UP101BHeldoutVerbTransferSchema,
		Experiment:"UP-101B-heldout-verb-transfer",
		SourceUP100BSeal:"0c91e700b58bc2873e932bfdd7cc2d0260beadf8",
		StateDimension:64,Epochs:20,LearningRate:0.08,ExactRecallCap:16,
		ExplicitClassAtLearnedInference:false,
		ClassMetrics:[]UP99BClassMetric{
			up101bEval(c,"train_verbs","all",-1,up101bTrainVerbs),
			up101bEval(c,"train_verbs","STORE",up97bStore,up101bTrainVerbs),
			up101bEval(c,"train_verbs","OBSERVE",up97bObserve,up101bTrainVerbs),
			up101bEval(c,"train_verbs","REPORT",up97bReport,up101bTrainVerbs),
			up101bEval(c,"heldout_verbs","all",-1,up101bHeldoutVerbs),
			up101bEval(c,"heldout_verbs","STORE",up97bStore,up101bHeldoutVerbs),
			up101bEval(c,"heldout_verbs","OBSERVE",up97bObserve,up101bHeldoutVerbs),
			up101bEval(c,"heldout_verbs","REPORT",up97bReport,up101bHeldoutVerbs),
		},
		Confusions:[]UP99BConfusion{up101bHeldoutConfusion(c)},
	}
	for vi:=0;vi<9;vi++ { result.VerbMetrics=append(result.VerbMetrics,up101bVerbEval(c,vi)) }
	seedBases:=[]int{159000000,160000000}
	for _,arm:=range []string{"explicit_event_class","learned_heldout_verb"} {
		for _,writes:=range []int{32,64,128,256} {
			for _,targets:=range []int{4,8,16} {
				result.RoutingPoints=append(result.RoutingPoints,up101bRoute(c,arm,writes,targets,seedBases))
			}
		}
	}
	return result,nil
}
