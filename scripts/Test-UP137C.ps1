$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up137c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up137c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP137C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP137C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP137C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP137C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up137c-mixed-history-age-selectivity -count=1;if($LASTEXITCODE-ne 0){throw 'UP137C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP137C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up137c-mixed-history-age-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP137C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up137c-mixed-history-age-selectivity)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP137C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP137C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up137c-mixed-history-age-selectivity.v1'){throw 'UP137C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 2 -or [int]$R.generations-ne 7 -or [int]$R.scheduled_writes-ne 224){throw 'UP137C_DESIGN'}
 foreach($P in $R.points){if([int]$P.path_metrics.Count-ne 4){throw 'UP137C_PATH_COUNT'}}
 if([int]$R.max_history_age-ne 2 -or [int]$R.exact_recall_cap-ne 16 -or [int]$R.history_entries-ne 32){throw 'UP137C_MEMORY'}
 if($R.semantic_priority_used -or $R.query_priority_used -or $R.future_oracle_used -or $R.memory_increased){throw 'UP137C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){foreach($M in $P.path_metrics){Write-Host "arm=$($P.arm) path=$($M.path) target=$($M.target) admit=$($M.admission_rate) acc=$($M.final_accuracy)"}}
 Write-Host 'WINGLESS_UP137_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
