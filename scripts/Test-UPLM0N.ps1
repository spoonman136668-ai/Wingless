$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0n-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0n-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0N_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0N_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0N_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0N_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0n-grounded-paraphrase-order -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0N_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0N_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0n-grounded-paraphrase-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0N_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0n-grounded-paraphrase-order)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0N_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0N_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0n-grounded-paraphrase-order.v1'){throw 'UPLM0N_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0N_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0N_RECALL_CAP'}
 if([int]$R.classifier_epochs-ne 20){throw 'UPLM0N_EPOCHS'}
 if([double]$R.classifier_learning_rate-ne 0.08){throw 'UPLM0N_LR'}
 if($R.explicit_type_at_learned_inference){throw 'UPLM0N_TYPE_LEAK'}
 if($R.value_bytes_at_classifier_inference){throw 'UPLM0N_VALUE_LEAK'}
 if($R.attention_used){throw 'UPLM0N_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0N_ORACLE_PRESENT'}
 if($R.language_retraining_on_paraphrases){throw 'UPLM0N_LANGUAGE_RETRAINING'}
 if([int]$R.classifier_metrics.Count-ne 8){throw 'UPLM0N_CLASSIFIER_COUNT'}
 if([int]$R.metrics.Count-ne 16){throw 'UPLM0N_METRIC_COUNT'}
 foreach($M in $R.metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM0N_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.classifier_metrics){Write-Host "classifier split=$($M.split) label=$($M.label) accuracy=$($M.accuracy) precision=$($M.precision) recall=$($M.recall)"}
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent=$($M.dependent_first_byte_accuracy) queryset_exact=$($M.query_set_exact_accuracy) event=$($M.event_routing_accuracy) report=$($M.report_routing_accuracy) entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP133_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
