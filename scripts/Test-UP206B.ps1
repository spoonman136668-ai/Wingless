$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up206b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up206b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP206B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP206B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP206B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP206B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up206b-two-anchor-affine-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP206B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP206B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up206b-two-anchor-affine-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP206B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up206b-two-anchor-affine-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP206B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP206B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up206b-two-anchor-affine-transfer.v1'){throw 'UP206B_SCHEMA_MISMATCH'}
 if([int]$R.anchors_per_phase-ne 2 -or [int]$R.metrics.Count-ne 210 -or [int]$R.scales.Count-ne 12 -or [int]$R.summaries.Count-ne 4){throw 'UP206B_DESIGN'}
 if($R.heldout_anchor_search_used -or $R.nonlinear_transform_used -or $R.heldout_feature_fitting_used -or $R.maintenance_triggered){throw 'UP206B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) perms=$($S.permutation_count) one_anchor_mae=$($S.one_anchor_mean_absolute_error) affine_mae=$($S.affine_mean_absolute_error) affine_vs_one=$($S.affine_better_than_one_anchor)/$($S.evaluation_points) one_max=$($S.one_anchor_max_absolute_error) affine_max=$($S.affine_max_absolute_error)"}
 foreach($X in $R.scales){Write-Host "scale composition=$($X.composition) phase=$($X.phase) value=$($X.scale) offset=$($X.offset)"}
 Write-Host 'WINGLESS_UP228_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
