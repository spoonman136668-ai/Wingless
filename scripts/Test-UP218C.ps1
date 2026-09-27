$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up218c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up218c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP218C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP218C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up218c-fresh-critical-ninth-action -count=1;if($LASTEXITCODE-ne 0){throw 'UP218C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP218C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up218c-fresh-critical-ninth-action)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP218C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up218c-fresh-critical-ninth-action)|Out-String).Trim();if($P1-cne $P2){throw 'UP218C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up218c-fresh-critical-ninth-action.v1' -or [int]$R.metrics.Count-ne 32 -or [int]$R.max_actions-ne 9 -or [int]$R.critical_horizon-ne 2){throw 'UP218C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.threshold_fitting_used -or $R.schedule_retuning_used -or $R.future_policy_input_used){throw 'UP218C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) horizon=$($M.horizon) cap8_losses=$($M.cap8_losses) cap9_losses=$($M.cap9_losses) ninth=$($M.ninth_actions) rescues=$($M.ninth_rescues) unnecessary=$($M.unnecessary_ninth) necessary=$($M.necessary_ninth_attempts) ineffective=$($M.ineffective_ninth) ratio=$($M.rescue_per_unnecessary)"}
 Write-Host 'WINGLESS_UP303_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
