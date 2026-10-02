package unitary

type wlmLmRawRepPredFreshHoldoutR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepPredFreshHoldoutR2() interface{} {
	base := RunWlmLmRawRepPredFreshHoldoutR1().(wlmLmRawRepPredFreshHoldoutR1Result)
	return wlmLmRawRepPredFreshHoldoutR2Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-PRED-FRESH-HOLDOUT-R2",
		Metrics: base.Metrics,
	}
}
