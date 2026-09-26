$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up87c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up87c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP87C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP87C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP87C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP87C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up87c-seenonce-admission -count=1;if($LASTEXITCODE-ne 0){throw 'UP87C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP87C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up87c-seenonce-admission)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP87C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up87c-seenonce-admission)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP87C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP87C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up87c-seenonce-admission.v1'){throw 'UP87C_SCHEMA_MISMATCH'}
 if([int]$R.active_keys-ne 32){throw 'UP87C_ACTIVE_KEYS'}
 if([int]$R.churn_writes-ne 24){throw 'UP87C_CHURN'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP87C_RECALL_CAP'}
 if([int]$R.entry_payload_bytes-ne 16){throw 'UP87C_ENTRY_BYTES'}
 if($R.future_oracle_used){throw 'UP87C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP87C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw 'UP87C_POINT_COUNT'}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP87C_RECALL_CAP_EXCEEDED'}
  if($M.policy-ceq 'seen_once_gate_plus_aging' -and [int]$M.policy_metadata_bytes-ne 13){throw 'UP87C_GATE_META'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "policy=$($M.policy) hot=$($M.hot_keys) hot_acc=$($M.hot_query_accuracy) hot_exact=$($M.hot_set_exact_accuracy) cold=$($M.cold_query_accuracy) rejected=$($M.rejected_one_shot_writes) admitted2=$($M.second_write_admissions) metadata=$($M.policy_metadata_bytes)"}
 Write-Host 'WINGLESS_UP129_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
