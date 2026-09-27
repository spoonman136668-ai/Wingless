package unitary

import (
	"fmt"
	"sort"
)

const UP246CCoordinatePuritySchema = "wingless.up246c-initial-native-coordinate-purity.v1"

type UP246CCoordinateSummary struct {
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	PureGroups int `json:"pure_groups"`
	MixedGroups int `json:"mixed_groups"`
	MaxClassesPerGroup int `json:"max_classes_per_group"`
}

type UP246CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP245CSeal string `json:"source_up245c_seal"`
	ConditionA string `json:"condition_a"`
	ConditionB string `json:"condition_b"`
	PairedInitialStates int `json:"paired_initial_states"`
	Coordinates []string `json:"coordinates"`
	OutcomeClasses []string `json:"outcome_classes"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP246CCoordinateSummary `json:"summaries"`
}

func up246cClass(a,b UP237CArm)string{
	af:=a.Outcome=="loss";bf:=b.Outcome=="loss"
	if af&&bf{return "both_fail"}
	if af{return "a_only_fail"}
	if bf{return "b_only_fail"}
	return "neither_fail"
}

func up246cValues(s UP223CSnapshot)map[string]string{
	return map[string]string{
		"endangered_age":fmt.Sprint(s.EndangeredAge),
		"predicted_is_endangered":fmt.Sprint(s.PredictedIsEndangered),
		"hand_distance":fmt.Sprint(s.HandDistance),
		"age0":fmt.Sprint(s.Age0),
		"age1":fmt.Sprint(s.Age1),
		"age2":fmt.Sprint(s.Age2),
		"age3":fmt.Sprint(s.Age3),
		"adversarial_horizon":fmt.Sprint(s.AdversarialHorizon),
		"no_query_horizon":fmt.Sprint(s.NoQueryHorizon),
	}
}

func RunUP246C()(UP246CResult,error){
	cohorts:=[][]int{{0,9,10,15},{1,4,11,14},{2,5,8,13},{3,6,7,12}}
	hands:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15}
	coords:=[]string{"endangered_age","predicted_is_endangered","hand_distance","age0","age1","age2","age3","adversarial_horizon","no_query_horizon"}
	res:=UP246CResult{
		Schema:UP246CCoordinatePuritySchema,
		Experiment:"UP-246C-initial-native-coordinate-purity",
		SourceUP245CSeal:"78bab66e7a4a9e784425a228cd2dac023e230774",
		ConditionA:"schedule2|advance22",
		ConditionB:"schedule3|advance0",
		Coordinates:coords,
		OutcomeClasses:[]string{"a_only_fail","b_only_fail","both_fail","neither_fail"},
		InterventionChanged:false,
		ThresholdFittingUsed:false,
		ClassifierTrainingUsed:false,
		AdaptiveFeatureSelectionUsed:false,
		NewNativeFieldUsed:false,
		LiveActivation:false,
	}
	type sample struct{values map[string]string;class string}
	samples:=[]sample{}
	for _,cohort:=range cohorts{
		for _,hand:=range hands{
			x,e,ok:=up161cTriggerState(hand,cohort);if !ok{continue}
			s:=up223cSnapshot(x,e,0)
			a,err:=up237cRun(x,e,22,"alternating_shield","fixed_offset_refresh");if err!=nil{return UP246CResult{},err}
			b,err:=up237cRun(x,e,0,"fixed_offset_refresh","alternating_shield");if err!=nil{return UP246CResult{},err}
			res.PairedInitialStates++
			samples=append(samples,sample{values:up246cValues(s),class:up246cClass(a,b)})
		}
	}
	for _,coord:=range coords{
		groups:=map[string]map[string]bool{}
		for _,s:=range samples{
			v:=s.values[coord]
			if groups[v]==nil{groups[v]=map[string]bool{}}
			groups[v][s.class]=true
		}
		keys:=make([]string,0,len(groups));for k:=range groups{keys=append(keys,k)};sort.Strings(keys)
		sm:=UP246CCoordinateSummary{Coordinate:coord,Groups:len(groups)}
		for _,k:=range keys{
			n:=len(groups[k]);if n==1{sm.PureGroups++}else{sm.MixedGroups++};if n>sm.MaxClassesPerGroup{sm.MaxClassesPerGroup=n}
		}
		res.Summaries=append(res.Summaries,sm)
	}
	return res,nil
}
