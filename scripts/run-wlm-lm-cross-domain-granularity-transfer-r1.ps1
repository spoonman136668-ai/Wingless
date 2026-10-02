$ErrorActionPreference='Stop'
$env:GOFLAGS='-buildvcs=false'
$env:GOMAXPROCS='1'
& go test ./... -count=1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_FULL_REGRESSION_FAILED'}
$generated=@('.ice/architecture-map.json','.ice/dependency-map.json','.ice/integration-seams.json','.ice/manifest.json','.ice/symbol-map.json','.ice/test-map.json')
foreach($path in $generated){
  & git checkout -- $path 2>$null
  if($LASTEXITCODE-ne 0){if(Test-Path -LiteralPath $path){Remove-Item -LiteralPath $path -Force -Recurse}}
}
$raw=& go run ./cmd/wlm-lm-cross-domain-granularity-transfer-r1
if($LASTEXITCODE-ne 0){throw 'WINGLESS_PROBE_FAILED'}
$probe=($raw-join[Environment]::NewLine)|ConvertFrom-Json
if([string]$probe.schema-cne'wingless.research-scientific-result.v1'){throw 'WINGLESS_RESULT_SCHEMA_INVALID'}
if([string]$probe.experiment-cne'WLM-LM-CROSS-DOMAIN-GRANULARITY-TRANSFER-R1'){throw 'WINGLESS_RESULT_EXPERIMENT_INVALID'}
if($null-eq$probe.metrics){throw 'WINGLESS_RESULT_METRICS_MISSING'}
$metrics=[ordered]@{}
foreach($property in $probe.metrics.PSObject.Properties){
  $value=[double]$property.Value
  if([double]::IsNaN($value)-or[double]::IsInfinity($value)){throw "WINGLESS_METRIC_NONFINITE:$($property.Name)"}
  $metrics[$property.Name]=$value
}
$result=[ordered]@{schema='wingless.research-scientific-result.v1';experiment='WLM-LM-CROSS-DOMAIN-GRANULARITY-TRANSFER-R1';metrics=$metrics}
Write-Output($result|ConvertTo-Json -Depth 8 -Compress)
Write-Host '=== SCIENTIFIC DIAGNOSIS'
Write-Host 'Evidence-preserving WBG-4 cross-domain transfer run completed.'
