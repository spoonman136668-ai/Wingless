$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up208b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up208b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP208B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP208B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up208b-family-holdout-native-offset -count=1;if($LASTEXITCODE-ne 0){throw 'UP208B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP208B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up208b-family-holdout-native-offset)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP208B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up208b-family-holdout-native-offset)|Out-String).Trim();if($P1-cne $P2){throw 'UP208B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up208b-family-holdout-native-offset.v1' -or [int]$R.metrics.Count-ne 222 -or [int]$R.summaries.Count-ne 4){throw 'UP208B_DESIGN'}
 foreach($S in $R.summaries){if([int]$S.training_points-ne 45){throw 'UP208B_TRAINING_SPLIT'}}
 if($R.heldout_fitting_used -or $R.adaptive_feature_selection_used -or $R.family_identity_input_used -or $R.phase_input_used -or $R.parity_input_used -or $R.heldout_anchor_measurement_used_by_native -or $R.maintenance_triggered){throw 'UP208B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "heldout=$($S.heldout_family) static_mae=$($S.static_mean_absolute_error) native_mae=$($S.native_offset_mean_absolute_error) oracle_mae=$($S.oracle_anchor_mean_absolute_error) native_vs_static=$($S.native_better_than_static)/$($S.evaluation_points) native_max=$($S.native_offset_max_absolute_error)"}
 Write-Host 'WINGLESS_UP236_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
