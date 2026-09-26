$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up142c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up142c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP142C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP142C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP142C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP142C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up142c-age-rank-replacement -count=1;if($LASTEXITCODE-ne 0){throw 'UP142C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP142C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up142c-age-rank-replacement)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP142C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up142c-age-rank-replacement)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP142C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP142C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up142c-age-rank-replacement.v1'){throw 'UP142C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 12 -or [int]$R.history_entries-ne 32 -or [int]$R.max_history_age-ne 2){throw 'UP142C_DESIGN'}
 foreach($P in $R.points){if(-not $P.target_present_before -or -not $P.incoming_present_after -or [int]$P.final_history_entries-ne 32){throw 'UP142C_FIXTURE'}}
 if(-not $R.direct_state_fixture -or $R.memory_increased -or $R.semantic_priority_used -or $R.future_oracle_used){throw 'UP142C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "target=$($P.target) age=$($P.target_age) idx=$($P.target_index) survives=$($P.target_present_after) evicted=$($P.evicted_key) evicted_idx=$($P.evicted_index)"}
 Write-Host 'WINGLESS_UP142_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
