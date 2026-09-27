$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up190b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up190b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP190B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP190B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP190B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP190B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up190b-endogenous-context-shift -count=1;if($LASTEXITCODE-ne 0){throw 'UP190B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP190B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up190b-endogenous-context-shift)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP190B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up190b-endogenous-context-shift)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP190B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP190B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up190b-endogenous-context-shift.v1'){throw 'UP190B_SCHEMA_MISMATCH'}
 if([int]$R.profiles.Count-ne 3 -or [int]$R.metrics.Count-ne 18 -or [int]$R.summaries.Count-ne 2){throw 'UP190B_DESIGN'}
 if($R.phase_input_used -or $R.parity_input_used -or $R.heldout_fitting_used -or $R.adaptive_shift_search_used -or $R.maintenance_triggered){throw 'UP190B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "model=$($S.model) mean_abs_mass_error=$($S.mean_absolute_mass_ratio_error) max_abs_mass_error=$($S.max_absolute_mass_ratio_error)"}
 foreach($M in $R.metrics|Where-Object{$_.model-eq 'pooled_plus_native_mass_correction'}){Write-Host "profile=$($M.profile) phase=$($M.phase) native_correct=$($M.native_correct_count) outside_known=$($M.outside_known_native_counts) abs_mass_error=$($M.absolute_mass_ratio_error)"}
 Write-Host 'WINGLESS_UP190_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
