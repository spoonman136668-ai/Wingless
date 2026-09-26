$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up145c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up145c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP145C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP145C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP145C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP145C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up145c-lexical-consolidation-integration -count=1;if($LASTEXITCODE-ne 0){throw 'UP145C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP145C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up145c-lexical-consolidation-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP145C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up145c-lexical-consolidation-integration)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP145C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP145C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up145c-lexical-consolidation-integration.v1'){throw 'UP145C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 2 -or [int]$R.candidate_generations-ne 3 -or [int]$R.candidate_calls-ne 96){throw 'UP145C_DESIGN'}
 if([int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32 -or [int]$R.hot_durable_facts-ne 12 -or [int]$R.cold_durable_facts-ne 4){throw 'UP145C_MEMORY'}
 foreach($P in $R.points){if([int]$P.candidates.Count-ne 4){throw 'UP145C_CANDIDATE_COUNT'}}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP145C_BOUNDARY'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){foreach($C in $P.candidates){Write-Host "assignment=$($P.assignment) subject=$($C.subject) verb=$($C.verb) class=$($C.class) role=$($C.role) admit=$($C.admission_count) retained=$($C.final_retained)"};Write-Host "assignment=$($P.assignment) hot=$($P.hot_durable_accuracy) cold=$($P.cold_durable_accuracy) fp=$($P.one_shot_false_admissions) panic=$($P.panic_rate)"}
 Write-Host 'WINGLESS_UP145_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
