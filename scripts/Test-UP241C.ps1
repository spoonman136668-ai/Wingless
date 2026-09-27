$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up241c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up241c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP241C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP241C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up241c-event-support-coverage -count=1;if($LASTEXITCODE-ne 0){throw 'UP241C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP241C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up241c-event-support-coverage)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP241C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up241c-event-support-coverage)|Out-String).Trim();if($P1-cne $P2){throw 'UP241C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up241c-event-support-coverage.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.schedules-ne 4 -or [int]$R.advance_conditions.Count-ne 2 -or [int]$R.summaries.Count-ne 8){throw 'UP241C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.classifier_training_used -or $R.adaptive_feature_selection_used -or $R.new_native_field_used -or $R.live_activation){throw 'UP241C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "schedule=$($S.schedule) advance=$($S.phase_advance_writes) arms=$($S.arms) failures=$($S.failures) survivors=$($S.survivors) event_arms=$($S.event_arms) event_failures=$($S.event_failures) event_survivors=$($S.event_survivors) silent_failures=$($S.silent_failures)"}
 Write-Host 'WINGLESS_UP350_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
