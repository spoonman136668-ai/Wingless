$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up192b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up192b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP192B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP192B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP192B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP192B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up192b-native-margin-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP192B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP192B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up192b-native-margin-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP192B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up192b-native-margin-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP192B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP192B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up192b-native-margin-geometry.v1'){throw 'UP192B_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 27 -or [int]$R.profiles.Count-ne 9){throw 'UP192B_DESIGN'}
 if($R.fitted_correction_used -or $R.heldout_tuning_used -or $R.adaptive_feature_selection_used -or $R.adaptive_profile_search_used -or $R.maintenance_triggered){throw 'UP192B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "corr_correct_count=$($R.correlations.native_correct_count) corr_mean_abs_margin=$($R.correlations.mean_absolute_margin) corr_near_zero=$($R.correlations.near_zero_margin_count) corr_min_abs_margin=$($R.correlations.min_absolute_margin)"
 Write-Host 'WINGLESS_UP192_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
