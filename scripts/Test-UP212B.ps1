$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up212b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up212b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP212B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP212B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up212b-native-regime-identification -count=1;if($LASTEXITCODE-ne 0){throw 'UP212B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP212B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up212b-native-regime-identification)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP212B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up212b-native-regime-identification)|Out-String).Trim();if($P1-cne $P2){throw 'UP212B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up212b-native-regime-identification.v1' -or [int]$R.training_points-ne 60 -or [int]$R.evaluation_points-ne 48 -or [int]$R.regimes.Count-ne 4 -or [int]$R.confusion.Count-ne 16){throw 'UP212B_DESIGN'}
 if($R.evaluation_label_fitting_used -or $R.adaptive_feature_selection_used -or $R.phase_input_used -or $R.parity_input_used -or $R.calibration_target_input_used -or $R.nonlinear_classifier_used){throw 'UP212B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "accuracy=$($R.accuracy) chance=$($R.chance_accuracy) mean_nearest=$($R.mean_nearest_distance) mean_margin=$($R.mean_confidence_margin)"
 foreach($S in $R.regimes){Write-Host "regime=$($S.regime) correct=$($S.correct)/$($S.total) accuracy=$($S.accuracy)"}
 Write-Host 'WINGLESS_UP254_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
