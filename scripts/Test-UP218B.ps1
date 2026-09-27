$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up218b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up218b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP218B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP218B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up218b-late-confidence-blend-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP218B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP218B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up218b-late-confidence-blend-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP218B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up218b-late-confidence-blend-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP218B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up218b-late-confidence-blend-transfer.v1' -or [int]$R.route_decisions-ne 96 -or [int]$R.metrics.Count-ne 1776 -or [int]$R.summaries.Count-ne 4){throw 'UP218B_DESIGN'}
 if($R.evaluation_derived_threshold_used -or $R.adaptive_weighting_used -or $R.retraining_used -or $R.phase_input_used -or $R.parity_input_used -or $R.nonlinear_model_used){throw 'UP218B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "threshold=$($R.training_derived_threshold) route=$($R.route_correct)/$($R.route_decisions) blends=$($R.blend_decisions)"
 foreach($S in $R.summaries){Write-Host "regime=$($S.regime) route=$($S.route_correct)/$($S.route_total) blends=$($S.blend_decisions) static_mae=$($S.static_mean_absolute_error) routed_mae=$($S.routed_mean_absolute_error) blend_mae=$($S.blend_mean_absolute_error) oracle_family_mae=$($S.oracle_family_mean_absolute_error) blend_vs_routed=$($S.blend_better_than_routed)/$($S.evaluation_points)"}
 Write-Host 'WINGLESS_UP290_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
