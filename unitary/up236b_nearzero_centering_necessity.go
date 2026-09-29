package unitary
import "math"
const UP236BCenteringSchema="wingless.up236b-nearzero-centering-necessity.v1"
type UP236BFamilySummary struct{Family string `json:"family"`;States int `json:"states"`;RawCorrect int `json:"raw_correct"`;CenteredCorrect int `json:"centered_correct"`}
type UP236BResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUP235BSeal string `json:"source_up235b_seal"`;ParentMinimumPerfectCardinality int `json:"parent_minimum_perfect_cardinality"`;ParentMinimumPerfectMasks []int `json:"parent_minimum_perfect_masks"`;Feature string `json:"feature"`;RawAccuracy float64 `json:"raw_accuracy"`;CenteredAccuracy float64 `json:"centered_accuracy"`;TrainingParityLabelsUsed bool `json:"training_parity_labels_used"`;PhaseInputAtInferenceUsed bool `json:"phase_input_at_inference_used"`;EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`;AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`;NonlinearClassifierUsed bool `json:"nonlinear_classifier_used"`;LiveActivation bool `json:"live_activation"`;Classification string `json:"classification"`;Summaries []UP236BFamilySummary `json:"summaries"`}
func RunUP236B()(UP236BResult,error){
 p,err:=RunUP235B();if err!=nil{return UP236BResult{},err}
 train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72};eval:=[]int{121,122,123,124,125,126,127,128,129,130,131,132,133,134,135,136,137,138,139,140,141,142,143,144,145,146,147,148,149,150,151,152}
 specs:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
 parity:=func(ph int)string{if ph%2==0{return "even"};return "odd"};feature:=func(ph int,f UP197BComposition)float64{return up193bFeature(ph,f.Canonical)[2]}
 type row struct{name string;phase int;x float64};rows:=[]row{};familySum:=map[string]float64{};familyN:=map[string]int{};rawMean:=0.0
 for _,f:=range specs{for _,ph:=range train{x:=feature(ph,f);rows=append(rows,row{f.Name,ph,x});familySum[f.Name]+=x;familyN[f.Name]++;rawMean+=x}}
 rawMean/=float64(len(rows));familyMean:=map[string]float64{};for _,f:=range specs{familyMean[f.Name]=familySum[f.Name]/float64(familyN[f.Name])}
 rawStd,ctrStd:=0.0,0.0;for _,r:=range rows{d:=r.x-rawMean;rawStd+=d*d;c:=r.x-familyMean[r.name];ctrStd+=c*c};rawStd=math.Sqrt(rawStd/float64(len(rows)));ctrStd=math.Sqrt(ctrStd/float64(len(rows)));if rawStd==0{rawStd=1};if ctrStd==0{ctrStd=1}
 rawSum:=map[string]float64{"even":0,"odd":0};ctrSum:=map[string]float64{"even":0,"odd":0};n:=map[string]int{"even":0,"odd":0}
 for _,r:=range rows{p:=parity(r.phase);rawSum[p]+=(r.x-rawMean)/rawStd;ctrSum[p]+=(r.x-familyMean[r.name])/ctrStd;n[p]++}
 rawCent:=map[string]float64{"even":rawSum["even"]/float64(n["even"]),"odd":rawSum["odd"]/float64(n["odd"])};ctrCent:=map[string]float64{"even":ctrSum["even"]/float64(n["even"]),"odd":ctrSum["odd"]/float64(n["odd"])}
 r:=UP236BResult{Schema:UP236BCenteringSchema,Experiment:"UP-236B-nearzero-centering-necessity",SourceUP235BSeal:"406089e7f12165a52f623c6c3148bc9862c04529",ParentMinimumPerfectCardinality:p.MinimumPerfectCardinality,ParentMinimumPerfectMasks:append([]int{},p.MinimumPerfectMasks...),Feature:"near_zero_margin_count",TrainingParityLabelsUsed:true}
 rawCorrect,ctrCorrect:=0,0
 for _,f:=range specs{sm:=UP236BFamilySummary{Family:f.Name};for _,ph:=range eval{x:=feature(ph,f);rz:=(x-rawMean)/rawStd;cz:=(x-familyMean[f.Name])/ctrStd;rb:="even";if math.Abs(rz-rawCent["odd"])<math.Abs(rz-rawCent["even"]){rb="odd"};cb:="even";if math.Abs(cz-ctrCent["odd"])<math.Abs(cz-ctrCent["even"]){cb="odd"};want:=parity(ph);sm.States++;if rb==want{sm.RawCorrect++;rawCorrect++};if cb==want{sm.CenteredCorrect++;ctrCorrect++}};r.Summaries=append(r.Summaries,sm)}
 r.RawAccuracy=float64(rawCorrect)/128.0;r.CenteredAccuracy=float64(ctrCorrect)/128.0
 if r.CenteredAccuracy<1.0{r.Classification="ANCHOR_NOT_REPRODUCED"}else if r.RawAccuracy==1.0{r.Classification="RAW_SINGLE_FEATURE_SUFFICIENT"}else if r.RawAccuracy<1.0{r.Classification="CENTERING_REQUIRED"}else{r.Classification="OTHER_VALID_PATTERN"}
 return r,nil
}
