$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm1a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm1a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM1A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM1A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM1A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM1A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm1a-initialization-history -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM1A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm1a-initialization-history)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm1a-initialization-history)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM1A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM1A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm1a-initialization-history.v1'){throw 'UPLM1A_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM1A_STATE_DIM'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM1A_LR'}
 if([int]$R.base_pretrain_epochs-ne 20){throw 'UPLM1A_BASE_EPOCHS'}
 if([int]$R.readout_parameter_count-ne 1430){throw "UPLM1A_PARAM_COUNT $($R.readout_parameter_count)"}
 if($R.recurrent_parameters_trained -or $R.router_used -or $R.exact_recall_used -or $R.attention_used){throw 'UPLM1A_FORBIDDEN_MECHANISM'}
 if([int]$R.metrics.Count-ne 16){throw "UPLM1A_METRIC_COUNT $($R.metrics.Count)"}
 foreach($A in @('fresh_4','warm_4','fresh_20','warm_20')){
  if(@($R.metrics|Where-Object{$_.arm-ceq $A}).Count-ne 4){throw "UPLM1A_ARM_$A"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity)"}
 Write-Host 'WINGLESS_UP163_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
