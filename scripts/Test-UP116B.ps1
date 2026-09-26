$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up116b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up116b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP116B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP116B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP116B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP116B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up116b-stores-margin -count=1;if($LASTEXITCODE-ne 0){throw 'UP116B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP116B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up116b-stores-margin)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP116B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up116b-stores-margin)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP116B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP116B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up116b-stores-margin.v1'){throw 'UP116B_SCHEMA_MISMATCH'}
 if($R.training_changed){throw 'UP116B_TRAINING_CHANGED'}
 if([int]$R.state_dimension-ne 64){throw 'UP116B_STATE_DIM'}
 if([int]$R.fixed_budget-ne 6){throw 'UP116B_BUDGET'}
 if([int]$R.points.Count-ne 84){throw 'UP116B_POINT_COUNT'}
 foreach($M in $R.points){
  if(([int]$M.stores.predicted_store+[int]$M.stores.predicted_observe+[int]$M.stores.predicted_report)-ne 6){throw 'UP116B_STORES_COUNT'}
  if(([int]$M.keeps.predicted_store+[int]$M.keeps.predicted_observe+[int]$M.keeps.predicted_report)-ne 6){throw 'UP116B_KEEPS_COUNT'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points|Where-Object{[double]$_.stores.accuracy-lt 1}){
  Write-Host "order=$($M.class_order) policy=$($M.policy) stage=$($M.stage) stores_acc=$($M.stores.accuracy) stores_mean_margin=$($M.stores.mean_target_margin) stores_min_margin=$($M.stores.min_target_margin) predS=$($M.stores.predicted_store) predO=$($M.stores.predicted_observe) predR=$($M.stores.predicted_report) keeps_margin=$($M.keeps.mean_target_margin)"
 }
 Write-Host 'WINGLESS_UP164_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
