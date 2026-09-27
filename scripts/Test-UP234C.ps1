$ErrorActionPreference='Stop'
$Repo=(Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$PriorGoCache=$env:GOCACHE;$PriorGoTmp=$env:GOTMPDIR
$GoCache=Join-Path $env:TEMP ("wingless-up234c-gocache-"+$PID);$GoTmp=Join-Path $env:TEMP ("wingless-up234c-gotmp-"+$PID)
New-Item -ItemType Directory -Force -Path $GoCache|Out-Null;New-Item -ItemType Directory -Force -Path $GoTmp|Out-Null
$env:GOCACHE=$GoCache;$env:GOTMPDIR=$GoTmp
Push-Location $Repo
try{
 go env -w GOFLAGS=-buildvcs=false;if($LASTEXITCODE-ne 0){throw 'UP234C_GO_ENV_FAILED'}
 go clean -cache -testcache;if($LASTEXITCODE-ne 0){throw 'UP234C_GO_CLEAN_FAILED'}
 go test ./unitary ./cmd/unitary-up234c-distance-metric-transfer -count=1;if($LASTEXITCODE-ne 0){throw 'UP234C_FOCUSED_TEST_FAILED'}
 go test ./... -count=1;if($LASTEXITCODE-ne 0){throw 'UP234C_FULL_REGRESSION_FAILED'}
 $P1=((go run ./cmd/unitary-up234c-distance-metric-transfer)|Out-String).Trim();if($LASTEXITCODE-ne 0){throw 'UP234C_PROBE1_FAILED'}
 $P2=((go run ./cmd/unitary-up234c-distance-metric-transfer)|Out-String).Trim();if($P1-cne $P2){throw 'UP234C_NONDETERMINISTIC_OUTPUT'}
 $R=$P1|ConvertFrom-Json
 if($R.schema-cne 'wingless.up234c-distance-metric-transfer.v1' -or [int]$R.summaries.Count-ne 2){throw 'UP234C_DESIGN'}
 if($R.heldout_fitting_used -or $R.new_native_field_used -or $R.learned_scaling_used -or $R.learned_weights_used -or $R.threshold_fitting_used -or $R.intervention_changed -or $R.live_activation){throw 'UP234C_BOUNDARY_LEAK'}
 Write-Host $P1;Write-Host '';Write-Host '=== SCIENTIFIC DIAGNOSIS (does not control harness acceptance) ==='
 foreach($S in $R.summaries){Write-Host "metric=$($S.metric) failures=$($S.failures) survivors=$($S.survivors) event_fail=$($S.event_failures) event_survive=$($S.event_survivors) high=$($S.high_risk) low=$($S.low_risk) unknown=$($S.unknown) noevent=$($S.no_event) high_fail=$($S.high_risk_failures) high_survive=$($S.high_risk_survivors) low_fail=$($S.low_risk_failures) unknown_fail=$($S.unknown_failures) recall=$($S.high_risk_failure_recall) precision=$($S.high_risk_precision)"}
 Write-Host 'WINGLESS_UP334_HARNESS_PASS'
}finally{Pop-Location;if($null-eq $PriorGoCache){Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue}else{$env:GOCACHE=$PriorGoCache};if($null-eq $PriorGoTmp){Remove-Item Env:GOTMPDIR -ErrorAction SilentlyContinue}else{$env:GOTMPDIR=$PriorGoTmp};Remove-Item -Recurse -Force $GoCache,$GoTmp -ErrorAction SilentlyContinue}
