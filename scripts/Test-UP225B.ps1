$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up225b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up225b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP225B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP225B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up225b-manifold-duplicate-attribution -count=1;if($LASTEXITCODE-ne 0){throw 'UP225B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP225B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up225b-manifold-duplicate-attribution)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP225B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up225b-manifold-duplicate-attribution)|Out-String).Trim();if($P1-cne $P2){throw 'UP225B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up225b-manifold-duplicate-attribution.v1' -or [int]$R.training_states-ne 60 -or [int]$R.family_summaries.Count-ne 4){throw 'UP225B_DESIGN'}
 if($R.evaluation_states_used -or $R.gate_changed -or $R.evaluation_derived_threshold_used -or $R.new_native_feature_used -or $R.live_activation){throw 'UP225B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "correct=$($R.training_correct_states) unique=$($R.unique_correct_geometries) duplicate_states=$($R.duplicate_correct_states) duplicate_groups=$($R.duplicate_geometry_groups) cross_family_groups=$($R.cross_family_duplicate_groups) within_family_groups=$($R.within_family_duplicate_groups) max_multiplicity=$($R.max_duplicate_multiplicity)"
 foreach($S in $R.family_summaries){Write-Host "family=$($S.family) correct=$($S.correct_states) distinct=$($S.distinct_geometries) duplicate_states=$($S.duplicate_states)"}
 Write-Host 'WINGLESS_UP343_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
