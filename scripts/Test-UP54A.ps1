$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up54a-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up54a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 Write-Host '=== WINGLESS UP-54A FOUR ROLE TRANSFER ==='
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP54A_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP54A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP54A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP54A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up54a-four-role -count=1;if($LASTEXITCODE-ne 0){throw 'UP54A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP54A_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up54a-four-role)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP54A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up54a-four-role)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP54A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP54A_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up54a-four-role-transfer.v1'){throw 'UP54A_SCHEMA_MISMATCH'}
 if([int]$R.roles-ne 4 -or [int]$R.joint_dimension-ne 81){throw 'UP54A_DIMENSION_MISMATCH'}
 if(-not $R.train_single_step_only){throw 'UP54A_SEQUENCE_TRAIN_LEAK'}
 if(-not $R.only_swap01_trained){throw 'UP54A_SWAP_TRAIN_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Unitary held:     $($R.unitary.heldout_accuracy)"
 Write-Host "Unseen swap12:    $($R.unseen_swap12_accuracy)"
 Write-Host "Unseen swap23:    $($R.unseen_swap23_accuracy)"
 Write-Host "Derived role1:    $($R.derived_role1_accuracy)"
 Write-Host "Derived role2:    $($R.derived_role2_accuracy)"
 Write-Host "Derived role3:    $($R.derived_role3_accuracy)"
 Write-Host "Long programs:    $($R.long_program_accuracy)"
 Write-Host "Four-role gate:   $($R.four_role_transfer_gate)"
 Write-Host "Matched control:  $($R.nonunitary_heldout_accuracy)"
 Write-Host 'WINGLESS_UP54_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}