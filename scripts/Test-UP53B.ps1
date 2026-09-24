$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up53b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up53b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-53B ORACLE RESET DECOMPOSITION ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP53B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP53B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP53B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP53B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up53b-oracle-reset -count=1;if($LASTEXITCODE-ne 0){throw 'UP53B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP53B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up53b-oracle-reset)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up53b-oracle-reset)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP53B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up53b-oracle-reset-decomposition.v1'){throw 'UP53B_SCHEMA_MISMATCH'}
 if(-not $R.oracle_reset_used_for_path -or $R.training_changed){throw 'UP53B_BOUNDARY_MISMATCH'}
 if([int]$R.points.Count-ne 6){throw 'UP53B_POINT_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){
  Write-Host "noise=$($P.memory_noise) writes=$($P.writes_per_scenario) closed_gate=$($P.closed_loop_gate) closed_commit=$($P.closed_loop.commit_decode_accuracy) oracle_commit=$($P.oracle_reset.commit_decode_accuracy) delta=$($P.oracle_minus_closed_commit) closed_final=$($P.closed_loop.exact_final_table_accuracy) oracle_final=$($P.oracle_reset.exact_final_table_accuracy)"
 }
 Write-Host 'WINGLESS_UP53_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}