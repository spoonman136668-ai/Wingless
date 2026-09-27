$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm3c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm3c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM3C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM3C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM3C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM3C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm3c-distinct-arm-throughput -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM3C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm3c-distinct-arm-throughput)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm3c-distinct-arm-throughput)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM3C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM3C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm3c-distinct-arm-throughput.v1'){throw 'UPLM3C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 72 -or [int]$R.metrics.Count-ne 6 -or [int]$R.total_budget-ne 6 -or -not $R.one_action_per_arm_per_round){throw 'UPLM3C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation){throw 'UPLM3C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "profile=$($M.profile) throughput=$($M.throughput) earliest_failed=$($M.earliest_failed) fixed_failed=$($M.fixed_failed) advantage=$($M.earliest_advantage)"}
 Write-Host 'WINGLESS_UPLM3C_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
