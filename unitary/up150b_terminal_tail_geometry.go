package unitary

import "math"

const UP150BGeometrySchema="wingless.up150b-terminal-tail-geometry.v1"

type UP150BTailMetric struct{
	Shift int `json:"shift"`
	PriorClass string `json:"prior_class"`
	TailExamples []string `json:"tail_examples"`
	TailUpdateNorm float64 `json:"tail_update_norm"`
	OldAnchorUpdateNorm float64 `json:"old_anchor_update_norm"`
	CosineTailVsOldAnchor float64 `json:"cosine_tail_vs_old_anchor"`
	MeanCorrectClassProbability float64 `json:"mean_correct_class_probability"`
}
type UP150BPairwise struct{
	ShiftA int `json:"shift_a"`
	ShiftB int `json:"shift_b"`
	Cosine float64 `json:"cosine"`
}
type UP150BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP149BSeal string `json:"source_up149b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TailExamplesPerShift int `json:"tail_examples_per_shift"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveTailSelection bool `json:"adaptive_tail_selection"`
	ExtraTrainingUsed bool `json:"extra_training_used"`
	Tails []UP150BTailMetric `json:"tails"`
	Pairwise []UP150BPairwise `json:"pairwise"`
}

func up150bGateVector(g *up129bGate)[]float64{
	out:=make([]float64,0,387)
	for k:=0;k<3;k++{for i:=0;i<128;i++{out=append(out,g.w[k][i])};out=append(out,g.b[k])}
	return out
}
func up150bDelta(before,after *up129bGate)[]float64{
	a,b:=up150bGateVector(before),up150bGateVector(after)
	out:=make([]float64,len(a));for i:=range out{out[i]=b[i]-a[i]};return out
}
func up150bAdd(dst,src []float64){for i:=range dst{dst[i]+=src[i]}}
func up150bNorm(v []float64)float64{s:=0.0;for _,x:=range v{s+=x*x};return math.Sqrt(s)}
func up150bCos(a,b []float64)float64{
	aa,bb,d:=0.0,0.0,0.0
	for i:=range a{d+=a[i]*b[i];aa+=a[i]*a[i];bb+=b[i]*b[i]}
	if aa<=1e-24||bb<=1e-24{return 0}
	return d/math.Sqrt(aa*bb)
}
func up150bCommonPrefix(o,r [64]float64)*up129bGate{
	g:=up129bTrainGate(o,r);epoch:=0
	for _,shift:=range []int{0,6,12}{
		for rep:=0;rep<5;rep++{up148bTrainEpoch(g,shift,epoch,o,r);epoch++}
	}
	return g
}
func up150bExampleDelta(base *up129bGate,e up135bExample,o,r [64]float64)[]float64{
	clone:=*base;up135bStep(&clone,e,o,r);return up150bDelta(base,&clone)
}
func up150bOldAnchorDelta(base *up129bGate,o,r [64]float64)[]float64{
	out:=make([]float64,387)
	old:=up132bOldSurfaces()
	for i,s:=range old{
		clone:=*base
		subjectIndex:=(i+15)%2
		up129bGateStep(&clone,up121bOriginalNames[subjectIndex],s.verb,s.class,o,r)
		up150bAdd(out,up150bDelta(base,&clone))
	}
	return out
}
func up150bPriorClass(shift int)string{if shift==6||shift==18{return "protective"};return "damaging"}

func RunUP150B()(UP150BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	oldDelta:=up150bOldAnchorDelta(common,o,r)
	oldNorm:=up150bNorm(oldDelta)
	examples:=up135bNewExamples()
	shifts:=[]int{0,6,12,18}
	vectors:=make([][]float64,0,4)
	res:=UP150BResult{
		Schema:UP150BGeometrySchema,Experiment:"UP-150B-terminal-tail-geometry",
		SourceUP149BSeal:"6d67cc6d4b3b582879e7771c44434fea3fca853f",
		CommonPrefixEpochs:15,TailExamplesPerShift:4,
		DiagnosticUpdatesRetained:false,AdaptiveTailSelection:false,ExtraTrainingUsed:false,
	}
	for _,shift:=range shifts{
		full:=up139bRotate(examples,shift)
		v:=make([]float64,387)
		names:=[]string{}
		prob:=0.0
		for i:=20;i<24;i++{
			e:=full[i]
			up150bAdd(v,up150bExampleDelta(common,e,o,r))
			names=append(names,e.name+" "+e.verb)
			raw,proj:=up129bViews(e.name,e.verb,o,r)
			p:=up129bProbs(common,raw,proj)
			prob+=p[e.class]
		}
		vectors=append(vectors,v)
		res.Tails=append(res.Tails,UP150BTailMetric{
			Shift:shift,PriorClass:up150bPriorClass(shift),TailExamples:names,
			TailUpdateNorm:up150bNorm(v),OldAnchorUpdateNorm:oldNorm,
			CosineTailVsOldAnchor:up150bCos(v,oldDelta),MeanCorrectClassProbability:prob/4,
		})
	}
	for i:=0;i<len(shifts);i++{for j:=i+1;j<len(shifts);j++{
		res.Pairwise=append(res.Pairwise,UP150BPairwise{ShiftA:shifts[i],ShiftB:shifts[j],Cosine:up150bCos(vectors[i],vectors[j])})
	}}
	return res,nil
}
