$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE
$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up51a-gocache-"+$PID)
$GoTmp=Join-Path $env:TEMP ("wingless-up51a-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null
New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache
$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try {
 Write-Host '=== WINGLESS UP-51A SPARSE BINDING GENERALIZATION ==='
 go env -w GOFLAGS=-buildvcs=false
 if($LASTEXITCODE-ne 0){throw 'UP51A_GO_ENV_FAILED'}
 go clean -cache -testcache
 if($LASTEXITCODE-ne 0){throw 'UP51A_GO_CLEAN_FAILED'}
 go run ./cmd/ice build
 if($LASTEXITCODE-ne 0){throw 'UP51A_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate
 if($LASTEXITCODE-ne 0){throw 'UP51A_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up51a-sparse-binding -count=1
 if($LASTEXITCODE-ne 0){throw 'UP51A_FOCUSED_TEST_FAILED'}
 go test ./... -count=1
 if($LASTEXITCODE-ne 0){throw 'UP51A_FULL_REGRESSION_FAILED'}

 $P1=((go run ./cmd/unitary-up51a-sparse-binding)|Out-String).Trim()
 if($LASTEXITCODE-ne 0){throw 'UP51A_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up51a-sparse-binding)|Out-String).Trim()
 if($LASTEXITCODE-ne 0){throw 'UP51A_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP51A_NONDETERMINISTIC_OUTPUT'}

 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up51a-sparse-binding-generalization.v1'){throw 'UP51A_SCHEMA_MISMATCH'}
 if(-not $R.train_single_step_only){throw 'UP51A_SEQUENCE_TRAIN_LEAK'}
 if([int]$R.diagnosis.training_basis_states-ne 9){throw 'UP51A_TRAIN_BASIS_COUNT'}
 if([int]$R.diagnosis.heldout_basis_states-ne 18){throw 'UP51A_HELD_BASIS_COUNT'}

 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "Training coverage:      $($R.diagnosis.training_coverage_fraction)"
 Write-Host "Unitary train:          $($R.diagnosis.unitary_train_accuracy)"
 Write-Host "Unseen primitive:       $($R.diagnosis.unseen_primitive_accuracy)"
 Write-Host "Unseen swap12:          $($R.diagnosis.unseen_swap12_accuracy)"
 Write-Host "Derived role1:          $($R.diagnosis.derived_role1_accuracy)"
 Write-Host "Derived role2:          $($R.diagnosis.derived_role2_accuracy)"
 Write-Host "Long programs:          $($R.diagnosis.long_program_accuracy)"
 Write-Host "Aggregate held-out:     $($R.diagnosis.aggregate_heldout_accuracy)"
 Write-Host "Sparse gate:            $($R.diagnosis.sparse_generalization_gate)"
 Write-Host "Matched nonunitary:     $($R.diagnosis.nonunitary_heldout_accuracy)"
 Write-Host 'WINGLESS_UP51_HARNESS_PASS'
} finally {
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache -ErrorAction SilentlyContinue
 Remove-Item -Recurse -Force $GoTmp -ErrorAction SilentlyContinue
}
