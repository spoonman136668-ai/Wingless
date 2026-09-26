$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up151c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up151c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP151C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP151C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP151C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP151C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up151c-state-eviction-forecast -count=1;if($LASTEXITCODE-ne 0){throw 'UP151C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP151C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up151c-state-eviction-forecast)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP151C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up151c-state-eviction-forecast)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP151C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP151C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up151c-state-eviction-forecast.v1'){throw 'UP151C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 64 -or [int]$R.total_predictions-ne 64){throw 'UP151C_DESIGN'}
 if(-not $R.predictor_uses_hand -or -not $R.predictor_uses_ages -or $R.predictor_uses_identity -or $R.future_oracle_used -or $R.memory_policy_changed){throw 'UP151C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "correct=$($R.correct_predictions) total=$($R.total_predictions) accuracy=$($R.prediction_accuracy)"
 Write-Host 'WINGLESS_UP151_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
