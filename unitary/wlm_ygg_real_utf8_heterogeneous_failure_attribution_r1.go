package unitary

import "sort"

type wlmYggRealUTF8HeterogeneousFailureAttributionR1Score struct {
	key   [4]uint8
	score float64
}

type wlmYggRealUTF8HeterogeneousFailureAttributionR1Result struct {
	Schema     string             `json:"schema"`
	Experiment string             `json:"experiment"`
	Metrics    map[string]float64 `json:"metrics"`
}

func wlmYggRealUTF8HeterogeneousFailureAttributionR1Rank(
	background map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Background,
	diagnostic map[[4]uint8]wlmYggRealUTF8CalibratedMaintenanceR1Diagnostic,
	target [4]uint8,
) (int,float64,float64,[4]uint8,bool) {
	rows:=make([]wlmYggRealUTF8HeterogeneousFailureAttributionR1Score,0)
	for key,bg:=range background {
		d:=diagnostic[key]
		if bg.occ<3 || d.occ<1 { continue }
		expected:=float64(bg.err)*float64(d.occ)/float64(bg.occ)
		rows=append(rows,wlmYggRealUTF8HeterogeneousFailureAttributionR1Score{
			key:key,score:float64(d.err)-expected,
		})
	}
	if len(rows)==0 { return 0,0,0,[4]uint8{},false }
	sort.Slice(rows,func(i,j int)bool{
		if rows[i].score!=rows[j].score { return rows[i].score>rows[j].score }
		return wlmLmRealUTF8BytePilotR2Less(rows[i].key,rows[j].key)
	})
	targetRank:=0
	targetScore:=0.0
	for i,row:=range rows {
		if row.key==target {
			targetRank=i+1
			targetScore=row.score
			break
		}
	}
	if targetRank==0 { return 0,0,0,rows[0].key,false }
	margin:=rows[0].score-targetScore
	return targetRank,targetScore,margin,rows[0].key,true
}

func wlmYggRealUTF8HeterogeneousFailureAttributionR1Min(a,b float64)float64{if b<a{return b};return a}
func wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(a,b float64)float64{if b>a{return b};return a}

