package unitary

import "math"

const UP209BNativeDeltaSchema="wingless.up209b-family-holdout-native-delta.v1"

type UP209BMetric struct{
	HeldoutFamily string `json:"heldout_family"`
	Permutation []int `json:"permutation"`
	Phase int `json:"phase"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	StaticPredictedFactor float64 `json:"static_predicted_factor"`
	NativeDeltaPredictedFactor float64 `json:"native_delta_predicted_factor"`
	OracleAnchorPredictedFactor float64 `json:"oracle_anchor_predicted_factor"`
	StaticAbsoluteError float64 `json:"static_absolute_error"`
	NativeDeltaAbsoluteError float64 `json:"native_delta_absolute_error"`
	OracleAnchorAbsoluteError float64 `json:"oracle_anchor_absolute_error"`
}
type UP209BSummary struct{
	HeldoutFamily string `json:"heldout_family"`
	TrainingFamilies []string `json:"training_families"`
	TrainingPoints int `json:"training_points"`
	EvaluationPoints int `json:"evaluation_points"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	NativeDeltaMeanAbsoluteError float64 `json:"native_delta_mean_absolute_error"`
	OracleAnchorMeanAbsoluteError float64 `json:"oracle_anchor_mean_absolute_error"`
	StaticMaxAbsoluteError float64 `json:"static_max_absolute_error"`
	NativeDeltaMaxAbsoluteError float64 `json:"native_delta_max_absolute_error"`
	OracleAnchorMaxAbsoluteError float64 `json:"oracle_anchor_max_absolute_error"`
	NativeBetterThanStatic int `json:"native_better_than_static"`
}
type UP209BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP208BSeal string `json:"source_up208b_seal"`
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
	FutureStateInputUsed bool `json:"future_state_input_used"`
	HeldoutAnchorMeasurementUsedByNative bool `json:"heldout_anchor_measurement_used_by_native"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP209BMetric `json:"metrics"`
	Summaries []UP209BSummary `json:"summaries"`
}
type up209bSample struct{x [8]float64;y float64}

func up209bFeature(phase int,indices []int)[8]float64{
	c:=up193bFeature(phase,indices)
	p:=up193bFeature(phase-1,indices)
	return [8]float64{c[0],c[1],c[2],c[3],c[0]-p[0],c[1]-p[1],c[2]-p[2],c[3]-p[3]}
}
func up209bSolve(a [9][9]float64,b [9]float64)([9]float64,bool){
	for i:=0;i<9;i++{
		p:=i;best:=math.Abs(a[i][i])
		for r:=i+1;r<9;r++{if v:=math.Abs(a[r][i]);v>best{best=v;p=r}}
		if best<1e-12{return [9]float64{},false}
		if p!=i{a[i],a[p]=a[p],a[i];b[i],b[p]=b[p],b[i]}
		d:=a[i][i];for c:=i;c<9;c++{a[i][c]/=d};b[i]/=d
		for r:=0;r<9;r++{if r==i{continue};f:=a[r][i];for c:=i;c<9;c++{a[r][c]-=f*a[i][c]};b[r]-=f*b[i]}
	}
	return b,true
}
func up209bFit(samples []up209bSample)(mean,std [8]float64,coef [9]float64,ok bool){
	for _,s:=range samples{for j:=0;j<8;j++{mean[j]+=s.x[j]}}
	for j:=0;j<8;j++{mean[j]/=float64(len(samples))}
	for _,s:=range samples{for j:=0;j<8;j++{d:=s.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<8;j++{std[j]=math.Sqrt(std[j]/float64(len(samples)));if std[j]==0{std[j]=1}}
	var a [9][9]float64;var b [9]float64
	for _,s:=range samples{
		var v [9]float64;v[0]=1
		for j:=0;j<8;j++{v[j+1]=(s.x[j]-mean[j])/std[j]}
		for i:=0;i<9;i++{b[i]+=v[i]*s.y;for j:=0;j<9;j++{a[i][j]+=v[i]*v[j]}}
	}
	for i:=1;i<9;i++{a[i][i]+=1e-6}
	c,o:=up209bSolve(a,b)
	return mean,std,c,o
}
func RunUP209B()(UP209BResult,error){
	template:=[]int{55,56,57};train:=[]int{59,60,61,62,63,64,65,66,67,68,69,70,71,72};eval:=[]int{94,95,96}
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
	res:=UP209BResult{
		Schema:UP209BNativeDeltaSchema,Experiment:"UP-209B-family-holdout-native-delta",SourceUP208BSeal:"322ceca5b89ae730260360ef0f87ce927e2da837",
		TemplatePhases:template,TrainingPhases:train,EvaluationPhases:eval,
		Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin","delta_native_correct_count","delta_mean_absolute_margin","delta_near_zero_margin_count","delta_min_absolute_margin"},
		RidgeLambda:1e-6,HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,FamilyIdentityInputUsed:false,PhaseInputUsed:false,ParityInputUsed:false,FutureStateInputUsed:false,HeldoutAnchorMeasurementUsedByNative:false,MaintenanceTriggered:false,
	}
	for hi,h:=range fams{
		samples:=[]up209bSample{};trainNames:=[]string{}
		for fi,f:=range fams{
			if fi==hi{continue};trainNames=append(trainNames,f.name)
			for _,ph:=range train{
				pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_train",SurfaceIndices:f.canonical},ph)
				samples=append(samples,up209bSample{x:up209bFeature(ph,f.canonical),y:pt.RequiredMassFactor-f.means[f.canon]})
			}
		}
		mean,std,coef,ok:=up209bFit(samples);if !ok{return UP209BResult{},nil}
		s:=UP209BSummary{HeldoutFamily:h.name,TrainingFamilies:trainNames,TrainingPoints:len(samples)}
		sumS,sumN,sumO:=0.0,0.0,0.0
		for _,ph:=range eval{
			x:=up209bFeature(ph,h.canonical);var z [9]float64;z[0]=1
			for j:=0;j<8;j++{z[j+1]=(x[j]-mean[j])/std[j]}
			shift:=0.0;for i:=0;i<9;i++{shift+=coef[i]*z[i]}
			canonPt:=up191bPoint(model,UP191BProfile{Name:h.name+"_eval",SurfaceIndices:h.canonical},ph)
			oracleShift:=canonPt.RequiredMassFactor-h.means[h.canon]
			for i,p:=range h.perms{
				if i==h.canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:h.name+"_eval",SurfaceIndices:p},ph)
				sp:=h.means[i];np:=sp+shift;op:=sp+oracleShift
				se:=math.Abs(sp-pt.RequiredMassFactor);ne:=math.Abs(np-pt.RequiredMassFactor);oe:=math.Abs(op-pt.RequiredMassFactor)
				res.Metrics=append(res.Metrics,UP209BMetric{HeldoutFamily:h.name,Permutation:append([]int(nil),p...),Phase:ph,ActualRequiredFactor:pt.RequiredMassFactor,StaticPredictedFactor:sp,NativeDeltaPredictedFactor:np,OracleAnchorPredictedFactor:op,StaticAbsoluteError:se,NativeDeltaAbsoluteError:ne,OracleAnchorAbsoluteError:oe})
				s.EvaluationPoints++;sumS+=se;sumN+=ne;sumO+=oe
				if se>s.StaticMaxAbsoluteError{s.StaticMaxAbsoluteError=se};if ne>s.NativeDeltaMaxAbsoluteError{s.NativeDeltaMaxAbsoluteError=ne};if oe>s.OracleAnchorMaxAbsoluteError{s.OracleAnchorMaxAbsoluteError=oe}
				if ne<se{s.NativeBetterThanStatic++}
			}
		}
		s.StaticMeanAbsoluteError=sumS/float64(s.EvaluationPoints);s.NativeDeltaMeanAbsoluteError=sumN/float64(s.EvaluationPoints);s.OracleAnchorMeanAbsoluteError=sumO/float64(s.EvaluationPoints)
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
