$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0o-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0o-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0O_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0O_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0O_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0O_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0o-paraphrase-byte-adaptation -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0O_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0O_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0o-paraphrase-byte-adaptation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0O_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0o-paraphrase-byte-adaptation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0O_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0O_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0o-paraphrase-byte-adaptation.v1'){throw 'UPLM0O_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0O_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0O_RECALL_CAP'}
 if([int]$R.base_epochs-ne 20){throw 'UPLM0O_BASE_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0O_LR'}
 if($R.router_retraining_used){throw 'UPLM0O_ROUTER_RETRAINING'}
 if($R.attention_used){throw 'UPLM0O_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0O_ORACLE_PRESENT'}
 if([int]$R.base_replay_passes_per_adapt_epoch-ne 1){throw 'UPLM0O_BASE_REPLAY'}
 if([int]$R.byte_metrics.Count-ne 8){throw 'UPLM0O_BYTE_METRIC_COUNT'}
 if([int]$R.routing_metrics.Count-ne 20){throw 'UPLM0O_ROUTING_METRIC_COUNT'}
 $Expected=@(0,1,2,4)
 foreach($E in $Expected){
   $B=@($R.byte_metrics|Where-Object{[int]$_.adaptation_epochs-eq $E})
   $M=@($R.routing_metrics|Where-Object{[int]$_.adaptation_epochs-eq $E})
   if($B.Count-ne 2){throw "UPLM0O_BYTE_DEPTH_$E"}
   if($M.Count-ne 5){throw "UPLM0O_ROUTE_DEPTH_$E"}
 }
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM0O_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "adapt=$($M.adaptation_epochs) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){Write-Host "adapt=$($M.adaptation_epochs) split=$($M.metric.split) accuracy=$($M.metric.top1_accuracy) perplexity=$($M.metric.perplexity) dependent=$($M.metric.dependent_first_byte_accuracy) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy)"}
 Write-Host 'WINGLESS_UP136_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
