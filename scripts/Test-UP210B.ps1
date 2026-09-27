$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up210b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up210b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP210B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP210B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up210b-family-holdout-drift -count=1;if($LASTEXITCODE-ne 0){throw 'UP210B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP210B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up210b-family-holdout-drift)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP210B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up210b-family-holdout-drift)|Out-String).Trim();if($P1-cne $P2){throw 'UP210B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up210b-family-holdout-drift.v1' -or [int]$R.metrics.Count-ne 48 -or [int]$R.summaries.Count-ne 4){throw 'UP210B_DESIGN'}
 foreach($S in $R.summaries){if([int]$S.training_points-ne 42 -or [int]$S.evaluation_points-ne 12){throw 'UP210B_SPLIT'}}
 if($R.heldout_fitting_used -or $R.adaptive_feature_selection_used -or $R.family_identity_input_used -or $R.absolute_state_input_used -or $R.phase_input_used -or $R.parity_input_used -or $R.future_state_input_used){throw 'UP210B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "heldout=$($S.heldout_family) zero_mae=$($S.zero_drift_mean_absolute_error) native_mae=$($S.native_drift_mean_absolute_error) native_vs_zero=$($S.native_better_than_zero)/$($S.evaluation_points) direction=$($S.correct_drift_direction)/$($S.nonzero_actual_drift_points) native_max=$($S.native_drift_max_absolute_error)"}
 Write-Host 'WINGLESS_UP245_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
