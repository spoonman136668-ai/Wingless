package unitary

const UP149BFourTailSchema="wingless.up149b-four-tail-causal-ranking.v1"

type UP149BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP148BSeal string `json:"source_up148b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	LearningRate float64 `json:"learning_rate"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	AdaptiveTailSelection bool `json:"adaptive_tail_selection"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP148BArm `json:"arms"`
}

func RunUP149B()(UP149BResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP149BResult{
		Schema:UP149BFourTailSchema,Experiment:"UP-149B-four-tail-causal-ranking",
		SourceUP148BSeal:"f5c8bb8c67260f266b74513a118ab4a09c7dce76",
		CommonPrefixEpochs:15,TerminalEpochs:5,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		AdaptiveTailSelection:false,ExtraUpdatesUsed:false,
	}
	for _,x:=range []struct{name string;shift int}{
		{"terminal_0",0},{"terminal_6",6},{"terminal_12",12},{"terminal_18",18},
	}{
		res.Arms=append(res.Arms,up148bRun(x.name,x.shift,o,r))
	}
	return res,nil
}
