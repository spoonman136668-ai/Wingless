$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up118b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up118b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP118B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP118B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP118B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP118B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up118b-stores-centroid-shift -count=1;if($LASTEXITCODE-ne 0){throw 'UP118B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP118B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up118b-stores-centroid-shift)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP118B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up118b-stores-centroid-shift)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP118B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP118B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up118b-stores-centroid-shift.v1'){throw 'UP118B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP118B_STATE_DIM'}
 if($R.adaptive_geometry){throw 'UP118B_ADAPTIVE_GEOMETRY'}
 if([int]$R.static_metrics.Count-ne 2){throw 'UP118B_STATIC_COUNT'}
 if([int]$R.points.Count-ne 168){throw "UP118B_POINT_COUNT $($R.points.Count)"}
 foreach($A in @('baseline','stores_centroid_shift')){
  $Pts=@($R.points|Where-Object{$_.representation_arm-ceq $A});if($Pts.Count-ne 84){throw "UP118B_ARM_$A"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.static_metrics){Write-Host "arm=$($S.arm) store_cos=$($S.store_centroid_cosine) observe_cos=$($S.observe_centroid_cosine) report_cos=$($S.report_centroid_cosine)"}
 foreach($M in $R.points){if([double]$M.stores_accuracy-lt 1 -or $M.representation_arm-ceq 'stores_centroid_shift'){Write-Host "arm=$($M.representation_arm) order=$($M.class_order) policy=$($M.replay_policy) stage=$($M.stage) stores=$($M.stores_accuracy) keeps=$($M.keeps_accuracy) base=$($M.base_lexicon_accuracy) acquired=$($M.acquired_aggregate_accuracy) margin=$($M.stores_mean_store_observe_margin)"}}
 Write-Host 'WINGLESS_UP159_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
