$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up146c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up146c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP146C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP146C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP146C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP146C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up146c-durable-access-recency -count=1;if($LASTEXITCODE-ne 0){throw 'UP146C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP146C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up146c-durable-access-recency)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up146c-durable-access-recency)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP146C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up146c-durable-access-recency.v1'){throw 'UP146C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 8 -or [int]$R.assignments-ne 2 -or [int]$R.cold_blocks-ne 4){throw 'UP146C_DESIGN'}
 foreach($P in $R.points){if([int]$P.durable.Count-ne 16 -or [int]$P.candidates.Count-ne 4 -or [int]$P.cold_indices.Count-ne 4){throw 'UP146C_COUNT'}}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32 -or $R.memory_policy_changed -or $R.adaptive_refresh_used -or $R.semantic_priority_used -or $R.query_priority_used){throw 'UP146C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "assignment=$($P.assignment) cold=$($P.cold_block) hot=$($P.hot_durable_accuracy) cold_acc=$($P.cold_durable_accuracy) recall=$($P.recall_entries_used) fp=$($P.one_shot_false_admissions)"}
 Write-Host 'WINGLESS_UP146_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
