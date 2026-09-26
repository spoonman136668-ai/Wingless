$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1c-unseen-lexical-integration -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1c-unseen-lexical-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1c-unseen-lexical-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1c-unseen-lexical-integration.v1'){throw 'UPLM1C_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1C_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1C_RECALL_CAP'}
 if([int]$R.base_pretrain_epochs-ne 20 -or [int]$R.joint_interleaved_epochs-ne 20){throw 'UPLM1C_BYTE_TRAINING'}
 if([int]$R.router_grounding_epochs-ne 20){throw 'UPLM1C_ROUTER_EPOCHS'}
 if($R.byte_model_third_family_training){throw 'UPLM1C_THIRD_BYTE_TRAINING'}
 if($R.recurrent_parameters_trained){throw 'UPLM1C_RECURRENT_TRAINED'}
 if($R.attention_used){throw 'UPLM1C_ATTENTION'}
 if($R.future_oracle_used){throw 'UPLM1C_ORACLE'}
 if([int]$R.router_metrics.Count-ne 1){throw 'UPLM1C_ROUTER_COUNT'}
 if([int]$R.metrics.Count-ne 14){throw "UPLM1C_METRIC_COUNT_$($R.metrics.Count)"}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM1C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.router_metrics){Write-Host "router split=$($M.split) accuracy=$($M.accuracy) examples=$($M.examples)"}
 foreach($M in $R.metrics){Write-Host "split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent=$($M.dependent_first_byte_accuracy) exact=$($M.query_set_exact_accuracy) event=$($M.event_routing_accuracy) report=$($M.report_routing_accuracy) entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP162_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
