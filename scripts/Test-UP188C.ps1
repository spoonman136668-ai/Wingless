$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up188c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up188c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP188C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP188C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up188c-risk-policy-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP188C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP188C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up188c-risk-policy-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP188C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up188c-risk-policy-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP188C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up188c-risk-policy-transfer.v1' -or [int]$R.metrics.Count-ne 24 -or [int]$R.policies.Count-ne 4 -or [int]$R.risk_buckets.Count-ne 3){throw 'UP188C_DESIGN'}
 if($R.intervention_used -or $R.adaptive_policy_selection_used -or $R.adaptive_bucketing_used -or $R.warning_threshold_changed -or $R.future_policy_schedule_used_by_warning -or $R.live_activation){throw 'UP188C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) bucket=$($M.risk_bucket) checkpoints=$($M.checkpoints) one_rate=$($M.one_cadence_failure_rate) two_rate=$($M.two_cadence_failure_rate) losses=$($M.trajectory_losses_by_64)"}
 Write-Host 'WINGLESS_UP234_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
