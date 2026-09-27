$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up230c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up230c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP230C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP230C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up230c-action9-hysteresis -count=1;if($LASTEXITCODE-ne 0){throw 'UP230C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP230C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up230c-action9-hysteresis)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP230C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up230c-action9-hysteresis)|Out-String).Trim();if($P1-cne $P2){throw 'UP230C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up230c-action9-hysteresis.v1' -or [int]$R.metrics.Count-ne 8 -or [int]$R.max_actions-ne 9 -or [int]$R.persistence_count-ne 2){throw 'UP230C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.threshold_changed -or $R.action_budget_increased -or $R.new_action_type_used -or $R.cohort_specific_tuning_used -or $R.future_policy_input_used -or $R.live_activation){throw 'UP230C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) cap8=$($M.cap8_losses) single_losses=$($M.single_critical_losses) hyst_losses=$($M.hysteresis_losses) single_ninth=$($M.single_ninth_actions) hyst_ninth=$($M.hysteresis_ninth_actions) single_rescues=$($M.single_rescues) hyst_rescues=$($M.hysteresis_rescues) single_unnecessary=$($M.single_unnecessary_ninth) hyst_unnecessary=$($M.hysteresis_unnecessary_ninth) single_ratio=$($M.single_rescue_per_unnecessary) hyst_ratio=$($M.hysteresis_rescue_per_unnecessary)"}
 Write-Host 'WINGLESS_UP328_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
