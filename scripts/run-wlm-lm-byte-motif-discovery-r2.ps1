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

$raw = & go run ./cmd/wlm-lm-byte-motif-discovery-r2
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-LM-BYTE-MOTIF-DISCOVERY-R2') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }
if ($null -eq $probe.metrics) { throw 'WINGLESS_RESULT_METRICS_MISSING' }

$required = @(
  'valid_seed_count',
  'completed_substrate_runs',
  'selected_motif_count',
  'minimum_selected_true_motif_recall',
  'maximum_false_selected_motif_count',
  'minimum_motif_enabled_successor_accuracy',
  'maximum_byte_only_ablation_successor_accuracy',
  'maximum_shuffled_motif_successor_accuracy',
  'minimum_causal_accuracy_gain_over_ablation',
  'minimum_effective_event_reduction_fraction',
  'substrate_selected_motif_mismatch_count',
  'substrate_prediction_mismatch_count',
  'substrate_event_accounting_mismatch_count',
  'duplicate_shuffled_control_key_count',
  'invalid_candidate_rows',
  'counter_overflow_rows'
)

$metrics = [ordered]@{}
foreach ($name in $required) {
  $property = $probe.metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  $metrics[$name] = $value
}

$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-LM-BYTE-MOTIF-DISCOVERY-R2'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving WBG-2 run completed; classify only from frozen preregistered thresholds.'
