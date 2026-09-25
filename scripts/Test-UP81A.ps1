$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up81a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up81a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP81A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP81A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP81A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP81A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up81a-six-role-coverage-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP81A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP81A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up81a-six-role-coverage-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP81A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up81a-six-role-coverage-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP81A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP81A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up81a-six-role-coverage-transfer.v1'){throw 'UP81A_SCHEMA_MISMATCH'}
 if([int]$R.roles-ne 6){throw 'UP81A_ROLE_COUNT'}
 if([int]$R.joint_states-ne 729){throw 'UP81A_JOINT_STATE_COUNT'}
 if([int]$R.heldout_states-ne 486){throw 'UP81A_HELDOUT_COUNT'}
 if([int]$R.points.Count-ne 2){throw "UP81A_POINT_COUNT $($R.points.Count)"}
 foreach($P in $R.points){
  if([int]$P.training_states-ne 108){throw "UP81A_TRAINING_STATES $($P.training_states)"}
  if([int]$P.training_steps-ne 1){throw 'UP81A_STEP_COUNT'}
  if([int]$P.heldout_states-ne 486){throw 'UP81A_POINT_HELDOUT_COUNT'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "representation=$($P.representation) training_states=$($P.training_states) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP81_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
