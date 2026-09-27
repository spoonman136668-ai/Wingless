$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up171c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up171c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP171C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP171C_GO_CLEAN_FAILED'}
 go run ./cmd/ice build;if($LASTEXITCODE-ne 0){throw 'UP171C_ICE_BUILD_FAILED'}
 go run ./cmd/ice validate;if($LASTEXITCODE-ne 0){throw 'UP171C_ICE_VALIDATE_FAILED'}
 go test ./unitary ./cmd/unitary-up171c-two-shot-state-space -count=1;if($LASTEXITCODE-ne 0){throw 'UP171C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP171C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up171c-two-shot-state-space)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP171C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up171c-two-shot-state-space)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP171C_PROBE2_FAILED'}
 if($P1-cne $P2){throw 'UP171C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up171c-two-shot-state-space.v1'){throw 'UP171C_SCHEMA_MISMATCH'}
 if([int]$R.initial_hands.Count-ne 16 -or [int]$R.metrics.Count-ne 16 -or [int]$R.no_refresh_cadence2_targeted_subsets.Count-ne 2 -or [int]$R.max_actions_per_arm-ne 2){throw 'UP171C_DESIGN'}
 if(-not $R.counterfactual_only -or $R.live_activation -or $R.adaptive_hand_selection_used -or $R.adaptive_action_budget_used -or $R.warning_threshold_changed){throw 'UP171C_BOUNDARY_LEAK'}
 Write-Host $P1
 Write-Host ''
 Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($M in $R.metrics){Write-Host "policy=$($M.policy) cadence=$($M.cadence) intervention=$($M.intervention) arms=$($M.arms) baseline=$($M.baseline_losses) treated=$($M.treated_losses) prevented=$($M.prevented_losses) efficiency=$($M.prevented_per_action) delay=$($M.mean_loss_delay_writes)"}
 foreach($S in $R.no_refresh_cadence2_targeted_subsets){Write-Host "subset=$($S.subset) arms=$($S.arms) baseline=$($S.baseline_losses) treated=$($S.treated_losses) prevented=$($S.prevented_losses)"}
 Write-Host 'WINGLESS_UP171_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
