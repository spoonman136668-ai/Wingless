package unitary

import "math"

const UP200BCombinedHistorySchema="wingless.up200b-combined-history-calibration.v1"

type UP200BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
	PairwiseAbsoluteError float64 `json:"pairwise_absolute_error"`
	PositionAbsoluteError float64 `json:"position_absolute_error"`
	CombinedAbsoluteError float64 `json:"combined_absolute_error"`
}
type UP200BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	PairwiseMeanAbsoluteError float64 `json:"pairwise_mean_absolute_error"`
	PositionMeanAbsoluteError float64 `json:"position_mean_absolute_error"`
	CombinedMeanAbsoluteError float64 `json:"combined_mean_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	PairwiseMaxAbsoluteError float64 `json:"pairwise_max_absolute_error"`
	PositionMaxAbsoluteError float64 `json:"position_max_absolute_error"`
	CombinedMaxAbsoluteError float64 `json:"combined_max_absolute_error"`
	CombinedBetterThanPairwise int `json:"combined_better_than_pairwise"`
	CombinedBetterThanPosition int `json:"combined_better_than_position"`
	CombinedBetterThanBaseline int `json:"combined_better_than_baseline"`
}
type UP200BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP199BSeal string `json:"source_up199b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	PairwiseFeatureCount int `json:"pairwise_feature_count"`
	PositionFeatureCount int `json:"position_feature_count"`
	CombinedFeatureCount int `json:"combined_feature_count"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ArchitectureSearchUsed bool `json:"architecture_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP200BMetric `json:"metrics"`
	Summaries []UP200BSummary `json:"summaries"`
}
func up200bCombined(canonical,perm []int)[]float64{
	x:=up198bPairSlice(canonical,perm)
	x=append(x,up199bPositionFeatures(canonical,perm)...)
	return x
}
func RunUP200B()(UP200BResult,error){
	train:=[]int{55,56,57};eval:=[]int{67,68,69}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP200BResult{Schema:UP200BCombinedHistorySchema,Experiment:"UP-200B-combined-history-calibration",SourceUP199BSeal:"0b107d58dded41b4c239a52e61160fa38aac1489",TrainingPhases:train,EvaluationPhases:eval,PairwiseFeatureCount:12,PositionFeatureCount:16,CombinedFeatureCount:28,RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,ArchitectureSearchUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		pairX,posX,combX:=[][]float64{},[][]float64{},[][]float64{};ys:=[]float64{};mean:=0.0
		for _,phase:=range train{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase)
			pairX=append(pairX,up198bPairSlice(c.Canonical,p));posX=append(posX,up199bPositionFeatures(c.Canonical,p));combX=append(combX,up200bCombined(c.Canonical,p));ys=append(ys,pt.RequiredMassFactor);mean+=pt.RequiredMassFactor
		}}
		mean/=float64(len(ys))
		pairCoef,ok:=up198bFit(pairX,ys,1e-6);if !ok{return UP200BResult{},nil}
		posCoef,ok:=up198bFit(posX,ys,1e-6);if !ok{return UP200BResult{},nil}
		combCoef,ok:=up198bFit(combX,ys,1e-6);if !ok{return UP200BResult{},nil}
		s:=UP200BSummary{Composition:c.Name};sumB,sumP,sumPos,sumC:=0.0,0.0,0.0,0.0
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
			pp:=up198bPredict(pairCoef,up198bPairSlice(c.Canonical,p));gp:=up198bPredict(posCoef,up199bPositionFeatures(c.Canonical,p));cp:=up198bPredict(combCoef,up200bCombined(c.Canonical,p))
			be:=math.Abs(mean-pt.RequiredMassFactor);pe:=math.Abs(pp-pt.RequiredMassFactor);ge:=math.Abs(gp-pt.RequiredMassFactor);ce:=math.Abs(cp-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP200BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,BaselineAbsoluteError:be,PairwiseAbsoluteError:pe,PositionAbsoluteError:ge,CombinedAbsoluteError:ce})
			s.EvaluationPoints++;sumB+=be;sumP+=pe;sumPos+=ge;sumC+=ce
			if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be};if pe>s.PairwiseMaxAbsoluteError{s.PairwiseMaxAbsoluteError=pe};if ge>s.PositionMaxAbsoluteError{s.PositionMaxAbsoluteError=ge};if ce>s.CombinedMaxAbsoluteError{s.CombinedMaxAbsoluteError=ce}
			if ce<pe{s.CombinedBetterThanPairwise++};if ce<ge{s.CombinedBetterThanPosition++};if ce<be{s.CombinedBetterThanBaseline++}
		}}
		s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints);s.PairwiseMeanAbsoluteError=sumP/float64(s.EvaluationPoints);s.PositionMeanAbsoluteError=sumPos/float64(s.EvaluationPoints);s.CombinedMeanAbsoluteError=sumC/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
