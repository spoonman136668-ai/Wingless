$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up226b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up226b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP226B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP226B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up226b-feature-component-degeneracy -count=1;if($LASTEXITCODE-ne 0){throw 'UP226B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP226B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up226b-feature-component-degeneracy)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP226B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up226b-feature-component-degeneracy)|Out-String).Trim();if($P1-cne $P2){throw 'UP226B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up226b-feature-component-degeneracy.v1' -or [int]$R.training_states-ne 60 -or [int]$R.features.Count-ne 4 -or [int]$R.family_summaries.Count-ne 4){throw 'UP226B_DESIGN'}
 if($R.evaluation_states_used -or $R.gate_changed -or $R.evaluation_derived_threshold_used -or $R.adaptive_feature_selection_used -or $R.new_native_feature_used -or $R.phase_input_used -or $R.live_activation){throw 'UP226B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.family_summaries){Write-Host "family=$($S.family) correct=$($S.correct_states) joint=$($S.distinct_joint_geometries) native_count=$($S.distinct_native_correct_count) mean_abs=$($S.distinct_mean_absolute_margin) near_zero=$($S.distinct_near_zero_margin_count) min_abs=$($S.distinct_min_absolute_margin)"}
 Write-Host 'WINGLESS_UP346_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
