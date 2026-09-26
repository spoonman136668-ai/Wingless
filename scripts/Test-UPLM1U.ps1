$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1u-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1u-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1U_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1U_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1U_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1U_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1u-integrated-fixed-mass -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1U_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1U_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1u-integrated-fixed-mass)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1U_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1u-integrated-fixed-mass)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1U_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1U_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1u-integrated-fixed-mass.v1'){throw 'UPLM1U_SCHEMA_MISMATCH'}
 if([int]$R.byte_metrics.Count-ne 15 -or [int]$R.byte_summaries.Count-ne 3){throw 'UPLM1U_BYTE_COUNT'}
 if([int]$R.router_metrics.Count-ne 7 -or [int]$R.integrated_metrics.Count-ne 30){throw 'UPLM1U_INTEGRATION_COUNT'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.adaptation_epochs-ne 4 -or [int]$R.router_epochs-ne 20){throw 'UPLM1U_BUDGET'}
 if([math]::Abs([double]$R.total_mass_per_example-0.40)-gt 1e-12){throw 'UPLM1U_MASS'}
 if($R.recurrent_parameters_trained -or $R.recall_cap_changed -or $R.router_modified -or $R.adaptive_weighting_used -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1U_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.byte_summaries){Write-Host "arm=$($S.arm) fifth=$($S.fifth_heldout_accuracy) prior=$($S.prior_mean_heldout_accuracy) min=$($S.min_heldout_accuracy)"}
 Write-Host "integrated_metrics=$($R.integrated_metrics.Count) router_metrics=$($R.router_metrics.Count)"
 Write-Host 'WINGLESS_UP153_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
