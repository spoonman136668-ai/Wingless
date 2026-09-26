$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up89c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up89c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP89C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP89C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP89C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP89C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up89c-filter-saturation -count=1;if($LASTEXITCODE-ne 0){throw 'UP89C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP89C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up89c-filter-saturation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP89C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up89c-filter-saturation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP89C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP89C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up89c-filter-saturation.v1'){throw 'UP89C_SCHEMA_MISMATCH'}
 if([int]$R.active_keys-ne 32){throw 'UP89C_ACTIVE_KEYS'}
 if([int]$R.hot_keys-ne 16){throw 'UP89C_HOT_KEYS'}
 if([int]$R.filter_bits-ne 256){throw 'UP89C_FILTER_BITS'}
 if([int]$R.hash_count-ne 2){throw 'UP89C_HASH_COUNT'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP89C_RECALL_CAP'}
 if([int]$R.entry_payload_bytes-ne 16){throw 'UP89C_ENTRY_BYTES'}
 if($R.future_oracle_used){throw 'UP89C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP89C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw 'UP89C_POINT_COUNT'}
 $Expected=@(24,48,96,192)
 for($i=0;$i-lt 4;$i++){
   if([int]$R.points[$i].churn_writes-ne $Expected[$i]){throw 'UP89C_CHURN_ORDER'}
   if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP89C_RECALL_CAP_EXCEEDED'}
   if([int]$R.points[$i].policy_metadata_bytes-ne 37){throw 'UP89C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "churn=$($M.churn_writes) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) rejected=$($M.rejected_one_shot_writes) fp_admit=$($M.false_positive_churn_admissions) metadata=$($M.policy_metadata_bytes)"}
 Write-Host 'WINGLESS_UP135_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
