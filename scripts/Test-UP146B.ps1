$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up146b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up146b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP146B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP146B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP146B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP146B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up146b-tail-habituation -count=1;if($LASTEXITCODE-ne 0){throw 'UP146B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP146B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up146b-tail-habituation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up146b-tail-habituation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP146B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up146b-tail-habituation.v1'){throw 'UP146B_SCHEMA_MISMATCH'}
 if([int]$R.epoch_metrics.Count-ne 20 -or [int]$R.position_summaries.Count-ne 5){throw 'UP146B_METRIC_COUNT'}
 if([int]$R.distinct_tail_sets-ne 4 -or [int]$R.uses_per_tail_set-ne 5 -or [int]$R.switches-ne 3){throw 'UP146B_SCHEDULE'}
 for($i=0;$i-lt 5;$i++){if([int]$R.position_summaries[$i].repetition-ne ($i+1) -or [int]$R.position_summaries[$i].samples-ne 4){throw 'UP146B_POSITION_SUMMARY'}}
 if($R.adaptive_ordering_used -or $R.extra_updates_used){throw 'UP146B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.position_summaries){Write-Host "rep=$($S.repetition) prefix=$($S.mean_prefix_damage) anchor=$($S.mean_anchor_recovery) tail=$($S.mean_tail_damage) net=$($S.mean_net_epoch_change)"}
 Write-Host 'WINGLESS_UP146_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
