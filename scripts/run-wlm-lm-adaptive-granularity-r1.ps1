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

$raw = & go run ./cmd/wlm-lm-adaptive-granularity-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-LM-ADAPTIVE-GRANULARITY-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }
if ($null -eq $probe.metrics) { throw 'WINGLESS_RESULT_METRICS_MISSING' }

$required = @(
  'valid_seed_count',
  'completed_substrate_runs',
  'minimum_adaptive_target_accuracy',
  'minimum_full_byte_target_accuracy',
  'maximum_adaptive_accuracy_gap_vs_full_byte',
  'maximum_static_fixed4_target_accuracy',
  'minimum_adaptive_accuracy_gain_vs_static_fixed4',
  'minimum_overall_event_reduction_fraction',
  'minimum_predictable_region_event_reduction_fraction',
  'maximum_high_entropy_false_coarse_activation_rate',
  'minimum_predictable_coarse_activation_recall',
  'maximum_adaptive_events_per_raw_byte',
  'substrate_selected_motif_mismatch_count',
  'substrate_prediction_mismatch_count',
  'substrate_event_accounting_mismatch_count',
  'invalid_stream_rows',
  'counter_overflow_rows'
)
$metrics=[ordered]@{}
foreach($name in $required){
  $property=$probe.metrics.PSObject.Properties[$name]
  if($null -eq $property){throw "WINGLESS_METRIC_MISSING:$name"}
  $value=[double]$property.Value
  if([double]::IsNaN($value)-or[double]::IsInfinity($value)){throw "WINGLESS_METRIC_NONFINITE:$name"}
  $metrics[$name]=$value
}
$result=[ordered]@{
  schema='wingless.research-scientific-result.v1'
  experiment='WLM-LM-ADAPTIVE-GRANULARITY-R1'
  metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving WBG-3 adaptive-granularity run completed.'
