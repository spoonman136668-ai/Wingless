package unitary

const UP86CHotsetTransitionSchema = "wingless.up86c-hotset-transition.v1"

type UP86CHotsetTransitionResult struct {
	Schema            string       `json:"schema"`
	Experiment        string       `json:"experiment"`
	SourceUP85CSeal   string       `json:"source_up85c_seal"`
	ActiveKeys        int          `json:"active_keys"`
	ChurnWrites       int          `json:"churn_writes"`
	ExactRecallCap    int          `json:"exact_recall_cap"`
	EntryPayloadBytes int          `json:"entry_payload_bytes"`
	FutureOracleUsed  bool         `json:"future_oracle_used"`
	Points            []UP85CPoint `json:"points"`
}

func RunUP86C()(UP86CHotsetTransitionResult,error){
	result:=UP86CHotsetTransitionResult{
		Schema:UP86CHotsetTransitionSchema,
		Experiment:"UP-86C-hotset-transition",
		SourceUP85CSeal:"8586ad159bde0d49ffad2e1f3e61e306670a4b3c",
		ActiveKeys:32,ChurnWrites:24,ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,
	}
	seedBases:=[]int{163000000,164000000}
	for _,policy:=range []string{"fifo","lru","two_bit_aging"} {
		for _,hot:=range []int{13,14,15,16} {
			result.Points=append(result.Points,up85cRun(policy,hot,seedBases))
		}
	}
	return result,nil
}
