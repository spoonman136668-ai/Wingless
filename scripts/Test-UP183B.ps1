$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up183b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up183b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP183B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP183B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP183B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP183B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up183b-calibration-drift-sweep -count=1;if($LASTEXITCODE-ne 0){throw 'UP183B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP183B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up183b-calibration-drift-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP183B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up183b-calibration-drift-sweep)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP183B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP183B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up183b-calibration-drift-sweep.v1' -or [int]$R.points.Count-ne 4){throw 'UP183B_DESIGN'}
 if(-not $R.model_frozen_before_evaluation -or $R.evaluation_fitting_used -or $R.maintenance_triggered){throw 'UP183B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "phase=$($P.evaluation_phase) actual=$($P.actual_crossings) predicted=$($P.predicted_crossings_sum) ratio=$($P.predicted_to_actual_ratio) brier=$($P.brier_score) ece=$($P.expected_calibration_error) auroc=$($P.auroc) ap=$($P.average_precision)"}
 Write-Host 'WINGLESS_UP183_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
