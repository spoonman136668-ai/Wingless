$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up158c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up158c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP158C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP158C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP158C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP158C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up158c-save-durability-signal -count=1;if($LASTEXITCODE-ne 0){throw 'UP158C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP158C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up158c-save-durability-signal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP158C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up158c-save-durability-signal)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP158C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP158C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up158c-save-durability-signal.v1'){throw 'UP158C_SCHEMA_MISMATCH'}
 if([int]$R.arms-ne 16 -or [int]$R.points.Count-ne 16){throw 'UP158C_DESIGN'}
 if($R.future_schedule_used -or $R.semantic_class_used -or -not $R.counterfactual_only -or $R.live_activation){throw 'UP158C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "permanent=$($R.permanent_saves) delayed=$($R.delayed_losses) perm_mean=$($R.permanent_mean_horizon) delayed_mean=$($R.delayed_mean_horizon) auroc=$($R.horizon_auroc)"
 Write-Host 'WINGLESS_UP158_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
