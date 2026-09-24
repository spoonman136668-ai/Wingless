$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up51b-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up51b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-51B PAIR05 NOISE BOUNDARY ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP51B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP51B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP51B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP51B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up51b-noise-boundary -count=1;if($LASTEXITCODE-ne 0){throw 'UP51B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP51B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up51b-noise-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP51B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up51b-noise-boundary)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP51B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP51B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up51b-noise-boundary.v1'){throw 'UP51B_SCHEMA_MISMATCH'}
 if($R.selection_run){throw 'UP51B_SELECTION_PRESENT'}
 if([int]$R.points.Count-ne 7){throw 'UP51B_POINT_COUNT'}
 if([int]$R.frozen_pair[0]-ne 0 -or [int]$R.frozen_pair[1]-ne 5){throw 'UP51B_PAIR_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Last passing noise:  $($R.last_passing_noise)"
 Write-Host "First failing noise: $($R.first_failing_noise)"
 foreach($P in $R.points){
  Write-Host "noise=$($P.memory_noise) gate=$($P.hard_gate) held=$($P.selected.static.held_out_accuracy) commit=$($P.selected.mutable_integration.commit_decode_accuracy) final=$($P.selected.mutable_integration.exact_final_table_accuracy) relation=$($P.selected.mutable_integration.relational_query_accuracy)"
 }
 Write-Host 'WINGLESS_UP51_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
