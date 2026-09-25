$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-uplm0e-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-uplm0e-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UPLM0E_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UPLM0E_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UPLM0E_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UPLM0E_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up-lm0e-mutable-binding-language -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0E_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UPLM0E_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up-lm0e-mutable-binding-language)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0E_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up-lm0e-mutable-binding-language)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UPLM0E_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UPLM0E_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up-lm0e-mutable-binding-language.v1'){throw 'UPLM0E_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64){throw 'UPLM0E_STATE_DIM'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UPLM0E_RECALL_CAP'}
 if([int]$R.train_paragraphs-ne 96){throw "UPLM0E_TRAIN_PARAGRAPHS $($R.train_paragraphs)"}
 if([int]$R.heldout_paragraphs-ne 48){throw "UPLM0E_HELD_PARAGRAPHS $($R.heldout_paragraphs)"}
 if($R.attention_used){throw 'UPLM0E_ATTENTION_PRESENT'}
 if($R.future_oracle_used){throw 'UPLM0E_ORACLE_PRESENT'}
 if([int]$R.metrics.Count-ne 6){throw "UPLM0E_METRIC_COUNT $($R.metrics.Count)"}
 foreach($M in $R.metrics){
  if([int]$M.max_recall_entries-gt 16){throw 'UPLM0E_RECALL_CAP_EXCEEDED'}
  if([int]$M.exact_recall_bytes-gt 256){throw 'UPLM0E_RECALL_BYTES_EXCEEDED'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "arm=$($M.arm) split=$($M.split) accuracy=$($M.top1_accuracy) perplexity=$($M.perplexity) dependent_accuracy=$($M.dependent_first_byte_accuracy) queryset_exact=$($M.query_set_exact_accuracy) update0=$($M.update0_exact_accuracy) update1=$($M.update1_exact_accuracy) update2=$($M.update2_exact_accuracy) update4=$($M.update4_exact_accuracy) recall_entries=$($M.max_recall_entries)"}
 Write-Host 'WINGLESS_UP102_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
