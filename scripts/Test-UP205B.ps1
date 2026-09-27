$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up205b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up205b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP205B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP205B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP205B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP205B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up205b-anchor-offset-family-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP205B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP205B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up205b-anchor-offset-family-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP205B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up205b-anchor-offset-family-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP205B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP205B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up205b-anchor-offset-family-transfer.v1'){throw 'UP205B_SCHEMA_MISMATCH'}
 if([int]$R.anchors_per_phase-ne 1 -or [int]$R.metrics.Count-ne 84 -or [int]$R.summaries.Count-ne 2){throw 'UP205B_DESIGN'}
 if($R.heldout_anchor_search_used -or $R.heldout_feature_fitting_used -or $R.maintenance_triggered){throw 'UP205B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "composition=$($S.composition) perms=$($S.permutation_count) identity_mae=$($S.identity_mean_absolute_error) anchor_mae=$($S.anchor_offset_mean_absolute_error) anchor_vs_identity=$($S.anchor_offset_better_than_identity)/$($S.evaluation_points) identity_max=$($S.identity_max_absolute_error) anchor_max=$($S.anchor_offset_max_absolute_error)"}
 Write-Host 'WINGLESS_UP225_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
