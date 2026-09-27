$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up189c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up189c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP189C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP189C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up189c-correction-policy-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP189C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP189C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up189c-correction-policy-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP189C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up189c-correction-policy-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP189C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up189c-correction-policy-transfer.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.max_actions-ne 2 -or [int]$R.stage_two_delay_intervals-ne 1){throw 'UP189C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.warning_threshold_changed -or $R.adaptive_delay_used -or $R.policy_specific_rule_used -or $R.future_policy_schedule_used_by_warning){throw 'UP189C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) baseline_losses=$($M.baseline_losses) treated_losses=$($M.treated_losses) prevented=$($M.prevented_losses) accelerated=$($M.accelerated_losses) stage1=$($M.stage_one_actions) stage2=$($M.stage_two_actions) efficiency=$($M.prevented_per_action) failure_delay=$($M.mean_loss_step_change_among_failures)"}
 Write-Host 'WINGLESS_UP237_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
