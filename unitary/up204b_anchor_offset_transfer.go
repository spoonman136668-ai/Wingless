package unitary

import "math"

const UP204BAnchorSchema="wingless.up204b-anchor-offset-transfer.v1"

type UP204BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	IdentityMeanPredictedFactor float64 `json:"identity_mean_predicted_factor"`
	AnchorOffsetPredictedFactor float64 `json:"anchor_offset_predicted_factor"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	AnchorOffsetAbsoluteError float64 `json:"anchor_offset_absolute_error"`
}
type UP204BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	IdentityMeanAbsoluteError float64 `json:"identity_mean_absolute_error"`
	AnchorOffsetMeanAbsoluteError float64 `json:"anchor_offset_mean_absolute_error"`
	IdentityMaxAbsoluteError float64 `json:"identity_max_absolute_error"`
	AnchorOffsetMaxAbsoluteError float64 `json:"anchor_offset_max_absolute_error"`
	AnchorOffsetBetterThanIdentity int `json:"anchor_offset_better_than_identity"`
}
type UP204BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP203BSeal string `json:"source_up203b_seal"`
	TemplatePhases []int `json:"template_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	AnchorRule string `json:"anchor_rule"`
	AnchorsPerPhase int `json:"anchors_per_phase"`
	HeldoutAnchorSearchUsed bool `json:"heldout_anchor_search_used"`
	HeldoutFeatureFittingUsed bool `json:"heldout_feature_fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP204BMetric `json:"metrics"`
	Summaries []UP204BSummary `json:"summaries"`
}
func RunUP204B()(UP204BResult,error){
	template:=[]int{55,56,57};eval:=[]int{79,80,81}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP204BResult{Schema:UP204BAnchorSchema,Experiment:"UP-204B-anchor-offset-transfer",SourceUP203BSeal:"8a379456d273d79b9410c47e1009b3ca05a2cafa",TemplatePhases:template,EvaluationPhases:eval,AnchorRule:"canonical_permutation",AnchorsPerPhase:1,HeldoutAnchorSearchUsed:false,HeldoutFeatureFittingUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		means:=make([]float64,len(perms));grand:=0.0
		anchorIndex:=-1
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
		summary:=UP204BSummary{Composition:c.Name}
		sumI,sumA:=0.0,0.0
		for _,phase:=range eval{
			anchorPt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:perms[anchorIndex]},phase)
			offset:=anchorPt.RequiredMassFactor-residual[anchorIndex]
			for i,p:=range perms{
				if i==anchorIndex{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
				ip:=means[i];ap:=offset+residual[i]
				ie:=math.Abs(ip-pt.RequiredMassFactor);ae:=math.Abs(ap-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP204BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,IdentityMeanPredictedFactor:ip,AnchorOffsetPredictedFactor:ap,IdentityMeanAbsoluteError:ie,AnchorOffsetAbsoluteError:ae})
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
