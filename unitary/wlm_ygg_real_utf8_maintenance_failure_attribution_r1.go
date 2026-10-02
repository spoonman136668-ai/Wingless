package unitary

import "sort"

type wlmYggRealUTF8MaintenanceFailureAttributionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8MaintenanceFailureAttributionR1Errors(
	data []byte,
	baseline *[256][256]uint32,
	active map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
	catalog map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) map[[4]uint8]int {
	out := make(map[[4]uint8]int)
	for i := 4; i < len(data); i++ {
		key := [4]uint8{data[i-4], data[i-3], data[i-2], data[i-1]}
		if _, ok := catalog[key]; !ok {
			continue
		}
		if wlmYggRealUTF8Level1MaintenanceR3Predict(baseline, active, key) != data[i] {
			out[key]++
		}
	}
	return out
}

func wlmYggRealUTF8MaintenanceFailureAttributionR1Rank(
	counts map[[4]uint8]int,
	target [4]uint8,
) (rank int, ratio float64) {
	keys := make([][4]uint8, 0)
	for key, count := range counts {
		if count > 0 {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]]
		}
		return wlmLmRealUTF8BytePilotR2Less(keys[i], keys[j])
	})
	rank = len(keys) + 1
	for i, key := range keys {
		if key == target {
			rank = i + 1
			break
		}
	}
	targetCount := counts[target]
	runner := 0
	for key, count := range counts {
		if key == target {
			continue
		}
		if count > runner {
			runner = count
		}
	}
	if targetCount <= 0 {
		return rank, 0
	}
	if runner == 0 {
		return rank, 1e9
	}
	return rank, float64(targetCount) / float64(runner)
}

func wlmYggRealUTF8MaintenanceFailureAttributionR1Excess(
	degraded, canonical map[[4]uint8]int,
	catalog map[[4]uint8]wlmLmRealUTF8BytePilotR2Selected,
) map[[4]uint8]int {
	out := make(map[[4]uint8]int)
	for key := range catalog {
		delta := degraded[key] - canonical[key]
		if delta > 0 {
			out[key] = delta
		}
	}
	return out
}

func wlmYggRealUTF8MaintenanceFailureAttributionR1Min(a, b float64) float64 {
	if b < a { return b }
	return a
}
func wlmYggRealUTF8MaintenanceFailureAttributionR1Max(a, b float64) float64 {
	if b > a { return b }
	return a
}

// RunWlmYggRealUTF8MaintenanceFailureAttributionR1 attributes the frozen R3 real-workload maintenance failure.
func RunWlmYggRealUTF8MaintenanceFailureAttributionR1() interface{} {
	metrics := map[string]float64{
		"valid_evaluation_file_count":0,
		"completed_target_count":0,
		"minimum_target_eligibility_count_per_file":256,
		"minimum_validation_target_occurrence_count":1e9,
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"tokenizer_use_count":0,
		"capacity_growth_event_count":0,
		"invalid_attribution_rows":0,
		"raw_monitor_top1_target_count":0,
		"excess_error_monitor_top1_target_count":0,
		"maximum_raw_target_rank":0,
		"minimum_raw_target_to_runner_up_ratio":1e9,
		"maximum_excess_target_rank":0,
		"minimum_excess_target_to_runner_up_ratio":1e9,
		"minimum_correct_restoration_selected_key_accuracy_gain":1,
		"maximum_correct_restoration_selected_key_accuracy_gain":-1,
		"correct_restoration_gain_ge_0_20_count":0,
		"maximum_correct_restoration_global_bpb_change":-1e9,
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
	canonical:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	if len(canonical)!=256{metrics["invalid_attribution_rows"]++}

	for _,f:=range wlmLmRealUTF8BytePilotR2Eval{
		if trainSet[f.sha]{metrics["invalid_attribution_rows"]++}
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok{metrics["evaluation_file_identity_mismatch_count"]++}
		quarters:=wlmYggRealUTF8Level1MaintenanceR3Split(data)
		targets:=wlmYggRealUTF8Level1MaintenanceR3Targets(quarters,&baseline,canonical)
		metrics["minimum_target_eligibility_count_per_file"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Min(
			metrics["minimum_target_eligibility_count_per_file"],float64(len(targets)),
		)
		if len(targets)<2{
			metrics["invalid_attribution_rows"]++
			continue
		}
		metrics["valid_evaluation_file_count"]++

		for targetIndex:=0;targetIndex<2;targetIndex++{
			target:=targets[targetIndex].key
			degraded:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(canonical)
			delete(degraded,target)

			rawErrors:=wlmYggRealUTF8MaintenanceFailureAttributionR1Errors(
				quarters[1],&baseline,degraded,canonical,
			)
			backgroundErrors:=wlmYggRealUTF8MaintenanceFailureAttributionR1Errors(
				quarters[1],&baseline,canonical,canonical,
			)
			excess:=wlmYggRealUTF8MaintenanceFailureAttributionR1Excess(rawErrors,backgroundErrors,canonical)

			rawRank,rawRatio:=wlmYggRealUTF8MaintenanceFailureAttributionR1Rank(rawErrors,target)
			excessRank,excessRatio:=wlmYggRealUTF8MaintenanceFailureAttributionR1Rank(excess,target)
			if rawRank==1{metrics["raw_monitor_top1_target_count"]++}
			if excessRank==1{metrics["excess_error_monitor_top1_target_count"]++}
			metrics["maximum_raw_target_rank"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Max(metrics["maximum_raw_target_rank"],float64(rawRank))
			metrics["minimum_raw_target_to_runner_up_ratio"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Min(metrics["minimum_raw_target_to_runner_up_ratio"],rawRatio)
			metrics["maximum_excess_target_rank"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Max(metrics["maximum_excess_target_rank"],float64(excessRank))
			metrics["minimum_excess_target_to_runner_up_ratio"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Min(metrics["minimum_excess_target_to_runner_up_ratio"],excessRatio)

			degradedAcc,count:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(
				quarters[2],target,&baseline,degraded,
			)
			canonicalAcc,_:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(
				quarters[2],target,&baseline,canonical,
			)
			metrics["minimum_validation_target_occurrence_count"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Min(
				metrics["minimum_validation_target_occurrence_count"],float64(count),
			)
			if count<3{metrics["invalid_attribution_rows"]++}
			gain:=canonicalAcc-degradedAcc
			metrics["minimum_correct_restoration_selected_key_accuracy_gain"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Min(
				metrics["minimum_correct_restoration_selected_key_accuracy_gain"],gain,
			)
			metrics["maximum_correct_restoration_selected_key_accuracy_gain"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Max(
				metrics["maximum_correct_restoration_selected_key_accuracy_gain"],gain,
			)
			if gain>=0.20{metrics["correct_restoration_gain_ge_0_20_count"]++}
			degradedBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,degraded)
			canonicalBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,canonical)
			bpbChange:=canonicalBPB-degradedBPB
			metrics["maximum_correct_restoration_global_bpb_change"]=wlmYggRealUTF8MaintenanceFailureAttributionR1Max(
				metrics["maximum_correct_restoration_global_bpb_change"],bpbChange,
			)
			metrics["completed_target_count"]++
		}
	}

	return wlmYggRealUTF8MaintenanceFailureAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-YGG-REAL-UTF8-MAINTENANCE-FAILURE-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
