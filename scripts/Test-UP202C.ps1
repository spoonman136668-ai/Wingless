$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up202c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up202c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP202C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP202C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up202c-policy-horizon-selectivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP202C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP202C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up202c-policy-horizon-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP202C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up202c-policy-horizon-selectivity)|Out-String).Trim();if($P1-cne $P2){throw 'UP202C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up202c-policy-horizon-selectivity.v1' -or [int]$R.metrics.Count-ne 40 -or [int]$R.policies.Count-ne 4 -or [int]$R.horizons.Count-ne 10 -or [int]$R.grace_writes-ne 8){throw 'UP202C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.policy_specific_rule_used -or $R.adaptive_selection_used){throw 'UP202C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) horizon=$($M.horizon) baseline_failures=$($M.baseline_failures) survivors=$($M.baseline_survivors) treated_failures=$($M.treated_failures) prevented=$($M.prevented_failures) actions=$($M.actions_taken) apparent_unnecessary=$($M.apparent_unnecessary_action_arms) near_future=$($M.near_future_action_arms) true_unnecessary=$($M.true_unnecessary_action_arms) true_fraction=$($M.true_unnecessary_fraction_among_beyond_grace_survivors)"}
 Write-Host 'WINGLESS_UP268_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
