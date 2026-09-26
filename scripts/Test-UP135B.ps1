$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up135b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up135b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP135B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP135B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP135B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP135B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up135b-final-new-refresh -count=1;if($LASTEXITCODE-ne 0){throw 'UP135B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP135B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up135b-final-new-refresh)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP135B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up135b-final-new-refresh)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP135B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP135B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up135b-final-new-refresh.v1'){throw 'UP135B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128){throw 'UP135B_DIM'}
 if([int]$R.grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP135B_TRAINING'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15){throw 'UP135B_COUNTS'}
 if($R.extra_updates_used -or $R.adaptive_split_used -or $R.projector_recomputed){throw 'UP135B_BOUNDARY'}
 if([int]$R.arm_metrics.Count-ne 4){throw "UP135B_ARM_COUNT $($R.arm_metrics.Count)"}
 $Expected=@(0,4,8,12)
 for($i=0;$i-lt 4;$i++){if([int]$R.arm_metrics[$i].final_new_updates-ne $Expected[$i]){throw 'UP135B_REFRESH_ORDER'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.arm_metrics){Write-Host "arm=$($M.arm) primary=$($M.primary_new.overall_accuracy) secondary=$($M.secondary_new.overall_accuracy) old_heldout=$($M.old_heldout_retained.class_accuracy) old_unseen=$($M.old_unseen_retained.class_accuracy)"}
 Write-Host 'WINGLESS_UP176_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
