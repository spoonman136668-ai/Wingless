$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up200c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up200c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP200C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP200C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up200c-hostile-trigger-oracle-gap -count=1;if($LASTEXITCODE-ne 0){throw 'UP200C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP200C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up200c-hostile-trigger-oracle-gap)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP200C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up200c-hostile-trigger-oracle-gap)|Out-String).Trim();if($P1-cne $P2){throw 'UP200C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up200c-hostile-trigger-oracle-gap.v1' -or [int]$R.points.Count-ne 64 -or [int]$R.horizon-ne 56 -or [int]$R.action_cap-ne 8){throw 'UP200C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or -not $R.oracle_future_outcome_access -or $R.oracle_feeds_candidate_rule -or $R.adaptive_selection_used){throw 'UP200C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "arms=$($R.arms) baseline_losses=$($R.baseline_losses) internal_prevented=$($R.internal_prevented) internal_actions=$($R.internal_actions) oracle_safe=$($R.oracle_safe_arms) oracle_actions=$($R.oracle_actions) internal_onset=$($R.mean_internal_onset) oracle_latest=$($R.mean_oracle_latest_safe_onset) headroom=$($R.mean_onset_headroom) action_savings=$($R.oracle_action_savings) internal_equals_latest=$($R.internal_equals_latest_safe)"
 Write-Host 'WINGLESS_UP265_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
