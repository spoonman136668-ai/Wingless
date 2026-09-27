package unitary

import "math"

const UP205BAnchorFamilySchema="wingless.up205b-anchor-offset-family-transfer.v1"

type UP205BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	IdentityMeanPredictedFactor float64 `json:"identity_mean_predicted_factor"`
	AnchorOffsetPredictedFactor float64 `json:"anchor_offset_predicted_factor"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	AnchorOffsetAbsoluteError float64 `json:"anchor_offset_absolute_error"`
}
type UP205BSummary struct{
	Composition string `json:"composition"`
	PermutationCount int `json:"permutation_count"`
	EvaluationPoints int `json:"evaluation_points"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	AnchorOffsetMeanAbsoluteError float64 `json:"anchor_offset_mean_absolute_error"`
	IdentityMaxAbsoluteError float64 `json:"identity_max_absolute_error"`
	AnchorOffsetMaxAbsoluteError float64 `json:"anchor_offset_max_absolute_error"`
	AnchorOffsetBetterThanIdentity int `json:"anchor_offset_better_than_identity"`
}
type UP205BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP204BSeal string `json:"source_up204b_seal"`
	TemplatePhases []int `json:"template_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	AnchorRule string `json:"anchor_rule"`
	AnchorsPerPhase int `json:"anchors_per_phase"`
	HeldoutAnchorSearchUsed bool `json:"heldout_anchor_search_used"`
	HeldoutFeatureFittingUsed bool `json:"heldout_feature_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP205BMetric `json:"metrics"`
	Summaries []UP205BSummary `json:"summaries"`
}
func RunUP205B()(UP205BResult,error){
	template:=[]int{55,56,57};eval:=[]int{82,83,84}
	comps:=[]UP197BComposition{{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP205BResult{Schema:UP205BAnchorFamilySchema,Experiment:"UP-205B-anchor-offset-family-transfer",SourceUP204BSeal:"d64d7b33463536887bf759a15215367eabd26800",TemplatePhases:template,EvaluationPhases:eval,AnchorRule:"canonical_permutation",AnchorsPerPhase:1,HeldoutAnchorSearchUsed:false,HeldoutFeatureFittingUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		means:=make([]float64,len(perms));grand:=0.0;anchorIndex:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(c.Canonical){anchorIndex=i}
			s:=0.0
			for _,phase:=range template{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_template",SurfaceIndices:p},phase)
				s+=pt.RequiredMassFactor
			}
			means[i]=s/float64(len(template));grand+=means[i]
		}
		grand/=float64(len(means))
		residual:=make([]float64,len(means))
		for i,v:=range means{residual[i]=v-grand}
		sumI,sumA:=0.0,0.0
		summary:=UP205BSummary{Composition:c.Name,PermutationCount:len(perms)}
		for _,phase:=range eval{
			anchorPt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:perms[anchorIndex]},phase)
			offset:=anchorPt.RequiredMassFactor-residual[anchorIndex]
			for i,p:=range perms{
				if i==anchorIndex{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
				ip:=means[i];ap:=offset+residual[i]
				ie:=math.Abs(ip-pt.RequiredMassFactor);ae:=math.Abs(ap-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP205BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,IdentityMeanPredictedFactor:ip,AnchorOffsetPredictedFactor:ap,IdentityMeanAbsoluteError:ie,AnchorOffsetAbsoluteError:ae})
				summary.EvaluationPoints++;sumI+=ie;sumA+=ae
				if ie>summary.IdentityMaxAbsoluteError{summary.IdentityMaxAbsoluteError=ie}
				if ae>summary.AnchorOffsetMaxAbsoluteError{summary.AnchorOffsetMaxAbsoluteError=ae}
				if ae<ie{summary.AnchorOffsetBetterThanIdentity++}
			}
		}
		summary.IdentityMeanAbsoluteError=sumI/float64(summary.EvaluationPoints)
		summary.AnchorOffsetMeanAbsoluteError=sumA/float64(summary.EvaluationPoints)
		res.Summaries=append(res.Summaries,summary)
	}
	return res,nil
}
