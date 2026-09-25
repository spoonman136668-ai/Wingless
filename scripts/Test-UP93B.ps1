$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up93b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up93b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP93B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP93B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP93B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP93B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up93b-confidence-change-routing -count=1;if($LASTEXITCODE-ne 0){throw 'UP93B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP93B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up93b-confidence-change-routing)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP93B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up93b-confidence-change-routing)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP93B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP93B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up93b-confidence-change-routing.v1'){throw 'UP93B_SCHEMA_MISMATCH'}
 if([int]$R.exact_recall_cap-ne 16){throw 'UP93B_RECALL_CAP'}
 if([int]$R.recurrent_state_bytes-ne 512){throw 'UP93B_STATE_BYTES'}
 if($R.threshold_selection){throw 'UP93B_THRESHOLD_SELECTION'}
 if([int]$R.points.Count-ne 60){throw "UP93B_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP93B_RECALL_CAP_EXCEEDED'}
  if([double]$M.admission_precision-lt 0 -or [double]$M.admission_precision-gt 1){throw 'UP93B_PRECISION_RANGE'}
  if([double]$M.admission_recall-lt 0 -or [double]$M.admission_recall-gt 1){throw 'UP93B_RECALL_RANGE'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) writes=$($M.total_writes) targets=$($M.target_keys) query=$($M.query_accuracy) exact=$($M.exact_target_set_accuracy) threshold=$($M.margin_threshold) precision=$($M.admission_precision) recall=$($M.admission_recall) false_positive=$($M.false_positive_admissions)"}
 Write-Host 'WINGLESS_UP97_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
