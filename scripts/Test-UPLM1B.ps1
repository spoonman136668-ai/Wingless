$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1b-joint-readout-reintegration -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1b-joint-readout-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1b-joint-readout-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1b-joint-readout-reintegration.v1'){throw 'UPLM1B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1B_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1B_RECALL_CAP'}
 if([int]$R.base_pretrain_epochs-ne 20){throw 'UPLM1B_BASE_EPOCHS'}
 if([int]$R.joint_interleaved_epochs-ne 20){throw 'UPLM1B_JOINT_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM1B_LR'}
 if($R.recurrent_parameters_trained){throw 'UPLM1B_RECURRENT_TRAINING'}
 if($R.router_retraining_used){throw 'UPLM1B_ROUTER_RETRAINING'}
 if($R.attention_used){throw 'UPLM1B_ATTENTION'}
 if($R.future_oracle_used){throw 'UPLM1B_ORACLE'}
 if([int]$R.byte_metrics.Count-ne 2){throw 'UPLM1B_BYTE_COUNT'}
 if([int]$R.integrated_metrics.Count-ne 12){throw "UPLM1B_INTEGRATED_COUNT $($R.integrated_metrics.Count)"}
 foreach($M in $R.integrated_metrics){if([int]$M.max_recall_entries-gt 16){throw 'UPLM1B_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "byte split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.integrated_metrics){Write-Host "integrated split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent=$($M.dependent_first_byte_accuracy) exact=$($M.query_set_exact_accuracy) admission_precision=$($M.admission_precision) admission_recall=$($M.admission_recall) event=$($M.event_routing_accuracy) report=$($M.report_routing_accuracy) entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP160_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
