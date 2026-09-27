package unitary

import "math"

const UP222BAbsoluteOODSchema="wingless.up222b-absolute-distance-ood-rejection.v1"

type UP222BDecision struct{
	Family string `json:"family"`
	Phase int `json:"phase"`
	RoutedRegime string `json:"routed_regime"`
	ConfidenceMargin float64 `json:"confidence_margin"`
	NearestCentroidDistance float64 `json:"nearest_centroid_distance"`
	MarginRejected bool `json:"margin_rejected"`
	DistanceRejected bool `json:"distance_rejected"`
	UnionRejected bool `json:"union_rejected"`
	HarmfulForcedRoute bool `json:"harmful_forced_route"`
	RoutedMAE float64 `json:"routed_mae"`
	StaticMAE float64 `json:"static_mae"`
	OracleAnchorMAE float64 `json:"oracle_anchor_mae"`
}
type UP222BSummary struct{
	Gate string `json:"gate"`
	Scope string `json:"scope"`
	Family string `json:"family,omitempty"`
	Decisions int `json:"decisions"`
	Rejected int `json:"rejected"`
	Accepted int `json:"accepted"`
	HarmfulForcedRoutes int `json:"harmful_forced_routes"`
	HarmfulRejected int `json:"harmful_rejected"`
	HarmfulAccepted int `json:"harmful_accepted"`
	BenignRejected int `json:"benign_rejected"`
	BenignRoutes int `json:"benign_routes"`
	RejectionRate float64 `json:"rejection_rate"`
	HarmfulRouteRejectionRate float64 `json:"harmful_route_rejection_rate"`
	BenignRejectionRate float64 `json:"benign_rejection_rate"`
}
type UP222BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP221BSeal string `json:"source_up221b_seal"`
	TrainingMarginThreshold float64 `json:"training_margin_threshold"`
	TrainingDistanceThreshold float64 `json:"training_distance_threshold"`
	TrainingCorrectStates int `json:"training_correct_states"`
	KnownRegimes []string `json:"known_regimes"`
	UnseenFamilies []string `json:"unseen_families"`
	EvaluationPhases []int `json:"evaluation_phases"`
	UnseenFamilyFittingUsed bool `json:"unseen_family_fitting_used"`
	NewCentroidUsed bool `json:"new_centroid_used"`
	NewCalibratorUsed bool `json:"new_calibrator_used"`
	EvaluationDerivedThresholdUsed bool `json:"evaluation_derived_threshold_used"`
	AdaptiveGateSelectionUsed bool `json:"adaptive_gate_selection_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	NonlinearModelUsed bool `json:"nonlinear_model_used"`
	Decisions []UP222BDecision `json:"decisions"`
	Summaries []UP222BSummary `json:"summaries"`
}
type up222bKnown struct{name string;canonical []int;centroid [4]float64;coef [5]float64}

func RunUP222B()(UP222BResult,error){
	template:=[]int{55,56,57}
	train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72}
	eval:=[]int{320,321,322,323,324,325,326,327,328,329,330,331,332,333,334,335,336,337,338,339,340,341,342,343,344,345,346,347,348,349,350,351}
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
	kms:=[]up222bKnown{}
	for _,f:=range known{
		km:=up222bKnown{name:f.Name,canonical:f.Canonical};n:=0
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
		for i:=1;i<5;i++{a[i][i]+=1e-6}
		c,ok:=up193bSolve(a,b);if !ok{return UP222BResult{},nil};km.coef=c;kms=append(kms,km)
	}
	marginThreshold:=math.Inf(1);distanceThreshold:=0.0;trainCorrect:=0
	for fi,f:=range known{for _,ph:=range train{
		x:=up193bFeature(ph,f.Canonical);z:=[4]float64{};for j:=0;j<4;j++{z[j]=(x[j]-mean[j])/std[j]}
		bestIdx:=-1;best,second:=math.Inf(1),math.Inf(1)
		for gi,g:=range kms{d:=up212bDist(z,g.centroid);if d<best{second=best;best=d;bestIdx=gi}else if d<second{second=d}}
		if bestIdx==fi{
			trainCorrect++;m:=second-best
			if m<marginThreshold{marginThreshold=m}
			if best>distanceThreshold{distanceThreshold=best}
		}
	}}
	if math.IsInf(marginThreshold,1){marginThreshold=0}
	res:=UP222BResult{
		Schema:UP222BAbsoluteOODSchema,Experiment:"UP-222B-absolute-distance-ood-rejection",
		SourceUP221BSeal:"4d9a0a4c40b2a65c6ada6a511b9e24c986e6aee9",
		TrainingMarginThreshold:marginThreshold,TrainingDistanceThreshold:distanceThreshold,TrainingCorrectStates:trainCorrect,
		KnownRegimes:[]string{"mixed4","observe4","store4","cross3"},
		UnseenFamilies:[]string{"mixed4b","observe3","store3","cross4"},EvaluationPhases:eval,
		UnseenFamilyFittingUsed:false,NewCentroidUsed:false,NewCalibratorUsed:false,
		EvaluationDerivedThresholdUsed:false,AdaptiveGateSelectionUsed:false,PhaseInputUsed:false,ParityInputUsed:false,NonlinearModelUsed:false,
	}
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
			mr:=margin<marginThreshold;dr:=best>distanceThreshold
			res.Decisions=append(res.Decisions,UP222BDecision{Family:f.Name,Phase:ph,RoutedRegime:kms[bestIdx].name,ConfidenceMargin:margin,NearestCentroidDistance:best,MarginRejected:mr,DistanceRejected:dr,UnionRejected:mr||dr,HarmfulForcedRoute:rmae>smae,RoutedMAE:rmae,StaticMAE:smae,OracleAnchorMAE:amae})
		}
	}
	gates:=[]string{"margin","distance","union"}
	families:=[]string{"mixed4b","observe3","store3","cross4",""}
	for _,gate:=range gates{for _,fam:=range families{
		scope:="family";if fam==""{scope="overall"}
		s:=UP222BSummary{Gate:gate,Scope:scope,Family:fam}
		for _,d:=range res.Decisions{
			if fam!=""&&d.Family!=fam{continue}
			reject:=d.MarginRejected;if gate=="distance"{reject=d.DistanceRejected};if gate=="union"{reject=d.UnionRejected}
			s.Decisions++;if reject{s.Rejected++}else{s.Accepted++}
			if d.HarmfulForcedRoute{
				s.HarmfulForcedRoutes++;if reject{s.HarmfulRejected++}else{s.HarmfulAccepted++}
			}else{
				s.BenignRoutes++;if reject{s.BenignRejected++}
			}
		}
		if s.Decisions>0{s.RejectionRate=float64(s.Rejected)/float64(s.Decisions)}
		if s.HarmfulForcedRoutes>0{s.HarmfulRouteRejectionRate=float64(s.HarmfulRejected)/float64(s.HarmfulForcedRoutes)}
		if s.BenignRoutes>0{s.BenignRejectionRate=float64(s.BenignRejected)/float64(s.BenignRoutes)}
		res.Summaries=append(res.Summaries,s)
	}}
	return res,nil
}
