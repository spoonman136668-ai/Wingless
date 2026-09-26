$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up177b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up177b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP177B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP177B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP177B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP177B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up177b-threshold-frontier -count=1;if($LASTEXITCODE-ne 0){throw 'UP177B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP177B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up177b-threshold-frontier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP177B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up177b-threshold-frontier)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP177B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP177B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up177b-threshold-frontier.v1'){throw 'UP177B_SCHEMA_MISMATCH'}
 if([int]$R.context_phase-ne 22 -or [int]$R.metrics.Count-ne 3 -or [int]$R.total_slots-ne 17280){throw 'UP177B_DESIGN'}
 if(-not $R.predictor_frozen_before_run -or $R.adaptive_threshold_used -or $R.maintenance_triggered){throw 'UP177B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "threshold=$($M.threshold) precision=$($M.precision) recall=$($M.recall) warnings=$($M.warning_count)"}
 Write-Host 'WINGLESS_UP177_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
