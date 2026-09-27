$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up184c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up184c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP184C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP184C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP184C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP184C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up184c-posthorizon-failure-distance -count=1;if($LASTEXITCODE-ne 0){throw 'UP184C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP184C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up184c-posthorizon-failure-distance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP184C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up184c-posthorizon-failure-distance)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP184C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP184C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up184c-posthorizon-failure-distance.v1'){throw 'UP184C_SCHEMA_MISMATCH'}
 if([int]$R.horizons.Count-ne 14 -or [int]$R.metrics.Count-ne 28 -or [int]$R.baseline_followup_limit-ne 64){throw 'UP184C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.warning_threshold_changed -or $R.adaptive_horizon_selection_used){throw 'UP184C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) horizon=$($M.horizon) survivors=$($M.baseline_survivors) unnecessary=$($M.unnecessary_action_arms) within1=$($M.fail_within_1) within2=$($M.fail_within_2) within4=$($M.fail_within_4) within8=$($M.fail_within_8) within16=$($M.fail_within_16) beyond64=$($M.survive_beyond_64) min=$($M.min_post_horizon_distance) max=$($M.max_post_horizon_distance) mean=$($M.mean_post_horizon_distance)"}
 Write-Host 'WINGLESS_UP223_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
