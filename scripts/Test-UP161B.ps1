$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up161b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up161b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP161B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP161B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP161B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP161B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up161b-balanced-cleanup-schedule -count=1;if($LASTEXITCODE-ne 0){throw 'UP161B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP161B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up161b-balanced-cleanup-schedule)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP161B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up161b-balanced-cleanup-schedule)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP161B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP161B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up161b-balanced-cleanup-schedule.v1'){throw 'UP161B_SCHEMA_MISMATCH'}
 if([int]$R.arms.Count-ne 18 -or [int]$R.summaries.Count-ne 3 -or [int]$R.subjects-ne 6 -or [int]$R.schedules-ne 3){throw 'UP161B_DESIGN'}
 if([int]$R.terminal_epochs-ne 5 -or [int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15){throw 'UP161B_BUDGET'}
 foreach($A in $R.arms){if([int]$A.epoch_orders.Count-ne 5 -or [int]$A.training_pair.Count-ne 2){throw 'UP161B_ORDER_OR_PAIR_COUNT'}}
 if($R.adaptive_ordering_used -or $R.extra_updates_used){throw 'UP161B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "schedule=$($S.schedule) mean_old=$($S.mean_old_retention) min_old=$($S.min_old_retention) spread=$($S.subject_retention_spread)"}
 Write-Host 'WINGLESS_UP161_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
