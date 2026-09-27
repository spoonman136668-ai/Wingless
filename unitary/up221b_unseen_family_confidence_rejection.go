package unitary

import "math"

const UP221BUnseenSchema="wingless.up221b-unseen-family-confidence-rejection.v1"

type UP221BDecision struct{
	Family string `json:"family"`
	Phase int `json:"phase"`
	RoutedRegime string `json:"routed_regime"`
	ConfidenceMargin float64 `json:"confidence_margin"`
	Rejected bool `json:"rejected"`
	HarmfulForcedRoute bool `json:"harmful_forced_route"`
	RoutedMAE float64 `json:"routed_mae"`
	StaticMAE float64 `json:"static_mae"`
	OracleAnchorMAE float64 `json:"oracle_anchor_mae"`
}
type UP221BSummary struct{
	Scope string `json:"scope"`
	Family string `json:"family,omitempty"`
	Decisions int `json:"decisions"`
	Rejected int `json:"rejected"`
	Accepted int `json:"accepted"`
	HarmfulForcedRoutes int `json:"harmful_forced_routes"`
	HarmfulRejected int `json:"harmful_rejected"`
	HarmfulAccepted int `json:"harmful_accepted"`
	BenignRejected int `json:"benign_rejected"`
	RejectionRate float64 `json:"rejection_rate"`
	HarmfulRouteRejectionRate float64 `json:"harmful_route_rejection_rate"`
	RejectedMeanRoutedMAE float64 `json:"rejected_mean_routed_mae"`
	RejectedMeanStaticMAE float64 `json:"rejected_mean_static_mae"`
	AcceptedMeanRoutedMAE float64 `json:"accepted_mean_routed_mae"`
	AcceptedMeanStaticMAE float64 `json:"accepted_mean_static_mae"`
}
type UP221BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP220BSeal string `json:"source_up220b_seal"`
	TrainingConfidenceThreshold float64 `json:"training_confidence_threshold"`
	KnownRegimes []string `json:"known_regimes"`
	UnseenFamilies []string `json:"unseen_families"`
	EvaluationPhases []int `json:"evaluation_phases"`
	UnseenFamilyFittingUsed bool `json:"unseen_family_fitting_used"`
	NewCentroidUsed bool `json:"new_centroid_used"`
	NewCalibratorUsed bool `json:"new_calibrator_used"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Decisions []UP221BDecision `json:"decisions"`
	Summaries []UP221BSummary `json:"summaries"`
}
func RunUP221B()(UP221BResult,error){
	template:=[]int{55,56,57}
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{288,289,290,291,292,293,294,295,296,297,298,299,300,301,302,303,304,305,306,307,308,309,310,311,312,313,314,315,316,317,318,319}
	known:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
	unseen:=[]UP197BComposition{{Name:"mixed4b",Canonical:[]int{2,7,3,8}},{Name:"observe3",Canonical:[]int{5,6,8}},{Name:"store3",Canonical:[]int{0,2,3}},{Name:"cross4",Canonical:[]int{0,5,8,13}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type raw struct{name string;x [4]float64};raws:=[]raw{}
	for _,f:=range known{for _,ph:=range train{raws=append(raws,raw{name:f.Name,x:up193bFeature(ph,f.Canonical)})}}
	var mean,std [4]float64
	for _,r:=range raws{for j:=0;j<4;j++{mean[j]+=r.x[j]}}
	for j:=0;j<4;j++{mean[j]/=float64(len(raws))}
	for _,r:=range raws{for j:=0;j<4;j++{d:=r.x[j]-mean[j];std[j]+=d*d}}
	for j:=0;j<4;j++{std[j]=math.Sqrt(std[j]/float64(len(raws)));if std[j]==0{std[j]=1}}
	type knownModel struct{name string;canonical []int;centroid [4]float64;coef [5]float64}
	kms:=[]knownModel{}
	for _,f:=range known{
		km:=knownModel{name:f.Name,canonical:f.Canonical};n:=0
		for _,r:=range raws{if r.name==f.Name{for j:=0;j<4;j++{km.centroid[j]+=(r.x[j]-mean[j])/std[j]};n++}}
		for j:=0;j<4;j++{km.centroid[j]/=float64(n)}
		perms:=up196bPermutations(f.Canonical);canon:=-1;means:=make([]float64,len(perms))
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(f.Canonical){canon=i}
			s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=s/float64(len(template))
		}
		var a [5][5]float64;var b [5]float64
		for _,ph:=range train{
			x:=up193bFeature(ph,f.Canonical);v:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
			pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_train",SurfaceIndices:f.Canonical},ph);y:=pt.RequiredMassFactor-means[canon]
			for i:=0;i<5;i++{b[i]+=v[i]*y;for j:=0;j<5;j++{a[i][j]+=v[i]*v[j]}}
		}
		for i:=1;i<5;i++{a[i][i]+=1e-6};c,ok:=up193bSolve(a,b);if !ok{return UP221BResult{},nil};km.coef=c;kms=append(kms,km)
	}
	threshold:=math.Inf(1)
	for fi,f:=range known{for _,ph:=range train{
		x:=up193bFeature(ph,f.Canonical);z:=[4]float64{};for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range kms{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		if bestIdx==fi{m:=second-best;if m<threshold{threshold=m}}
	}}
	if math.IsInf(threshold,1){threshold=0}
	res:=UP221BResult{Schema:UP221BUnseenSchema,Experiment:"UP-221B-unseen-family-confidence-rejection",SourceUP220BSeal:"9cd5ec9eb8bd3b6da3790887b603c4427f0edb8b",TrainingConfidenceThreshold:threshold,KnownRegimes:[]string{"mixed4","observe4","store4","cross3"},UnseenFamilies:[]string{"mixed4b","observe3","store3","cross4"},EvaluationPhases:eval,UnseenFamilyFittingUsed:false,NewCentroidUsed:false,NewCalibratorUsed:false,EvaluationDerivedThresholdUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearModelUsed:false}
	for _,f:=range unseen{
		perms:=up196bPermutations(f.Canonical);means:=make([]float64,len(perms));canon:=-1
		for i,p:=range perms{
			if up201bKey(p)==up201bKey(f.Canonical){canon=i}
			s:=0.0;for _,ph:=range template{s+=up191bPoint(model,UP191BProfile{Name:f.Name+"_template",SurfaceIndices:p},ph).RequiredMassFactor}
			means[i]=s/float64(len(template))
		}
		for _,ph:=range eval{
			x:=up193bFeature(ph,f.Canonical);z4:=[4]float64{};v:=[5]float64{1}
			for j:=0;j<4;j++{z4[j]=(x[j]-mean[j])/std[j];v[j+1]=z4[j]}
			bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
			for gi,g:=range kms{d:=up212bDist(z4,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
			margin:=second-best;shift:=0.0;for i:=0;i<5;i++{shift+=kms[bestIdx].coef[i]*v[i]}
			canonPt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_eval",SurfaceIndices:f.Canonical},ph)
			anchorShift:=canonPt.RequiredMassFactor-means[canon]
			sumR,sumS,sumA:=0.0,0.0,0.0;n:=0
			for i,p:=range perms{
				if i==canon{continue}
				pt:=up191bPoint(model,UP191BProfile{Name:f.Name+"_eval",SurfaceIndices:p},ph);sp:=means[i]
				sumR+=math.Abs(sp+shift-pt.RequiredMassFactor);sumS+=math.Abs(sp-pt.RequiredMassFactor);sumA+=math.Abs(sp+anchorShift-pt.RequiredMassFactor);n++
			}
			rmae,smae,amae:=sumR/float64(n),sumS/float64(n),sumA/float64(n)
			res.Decisions=append(res.Decisions,UP221BDecision{Family:f.Name,Phase:ph,RoutedRegime:kms[bestIdx].name,ConfidenceMargin:margin,Rejected:margin<threshold,HarmfulForcedRoute:rmae>smae,RoutedMAE:rmae,StaticMAE:smae,OracleAnchorMAE:amae})
		}
	}
	families:=[]string{"mixed4b","observe3","store3","cross4",""}
	for _,fam:=range families{
		scope:="family";if fam==""{scope="overall"}
		s:=UP221BSummary{Scope:scope,Family:fam};rr,rs,ar,as:=0.0,0.0,0.0,0.0;rn,an:=0,0
		for _,d:=range res.Decisions{
			if fam!=""&&d.Family!=fam{continue}
			s.Decisions++;if d.Rejected{s.Rejected++;rr+=d.RoutedMAE;rs+=d.StaticMAE;rn++}else{s.Accepted++;ar+=d.RoutedMAE;as+=d.StaticMAE;an++}
			if d.HarmfulForcedRoute{s.HarmfulForcedRoutes++;if d.Rejected{s.HarmfulRejected++}else{s.HarmfulAccepted++}}else if d.Rejected{s.BenignRejected++}
		}
		if s.Decisions>0{s.RejectionRate=float64(s.Rejected)/float64(s.Decisions)}
		if s.HarmfulForcedRoutes>0{s.HarmfulRouteRejectionRate=float64(s.HarmfulRejected)/float64(s.HarmfulForcedRoutes)}
		if rn>0{s.RejectedMeanRoutedMAE=rr/float64(rn);s.RejectedMeanStaticMAE=rs/float64(rn)}
		if an>0{s.AcceptedMeanRoutedMAE=ar/float64(an);s.AcceptedMeanStaticMAE=as/float64(an)}
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
