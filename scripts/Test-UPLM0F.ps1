$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0f-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0f-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0F_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0F_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0F_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0F_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0f-learned-admission-language -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0F_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0F_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0f-learned-admission-language)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0F_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0f-learned-admission-language)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0F_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0F_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0f-learned-admission-language.v1'){throw 'UPLM0F_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0F_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0F_RECALL_CAP'}
 if($R.explicit_type_at_learned_inference){throw 'UPLM0F_TYPE_LEAK'}
 if($R.attention_used){throw 'UPLM0F_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0F_ORACLE_PRESENT'}
 if([int]$R.metrics.Count-ne 6){throw "UPLM0F_METRIC_COUNT $($R.metrics.Count)"}
 foreach($M in $R.metrics){
  if([int]$M.max_recall_entries-gt 16){throw 'UPLM0F_RECALL_CAP_EXCEEDED'}
  if([double]$M.admission_precision-lt 0 -or [double]$M.admission_precision-gt 1){throw 'UPLM0F_PRECISION_RANGE'}
  if([double]$M.admission_recall-lt 0 -or [double]$M.admission_recall-gt 1){throw 'UPLM0F_RECALL_RANGE'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent_accuracy=$($M.dependent_first_byte_accuracy) queryset_exact=$($M.query_set_exact_accuracy) precision=$($M.admission_precision) recall=$($M.admission_recall) recall_entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP105_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
