$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1s-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1s-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1S_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1S_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1S_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1S_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1s-fifth-isolation-attribution -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1S_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1S_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1s-fifth-isolation-attribution)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1S_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1s-fifth-isolation-attribution)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1S_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1S_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1s-fifth-isolation-attribution.v1'){throw 'UPLM1S_SCHEMA_MISMATCH'}
 if([int]$R.adaptation_epochs-ne 4 -or [double]$R.learning_rate-ne 0.08){throw 'UPLM1S_BUDGET'}
 if([int]$R.metrics.Count-ne 15 -or [int]$R.summaries.Count-ne 3){throw 'UPLM1S_METRIC_COUNT'}
 if($R.recurrent_parameters_trained -or $R.router_used -or $R.exact_recall_used -or $R.adaptive_schedule_used -or $R.sixth_family_used){throw 'UPLM1S_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "fifth_only_minus_cyclic_fifth=$($R.fifth_only_minus_cyclic_fifth_accuracy)"
 Write-Host "fifth_only_minus_cyclic_prior_mean=$($R.fifth_only_minus_cyclic_prior_mean_accuracy)"
 Write-Host "fifth_only_gain_over_no_adaptation=$($R.fifth_only_gain_over_no_adaptation)"
 Write-Host 'WINGLESS_UPLM1S_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
