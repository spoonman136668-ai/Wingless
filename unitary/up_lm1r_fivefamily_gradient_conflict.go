package unitary

import (
	"fmt"
	"math"
)

const UPLM1RGradientConflictSchema = "wingless.up-lm1r-fivefamily-gradient-conflict.v1"

type UPLM1RFamilyMetric struct {
	Family          string  `json:"family"`
	TrainExamples   int     `json:"train_examples"`
	GradientNorm    float64 `json:"aggregate_gradient_norm"`
	HeldoutAccuracy float64 `json:"heldout_accuracy"`
	HeldoutPerplexity float64 `json:"heldout_perplexity"`
}

type UPLM1RPairMetric struct {
	FamilyA                    string  `json:"family_a"`
	FamilyB                    string  `json:"family_b"`
	MatchedPairs               int     `json:"matched_pairs"`
	AggregateCosine            float64 `json:"aggregate_cosine"`
	AggregateDotProduct        float64 `json:"aggregate_dot_product"`
	MeanMatchedCosine          float64 `json:"mean_matched_cosine"`
	MinMatchedCosine           float64 `json:"min_matched_cosine"`
	MaxMatchedCosine           float64 `json:"max_matched_cosine"`
	NegativeMatchedFraction    float64 `json:"negative_matched_fraction"`
	NonPositiveMatchedFraction float64 `json:"non_positive_matched_fraction"`
}

type UPLM1RGradientConflictResult struct {
	Schema                        string                `json:"schema"`
	Experiment                    string                `json:"experiment"`
	SourceUPLM1QSeal              string                `json:"source_up_lm1q_seal"`
	StateDimension                int                   `json:"state_dimension"`
	FamilyCount                   int                   `json:"family_count"`
	PairCount                     int                   `json:"pair_count"`
	ParameterUpdatesAfterStart    bool                  `json:"parameter_updates_after_start"`
	RecurrentGradientsComputed    bool                  `json:"recurrent_gradients_computed"`
	RouterUsed                    bool                  `json:"router_used"`
	ExactRecallUsed               bool                  `json:"exact_recall_used"`
	SixthFamilyUsed               bool                  `json:"sixth_family_used"`
	FamilyMetrics                 []UPLM1RFamilyMetric  `json:"family_metrics"`
	PairMetrics                   []UPLM1RPairMetric    `json:"pair_metrics"`
	FifthVsPriorMeanAggregateCosine float64             `json:"fifth_vs_prior_mean_aggregate_cosine"`
	FifthVsPriorMinAggregateCosine  float64             `json:"fifth_vs_prior_min_aggregate_cosine"`
	PriorVsPriorMeanAggregateCosine float64             `json:"prior_vs_prior_mean_aggregate_cosine"`
	FifthGradientNormRelativeToPriorMean float64         `json:"fifth_gradient_norm_relative_to_prior_mean"`
}

func uplm1rAddGradient(dst *uplm0xGradient,src uplm0xGradient) {
	for c:=range dst.w {
		for i:=0;i<64;i++ { dst.w[c][i]+=src.w[c][i] }
		dst.b[c]+=src.b[c]
	}
}

func uplm1rGradientNorm(g uplm0xGradient) float64 {
	_,n2,_:=uplm0xPairStats(g,g)
	return math.Sqrt(n2)
}

func uplm1rPairMetric(aName,bName string,aAgg,bAgg uplm0xGradient,aRows,bRows []uplm0xGradient) UPLM1RPairMetric {
	dot,a2,b2:=uplm0xPairStats(aAgg,bAgg)
	minCos,maxCos:=math.Inf(1),math.Inf(-1)
	sum:=0.0
	negative,nonPositive:=0,0
	for i:=range aRows {
		d,aa,bb:=uplm0xPairStats(aRows[i],bRows[i])
		c:=uplm0xCosine(d,aa,bb)
		if c<minCos { minCos=c }
		if c>maxCos { maxCos=c }
		sum+=c
		if c<0 { negative++ }
		if c<=0 { nonPositive++ }
	}
	n:=len(aRows)
	mean:=0.0
	if n>0 { mean=sum/float64(n) }
	return UPLM1RPairMetric{
		FamilyA:aName,FamilyB:bName,MatchedPairs:n,
		AggregateCosine:uplm0xCosine(dot,a2,b2),AggregateDotProduct:dot,
		MeanMatchedCosine:mean,MinMatchedCosine:minCos,MaxMatchedCosine:maxCos,
		NegativeMatchedFraction:float64(negative)/float64(n),
		NonPositiveMatchedFraction:float64(nonPositive)/float64(n),
	}
}

