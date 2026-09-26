$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0k-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0k-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0K_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0K_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0K_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0K_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0k-stream-scaling -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0K_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0K_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0k-stream-scaling)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0K_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0k-stream-scaling)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0K_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0K_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0k-stream-scaling.v1'){throw 'UPLM0K_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0K_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0K_RECALL_CAP'}
 if([int]$R.classifier_epochs-ne 20){throw 'UPLM0K_EPOCHS'}
 if([double]$R.classifier_learning_rate-ne 0.08){throw 'UPLM0K_LR'}
 if($R.explicit_type_at_learned_inference){throw 'UPLM0K_TYPE_LEAK'}
 if($R.value_bytes_at_classifier_inference){throw 'UPLM0K_VALUE_LEAK'}
 if($R.attention_used){throw 'UPLM0K_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0K_ORACLE_PRESENT'}
 if([int]$R.classifier_metrics.Count-ne 8){throw 'UPLM0K_CLASSIFIER_COUNT'}
 if([int]$R.metrics.Count-ne 8){throw 'UPLM0K_METRIC_COUNT'}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM0K_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.classifier_metrics){Write-Host "classifier split=$($M.split) label=$($M.label) accuracy=$($M.accuracy) precision=$($M.precision) recall=$($M.recall)"}
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent=$($M.dependent_first_byte_accuracy) queryset_exact=$($M.query_set_exact_accuracy) event=$($M.event_routing_accuracy) report=$($M.report_routing_accuracy) entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP124_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
