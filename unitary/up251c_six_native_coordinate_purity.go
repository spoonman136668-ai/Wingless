package unitary
import("fmt";"sort")
const UP251CSixPuritySchema="wingless.up251c-six-native-coordinate-purity.v1"
type UP251CSixSummary struct{Coordinates []string `json:"coordinates"`;Groups int `json:"groups"`;PureGroups int `json:"pure_groups"`;MixedGroups int `json:"mixed_groups"`;MaxClassesPerGroup int `json:"max_classes_per_group"`}
type UP251CResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUP250CSeal string `json:"source_up250c_seal"`;PairedInitialStates int `json:"paired_initial_states"`;CoordinateCount int `json:"coordinate_count"`;ParentCoordinateQuintuples int `json:"parent_coordinate_quintuples"`;ParentFullyPureQuintuples int `json:"parent_fully_pure_quintuples"`;CoordinateSextuples int `json:"coordinate_sextuples"`;FullyPureSextuples int `json:"fully_pure_sextuples"`;InterventionChanged bool `json:"intervention_changed"`;ThresholdFittingUsed bool `json:"threshold_fitting_used"`;ClassifierTrainingUsed bool `json:"classifier_training_used"`;AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`;NewNativeFieldUsed bool `json:"new_native_field_used"`;LiveActivation bool `json:"live_activation"`;Summaries []UP251CSixSummary `json:"summaries"`}
func RunUP251C()(UP251CResult,error){
 p,err:=RunUP250C();if err!=nil{return UP251CResult{},err}
 cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}};hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
 coords:=[]string{"endangered_age","predicted_is_endangered","hand_distance","age0","age1","age2","age3","adversarial_horizon","no_query_horizon"}
 type sample struct{values map[string]string;class string};samples:=[]sample{}
 r:=UP251CResult{Schema:UP251CSixPuritySchema,Experiment:"UP-251C-six-native-coordinate-purity",SourceUP250CSeal:"7fc45861274d04f44b98585870b61554db70a8d5",CoordinateCount:len(coords),ParentCoordinateQuintuples:p.CoordinateQuintuples,ParentFullyPureQuintuples:p.FullyPureQuintuples}
 for _,co:=range cohorts{for _,hand:=range hands{x,e,ok:=up161cTriggerState(hand,co);if !ok{continue};s:=up223cSnapshot(x,e,0);a,er:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if er!=nil{return UP251CResult{},er};b,er:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if er!=nil{return UP251CResult{},er};r.PairedInitialStates++;samples=append(samples,sample{up246cValues(s),up246cClass(a,b)})}}
 for a:=0;a<len(coords);a++{for b:=a+1;b<len(coords);b++{for c:=b+1;c<len(coords);c++{for d:=c+1;d<len(coords);d++{for e:=d+1;e<len(coords);e++{for f:=e+1;f<len(coords);f++{
  sel:=[]string{coords[a],coords[b],coords[c],coords[d],coords[e],coords[f]};groups:=map[string]map[string]bool{}
  for _,s:=range samples{key:=fmt.Sprintf("%s=%s|%s=%s|%s=%s|%s=%s|%s=%s|%s=%s",sel[0],s.values[sel[0]],sel[1],s.values[sel[1]],sel[2],s.values[sel[2]],sel[3],s.values[sel[3]],sel[4],s.values[sel[4]],sel[5],s.values[sel[5]]);if groups[key]==nil{groups[key]=map[string]bool{}};groups[key][s.class]=true}
  keys:=make([]string,0,len(groups));for k:=range groups{keys=append(keys,k)};sort.Strings(keys);sm:=UP251CSixSummary{Coordinates:sel,Groups:len(groups)}
  for _,k:=range keys{n:=len(groups[k]);if n==1{sm.PureGroups++}else{sm.MixedGroups++};if n>sm.MaxClassesPerGroup{sm.MaxClassesPerGroup=n}}
  if sm.MixedGroups==0{r.FullyPureSextuples++};r.Summaries=append(r.Summaries,sm)
 }}}}}}
 r.CoordinateSextuples=len(r.Summaries);return r,nil
}
