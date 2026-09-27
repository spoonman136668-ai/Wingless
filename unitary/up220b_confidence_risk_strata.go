package unitary

import (
	"math"
	"sort"
)

const UP220BConfidenceStrataSchema="wingless.up220b-confidence-risk-strata.v1"

type UP220BDecision struct{
	Regime string `json:"regime"`
	Phase int `json:"phase"`
	ConfidenceMargin float64 `json:"confidence_margin"`
	Stratum string `json:"stratum"`
	RouteCorrect bool `json:"route_correct"`
	RoutedMeanAbsoluteError float64 `json:"routed_mean_absolute_error"`
	StaticMeanAbsoluteError float64 `json:"static_mean_absolute_error"`
	OracleFamilyMeanAbsoluteError float64 `json:"oracle_family_mean_absolute_error"`
}
type UP220BSummary struct{
	Stratum string `json:"stratum"`
	Decisions int `json:"decisions"`
	Misroutes int `json:"misroutes"`
	MisrouteRate float64 `json:"misroute_rate"`
	MeanRoutedMAE float64 `json:"mean_routed_mae"`
	MeanStaticMAE float64 `json:"mean_static_mae"`
	MeanOracleFamilyMAE float64 `json:"mean_oracle_family_mae"`
}
type UP220BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP219BSeal string `json:"source_up219b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	TrainingCorrectMargins int `json:"training_correct_margins"`
	TrainingMin float64 `json:"training_min"`
	TrainingQ25 float64 `json:"training_q25"`
	TrainingQ50 float64 `json:"training_q50"`
	TrainingQ75 float64 `json:"training_q75"`
	Strata []string `json:"strata"`
	EvaluationDerivedBoundaryUsed bool `json:"evaluation_derived_boundary_used"`
	AdaptiveBinningUsed bool `json:"adaptive_binning_used"`
	RetrainingUsed bool `json:"retraining_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Decisions []UP220BDecision `json:"decisions"`
	Summaries []UP220BSummary `json:"summaries"`
}

func up220bStratum(m,min,q25,q50,q75 float64)string{
	if m<min{return "below_training_min"}
	if m<q25{return "training_low"}
	if m<q50{return "training_mid"}
	if m<q75{return "training_high"}
	return "training_top"
}

func RunUP220B()(UP220BResult,error){
	template:=[]int{55,56,57}
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{256,257,258,259,260,261,262,263,264,265,266,267,268,269,270,271,272,273,274,275,276,277,278,279,280,281,282,283,284,285,286,287}
	specs:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type raw struct{name string;x [4]float64}
	raws:=[]raw{}
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
		for i:=1;i<5;i++{a[i][i]+=1e-6};c,ok:=up193bSolve(a,b);if !ok{return UP220BResult{},nil};ff.coef=c
		fams=append(fams,ff)
	}
	margins:=[]float64{}
	for fi,f:=range fams{for _,ph:=range train{
		x:=up193bFeature(ph,f.canonical);z:=[4]float64{};for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range fams{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		if bestIdx==fi{margins=append(margins,second-best)}
	}}
	sort.Float64s(margins)
	if len(margins)==0{return UP220BResult{},nil}
	at:=func(p int)float64{return margins[(len(margins)-1)*p/100]}
	minv,q25,q50,q75:=margins[0],at(25),at(50),at(75)
	strata:=[]string{"below_training_min","training_low","training_mid","training_high","training_top"}
	res:=UP220BResult{Schema:UP220BConfidenceStrataSchema,Experiment:"UP-220B-confidence-risk-strata",SourceUP219BSeal:"f02efe61c0a43ad1f1c4fc2abc8f34e2d559bcd8",TrainingPhases:train,EvaluationPhases:eval,TrainingCorrectMargins:len(margins),TrainingMin:minv,TrainingQ25:q25,TrainingQ50:q50,TrainingQ75:q75,Strata:strata,EvaluationDerivedBoundaryUsed:false,AdaptiveBinningUsed:false,RetrainingUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearModelUsed:false}
	for fi,f:=range fams{for _,ph:=range eval{
		x:=up193bFeature(ph,f.canonical);z4:=[4]float64{};v:=[5]float64{1}
		for j:=0;j<4;j++{z4[j]=(x[j]-mean[j])/std[j];v[j+1]=z4[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range fams{d:=up212bDist(z4,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		margin:=second-best
		routedShift,oracleShift:=0.0,0.0
		for i:=0;i<5;i++{routedShift+=fams[bestIdx].coef[i]*v[i];oracleShift+=f.coef[i]*v[i]}
		sumR,sumS,sumO:=0.0,0.0,0.0;n:=0
		for i,p:=range f.perms{
			if i==f.canon{continue}
			pt:=up191bPoint(model,UP191BProfile{Name:f.name+"_eval",SurfaceIndices:p},ph);sp:=f.means[i]
			sumR+=math.Abs(sp+routedShift-pt.RequiredMassFactor);sumS+=math.Abs(sp-pt.RequiredMassFactor);sumO+=math.Abs(sp+oracleShift-pt.RequiredMassFactor);n++
		}
		res.Decisions=append(res.Decisions,UP220BDecision{Regime:f.name,Phase:ph,ConfidenceMargin:margin,Stratum:up220bStratum(margin,minv,q25,q50,q75),RouteCorrect:bestIdx==fi,RoutedMeanAbsoluteError:sumR/float64(n),StaticMeanAbsoluteError:sumS/float64(n),OracleFamilyMeanAbsoluteError:sumO/float64(n)})
	}}
	for _,stratum:=range strata{
		s:=UP220BSummary{Stratum:stratum};sr,ss,so:=0.0,0.0,0.0
		for _,d:=range res.Decisions{
			if d.Stratum!=stratum{continue}
			s.Decisions++;if !d.RouteCorrect{s.Misroutes++};sr+=d.RoutedMeanAbsoluteError;ss+=d.StaticMeanAbsoluteError;so+=d.OracleFamilyMeanAbsoluteError
		}
		if s.Decisions>0{s.MisrouteRate=float64(s.Misroutes)/float64(s.Decisions);s.MeanRoutedMAE=sr/float64(s.Decisions);s.MeanStaticMAE=ss/float64(s.Decisions);s.MeanOracleFamilyMAE=so/float64(s.Decisions)}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