// RunWlmYggRealUTF8HeterogeneousFailureAttributionR1 diagnoses the frozen heterogeneous-maintenance failure.
func RunWlmYggRealUTF8HeterogeneousFailureAttributionR1() interface{} {
	metrics:=map[string]float64{
		"valid_evaluation_file_count":0,
		"completed_target_count":0,
		"minimum_target_eligibility_count_per_file":256,
		"calibrated_monitor_top1_target_count":0,
		"maximum_true_target_rank":0,
		"maximum_selected_minus_true_target_score_margin":0,
		"minimum_true_target_calibrated_score":1e9,
		"minimum_exact_target_repair_accuracy_gain":1,
		"maximum_exact_target_repair_global_bpb_change":-1e9,
		"exact_target_repair_retain_count":0,
		"maximum_counterfactual_audit_prediction_mismatch_count":0,
		"maximum_counterfactual_audit_probability_abs_delta":0,
		"minimum_per_class_target_count":2,
		"record_missing_top1_target_count":0,
		"successor_distribution_rotated_top1_target_count":0,
		"canonical_best_bucket_erased_top1_target_count":0,
		"nonzero_counts_flattened_top1_target_count":0,
		"record_missing_max_target_rank":0,
		"successor_distribution_rotated_max_target_rank":0,
		"canonical_best_bucket_erased_max_target_rank":0,
		"nonzero_counts_flattened_max_target_rank":0,
		"record_missing_exact_repair_pass_count":0,
		"successor_distribution_rotated_exact_repair_pass_count":0,
		"canonical_best_bucket_erased_exact_repair_pass_count":0,
		"nonzero_counts_flattened_exact_repair_pass_count":0,
		"invalid_attribution_rows":0,
		"training_file_identity_mismatch_count":0,
		"evaluation_file_identity_mismatch_count":0,
		"fresh_eval_blob_overlap_count":0,
		"tokenizer_use_count":0,
		"capacity_growth_event_count":0,
	}
	classCounts:=[4]int{}

	train:=make([][]byte,0,len(wlmLmRealUTF8BytePilotR2Train))
	known:=make(map[string]bool)
	for _,f:=range wlmLmRealUTF8BytePilotR2Train {
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok { metrics["training_file_identity_mismatch_count"]++ }
		train=append(train,data)
		known[f.sha]=true
	}
	for _,f:=range wlmLmRealUTF8BytePilotR2Eval { known[f.sha]=true }
	baseline,_,selected:=wlmLmRealUTF8BytePilotR2TrainModel(train,metrics)
	canonical:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(selected)
	if len(canonical)!=256 { metrics["invalid_attribution_rows"]++ }

	for _,f:=range wlmYggRealUTF8HeterogeneousMaintenanceR1Eval {
		if known[f.sha] { metrics["fresh_eval_blob_overlap_count"]++ }
		data,ok:=wlmLmRealUTF8BytePilotR2LoadBlob(f)
		if !ok { metrics["evaluation_file_identity_mismatch_count"]++ }
		quarters:=wlmYggRealUTF8Level1MaintenanceR3Split(data)
		schedule,eligibilityCount:=wlmYggRealUTF8HeterogeneousMaintenanceR1Schedule(quarters,&baseline,canonical)
		metrics["minimum_target_eligibility_count_per_file"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Min(
			metrics["minimum_target_eligibility_count_per_file"],float64(eligibilityCount),
		)
		if len(schedule)!=4 {
			metrics["invalid_attribution_rows"]++
			continue
		}
		background:=wlmYggRealUTF8CalibratedMaintenanceR1BackgroundState(quarters[0],&baseline,canonical)
		if len(background)!=256 { metrics["invalid_attribution_rows"]++ }
		metrics["valid_evaluation_file_count"]++

		for classIndex,className:=range wlmYggRealUTF8HeterogeneousMaintenanceR1Classes {
			targetKey:=schedule[classIndex].key
			degraded,ok:=wlmYggRealUTF8HeterogeneousMaintenanceR1Degrade(canonical,targetKey,className)
			if !ok {
				metrics["invalid_attribution_rows"]++
				continue
			}
			diagnostic:=wlmYggRealUTF8CalibratedMaintenanceR1DiagnosticState(quarters[1],&baseline,degraded,canonical)
			rank,targetScore,margin,selectedKey,rankOK:=wlmYggRealUTF8HeterogeneousFailureAttributionR1Rank(background,diagnostic,targetKey)
			if !rankOK {
				metrics["invalid_attribution_rows"]++
				continue
			}
			if selectedKey==targetKey {
				metrics["calibrated_monitor_top1_target_count"]++
				switch className {
				case "record_missing":
					metrics["record_missing_top1_target_count"]++
				case "successor_distribution_rotated":
					metrics["successor_distribution_rotated_top1_target_count"]++
				case "canonical_best_bucket_erased":
					metrics["canonical_best_bucket_erased_top1_target_count"]++
				case "nonzero_counts_flattened":
					metrics["nonzero_counts_flattened_top1_target_count"]++
				}
			}
			metrics["maximum_true_target_rank"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["maximum_true_target_rank"],float64(rank))
			metrics["maximum_selected_minus_true_target_score_margin"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["maximum_selected_minus_true_target_score_margin"],margin)
			metrics["minimum_true_target_calibrated_score"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Min(metrics["minimum_true_target_calibrated_score"],targetScore)
			switch className {
			case "record_missing":
				metrics["record_missing_max_target_rank"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["record_missing_max_target_rank"],float64(rank))
			case "successor_distribution_rotated":
				metrics["successor_distribution_rotated_max_target_rank"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["successor_distribution_rotated_max_target_rank"],float64(rank))
			case "canonical_best_bucket_erased":
				metrics["canonical_best_bucket_erased_max_target_rank"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["canonical_best_bucket_erased_max_target_rank"],float64(rank))
			case "nonzero_counts_flattened":
				metrics["nonzero_counts_flattened_max_target_rank"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["nonzero_counts_flattened_max_target_rank"],float64(rank))
			}

			preAcc,preCount:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[2],targetKey,&baseline,degraded)
			if preCount<3 { metrics["invalid_attribution_rows"]++ }
			preBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,degraded)
			repaired:=wlmYggRealUTF8Level1MaintenanceR3CloneMap(degraded)
			repaired[targetKey]=wlmYggRealUTF8Level1MaintenanceR3CloneEntry(canonical[targetKey])
			postAcc,_:=wlmYggRealUTF8Level1MaintenanceR3KeyAccuracy(quarters[2],targetKey,&baseline,repaired)
			postBPB:=wlmYggRealUTF8Level1MaintenanceR3BPB(quarters[2],&baseline,repaired)
			gain:=postAcc-preAcc
			bpbDelta:=postBPB-preBPB
			metrics["minimum_exact_target_repair_accuracy_gain"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Min(metrics["minimum_exact_target_repair_accuracy_gain"],gain)
			metrics["maximum_exact_target_repair_global_bpb_change"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["maximum_exact_target_repair_global_bpb_change"],bpbDelta)
			retain:=gain>=0.20 && bpbDelta<=0.005
			if retain {
				metrics["exact_target_repair_retain_count"]++
				switch className {
				case "record_missing":
					metrics["record_missing_exact_repair_pass_count"]++
				case "successor_distribution_rotated":
					metrics["successor_distribution_rotated_exact_repair_pass_count"]++
				case "canonical_best_bucket_erased":
					metrics["canonical_best_bucket_erased_exact_repair_pass_count"]++
				case "nonzero_counts_flattened":
					metrics["nonzero_counts_flattened_exact_repair_pass_count"]++
				}
			}
			predMismatch,probDelta:=wlmYggRealUTF8Level1MaintenanceR3Audit(quarters[3],&baseline,repaired,canonical)
			metrics["maximum_counterfactual_audit_prediction_mismatch_count"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["maximum_counterfactual_audit_prediction_mismatch_count"],float64(predMismatch))
			metrics["maximum_counterfactual_audit_probability_abs_delta"]=wlmYggRealUTF8HeterogeneousFailureAttributionR1Max(metrics["maximum_counterfactual_audit_probability_abs_delta"],probDelta)
			classCounts[classIndex]++
			metrics["completed_target_count"]++
		}
	}
	minClass:=classCounts[0]
	for _,count:=range classCounts[1:] { if count<minClass { minClass=count } }
	metrics["minimum_per_class_target_count"]=float64(minClass)
	return wlmYggRealUTF8HeterogeneousFailureAttributionR1Result{
		Schema:"wingless.research-scientific-result.v1",
		Experiment:"WLM-YGG-REAL-UTF8-HETEROGENEOUS-FAILURE-ATTRIBUTION-R1",
		Metrics:metrics,
	}
}
