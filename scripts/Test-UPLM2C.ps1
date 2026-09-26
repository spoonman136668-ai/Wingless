$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2c-balanced-prior-sequence-integration -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2c-balanced-prior-sequence-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2c-balanced-prior-sequence-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2c-balanced-prior-sequence-integration.v1'){throw 'UPLM2C_SCHEMA_MISMATCH'}
 if([int]$R.allocations-ne 5 -or [int]$R.training_schedules-ne 2 -or [int]$R.evaluation_orders-ne 5 -or [int]$R.families-ne 5){throw 'UPLM2C_DESIGN'}
 if([int]$R.metrics.Count-ne 500 -or [int]$R.family_summaries.Count-ne 250 -or [int]$R.order_summaries.Count-ne 50){throw 'UPLM2C_COUNT'}
 if(-not $R.fifth_always_last -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM2C_BUDGET'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM2C_MASS'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.router_modified -or $R.adaptive_weighting_used -or $R.adaptive_order_used -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.order_summaries){Write-Host "allocation=$($S.allocation) training=$($S.training_schedule) eval=$($S.evaluation_order) fifth=$($S.fifth_integrated_mean) prior=$($S.prior_integrated_mean) min=$($S.min_family_integrated_mean) exact=$($S.structural_exact_all_families)"}
 Write-Host 'WINGLESS_UP160_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
