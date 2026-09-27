$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3f-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3f-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3F_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3F_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3F_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3F_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3f-budget-reachability -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3F_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3F_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3f-budget-reachability)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3F_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3f-budget-reachability)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3F_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3F_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3f-budget-reachability.v1'){throw 'UPLM3F_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 432 -or [int]$R.metrics.Count-ne 36 -or -not $R.distinct_arm_per_round){throw 'UPLM3F_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_resources_used -or $R.future_schedule_used){throw 'UPLM3F_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "profile=$($M.profile) budget=$($M.budget) start=$($M.action_start_round) throughput=$($M.throughput) reachable=$($M.reachable_budget) earliest_actions=$($M.earliest_actions) earliest_failed=$($M.earliest_failed)"}
 Write-Host 'WINGLESS_UP194_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
