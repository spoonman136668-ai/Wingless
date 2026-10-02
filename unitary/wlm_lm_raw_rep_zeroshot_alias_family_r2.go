package unitary

type wlmLmRawRepZeroshotAliasFamilyR2Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

func RunWlmLmRawRepZeroshotAliasFamilyR2() interface{} {
	base:=RunWlmLmRawRepZeroshotAliasFamilyR1().(wlmLmRawRepZeroshotAliasFamilyR1Result)
	return wlmLmRawRepZeroshotAliasFamilyR2Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R2",
		Metrics:base.Metrics,
	}
}
