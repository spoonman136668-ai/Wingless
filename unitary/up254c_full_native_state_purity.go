package unitary
import("fmt";"sort")
const UP254CFullPuritySchema="wingless.up254c-full-native-state-purity.v1"
type UP254CResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUP253CSeal string `json:"source_up253c_seal"`;PairedInitialStates int `json:"paired_initial_states"`;ParentCoordinateOctuples int `json:"parent_coordinate_octuples"`;ParentFullyPureOctuples int `json:"parent_fully_pure_octuples"`;Coordinates []string `json:"coordinates"`;Groups int `json:"groups"`;PureGroups int `json:"pure_groups"`;MixedGroups int `json:"mixed_groups"`;MaxClassesPerGroup int `json:"max_classes_per_group"`;Classification string `json:"classification"`;ThresholdFittingUsed bool `json:"threshold_fitting_used"`;ClassifierTrainingUsed bool `json:"classifier_training_used"`;AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`;NewNativeFieldUsed bool `json:"new_native_field_used"`;LiveActivation bool `json:"live_activation"`}
func RunUP254C()(UP254CResult,error){
 p,err:=RunUP253C();if err!=nil{return UP254CResult{},err}
 cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}};hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};coords:=[]string{"endangered_age","predicted_is_endangered","hand_distance","age0","age1","age2","age3","adversarial_horizon","no_query_horizon"}
 type sample struct{values map[string]string;class string};samples:=[]sample{}
 r:=UP254CResult{Schema:UP254CFullPuritySchema,Experiment:"UP-254C-full-native-state-purity",SourceUP253CSeal:"8ffd93eab4b09801a86ea0f29855e6062fc23500",ParentCoordinateOctuples:p.CoordinateOctuples,ParentFullyPureOctuples:p.FullyPureOctuples,Coordinates:coords}
 for _,co:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,co);if !ok{continue};s:=up223cSnapshot(x,e,0);a,er:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if er!=nil{return UP254CResult{},er};b,er:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if er!=nil{return UP254CResult{},er};r.PairedInitialStates++;samples=append(samples,sample{up246cValues(s),up246cClass(a,b)})}}
 groups:=map[string]map[string]bool{}
 for _,s:=range samples{key:=fmt.Sprintf("%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s",coords[0],s.values[coords[0]],coords[1],s.values[coords[1]],coords[2],s.values[coords[2]],coords[3],s.values[coords[3]],coords[4],s.values[coords[4]],coords[5],s.values[coords[5]],coords[6],s.values[coords[6]],coords[7],s.values[coords[7]],coords[8],s.values[coords[8]]);if groups[key]==nil{groups[key]=map[string]bool{}};groups[key][s.class]=true}
 keys:=make([]string,0,len(groups));for k:=range groups{keys=append(keys,k)};sort.Strings(keys);r.Groups=len(groups)
 for _,k:=range keys{n:=len(groups[k]);if n==1{r.PureGroups++}else{r.MixedGroups++};if n>r.MaxClassesPerGroup{r.MaxClassesPerGroup=n}}
 if p.FullyPureOctuples!=0{r.Classification="ANCHOR_NOT_REPRODUCED"}else if r.MixedGroups==0{r.Classification="FULL_NATIVE_STATE_PURE"}else{r.Classification="FULL_NATIVE_STATE_NONSEPARABLE"}
 return r,nil
}
