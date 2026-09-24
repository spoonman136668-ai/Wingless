$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR;$GoCache=Join-Path $env:TEMP ("wingless-up56a-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up56a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null;$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-56A DOUBLE CONTROLLED ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP56A_GO_ENV_FAILED'};go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP56A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP56A_ICE_BUILD_FAILED'};go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP56A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up56a-double-controlled -count=1;if($LASTEXITCODE-ne 0){throw 'UP56A_FOCUSED_TEST_FAILED'};go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP56A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up56a-double-controlled)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP56A_PROBE1_FAILED'};$P2=((go run ./cmd/unitary-up56a-double-controlled)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP56A_PROBE2_FAILED'};if($P1-cne $P2){throw 'UP56A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json;if($R.schema-cne 'wingless.up56a-double-controlled.v1'){throw 'UP56A_SCHEMA_MISMATCH'};if(-not $R.train_single_step_only){throw 'UP56A_TRAIN_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ===';Write-Host "Held=$($R.unitary.heldout_accuracy) unseen_tuple=$($R.unseen_control_tuple_accuracy) moved_target=$($R.moved_target_accuracy) long=$($R.long_program_accuracy) gate=$($R.double_controlled_gate) control=$($R.nonunitary_heldout_accuracy)";Write-Host 'WINGLESS_UP56_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue;Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue}
