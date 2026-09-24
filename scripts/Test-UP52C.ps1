$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up52c-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up52c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-52C MEMORY LOAD INTERFERENCE ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP52C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP52C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP52C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP52C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up52c-memory-load -count=1;if($LASTEXITCODE-ne 0){throw 'UP52C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP52C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up52c-memory-load)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up52c-memory-load)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP52C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up52c-memory-load-interference.v1'){throw 'UP52C_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 5){throw 'UP52C_POINT_COUNT'}
 if([double]$R.memory_noise-ne 0.05){throw 'UP52C_NOISE_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Last passing writes:  $($R.last_passing_writes)"
 Write-Host "First failing writes: $($R.first_failing_writes)"
 foreach($P in $R.points){
  Write-Host "writes=$($P.writes_per_scenario) gate=$($P.hard_gate) commit=$($P.integration.commit_decode_accuracy) final=$($P.integration.exact_final_table_accuracy) relation=$($P.integration.relational_query_accuracy) drift=$($P.integration.max_norm_drift)"
 }
 Write-Host 'WINGLESS_UP52_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
