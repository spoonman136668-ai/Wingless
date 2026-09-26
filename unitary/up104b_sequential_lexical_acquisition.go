package unitary

const UP104BSequentialLexicalSchema = "wingless.up104b-sequential-lexical-acquisition.v1"

type UP104BStagePoint struct {
	Stage                    int     `json:"stage"`
	AcquiredVerbs            int     `json:"acquired_verbs"`
	BaseSeenAccuracy         float64 `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64 `json:"acquired_aggregate_accuracy"`
	HoldsAccuracy            float64 `json:"holds_accuracy"`
	NotesAccuracy            float64 `json:"notes_accuracy"`
	TellsAccuracy            float64 `json:"tells_accuracy"`
}

type UP104BSequentialLexicalResult struct {
	Schema           string             `json:"schema"`
	Experiment       string             `json:"experiment"`
	SourceUP103BSeal string             `json:"source_up103b_seal"`
	StateDimension   int                `json:"state_dimension"`
	EpochsPerBatch   int                `json:"epochs_per_batch"`
	LearningRate     float64            `json:"learning_rate"`
	GroundingPerVerb int                `json:"grounding_per_verb"`
	Points           []UP104BStagePoint `json:"points"`
}

func up104bVerbAccuracy(c *up97bClassifier, vi int) float64 {
	hits,total:=0,0
	for ni:=4;ni<6;ni++ {
		target:=up97bVerbClass(vi)
		pred:=up97bArgmax(c.probs(up99bPrefixEncoder(ni,vi,0)))
		total++
		if pred==target { hits++ }
	}
	return float64(hits)/float64(total)
}

func up104bPoint(stage int,c *up97bClassifier) UP104BStagePoint {
	holds:=up104bVerbAccuracy(c,2)
	notes:=up104bVerbAccuracy(c,5)
	tells:=up104bVerbAccuracy(c,8)
	agg:=0.0
	if stage>0 {
		sum:=0.0
		for i,vi:=range []int{2,5,8} {
			if i>=stage { break }
			sum+=up104bVerbAccuracy(c,vi)
		}
		agg=sum/float64(stage)
	}
	return UP104BStagePoint{
		Stage:stage,AcquiredVerbs:stage,
		BaseSeenAccuracy:up102bSeenAccuracy(c),
		AcquiredAggregateAccuracy:agg,
		HoldsAccuracy:holds,NotesAccuracy:notes,TellsAccuracy:tells,
	}
}

func up104bAcquire(c *up97bClassifier,currentVI int,previous []int) {
	for epoch:=0;epoch<20;epoch++ {
		for ni:=0;ni<4;ni++ { up103bTrainStep(c,ni,currentVI) }
		for _,vi:=range up101bTrainVerbs { up103bTrainStep(c,0,vi) }
		for _,vi:=range previous { up103bTrainStep(c,0,vi) }
	}
}

func RunUP104B()(UP104BSequentialLexicalResult,error){
	c:=up101bTrain()
	result:=UP104BSequentialLexicalResult{
		Schema:UP104BSequentialLexicalSchema,
		Experiment:"UP-104B-sequential-lexical-acquisition",
		SourceUP103BSeal:"dfb9e9378ed50e3075dc2ee3dd4697a941f66a1b",
		StateDimension:64,EpochsPerBatch:20,LearningRate:0.08,GroundingPerVerb:4,
	}
	result.Points=append(result.Points,up104bPoint(0,c))
	acquired:=[]int{}
	for stage,vi:=range []int{2,5,8} {
		up104bAcquire(c,vi,acquired)
		acquired=append(acquired,vi)
		result.Points=append(result.Points,up104bPoint(stage+1,c))
	}
	return result,nil
}
