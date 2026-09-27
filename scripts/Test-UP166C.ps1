$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up166c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up166c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP166C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP166C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP166C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP166C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up166c-warning-specificity -count=1;if($LASTEXITCODE-ne 0){throw 'UP166C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP166C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up166c-warning-specificity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP166C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up166c-warning-specificity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP166C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP166C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up166c-warning-specificity.v1'){throw 'UP166C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 10 -or -not $R.shadow_only -or $R.corrective_action_used -or $R.warning_model_changed_across_policies -or $R.future_real_policy_schedule_used_by_warning -or -not $R.safe_control_included){throw 'UP166C_DESIGN'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) losses=$($M.eventual_losses) warnings=$($M.warning_intervals) expired=$($M.expired_warning_intervals) precision=$($M.warning_precision) recall=$($M.warning_recall) lead=$($M.mean_first_warning_lead_writes)"}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
