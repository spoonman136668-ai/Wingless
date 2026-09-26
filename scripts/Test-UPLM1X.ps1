$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1x-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1x-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1X_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1X_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1X_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1X_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1x-mild-mass-order-robustness -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1X_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1X_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1x-mild-mass-order-robustness)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1X_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1x-mild-mass-order-robustness)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1X_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1X_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1x-mild-mass-order-robustness.v1'){throw 'UPLM1X_SCHEMA_MISMATCH'}
 if([int]$R.arms-ne 5 -or [int]$R.families-ne 5 -or [int]$R.orders-ne 5){throw 'UPLM1X_DESIGN'}
 if([int]$R.integrated_metrics.Count-ne 250 -or [int]$R.family_order_summaries.Count-ne 125 -or [int]$R.order_summaries.Count-ne 25){throw 'UPLM1X_COUNT'}
 if([int]$R.router_metrics.Count-ne 5 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM1X_BUDGET'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM1X_MASS'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.router_modified -or $R.adaptive_weighting_used -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1X_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.order_summaries){Write-Host "arm=$($S.arm) order=$($S.order) fifth=$($S.fifth_integrated_mean) prior=$($S.prior_integrated_mean) all=$($S.all_family_integrated_mean) min=$($S.min_family_integrated_mean) exact=$($S.structural_exact_all_families)"}
 Write-Host 'WINGLESS_UP156_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
