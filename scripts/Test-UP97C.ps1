$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up97c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up97c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP97C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP97C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP97C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP97C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up97c-generation-phase -count=1;if($LASTEXITCODE-ne 0){throw 'UP97C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP97C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up97c-generation-phase)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP97C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up97c-generation-phase)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP97C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP97C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up97c-generation-phase.v1'){throw 'UP97C_SCHEMA_MISMATCH'}
 if([int]$R.filter_bits-ne 1024){throw 'UP97C_FILTER_BITS'}
 if([int]$R.hash_count-ne 2){throw 'UP97C_HASH_COUNT'}
 if([int]$R.generation_interval-ne 32){throw 'UP97C_INTERVAL'}
 if([int]$R.churn_writes-ne 12288){throw 'UP97C_CHURN'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP97C_RECALL_CAP'}
 if([int]$R.points.Count-ne 4){throw 'UP97C_POINT_COUNT'}
 $Expected=@(0,8,16,24)
 for($i=0;$i-lt 4;$i++){
  if([int]$R.points[$i].phase_offset-ne $Expected[$i]){throw 'UP97C_PHASE_ORDER'}
  if([int]$R.points[$i].recall_entries_used-gt 16){throw 'UP97C_RECALL_CAP_EXCEEDED'}
  if([int]$R.points[$i].policy_metadata_bytes-ne 134){throw 'UP97C_METADATA'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "offset=$($M.phase_offset) hot=$($M.hot_query_accuracy) exact=$($M.hot_set_exact_accuracy) fp=$($M.false_positive_churn_admissions) first_fp=$($M.first_false_positive_churn_index) episodes_fp=$($M.episodes_with_false_positive)"}
 Write-Host 'WINGLESS_UP158_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
