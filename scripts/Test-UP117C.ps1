$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up117c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up117c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP117C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP117C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP117C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP117C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up117c-crossgeneration-confirmation -count=1;if($LASTEXITCODE-ne 0){throw 'UP117C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP117C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up117c-crossgeneration-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP117C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up117c-crossgeneration-confirmation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP117C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP117C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up117c-crossgeneration-confirmation.v1'){throw 'UP117C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.target_candidates-ne 4 -or [int]$R.burst_distractors-ne 4){throw 'UP117C_WORKLOAD'}
 if([int]$R.generation_interval-ne 32 -or [int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16){throw 'UP117C_MEMORY'}
 if([int]$R.episodes_per_seed-ne 32){throw 'UP117C_EPISODES'}
 if($R.query_evidence_used_for_admission -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.phase_label_used){throw 'UP117C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 2){throw "UP117C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){
  if([int]$M.recall_entries_used-gt 16){throw 'UP117C_RECALL_CAP_EXCEEDED'}
  if([int]$M.admission_memory_bytes-ne 128){throw 'UP117C_ADMISSION_MEMORY'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) target_admit=$($M.target_admission_rate) burst_admit=$($M.burst_admission_rate) hot=$($M.hot_accuracy) target=$($M.target_accuracy) burst=$($M.burst_accuracy) target16=$($M.target16_accuracy) exact=$($M.target16_exact_accuracy) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP175_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
