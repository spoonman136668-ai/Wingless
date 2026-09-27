$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up215c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up215c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP215C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP215C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up215c-native-tail-selectivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP215C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP215C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up215c-native-tail-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP215C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up215c-native-tail-selectivity)|Out-String).Trim();if($P1-cne $P2){throw 'UP215C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up215c-native-tail-selectivity.v1' -or [int]$R.metrics.Count-ne 32 -or [int]$R.near_risk_horizon-ne 4){throw 'UP215C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.threshold_fitting_used -or $R.future_policy_input_used){throw 'UP215C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "schedule=$($M.schedule) advance=$($M.phase_advance_writes) horizon=$($M.horizon) cap7_losses=$($M.cap7_losses) native_losses=$($M.native_tail_losses) eighth=$($M.eighth_actions) rescues=$($M.tail_rescues) unnecessary=$($M.unnecessary_eighth) necessary_attempts=$($M.necessary_eighth_attempts) ineffective=$($M.ineffective_eighth) no_eighth_survivors=$($M.no_eighth_survivors) ratio=$($M.rescue_per_unnecessary)"}
 Write-Host 'WINGLESS_UP297_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
