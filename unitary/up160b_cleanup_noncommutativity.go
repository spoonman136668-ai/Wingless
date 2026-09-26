package unitary

const UP160BNoncommSchema="wingless.up160b-cleanup-noncommutativity.v1"

type UP160BMetric struct{
	Subject string `json:"subject"`
	TrainingPair []string `json:"training_pair"`
	Comparison string `json:"comparison"`
	OrderA string `json:"order_a"`
	OrderB string `json:"order_b"`
	OldRetentionA float64 `json:"old_retention_a"`
	OldRetentionB float64 `json:"old_retention_b"`
	RetentionDeltaAminusB float64 `json:"retention_delta_a_minus_b"`
	GateDistanceAtoB float64 `json:"gate_distance_a_to_b"`
	DisplacementNormA float64 `json:"displacement_norm_a"`
	DisplacementNormB float64 `json:"displacement_norm_b"`
}
type UP160BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP159BSeal string `json:"source_up159b_seal"`
	Subjects int `json:"subjects"`
	ComparisonsPerSubject int `json:"comparisons_per_subject"`
	CleanupUpdatesPerArm int `json:"cleanup_updates_per_arm"`
	DiagnosticUpdatesRetained bool `json:"diagnostic_updates_retained"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Metrics []UP160BMetric `json:"metrics"`
}

func up160bPostReport(common *up129bGate,subject string,o,r [64]float64)*up129bGate{
	g:=*common
	pair:=up159bPair(subject)
	reportBlock:=up154bReportBlock(subject)
	count:=0
	for _,e:=range up159bExamples(pair){
		if up153bInBlock(e,reportBlock){continue}
		up135bStep(&g,e,o,r);count++
	}
	if count!=20{panic("UP160B_PREFIX_COUNT")}
	pre:=up158bItems(15,[]int{10,11,12})
	up156bApplyOld(&g,pre,0,len(pre),o,r)
	for _,e:=range reportBlock{up135bStep(&g,e,o,r)}
	return &g
}
func up160bApplyOrder(base *up129bGate,order string,o,r [64]float64)*up129bGate{
	g:=*base
	items:=up158bItems(15,up158bIndicesForOrder(order))
	up156bApplyOld(&g,items,0,len(items),o,r)
	return &g
}
func up160bDistance(a,b *up129bGate)float64{return up150bNorm(up150bDelta(a,b))}
func up160bOld(g *up129bGate,o,r [64]float64)float64{return up143bOldMean(g,o,r)}

func RunUP160B()(UP160BResult,error){
	o,r:=up124bCompetitorDirections()
	common:=up150bCommonPrefix(o,r)
	res:=UP160BResult{
		Schema:UP160BNoncommSchema,Experiment:"UP-160B-cleanup-noncommutativity",
		SourceUP159BSeal:"c1e458477a74c553edbe5f28f8281b9575beced7",
		Subjects:6,ComparisonsPerSubject:3,CleanupUpdatesPerArm:12,
		DiagnosticUpdatesRetained:false,AdaptiveOrderingUsed:false,
	}
	comps:=[]struct{name,a,b string}{
		{"swap_store_observe","STORE_OBSERVE_REPORT","OBSERVE_STORE_REPORT"},
		{"swap_store_report","STORE_REPORT_OBSERVE","REPORT_STORE_OBSERVE"},
		{"swap_observe_report","OBSERVE_REPORT_STORE","REPORT_OBSERVE_STORE"},
	}
	for _,subject:=range up130bNewNames{
		base:=up160bPostReport(common,subject,o,r)
		for _,c:=range comps{
			a:=up160bApplyOrder(base,c.a,o,r)
			b:=up160bApplyOrder(base,c.b,o,r)
			oa,ob:=up160bOld(a,o,r),up160bOld(b,o,r)
			res.Metrics=append(res.Metrics,UP160BMetric{
				Subject:subject,TrainingPair:append([]string(nil),up159bPair(subject)...),
				Comparison:c.name,OrderA:c.a,OrderB:c.b,
				OldRetentionA:oa,OldRetentionB:ob,RetentionDeltaAminusB:oa-ob,
				GateDistanceAtoB:up160bDistance(a,b),
				DisplacementNormA:up160bDistance(base,a),DisplacementNormB:up160bDistance(base,b),
			})
		}
	}
	return res,nil
}
