package unitary

const UP60CSixBankHighNoiseSchema = "wingless.up60c-six-bank-high-noise.v1"

type UP60CSixBankHighNoiseResult struct {
	Schema             string             `json:"schema"`
	Experiment         string             `json:"experiment"`
	SourceUP59CSeal    string             `json:"source_up59c_seal"`
	Banks              int                `json:"banks"`
	StateDimension     int                `json:"state_dimension"`
	ScheduleBases      []int              `json:"schedule_bases"`
	NoiseLevels        []float64          `json:"noise_levels"`
	Families           []string           `json:"families"`
	SelectionPerformed bool               `json:"selection_performed"`
	Metrics            []UP59CNoiseMetric `json:"metrics"`
}

func RunUP60C()(UP60CSixBankHighNoiseResult,error){
	schedules:=[]int{75000000,76000000}
	noises:=[]float64{0.01,0.015,0.02,0.03,0.04}
	all:=up58cFamilies()
	var families []UP58CTagFamily
	for _,f:=range all{
		if f.Name=="golden_rotation"||f.Name=="fixed_irregular"{families=append(families,f)}
	}
	result:=UP60CSixBankHighNoiseResult{
		Schema:UP60CSixBankHighNoiseSchema,
		Experiment:"UP-60C-six-bank-high-noise",
		SourceUP59CSeal:"55f1725554d4f7493f206852385545c74c2d5d69",
		Banks:6,StateDimension:16,
		ScheduleBases:append([]int(nil),schedules...),
		NoiseLevels:append([]float64(nil),noises...),
		SelectionPerformed:false,
	}
	for _,f:=range families{result.Families=append(result.Families,f.Name)}
	for _,seedBase:=range schedules{
		for _,f:=range families{
			for _,noise:=range noises{
				m,err:=up59cEvaluate(f.Tags,f.Name,noise,seedBase)
				if err!=nil{return UP60CSixBankHighNoiseResult{},err}
				result.Metrics=append(result.Metrics,m)
			}
		}
	}
	return result,nil
}
