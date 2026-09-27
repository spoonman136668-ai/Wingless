$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up191c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up191c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP191C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP191C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up191c-critical-refresh-budget-sweep -count=1;if($LASTEXITCODE-ne 0){throw 'UP191C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP191C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up191c-critical-refresh-budget-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP191C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up191c-critical-refresh-budget-sweep)|Out-String).Trim();if($P1-cne $P2){throw 'UP191C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up191c-critical-refresh-budget-sweep.v1' -or [int]$R.metrics.Count-ne 64 -or [int]$R.action_caps.Count-ne 8){throw 'UP191C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_cap_used -or $R.warning_threshold_changed -or $R.policy_specific_rule_used -or $R.new_action_type_used){throw 'UP191C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) cap=$($M.action_cap) baseline=$($M.baseline_losses) treated=$($M.treated_losses) prevented=$($M.prevented_losses) accelerated=$($M.accelerated_losses) actions=$($M.actions_taken) efficiency=$($M.prevented_per_action) failure_delay=$($M.mean_loss_step_change_among_failures)"}
 Write-Host 'WINGLESS_UP242_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
