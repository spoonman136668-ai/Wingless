$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up205c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up205c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP205C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP205C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up205c-repeated-switch-failure-localization -count=1;if($LASTEXITCODE-ne 0){throw 'UP205C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP205C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up205c-repeated-switch-failure-localization)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP205C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up205c-repeated-switch-failure-localization)|Out-String).Trim();if($P1-cne $P2){throw 'UP205C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up205c-repeated-switch-failure-localization.v1' -or [int]$R.points.Count-ne 256 -or [int]$R.summaries.Count-ne 4 -or [int]$R.action_cap-ne 8){throw 'UP205C_DESIGN'}
 if($R.rule_changed -or -not $R.counterfactual_only -or $R.live_activation){throw 'UP205C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "schedule=$($S.schedule) failures=$($S.treated_failures) failures_at_cap=$($S.failures_at_action_cap) mean_failure_actions=$($S.mean_actions_among_failures) mean_last_action_gap=$($S.mean_last_action_to_failure) min_failure=$($S.min_treated_failure_step) max_failure=$($S.max_treated_failure_step)"}
 foreach($P in $R.points){if(-not $P.survives_horizon){Write-Host "failure schedule=$($P.schedule) cohort=$($P.cohort_index) hand=$($P.initial_hand) endangered=$($P.endangered_key) baseline=$($P.baseline_loss_step) treated=$($P.treated_loss_step) onset=$($P.trigger_onset) actions=$($P.action_count) last_action=$($P.last_action_write) cap=$($P.action_cap_reached) delay=$($P.loss_delay)"}}
 Write-Host 'WINGLESS_UP276_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
