$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up225c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up225c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP225C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP225C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up225c-residual-stage-localization -count=1;if($LASTEXITCODE-ne 0){throw 'UP225C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP225C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up225c-residual-stage-localization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP225C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up225c-residual-stage-localization)|Out-String).Trim();if($P1-cne $P2){throw 'UP225C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up225c-residual-stage-localization.v1' -or [int]$R.arms.Count-ne 64 -or [int]$R.stage_summaries.Count-ne 9){throw 'UP225C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.adaptive_feature_selection_used -or $R.live_activation){throw 'UP225C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.stage_summaries){Write-Host "actions=$($S.actions_taken) losses=$($S.loss_arms) survivors=$($S.survivor_arms)"}
 foreach($A in $R.silent_loss_arms){Write-Host "silent cohort=$($A.cohort) hand=$($A.initial_hand) key=$($A.endangered_key) loss=$($A.loss_step) actions=$($A.actions_taken) steps=$($A.action_steps -join ',') adv_min=$($A.min_adversarial_horizon) nq_min=$($A.min_no_query_horizon) age_min=$($A.min_endangered_age) dist_min=$($A.min_hand_distance) pred_count=$($A.predicted_endangered_count)"}
 Write-Host "silent_snapshots=$($R.silent_loss_snapshots.Count)"
 Write-Host 'WINGLESS_UP317_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
