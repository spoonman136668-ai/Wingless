$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up206c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up206c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP206C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP206C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up206c-repeated-switch-warning-trace -count=1;if($LASTEXITCODE-ne 0){throw 'UP206C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP206C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up206c-repeated-switch-warning-trace)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP206C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up206c-repeated-switch-warning-trace)|Out-String).Trim();if($P1-cne $P2){throw 'UP206C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up206c-repeated-switch-warning-trace.v1' -or [int]$R.points.Count-ne 256 -or [int]$R.summaries.Count-ne 4){throw 'UP206C_DESIGN'}
 if($R.rule_changed -or -not $R.diagnostic_observation_only -or -not $R.counterfactual_only -or $R.live_activation){throw 'UP206C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "schedule=$($S.schedule) failures=$($S.treated_failures) critical_final=$($S.failures_critical_at_final_monitor) posttrigger_critical=$($S.failures_with_post_trigger_critical) mean_lead=$($S.mean_critical_lead_writes) min_lead=$($S.min_critical_lead_writes) max_lead=$($S.max_critical_lead_writes)"}
 foreach($P in $R.points){if(-not $P.survives_horizon){Write-Host "failure schedule=$($P.schedule) cohort=$($P.cohort_index) hand=$($P.initial_hand) baseline=$($P.baseline_loss_step) treated=$($P.treated_loss_step) onset=$($P.trigger_onset) last_monitor=$($P.last_monitor_before_loss) horizon_at_last=$($P.horizon_at_last_monitor) last_critical=$($P.last_critical_before_loss) lead=$($P.critical_lead_writes) postcritical=$($P.post_trigger_critical_count)"}}
 Write-Host 'WINGLESS_UP278_HARNESS_PASS'
}finally{
 Pop-Location
 if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache}
 if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp}
 Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue
}
