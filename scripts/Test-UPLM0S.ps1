$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0s-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0s-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0S_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0S_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0S_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0S_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0s-lexical-loss-mask -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0S_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0S_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0s-lexical-loss-mask)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0S_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0s-lexical-loss-mask)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0S_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0S_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0s-lexical-loss-mask.v1'){throw 'UPLM0S_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0S_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0S_RECALL_CAP'}
 if([int]$R.base_epochs-ne 20 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM0S_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0S_LR'}
 if($R.base_replay_used){throw 'UPLM0S_BASE_REPLAY'}
 if($R.router_retraining_used){throw 'UPLM0S_ROUTER_RETRAINING'}
 if($R.attention_used){throw 'UPLM0S_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0S_ORACLE_PRESENT'}
 if([int]$R.byte_metrics.Count-ne 4){throw 'UPLM0S_BYTE_COUNT'}
 if([int]$R.routing_metrics.Count-ne 10){throw 'UPLM0S_ROUTE_COUNT'}
 $Full=@($R.byte_metrics|Where-Object{$_.arm-ceq 'full_sentence'})
 $Mask=@($R.byte_metrics|Where-Object{$_.arm-ceq 'verb_targets_only'})
 if($Full.Count-ne 2 -or $Mask.Count-ne 2){throw 'UPLM0S_ARM_COUNT'}
 if([int]$Full[0].updated_targets_per_epoch-le [int]$Mask[0].updated_targets_per_epoch){throw 'UPLM0S_MASK_NOT_SPARSE'}
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM0S_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "arm=$($M.arm) split=$($M.split) updates=$($M.updated_targets_per_epoch) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 Write-Host 'WINGLESS_UP148_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
