$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up220b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up220b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP220B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP220B_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up220b-confidence-risk-strata -count=1;if($LASTEXITCODE-ne 0){throw 'UP220B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP220B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up220b-confidence-risk-strata)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP220B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up220b-confidence-risk-strata)|Out-String).Trim();if($P1-cne $P2){throw 'UP220B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up220b-confidence-risk-strata.v1' -or [int]$R.decisions.Count-ne 128 -or [int]$R.summaries.Count-ne 5 -or [int]$R.strata.Count-ne 5){throw 'UP220B_DESIGN'}
 if($R.evaluation_derived_boundary_used -or $R.adaptive_binning_used -or $R.retraining_used -or $R.phase_input_used -or $R.parity_input_used -or $R.nonlinear_model_used){throw 'UP220B_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "train_n=$($R.training_correct_margins) min=$($R.training_min) q25=$($R.training_q25) q50=$($R.training_q50) q75=$($R.training_q75)"
 foreach($S in $R.summaries){Write-Host "stratum=$($S.stratum) decisions=$($S.decisions) misroutes=$($S.misroutes) misroute_rate=$($S.misroute_rate) routed_mae=$($S.mean_routed_mae) static_mae=$($S.mean_static_mae) oracle_mae=$($S.mean_oracle_family_mae)"}
 Write-Host 'WINGLESS_UP306_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
