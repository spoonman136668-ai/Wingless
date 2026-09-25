package unitary

import "math"

const UP98BDecomposedThreeWaySchema = "wingless.up98b-decomposed-threeway.v1"

type UP98BClassMetric struct {
	Arm string `json:"arm"`
	Split string `json:"split"`
	Label string `json:"label"`
	Accuracy float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall float64 `json:"recall"`
	Examples int `json:"examples"`
}

type UP98BVerbMetric struct {
	Arm string `json:"arm"`
	Verb string `json:"verb"`
	Accuracy float64 `json:"accuracy"`
	Examples int `json:"examples"`
}

type UP98BConfusion struct {
	Arm string `json:"arm"`
	Counts [3][3]int `json:"counts"`
}

type UP98BDecomposedThreeWayResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP97BSeal string `json:"source_up97b_seal"`
	StateDimension int `json:"state_dimension"`
	Epochs int `json:"epochs"`
	LearningRate float64 `json:"learning_rate"`
	ThresholdSearchUsed bool `json:"threshold_search_used"`
	ClassWeightsUsed bool `json:"class_weights_used"`
	ClassMetrics []UP98BClassMetric `json:"class_metrics"`
	VerbMetrics []UP98BVerbMetric `json:"verb_metrics"`
	Confusions []UP98BConfusion `json:"confusions"`
}

type up98bOVR struct {
	w [3][64]float64
	b [3]float64
}

func up98bSigmoid(x float64) float64 {
	if x >= 0 {
		z := math.Exp(-x)
		return 1 / (1 + z)
	}
	z := math.Exp(x)
	return z / (1 + z)
}

func (c *up98bOVR) probs(h [64]float64) [3]float64 {
	var p [3]float64
	for k:=0;k<3;k++{
		s:=c.b[k]
		for i:=0;i<64;i++{s+=c.w[k][i]*h[i]}
		p[k]=up98bSigmoid(s)
	}
	return p
}

func up98bArgmax(p [3]float64) int {
	best:=0
	for k:=1;k<3;k++{if p[k]>p[best]{best=k}}
	return best
}

func up98bTrainOVR() *up98bOVR {
	c:=&up98bOVR{}
	for epoch:=0;epoch<20;epoch++{
		for ni:=0;ni<6;ni++{
			for vali:=0;vali<6;vali++{
				for vi:=0;vi<9;vi++{
					if !up97bTrainSplit(ni,vali,vi){continue}
					h:=up95bEncode(up97bSurface(ni,vi,vali))
					target:=up97bVerbClass(vi)
					p:=c.probs(h)
					for k:=0;k<3;k++{
						y:=0.0
						if k==target{y=1}
						g:=p[k]-y
						for i:=0;i<64;i++{c.w[k][i]-=0.08*g*h[i]}
						c.b[k]-=0.08*g
					}
				}
			}
		}
	}
	return c
}

func up98bEval(arm,split,label string,class int,predict func([64]float64)int) UP98BClassMetric {
	total,hits,tp,fp,fn:=0,0,0,0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			for vi:=0;vi<9;vi++{
				isTrain:=up97bTrainSplit(ni,vali,vi)
				if split=="train"&&!isTrain{continue}
				if split=="heldout_recombination"&&isTrain{continue}
				target:=up97bVerbClass(vi)
				pred:=predict(up95bEncode(up97bSurface(ni,vi,vali)))
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
	return UP98BClassMetric{Arm:arm,Split:split,Label:label,Accuracy:float64(hits)/float64(total),Precision:precision,Recall:recall,Examples:total}
}

func up98bVerbEval(arm string,vi int,predict func([64]float64)int) UP98BVerbMetric {
	total,hits:=0,0
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			if up97bTrainSplit(ni,vali,vi){continue}
			target:=up97bVerbClass(vi)
			pred:=predict(up95bEncode(up97bSurface(ni,vi,vali)))
			total++
			if pred==target{hits++}
		}
	}
	return UP98BVerbMetric{Arm:arm,Verb:up97bVerbs[vi],Accuracy:float64(hits)/float64(total),Examples:total}
}

func up98bConfusion(arm string,predict func([64]float64)int) UP98BConfusion {
	var counts [3][3]int
	for ni:=0;ni<6;ni++{
		for vali:=0;vali<6;vali++{
			for vi:=0;vi<9;vi++{
				if up97bTrainSplit(ni,vali,vi){continue}
				target:=up97bVerbClass(vi)
				pred:=predict(up95bEncode(up97bSurface(ni,vi,vali)))
				counts[target][pred]++
			}
		}
	}
	return UP98BConfusion{Arm:arm,Counts:counts}
}

func RunUP98B()(UP98BDecomposedThreeWayResult,error){
	base:=up97bTrainClassifier()
	ovr:=up98bTrainOVR()
	basePredict:=func(h [64]float64)int{return up97bArgmax(base.probs(h))}
	ovrPredict:=func(h [64]float64)int{return up98bArgmax(ovr.probs(h))}
	result:=UP98BDecomposedThreeWayResult{
		Schema:UP98BDecomposedThreeWaySchema,
		Experiment:"UP-98B-decomposed-threeway",
		SourceUP97BSeal:"aca2c8ab1c465c3978f5bcb5078c663e789f4d02",
		StateDimension:64,Epochs:20,LearningRate:0.08,
		ThresholdSearchUsed:false,ClassWeightsUsed:false,
	}
	for _,arm:=range []struct{name string; predict func([64]float64)int}{
		{"softmax_baseline",basePredict},
		{"one_vs_rest",ovrPredict},
	}{
		result.ClassMetrics=append(result.ClassMetrics,
			up98bEval(arm.name,"train","all",-1,arm.predict),
			up98bEval(arm.name,"heldout_recombination","all",-1,arm.predict),
			up98bEval(arm.name,"heldout_recombination","STORE",up97bStore,arm.predict),
			up98bEval(arm.name,"heldout_recombination","OBSERVE",up97bObserve,arm.predict),
			up98bEval(arm.name,"heldout_recombination","REPORT",up97bReport,arm.predict),
		)
		for vi:=0;vi<9;vi++{result.VerbMetrics=append(result.VerbMetrics,up98bVerbEval(arm.name,vi,arm.predict))}
		result.Confusions=append(result.Confusions,up98bConfusion(arm.name,arm.predict))
	}
	return result,nil
}
