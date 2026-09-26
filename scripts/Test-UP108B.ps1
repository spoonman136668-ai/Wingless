$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up108b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up108b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP108B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP108B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP108B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP108B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up108b-stratified-replay -count=1;if($LASTEXITCODE-ne 0){throw 'UP108B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP108B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up108b-stratified-replay)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP108B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up108b-stratified-replay)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP108B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP108B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up108b-stratified-replay.v1'){throw 'UP108B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP108B_STATE_DIM'}
 if([int]$R.epochs_per_batch-ne 20){throw 'UP108B_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UP108B_LR'}
 if([int]$R.grounding_per_verb-ne 4){throw 'UP108B_GROUNDING'}
 if([int]$R.fixed_budget-ne 6){throw 'UP108B_BUDGET'}
 if([int]$R.points.Count-ne 21){throw 'UP108B_POINT_COUNT'}
 foreach($A in @('full_two_per_prior','pair_major_fixed6','verb_stratified_fixed6')){
  $Pts=@($R.points|Where-Object{$_.arm-ceq $A});if($Pts.Count-ne 7){throw "UP108B_ARM_$A"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) stage=$($M.stage) base=$($M.base_seen_accuracy) acquired=$($M.acquired_aggregate_accuracy)"}
 Write-Host 'WINGLESS_UP140_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
