package unitary
import("fmt";"sort")
const UP253CEightPuritySchema="wingless.up253c-eight-native-coordinate-purity.v1"
type UP253CEightSummary struct{Coordinates []string `json:"coordinates"`;Groups int `json:"groups"`;PureGroups int `json:"pure_groups"`;MixedGroups int `json:"mixed_groups"`;MaxClassesPerGroup int `json:"max_classes_per_group"`}
type UP253CResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUP252CSeal string `json:"source_up252c_seal"`;PairedInitialStates int `json:"paired_initial_states"`;ParentCoordinateSeptuples int `json:"parent_coordinate_septuples"`;ParentFullyPureSeptuples int `json:"parent_fully_pure_septuples"`;CoordinateOctuples int `json:"coordinate_octuples"`;FullyPureOctuples int `json:"fully_pure_octuples"`;InterventionChanged bool `json:"intervention_changed"`;ThresholdFittingUsed bool `json:"threshold_fitting_used"`;ClassifierTrainingUsed bool `json:"classifier_training_used"`;AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`;NewNativeFieldUsed bool `json:"new_native_field_used"`;LiveActivation bool `json:"live_activation"`;Summaries []UP253CEightSummary `json:"summaries"`}
func RunUP253C()(UP253CResult,error){
 p,err:=RunUP252C();if err!=nil{return UP253CResult{},err}
 cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}};hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15};coords:=[]string{"endangered_age","predicted_is_endangered","hand_distance","age0","age1","age2","age3","adversarial_horizon","no_query_horizon"}
 type sample struct{values map[string]string;class string};samples:=[]sample{}
 r:=UP253CResult{Schema:UP253CEightPuritySchema,Experiment:"UP-253C-eight-native-coordinate-purity",SourceUP252CSeal:"cd3e7cb56c9df5163238db647b9d9e7255f56138",ParentCoordinateSeptuples:p.CoordinateSeptuples,ParentFullyPureSeptuples:p.FullyPureSeptuples}
 for _,co:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,co);if !ok{continue};s:=up223cSnapshot(x,e,0);a,er:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if er!=nil{return UP253CResult{},er};b,er:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if er!=nil{return UP253CResult{},er};r.PairedInitialStates++;samples=append(samples,sample{up246cValues(s),up246cClass(a,b)})}}
 for omit:=0;omit<len(coords);omit++{sel:=[]string{};for i,c:=range coords{if i!=omit{sel=append(sel,c)}};groups:=map[string]map[string]bool{}
  for _,s:=range samples{key:=fmt.Sprintf("%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s",sel[0],s.values[sel[0]],sel[1],s.values[sel[1]],sel[2],s.values[sel[2]],sel[3],s.values[sel[3]],sel[4],s.values[sel[4]],sel[5],s.values[sel[5]],sel[6],s.values[sel[6]],sel[7],s.values[sel[7]]);if groups[key]==nil{groups[key]=map[string]bool{}};groups[key][s.class]=true}
  keys:=make([]string,0,len(groups));for k:=range groups{keys=append(keys,k)};sort.Strings(keys);sm:=UP253CEightSummary{Coordinates:sel,Groups:len(groups)}
  for _,k:=range keys{n:=len(groups[k]);if n==1{sm.PureGroups++}else{sm.MixedGroups++};if n>sm.MaxClassesPerGroup{sm.MaxClassesPerGroup=n}}
  if sm.MixedGroups==0{r.FullyPureOctuples++};r.Summaries=append(r.Summaries,sm)
 }
 r.CoordinateOctuples=len(r.Summaries);return r,nil
}
