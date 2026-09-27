package unitary

import "math"

const UP208BFamilyHoldoutSchema="wingless.up208b-family-holdout-native-offset.v1"

type UP208BMetric struct{
	HeldoutFamily string `json:"heldout_family"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	StaticPredictedFactor float64 `json:"static_predicted_factor"`
	NativeOffsetPredictedFactor float64 `json:"native_offset_predicted_factor"`
	OracleAnchorPredictedFactor float64 `json:"oracle_anchor_predicted_factor"`
	StaticAbsoluteError float64 `json:"static_absolute_error"`
	NativeOffsetAbsoluteError float64 `json:"native_offset_absolute_error"`
	OracleAnchorAbsoluteError float64 `json:"oracle_anchor_absolute_error"`
}
type UP208BSummary struct{
	HeldoutFamily string `json:"heldout_family"`
	TrainingFamilies []string `json:"training_families"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	NativeOffsetMeanAbsoluteError float64 `json:"native_offset_mean_absolute_error"`
	OracleAnchorMeanAbsoluteError float64 `json:"oracle_anchor_mean_absolute_error"`
	StaticMaxAbsoluteError float64 `json:"static_max_absolute_error"`
	NativeOffsetMaxAbsoluteError float64 `json:"native_offset_max_absolute_error"`
	OracleAnchorMaxAbsoluteError float64 `json:"oracle_anchor_max_absolute_error"`
	NativeBetterThanStatic int `json:"native_better_than_static"`
}
type UP208BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP207BSeal string `json:"source_up207b_seal"`
	TemplatePhases []int `json:"template_phases"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	RidgeLambda float64 `json:"ridge_lambda"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	FamilyIdentityInputUsed bool `json:"family_identity_input_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutAnchorMeasurementUsedByNative bool `json:"heldout_anchor_measurement_used_by_native"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP208BMetric `json:"metrics"`
	Summaries []UP208BSummary `json:"summaries"`
}
type up208bFamily struct{
	name string
	canonical []int
	perms [][]int
	means []float64
	canon int
}
type up208bSample struct{x [4]float64;y float64}
func up208bFit(samples []up208bSample)(mean,std [4]float64,coef [5]float64,ok bool){
	for _,s:=range samples{for j:=0;j<4;j++{mean[j]+=s.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(samples))}
	for _,s:=range samples{for j:=0;j<4;j++{d:=s.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(samples)));if std[j]==0{std[j]=1}}
	var a [5][5]float64;var b [5]float64
	for _,s:=range samples{
		v:=[5]float64{1,(s.x[0]-mean[0])/std[0],(s.x[1]-mean[1])/std[1],(s.x[2]-mean[2])/std[2],(s.x[3]-mean[3])/std[3]}
		for i:=0;i<5;i++{b[i]+=v[i]*s.y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
	}
	for i:=1;i<5;i++{a[i][i]+=1e-6}
	c,o:=up193bSolve(a,b)
	return mean,std,c,o
}
func RunUP208B()(UP208BResult,error){
	template:=[]int{55,56,57};train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72};eval:=[]int{91,92,93}
	specs:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	fams:=[]up208bFamily{}
	for _,c:=range specs{
		perms:=up196bPermutations(c.Canonical);means:=make([]float64,len(perms));canon:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(c.Canonical){canon=i}
			sum:=0.0
			for _,ph:=range template{sum+=up191bPoint(model,UP191BProfile{Name:c.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=sum/float64(len(template))
		}
		fams=append(fams,up208bFamily{name:c.Name,canonical:c.Canonical,perms:perms,means:means,canon:canon})
	}
	res:=UP208BResult{
		Schema:UP208BFamilyHoldoutSchema,Experiment:"UP-208B-family-holdout-native-offset",SourceUP207BSeal:"8ce15c4234206ebb1bab899a88ea7d822df5a87e",
		TemplatePhases:template,TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},RidgeLambda:1e-6,
		HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,FamilyIdentityInputUsed:false,PhaseInputUsed:false,ParityInputUsed:false,HeldoutAnchorMeasurementUsedByNative:false,MaintenanceTriggered:false,
	}
	for hi,h:=range fams{
		samples:=[]up208bSample{};trainNames:=[]string{}
		for fi,f:=range fams{
			if fi==hi{continue};trainNames=append(trainNames,f.name)
			for _,ph:=range train{
				pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_train",SurfaceIndices:f.canonical},ph)
				samples=append(samples,up208bSample{x:up193bFeature(ph,f.canonical),y:pt.RequiredMassFactor-f.means[f.canon]})
			}
		}
		mean,std,coef,ok:=up208bFit(samples);if !ok{return UP208BResult{},nil}
		s:=UP208BSummary{HeldoutFamily:h.name,TrainingFamilies:trainNames,TrainingPoints:len(samples)}
		sumS,sumN,sumO:=0.0,0.0,0.0
		for _,ph:=range eval{
			x:=up193bFeature(ph,h.canonical)
			z:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			shift:=0.0;for i:=0;i<5;i++{shift+=coef[i]*z[i]}
			canonPt:=up191bPoint(model,UP191BProfile{Name:h.name+"_eval",SurfaceIndices:h.canonical},ph)
			oracleShift:=canonPt.RequiredMassFactor-h.means[h.canon]
			for i,p:=range h.perms{
				if i==h.canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:h.name+"_eval",SurfaceIndices:p},ph)
				sp:=h.means[i];np:=sp+shift;op:=sp+oracleShift
				se:=math.Abs(sp-pt.RequiredMassFactor);ne:=math.Abs(np-pt.RequiredMassFactor);oe:=math.Abs(op-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP208BMetric{HeldoutFamily:h.name,Permutation:append([]int(nil),p...),Phase:ph,ActualRequiredFactor:pt.RequiredMassFactor,StaticPredictedFactor:sp,NativeOffsetPredictedFactor:np,OracleAnchorPredictedFactor:op,StaticAbsoluteError:se,NativeOffsetAbsoluteError:ne,OracleAnchorAbsoluteError:oe})
				s.EvaluationPoints++;sumS+=se;sumN+=ne;sumO+=oe
				if se>s.StaticMaxAbsoluteError{s.StaticMaxAbsoluteError=se};if ne>s.NativeOffsetMaxAbsoluteError{s.NativeOffsetMaxAbsoluteError=ne};if oe>s.OracleAnchorMaxAbsoluteError{s.OracleAnchorMaxAbsoluteError=oe}
				if ne<se{s.NativeBetterThanStatic++}
			}
		}
		s.StaticMeanAbsoluteError=sumS/float64(s.EvaluationPoints);s.NativeOffsetMeanAbsoluteError=sumN/float64(s.EvaluationPoints);s.OracleAnchorMeanAbsoluteError=sumO/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
