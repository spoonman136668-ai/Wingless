$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up221c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up221c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP221C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP221C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up221c-residual-dual-horizon-diagnostic -count=1;if($LASTEXITCODE-ne 0){throw 'UP221C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP221C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up221c-residual-dual-horizon-diagnostic)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP221C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up221c-residual-dual-horizon-diagnostic)|Out-String).Trim();if($P1-cne $P2){throw 'UP221C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up221c-residual-dual-horizon-diagnostic.v1' -or [int]$R.arms.Count-ne 64 -or [int]$R.summaries.Count-ne 2){throw 'UP221C_DESIGN'}
 if($R.intervention_changed -or $R.threshold_fitting_used -or $R.adaptive_signal_selection_used -or $R.live_activation){throw 'UP221C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 Write-Host "silent_failures=$($R.silent_cap8_failures) caught_noquery_critical=$($R.silent_caught_by_no_query_critical) caught_noquery_near=$($R.silent_caught_by_no_query_near)"
 foreach($S in $R.summaries){Write-Host "outcome=$($S.outcome) arms=$($S.arms) adv_critical=$($S.adversarial_critical) adv_rate=$($S.adversarial_critical_rate) nq_critical=$($S.no_query_critical) nq_critical_rate=$($S.no_query_critical_rate) nq_near=$($S.no_query_near) nq_near_rate=$($S.no_query_near_rate)"}
 Write-Host 'WINGLESS_UP310_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
