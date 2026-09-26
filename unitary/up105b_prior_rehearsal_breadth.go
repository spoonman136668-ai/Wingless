package unitary

const UP105BPriorRehearsalBreadthSchema = "wingless.up105b-prior-rehearsal-breadth.v1"

type UP105BPoint struct {
	ReplaySubjects            int     `json:"replay_subjects"`
	Stage                     int     `json:"stage"`
	AcquiredVerbs             int     `json:"acquired_verbs"`
	BaseSeenAccuracy          float64 `json:"base_seen_accuracy"`
	AcquiredAggregateAccuracy float64 `json:"acquired_aggregate_accuracy"`
	HoldsAccuracy             float64 `json:"holds_accuracy"`
	NotesAccuracy             float64 `json:"notes_accuracy"`
	TellsAccuracy             float64 `json:"tells_accuracy"`
}

type UP105BPriorRehearsalBreadthResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP104BSeal string        `json:"source_up104b_seal"`
	StateDimension   int           `json:"state_dimension"`
	EpochsPerBatch   int           `json:"epochs_per_batch"`
	LearningRate     float64       `json:"learning_rate"`
	GroundingPerVerb int           `json:"grounding_per_verb"`
	BaseReplaySubjects int         `json:"base_replay_subjects"`
	Points           []UP105BPoint `json:"points"`
}

func up105bPoint(replaySubjects, stage int, c *up97bClassifier) UP105BPoint {
	p := up104bPoint(stage, c)
	return UP105BPoint{
		ReplaySubjects: replaySubjects,
		Stage: p.Stage,
		AcquiredVerbs: p.AcquiredVerbs,
		BaseSeenAccuracy: p.BaseSeenAccuracy,
		AcquiredAggregateAccuracy: p.AcquiredAggregateAccuracy,
		HoldsAccuracy: p.HoldsAccuracy,
		NotesAccuracy: p.NotesAccuracy,
		TellsAccuracy: p.TellsAccuracy,
	}
}

func up105bAcquire(c *up97bClassifier, currentVI int, previous []int, replaySubjects int) {
	for epoch := 0; epoch < 20; epoch++ {
		for ni := 0; ni < 4; ni++ {
			up103bTrainStep(c, ni, currentVI)
		}
		for _, vi := range up101bTrainVerbs {
			up103bTrainStep(c, 0, vi)
		}
		for _, vi := range previous {
			for ni := 0; ni < replaySubjects; ni++ {
				up103bTrainStep(c, ni, vi)
			}
		}
	}
}

func RunUP105B() (UP105BPriorRehearsalBreadthResult, error) {
	result := UP105BPriorRehearsalBreadthResult{
		Schema: UP105BPriorRehearsalBreadthSchema,
		Experiment: "UP-105B-prior-rehearsal-breadth",
		SourceUP104BSeal: "9a54896c6903f5c06f86263ec050af53cb8985a6",
		StateDimension: 64,
		EpochsPerBatch: 20,
		LearningRate: 0.08,
		GroundingPerVerb: 4,
		BaseReplaySubjects: 1,
	}

	for _, replaySubjects := range []int{1, 2, 4} {
		c := up101bTrain()
		result.Points = append(result.Points, up105bPoint(replaySubjects, 0, c))
		acquired := []int{}
		for stage, vi := range []int{2, 5, 8} {
			up105bAcquire(c, vi, acquired, replaySubjects)
			acquired = append(acquired, vi)
			result.Points = append(result.Points, up105bPoint(replaySubjects, stage+1, c))
		}
	}
	return result, nil
}
