package unitary

import (
	"fmt"
	"sort"
)

const UPLM1DThirdFamilyAdaptationSchema = "wingless.up-lm1d-third-family-byte-adaptation.v1"

type UPLM1DMetric struct {
	AdaptationEpochs int          `json:"adaptation_epochs"`
	Metric           UPLM0JMetric `json:"metric"`
}

type UPLM1DThirdFamilyAdaptationResult struct {
	Schema                       string          `json:"schema"`
	Experiment                   string          `json:"experiment"`
	SourceUPLM1CSeal             string          `json:"source_up_lm1c_seal"`
	StateDimension               int             `json:"state_dimension"`
	ExactRecallCap               int             `json:"exact_recall_cap"`
	BasePretrainEpochs           int             `json:"base_pretrain_epochs"`
	JointInterleavedEpochs       int             `json:"joint_interleaved_epochs"`
	LearningRate                 float64         `json:"learning_rate"`
	RouterRetrainingUsed         bool            `json:"router_retraining_used"`
	RecurrentParametersTrained   bool            `json:"recurrent_parameters_trained"`
	AttentionUsed                bool            `json:"attention_used"`
	FutureOracleUsed             bool            `json:"future_oracle_used"`
	OutputAlphabetExtended       bool            `json:"output_alphabet_extended"`
	NewOutputBytes               []int           `json:"new_output_bytes"`
	Metrics                      []UPLM1DMetric  `json:"metrics"`
}

func uplm1dCorpus()(train,held []uplm0fExample){
	for n:=0;n<6;n++ {
		for v:=0;v<6;v++ {
			for p:=0;p<4;p++ {
				ex:=uplm1cExample(n,v,p,"block")
				if (n+2*v+p)%3!=2 { train=append(train,ex) } else { held=append(held,ex) }
			}
		}
	}
	return
}

func uplm1dExtendAlphabet(src *uplm0aModel,examples ...[]uplm0fExample)(*uplm0aModel,[]int){
	seen:=map[byte]bool{}
	for _,b:=range src.alphabet { seen[b]=true }
	added:=[]int{}
	for _,set:=range examples {
		for _,ex:=range set {
			for i:=0;i<len(ex.text);i++ {
				b:=ex.text[i]
				if !seen[b] {
					seen[b]=true
					added=append(added,int(b))
				}
			}
		}
	}
	sort.Ints(added)
	alphabet:=append([]byte(nil),src.alphabet...)
	for _,v:=range added { alphabet=append(alphabet,byte(v)) }
	dst:=newUPLM0AModel(alphabet)
	for oldIndex,b:=range src.alphabet {
		newIndex:=dst.index[int(b)]
		copy(dst.w[newIndex],src.w[oldIndex])
		dst.b[newIndex]=src.b[oldIndex]
	}
	return dst,added
}

func uplm1dEval(model *uplm0aModel,classifier *uplm0jClassifier,examples []uplm0fExample,stream int,split string) UPLM0JMetric {
	return uplm1cEval(model,classifier,examples,stream,split)
}

func RunUPLM1D()(UPLM1DThirdFamilyAdaptationResult,error){
	baseTrain,baseHeld,alphabet:=uplm0fCorpus()
	paraTrain,paraHeld:=uplm0oParaphraseCorpus()
	thirdTrain,thirdHeld:=uplm1dCorpus()
	if len(baseTrain)!=len(paraTrain)||len(baseTrain)!=len(thirdTrain) {
		return UPLM1DThirdFamilyAdaptationResult{},fmt.Errorf("training corpus count mismatch: base=%d para=%d third=%d",len(baseTrain),len(paraTrain),len(thirdTrain))
	}

	start:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { start.trainSentence(ex.text,0.08) }
	}
	uplm1aJointTrain(start,baseTrain,paraTrain,20)
	start,added:=uplm1dExtendAlphabet(start,thirdTrain,thirdHeld)
	classifier:=uplm1cTrainClassifier()

	result:=UPLM1DThirdFamilyAdaptationResult{
		Schema:UPLM1DThirdFamilyAdaptationSchema,
		Experiment:"UP-LM1D-third-family-byte-adaptation",
		SourceUPLM1CSeal:"4ac7668597fbe21ab7d2e7b17f65267cb3912e48",
		StateDimension:64,ExactRecallCap:16,BasePretrainEpochs:20,JointInterleavedEpochs:20,LearningRate:0.08,
		RouterRetrainingUsed:false,RecurrentParametersTrained:false,AttentionUsed:false,FutureOracleUsed:false,
		OutputAlphabetExtended:len(added)>0,NewOutputBytes:append([]int(nil),added...),
	}

	for _,adapt:=range []int{0,1,2,4} {
		model:=uplm0oCloneModel(start)
		for epoch:=0;epoch<adapt;epoch++ {
			for i:=range thirdTrain {
				model.trainSentence(thirdTrain[i].text,0.08)
				model.trainSentence(baseTrain[i].text,0.08)
				model.trainSentence(paraTrain[i].text,0.08)
			}
		}
		result.Metrics=append(result.Metrics,
			UPLM1DMetric{AdaptationEpochs:adapt,Metric:uplm1dEval(model,classifier,baseHeld,1,"base_block_stream1")},
			UPLM1DMetric{AdaptationEpochs:adapt,Metric:uplm1dEval(model,classifier,paraHeld,1,"paraphrase_block_stream1")},
		)
		for _,order:=range []string{"block","per_name","paired_names","stores_then_local_reports","reverse_report_tail"} {
			held:=uplm1cHeldout(order)
			for _,stream:=range []int{1,4} {
				result.Metrics=append(result.Metrics,
					UPLM1DMetric{AdaptationEpochs:adapt,Metric:uplm1dEval(model,classifier,held,stream,"third_"+order+"_stream"+itoa(stream))},
				)
			}
		}
	}
	return result,nil
}
