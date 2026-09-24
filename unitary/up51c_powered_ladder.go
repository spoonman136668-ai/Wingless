package unitary

const UP51CPoweredLadderSchema = "wingless.up51c-powered-depth-ladder.v1"

type UP51CPoweredLadderResult struct {
	Schema              string                 `json:"schema"`
	Experiment          string                 `json:"experiment"`
	SourceUP50CSeal     string                 `json:"source_up50c_seal"`
	Dimension           int                    `json:"dimension"`
	Depths              []int                  `json:"depths"`
	Metrics             []UP50CPoweredMetric   `json:"metrics"`
	FirstFailingDepth   int                    `json:"first_failing_depth"`
	LastPassingDepth    int                    `json:"last_passing_depth"`
}

func RunUP51C() (UP51CPoweredLadderResult,error) {
	depths:=[]int{8192,32768,131072,524288,1048576}
	result:=UP51CPoweredLadderResult{
		Schema:UP51CPoweredLadderSchema,
		Experiment:"UP-51C-dim128-powered-depth-ladder",
		SourceUP50CSeal:"7eb1bfb8c13c139e732a72c08763e5a1043bbc4e",
		Dimension:128,
		Depths:append([]int(nil),depths...),
		FirstFailingDepth:-1,
		LastPassingDepth:-1,
	}
	for _,depth:=range depths{
		metric,err:=up50cPoweredMetric(128,depth)
		if err!=nil{return UP51CPoweredLadderResult{},err}
		result.Metrics=append(result.Metrics,metric)
		if metric.Gate{
			result.LastPassingDepth=depth
		}else if result.FirstFailingDepth<0{
			result.FirstFailingDepth=depth
		}
	}
	return result,nil
}
