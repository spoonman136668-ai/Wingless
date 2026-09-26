$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0x-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0x-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0X_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0X_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0X_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0X_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0x-gradient-conflict -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0X_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0X_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0x-gradient-conflict)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0X_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0x-gradient-conflict)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0X_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0X_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0x-gradient-conflict.v1'){throw 'UPLM0X_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0X_STATE_DIM'}
 if([int]$R.base_epochs-ne 20){throw 'UPLM0X_BASE_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0X_LR'}
 if($R.adaptation_used){throw 'UPLM0X_ADAPTATION_PRESENT'}
 if($R.parameter_updates_after_base){throw 'UPLM0X_POST_BASE_UPDATE_PRESENT'}
 if([int]$R.pair_count-le 0){throw 'UPLM0X_PAIR_COUNT'}
 if([int]$R.row_metrics.Count-ne [int]$R.alphabet_bytes.Count){throw 'UPLM0X_ROW_COUNT'}
 if([double]$R.negative_pair_fraction-lt 0 -or [double]$R.negative_pair_fraction-gt 1){throw 'UPLM0X_NEGATIVE_FRACTION'}
 if([double]$R.non_positive_pair_fraction-lt 0 -or [double]$R.non_positive_pair_fraction-gt 1){throw 'UPLM0X_NONPOSITIVE_FRACTION'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "pairs=$($R.pair_count) mean_cos=$($R.mean_pair_cosine) min_cos=$($R.min_pair_cosine) max_cos=$($R.max_pair_cosine) negative_fraction=$($R.negative_pair_fraction) global_cos=$($R.global_cosine)"
 foreach($M in $R.row_metrics|Where-Object{[double]$_.cosine-lt 0}){Write-Host "negative_row byte=$($M.byte) cosine=$($M.cosine) dot=$($M.dot_product)"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
