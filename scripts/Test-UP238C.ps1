$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up238c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up238c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP238C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP238C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up238c-trajectory-signature-holdout -count=1;if($LASTEXITCODE-ne 0){throw 'UP238C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP238C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up238c-trajectory-signature-holdout)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP238C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up238c-trajectory-signature-holdout)|Out-String).Trim();if($P1-cne $P2){throw 'UP238C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up238c-trajectory-signature-holdout.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.holdout_folds-ne 4 -or [int]$R.summaries.Count-ne 8){throw 'UP238C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or -not $R.signature_lookup_used -or $R.adaptive_feature_selection_used -or -not $R.holdout_schedule_excluded_from_training -or $R.new_native_field_used -or $R.live_activation){throw 'UP238C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "holdout=$($S.holdout_schedule) mode=$($S.signature_mode) train_fail_sig=$($S.train_failure_signatures) train_survive_sig=$($S.train_survivor_signatures) shared=$($S.train_shared_signatures) test_event=$($S.test_event_arms) failures=$($S.test_failures) survivors=$($S.test_survivors) tp=$($S.true_positive) fn=$($S.false_negative) fp=$($S.false_positive) tn=$($S.true_negative) silent=$($S.holdout_silent_failures) sensitivity=$($S.sensitivity) specificity=$($S.specificity)"}
 Write-Host 'WINGLESS_UP341_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
