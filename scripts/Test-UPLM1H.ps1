$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1h-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1h-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1H_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1H_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1H_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1H_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1h-symmetric-operator-split -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1H_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1H_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1h-symmetric-operator-split)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1H_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1h-symmetric-operator-split)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1H_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1H_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1h-symmetric-operator-split.v1'){throw 'UPLM1H_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1H_STATE_DIM'}
 if([int]$R.adaptation_epochs-ne 4){throw 'UPLM1H_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08 -or [double]$R.half_step_learning_rate-ne 0.04){throw 'UPLM1H_LR'}
 if($R.recurrent_parameters_trained){throw 'UPLM1H_RECURRENT_TRAINED'}
 if($R.router_used -or $R.exact_recall_used -or $R.attention_used){throw 'UPLM1H_EXTRA_MECHANISM'}
 if([int]$R.metrics.Count-ne 24){throw "UPLM1H_METRIC_COUNT $($R.metrics.Count)"}
 if([int]$R.summaries.Count-ne 4){throw 'UPLM1H_SUMMARY_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.summaries){Write-Host "arm=$($M.arm) min=$($M.min_heldout_accuracy) mean=$($M.mean_heldout_accuracy) spread=$($M.heldout_accuracy_spread)"}
 Write-Host 'WINGLESS_UP164_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
