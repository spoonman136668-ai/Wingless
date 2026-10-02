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

$raw = & go run ./cmd/wlm-si-priority-container-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$names = @(
  'paired_case_count',
  'completed_substrate_case_runs',
  'total_emitted_rows',
  'paired_observable_mismatch_cases',
  'row_count_mismatch_cases',
  'state_cardinality_mismatch_cases',
  'invalid_operation_schedules',
  'budget_overrun_rows',
  'law_violation_cases',
  'max_conservation_error',
  'max_residual_law_error'
)

$metrics = [ordered]@{}
foreach ($name in $names) {
  $property = $probe.Metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  $metrics[$name] = $value
}

$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-SI-PRIORITY-CONTAINER-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen priority-container matrix completed; classification is controller-owned from preregistered thresholds.'
