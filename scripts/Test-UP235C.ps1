$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up235c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up235c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP235C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP235C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up235c-risk-gated-action9 -count=1;if($LASTEXITCODE-ne 0){throw 'UP235C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP235C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up235c-risk-gated-action9)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP235C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up235c-risk-gated-action9)|Out-String).Trim();if($P1-cne $P2){throw 'UP235C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up235c-risk-gated-action9.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.max_actions-ne 9){throw 'UP235C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.threshold_fitting_used -or $R.prototype_mutation_used -or $R.learned_weights_used -or $R.new_native_field_used -or $R.retry_after_rejected_first_event -or $R.action_budget_increased -or $R.new_action_type_used -or $R.live_activation){throw 'UP235C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) cap8=$($M.cap8_losses) single=$($M.single_losses) risk=$($M.risk_losses) single9=$($M.single_ninth_actions) risk9=$($M.risk_ninth_actions) single_rescue=$($M.single_rescues) risk_rescue=$($M.risk_rescues) single_unnecessary=$($M.single_unnecessary) risk_unnecessary=$($M.risk_unnecessary) single_ratio=$($M.single_rescue_per_unnecessary) risk_ratio=$($M.risk_rescue_per_unnecessary)"}
 Write-Host 'WINGLESS_UP335_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
