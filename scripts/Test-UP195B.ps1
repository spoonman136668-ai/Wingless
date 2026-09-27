$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up195b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up195b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP195B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP195B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP195B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP195B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up195b-adjacent-swap-history -count=1;if($LASTEXITCODE-ne 0){throw 'UP195B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP195B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up195b-adjacent-swap-history)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP195B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up195b-adjacent-swap-history)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP195B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP195B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up195b-adjacent-swap-history.v1'){throw 'UP195B_SCHEMA_MISMATCH'}
 if([int]$R.compositions.Count-ne 4 -or [int]$R.metrics.Count-ne 33){throw 'UP195B_DESIGN'}
 if($R.new_correction_fit -or $R.adaptive_swap_search_used -or $R.maintenance_triggered){throw 'UP195B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.group_summaries){Write-Host "group=$($S.group) value=$($S.value) observations=$($S.observations) mean_delta=$($S.mean_absolute_factor_difference) max_delta=$($S.max_absolute_factor_difference)"}
 Write-Host "same_count=$($R.same_count_swap_observations) same_count_nonzero=$($R.same_count_nonzero_factor_observations)"
 Write-Host 'WINGLESS_UP195_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
