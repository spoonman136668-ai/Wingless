$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up53a-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up53a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-53A THREE-BASIS IDENTIFIABILITY ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP53A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP53A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP53A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP53A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up53a-identifiability -count=1;if($LASTEXITCODE-ne 0){throw 'UP53A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP53A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up53a-identifiability)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up53a-identifiability)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP53A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP53A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up53a-three-basis-identifiability.v1'){throw 'UP53A_SCHEMA_MISMATCH'}
 if([int]$R.triads.Count-ne 3){throw 'UP53A_TRIAD_COUNT'}
 if([int]$R.training_basis_states-ne 3){throw 'UP53A_BASIS_COUNT'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Invariant triad pass:    $($R.invariant_triad_pass)"
 Write-Host "Informative triads pass: $($R.informative_triads_pass)"
 foreach($T in $R.triads){
  Write-Host "$($T.spec.name) gate=$($T.gate) held=$($T.unitary.heldout_accuracy) params=$($T.unitary.parameters -join ',') role1=$($T.derived_role1_accuracy) role2=$($T.derived_role2_accuracy) long=$($T.long_program_accuracy)"
 }
 Write-Host 'WINGLESS_UP53_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
