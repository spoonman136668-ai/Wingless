package unitary

import "math"

const UP119BShiftedMarginTopologySchema = "wingless.up119b-shifted-margin-topology.v1"

type UP119BPoint struct {
	RepresentationArm       string  `json:"representation_arm"`
	ClassOrder              string  `json:"class_order"`
	ReplayPolicy            string  `json:"replay_policy"`
	Stage                   int     `json:"stage"`
	StoresAccuracy          float64 `json:"stores_accuracy"`
	KeepsAccuracy           float64 `json:"keeps_accuracy"`
	MeanStoreLogit          float64 `json:"mean_store_logit"`
	MeanObserveLogit        float64 `json:"mean_observe_logit"`
	MeanReportLogit         float64 `json:"mean_report_logit"`
	MeanStoreObserveMargin  float64 `json:"mean_store_observe_margin"`
	MinStoreObserveMargin   float64 `json:"min_store_observe_margin"`
	MeanStoreReportMargin   float64 `json:"mean_store_report_margin"`
	MinStoreReportMargin    float64 `json:"min_store_report_margin"`
	WrongObserveCount       int     `json:"wrong_observe_count"`
	WrongReportCount        int     `json:"wrong_report_count"`
}

type UP119BShiftedMarginTopologyResult struct {
	Schema           string        `json:"schema"`
	Experiment       string        `json:"experiment"`
	SourceUP118BSeal string        `json:"source_up118b_seal"`
	StateDimension   int           `json:"state_dimension"`
	TrainingChanged  bool          `json:"training_changed"`
	Points           []UP119BPoint `json:"points"`
}

func up119bPoint(arm,order,policy string,stage int,c *up97bClassifier,delta [64]float64)UP119BPoint{
	keeps:=up106bVerbSpec{verb:"keeps",class:up97bStore}
	sumS,sumO,sumR,sumSO,sumSR:=0.0,0.0,0.0,0.0,0.0
	minSO,minSR:=math.Inf(1),math.Inf(1)
	hits,wrongO,wrongR:=0,0,0
	for ni:=0;ni<6;ni++{
		z:=up117bLogits(c,up118bEncode(arm,ni,"stores",delta))
		pred:=0
		if z[1]>z[pred]{pred=1}
		if z[2]>z[pred]{pred=2}
		if pred==up97bStore{hits++}
		if pred==up97bObserve{wrongO++}
		if pred==up97bReport{wrongR++}
		so:=z[up97bStore]-z[up97bObserve]
		sr:=z[up97bStore]-z[up97bReport]
		sumS+=z[0];sumO+=z[1];sumR+=z[2];sumSO+=so;sumSR+=sr
		if so<minSO{minSO=so};if sr<minSR{minSR=sr}
	}
	return UP119BPoint{
		RepresentationArm:arm,ClassOrder:order,ReplayPolicy:policy,Stage:stage,
		StoresAccuracy:float64(hits)/6.0,
		KeepsAccuracy:up118bVerbAccuracy(c,arm,keeps,0,6,delta),
		MeanStoreLogit:sumS/6,MeanObserveLogit:sumO/6,MeanReportLogit:sumR/6,
		MeanStoreObserveMargin:sumSO/6,MinStoreObserveMargin:minSO,
		MeanStoreReportMargin:sumSR/6,MinStoreReportMargin:minSR,
		WrongObserveCount:wrongO,WrongReportCount:wrongR,
	}
}

func RunUP119B()(UP119BShiftedMarginTopologyResult,error){
	delta:=up118bDelta()
	result:=UP119BShiftedMarginTopologyResult{
		Schema:UP119BShiftedMarginTopologySchema,Experiment:"UP-119B-shifted-margin-topology",
		SourceUP118BSeal:"3db20b638d492fb4e1d70f3765e8cf7f82ed31b5",
		StateDimension:64,TrainingChanged:false,
	}
	orders:=[][]int{
		{up97bStore,up97bObserve,up97bReport},{up97bStore,up97bReport,up97bObserve},
		{up97bObserve,up97bStore,up97bReport},{up97bObserve,up97bReport,up97bStore},
		{up97bReport,up97bStore,up97bObserve},{up97bReport,up97bObserve,up97bStore},
	}
	for _,arm:=range []string{"baseline","stores_centroid_shift"}{
		for _,order:=range orders{
			name:=up114bOrderName(order)
			seq:=up114bSequence(order)
			for _,policy:=range []string{"current_class_excluded","stage3_anchor2"}{
				c:=up118bTrainBase(arm,delta)
				result.Points=append(result.Points,up119bPoint(arm,name,policy,0,c,delta))
				acquired:=[]up106bVerbSpec{}
				for idx,spec:=range seq{
					stage:=idx+1
					up118bAcquire(c,arm,spec,acquired,policy,stage,delta)
					acquired=append(acquired,spec)
					result.Points=append(result.Points,up119bPoint(arm,name,policy,stage,c,delta))
				}
			}
		}
	}
	return result,nil
}
