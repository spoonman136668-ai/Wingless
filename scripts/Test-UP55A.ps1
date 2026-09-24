$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up55a-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up55a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-55A CONTROLLED COMPOSITION ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP55A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP55A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP55A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP55A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up55a-controlled -count=1;if($LASTEXITCODE-ne 0){throw 'UP55A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP55A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up55a-controlled)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP55A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up55a-controlled)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP55A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP55A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up55a-controlled-composition.v1'){throw 'UP55A_SCHEMA_MISMATCH'}
 if(-not $R.train_single_step_only){throw 'UP55A_SEQUENCE_TRAIN_LEAK'}
 if([int]$R.trained_control_value-ne 0 -or [int]$R.trained_target_role-ne 0 -or [int]$R.trained_control_role-ne 1){throw 'UP55A_TRAIN_SCOPE_MISMATCH'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Unitary held:          $($R.unitary.heldout_accuracy)"
 Write-Host "Unseen control value:  $($R.unseen_control_value_accuracy)"
 Write-Host "Moved control:         $($R.moved_control_accuracy)"
 Write-Host "Moved target:          $($R.moved_target_accuracy)"
 Write-Host "Long program:          $($R.long_program_accuracy)"
 Write-Host "Controlled gate:       $($R.controlled_composition_gate)"
 Write-Host "Matched control held:  $($R.nonunitary_heldout_accuracy)"
 Write-Host 'WINGLESS_UP55_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}