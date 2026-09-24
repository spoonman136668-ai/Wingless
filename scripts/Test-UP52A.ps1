$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up52a-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up52a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-52A BASIS COVERAGE LADDER ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP52A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP52A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP52A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP52A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up52a-basis-coverage -count=1;if($LASTEXITCODE-ne 0){throw 'UP52A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP52A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up52a-basis-coverage)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up52a-basis-coverage)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP52A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP52A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up52a-basis-coverage-ladder.v1'){throw 'UP52A_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 4){throw 'UP52A_POINT_COUNT'}
 if(-not $R.train_single_step_only){throw 'UP52A_SEQUENCE_TRAIN_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Minimum passing basis: $($R.minimum_passing_basis)"
 foreach($P in $R.points){
  Write-Host "basis=$($P.training_basis_states) coverage=$($P.training_coverage_fraction) gate=$($P.gate) held=$($P.unitary.heldout_accuracy) primitive=$($P.unseen_primitive_accuracy) role1=$($P.derived_role1_accuracy) role2=$($P.derived_role2_accuracy) long=$($P.long_program_accuracy)"
 }
 Write-Host 'WINGLESS_UP52_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
