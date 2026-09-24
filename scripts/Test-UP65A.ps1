$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up65a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up65a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP65A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP65A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP65A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP65A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up65a-distributed-readout -count=1;if($LASTEXITCODE-ne 0){throw 'UP65A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP65A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up65a-distributed-readout)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP65A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up65a-distributed-readout)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP65A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP65A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up65a-distributed-readout.v1'){throw 'UP65A_SCHEMA_MISMATCH'}
 if(-not $R.observer_only){throw 'UP65A_SCOPE_MISMATCH'}
 if([int]$R.factorized_dense.train_states-ne 81 -or [int]$R.factorized_dense.heldout_states-ne 162){throw 'UP65A_SPLIT_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "factorized held=$($R.factorized_dense.mean_heldout_accuracy) gate=$($R.factorized_dense.gate)"
 Write-Host "entangled held=$($R.entangled_joint.mean_heldout_accuracy) gate=$($R.entangled_joint.gate)"
 foreach($M in $R.factorized_dense.per_role){Write-Host "factorized role=$($M.role) train=$($M.train_accuracy) held=$($M.heldout_accuracy)"}
 foreach($M in $R.entangled_joint.per_role){Write-Host "entangled role=$($M.role) train=$($M.train_accuracy) held=$($M.heldout_accuracy)"}
 Write-Host 'WINGLESS_UP65_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
