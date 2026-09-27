$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up236c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up236c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP236C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP236C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up236c-risk-recheck-action9 -count=1;if($LASTEXITCODE-ne 0){throw 'UP236C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP236C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up236c-risk-recheck-action9)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP236C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up236c-risk-recheck-action9)|Out-String).Trim();if($P1-cne $P2){throw 'UP236C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up236c-risk-recheck-action9.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.max_actions-ne 9){throw 'UP236C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.classifier_changed -or $R.prototype_mutation_used -or $R.threshold_fitting_used -or $R.new_native_field_used -or $R.action_budget_increased -or $R.live_activation){throw 'UP236C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) cap8=$($M.cap8_losses) single=$($M.single_losses) first=$($M.first_risk_losses) recheck=$($M.recheck_risk_losses) first9=$($M.first_risk_ninth) recheck9=$($M.recheck_risk_ninth) first_rescue=$($M.first_risk_rescues) recheck_rescue=$($M.recheck_risk_rescues) first_unnecessary=$($M.first_risk_unnecessary) recheck_unnecessary=$($M.recheck_risk_unnecessary) delayed=$($M.delayed_authorizations)"}
 Write-Host 'WINGLESS_UP337_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
