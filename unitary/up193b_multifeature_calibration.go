package unitary

import "math"

const UP193BMultiFeatureSchema="wingless.up193b-multifeature-calibration.v1"

type UP193BProfile struct{
	Name string `json:"name"`
	SurfaceIndices []int `json:"surface_indices"`
}
type UP193BMetric struct{
	Profile string `json:"profile"`
	Phase int `json:"phase"`
	NativeCorrectCount int `json:"native_correct_count"`
	ActualRequiredFactor float64 `json:"actual_required_factor"`
	PredictedFactor float64 `json:"predicted_factor"`
	PooledAbsoluteMassRatioError float64 `json:"pooled_absolute_mass_ratio_error"`
	CorrectedAbsoluteMassRatioError float64 `json:"corrected_absolute_mass_ratio_error"`
}
type UP193BSummary struct{
	Model string `json:"model"`
	MeanAbsoluteMassRatioError float64 `json:"mean_absolute_mass_ratio_error"`
	MaxAbsoluteMassRatioError float64 `json:"max_absolute_mass_ratio_error"`
}
type UP193BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP192BSeal string `json:"source_up192b_seal"`
	TrainingPhases []int `json:"training_phases"`
	EvaluationPhases []int `json:"evaluation_phases"`
	TrainingProfiles []UP193BProfile `json:"training_profiles"`
	EvaluationProfiles []UP193BProfile `json:"evaluation_profiles"`
	Features []string `json:"features"`
	RidgeLambda float64 `json:"ridge_lambda"`
	FeatureMeans []float64 `json:"feature_means"`
	FeatureStdDevs []float64 `json:"feature_std_devs"`
	Coefficients []float64 `json:"coefficients"`
	TrainingPoints int `json:"training_points"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	AdaptiveProfileSearchUsed bool `json:"adaptive_profile_search_used"`
	PhaseInputUsed bool `json:"phase_input_used"`
	ParityInputUsed bool `json:"parity_input_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Metrics []UP193BMetric `json:"metrics"`
	Summaries []UP193BSummary `json:"summaries"`
}
func up193bFeature(phase int,indices []int)[4]float64{
	c,meanAbs,nz,minAbs:=up192bGeometry(phase,indices)
	return [4]float64{float64(c),meanAbs,float64(nz),minAbs}
}
func up193bSolve(a [5][5]float64,b [5]float64)([5]float64,bool){
	for i:=0;i<5;i++{
		p:=i;best:=math.Abs(a[i][i]);for r:=i+1;r<5;r++{if v:=math.Abs(a[r][i]);v>best{best=v;p=r}}
		if best<1e-12{return [5]float64{},false}
		if p!=i{a[i],a[p]=a[p],a[i];b[i],b[p]=b[p],b[i]}
		d:=a[i][i];for c:=i;c<5;c++{a[i][c]/=d};b[i]/=d
		for r:=0;r<5;r++{if r==i{continue};f:=a[r][i];for c:=i;c<5;c++{a[r][c]-=f*a[i][c]};b[r]-=f*b[i]}
	}
	return b,true
}
func up193bRows(phase int,indices []int)[]up182bRow{if len(indices)==0{return up182bRows(phase)};return up190bRows(phase,indices)}
func RunUP193B()(UP193BResult,error){
	trainPhases:=[]int{43,44,45};evalPhases:=[]int{46,47,48}
	trainProfiles:=[]UP193BProfile{{Name:"control"},{Name:"store1",SurfaceIndices:[]int{0}},{Name:"store2",SurfaceIndices:[]int{0,1}},{Name:"store4",SurfaceIndices:[]int{0,1,2,3}},{Name:"observe1",SurfaceIndices:[]int{5}},{Name:"observe2",SurfaceIndices:[]int{5,6}},{Name:"observe4",SurfaceIndices:[]int{5,6,7,8}},{Name:"mixed2",SurfaceIndices:[]int{0,5}},{Name:"mixed4",SurfaceIndices:[]int{0,5,1,6}}}
	evalProfiles:=[]UP193BProfile{{Name:"store3",SurfaceIndices:[]int{0,1,2}},{Name:"observe3",SurfaceIndices:[]int{5,6,7}},{Name:"mixed3",SurfaceIndices:[]int{0,5,1}},{Name:"report2",SurfaceIndices:[]int{13,14}},{Name:"cross3",SurfaceIndices:[]int{0,5,13}}}
	pooled:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	type sample struct{x [4]float64;y float64}
	samples:=[]sample{}
	for _,p:=range trainProfiles{for _,phase:=range trainPhases{
		pp:=UP191BProfile{Name:p.Name,SurfaceIndices:p.SurfaceIndices};pt:=up191bPoint(pooled,pp,phase)
		samples=append(samples,sample{x:up193bFeature(phase,p.SurfaceIndices),y:pt.RequiredMassFactor})
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
	coef,ok:=up193bSolve(a,b);if !ok{return UP193BResult{},nil}
	res:=UP193BResult{Schema:UP193BMultiFeatureSchema,Experiment:"UP-193B-multifeature-calibration",SourceUP192BSeal:"f60d940777f04fdd5aee86b2f8483a0d7faa8f6c",TrainingPhases:trainPhases,EvaluationPhases:evalPhases,TrainingProfiles:trainProfiles,EvaluationProfiles:evalProfiles,Features:[]string{"native_correct_count","mean_absolute_margin","near_zero_margin_count","min_absolute_margin"},RidgeLambda:1e-6,FeatureMeans:[]float64{mean[0],mean[1],mean[2],mean[3]},FeatureStdDevs:[]float64{std[0],std[1],std[2],std[3]},Coefficients:[]float64{coef[0],coef[1],coef[2],coef[3],coef[4]},TrainingPoints:len(samples),HeldoutFittingUsed:false,AdaptiveFeatureSelectionUsed:false,AdaptiveProfileSearchUsed:false,PhaseInputUsed:false,ParityInputUsed:false,MaintenanceTriggered:false}
	sumPool,sumCorr,maxPool,maxCorr:=0.0,0.0,0.0,0.0;n:=0
	for _,p:=range evalProfiles{for _,phase:=range evalPhases{
		x:=up193bFeature(phase,p.SurfaceIndices);z:=[5]float64{1,(x[0]-mean[0])/std[0],(x[1]-mean[1])/std[1],(x[2]-mean[2])/std[2],(x[3]-mean[3])/std[3]}
		f:=0.0;for i:=0;i<5;i++{f+=coef[i]*z[i]}
		rows:=up193bRows(phase,p.SurfaceIndices);actual:=0;predPool,predCorr:=0.0,0.0
		for _,row:=range rows{q:=up188bP(pooled,row);predPool+=q;qc:=q*f;if qc<0{qc=0};if qc>1{qc=1};predCorr+=qc;if row.positive{actual++}}
		actualFactor:=0.0;poolErr,corrErr:=0.0,0.0
		if predPool>0{actualFactor=float64(actual)/predPool};if actual>0{poolErr=math.Abs(predPool/float64(actual)-1);corrErr=math.Abs(predCorr/float64(actual)-1)}
		res.Metrics=append(res.Metrics,UP193BMetric{Profile:p.Name,Phase:phase,NativeCorrectCount:int(x[0]),ActualRequiredFactor:actualFactor,PredictedFactor:f,PooledAbsoluteMassRatioError:poolErr,CorrectedAbsoluteMassRatioError:corrErr})
		sumPool+=poolErr;sumCorr+=corrErr;if poolErr>maxPool{maxPool=poolErr};if corrErr>maxCorr{maxCorr=corrErr};n++
	}}
	res.Summaries=[]UP193BSummary{{Model:"pooled_26_30",MeanAbsoluteMassRatioError:sumPool/float64(n),MaxAbsoluteMassRatioError:maxPool},{Model:"ridge_native_geometry",MeanAbsoluteMassRatioError:sumCorr/float64(n),MaxAbsoluteMassRatioError:maxCorr}}
	return res,nil
}
