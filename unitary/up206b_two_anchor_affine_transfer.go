package unitary

import "math"

const UP206BAffineSchema="wingless.up206b-two-anchor-affine-transfer.v1"

type UP206BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	OneAnchorPredictedFactor float64 `json:"one_anchor_predicted_factor"`
	AffinePredictedFactor float64 `json:"affine_predicted_factor"`
	OneAnchorAbsoluteError float64 `json:"one_anchor_absolute_error"`
	AffineAbsoluteError float64 `json:"affine_absolute_error"`
}
type UP206BScale struct{
	Composition string `json:"composition"`
	Phase int `json:"phase"`
	Scale float64 `json:"scale"`
	Offset float64 `json:"offset"`
}
type UP206BSummary struct{
	Composition string `json:"composition"`
	PermutationCount int `json:"permutation_count"`
	EvaluationPoints int `json:"evaluation_points"`
	OneAnchorMeanAbsoluteError float64 `json:"one_anchor_mean_absolute_error"`
	AffineMeanAbsoluteError float64 `json:"affine_mean_absolute_error"`
	OneAnchorMaxAbsoluteError float64 `json:"one_anchor_max_absolute_error"`
	AffineMaxAbsoluteError float64 `json:"affine_max_absolute_error"`
	AffineBetterThanOneAnchor int `json:"affine_better_than_one_anchor"`
}
type UP206BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP205BSeal string `json:"source_up205b_seal"`
	TemplatePhases []int `json:"template_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	AnchorRule string `json:"anchor_rule"`
	AnchorsPerPhase int `json:"anchors_per_phase"`
	HeldoutAnchorSearchUsed bool `json:"heldout_anchor_search_used"`
	NonlinearTransformUsed bool `json:"nonlinear_transform_used"`
	HeldoutFeatureFittingUsed bool `json:"heldout_feature_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP206BMetric `json:"metrics"`
	Scales []UP206BScale `json:"scales"`
	Summaries []UP206BSummary `json:"summaries"`
}
func up206bReverse(in []int)[]int{
	out:=append([]int(nil),in...)
	for i,j:=0,len(out)-1;i<j;i,j=i+1,j-1{out[i],out[j]=out[j],out[i]}
	return out
}
func RunUP206B()(UP206BResult,error){
	template:=[]int{55,56,57};eval:=[]int{85,86,87}
	comps:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP206BResult{
		Schema:UP206BAffineSchema,Experiment:"UP-206B-two-anchor-affine-transfer",
		SourceUP205BSeal:"57554c8fe3a98d9ed21ffeb5f0c1a15d23bee3c7",
		TemplatePhases:template,EvaluationPhases:eval,
		AnchorRule:"canonical_and_reverse",AnchorsPerPhase:2,
		HeldoutAnchorSearchUsed:false,NonlinearTransformUsed:false,HeldoutFeatureFittingUsed:false,MaintenanceTriggered:false,
	}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		means:=make([]float64,len(perms));grand:=0.0
		canonKey:=up201bKey(c.Canonical);rev:=up206bReverse(c.Canonical);revKey:=up201bKey(rev)
		canonIndex,revIndex:=-1,-1
		for i,p:=range perms{
			k:=up201bKey(p);if k==canonKey{canonIndex=i};if k==revKey{revIndex=i}
			s:=0.0
			for _,phase:=range template{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_template",SurfaceIndices:p},phase)
				s+=pt.RequiredMassFactor
			}
			means[i]=s/float64(len(template));grand+=means[i]
		}
		grand/=float64(len(means))
		resid:=make([]float64,len(means));for i,v:=range means{resid[i]=v-grand}
		sumO,sumA:=0.0,0.0
		summary:=UP206BSummary{Composition:c.Name,PermutationCount:len(perms)}
		for _,phase:=range eval{
			canonPt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:perms[canonIndex]},phase)
			revPt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:perms[revIndex]},phase)
			oneOffset:=canonPt.RequiredMassFactor-resid[canonIndex]
			den:=resid[revIndex]-resid[canonIndex]
			scale:=1.0
			if math.Abs(den)>1e-12{scale=(revPt.RequiredMassFactor-canonPt.RequiredMassFactor)/den}
			offset:=canonPt.RequiredMassFactor-scale*resid[canonIndex]
			res.Scales=append(res.Scales,UP206BScale{Composition:c.Name,Phase:phase,Scale:scale,Offset:offset})
			for i,p:=range perms{
				if i==canonIndex||i==revIndex{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
				op:=oneOffset+resid[i];ap:=offset+scale*resid[i]
				oe:=math.Abs(op-pt.RequiredMassFactor);ae:=math.Abs(ap-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP206BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,OneAnchorPredictedFactor:op,AffinePredictedFactor:ap,OneAnchorAbsoluteError:oe,AffineAbsoluteError:ae})
				summary.EvaluationPoints++;sumO+=oe;sumA+=ae
				if oe>summary.OneAnchorMaxAbsoluteError{summary.OneAnchorMaxAbsoluteError=oe}
				if ae>summary.AffineMaxAbsoluteError{summary.AffineMaxAbsoluteError=ae}
				if ae<oe{summary.AffineBetterThanOneAnchor++}
			}
		}
		summary.OneAnchorMeanAbsoluteError=sumO/float64(summary.EvaluationPoints)
		summary.AffineMeanAbsoluteError=sumA/float64(summary.EvaluationPoints)
		res.Summaries=append(res.Summaries,summary)
	}
	return res,nil
}
