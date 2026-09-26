$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0z-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0z-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0z-joint-readout -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0Z_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0z-joint-readout)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Z_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0z-joint-readout)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0Z_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0Z_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0z-joint-readout.v1'){throw 'UPLM0Z_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0Z_STATE_DIM'}
 if([double]$R.learning_rate-ne 0.08){throw 'UPLM0Z_LR'}
 if([int]$R.base_only_epochs-ne 20 -or [int]$R.paraphrase_only_epochs-ne 20 -or [int]$R.joint_interleaved_epochs-ne 20 -or [int]$R.joint_fullbatch_epochs-ne 100){throw 'UPLM0Z_EPOCHS'}
 if($R.recurrent_parameters_trained){throw 'UPLM0Z_RECURRENT_TRAINING'}
 if($R.router_used -or $R.exact_recall_used -or $R.attention_used){throw 'UPLM0Z_FORBIDDEN_MECHANISM'}
 if([int]$R.readout_parameter_count-le 0){throw 'UPLM0Z_PARAM_COUNT'}
 if([int]$R.metrics.Count-ne 16){throw "UPLM0Z_METRIC_COUNT $($R.metrics.Count)"}
 foreach($A in @('base_only','paraphrase_only','joint_interleaved','joint_fullbatch')){
  $Pts=@($R.metrics|Where-Object{$_.arm-ceq $A});if($Pts.Count-ne 4){throw "UPLM0Z_ARM_$A"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) tokens=$($M.tokens)"}
 Write-Host 'WINGLESS_UP158_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
