package unitary

import "math"

const UP207BNativeOffsetSchema="wingless.up207b-native-offset-proxy.v1"

type UP207BMetric struct{
	Composition string `json:"composition"`
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
type UP207BSummary struct{
	Composition string `json:"composition"`
	EvaluationPoints int `json:"evaluation_points"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	NativeOffsetMeanAbsoluteError float64 `json:"native_offset_mean_absolute_error"`
	OracleAnchorMeanAbsoluteError float64 `json:"oracle_anchor_mean_absolute_error"`
	StaticMaxAbsoluteError float64 `json:"static_max_absolute_error"`
	NativeOffsetMaxAbsoluteError float64 `json:"native_offset_max_absolute_error"`
	OracleAnchorMaxAbsoluteError float64 `json:"oracle_anchor_max_absolute_error"`
	NativeBetterThanStatic int `json:"native_better_than_static"`
}
type UP207BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP206BSeal string `json:"source_up206b_seal"`
	TemplatePhases []int `json:"template_phases"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	Features []string `json:"features"`
	RidgeLambda float64 `json:"ridge_lambda"`
	TrainingPoints int `json:"training_points"`
	FeatureMeans []float64 `json:"feature_means"`
	FeatureStdDevs []float64 `json:"feature_std_devs"`
	Coefficients []float64 `json:"coefficients"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	HeldoutAnchorMeasurementUsed bool `json:"heldout_anchor_measurement_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP207BMetric `json:"metrics"`
	Summaries []UP207BSummary `json:"summaries"`
}
func RunUP207B()(UP207BResult,error){
	template:=[]int{55,56,57}
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{88,89,90}
	comps:=[]UP197BComposition{
		{Name:"mixed4",Canonical:[]int{0,5,1,6}},
		{Name:"observe4",Canonical:[]int{5,6,7,8}},
		{Name:"store4",Canonical:[]int{0,1,2,3}},
		{Name:"cross3",Canonical:[]int{0,5,13}},
	}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type family struct{c UP197BComposition;perms [][]int;means []float64;canon int}
	fams:=[]family{}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical);means:=make([]float64,len(perms));canon:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(c.Canonical){canon=i}
			s:=0.0
			for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:c.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=s/float64(len(template))
		}
		fams=append(fams,family{c:c,perms:perms,means:means,canon:canon})
	}
	type sample struct{x [4]float64;y float64}
	samples:=[]sample{}
	for _,f:=range fams{for _,ph:=range train{
		pt:=up191bPoint(model,UP191BProfile{Name:f.c.Name+"_train",SurfaceIndices:f.c.Canonical},ph)
		samples=append(samples,sample{x:up193bFeature(ph,f.c.Canonical),y:pt.RequiredMassFactor-f.means[f.canon]})
	}}
	var mean,std [4]float64
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
	coef,ok:=up193bSolve(a,b);if !ok{return UP207BResult{},nil}
	res:=UP207BResult{
		Schema:UP207BNativeOffsetSchema,Experiment:"UP-207B-native-offset-proxy",SourceUP206BSeal:"dc08d20b1bde815dbb69921705519521dbca58bb",
		TemplatePhases:template,TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},
		RidgeLambda:1e-6,TrainingPoints:len(samples),
		FeatureMeans:[]float64{mean[0],mean[1],mean[2],mean[3]},FeatureStdDevs:[]float64{std[0],std[1],std[2],std[3]},
		Coefficients:[]float64{coef[0],coef[1],coef[2],coef[3],coef[4]},
		HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,PhaseInputUsed:false,ParityInputUsed:false,HeldoutAnchorMeasurementUsed:false,MaintenanceTriggered:false,
	}
	for _,f:=range fams{
		s:=UP207BSummary{Composition:f.c.Name};sumS,sumN,sumO:=0.0,0.0,0.0
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.c.Canonical)
			z:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			shift:=0.0;for i:=0;i<5;i++{shift+=coef[i]*z[i]}
			canonPt:=up191bPoint(model,UP191BProfile{Name:f.c.Name+"_eval",SurfaceIndices:f.c.Canonical},ph)
			oracleShift:=canonPt.RequiredMassFactor-f.means[f.canon]
			for i,p:=range f.perms{
				if i==f.canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:f.c.Name+"_eval",SurfaceIndices:p},ph)
				sp:=f.means[i];np:=sp+shift;op:=sp+oracleShift
				se:=math.Abs(sp-pt.RequiredMassFactor);ne:=math.Abs(np-pt.RequiredMassFactor);oe:=math.Abs(op-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP207BMetric{Composition:f.c.Name,Permutation:append([]int(nil),p...),Phase:ph,ActualRequiredFactor:pt.RequiredMassFactor,StaticPredictedFactor:sp,NativeOffsetPredictedFactor:np,OracleAnchorPredictedFactor:op,StaticAbsoluteError:se,NativeOffsetAbsoluteError:ne,OracleAnchorAbsoluteError:oe})
				s.EvaluationPoints++;sumS+=se;sumN+=ne;sumO+=oe
				if se>s.StaticMaxAbsoluteError{s.StaticMaxAbsoluteError=se}
				if ne>s.NativeOffsetMaxAbsoluteError{s.NativeOffsetMaxAbsoluteError=ne}
				if oe>s.OracleAnchorMaxAbsoluteError{s.OracleAnchorMaxAbsoluteError=oe}
				if ne<se{s.NativeBetterThanStatic++}
			}
		}
		s.StaticMeanAbsoluteError=sumS/float64(s.EvaluationPoints);s.NativeOffsetMeanAbsoluteError=sumN/float64(s.EvaluationPoints);s.OracleAnchorMeanAbsoluteError=sumO/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
