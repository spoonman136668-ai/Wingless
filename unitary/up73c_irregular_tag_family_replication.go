package unitary

import "math"

const UP73CIrregularTagFamilySchema = "wingless.up73c-irregular-tag-family-replication.v1"

type UP73CIrregularTagFamilyResult struct {
	Schema             string                  `json:"schema"`
	Experiment         string                  `json:"experiment"`
	SourceUP72CSeal    string                  `json:"source_up72c_seal"`
	StateDimension     int                     `json:"state_dimension"`
	Banks              int                     `json:"banks"`
	ScheduleBases      []int                   `json:"schedule_bases"`
	NoiseLevels        []float64               `json:"noise_levels"`
	TagFamilies        []string                `json:"tag_families"`
	SelectionPerformed bool                    `json:"selection_performed"`
	Points             []UP72CTagGeometryPoint `json:"points"`
}

func up73cRotationTags(banks int,angle float64)[]float64{
	tags:=make([]float64,banks)
	for i:=0;i<banks;i++{tags[i]=angle*float64(i)}
	return tags
}

func up73cEval(family string,tags []float64,noise float64,seedBase int)UP72CTagGeometryPoint{
	p,err:=up72cEvaluate(family,tags,noise,seedBase)
	if err!=nil{
		return UP72CTagGeometryPoint{TagFamily:family,ScheduleBase:seedBase,Banks:10,MemoryNoise:noise,Gate:false,EvaluationError:err.Error()}
	}
	return p
}

func RunUP73C()(UP73CIrregularTagFamilyResult,error){
	schedules:=[]int{117000000,118000000}
	noises:=[]float64{0.002,0.003}
	families:=[]string{"golden_rotation","sqrt2_rotation","sqrt3_rotation"}
	result:=UP73CIrregularTagFamilyResult{
		Schema:UP73CIrregularTagFamilySchema,
		Experiment:"UP-73C-irregular-tag-family-replication",
		SourceUP72CSeal:"cdcd3c70b8b7c334df3fc8ee05168b63ababf052",
		StateDimension:16,Banks:10,
		ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),
		TagFamilies:append([]string(nil),families...),
		SelectionPerformed:false,
	}
	for _,seedBase:=range schedules{
		for _,noise:=range noises{
			result.Points=append(result.Points,up73cEval("golden_rotation",up57cGoldenTags(10),noise,seedBase))
			result.Points=append(result.Points,up73cEval("sqrt2_rotation",up73cRotationTags(10,math.Sqrt(2)),noise,seedBase))
			result.Points=append(result.Points,up73cEval("sqrt3_rotation",up73cRotationTags(10,math.Sqrt(3)),noise,seedBase))
		}
	}
	return result,nil
}
