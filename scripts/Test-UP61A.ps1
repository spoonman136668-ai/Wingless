$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up61a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up61a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP61A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP61A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP61A_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP61A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up61a-training-budget -count=1;if($LASTEXITCODE-ne 0){throw 'UP61A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP61A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up61a-training-budget)|Out-String).Trim();$P2=((go run ./cmd/unitary-up61a-training-budget)|Out-String).Trim();if($P1-cne $P2){throw 'UP61A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up61a-training-budget.v1'){throw 'UP61A_SCHEMA_MISMATCH'};if($R.other_hyperparameters_changed){throw 'UP61A_HYPERPARAMETER_CHANGE'};if([int]$R.points.Count-ne 5){throw 'UP61A_POINT_COUNT'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';foreach($P in $R.points){Write-Host "steps=$($P.steps) double_error=$($P.double_angle_error) train=$($P.train_accuracy) long128=$($P.long128_accuracy) gate=$($P.gate)"};Write-Host 'WINGLESS_UP61_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue}
