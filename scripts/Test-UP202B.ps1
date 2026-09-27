$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up202b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up202b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP202B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP202B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP202B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP202B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up202b-permutation-trend-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP202B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP202B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up202b-permutation-trend-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP202B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up202b-permutation-trend-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP202B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP202B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up202b-permutation-trend-transfer.v1'){throw 'UP202B_SCHEMA_MISMATCH'}
 if([int]$R.permutation_count_per_composition-ne 24 -or [int]$R.metrics.Count-ne 144 -or [int]$R.summaries.Count-ne 2){throw 'UP202B_DESIGN'}
 if($R.heldout_fitting_used -or $R.adaptive_trend_selection_used -or $R.nonlinear_trend_used -or $R.maintenance_triggered){throw 'UP202B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) baseline_mae=$($S.baseline_mean_absolute_error) identity_mae=$($S.identity_mean_mean_absolute_error) trend_mae=$($S.identity_trend_mean_absolute_error) trend_vs_identity=$($S.trend_better_than_identity_mean)/$($S.evaluation_points) trend_vs_baseline=$($S.trend_better_than_baseline)/$($S.evaluation_points)"}
 Write-Host 'WINGLESS_UP214_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
