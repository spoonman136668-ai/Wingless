$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up95c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up95c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP95C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP95C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP95C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP95C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up95c-filter-width-generational -count=1;if($LASTEXITCODE-ne 0){throw 'UP95C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP95C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up95c-filter-width-generational)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP95C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up95c-filter-width-generational)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP95C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP95C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up95c-filter-width-generational.v1'){throw 'UP95C_SCHEMA_MISMATCH'}
 if([int]$R.hash_count-ne 2 -or [int]$R.generation_interval-ne 32 -or [int]$R.churn_writes-ne 1536){throw 'UP95C_CONFIG'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP95C_RECALL_CAP'}
 if([int]$R.points.Count-ne 3){throw 'UP95C_POINT_COUNT'}
 $Expected=@(256,512,1024)
 for($i=0;$i-lt 3;$i++){
  if([int]$R.points[$i].filter_bits-ne $Expected[$i]){throw 'UP95C_WIDTH_ORDER'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP95C_RECALL_CAP_EXCEEDED'}
  $Meta=6+[int]($Expected[$i]/8)
  if([int]$R.points[$i].policy_metadata_bytes-ne $Meta){throw 'UP95C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "bits=$($M.filter_bits) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) metadata=$($M.policy_metadata_bytes)"}
 Write-Host 'WINGLESS_UP153_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
