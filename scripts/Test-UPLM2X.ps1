$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2x-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2x-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2X_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2X_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2X_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2X_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2x-global-budget-allocation -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2X_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2X_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2x-global-budget-allocation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2X_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2x-global-budget-allocation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2X_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2X_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2x-global-budget-allocation.v1'){throw 'UPLM2X_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 6 -or [int]$R.arms_per_scenario-ne 6 -or [int]$R.global_rounds-ne 12 -or [int]$R.max_actions_per_round-ne 1){throw 'UPLM2X_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.future_schedule_used){throw 'UPLM2X_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "rot=$($P.identity_rotation) policy=$($P.policy) actions=$($P.actions) completed=$($P.completed) failed=$($P.failed)"}
 Write-Host 'WINGLESS_UP185_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
