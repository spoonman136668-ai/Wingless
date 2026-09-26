$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1i-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1i-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1I_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1I_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1I_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1I_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1i-operator-split-reintegration -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1I_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1I_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1i-operator-split-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1I_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1i-operator-split-reintegration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1I_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1I_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1i-operator-split-reintegration.v1'){throw 'UPLM1I_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1I_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM1I_RECALL_CAP'}
 if([int]$R.base_pretrain_epochs-ne 20 -or [int]$R.joint_interleaved_epochs-ne 20 -or [int]$R.adaptation_epochs-ne 4){throw 'UPLM1I_EPOCHS'}
 if([double]$R.learning_rate-ne 0.08 -or [double]$R.half_step_learning_rate-ne 0.04){throw 'UPLM1I_LR'}
 if($R.router_retraining_used -or $R.recurrent_parameters_trained -or $R.attention_used -or $R.future_oracle_used){throw 'UPLM1I_BOUNDARY_VIOLATION'}
 if([int]$R.byte_metrics.Count-ne 6){throw "UPLM1I_BYTE_COUNT $($R.byte_metrics.Count)"}
 if([int]$R.routing_metrics.Count-ne 28){throw "UPLM1I_ROUTING_COUNT $($R.routing_metrics.Count)"}
 foreach($M in $R.routing_metrics){if([int]$M.metric.max_recall_entries-gt 16){throw 'UPLM1I_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.byte_metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 foreach($M in $R.routing_metrics){Write-Host "arm=$($M.arm) split=$($M.metric.split) accuracy=$($M.metric.top1_accuracy) exact=$($M.metric.query_set_exact_accuracy) event=$($M.metric.event_routing_accuracy) report=$($M.metric.report_routing_accuracy) entries=$($M.metric.max_recall_entries)"}
 Write-Host 'WINGLESS_UP168_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
