$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3j-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3j-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3J_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3J_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3J_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3J_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3j-burst-timing-sweep -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3J_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3J_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3j-burst-timing-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3J_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3j-burst-timing-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3J_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3J_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3j-burst-timing-sweep.v1'){throw 'UPLM3J_SCHEMA_MISMATCH'}
 if([int]$R.burst_rounds.Count-ne 6 -or [int]$R.metrics.Count-ne 108 -or [int]$R.reach_summaries.Count-ne 30){throw 'UPLM3J_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_burst_used -or $R.adaptive_resources_used -or $R.future_schedule_used){throw 'UPLM3J_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.reach_summaries){Write-Host "burst=$($S.burst_round) reachable=$($S.reachable_budget) observations=$($S.observations) spread=$($S.failure_spread) min_failed=$($S.min_earliest_failed) max_failed=$($S.max_earliest_failed)"}
 Write-Host 'WINGLESS_UP199_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
