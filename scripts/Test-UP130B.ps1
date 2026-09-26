$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up130b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up130b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP130B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP130B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP130B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP130B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up130b-dualview-zero-shot-lexical -count=1;if($LASTEXITCODE-ne 0){throw 'UP130B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP130B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up130b-dualview-zero-shot-lexical)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP130B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up130b-dualview-zero-shot-lexical)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP130B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP130B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up130b-dualview-zero-shot-lexical.v1'){throw 'UP130B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP130B_DIM'}
 if([int]$R.gate_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP130B_GATE_CONFIG'}
 if($R.new_surface_training_used -or $R.new_subject_training_used -or $R.projector_recomputed -or $R.explicit_class_at_inference){throw 'UP130B_BOUNDARY'}
 if([int]$R.split_metrics.Count-ne 3){throw 'UP130B_SPLIT_COUNT'}
 if([int]$R.surface_metrics.Count-ne 36){throw "UP130B_SURFACE_COUNT $($R.surface_metrics.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.split_metrics){Write-Host "split=$($M.split) overall=$($M.overall_accuracy) store=$($M.store_accuracy) observe=$($M.observe_accuracy) report=$($M.report_accuracy) precision=$($M.store_precision) recall=$($M.store_recall) worst=$($M.worst_surface_accuracy)"}
 Write-Host 'WINGLESS_UP169_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
