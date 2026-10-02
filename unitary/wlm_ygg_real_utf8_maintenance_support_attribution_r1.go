package unitary

import "fmt"

type wlmYggRealUTF8MaintenanceSupportAttributionR1Tier struct {
	calibrationMin int
	diagnosticMin  int
	validationMin  int
	auditMin       int
	gainMin        float64
}

type wlmYggRealUTF8MaintenanceSupportAttributionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8MaintenanceSupportAttributionR1EligibleCount(
	quarters [4][]byte,
	baseline *[256][256]uint32,
	canonical map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	tier wlmYggRealUTF8MaintenanceSupportAttributionR1Tier,
) int {
	cal := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[0], baseline, canonical)
	diag := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[1], baseline, canonical)
	val := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[2], baseline, canonical)
	audit := wlmYggRealUTF8Level1MaintenanceR1SupportFor(quarters[3], baseline, canonical)
	count := 0
	for key, row := range cal {
		if row.occ < tier.calibrationMin ||
			diag[key].occ < tier.diagnosticMin ||
			val[key].occ < tier.validationMin ||
			audit[key].occ < tier.auditMin {
			continue
		}
		gain := float64(row.canonicalHit-row.baselineHit) / float64(row.occ)
		if gain < tier.gainMin {
			continue
		}
		count++
	}
	return count
}

// RunWlmYggRealUTF8MaintenanceSupportAttributionR1 evaluates the frozen target-support tiers.
func RunWlmYggRealUTF8MaintenanceSupportAttributionR1() interface{} {
	tiers := [...]wlmYggRealUTF8MaintenanceSupportAttributionR1Tier{
		{8,4,4,4,0.25},
		{6,3,3,3,0.20},
		{4,2,2,2,0.15},
		{3,2,2,2,0.10},
		{2,1,1,1,0.10},
	}
	metrics := map[string]float64{
		"valid_evaluation_file_count":0,
		"canonical_selected_motif_count":0,
		"strictest_feasible_tier_index":-1,
		"minimum_selected_tier_eligible_count_per_file":0,
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"tokenizer_use_count":0,
		"capacity_growth_event_count":0,
		"invalid_attribution_rows":0,
	}
	for fileIndex:=0;fileIndex<2;fileIndex++ {
		for tierIndex:=0;tierIndex<5;tierIndex++ {
			metrics[fmt.Sprintf("file%d_tier%d_eligible_count",fileIndex,tierIndex)] = 0
		}
	}
	for tierIndex:=0;tierIndex<5;tierIndex++ {
		metrics[fmt.Sprintf("minimum_tier%d_eligible_count_per_file",tierIndex)] = 1e9
	}

	train:=make([][]byte,0,len(wlmLmRealUTF8BytePilotR2Train))
	trainSet:=make(map[string]bool)
	for _,f:=range wlmLmRealUTF8BytePilotR2Train{
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["training_file_identity_mismatch_count"]++}
		train=append(train,data)
		trainSet[f.sha]=true
	}
	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	canonical:=wlmYggRealUTF8Level1MaintenanceR1CloneMap(selected)
	metrics["canonical_selected_motif_count"]=float64(len(canonical))
	if len(canonical)!=256{metrics["invalid_attribution_rows"]++}

	fileTierCounts := make([][5]int,0,2)
	for fileIndex,f:=range wlmLmRealUTF8BytePilotR2Eval {
		if trainSet[f.sha]{metrics["invalid_attribution_rows"]++}
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		quarters:=wlmYggRealUTF8Level1MaintenanceR1Split(data)
		var counts [5]int
		for tierIndex,tier:=range tiers {
			count:=wlmYggRealUTF8MaintenanceSupportAttributionR1EligibleCount(quarters,&baseline,canonical,tier)
			counts[tierIndex]=count
			metrics[fmt.Sprintf("file%d_tier%d_eligible_count",fileIndex,tierIndex)] = float64(count)
			name:=fmt.Sprintf("minimum_tier%d_eligible_count_per_file",tierIndex)
			if float64(count)<metrics[name]{metrics[name]=float64(count)}
		}
		fileTierCounts=append(fileTierCounts,counts)
		metrics["valid_evaluation_file_count"]++
	}
	if len(fileTierCounts)!=2 {
		metrics["invalid_attribution_rows"]++
	} else {
		for tierIndex:=0;tierIndex<5;tierIndex++ {
			minCount:=fileTierCounts[0][tierIndex]
			if fileTierCounts[1][tierIndex]<minCount{minCount=fileTierCounts[1][tierIndex]}
			if minCount>=2 {
				metrics["strictest_feasible_tier_index"]=float64(tierIndex)
				metrics["minimum_selected_tier_eligible_count_per_file"]=float64(minCount)
				break
			}
		}
	}
	return wlmYggRealUTF8MaintenanceSupportAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-YGG-REAL-UTF8-MAINTENANCE-SUPPORT-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
