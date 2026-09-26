$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2q-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2q-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2q-dynamic-pressure-countdown -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2Q_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2q-dynamic-pressure-countdown)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2Q_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2q-dynamic-pressure-countdown)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2Q_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2Q_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2q-dynamic-pressure-countdown.v1'){throw 'UPLM2Q_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 72 -or [int]$R.total_predictions-ne 72 -or [int]$R.snapshots.Count-ne 3){throw 'UPLM2Q_DESIGN'}
 if($R.future_oracle_used -or $R.intervention_triggered){throw 'UPLM2Q_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "exact=$($R.exact_predictions)/$($R.total_predictions) rate=$($R.exact_prediction_rate)"
 Write-Host 'WINGLESS_UP179_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
