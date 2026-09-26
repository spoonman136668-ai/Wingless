$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up119c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up119c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP119C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP119C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP119C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP119C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up119c-mixed-persistence-horizon -count=1;if($LASTEXITCODE-ne 0){throw 'UP119C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP119C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up119c-mixed-persistence-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP119C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up119c-mixed-persistence-horizon)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP119C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP119C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up119c-mixed-persistence-horizon.v1'){throw 'UP119C_SCHEMA_MISMATCH'}
 if([int]$R.hot_keys-ne 12 -or [int]$R.target_candidates-ne 4 -or [int]$R.burst_distractors-ne 4){throw 'UP119C_WORKLOAD'}
 if([int]$R.admission_memory_bytes-ne 128 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.key_bits-ne 14){throw 'UP119C_MEMORY'}
 if([int]$R.churn_writes-ne 12288 -or [int]$R.episodes_per_seed-ne 32){throw 'UP119C_CHURN'}
 if($R.query_evidence_used_for_admission -or $R.semantic_priority_used -or $R.future_oracle_used -or $R.adaptive_horizon_used){throw 'UP119C_LABEL_LEAK'}
 if([int]$R.points.Count-ne 4){throw "UP119C_POINT_COUNT $($R.points.Count)"}
 foreach($M in $R.points){if([int]$M.recall_entries_used-gt 16){throw 'UP119C_RECALL_CAP_EXCEEDED'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.points){Write-Host "arm=$($M.arm) workload=$($M.workload) adjacent_admit=$($M.adjacent_admission_rate) skipped_admit=$($M.skipped_admission_rate) burst_admit=$($M.burst_admission_rate) hot=$($M.hot_accuracy) target=$($M.target_accuracy) exact=$($M.target16_exact_accuracy) fp=$($M.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP180_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
