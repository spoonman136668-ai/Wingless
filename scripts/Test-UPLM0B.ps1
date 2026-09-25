$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0b-long-dependency-byte -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0b-long-dependency-byte)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0b-long-dependency-byte)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0b-long-dependency-byte.v1'){throw 'UPLM0B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0B_STATE_DIM'}
 if([int]$R.epochs-ne 20){throw 'UPLM0B_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0B_LR'}
 if([int]$R.train_sentences-ne 96){throw "UPLM0B_TRAIN_SENTENCES $($R.train_sentences)"}
 if([int]$R.heldout_sentences-ne 48){throw "UPLM0B_HELD_SENTENCES $($R.heldout_sentences)"}
 if($R.attention_used){throw 'UPLM0B_ATTENTION_PRESENT'}
 if($R.exact_recall_used){throw 'UPLM0B_EXACT_RECALL_PRESENT'}
 if($R.pretrained_weights){throw 'UPLM0B_PRETRAINED_PRESENT'}
 if([int]$R.metrics.Count-ne 6){throw "UPLM0B_METRIC_COUNT $($R.metrics.Count)"}
 foreach($M in $R.metrics){
   if([double]$M.top1_accuracy-lt 0 -or [double]$M.top1_accuracy-gt 1){throw 'UPLM0B_ACCURACY_RANGE'}
   if([double]$M.dependent_first_byte_accuracy-lt 0 -or [double]$M.dependent_first_byte_accuracy-gt 1){throw 'UPLM0B_DEP_RANGE'}
   if([int]$M.dependent_tokens-le 0){throw 'UPLM0B_DEP_TOKENS'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent_accuracy=$($M.dependent_first_byte_accuracy) dependent_perplexity=$($M.dependent_perplexity)"}
 Write-Host 'WINGLESS_UP93_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
