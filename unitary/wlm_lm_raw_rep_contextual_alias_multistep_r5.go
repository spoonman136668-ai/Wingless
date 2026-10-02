package unitary

type wlmLmRawRepContextualAliasMultistepR5Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepContextualAliasMultistepR5() interface{} {
	base := RunWlmLmRawRepContextualAliasMultistepR2().(wlmLmRawRepContextualAliasMultistepR2Result)
	return wlmLmRawRepContextualAliasMultistepR5Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R5",
		Metrics: base.Metrics,
	}
}
