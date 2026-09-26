$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up178b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up178b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP178B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP178B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP178B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP178B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up178b-score-discrimination -count=1;if($LASTEXITCODE-ne 0){throw 'UP178B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP178B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up178b-score-discrimination)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP178B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up178b-score-discrimination)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP178B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP178B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up178b-score-discrimination.v1'){throw 'UP178B_SCHEMA_MISMATCH'}
 if([int]$R.context_phase-ne 23 -or [int]$R.metrics.Count-ne 2){throw 'UP178B_DESIGN'}
 foreach($M in $R.metrics){if([int]$M.total_slots-ne 17280){throw 'UP178B_SLOT_COUNT'}}
 if(-not $R.predictors_frozen_before_run -or $R.threshold_tuning_used -or $R.maintenance_triggered){throw 'UP178B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "predictor=$($M.predictor) auroc=$($M.auroc) ap=$($M.average_precision) base=$($M.base_crossing_rate)"}
 Write-Host 'WINGLESS_UP178_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
