$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up115b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up115b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP115B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP115B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP115B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP115B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up115b-base-class-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP115B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP115B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up115b-base-class-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP115B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up115b-base-class-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP115B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP115B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up115b-base-class-diagnostic.v1'){throw 'UP115B_SCHEMA_MISMATCH'}
 if($R.training_changed){throw 'UP115B_TRAINING_CHANGED'}
 if([int]$R.state_dimension-ne 64){throw 'UP115B_STATE_DIM'}
 if([int]$R.fixed_budget-ne 6){throw 'UP115B_BUDGET'}
 if([int]$R.points.Count-ne 84){throw 'UP115B_POINT_COUNT'}
 foreach($M in $R.points){
  if(@($M.base_verb_accuracy.PSObject.Properties).Count-ne 6){throw 'UP115B_BASE_VERB_COUNT'}
  if(@($M.acquired_verb_accuracy.PSObject.Properties).Count-ne 6){throw 'UP115B_ACQUIRED_VERB_COUNT'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points|Where-Object{[double]$_.base_seen_accuracy-lt 1}){
  Write-Host "order=$($M.class_order) policy=$($M.policy) stage=$($M.stage) base=$($M.base_seen_accuracy) STORE=$($M.base_store_accuracy) OBSERVE=$($M.base_observe_accuracy) REPORT=$($M.base_report_accuracy)"
 }
 Write-Host 'WINGLESS_UP161_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
