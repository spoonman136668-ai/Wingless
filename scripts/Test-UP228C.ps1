$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up228c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up228c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP228C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP228C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up228c-bounded-correction-composition -count=1;if($LASTEXITCODE-ne 0){throw 'UP228C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP228C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up228c-bounded-correction-composition)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP228C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up228c-bounded-correction-composition)|Out-String).Trim();if($P1-cne $P2){throw 'UP228C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up228c-bounded-correction-composition.v1' -or [int]$R.metrics.Count-ne 16 -or [int]$R.policy_arms.Count-ne 4 -or [int]$R.max_actions-ne 9){throw 'UP228C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.new_trigger_used -or $R.new_action_type_used -or $R.threshold_fitting_used -or $R.schedule_compression_used -or $R.future_policy_input_used -or $R.live_activation){throw 'UP228C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) horizon=$($M.horizon) baseline=$($M.baseline_losses) retime=$($M.retime_only_losses) tail9=$($M.tail9_only_losses) combined=$($M.combined_losses) retime_adv=$($M.retime_only_advance_events) combined_adv=$($M.combined_advance_events) tail9_actions=$($M.tail9_only_ninth_actions) combined_ninth=$($M.combined_ninth_actions) rescue_base=$($M.combined_rescue_beyond_baseline) rescue_retime=$($M.combined_rescue_beyond_retime_only) rescue_tail9=$($M.combined_rescue_beyond_tail9_only)"}
 Write-Host 'WINGLESS_UP323_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
