$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up159b-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up159b-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP159B_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP159B_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP159B_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP159B_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up159b-cleanup-order-subject-generalization -count=1;if($LASTEXITCODE-ne 0){throw 'UP159B_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP159B_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up159b-cleanup-order-subject-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP159B_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up159b-cleanup-order-subject-generalization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP159B_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP159B_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up159b-cleanup-order-subject-generalization.v1'){throw 'UP159B_SCHEMA_MISMATCH'}
 if([int]$R.metrics.Count-ne 36 -or [int]$R.subjects-ne 6 -or [int]$R.orders-ne 6 -or [int]$R.training_pairs-ne 3){throw 'UP159B_DESIGN'}
 if([int]$R.new_updates_per_epoch-ne 24 -or [int]$R.old_updates_per_epoch-ne 15 -or [int]$R.rehearsal_before-ne 3 -or [int]$R.rehearsal_after-ne 12){throw 'UP159B_BUDGET'}
 foreach($M in $R.metrics){
  if([int]$M.training_pair.Count-ne 2 -or [int]$M.heldout_new_names.Count-ne 4){throw 'UP159B_PAIR_OR_HOLDOUT_COUNT'}
  if($M.training_pair -notcontains $M.subject){throw 'UP159B_SUBJECT_NOT_IN_PAIR'}
 }
 if($R.extra_updates_used -or $R.adaptive_ordering_used){throw 'UP159B_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "subject=$($M.subject) order=$($M.cleanup_order) recovery=$($M.mean_post_report_cleanup_recovery) old=$($M.final_mean_old_retention) heldout_new=$($M.final_heldout_new_accuracy)"}
 Write-Host 'WINGLESS_UP159_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
