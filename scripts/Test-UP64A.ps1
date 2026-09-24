$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up64a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up64a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP64A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP64A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP64A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP64A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up64a-five-role -count=1;if($LASTEXITCODE-ne 0){throw 'UP64A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP64A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up64a-five-role)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP64A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up64a-five-role)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP64A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP64A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up64a-five-role-zero-shot.v1'){throw 'UP64A_SCHEMA_MISMATCH'}
 if($R.five_role_training_used){throw 'UP64A_TRAINING_LEAK'}
 if([int]$R.joint_dimension-ne 243){throw 'UP64A_DIMENSION_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "aggregate=$($R.aggregate_accuracy) gate=$($R.gate) control=$($R.matched_control_accuracy)"
 foreach($M in $R.category_metrics){Write-Host "category=$($M.category) samples=$($M.samples) accuracy=$($M.accuracy)"}
 Write-Host 'WINGLESS_UP64_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
