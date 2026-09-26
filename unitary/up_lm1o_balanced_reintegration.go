package unitary

const UPLM1OBalancedReintegrationSchema = "wingless.up-lm1o-balanced-reintegration.v1"

type UPLM1OByteMetric struct {
	Arm          string  `json:"arm"`
	Split        string  `json:"split"`
	Top1Accuracy float64 `json:"top1_accuracy"`
	Perplexity   float64 `json:"perplexity"`
}

type UPLM1ORoutingMetric struct {
	Arm    string       `json:"arm"`
	Metric UPLM0JMetric `json:"metric"`
}

type UPLM1OBalancedReintegrationResult struct {
	Schema                     string                 `json:"schema"`
	Experiment                 string                 `json:"experiment"`
	SourceUPLM1NSeal           string                 `json:"source_up_lm1n_seal"`
	StateDimension             int                    `json:"state_dimension"`
	ExactRecallCap             int                    `json:"exact_recall_cap"`
	AdaptationEpochs           int                    `json:"adaptation_epochs"`
	FullStepLearningRate       float64                `json:"full_step_learning_rate"`
	HalfStepLearningRate       float64                `json:"half_step_learning_rate"`
	RouterRetrainingUsed       bool                   `json:"router_retraining_used"`
	RecurrentParametersTrained bool                   `json:"recurrent_parameters_trained"`
	RecallCapChanged           bool                   `json:"recall_cap_changed"`
	AttentionUsed              bool                   `json:"attention_used"`
	FutureOracleUsed           bool                   `json:"future_oracle_used"`
	ByteMetrics                []UPLM1OByteMetric     `json:"byte_metrics"`
	RoutingMetrics             []UPLM1ORoutingMetric  `json:"routing_metrics"`
}

func RunUPLM1O()(UPLM1OBalancedReintegrationResult,error) {
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1OBalancedReintegrationResult{},err }

	d:=uplm1lStoreDirection()
	classifier:=uplm1lTrainClassifier("states_store_orthogonal",d)

	result:=UPLM1OBalancedReintegrationResult{
		Schema:UPLM1OBalancedReintegrationSchema,
		Experiment:"UP-LM1O-balanced-reintegration",
		SourceUPLM1NSeal:"cf646c9a3da3da131ea245ca757f6ab3236a9014",
		StateDimension:64,ExactRecallCap:16,AdaptationEpochs:4,
		FullStepLearningRate:0.08,HalfStepLearningRate:0.04,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,RecallCapChanged:false,AttentionUsed:false,FutureOracleUsed:false,
	}

	families:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	heldSets:=[]struct{name string;data []uplm0fExample}{
		{"base_heldout",baseHeld},{"paraphrase_heldout",paraHeld},{"third_heldout",thirdHeld},{"fourth_heldout",fourthHeld},
	}

	for _,arm:=range []string{"cyclic_control","mirrored_by_index","rotating_palindromic_split"} {
		model:=uplm0oCloneModel(start)
		uplm1nTrainArm(model,arm,families)

		for _,x:=range heldSets {
			m:=uplm0oByteEval(model,x.data,0,x.name)
			result.ByteMetrics=append(result.ByteMetrics,UPLM1OByteMetric{
				Arm:arm,Split:x.name,Top1Accuracy:m.Top1Accuracy,Perplexity:m.Perplexity,
			})
		}

		for _,stream:=range []int{1,4} {
			result.RoutingMetrics=append(result.RoutingMetrics,
				UPLM1ORoutingMetric{Arm:arm,Metric:uplm1mEvaluate(model,classifier,d,baseHeld,stream,"base_block_stream"+itoa(stream))},
				UPLM1ORoutingMetric{Arm:arm,Metric:uplm1mEvaluate(model,classifier,d,paraHeld,stream,"paraphrase_block_stream"+itoa(stream))},
				UPLM1ORoutingMetric{Arm:arm,Metric:uplm1mEvaluate(model,classifier,d,thirdHeld,stream,"third_block_stream"+itoa(stream))},
			)
		}

		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1jHeldout(order)
			for _,stream:=range []int{1,4} {
				result.RoutingMetrics=append(result.RoutingMetrics,
					UPLM1ORoutingMetric{
						Arm:arm,
						Metric:uplm1mEvaluate(model,classifier,d,held,stream,"fourth_"+order+"_stream"+itoa(stream)),
					},
				)
			}
		}
	}
	return result,nil
}
