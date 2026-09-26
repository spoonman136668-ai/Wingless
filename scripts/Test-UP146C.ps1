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
 go test ./unitary ./cmd/unitary-up146c-durable-reuse-pressure -count=1;if($LASTEXITCODE-ne 0){throw 'UP146C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP146C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up146c-durable-reuse-pressure)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up146c-durable-reuse-pressure)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP146C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP146C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up146c-durable-reuse-pressure.v1'){throw 'UP146C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 4 -or [int]$R.candidate_calls-ne 96 -or [int]$R.queries_per_refresh-ne 4 -or [int]$R.refresh_every_candidate_calls-ne 4){throw 'UP146C_DESIGN'}
 foreach($P in $R.points){if([int]$P.candidates.Count-ne 4){throw 'UP146C_CANDIDATE_COUNT'}}
 if($R.memory_policy_changed -or $R.adaptive_query_selection -or $R.memory_increased -or [int]$R.exact_recall_cap-ne 16){throw 'UP146C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "assignment=$($P.assignment) arm=$($P.protection_arm) quin4=$($P.quin_store4_accuracy) rue4=$($P.rue_store4_accuracy) protected=$($P.protected_group_accuracy) unprotected=$($P.unprotected_matched_group_accuracy) background=$($P.background_quin_accuracy)"}
 Write-Host 'WINGLESS_UP146_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
