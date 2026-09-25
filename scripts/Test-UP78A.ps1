$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up78a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up78a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP78A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP78A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP78A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP78A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up78a-one-step-coverage-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UP78A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP78A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up78a-one-step-coverage-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP78A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up78a-one-step-coverage-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP78A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP78A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up78a-one-step-coverage-boundary.v1'){throw 'UP78A_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 8){throw "UP78A_POINT_COUNT $($R.points.Count)"}
 $Levels=@($R.points|ForEach-Object{[int]$_.training_states}|Sort-Object -Unique)
 if(($Levels -join ',') -cne '27,36,45,54'){throw "UP78A_LEVELS $($Levels -join ',')"}
 foreach($P in $R.points){
  if([int]$P.training_steps-ne 1){throw 'UP78A_STEP_COUNT'}
  if([int]$P.heldout_states-ne 162){throw 'UP78A_HELDOUT_COUNT'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "training_states=$($P.training_states) steps=$($P.training_steps) representation=$($P.representation) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP78_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
