$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up147b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up147b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP147B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP147B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP147B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP147B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up147b-tail-block-permutation -count=1;if($LASTEXITCODE-ne 0){throw 'UP147B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP147B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up147b-tail-block-permutation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP147B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up147b-tail-block-permutation)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP147B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP147B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up147b-tail-block-permutation.v1'){throw 'UP147B_SCHEMA_MISMATCH'}
 if([int]$R.epoch_metrics.Count-ne 80 -or [int]$R.cells.Count-ne 16 -or [int]$R.schedule_summaries.Count-ne 4){throw 'UP147B_METRIC_COUNT'}
 if([int]$R.schedules-ne 4 -or [int]$R.blocks_per_schedule-ne 4 -or [int]$R.epochs_per_block-ne 5 -or [int]$R.distinct_tail_sets-ne 4 -or [int]$R.switches_per_schedule-ne 3){throw 'UP147B_DESIGN'}
 if($R.adaptive_ordering_used -or $R.extra_updates_used){throw 'UP147B_BOUNDARY_LEAK'}
 foreach($C in $R.cells){if([int]$C.samples-ne 5){throw 'UP147B_CELL_SAMPLES'}}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($C in $R.cells){Write-Host "schedule=$($C.schedule) pos=$($C.block_position) shift=$($C.shift) tail=$($C.mean_tail_damage) net=$($C.mean_net_epoch_change)"}
 foreach($S in $R.schedule_summaries){Write-Host "schedule=$($S.schedule) final_old=$($S.final_mean_old_retention) final_new=$($S.final_mean_new_accuracy)"}
 Write-Host 'WINGLESS_UP147_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
