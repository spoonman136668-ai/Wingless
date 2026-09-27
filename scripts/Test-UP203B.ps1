$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up203b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up203b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP203B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP203B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP203B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP203B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up203b-order-signal-persistence -count=1;if($LASTEXITCODE-ne 0){throw 'UP203B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP203B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up203b-order-signal-persistence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP203B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up203b-order-signal-persistence)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP203B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP203B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up203b-order-signal-persistence.v1'){throw 'UP203B_SCHEMA_MISMATCH'}
 if([int]$R.permutation_count_per_composition-ne 24 -or [int]$R.phase_metrics.Count-ne 6 -or [int]$R.summaries.Count-ne 2){throw 'UP203B_DESIGN'}
 if($R.heldout_fitting_used -or $R.feature_search_used -or $R.maintenance_triggered){throw 'UP203B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) template_spread=$($S.template_factor_spread) mean_corr=$($S.mean_residual_correlation) min_corr=$($S.min_residual_correlation) max_corr=$($S.max_residual_correlation) heldout_spread=$($S.mean_heldout_factor_spread)"}
 Write-Host 'WINGLESS_UP217_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
