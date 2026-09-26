$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up182b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up182b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP182B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP182B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP182B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP182B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up182b-continuous-risk-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP182B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP182B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up182b-continuous-risk-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP182B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up182b-continuous-risk-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP182B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP182B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up182b-continuous-risk-calibration.v1'){throw 'UP182B_SCHEMA_MISMATCH'}
 if([int]$R.calibration_phase-ne 26 -or [int]$R.evaluation_phase-ne 27 -or [int]$R.calibration.total_slots-ne 17280){throw 'UP182B_DESIGN'}
 if(-not $R.feature_bins_frozen_before_evaluation -or $R.evaluation_fitting_used -or $R.maintenance_triggered){throw 'UP182B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "brier=$($R.calibration.brier_score) ece=$($R.calibration.expected_calibration_error) predicted=$($R.calibration.predicted_crossings_sum) actual=$($R.calibration.actual_crossings)"
 foreach($M in $R.metrics){Write-Host "predictor=$($M.predictor) auroc=$($M.auroc) ap=$($M.average_precision)"}
 Write-Host 'WINGLESS_UP182_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
