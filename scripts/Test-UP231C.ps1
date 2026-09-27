$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up231c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up231c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP231C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP231C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up231c-first-critical-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP231C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP231C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up231c-first-critical-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP231C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up231c-first-critical-geometry)|Out-String).Trim();if($P1-cne $P2){throw 'UP231C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up231c-first-critical-geometry.v1' -or [int]$R.arms_total-ne 512 -or [int]$R.arms.Count-ne 512 -or [int]$R.monitoring_cadence-ne 2){throw 'UP231C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.adaptive_feature_selection_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP231C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "failures=$($R.failures) survivors=$($R.survivors) event_failures=$($R.event_bearing_failures) event_survivors=$($R.event_bearing_survivors) silent_failures=$($R.silent_failures) noevent_survivors=$($R.no_event_survivors) failure_signatures=$($R.distinct_failure_signatures) survivor_signatures=$($R.distinct_survivor_signatures) shared_signatures=$($R.shared_signatures)"
 foreach($A in @($R.arms|Where-Object{$_.event_present})){Write-Host "schedule=$($A.schedule) advance=$($A.phase_advance_writes) cohort=$($A.cohort) hand=$($A.initial_hand) outcome=$($A.outcome) event=$($A.event_step) loss=$($A.loss_step) signature=$($A.signature)"}
 Write-Host 'WINGLESS_UP329_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
