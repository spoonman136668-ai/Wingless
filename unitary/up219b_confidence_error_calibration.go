package unitary

import "math"

const UP219BConfidenceErrorSchema="wingless.up219b-confidence-error-calibration.v1"

type UP219BDecision struct{
	Regime string `json:"regime"`
	Phase int `json:"phase"`
	ConfidenceMargin float64 `json:"confidence_margin"`
	ConfidenceStratum string `json:"confidence_stratum"`
	RouteCorrect bool `json:"route_correct"`
	RoutedMeanAbsoluteError float64 `json:"routed_mean_absolute_error"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	OracleFamilyMeanAbsoluteError float64 `json:"oracle_family_mean_absolute_error"`
}
type UP219BSummary struct{
	Scope string `json:"scope"`
	Regime string `json:"regime,omitempty"`
	ConfidenceStratum string `json:"confidence_stratum"`
	Decisions int `json:"decisions"`
	Misroutes int `json:"misroutes"`
	MisrouteRate float64 `json:"misroute_rate"`
	MeanRoutedMAE float64 `json:"mean_routed_mae"`
	MeanStaticMAE float64 `json:"mean_static_mae"`
	MeanOracleFamilyMAE float64 `json:"mean_oracle_family_mae"`
}
type UP219BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP218BSeal string `json:"source_up218b_seal"`
	TrainingDerivedThreshold float64 `json:"training_derived_threshold"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	RouteDecisions int `json:"route_decisions"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	RetrainingUsed bool `json:"retraining_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Decisions []UP219BDecision `json:"decisions"`
	Summaries []UP219BSummary `json:"summaries"`
}
func RunUP219B()(UP219BResult,error){
	template:=[]int{55,56,57}
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{232,233,234,235,236,237,238,239,240,241,242,243,244,245,246,247,248,249,250,251,252,253,254,255}
	specs:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type raw struct{name string;x [4]float64};raws:=[]raw{}
	for _,f:=range specs{for _,ph:=range train{raws=append(raws,raw{name:f.Name,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range raws{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(raws))}
	for _,r:=range raws{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(raws)));if std[j]==0{std[j]=1}}
	fams:=[]up213bFamily{}
	for _,f:=range specs{
		perms:=up196bPermutations(f.Canonical);means:=make([]float64,len(perms));canon:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(f.Canonical){canon=i}
			s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=s/float64(len(template))
		}
		ff:=up213bFamily{name:f.Name,canonical:f.Canonical,perms:perms,means:means,canon:canon};n:=0
		for _,r:=range raws{if r.name==f.Name{for j:=0;j<4;j++{ff.centroid[j]+=(r.x[j]-mean[j])/std[j]};n++}}
		for j:=0;j<4;j++{ff.centroid[j]/=float64(n)}
		var a [5][5]float64;var b [5]float64
		for _,ph:=range train{
			x:=up193bFeature(ph,f.Canonical);v:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph);y:=pt.RequiredMassFactor-means[canon]
			for i:=0;i<5;i++{b[i]+=v[i]*y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
		}
		for i:=1;i<5;i++{a[i][i]+=1e-6};c,ok:=up193bSolve(a,b);if !ok{return UP219BResult{},nil};ff.coef=c;fams=append(fams,ff)
	}
	threshold:=math.Inf(1)
	for fi,f:=range fams{for _,ph:=range train{
		x:=up193bFeature(ph,f.canonical);z:=[4]float64{};for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range fams{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		if bestIdx==fi{m:=second-best;if m<threshold{threshold=m}}
	}}
	if math.IsInf(threshold,1){threshold=0}
	res:=UP219BResult{Schema:UP219BConfidenceErrorSchema,Experiment:"UP-219B-confidence-error-calibration",SourceUP218BSeal:"94aaf811a0e2ab0bdb62a197691e028360bac08d",TrainingDerivedThreshold:threshold,TrainingPhases:train,EvaluationPhases:eval,EvaluationDerivedThresholdUsed:false,RetrainingUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearModelUsed:false}
	for fi,f:=range fams{for _,ph:=range eval{
		x:=up193bFeature(ph,f.canonical);z4:=[4]float64{};v:=[5]float64{1}
		for j:=0;j<4;j++{z4[j]=(x[j]-mean[j])/std[j];v[j+1]=z4[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range fams{d:=up212bDist(z4,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		margin:=second-best;stratum:="confident";if margin<threshold{stratum="uncertain"}
		routedShift,oracleShift:=0.0,0.0
		for i:=0;i<5;i++{routedShift+=fams[bestIdx].coef[i]*v[i];oracleShift+=f.coef[i]*v[i]}
		sumR,sumS,sumO:=0.0,0.0,0.0;n:=0
		for i,p:=range f.perms{
			if i==f.canon{continue}
			pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:p},ph)
			sp:=f.means[i]
			sumS+=math.Abs(sp-pt.RequiredMassFactor)
			sumR+=math.Abs(sp+routedShift-pt.RequiredMassFactor)
			sumO+=math.Abs(sp+oracleShift-pt.RequiredMassFactor);n++
		}
		res.Decisions=append(res.Decisions,UP219BDecision{Regime:f.name,Phase:ph,ConfidenceMargin:margin,ConfidenceStratum:stratum,RouteCorrect:bestIdx==fi,RoutedMeanAbsoluteError:sumR/float64(n),StaticMeanAbsoluteError:sumS/float64(n),OracleFamilyMeanAbsoluteError:sumO/float64(n)})
		res.RouteDecisions++
	}}
	scopes:=[]struct{scope,regime,stratum string}{
		{"overall","","uncertain"},{"overall","","confident"},
		{"regime","mixed4","uncertain"},{"regime","mixed4","confident"},
		{"regime","observe4","uncertain"},{"regime","observe4","confident"},
		{"regime","store4","uncertain"},{"regime","store4","confident"},
		{"regime","cross3","uncertain"},{"regime","cross3","confident"},
	}
	for _,q:=range scopes{
		s:=UP219BSummary{Scope:q.scope,Regime:q.regime,ConfidenceStratum:q.stratum};sr,ss,so:=0.0,0.0,0.0
		for _,d:=range res.Decisions{
			if d.ConfidenceStratum!=q.stratum{continue};if q.regime!=""&&d.Regime!=q.regime{continue}
			s.Decisions++;if !d.RouteCorrect{s.Misroutes++};sr+=d.RoutedMeanAbsoluteError;ss+=d.StaticMeanAbsoluteError;so+=d.OracleFamilyMeanAbsoluteError
		}
		if s.Decisions>0{s.MisrouteRate=float64(s.Misroutes)/float64(s.Decisions);s.MeanRoutedMAE=sr/float64(s.Decisions);s.MeanStaticMAE=ss/float64(s.Decisions);s.MeanOracleFamilyMAE=so/float64(s.Decisions)}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
