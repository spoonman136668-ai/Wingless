package unitary

const UPLM1QFivePositionSchema = "wingless.up-lm1q-fivefamily-cyclic-position.v1"

type UPLM1QMetric struct {
	Arm          string  `json:"arm"`
	Family       string  `json:"family"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
	Position     int     `json:"update_position"`
}

type UPLM1QSummary struct {
	Arm                   string  `json:"arm"`
	Rotation              int     `json:"rotation"`
	MinHeldoutAccuracy    float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy   float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
	MinFamily             string  `json:"min_family"`
	MaxFamily             string  `json:"max_family"`
}

type UPLM1QFivePositionResult struct {
	Schema                     string          `json:"schema"`
	Experiment                 string          `json:"experiment"`
	SourceUPLM1PSeal           string          `json:"source_up_lm1p_seal"`
	StateDimension             int             `json:"state_dimension"`
	AdaptationEpochs           int             `json:"adaptation_epochs"`
	LearningRate               float64         `json:"learning_rate"`
	RecurrentParametersTrained bool            `json:"recurrent_parameters_trained"`
	RouterUsed                 bool            `json:"router_used"`
	ExactRecallUsed            bool            `json:"exact_recall_used"`
	AdaptiveOrderingUsed       bool            `json:"adaptive_ordering_used"`
	SixthFamilyUsed            bool            `json:"sixth_family_used"`
	Metrics                    []UPLM1QMetric  `json:"metrics"`
	Summaries                  []UPLM1QSummary `json:"summaries"`
}

var uplm1qFamilyNames=[5]string{"base","paraphrase","third","fourth","fifth"}

func uplm1qTrain(model *uplm0aModel,families [5][]uplm0fExample,rotation int) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range families[0] {
			for j:=0;j<5;j++ {
				f:=(rotation+j)%5
				model.trainSentence(families[f][i].text,0.08)
			}
		}
	}
}

func uplm1qPosition(rotation,family int) int {
	for pos:=0;pos<5;pos++ {
		if (rotation+pos)%5==family { return pos }
	}
	return -1
}

func uplm1qSummary(arm string,rotation int,metrics []UPLM1QMetric) UPLM1QSummary {
	minVal,maxVal,sum:=1.0,-1.0,0.0
	minFamily,maxFamily:="",""
	count:=0
	for _,m:=range metrics {
		if m.Arm!=arm { continue }
		count++
		sum+=m.Top1Accuracy
		if m.Top1Accuracy<minVal { minVal=m.Top1Accuracy;minFamily=m.Family }
		if m.Top1Accuracy>maxVal { maxVal=m.Top1Accuracy;maxFamily=m.Family }
	}
	mean:=0.0
	if count>0 { mean=sum/float64(count) }
	return UPLM1QSummary{
		Arm:arm,Rotation:rotation,MinHeldoutAccuracy:minVal,MeanHeldoutAccuracy:mean,
		HeldoutAccuracySpread:maxVal-minVal,MinFamily:minFamily,MaxFamily:maxFamily,
	}
}

func RunUPLM1Q()(UPLM1QFivePositionResult,error) {
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1QFivePositionResult{},err }

	fourFamilies:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",fourFamilies)

	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	families:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}

	result:=UPLM1QFivePositionResult{
		Schema:UPLM1QFivePositionSchema,Experiment:"UP-LM1Q-fivefamily-cyclic-position",
		SourceUPLM1PSeal:"b34310666d604423be4a98a46cfc0a8d49fec26e",
		StateDimension:64,AdaptationEpochs:4,LearningRate:0.08,
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,
		AdaptiveOrderingUsed:false,SixthFamilyUsed:false,
	}

	for rotation:=0;rotation<5;rotation++ {
		arm:="rotate_"+itoa(rotation)
		model:=uplm0oCloneModel(start)
		uplm1qTrain(model,families,rotation)
		for f:=0;f<5;f++ {
			m:=uplm0oByteEval(model,held[f],0,uplm1qFamilyNames[f]+"_heldout")
			result.Metrics=append(result.Metrics,UPLM1QMetric{
				Arm:arm,Family:uplm1qFamilyNames[f],Top1Accuracy:m.Top1Accuracy,
				Perplexity:m.Perplexity,Position:uplm1qPosition(rotation,f),
			})
		}
		result.Summaries=append(result.Summaries,uplm1qSummary(arm,rotation,result.Metrics))
	}
	return result,nil
}
