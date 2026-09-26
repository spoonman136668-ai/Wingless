$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up138b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up138b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP138B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP138B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP138B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP138B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up138b-rotation-refresh-factorial -count=1;if($LASTEXITCODE-ne 0){throw 'UP138B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP138B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up138b-rotation-refresh-factorial)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP138B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up138b-rotation-refresh-factorial)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP138B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP138B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up138b-rotation-refresh-factorial.v1'){throw 'UP138B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP138B_DIM'}
 if([int]$R.grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP138B_TRAINING'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15){throw 'UP138B_COUNTS'}
 if($R.extra_updates_used -or $R.adaptive_ordering_used -or $R.projector_recomputed){throw 'UP138B_BOUNDARY'}
 if([int]$R.arm_metrics.Count-ne 4){throw "UP138B_ARM_COUNT $($R.arm_metrics.Count)"}
 if([int]$R.contrasts.Count-ne 4){throw "UP138B_CONTRAST_COUNT $($R.contrasts.Count)"}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.arm_metrics){Write-Host "arm=$($M.arm) primary=$($M.primary_new.overall_accuracy) secondary=$($M.secondary_new.overall_accuracy) old_heldout=$($M.old_heldout_retained.class_accuracy) old_unseen=$($M.old_unseen_retained.class_accuracy)"}
 foreach($C in $R.contrasts){Write-Host "contrast=$($C.metric) rotation=$($C.rotation_main_effect) refresh=$($C.refresh_main_effect) interaction=$($C.rotation_refresh_interaction)"}
 Write-Host 'WINGLESS_UP184_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
