$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up70a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up70a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP70A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP70A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP70A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP70A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up70a-training-coverage-ladder -count=1;if($LASTEXITCODE-ne 0){throw 'UP70A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP70A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up70a-training-coverage-ladder)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP70A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up70a-training-coverage-ladder)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP70A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP70A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up70a-training-coverage-ladder.v1'){throw 'UP70A_SCHEMA_MISMATCH'}
 if($R.observer_class_changed){throw 'UP70A_OBSERVER_CHANGED'}
 if([int]$R.points.Count-ne 6){throw 'UP70A_POINT_COUNT'}
 $Counts=@($R.points|ForEach-Object{[int]$_.training_states}|Sort-Object -Unique)
 if(($Counts -join ',') -cne '27,54,81'){throw "UP70A_TRAIN_COUNTS $($Counts -join ',')"}
 foreach($P in $R.points){if([int]$P.heldout_states-ne 162){throw 'UP70A_HELDOUT_CHANGED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "train=$($P.training_states) representation=$($P.representation) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP70_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
