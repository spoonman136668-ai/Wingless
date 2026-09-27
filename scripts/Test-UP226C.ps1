$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up226c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up226c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP226C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP226C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up226c-critical-advanced-action5 -count=1;if($LASTEXITCODE-ne 0){throw 'UP226C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP226C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up226c-critical-advanced-action5)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP226C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up226c-critical-advanced-action5)|Out-String).Trim();if($P1-cne $P2){throw 'UP226C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up226c-critical-advanced-action5.v1' -or [int]$R.arms.Count-ne 64 -or [int]$R.max_actions-ne 8 -or [int]$R.monitoring_cadence-ne 2){throw 'UP226C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.threshold_fitting_used -or $R.schedule_compression_used -or $R.action_budget_increased -or $R.future_policy_input_used -or $R.live_activation){throw 'UP226C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "comparator_losses=$($R.comparator_losses) primary_losses=$($R.primary_losses) advanced=$($R.advanced_action5_arms) rescued=$($R.rescued_losses) unnecessary=$($R.unnecessary_advances) necessary=$($R.necessary_advance_attempts) ineffective=$($R.ineffective_advance_attempts)"
 foreach($A in @($R.arms|Where-Object{$_.action5_advanced -or [int]$_.comparator_loss_step -le 80})){Write-Host "cohort=$($A.cohort) hand=$($A.initial_hand) endangered=$($A.endangered_key) base_loss=$($A.comparator_loss_step) primary_loss=$($A.primary_loss_step) advance=$($A.action5_advance_writes)"}
 Write-Host 'WINGLESS_UP319_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
