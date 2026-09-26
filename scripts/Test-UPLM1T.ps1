$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1t-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1t-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1T_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1T_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1T_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1T_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1t-fixed-mass-plasticity-map -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1T_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1T_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1t-fixed-mass-plasticity-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1T_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1t-fixed-mass-plasticity-map)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1T_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1T_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1t-fixed-mass-plasticity-map.v1'){throw 'UPLM1T_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 15 -or [int]$R.summaries.Count-ne 3){throw 'UPLM1T_METRIC_COUNT'}
 if([int]$R.adaptation_epochs-ne 4 -or [int]$R.updates_per_family_per_example-ne 1){throw 'UPLM1T_BUDGET'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM1T_TOTAL_MASS'}
 foreach($S in $R.summaries){if([math]::Abs([double]$S.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM1T_ARM_MASS'}}
 if($R.recurrent_parameters_trained -or $R.router_used -or $R.exact_recall_used -or $R.adaptive_weighting_used -or $R.sixth_family_used){throw 'UPLM1T_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "arm=$($S.arm) fifth=$($S.fifth_heldout_accuracy) prior=$($S.prior_mean_heldout_accuracy) min=$($S.min_heldout_accuracy) fifth_delta=$($S.fifth_accuracy_delta_vs_equal) prior_delta=$($S.prior_mean_delta_vs_equal)"}
 Write-Host 'WINGLESS_UP152_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
