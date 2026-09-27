$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up209b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up209b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP209B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP209B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up209b-family-holdout-native-delta -count=1;if($LASTEXITCODE-ne 0){throw 'UP209B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP209B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up209b-family-holdout-native-delta)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP209B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up209b-family-holdout-native-delta)|Out-String).Trim();if($P1-cne $P2){throw 'UP209B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up209b-family-holdout-native-delta.v1' -or [int]$R.metrics.Count-ne 222 -or [int]$R.summaries.Count-ne 4){throw 'UP209B_DESIGN'}
 foreach($S in $R.summaries){if([int]$S.training_points-ne 42){throw 'UP209B_TRAINING_SPLIT'}}
 if($R.heldout_fitting_used -or $R.adaptive_feature_selection_used -or $R.family_identity_input_used -or $R.phase_input_used -or $R.parity_input_used -or $R.future_state_input_used -or $R.heldout_anchor_measurement_used_by_native -or $R.maintenance_triggered){throw 'UP209B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "heldout=$($S.heldout_family) static_mae=$($S.static_mean_absolute_error) native_delta_mae=$($S.native_delta_mean_absolute_error) oracle_mae=$($S.oracle_anchor_mean_absolute_error) native_vs_static=$($S.native_better_than_static)/$($S.evaluation_points) native_max=$($S.native_delta_max_absolute_error)"}
 Write-Host 'WINGLESS_UP239_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
