$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up145b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up145b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP145B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP145B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP145B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP145B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up145b-tail-switch-frequency -count=1;if($LASTEXITCODE-ne 0){throw 'UP145B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP145B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up145b-tail-switch-frequency)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP145B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up145b-tail-switch-frequency)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP145B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP145B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up145b-tail-switch-frequency.v1'){throw 'UP145B_SCHEMA_MISMATCH'}
 if([int]$R.summaries.Count-ne 3){throw 'UP145B_SUMMARY_COUNT'}
 $Expected=@(3,11,19)
 for($i=0;$i-lt 3;$i++){if([int]$R.summaries[$i].switches-ne $Expected[$i]){throw 'UP145B_SWITCH_COUNT'}}
 if([int]$R.distinct_tail_sets-ne 4 -or [int]$R.uses_per_tail_set-ne 5){throw 'UP145B_COVERAGE'}
 if($R.adaptive_ordering_used -or $R.extra_updates_used){throw 'UP145B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "arm=$($S.arm) switches=$($S.switches) tail=$($S.mean_tail_damage) net=$($S.mean_net_epoch_change) old=$($S.final_mean_old_retention) new=$($S.final_mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP145_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
