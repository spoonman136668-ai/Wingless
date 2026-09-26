$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up98c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up98c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP98C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP98C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP98C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP98C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up98c-promotion-churn-separation -count=1;if($LASTEXITCODE-ne 0){throw 'UP98C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP98C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up98c-promotion-churn-separation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP98C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up98c-promotion-churn-separation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP98C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP98C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up98c-promotion-churn-separation.v1'){throw 'UP98C_SCHEMA_MISMATCH'}
 if([int]$R.filter_bits-ne 1024){throw 'UP98C_FILTER_BITS'}
 if([int]$R.hash_count-ne 2){throw 'UP98C_HASH_COUNT'}
 if([int]$R.generation_interval-ne 32){throw 'UP98C_INTERVAL'}
 if([int]$R.churn_writes-ne 12288){throw 'UP98C_CHURN'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP98C_RECALL_CAP'}
 if(-not $R.diagnostic_phase_label_used){throw 'UP98C_PHASE_LABEL_EXPECTED'}
 if($R.future_oracle_used){throw 'UP98C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP98C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 5){throw 'UP98C_POINT_COUNT'}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP98C_RECALL_CAP_EXCEEDED'}
  if([int]$M.policy_metadata_bytes-ne 134){throw 'UP98C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) offset=$($M.phase_offset) clear=$($M.cleared_after_promotion) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) first=$($M.first_false_positive_churn_index)"}
 Write-Host 'WINGLESS_UP162_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
