package unitary

import "math"

const UP203BOrderSignalSchema="wingless.up203b-order-signal-persistence.v1"

type UP203BPhaseMetric struct{
	Composition string `json:"composition"`
	Phase int `json:"phase"`
	ResidualCorrelation float64 `json:"residual_correlation"`
	FactorSpread float64 `json:"factor_spread"`
}
type UP203BSummary struct{
	Composition string `json:"composition"`
	TemplateFactorSpread float64 `json:"template_factor_spread"`
	MeanResidualCorrelation float64 `json:"mean_residual_correlation"`
	MinResidualCorrelation float64 `json:"min_residual_correlation"`
	MaxResidualCorrelation float64 `json:"max_residual_correlation"`
	MeanHeldoutFactorSpread float64 `json:"mean_heldout_factor_spread"`
}
type UP203BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP202BSeal string `json:"source_up202b_seal"`
	TemplatePhases []int `json:"template_phases"`
	DiagnosticPhases []int `json:"diagnostic_phases"`
	PermutationCountPerComposition int `json:"permutation_count_per_composition"`
	HeldoutFittingUsed bool `json:"heldout_fitting_used"`
	FeatureSearchUsed bool `json:"feature_search_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	PhaseMetrics []UP203BPhaseMetric `json:"phase_metrics"`
	Summaries []UP203BSummary `json:"summaries"`
}
func up203bCorr(a,b []float64)float64{
	if len(a)!=len(b)||len(a)<2{return 0}
	ma,mb:=0.0,0.0
	for i:=range a{ma+=a[i];mb+=b[i]}
	ma/=float64(len(a));mb/=float64(len(b))
	num,da,db:=0.0,0.0,0.0
	for i:=range a{x:=a[i]-ma;y:=b[i]-mb;num+=x*y;da+=x*x;db+=y*y}
	if da==0||db==0{return 0}
	return num/math.Sqrt(da*db)
}
func up203bSpread(v []float64)float64{
	if len(v)==0{return 0}
	min,max:=v[0],v[0]
	for _,x:=range v[1:]{if x<min{min=x};if x>max{max=x}}
	return max-min
}
func RunUP203B()(UP203BResult,error){
	template:=[]int{55,56,57};diag:=[]int{76,77,78}
	comps:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}}}
	model:=up184bBuild("pooled_26_30",[]int{26,27,28,29,30})
	res:=UP203BResult{Schema:UP203BOrderSignalSchema,Experiment:"UP-203B-order-signal-persistence",SourceUP202BSeal:"fc4846afbf855323c2c90e160c2bccf545a1661a",TemplatePhases:template,DiagnosticPhases:diag,PermutationCountPerComposition:24,HeldoutFittingUsed:false,FeatureSearchUsed:false,MaintenanceTriggered:false}
	for _,c:=range comps{
		perms:=up196bPermutations(c.Canonical)
		templateMeans:=make([]float64,len(perms));grand:=0.0
		for i,p:=range perms{
			s:=0.0
			for _,phase:=range template{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_template",SurfaceIndices:p},phase)
				s+=pt.RequiredMassFactor
			}
			templateMeans[i]=s/float64(len(template));grand+=templateMeans[i]
		}
		grand/=float64(len(templateMeans))
		templateResidual:=make([]float64,len(templateMeans))
		for i,v:=range templateMeans{templateResidual[i]=v-grand}
		s:=UP203BSummary{Composition:c.Name,TemplateFactorSpread:up203bSpread(templateMeans),MinResidualCorrelation:1,MaxResidualCorrelation:-1}
		corrSum,spreadSum:=0.0,0.0
		for _,phase:=range diag{
			vals:=make([]float64,len(perms));mean:=0.0
			for i,p:=range perms{
				pt:=up191bPoint(model,UP191BProfile{Name:c.Name+"_diag",SurfaceIndices:p},phase)
				vals[i]=pt.RequiredMassFactor;mean+=vals[i]
			}
			mean/=float64(len(vals))
			resid:=make([]float64,len(vals))
			for i,v:=range vals{resid[i]=v-mean}
			corr:=up203bCorr(templateResidual,resid);spread:=up203bSpread(vals)
			res.PhaseMetrics=append(res.PhaseMetrics,UP203BPhaseMetric{Composition:c.Name,Phase:phase,ResidualCorrelation:corr,FactorSpread:spread})
			corrSum+=corr;spreadSum+=spread
			if corr<s.MinResidualCorrelation{s.MinResidualCorrelation=corr}
			if corr>s.MaxResidualCorrelation{s.MaxResidualCorrelation=corr}
		}
		s.MeanResidualCorrelation=corrSum/float64(len(diag));s.MeanHeldoutFactorSpread=spreadSum/float64(len(diag))
		res.Summaries=append(res.Summaries,s)
	}
	return res,nil
}
