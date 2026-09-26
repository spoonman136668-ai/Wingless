package unitary

const UP99BVerbWindowSchema = "wingless.up99b-verb-window-separability.v1"

type UP99BClassMetric struct {
	Arm string `json:"arm"`
	Split string `json:"split"`
	Label string `json:"label"`
	Accuracy float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Examples int `json:"examples"`
}

type UP99BVerbMetric struct {
	Arm string `json:"arm"`
	Verb string `json:"verb"`
	Accuracy float64 `json:"accuracy"`
	Examples int `json:"examples"`
}

type UP99BConfusion struct {
	Arm string `json:"arm"`
	Counts [3][3]int `json:"counts"`
}

type UP99BVerbWindowResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP98BSeal string `json:"source_up98b_seal"`
	StateDimension int `json:"state_dimension"`
	Epochs int `json:"epochs"`
	LearningRate float64 `json:"learning_rate"`
	ExplicitClassAtInference bool `json:"explicit_class_at_inference"`
	ClassMetrics []UP99BClassMetric `json:"class_metrics"`
	VerbMetrics []UP99BVerbMetric `json:"verb_metrics"`
	Confusions []UP99BConfusion `json:"confusions"`
}

type up99bEncoder func(ni,vi,vali int)[64]float64

func up99bFullEncoder(ni,vi,vali int)[64]float64 {
	return up95bEncode(up97bSurface(ni,vi,vali))
}

func up99bPrefixEncoder(ni,vi,vali int)[64]float64 {
	_ = vali
	return up95bEncode(up97bNames[ni]+" "+up97bVerbs[vi])
}

func up99bTrain(enc up99bEncoder)*up97bClassifier {
	c:=&up97bClassifier{}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<6;ni++{
			for vali:=0;vali<6;vali++{
				for vi:=0;vi<9;vi++{
					if !up97bTrainSplit(ni,vali,vi){continue}
					h:=enc(ni,vi,vali)
					p:=c.probs(h)
					target:=up97bVerbClass(vi)
					for k:=0;k<3;k++{
						g:=p[k]
						if k==target{g-=1}
						for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
						c.b[k]-=0.08*g
					}
				}
			}
		}
	}
	return c
}

func up99bEval(arm,split,label string,class int,c *up97bClassifier,enc up99bEncoder) UP99BClassMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			for vi:=0;vi<9;vi++{
				isTrain:=up97bTrainSplit(ni,vali,vi)
				if split=="train"&&!isTrain{continue}
				if split=="heldout_recombination"&&isTrain{continue}
				target:=up97bVerbClass(vi)
				pred:=up97bArgmax(c.probs(enc(ni,vi,vali)))
				total++
				if pred==target{hits++}
				if class>=0{
					if pred==class&&target==class{tp++}
					if pred==class&&target!=class{fp++}
					if pred!=class&&target==class{fn++}
				}
			}
		}
	}
	precision,recall:=1.0,1.0
	if class>=0{
		if tp+fp>0{precision=float64(tp)/float64(tp+fp)}
		if tp+fn>0{recall=float64(tp)/float64(tp+fn)}
	}
	return UP99BClassMetric{Arm:arm,Split:split,Label:label,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func up99bVerbEval(arm string,vi int,c *up97bClassifier,enc up99bEncoder) UP99BVerbMetric {
	total,hits:=0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			if up97bTrainSplit(ni,vali,vi){continue}
			target:=up97bVerbClass(vi)
			pred:=up97bArgmax(c.probs(enc(ni,vi,vali)))
			total++
			if pred==target{hits++}
		}
	}
	return UP99BVerbMetric{Arm:arm,Verb:up97bVerbs[vi],Accuracy:float64(hits)/float64(total),Examples:total}
}

func up99bConfusion(arm string,c *up97bClassifier,enc up99bEncoder) UP99BConfusion {
	var counts [3][3]int
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			for vi:=0;vi<9;vi++{
				if up97bTrainSplit(ni,vali,vi){continue}
				target:=up97bVerbClass(vi)
				pred:=up97bArgmax(c.probs(enc(ni,vi,vali)))
				counts[target][pred]++
			}
		}
	}
	return UP99BConfusion{Arm:arm,Counts:counts}
}

func RunUP99B()(UP99BVerbWindowResult,error){
	result:=UP99BVerbWindowResult{
		Schema:UP99BVerbWindowSchema,Experiment:"UP-99B-verb-window-separability",
		SourceUP98BSeal:"dde36e56943a4e7af2c39d0aa33c7f6b9a9bb23b",
		StateDimension:64,Epochs:20,LearningRate:0.08,ExplicitClassAtInference:false,
	}
	for _,arm:=range []struct{name string; enc up99bEncoder}{
		{"full_clause",up99bFullEncoder},
		{"prefix_through_verb",up99bPrefixEncoder},
	}{
		c:=up99bTrain(arm.enc)
		result.ClassMetrics=append(result.ClassMetrics,
			up99bEval(arm.name,"train","all",-1,c,arm.enc),
			up99bEval(arm.name,"heldout_recombination","all",-1,c,arm.enc),
			up99bEval(arm.name,"heldout_recombination","STORE",up97bStore,c,arm.enc),
			up99bEval(arm.name,"heldout_recombination","OBSERVE",up97bObserve,c,arm.enc),
			up99bEval(arm.name,"heldout_recombination","REPORT",up97bReport,c,arm.enc),
		)
		for vi:=0;vi<9;vi++{result.VerbMetrics=append(result.VerbMetrics,up99bVerbEval(arm.name,vi,c,arm.enc))}
		result.Confusions=append(result.Confusions,up99bConfusion(arm.name,c,arm.enc))
	}
	return result,nil
}
