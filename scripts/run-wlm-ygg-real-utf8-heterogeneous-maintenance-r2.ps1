$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
$env:GOMAXPROCS = '1'

& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }

$generated = @(
  '.ice/architecture-map.json',
  '.ice/dependency-map.json',
  '.ice/integration-seams.json',
  '.ice/manifest.json',
  '.ice/symbol-map.json',
  '.ice/test-map.json'
)
foreach ($path in $generated) {
  & git checkout -- $path 2>$null
  if ($LASTEXITCODE -ne 0) {
    if (Test-Path -LiteralPath $path) { Remove-Item -LiteralPath $path -Force -Recurse }
  }
}

$raw = & go run ./cmd/wlm-ygg-real-utf8-heterogeneous-maintenance-r2
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-HETEROGENEOUS-MAINTENANCE-R2') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required=@(
 'valid_evaluation_file_count','completed_cycle_count','minimum_target_eligibility_count_per_file',
 'minimum_monitor_target_selection_accuracy','minimum_selected_paired_delta_score',
 'minimum_selected_minus_runner_up_paired_delta','maximum_wrong_candidate_selected_key_accuracy_gain',
 'wrong_candidate_retain_count','wrong_candidate_revert_count','minimum_correct_candidate_selected_key_accuracy_gain',
 'correct_candidate_retain_count','correct_candidate_revert_count','minimum_per_class_completed_cycle_count',
 'maximum_post_repair_audit_prediction_mismatch_count','maximum_post_repair_audit_probability_abs_delta',
 'minimum_structure_jaccard_after_each_cycle','maximum_candidate_read_bytes','maximum_candidate_repair_operations',
 'maximum_final_selected_motif_count','minimum_final_selected_motif_count','lesion_label_access_count',
 'class_label_access_in_policy_count','healthy_profile_recompute_count','live_canonical_shadow_access_count',
 'evaluator_schedule_access_in_policy_count','tokenizer_use_count','capacity_growth_event_count',
 'invalid_cycle_rows','training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count',
 'fresh_eval_blob_overlap_count','canonical_selected_motif_count'
)
$metrics=[ordered]@{}
foreach($name in $required){
 $p=$probe.metrics.PSObject.Properties[$name]
 if($null -eq $p){throw "WINGLESS_METRIC_MISSING:$name"}
 $v=[double]$p.Value
 if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
 $metrics[$name]=$v
}
if($metrics.training_file_identity_mismatch_count -ne 0){throw 'WINGLESS_TRAIN_BLOB_IDENTITY_INVALID'}
if($metrics.evaluation_file_identity_mismatch_count -ne 0){throw 'WINGLESS_EVAL_BLOB_IDENTITY_INVALID'}
if($metrics.fresh_eval_blob_overlap_count -ne 0){throw 'WINGLESS_FRESH_EVAL_OVERLAP_INVALID'}
if($metrics.canonical_selected_motif_count -ne 256){throw 'WINGLESS_CANONICAL_CAPACITY_INVALID'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-YGG-REAL-UTF8-HETEROGENEOUS-MAINTENANCE-R2'
 metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving paired-delta heterogeneous maintenance run completed; classify only from frozen thresholds.'
