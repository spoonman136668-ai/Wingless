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

$raw = & go run ./cmd/wlm-ygg-real-utf8-heterogeneous-failure-attribution-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-HETEROGENEOUS-FAILURE-ATTRIBUTION-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required=@(
 'valid_evaluation_file_count','completed_target_count','minimum_target_eligibility_count_per_file',
 'calibrated_monitor_top1_target_count','maximum_true_target_rank',
 'maximum_selected_minus_true_target_score_margin','minimum_true_target_calibrated_score',
 'minimum_exact_target_repair_accuracy_gain','maximum_exact_target_repair_global_bpb_change',
 'exact_target_repair_retain_count','maximum_counterfactual_audit_prediction_mismatch_count',
 'maximum_counterfactual_audit_probability_abs_delta','minimum_per_class_target_count',
 'record_missing_top1_target_count','successor_distribution_rotated_top1_target_count',
 'canonical_best_bucket_erased_top1_target_count','nonzero_counts_flattened_top1_target_count',
 'record_missing_max_target_rank','successor_distribution_rotated_max_target_rank',
 'canonical_best_bucket_erased_max_target_rank','nonzero_counts_flattened_max_target_rank',
 'record_missing_exact_repair_pass_count','successor_distribution_rotated_exact_repair_pass_count',
 'canonical_best_bucket_erased_exact_repair_pass_count','nonzero_counts_flattened_exact_repair_pass_count',
 'invalid_attribution_rows','training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count',
 'fresh_eval_blob_overlap_count','tokenizer_use_count','capacity_growth_event_count'
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
if($metrics.invalid_attribution_rows -ne 0){throw 'WINGLESS_ATTRIBUTION_CONSTRUCTION_INVALID'}

$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-YGG-REAL-UTF8-HETEROGENEOUS-FAILURE-ATTRIBUTION-R1'
 metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen heterogeneous-maintenance attribution completed; classify from preregistered attribution rules.'
