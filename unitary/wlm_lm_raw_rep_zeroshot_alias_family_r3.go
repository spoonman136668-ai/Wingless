package unitary

type wlmLmRawRepZeroshotAliasFamilyR3Result struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	Metrics map[string]float64 `json:"metrics"`
}

// RunWlmLmRawRepZeroshotAliasFamilyR3 delegates to the exact frozen R1 scientific function.
// R3 exists only to restore a valid two-commit qualification boundary.
func RunWlmLmRawRepZeroshotAliasFamilyR3() interface{} {
	base:=RunWlmLmRawRepZeroshotAliasFamilyR1().(wlmLmRawRepZeroshotAliasFamilyR1Result)
	return wlmLmRawRepZeroshotAliasFamilyR3Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-LM-RAW-REP-ZEROSHOT-ALIAS-FAMILY-R3",
		Metrics:base.Metrics,
	}
}
