$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up224c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up224c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP224C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP224C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up224c-pre-action8-geometry -count=1;if($LASTEXITCODE-ne 0){throw 'UP224C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP224C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up224c-pre-action8-geometry)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP224C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up224c-pre-action8-geometry)|Out-String).Trim();if($P1-cne $P2){throw 'UP224C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up224c-pre-action8-geometry.v1' -or [int]$R.arms.Count-ne 64 -or [int]$R.group_summaries.Count-ne 3){throw 'UP224C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.adaptive_feature_selection_used -or $R.live_activation){throw 'UP224C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.group_summaries){Write-Host "group=$($S.group) arms=$($S.arms) obs=$($S.observation_min)..$($S.observation_max) adv=$($S.min_adversarial_horizon_min)..$($S.min_adversarial_horizon_max) nq=$($S.min_no_query_horizon_min)..$($S.min_no_query_horizon_max) age=$($S.min_endangered_age_min)..$($S.min_endangered_age_max) dist=$($S.min_hand_distance_min)..$($S.min_hand_distance_max) pred=$($S.predicted_endangered_count_min)..$($S.predicted_endangered_count_max)"}
 Write-Host "pre_action8_loss_snapshots=$($R.pre_action8_loss_snapshots.Count)"
 Write-Host 'WINGLESS_UP315_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
