package unitary
import "math"
const UP237BFarParitySchema="wingless.up237b-far-centered-nearzero-parity.v1"
type UP237BFamilySummary struct{Family string `json:"family"`;States int `json:"states"`;Correct int `json:"correct"`}
type UP237BResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUP236BSeal string `json:"source_up236b_seal"`;ParentClassification string `json:"parent_classification"`;ParentCenteredAccuracy float64 `json:"parent_centered_accuracy"`;TrainingPhases []int `json:"training_phases"`;FarEvaluationPhases []int `json:"far_evaluation_phases"`;Feature string `json:"feature"`;FarAccuracy float64 `json:"far_accuracy"`;Classification string `json:"classification"`;EvaluationLabelFittingUsed bool `json:"evaluation_label_fitting_used"`;PhaseInputAtInferenceUsed bool `json:"phase_input_at_inference_used"`;AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`;NonlinearClassifierUsed bool `json:"nonlinear_classifier_used"`;LiveActivation bool `json:"live_activation"`;Summaries []UP237BFamilySummary `json:"summaries"`}
func RunUP237B()(UP237BResult,error){
 p,err:=RunUP236B();if err!=nil{return UP237BResult{},err}
 train:=[]int{58,59,60,61,62,63,64,65,66,67,68,69,70,71,72};far:=[]int{185,186,187,188,189,190,191,192,193,194,195,196,197,198,199,200,201,202,203,204,205,206,207,208,209,210,211,212,213,214,215,216}
 specs:=[]UP197BComposition{{Name:"mixed4",Canonical:[]int{0,5,1,6}},{Name:"observe4",Canonical:[]int{5,6,7,8}},{Name:"store4",Canonical:[]int{0,1,2,3}},{Name:"cross3",Canonical:[]int{0,5,13}}}
 feat:=func(ph int,f UP197BComposition)float64{return up193bFeature(ph,f.Canonical)[2]};par:=func(ph int)string{if ph%2==0{return "even"};return "odd"}
 famSum:=map[string]float64{};famN:=map[string]int{};type row struct{name string;ph int;x float64};rows:=[]row{}
 for _,f:=range specs{for _,ph:=range train{x:=feat(ph,f);rows=append(rows,row{f.Name,ph,x});famSum[f.Name]+=x;famN[f.Name]++}}
 famMean:=map[string]float64{};for _,f:=range specs{famMean[f.Name]=famSum[f.Name]/float64(famN[f.Name])}
 std:=0.0;for _,x:=range rows{d:=x.x-famMean[x.name];std+=d*d};std=math.Sqrt(std/float64(len(rows)));if std==0{std=1}
 sum:=map[string]float64{"even":0,"odd":0};n:=map[string]int{"even":0,"odd":0}
 for _,x:=range rows{q:=par(x.ph);sum[q]+=(x.x-famMean[x.name])/std;n[q]++}
 cent:=map[string]float64{"even":sum["even"]/float64(n["even"]),"odd":sum["odd"]/float64(n["odd"])}
 r:=UP237BResult{Schema:UP237BFarParitySchema,Experiment:"UP-237B-far-centered-nearzero-parity",SourceUP236BSeal:"a76db34ca7373b5df22e72069d3e9b9a7ae7c6d1",ParentClassification:p.Classification,ParentCenteredAccuracy:p.CenteredAccuracy,TrainingPhases:train,FarEvaluationPhases:far,Feature:"near_zero_margin_count"}
 correct:=0
 for _,f:=range specs{sm:=UP237BFamilySummary{Family:f.Name};for _,ph:=range far{z:=(feat(ph,f)-famMean[f.Name])/std;pred:="even";if math.Abs(z-cent["odd"])<math.Abs(z-cent["even"]){pred="odd"};sm.States++;if pred==par(ph){sm.Correct++;correct++}};r.Summaries=append(r.Summaries,sm)}
 r.FarAccuracy=float64(correct)/128.0
 if p.Classification!="CENTERING_REQUIRED"||p.CenteredAccuracy!=1{r.Classification="ANCHOR_NOT_REPRODUCED"}else if r.FarAccuracy==1{r.Classification="FAR_TRANSFER_PERFECT"}else{r.Classification="FAR_TRANSFER_DECAY"}
 return r,nil
}
