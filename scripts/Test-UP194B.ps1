$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up194b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up194b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP194B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP194B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP194B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP194B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up194b-order-history-effect -count=1;if($LASTEXITCODE-ne 0){throw 'UP194B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP194B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up194b-order-history-effect)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP194B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up194b-order-history-effect)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP194B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP194B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up194b-order-history-effect.v1'){throw 'UP194B_SCHEMA_MISMATCH'}
 if([int]$R.pairs.Count-ne 4 -or [int]$R.metrics.Count-ne 12){throw 'UP194B_DESIGN'}
 if($R.new_correction_fit -or $R.correction_applied -or $R.adaptive_pair_search_used -or $R.maintenance_triggered){throw 'UP194B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "pair=$($M.pair) phase=$($M.phase) count_delta=$($M.absolute_native_count_difference) same_count=$($M.same_native_correct_count) factor_delta=$($M.absolute_factor_difference)"}
 Write-Host "mean_factor_delta=$($R.summary.mean_absolute_factor_difference) max_factor_delta=$($R.summary.max_absolute_factor_difference) same_count_pairs=$($R.summary.same_count_pairs) same_count_nonzero=$($R.summary.same_count_nonzero_factor_pairs)"
 Write-Host 'WINGLESS_UP194_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
