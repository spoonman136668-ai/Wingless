package unitary

type wlmLmRawRepZeroshotAliasFamilyR4Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepZeroshotAliasFamilyR4() interface{} {
	base := RunWlmLmRawRepZeroshotAliasFamilyR1().(wlmLmRawRepZeroshotAliasFamilyR1Result)
	return wlmLmRawRepZeroshotAliasFamilyR4Result{
		Schema: "wingless.research-scientific-result.v1",
		Experiment: "WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R4",
		Metrics: base.Metrics,
	}
}
