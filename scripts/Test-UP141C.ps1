$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up141c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up141c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP141C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP141C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP141C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP141C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up141c-history-pressure-threshold -count=1;if($LASTEXITCODE-ne 0){throw 'UP141C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP141C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up141c-history-pressure-threshold)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP141C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up141c-history-pressure-threshold)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP141C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP141C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up141c-history-pressure-threshold.v1'){throw 'UP141C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 10 -or [int]$R.unresolved_calls_per_arm-ne 64 -or [int]$R.candidate_boundaries_per_arm-ne 2){throw 'UP141C_DESIGN'}
 if([int]$R.qualified_levels.Count-ne 5 -or [int]$R.history_entries-ne 32 -or [int]$R.max_history_age-ne 2){throw 'UP141C_MEMORY'}
 foreach($P in $R.points){if([int]$P.unresolved_calls-ne 64){throw 'UP141C_EVENT_COUNT'};if([int]$P.unexpected_manipulation_admissions-ne 0){throw 'UP141C_UNEXPECTED_ADMISSION'}}
 if($R.memory_increased -or $R.semantic_priority_used -or $R.future_oracle_used){throw 'UP141C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "target=$($P.target) q=$($P.qualified_per_window) h2=$($P.history_present_after_boundary2_rate) entries=$($P.mean_history_entries_after_boundary2) age2ev=$($P.age2_evictions) admit=$($P.admission_rate)"}
 Write-Host 'WINGLESS_UP141_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
