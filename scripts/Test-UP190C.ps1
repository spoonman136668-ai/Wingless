$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up190c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up190c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP190C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP190C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up190c-two-action-oracle-ceiling -count=1;if($LASTEXITCODE-ne 0){throw 'UP190C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP190C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up190c-two-action-oracle-ceiling)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP190C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up190c-two-action-oracle-ceiling)|Out-String).Trim();if($P1-cne $P2){throw 'UP190C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up190c-two-action-oracle-ceiling.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.max_actions-ne 2){throw 'UP190C_DESIGN'}
 if(-not $R.offline_oracle_only -or $R.live_activation -or $R.action_budget_changed -or $R.arbitrary_write_timing_used){throw 'UP190C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) baseline=$($M.baseline_losses) frozen_prevented=$($M.frozen_rule_prevented) oracle_preventable=$($M.oracle_preventable) oracle_missed=$($M.oracle_preventable_frozen_missed) not_preventable=$($M.not_preventable_with_two_actions) mean_best_extension=$($M.mean_best_oracle_loss_extension) max_best_extension=$($M.max_best_oracle_loss_extension)"}
 Write-Host 'WINGLESS_UP240_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
