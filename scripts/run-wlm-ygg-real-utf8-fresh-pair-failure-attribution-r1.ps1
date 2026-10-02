$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$raw=& go run ./cmd/wlm-ygg-real-utf8-fresh-pair-failure-attribution-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-YGG-REAL-UTF8-FRESH-PAIR-FAILURE-ATTRIBUTION-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
$required=@('valid_evaluation_file_count','completed_independent_row_count','minimum_target_eligibility_count_per_file','minimum_per_class_completed_row_count','monitor_top1_target_count','minimum_true_target_rank','maximum_true_target_rank','minimum_selected_paired_delta_score','minimum_selected_minus_runner_up_paired_delta','exact_target_repair_retain_count','exact_target_repair_revert_count','exact_target_repair_retain_count_plus_revert_count','minimum_exact_target_selected_key_accuracy_gain','maximum_exact_target_global_bpb_change','maximum_exact_target_post_repair_audit_prediction_mismatch_count','maximum_exact_target_post_repair_audit_probability_abs_delta','sequential_control_valid_evaluation_file_count','sequential_control_completed_cycle_count','sequential_control_minimum_monitor_target_selection_accuracy','sequential_control_wrong_candidate_retain_count','sequential_control_wrong_candidate_revert_count','sequential_control_correct_candidate_retain_count','sequential_control_correct_candidate_revert_count','sequential_control_maximum_post_repair_audit_prediction_mismatch_count','sequential_control_maximum_post_repair_audit_probability_abs_delta','sequential_control_minimum_structure_jaccard_after_each_cycle','maximum_candidate_read_bytes','maximum_candidate_repair_operations','canonical_selected_motif_count','training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count','fresh_eval_blob_overlap_count','invalid_attribution_rows','lesion_label_access_in_monitor_count','class_label_access_in_monitor_count','tokenizer_use_count','capacity_growth_event_count')
$metrics=[ordered]@{}
foreach($name in $required){$p=$probe.metrics.PSObject.Properties[$name];if($null-eq$p){throw "WINGLESS_METRIC_MISSING:$name"};$v=[double]$p.Value;if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"};$metrics[$name]=$v}
if($metrics.valid_evaluation_file_count-ne 2-or$metrics.completed_independent_row_count-ne 8-or$metrics.minimum_per_class_completed_row_count-ne 2){throw 'WINGLESS_DIAGNOSTIC_CARDINALITY_INVALID'}
if($metrics.minimum_target_eligibility_count_per_file-lt 4-or$metrics.minimum_true_target_rank-lt 1-or$metrics.maximum_true_target_rank-gt 256){throw 'WINGLESS_DIAGNOSTIC_RANGE_INVALID'}
if($metrics.monitor_top1_target_count-lt 0-or$metrics.monitor_top1_target_count-gt 8-or$metrics.exact_target_repair_retain_count_plus_revert_count-ne 8){throw 'WINGLESS_ATTRIBUTION_CARDINALITY_INVALID'}
if($metrics.sequential_control_valid_evaluation_file_count-ne 2-or$metrics.sequential_control_completed_cycle_count-ne 8-or$metrics.sequential_control_minimum_monitor_target_selection_accuracy-ne 0){throw 'WINGLESS_CONTROL_CARDINALITY_INVALID'}
if($metrics.sequential_control_wrong_candidate_retain_count-ne 0-or$metrics.sequential_control_wrong_candidate_revert_count-ne 8-or$metrics.sequential_control_correct_candidate_retain_count-ne 4-or$metrics.sequential_control_correct_candidate_revert_count-ne 4){throw 'WINGLESS_CONTROL_REPAIR_INVALID'}
if($metrics.sequential_control_maximum_post_repair_audit_prediction_mismatch_count-ne 113-or$metrics.sequential_control_maximum_post_repair_audit_probability_abs_delta-ne 0.40425531914893614-or$metrics.sequential_control_minimum_structure_jaccard_after_each_cycle-ne 1){throw 'WINGLESS_CONTROL_AUDIT_INVALID'}
if($metrics.maximum_candidate_read_bytes-gt 1028-or$metrics.maximum_candidate_repair_operations-gt 1-or$metrics.canonical_selected_motif_count-ne 256){throw 'WINGLESS_BUDGET_INVALID'}
if($metrics.training_file_identity_mismatch_count-ne 0-or$metrics.evaluation_file_identity_mismatch_count-ne 0-or$metrics.fresh_eval_blob_overlap_count-ne 0-or$metrics.invalid_attribution_rows-ne 0){throw 'WINGLESS_VALIDITY_INVALID'}
if($metrics.lesion_label_access_in_monitor_count-ne 0-or$metrics.class_label_access_in_monitor_count-ne 0-or$metrics.tokenizer_use_count-ne 0-or$metrics.capacity_growth_event_count-ne 0){throw 'WINGLESS_CONTROL_INVALID'}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-YGG-REAL-UTF8-FRESH-PAIR-FAILURE-ATTRIBUTION-R1';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Fresh-pair failure attribution completed; classify only from frozen rules.'
