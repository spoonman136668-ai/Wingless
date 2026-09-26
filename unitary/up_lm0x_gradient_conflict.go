package unitary

import (
	"fmt"
	"math"
)

const UPLM0XGradientConflictSchema = "wingless.up-lm0x-gradient-conflict.v1"

type UPLM0XRowMetric struct {
	Byte             int     `json:"byte"`
	Cosine           float64 `json:"cosine"`
	DotProduct       float64 `json:"dot_product"`
	BaseGradientNorm float64 `json:"base_gradient_norm"`
	ParaGradientNorm float64 `json:"paraphrase_gradient_norm"`
}

type UPLM0XGradientConflictResult struct {
	Schema                    string              `json:"schema"`
	Experiment                string              `json:"experiment"`
	SourceUPLM0WSeal          string              `json:"source_up_lm0w_seal"`
	StateDimension            int                 `json:"state_dimension"`
	BaseEpochs                int                 `json:"base_epochs"`
	LearningRate              float64             `json:"learning_rate"`
	AdaptationUsed            bool                `json:"adaptation_used"`
	ParameterUpdatesAfterBase bool                `json:"parameter_updates_after_base"`
	PairCount                 int                 `json:"pair_count"`
	AlphabetBytes             []int               `json:"alphabet_bytes"`
	MeanPairCosine            float64             `json:"mean_pair_cosine"`
	MinPairCosine             float64             `json:"min_pair_cosine"`
	MaxPairCosine             float64             `json:"max_pair_cosine"`
	NegativePairFraction      float64             `json:"negative_pair_fraction"`
	NonPositivePairFraction   float64             `json:"non_positive_pair_fraction"`
	GlobalCosine              float64             `json:"global_cosine"`
	RowMetrics                []UPLM0XRowMetric   `json:"row_metrics"`
}

type uplm0xGradient struct {
	w [][]float64
	b []float64
}

func uplm0xNewGradient(rows int) uplm0xGradient {
	g:=uplm0xGradient{w:make([][]float64,rows),b:make([]float64,rows)}
	for c:=range g.w { g.w[c]=make([]float64,64) }
	return g
}

func uplm0xSentenceGradient(model *uplm0aModel,s string) uplm0xGradient {
	g:=uplm0xNewGradient(len(model.alphabet))
	var h [64]float64
	for t:=0;t<len(s)-1;t++ {
		h=uplm0aStep(h,s[t])
		target:=model.index[int(s[t+1])]
		p:=model.probs(h)
		for c:=range p {
			err:=p[c]
			if c==target { err-=1 }
			for i:=0;i<64;i++ { g.w[c][i]+=err*h[i] }
			g.b[c]+=err
		}
	}
	return g
}

func uplm0xPairStats(a,b uplm0xGradient)(dot,normA2,normB2 float64) {
	for c:=range a.w {
		for i:=0;i<64;i++ {
			dot+=a.w[c][i]*b.w[c][i]
			normA2+=a.w[c][i]*a.w[c][i]
			normB2+=b.w[c][i]*b.w[c][i]
		}
		dot+=a.b[c]*b.b[c]
		normA2+=a.b[c]*a.b[c]
		normB2+=b.b[c]*b.b[c]
	}
	return
}

func uplm0xCosine(dot,normA2,normB2 float64) float64 {
	if normA2<=0 || normB2<=0 { return 0 }
	return dot/math.Sqrt(normA2*normB2)
}

func RunUPLM0X()(UPLM0XGradientConflictResult,error){
	baseTrain,_,alphabet:=uplm0fCorpus()
	paraTrain,_:=uplm0oParaphraseCorpus()
	result:=UPLM0XGradientConflictResult{
		Schema:UPLM0XGradientConflictSchema,
		Experiment:"UP-LM0X-gradient-conflict",
		SourceUPLM0WSeal:"205791da63c3ce4a214d44006ab77564e93c305e",
		StateDimension:64,BaseEpochs:20,LearningRate:0.08,
		AdaptationUsed:false,ParameterUpdatesAfterBase:false,
	}
	if len(baseTrain)!=len(paraTrain) {
		return result,fmt.Errorf("matched corpus count mismatch: base=%d paraphrase=%d",len(baseTrain),len(paraTrain))
	}

	model:=newUPLM0AModel(alphabet)
	for epoch:=0;epoch<20;epoch++ {
		for _,ex:=range baseTrain { model.trainSentence(ex.text,0.08) }
	}
	for _,b:=range model.alphabet { result.AlphabetBytes=append(result.AlphabetBytes,int(b)) }

	rows:=len(model.alphabet)
	rowDot:=make([]float64,rows)
	rowBase2:=make([]float64,rows)
	rowPara2:=make([]float64,rows)
	minCos:=math.Inf(1)
	maxCos:=math.Inf(-1)
	sumCos:=0.0
	negative,nonPositive:=0,0
	globalDot,globalBase2,globalPara2:=0.0,0.0,0.0

	for i:=range baseTrain {
		bg:=uplm0xSentenceGradient(model,baseTrain[i].text)
		pg:=uplm0xSentenceGradient(model,paraTrain[i].text)
		dot,bn2,pn2:=uplm0xPairStats(bg,pg)
		cos:=uplm0xCosine(dot,bn2,pn2)
		if cos<minCos { minCos=cos }
		if cos>maxCos { maxCos=cos }
		sumCos+=cos
		if cos<0 { negative++ }
		if cos<=0 { nonPositive++ }
		globalDot+=dot
		globalBase2+=bn2
		globalPara2+=pn2

		for c:=0;c<rows;c++ {
			for j:=0;j<64;j++ {
				rowDot[c]+=bg.w[c][j]*pg.w[c][j]
				rowBase2[c]+=bg.w[c][j]*bg.w[c][j]
				rowPara2[c]+=pg.w[c][j]*pg.w[c][j]
			}
			rowDot[c]+=bg.b[c]*pg.b[c]
			rowBase2[c]+=bg.b[c]*bg.b[c]
			rowPara2[c]+=pg.b[c]*pg.b[c]
		}
	}

	result.PairCount=len(baseTrain)
	if result.PairCount>0 {
		result.MeanPairCosine=sumCos/float64(result.PairCount)
		result.MinPairCosine=minCos
		result.MaxPairCosine=maxCos
		result.NegativePairFraction=float64(negative)/float64(result.PairCount)
		result.NonPositivePairFraction=float64(nonPositive)/float64(result.PairCount)
	}
	result.GlobalCosine=uplm0xCosine(globalDot,globalBase2,globalPara2)

	for c,b:=range model.alphabet {
		result.RowMetrics=append(result.RowMetrics,UPLM0XRowMetric{
			Byte:int(b),
			Cosine:uplm0xCosine(rowDot[c],rowBase2[c],rowPara2[c]),
			DotProduct:rowDot[c],
			BaseGradientNorm:math.Sqrt(rowBase2[c]),
			ParaGradientNorm:math.Sqrt(rowPara2[c]),
		})
	}
	return result,nil
}
