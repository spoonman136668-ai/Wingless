$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up227c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up227c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP227C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP227C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up227c-action5-retiming-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP227C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP227C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up227c-action5-retiming-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP227C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up227c-action5-retiming-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP227C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up227c-action5-retiming-transfer.v1' -or [int]$R.metrics.Count-ne 16 -or [int]$R.max_actions-ne 8 -or [int]$R.monitoring_cadence-ne 2){throw 'UP227C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.threshold_fitting_used -or $R.schedule_compression_used -or $R.action_budget_increased -or $R.future_policy_input_used -or $R.live_activation){throw 'UP227C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) horizon=$($M.horizon) base_losses=$($M.comparator_losses) primary_losses=$($M.primary_losses) advanced=$($M.advanced_action5_arms) rescued=$($M.rescued_losses) unnecessary=$($M.unnecessary_advances) necessary=$($M.necessary_advance_attempts) ineffective=$($M.ineffective_advance_attempts) advance_writes=$($M.total_advance_writes)"}
 Write-Host 'WINGLESS_UP321_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
