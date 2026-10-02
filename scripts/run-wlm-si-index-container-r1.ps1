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

$raw = & go run ./cmd/wlm-si-index-container-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$expected = [ordered]@{
  paired_case_count = 512
  completed_substrate_case_runs = 1024
  total_emitted_rows = 65536
  paired_observable_mismatch_cases = 0
  row_count_mismatch_cases = 0
  state_cardinality_mismatch_cases = 0
  invalid_operation_schedules = 0
  budget_overrun_rows = 0
  law_violation_cases = 0
  max_conservation_error = 0
  max_residual_law_error = 0
}

$metrics = [ordered]@{}
foreach ($name in $expected.Keys) {
  $property = $probe.Metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  if ($value -ne [double]$expected[$name]) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:$name actual=$value expected=$($expected[$name])" }
  $metrics[$name] = $value
}

$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-SI-INDEX-CONTAINER-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen exact-invariance gates passed for the paired index substrate matrix.'
