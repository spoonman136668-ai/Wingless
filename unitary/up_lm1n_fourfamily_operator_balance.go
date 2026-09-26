package unitary

const UPLM1NOperatorBalanceSchema = "wingless.up-lm1n-fourfamily-operator-balance.v1"

type UPLM1NMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1NSummary struct {
	Arm                   string  `json:"arm"`
	MinHeldoutAccuracy    float64 `json:"min_heldout_accuracy"`
	MeanHeldoutAccuracy   float64 `json:"mean_heldout_accuracy"`
	HeldoutAccuracySpread float64 `json:"heldout_accuracy_spread"`
}

type UPLM1NOperatorBalanceResult struct {
	Schema                     string          `json:"schema"`
	Experiment                 string          `json:"experiment"`
	SourceUPLM1MSeal           string          `json:"source_up_lm1m_seal"`
	StateDimension             int             `json:"state_dimension"`
	AdaptationEpochs           int             `json:"adaptation_epochs"`
	FullStepLearningRate       float64         `json:"full_step_learning_rate"`
	HalfStepLearningRate       float64         `json:"half_step_learning_rate"`
	PerFamilyNominalLRPerIndex float64         `json:"per_family_nominal_lr_per_index"`
	RecurrentParametersTrained bool            `json:"recurrent_parameters_trained"`
	RouterUsed                 bool            `json:"router_used"`
	ExactRecallUsed            bool            `json:"exact_recall_used"`
	AdaptiveOrderingUsed       bool            `json:"adaptive_ordering_used"`
	Metrics                    []UPLM1NMetric  `json:"metrics"`
	Summaries                  []UPLM1NSummary `json:"summaries"`
}

func uplm1nTrainArm(model *uplm0aModel,arm string,families [4][]uplm0fExample) {
	for epoch:=0;epoch<4;epoch++ {
		for i:=range families[0] {
			switch arm {
			case "cyclic_control":
				for f:=0;f<4;f++ {
					model.trainSentence(families[f][i].text,0.08)
				}
			case "mirrored_by_index":
				if i%2==0 {
					for f:=0;f<4;f++ {
						model.trainSentence(families[f][i].text,0.08)
					}
				}else{
					for f:=3;f>=0;f-- {
						model.trainSentence(families[f][i].text,0.08)
					}
				}
			case "rotating_palindromic_split":
				center:=i%4
				others:=[3]int{(center+1)%4,(center+2)%4,(center+3)%4}
				for j:=0;j<3;j++ {
					model.trainSentence(families[others[j]][i].text,0.04)
				}
				model.trainSentence(families[center][i].text,0.08)
				for j:=2;j>=0;j-- {
					model.trainSentence(families[others[j]][i].text,0.04)
				}
			}
		}
	}
}

func uplm1nSummary(arm string,metrics []UPLM1NMetric) UPLM1NSummary {
	vals:=[]float64{}
	for _,m:=range metrics {
		if m.Arm!=arm { continue }
		if m.Split=="base_heldout"||m.Split=="paraphrase_heldout"||m.Split=="third_heldout"||m.Split=="fourth_heldout" {
			vals=append(vals,m.Top1Accuracy)
		}
	}
	min,max,sum:=1.0,0.0,0.0
	for _,v:=range vals {
		if v<min { min=v }
		if v>max { max=v }
		sum+=v
	}
	mean:=0.0
	if len(vals)>0 { mean=sum/float64(len(vals)) }
	return UPLM1NSummary{
		Arm:arm,MinHeldoutAccuracy:min,MeanHeldoutAccuracy:mean,HeldoutAccuracySpread:max-min,
	}
}

func RunUPLM1N()(UPLM1NOperatorBalanceResult,error){
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1NOperatorBalanceResult{},err }

	if len(baseTrain)!=len(paraTrain)||len(baseTrain)!=len(thirdTrain)||len(baseTrain)!=len(fourthTrain) {
		return UPLM1NOperatorBalanceResult{},fmt.Errorf("matched family counts differ")
	}

	result:=UPLM1NOperatorBalanceResult{
		Schema:UPLM1NOperatorBalanceSchema,
		Experiment:"UP-LM1N-fourfamily-operator-balance",
		SourceUPLM1MSeal:"a87c788a0690caae6bb47fe84b9ba74f1935a61b",
		StateDimension:64,AdaptationEpochs:4,
		FullStepLearningRate:0.08,HalfStepLearningRate:0.04,PerFamilyNominalLRPerIndex:0.08,
		RecurrentParametersTrained:false,RouterUsed:false,ExactRecallUsed:false,AdaptiveOrderingUsed:false,
	}

	families:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	evalSets:=[]struct{name string;train,held []uplm0fExample}{
		{"base",baseTrain,baseHeld},
		{"paraphrase",paraTrain,paraHeld},
		{"third",thirdTrain,thirdHeld},
		{"fourth",fourthTrain,fourthHeld},
	}
	arms:=[]string{"cyclic_control","mirrored_by_index","rotating_palindromic_split"}
	for _,arm:=range arms {
		model:=uplm0oCloneModel(start)
		uplm1nTrainArm(model,arm,families)
		for _,x:=range evalSets {
			mt:=uplm0oByteEval(model,x.train,0,x.name+"_train")
			mh:=uplm0oByteEval(model,x.held,0,x.name+"_heldout")
			result.Metrics=append(result.Metrics,
				UPLM1NMetric{Arm:arm,Split:x.name+"_train",Top1Accuracy:mt.Top1Accuracy,Perplexity:mt.Perplexity},
				UPLM1NMetric{Arm:arm,Split:x.name+"_heldout",Top1Accuracy:mh.Top1Accuracy,Perplexity:mh.Perplexity},
			)
		}
		result.Summaries=append(result.Summaries,uplm1nSummary(arm,result.Metrics))
	}
	return result,nil
}