func RunUPLM1R()(UPLM1RGradientConflictResult,error) {
	start,baseTrain,baseHeld,paraTrain,paraHeld,thirdTrain,thirdHeld,fourthTrain,fourthHeld,err:=uplm1mStartModel()
	if err!=nil { return UPLM1RGradientConflictResult{},err }

	fourFamilies:=[4][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain}
	uplm1nTrainArm(start,"rotating_palindromic_split",fourFamilies)

	fifthTrain,fifthHeld:=uplm1pCorpus()
	start,_=uplm1dExtendAlphabet(start,fifthTrain,fifthHeld)

	trains:=[5][]uplm0fExample{baseTrain,paraTrain,thirdTrain,fourthTrain,fifthTrain}
	held:=[5][]uplm0fExample{baseHeld,paraHeld,thirdHeld,fourthHeld,fifthHeld}
	for i:=1;i<5;i++ {
		if len(trains[i])!=len(trains[0]) {
			return UPLM1RGradientConflictResult{},fmt.Errorf("matched family counts differ: family0=%d family%d=%d",len(trains[0]),i,len(trains[i]))
		}
	}

	result:=UPLM1RGradientConflictResult{
		Schema:UPLM1RGradientConflictSchema,Experiment:"UP-LM1R-fivefamily-gradient-conflict",
		SourceUPLM1QSeal:"6c2111454a22feb60a36448effcc21997b4403e3",
		StateDimension:64,FamilyCount:5,PairCount:10,
		ParameterUpdatesAfterStart:false,RecurrentGradientsComputed:false,RouterUsed:false,ExactRecallUsed:false,SixthFamilyUsed:false,
	}

	var rows [5][]uplm0xGradient
	var aggs [5]uplm0xGradient
	norms:=[5]float64{}
	for f:=0;f<5;f++ {
		aggs[f]=uplm0xNewGradient(len(start.alphabet))
		rows[f]=make([]uplm0xGradient,len(trains[f]))
		for i,ex:=range trains[f] {
			g:=uplm0xSentenceGradient(start,ex.text)
			rows[f][i]=g
			uplm1rAddGradient(&aggs[f],g)
		}
		norms[f]=uplm1rGradientNorm(aggs[f])
		e:=uplm0oByteEval(start,held[f],0,uplm1qFamilyNames[f]+"_heldout")
		result.FamilyMetrics=append(result.FamilyMetrics,UPLM1RFamilyMetric{
			Family:uplm1qFamilyNames[f],TrainExamples:len(trains[f]),GradientNorm:norms[f],
			HeldoutAccuracy:e.Top1Accuracy,HeldoutPerplexity:e.Perplexity,
		})
	}

	fifthCos:=[]float64{}
	priorCos:=[]float64{}
	for a:=0;a<5;a++ {
		for b:=a+1;b<5;b++ {
			m:=uplm1rPairMetric(uplm1qFamilyNames[a],uplm1qFamilyNames[b],aggs[a],aggs[b],rows[a],rows[b])
			result.PairMetrics=append(result.PairMetrics,m)
			if b==4 { fifthCos=append(fifthCos,m.AggregateCosine) } else { priorCos=append(priorCos,m.AggregateCosine) }
		}
	}
	sum:=0.0
	minF:=math.Inf(1)
	for _,v:=range fifthCos { sum+=v;if v<minF{minF=v} }
	if len(fifthCos)>0 {
		result.FifthVsPriorMeanAggregateCosine=sum/float64(len(fifthCos))
		result.FifthVsPriorMinAggregateCosine=minF
	}
	sum=0
	for _,v:=range priorCos { sum+=v }
	if len(priorCos)>0 { result.PriorVsPriorMeanAggregateCosine=sum/float64(len(priorCos)) }
	priorNorm:=0.0
	for i:=0;i<4;i++ { priorNorm+=norms[i] }
	priorNorm/=4
	if priorNorm>0 { result.FifthGradientNormRelativeToPriorMean=norms[4]/priorNorm }
	return result,nil
}
