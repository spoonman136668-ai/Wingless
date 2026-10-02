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

$raw = & go run ./cmd/wlm-lm-byte-motif-attribution-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json

$exact = [ordered]@{
  valid_seed_count = 4
  minimum_unique_hidden_motif_ids_per_seed = 2
  maximum_unique_hidden_motif_ids_per_seed = 2
  union_hidden_motif_id_count_across_seeds = 4
  minimum_present_motif_occurrence_count = 2048
  maximum_present_motif_occurrence_count = 2048
  minimum_absent_motif_count_per_seed = 6
  maximum_absent_motif_count_per_seed = 6
  minimum_present_true_motif_selector_recall = 1
  maximum_present_true_motif_selector_recall = 1
  minimum_selected_true_motif_count_per_seed = 2
  maximum_selected_true_motif_count_per_seed = 2
  minimum_selected_false_motif_count_per_seed = 6
  maximum_selected_false_motif_count_per_seed = 6
  substrate_selected_motif_mismatch_count = 0
  counter_overflow_rows = 0
}
foreach ($name in $exact.Keys) {
  $property = $probe.metrics.PSObject.Properties[$name]
  if ($null -eq $property) { throw "WINGLESS_METRIC_MISSING:$name" }
  $value = [double]$property.Value
  if ([double]::IsNaN($value) -or [double]::IsInfinity($value)) { throw "WINGLESS_METRIC_NONFINITE:$name" }
  if ($value -ne [double]$exact[$name]) { throw "WINGLESS_ATTRIBUTION_GATE_FAILED:$name actual=$value expected=$($exact[$name])" }
}

$metrics=[ordered]@{}
foreach($property in $probe.metrics.PSObject.Properties){$metrics[$property.Name]=[double]$property.Value}
$result=[ordered]@{
  schema='wingless.research-scientific-result.v1'
  experiment='WLM-LM-BYTE-MOTIF-ATTRIBUTION-R1'
  metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen WBG-2 corpus-support attribution gates passed.'
