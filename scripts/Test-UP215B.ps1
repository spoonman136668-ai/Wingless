$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up215b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up215b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP215B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP215B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up215b-route-confidence-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP215B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP215B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up215b-route-confidence-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP215B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up215b-route-confidence-diagnostic)|Out-String).Trim();if($P1-cne $P2){throw 'UP215B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up215b-route-confidence-diagnostic.v1' -or [int]$R.route_decisions-ne 96 -or [int]$R.points.Count-ne 96 -or [int]$R.confusion.Count-ne 16 -or [int]$R.summaries.Count-ne 4){throw 'UP215B_DESIGN'}
 if($R.evaluation_label_fitting_used -or $R.retraining_used -or $R.adaptive_feature_selection_used -or $R.confidence_threshold_used -or $R.fallback_route_used -or $R.phase_input_used -or $R.parity_input_used -or $R.nonlinear_classifier_used){throw 'UP215B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "route_accuracy=$($R.route_accuracy) route_correct=$($R.route_correct)/$($R.route_decisions)"
 foreach($S in $R.summaries){Write-Host "regime=$($S.regime) route=$($S.correct)/$($S.total) margin_correct=$($S.mean_margin_correct) margin_incorrect=$($S.mean_margin_incorrect) nearest_correct=$($S.mean_nearest_correct) nearest_incorrect=$($S.mean_nearest_incorrect) min_margin=$($S.min_margin) max_margin=$($S.max_margin)"}
 foreach($C in $R.confusion){if([int]$C.count-gt 0){Write-Host "confusion true=$($C.true_regime) predicted=$($C.predicted_regime) count=$($C.count)"}}
 Write-Host 'WINGLESS_UP272_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
