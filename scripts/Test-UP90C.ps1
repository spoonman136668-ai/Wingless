$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up90c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up90c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP90C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP90C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP90C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP90C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up90c-generational-filter -count=1;if($LASTEXITCODE-ne 0){throw 'UP90C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP90C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up90c-generational-filter)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP90C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up90c-generational-filter)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP90C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP90C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up90c-generational-filter.v1'){throw 'UP90C_SCHEMA_MISMATCH'}
 if([int]$R.active_keys-ne 32){throw 'UP90C_ACTIVE_KEYS'}
 if([int]$R.hot_keys-ne 16){throw 'UP90C_HOT_KEYS'}
 if([int]$R.filter_bits-ne 256){throw 'UP90C_FILTER_BITS'}
 if([int]$R.hash_count-ne 2){throw 'UP90C_HASH_COUNT'}
 if([int]$R.churn_writes-ne 192){throw 'UP90C_CHURN'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP90C_RECALL_CAP'}
 if([int]$R.entry_payload_bytes-ne 16){throw 'UP90C_ENTRY_BYTES'}
 if($R.future_oracle_used){throw 'UP90C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP90C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw 'UP90C_POINT_COUNT'}
 $Expected=@(0,32,64,96)
 for($i=0;$i-lt 4;$i++){
   if([int]$R.points[$i].reset_interval-ne $Expected[$i]){throw 'UP90C_INTERVAL_ORDER'}
   if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP90C_RECALL_CAP_EXCEEDED'}
   $Meta=37;if($Expected[$i]-gt 0){$Meta=38}
   if([int]$R.points[$i].policy_metadata_bytes-ne $Meta){throw 'UP90C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "reset=$($M.reset_interval) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) rejected=$($M.rejected_one_shot_writes) fp_admit=$($M.false_positive_churn_admissions) resets=$($M.filter_resets) metadata=$($M.policy_metadata_bytes)"}
 Write-Host 'WINGLESS_UP138_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
