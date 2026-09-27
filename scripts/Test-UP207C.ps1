$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up207c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up207c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP207C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP207C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up207c-warning-pullforward-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP207C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP207C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up207c-warning-pullforward-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP207C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up207c-warning-pullforward-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP207C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up207c-warning-pullforward-transfer.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.phase_shift_writes-ne 4 -or [int]$R.action_cap-ne 8){throw 'UP207C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.extra_budget_used -or $R.policy_specific_rule_used -or $R.future_schedule_used_by_warning -or $R.adaptive_selection_used){throw 'UP207C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) mode=$($M.mode) baseline=$($M.baseline_losses) treated=$($M.treated_losses) prevented=$($M.prevented_losses) actions=$($M.actions_taken) pullforwards=$($M.pull_forward_actions) accelerated=$($M.accelerated_losses) failure_extension=$($M.mean_loss_step_extension_among_failures)"}
 Write-Host 'WINGLESS_UP280_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
