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

$raw = & go run ./cmd/wlm-ygg-real-utf8-level1-maintenance-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-LEVEL1-MAINTENANCE-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required=@(
 'valid_evaluation_file_count','completed_cycle_count','minimum_target_eligibility_count_per_file',
 'minimum_monitor_target_selection_accuracy','minimum_selected_error_count_ratio_vs_runner_up',
 'maximum_wrong_candidate_selected_key_accuracy_gain','wrong_candidate_retain_count','wrong_candidate_revert_count',
 'minimum_correct_candidate_selected_key_accuracy_gain','correct_candidate_retain_count','correct_candidate_revert_count',
 'maximum_post_repair_audit_prediction_mismatch_count','maximum_post_repair_audit_probability_abs_delta',
 'minimum_structure_jaccard_after_each_cycle','maximum_candidate_read_bytes','maximum_candidate_repair_operations',
 'maximum_final_selected_motif_count','minimum_final_selected_motif_count','lesion_label_access_count',
 'evaluator_schedule_access_in_policy_count','tokenizer_use_count','capacity_growth_event_count','invalid_cycle_rows',
 'training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count','canonical_selected_motif_count'
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
if($metrics.canonical_selected_motif_count -ne 256){throw 'WINGLESS_CANONICAL_CAPACITY_INVALID'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-YGG-REAL-UTF8-LEVEL1-MAINTENANCE-R1'
 metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving real UTF8 level1 maintenance run completed; classify only from frozen thresholds.'
