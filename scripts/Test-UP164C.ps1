$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up164c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up164c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP164C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP164C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP164C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP164C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up164c-robust-horizon-warning -count=1;if($LASTEXITCODE-ne 0){throw 'UP164C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP164C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up164c-robust-horizon-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP164C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up164c-robust-horizon-warning)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP164C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP164C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up164c-robust-horizon-warning.v1'){throw 'UP164C_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 2 -or -not $R.shadow_only -or $R.corrective_action_used -or $R.exact_future_schedule_used -or -not $R.current_state_robust_horizon_used){throw 'UP164C_DESIGN'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "cadence=$($M.cadence) warned=$($M.warned_losses) missed=$($M.missed_losses) expired=$($M.expired_warning_intervals) recall=$($M.warning_recall)"}
 Write-Host 'WINGLESS_UP164_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
