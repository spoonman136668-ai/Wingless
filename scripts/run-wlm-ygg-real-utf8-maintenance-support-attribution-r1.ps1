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

$raw = & go run ./cmd/wlm-ygg-real-utf8-maintenance-support-attribution-r1
if ($LASTEXITCODE -ne 0) { throw 'WINGLESS_PROBE_FAILED' }
$probe = ($raw -join [Environment]::NewLine) | ConvertFrom-Json
if ([string]$probe.schema -cne 'wingless.research-scientific-result.v1') { throw 'WINGLESS_RESULT_SCHEMA_INVALID' }
if ([string]$probe.experiment -cne 'WLM-YGG-REAL-UTF8-MAINTENANCE-SUPPORT-ATTRIBUTION-R1') { throw 'WINGLESS_RESULT_EXPERIMENT_INVALID' }

$required=@(
 'valid_evaluation_file_count','canonical_selected_motif_count','strictest_feasible_tier_index',
 'minimum_selected_tier_eligible_count_per_file','training_file_identity_mismatch_count',
 'evaluation_file_identity_mismatch_count','tokenizer_use_count','capacity_growth_event_count',
 'invalid_attribution_rows'
)
for($file=0;$file -lt 2;$file++){for($tier=0;$tier -lt 5;$tier++){$required+=("file{0}_tier{1}_eligible_count" -f $file,$tier)}}
for($tier=0;$tier -lt 5;$tier++){$required+=("minimum_tier{0}_eligible_count_per_file" -f $tier)}

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
 experiment='WLM-YGG-REAL-UTF8-MAINTENANCE-SUPPORT-ATTRIBUTION-R1'
 metrics=$metrics
}
Write-Output ($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Frozen real UTF8 maintenance support attribution completed.'
