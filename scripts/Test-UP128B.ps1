$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up128b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up128b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP128B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP128B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP128B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP128B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up128b-threeway-store-gate -count=1;if($LASTEXITCODE-ne 0){throw 'UP128B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP128B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up128b-threeway-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP128B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up128b-threeway-store-gate)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP128B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP128B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up128b-threeway-store-gate.v1'){throw 'UP128B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP128B_PARAMS'}
 if($R.explicit_class_at_inference -or $R.class_weighting_used -or $R.threshold_used){throw 'UP128B_BOUNDARY'}
 if([int]$R.control_balanced_binary.Count-ne 3){throw 'UP128B_CONTROL_COUNT'}
 if([int]$R.split_metrics.Count-ne 3){throw 'UP128B_SPLIT_COUNT'}
 if([int]$R.surface_metrics.Count-ne 45){throw "UP128B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.split_metrics){Write-Host "split=$($M.split) class_acc=$($M.class_accuracy) store_precision=$($M.store_gate_precision) store_recall=$($M.store_gate_recall)"}
 foreach($M in $R.surface_metrics){if($M.accuracy-lt 1){Write-Host "surface split=$($M.split) surface=$($M.surface) target=$($M.target_class) acc=$($M.accuracy) store=$($M.mean_store_probability) observe=$($M.mean_observe_probability) report=$($M.mean_report_probability)"}}
 Write-Host 'WINGLESS_UP172_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
