$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up224b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up224b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP224B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP224B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up224b-training-manifold-degeneracy -count=1;if($LASTEXITCODE-ne 0){throw 'UP224B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP224B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up224b-training-manifold-degeneracy)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP224B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up224b-training-manifold-degeneracy)|Out-String).Trim();if($P1-cne $P2){throw 'UP224B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up224b-training-manifold-degeneracy.v1' -or [int]$R.training_states-ne 60){throw 'UP224B_DESIGN'}
 if($R.evaluation_states_used -or $R.gate_changed -or $R.evaluation_derived_threshold_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP224B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "train_states=$($R.training_states) train_correct=$($R.training_correct_states) unique_correct=$($R.unique_correct_states) duplicate_correct=$($R.duplicate_correct_states) max_multiplicity=$($R.max_duplicate_multiplicity) zero_pairs=$($R.zero_distance_pairs) original_manifold=$($R.original_manifold_threshold) unique_min=$($R.unique_nearest_min) unique_max=$($R.unique_nearest_max)"
 Write-Host 'WINGLESS_UP340_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
