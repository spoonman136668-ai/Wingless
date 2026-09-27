$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up191b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up191b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP191B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP191B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP191B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP191B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up191b-native-response-curve -count=1;if($LASTEXITCODE-ne 0){throw 'UP191B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP191B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up191b-native-response-curve)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP191B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up191b-native-response-curve)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP191B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP191B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up191b-native-response-curve.v1'){throw 'UP191B_SCHEMA_MISMATCH'}
 if([int]$R.profiles.Count-ne 9 -or [int]$R.points.Count-ne 27){throw 'UP191B_DESIGN'}
 if($R.fitted_correction_used -or $R.heldout_tuning_used -or $R.adaptive_profile_search_used -or $R.phase_input_used -or $R.parity_input_used -or $R.maintenance_triggered){throw 'UP191B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "pearson_native_count_required_factor=$($R.pearson_native_count_required_factor) max_within_count_factor_spread=$($R.max_within_count_factor_spread)"
 foreach($G in $R.count_groups){Write-Host "native_correct=$($G.native_correct_count) n=$($G.observations) mean_factor=$($G.mean_required_mass_factor) spread=$($G.factor_spread)"}
 Write-Host 'WINGLESS_UP191_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
