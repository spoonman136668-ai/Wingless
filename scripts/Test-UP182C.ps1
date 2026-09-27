$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up182c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up182c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP182C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP182C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP182C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP182C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up182c-horizon-selectivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP182C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP182C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up182c-horizon-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP182C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up182c-horizon-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP182C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP182C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up182c-horizon-selectivity.v1'){throw 'UP182C_SCHEMA_MISMATCH'}
 if([int]$R.horizons.Count-ne 7 -or [int]$R.metrics.Count-ne 14 -or [int]$R.max_actions-ne 2){throw 'UP182C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_horizon_used -or $R.warning_threshold_changed){throw 'UP182C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) horizon=$($M.horizon) baseline_failures=$($M.baseline_failures) survivors=$($M.baseline_survivors) actions=$($M.intervention_actions) action_arms=$($M.arms_with_any_action) prevented=$($M.prevented_failures) unnecessary=$($M.unnecessary_action_arms) necessary=$($M.necessary_action_arms) no_action_survivors=$($M.no_action_survivors) ratio=$($M.prevented_per_unnecessary_action_arm)"}
 Write-Host 'WINGLESS_UP218_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
