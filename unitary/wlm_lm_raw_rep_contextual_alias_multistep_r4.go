package unitary

type wlmLmRawRepContextualAliasMultistepR4Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepContextualAliasMultistepR4() interface{} {
	base := RunWlmLmRawRepContextualAliasMultistepR2().(wlmLmRawRepContextualAliasMultistepR2Result)
	return wlmLmRawRepContextualAliasMultistepR4Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-CONTEXTUAL-ALIAS-MULTISTEP-R4",
		Metrics: base.Metrics,
	}
}
