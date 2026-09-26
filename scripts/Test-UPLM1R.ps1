$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1r-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1r-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1R_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1R_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1R_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1R_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1r-fivefamily-gradient-conflict -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1R_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1R_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1r-fivefamily-gradient-conflict)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1R_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1r-fivefamily-gradient-conflict)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1R_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1R_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1r-fivefamily-gradient-conflict.v1'){throw 'UPLM1R_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.family_count-ne 5 -or [int]$R.pair_count-ne 10){throw 'UPLM1R_CONFIG'}
 if($R.parameter_updates_after_start -or $R.recurrent_gradients_computed -or $R.router_used -or $R.exact_recall_used -or $R.sixth_family_used){throw 'UPLM1R_BOUNDARY'}
 if([int]$R.family_metrics.Count-ne 5){throw "UPLM1R_FAMILY_COUNT $($R.family_metrics.Count)"}
 if([int]$R.pair_metrics.Count-ne 10){throw "UPLM1R_PAIR_COUNT $($R.pair_metrics.Count)"}
 foreach($M in $R.pair_metrics){if([int]$M.matched_pairs-le 0){throw 'UPLM1R_EMPTY_PAIR'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "fifth_vs_prior_mean=$($R.fifth_vs_prior_mean_aggregate_cosine) fifth_min=$($R.fifth_vs_prior_min_aggregate_cosine) prior_mean=$($R.prior_vs_prior_mean_aggregate_cosine) fifth_norm_ratio=$($R.fifth_gradient_norm_relative_to_prior_mean)"
 foreach($M in $R.pair_metrics){Write-Host "pair=$($M.family_a)-$($M.family_b) global=$($M.aggregate_cosine) matched_mean=$($M.mean_matched_cosine) neg=$($M.negative_matched_fraction)"}
 Write-Host 'WINGLESS_UP189_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
