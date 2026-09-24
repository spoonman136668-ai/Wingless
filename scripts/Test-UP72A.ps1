$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up72a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up72a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP72A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP72A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP72A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP72A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up72a-compact-factorized-54-split-replication -count=1;if($LASTEXITCODE-ne 0){throw 'UP72A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP72A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up72a-compact-factorized-54-split-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP72A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up72a-compact-factorized-54-split-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP72A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP72A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up72a-compact-factorized-54-split-replication.v1'){throw 'UP72A_SCHEMA_MISMATCH'}
 if([int]$R.points.Count-ne 6){throw 'UP72A_POINT_COUNT'}
 foreach($P in $R.points){
  if([int]$P.training_states-ne 54){throw "UP72A_TRAIN_COUNT residue=$($P.training_residue) actual=$($P.training_states)"}
  if([int]$P.heldout_states-ne 162){throw "UP72A_HELDOUT_COUNT residue=$($P.training_residue) actual=$($P.heldout_states)"}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($P in $R.points){Write-Host "residue=$($P.training_residue) representation=$($P.representation) train=$($P.training_states) held=$($P.mean_heldout_accuracy) gate=$($P.gate)"}
 Write-Host 'WINGLESS_UP72_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
