package unitary

type wlmLmRawRepContextualAliasMultistepR3Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepContextualAliasMultistepR3() interface{} {
	base := RunWlmLmRawRepContextualAliasMultistepR2().(wlmLmRawRepContextualAliasMultistepR2Result)
	return wlmLmRawRepContextualAliasMultistepR3Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R3",
		Metrics: base.Metrics,
	}
}
