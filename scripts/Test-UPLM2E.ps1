$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm2e-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm2e-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM2E_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM2E_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM2E_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM2E_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm2e-pername-locality-ablation -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2E_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM2E_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm2e-pername-locality-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2E_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm2e-pername-locality-ablation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM2E_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM2E_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm2e-pername-locality-ablation.v1'){throw 'UPLM2E_SCHEMA_MISMATCH'}
 if([int]$R.allocations-ne 5 -or [int]$R.training_schedules-ne 2 -or [int]$R.locality_variants-ne 4 -or [int]$R.families-ne 5){throw 'UPLM2E_DESIGN'}
 if([int]$R.metrics.Count-ne 400 -or [int]$R.family_summaries.Count-ne 200 -or [int]$R.variant_summaries.Count-ne 40){throw 'UPLM2E_COUNT'}
 if(-not $R.fifth_always_last -or $R.evaluation_retraining_used -or $R.content_changed_across_variants -or [int]$R.exact_recall_cap-ne 16){throw 'UPLM2E_BOUNDARY'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM2E_MASS'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.router_modified -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM2E_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.variant_summaries){Write-Host "allocation=$($S.allocation) schedule=$($S.training_schedule) variant=$($S.locality_variant) fifth=$($S.fifth_integrated_mean) prior=$($S.prior_integrated_mean) min=$($S.min_family_integrated_mean) exact=$($S.structural_exact_all_families)"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
