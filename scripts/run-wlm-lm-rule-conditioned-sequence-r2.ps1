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

$raw = & go run ./cmd/wlm-lm-rule-conditioned-sequence-r2
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$exact = [ordered]@{
  valid_seed_count = 8
  completed_substrate_runs = 16
  training_transition_count = 49152
  heldout_one_step_evaluation_count = 1024
  autoregressive_evaluation_transition_count = 16384
  substrate_prediction_mismatch_count = 0
  substrate_factorized_state_mismatch_count = 0
  substrate_context_observable_mismatch_count = 0
  heldout_rule_context_leak_count = 0
  factorized_bit_coverage_missing_count = 0
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

$one = [double]$probe.metrics.minimum_factorized_heldout_one_step_accuracy
$lookup = [double]$probe.metrics.maximum_pooled_lookup_heldout_one_step_accuracy
$roll = [double]$probe.metrics.minimum_factorized_autoregressive_accuracy
if ($one -lt 0.95) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:minimum_factorized_heldout_one_step_accuracy actual=$one threshold=0.95" }
if ($lookup -gt 0.25) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:maximum_pooled_lookup_heldout_one_step_accuracy actual=$lookup threshold=0.25" }
if ($roll -lt 0.95) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:minimum_factorized_autoregressive_accuracy actual=$roll threshold=0.95" }

$metrics = [ordered]@{}
foreach ($property in $probe.metrics.PSObject.Properties) {
  $metrics[$property.Name] = [double]$property.Value
}
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-LM-RULE-CONDITIONED-SEQUENCE-R2'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen rule-conditioned held-out sequence gates passed.'
