$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up141b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up141b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP141B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP141B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP141B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP141B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up141b-tail-diversity-replication -count=1;if($LASTEXITCODE-ne 0){throw 'UP141B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP141B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up141b-tail-diversity-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP141B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up141b-tail-diversity-replication)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP141B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP141B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up141b-tail-diversity-replication.v1'){throw 'UP141B_SCHEMA_MISMATCH'}
 if([int]$R.state_dimension-ne 64 -or [int]$R.gate_input_dimension-ne 128 -or [int]$R.replication_surface_count-ne 12){throw 'UP141B_CONFIG'}
 if([int]$R.grounding_epochs-ne 20 -or [double]$R.learning_rate-ne 0.08){throw 'UP141B_TRAINING'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.final_new_updates-ne 4){throw 'UP141B_COUNTS'}
 if($R.extra_updates_used -or $R.adaptive_shift_choice -or $R.projector_recomputed){throw 'UP141B_BOUNDARY'}
 if([int]$R.arm_metrics.Count-ne 3){throw "UP141B_ARM_COUNT $($R.arm_metrics.Count)"}
 $Expected=@(1,8,20)
 for($i=0;$i-lt 3;$i++){if([int]$R.arm_metrics[$i].distinct_tail_sets-ne $Expected[$i]){throw 'UP141B_DIVERSITY_ORDER'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "zero_primary=$($R.zero_shot_primary.overall_accuracy) zero_secondary=$($R.zero_shot_secondary.overall_accuracy)"
 foreach($M in $R.arm_metrics){Write-Host "arm=$($M.arm) diversity=$($M.distinct_tail_sets) primary=$($M.primary_new.overall_accuracy) secondary=$($M.secondary_new.overall_accuracy) old_mean=$($M.mean_old_retention) new_mean=$($M.mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP190_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
