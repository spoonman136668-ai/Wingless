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

$raw = & go run ./cmd/wlm-ygg-real-utf8-maintenance-failure-attribution-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-MAINTENANCE-FAILURE-ATTRIBUTION-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required=@(
 'valid_evaluation_file_count','completed_target_count','minimum_target_eligibility_count_per_file',
 'minimum_validation_target_occurrence_count','training_file_identity_mismatch_count',
 'evaluation_file_identity_mismatch_count','tokenizer_use_count','capacity_growth_event_count',
 'invalid_attribution_rows','raw_monitor_top1_target_count','excess_error_monitor_top1_target_count',
 'maximum_raw_target_rank','minimum_raw_target_to_runner_up_ratio','maximum_excess_target_rank',
 'minimum_excess_target_to_runner_up_ratio','minimum_correct_restoration_selected_key_accuracy_gain',
 'maximum_correct_restoration_selected_key_accuracy_gain','correct_restoration_gain_ge_0_20_count',
 'maximum_correct_restoration_global_bpb_change'
)
$metrics=[ordered]@{}
foreach($name in $required){
 $p=$probe.metrics.PSObject.Properties[$name]
 if($null -eq $p){throw "WINGLESS_METRIC_MISSING:$name"}
 $v=[double]$p.Value
 if([double]::IsNaN($v)-or[double]::IsInfinity($v)){throw "WINGLESS_METRIC_NONFINITE:$name"}
 $metrics[$name]=$v
}
$result=[ordered]@{
 schema='wingless.research-scientific-result.v1'
 experiment='WLM-YGG-REAL-UTF8-MAINTENANCE-FAILURE-ATTRIBUTION-R1'
 metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen real UTF8 maintenance failure attribution completed.'
