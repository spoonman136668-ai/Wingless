$ErrorActionPreference = 'Stop'
$env:GOFLAGS = '-buildvcs=false'
$env:GOMAXPROCS = '1'

& go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_FULL_REGRESSION_FAILED' }

$generated = @('.ice/architecture-map.json','.ice/dependency-map.json','.ice/integration-seams.json','.ice/manifest.json','.ice/symbol-map.json','.ice/test-map.json')
foreach ($path in $generated) {
  & git checkout -- $path 2>$null
  if ($LASTEXITCODE -ne 0) { if (Test-Path -LiteralPath $path) { Remove-Item -LiteralPath $path -Force -Recurse } }
}

$raw = & go run ./cmd/wlm-ygg-real-utf8-paired-delta-monitor-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-PAIRED-DELTA-MONITOR-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required = @('valid_evaluation_file_count','completed_row_count','minimum_target_eligibility_count_per_file','legacy_misrank_count','legacy_nonzero_counts_flattened_misrank_count','legacy_other_class_misrank_count','paired_delta_top1_target_count','paired_delta_nonpositive_target_count','paired_delta_positive_off_target_count','minimum_paired_delta_target_margin','invalid_row_count','training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count','fresh_eval_blob_overlap_count','tokenizer_use_count','capacity_growth_event_count')
$metrics = [ordered]@{}
foreach ($name in $required) {
  $property = $probe.metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  $metrics[$name] = $value
}
if ($metrics.training_file_identity_mismatch_count -ne 0) { throw 'WINGLESS_TRAIN_BLOB_IDENTITY_INVALID' }
if ($metrics.evaluation_file_identity_mismatch_count -ne 0) { throw 'WINGLESS_EVAL_BLOB_IDENTITY_INVALID' }
if ($metrics.fresh_eval_blob_overlap_count -ne 0) { throw 'WINGLESS_FRESH_EVAL_OVERLAP_INVALID' }
if ($metrics.invalid_row_count -ne 0) { throw 'WINGLESS_PAIRED_DELTA_CONSTRUCTION_INVALID' }

$result = [ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-YGG-REAL-UTF8-PAIRED-DELTA-MONITOR-R1';metrics=$metrics}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen legacy and paired-delta monitor readouts completed; classify from preregistered rules.'
