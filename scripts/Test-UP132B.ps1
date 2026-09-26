$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up132b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up132b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP132B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP132B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP132B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP132B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up132b-rehearsal-coverage -count=1;if($LASTEXITCODE-ne 0){throw 'UP132B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP132B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up132b-rehearsal-coverage)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP132B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up132b-rehearsal-coverage)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP132B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP132B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up132b-rehearsal-coverage.v1'){throw 'UP132B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP132B_DIM'}
 if([int]$R.grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP132B_TRAINING'}
 if([int]$R.new_grounding_subjects-ne 2 -or [int]$R.old_updates_per_epoch-ne 15){throw 'UP132B_UPDATE_COUNTS'}
 if($R.projector_recomputed -or $R.adaptive_subject_choice){throw 'UP132B_BOUNDARY'}
 if([int]$R.arm_metrics.Count-ne 3){throw "UP132B_ARM_COUNT $($R.arm_metrics.Count)"}
 foreach($M in $R.arm_metrics){
  if([int]$M.old_updates_per_epoch-ne 15 -or [int]$M.new_updates_per_epoch-ne 24){throw 'UP132B_ARM_UPDATE_COUNTS'}
 }
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.arm_metrics){Write-Host "arm=$($M.arm) primary=$($M.primary_new.overall_accuracy) secondary=$($M.secondary_new.overall_accuracy) old_h=$($M.old_heldout_retained.class_accuracy) old_u=$($M.old_unseen_retained.class_accuracy) worst=$($M.primary_new.worst_surface_accuracy)"}
 Write-Host 'WINGLESS_UP175_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
