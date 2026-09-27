package unitary

import "math"

const UP198BTripletSchema="wingless.up198b-triplet-transition-calibration.v1"

type UP198BMetric struct{
	Composition string `json:"composition"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	BaselinePredictedFactor float64 `json:"baseline_predicted_factor"`
	PairwisePredictedFactor float64 `json:"pairwise_predicted_factor"`
	TripletPredictedFactor float64 `json:"triplet_predicted_factor"`
	BaselineAbsoluteError float64 `json:"baseline_absolute_error"`
	PairwiseAbsoluteError float64 `json:"pairwise_absolute_error"`
	TripletAbsoluteError float64 `json:"triplet_absolute_error"`
}
type UP198BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	BaselineMeanAbsoluteError float64 `json:"baseline_mean_absolute_error"`
	PairwiseMeanAbsoluteError float64 `json:"pairwise_mean_absolute_error"`
	TripletMeanAbsoluteError float64 `json:"triplet_mean_absolute_error"`
	BaselineMaxAbsoluteError float64 `json:"baseline_max_absolute_error"`
	PairwiseMaxAbsoluteError float64 `json:"pairwise_max_absolute_error"`
	TripletMaxAbsoluteError float64 `json:"triplet_max_absolute_error"`
	TripletBetterThanPairwise int `json:"triplet_better_than_pairwise"`
	TripletBetterThanBaseline int `json:"triplet_better_than_baseline"`
}
type UP198BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP197BSeal string `json:"source_up197b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Compositions []UP197BComposition `json:"compositions"`
	PairwiseFeatureCount int `json:"pairwise_feature_count"`
	TripletFeatureCount int `json:"triplet_feature_count"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	ArchitectureSearchUsed bool `json:"architecture_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP198BMetric `json:"metrics"`
	Summaries []UP198BSummary `json:"summaries"`
}

func up198bTripletIndex(a,b,c int)int{
	k:=0
	for i:=0;i<4;i++{for j:=0;j<4;j++{if j==i{continue};for l:=0;l<4;l++{if l==i||l==j{continue};if i==a&&j==b&&l==c{return k};k++}}}
	return -1
}
func up198bTripletFeatures(canonical,perm []int)[24]float64{
	pos:=map[int]int{};for i,v:=range canonical{pos[v]=i}
	var x [24]float64
	for i:=0;i+2<len(perm);i++{
		idx:=up198bTripletIndex(pos[perm[i]],pos[perm[i+1]],pos[perm[i+2]])
		if idx>=0{x[idx]=1}
	}
	return x
}
func up198bSolve(a [][]float64,b []float64)([]float64,bool){
	n:=len(b)
	for i:=0;i<n;i++{
		p:=i;best:=math.Abs(a[i][i])
		for r:=i+1;r<n;r++{if v:=math.Abs(a[r][i]);v>best{best=v;p=r}}
		if best<1e-12{return nil,false}
		if p!=i{a[i],a[p]=a[p],a[i];b[i],b[p]=b[p],b[i]}
		d:=a[i][i];for c:=i;c<n;c++{a[i][c]/=d};b[i]/=d
		for r:=0;r<n;r++{if r==i{continue};f:=a[r][i];for c:=i;c<n;c++{a[r][c]-=f*a[i][c]};b[r]-=f*b[i]}
	}
	return b,true
}
func up198bFit(samplesX [][]float64,samplesY []float64,lambda float64)([]float64,bool){
	nf:=len(samplesX[0])+1
	a:=make([][]float64,nf);for i:=range a{a[i]=make([]float64,nf)}
	b:=make([]float64,nf)
	for n,x:=range samplesX{
		v:=make([]float64,nf);v[0]=1;copy(v[1:],x)
		for i:=0;i<nf;i++{b[i]+=v[i]*samplesY[n];for j:=0;j<nf;j++{a[i][j]+=v[i]*v[j]}}
	}
	for i:=1;i<nf;i++{a[i][i]+=lambda}
	return up198bSolve(a,b)
}
func up198bPredict(coef,x []float64)float64{y:=coef[0];for i,v:=range x{y+=coef[i+1]*v};return y}
func up198bPairSlice(canonical,perm []int)[]float64{x:=up197bFeatures(canonical,perm);return append([]float64(nil),x[:]...)}
func up198bTripletSlice(canonical,perm []int)[]float64{x:=up198bTripletFeatures(canonical,perm);return append([]float64(nil),x[:]...)}

func RunUP198B()(UP198BResult,error){
	train:=[]int{55,56,57};eval:=[]int{61,62,63}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP198BResult{Schema:UP198BTripletSchema,Experiment:"UP-198B-triplet-transition-calibration",SourceUP197BSeal:"f5728ed889c5e6ef479952e18a2a81f01090c17a",TrainingPhases:train,EvaluationPhases:eval,Compositions:comps,PairwiseFeatureCount:12,TripletFeatureCount:24,RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,ArchitectureSearchUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		pairX,tripX:=[][]float64{},[][]float64{};ys:=[]float64{};mean:=0.0
		for _,phase:=range train{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_train",SurfaceIndices:p},phase)
			pairX=append(pairX,up198bPairSlice(c.Canonical,p));tripX=append(tripX,up198bTripletSlice(c.Canonical,p));ys=append(ys,pt.RequiredMassFactor);mean+=pt.RequiredMassFactor
		}}
		mean/=float64(len(ys))
		pairCoef,ok:=up198bFit(pairX,ys,1e-6);if !ok{return UP198BResult{},nil}
		tripCoef,ok:=up198bFit(tripX,ys,1e-6);if !ok{return UP198BResult{},nil}
		s:=UP198BSummary{Composition:c.Name};sumB,sumP,sumT:=0.0,0.0,0.0
		for _,phase:=range eval{for _,p:=range perms{
			pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_eval",SurfaceIndices:p},phase)
			pp:=up198bPredict(pairCoef,up198bPairSlice(c.Canonical,p));tp:=up198bPredict(tripCoef,up198bTripletSlice(c.Canonical,p))
			be:=math.Abs(mean-pt.RequiredMassFactor);pe:=math.Abs(pp-pt.RequiredMassFactor);te:=math.Abs(tp-pt.RequiredMassFactor)
			res.Metrics=append(res.Metrics,UP198BMetric{Composition:c.Name,Permutation:append([]int(nil),p...),Phase:phase,ActualRequiredFactor:pt.RequiredMassFactor,BaselinePredictedFactor:mean,PairwisePredictedFactor:pp,TripletPredictedFactor:tp,BaselineAbsoluteError:be,PairwiseAbsoluteError:pe,TripletAbsoluteError:te})
			s.EvaluationPoints++;sumB+=be;sumP+=pe;sumT+=te
			if be>s.BaselineMaxAbsoluteError{s.BaselineMaxAbsoluteError=be};if pe>s.PairwiseMaxAbsoluteError{s.PairwiseMaxAbsoluteError=pe};if te>s.TripletMaxAbsoluteError{s.TripletMaxAbsoluteError=te}
			if te<pe{s.TripletBetterThanPairwise++};if te<be{s.TripletBetterThanBaseline++}
		}}
		s.BaselineMeanAbsoluteError=sumB/float64(s.EvaluationPoints);s.PairwiseMeanAbsoluteError=sumP/float64(s.EvaluationPoints);s.TripletMeanAbsoluteError=sumT/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
