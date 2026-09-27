package unitary

import "math"

const UP199BGlobalPositionSchema="wingless.up199b-global-position-calibration.v1"

type UP199BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	BaselinePredictedFactor float64 `json:"baseline_predicted_factor"`
	PairwisePredictedFactor float64 `json:"pairwise_predicted_factor"`
	PositionPredictedFactor float64 `json:"position_predicted_factor"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
	PairwiseAbsoluteError float64 `json:"pairwise_absolute_error"`
	PositionAbsoluteError float64 `json:"position_absolute_error"`
}
type UP199BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	PairwiseMeanAbsoluteError float64 `json:"pairwise_mean_absolute_error"`
	PositionMeanAbsoluteError float64 `json:"position_mean_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	PairwiseMaxAbsoluteError float64 `json:"pairwise_max_absolute_error"`
	PositionMaxAbsoluteError float64 `json:"position_max_absolute_error"`
	PositionBetterThanPairwise int `json:"position_better_than_pairwise"`
	PositionBetterThanBaseline int `json:"position_better_than_baseline"`
}
type UP199BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP198BSeal string `json:"source_up198b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Compositions []UP197BComposition `json:"compositions"`
	PairwiseFeatureCount int `json:"pairwise_feature_count"`
	PositionFeatureCount int `json:"position_feature_count"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ArchitectureSearchUsed bool `json:"architecture_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP199BMetric `json:"metrics"`
	Summaries []UP199BSummary `json:"summaries"`
}
func up199bPositionFeatures(canonical,perm []int)[]float64{
	pos:=map[int]int{};for i,v:=range canonical{pos[v]=i}
	x:=make([]float64,16)
	for seqPos,v:=range perm{x[pos[v]*4+seqPos]=1}
	return x
}
func RunUP199B()(UP199BResult,error){
	train:=[]int{55,56,57};eval:=[]int{64,65,66}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP199BResult{Schema:UP199BGlobalPositionSchema,Experiment:"UP-199B-global-position-calibration",SourceUP198BSeal:"cb538e82e05c087d2939ba130339b04870aeb953",TrainingPhases:train,EvaluationPhases:eval,Compositions:comps,PairwiseFeatureCount:12,PositionFeatureCount:16,RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,ArchitectureSearchUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		pairX,posX:=[][]float64{},[][]float64{};ys:=[]float64{};mean:=0.0
		for _,phase:=range train{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase)
			pairX=append(pairX,up198bPairSlice(c.Canonical,p));posX=append(posX,up199bPositionFeatures(c.Canonical,p));ys=append(ys,pt.RequiredMassFactor);mean+=pt.RequiredMassFactor
		}}
		mean/=float64(len(ys))
		pairCoef,ok:=up198bFit(pairX,ys,1e-6);if !ok{return UP199BResult{},nil}
		posCoef,ok:=up198bFit(posX,ys,1e-6);if !ok{return UP199BResult{},nil}
		s:=UP199BSummary{Composition:c.Name};sumB,sumP,sumPos:=0.0,0.0,0.0
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
			pp:=up198bPredict(pairCoef,up198bPairSlice(c.Canonical,p));gp:=up198bPredict(posCoef,up199bPositionFeatures(c.Canonical,p))
			be:=math.Abs(mean-pt.RequiredMassFactor);pe:=math.Abs(pp-pt.RequiredMassFactor);ge:=math.Abs(gp-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP199BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,BaselinePredictedFactor:mean,PairwisePredictedFactor:pp,PositionPredictedFactor:gp,BaselineAbsoluteError:be,PairwiseAbsoluteError:pe,PositionAbsoluteError:ge})
			s.EvaluationPoints++;sumB+=be;sumP+=pe;sumPos+=ge
			if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be};if pe>s.PairwiseMaxAbsoluteError{s.PairwiseMaxAbsoluteError=pe};if ge>s.PositionMaxAbsoluteError{s.PositionMaxAbsoluteError=ge}
			if ge<pe{s.PositionBetterThanPairwise++};if ge<be{s.PositionBetterThanBaseline++}
		}}
		s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints);s.PairwiseMeanAbsoluteError=sumP/float64(s.EvaluationPoints);s.PositionMeanAbsoluteError=sumPos/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
