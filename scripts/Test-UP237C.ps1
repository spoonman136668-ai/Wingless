$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up237c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up237c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP237C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP237C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up237c-critical-trajectory-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP237C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP237C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up237c-critical-trajectory-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP237C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up237c-critical-trajectory-diagnostic)|Out-String).Trim();if($P1-cne $P2){throw 'UP237C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up237c-critical-trajectory-diagnostic.v1' -or [int]$R.arms_total-ne 512){throw 'UP237C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP237C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "failures=$($R.failures) survivors=$($R.survivors) event_failures=$($R.event_failures) event_survivors=$($R.event_survivors) silent=$($R.silent_failures) current_fail_sig=$($R.current_failure_signatures) current_survive_sig=$($R.current_survivor_signatures) current_shared=$($R.current_shared_signatures) traj_fail_sig=$($R.trajectory_failure_signatures) traj_survive_sig=$($R.trajectory_survivor_signatures) traj_shared=$($R.trajectory_shared_signatures)"
 Write-Host 'WINGLESS_UP338_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
