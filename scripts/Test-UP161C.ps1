$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up161c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up161c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP161C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP161C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP161C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP161C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up161c-bounded-protection-horizon -count=1;if($LASTEXITCODE-ne 0){throw 'UP161C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP161C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up161c-bounded-protection-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP161C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up161c-bounded-protection-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP161C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP161C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up161c-bounded-protection-horizon.v1'){throw 'UP161C_SCHEMA_MISMATCH'}
 if([int]$R.arms-ne 16 -or -not $R.counterfactual_only -or $R.live_activation -or $R.future_schedule_used){throw 'UP161C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "min=$($R.min_robust_horizon) mean=$($R.mean_robust_horizon) max=$($R.max_robust_horizon)"
 Write-Host 'WINGLESS_UP161_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
