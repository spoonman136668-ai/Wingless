$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up185c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up185c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP185C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP185C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP185C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP185C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up185c-grace-window-selectivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP185C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP185C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up185c-grace-window-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP185C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up185c-grace-window-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP185C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP185C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up185c-grace-window-selectivity.v1'){throw 'UP185C_SCHEMA_MISMATCH'}
 if([int]$R.horizons.Count-ne 14 -or [int]$R.grace_windows.Count-ne 4 -or [int]$R.metrics.Count-ne 112){throw 'UP185C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.warning_threshold_changed -or $R.adaptive_grace_selection_used){throw 'UP185C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) horizon=$($M.horizon) grace=$($M.grace) prevented=$($M.prevented_failures_by_horizon) apparent_unnecessary=$($M.apparent_unnecessary_action_arms) near_future=$($M.near_future_action_arms) true_unnecessary=$($M.true_unnecessary_action_arms) ratio=$($M.prevented_per_true_unnecessary) true_fraction=$($M.true_unnecessary_fraction_among_beyond_grace_survivors)"}
 Write-Host 'WINGLESS_UP226_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
