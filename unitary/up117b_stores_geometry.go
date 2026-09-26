package unitary

import "math"

const UP117BStoresGeometrySchema = "wingless.up117b-stores-geometry.v1"

type UP117BStaticRelation struct {
	SourceVerb string  `json:"source_verb"`
	Target     string  `json:"target"`
	Cosine     float64 `json:"cosine"`
}

type UP117BDynamicPoint struct {
	ClassOrder                 string  `json:"class_order"`
	Policy                     string  `json:"policy"`
	Stage                      int     `json:"stage"`
	Verb                       string  `json:"verb"`
	Accuracy                   float64 `json:"accuracy"`
	MeanStoreLogit             float64 `json:"mean_store_logit"`
	MeanObserveLogit           float64 `json:"mean_observe_logit"`
	MeanReportLogit            float64 `json:"mean_report_logit"`
	MeanStoreMinusObserve      float64 `json:"mean_store_minus_observe_logit_margin"`
	MinStoreMinusObserve       float64 `json:"min_store_minus_observe_logit_margin"`
	MeanDotContribution        float64 `json:"mean_store_minus_observe_dot_contribution"`
	BiasContribution           float64 `json:"store_minus_observe_bias_contribution"`
}

type UP117BStoresGeometryResult struct {
	Schema           string                  `json:"schema"`
	Experiment       string                  `json:"experiment"`
	SourceUP116BSeal string                  `json:"source_up116b_seal"`
	StateDimension   int                     `json:"state_dimension"`
	TrainingChanged  bool                    `json:"training_changed"`
	StaticRelations  []UP117BStaticRelation  `json:"static_relations"`
	DynamicPoints    []UP117BDynamicPoint    `json:"dynamic_points"`
}

var up117bVerbs=[]up106bVerbSpec{
	{verb:"stores",class:up97bStore},{verb:"keeps",class:up97bStore},
	{verb:"holds",class:up97bStore},{verb:"saves",class:up97bStore},
	{verb:"observes",class:up97bObserve},{verb:"sees",class:up97bObserve},
	{verb:"notes",class:up97bObserve},{verb:"watches",class:up97bObserve},
	{verb:"reports",class:up97bReport},{verb:"recalls",class:up97bReport},
	{verb:"tells",class:up97bReport},{verb:"remembers",class:up97bReport},
}

func up117bMeanVector(verb string)[64]float64{
	var out [64]float64
	for ni:=0;ni<6;ni++{
		h:=up106bEncode(ni,verb)
		for i:=0;i<64;i++{out[i]+=h[i]/6.0}
	}
	return out
}

func up117bCos(a,b [64]float64)float64{
	dot,aa,bb:=0.0,0.0,0.0
	for i:=0;i<64;i++{dot+=a[i]*b[i];aa+=a[i]*a[i];bb+=b[i]*b[i]}
	if aa==0||bb==0{return 0}
	return dot/math.Sqrt(aa*bb)
}

func up117bCentroid(class int,exclude string)[64]float64{
	var out [64]float64
	n:=0.0
	for _,spec:=range up117bVerbs{
		if spec.class!=class||spec.verb==exclude{continue}
		v:=up117bMeanVector(spec.verb)
		for i:=0;i<64;i++{out[i]+=v[i]}
		n++
	}
	if n>0{for i:=0;i<64;i++{out[i]/=n}}
	return out
}

func up117bStatic()[]UP117BStaticRelation{
	out:=[]UP117BStaticRelation{}
	for _,source:=range []string{"stores","keeps"}{
		sv:=up117bMeanVector(source)
		for _,spec:=range up117bVerbs{
			if spec.verb==source{continue}
			out=append(out,UP117BStaticRelation{SourceVerb:source,Target:spec.verb,Cosine:up117bCos(sv,up117bMeanVector(spec.verb))})
		}
		for _,x:=range []struct{name string;class int}{{"STORE_centroid",up97bStore},{"OBSERVE_centroid",up97bObserve},{"REPORT_centroid",up97bReport}}{
			out=append(out,UP117BStaticRelation{SourceVerb:source,Target:x.name,Cosine:up117bCos(sv,up117bCentroid(x.class,source))})
		}
	}
	return out
}

func up117bLogits(c *up97bClassifier,h [64]float64)[3]float64{
	var z [3]float64
	for k:=0;k<3;k++{
		z[k]=c.b[k]
		for i:=0;i<64;i++{z[k]+=c.w[k][i]*h[i]}
	}
	return z
}

func up117bDynamic(orderName,policy string,stage int,c *up97bClassifier,verb string)UP117BDynamicPoint{
	target:=up97bStore
	hits:=0
	sumS,sumO,sumR,sumMargin,sumDot:=0.0,0.0,0.0,0.0,0.0
	minMargin:=math.Inf(1)
	for ni:=0;ni<6;ni++{
		h:=up106bEncode(ni,verb)
		z:=up117bLogits(c,h)
		pred:=0
		if z[1]>z[pred]{pred=1}
		if z[2]>z[pred]{pred=2}
		if pred==target{hits++}
		margin:=z[up97bStore]-z[up97bObserve]
		dot:=0.0
		for i:=0;i<64;i++{dot+=(c.w[up97bStore][i]-c.w[up97bObserve][i])*h[i]}
		sumS+=z[0];sumO+=z[1];sumR+=z[2];sumMargin+=margin;sumDot+=dot
		if margin<minMargin{minMargin=margin}
	}
	return UP117BDynamicPoint{
		ClassOrder:orderName,Policy:policy,Stage:stage,Verb:verb,Accuracy:float64(hits)/6.0,
		MeanStoreLogit:sumS/6,MeanObserveLogit:sumO/6,MeanReportLogit:sumR/6,
		MeanStoreMinusObserve:sumMargin/6,MinStoreMinusObserve:minMargin,
		MeanDotContribution:sumDot/6,BiasContribution:c.b[up97bStore]-c.b[up97bObserve],
	}
}

func RunUP117B()(UP117BStoresGeometryResult,error){
	result:=UP117BStoresGeometryResult{
		Schema:UP117BStoresGeometrySchema,Experiment:"UP-117B-stores-geometry",
		SourceUP116BSeal:"73e575a82532a0886309a112c7501467d7bc5ebe",
		StateDimension:64,TrainingChanged:false,StaticRelations:up117bStatic(),
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,order:=range orders{
		name:=up114bOrderName(order)
		seq:=up114bSequence(order)
		for _,policy:=range []string{"current_class_excluded","stage3_anchor2"}{
			c:=up106bTrainBase()
			for _,verb:=range []string{"stores","keeps"}{result.DynamicPoints=append(result.DynamicPoints,up117bDynamic(name,policy,0,c,verb))}
			acquired:=[]up106bVerbSpec{}
			for idx,spec:=range seq{
				stage:=idx+1
				up113bAcquire(c,spec,acquired,policy,stage)
				acquired=append(acquired,spec)
				for _,verb:=range []string{"stores","keeps"}{result.DynamicPoints=append(result.DynamicPoints,up117bDynamic(name,policy,stage,c,verb))}
			}
		}
	}
	return result,nil
}
