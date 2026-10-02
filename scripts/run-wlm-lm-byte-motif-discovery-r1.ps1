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

$raw = & go run ./cmd/wlm-lm-byte-motif-discovery-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$exact = [ordered]@{
  valid_seed_count = 4
  completed_substrate_runs = 8
  selected_motif_count = 8
  substrate_selected_motif_mismatch_count = 0
  substrate_prediction_mismatch_count = 0
  substrate_event_accounting_mismatch_count = 0
  duplicate_shuffled_control_key_count = 0
  invalid_candidate_rows = 0
  counter_overflow_rows = 0
}
foreach ($name in $exact.Keys) {
  $property = $probe.metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  if ($value -ne [double]$exact[$name]) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:$name actual=$value expected=$($exact[$name])" }
}

$thresholds = @(
  @('minimum_selected_true_motif_recall','ge',0.875),
  @('maximum_false_selected_motif_count','le',1),
  @('minimum_motif_enabled_successor_accuracy','ge',0.95),
  @('maximum_byte_only_ablation_successor_accuracy','le',0.20),
  @('maximum_shuffled_motif_successor_accuracy','le',0.20),
  @('minimum_causal_accuracy_gain_over_ablation','ge',0.70),
  @('minimum_effective_event_reduction_fraction','ge',0.30)
)
foreach ($row in $thresholds) {
  $name=[string]$row[0]; $op=[string]$row[1]; $threshold=[double]$row[2]
  $value=[double]$probe.metrics.PSObject.Properties[$name].Value
  if ($op -eq 'ge' -and $value -lt $threshold) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:$name actual=$value threshold=$threshold" }
  if ($op -eq 'le' -and $value -gt $threshold) { throw "WINGLESS_SCIENTIFIC_GATE_FAILED:$name actual=$value threshold=$threshold" }
}

$metrics = [ordered]@{}
foreach ($property in $probe.metrics.PSObject.Properties) {
  $metrics[$property.Name] = [double]$property.Value
}
$result = [ordered]@{
  schema = 'wingless.research-scientific-result.v1'
  experiment = 'WLM-LM-BYTE-MOTIF-DISCOVERY-R1'
  metrics = $metrics
}
Write-Output ($result | ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen WBG-2 byte motif discovery gates passed.'
