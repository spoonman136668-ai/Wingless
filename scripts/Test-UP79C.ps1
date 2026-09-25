$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up79c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up79c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP79C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP79C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP79C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP79C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up79c-bounded-eviction-policy -count=1;if($LASTEXITCODE-ne 0){throw 'UP79C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP79C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up79c-bounded-eviction-policy)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP79C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up79c-bounded-eviction-policy)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP79C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP79C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up79c-bounded-eviction-policy.v1'){throw 'UP79C_SCHEMA_MISMATCH'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP79C_RECALL_CAP'}
 if([int]$R.entry_payload_bytes-ne 16){throw 'UP79C_ENTRY_BYTES'}
 if($R.future_oracle_used){throw 'UP79C_ORACLE_PRESENT'}
 if([int]$R.points.Count-ne 8){throw "UP79C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP79C_RECALL_CAP_EXCEEDED'}
  if($M.policy-ceq 'fifo' -and [int]$M.policy_metadata_bytes-ne 16){throw 'UP79C_FIFO_META'}
  if($M.policy-ceq 'lru' -and [int]$M.policy_metadata_bytes-ne 136){throw 'UP79C_LRU_META'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "policy=$($M.policy) active=$($M.active_keys) churn=$($M.churn_writes) hot=$($M.hot_query_accuracy) hot_exact=$($M.hot_set_exact_accuracy) cold=$($M.cold_query_accuracy) cold_exact=$($M.cold_set_exact_accuracy) entries=$($M.recall_entries_used) metadata_bytes=$($M.policy_metadata_bytes) total_bytes=$($M.total_bounded_memory_bytes)"}
 Write-Host 'WINGLESS_UP107_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
