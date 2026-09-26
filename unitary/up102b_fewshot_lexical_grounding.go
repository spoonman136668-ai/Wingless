package unitary

const UP102BFewshotGroundingSchema = "wingless.up102b-fewshot-lexical-grounding.v1"

type UP102BGroundingPoint struct {
	ExamplesPerVerb int     `json:"examples_per_verb"`
	HeldoutAccuracy float64 `json:"heldout_accuracy"`
	StorePrecision  float64 `json:"store_precision"`
	StoreRecall     float64 `json:"store_recall"`
	ObservePrecision float64 `json:"observe_precision"`
	ObserveRecall    float64 `json:"observe_recall"`
	ReportPrecision  float64 `json:"report_precision"`
	ReportRecall     float64 `json:"report_recall"`
	SeenVerbAccuracy float64 `json:"seen_verb_accuracy"`
}

type UP102BFewshotGroundingResult struct {
	Schema          string                 `json:"schema"`
	Experiment      string                 `json:"experiment"`
	SourceUP101BSeal string                `json:"source_up101b_seal"`
	StateDimension  int                    `json:"state_dimension"`
	Epochs          int                    `json:"epochs"`
	LearningRate    float64                `json:"learning_rate"`
	Points          []UP102BGroundingPoint `json:"points"`
}

func up102bAdapt(base *up97bClassifier, examplesPerVerb int) *up97bClassifier {
	copyModel := *base
	c := &copyModel
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<examplesPerVerb;ni++ {
			for _,vi:=range up101bHeldoutVerbs {
				h:=up99bPrefixEncoder(ni,vi,0)
				p:=c.probs(h)
				target:=up97bVerbClass(vi)
				for k:=0;k<3;k++ {
					g:=p[k]
					if k==target { g-=1 }
					for i:=0;i<64;i++ { c.w[k][i]-=0.08*g*h[i] }
					c.b[k]-=0.08*g
				}
			}
		}
	}
	return c
}

func up102bHeldoutEval(c *up97bClassifier)(accuracy float64, precision [3]float64, recall [3]float64){
	total,hits:=0,0
	var tp,fp,fn [3]int
	for ni:=4;ni<6;ni++ {
		for _,vi:=range up101bHeldoutVerbs {
			target:=up97bVerbClass(vi)
			pred:=up97bArgmax(c.probs(up99bPrefixEncoder(ni,vi,0)))
			total++
			if pred==target { hits++ }
			for class:=0;class<3;class++ {
				if pred==class&&target==class { tp[class]++ }
				if pred==class&&target!=class { fp[class]++ }
				if pred!=class&&target==class { fn[class]++ }
			}
		}
	}
	for class:=0;class<3;class++ {
		precision[class]=1
		recall[class]=1
		if tp[class]+fp[class]>0 { precision[class]=float64(tp[class])/float64(tp[class]+fp[class]) }
		if tp[class]+fn[class]>0 { recall[class]=float64(tp[class])/float64(tp[class]+fn[class]) }
	}
	return float64(hits)/float64(total),precision,recall
}

func up102bSeenAccuracy(c *up97bClassifier) float64 {
	total,hits:=0,0
	for ni:=0;ni<6;ni++ {
		for _,vi:=range up101bTrainVerbs {
			target:=up97bVerbClass(vi)
			pred:=up97bArgmax(c.probs(up99bPrefixEncoder(ni,vi,0)))
			total++
			if pred==target { hits++ }
		}
	}
	return float64(hits)/float64(total)
}

func RunUP102B()(UP102BFewshotGroundingResult,error){
	base:=up101bTrain()
	result:=UP102BFewshotGroundingResult{
		Schema:UP102BFewshotGroundingSchema,
		Experiment:"UP-102B-fewshot-lexical-grounding",
		SourceUP101BSeal:"397840ebed9f170741980ee2d4c2f393bbe3a748",
		StateDimension:64,Epochs:20,LearningRate:0.08,
	}
	for _,k:=range []int{0,1,2,4} {
		c:=up102bAdapt(base,k)
		acc,p,r:=up102bHeldoutEval(c)
		result.Points=append(result.Points,UP102BGroundingPoint{
			ExamplesPerVerb:k,HeldoutAccuracy:acc,
			StorePrecision:p[up97bStore],StoreRecall:r[up97bStore],
			ObservePrecision:p[up97bObserve],ObserveRecall:r[up97bObserve],
			ReportPrecision:p[up97bReport],ReportRecall:r[up97bReport],
			SeenVerbAccuracy:up102bSeenAccuracy(c),
		})
	}
	return result,nil
}
