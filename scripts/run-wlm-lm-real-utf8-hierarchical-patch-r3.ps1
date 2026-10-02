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

$raw = & go run ./cmd/wlm-lm-real-utf8-hierarchical-patch-r3
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-LM-REAL-UTF8-HIERARCHICAL-PATCH-R3') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }
if ($null -eq $probe.metrics) { throw 'WINGLESS_RESULT_METRICS_MISSING' }

$required = @(
  'training_file_identity_mismatch_count','evaluation_file_identity_mismatch_count',
  'train_eval_blob_overlap_count','selected_level1_count','selected_level2_count',
  'minimum_eval_level1_coverage_fraction','minimum_eval_level2_coverage_fraction',
  'minimum_level2_covered_accuracy_gain','minimum_effective_event_reduction_fraction',
  'maximum_hierarchical_minus_level1_bits_per_byte','maximum_hierarchical_minus_baseline_bits_per_byte',
  'maximum_total_learned_structure_count','capacity_growth_event_count','tokenizer_use_count',
  'invalid_byte_rows','counter_overflow_rows'
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
  experiment='WLM-LM-REAL-UTF8-HIERARCHICAL-PATCH-R3'
  metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving real UTF8 hierarchical patch R3 completed; classify only from the original frozen thresholds.'
