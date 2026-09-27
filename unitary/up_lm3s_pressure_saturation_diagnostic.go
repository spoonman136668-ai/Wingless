package unitary

const UPLM3SPressureSaturationSchema="wingless.up-lm3s-pressure-saturation-diagnostic.v1"

type UPLM3SMetric struct{
	PressureLevel int `json:"pressure_level"`
	ResourceConfigurations int `json:"resource_configurations"`
	MinEarliestFailed int `json:"min_earliest_failed"`
	MaxEarliestFailed int `json:"max_earliest_failed"`
	OutcomeRange int `json:"outcome_range"`
	DistinctOutcomes int `json:"distinct_outcomes"`
	ConfigurationsAtMin int `json:"configurations_at_min"`
	ConfigurationsAtMax int `json:"configurations_at_max"`
}
type UPLM3SResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3RSeal string `json:"source_up_lm3r_seal"`
	PressureLevels []int `json:"pressure_levels"`
	ResourceConfigurationsPerPressure int `json:"resource_configurations_per_pressure"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptivePressureSelectionUsed bool `json:"adaptive_pressure_selection_used"`
	ResourceTuningUsed bool `json:"resource_tuning_used"`
	Metrics []UPLM3SMetric `json:"metrics"`
}
func RunUPLM3S()(UPLM3SResult,error){
	pressures:=[]int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16}
	budgets:=[]int{4,5,6};starts:=[]int{3,4};tps:=[]int{1,2,3};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	res:=UPLM3SResult{Schema:UPLM3SPressureSaturationSchema,Experiment:"UP-LM3S-pressure-saturation-diagnostic",SourceUPLM3RSeal:"4f546af41f8e85af2267c13d0ac8c4851efab798",PressureLevels:pressures,ResourceConfigurationsPerPressure:18,CounterfactualOnly:true,LiveActivation:false,AdaptivePressureSelectionUsed:false,ResourceTuningUsed:false}
	for _,pressure:=range pressures{
		vals:=make([]int,0,18)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			v:=0
			for _,rot:=range rots{for _,perm:=range perms{_,ef,_:=uplm3oRun(rot,perm,"earliest_deadline",budget,start,tp,pressure);v+=ef}}
			vals=append(vals,v)
		}}}
		minv,maxv:=vals[0],vals[0];distinct:=map[int]bool{}
		for _,v:=range vals{
			distinct[v]=true
			if v<minv{minv=v}
			if v>maxv{maxv=v}
		}
		atMin,atMax:=0,0
		for _,v:=range vals{if v==minv{atMin++};if v==maxv{atMax++}}
		res.Metrics=append(res.Metrics,UPLM3SMetric{PressureLevel:pressure,ResourceConfigurations:len(vals),MinEarliestFailed:minv,MaxEarliestFailed:maxv,OutcomeRange:maxv-minv,DistinctOutcomes:len(distinct),ConfigurationsAtMin:atMin,ConfigurationsAtMax:atMax})
	}
	return res,nil
}
