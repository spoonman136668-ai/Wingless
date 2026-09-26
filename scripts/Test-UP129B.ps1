$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up129b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up129b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP129B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP129B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP129B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP129B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up129b-dualview-store-gate -count=1;if($LASTEXITCODE-ne 0){throw 'UP129B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP129B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up129b-dualview-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP129B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up129b-dualview-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP129B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP129B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up129b-dualview-store-gate.v1'){throw 'UP129B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP129B_DIM'}
 if([int]$R.epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP129B_TRAINING'}
 if($R.explicit_class_at_inference){throw 'UP129B_CLASS_LEAK'}
 if($R.projector_retrained -or $R.adaptive_geometry){throw 'UP129B_PROJECTOR_CHANGED'}
 if([int]$R.split_metrics.Count-ne 3){throw 'UP129B_SPLIT_COUNT'}
 if([int]$R.surface_metrics.Count-ne 45){throw "UP129B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 if([int]$R.points.Count-ne 84){throw "UP129B_POINT_COUNT $($R.points.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.split_metrics){Write-Host "split=$($M.split) acc=$($M.class_accuracy) store_p=$($M.store_precision) store_r=$($M.store_recall)"}
 foreach($M in $R.surface_metrics){if($M.surface-ceq 'stores'){Write-Host "stores split=$($M.split) acc=$($M.accuracy) store=$($M.mean_store_probability) obs=$($M.mean_observe_probability) raw=$($M.mean_store_raw_view_logit_contribution) proj=$($M.mean_store_projected_view_logit_contribution)"}}
 Write-Host 'WINGLESS_UP166_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
