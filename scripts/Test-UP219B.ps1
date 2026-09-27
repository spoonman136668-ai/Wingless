$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up219b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up219b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP219B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP219B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up219b-confidence-error-calibration -count=1;if($LASTEXITCODE-ne 0){throw 'UP219B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP219B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up219b-confidence-error-calibration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP219B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up219b-confidence-error-calibration)|Out-String).Trim();if($P1-cne $P2){throw 'UP219B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up219b-confidence-error-calibration.v1' -or [int]$R.route_decisions-ne 96 -or [int]$R.decisions.Count-ne 96 -or [int]$R.summaries.Count-ne 10){throw 'UP219B_DESIGN'}
 if($R.evaluation_derived_threshold_used -or $R.retraining_used -or $R.phase_input_used -or $R.parity_input_used -or $R.nonlinear_model_used){throw 'UP219B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "threshold=$($R.training_derived_threshold)"
 foreach($S in $R.summaries){Write-Host "scope=$($S.scope) regime=$($S.regime) stratum=$($S.confidence_stratum) decisions=$($S.decisions) misroutes=$($S.misroutes) misroute_rate=$($S.misroute_rate) routed_mae=$($S.mean_routed_mae) static_mae=$($S.mean_static_mae) oracle_mae=$($S.mean_oracle_family_mae)"}
 Write-Host 'WINGLESS_UP298_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
