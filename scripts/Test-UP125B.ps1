$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up125b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up125b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP125B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP125B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP125B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP125B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up125b-learned-store-gate -count=1;if($LASTEXITCODE-ne 0){throw 'UP125B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP125B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up125b-learned-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP125B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up125b-learned-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP125B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP125B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up125b-learned-store-gate.v1'){throw 'UP125B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UP125B_STATE_DIM'}
 if([int]$R.gate_epochs-ne 20){throw 'UP125B_GATE_EPOCHS'}
 if([double]$R.gate_learning_rate-ne 0.08 -or [double]$R.gate_threshold-ne 0.5){throw 'UP125B_GATE_PARAMS'}
 if($R.explicit_class_at_inference){throw 'UP125B_CLASS_LEAK'}
 if($R.projector_retrained -or $R.adaptive_geometry){throw 'UP125B_PROJECTOR_CHANGED'}
 if([int]$R.gate_metrics.Count-ne 2){throw 'UP125B_GATE_METRIC_COUNT'}
 if([int]$R.points.Count-ne 168){throw "UP125B_POINT_COUNT $($R.points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.gate_metrics){Write-Host "gate split=$($M.split) acc=$($M.accuracy) precision=$($M.precision) recall=$($M.recall)"}
 foreach($M in $R.points){if([int]$M.stage-eq 6){Write-Host "arm=$($M.arm) order=$($M.class_order) policy=$($M.replay_policy) base_o=$($M.base_original_accuracy) base_u=$($M.base_unseen_accuracy) acquired=$($M.acquired_aggregate_accuracy) nonstore_error=$($M.non_store_error_rate)"}}
 Write-Host 'WINGLESS_UP165_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
