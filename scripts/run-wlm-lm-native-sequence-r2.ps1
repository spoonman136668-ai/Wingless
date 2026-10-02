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

$raw = & go run ./cmd/wlm-lm-native-sequence-r2
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$exact = [ordered]@{
  valid_seed_count = 8
  completed_substrate_runs = 16
  training_transition_count = 32768
  evaluation_transition_count = 16384
  substrate_prediction_mismatch_count = 0
  substrate_learned_state_mismatch_count = 0
  substrate_context_observable_mismatch_count = 0
  invalid_context_rows = 0
  counter_overflow_rows = 0
}
foreach ($name in $exact.Keys) {
  $property = $probe.metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  if ($value -ne [double]$exact[$name]) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:$name actual=$value expected=$($exact[$name])" }
}

$trained = [double]$probe.metrics.minimum_trained_next_token_accuracy
$baseline = [double]$probe.metrics.maximum_untrained_baseline_accuracy
$improvement = [double]$probe.metrics.minimum_accuracy_improvement_over_untrained
if ($trained -lt 0.95) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:minimum_trained_next_token_accuracy actual=$trained threshold=0.95" }
if ($baseline -gt 0.20) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:maximum_untrained_baseline_accuracy actual=$baseline threshold=0.20" }
if ($improvement -lt 0.70) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:minimum_accuracy_improvement_over_untrained actual=$improvement threshold=0.70" }

$metrics = [ordered]@{}
foreach ($property in $probe.metrics.PSObject.Properties) {
  $metrics[$property.Name] = [double]$property.Value
}
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-LM-NATIVE-SEQUENCE-R2'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen native sequence-learning and autoregressive-generation gates passed.'
