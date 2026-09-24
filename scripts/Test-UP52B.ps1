$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up52b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up52b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-52B BOUNDARY WRITE CHAIN ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP52B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP52B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP52B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP52B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up52b-boundary-write -count=1;if($LASTEXITCODE-ne 0){throw 'UP52B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP52B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up52b-boundary-write)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up52b-boundary-write)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP52B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up52b-boundary-write-chain.v1'){throw 'UP52B_SCHEMA_MISMATCH'}
 if($R.selection_run){throw 'UP52B_SELECTION_PRESENT'}
 if([int]$R.series.Count-ne 2){throw 'UP52B_SERIES_COUNT'}
 foreach($S in $R.series){if([int]$S.points.Count-ne 7){throw 'UP52B_POINT_COUNT'}}
 $Low=@($R.series|Where-Object {[double]$_.memory_noise-eq 0.065})[0]
 $High=@($R.series|Where-Object {[double]$_.memory_noise-eq 0.07})[0]
 $Low16=@($Low.points|Where-Object {[int]$_.writes_per_scenario-eq 16})[0]
 $High16=@($High.points|Where-Object {[int]$_.writes_per_scenario-eq 16})[0]
 if(-not $Low16.hard_gate){throw 'UP52B_REPRO_LOW16_FAILED'}
 if($High16.hard_gate){throw 'UP52B_REPRO_HIGH16_FAILED'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.series){
  Write-Host "noise=$($S.memory_noise) last_pass=$($S.last_passing_writes) first_fail=$($S.first_failing_writes)"
  foreach($P in $S.points){
   Write-Host "  writes=$($P.writes_per_scenario) gate=$($P.hard_gate) commit=$($P.selected.commit_decode_accuracy) final=$($P.selected.exact_final_table_accuracy) relation=$($P.selected.relational_query_accuracy) delta=$($P.commit_delta)"
  }
 }
 Write-Host 'WINGLESS_UP52_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
