package unitary

const UP71CBankNoiseBoundarySchema = "wingless.up71c-bank-noise-boundary.v1"

type UP71CBankNoiseBoundaryResult struct {
	Schema             string                  `json:"schema"`
	Experiment         string                  `json:"experiment"`
	SourceUP70CSeal    string                  `json:"source_up70c_seal"`
	StateDimension     int                     `json:"state_dimension"`
	ScheduleBases      []int                   `json:"schedule_bases"`
	BankCounts         []int                   `json:"bank_counts"`
	NoiseLevels        []float64               `json:"noise_levels"`
	SelectionPerformed bool                    `json:"selection_performed"`
	Points             []UP57CGoldenScalePoint `json:"points"`
}

func RunUP71C()(UP71CBankNoiseBoundaryResult,error){
	schedules:=[]int{113000000,114000000};banks:=[]int{9,10};noises:=[]float64{0,0.001,0.002,0.003,0.004,0.008}
	result:=UP71CBankNoiseBoundaryResult{Schema:UP71CBankNoiseBoundarySchema,Experiment:"UP-71C-bank-noise-boundary",SourceUP70CSeal:"52af265b43b4074e62bdc3836ab68c7d90e58452",StateDimension:16,ScheduleBases:append([]int(nil),schedules...),BankCounts:append([]int(nil),banks...),NoiseLevels:append([]float64(nil),noises...),SelectionPerformed:false}
	for _,seedBase:=range schedules { for _,bankCount:=range banks { for _,noise:=range noises {
		p,err:=up57cEvaluate(bankCount,noise,seedBase);if err!=nil{return UP71CBankNoiseBoundaryResult{},err};result.Points=append(result.Points,p)
	}}}
	return result,nil
}
