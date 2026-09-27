$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up167c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up167c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP167C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP167C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP167C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP167C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up167c-warning-confirmation -count=1;if($LASTEXITCODE-ne 0){throw 'UP167C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP167C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up167c-warning-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP167C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up167c-warning-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP167C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP167C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up167c-warning-confirmation.v1'){throw 'UP167C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 24 -or [int]$R.confirmations.Count-ne 3 -or -not $R.shadow_only){throw 'UP167C_DESIGN'}
 if($R.corrective_action_used -or $R.adaptive_confirmation_used -or $R.warning_threshold_changed -or $R.future_real_policy_schedule_used_by_warning){throw 'UP167C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) confirm=$($M.confirmation) precision=$($M.trigger_precision) recall=$($M.trigger_recall) eligible=$($M.eligible_intervals) expired=$($M.expired_eligible_intervals) lead=$($M.mean_first_trigger_lead_writes)"}
 Write-Host 'WINGLESS_UP167_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
