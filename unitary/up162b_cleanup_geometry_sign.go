package unitary

const UP162BGeometrySchema="wingless.up162b-cleanup-geometry-sign.v1"

type UP162BMetric struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	ParentRetentionDeltaSORminusOSR float64 `json:"parent_retention_delta_sor_minus_osr"`
	StoreNorm float64 `json:"store_norm"`
	ObserveNorm float64 `json:"observe_norm"`
	ReportNorm float64 `json:"report_norm"`
	OldReferenceNorm float64 `json:"old_reference_norm"`
	CosineStoreObserve float64 `json:"cosine_store_observe"`
	CosineStoreReport float64 `json:"cosine_store_report"`
	CosineObserveReport float64 `json:"cosine_observe_report"`
	CosineStoreOldReference float64 `json:"cosine_store_old_reference"`
	CosineObserveOldReference float64 `json:"cosine_observe_old_reference"`
	CosineReportOldReference float64 `json:"cosine_report_old_reference"`
	StoreMinusObserveOldAlignment float64 `json:"store_minus_observe_old_alignment"`
	CommutatorNorm float64 `json:"commutator_norm"`
	CommutatorCosineOldReference float64 `json:"commutator_cosine_old_reference"`
}
type UP162BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP161BSeal string `json:"source_up161b_seal"`
	MechanismUP160BSeal string `json:"mechanism_up160b_seal"`
	Subjects int `json:"subjects"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	PostResultClassifierUsed bool `json:"post_result_classifier_used"`
	Metrics []UP162BMetric `json:"metrics"`
}

func up162bFrozenItemDelta(base *up129bGate,item up156bOldItem,o,r [64]float64)[]float64{
	g:=*base
	up129bGateStep(&g,up121bOriginalNames[item.subjectIndex],item.surface.verb,item.surface.class,o,r)
	return up150bDelta(base,&g)
}
func up162bAggregate(base *up129bGate,indices []int,o,r [64]float64)[]float64{
	items:=up156bOldItems(15)
	out:=make([]float64,len(up150bGateVector(base)))
	for _,idx:=range indices{up150bAdd(out,up162bFrozenItemDelta(base,items[idx],o,r))}
	return out
}
func up162bAllIndices()[]int{out:=make([]int,15);for i:=range out{out[i]=i};return out}
func up162bParentDelta(subject string)float64{
	switch subject{
	case "mia","noah","pax":return 0.005555555555555536
	case "opal","rue":return -0.011111111111111072
	default:return -0.022222222222222254
	}
}
func RunUP162B()(UP162BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	res:=UP162BResult{
		Schema:UP162BGeometrySchema,Experiment:"UP-162B-cleanup-geometry-sign",
		SourceUP161BSeal:"4678afd1a0338aa849a4185d7d02d90166cd7b10",
		MechanismUP160BSeal:"d4a12a5ecd2d5f24671d1af3a2a7267732f10fdd",
		Subjects:6,DiagnosticUpdatesRetained:false,AdaptiveOrderingUsed:false,PostResultClassifierUsed:false,
	}
	storeIdx:=[]int{0,1,2,3,4};observeIdx:=[]int{5,6,7,8,9};reportIdx:=[]int{13,14}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r)
		sv:=up162bAggregate(base,storeIdx,o,r)
		ov:=up162bAggregate(base,observeIdx,o,r)
		rv:=up162bAggregate(base,reportIdx,o,r)
		oldv:=up162bAggregate(base,up162bAllIndices(),o,r)
		sor:=up160bApplyOrder(base,"STORE_OBSERVE_REPORT",o,r)
		osr:=up160bApplyOrder(base,"OBSERVE_STORE_REPORT",o,r)
		comm:=up150bDelta(osr,sor)
		cs,co,cr:=up150bCos(sv,oldv),up150bCos(ov,oldv),up150bCos(rv,oldv)
		res.Metrics=append(res.Metrics,UP162BMetric{
			Subject:subject,TrainingPair:append([]string(nil),up159bPair(subject)...),
			ParentRetentionDeltaSORminusOSR:up162bParentDelta(subject),
			StoreNorm:up150bNorm(sv),ObserveNorm:up150bNorm(ov),ReportNorm:up150bNorm(rv),OldReferenceNorm:up150bNorm(oldv),
			CosineStoreObserve:up150bCos(sv,ov),CosineStoreReport:up150bCos(sv,rv),CosineObserveReport:up150bCos(ov,rv),
			CosineStoreOldReference:cs,CosineObserveOldReference:co,CosineReportOldReference:cr,
			StoreMinusObserveOldAlignment:cs-co,
			CommutatorNorm:up150bNorm(comm),CommutatorCosineOldReference:up150bCos(comm,oldv),
		})
	}
	return res,nil
}
