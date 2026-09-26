$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up96c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up96c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP96C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP96C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP96C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP96C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up96c-width1024-horizon -count=1;if($LASTEXITCODE-ne 0){throw 'UP96C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP96C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up96c-width1024-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP96C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up96c-width1024-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP96C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP96C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up96c-width1024-horizon.v1'){throw 'UP96C_SCHEMA_MISMATCH'}
 if([int]$R.filter_bits-ne 1024){throw 'UP96C_FILTER_BITS'}
 if([int]$R.hash_count-ne 2){throw 'UP96C_HASH_COUNT'}
 if([int]$R.generation_interval-ne 32){throw 'UP96C_INTERVAL'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP96C_RECALL_CAP'}
 if($R.future_oracle_used){throw 'UP96C_ORACLE_PRESENT'}
 if($R.query_labels_used_for_admission){throw 'UP96C_QUERY_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw 'UP96C_POINT_COUNT'}
 $Expected=@(1536,3072,6144,12288)
 for($i=0;$i-lt 4;$i++){
  if([int]$R.points[$i].churn_writes-ne $Expected[$i]){throw 'UP96C_CHURN_ORDER'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP96C_RECALL_CAP_EXCEEDED'}
  if([int]$R.points[$i].policy_metadata_bytes-ne 134){throw 'UP96C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "churn=$($M.churn_writes) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) resets=$($M.filter_resets)"}
 Write-Host 'WINGLESS_UP155_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
