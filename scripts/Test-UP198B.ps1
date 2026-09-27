$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up198b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up198b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP198B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP198B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP198B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP198B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up198b-triplet-transition-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP198B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP198B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up198b-triplet-transition-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP198B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up198b-triplet-transition-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP198B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP198B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up198b-triplet-transition-calibration.v1'){throw 'UP198B_SCHEMA_MISMATCH'}
 if([int]$R.pairwise_feature_count-ne 12 -or [int]$R.triplet_feature_count-ne 24 -or [int]$R.metrics.Count-ne 144 -or [int]$R.summaries.Count-ne 2){throw 'UP198B_DESIGN'}
 if($R.heldout_fitting_used -or $R.adaptive_feature_selection_used -or $R.architecture_search_used -or $R.phase_input_used -or $R.parity_input_used -or $R.maintenance_triggered){throw 'UP198B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) baseline_mae=$($S.baseline_mean_absolute_error) pairwise_mae=$($S.pairwise_mean_absolute_error) triplet_mae=$($S.triplet_mean_absolute_error) triplet_vs_pairwise=$($S.triplet_better_than_pairwise)/$($S.evaluation_points) triplet_vs_baseline=$($S.triplet_better_than_baseline)/$($S.evaluation_points)"}
 Write-Host 'WINGLESS_UP202_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
